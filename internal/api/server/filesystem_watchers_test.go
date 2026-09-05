package server

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/events"
)

type watcherSchedulerStub struct{ reloads int }

func (*watcherSchedulerStub) Reload()             {}
func (*watcherSchedulerStub) RunNow(string) error { return nil }
func (s *watcherSchedulerStub) ReloadWatchers()   { s.reloads++ }
func (*watcherSchedulerStub) WatcherHealth(string) domain.WatcherHealth {
	return domain.WatcherHealth{State: domain.WatcherActive}
}

func TestFilesystemWatcherLifecycleAndDurationContract(t *testing.T) {
	server := newTestServer(t)
	task := newTaskFor(t, server, TaskCreateRequest{Command: "echo"})
	request := FilesystemWatcherCreateRequest{Name: "incoming", Kind: domain.WatcherDirectory, Path: filepath.Join(t.TempDir(), "missing"), Pattern: "*.json", Recursive: true, Debounce: "300ms", Stability: "750ms", TargetTaskID: task.Task.ID}
	recorder := doJSON(t, server, http.MethodPost, "/v1/filesystem-watchers", request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	var created FilesystemWatcherResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.Debounce != "300ms" || created.Stability != "750ms" || created.Health.State != domain.WatcherDegraded {
		t.Fatalf("created = %+v", created)
	}
	recorder = doJSON(t, server, http.MethodGet, "/v1/filesystem-watchers", nil)
	if recorder.Code != http.StatusOK {
		t.Fatalf("list status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	name := "renamed"
	recorder = doJSON(t, server, http.MethodPatch, "/v1/filesystem-watchers/"+created.ID, FilesystemWatcherUpdateRequest{Name: &name})
	if recorder.Code != http.StatusOK {
		t.Fatalf("update status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	recorder = doJSON(t, server, http.MethodPost, "/v1/filesystem-watchers/"+created.ID+"/disable", nil)
	if recorder.Code != http.StatusOK {
		t.Fatalf("disable status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	recorder = doJSON(t, server, http.MethodDelete, "/v1/filesystem-watchers/"+created.ID, nil)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("delete status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestFilesystemWatcherRejectsInvalidDuration(t *testing.T) {
	server := newTestServer(t)
	task := newTaskFor(t, server, TaskCreateRequest{Command: "echo"})
	recorder := doJSON(t, server, http.MethodPost, "/v1/filesystem-watchers", FilesystemWatcherCreateRequest{Name: "bad", Kind: domain.WatcherFile, Path: "file", Debounce: "soon", TargetTaskID: task.Task.ID})
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestFilesystemWatcherMutationReloadsRuntimeAndPublishesOneEvent(t *testing.T) {
	server := newTestServer(t)
	scheduler := &watcherSchedulerStub{}
	broker := events.NewBroker()
	server.sched = scheduler
	server.broker = broker
	task := newTaskFor(t, server, TaskCreateRequest{Command: "echo"})
	stream, cancel := broker.Subscribe()
	defer cancel()
	recorder := doJSON(t, server, http.MethodPost, "/v1/filesystem-watchers", FilesystemWatcherCreateRequest{Name: "incoming", Kind: domain.WatcherDirectory, Path: t.TempDir(), TargetTaskID: task.Task.ID})
	if recorder.Code != http.StatusCreated || scheduler.reloads != 1 {
		t.Fatalf("status=%d reloads=%d body=%s", recorder.Code, scheduler.reloads, recorder.Body.String())
	}
	select {
	case event := <-stream:
		if event.Kind != events.KindWatcher || event.Watcher == nil || event.Watcher.Verb != events.VerbCreated {
			t.Fatalf("event = %+v", event)
		}
	default:
		t.Fatal("watcher lifecycle event was not published")
	}
}
