package store

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestS143CancellationPreservesNotificationStreakAndRecovery(t *testing.T) {
	for _, canceled := range []domain.RunOutcome{domain.OutcomeCanceled, domain.OutcomeCancellationFailed} {
		t.Run(string(canceled), func(t *testing.T) {
			st := openMem(t)
			task := createNotificationTestTask(t, st, "")
			channel := createNotificationTestChannel(t, st, "cancellation fixture")
			assignment := domain.NotificationAssignment{ChannelID: channel.ID, OnSuccess: true, OnFailure: true, FailureThreshold: 2, OnRecovery: true, DurationThresholdSeconds: 1}
			if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, []domain.NotificationAssignment{assignment}); err != nil {
				t.Fatal(err)
			}
			base := time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)
			record := func(id string, outcome domain.RunOutcome, duration time.Duration) {
				t.Helper()
				start := base
				end := start.Add(duration)
				r := domain.Run{ID: id, TaskID: task.ID, ScheduledFor: start, StartedAt: &start, EndedAt: &end, Outcome: outcome, Trigger: domain.TriggerManual}
				if err := st.RecordRunAndCreateDeliveries(&r, ""); err != nil {
					t.Fatal(err)
				}
				base = base.Add(time.Minute)
			}
			count := func(want int) {
				t.Helper()
				d, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{})
				if err != nil {
					t.Fatal(err)
				}
				if len(d) != want {
					t.Fatalf("deliveries=%+v want %d", d, want)
				}
			}
			record("failure-one", domain.OutcomeFailure, 100*time.Millisecond)
			count(0)
			record("canceled-one", canceled, 60*time.Second)
			count(0)
			record("failure-two", domain.OutcomeFailure, 100*time.Millisecond)
			count(1)
			record("canceled-two", canceled, 60*time.Second)
			count(1)
			record("recovered", domain.OutcomeSuccess, 100*time.Millisecond)
			count(2)
			deliveries, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{})
			if err != nil {
				t.Fatal(err)
			}
			observed := map[string]domain.NotificationDelivery{}
			for _, d := range deliveries {
				observed[d.RunID] = d
			}
			if d, ok := observed["failure-two"]; !ok || d.ConditionKind != domain.NotificationConditionConsecutiveFailure {
				t.Fatalf("failure streak was reset by cancellation: %+v", deliveries)
			}
			if d, ok := observed["recovered"]; !ok || d.ConditionKind != domain.NotificationConditionRecovery {
				t.Fatalf("real success did not retain recovery: %+v", deliveries)
			}
		})
	}
}

func TestS143CancellationNeverCreatesCompletionDeliveries(t *testing.T) {
	for _, outcome := range []domain.RunOutcome{domain.OutcomeCanceled, domain.OutcomeCancellationFailed} {
		t.Run(string(outcome), func(t *testing.T) {
			st := openMem(t)
			source := createNotificationTestTask(t, st, "")
			for _, on := range []domain.CompletionOutcome{domain.CompletionOnSuccess, domain.CompletionOnFailure, domain.CompletionOnAny} {
				target := createNotificationTestTask(t, st, "")
				chain := domain.CompletionChain{SourceTaskID: source.ID, TargetTaskID: target.ID, OnOutcome: on}
				if err := st.CreateCompletionChain(&chain); err != nil {
					t.Fatal(err)
				}
			}
			now := time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)
			code := 1
			r := domain.Run{ID: "cancel-source", TaskID: source.ID, ScheduledFor: now, StartedAt: &now, EndedAt: &now, Outcome: outcome, ExitCode: &code, Output: "partial fixture output", OutputTruncated: true, Trigger: domain.TriggerManual}
			if err := st.RecordRunAndCreateDeliveries(&r, ""); err != nil {
				t.Fatal(err)
			}
			deliveries, err := st.ListCompletionDeliveries()
			if err != nil {
				t.Fatal(err)
			}
			if len(deliveries) != 0 {
				t.Fatalf("cancellation dispatched completion work: %+v", deliveries)
			}
			stored, err := st.GetRun(r.ID)
			if err != nil {
				t.Fatal(err)
			}
			if stored.Outcome != outcome || stored.Output != r.Output || !stored.OutputTruncated || stored.ExitCode == nil || *stored.ExitCode != code {
				t.Fatalf("cancellation native/output data lost: %+v", stored)
			}
		})
	}
}

func TestS143CancellationHistoryReopensCurrentStore(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cancellation.db")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	task := createNotificationTestTask(t, st, "")
	now := time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)
	for i, outcome := range []domain.RunOutcome{domain.OutcomeCanceled, domain.OutcomeCancellationFailed} {
		r := domain.Run{ID: fmt.Sprintf("cancel-%d", i), TaskID: task.ID, ScheduledFor: now, StartedAt: &now, EndedAt: &now, Outcome: outcome, Output: "retained fixture output", Trigger: domain.TriggerManual}
		if err := st.RecordRunAndCreateDeliveries(&r, ""); err != nil {
			t.Fatal(err)
		}
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := reopened.Close(); err != nil {
			t.Error(err)
		}
	}()
	for i, outcome := range []domain.RunOutcome{domain.OutcomeCanceled, domain.OutcomeCancellationFailed} {
		r, err := reopened.GetRun(fmt.Sprintf("cancel-%d", i))
		if err != nil {
			t.Fatal(err)
		}
		if r.Outcome != outcome || r.Output != "retained fixture output" {
			t.Fatalf("cancellation history changed after reopen: %+v", r)
		}
	}
}
