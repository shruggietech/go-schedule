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

// TaskCreateRequest is the body for POST /v1/tasks. Provide either Schedule
// (human-readable recurrence) or At (one-off instant), not both.
type TaskCreateRequest struct {
	Name          string            `json:"name"`
	GroupID       string            `json:"group_id,omitempty"`
	Command       string            `json:"command"`
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

// TaskResponse is the detail returned for a task.
type TaskResponse struct {
	Task     domain.Task     `json:"task"`
	Schedule domain.Schedule `json:"schedule"`
	NextRuns []time.Time     `json:"next_runs"`
}

func (s *Server) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req TaskCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON: "+err.Error())
		return
	}
	now := time.Now().UTC()

	// Validate basics.
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, CodeValidation, "name", "name is required")
		return
	}
	if req.Command == "" {
		writeError(w, http.StatusBadRequest, CodeValidation, "command", "command is required")
		return
	}
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

	// Build the schedule: one-off (At) or recurring (Schedule).
	var sch domain.Schedule
	switch {
	case req.At != nil:
		if !req.At.After(now) {
			writeError(w, http.StatusBadRequest, CodeValidation, "at", "one-off time is in the past")
			return
		}
		sch = schedule.NewOneOff(*req.At)
	case req.Schedule != "":
		parsed, err := schedule.Parse(req.Schedule, tz, now)
		if err != nil {
			writeError(w, http.StatusBadRequest, CodeValidation, "schedule", err.Error())
			return
		}
		sch = parsed
	default:
		writeError(w, http.StatusBadRequest, CodeValidation, "schedule", "provide either 'schedule' or 'at'")
		return
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

	if err := s.store.CreateSchedule(&sch); err != nil {
		s.internal(w, err)
		return
	}
	task := &domain.Task{
		Name: req.Name, GroupID: req.GroupID, Command: req.Command, Args: req.Args,
		WorkingDir: req.WorkingDir, Env: req.Env, RunAs: req.RunAs, Enabled: true,
		Timezone: tz, ScheduleID: sch.ID, OverlapPolicy: overlap, CatchupPolicy: catchup,
		State: domain.TaskActive,
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
	sch, err := s.store.GetSchedule(task.ScheduleID)
	if err != nil {
		s.internal(w, err)
		return
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
		s.notFoundOr(w, err)
		return
	}
	s.reload()
	s.publishTaskUpdated(id)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleRunNow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := s.store.GetTask(id); err != nil {
		s.notFoundOr(w, err)
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
	Schedule string `json:"schedule"`
	Timezone string `json:"timezone,omitempty"`
}

type PreviewResponse struct {
	RRULE        string      `json:"rrule"`
	HumanSummary string      `json:"human_summary"`
	NextRuns     []time.Time `json:"next_runs"`
}

func (s *Server) handlePreview(w http.ResponseWriter, r *http.Request) {
	var req PreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	tz := orDefault(req.Timezone, "Local")
	now := time.Now().UTC()
	sch, err := schedule.Parse(req.Schedule, tz, now)
	if err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "schedule", err.Error())
		return
	}
	runs, _ := schedule.UpcomingRuns(sch, tz, now, 5)
	writeJSON(w, http.StatusOK, PreviewResponse{RRULE: sch.RRULE, HumanSummary: sch.HumanSummary, NextRuns: runs})
}

// ---- helpers ------------------------------------------------------------

func (s *Server) taskDetail(task domain.Task, sch domain.Schedule, now time.Time) TaskResponse {
	runs, _ := schedule.UpcomingRuns(sch, task.Timezone, now, 5)
	// Schedules stored before migration v4 carry no phrase, so an editing client
	// would have nothing to put in the schedule field. Derive an equivalent one
	// from the rule for the response only — never written back, so the stored
	// row (and the user's original wording, once it exists) stays authoritative.
	if sch.Expression == "" {
		sch.Expression = schedule.Render(sch, task.Timezone)
	}
	return TaskResponse{Task: task, Schedule: sch, NextRuns: runs}
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

func validOverlap(p domain.OverlapPolicy) bool {
	switch p {
	case domain.OverlapQueueOne, domain.OverlapSkip, domain.OverlapAllowConcurrent:
		return true
	}
	return false
}
