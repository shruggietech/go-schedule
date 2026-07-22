package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/executor"
	"github.com/shruggietech/go-schedule/internal/schedule"
	"github.com/shruggietech/go-schedule/internal/store"
	"github.com/shruggietech/go-schedule/internal/timezone"
)

// TaskUpdateRequest carries optional task fields. Empty/nil fields are left
// unchanged. Providing Schedule or At replaces the task's schedule.
//
// GroupID is a pointer because group membership needs three distinct intents,
// and an empty string cannot carry two of them: nil leaves membership
// unchanged, a pointer to "" removes the task from all groups, and a pointer to
// an id assigns it. (Same convention as GroupUpdateRequest.Parent.)
type TaskUpdateRequest struct {
	Name          string            `json:"name,omitempty"`
	GroupID       *string           `json:"group_id,omitempty"`
	Command       string            `json:"command,omitempty"`
	Args          []string          `json:"args,omitempty"`
	WorkingDir    string            `json:"working_dir,omitempty"`
	Env           map[string]string `json:"env,omitempty"`
	RunAs         string            `json:"run_as,omitempty"`
	Timezone      string            `json:"timezone,omitempty"`
	Schedule      string            `json:"schedule,omitempty"`
	At            *time.Time        `json:"at,omitempty"`
	OverlapPolicy string            `json:"overlap_policy,omitempty"`
	CatchupPolicy string            `json:"catchup_policy,omitempty"`
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
		parsed, err := schedule.Parse(req.Schedule, task.Timezone, now)
		if err != nil {
			writeError(w, http.StatusBadRequest, CodeValidation, "schedule", err.Error())
			return
		}
		sch = parsed
	}
	if sch.Kind != "" {
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
	}

	if err := s.store.UpdateTask(&task); err != nil {
		s.internal(w, err)
		return
	}
	s.reload()
	s.publishTaskUpdated(task.ID)
	writeJSON(w, http.StatusOK, s.taskDetail(task, sch, now))
}
