package gui

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestScheduleColumnsNameEveryField(t *testing.T) {
	want := []string{"When", "Task", "Event", "Outcome"}
	got := make([]string, len(scheduleColumns))
	for i, column := range scheduleColumns {
		got[i] = column.Header
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Schedule headers=%v, want %v", got, want)
	}
}

func TestScheduleRowModelsNormalizeEventsAndOutcomes(t *testing.T) {
	base := time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name             string
		occurrence       server.Occurrence
		wantEvent        string
		wantOutcome      string
		wantEventStyle   widget.Importance
		wantOutcomeStyle widget.Importance
	}{
		{name: "future", occurrence: server.Occurrence{TaskID: "a", TaskName: "Backup", Time: base, Kind: "scheduled"}, wantEvent: "▷ SCHEDULED", wantOutcome: "— NOT AVAILABLE", wantEventStyle: widget.HighImportance, wantOutcomeStyle: widget.LowImportance},
		{name: "success", occurrence: server.Occurrence{TaskID: "a", TaskName: "Backup", Time: base, Kind: "past", Outcome: domain.OutcomeSuccess}, wantEvent: "COMPLETED", wantOutcome: "✓ SUCCESS", wantOutcomeStyle: widget.SuccessImportance},
		{name: "failure", occurrence: server.Occurrence{TaskID: "a", TaskName: "Backup", Time: base, Kind: "past", Outcome: domain.OutcomeFailure}, wantEvent: "COMPLETED", wantOutcome: "✗ FAILURE", wantOutcomeStyle: widget.DangerImportance},
		{name: "skipped", occurrence: server.Occurrence{TaskID: "a", TaskName: "Backup", Time: base, Kind: "past", Outcome: domain.OutcomeSkipped}, wantEvent: "COMPLETED", wantOutcome: "↷ SKIPPED", wantOutcomeStyle: widget.LowImportance},
		{name: "caught up", occurrence: server.Occurrence{TaskID: "a", TaskName: "Backup", Time: base, Kind: "past", Outcome: domain.OutcomeCaughtUp}, wantEvent: "COMPLETED", wantOutcome: "↻ CAUGHT UP", wantOutcomeStyle: widget.HighImportance},
		{name: "queued", occurrence: server.Occurrence{TaskID: "a", TaskName: "Backup", Time: base, Kind: "past", Outcome: domain.OutcomeQueued}, wantEvent: "COMPLETED", wantOutcome: "⋯ QUEUED", wantOutcomeStyle: widget.WarningImportance},
		{name: "missing", occurrence: server.Occurrence{TaskID: "a", TaskName: "Backup", Time: base, Kind: "past"}, wantEvent: "COMPLETED", wantOutcome: "• NOT AVAILABLE", wantOutcomeStyle: widget.MediumImportance},
		{name: "unknown", occurrence: server.Occurrence{TaskID: "a", TaskName: "Backup", Time: base, Kind: "past", Outcome: domain.RunOutcome("awaiting_review")}, wantEvent: "COMPLETED", wantOutcome: "? AWAITING REVIEW", wantOutcomeStyle: widget.MediumImportance},
		{name: "unknown event", occurrence: server.Occurrence{TaskID: "a", TaskName: "Backup", Time: base, Kind: "deferred"}, wantEvent: "? DEFERRED", wantOutcome: "— NOT AVAILABLE", wantEventStyle: widget.MediumImportance, wantOutcomeStyle: widget.LowImportance},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := scheduleRowModels([]server.Occurrence{tt.occurrence})[0]
			if row.Cells[2].Text != tt.wantEvent || row.Cells[3].Text != tt.wantOutcome {
				t.Fatalf("event/outcome=%q/%q, want %q/%q", row.Cells[2].Text, row.Cells[3].Text, tt.wantEvent, tt.wantOutcome)
			}
			if row.Cells[2].Importance != tt.wantEventStyle || row.Cells[3].Importance != tt.wantOutcomeStyle {
				t.Fatalf("importance=%v/%v, want %v/%v", row.Cells[2].Importance, row.Cells[3].Importance, tt.wantEventStyle, tt.wantOutcomeStyle)
			}
			if !strings.Contains(row.Summary, "Task: Backup") || row.Identity == "" {
				t.Fatalf("row summary/identity=%q/%q", row.Summary, row.Identity)
			}
		})
	}
}

