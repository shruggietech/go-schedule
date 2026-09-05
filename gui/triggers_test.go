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
	if rows[0].Cells[0].Text != "Build hook" || rows[0].Cells[1].Text != "Standalone" || rows[0].Cells[2].Text != "-" || rows[0].Cells[3].Text != "Build" || rows[0].Cells[4].Text != "Yes" || rows[0].Cells[5].Text != "Ready" {
		t.Fatalf("cells = %+v", rows[0].Cells)
	}
}

func TestTriggerRowsExposeSetMembershipWithoutSecrets(t *testing.T) {
	rows := triggerRows([]server.TriggerResponse{{ID: "trigger-id", Name: "Deploy set 07", SetID: "set-id", SetName: "Deploy set", SetPosition: 7, TargetTaskID: "task-id", TargetTaskName: "Deploy"}})
	if rows[0].Cells[1].Text != "Deploy set" || rows[0].Cells[2].Text != "7" {
		t.Fatalf("cells = %+v", rows[0].Cells)
	}
}

func TestTriggerSetCommandsUseOneOrderedCommandPerLine(t *testing.T) {
	result := server.TriggerSetSecretResponse{Members: []server.TriggerSetSecretMember{{Position: 1, Command: "gosched trigger fire first"}, {Position: 2, Command: "gosched trigger fire second"}}}
	if got, want := triggerSetCommands(result), "gosched trigger fire first\ngosched trigger fire second\n"; got != want {
		t.Fatalf("commands = %q, want %q", got, want)
	}
}

func TestUIIncludesTriggerViewAfterChains(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	labels := ui.navigation.labels()
	if len(labels) < 4 || labels[2] != "Chains" || labels[3] != "Triggers" || ui.navigation.contentFor(navigationTriggers) == nil {
		t.Fatalf("navigation = %v", labels)
	}
}
