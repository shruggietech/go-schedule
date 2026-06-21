package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// LogsResponse is returned by GET /v1/logs.
type LogsResponse struct {
	Logs []domain.LogRecord `json:"logs"`
}

// handleListLogs returns recent log records from the in-memory ring, newest
// first. Optional query params: severity (info|warning|error), limit, since
// (RFC3339).
func (s *Server) handleListLogs(w http.ResponseWriter, r *http.Request) {
	if s.logs == nil {
		writeJSON(w, http.StatusOK, LogsResponse{Logs: []domain.LogRecord{}})
		return
	}

	q := r.URL.Query()
	severity := domain.AlertSeverity(q.Get("severity"))
	switch severity {
	case "", domain.SeverityInfo, domain.SeverityWarning, domain.SeverityError:
		// ok
	default:
		writeError(w, http.StatusBadRequest, CodeValidation, "severity",
			"severity must be one of: info, warning, error")
		return
	}

	limit := 0
	if v := q.Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			writeError(w, http.StatusBadRequest, CodeValidation, "limit", "limit must be a non-negative integer")
			return
		}
		limit = n
	}

	recs := s.logs.Snapshot(severity, limit)

	if v := q.Get("since"); v != "" {
		since, err := time.Parse(time.RFC3339, v)
		if err != nil {
			writeError(w, http.StatusBadRequest, CodeValidation, "since", "since must be an RFC3339 timestamp")
			return
		}
		filtered := recs[:0]
		for _, rec := range recs {
			if !rec.Time.Before(since) {
				filtered = append(filtered, rec)
			}
		}
		recs = filtered
	}

	if recs == nil {
		recs = []domain.LogRecord{}
	}
	writeJSON(w, http.StatusOK, LogsResponse{Logs: recs})
}
