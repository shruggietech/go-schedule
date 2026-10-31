package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestRunSourceColumnsPreserveCompletionCorrelationAndCompatibility(t *testing.T) {
	task, run := runSourceColumns(domain.Run{Trigger: domain.TriggerCompletion, SourceTaskID: "source-task", SourceRunID: "source-run"})
	if task != "source-task" || run != "source-run" {
		t.Fatalf("completion columns=%q %q", task, run)
	}
	task, run = runSourceColumns(domain.Run{Trigger: domain.TriggerSchedule})
	if task != "-" || run != "-" {
		t.Fatalf("schedule compatibility columns=%q %q", task, run)
	}
}

func TestRunHistoryKeepsLegacyTableUnlessDetailsRequested(t *testing.T) {
	now := time.Date(2026, 10, 19, 12, 0, 0, 0, time.UTC)
	value := domain.Run{ID: "exact-id", TaskID: "task-id", ScheduledFor: now, StartedAt: &now, Output: "captured text", OutputTruncated: true, Trigger: domain.TriggerCompletion, SourceTaskID: "upstream-task", SourceRunID: "upstream-run", SourceTriggerID: "source-trigger", SourceWatcherID: "source-watcher"}
	var plain, detail bytes.Buffer
	if err := writeRunHistory(&plain, []domain.Run{value}, false); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(plain.String(), "exact-id") || strings.Contains(plain.String(), "captured text") || !strings.HasPrefix(plain.String(), "SCHEDULED") {
		t.Fatalf("legacy output changed: %s", plain.String())
	}
	if err := writeRunHistory(&detail, []domain.Run{value}, true); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"exact-id", "task-id", "running", "captured text", "Output truncated: true", "upstream-task", "upstream-run", "source-trigger", "source-watcher"} {
		if !strings.Contains(detail.String(), want) {
			t.Errorf("missing %q in %s", want, detail.String())
		}
	}
}
