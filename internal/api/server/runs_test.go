package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

type activeSchedulerStub struct{ runs []domain.Run }

func (activeSchedulerStub) Reload()                    {}
func (activeSchedulerStub) RunNow(string) error        { return nil }
func (s activeSchedulerStub) ActiveRuns() []domain.Run { return s.runs }

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

func TestActiveRunsAndBoundedAlertsAPI(t *testing.T) {
	s := newTestServer(t)
	now := time.Now().UTC()
	s.sched = activeSchedulerStub{runs: []domain.Run{{ID: "active-1", TaskID: "task-1", ScheduledFor: now, StartedAt: &now, Trigger: domain.TriggerManual}}}
	active := doJSON(t, s, http.MethodGet, "/v1/runs/active", nil)
	if active.Code != http.StatusOK || !strings.Contains(active.Body.String(), `"id":"active-1"`) {
		t.Fatalf("active status=%d body=%s", active.Code, active.Body.String())
	}
	for index := 0; index < 3; index++ {
		alert := domain.Alert{Severity: domain.SeverityInfo, Kind: domain.AlertMissedRun, Message: "notice", CreatedAt: now.Add(time.Duration(index) * time.Second)}
		if err := s.store.CreateAlert(&alert); err != nil {
			t.Fatal(err)
		}
	}
	limited := doJSON(t, s, http.MethodGet, "/v1/alerts?limit=2", nil)
	var decoded struct {
		Alerts []domain.Alert `json:"alerts"`
	}
	if err := json.Unmarshal(limited.Body.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if limited.Code != http.StatusOK || len(decoded.Alerts) != 2 {
		t.Fatalf("limited status=%d alerts=%+v", limited.Code, decoded.Alerts)
	}
	invalid := doJSON(t, s, http.MethodGet, "/v1/alerts?limit=-1", nil)
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("invalid status=%d body=%s", invalid.Code, invalid.Body.String())
	}
}
