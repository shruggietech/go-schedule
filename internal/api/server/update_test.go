package server

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// newTaskFor creates a task and returns its detail, for tests that need an
// existing task to update.
func newTaskFor(t *testing.T, s *Server, req TaskCreateRequest) TaskResponse {
	t.Helper()
	rec := doJSON(t, s, http.MethodPost, "/v1/tasks", req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create task: status %d, body=%s", rec.Code, rec.Body.String())
	}
	var resp TaskResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	return resp
}

// getTask fetches a task's detail through the API.
func getTask(t *testing.T, s *Server, id string) TaskResponse {
	t.Helper()
	rec := doJSON(t, s, http.MethodGet, "/v1/tasks/"+id, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("get task: status %d, body=%s", rec.Code, rec.Body.String())
	}
	var resp TaskResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	return resp
}

func newGroupFor(t *testing.T, s *Server, name string) domain.Group {
	t.Helper()
	rec := doJSON(t, s, http.MethodPost, "/v1/groups", GroupCreateRequest{Name: name})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create group: status %d, body=%s", rec.Code, rec.Body.String())
	}
	var g domain.Group
	if err := json.Unmarshal(rec.Body.Bytes(), &g); err != nil {
		t.Fatal(err)
	}
	return g
}

func ptr(s string) *string { return &s }

// TestUpdateTask_GroupTriState pins FR-014: three distinct intents must be
// expressible. Before this feature an empty group value meant "unchanged", so
// no client could take a task back out of a group.
func TestUpdateTask_GroupTriState(t *testing.T) {
	s := newTestServer(t)
	g1 := newGroupFor(t, s, "alpha")
	g2 := newGroupFor(t, s, "beta")
	task := newTaskFor(t, s, TaskCreateRequest{
		Name: "t", Command: "/bin/true", Schedule: "every day at 09:00", Timezone: "UTC", GroupID: g1.ID,
	})
	path := "/v1/tasks/" + task.Task.ID

	groupNow := func() string {
		t.Helper()
		rec := doJSON(t, s, http.MethodGet, path, nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("get task: status %d", rec.Code)
		}
		var resp TaskResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}
		return resp.Task.GroupID
	}

	// Omitted: membership unchanged.
	if rec := doJSON(t, s, http.MethodPatch, path, TaskUpdateRequest{Name: "renamed"}); rec.Code != http.StatusOK {
		t.Fatalf("omitted group: status %d, body=%s", rec.Code, rec.Body.String())
	}
	if got := groupNow(); got != g1.ID {
		t.Errorf("omitted group_id changed membership: got %q, want %q", got, g1.ID)
	}

	// Named: reassign.
	if rec := doJSON(t, s, http.MethodPatch, path, TaskUpdateRequest{GroupID: ptr(g2.ID)}); rec.Code != http.StatusOK {
		t.Fatalf("assign group: status %d, body=%s", rec.Code, rec.Body.String())
	}
	if got := groupNow(); got != g2.ID {
		t.Errorf("assign: got %q, want %q", got, g2.ID)
	}

	// Empty: clear.
	if rec := doJSON(t, s, http.MethodPatch, path, TaskUpdateRequest{GroupID: ptr("")}); rec.Code != http.StatusOK {
		t.Fatalf("clear group: status %d, body=%s", rec.Code, rec.Body.String())
	}
	if got := groupNow(); got != "" {
		t.Errorf("clear: got %q, want empty (task should be ungrouped)", got)
	}
}

// TestUpdateTask_UnknownGroupIsValidationError pins FR-016: a bad group id is
// the caller's mistake and must name the field, not surface as a 500.
func TestUpdateTask_UnknownGroupIsValidationError(t *testing.T) {
	s := newTestServer(t)
	task := newTaskFor(t, s, TaskCreateRequest{
		Name: "t", Command: "/bin/true", Schedule: "every day at 09:00", Timezone: "UTC",
	})

	rec := doJSON(t, s, http.MethodPatch, "/v1/tasks/"+task.Task.ID, TaskUpdateRequest{GroupID: ptr("no-such-group")})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	var env APIError
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	if env.Error.Code != CodeValidation || env.Error.Field != "group_id" {
		t.Errorf("error = %+v, want validation error naming group_id", env.Error)
	}
}

// TestUpdateTask_UntouchedSaveLeavesScheduleIdentical is SC-002 — the central
// promise of the schedule-fidelity fix. Opening a task and saving without
// changing anything must not move when it runs.
func TestUpdateTask_UntouchedSaveLeavesScheduleIdentical(t *testing.T) {
	s := newTestServer(t)
	task := newTaskFor(t, s, TaskCreateRequest{
		Name: "nightly", Command: "/bin/true", Schedule: "weekdays at 09:00", Timezone: "UTC",
	})
	beforeSchedID := task.Task.ScheduleID
	beforeRuns := task.NextRuns
	beforeRule := task.Schedule.RRULE

	// The shape the editor submits when nothing about the timing was touched:
	// every other field present, schedule blank and no one-off instant.
	rec := doJSON(t, s, http.MethodPatch, "/v1/tasks/"+task.Task.ID, TaskUpdateRequest{
		Name: task.Task.Name, Command: task.Task.Command, Timezone: task.Task.Timezone,
		OverlapPolicy: string(task.Task.OverlapPolicy), CatchupPolicy: string(task.Task.CatchupPolicy),
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	var after TaskResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &after); err != nil {
		t.Fatal(err)
	}

	if after.Task.ScheduleID != beforeSchedID {
		t.Errorf("schedule was replaced: %q -> %q", beforeSchedID, after.Task.ScheduleID)
	}
	if after.Schedule.RRULE != beforeRule {
		t.Errorf("recurrence changed: %q -> %q", beforeRule, after.Schedule.RRULE)
	}
	if len(after.NextRuns) != len(beforeRuns) {
		t.Fatalf("next-run count changed: %d -> %d", len(beforeRuns), len(after.NextRuns))
	}
	for i := range beforeRuns {
		if !after.NextRuns[i].Equal(beforeRuns[i]) {
			t.Errorf("next run %d moved: %s -> %s", i, beforeRuns[i], after.NextRuns[i])
		}
	}
}

// TestUpdateTask_TimezoneChangeReanchors pins FR-011. Once the editor prefills
// the phrase, a normal save resubmits it and the server re-parses it in the new
// zone; this test stops that from silently regressing.
func TestUpdateTask_TimezoneChangeReanchors(t *testing.T) {
	s := newTestServer(t)
	task := newTaskFor(t, s, TaskCreateRequest{
		Name: "t", Command: "/bin/true", Schedule: "every day at 09:00", Timezone: "UTC",
	})
	before := task.NextRuns[0]

	rec := doJSON(t, s, http.MethodPatch, "/v1/tasks/"+task.Task.ID, TaskUpdateRequest{
		Timezone: "America/New_York", Schedule: "every day at 09:00",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	var after TaskResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &after); err != nil {
		t.Fatal(err)
	}
	if len(after.NextRuns) == 0 {
		t.Fatal("no next runs after update")
	}
	if after.NextRuns[0].Equal(before) {
		t.Errorf("09:00 UTC and 09:00 America/New_York resolved to the same instant %s; "+
			"the recurrence was not re-interpreted in the new zone", before)
	}
}
