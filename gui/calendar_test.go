package gui

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func TestOccurrencesByDay_Buckets(t *testing.T) {
	base := time.Date(2026, 6, 10, 9, 0, 0, 0, time.Local)
	occ := []server.Occurrence{
		{Time: base, TaskName: "a"},
		{Time: base.Add(2 * time.Hour), TaskName: "b"},
		{Time: base.AddDate(0, 0, 1), TaskName: "c"},
	}
	byDay := occurrencesByDay(occ)
	if got := len(byDay["2026-06-10"]); got != 2 {
		t.Fatalf("2026-06-10 = %d, want 2", got)
	}
	if got := len(byDay["2026-06-11"]); got != 1 {
		t.Fatalf("2026-06-11 = %d, want 1", got)
	}
}

func TestBuildCalendarGrid_PlacesAndMarks(t *testing.T) {
	month := time.Date(2026, 6, 1, 0, 0, 0, 0, time.Local)
	occ := []server.Occurrence{{Time: time.Date(2026, 6, 15, 8, 0, 0, 0, time.Local), TaskName: "x"}}
	var selected time.Time
	grid := buildCalendarGrid(occ, month, func(d time.Time) { selected = d })

	// 7 weekday headers + leading blanks + 30 day cells.
	if len(grid.Objects) < 7+30 {
		t.Fatalf("grid has %d objects, want >= 37", len(grid.Objects))
	}

	// Find the marked day cell (15) and tap it.
	var tapped bool
	for _, o := range grid.Objects {
		if b, ok := o.(*cursorButton); ok && b.Text == "15 •1" {
			b.OnTapped()
			tapped = true
		}
	}
	if !tapped {
		t.Fatal("did not find marked day cell '15 •1'")
	}
	if selected.Day() != 15 {
		t.Fatalf("selected day = %d, want 15", selected.Day())
	}
}

func TestBuildCalendarGrid_EmptyMonthNoError(t *testing.T) {
	month := time.Date(2026, 2, 1, 0, 0, 0, 0, time.Local)
	grid := buildCalendarGrid(nil, month, nil)
	if grid == nil || len(grid.Objects) < 7 {
		t.Fatal("empty calendar should still render headers")
	}
}

func TestOccurrencesOnDayText(t *testing.T) {
	day := time.Date(2026, 6, 15, 0, 0, 0, 0, time.Local)
	occ := []server.Occurrence{{Time: time.Date(2026, 6, 15, 8, 0, 0, 0, time.Local), TaskName: "x"}}
	if got := occurrencesOnDayText(occ, day); got == "" {
		t.Fatal("expected non-empty text for a day with runs")
	}
	empty := occurrencesOnDayText(nil, day)
	if empty == "" {
		t.Fatal("expected a 'no runs' message for an empty day")
	}
}

func TestUI_ScheduleTabBuilds(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	// Schedule is the 2nd tab; it must build (List + Calendar toggle) without panic.
	if ui.tabs.Items[1].Text != "Schedule" || ui.tabs.Items[1].Content == nil {
		t.Fatalf("schedule tab missing: %+v", ui.tabs.Items[1])
	}
}
