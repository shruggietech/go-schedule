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
	"github.com/shruggietech/go-schedule/internal/timezone"
)

// newTestEditor builds a wired editor (ready == true) against a fake backend.
// existing is nil for a create; a task alone stands in for an edit whose
// schedule is irrelevant to the test.
func newTestEditor(t *testing.T, existing *domain.Task) (*taskEditor, *fakeBackend) {
	t.Helper()
	var detail *server.TaskResponse
	if existing != nil {
		detail = &server.TaskResponse{Task: *existing}
	}
	return newTestEditorDetail(t, detail)
}

// newTestEditorDetail builds a wired editor from full task detail, for tests
// that care about the stored schedule.
func newTestEditorDetail(t *testing.T, detail *server.TaskResponse) (*taskEditor, *fakeBackend) {
	t.Helper()
	fb := &fakeBackend{}
	ui := NewUI(testApp, fb)
	e := newTaskEditor(ui, detail)
	e.previewSync = true // deterministic: no cross-test goroutines/fyne.Do
	e.build()            // wires layout, sets ready
	return e, fb
}

func whenLabels(e *taskEditor) []string {
	out := make([]string, len(e.whenForm.Items))
	for i, it := range e.whenForm.Items {
		out[i] = it.Text
	}
	return out
}

func hasLabel(labels []string, want string) bool {
	for _, l := range labels {
		if l == want {
			return true
		}
	}
	return false
}

// --- US1: mode-driven visibility -----------------------------------------

func TestEditor_ModeVisibility(t *testing.T) {
	e, _ := newTestEditor(t, nil)

	labels := whenLabels(e)
	if !hasLabel(labels, "Schedule") {
		t.Fatalf("Recurring mode missing Schedule row: %v", labels)
	}
	if hasLabel(labels, "Date *") || hasLabel(labels, "Time *") {
		t.Fatalf("Recurring mode should not show one-off rows: %v", labels)
	}

	e.mode.SetSelected(modeOneOff)
	labels = whenLabels(e)
	if !hasLabel(labels, "Date *") || !hasLabel(labels, "Time *") {
		t.Fatalf("One-off mode missing date/time rows: %v", labels)
	}
	if hasLabel(labels, "Schedule") {
		t.Fatalf("One-off mode should not show Schedule row: %v", labels)
	}
}

func TestEditor_ModeTogglePreservesValues(t *testing.T) {
	e, _ := newTestEditor(t, nil)
	e.schedule.SetText("every 15 minutes")
	e.oneOffDate.SetText("2099-01-02")

	e.mode.SetSelected(modeOneOff)
	e.mode.SetSelected(modeRecurring)

	if e.schedule.Text != "every 15 minutes" {
		t.Fatalf("schedule lost on toggle: %q", e.schedule.Text)
	}
	if e.oneOffDate.Text != "2099-01-02" {
		t.Fatalf("one-off date lost on toggle: %q", e.oneOffDate.Text)
	}
}

// --- US2: validation gating ----------------------------------------------

func TestEditor_SaveGating(t *testing.T) {
	e, _ := newTestEditor(t, nil)

	if e.save.Disabled() {
		t.Fatal("Save should allow an empty draft")
	}

	e.name.SetText("nightly")
	e.commandLine.SetText("cmd")
	if e.save.Disabled() {
		t.Fatal("Save should allow a manual-only draft without a schedule")
	}

	e.schedule.SetText("every 15 minutes")
	if e.save.Disabled() {
		t.Fatal("Save should be enabled with name+command+valid schedule")
	}

	e.schedule.SetText("nonsense gibberish")
	if !e.save.Disabled() {
		t.Fatal("Save should disable for an unparseable schedule")
	}
}

