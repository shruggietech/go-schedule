package gui

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestActivityColumnsNameEveryField(t *testing.T) {
	want := []string{"When", "Severity", "Source", "Summary"}
	got := make([]string, len(activityColumns))
	for i, column := range activityColumns {
		got[i] = column.Header
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Activity headers=%v, want %v", got, want)
	}
}

func TestActivityRowModelsNormalizeSeverityAndFallbacks(t *testing.T) {
	base := time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name       string
		entry      logEntry
		wantLevel  string
		wantSource string
		wantText   string
		importance widget.Importance
	}{
		{name: "info", entry: logEntry{id: "i", time: base, severity: domain.SeverityInfo, source: "daemon", message: "ready"}, wantLevel: "• INFO", wantSource: "daemon", wantText: "ready", importance: widget.HighImportance},
		{name: "empty is info", entry: logEntry{id: "e", time: base, source: "", message: ""}, wantLevel: "• INFO", wantSource: "daemon", wantText: "No message", importance: widget.HighImportance},
		{name: "warning", entry: logEntry{id: "w", time: base, severity: domain.SeverityWarning, source: "scheduler", message: "late"}, wantLevel: "⚠ WARNING", wantSource: "scheduler", wantText: "late", importance: widget.WarningImportance},
		{name: "error", entry: logEntry{id: "x", time: base, severity: domain.SeverityError, source: "alert: run_failed", message: "boom"}, wantLevel: "✗ ERROR", wantSource: "alert: run_failed", wantText: "boom", importance: widget.DangerImportance},
		{name: "unknown", entry: logEntry{id: "u", time: base, severity: domain.AlertSeverity("notice_level"), source: "agent", message: "hello"}, wantLevel: "? NOTICE LEVEL", wantSource: "agent", wantText: "hello", importance: widget.MediumImportance},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := activityRowModels([]logEntry{tt.entry})[0]
			if row.Cells[1].Text != tt.wantLevel || row.Cells[1].Importance != tt.importance || row.Cells[2].Text != tt.wantSource || row.Cells[3].Text != tt.wantText {
				t.Fatalf("row cells=%+v", row.Cells)
			}
			if row.Identity == "" || !strings.Contains(row.Summary, "Severity: "+tt.wantLevel) {
				t.Fatalf("identity/summary=%q/%q", row.Identity, row.Summary)
			}
		})
	}
}

func TestActivityRowModelsKeepLogAlertAndDuplicateIdentityDistinct(t *testing.T) {
	base := time.Date(2026, 9, 3, 12, 0, 0, 123, time.UTC)
	entries := []logEntry{
		{id: "same", time: base, severity: domain.SeverityInfo, message: "備份 ✅"},
		{id: "same", time: base, severity: domain.SeverityInfo, message: "備份 ✅", isAlert: true},
		{id: "same", time: base, severity: domain.SeverityInfo, message: "備份 ✅"},
	}
	rows := activityRowModels(entries)
	seen := make(map[string]bool)
	for _, row := range rows {
		if seen[row.Identity] {
			t.Fatalf("duplicate identity %q", row.Identity)
		}
		seen[row.Identity] = true
		if row.Cells[3].Text != "備份 ✅" {
			t.Fatalf("Unicode message=%q", row.Cells[3].Text)
		}
	}
	for index, row := range rows {
		entry, ok := activityEntryForIdentity(entries, rows, row.Identity)
		if !ok || entry.isAlert != entries[index].isAlert {
			t.Fatalf("identity %d resolved %+v, %v", index, entry, ok)
		}
	}
	if _, ok := activityEntryForIdentity(entries, rows, "removed"); ok {
		t.Fatal("stale Activity identity resolved")
	}
	updated := entries[0]
	updated.message = "updated presentation"
	updated.severity = domain.SeverityWarning
	if got := activityRowModels([]logEntry{updated})[0].Identity; got != rows[0].Identity {
		t.Fatalf("same backend event changed identity after presentation update: %q != %q", got, rows[0].Identity)
	}
}

