package gui

import (
	"context"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"

	"github.com/shruggietech/go-scheduler/internal/api/server"
	"github.com/shruggietech/go-scheduler/internal/domain"
	"github.com/shruggietech/go-scheduler/internal/schedule"
	"github.com/shruggietech/go-scheduler/internal/timezone"
)

// taskEditor holds the widgets and live state of the New Task / Edit Task dialog.
// Its behavior methods (visibility, validation, previews, submission) are split
// out so they can be unit-tested headlessly without showing the dialog.
type taskEditor struct {
	app      *App
	existing *domain.Task // nil for a new task

	// What to run
	name    *widget.Entry
	command *widget.Entry
	args    *widget.Entry

	// When
	tz       *widget.SelectEntry
	mode     *widget.Select
	schedule *widget.Entry
	startAt  *widget.Entry // anchor time-of-day, sub-daily intervals only

	oneOffDate *widget.Entry
	oneOffTime *widget.Entry
	oneOffEcho *widget.Label

	schedPreview *widget.Label    // schedule summary + next runs (recurring only)
	cmdPreview   *widget.RichText // resolved command line as a code block
	whenForm     *widget.Form     // current When form (rebuilt fresh each time)
	whenHolder   *fyne.Container  // stable parent that holds the current whenForm

	// Right pane: shows the Preview by default, or Help when toggled.
	rightHolder    *fyne.Container
	previewContent fyne.CanvasObject
	helpContent    fyne.CanvasObject
	rightTitle     *widget.Label
	helpToggle     *cursorButton
	helpVisible    bool

	// Advanced
	overlap *widget.Select
	catchup *widget.Select

	save          *cursorButton
	cancelHandler func() // dismisses the dialog; nil in tests

	baseline    editorSnapshot // field values at open, for dirty detection
	ready       bool           // true once build() has wired the layout; gates OnChanged callbacks
	previewSync bool           // tests set this to fetch the schedule preview synchronously
}

// editorSnapshot captures the editor's field values so Cancel can detect unsaved
// changes (FR-011/FR-012).
type editorSnapshot struct {
	name, command, args, tz, mode string
	schedule, startAt             string
	oneOffDate, oneOffTime        string
	overlap, catchup              string
}

const (
	modeRecurring = "Recurring"
	modeOneOff    = "One-off"
)

// showTaskEditor opens the guided create/edit dialog. A live preview of both the
// schedule (plain-language summary + next runs) and the resolved command line is
// shown as the user types (FR-006/FR-007/FR-008). existing is nil for a new task.
func (a *App) showTaskEditor(existing *domain.Task) {
	e := newTaskEditor(a, existing)
	body := e.build()

	title := "New Task"
	if existing != nil {
		title = "Edit Task"
	}
	d := dialog.NewCustomWithoutButtons(title, body, a.win)
	e.save.OnTapped = func() {
		e.submit()
		d.Hide()
	}
	e.cancelHandler = d.Hide
	d.Resize(fyne.NewSize(1180, 720)) // ~2× width for the two-pane layout (FR-002)
	d.Show()
}

func newTaskEditor(a *App, existing *domain.Task) *taskEditor {
	e := &taskEditor{app: a, existing: existing}

	e.name = widget.NewEntry()
	e.command = widget.NewEntry()
	e.args = widget.NewMultiLineEntry()
	e.args.SetPlaceHolder("one argument per line")

	e.tz = widget.NewSelectEntry(commonZones)
	e.tz.SetText("Local")

	e.mode = widget.NewSelect([]string{modeRecurring, modeOneOff}, nil)

	e.schedule = widget.NewEntry()
	e.schedule.SetPlaceHolder(`e.g. "every 15 minutes" or "3rd wednesday monthly at 14:00"`)

	e.startAt = widget.NewEntry()
	e.startAt.SetPlaceHolder("e.g. 09:00 — aligns the first cycle")

	e.oneOffDate = widget.NewEntry()
	e.oneOffDate.SetPlaceHolder("2026-08-04")
	e.oneOffTime = widget.NewEntry()
	e.oneOffTime.SetPlaceHolder("09:00")
	e.oneOffEcho = widget.NewLabel("")
	e.oneOffEcho.Wrapping = fyne.TextWrapWord

	e.schedPreview = widget.NewLabel("")
	e.schedPreview.Wrapping = fyne.TextWrapWord
	e.cmdPreview = widget.NewRichText()
	e.cmdPreview.Wrapping = fyne.TextWrapWord

	e.overlap = widget.NewSelect(overlapLabels(), nil)
	e.overlap.SetSelected(overlapLabel(domain.OverlapQueueOne))
	e.catchup = widget.NewSelect(catchupLabels(), nil)
	e.catchup.SetSelected(catchupLabel(domain.CatchupOne))

	e.save = newCursorButton("Save", theme.ConfirmIcon(), widget.HighImportance, nil)

	e.wireValidators()
	e.prefill()
	return e
}

