package gui

import (
	"reflect"
	"testing"

	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestTaskColumnsNameEveryField(t *testing.T) {
	want := []string{"Task", "Enabled", "Lifecycle", "Time zone", "Group"}
	got := make([]string, len(taskColumns))
	for i, column := range taskColumns {
		got[i] = column.Header
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("task headers=%v, want %v", got, want)
	}
}

func TestTaskRowModelSeparatesEnabledFromLifecycle(t *testing.T) {
	groups := []domain.Group{{ID: "nightly", Name: "Nightly"}}
	tests := []struct {
		name       string
		task       domain.Task
		wantCells  []string
		importance []widget.Importance
	}{
		{
			name:       "enabled active",
			task:       domain.Task{ID: "a", Name: "Backup", Enabled: true, State: domain.TaskActive, Timezone: "UTC", GroupID: "nightly"},
			wantCells:  []string{"Backup", "Enabled", "Active", "UTC", "Nightly"},
			importance: []widget.Importance{widget.MediumImportance, widget.SuccessImportance, widget.HighImportance, widget.MediumImportance, widget.MediumImportance},
		},
		{
			name:       "disabled completed",
			task:       domain.Task{ID: "b", Name: "One shot", Enabled: false, State: domain.TaskCompleted, Timezone: "America/New_York"},
			wantCells:  []string{"One shot", "Disabled", "Completed", "America/New_York", "Not assigned"},
			importance: []widget.Importance{widget.MediumImportance, widget.LowImportance, widget.MediumImportance, widget.MediumImportance, widget.LowImportance},
		},
		{
			name:       "enabled lifecycle disabled remains explicit",
			task:       domain.Task{ID: "c", Name: "Legacy", Enabled: true, State: domain.TaskDisabled, Timezone: "UTC"},
			wantCells:  []string{"Legacy", "Enabled", "Disabled", "UTC", "Not assigned"},
			importance: []widget.Importance{widget.MediumImportance, widget.SuccessImportance, widget.LowImportance, widget.MediumImportance, widget.LowImportance},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := taskRowModel(tt.task, groups)
			if row.Identity != tt.task.ID || len(row.Cells) != len(tt.wantCells) {
				t.Fatalf("row=%+v", row)
			}
			for i, cell := range row.Cells {
				if cell.Text != tt.wantCells[i] || cell.Importance != tt.importance[i] {
					t.Errorf("cell %d=(%q,%v), want (%q,%v)", i, cell.Text, cell.Importance, tt.wantCells[i], tt.importance[i])
				}
			}
			if row.Summary == "" {
				t.Fatal("task row omitted full-value summary")
			}
		})
	}
}

func TestTaskRowModelNormalizesMissingUnicodeAndUnknownValues(t *testing.T) {
	row := taskRowModel(domain.Task{
		ID: "unicode", Name: "備份 ✅", Enabled: false,
		State: domain.TaskState("awaiting_review"), GroupID: "missing",
	}, nil)
	want := []string{"備份 ✅", "Disabled", "Awaiting review", "Unknown", "Not assigned"}
	for i, cell := range row.Cells {
		if cell.Text != want[i] {
			t.Errorf("cell %d=%q, want %q", i, cell.Text, want[i])
		}
	}
	if row.Cells[2].Importance != widget.MediumImportance {
		t.Fatalf("unknown lifecycle importance=%v, want neutral", row.Cells[2].Importance)
	}

	empty := taskRowModel(domain.Task{ID: "empty"}, nil)
	if empty.Cells[0].Text != "Unnamed task" || empty.Cells[2].Text != "Unknown" {
		t.Fatalf("empty fallbacks=%q/%q", empty.Cells[0].Text, empty.Cells[2].Text)
	}
}

func TestTaskRowTextNoLongerProducesConcatenatedStatus(t *testing.T) {
	// The old formatter was the source of unlabeled "[active] disabled" rows.
	// Its absence is part of the contract: task data now enters the structured
	// model exclusively.
	if got := taskRowModel(domain.Task{ID: "x", Name: "X", Enabled: false, State: domain.TaskActive}, nil).Summary; got == "X   [active]   disabled" {
		t.Fatal("task row fell back to ambiguous concatenated status")
	}
}
