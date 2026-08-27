package gui

import (
	"context"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/domain"
)

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
	if ui.logsTab == nil || ui.logsTab.Text != "Activity" {
		t.Fatalf("activity tab missing or mislabeled: %+v", ui.logsTab)
	}
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
	walk(ui.logsTab.Content)

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
