package gui

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// logEntry is a unified row in the Logs view: either a daemon log record or a
// scheduler alert, normalized to a common shape.
type logEntry struct {
	time     time.Time
	severity domain.AlertSeverity
	source   string
	message  string
	detail   string
	isAlert  bool
	alertID  string
}

// buildLogsTab shows a unified, filterable Logs view that merges daemon log
// records and scheduler alerts (FR-011). It supports severity filters (FR-013),
// click-through detail (FR-014), and Dismiss All (FR-015), and updates live from
// the event stream (FR-018).
func (a *App) buildLogsTab() fyne.CanvasObject {
	var rows []logEntry
	filter := domain.AlertSeverity("") // "" = all
	var clearedAt time.Time            // Dismiss All cutoff

	list := widget.NewList(
		func() int { return len(rows) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			e := rows[i]
			o.(*widget.Label).SetText(fmt.Sprintf("%s  %s  [%s]  %s",
				severityMark(e.severity), fmtTime(e.time), e.source, e.message))
		},
	)

	rebuild := func() {
		snap := a.model.Snapshot()
		rows = mergeLogEntries(snap.Logs, snap.Alerts, filter, clearedAt)
		list.Refresh()
	}
	a.registerRefresher(rebuild)

	list.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(rows) {
			return
		}
		a.showLogDetail(rows[id])
		list.UnselectAll()
	}

	severitySel := widget.NewSelect(
		[]string{"All", "Errors", "Warnings", "Info"},
		func(s string) {
			switch s {
			case "Errors":
				filter = domain.SeverityError
			case "Warnings":
				filter = domain.SeverityWarning
			case "Info":
				filter = domain.SeverityInfo
			default:
				filter = ""
			}
			rebuild()
		},
	)
	severitySel.SetSelected("All")

	dismissBtn := newToolbarButton("Dismiss All", theme.DeleteIcon(), func() {
		clearedAt = time.Now()
		// Acknowledge the alerts currently shown so they don't reappear unacked.
		ids := make([]string, 0)
		for _, e := range rows {
			if e.isAlert && e.alertID != "" {
				ids = append(ids, e.alertID)
			}
		}
		a.run(func(ctx context.Context) error {
			for _, id := range ids {
				if err := a.backend.AckAlert(ctx, id); err != nil {
					return err
				}
			}
			return nil
		})
		rebuild()
	})

	toolbar := container.NewHBox(widget.NewLabel("Severity:"), severitySel, dismissBtn)
	return container.NewBorder(toolbar, nil, nil, nil, list)
}

// showLogDetail opens a dialog with the full message and cause/context of an entry.
func (a *App) showLogDetail(e logEntry) {
	var b strings.Builder
	fmt.Fprintf(&b, "Severity: %s\n", e.severity)
	fmt.Fprintf(&b, "Time:     %s\n", e.time.Format(time.RFC3339))
	if e.source != "" {
		fmt.Fprintf(&b, "Source:   %s\n", e.source)
	}
	fmt.Fprintf(&b, "\n%s\n", e.message)
	if e.detail != "" {
		fmt.Fprintf(&b, "\n%s\n", e.detail)
	}
	entry := widget.NewMultiLineEntry()
	entry.SetText(b.String())
	entry.Wrapping = fyne.TextWrapWord
	d := dialog.NewCustom("Log detail", "Close", entry, a.win)
	d.Resize(fyne.NewSize(560, 360))
	d.Show()
}

// mergeLogEntries combines log records and alerts into a single severity-filtered,
// newest-first list, dropping entries at or before clearedAt (Dismiss All cutoff).
func mergeLogEntries(logs []domain.LogRecord, alerts []domain.Alert, filter domain.AlertSeverity, clearedAt time.Time) []logEntry {
	out := make([]logEntry, 0, len(logs)+len(alerts))
	for _, l := range logs {
		out = append(out, logEntry{
			time: l.Time, severity: l.Severity, source: srcOr(l.Source, "daemon"),
			message: l.Message, detail: attrsDetail(l),
		})
	}
	for _, al := range alerts {
		out = append(out, logEntry{
			time: al.CreatedAt, severity: al.Severity, source: "alert: " + string(al.Kind),
			message: al.Message, isAlert: true, alertID: al.ID,
		})
	}
	filtered := out[:0]
	for _, e := range out {
		if !clearedAt.IsZero() && !e.time.After(clearedAt) {
			continue
		}
		if filter != "" && e.severity != filter {
			continue
		}
		filtered = append(filtered, e)
	}
	sort.SliceStable(filtered, func(i, j int) bool { return filtered[i].time.After(filtered[j].time) })
	return filtered
}

func severityMark(s domain.AlertSeverity) string {
	switch s {
	case domain.SeverityError:
		return "✗ ERROR  "
	case domain.SeverityWarning:
		return "⚠ WARN   "
	default:
		return "• info   "
	}
}

func srcOr(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

// attrsDetail renders a log record's structured attributes as a readable block.
func attrsDetail(l domain.LogRecord) string {
	if len(l.Attrs) == 0 && l.TaskID == "" && l.RunID == "" {
		return ""
	}
	var b strings.Builder
	if l.TaskID != "" {
		fmt.Fprintf(&b, "task: %s\n", l.TaskID)
	}
	if l.RunID != "" {
		fmt.Fprintf(&b, "run: %s\n", l.RunID)
	}
	for k, v := range l.Attrs {
		fmt.Fprintf(&b, "%s: %v\n", k, v)
	}
	return strings.TrimRight(b.String(), "\n")
}