func TestScheduleRowModelsPreserveUnicodeFallbacksAndDuplicateIdentity(t *testing.T) {
	occurrence := server.Occurrence{Time: time.Date(2026, 9, 3, 12, 0, 0, 123, time.UTC), Kind: "scheduled", TaskName: "備份 ✅"}
	rows := scheduleRowModels([]server.Occurrence{occurrence, occurrence})
	if rows[0].Identity == rows[1].Identity {
		t.Fatalf("duplicate occurrence identities collide: %q", rows[0].Identity)
	}
	if rows[0].Cells[1].Text != "備份 ✅" {
		t.Fatalf("Unicode task=%q", rows[0].Cells[1].Text)
	}
	empty := scheduleRowModels([]server.Occurrence{{Time: occurrence.Time, Kind: "scheduled"}})[0]
	if empty.Cells[1].Text != "Unnamed task" {
		t.Fatalf("missing task fallback=%q", empty.Cells[1].Text)
	}
	updated := occurrence
	updated.Kind = "past"
	updated.Outcome = domain.OutcomeSuccess
	if got := scheduleRowModels([]server.Occurrence{updated})[0].Identity; got != rows[0].Identity {
		t.Fatalf("same occurrence changed identity after outcome update: %q != %q", got, rows[0].Identity)
	}
}

func TestScheduleRowModelsUseRunIDForEqualTimeRecords(t *testing.T) {
	at := time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC)
	queued := server.Occurrence{TaskID: "task", RunID: "run-queued", Time: at, Kind: "past", Outcome: domain.OutcomeQueued}
	completed := server.Occurrence{TaskID: "task", RunID: "run-completed", Time: at, Kind: "past", Outcome: domain.OutcomeSuccess}
	first := scheduleRowModels([]server.Occurrence{queued, completed})
	reordered := scheduleRowModels([]server.Occurrence{completed, queued})
	if first[0].Identity == first[1].Identity {
		t.Fatal("equal-time runs share an identity")
	}
	if first[0].Identity != reordered[1].Identity || first[1].Identity != reordered[0].Identity {
		t.Fatalf("run identities changed with equal-time order: %q/%q vs %q/%q", first[0].Identity, first[1].Identity, reordered[0].Identity, reordered[1].Identity)
	}
}

func TestScheduleTableHasFixedHeadersAndDisclosure(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	if ui.scheduleTable == nil || ui.scheduleTable.header == nil {
		t.Fatal("Schedule did not expose the shared fixed-header table")
	}
	rows := scheduleRowModels([]server.Occurrence{{
		TaskID: "a", TaskName: "Long running backup", Kind: "past",
		Time: time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC), Outcome: domain.OutcomeFailure,
	}})
	ui.scheduleTable.setRows(rows)
	ui.scheduleTable.list.Select(0)
	for _, want := range []string{"Task: Long running backup", "Event: COMPLETED", "Outcome: ✗ FAILURE"} {
		if !strings.Contains(ui.scheduleTable.disclosure.Text, want) {
			t.Errorf("Schedule disclosure %q missing %q", ui.scheduleTable.disclosure.Text, want)
		}
	}
}

