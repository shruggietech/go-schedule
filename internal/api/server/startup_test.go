package server

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestStartupPreviewCreateAndUpdate(t *testing.T) {
	s := newTestServer(t)
	for _, tc := range []struct{ expression, syntax string }{{"@reboot", "cron"}, {"at scheduler startup", "human"}} {
		rec := doJSON(t, s, http.MethodPost, "/v1/schedules/preview", PreviewRequest{Schedule: tc.expression, ScheduleSyntax: tc.syntax, Timezone: "UTC"})
		if rec.Code != http.StatusOK {
			t.Fatalf("preview %q: status=%d body=%s", tc.expression, rec.Code, rec.Body.String())
		}
		var preview PreviewResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &preview); err != nil {
			t.Fatal(err)
		}
		if preview.HumanSummary != "At scheduler startup" || preview.SourceSyntax != tc.syntax || preview.PolicySummary != "" || len(preview.NextRuns) != 0 {
			t.Fatalf("preview %q = %+v", tc.expression, preview)
		}
	}
	created := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{Name: "boot", Command: "x", Schedule: "@reboot", ScheduleSyntax: "cron", Timezone: "UTC"})
	if created.Code != http.StatusCreated {
		t.Fatalf("create: status=%d body=%s", created.Code, created.Body.String())
	}
	var detail TaskResponse
	if err := json.Unmarshal(created.Body.Bytes(), &detail); err != nil {
		t.Fatal(err)
	}
	if detail.Schedule.Kind != domain.ScheduleEvent || detail.Schedule.TriggerID != domain.StartupEventID || detail.Schedule.SourceSyntax != "cron" || detail.PolicySummary != "" || len(detail.NextRuns) != 0 {
		t.Fatalf("created detail = %+v", detail)
	}
	updated := doJSON(t, s, http.MethodPatch, "/v1/tasks/"+detail.Task.ID, TaskUpdateRequest{Schedule: "at scheduler startup", ScheduleSyntax: "human"})
	if updated.Code != http.StatusOK {
		t.Fatalf("update: status=%d body=%s", updated.Code, updated.Body.String())
	}
	if err := json.Unmarshal(updated.Body.Bytes(), &detail); err != nil {
		t.Fatal(err)
	}
	if detail.Schedule.SourceSyntax != "human" || detail.Schedule.Expression != "at scheduler startup" {
		t.Fatalf("updated schedule = %+v", detail.Schedule)
	}
}
