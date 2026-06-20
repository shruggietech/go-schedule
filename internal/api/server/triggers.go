package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/shruggietech/go-scheduler/internal/domain"
)

// TriggerCreateRequest is the body for POST /v1/triggers.
type TriggerCreateRequest struct {
	SourceTaskID string `json:"source_task_id"`
	TargetTaskID string `json:"target_task_id"`
	OnOutcome    string `json:"on_outcome,omitempty"`   // success|failure|any (default success)
	DedupKey     string `json:"dedup_key,omitempty"`    // empty => per source-run
	DedupWindow  string `json:"dedup_window,omitempty"` // Go duration, e.g. "5m"
}

func (s *Server) handleCreateTrigger(w http.ResponseWriter, r *http.Request) {
	var req TriggerCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	if req.SourceTaskID == "" || req.TargetTaskID == "" {
		writeError(w, http.StatusBadRequest, CodeValidation, "source_task_id", "source_task_id and target_task_id are required")
		return
	}
	// Both tasks must exist.
	for field, id := range map[string]string{"source_task_id": req.SourceTaskID, "target_task_id": req.TargetTaskID} {
		if _, err := s.store.GetTask(id); err != nil {
			writeError(w, http.StatusBadRequest, CodeValidation, field, "task does not exist")
			return
		}
	}
	outcome := domain.TriggerOutcome(orDefault(req.OnOutcome, string(domain.OnSuccess)))
	switch outcome {
	case domain.OnSuccess, domain.OnFailure, domain.OnAny:
	default:
		writeError(w, http.StatusBadRequest, CodeValidation, "on_outcome", "must be success, failure, or any")
		return
	}
	var window time.Duration
	if req.DedupWindow != "" {
		d, err := time.ParseDuration(req.DedupWindow)
		if err != nil || d < 0 {
			writeError(w, http.StatusBadRequest, CodeValidation, "dedup_window", "must be a non-negative Go duration (e.g. 5m)")
			return
		}
		window = d
	}
	tr := &domain.Trigger{
		SourceTaskID: req.SourceTaskID, TargetTaskID: req.TargetTaskID,
		OnOutcome: outcome, DedupKey: req.DedupKey, DedupWindow: window,
	}
	if err := s.store.CreateTrigger(tr); err != nil {
		s.internal(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, tr)
}

func (s *Server) handleListTriggers(w http.ResponseWriter, _ *http.Request) {
	triggers, err := s.store.ListTriggers()
	if err != nil {
		s.internal(w, err)
		return
	}
	if triggers == nil {
		triggers = []domain.Trigger{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"triggers": triggers})
}

func (s *Server) handleDeleteTrigger(w http.ResponseWriter, r *http.Request) {
	if err := s.store.DeleteTrigger(r.PathValue("id")); err != nil {
		s.notFoundOr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
