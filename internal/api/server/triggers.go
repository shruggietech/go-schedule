package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/engine"
	"github.com/shruggietech/go-schedule/internal/events"
	"github.com/shruggietech/go-schedule/internal/store"
)

type TriggerCreateRequest struct {
	Name         string `json:"name"`
	TargetTaskID string `json:"target_task_id"`
	Enabled      *bool  `json:"enabled,omitempty"`
}

type TriggerUpdateRequest struct {
	Name         *string `json:"name,omitempty"`
	TargetTaskID *string `json:"target_task_id,omitempty"`
}

type TriggerResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	TargetTaskID   string    `json:"target_task_id"`
	TargetTaskName string    `json:"target_task_name,omitempty"`
	Enabled        bool      `json:"enabled"`
	Readiness      string    `json:"readiness"`
	Reason         string    `json:"reason,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type TriggerSecretResponse struct {
	Trigger TriggerResponse `json:"trigger"`
	Key     string          `json:"key"`
	Command string          `json:"command"`
}

func (s *Server) triggerResponse(t domain.ExternalTrigger) TriggerResponse {
	response := TriggerResponse{ID: t.ID, Name: t.Name, TargetTaskID: t.TargetTaskID, TargetTaskName: t.TargetTaskName, Enabled: t.Enabled, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt}
	if !t.Enabled {
		response.Readiness, response.Reason = "disabled", "Trigger is disabled."
		return response
	}
	task, err := s.store.GetTask(t.TargetTaskID)
	if err != nil {
		response.Readiness, response.Reason = "target_missing", "Target task is missing."
		return response
	}
	switch {
	case strings.TrimSpace(task.Command) == "":
		response.Readiness, response.Reason = "command_incomplete", "Target task has no command."
	case task.State != domain.TaskActive:
		response.Readiness, response.Reason = "task_inactive", "Target task is not active."
	case !task.Enabled:
		response.Readiness, response.Reason = "task_disabled", "Target task is disabled."
	default:
		enabled, err := s.store.GroupChainEnabled(task.GroupID)
		if err != nil || !enabled {
			response.Readiness, response.Reason = "group_blocked", "An ancestor group disables execution."
		} else {
			response.Readiness = "ready"
		}
	}
	return response
}

func (s *Server) handleCreateTrigger(w http.ResponseWriter, r *http.Request) {
	var req TriggerCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, CodeValidation, "name", "name is required")
		return
	}
	if req.TargetTaskID == "" {
		writeError(w, http.StatusBadRequest, CodeValidation, "target_task_id", "target task is required")
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	t := domain.ExternalTrigger{Name: req.Name, TargetTaskID: req.TargetTaskID, Enabled: enabled}
	if err := s.store.CreateExternalTrigger(&t); err != nil {
		s.triggerStoreError(w, err)
		return
	}
	created, err := s.store.GetExternalTrigger(t.ID)
	if err != nil {
		s.internal(w, err)
		return
	}
	s.publishTrigger(events.VerbCreated, created)
	writeJSON(w, http.StatusCreated, TriggerSecretResponse{Trigger: s.triggerResponse(created), Key: created.Key, Command: "gosched trigger fire " + created.Key})
}

func (s *Server) handleListTriggers(w http.ResponseWriter, _ *http.Request) {
	items, err := s.store.ListExternalTriggers()
	if err != nil {
		s.internal(w, err)
		return
	}
	responses := make([]TriggerResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, s.triggerResponse(item))
	}
	writeJSON(w, http.StatusOK, map[string]any{"triggers": responses})
}

func (s *Server) handleGetTrigger(w http.ResponseWriter, r *http.Request) {
	t, err := s.store.GetExternalTrigger(r.PathValue("id"))
	if err != nil {
		s.notFoundOr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, s.triggerResponse(t))
}

func (s *Server) handleUpdateTrigger(w http.ResponseWriter, r *http.Request) {
	t, err := s.store.GetExternalTrigger(r.PathValue("id"))
	if err != nil {
		s.notFoundOr(w, err)
		return
	}
	var req TriggerUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	if req.Name == nil && req.TargetTaskID == nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "provide at least one trigger field")
		return
	}
	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.TargetTaskID != nil {
		t.TargetTaskID = *req.TargetTaskID
	}
	if err := s.store.UpdateExternalTrigger(&t); err != nil {
		s.triggerStoreError(w, err)
		return
	}
	updated, _ := s.store.GetExternalTrigger(t.ID)
	s.publishTrigger(events.VerbUpdated, updated)
	writeJSON(w, http.StatusOK, s.triggerResponse(updated))
}

func (s *Server) handleDeleteTrigger(w http.ResponseWriter, r *http.Request) {
	if err := s.store.DeleteExternalTrigger(r.PathValue("id")); err != nil {
		s.notFoundOr(w, err)
		return
	}
	s.publishTriggerDeleted(r.PathValue("id"))
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleEnableTrigger(w http.ResponseWriter, r *http.Request) {
	s.setTriggerEnabled(w, r.PathValue("id"), true)
}

func (s *Server) handleDisableTrigger(w http.ResponseWriter, r *http.Request) {
	s.setTriggerEnabled(w, r.PathValue("id"), false)
}

func (s *Server) setTriggerEnabled(w http.ResponseWriter, id string, enabled bool) {
	if err := s.store.SetExternalTriggerEnabled(id, enabled); err != nil {
		s.notFoundOr(w, err)
		return
	}
	t, _ := s.store.GetExternalTrigger(id)
	s.publishTrigger(events.VerbUpdated, t)
	writeJSON(w, http.StatusOK, s.triggerResponse(t))
}

func (s *Server) handleRotateTrigger(w http.ResponseWriter, r *http.Request) {
	key, err := s.store.RotateExternalTrigger(r.PathValue("id"))
	if err != nil {
		s.notFoundOr(w, err)
		return
	}
	t, _ := s.store.GetExternalTrigger(r.PathValue("id"))
	s.publishTrigger(events.VerbUpdated, t)
	writeJSON(w, http.StatusOK, TriggerSecretResponse{Trigger: s.triggerResponse(t), Key: key, Command: "gosched trigger fire " + key})
}

func (s *Server) handleRevealTrigger(w http.ResponseWriter, r *http.Request) {
	t, err := s.store.GetExternalTrigger(r.PathValue("id"))
	if err != nil {
		s.notFoundOr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, TriggerSecretResponse{Trigger: s.triggerResponse(t), Key: t.Key, Command: "gosched trigger fire " + t.Key})
}

func (s *Server) handleFireTrigger(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key string `json:"key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Key == "" {
		writeError(w, http.StatusBadRequest, CodeValidation, "key", "key is required")
		return
	}
	firer, ok := s.sched.(interface{ FireExternalTrigger(string) (string, error) })
	if !ok || firer == nil {
		writeError(w, http.StatusServiceUnavailable, "trigger_dispatch_unavailable", "", "scheduler is unavailable")
		return
	}
	id, err := firer.FireExternalTrigger(req.Key)
	if err != nil {
		s.triggerFireError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"trigger_id": id, "status": "accepted"})
}

