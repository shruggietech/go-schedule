package server

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// TestTaskDetail_ServesStoredExpression verifies the phrase a schedule was
// created from is served back on task detail, which is what lets an editing
// client show the user their own wording again.
func TestTaskDetail_ServesStoredExpression(t *testing.T) {
	s := newTestServer(t)
	task := newTaskFor(t, s, TaskCreateRequest{
		Name: "t", Command: "/bin/true", Schedule: "weekdays at 9:00 AM", Timezone: "UTC",
	})
	if got := task.Schedule.Expression; got != "weekdays at 9:00 AM" {
		t.Errorf("Expression = %q, want the phrase as typed", got)
	}

	// It survives a round trip through storage, not just the create response.
	fetched := getTask(t, s, task.Task.ID)
	if got := fetched.Schedule.Expression; got != "weekdays at 9:00 AM" {
		t.Errorf("Expression after reload = %q, want the phrase as typed", got)
	}

	// And it is re-submittable: sending it back reproduces the same recurrence.
	if fetched.Schedule.RRULE != task.Schedule.RRULE {
		t.Errorf("recurrence changed across reload: %q -> %q", task.Schedule.RRULE, fetched.Schedule.RRULE)
	}
}

func TestTaskDetail_ServesStoredCronExpressionAndIdentity(t *testing.T) {
	s := newTestServer(t)
	task := newTaskFor(t, s, TaskCreateRequest{
		Name: "cron", Command: "/bin/true", Schedule: "0 9 * * 1-5", Timezone: "UTC",
	})
	fetched := getTask(t, s, task.Task.ID)
	if fetched.Schedule.Expression != "0 9 * * 1-5" {
		t.Fatalf("Expression after reload = %q", fetched.Schedule.Expression)
	}
	if fetched.Schedule.SourceSyntax != "cron" {
		t.Fatalf("SourceSyntax after reload = %q, want cron", fetched.Schedule.SourceSyntax)
	}
}

// TestTaskDetail_NoExpressionForOneOff verifies one-off schedules are served
// without a phrase — their date and time come from RunAt, and a fabricated
// phrase would be meaningless.
func TestTaskDetail_NoExpressionForOneOff(t *testing.T) {
	s := newTestServer(t)
	at := time.Now().UTC().Add(24 * time.Hour)
	task := newTaskFor(t, s, TaskCreateRequest{
		Name: "once", Command: "/bin/true", Timezone: "UTC", At: &at,
	})
	if got := task.Schedule.Expression; got != "" {
		t.Errorf("one-off Expression = %q, want empty", got)
	}
	if task.Schedule.RunAt == nil {
		t.Error("one-off served without RunAt")
	}
	if task.Schedule.SourceSyntax != "" {
		t.Errorf("one-off SourceSyntax = %q, want empty", task.Schedule.SourceSyntax)
	}
}

func TestTaskDetail_LegacyExpressionlessRecurrenceHasNoSourceIdentity(t *testing.T) {
	s := newTestServer(t)
	anchor := time.Now().UTC()
	sch := domain.Schedule{
		Kind: domain.ScheduleRecurring, RRULE: "FREQ=DAILY;BYHOUR=9;BYMINUTE=0;BYSECOND=0",
		Anchor: &anchor, HumanSummary: "Every day at 09:00",
	}
	if err := s.store.CreateSchedule(&sch); err != nil {
		t.Fatal(err)
	}
	task := domain.Task{
		Name: "legacy", Command: "/bin/true", Enabled: true, Timezone: "UTC", ScheduleID: sch.ID,
		OverlapPolicy: domain.OverlapQueueOne, CatchupPolicy: domain.CatchupOne,
		MissingDatePolicy: domain.MissingDateSkip, State: domain.TaskActive,
	}
	if err := s.store.CreateTask(&task); err != nil {
		t.Fatal(err)
	}
	fetched := getTask(t, s, task.ID)
	if fetched.Schedule.SourceSyntax != "" {
		t.Fatalf("legacy SourceSyntax = %q, want empty", fetched.Schedule.SourceSyntax)
	}
}