// --- construction --------------------------------------------------------

func (e *taskEditor) build() *fyne.Container {
	// Left pane: the form sections.
	runForm := widget.NewForm(
		requiredItem("Name", e.name),
		requiredItem("Command", e.command),
	)
	argsItem := widget.NewFormItem("Arguments", e.args)
	argsItem.HintText = "One argument per line" // persistent caption (FR-020)
	runForm.AppendItem(argsItem)

	// rebuildWhen swaps a freshly-built form into whenHolder on every change; a
	// fresh widget.Form (rather than mutating Items in place) guarantees every row
	// — including conditionally-shown ones like Start at — gets a renderer.
	e.whenHolder = container.NewStack()
	e.ready = true
	e.rebuildWhen()

	// Advanced Settings: custom collapsible (▶ collapsed / ▼ expanded), FR-009.
	advForm := widget.NewForm(
		widget.NewFormItem("Overlap", e.overlap),
		widget.NewFormItem("Catch-up", e.catchup),
	)
	advanced := newCollapsible("Advanced Settings", advForm)

	left := container.NewVScroll(container.NewVBox(
		sectionHeader("What to run"),
		runForm,
		widget.NewSeparator(),
		sectionHeader("When"),
		e.whenHolder,
		widget.NewSeparator(),
		advanced,
	))

	// Right pane: Preview (default) / Help.
	right := e.buildRightPane()

	// Two equal halves (FR-002/FR-003).
	split := container.NewGridWithColumns(2, left, right)

	// Footer: right-aligned Save/Cancel (FR-010); Cancel guards unsaved input.
	cancel := newCursorButton("Cancel", theme.CancelIcon(), widget.MediumImportance, e.requestCancel)
	footer := container.NewBorder(nil, nil, nil, container.NewHBox(cancel, e.save))

	e.baseline = e.snapshot()
	e.updatePreview()
	e.revalidate()
	return container.NewBorder(nil, footer, nil, nil, split)
}

// buildRightPane assembles the right half: a header with the pane title and the
// Help/Preview toggle, over a holder that swaps between the live Preview and the
// Help guidance (FR-003/FR-004).
func (e *taskEditor) buildRightPane() fyne.CanvasObject {
	e.previewContent = container.NewVScroll(container.NewVBox(e.schedPreview, e.cmdPreview))
	e.helpContent = helpView()
	e.helpContent.Hide()
	e.rightHolder = container.NewStack(e.previewContent, e.helpContent)

	e.rightTitle = sectionHeader("Preview")
	e.helpToggle = newCursorButton("Help", theme.HelpIcon(), widget.LowImportance, e.toggleHelp)
	header := container.NewBorder(nil, nil, e.rightTitle, e.helpToggle)
	return container.NewBorder(header, nil, nil, nil, e.rightHolder)
}

// toggleHelp swaps the right pane between Preview and Help without rebuilding the
// form, so inputs and the computed preview persist (FR-005).
func (e *taskEditor) toggleHelp() {
	e.helpVisible = !e.helpVisible
	if e.helpVisible {
		e.rightTitle.SetText("Help")
		e.helpToggle.SetText("Preview")
		e.previewContent.Hide()
		e.helpContent.Show()
	} else {
		e.rightTitle.SetText("Preview")
		e.helpToggle.SetText("Help")
		e.helpContent.Hide()
		e.previewContent.Show()
	}
}