func (s *Server) triggerStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, CodeNotFound, "target_task_id", "trigger or target task was not found")
	case errors.Is(err, store.ErrInvalidTrigger):
		writeError(w, http.StatusBadRequest, CodeValidation, "trigger", "name and target task are required")
	default:
		s.internal(w, err)
	}
}

func (s *Server) triggerFireError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, engine.ErrTriggerUnknown):
		writeError(w, http.StatusNotFound, "trigger_unknown", "", "trigger key is unknown")
	case errors.Is(err, engine.ErrTriggerDisabled):
		writeError(w, http.StatusConflict, "trigger_disabled", "", "trigger is disabled")
	case errors.Is(err, engine.ErrTriggerTargetMissing):
		writeError(w, http.StatusConflict, "trigger_target_missing", "", "target task is missing")
	case errors.Is(err, engine.ErrTriggerCommandIncomplete):
		writeError(w, http.StatusConflict, "trigger_command_incomplete", "", "target task has no command")
	case errors.Is(err, engine.ErrTriggerTaskInactive):
		writeError(w, http.StatusConflict, "trigger_task_inactive", "", "target task is not active")
	case errors.Is(err, engine.ErrTriggerTaskDisabled):
		writeError(w, http.StatusConflict, "trigger_task_disabled", "", "target task is disabled")
	case errors.Is(err, engine.ErrTriggerGroupBlocked):
		writeError(w, http.StatusConflict, "trigger_group_blocked", "", "an ancestor group disables execution")
	default:
		writeError(w, http.StatusServiceUnavailable, "trigger_dispatch_unavailable", "", "scheduler cannot accept the run")
	}
}
