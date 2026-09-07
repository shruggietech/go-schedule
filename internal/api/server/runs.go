package server

import (
	"net/http"
	"strconv"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func (s *Server) handleListRuns(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			limit = n
		}
	}
	runs, err := s.store.ListRuns(r.URL.Query().Get("task"), limit)
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
	limit := 0
	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			writeError(w, http.StatusBadRequest, CodeValidation, "limit", "limit must be a non-negative integer")
			return
		}
		limit = parsed
	}
	alerts, err := s.store.ListAlertsLimited(unacked, limit)
	if err != nil {
		s.internal(w, err)
		return
	}
	if alerts == nil {
		alerts = []domain.Alert{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"alerts": alerts})
}

func (s *Server) handleAckAlert(w http.ResponseWriter, r *http.Request) {
	if err := s.store.AckAlert(r.PathValue("id")); err != nil {
		s.notFoundOr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
