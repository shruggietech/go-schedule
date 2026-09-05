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
	tasklogic "github.com/shruggietech/go-schedule/internal/task"
	"github.com/shruggietech/go-schedule/internal/timezone"
)

// TaskCreateRequest is the body for POST /v1/tasks. Provide either Schedule
// (human-readable recurrence or supported cron) or At (one-off instant), not
// both.
type TaskCreateRequest struct {
	Name           string            `json:"name"`
	GroupID        string            `json:"group_id,omitempty"`
	Command        string            `json:"command"`
	Args           []string          `json:"args,omitempty"`
	WorkingDir     string            `json:"working_dir,omitempty"`
	Env            map[string]string `json:"env,omitempty"`
	Stdin          string            `json:"stdin,omitempty"`
	RunAs          string            `json:"run_as,omitempty"`
	Timezone       string            `json:"timezone,omitempty"`
	Schedule       string            `json:"schedule,omitempty"`
	ScheduleSyntax string            `json:"schedule_syntax,omitempty"`
	At             *time.Time        `json:"at,omitempty"`
	OverlapPolicy  string            `json:"overlap_policy,omitempty"`
	CatchupPolicy  string            `json:"catchup_policy,omitempty"`
	// MissingDatePolicy defaults to skip, which is the behavior of every task
	// created before the policy existed.
	MissingDatePolicy string `json:"missing_date_policy,omitempty"`
	TimeBasis         string `json:"time_basis,omitempty"`
	DSTGapPolicy      string `json:"dst_gap_policy,omitempty"`
	DSTOverlapPolicy  string `json:"dst_overlap_policy,omitempty"`
	// Enabled is optional for wire compatibility. Omitted requests retain the
	// historical enabled default; explicit false supports atomic draft creation.
	Enabled *bool `json:"enabled,omitempty"`
}

// TaskResponse is the detail returned for a task.
type TaskResponse struct {
	Task          domain.Task         `json:"task"`
	Schedule      *domain.Schedule    `json:"schedule"`
	Readiness     tasklogic.Readiness `json:"readiness"`
	PolicySummary string              `json:"policy_summary"`
	NextRuns      []time.Time         `json:"next_runs"`
}