// snapshot captures current field values for dirty detection.
func (e *taskEditor) snapshot() editorSnapshot {
	return editorSnapshot{
		name: e.name.Text, command: e.command.Text, args: e.args.Text, tz: e.tz.Text,
		mode: e.mode.Selected, schedule: e.schedule.Text, startAt: e.startAt.Text,
		oneOffDate: e.oneOffDate.Text, oneOffTime: e.oneOffTime.Text,
		overlap: e.overlap.Selected, catchup: e.catchup.Selected,
	}
}

// isDirty reports whether any field changed from its baseline at open (FR-011).
func (e *taskEditor) isDirty() bool { return e.snapshot() != e.baseline }

// requestCancel closes immediately for an untouched form, else confirms first
// (FR-011/FR-012).
func (e *taskEditor) requestCancel() {
	if !e.isDirty() {
		e.doCancel()
		return
	}
	dialog.NewConfirm("Discard changes?", "You have unsaved changes. Discard them?",
		func(ok bool) {
			if ok {
				e.doCancel()
			}
		}, e.app.win).Show()
}

func (e *taskEditor) doCancel() {
	if e.cancelHandler != nil {
		e.cancelHandler()
	}
}

// rebuildWhen recomputes the "When" form rows for the current Mode, showing only
// the relevant time inputs (FR-001) while preserving entered values (FR-002).
func (e *taskEditor) rebuildWhen() {
	items := []*widget.FormItem{
		widget.NewFormItem("Timezone", e.tz),
		widget.NewFormItem("Mode", e.mode),
	}
	if e.mode.Selected == modeOneOff {
		e.schedPreview.Hide()
		dateRow := container.NewBorder(nil, nil, nil, e.datePickerButton(), e.oneOffDate)
		items = append(items,
			requiredItem("Date", dateRow),
			requiredItem("Time", e.oneOffTime),
			withHint(widget.NewFormItem("", e.oneOffEcho), "Interpreted in the task's timezone"),
		)
	} else {
		e.schedPreview.Show()
		scheduleRow := container.NewBorder(nil, nil, nil, nil, e.schedule)
		schedItem := requiredItem("Schedule", scheduleRow)
		if e.existing != nil {
			schedItem = widget.NewFormItem("Schedule", scheduleRow) // optional on edit (blank = keep)
		}
		items = append(items, schedItem)
		if schedule.IsSubDailyInterval(e.effectiveScheduleRaw()) {
			startRow := container.NewBorder(nil, nil, nil, nil, e.startAt)
			items = append(items, withHint(widget.NewFormItem("Start at", startRow),
				"Optional anchor for the first cycle, e.g. 09:00"))
		}
	}
	e.whenForm = widget.NewForm(items...)
	e.whenHolder.Objects = []fyne.CanvasObject{e.whenForm}
	e.whenHolder.Refresh()
}

// --- validators & wiring -------------------------------------------------

func (e *taskEditor) wireValidators() {
	e.name.Validator = nonEmptyValidator("name")
	e.command.Validator = nonEmptyValidator("command")
	e.tz.Validator = func(s string) error {
		if _, err := timezone.Resolve(tzOrLocal(s)); err != nil {
			return err
		}
		return nil
	}

	e.mode.OnChanged = func(string) { e.onChange(true) }
	e.schedule.OnChanged = func(string) { e.onChange(true) }
	e.startAt.OnChanged = func(string) { e.onChange(false) }
	e.name.OnChanged = func(string) { e.onChange(false) }
	e.command.OnChanged = func(string) { e.updateCmdPreview(); e.onChange(false) }
	e.args.OnChanged = func(string) { e.updateCmdPreview(); e.onChange(false) }
	e.tz.OnChanged = func(string) { e.onChange(false) }
	e.oneOffDate.OnChanged = func(string) { e.updateOneOffEcho(); e.onChange(false) }
	e.oneOffTime.OnChanged = func(string) { e.updateOneOffEcho(); e.onChange(false) }
}

// onChange is the shared field-change handler. rebuild is true for changes that
// can alter which rows the When form shows (Mode, Schedule). It is a no-op until
// build() has wired the layout.
func (e *taskEditor) onChange(rebuild bool) {
	if !e.ready {
		return
	}
	if rebuild {
		e.rebuildWhen()
	}
	e.updatePreview()
	e.revalidate()
}

