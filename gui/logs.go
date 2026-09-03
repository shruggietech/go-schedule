package gui

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// logEntry is a unified row in the Activity view: either a daemon log record or a
// scheduler alert, normalized to a common shape.
type logEntry struct {
	id       string
	time     time.Time
	severity domain.AlertSeverity
	source   string
	message  string
	detail   string
	isAlert  bool
	alertID  string
}

var activityColumns = []structuredColumn{
	{Header: "When", Minimum: 145},
	{Header: "Severity", Minimum: 110, Alignment: fyne.TextAlignCenter},
	{Header: "Source", Minimum: 150, Weight: 1},
	{Header: "Summary", Minimum: 200, Weight: 3},
}

// buildLogsTab shows a unified, filterable Activity view that merges daemon log
// records and scheduler alerts (FR-011). It supports severity filters (FR-013),
// click-through detail (FR-014), and clearing the current view (FR-015), and updates live from
// the event stream (FR-018).
func (a *App) buildLogsTab() fyne.CanvasObject {
	var rows []logEntry
	filter := domain.AlertSeverity("") // "" = all
	var clearedAt time.Time            // Clear View cutoff

	var table *structuredList
	table = newStructuredList(activityColumns, "", func(identity string) {
		entry, ok := activityEntryForIdentity(rows, table.rows, identity)
		if !ok {
			return
		}
		a.showLogDetail(entry)
		table.list.UnselectAll()
	}, nil)
	a.activityTable = table
	diagnostics := widget.NewLabel(activityDiagnosticsText(""))
	diagnostics.Wrapping = fyne.TextWrapBreak

	rebuild := func() {
		snap := a.model.Snapshot()
		rows = mergeLogEntries(snap.Logs, snap.Alerts, filter, clearedAt)
		diagnostics.SetText(activityDiagnosticsText(snap.LogPath))
		table.setRows(activityRowModels(rows))
	}
	a.registerRefresher(rebuild)

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

	clearBtn := newToolbarButton("Clear View", theme.ContentClearIcon(), func() {
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

	toolbar := container.NewHBox(widget.NewLabel("Severity:"), severitySel, clearBtn)
	help := widget.NewLabel("Hides current activity and acknowledges visible alerts. Records are not deleted.")
	help.Wrapping = fyne.TextWrapWord
	return container.NewBorder(container.NewVBox(toolbar, diagnostics, help), nil, nil, nil, table.root)
}

func activityDiagnosticsText(logPath string) string {
	const summary = "This view shows a limited set of recent daemon log records plus scheduler alerts; older daemon records remain in the full log.\n"
	if logPath == "" {
		return summary + "Full daemon log: unavailable until daemon responds."
	}
	return summary + "Full daemon log: " + logPath
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
	d := dialog.NewCustom("Activity detail", "Close", entry, a.win)
	d.Resize(fyne.NewSize(560, 360))
	d.Show()
}

// mergeLogEntries combines log records and alerts into a single severity-filtered,
// newest-first list, dropping entries at or before the Clear View cutoff.
func mergeLogEntries(logs []domain.LogRecord, alerts []domain.Alert, filter domain.AlertSeverity, clearedAt time.Time) []logEntry {
	out := make([]logEntry, 0, len(logs)+len(alerts))
	for _, l := range logs {
		out = append(out, logEntry{
			id:   l.ID,
			time: l.Time, severity: l.Severity, source: srcOr(l.Source, "daemon"),
			message: l.Message, detail: attrsDetail(l),
		})
	}
	for _, al := range alerts {
		out = append(out, logEntry{
			id:   al.ID,
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

func activityRowModels(entries []logEntry) []structuredRowModel {
	rows := make([]structuredRowModel, len(entries))
	identities := make(map[string]int, len(entries))
	for index, entry := range entries {
		severityText, importance := activitySeverityPresentation(entry.severity)
		source := srcOr(entry.source, "daemon")
		message := entry.message
		if strings.TrimSpace(message) == "" {
			message = "No message"
		}
		kind := "log"
		if entry.isAlert {
			kind = "alert"
		}
		identityParts := []string{kind, entry.id}
		if entry.id == "" {
			identityParts = append(identityParts,
				entry.time.Format(time.RFC3339Nano),
				string(entry.severity),
				source,
				message,
			)
		}
		baseIdentity := stablePresentationIdentity(identityParts...)
		ordinal := identities[baseIdentity]
		identities[baseIdentity] = ordinal + 1
		identity := stablePresentationIdentity(baseIdentity, strconv.Itoa(ordinal))
		cells := []structuredCell{
			{Text: fmtTime(entry.time), TextStyle: fyne.TextStyle{Monospace: true}},
			{Text: severityText, Importance: importance},
			{Text: source, FullText: source},
			{Text: message, FullText: message},
		}
		rows[index] = structuredRowModel{
			Identity: identity,
			Cells:    cells,
			Summary:  structuredRowSummary(activityColumns, cells),
		}
	}
	return rows
}

func activitySeverityPresentation(severity domain.AlertSeverity) (string, widget.Importance) {
	switch severity {
	case "", domain.SeverityInfo:
		return "• INFO", widget.HighImportance
	case domain.SeverityWarning:
		return "⚠ WARNING", widget.WarningImportance
	case domain.SeverityError:
		return "✗ ERROR", widget.DangerImportance
	default:
		return "? " + normalizedWords(string(severity), "UNKNOWN", true), widget.MediumImportance
	}
}

func activityEntryForIdentity(entries []logEntry, rows []structuredRowModel, identity string) (logEntry, bool) {
	for index, row := range rows {
		if row.Identity == identity && index < len(entries) {
			return entries[index], true
		}
	}
	return logEntry{}, false
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