func (s *Server) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req TaskCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON: "+err.Error())
		return
	}
	now := time.Now().UTC()

	tz := req.Timezone
	if tz == "" {
		tz = "Local"
	}
	if _, err := timezone.Resolve(tz); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "timezone", err.Error())
		return
	}
	if err := executor.ValidateRunAs(req.RunAs); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "run_as", err.Error())
		return
	}
	if req.ScheduleSyntax != "" && req.Schedule == "" {
		writeError(w, http.StatusBadRequest, CodeValidation, "schedule_syntax", "schedule_syntax requires schedule")
		return
	}
	// Build the schedule: one-off (At) or recurring (Schedule).
	var sch *domain.Schedule
	switch {
	case req.At != nil:
		if !req.At.After(now) {
			writeError(w, http.StatusBadRequest, CodeValidation, "at", "one-off time is in the past")
			return
		}
		oneOff := schedule.NewOneOff(*req.At)
		sch = &oneOff
	case req.Schedule != "":
		input, err := scheduleinput.Parse(req.Schedule, scheduleinput.Syntax(req.ScheduleSyntax), tz, now)
		if err != nil {
			field := "schedule"
			if errors.Is(err, scheduleinput.ErrInvalidSyntax) {
				field = "schedule_syntax"
			}
			writeError(w, http.StatusBadRequest, CodeValidation, field, err.Error())
			return
		}
		sch = &input.Schedule
	}

	overlap := domain.OverlapPolicy(orDefault(req.OverlapPolicy, string(domain.OverlapQueueOne)))
	catchup := domain.CatchupPolicy(orDefault(req.CatchupPolicy, string(domain.CatchupOne)))
	if !validOverlap(overlap) {
		writeError(w, http.StatusBadRequest, CodeValidation, "overlap_policy", "must be queue_one, skip, or allow_concurrent")
		return
	}
	if catchup != domain.CatchupOne && catchup != domain.CatchupNone {
		writeError(w, http.StatusBadRequest, CodeValidation, "catchup_policy", "must be one or none")
		return
	}
	missingDate := domain.MissingDatePolicy(orDefault(req.MissingDatePolicy, string(domain.MissingDateSkip)))
	if !validMissingDate(missingDate) {
		writeError(w, http.StatusBadRequest, CodeValidation, "missing_date_policy", "must be skip, last_valid, or next_valid")
		return
	}
	timeBasis := domain.TimeBasis(orDefault(req.TimeBasis, string(domain.TimeBasisWallClock)))
	if !validTimeBasis(timeBasis) {
		writeError(w, http.StatusBadRequest, CodeValidation, "time_basis", "must be wall_clock, elapsed, or utc")
		return
	}
	dstGap := domain.DSTGapPolicy(orDefault(req.DSTGapPolicy, string(domain.DSTGapNextValid)))
	if !validDSTGap(dstGap) {
		writeError(w, http.StatusBadRequest, CodeValidation, "dst_gap_policy", "must be next_valid or skip")
		return
	}
	dstOverlap := domain.DSTOverlapPolicy(orDefault(req.DSTOverlapPolicy, string(domain.DSTOverlapFirst)))
	if !validDSTOverlap(dstOverlap) {
		writeError(w, http.StatusBadRequest, CodeValidation, "dst_overlap_policy", "must be first, both, or last")
		return
	}
	policy := domain.SchedulePolicy{MissingDate: missingDate, TimeBasis: timeBasis, DSTGap: dstGap, DSTOverlap: dstOverlap}
	if sch != nil {
		if err := schedule.ValidatePolicy(*sch, policy); err != nil {
			writeError(w, http.StatusBadRequest, CodeValidation, "time_basis", err.Error())
			return
		}
		if err := schedule.PrepareForPolicy(sch, tz, policy); err != nil {
			writeError(w, http.StatusBadRequest, CodeValidation, "time_basis", err.Error())
			return
		}
	}

	if sch != nil {
		if err := s.store.CreateSchedule(sch); err != nil {
			s.internal(w, err)
			return
		}
	}
	enabled := sch != nil && req.Command != ""
	if req.Enabled != nil {
		enabled = *req.Enabled && enabled
	}
	scheduleID := ""
	if sch != nil {
		scheduleID = sch.ID
	}
	task := &domain.Task{
		Name: req.Name, GroupID: req.GroupID, Command: req.Command, Args: req.Args,
		WorkingDir: req.WorkingDir, Env: req.Env, Stdin: req.Stdin, RunAs: req.RunAs, Enabled: enabled,
		Timezone: tz, ScheduleID: scheduleID, OverlapPolicy: overlap, CatchupPolicy: catchup,
		MissingDatePolicy: missingDate, TimeBasis: timeBasis, DSTGapPolicy: dstGap,
		DSTOverlapPolicy: dstOverlap, State: domain.TaskActive,
	}
	if err := s.store.CreateTask(task); err != nil {
		s.internal(w, err)
		return
	}
	s.reload()
	s.publishTaskCreated(*task)
	writeJSON(w, http.StatusCreated, s.taskDetail(*task, sch, now))
}

func (s *Server) handleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.store.ListTasks(r.URL.Query().Get("group"), r.URL.Query().Get("state"))
	if err != nil {
		s.internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tasks": tasks})
}

func (s *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
	task, err := s.store.GetTask(r.PathValue("id"))
	if err != nil {
		s.notFoundOr(w, err)
		return
	}
	var sch *domain.Schedule
	if task.ScheduleID != "" {
		stored, err := s.store.GetSchedule(task.ScheduleID)
		if err != nil {
			s.internal(w, err)
			return
		}
		sch = &stored
	}
	writeJSON(w, http.StatusOK, s.taskDetail(task, sch, time.Now().UTC()))
}

func (s *Server) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.store.DeleteTask(id); err != nil {
		s.notFoundOr(w, err)
		return
	}
	s.reload()
	s.publishTaskDeleted(id)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleEnableTask(w http.ResponseWriter, r *http.Request)  { s.setEnabled(w, r, true) }
func (s *Server) handleDisableTask(w http.ResponseWriter, r *http.Request) { s.setEnabled(w, r, false) }

