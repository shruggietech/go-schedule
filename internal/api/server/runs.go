package server

import (
	"net/http"
	"strconv"

	"github.com/shruggietech/go-schedule/internal/domain"
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