func TestEditor_SaveGating_OneOff(t *testing.T) {
	e, _ := newTestEditor(t, nil)
	e.name.SetText("once")
	e.commandLine.SetText("cmd")
	e.mode.SetSelected(modeOneOff)

	e.oneOffDate.SetText("2000-01-01")
	e.oneOffTime.SetText("09:00")
	if !e.save.Disabled() {
		t.Fatal("Save should disable for a past one-off time")
	}

	e.oneOffDate.SetText("2099-01-01")
	if e.save.Disabled() {
		t.Fatal("Save should enable for a future one-off time")
	}
}

func TestEditor_SaveGating_BadTimezone(t *testing.T) {
	e, _ := newTestEditor(t, nil)
	e.name.SetText("x")
	e.commandLine.SetText("cmd")
	e.schedule.SetText("every 15 minutes")
	if e.save.Disabled() {
		t.Fatal("precondition: Save enabled")
	}
	e.tz.SetText("Mars/Phobos")
	if !e.save.Disabled() {
		t.Fatal("Save should disable for an unknown timezone")
	}
}

func TestEditor_SaveGating_ElapsedRejectsCalendarSchedule(t *testing.T) {
	e, _ := newTestEditor(t, nil)
	e.name.SetText("monthly")
	e.commandLine.SetText("cmd")
	e.schedule.SetText("3rd wednesday monthly at 14:00")
	if e.save.Disabled() {
		t.Fatal("precondition: wall-clock monthly schedule should be valid")
	}
	e.timeBasis.SetSelected(timeBasisLabel(domain.TimeBasisElapsed))
	if !e.save.Disabled() {
		t.Fatal("Save should disable when elapsed time is incompatible with the recurrence")
	}
}

func TestEditor_CommandFieldIsOneRoomyGrowingMultilineInput(t *testing.T) {
	e, _ := newTestEditor(t, nil)
	if e.commandLine == nil || !e.commandLine.MultiLine {
		t.Fatal("Command line must be one multiline entry")
	}
	defaultThreeRows := widget.NewMultiLineEntry().MinSize().Height
	if got := e.commandLine.MinSize().Height; got < defaultThreeRows*1.7 {
		t.Fatalf("Command line minimum height = %.1f, want at least six visible rows (default three = %.1f)", got, defaultThreeRows)
	}
	before := e.commandLine.Size().Height
	e.leftContent.Resize(fyne.NewSize(e.leftContent.MinSize().Width, e.leftContent.MinSize().Height+120))
	if after := e.commandLine.Size().Height; after <= before {
		t.Fatalf("Command line height did not grow with available dialog content: %.1f -> %.1f", before, after)
	}
}

func TestEditor_CommandValidationAndExactCreateSubmission(t *testing.T) {
	e, fb := newTestEditor(t, nil)
	e.name.SetText("portable")
	e.schedule.SetText("every day at 09:00")
	e.commandLine.SetText(`"C:\Program Files\Tool\tool.exe" --tag one --tag two '' "héllo 世界"`)
	if e.save.Disabled() {
		t.Fatal("valid portable command line should enable Save when the rest of the form is valid")
	}

	e.submit()
	waitFor(t, func() bool { n, _ := fb.lastCreateCall(); return n == 1 })
	_, req := fb.lastCreateCall()
	if req.Command != `C:\Program Files\Tool\tool.exe` {
		t.Fatalf("created command = %q", req.Command)
	}
	wantArgs := []string{"--tag", "one", "--tag", "two", "", "héllo 世界"}
	if !reflect.DeepEqual(req.Args, wantArgs) {
		t.Fatalf("created args = %#v, want %#v", req.Args, wantArgs)
	}
}

func TestEditorCreationActivationDefaultsClearedAndPersistsThroughValidation(t *testing.T) {
	e, _ := newTestEditor(t, nil)
	if e.activateAfterSave == nil || e.activateAfterSave.Checked {
		t.Fatalf("fresh activation choice=%v, want visible and cleared", e.activateAfterSave)
	}
	e.activateAfterSave.SetChecked(true)
	e.name.SetText("")
	e.commandLine.SetText(`program "unclosed`)
	if !e.activateAfterSave.Checked {
		t.Fatal("validation changes reset the activation choice")
	}
	if !e.isDirty() {
		t.Fatal("changing activation intent did not mark the creation draft dirty")
	}
}