func (s *Server) setEnabled(w http.ResponseWriter, r *http.Request, enabled bool) {
	id := r.PathValue("id")
	if err := s.store.SetTaskEnabled(id, enabled); err != nil {
		if errors.Is(err, store.ErrTaskNotRunnable) {
			writeError(w, http.StatusBadRequest, CodeValidation, "enabled", "configure a command before enabling this task")
			return
		}
		if errors.Is(err, store.ErrTaskNotActivationReady) {
			writeError(w, http.StatusBadRequest, CodeValidation, "enabled", "configure an automatic activation source before enabling this task")
			return
		}
		if errors.Is(err, store.ErrTaskTerminal) {
			writeError(w, http.StatusBadRequest, CodeValidation, "enabled", "only active tasks can be enabled")
			return
		}
		s.notFoundOr(w, err)
		return
	}
	s.reload()
	s.publishTaskUpdated(id)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleRunNow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := s.store.GetTask(id)
	if err != nil {
		s.notFoundOr(w, err)
		return
	}
	if !tasklogic.EvaluateReadiness(task, false).CommandReady {
		writeError(w, http.StatusBadRequest, CodeValidation, "command", "configure a command before running this task")
		return
	}
	if s.sched != nil {
		if err := s.sched.RunNow(id); err != nil {
			s.internal(w, err)
			return
		}
	}
	w.WriteHeader(http.StatusAccepted)
}

// PreviewRequest/Response back POST /v1/schedules/preview.
type PreviewRequest struct {
	Schedule       string `json:"schedule"`
	ScheduleSyntax string `json:"schedule_syntax,omitempty"`
	Timezone       string `json:"timezone,omitempty"`
	// MissingDatePolicy makes the preview honest for a by-date rule: the same
	// phrase produces different run times under different policies, so a preview
	// that ignored it would show times the created task would not keep.
	MissingDatePolicy string `json:"missing_date_policy,omitempty"`
	TimeBasis         string `json:"time_basis,omitempty"`
	DSTGapPolicy      string `json:"dst_gap_policy,omitempty"`
	DSTOverlapPolicy  string `json:"dst_overlap_policy,omitempty"`
}

type PreviewResponse struct {
	RRULE              string                    `json:"rrule"`
	CalendarAdjustment domain.CalendarAdjustment `json:"calendar_adjustment,omitempty"`
	HumanSummary       string                    `json:"human_summary"`
	NextRuns           []time.Time               `json:"next_runs"`
	SourceSyntax       string                    `json:"source_syntax,omitempty"`
	PolicySummary      string                    `json:"policy_summary"`
}

func (s *Server) handlePreview(w http.ResponseWriter, r *http.Request) {
	var req PreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	if req.ScheduleSyntax != "" && req.Schedule == "" {
		writeError(w, http.StatusBadRequest, CodeValidation, "schedule_syntax", "schedule_syntax requires schedule")
		return
	}
	if req.MissingDatePolicy != "" && !validMissingDate(domain.MissingDatePolicy(req.MissingDatePolicy)) {
		writeError(w, http.StatusBadRequest, CodeValidation, "missing_date_policy", "invalid policy")
		return
	}
	if req.TimeBasis != "" && !validTimeBasis(domain.TimeBasis(req.TimeBasis)) {
		writeError(w, http.StatusBadRequest, CodeValidation, "time_basis", "invalid policy")
		return
	}
	if req.DSTGapPolicy != "" && !validDSTGap(domain.DSTGapPolicy(req.DSTGapPolicy)) {
		writeError(w, http.StatusBadRequest, CodeValidation, "dst_gap_policy", "invalid policy")
		return
	}
	if req.DSTOverlapPolicy != "" && !validDSTOverlap(domain.DSTOverlapPolicy(req.DSTOverlapPolicy)) {
		writeError(w, http.StatusBadRequest, CodeValidation, "dst_overlap_policy", "invalid policy")
		return
	}
	tz := orDefault(req.Timezone, "Local")
	now := time.Now().UTC()
	input, err := scheduleinput.Parse(req.Schedule, scheduleinput.Syntax(req.ScheduleSyntax), tz, now)
	if err != nil {
		field := "schedule"
		if errors.Is(err, scheduleinput.ErrInvalidSyntax) {
			field = "schedule_syntax"
		}
		writeError(w, http.StatusBadRequest, CodeValidation, field, err.Error())
		return
	}
	sch := input.Schedule
	policy := (domain.SchedulePolicy{
		MissingDate: domain.MissingDatePolicy(req.MissingDatePolicy), TimeBasis: domain.TimeBasis(req.TimeBasis),
		DSTGap: domain.DSTGapPolicy(req.DSTGapPolicy), DSTOverlap: domain.DSTOverlapPolicy(req.DSTOverlapPolicy),
	}).Effective()
	if err := schedule.ValidatePolicy(sch, policy); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "time_basis", err.Error())
		return
	}
	if err := schedule.PrepareForPolicy(&sch, tz, policy); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "time_basis", err.Error())
		return
	}
	runs, _ := schedule.UpcomingRunsWithPolicy(sch, tz, policy, now, 5)
	policySummary := ""
	if sch.Kind == domain.ScheduleRecurring {
		policySummary = schedule.DescribePolicy(policy)
	}
	writeJSON(w, http.StatusOK, PreviewResponse{
		RRULE:              sch.RRULE,
		CalendarAdjustment: sch.CalendarAdjustment,
		HumanSummary:       schedule.Describe(sch, domain.MissingDatePolicy(req.MissingDatePolicy)),
		NextRuns:           runs,
		SourceSyntax:       string(input.Syntax),
		PolicySummary:      policySummary,
	})
}

