package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/events"
	"github.com/shruggietech/go-schedule/internal/store"
)

// ChainCreateRequest is the body for POST /v1/chains.
type ChainCreateRequest struct {
	SourceTaskID string `json:"source_task_id"`
	TargetTaskID string `json:"target_task_id"`
	OnOutcome    string `json:"on_outcome"`
}

// ChainUpdateRequest is the body for PATCH /v1/chains/{id}. Nil fields are
// unchanged, while provided empty values are validated and rejected.
type ChainUpdateRequest struct {
	SourceTaskID *string `json:"source_task_id,omitempty"`
	TargetTaskID *string `json:"target_task_id,omitempty"`
	OnOutcome    *string `json:"on_outcome,omitempty"`
}

func (s *Server) handleCreateChain(w http.ResponseWriter, r *http.Request) {
	var req ChainCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	chain := domain.CompletionChain{
		SourceTaskID: req.SourceTaskID,
		TargetTaskID: req.TargetTaskID,
		OnOutcome:    domain.CompletionOutcome(req.OnOutcome),
	}
	if field, message := validateChainInput(chain); field != "" {
		writeError(w, http.StatusBadRequest, CodeValidation, field, message)
		return
	}
	if err := s.store.CreateCompletionChain(&chain); err != nil {
		s.chainError(w, err)
		return
	}
	created, err := s.store.GetCompletionChain(chain.ID)
	if err != nil {
		s.internal(w, err)
		return
	}
	s.publishChain(events.VerbCreated, created)
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleListChains(w http.ResponseWriter, _ *http.Request) {
	chains, err := s.store.ListCompletionChains()
	if err != nil {
		s.internal(w, err)
		return
	}
	if chains == nil {
		chains = []domain.CompletionChain{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"chains": chains})
}

func (s *Server) handleGetChain(w http.ResponseWriter, r *http.Request) {
	chain, err := s.store.GetCompletionChain(r.PathValue("id"))
	if err != nil {
		s.notFoundOr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, chain)
}

func (s *Server) handleUpdateChain(w http.ResponseWriter, r *http.Request) {
	chain, err := s.store.GetCompletionChain(r.PathValue("id"))
	if err != nil {
		s.notFoundOr(w, err)
		return
	}
	var req ChainUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	if req.SourceTaskID == nil && req.TargetTaskID == nil && req.OnOutcome == nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "provide at least one chain field")
		return
	}
	if req.SourceTaskID != nil {
		chain.SourceTaskID = *req.SourceTaskID
	}
	if req.TargetTaskID != nil {
		chain.TargetTaskID = *req.TargetTaskID
	}
	if req.OnOutcome != nil {
		chain.OnOutcome = domain.CompletionOutcome(*req.OnOutcome)
	}
	if field, message := validateChainInput(chain); field != "" {
		writeError(w, http.StatusBadRequest, CodeValidation, field, message)
		return
	}
	if err := s.store.UpdateCompletionChain(&chain); err != nil {
		s.chainError(w, err)
		return
	}
	updated, err := s.store.GetCompletionChain(chain.ID)
	if err != nil {
		s.internal(w, err)
		return
	}
	s.publishChain(events.VerbUpdated, updated)
	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) handleDeleteChain(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.store.DeleteCompletionChain(id); err != nil {
		s.notFoundOr(w, err)
		return
	}
	s.publishChainDeleted(id)
	w.WriteHeader(http.StatusNoContent)
}

func validateChainInput(chain domain.CompletionChain) (string, string) {
	if chain.SourceTaskID == "" {
		return "source_task_id", "source task is required"
	}
	if chain.TargetTaskID == "" {
		return "target_task_id", "target task is required"
	}
	if chain.SourceTaskID == chain.TargetTaskID {
		return "target_task_id", "source and target must be different tasks"
	}
	if chain.OnOutcome != domain.CompletionOnSuccess && chain.OnOutcome != domain.CompletionOnFailure && chain.OnOutcome != domain.CompletionOnAny {
		return "on_outcome", "outcome must be success, failure, or any"
	}
	return "", ""
}

func (s *Server) chainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, CodeNotFound, "", "source task, target task, or chain was not found")
	case errors.Is(err, store.ErrChainCycle):
		writeError(w, http.StatusBadRequest, CodeValidation, "target_task_id", "chain would create a cycle")
	case errors.Is(err, store.ErrDuplicateChain):
		writeError(w, http.StatusBadRequest, CodeValidation, "on_outcome", "this source, target, and outcome chain already exists")
	case errors.Is(err, store.ErrInvalidOutcome):
		writeError(w, http.StatusBadRequest, CodeValidation, "on_outcome", "outcome must be success, failure, or any")
	default:
		s.internal(w, err)
	}
}