func TestEditorCreationSubmitsExplicitActivationIntentAndFreshEditorResets(t *testing.T) {
	for _, test := range []struct {
		name string
		set  bool
	}{
		{name: "inactive default", set: false},
		{name: "active opt in", set: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			e, backend := newTestEditor(t, nil)
			e.name.SetText("safe")
			e.commandLine.SetText("program")
			e.schedule.SetText("every day at 09:00")
			e.activateAfterSave.SetChecked(test.set)
			e.submit()
			waitFor(t, func() bool { count, _ := backend.lastCreateCall(); return count == 1 })
			_, request := backend.lastCreateCall()
			if request.Enabled == nil || *request.Enabled != test.set {
				t.Fatalf("create enabled intent=%v, want %v", request.Enabled, test.set)
			}
		})
	}
	fresh, _ := newTestEditor(t, nil)
	if fresh.activateAfterSave.Checked {
		t.Fatal("fresh creation editor retained prior opt-in")
	}
}

func TestEditorExistingTaskHasNoCreationActivationChoice(t *testing.T) {
	existing := &domain.Task{ID: "task", Name: "existing", Enabled: true, State: domain.TaskActive}
	e, backend := newTestEditor(t, existing)
	if e.activateAfterSave != nil {
		t.Fatal("edit editor exposed creation-only activation choice")
	}
	e.name.SetText("renamed")
	e.submit()
	waitFor(t, func() bool { count, _, _ := backend.lastUpdateCall(); return count == 1 })
	if existing.Enabled != true {
		t.Fatal("editing mutated enabled state")
	}
}

func TestEditor_InvalidCommandDisablesSave(t *testing.T) {
	e, _ := newTestEditor(t, nil)
	e.name.SetText("invalid")
	e.schedule.SetText("every day at 09:00")
	e.commandLine.SetText(`program "unclosed`)
	if !e.save.Disabled() {
		t.Fatal("unmatched quote must disable Save")
	}
}

// --- US3: combined preview -----------------------------------------------

