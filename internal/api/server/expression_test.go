package server

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// TestTaskDetail_BackfillsExpressionForLegacySchedules pins FR-003: schedules
// written before store migration v4 carry no phrase, and an editing client
// showing a blank field is exactly the reported defect. The server derives an
// equivalent phrase on read. Nothing is written back to the store.
func TestTaskDetail_BackfillsExpressionForLegacySchedules(t *testing.T) {
	s := newTestServer(t)

	// A schedule as a pre-v4 database holds it: RRULE and summary, no phrase.
	legacy := &domain.Schedule{
		Kind:         domain.ScheduleRecurring,
		RRULE:        "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR;BYHOUR=9;BYMINUTE=0;BYSECOND=0",
		HumanSummary: "Every weekday at 09:00",
	}
	if err := s.store.CreateSchedule(legacy); err != nil {
		t.Fatal(err)
	}
	task := &domain.Task{
		Name: "legacy", Command: "/bin/true", Timezone: "UTC",
		ScheduleID: legacy.ID, State: domain.TaskActive, Enabled: true,
		OverlapPolicy: domain.OverlapQueueOne, CatchupPolicy: domain.CatchupOne,
	}
	if err := s.store.CreateTask(task); err != nil {
		t.Fatal(err)
	}

	rec := doJSON(t, s, http.MethodGet, "/v1/tasks/"+task.ID, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	var resp TaskResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Schedule.Expression == "" {
		t.Fatal("legacy schedule served with no expression: an editing client has nothing to show")
	}

	// The derived phrase must describe the same recurrence — resubmitting it
	// must not move the task.
	rec = doJSON(t, s, http.MethodPatch, "/v1/tasks/"+task.ID, TaskUpdateRequest{
		Schedule: resp.Schedule.Expression,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("resubmitting the derived phrase %q failed: %d %s",
			resp.Schedule.Expression, rec.Code, rec.Body.String())
	}
	var after TaskResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &after); err != nil {
		t.Fatal(err)
	}
	if after.Schedule.RRULE != legacy.RRULE {
		t.Errorf("derived phrase %q changed the recurrence:\n before: %s\n after:  %s",
			resp.Schedule.Expression, legacy.RRULE, after.Schedule.RRULE)
	}

	// Nothing was written back to storage by the read.
	stored, err := s.store.GetSchedule(legacy.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Expression != "" {
		t.Errorf("read path wrote %q back to the store; derivation must be read-only", stored.Expression)
	}
}

// TestTaskDetail_KeepsStoredExpression verifies a stored phrase is served as-is:
// the user's own wording wins over the canonical derived form.
func TestTaskDetail_KeepsStoredExpression(t *testing.T) {
	s := newTestServer(t)
	task := newTaskFor(t, s, TaskCreateRequest{
		Name: "t", Command: "/bin/true", Schedule: "weekdays at 9:00 AM", Timezone: "UTC",
	})
	if got := task.Schedule.Expression; got != "weekdays at 9:00 AM" {
		t.Errorf("Expression = %q, want the phrase as typed", got)
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
