package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestRunsAPIIncludesOptionalCompletionCorrelation(t *testing.T) {
	s := newTestServer(t)
	target := apiTask(t, s, "target")
	for _, run := range []domain.Run{
		{TaskID: target.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerSchedule},
		{TaskID: target.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerCompletion, SourceTaskID: "source", SourceRunID: "source-run"},
	} {
		if err := s.store.CreateRun(&run); err != nil {
			t.Fatal(err)
		}
	}
	response := doJSON(t, s, http.MethodGet, "/v1/runs?task="+target.ID, nil)
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	var decoded struct {
		Runs []map[string]any `json:"runs"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Runs) != 2 {
		t.Fatalf("runs=%+v", decoded.Runs)
	}
	correlated, plain := decoded.Runs[0], decoded.Runs[1]
	if correlated["trigger"] != string(domain.TriggerCompletion) {
		correlated, plain = plain, correlated
	}
	if correlated["source_task_id"] != "source" || correlated["source_run_id"] != "source-run" {
		t.Fatalf("completion run=%+v", correlated)
	}
	if _, exists := plain["source_task_id"]; exists {
		t.Fatalf("plain run exposed empty source_task_id: %+v", plain)
	}
	if _, exists := plain["source_run_id"]; exists {
		t.Fatalf("plain run exposed empty source_run_id: %+v", plain)
	}
}

func TestGetRunAPIUsesExactIdentityAndIncludesDiagnostics(t *testing.T) {
	s := newTestServer(t)
	task := apiTask(t, s, "diagnostic")
	run := domain.Run{
		TaskID: task.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeFailure,
		Trigger: domain.TriggerManual, Output: "partial", OutputTruncated: true,
	}
	if err := s.store.CreateRun(&run); err != nil {
		t.Fatal(err)
	}
	alert := domain.Alert{TaskID: task.ID, RunID: run.ID, Severity: domain.SeverityError, Kind: domain.AlertRunFailed, Message: "task run failed"}
	if err := s.store.CreateAlert(&alert); err != nil {
		t.Fatal(err)
	}

	response := doJSON(t, s, http.MethodGet, "/v1/runs/"+run.ID, nil)
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	var got domain.Run
	if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != run.ID || got.Output != "partial" || !got.OutputTruncated {
		t.Fatalf("run=%+v", got)
	}

	missing := doJSON(t, s, http.MethodGet, "/v1/runs/missing", nil)
	if missing.Code != http.StatusNotFound {
		t.Fatalf("missing status=%d body=%s", missing.Code, missing.Body.String())
	}
	alerts := doJSON(t, s, http.MethodGet, "/v1/alerts", nil)
	if alerts.Code != http.StatusOK || !strings.Contains(alerts.Body.String(), `"run_id":"`+run.ID+`"`) {
		t.Fatalf("alerts status=%d body=%s", alerts.Code, alerts.Body.String())
	}
}
