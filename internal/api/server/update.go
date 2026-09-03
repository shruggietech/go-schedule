package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/executor"
	"github.com/shruggietech/go-schedule/internal/schedule"
	"github.com/shruggietech/go-schedule/internal/scheduleinput"
	"github.com/shruggietech/go-schedule/internal/store"
	"github.com/shruggietech/go-schedule/internal/timezone"
)

// TaskUpdateRequest carries optional task fields. Empty/nil fields are left
// unchanged, except Args where nil means unchanged and a non-nil empty slice
// explicitly clears every argument. Providing Schedule or At replaces the
// task's schedule.
//
// GroupID is a pointer because group membership needs three distinct intents,
// and an empty string cannot carry two of them: nil leaves membership
// unchanged, a pointer to "" removes the task from all groups, and a pointer to
// an id assigns it. (Same convention as GroupUpdateRequest.Parent.)
type TaskUpdateRequest struct {
	Name           string            `json:"name,omitempty"`
	GroupID        *string           `json:"group_id,omitempty"`
	Command        string            `json:"command,omitempty"`
	Args           []string          `json:"args"`
	WorkingDir     string            `json:"working_dir,omitempty"`
	Env            map[string]string `json:"env,omitempty"`
	Stdin          *string           `json:"stdin,omitempty"`
	RunAs          string            `json:"run_as,omitempty"`
	Timezone       string            `json:"timezone,omitempty"`
	Schedule       string            `json:"schedule,omitempty"`
	ScheduleSyntax string            `json:"schedule_syntax,omitempty"`
	At             *time.Time        `json:"at,omitempty"`
	OverlapPolicy  string            `json:"overlap_policy,omitempty"`
	CatchupPolicy  string            `json:"catchup_policy,omitempty"`
	// MissingDatePolicy is independent of Schedule: replacing the phrase leaves
	// the policy alone and vice versa, because the policy states the operator's
	// intent for calendar anomalies rather than anything about the phrase.
	MissingDatePolicy string `json:"missing_date_policy,omitempty"`
	TimeBasis         string `json:"time_basis,omitempty"`
	DSTGapPolicy      string `json:"dst_gap_policy,omitempty"`
	DSTOverlapPolicy  string `json:"dst_overlap_policy,omitempty"`
}

