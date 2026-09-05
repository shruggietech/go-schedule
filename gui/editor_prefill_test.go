package gui

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

// These pin issue #4: opening a task for editing used to reset Mode to
// Recurring and blank the timing fields regardless of how the task was actually
// scheduled, so the dialog asserted something false about every one-off task.

func recurringDetail(expression string) *server.TaskResponse {
	return &server.TaskResponse{
		Task: domain.Task{
			ID: "t1", Name: "nightly", Command: "/bin/true", Timezone: "UTC",
			OverlapPolicy: domain.OverlapQueueOne, CatchupPolicy: domain.CatchupOne,
		},
		Schedule: domain.Schedule{
			Kind:         domain.ScheduleRecurring,
			RRULE:        "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR;BYHOUR=9;BYMINUTE=0;BYSECOND=0",
			HumanSummary: "Every weekday at 09:00",
			Expression:   expression,
		},
	}
}

func oneOffDetail(t *testing.T, tz string, at time.Time) *server.TaskResponse {
	t.Helper()
	utc := at.UTC()
	return &server.TaskResponse{
		Task: domain.Task{
			ID: "t2", Name: "once", Command: "/bin/true", Timezone: tz,
			OverlapPolicy: domain.OverlapQueueOne, CatchupPolicy: domain.CatchupOne,
		},
		Schedule: domain.Schedule{
			Kind:         domain.ScheduleOneOff,
			RunAt:        &utc,
			HumanSummary: "Once at " + utc.Format("2006-01-02 15:04 MST"),
		},
	}
}

// TestEditor_PrefillsOneOffModeAndInstant covers FR-006/FR-007.
func TestEditor_PrefillsOneOffModeAndInstant(t *testing.T) {
	// 14:30 UTC is 09:30 in New York on this date, a zone that actually
	// differs, so a naive UTC render would be visibly wrong.
	at := time.Date(2026, 8, 4, 14, 30, 0, 0, time.UTC)
	e, _ := newTestEditorDetail(t, oneOffDetail(t, "America/New_York", at))

	if e.mode.Selected != modeOneOff {
		t.Fatalf("Mode = %q, want %q", e.mode.Selected, modeOneOff)
	}
	if e.oneOffDate.Text != "2026-08-04" {
		t.Errorf("Date = %q, want 2026-08-04", e.oneOffDate.Text)
	}
	if e.oneOffTime.Text != "10:30" {
		t.Errorf("Time = %q, want 10:30 (14:30 UTC expressed in America/New_York)", e.oneOffTime.Text)
	}
}

// TestEditor_PrefillsRecurringSchedule covers FR-006.
func TestEditor_PrefillsRecurringSchedule(t *testing.T) {
	e, _ := newTestEditorDetail(t, recurringDetail("weekdays at 09:00"))

	if e.mode.Selected != modeRecurring {
		t.Fatalf("Mode = %q, want %q", e.mode.Selected, modeRecurring)
	}
	if e.schedule.Text != "weekdays at 09:00" {
		t.Errorf("Schedule = %q, want the stored phrase", e.schedule.Text)
	}
}

