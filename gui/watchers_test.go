package gui

import (
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestWatcherRowsExposeSelectionHealthAndReadiness(t *testing.T) {
	rows := watcherRows([]server.FilesystemWatcherResponse{{ID: "watcher-id", Name: "Incoming", Kind: domain.WatcherDirectory, Path: "incoming", Pattern: "*.json", TargetTaskName: "Import", Enabled: true, Health: domain.WatcherHealth{State: domain.WatcherActive}, Readiness: "ready"}})
	if len(rows) != 1 || rows[0].Identity != "watcher-id" || len(rows[0].Cells) != len(watcherColumns) {
		t.Fatalf("rows = %+v", rows)
	}
	if rows[0].Cells[0].Text != "Incoming" || rows[0].Cells[1].Text != "incoming (*.json)" || rows[0].Cells[2].Text != "Import" || rows[0].Cells[3].Text != "Yes" || rows[0].Cells[4].Text != "Active" || rows[0].Cells[5].Text != "Ready" {
		t.Fatalf("cells = %+v", rows[0].Cells)
	}
}

func TestTriggersViewContainsFilesystemWatcherTable(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	if ui.watcherTable == nil || len(ui.watcherTable.columns) != len(watcherColumns) {
		t.Fatal("filesystem watcher table is missing")
	}
}

func TestWatcherOnlyTaskIsNotLabeledManualOnly(t *testing.T) {
	task := domain.Task{ID: "task-1", Name: "Import", Command: "echo", State: domain.TaskActive, Enabled: true}
	watchers := []server.FilesystemWatcherResponse{{ID: "watcher-1", TargetTaskID: task.ID, Enabled: true}}
	row := taskRowModelWithSources(task, nil, nil, nil, watchers)
	if row.Cells[2].Text != "Runnable" {
		t.Fatalf("effective state = %q", row.Cells[2].Text)
	}
	groupModel := newGroupTreeModelWithSources(nil, []domain.Task{task}, nil, nil, watchers)
	if label := groupModel.label(taskNodeID(task.ID)); label == "" || label == taskRowMarker+"Import   Manual only" {
		t.Fatalf("group task label = %q", label)
	}
}
