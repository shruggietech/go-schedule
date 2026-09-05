package gui

import (
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func TestTriggerRowsExposeStructuredLifecycleFields(t *testing.T) {
	rows := triggerRows([]server.TriggerResponse{{ID: "trigger-id", Name: "Build hook", TargetTaskID: "task-id", TargetTaskName: "Build", Enabled: true, Readiness: "ready"}})
	if len(rows) != 1 || rows[0].Identity != "trigger-id" || len(rows[0].Cells) != len(triggerColumns) {
		t.Fatalf("rows = %+v", rows)
	}
	if rows[0].Cells[0].Text != "Build hook" || rows[0].Cells[1].Text != "Build" || rows[0].Cells[2].Text != "Yes" || rows[0].Cells[3].Text != "Ready" {
		t.Fatalf("cells = %+v", rows[0].Cells)
	}
}

func TestUIIncludesTriggerViewAfterChains(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	labels := ui.navigation.labels()
	if len(labels) < 4 || labels[2] != "Chains" || labels[3] != "Triggers" || ui.navigation.contentFor(navigationTriggers) == nil {
		t.Fatalf("navigation = %v", labels)
	}
}