func TestEditor_ExistingCommandRoundTripsWithoutLoss(t *testing.T) {
	detail := recurringDetail("weekdays at 09:00")
	detail.Task.Command = `C:\Program Files\Tool's "runner".exe`
	detail.Task.Args = []string{"", "--tag", "one", "--tag", "two", "héllo 世界", "tabs\there", "lines\r\nhere", `C:\trail\`, "$HOME", "|", "*.txt"}

	e, fb := newTestEditorDetail(t, detail)
	got, err := e.invocation()
	if err != nil {
		t.Fatalf("canonical prefill does not parse: %v", err)
	}
	if got.Program != detail.Task.Command || !reflect.DeepEqual(got.Args, detail.Task.Args) {
		t.Fatalf("prefilled text %q produced %#v, want command %q args %#v", e.commandLine.Text, got, detail.Task.Command, detail.Task.Args)
	}
	if e.isDirty() {
		t.Fatal("canonical command prefill must be the unchanged edit baseline")
	}

	e.submit()
	waitFor(t, func() bool { n, _, _ := fb.lastUpdateCall(); return n == 1 })
	_, _, req := fb.lastUpdateCall()
	if req.Command != detail.Task.Command || !reflect.DeepEqual(req.Args, detail.Task.Args) {
		t.Fatalf("unchanged update = command %q args %#v", req.Command, req.Args)
	}
}

func TestEditor_ExistingCommandCanClearAllArguments(t *testing.T) {
	detail := recurringDetail("weekdays at 09:00")
	detail.Task.Command = "program"
	detail.Task.Args = []string{"stale", "target"}

	e, fb := newTestEditorDetail(t, detail)
	e.commandLine.SetText("program")
	e.submit()
	waitFor(t, func() bool { n, _, _ := fb.lastUpdateCall(); return n == 1 })
	_, _, req := fb.lastUpdateCall()
	if req.Args == nil || len(req.Args) != 0 {
		t.Fatalf("cleared update args = %#v, want explicit non-nil empty slice", req.Args)
	}
}

func TestEditor_PrefillsStartupSchedule(t *testing.T) {
	detail := recurringDetail("@reboot")
	detail.Schedule.Kind = domain.ScheduleEvent
	detail.Schedule.TriggerID = domain.StartupEventID
	detail.Schedule.RRULE = ""
	e, fb := newTestEditorDetail(t, detail)
	if e.mode.Selected != modeRecurring || e.schedule.Text != "@reboot" || e.scheduleUnreadable {
		t.Fatalf("startup prefill: mode=%q schedule=%q unreadable=%v", e.mode.Selected, e.schedule.Text, e.scheduleUnreadable)
	}
	if _, req := fb.lastPreviewCall(); req.Schedule != "@reboot" || req.ScheduleSyntax != "cron" {
		t.Fatalf("startup preview request = %+v", req)
	}
}

func TestEditor_CronPrefillAndSyntaxSwitchFollowCurrentText(t *testing.T) {
	detail := recurringDetail("0 9 * * 1-5")
	detail.Schedule.SourceSyntax = "human" // Deliberately stale metadata.
	e, fb := newTestEditorDetail(t, detail)

	if got := e.schedule.Text; got != "0 9 * * 1-5" {
		t.Fatalf("Schedule = %q, want exact retained cron expression", got)
	}
	if e.isDirty() {
		t.Fatal("untouched cron edit should not be dirty")
	}
	if _, req := fb.lastPreviewCall(); req.ScheduleSyntax != "cron" {
		t.Fatalf("initial preview syntax = %q, want cron from current text", req.ScheduleSyntax)
	}

	e.submit()
	waitFor(t, func() bool { n, _, _ := fb.lastUpdateCall(); return n == 1 })
	if _, _, req := fb.lastUpdateCall(); req.Schedule != "0 9 * * 1-5" || req.ScheduleSyntax != "cron" {
		t.Fatalf("untouched update = %+v, want retained cron", req)
	}

	e.schedule.SetText("weekdays at 09:00")
	if _, req := fb.lastPreviewCall(); req.ScheduleSyntax != "human" {
		t.Fatalf("switched preview syntax = %q, want human", req.ScheduleSyntax)
	}
	e.submit()
	waitFor(t, func() bool { n, _, _ := fb.lastUpdateCall(); return n == 2 })
	if _, _, req := fb.lastUpdateCall(); req.ScheduleSyntax != "human" {
		t.Fatalf("switched update syntax = %q, want human", req.ScheduleSyntax)
	}
}

func TestEditor_HumanPrefillCanSwitchToCron(t *testing.T) {
	detail := recurringDetail("weekdays at 09:00")
	detail.Schedule.SourceSyntax = "human"
	e, fb := newTestEditorDetail(t, detail)

	e.schedule.SetText("0 9 * * 1-5")
	if _, req := fb.lastPreviewCall(); req.ScheduleSyntax != "cron" {
		t.Fatalf("switched preview syntax = %q, want cron", req.ScheduleSyntax)
	}
	e.submit()
	waitFor(t, func() bool { n, _, _ := fb.lastUpdateCall(); return n == 1 })
	if _, _, req := fb.lastUpdateCall(); req.Schedule != "0 9 * * 1-5" || req.ScheduleSyntax != "cron" {
		t.Fatalf("switched update = %+v, want cron", req)
	}
}

func TestEditor_ExpressionlessLegacyUpdateOmitsSyntax(t *testing.T) {
	detail := recurringDetail("")
	e, fb := newTestEditorDetail(t, detail)

	if !e.valid() {
		t.Fatal("blank expressionless legacy edit should preserve the current schedule")
	}
	e.submit()
	waitFor(t, func() bool { n, _, _ := fb.lastUpdateCall(); return n == 1 })
	if _, _, req := fb.lastUpdateCall(); req.Schedule != "" || req.ScheduleSyntax != "" {
		t.Fatalf("legacy update = %+v, want recurring fields omitted", req)
	}
}

// TestEditor_PrefillsAnchorIntoStartAt covers FR-006: an anchored phrase is
// split so effectiveSchedule() rebuilds exactly the same phrase, rather than
// appending a second anchor clause.
func TestEditor_PrefillsAnchorIntoStartAt(t *testing.T) {
	d := recurringDetail("every 15 minutes starting at 09:00")
	d.Schedule.RRULE = "FREQ=MINUTELY;INTERVAL=15"
	e, _ := newTestEditorDetail(t, d)

	if e.schedule.Text != "every 15 minutes" {
		t.Errorf("Schedule = %q, want the phrase without its anchor clause", e.schedule.Text)
	}
	if e.startAt.Text != "09:00" {
		t.Errorf("Start at = %q, want 09:00", e.startAt.Text)
	}
	if got := e.effectiveSchedule(); got != "every 15 minutes starting at 09:00" {
		t.Errorf("effectiveSchedule() = %q, want the original phrase reconstructed exactly", got)
	}
}

// TestEditor_UntouchedEditIsNotDirty covers FR-008. Prefilled values must form
// the dirty-detection baseline, or Cancel prompts about changes nobody made.
func TestEditor_UntouchedEditIsNotDirty(t *testing.T) {
	for _, tc := range []struct {
		name   string
		detail *server.TaskResponse
	}{
		{"recurring", recurringDetail("weekdays at 09:00")},
		{"one-off", oneOffDetail(t, "UTC", time.Date(2026, 8, 4, 9, 0, 0, 0, time.UTC))},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e, _ := newTestEditorDetail(t, tc.detail)
			if e.isDirty() {
				t.Error("a freshly opened, untouched edit reports unsaved changes")
			}
		})
	}
}

// TestEditor_ModeSwitchRequiresNewTiming covers FR-011b. "Blank keeps the
// existing schedule" is only meaningful while the mode is unchanged: after a
// switch there is no existing schedule of the new kind to keep.
func TestEditor_ModeSwitchRequiresNewTiming(t *testing.T) {
	e, _ := newTestEditorDetail(t, recurringDetail("weekdays at 09:00"))

	// Unchanged mode with a blank schedule still means "keep existing".
	e.schedule.SetText("")
	if !e.valid() {
		t.Error("blank schedule with unchanged mode should stay valid (keeps existing)")
	}

	// Switching to One-off with no date/time must block saving.
	e.mode.SetSelected(modeOneOff)
	if e.valid() {
		t.Error("switching to One-off with empty date/time must not be saveable")
	}

	// Supplying the new mode's timing makes it valid again.
	future := time.Now().Add(48 * time.Hour)
	e.oneOffDate.SetText(future.Format("2006-01-02"))
	e.oneOffTime.SetText("09:00")
	if !e.valid() {
		t.Error("One-off with a future date and time should be saveable")
	}

	// And switching back to Recurring with a blank schedule is valid again,
	// because the stored schedule is recurring.
	e.mode.SetSelected(modeRecurring)
	if !e.valid() {
		t.Error("returning to the stored mode should restore the keep-existing allowance")
	}
}

// TestTasksTab_EditFetchesDetail covers FR-009: the editor is populated from a
// detail fetch, and a failed fetch does not block editing.
func TestTasksTab_EditFetchesDetail(t *testing.T) {
	fb := &fakeBackend{
		tasks:   []domain.Task{{ID: "t1", Name: "nightly", Command: "/bin/true", Timezone: "UTC", Enabled: true, State: domain.TaskActive}},
		details: map[string]server.TaskResponse{"t1": *recurringDetail("weekdays at 09:00")},
	}
	ui := NewUI(testApp, fb)
	task := fb.tasks[0]

	detail := ui.taskDetailFor(task)
	if detail == nil || detail.Schedule.Expression != "weekdays at 09:00" {
		t.Fatalf("expected fetched detail carrying the phrase, got %+v", detail)
	}

	// A failing lookup must still yield an openable editor, not a blocked one.
	fb.getTaskErr = errors.New("daemon unreachable")
	degraded := ui.taskDetailFor(task)
	if degraded == nil {
		t.Fatal("a failed detail fetch must still produce a task to edit, not nil")
	}
	if degraded.Task.ID != task.ID {
		t.Errorf("degraded detail lost the task: %+v", degraded.Task)
	}
	if degraded.Schedule.Kind != "" {
		t.Errorf("degraded detail should carry no schedule, got %+v", degraded.Schedule)
	}

	// The editor opens, and says the schedule could not be read rather than
	// implying the task has none (FR-009).
	e, _ := newTestEditorDetail(t, degraded)
	if !e.valid() {
		t.Error("a degraded edit must still be saveable (blank schedule keeps the existing one)")
	}
	if got := e.schedPreview.Text; !strings.Contains(got, "Could not read") {
		t.Errorf("preview = %q, want it to report that the schedule could not be read", got)
	}
}

// TestEditor_PrefillsAndSubmitsMissingDatePolicy is the GUI half of FR-022: the
// setting must be visible where the operator already looks for execution
// policies, prefilled from the task, and carried back on save. A selector that
// displays correctly but submits the default would silently revert the
// operator's choice on every unrelated edit.
func TestEditor_PrefillsAndSubmitsMissingDatePolicy(t *testing.T) {
	detail := recurringDetail("on the 31st of every month at 09:00")
	detail.Task.MissingDatePolicy = domain.MissingDateLastValid

	e, fb := newTestEditorDetail(t, detail)

	if got, want := e.missingDate.Selected, missingDateLabel(domain.MissingDateLastValid); got != want {
		t.Fatalf("prefilled selection = %q, want %q", got, want)
	}

	// Saving without touching the field keeps the task's policy. App.run
	// dispatches the call on a goroutine, so the assertion waits for it rather
	// than racing it.
	e.submit()
	waitFor(t, func() bool { n, _, _ := fb.lastUpdateCall(); return n == 1 })
	if _, _, req := fb.lastUpdateCall(); req.MissingDatePolicy != string(domain.MissingDateLastValid) || req.ScheduleSyntax != "human" {
		t.Errorf("submitted update = %+v, want last_valid with human syntax", req)
	}

	// Changing it submits the new value.
	e.missingDate.SetSelected(missingDateLabel(domain.MissingDateNextValid))
	e.submit()
	waitFor(t, func() bool { n, _, _ := fb.lastUpdateCall(); return n == 2 })
	if _, _, req := fb.lastUpdateCall(); req.MissingDatePolicy != string(domain.MissingDateNextValid) {
		t.Errorf("submitted policy after change = %q, want next_valid", req.MissingDatePolicy)
	}
}

// TestEditor_NewTaskDefaultsMissingDatePolicy pins the create path's default.
func TestEditor_NewTaskDefaultsMissingDatePolicy(t *testing.T) {
	e, fb := newTestEditor(t, nil)
	e.name.SetText("new")
	e.commandLine.SetText("/bin/true")
	e.schedule.SetText("every day at 09:00")
	e.submit()
	waitFor(t, func() bool { n, _ := fb.lastCreateCall(); return n == 1 })
	if _, req := fb.lastCreateCall(); req.MissingDatePolicy != string(domain.MissingDateSkip) || req.ScheduleSyntax != "human" {
		t.Errorf("new task request = %+v, want skip with human syntax", req)
	}
}

func TestEditor_PrefillsAndSubmitsDSTPolicy(t *testing.T) {
	detail := recurringDetail("every 6 hours")
	detail.Task.TimeBasis = domain.TimeBasisElapsed
	detail.Task.DSTGapPolicy = domain.DSTGapSkip
	detail.Task.DSTOverlapPolicy = domain.DSTOverlapBoth
	e, fb := newTestEditorDetail(t, detail)

	if e.timeBasis.Selected != timeBasisLabel(domain.TimeBasisElapsed) ||
		e.dstGap.Selected != dstGapLabel(domain.DSTGapSkip) ||
		e.dstOverlap.Selected != dstOverlapLabel(domain.DSTOverlapBoth) {
		t.Fatalf("prefill = %q/%q/%q", e.timeBasis.Selected, e.dstGap.Selected, e.dstOverlap.Selected)
	}
	e.submit()
	waitFor(t, func() bool { n, _, _ := fb.lastUpdateCall(); return n == 1 })
	_, _, req := fb.lastUpdateCall()
	if req.TimeBasis != "elapsed" || req.DSTGapPolicy != "skip" || req.DSTOverlapPolicy != "both" {
		t.Fatalf("submitted = %+v", req)
	}
}
