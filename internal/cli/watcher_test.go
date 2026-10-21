package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestWatcherCommandRegistersCompleteLifecycle(t *testing.T) {
	command := newWatcherCmd()
	want := map[string]bool{"create": true, "list": true, "show": true, "update": true, "enable": true, "disable": true, "rm": true}
	for _, child := range command.Commands() {
		delete(want, child.Name())
	}
	if len(want) != 0 {
		t.Fatalf("missing commands: %v", want)
	}
}

func TestWatcherHumanDetailSeparatesSavedAndObservedState(t *testing.T) {
	for _, tc := range []struct {
		pending *bool
		want    string
	}{
		{nil, "Reconciliation: unknown"},
		{boolPointer(true), "Reconciliation: pending"},
		{boolPointer(false), "Reconciliation: current"},
	} {
		var output bytes.Buffer
		item := server.FilesystemWatcherResponse{ID: "owned", Enabled: true, Health: domain.WatcherHealth{State: domain.WatcherDisabled, Reason: "watcher is disabled", ObservedAt: time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)}, Readiness: "ready", ReconciliationPending: tc.pending}
		if err := printWatcherTo(&output, item); err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{"Saved enabled: true", "Last observed health: disabled", "Last observed at: 2026-10-21T12:00:00Z", "Target readiness: ready", tc.want} {
			if !strings.Contains(output.String(), want) {
				t.Errorf("missing %q: %s", want, output.String())
			}
		}
		if tc.pending != nil && *tc.pending && !strings.Contains(output.String(), "gosched watcher show owned") {
			t.Error("pending observation needs useful follow-up guidance")
		}
	}
}

func boolPointer(value bool) *bool { return &value }

func TestWatcherKindRejectsUnknownValue(t *testing.T) {
	var value watcherKindValue
	if err := value.Set("socket"); err == nil {
		t.Fatal("expected invalid kind error")
	}
}

func TestWatcherHumanDetailIncludesConfigurationAndHealthReason(t *testing.T) {
	var output bytes.Buffer
	item := server.FilesystemWatcherResponse{ID: "watcher-1", Name: "incoming", Kind: domain.WatcherDirectory, Path: "/incoming", Pattern: "*.json", Recursive: true, Debounce: "250ms", Stability: "500ms", TargetTaskID: "task-1", TargetTaskName: "Import", Enabled: true, Health: domain.WatcherHealth{State: domain.WatcherDegraded, Reason: "root is missing"}, Readiness: "degraded", Reason: "root is missing"}
	if err := printWatcherTo(&output, item); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"Path: /incoming", "Pattern: *.json", "Recursive: true", "Debounce: 250ms", "Stability: 500ms", "Target: Import (task-1)", "Last observed health: degraded", "Health reason: root is missing"} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("output %q does not contain %q", output.String(), expected)
		}
	}
}
