package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestSystemSummaryReturnsBoundedSafeProjection(t *testing.T) {
	s := newTestServer(t)
	created := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{Name: "Daily", Command: "secret-command", Args: []string{"secret-argument"}, Env: map[string]string{"TOKEN": "secret-environment"}, Schedule: "every day at 09:00", Timezone: "UTC"})
	var task TaskResponse
	if err := json.Unmarshal(created.Body.Bytes(), &task); err != nil {
		t.Fatal(err)
	}
	ended := time.Now().UTC().Add(-time.Hour)
	run := &domain.Run{TaskID: task.Task.ID, ScheduledFor: ended, EndedAt: &ended, Outcome: domain.OutcomeFailure, Output: "secret-output", Trigger: domain.TriggerSchedule}
	if err := s.store.CreateRun(run); err != nil {
		t.Fatal(err)
	}
	if err := s.store.CreateAlert(&domain.Alert{TaskID: task.Task.ID, RunID: run.ID, Severity: domain.SeverityError, Kind: domain.AlertRunFailed, Message: "secret-alert"}); err != nil {
		t.Fatal(err)
	}

	response := doJSON(t, s, http.MethodGet, "/v1/system-summary", nil)
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	for _, secret := range []string{"secret-command", "secret-argument", "secret-environment", "secret-output", "secret-alert"} {
		if strings.Contains(response.Body.String(), secret) {
			t.Fatalf("response exposed %q: %s", secret, response.Body.String())
		}
	}
	var summary domain.SystemSummary
	if err := json.Unmarshal(response.Body.Bytes(), &summary); err != nil {
		t.Fatal(err)
	}
	if summary.Schema != domain.SystemSummarySchema || summary.ActiveTaskCount != 1 || summary.RecentFailureCount != 1 || summary.UnacknowledgedAlertCount != 1 {
		t.Fatalf("summary=%+v", summary)
	}
}