func (e *taskEditor) prefill() {
	if e.existing == nil {
		e.mode.SetSelected(modeRecurring)
		return
	}
	t := e.existing
	e.name.SetText(t.Name)
	e.command.SetText(t.Command)
	e.args.SetText(strings.Join(t.Args, "\n"))
	if t.Timezone != "" {
		e.tz.SetText(t.Timezone)
	}
	e.overlap.SetSelected(overlapLabel(t.OverlapPolicy))
	e.catchup.SetSelected(catchupLabel(t.CatchupPolicy))
	// The task row carries no schedule phrase; leave Schedule blank (an empty
	// schedule on update keeps the existing one) and default to Recurring.
	e.mode.SetSelected(modeRecurring)
}

// --- previews ------------------------------------------------------------

// updatePreview refreshes both the command-line preview and, in Recurring mode,
// the schedule summary. Invalid schedules render a warning synchronously; valid
// ones fetch the human summary and next runs from the backend asynchronously.
func (e *taskEditor) updatePreview() {
	e.updateCmdPreview()
	if e.mode.Selected != modeRecurring {
		return
	}
	s := e.effectiveSchedule()
	if s == "" {
		e.schedPreview.SetText("Type a schedule above to see upcoming runs")
		return
	}
	if _, err := schedule.Parse(s, e.tzName(), time.Now()); err != nil {
		e.schedPreview.SetText("⚠ " + cleanScheduleErr(err))
		return
	}
	if e.previewSync {
		e.fetchSchedulePreview(s)
		return
	}
	go e.fetchSchedulePreview(s)
}

// fetchSchedulePreview asks the backend for the human summary and next runs and
// renders them. Off the UI thread it marshals the update back via fyne.Do; when
// run synchronously (tests) it writes directly.
func (e *taskEditor) fetchSchedulePreview(s string) {
	ctx, cancel := e.app.bgCtx()
	defer cancel()
	resp, err := e.app.backend.Preview(ctx, server.PreviewRequest{Schedule: s, Timezone: e.tzName()})
	set := func() {
		if err != nil {
			e.schedPreview.SetText("⚠ " + cleanScheduleErr(err))
			return
		}
		txt := resp.HumanSummary
		for _, r := range resp.NextRuns {
			txt += "\n  • " + fmtTime(r)
		}
		e.schedPreview.SetText(txt)
	}
	if e.previewSync {
		set()
		return
	}
	fyne.Do(set)
}

// updateCmdPreview renders the resolved command line as a monospace code block
// with no prefix (FR-007/FR-008), or muted guidance text when empty.
func (e *taskEditor) updateCmdPreview() {
	line := commandLinePreview(e.command.Text, splitArgs(e.args.Text))
	if line == "" {
		e.cmdPreview.Segments = []widget.RichTextSegment{
			&widget.TextSegment{Style: widget.RichTextStyleInline, Text: "Enter a command to see what will run"},
		}
	} else {
		e.cmdPreview.Segments = []widget.RichTextSegment{
			&widget.TextSegment{Style: widget.RichTextStyleCodeBlock, Text: line},
		}
	}
	e.cmdPreview.Refresh()
}

// cmdPreviewString returns the current command-preview text (for tests).
func (e *taskEditor) cmdPreviewString() string {
	var b strings.Builder
	for _, s := range e.cmdPreview.Segments {
		if ts, ok := s.(*widget.TextSegment); ok {
			b.WriteString(ts.Text)
		}
	}
	return b.String()
}

func (e *taskEditor) updateOneOffEcho() {
	t, err := e.oneOffInstant()
	if err != nil {
		e.oneOffEcho.SetText("")
		return
	}
	e.oneOffEcho.SetText("Runs " + fmtTime(t))
}

// --- validation gating ---------------------------------------------------

// revalidate enables Save only when every currently-relevant field is valid
// (FR-003/FR-004/FR-005/FR-006). The relevant set depends on Mode and on whether
// this is a create (time field required) or an edit (time field optional).
func (e *taskEditor) revalidate() {
	if e.save == nil {
		return
	}
	if e.valid() {
		e.save.Enable()
	} else {
		e.save.Disable()
	}
}

