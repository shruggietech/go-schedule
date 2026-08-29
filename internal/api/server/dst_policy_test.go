package server

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestTaskDSTPolicyDefaultsExplicitValuesAndPersistence(t *testing.T) {
	s := newTestServer(t)
	defaults := newTaskFor(t, s, TaskCreateRequest{
		Name: "defaults", Command: "/bin/true", Schedule: "every day at 09:00", Timezone: "UTC",
	})
	if got := defaults.Task.SchedulePolicy(); got.TimeBasis != domain.TimeBasisWallClock || got.DSTGap != domain.DSTGapNextValid || got.DSTOverlap != domain.DSTOverlapFirst {
		t.Fatalf("defaults = %#v", got)
	}
	if defaults.PolicySummary == "" {
		t.Fatal("task detail omitted policy summary")
	}

	explicit := newTaskFor(t, s, TaskCreateRequest{
		Name: "explicit", Command: "/bin/true", Schedule: "every 6 hours", Timezone: "America/New_York",
		TimeBasis: "elapsed", DSTGapPolicy: "skip", DSTOverlapPolicy: "both",
	})
	if got := explicit.Task.SchedulePolicy(); got.TimeBasis != domain.TimeBasisElapsed || got.DSTGap != domain.DSTGapSkip || got.DSTOverlap != domain.DSTOverlapBoth {
		t.Fatalf("explicit = %#v", got)
	}

	rec := doJSON(t, s, http.MethodPatch, "/v1/tasks/"+explicit.Task.ID, TaskUpdateRequest{Name: "renamed"})
	if rec.Code != http.StatusOK {
		t.Fatalf("update status=%d body=%s", rec.Code, rec.Body.String())
	}
	got := getTask(t, s, explicit.Task.ID).Task.SchedulePolicy()
	if got.TimeBasis != domain.TimeBasisElapsed || got.DSTGap != domain.DSTGapSkip || got.DSTOverlap != domain.DSTOverlapBoth {
		t.Fatalf("policy changed on unrelated edit: %#v", got)
	}
	var preview PreviewResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &preview); err != nil || preview.PolicySummary == "" {
		t.Fatalf("preview policy summary=%q err=%v", preview.PolicySummary, err)
	}
}

func TestTaskDSTPolicyValidationIsFieldSpecificAndNonMutating(t *testing.T) {
	cases := []struct {
		name  string
		req   TaskCreateRequest
		field string
	}{
		{"basis", TaskCreateRequest{Name: "x", Command: "true", Schedule: "every hour", TimeBasis: "solar"}, "time_basis"},
		{"gap", TaskCreateRequest{Name: "x", Command: "true", Schedule: "every hour", DSTGapPolicy: "delay"}, "dst_gap_policy"},
		{"overlap", TaskCreateRequest{Name: "x", Command: "true", Schedule: "every hour", DSTOverlapPolicy: "random"}, "dst_overlap_policy"},
		{"elapsed calendar", TaskCreateRequest{Name: "x", Command: "true", Schedule: "3rd wednesday monthly at 14:00", TimeBasis: "elapsed"}, "time_basis"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestServer(t)
			rec := doJSON(t, s, http.MethodPost, "/v1/tasks", tc.req)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
			}
			var apiErr APIError
			_ = json.Unmarshal(rec.Body.Bytes(), &apiErr)
			if apiErr.Error.Field != tc.field {
				t.Fatalf("field=%q, want %q", apiErr.Error.Field, tc.field)
			}
			tasks, err := s.store.ListTasks("", "")
			if err != nil || len(tasks) != 0 {
				t.Fatalf("failure mutated tasks: len=%d err=%v", len(tasks), err)
			}
		})
	}
}

func TestPreviewAcceptsDSTPolicy(t *testing.T) {
	s := newTestServer(t)
	rec := doJSON(t, s, http.MethodPost, "/v1/schedules/preview", PreviewRequest{
		Schedule: "every 6 hours", Timezone: "America/New_York",
		TimeBasis: "elapsed", DSTGapPolicy: "skip", DSTOverlapPolicy: "both",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestDSTPolicyPreviewUpdateAndCalendarParity(t *testing.T) {
	s := newTestServer(t)
	badPreview := doJSON(t, s, http.MethodPost, "/v1/schedules/preview", PreviewRequest{
		Schedule: "every hour", DSTGapPolicy: "later",
	})
	if badPreview.Code != http.StatusBadRequest || !jsonErrorField(badPreview.Body.Bytes(), "dst_gap_policy") {
		t.Fatalf("bad preview status=%d body=%s", badPreview.Code, badPreview.Body.String())
	}

	monthly := newTaskFor(t, s, TaskCreateRequest{
		Name: "monthly", Command: "true", Schedule: "3rd wednesday monthly at 14:00", Timezone: "UTC",
	})
	badUpdate := doJSON(t, s, http.MethodPatch, "/v1/tasks/"+monthly.Task.ID, TaskUpdateRequest{TimeBasis: "elapsed"})
	if badUpdate.Code != http.StatusBadRequest || !jsonErrorField(badUpdate.Body.Bytes(), "time_basis") {
		t.Fatalf("bad update status=%d body=%s", badUpdate.Code, badUpdate.Body.String())
	}
	if got := getTask(t, s, monthly.Task.ID).Task.TimeBasis; got != domain.TimeBasisWallClock {
		t.Fatalf("invalid update mutated basis to %q", got)
	}

	utcTask := newTaskFor(t, s, TaskCreateRequest{
		Name: "utc", Command: "true", Schedule: "every day at 09:00",
		Timezone: "America/New_York", TimeBasis: "utc",
	})
	from := time.Now().UTC().Add(-time.Minute)
	to := from.Add(72 * time.Hour)
	path := "/v1/calendar?from=" + from.Format(time.RFC3339) + "&to=" + to.Format(time.RFC3339)
	rec := doJSON(t, s, http.MethodGet, path, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("calendar status=%d body=%s", rec.Code, rec.Body.String())
	}
	var calendar CalendarResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &calendar); err != nil {
		t.Fatal(err)
	}
	var first time.Time
	for _, occurrence := range calendar.Occurrences {
		if occurrence.TaskID == utcTask.Task.ID && occurrence.Kind == "scheduled" {
			first = occurrence.Time
			break
		}
	}
	if first.IsZero() || len(utcTask.NextRuns) == 0 || !first.Equal(utcTask.NextRuns[0]) {
		t.Fatalf("calendar first=%s detail=%v", first, utcTask.NextRuns)
	}
}

func jsonErrorField(body []byte, want string) bool {
	var apiErr APIError
	return json.Unmarshal(body, &apiErr) == nil && apiErr.Error.Field == want
}