func TestActivityTableHasFixedHeadersAndSemanticCells(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	if ui.activityTable == nil || ui.activityTable.header == nil {
		t.Fatal("Activity did not expose the shared fixed-header table")
	}
	rows := activityRowModels([]logEntry{{
		id: "error", time: time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC),
		severity: domain.SeverityError, source: "alert: run_failed", message: "task run failed",
	}})
	ui.activityTable.setRows(rows)
	row := ui.activityTable.list.CreateItem().(*structuredRow)
	ui.activityTable.list.UpdateItem(0, row)
	if row.labels[1].Text != "✗ ERROR" || row.labels[1].Importance != widget.DangerImportance {
		t.Fatalf("severity cell=%q/%v", row.labels[1].Text, row.labels[1].Importance)
	}
}

func TestMergeLogEntries_MergesSortsAndFilters(t *testing.T) {
	t0 := time.Unix(100, 0)
	logs := []domain.LogRecord{
		{ID: "l1", Time: t0.Add(1 * time.Second), Severity: domain.SeverityInfo, Message: "info"},
		{ID: "l2", Time: t0.Add(3 * time.Second), Severity: domain.SeverityError, Message: "err"},
	}
	alerts := []domain.Alert{
		{ID: "a1", CreatedAt: t0.Add(2 * time.Second), Severity: domain.SeverityWarning, Kind: domain.AlertMissedRun, Message: "missed"},
	}

	// No filter: all three, newest first.
	all := mergeLogEntries(logs, alerts, "", time.Time{})
	if len(all) != 3 {
		t.Fatalf("merged len = %d, want 3", len(all))
	}
	if all[0].message != "err" || all[2].message != "info" {
		t.Fatalf("not newest-first: %v", []string{all[0].message, all[1].message, all[2].message})
	}
	if all[0].id != "l2" || all[1].id != "a1" || all[2].id != "l1" {
		t.Fatalf("source identities not preserved: %q/%q/%q", all[0].id, all[1].id, all[2].id)
	}

	// Error filter: only the error log.
	errs := mergeLogEntries(logs, alerts, domain.SeverityError, time.Time{})
	if len(errs) != 1 || errs[0].message != "err" {
		t.Fatalf("error filter = %+v", errs)
	}

	// The alert entry carries its ID for acknowledgement.
	warns := mergeLogEntries(logs, alerts, domain.SeverityWarning, time.Time{})
	if len(warns) != 1 || !warns[0].isAlert || warns[0].alertID != "a1" {
		t.Fatalf("alert entry = %+v", warns)
	}
}

func TestMergeLogEntries_DismissCutoff(t *testing.T) {
	t0 := time.Unix(100, 0)
	logs := []domain.LogRecord{
		{ID: "old", Time: t0, Severity: domain.SeverityInfo, Message: "old"},
		{ID: "new", Time: t0.Add(10 * time.Second), Severity: domain.SeverityInfo, Message: "new"},
	}
	// Dismiss everything up to t0+5s: only "new" survives.
	got := mergeLogEntries(logs, nil, "", t0.Add(5*time.Second))
	if len(got) != 1 || got[0].message != "new" {
		t.Fatalf("dismiss cutoff = %+v", got)
	}
}

func TestUI_ActivityTabBuilds(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{
		logs:   []domain.LogRecord{{ID: "l1", Severity: domain.SeverityError, Message: "boom"}},
		alerts: []domain.Alert{{ID: "a1", Severity: domain.SeverityWarning, Kind: domain.AlertRunFailed, Message: "warn"}},
	})
	if ui.navigation.contentFor(navigationActivity) == nil || ui.navigation.label(navigationActivity) != "Activity" {
		t.Fatal("Activity destination missing or mislabeled")
	}
}

func TestActivityDiagnosticsText(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "Windows path", path: `C:\Schedule Data\日志\goschedule.log`, want: `Full daemon log: C:\Schedule Data\日志\goschedule.log`},
		{name: "custom path", path: `/var/lib/go schedule/activity.jsonl`, want: `Full daemon log: /var/lib/go schedule/activity.jsonl`},
		{name: "unavailable", want: "Full daemon log: unavailable until daemon responds."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := activityDiagnosticsText(tt.path)
			for _, phrase := range []string{"limited set", "recent daemon log records", "older daemon records", tt.want} {
				if !strings.Contains(got, phrase) {
					t.Errorf("diagnostics %q does not contain %q", got, phrase)
				}
			}
		})
	}
}