func (e *taskEditor) valid() bool {
	if strings.TrimSpace(e.name.Text) == "" || strings.TrimSpace(e.command.Text) == "" {
		return false
	}
	if _, err := timezone.Resolve(e.tzName()); err != nil {
		return false
	}
	creating := e.existing == nil
	if e.mode.Selected == modeOneOff {
		t, err := e.oneOffInstant()
		switch {
		case err != nil:
			return !creating && e.oneOffBlank() // edit + blank = keep existing
		default:
			return t.After(time.Now())
		}
	}
	// Recurring.
	s := e.effectiveSchedule()
	if s == "" {
		return !creating // edit + blank = keep existing
	}
	_, err := schedule.Parse(s, e.tzName(), time.Now())
	return err == nil
}

// --- submission ----------------------------------------------------------

func (e *taskEditor) submit() { e.app.submitTask(e.existing, e.buildForm()) }

// buildForm collects the editor's current values into a taskForm, mapping the
// friendly overlap/catch-up labels back to their wire values and appending any
// interval anchor to the schedule phrase.
func (e *taskEditor) buildForm() taskForm {
	f := taskForm{
		name: e.name.Text, command: e.command.Text, args: splitArgs(e.args.Text),
		tz: e.tzName(), mode: e.mode.Selected, schedule: e.effectiveSchedule(),
		overlap: string(overlapValue(e.overlap.Selected)),
		catchup: string(catchupValue(e.catchup.Selected)),
	}
	if e.mode.Selected == modeOneOff {
		if t, err := e.oneOffInstant(); err == nil {
			f.at = t.Format(time.RFC3339)
		}
	}
	return f
}

// --- helpers -------------------------------------------------------------

func (e *taskEditor) tzName() string { return tzOrLocal(e.tz.Text) }

// effectiveScheduleRaw is the typed schedule without the GUI anchor appended;
// used to decide whether to offer the Start-at field.
func (e *taskEditor) effectiveScheduleRaw() string { return strings.TrimSpace(e.schedule.Text) }

// effectiveSchedule appends the optional "starting at <time>" anchor when the
// Start-at field is filled and applicable, giving one schedule phrase that both
// the preview and submit use (FR-010).
func (e *taskEditor) effectiveSchedule() string {
	s := e.effectiveScheduleRaw()
	at := strings.TrimSpace(e.startAt.Text)
	if at == "" || !schedule.IsSubDailyInterval(s) || containsAnchorClause(s) {
		return s
	}
	return s + " starting at " + at
}

func (e *taskEditor) oneOffBlank() bool {
	return strings.TrimSpace(e.oneOffDate.Text) == "" && strings.TrimSpace(e.oneOffTime.Text) == ""
}

// oneOffInstant assembles the date + time entries into an instant in the task's
// timezone (FR-015). Both fields must be present and well-formed.
func (e *taskEditor) oneOffInstant() (time.Time, error) {
	loc, err := timezone.Resolve(e.tzName())
	if err != nil {
		return time.Time{}, err
	}
	date := strings.TrimSpace(e.oneOffDate.Text)
	tod := strings.TrimSpace(e.oneOffTime.Text)
	return time.ParseInLocation("2006-01-02 15:04", date+" "+tod, loc)
}

// datePickerButton opens a graphical month calendar; choosing a day fills the
// one-off Date field, so the user need not type the date by hand (FR-015).
func (e *taskEditor) datePickerButton() *cursorButton {
	return newCursorButton("Pick…", theme.MoreHorizontalIcon(), widget.LowImportance, e.showDatePicker)
}

func (e *taskEditor) showDatePicker() {
	start := time.Now()
	if t, err := e.oneOffInstant(); err == nil {
		start = t
	}
	var d dialog.Dialog
	cal := xwidget.NewCalendar(start, func(t time.Time) {
		e.oneOffDate.SetText(t.Format("2006-01-02"))
		if d != nil {
			d.Hide()
		}
	})
	d = dialog.NewCustom("Pick a date", "Close", cal, e.app.win)
	d.Show()
}

// containsAnchorClause reports whether the phrase already carries an anchor, so
// the GUI doesn't append a second one.
func containsAnchorClause(s string) bool {
	l := strings.ToLower(s)
	return strings.Contains(l, "starting at") || strings.Contains(l, " from ")
}

func cleanScheduleErr(err error) string {
	return strings.TrimPrefix(err.Error(), "schedule: ")
}

func sectionHeader(text string) *widget.Label {
	return widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
}

func requiredItem(label string, w fyne.CanvasObject) *widget.FormItem {
	return widget.NewFormItem(label+" *", w)
}