// ---- helpers ------------------------------------------------------------

func (s *Server) taskDetail(task domain.Task, sch *domain.Schedule, now time.Time) TaskResponse {
	var runs []time.Time
	var policySummary string
	// The summary is rendered against the task's policy on the way out rather
	// than stored, because the policy can change without the phrase changing. The
	// stored HumanSummary stays the phrase-level sentence; every client reads
	// this one, so none of them can claim a rule fires in a period it skips.
	if sch != nil {
		runs, _ = schedule.UpcomingRunsWithPolicy(*sch, task.Timezone, task.SchedulePolicy(), now, 5)
		sch.HumanSummary = schedule.Describe(*sch, task.MissingDatePolicy)
		sch.SourceSyntax = string(scheduleinput.SourceSyntax(*sch))
		if !schedule.IsStartup(*sch) {
			policySummary = schedule.DescribePolicy(task.SchedulePolicy())
		}
	}
	hasCompletion, _ := s.store.TaskHasIncomingCompletion(task.ID)
	return TaskResponse{Task: task, Schedule: sch, Readiness: tasklogic.EvaluateReadiness(task, hasCompletion), PolicySummary: policySummary, NextRuns: runs}
}

func (s *Server) reload() {
	if s.sched != nil {
		s.sched.Reload()
	}
}

func (s *Server) internal(w http.ResponseWriter, err error) {
	s.log.Error("api: internal error", "err", err)
	writeError(w, http.StatusInternalServerError, CodeInternal, "", "internal error")
}

func (s *Server) notFoundOr(w http.ResponseWriter, err error) {
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, CodeNotFound, "", "not found")
		return
	}
	s.internal(w, err)
}

func orDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

// validMissingDate reports whether p is one of the three known policies. An
// empty value never reaches here, callers default it first, so an unknown
// value is genuinely the caller's mistake.
func validMissingDate(p domain.MissingDatePolicy) bool {
	switch p {
	case domain.MissingDateSkip, domain.MissingDateLastValid, domain.MissingDateNextValid:
		return true
	}
	return false
}

func validOverlap(p domain.OverlapPolicy) bool {
	switch p {
	case domain.OverlapQueueOne, domain.OverlapSkip, domain.OverlapAllowConcurrent:
		return true
	}
	return false
}

func validTimeBasis(p domain.TimeBasis) bool {
	switch p {
	case domain.TimeBasisWallClock, domain.TimeBasisElapsed, domain.TimeBasisUTC:
		return true
	}
	return false
}

func validDSTGap(p domain.DSTGapPolicy) bool {
	return p == domain.DSTGapNextValid || p == domain.DSTGapSkip
}

func validDSTOverlap(p domain.DSTOverlapPolicy) bool {
	switch p {
	case domain.DSTOverlapFirst, domain.DSTOverlapBoth, domain.DSTOverlapLast:
		return true
	}
	return false
}