func TestScheduleColumnsPersistAndResetIndependently(t *testing.T) {
	prefs := testApp.Preferences()
	oldSchedule := prefs.String(scheduleColumnLayoutPreferenceKey)
	oldActivity := prefs.String(activityColumnLayoutPreferenceKey)
	t.Cleanup(func() {
		prefs.SetString(scheduleColumnLayoutPreferenceKey, oldSchedule)
		prefs.SetString(activityColumnLayoutPreferenceKey, oldActivity)
	})
	prefs.RemoveValue(scheduleColumnLayoutPreferenceKey)
	prefs.RemoveValue(activityColumnLayoutPreferenceKey)
	ui := NewUI(testApp, &fakeBackend{})
	ui.scheduleTable.profile.setProportions([]float32{0.4, 0.25, 0.15, 0.2}, true)
	ui.activityTable.profile.setProportions([]float32{0.2, 0.2, 0.25, 0.35}, true)

	rebuilt := NewUI(testApp, &fakeBackend{})
	if got := rebuilt.scheduleTable.profile.proportions(); !closeFloat32s(got, []float32{0.4, 0.25, 0.15, 0.2}) {
		t.Fatalf("Schedule restore=%v", got)
	}
	if got := rebuilt.activityTable.profile.proportions(); !closeFloat32s(got, []float32{0.2, 0.2, 0.25, 0.35}) {
		t.Fatalf("Activity restore=%v", got)
	}
	rebuilt.scheduleTable.resetColumns()
	if got := rebuilt.scheduleTable.profile.proportions(); !closeFloat32s(got, defaultColumnProportions(scheduleColumns)) {
		t.Fatalf("Schedule reset=%v", got)
	}
	if got := rebuilt.activityTable.profile.proportions(); !closeFloat32s(got, []float32{0.2, 0.2, 0.25, 0.35}) {
		t.Fatalf("Activity changed by Schedule reset=%v", got)
	}
}

func TestScheduleExposesResetColumnsActionAndPracticalWhenDefault(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	root := ui.navigation.contentFor(navigationSchedule)
	found := false
	walkInfoObjects(root, func(object fyne.CanvasObject) {
		if button, ok := object.(*widget.Button); ok && button.Text == "Reset columns" {
			found = true
		}
	})
	if !found {
		t.Fatal("Schedule Reset columns action missing")
	}
	defaults := defaultColumnProportions(scheduleColumns)
	if defaults[0] <= 0.2 {
		t.Fatalf("Schedule When default=%v, want practical share above 20%%", defaults[0])
	}
}

func TestScheduleControlsPreserveRangeAndListCalendarRoundTrip(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	root := ui.navigation.contentFor(navigationSchedule)
	var viewSelect, rangeSelect *widget.Select
	var walk func(fyne.CanvasObject)
	walk = func(object fyne.CanvasObject) {
		switch typed := object.(type) {
		case *widget.Select:
			switch {
			case reflect.DeepEqual(typed.Options, []string{"List", "Calendar"}):
				viewSelect = typed
			case reflect.DeepEqual(typed.Options, []string{"1 day", "7 days", "30 days"}):
				rangeSelect = typed
			}
		case *fyne.Container:
			for _, child := range typed.Objects {
				walk(child)
			}
		}
	}
	walk(root)
	if viewSelect == nil || rangeSelect == nil {
		t.Fatalf("Schedule selectors missing: view=%v range=%v", viewSelect, rangeSelect)
	}
	if viewSelect.Selected != "List" || rangeSelect.Selected != "7 days" {
		t.Fatalf("initial controls=%q/%q", viewSelect.Selected, rangeSelect.Selected)
	}

	viewSelect.SetSelected("Calendar")
	if got := findLabelText(root, "Select a day"); !strings.Contains(got, "Select a day") {
		t.Fatalf("Calendar view did not render detail guidance: %q", got)
	}
	rangeSelect.SetSelected("30 days")
	if rangeSelect.Selected != "30 days" {
		t.Fatalf("range selection=%q", rangeSelect.Selected)
	}
	viewSelect.SetSelected("List")
	if !canvasTreeContains(root, ui.scheduleTable.root) {
		t.Fatal("Schedule List table was not restored after Calendar round trip")
	}
}

func canvasTreeContains(root, target fyne.CanvasObject) bool {
	if root == target {
		return true
	}
	container, ok := root.(*fyne.Container)
	if !ok {
		return false
	}
	for _, child := range container.Objects {
		if canvasTreeContains(child, target) {
			return true
		}
	}
	return false
}

func TestSortByTimeKeepsChronologicalOrder(t *testing.T) {
	base := time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC)
	rows := sortByTime([]server.Occurrence{{TaskName: "late", Time: base.Add(time.Hour)}, {TaskName: "early", Time: base}})
	if rows[0].TaskName != "early" || rows[1].TaskName != "late" {
		t.Fatalf("order=%q/%q", rows[0].TaskName, rows[1].TaskName)
	}
}

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
	// Schedule must build (List + Calendar toggle) without panic.
	if ui.navigation.contentFor(navigationSchedule) == nil {
		t.Fatal("Schedule destination missing")
	}
}
