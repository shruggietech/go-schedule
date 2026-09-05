package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/events"
	"github.com/shruggietech/go-schedule/internal/store"
)

type FilesystemWatcherCreateRequest struct {
	Name         string             `json:"name"`
	Kind         domain.WatcherKind `json:"kind"`
	Path         string             `json:"path"`
	Pattern      string             `json:"pattern,omitempty"`
	Recursive    bool               `json:"recursive"`
	Debounce     string             `json:"debounce,omitempty"`
	Stability    string             `json:"stability,omitempty"`
	TargetTaskID string             `json:"target_task_id"`
	Enabled      *bool              `json:"enabled,omitempty"`
}

type FilesystemWatcherUpdateRequest struct {
	Name         *string             `json:"name,omitempty"`
	Kind         *domain.WatcherKind `json:"kind,omitempty"`
	Path         *string             `json:"path,omitempty"`
	Pattern      *string             `json:"pattern,omitempty"`
	Recursive    *bool               `json:"recursive,omitempty"`
	Debounce     *string             `json:"debounce,omitempty"`
	Stability    *string             `json:"stability,omitempty"`
	TargetTaskID *string             `json:"target_task_id,omitempty"`
}

type FilesystemWatcherResponse struct {
	ID             string               `json:"id"`
	Name           string               `json:"name"`
	Kind           domain.WatcherKind   `json:"kind"`
	Path           string               `json:"path"`
	Pattern        string               `json:"pattern,omitempty"`
	Recursive      bool                 `json:"recursive"`
	Debounce       string               `json:"debounce"`
	Stability      string               `json:"stability"`
	TargetTaskID   string               `json:"target_task_id"`
	TargetTaskName string               `json:"target_task_name,omitempty"`
	Enabled        bool                 `json:"enabled"`
	Health         domain.WatcherHealth `json:"health"`
	Readiness      string               `json:"readiness"`
	Reason         string               `json:"reason,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
}

type watcherRuntime interface {
	ReloadWatchers()
	WatcherHealth(string) domain.WatcherHealth
}

func (s *Server) watcherResponse(watcher domain.FilesystemWatcher) FilesystemWatcherResponse {
	response := FilesystemWatcherResponse{ID: watcher.ID, Name: watcher.Name, Kind: watcher.Kind, Path: watcher.Path, Pattern: watcher.Pattern, Recursive: watcher.Recursive, Debounce: watcher.Debounce.String(), Stability: watcher.Stability.String(), TargetTaskID: watcher.TargetTaskID, TargetTaskName: watcher.TargetTaskName, Enabled: watcher.Enabled, CreatedAt: watcher.CreatedAt, UpdatedAt: watcher.UpdatedAt}
	if runtime, ok := s.sched.(watcherRuntime); ok {
		response.Health = runtime.WatcherHealth(watcher.ID)
	} else if watcher.Enabled {
		response.Health = domain.WatcherHealth{State: domain.WatcherDegraded, Reason: "watcher runtime is unavailable"}
	} else {
		response.Health = domain.WatcherHealth{State: domain.WatcherDisabled, Reason: "watcher is disabled"}
	}
	if !watcher.Enabled {
		response.Readiness, response.Reason = "disabled", "Watcher is disabled."
		return response
	}
	task, err := s.store.GetTask(watcher.TargetTaskID)
	if err != nil {
		response.Readiness, response.Reason = "target_missing", "Target task is missing."
		return response
	}
	switch {
	case strings.TrimSpace(task.Command) == "":
		response.Readiness, response.Reason = "command_incomplete", "Target task has no command."
	case task.State != domain.TaskActive:
		response.Readiness, response.Reason = "task_inactive", "Target task is not active."
	case !task.Enabled:
		response.Readiness, response.Reason = "task_disabled", "Target task is disabled."
	default:
		enabled, err := s.store.GroupChainEnabled(task.GroupID)
		if err != nil || !enabled {
			response.Readiness, response.Reason = "group_blocked", "An ancestor group disables execution."
		} else if response.Health.State == domain.WatcherDegraded {
			response.Readiness, response.Reason = "degraded", response.Health.Reason
		} else {
			response.Readiness = "ready"
		}
	}
	return response
}

func (s *Server) handleCreateFilesystemWatcher(w http.ResponseWriter, r *http.Request) {
	var req FilesystemWatcherCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	watcher := domain.FilesystemWatcher{Name: req.Name, Kind: req.Kind, Path: req.Path, Pattern: req.Pattern, Recursive: req.Recursive, TargetTaskID: req.TargetTaskID, Enabled: true}
	if req.Enabled != nil {
		watcher.Enabled = *req.Enabled
	}
	var err error
	if watcher.Debounce, err = parseWatcherDuration(req.Debounce); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "debounce", err.Error())
		return
	}
	if watcher.Stability, err = parseWatcherDuration(req.Stability); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "stability", err.Error())
		return
	}
	if err := s.store.CreateFilesystemWatcher(&watcher); err != nil {
		s.watcherStoreError(w, err)
		return
	}
	s.reloadWatchers()
	created, _ := s.store.GetFilesystemWatcher(watcher.ID)
	s.publishFilesystemWatcher(events.VerbCreated, created)
	writeJSON(w, http.StatusCreated, s.watcherResponse(created))
}

func (s *Server) handleListFilesystemWatchers(w http.ResponseWriter, _ *http.Request) {
	items, err := s.store.ListFilesystemWatchers()
	if err != nil {
		s.internal(w, err)
		return
	}
	responses := make([]FilesystemWatcherResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, s.watcherResponse(item))
	}
	writeJSON(w, http.StatusOK, map[string]any{"filesystem_watchers": responses})
}

func (s *Server) handleGetFilesystemWatcher(w http.ResponseWriter, r *http.Request) {
	watcher, err := s.store.GetFilesystemWatcher(r.PathValue("id"))
	if err != nil {
		s.notFoundOr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, s.watcherResponse(watcher))
}

func (s *Server) handleUpdateFilesystemWatcher(w http.ResponseWriter, r *http.Request) {
	watcher, err := s.store.GetFilesystemWatcher(r.PathValue("id"))
	if err != nil {
		s.notFoundOr(w, err)
		return
	}
	var req FilesystemWatcherUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	if req.Name != nil {
		watcher.Name = *req.Name
	}
	if req.Kind != nil {
		watcher.Kind = *req.Kind
	}
	if req.Path != nil {
		watcher.Path = *req.Path
	}
	if req.Pattern != nil {
		watcher.Pattern = *req.Pattern
	}
	if req.Recursive != nil {
		watcher.Recursive = *req.Recursive
	}
	if req.TargetTaskID != nil {
		watcher.TargetTaskID = *req.TargetTaskID
	}
	if req.Debounce != nil {
		watcher.Debounce, err = parseWatcherDuration(*req.Debounce)
		if err != nil {
			writeError(w, http.StatusBadRequest, CodeValidation, "debounce", err.Error())
			return
		}
	}
	if req.Stability != nil {
		watcher.Stability, err = parseWatcherDuration(*req.Stability)
		if err != nil {
			writeError(w, http.StatusBadRequest, CodeValidation, "stability", err.Error())
			return
		}
	}
	if err := s.store.UpdateFilesystemWatcher(&watcher); err != nil {
		s.watcherStoreError(w, err)
		return
	}
	s.reloadWatchers()
	updated, _ := s.store.GetFilesystemWatcher(watcher.ID)
	s.publishFilesystemWatcher(events.VerbUpdated, updated)
	writeJSON(w, http.StatusOK, s.watcherResponse(updated))
}

func (s *Server) handleDeleteFilesystemWatcher(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.store.DeleteFilesystemWatcher(id); err != nil {
		s.notFoundOr(w, err)
		return
	}
	s.reloadWatchers()
	s.publishFilesystemWatcherDeleted(id)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleEnableFilesystemWatcher(w http.ResponseWriter, r *http.Request) {
	s.setFilesystemWatcherEnabled(w, r.PathValue("id"), true)
}
func (s *Server) handleDisableFilesystemWatcher(w http.ResponseWriter, r *http.Request) {
	s.setFilesystemWatcherEnabled(w, r.PathValue("id"), false)
}

func (s *Server) setFilesystemWatcherEnabled(w http.ResponseWriter, id string, enabled bool) {
	if err := s.store.SetFilesystemWatcherEnabled(id, enabled); err != nil {
		s.notFoundOr(w, err)
		return
	}
	s.reloadWatchers()
	watcher, _ := s.store.GetFilesystemWatcher(id)
	s.publishFilesystemWatcher(events.VerbUpdated, watcher)
	writeJSON(w, http.StatusOK, s.watcherResponse(watcher))
}

func (s *Server) reloadWatchers() {
	if runtime, ok := s.sched.(watcherRuntime); ok {
		runtime.ReloadWatchers()
	}
}

func parseWatcherDuration(value string) (time.Duration, error) {
	if value == "" {
		return 0, nil
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, errors.New("use a duration such as 250ms or 2s")
	}
	return duration, nil
}

func (s *Server) watcherStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, CodeNotFound, "target_task_id", "watcher or target task was not found")
	case errors.Is(err, store.ErrInvalidWatcher):
		writeError(w, http.StatusBadRequest, CodeValidation, "filesystem_watcher", "name, kind, path, timing, pattern, or target task is invalid")
	default:
		s.internal(w, err)
	}
}