func withHint(item *widget.FormItem, hint string) *widget.FormItem {
	item.HintText = hint
	return item
}

func nonEmptyValidator(field string) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return errEmptyField(field)
		}
		return nil
	}
}

// helpView is the in-modal Help content: a field-by-field guide with examples
// (FR-004). It replaces the old per-field Examples popup.
func helpView() fyne.CanvasObject {
	md := `## Task editor help

**Name** — a label to identify the task. _e.g._ ` + "`nightly-backup`" + `

**Command** — the program to run (just the executable, not a full command line).
_e.g._ ` + "`cmd`" + `, ` + "`python`" + `, ` + "`C:\\Windows\\System32\\notepad.exe`" + `

**Arguments** — one argument per line; each line is passed as a separate argument.
For ` + "`cmd /c echo hi`" + ` enter ` + "`/c`" + ` on one line and ` + "`echo hi`" + ` on the next.

**Timezone** — ` + "`Local`" + ` or an IANA name (` + "`UTC`" + `, ` + "`America/New_York`" + `).
Schedules are interpreted here; storage is UTC with DST handled.

**Mode** — _Recurring_ fires repeatedly on a Schedule; _One-off_ fires once at a date+time.

**Schedule** _(Recurring)_ — a plain-language phrase:
- Intervals: ` + "`every 15 minutes`" + `, ` + "`every 30s`" + `, ` + "`every 2 hours`" + `, ` + "`every 3 days`" + `, ` + "`every week`" + `
- Daily with a time: ` + "`every day at 09:00`" + `
- Weekday/weekend sets: ` + "`weekdays at 9:00 AM`" + `, ` + "`weekends at 18:00`" + `
- A single weekday: ` + "`every monday at 9am`" + `
- Monthly ordinals: ` + "`3rd wednesday monthly at 14:00`" + `, ` + "`last friday of the month`" + `

**Start at** _(sub-daily intervals)_ — aligns the first cycle. ` + "`every 15 minutes`" + ` with
Start at ` + "`09:00`" + ` runs at :00/:15/:30/:45. You can also type it inline:
` + "`every 15 minutes starting at 09:00`" + `.

**One-off date / time** — pick a future date and time (in the task's timezone); use the calendar
button to choose the date.

**Overlap** — what to do if a run is still going when the next is due: _Queue one run_ (default),
_Skip this run_, or _Allow concurrent runs_.

**Catch-up** — after downtime: _Run once to catch up_ (default) or _Skip missed runs_.`
	r := widget.NewRichTextFromMarkdown(md)
	r.Wrapping = fyne.TextWrapWord
	return container.NewVScroll(r)
}

// taskForm carries the submitted values from the editor to submitTask.
type taskForm struct {
	name, command, tz, mode, schedule, at, overlap, catchup string
	args                                                    []string
}

func (a *App) submitTask(existing *domain.Task, f taskForm) {
	var atPtr *time.Time
	if f.mode == modeOneOff {
		ts, err := time.Parse(time.RFC3339, strings.TrimSpace(f.at))
		if err != nil {
			a.showError(errInvalidOneOff)
			return
		}
		atPtr = &ts
	}

	a.run(func(ctx context.Context) error {
		if existing == nil {
			req := server.TaskCreateRequest{
				Name: f.name, Command: f.command, Args: f.args, Timezone: f.tz,
				OverlapPolicy: f.overlap, CatchupPolicy: f.catchup,
			}
			if atPtr != nil {
				req.At = atPtr
			} else {
				req.Schedule = f.schedule
			}
			_, err := a.backend.CreateTask(ctx, req)
			return err
		}
		req := server.TaskUpdateRequest{
			Name: f.name, Command: f.command, Args: f.args, Timezone: f.tz,
			OverlapPolicy: f.overlap, CatchupPolicy: f.catchup,
		}
		if atPtr != nil {
			req.At = atPtr
		} else {
			req.Schedule = f.schedule
		}
		_, err := a.backend.UpdateTask(ctx, existing.ID, req)
		return err
	})
}

func splitArgs(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		if t := strings.TrimSpace(line); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func tzOrLocal(s string) string {
	if strings.TrimSpace(s) == "" {
		return "Local"
	}
	return s
}
