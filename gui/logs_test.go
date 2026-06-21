package gui

import (
	"testing"
	"time"

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

func TestUI_LogsTabBuilds(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{
		logs:   []domain.LogRecord{{ID: "l1", Severity: domain.SeverityError, Message: "boom"}},
		alerts: []domain.Alert{{ID: "a1", Severity: domain.SeverityWarning, Kind: domain.AlertRunFailed, Message: "warn"}},
	})
	if ui.logsTab == nil || ui.logsTab.Text != "Logs" {
		t.Fatalf("logs tab missing or mislabeled: %+v", ui.logsTab)
	}
}
