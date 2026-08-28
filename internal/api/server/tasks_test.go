package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func doJSON(t *testing.T, s *Server, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatal(err)
		}
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, &buf)
	s.Handler().ServeHTTP(rec, req)
	return rec
}

func TestCreateTask_Recurring(t *testing.T) {
	s := newTestServer(t)
	rec := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "nightly", Command: "/bin/true", Schedule: "every day at 09:00", Timezone: "UTC",
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	var resp TaskResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Task.ID == "" || resp.Schedule.HumanSummary != "Every day at 09:00" {
		t.Fatalf("unexpected task detail: %+v", resp)
	}
	if len(resp.NextRuns) == 0 {
		t.Fatal("expected computed next runs")
	}
}

func TestCreateTask_OneOffPastRejected(t *testing.T) {
	s := newTestServer(t)
	past := time.Now().UTC().Add(-time.Hour)
	rec := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "bday", Command: "/bin/true", At: &past,
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	var e APIError
	_ = json.Unmarshal(rec.Body.Bytes(), &e)
	if e.Error.Field != "at" {
		t.Fatalf("expected field 'at', got %+v", e)
	}
}

func TestCreateTask_OneOffFuture(t *testing.T) {
	s := newTestServer(t)
	future := time.Now().UTC().Add(24 * time.Hour)
	rec := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "bday", Command: "/bin/true", At: &future,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestCreateTask_InvalidScheduleAndTimezone(t *testing.T) {
	s := newTestServer(t)
	if rec := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "x", Command: "/bin/true", Schedule: "every blorp",
	}); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad schedule: status %d", rec.Code)
	}
	if rec := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "x", Command: "/bin/true", Schedule: "every day", Timezone: "Mars/Phobos",
	}); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad timezone: status %d", rec.Code)
	}
}

func TestPreview(t *testing.T) {
	s := newTestServer(t)
	rec := doJSON(t, s, http.MethodPost, "/v1/schedules/preview", PreviewRequest{
		Schedule: "3rd wednesday monthly at 14:00", Timezone: "UTC",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	var resp PreviewResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.HumanSummary != "The 3rd Wednesday of every month at 14:00" || len(resp.NextRuns) == 0 {
		t.Fatalf("unexpected preview: %+v", resp)
	}
}

func TestPreview_DualSyntaxParity(t *testing.T) {
	s := newTestServer(t)
	preview := func(schedule, syntax string) PreviewResponse {
		t.Helper()
		rec := doJSON(t, s, http.MethodPost, "/v1/schedules/preview", PreviewRequest{
			Schedule: schedule, ScheduleSyntax: syntax, Timezone: "UTC",
		})
		if rec.Code != http.StatusOK {
			t.Fatalf("preview %q (%q): status=%d body=%s", schedule, syntax, rec.Code, rec.Body.String())
		}
		var resp PreviewResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}
		return resp
	}

	autoCron := preview("0 9 * * 1-5", "")
	forcedCron := preview("0 9 * * 1-5", "cron")
	human := preview("weekdays at 09:00", "human")
	if autoCron.SourceSyntax != "cron" || forcedCron.SourceSyntax != "cron" || human.SourceSyntax != "human" {
		t.Fatalf("source identities = auto %q, forced %q, human %q", autoCron.SourceSyntax, forcedCron.SourceSyntax, human.SourceSyntax)
	}
	if autoCron.RRULE != human.RRULE || forcedCron.RRULE != human.RRULE {
		t.Fatalf("RRULE mismatch: auto=%q forced=%q human=%q", autoCron.RRULE, forcedCron.RRULE, human.RRULE)
	}
	if len(autoCron.NextRuns) != len(human.NextRuns) {
		t.Fatalf("run count mismatch: cron=%d human=%d", len(autoCron.NextRuns), len(human.NextRuns))
	}
	for i := range autoCron.NextRuns {
		if !autoCron.NextRuns[i].Equal(human.NextRuns[i]) {
			t.Fatalf("run %d mismatch: cron=%s human=%s", i, autoCron.NextRuns[i], human.NextRuns[i])
		}
	}
}

func TestPreviewAndCreate_OrdinalCronParity(t *testing.T) {
	s := newTestServer(t)
	preview := doJSON(t, s, http.MethodPost, "/v1/schedules/preview", PreviewRequest{
		Schedule: "0 9 * * 5#3", ScheduleSyntax: "cron", Timezone: "UTC",
	})
	if preview.Code != http.StatusOK {
		t.Fatalf("preview status=%d body=%s", preview.Code, preview.Body.String())
	}
	created := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "third-friday", Command: "/bin/true", Schedule: "0 9 * * 5#3",
		ScheduleSyntax: "cron", Timezone: "UTC",
	})
	if created.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", created.Code, created.Body.String())
	}
	var resp TaskResponse
	if err := json.Unmarshal(created.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Schedule.Expression != "0 9 * * 5#3" || resp.Schedule.SourceSyntax != "cron" {
		t.Fatalf("created schedule = %+v", resp.Schedule)
	}
}

func TestPreviewAndCreate_LastWeekdayCronParity(t *testing.T) {
	s := newTestServer(t)
	preview := doJSON(t, s, http.MethodPost, "/v1/schedules/preview", PreviewRequest{
		Schedule: "0 9 * * 5L", ScheduleSyntax: "cron", Timezone: "UTC",
	})
	if preview.Code != http.StatusOK {
		t.Fatalf("preview status=%d body=%s", preview.Code, preview.Body.String())
	}
	created := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "last-friday", Command: "/bin/true", Schedule: "0 9 * * 5L",
		ScheduleSyntax: "cron", Timezone: "UTC",
	})
	if created.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", created.Code, created.Body.String())
	}
	var resp TaskResponse
	if err := json.Unmarshal(created.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Schedule.Expression != "0 9 * * 5L" || resp.Schedule.SourceSyntax != "cron" {
		t.Fatalf("created schedule = %+v", resp.Schedule)
	}
}

