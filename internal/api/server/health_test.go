package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/store"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	log := config.NewLogger(config.Default(), discard{})
	return New(st, nil, nil, nil, "", log)
}

func TestHealth_OK(t *testing.T) {
	s := newTestServer(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	s.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("content-type = %q, want application/json", ct)
	}
	var h HealthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &h); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if h.Status != "ok" || h.Version == "" {
		t.Fatalf("unexpected health: %+v", h)
	}
}

func TestRuntimeInfoReturnsEffectiveDaemonPaths(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	want := RuntimeInfoResponse{
		DataDir:      `/srv/custom schedule`,
		DatabasePath: `/srv/custom schedule/goschedule.db`,
		ConfigPath:   `/etc/go-schedule/custom.json`,
		LogPath:      `/var/log/go-schedule/events.log`,
		LockPath:     `/srv/custom schedule/goschedd.lock`,
	}
	s := NewWithRuntimeInfo(st, nil, nil, nil, want.LogPath, want, config.NewLogger(config.Default(), discard{}))
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/runtime-info", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var got RuntimeInfoResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("runtime info = %+v, want %+v", got, want)
	}
}

func TestErrorEnvelope_OnUnknownRoute(t *testing.T) {
	s := newTestServer(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/nope", nil)
	s.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	var e APIError
	if err := json.Unmarshal(rec.Body.Bytes(), &e); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	if e.Error.Code != CodeNotFound || e.Error.Message == "" {
		t.Fatalf("unexpected error envelope: %+v", e)
	}
}

// discard is an io.Writer that drops log output during tests.
type discard struct{}

func (discard) Write(p []byte) (int, error) { return len(p), nil }