func TestUI_ActivityDiagnosticsRefreshesExactPath(t *testing.T) {
	want := `C:\Schedule Data\日志\goschedule.log`
	ui := NewUI(testApp, &fakeBackend{logPath: want})
	ui.model.OnChange = nil
	if err := ui.model.Refresh(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, refresh := range ui.refreshers {
		refresh()
	}
	if got := findLabelText(ui.navigation.contentFor(navigationActivity), "Full daemon log:"); !strings.Contains(got, want) {
		t.Fatalf("Activity diagnostics = %q, want exact path %q", got, want)
	}
}

func TestUI_ActivityDiagnosticsInitiallyUnavailable(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	if got := findLabelText(ui.navigation.contentFor(navigationActivity), "Full daemon log:"); !strings.Contains(got, "unavailable until daemon responds") {
		t.Fatalf("initial Activity diagnostics = %q", got)
	}
}

func findLabelText(root fyne.CanvasObject, prefix string) string {
	var found string
	var walk func(fyne.CanvasObject)
	walk = func(o fyne.CanvasObject) {
		switch w := o.(type) {
		case *widget.Label:
			if strings.Contains(w.Text, prefix) {
				found = w.Text
			}
		case *fyne.Container:
			for _, child := range w.Objects {
				walk(child)
			}
		}
	}
	walk(root)
	return found
}

func TestUI_ActivityClearControlExplainsNonDestructiveBehavior(t *testing.T) {
	backend := &fakeBackend{alerts: []domain.Alert{{
		ID: "visible-alert", CreatedAt: time.Now(), Severity: domain.SeverityWarning,
		Kind: domain.AlertRunFailed, Message: "warn",
	}}}
	ui := NewUI(testApp, backend)
	ui.model.OnChange = nil
	if err := ui.model.Refresh(context.Background()); err != nil {
		t.Fatalf("refresh Activity model: %v", err)
	}
	for _, refresh := range ui.refreshers {
		refresh()
	}
	var clearButton *cursorButton
	var helpText string

	var walk func(fyne.CanvasObject)
	walk = func(o fyne.CanvasObject) {
		switch w := o.(type) {
		case *cursorButton:
			if w.Text == "Clear View" {
				clearButton = w
			}
		case *widget.Label:
			if strings.Contains(w.Text, "Records are not deleted") {
				helpText = w.Text
			}
		case *fyne.Container:
			for _, child := range w.Objects {
				walk(child)
			}
		}
	}
	walk(ui.navigation.contentFor(navigationActivity))

	if clearButton == nil {
		t.Fatal("Activity view has no Clear View control")
	}
	if clearButton.Icon == nil || clearButton.Icon.Name() != theme.ContentClearIcon().Name() {
		t.Fatalf("Clear View icon = %v, want %q", clearButton.Icon, theme.ContentClearIcon().Name())
	}
	for _, phrase := range []string{"Hides current activity", "acknowledges visible alerts", "Records are not deleted"} {
		if !strings.Contains(helpText, phrase) {
			t.Errorf("Activity help %q does not contain %q", helpText, phrase)
		}
	}

	clearButton.OnTapped()
	waitFor(t, func() bool {
		acked := backend.acknowledgedAlerts()
		return len(acked) == 1 && acked[0] == "visible-alert"
	})
}

func TestActivityDetailRendersCompletionCorrelation(t *testing.T) {
	detail := attrsDetail(domain.LogRecord{
		TaskID: "target", RunID: "target-run",
		Attrs: map[string]any{"source_task": "source", "source_run": "source-run", "delivery": "delivery-id"},
	})
	for _, want := range []string{"task: target", "run: target-run", "source_task: source", "source_run: source-run", "delivery: delivery-id"} {
		if !strings.Contains(detail, want) {
			t.Fatalf("detail %q missing %q", detail, want)
		}
	}
}

func TestMergeLogEntriesPreservesFailedRunCorrelation(t *testing.T) {
	alert := domain.Alert{
		ID: "alert", TaskID: "task", RunID: "run", CreatedAt: time.Now(),
		Severity: domain.SeverityError, Kind: domain.AlertRunFailed, Message: "task run failed",
	}
	entries := mergeLogEntries(nil, []domain.Alert{alert}, "", time.Time{})
	if len(entries) != 1 || entries[0].taskID != "task" || entries[0].runID != "run" || !entries[0].runFailure {
		t.Fatalf("entry=%+v", entries)
	}
}

func TestFailedRunDetailIncludesExactActionableFields(t *testing.T) {
	exitCode := 7
	run := domain.Run{
		ID: "run-7", TaskID: "task-1", Trigger: domain.TriggerCompletion,
		Outcome: domain.OutcomeFailure, ExitCode: &exitCode,
		Output: "stdout line\nstderr line\n", OutputTruncated: true,
	}
	entry := logEntry{
		time: time.Date(2026, 9, 5, 1, 2, 3, 0, time.UTC), severity: domain.SeverityError,
		source: "alert: run_failed", message: "task run failed", taskID: "task-1",
		runID: "run-7", runFailure: true,
	}
	text := activityDetailText(entry, activityDiagnostic{run: &run, taskName: "Nightly backup"})
	for _, want := range []string{
		"Task: Nightly backup (task-1)", "Run: run-7", "Trigger: completion",
		"Outcome: failure", "Exit status: 7", "Output truncated: Yes",
		"Combined stdout/stderr:", "stdout line\nstderr line",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("detail %q missing %q", text, want)
		}
	}
}

