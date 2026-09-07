package operations

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type stubBackend struct {
	calendar server.CalendarResponse
	runs     []domain.Run
	active   []domain.Run
	logs     server.LogsResponse
	alerts   []domain.Alert
	errAt    string
	acked    []string
}

func (s *stubBackend) GetCalendar(context.Context, time.Time, time.Time) (server.CalendarResponse, error) {
	if s.errAt == "calendar" {
		return server.CalendarResponse{}, errors.New("boom")
	}
	return s.calendar, nil
}
func (s *stubBackend) ListRuns(context.Context, string, int) ([]domain.Run, error) {
	if s.errAt == "runs" {
		return nil, errors.New("boom")
	}
	return s.runs, nil
}
func (s *stubBackend) ListActiveRuns(context.Context) ([]domain.Run, error) {
	if s.errAt == "active" {
		return nil, errors.New("boom")
	}
	return s.active, nil
}
func (s *stubBackend) ListLogs(context.Context, string, int) (server.LogsResponse, error) {
	if s.errAt == "logs" {
		return server.LogsResponse{}, errors.New("boom")
	}
	return s.logs, nil
}
func (s *stubBackend) ListAlertsLimited(context.Context, bool, int) ([]domain.Alert, error) {
	if s.errAt == "alerts" {
		return nil, errors.New("boom")
	}
	return s.alerts, nil
}
func (s *stubBackend) AckAlert(_ context.Context, id string) error {
	s.acked = append(s.acked, id)
	return nil
}

func TestScheduleWindowDistinguishesPredictionsAndRecordedRuns(t *testing.T) {
	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	backend := &stubBackend{calendar: server.CalendarResponse{From: now.Add(-time.Hour), To: now.Add(time.Hour), Occurrences: []server.Occurrence{{TaskID: "task-1", TaskName: "Future", Time: now.Add(time.Hour), Kind: "scheduled"}, {TaskID: "task-2", TaskName: "Past", RunID: "run-1", Time: now.Add(-time.Minute), Kind: "past", Outcome: domain.OutcomeFailure}}}}
	service := NewService(backend)
	service.now = func() time.Time { return now }
	result := service.ScheduleWindow(context.Background(), 7)
	if result.Outcome != "accepted" || len(result.Schedule.Occurrences) != 2 {
		t.Fatalf("result=%+v", result)
	}
	if got := result.Schedule.Occurrences[0]; got.Kind != "recorded" || got.State != "failure" || got.ID != "run:run-1:0" {
		t.Fatalf("recorded=%+v", got)
	}
	if got := result.Schedule.Occurrences[1]; got.Kind != "prediction" || got.State != "upcoming" || got.RunID != "" {
		t.Fatalf("prediction=%+v", got)
	}
	if rejected := service.ScheduleWindow(context.Background(), 14); rejected.Outcome != "rejected" || rejected.Field != "days" {
		t.Fatalf("rejected=%+v", rejected)
	}
}

func TestActivityWorkspaceIsCompleteAndPreservesDiagnostics(t *testing.T) {
	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	exit := 7
	backend := &stubBackend{runs: []domain.Run{{ID: "run-1", TaskID: "task-1", ScheduledFor: now, Outcome: domain.OutcomeFailure, ExitCode: &exit, Output: "failed", OutputTruncated: true, Trigger: domain.RunTrigger("watcher"), SourceWatcherID: "watcher-1"}}, logs: server.LogsResponse{LogPath: `C:\ProgramData\goschedule\logs\goschedule.log`, Logs: []domain.LogRecord{{ID: "log-1", Time: now, Severity: domain.SeverityWarning, Message: "watch", Attrs: map[string]any{"z": 2, "a": 1}}}}, alerts: []domain.Alert{{ID: "alert-1", CreatedAt: now, Severity: domain.SeverityError, Kind: domain.AlertRunFailed, Message: "failed"}}}
	service := NewService(backend)
	service.now = func() time.Time { return now }
	result := service.ActivityWorkspace(context.Background())
	if result.Outcome != "accepted" || result.Activity.LogPath != backend.logs.LogPath {
		t.Fatalf("result=%+v", result)
	}
	if got := result.Activity.Runs[0]; got.State != "failure" || got.ExitCode == nil || *got.ExitCode != 7 || !got.OutputTruncated || got.SourceWatcherID != "watcher-1" {
		t.Fatalf("run=%+v", got)
	}
	if got := result.Activity.Logs[0].Detail; got != "a: 1\nz: 2" {
		t.Fatalf("detail=%q", got)
	}
	backend.errAt = "alerts"
	if failed := service.ActivityWorkspace(context.Background()); failed.Outcome != "unavailable" || failed.Activity != nil {
		t.Fatalf("partial=%+v", failed)
	}
}

func TestActivityWorkspaceIncludesAuthoritativeActiveRuns(t *testing.T) {
	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	backend := &stubBackend{active: []domain.Run{{ID: "active-1", TaskID: "task-1", ScheduledFor: now.Add(-time.Hour), StartedAt: &now, Trigger: domain.TriggerManual}}}
	service := NewService(backend)
	service.now = func() time.Time { return now }
	result := service.ActivityWorkspace(context.Background())
	if result.Outcome != "accepted" || len(result.Activity.Runs) != 1 || result.Activity.Runs[0].State != "running" {
		t.Fatalf("result=%+v", result)
	}
}

func TestActivityWorkspacePrefersPersistedCompletionDuringHandoff(t *testing.T) {
	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	backend := &stubBackend{runs: []domain.Run{{ID: "run-1", TaskID: "task-1", ScheduledFor: now, EndedAt: &now, Outcome: domain.OutcomeSuccess}}, active: []domain.Run{{ID: "run-1", TaskID: "task-1", ScheduledFor: now, StartedAt: &now}}}
	result := NewService(backend).ActivityWorkspace(context.Background())
	if len(result.Activity.Runs) != 1 || result.Activity.Runs[0].State != "success" {
		t.Fatalf("result=%+v", result)
	}
}

func TestRunStateUsesBoundedNonColorVocabulary(t *testing.T) {
	cases := []struct {
		outcome domain.RunOutcome
		started bool
		ended   bool
		want    string
	}{{domain.OutcomeSuccess, true, true, "success"}, {domain.OutcomeFailure, true, true, "failure"}, {domain.OutcomeSkipped, false, true, "skipped"}, {domain.OutcomeCaughtUp, true, true, "caught_up"}, {domain.OutcomeQueued, false, false, "queued"}, {"", true, false, "running"}, {"future_value", true, true, "unavailable"}}
	for _, test := range cases {
		if got := runState(test.outcome, test.started, test.ended); got != test.want {
			t.Errorf("runState(%q, %v, %v)=%q want %q", test.outcome, test.started, test.ended, got, test.want)
		}
	}
}

func TestAcknowledgeAlertsDeduplicatesAndRefreshes(t *testing.T) {
	backend := &stubBackend{}
	result := NewService(backend).AcknowledgeAlerts(context.Background(), []string{"alert-1", "", "alert-1", "alert-2"})
	if result.Outcome != "accepted" || result.Activity == nil {
		t.Fatalf("result=%+v", result)
	}
	if !reflect.DeepEqual(backend.acked, []string{"alert-1", "alert-2"}) {
		t.Fatalf("acked=%v", backend.acked)
	}
}
