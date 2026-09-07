package gui

import (
	"testing"

	"fyne.io/fyne/v2"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestSuggestionEntryInsertsOnlyOnEmptyForwardTab(t *testing.T) {
	entry := newSuggestionEntry("uname -a")
	entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyTab})
	if entry.Text != "uname -a" {
		t.Fatalf("text=%q", entry.Text)
	}
	if entry.CursorColumn != len([]rune(entry.Text)) {
		t.Fatalf("cursor=%d want=%d", entry.CursorColumn, len([]rune(entry.Text)))
	}
	if entry.AcceptsTab() {
		t.Fatal("next Tab should return to traversal")
	}
	entry.SetText("")
	if !entry.AcceptsTab() {
		t.Fatal("clearing should offer the suggestion again")
	}
}

func TestExistingBlankTaskNeverOffersSuggestion(t *testing.T) {
	detail := server.TaskResponse{Task: domain.Task{ID: "draft", Timezone: "Local"}}
	editor, _ := newTestEditorDetail(t, &detail)
	if editor.commandLine.suggestion != "" {
		t.Fatalf("edit suggestion=%q", editor.commandLine.suggestion)
	}
}
