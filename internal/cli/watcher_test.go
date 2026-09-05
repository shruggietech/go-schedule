package cli

import (
	"bytes"
	"strings"
	"testing"

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
	for _, expected := range []string{"Path: /incoming", "Pattern: *.json", "Recursive: true", "Debounce: 250ms", "Stability: 500ms", "Target: Import (task-1)", "Health: degraded", "Health reason: root is missing"} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("output %q does not contain %q", output.String(), expected)
		}
	}
}