func (s *Server) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var req TaskUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	task, err := s.store.GetTask(r.PathValue("id"))
	if err != nil {
		s.notFoundOr(w, err)
		return
	}
	now := time.Now().UTC()
	if req.ScheduleSyntax != "" && req.Schedule == "" {
		writeError(w, http.StatusBadRequest, CodeValidation, "schedule_syntax", "schedule_syntax requires schedule")
		return
	}

	if req.Name != "" {
		task.Name = req.Name
	}
	if req.Command != "" {
		task.Command = req.Command
	}
	if req.Args != nil {
		task.Args = req.Args
	}
	if req.WorkingDir != "" {
		task.WorkingDir = req.WorkingDir
	}
	if req.Env != nil {
		task.Env = req.Env
	}
	if req.Stdin != nil {
		task.Stdin = *req.Stdin
	}
	if req.GroupID != nil {
		// "" clears membership; a named group must exist, or this is the
		// caller's mistake and belongs in a validation error rather than a
		// foreign-key failure.
		if *req.GroupID != "" {
			if _, err := s.store.GetGroup(*req.GroupID); err != nil {
				if errors.Is(err, store.ErrNotFound) {
					writeError(w, http.StatusBadRequest, CodeValidation, "group_id", "group not found")
					return
				}
				s.internal(w, err)
				return
			}
		}
		task.GroupID = *req.GroupID
	}
	if req.RunAs != "" {
		if err := executor.ValidateRunAs(req.RunAs); err != nil {
			writeError(w, http.StatusBadRequest, CodeValidation, "run_as", err.Error())
			return
		}
		task.RunAs = req.RunAs
	}
	if req.Timezone != "" {
		if _, err := timezone.Resolve(req.Timezone); err != nil {
			writeError(w, http.StatusBadRequest, CodeValidation, "timezone", err.Error())
			return
		}
		task.Timezone = req.Timezone
	}
	if req.OverlapPolicy != "" {
		p := domain.OverlapPolicy(req.OverlapPolicy)
		if !validOverlap(p) {
			writeError(w, http.StatusBadRequest, CodeValidation, "overlap_policy", "invalid policy")
			return
		}
		task.OverlapPolicy = p
	}
	if req.CatchupPolicy != "" {
		c := domain.CatchupPolicy(req.CatchupPolicy)
		if c != domain.CatchupOne && c != domain.CatchupNone {
			writeError(w, http.StatusBadRequest, CodeValidation, "catchup_policy", "invalid policy")
			return
		}
		task.CatchupPolicy = c
	}

	if req.MissingDatePolicy != "" {
		m := domain.MissingDatePolicy(req.MissingDatePolicy)
		if !validMissingDate(m) {
			writeError(w, http.StatusBadRequest, CodeValidation, "missing_date_policy", "invalid policy")
			return
		}
		task.MissingDatePolicy = m
	}
	if req.TimeBasis != "" {
		p := domain.TimeBasis(req.TimeBasis)
		if !validTimeBasis(p) {
			writeError(w, http.StatusBadRequest, CodeValidation, "time_basis", "invalid policy")
			return
		}
		task.TimeBasis = p
	}
	if req.DSTGapPolicy != "" {
		p := domain.DSTGapPolicy(req.DSTGapPolicy)
		if !validDSTGap(p) {
			writeError(w, http.StatusBadRequest, CodeValidation, "dst_gap_policy", "invalid policy")
			return
		}
		task.DSTGapPolicy = p
	}
	if req.DSTOverlapPolicy != "" {
		p := domain.DSTOverlapPolicy(req.DSTOverlapPolicy)
		if !validDSTOverlap(p) {
			writeError(w, http.StatusBadRequest, CodeValidation, "dst_overlap_policy", "invalid policy")
			return
		}
		task.DSTOverlapPolicy = p
	}

	// Optional schedule replacement.
	var sch domain.Schedule
	switch {
	case req.At != nil:
		if !req.At.After(now) {
			writeError(w, http.StatusBadRequest, CodeValidation, "at", "one-off time is in the past")
			return
		}
		sch = schedule.NewOneOff(*req.At)
	case req.Schedule != "":
		input, err := scheduleinput.Parse(req.Schedule, scheduleinput.Syntax(req.ScheduleSyntax), task.Timezone, now)
		if err != nil {
			field := "schedule"
			if errors.Is(err, scheduleinput.ErrInvalidSyntax) {
				field = "schedule_syntax"
			}
			writeError(w, http.StatusBadRequest, CodeValidation, field, err.Error())
			return
		}
		sch = input.Schedule
	}
	if sch.Kind != "" {
		if err := schedule.PrepareForPolicy(&sch, task.Timezone, task.SchedulePolicy()); err != nil {
			writeError(w, http.StatusBadRequest, CodeValidation, "time_basis", err.Error())
			return
		}
		if err := s.store.CreateSchedule(&sch); err != nil {
			s.internal(w, err)
			return
		}
		task.ScheduleID = sch.ID
		// A revived one-off/recurring task becomes active again.
		if task.State == domain.TaskCompleted {
			task.State = domain.TaskActive
		}
	} else {
		sch, err = s.store.GetSchedule(task.ScheduleID)
		if err != nil {
			s.internal(w, err)
			return
		}
		needsElapsedEpoch := task.SchedulePolicy().TimeBasis == domain.TimeBasisElapsed && sch.ElapsedEpoch == nil
		if err := schedule.PrepareForPolicy(&sch, task.Timezone, task.SchedulePolicy()); err != nil {
			writeError(w, http.StatusBadRequest, CodeValidation, "time_basis", err.Error())
			return
		}
		if needsElapsedEpoch {
			// Schedules are immutable once referenced. Persist the newly bound
			// absolute phase as a replacement row when an existing task first
			// switches to elapsed mode.
			sch.ID = ""
			if err := s.store.CreateSchedule(&sch); err != nil {
				s.internal(w, err)
				return
			}
			task.ScheduleID = sch.ID
		}
	}

	if err := s.store.UpdateTask(&task); err != nil {
		s.internal(w, err)
		return
	}
	s.reload()
	s.publishTaskUpdated(task.ID)
	writeJSON(w, http.StatusOK, s.taskDetail(task, sch, now))
}
