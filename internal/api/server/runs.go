package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/runcontrol"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

func (s *Server) handleListRuns(w http.ResponseWriter, r *http.Request) {
	limit, ok := nonNegativeQuery(w, r, "limit", 100)
	if !ok {
		return
	}
	offset, ok := nonNegativeQuery(w, r, "offset", 0)
	if !ok {
		return
	}
	outputLimit, ok := nonNegativeQuery(w, r, "output_limit", 0)
	if !ok {
		return
	}
	runs, err := s.store.ListRunsPage(r.URL.Query().Get("task"), offset, limit, outputLimit)
	if err != nil {
		s.internal(w, err)
		return
	}
	if runs == nil {
		runs = []domain.Run{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"runs": runs})
}

func (s *Server) handleGetRun(w http.ResponseWriter, r *http.Request) {
	run, err := s.store.GetRun(r.PathValue("id"))
	if err != nil {
		s.notFoundOr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, run)
}

type activeRunSource interface{ ActiveRuns() []domain.Run }

func (s *Server) handleListActiveRuns(w http.ResponseWriter, _ *http.Request) {
	runs := []domain.Run{}
	if source, ok := s.sched.(activeRunSource); ok {
		runs = source.ActiveRuns()
	}
	writeJSON(w, http.StatusOK, map[string]any{"runs": runs})
}

func (s *Server) handleListAlerts(w http.ResponseWriter, r *http.Request) {
	unacked := r.URL.Query().Get("unacked") == "true" || r.URL.Query().Get("unacked") == "1"
	limit, ok := nonNegativeQuery(w, r, "limit", 0)
	if !ok {
		return
	}
	offset, ok := nonNegativeQuery(w, r, "offset", 0)
	if !ok {
		return
	}
	messageLimit, ok := nonNegativeQuery(w, r, "message_limit", 0)
	if !ok {
		return
	}
	alerts, err := s.store.ListAlertsPage(unacked, offset, limit, messageLimit)
	if err != nil {
		s.internal(w, err)
		return
	}
	if alerts == nil {
		alerts = []domain.Alert{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"alerts": alerts})
}

func nonNegativeQuery(w http.ResponseWriter, r *http.Request, name string, fallback int) (int, bool) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return fallback, true
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		writeError(w, http.StatusBadRequest, CodeValidation, name, name+" must be a non-negative integer")
		return 0, false
	}
	return parsed, true
}

func (s *Server) handleAckAlert(w http.ResponseWriter, r *http.Request) {
	if err := s.store.AckAlert(r.PathValue("id")); err != nil {
		s.notFoundOr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// runCanceler is optional for persistence-only API consumers.
type runCanceler interface {
	CancelRun(context.Context, string) (domain.RunCancellation, error)
}

func (s *Server) handleCancelRun(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength != 0 {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "run cancellation takes only the Run ID; send an empty request body")
		return
	}
	source, ok := s.sched.(runCanceler)
	if !ok {
		writeError(w, http.StatusServiceUnavailable, CodeInternal, "run", "this daemon does not support active-run cancellation")
		return
	}
	if err := r.Context().Err(); err != nil {
		writeError(w, http.StatusRequestTimeout, CodeConflict, "run", "cancellation request ended before admission")
		return
	}
	result, err := source.CancelRun(r.Context(), r.PathValue("id"))
	switch {
	case err == nil:
		status := http.StatusAccepted
		if result.Status == "already_finished" {
			status = http.StatusOK
		}
		writeJSON(w, status, result)
	case errors.Is(err, runcontrol.ErrNotActive):
		writeError(w, http.StatusConflict, CodeConflict, "run", "this history record is not an active execution; select an active Run ID")
	case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
		writeError(w, http.StatusRequestTimeout, CodeConflict, "run", "cancellation request ended before admission; refresh the run before retrying")
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, CodeNotFound, "run", "Run ID was not found; refresh active executions or history")
	default:
		s.notFoundOr(w, err)
	}
}
