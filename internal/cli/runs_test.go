package cli

import (
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
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