func TestFailedRunDetailDistinguishesLaunchEmptyLegacyAndUnavailable(t *testing.T) {
	launch := domain.Run{ID: "run", TaskID: "task", Outcome: domain.OutcomeFailure, Trigger: domain.TriggerManual}
	text := activityDetailText(logEntry{taskID: "task", runID: "run", runFailure: true}, activityDiagnostic{run: &launch})
	for _, want := range []string{"No process exit status (launch or setup failed)", "Combined stdout/stderr:\n(empty)", "Output truncated: No"} {
		if !strings.Contains(text, want) {
			t.Errorf("launch detail %q missing %q", text, want)
		}
	}

	legacy := activityDetailText(logEntry{runFailure: true}, activityDiagnostic{})
	if !strings.Contains(legacy, "Run: Unavailable (legacy alert has no run identity)") {
		t.Fatalf("legacy detail=%q", legacy)
	}
	missing := activityDetailText(logEntry{taskID: "task", runID: "gone", runFailure: true}, activityDiagnostic{runUnavailable: true, taskUnavailable: true})
	for _, want := range []string{"Task: Unavailable (task may have been deleted) (task)", "Run: gone", "Run diagnostics: Unavailable"} {
		if !strings.Contains(missing, want) {
			t.Errorf("missing detail %q missing %q", missing, want)
		}
	}
}

func TestLoadActivityDiagnosticUsesOnlyExactIdentifiers(t *testing.T) {
	backend := &fakeBackend{
		tasks: []domain.Task{{ID: "task", Name: "Exact task"}},
		runs:  map[string]domain.Run{"run": {ID: "run", TaskID: "task", Outcome: domain.OutcomeFailure}},
	}
	ui := NewUI(testApp, backend)
	diagnostic := ui.loadActivityDiagnostic(context.Background(), logEntry{taskID: "task", runID: "run", runFailure: true})
	if diagnostic.run == nil || diagnostic.run.ID != "run" || diagnostic.taskName != "Exact task" {
		t.Fatalf("diagnostic=%+v", diagnostic)
	}
	backend.mu.Lock()
	gotIDs := append([]string(nil), backend.getRunIDs...)
	backend.mu.Unlock()
	if !reflect.DeepEqual(gotIDs, []string{"run"}) {
		t.Fatalf("run lookups=%v", gotIDs)
	}
}
