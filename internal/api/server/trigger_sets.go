package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/events"
	"github.com/shruggietech/go-schedule/internal/store"
)

type TriggerSetCreateRequest struct {
	Name         string `json:"name"`
	TargetTaskID string `json:"target_task_id"`
	Count        int    `json:"count"`
	Enabled      *bool  `json:"enabled,omitempty"`
}

type TriggerSetRetargetRequest struct {
	TargetTaskID string `json:"target_task_id"`
}

type TriggerSetResponse struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	TargetTaskID   string            `json:"target_task_id"`
	TargetTaskName string            `json:"target_task_name,omitempty"`
	MemberCount    int               `json:"member_count"`
	EnabledCount   int               `json:"enabled_count"`
	Members        []TriggerResponse `json:"members"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type TriggerSetSecretMember struct {
	Position  int    `json:"position"`
	TriggerID string `json:"trigger_id"`
	Key       string `json:"key"`
	Command   string `json:"command"`
}

type TriggerSetSecretResponse struct {
	TriggerSet TriggerSetResponse       `json:"trigger_set"`
	Members    []TriggerSetSecretMember `json:"members"`
}

func (s *Server) triggerSetResponse(set domain.TriggerSet) TriggerSetResponse {
	response := TriggerSetResponse{ID: set.ID, Name: set.Name, TargetTaskID: set.TargetTaskID, TargetTaskName: set.TargetTaskName, MemberCount: len(set.Members), Members: make([]TriggerResponse, 0, len(set.Members)), CreatedAt: set.CreatedAt, UpdatedAt: set.UpdatedAt}
	for _, member := range set.Members {
		response.Members = append(response.Members, s.triggerResponse(member))
		if member.Enabled {
			response.EnabledCount++
		}
	}
	return response
}

func (s *Server) triggerSetSecretResponse(set domain.TriggerSet) TriggerSetSecretResponse {
	response := TriggerSetSecretResponse{TriggerSet: s.triggerSetResponse(set), Members: make([]TriggerSetSecretMember, 0, len(set.Members))}
	for _, member := range set.Members {
		response.Members = append(response.Members, TriggerSetSecretMember{Position: member.SetPosition, TriggerID: member.ID, Key: member.Key, Command: "gosched trigger fire " + member.Key})
	}
	return response
}

func (s *Server) handleCreateTriggerSet(w http.ResponseWriter, r *http.Request) {
	var req TriggerSetCreateRequest
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
	if req.Count < 1 || req.Count > 99 {
		writeError(w, http.StatusBadRequest, CodeValidation, "count", "count must be between 1 and 99")
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	set := domain.TriggerSet{Name: req.Name, TargetTaskID: req.TargetTaskID}
	if err := s.store.CreateTriggerSet(&set, req.Count, enabled); err != nil {
		s.triggerSetStoreError(w, err)
		return
	}
	loaded, err := s.store.GetTriggerSet(set.ID)
	if err != nil {
		s.internal(w, err)
		return
	}
	s.publishTriggerSet(events.VerbCreated, loaded)
	writeJSON(w, http.StatusCreated, s.triggerSetSecretResponse(loaded))
}

func (s *Server) handleListTriggerSets(w http.ResponseWriter, _ *http.Request) {
	sets, err := s.store.ListTriggerSets()
	if err != nil {
		s.internal(w, err)
		return
	}
	responses := make([]TriggerSetResponse, 0, len(sets))
	for _, set := range sets {
		responses = append(responses, s.triggerSetResponse(set))
	}
	writeJSON(w, http.StatusOK, map[string]any{"trigger_sets": responses})
}

func (s *Server) handleGetTriggerSet(w http.ResponseWriter, r *http.Request) {
	set, err := s.store.GetTriggerSet(r.PathValue("id"))
	if err != nil {
		s.notFoundOr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, s.triggerSetResponse(set))
}

func (s *Server) handleRevealTriggerSet(w http.ResponseWriter, r *http.Request) {
	set, err := s.store.GetTriggerSet(r.PathValue("id"))
	if err != nil {
		s.notFoundOr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, s.triggerSetSecretResponse(set))
}

func (s *Server) handleRetargetTriggerSet(w http.ResponseWriter, r *http.Request) {
	var req TriggerSetRetargetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TargetTaskID == "" {
		writeError(w, http.StatusBadRequest, CodeValidation, "target_task_id", "target task is required")
		return
	}
	if err := s.store.RetargetTriggerSet(r.PathValue("id"), req.TargetTaskID); err != nil {
		s.triggerSetStoreError(w, err)
		return
	}
	s.writeUpdatedTriggerSet(w, r.PathValue("id"))
}

func (s *Server) handleEnableTriggerSet(w http.ResponseWriter, r *http.Request) {
	s.setTriggerSetEnabled(w, r.PathValue("id"), true)
}

func (s *Server) handleDisableTriggerSet(w http.ResponseWriter, r *http.Request) {
	s.setTriggerSetEnabled(w, r.PathValue("id"), false)
}

func (s *Server) setTriggerSetEnabled(w http.ResponseWriter, id string, enabled bool) {
	if err := s.store.SetTriggerSetEnabled(id, enabled); err != nil {
		s.triggerSetStoreError(w, err)
		return
	}
	s.writeUpdatedTriggerSet(w, id)
}

func (s *Server) writeUpdatedTriggerSet(w http.ResponseWriter, id string) {
	set, err := s.store.GetTriggerSet(id)
	if err != nil {
		s.internal(w, err)
		return
	}
	s.publishTriggerSet(events.VerbUpdated, set)
	writeJSON(w, http.StatusOK, s.triggerSetResponse(set))
}

func (s *Server) handleRotateTriggerSet(w http.ResponseWriter, r *http.Request) {
	set, err := s.store.RotateTriggerSet(r.PathValue("id"))
	if err != nil {
		s.triggerSetStoreError(w, err)
		return
	}
	s.publishTriggerSet(events.VerbUpdated, set)
	writeJSON(w, http.StatusOK, s.triggerSetSecretResponse(set))
}

func (s *Server) handleDeleteTriggerSet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	set, err := s.store.DeleteTriggerSet(id)
	if err != nil {
		s.triggerSetStoreError(w, err)
		return
	}
	s.publishTriggerSet(events.VerbDeleted, set)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) triggerSetStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, CodeNotFound, "trigger_set", "Trigger Set or target task was not found")
	case errors.Is(err, store.ErrInvalidTriggerSet):
		writeError(w, http.StatusBadRequest, CodeValidation, "trigger_set", "name, target task, and a count from 1 through 99 are required")
	default:
		s.internal(w, fmt.Errorf("trigger set: %w", err))
	}
}
