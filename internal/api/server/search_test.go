package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestSearchReturnsBoundedSafeProjectionAndUpcomingOccurrence(t *testing.T) {
	s := newTestServer(t)
	created := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{Name: "Archive", Command: "secret-command", Args: []string{"secret-argument"}, Env: map[string]string{"TOKEN": "secret-environment"}, Schedule: "every day at 09:00", Timezone: "UTC"})
	var task TaskResponse
	if err := json.Unmarshal(created.Body.Bytes(), &task); err != nil {
		t.Fatal(err)
	}
	ended := time.Now().UTC().Add(-time.Hour)
	run := &domain.Run{ID: "run-archive", TaskID: task.Task.ID, ScheduledFor: ended, EndedAt: &ended, Outcome: domain.OutcomeFailure, Output: "secret-output", Trigger: domain.TriggerSchedule}
	if err := s.store.CreateRun(run); err != nil {
		t.Fatal(err)
	}
	if err := s.store.CreateAlert(&domain.Alert{ID: "alert-archive", TaskID: task.Task.ID, RunID: run.ID, Severity: domain.SeverityError, Kind: domain.AlertRunFailed, Message: "secret-alert"}); err != nil {
		t.Fatal(err)
	}

	response := doJSON(t, s, http.MethodGet, "/v1/search?q=archive&limit=50", nil)
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	for _, secret := range []string{"secret-command", "secret-argument", "secret-environment", "secret-output", "secret-alert"} {
		if strings.Contains(response.Body.String(), secret) {
			t.Fatalf("response exposed %q: %s", secret, response.Body.String())
		}
	}
	var result domain.DaemonSearch
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Schema != domain.DaemonSearchSchema || result.Query != "archive" || result.ObservedAt.IsZero() {
		t.Fatalf("search metadata = %+v", result)
	}
	kinds := map[domain.SearchKind]bool{}
	for _, match := range result.Results {
		kinds[match.Kind] = true
		if match.Kind == domain.SearchKindSchedule && match.OccurredAt == nil {
			t.Fatalf("schedule result has no occurrence: %+v", match)
		}
	}
	for _, kind := range []domain.SearchKind{domain.SearchKindTask, domain.SearchKindFailure, domain.SearchKindSchedule, domain.SearchKindAlert} {
		if !kinds[kind] {
			t.Fatalf("missing %q in %+v", kind, result.Results)
		}
	}
}

func TestSearchValidatesQueryKindsAndLimit(t *testing.T) {
	s := newTestServer(t)
	for _, path := range []string{
		"/v1/search",
		"/v1/search?q=%20%20",
		"/v1/search?q=valid&kind=output",
		"/v1/search?q=valid&limit=0",
		"/v1/search?q=valid&limit=51",
		"/v1/search?q=" + strings.Repeat("x", 201),
	} {
		response := doJSON(t, s, http.MethodGet, path, nil)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("%s status=%d body=%s", path, response.Code, response.Body.String())
		}
	}
}

func TestSearchReportsTruncationAndKindFilter(t *testing.T) {
	s := newTestServer(t)
	for _, name := range []string{"Archive A", "Archive B"} {
		created := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{Name: name, Command: "echo", Schedule: "every day at 09:00", Timezone: "UTC"})
		if created.Code != http.StatusCreated {
			t.Fatalf("create status=%d body=%s", created.Code, created.Body.String())
		}
	}
	response := doJSON(t, s, http.MethodGet, "/v1/search?q=archive&kind=task&limit=1", nil)
	var result domain.DaemonSearch
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if response.Code != http.StatusOK || !result.Truncated || len(result.Results) != 1 || result.Results[0].Kind != domain.SearchKindTask {
		t.Fatalf("status=%d result=%+v", response.Code, result)
	}
}