func TestCreateTask_CronSourceIdentity(t *testing.T) {
	s := newTestServer(t)
	for _, syntax := range []string{"", "cron"} {
		rec := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
			Name: "cron-" + syntax, Command: "/bin/true", Schedule: "0 9 * * 1-5",
			ScheduleSyntax: syntax, Timezone: "UTC",
		})
		if rec.Code != http.StatusCreated {
			t.Fatalf("syntax %q: status=%d body=%s", syntax, rec.Code, rec.Body.String())
		}
		var resp TaskResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}
		if resp.Schedule.Expression != "0 9 * * 1-5" || resp.Schedule.SourceSyntax != "cron" {
			t.Fatalf("syntax %q returned schedule %+v", syntax, resp.Schedule)
		}
	}
}

func TestScheduleSyntaxValidation(t *testing.T) {
	s := newTestServer(t)
	tests := []struct {
		name  string
		path  string
		body  any
		field string
	}{
		{"invalid hint", "/v1/schedules/preview", PreviewRequest{Schedule: "every day", ScheduleSyntax: "natural"}, "schedule_syntax"},
		{"forced cron never falls back", "/v1/schedules/preview", PreviewRequest{Schedule: "every day at 09:00", ScheduleSyntax: "cron"}, "schedule"},
		{"invalid cron never falls back", "/v1/schedules/preview", PreviewRequest{Schedule: "61 9 * * *"}, "schedule"},
		{"preview hint without schedule", "/v1/schedules/preview", PreviewRequest{ScheduleSyntax: "cron"}, "schedule_syntax"},
		{"unsupported named cron", "/v1/tasks", TaskCreateRequest{Name: "x", Command: "/bin/true", Schedule: "@reboot"}, "schedule"},
		{"lossy calendar step", "/v1/tasks", TaskCreateRequest{Name: "x", Command: "/bin/true", Schedule: "0 9 */2 * *"}, "schedule"},
		{"malformed ordinal", "/v1/tasks", TaskCreateRequest{Name: "x", Command: "/bin/true", Schedule: "0 9 * * 5#6"}, "schedule"},
		{"month-restricted ordinal", "/v1/tasks", TaskCreateRequest{Name: "x", Command: "/bin/true", Schedule: "0 9 * JAN 5#3"}, "schedule"},
		{"malformed last weekday", "/v1/tasks", TaskCreateRequest{Name: "x", Command: "/bin/true", Schedule: "0 9 * * 8L"}, "schedule"},
		{"month-restricted last weekday", "/v1/tasks", TaskCreateRequest{Name: "x", Command: "/bin/true", Schedule: "0 9 * JAN 5L"}, "schedule"},
		{"multiple last weekdays", "/v1/tasks", TaskCreateRequest{Name: "x", Command: "/bin/true", Schedule: "0 9 * * 5L,2L"}, "schedule"},
		{"hint without schedule", "/v1/tasks", TaskCreateRequest{Name: "x", Command: "/bin/true", ScheduleSyntax: "cron", At: futureTime()}, "schedule_syntax"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doJSON(t, s, http.MethodPost, tt.path, tt.body)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
			}
			var apiErr APIError
			if err := json.Unmarshal(rec.Body.Bytes(), &apiErr); err != nil {
				t.Fatal(err)
			}
			if apiErr.Error.Field != tt.field {
				t.Fatalf("field=%q want %q: %s", apiErr.Error.Field, tt.field, rec.Body.String())
			}
		})
	}
	list := doJSON(t, s, http.MethodGet, "/v1/tasks", nil)
	var listed struct {
		Tasks []domain.Task `json:"tasks"`
	}
	if err := json.Unmarshal(list.Body.Bytes(), &listed); err != nil {
		t.Fatal(err)
	}
	if len(listed.Tasks) != 0 {
		t.Fatalf("refused create mutated tasks: %s", list.Body.String())
	}
}

func futureTime() *time.Time {
	t := time.Now().UTC().Add(24 * time.Hour)
	return &t
}

func TestRunNow_NotFound(t *testing.T) {
	s := newTestServer(t)
	rec := doJSON(t, s, http.MethodPost, "/v1/tasks/missing/run-now", nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestListAndDeleteTask(t *testing.T) {
	s := newTestServer(t)
	create := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "x", Command: "/bin/true", Schedule: "every 15 minutes",
	})
	var resp TaskResponse
	_ = json.Unmarshal(create.Body.Bytes(), &resp)

	list := doJSON(t, s, http.MethodGet, "/v1/tasks", nil)
	if list.Code != http.StatusOK {
		t.Fatalf("list status %d", list.Code)
	}

	del := doJSON(t, s, http.MethodDelete, "/v1/tasks/"+resp.Task.ID, nil)
	if del.Code != http.StatusNoContent {
		t.Fatalf("delete status %d", del.Code)
	}
	get := doJSON(t, s, http.MethodGet, "/v1/tasks/"+resp.Task.ID, nil)
	if get.Code != http.StatusNotFound {
		t.Fatalf("get after delete status %d, want 404", get.Code)
	}
}
