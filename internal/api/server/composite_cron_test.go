package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestCompositeCronPreviewCreateEditAndReload(t *testing.T) {
	s := newTestServer(t)
	expr := "*/10 9-17 * * MON,WED,FRI"
	preview := doJSON(t, s, http.MethodPost, "/v1/schedules/preview", PreviewRequest{
		Schedule: expr, ScheduleSyntax: "cron", Timezone: "America/New_York",
	})
	if preview.Code != http.StatusOK {
		t.Fatalf("preview status=%d body=%s", preview.Code, preview.Body.String())
	}
	var p PreviewResponse
	if err := json.Unmarshal(preview.Body.Bytes(), &p); err != nil {
		t.Fatal(err)
	}
	if p.SourceSyntax != "cron" || !strings.Contains(strings.ToLower(p.HumanSummary), "every 10 minutes") || len(p.NextRuns) == 0 {
		t.Fatalf("preview=%+v", p)
	}

	created := newTaskFor(t, s, TaskCreateRequest{
		Name: "composite", Command: "/bin/true", Schedule: expr,
		ScheduleSyntax: "cron", Timezone: "America/New_York",
	})
	if created.Schedule.Expression != expr || created.Schedule.SourceSyntax != "cron" {
		t.Fatalf("created schedule=%+v", created.Schedule)
	}

	updatedExpr := "0 9,17 1,15 JAN,MAR *"
	path := "/v1/tasks/" + created.Task.ID
	updated := doJSON(t, s, http.MethodPatch, path, TaskUpdateRequest{
		Schedule: updatedExpr, ScheduleSyntax: "cron",
	})
	if updated.Code != http.StatusOK {
		t.Fatalf("update status=%d body=%s", updated.Code, updated.Body.String())
	}
	loaded := getTask(t, s, created.Task.ID)
	if loaded.Schedule.Expression != updatedExpr || loaded.Schedule.SourceSyntax != "cron" || len(loaded.NextRuns) == 0 {
		t.Fatalf("loaded=%+v", loaded)
	}
}

func TestSecondsCronPreviewCreateAndReload(t *testing.T) {
	s := newTestServer(t)
	expr := "5/15 * * * * *"
	preview := doJSON(t, s, http.MethodPost, "/v1/schedules/preview", PreviewRequest{
		Schedule: expr, ScheduleSyntax: "cron", Timezone: "UTC",
	})
	if preview.Code != http.StatusOK {
		t.Fatalf("preview status=%d body=%s", preview.Code, preview.Body.String())
	}
	var p PreviewResponse
	if err := json.Unmarshal(preview.Body.Bytes(), &p); err != nil {
		t.Fatal(err)
	}
	for _, run := range p.NextRuns {
		if run.Second() != 5 && run.Second() != 20 && run.Second() != 35 && run.Second() != 50 {
			t.Fatalf("preview run %s has unexpected second", run)
		}
	}
	created := newTaskFor(t, s, TaskCreateRequest{
		Name: "seconds", Command: "/bin/true", Schedule: expr, ScheduleSyntax: "cron", Timezone: "UTC",
	})
	loaded := getTask(t, s, created.Task.ID)
	if loaded.Schedule.Expression != expr || loaded.Schedule.SourceSyntax != "cron" || !strings.Contains(loaded.Schedule.RRULE, "BYSECOND=5,20,35,50") {
		t.Fatalf("loaded schedule=%+v", loaded.Schedule)
	}
}

func TestCompositeCronRefusedUpdateDoesNotMutate(t *testing.T) {
	s := newTestServer(t)
	created := newTaskFor(t, s, TaskCreateRequest{
		Name: "stable", Command: "/bin/true", Schedule: "0 9,17 * * *", ScheduleSyntax: "cron", Timezone: "UTC",
	})
	path := "/v1/tasks/" + created.Task.ID
	rec := doJSON(t, s, http.MethodPatch, path, TaskUpdateRequest{Schedule: "0 0 13 * 5", ScheduleSyntax: "cron"})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	loaded := getTask(t, s, created.Task.ID)
	if loaded.Schedule.Expression != "0 9,17 * * *" {
		t.Fatalf("refused update changed schedule to %q", loaded.Schedule.Expression)
	}
}

func TestCompositeCronPreviewAndTaskNameMissingDatePolicy(t *testing.T) {
	s := newTestServer(t)
	preview := doJSON(t, s, http.MethodPost, "/v1/schedules/preview", PreviewRequest{
		Schedule: "0 9 1,31 * *", ScheduleSyntax: "cron", Timezone: "UTC",
		MissingDatePolicy: "last_valid",
	})
	if preview.Code != http.StatusOK {
		t.Fatalf("preview status=%d body=%s", preview.Code, preview.Body.String())
	}
	var p PreviewResponse
	if err := json.Unmarshal(preview.Body.Bytes(), &p); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(p.HumanSummary, "last day of the month") {
		t.Fatalf("preview summary=%q, want effective missing-date policy", p.HumanSummary)
	}

	created := newTaskFor(t, s, TaskCreateRequest{
		Name: "date-set", Command: "/bin/true", Schedule: "0 9 1,31 * *",
		ScheduleSyntax: "cron", Timezone: "UTC", MissingDatePolicy: "next_valid",
	})
	loaded := getTask(t, s, created.Task.ID)
	if !strings.Contains(loaded.Schedule.HumanSummary, "rolling into the next period") {
		t.Fatalf("task summary=%q, want effective missing-date policy", loaded.Schedule.HumanSummary)
	}
}
