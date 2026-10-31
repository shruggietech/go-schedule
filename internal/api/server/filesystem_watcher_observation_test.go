package server

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

type observedWatcherScheduler struct {
	watcherSchedulerStub
	health domain.WatcherHealth
}

func (s *observedWatcherScheduler) WatcherHealth(string) domain.WatcherHealth { return s.health }

func TestWatcherPendingComparesActualDefinitionAndPreservesReadiness(t *testing.T) {
	server := newTestServer(t)
	enabled := true
	task := newTaskFor(t, server, TaskCreateRequest{Command: "echo", Schedule: "every hour", Enabled: &enabled})
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	definition := domain.FilesystemWatcher{ID: "watcher-1", Name: "Inbox", Kind: domain.WatcherDirectory, Path: t.TempDir(), Pattern: "*.csv", TargetTaskID: task.Task.ID, Enabled: true, UpdatedAt: now}
	scheduler := &observedWatcherScheduler{health: domain.WatcherHealth{State: domain.WatcherActive, ChangedAt: now}}
	server.sched = scheduler
	// The runtime has not yet observed a newly saved definition.
	assertPending := func(want *bool) map[string]json.RawMessage {
		t.Helper()
		data, err := json.Marshal(server.watcherResponse(definition))
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(data, &fields); err != nil {
			t.Fatal(err)
		}
		if want == nil {
			if _, ok := fields["reconciliation_pending"]; ok {
				t.Fatalf("unavailable runtime claimed reconciliation: %s", data)
			}
			return fields
		}
		var got *bool
		if err := json.Unmarshal(fields["reconciliation_pending"], &got); err != nil || got == nil || *got != *want {
			t.Fatalf("pending want %v: %s", *want, data)
		}
		return fields
	}
	pending, current := true, false
	fields := assertPending(&pending)
	var revision string
	if err := json.Unmarshal(fields["definition_revision"], &revision); err != nil || revision == "" {
		t.Fatalf("missing saved revision: %s", fields["definition_revision"])
	}
	setObserved := func(revision string, state domain.WatcherHealthState) {
		t.Helper()
		data, _ := json.Marshal(map[string]any{"state": state, "observed_at": now, "observed_definition_revision": revision})
		if err := json.Unmarshal(data, &scheduler.health); err != nil {
			t.Fatal(err)
		}
	}
	setObserved(revision, domain.WatcherActive)
	assertPending(&current)
	// Disable retains the previous active observation while readiness reflects the saved definition.
	definition.Enabled = false
	assertPending(&pending)
	response := server.watcherResponse(definition)
	if response.Health.State != domain.WatcherActive || response.Readiness != "disabled" {
		t.Fatalf("saved and observed state conflated: %+v", response)
	}
	fields = assertPending(&pending)
	_ = json.Unmarshal(fields["definition_revision"], &revision)
	setObserved(revision, domain.WatcherDisabled)
	assertPending(&current)
	definition.Enabled = true
	assertPending(&pending)
	if got := server.watcherResponse(definition); got.Health.State != domain.WatcherDisabled || got.Readiness != "ready" {
		t.Fatalf("enable fabricated health or lost readiness: %+v", got)
	}
	// A stale disabled observation cannot acknowledge a newer same-health pattern edit.
	fields = assertPending(&pending)
	_ = json.Unmarshal(fields["definition_revision"], &revision)
	definition.Pattern = "*.txt"
	setObserved(revision, domain.WatcherActive)
	assertPending(&pending)
	fields = assertPending(&pending)
	_ = json.Unmarshal(fields["definition_revision"], &revision)
	setObserved(revision, domain.WatcherActive)
	assertPending(&current)
	setObserved(revision, domain.WatcherDegraded)
	scheduler.health.Reason = "permission denied"
	if got := server.watcherResponse(definition); got.Readiness != "degraded" || got.Reason != "permission denied" {
		t.Fatalf("degraded readiness lost: %+v", got)
	}
	server.sched = nil
	assertPending(nil)
	if got := server.watcherResponse(definition); got.Health.State != domain.WatcherDegraded || got.Readiness != "degraded" {
		t.Fatalf("runtime unavailable context lost: %+v", got)
	}
}
