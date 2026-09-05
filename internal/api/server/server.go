// Package server implements the daemon's local HTTP/JSON API. It is served over
// the IPC transport (Unix socket / named pipe), not a network port. All
// responses are JSON; errors use a consistent envelope.
package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/shruggietech/go-schedule/internal/buildinfo"
	"github.com/shruggietech/go-schedule/internal/events"
	"github.com/shruggietech/go-schedule/internal/logbus"
	"github.com/shruggietech/go-schedule/internal/store"
)

// Scheduler is the subset of the engine the API needs. It may be nil in tests
// that exercise only persistence-backed endpoints.
type Scheduler interface {
	Reload()
	RunNow(taskID string) error
}

// Server holds dependencies and the route mux.
type Server struct {
	store   *store.Store
	sched   Scheduler
	broker  *events.Broker
	logs    *logbus.Ring
	logPath string
	runtime RuntimeInfoResponse
	log     *slog.Logger
	mux     *http.ServeMux
}

// RuntimeInfoResponse identifies the daemon's effective local storage paths.
// These values come from the configuration used by the running daemon, not
// from client-side defaults.
type RuntimeInfoResponse struct {
	DataDir      string `json:"data_dir"`
	DatabasePath string `json:"database_path"`
	ConfigPath   string `json:"config_path,omitempty"`
	LogPath      string `json:"log_path"`
	LockPath     string `json:"lock_path"`
}

// New constructs a Server and registers routes. sched, broker, and logs may be
// nil (e.g. in tests that exercise only persistence-backed endpoints).
func New(st *store.Store, sched Scheduler, broker *events.Broker, logs *logbus.Ring, logPath string, log *slog.Logger) *Server {
	return NewWithRuntimeInfo(st, sched, broker, logs, logPath, RuntimeInfoResponse{}, log)
}

// NewWithRuntimeInfo constructs a Server with authoritative daemon storage
// metadata for GET /v1/runtime-info.
func NewWithRuntimeInfo(st *store.Store, sched Scheduler, broker *events.Broker, logs *logbus.Ring, logPath string, runtime RuntimeInfoResponse, log *slog.Logger) *Server {
	s := &Server{store: st, sched: sched, broker: broker, logs: logs, logPath: logPath, runtime: runtime, log: log, mux: http.NewServeMux()}
	s.routes()
	return s
}

// Handler returns the HTTP handler for the API.
func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) routes() {
	s.mux.HandleFunc("GET /v1/health", s.handleHealth)
	s.mux.HandleFunc("GET /v1/runtime-info", s.handleRuntimeInfo)

	s.mux.HandleFunc("GET /v1/tasks", s.handleListTasks)
	s.mux.HandleFunc("POST /v1/tasks", s.handleCreateTask)
	s.mux.HandleFunc("GET /v1/tasks/{id}", s.handleGetTask)
	s.mux.HandleFunc("PATCH /v1/tasks/{id}", s.handleUpdateTask)
	s.mux.HandleFunc("DELETE /v1/tasks/{id}", s.handleDeleteTask)
	s.mux.HandleFunc("POST /v1/tasks/{id}/enable", s.handleEnableTask)
	s.mux.HandleFunc("POST /v1/tasks/{id}/disable", s.handleDisableTask)
	s.mux.HandleFunc("POST /v1/tasks/{id}/run-now", s.handleRunNow)

	s.mux.HandleFunc("GET /v1/chains", s.handleListChains)
	s.mux.HandleFunc("POST /v1/chains", s.handleCreateChain)
	s.mux.HandleFunc("GET /v1/chains/{id}", s.handleGetChain)
	s.mux.HandleFunc("PATCH /v1/chains/{id}", s.handleUpdateChain)
	s.mux.HandleFunc("DELETE /v1/chains/{id}", s.handleDeleteChain)

	s.mux.HandleFunc("GET /v1/groups", s.handleListGroups)
	s.mux.HandleFunc("POST /v1/groups", s.handleCreateGroup)
	s.mux.HandleFunc("GET /v1/groups/{id}", s.handleGetGroup)
	s.mux.HandleFunc("PATCH /v1/groups/{id}", s.handleUpdateGroup)
	s.mux.HandleFunc("DELETE /v1/groups/{id}", s.handleDeleteGroup)
	s.mux.HandleFunc("POST /v1/groups/{id}/enable", s.handleEnableGroup)
	s.mux.HandleFunc("POST /v1/groups/{id}/disable", s.handleDisableGroup)

	s.mux.HandleFunc("POST /v1/schedules/preview", s.handlePreview)

	s.mux.HandleFunc("GET /v1/runs", s.handleListRuns)
	s.mux.HandleFunc("GET /v1/runs/{id}", s.handleGetRun)
	s.mux.HandleFunc("GET /v1/alerts", s.handleListAlerts)
	s.mux.HandleFunc("POST /v1/alerts/{id}/ack", s.handleAckAlert)
	s.mux.HandleFunc("GET /v1/logs", s.handleListLogs)

	s.mux.HandleFunc("GET /v1/calendar", s.handleCalendar)
	s.mux.HandleFunc("GET /v1/events", s.handleEvents)

	// Fallback: unmatched routes return the consistent error envelope.
	s.mux.HandleFunc("/", s.handleNotFound)
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotFound, CodeNotFound, "", "no such endpoint: "+r.Method+" "+r.URL.Path)
}

// HealthResponse is returned by GET /v1/health.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok", Version: buildinfo.Version})
}

func (s *Server) handleRuntimeInfo(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.runtime)
}

// ---- response helpers ---------------------------------------------------

// APIError is the consistent error envelope: {"error": {...}}.
type APIError struct {
	Error ErrorBody `json:"error"`
}

// ErrorBody describes a single error.
type ErrorBody struct {
	Code    string `json:"code"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// Error codes used across the API.
const (
	CodeValidation = "validation_failed"
	CodeNotFound   = "not_found"
	CodeConflict   = "conflict"
	CodeInternal   = "internal"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError emits the error envelope with the given HTTP status.
func writeError(w http.ResponseWriter, status int, code, field, msg string) {
	writeJSON(w, status, APIError{Error: ErrorBody{Code: code, Field: field, Message: msg}})
}