func TestEditor_CommandPreview(t *testing.T) {
	e, _ := newTestEditor(t, nil)
	e.commandLine.SetText(`cmd /c "echo hello world"`)
	got := e.cmdPreviewString()
	for _, want := range []string{"Program", `"cmd"`, "Arguments in order (2)", `1  "/c"`, `2  "echo hello world"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("cmd preview = %q, missing %q", got, want)
		}
	}
	if strings.Contains(got, `cmd /c "echo hello world"`) {
		t.Fatalf("cmd preview = %q", got)
	}
	if strings.Contains(got, "Will run") {
		t.Fatalf("cmd preview should not have a 'Will run:' prefix: %q", got)
	}
}

func TestEditor_InvalidCommandClearsStalePreviewAndRecovers(t *testing.T) {
	e, _ := newTestEditor(t, nil)
	e.commandLine.SetText("program valid")
	if got := e.cmdPreviewString(); !strings.Contains(got, "Arguments in order (1)") {
		t.Fatalf("valid preview = %q", got)
	}
	e.commandLine.SetText("program\n'unclosed")
	got := e.cmdPreviewString()
	if strings.Contains(got, "Arguments in order") || !strings.Contains(got, "single quote opened at line 2, column 1") {
		t.Fatalf("invalid preview = %q", got)
	}
	e.commandLine.SetText(`program '$HOME' '|' '*.txt'`)
	got = e.cmdPreviewString()
	for _, literal := range []string{`"$HOME"`, `"|"`, `"*.txt"`} {
		if !strings.Contains(got, literal) {
			t.Fatalf("recovered preview = %q, missing literal %q", got, literal)
		}
	}
}

func TestEditor_CommandHelpDocumentsPortableDirectExecution(t *testing.T) {
	for _, want := range []string{
		"Command line", "portable", "same on Windows, macOS, and Linux", "empty argument",
		"literal newline", `C:\Program Files\Tool`, "--tag one --tag two", "--empty ''",
		"héllo 世界", "press Enter", "does not run a shell", "cmd /c", "sh -c", "expansion and security",
	} {
		if !strings.Contains(editorHelpMarkdown, want) {
			t.Errorf("editor help missing %q", want)
		}
	}
}

func TestEditor_EmptyScheduleShowsGuidance(t *testing.T) {
	e, _ := newTestEditor(t, nil)
	if got := e.schedPreview.Text; !strings.Contains(strings.ToLower(got), "type a schedule") {
		t.Fatalf("empty schedule preview = %q, want guidance", got)
	}
}

func TestEditor_PreviewSelectsHumanSyntax(t *testing.T) {
	e, fb := newTestEditor(t, nil)
	e.schedule.SetText("every day at 09:00")
	if _, req := fb.lastPreviewCall(); req.ScheduleSyntax != "human" {
		t.Fatalf("preview ScheduleSyntax = %q, want human", req.ScheduleSyntax)
	}
}

func TestEditor_CronPreviewAndCreateRetainSource(t *testing.T) {
	e, fb := newTestEditor(t, nil)
	e.name.SetText("weekday-cron")
	e.commandLine.SetText("cmd")
	e.schedule.SetText("  0 9 * * 1-5  ")

	if e.save.Disabled() {
		t.Fatal("Save should be enabled for supported cron")
	}
	if _, req := fb.lastPreviewCall(); req.Schedule != "0 9 * * 1-5" || req.ScheduleSyntax != "cron" {
		t.Fatalf("preview request = %+v, want retained normalized cron input", req)
	}

	e.submit()
	waitFor(t, func() bool { n, _ := fb.lastCreateCall(); return n == 1 })
	if _, req := fb.lastCreateCall(); req.Schedule != "0 9 * * 1-5" || req.ScheduleSyntax != "cron" {
		t.Fatalf("create request = %+v, want retained normalized cron input", req)
	}
}

func TestEditor_MonthlyCalendarCronAndPolicyPreview(t *testing.T) {
	e, fb := newTestEditor(t, nil)
	e.name.SetText("calendar-cron")
	e.commandLine.SetText("cmd")
	e.schedule.SetText("0 9 31W * *")
	if e.save.Disabled() {
		t.Fatal("Save should be enabled for nearest-weekday cron")
	}
	before, req := fb.lastPreviewCall()
	if req.Schedule != "0 9 31W * *" || req.ScheduleSyntax != "cron" || req.MissingDatePolicy != "skip" {
		t.Fatalf("initial preview = %+v", req)
	}
	e.missingDate.SetSelected(missingDateLabel(domain.MissingDateNextValid))
	after, req := fb.lastPreviewCall()
	if after <= before || req.MissingDatePolicy != "next_valid" {
		t.Fatalf("policy preview count %d -> %d, req=%+v", before, after, req)
	}
	e.submit()
	waitFor(t, func() bool { n, _ := fb.lastCreateCall(); return n == 1 })
	if _, created := fb.lastCreateCall(); created.Schedule != "0 9 31W * *" || created.MissingDatePolicy != "next_valid" {
		t.Fatalf("create request = %+v", created)
	}
}

func TestEditor_FiveWordHumanInputRemainsHuman(t *testing.T) {
	e, fb := newTestEditor(t, nil)
	e.schedule.SetText("3rd wednesday monthly at 14:00")

	if _, req := fb.lastPreviewCall(); req.ScheduleSyntax != "human" {
		t.Fatalf("preview ScheduleSyntax = %q, want human", req.ScheduleSyntax)
	}
}

func TestEditor_InvalidOrRefusedCronDisablesSaveWithoutPreview(t *testing.T) {
	for _, tc := range []struct {
		name, input, reason string
	}{
		{name: "invalid field", input: "61 9 * * *", reason: "minute"},
		{name: "fidelity refusal", input: "0 9 1 * 1", reason: "either"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e, fb := newTestEditor(t, nil)
			e.name.SetText("invalid-cron")
			e.commandLine.SetText("cmd")
			before, _ := fb.lastPreviewCall()
			e.schedule.SetText(tc.input)

			if !e.save.Disabled() {
				t.Fatal("Save should stay disabled for invalid or refused cron")
			}
			if !strings.Contains(strings.ToLower(e.schedPreview.Text), tc.reason) {
				t.Fatalf("preview error = %q, want named reason %q", e.schedPreview.Text, tc.reason)
			}
			if after, _ := fb.lastPreviewCall(); after != before {
				t.Fatalf("backend previews changed from %d to %d for refused local input", before, after)
			}
		})
	}
}

func TestEditor_OneOffSubmissionOmitsRecurringSyntax(t *testing.T) {
	e, fb := newTestEditor(t, nil)
	e.name.SetText("once")
	e.commandLine.SetText("cmd")
	e.schedule.SetText("0 9 * * 1-5")
	e.mode.SetSelected(modeOneOff)
	e.oneOffDate.SetText("2099-01-01")
	e.oneOffTime.SetText("09:00")

	e.submit()
	waitFor(t, func() bool { n, _ := fb.lastCreateCall(); return n == 1 })
	if _, req := fb.lastCreateCall(); req.Schedule != "" || req.ScheduleSyntax != "" || req.At == nil {
		t.Fatalf("one-off create = %+v, want At only", req)
	}
}

func TestEditor_HelpDocumentsDualSyntax(t *testing.T) {
	for _, want := range []string{
		"plain-language phrase",
		"five- or six-field cron",
		"0 9 * * 1-5",
		"docs/cron.md#fidelity",
	} {
		if !strings.Contains(editorHelpMarkdown, want) {
			t.Errorf("editor help missing %q", want)
		}
	}
}

// --- 003: Help toggle ----------------------------------------------------

func TestEditor_HelpToggle(t *testing.T) {
	for _, mode := range []string{modeRecurring, modeOneOff, modeManual} {
		t.Run(mode, func(t *testing.T) {
			e, _ := newTestEditor(t, nil)
			e.mode.SetSelected(mode)
			e.schedule.SetText("every 15 minutes")
			if !e.previewContent.Visible() || e.helpContent.Visible() {
				t.Fatal("right pane should start on Preview")
			}
			e.toggleHelp()
			if e.previewContent.Visible() || !e.helpContent.Visible() {
				t.Fatal("toggle should show Help, hide Preview")
			}
			e.toggleHelp()
			if !e.previewContent.Visible() || e.helpContent.Visible() {
				t.Fatal("toggle back should restore Preview")
			}
			if e.schedule.Text != "every 15 minutes" {
				t.Fatalf("input lost across Help toggle: %q", e.schedule.Text)
			}
		})
	}
}

// --- 003: dirty detection drives Cancel confirm --------------------------

func TestEditor_DirtyDetection(t *testing.T) {
	e, _ := newTestEditor(t, nil)
	if e.isDirty() {
		t.Fatal("fresh New form should not be dirty")
	}
	e.name.SetText("something")
	if !e.isDirty() {
		t.Fatal("editing a field should mark the form dirty")
	}
}

func TestEditor_DirtyDetection_EditBaseline(t *testing.T) {
	task := &domain.Task{ID: "t1", Name: "nightly", Command: "cmd", Timezone: "UTC"}
	e, _ := newTestEditor(t, task)
	if e.isDirty() {
		t.Fatal("unchanged Edit form should not be dirty (baseline = prefilled values)")
	}
	e.commandLine.SetText("python")
	if !e.isDirty() {
		t.Fatal("changing a prefilled field should mark the form dirty")
	}
}

// --- US4: interval anchor ------------------------------------------------

func TestEditor_StartAtVisibilityAndPhrase(t *testing.T) {
	e, _ := newTestEditor(t, nil)

	e.schedule.SetText("every 15 minutes")
	if !hasLabel(whenLabels(e), "Start at") {
		t.Fatalf("Start at should appear for sub-daily interval: %v", whenLabels(e))
	}

	e.schedule.SetText("every day at 09:00")
	if hasLabel(whenLabels(e), "Start at") {
		t.Fatalf("Start at should be hidden for daily schedule: %v", whenLabels(e))
	}

	e.schedule.SetText("every 15 minutes")
	e.startAt.SetText("09:00")
	if got := e.effectiveSchedule(); got != "every 15 minutes starting at 09:00" {
		t.Fatalf("effectiveSchedule = %q", got)
	}
	e.name.SetText("x")
	e.commandLine.SetText("cmd")
	if got := e.buildForm().schedule; got != "every 15 minutes starting at 09:00" {
		t.Fatalf("submitted schedule = %q", got)
	}
}

// --- US5: timezone combo + one-off assembly ------------------------------

func TestEditor_TimezoneComboAndOneOffAssembly(t *testing.T) {
	e, _ := newTestEditor(t, nil)

	e.tz.SetText("UTC")
	if _, err := timezone.Resolve("UTC"); err != nil {
		t.Fatalf("UTC should resolve: %v", err)
	}
	e.mode.SetSelected(modeOneOff)
	e.oneOffDate.SetText("2099-08-04")
	e.oneOffTime.SetText("09:00")
	got, err := e.oneOffInstant()
	if err != nil {
		t.Fatalf("oneOffInstant: %v", err)
	}
	if got.UTC() != time.Date(2099, 8, 4, 9, 0, 0, 0, time.UTC) {
		t.Fatalf("assembled instant = %v", got.UTC())
	}
}

// --- US6: advanced labels submit correct wire values ---------------------

func TestEditor_AdvancedLabelsMapToWire(t *testing.T) {
	e, _ := newTestEditor(t, nil)
	e.name.SetText("x")
	e.commandLine.SetText("cmd")
	e.schedule.SetText("every 15 minutes")
	e.overlap.SetSelected("Allow concurrent runs")
	e.catchup.SetSelected("Skip missed runs")
	e.timeBasis.SetSelected(timeBasisLabel(domain.TimeBasisElapsed))
	e.dstGap.SetSelected(dstGapLabel(domain.DSTGapSkip))
	e.dstOverlap.SetSelected(dstOverlapLabel(domain.DSTOverlapBoth))

	f := e.buildForm()
	if f.overlap != string(domain.OverlapAllowConcurrent) {
		t.Fatalf("overlap wire = %q, want %q", f.overlap, domain.OverlapAllowConcurrent)
	}
	if f.catchup != string(domain.CatchupNone) {
		t.Fatalf("catchup wire = %q, want %q", f.catchup, domain.CatchupNone)
	}
	if f.timeBasis != string(domain.TimeBasisElapsed) || f.dstGap != string(domain.DSTGapSkip) || f.dstOverlap != string(domain.DSTOverlapBoth) {
		t.Fatalf("DST policy wire = %q/%q/%q", f.timeBasis, f.dstGap, f.dstOverlap)
	}
}

func TestEditor_EditPrefillMapsPolicyLabels(t *testing.T) {
	task := &domain.Task{
		ID: "t1", Name: "nightly", Command: "cmd", Timezone: "UTC",
		OverlapPolicy: domain.OverlapSkip, CatchupPolicy: domain.CatchupNone,
	}
	e, _ := newTestEditor(t, task)
	if e.overlap.Selected != "Skip this run" {
		t.Fatalf("overlap label = %q, want 'Skip this run'", e.overlap.Selected)
	}
	if e.catchup.Selected != "Skip missed runs" {
		t.Fatalf("catchup label = %q", e.catchup.Selected)
	}
	// On edit, a blank schedule is allowed (keeps the existing one).
	if !e.valid() {
		t.Fatal("edit with name+command should be valid even with blank schedule")
	}
}
