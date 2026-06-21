package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/logbus"
	"github.com/shruggietech/go-schedule/internal/store"
)

func newLogServer(t *testing.T) (*Server, *logbus.Ring) {
	t.Helper()
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	ring := logbus.NewRing(100)
	s := New(st, nil, nil, ring, config.NewLogger(config.Default(), discard{}))
	return s, ring
}

func getLogs(t *testing.T, s *Server, query string) LogsResponse {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/logs"+query, nil)
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var out LogsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestHandleListLogs_SeverityFilter(t *testing.T) {
	s, ring := newLogServer(t)
	ring.Add(domain.LogRecord{ID: "1", Severity: domain.SeverityInfo, Message: "i"})
	ring.Add(domain.LogRecord{ID: "2", Severity: domain.SeverityError, Message: "e"})
	ring.Add(domain.LogRecord{ID: "3", Severity: domain.SeverityWarning, Message: "w"})

	all := getLogs(t, s, "")
	if len(all.Logs) != 3 {
		t.Fatalf("all = %d, want 3", len(all.Logs))
	}

	errs := getLogs(t, s, "?severity=error")
	if len(errs.Logs) != 1 || errs.Logs[0].Severity != domain.SeverityError {
		t.Fatalf("error filter = %+v", errs.Logs)
	}
}

func TestHandleListLogs_BadSeverity(t *testing.T) {
	s, _ := newLogServer(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/logs?severity=bogus", nil)
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestHandleListLogs_NilRingEmpty(t *testing.T) {
	st, _ := store.Open(":memory:")
	t.Cleanup(func() { _ = st.Close() })
	s := New(st, nil, nil, nil, config.NewLogger(config.Default(), discard{}))
	resp := getLogs(t, s, "")
	if len(resp.Logs) != 0 {
		t.Fatalf("nil ring should yield empty logs, got %d", len(resp.Logs))
	}
}
