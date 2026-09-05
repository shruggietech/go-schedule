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
	id         string
	time       time.Time
	severity   domain.AlertSeverity
	source     string
	message    string
	detail     string
	isAlert    bool
	alertID    string
	taskID     string
	runID      string
	runFailure bool
}

type activityDiagnostic struct {
	run             *domain.Run
	taskName        string
	runUnavailable  bool
	taskUnavailable bool
	loading         bool
}

var activityColumns = []structuredColumn{
	{Header: "When", Minimum: 145, Preferred: 26},
	{Header: "Severity", Minimum: 110, Preferred: 16, Alignment: fyne.TextAlignCenter},
	{Header: "Source", Minimum: 150, Weight: 1, Preferred: 22},
	{Header: "Summary", Minimum: 200, Weight: 3, Preferred: 36},
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
	table = newAdjustableStructuredList(activityColumns, "", func(identity string) {
		entry, ok := activityEntryForIdentity(rows, table.rows, identity)
		if !ok {
			return
		}
		a.showLogDetail(entry)
		table.list.UnselectAll()
	}, nil, a.fyne.Preferences(), activityColumnLayoutPreferenceKey)
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

	toolbar := container.NewHBox(widget.NewLabel("Severity:"), severitySel, clearBtn, widget.NewButton("Reset columns", table.resetColumns))
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
	diagnostic := activityDiagnostic{loading: e.runFailure && e.runID != ""}
	entry := widget.NewMultiLineEntry()
	entry.SetText(activityDetailText(e, diagnostic))
	entry.Wrapping = fyne.TextWrapWord
	d := dialog.NewCustom("Activity detail", "Close", entry, a.win)
	d.Resize(fyne.NewSize(620, 440))
	d.Show()
	if !e.runFailure {
		return
	}
	go func() {
		ctx, cancel := a.bgCtx()
		defer cancel()
		loaded := a.loadActivityDiagnostic(ctx, e)
		fyne.Do(func() { entry.SetText(activityDetailText(e, loaded)) })
	}()
}

func (a *App) loadActivityDiagnostic(ctx context.Context, e logEntry) activityDiagnostic {
	var diagnostic activityDiagnostic
	if e.taskID != "" {
		task, err := a.backend.GetTask(ctx, e.taskID)
		if err != nil {
			diagnostic.taskUnavailable = true
		} else {
			diagnostic.taskName = task.Task.Name
		}
	}
	if e.runID != "" {
		run, err := a.backend.GetRun(ctx, e.runID)
		if err != nil || run.ID != e.runID || e.taskID != "" && run.TaskID != e.taskID {
			diagnostic.runUnavailable = true
		} else {
			diagnostic.run = &run
		}
	}
	return diagnostic
}

func activityDetailText(e logEntry, diagnostic activityDiagnostic) string {
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
	if !e.runFailure {
		return b.String()
	}
	b.WriteString("\nFailed run\n")
	switch {
	case diagnostic.taskName != "":
		fmt.Fprintf(&b, "Task: %s (%s)\n", diagnostic.taskName, valueOr(e.taskID, "Unavailable"))
	case diagnostic.taskUnavailable:
		fmt.Fprintf(&b, "Task: Unavailable (task may have been deleted) (%s)\n", valueOr(e.taskID, "legacy alert has no task identity"))
	default:
		fmt.Fprintf(&b, "Task: %s\n", valueOr(e.taskID, "Unavailable (legacy alert has no task identity)"))
	}
	if e.runID == "" {
		b.WriteString("Run: Unavailable (legacy alert has no run identity)\n")
		return b.String()
	}
	fmt.Fprintf(&b, "Run: %s\n", e.runID)
	switch {
	case diagnostic.loading:
		b.WriteString("Run diagnostics: Loading…\n")
		return b.String()
	case diagnostic.runUnavailable || diagnostic.run == nil:
		b.WriteString("Run diagnostics: Unavailable (run record may have been deleted)\n")
		return b.String()
	}
	run := diagnostic.run
	fmt.Fprintf(&b, "Trigger: %s\n", run.Trigger)
	if run.SourceTriggerID != "" {
		fmt.Fprintf(&b, "Source trigger: %s\n", run.SourceTriggerID)
	}
	fmt.Fprintf(&b, "Outcome: %s\n", run.Outcome)
	if run.ExitCode == nil {
		b.WriteString("Exit status: No process exit status (launch or setup failed)\n")
	} else {
		fmt.Fprintf(&b, "Exit status: %d\n", *run.ExitCode)
	}
	fmt.Fprintf(&b, "Output truncated: %s\n", yesNo(run.OutputTruncated))
	b.WriteString("\nCombined stdout/stderr:\n")
	if run.Output == "" {
		b.WriteString("(empty)\n")
	} else {
		b.WriteString(run.Output)
		if !strings.HasSuffix(run.Output, "\n") {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func valueOr(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func yesNo(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

// mergeLogEntries combines log records and alerts into a single severity-filtered,
// newest-first list, dropping entries at or before the Clear View cutoff.
func mergeLogEntries(logs []domain.LogRecord, alerts []domain.Alert, filter domain.AlertSeverity, clearedAt time.Time) []logEntry {
	out := make([]logEntry, 0, len(logs)+len(alerts))
	for _, l := range logs {
		out = append(out, logEntry{
			id:   l.ID,
			time: l.Time, severity: l.Severity, source: srcOr(l.Source, "daemon"),
			message: l.Message, detail: attrsDetail(l), taskID: l.TaskID, runID: l.RunID,
		})
	}
	for _, al := range alerts {
		out = append(out, logEntry{
			id:   al.ID,
			time: al.CreatedAt, severity: al.Severity, source: "alert: " + string(al.Kind),
			message: al.Message, isAlert: true, alertID: al.ID, taskID: al.TaskID, runID: al.RunID,
			runFailure: al.Kind == domain.AlertRunFailed,
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
