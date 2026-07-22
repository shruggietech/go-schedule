package server

import (
	"testing"
	"time"
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
}
