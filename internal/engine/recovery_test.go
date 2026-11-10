package engine

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

func TestS145NoStartResolutionsWakeRecoveryAndCleanupWithoutNativeResults(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	entry := setupTask(t, st, domain.OverlapAllowConcurrent)
	branch := setupTask(t, st, domain.OverlapAllowConcurrent)
	target := setupTask(t, st, domain.OverlapAllowConcurrent)
	recovery := setupTask(t, st, domain.OverlapAllowConcurrent)
	cleanup := setupTask(t, st, domain.OverlapAllowConcurrent)
	join := domain.ParallelJoin{Name: "no start", EntryTaskID: entry.ID, TargetTaskID: target.ID, RecoveryTaskID: recovery.ID, RecoveryOn: []domain.JoinRecoveryCondition{domain.JoinRecoveryOnNotStarted}, CleanupTaskID: cleanup.ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "A", TaskID: branch.ID}}}
	if err := st.CreateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	if err := st.SetTaskEnabled(branch.ID, false); err != nil {
		t.Fatal(err)
	}
	if err := st.SetTaskEnabled(recovery.ID, false); err != nil {
		t.Fatal(err)
	}
	runner := &completionRunner{outcomes: map[string]domain.RunOutcome{}, runs: make(chan domain.Run, 8)}
	e, stop := startCompletionEngine(t, st, runner)
	defer stop()
	e.dispatch(entry, e.clk.Now(), domain.TriggerManual)
	waitRun(t, runner.runs, cleanup.ID)
	waitFor(t, func() bool {
		v, er := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{})
		return er == nil && len(v) == 1 && v[0].State == domain.JoinInstanceCompleted
	}, "no-start-only continuation chain")
	v, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if v[0].SelectedRoute != domain.JoinRouteRecovery || v[0].Branches[0].Settlement != domain.JoinSettlementNotStarted || v[0].RecoverySettlement != domain.JoinSettlementNotStarted || v[0].RecoveryRunID != "" || v[0].RecoveryReason == "" || v[0].CleanupOutcome != domain.OutcomeSuccess {
		t.Fatalf("no-start route=%+v", v[0])
	}
	runs, err := st.ListRuns("", 0)
	if err != nil || len(runs) != 2 {
		t.Fatalf("fabricated native runs=%+v err=%v", runs, err)
	}
}

type recoveryGatedRunner struct {
	gates    map[string]chan struct{}
	outcomes map[string]domain.RunOutcome
	started  chan domain.Run
}

func (r *recoveryGatedRunner) Run(_ context.Context, task domain.Task, sf time.Time, trigger domain.RunTrigger) domain.Run {
	run := domain.Run{TaskID: task.ID, ScheduledFor: sf, Trigger: trigger, Outcome: r.outcomes[task.ID]}
	if run.Outcome == "" {
		run.Outcome = domain.OutcomeSuccess
	}
	r.started <- run
	if gate := r.gates[task.ID]; gate != nil {
		<-gate
	}
	ended := sf.Add(time.Second)
	run.EndedAt = &ended
	return run
}
func TestS145RecoveryWaitsForSiblingsAndCleanupWaitsForFollowup(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	entry := setupTask(t, st, domain.OverlapAllowConcurrent)
	a := setupTask(t, st, domain.OverlapAllowConcurrent)
	b := setupTask(t, st, domain.OverlapAllowConcurrent)
	target := setupTask(t, st, domain.OverlapAllowConcurrent)
	recovery := setupTask(t, st, domain.OverlapAllowConcurrent)
	cleanup := setupTask(t, st, domain.OverlapAllowConcurrent)
	join := domain.ParallelJoin{Name: "owned ordering", EntryTaskID: entry.ID, TargetTaskID: target.ID, RecoveryTaskID: recovery.ID, CleanupTaskID: cleanup.ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "A", TaskID: a.ID}, {ID: "b", Name: "B", TaskID: b.ID}}}
	if err := st.CreateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	gates := map[string]chan struct{}{a.ID: make(chan struct{}), b.ID: make(chan struct{}), recovery.ID: make(chan struct{})}
	release := func(id string) {
		select {
		case <-gates[id]:
		default:
			close(gates[id])
		}
	}
	runner := &recoveryGatedRunner{gates: gates, outcomes: map[string]domain.RunOutcome{a.ID: domain.OutcomeFailure, recovery.ID: domain.OutcomeFailure}, started: make(chan domain.Run, 12)}
	e, stop := startCompletionEngine(t, st, runner)
	defer stop()
	defer func() {
		for id := range gates {
			release(id)
		}
	}()
	committed := make(chan domain.Run, 12)
	e.SetOnRun(func(run domain.Run) { committed <- run })
	e.dispatch(entry, e.clk.Now(), domain.TriggerManual)
	seen := map[string]bool{}
	for !seen[a.ID] || !seen[b.ID] {
		select {
		case run := <-runner.started:
			seen[run.TaskID] = true
		case <-time.After(5 * time.Second):
			t.Fatal("parallel branches did not start")
		}
	}
	release(a.ID)
	waitRun(t, committed, a.ID)
	var heldID string
	for _, active := range e.ActiveRuns() {
		if active.TaskID == b.ID {
			heldID = active.ID
		}
	}
	if heldID == "" {
		t.Fatal("held sibling lost its native ownership")
	}
	if _, err := e.CancelRun(context.Background(), heldID); err != nil {
		t.Fatal(err)
	}
	e.cancellationFailed(heldID, errors.New("synthetic termination remains unconfirmed"))
	failedFeedback := false
	for _, active := range e.ActiveRuns() {
		failedFeedback = failedFeedback || (active.ID == heldID && active.CancellationStatus == "failed" && active.Outcome == "")
	}
	if !failedFeedback {
		t.Fatal("failed cancellation feedback did not retain active native ownership")
	}
	instances, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{})
	if err != nil || len(instances) != 1 {
		t.Fatalf("instances=%+v err=%v", instances, err)
	}
	if instances[0].RecoveryDeliveryID != "" || instances[0].CleanupDeliveryID != "" || instances[0].SelectedRoute != domain.JoinRoutePending {
		t.Fatalf("follow-up before sibling settled: %+v", instances[0])
	}
	release(b.ID)
	waitRun(t, runner.started, recovery.ID)
	instances, err = st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if instances[0].SelectedRoute != domain.JoinRouteRecovery || instances[0].CleanupDeliveryID != "" {
		t.Fatalf("cleanup before recovery terminated: %+v", instances[0])
	}
	release(recovery.ID)
	waitRun(t, committed, cleanup.ID)
	waitFor(t, func() bool {
		v, er := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{})
		return er == nil && len(v) == 1 && v[0].State == domain.JoinInstanceCompleted
	}, "cleanup final result")
	instances, err = st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if instances[0].TargetRunID != "" || instances[0].RecoveryOutcome != domain.OutcomeFailure || instances[0].RecoverySettlement != domain.JoinSettlementFailure || instances[0].CleanupOutcome != domain.OutcomeSuccess || instances[0].RecoveryRunID == instances[0].CleanupRunID {
		t.Fatalf("stage evidence=%+v", instances[0])
	}
}

func TestS145DelayedRecoveryRejectionWakesCleanup(t *testing.T) {
	for _, boundary := range []string{"queue", "worker"} {
		t.Run(boundary, func(t *testing.T) {
			st, err := store.Open(":memory:")
			if err != nil {
				t.Fatal(err)
			}
			defer st.Close()
			entry := setupTask(t, st, domain.OverlapAllowConcurrent)
			branch := setupTask(t, st, domain.OverlapAllowConcurrent)
			target := setupTask(t, st, domain.OverlapAllowConcurrent)
			recovery := setupTask(t, st, domain.OverlapQueueOne)
			cleanup := setupTask(t, st, domain.OverlapAllowConcurrent)
			join := domain.ParallelJoin{Name: "delayed recovery", EntryTaskID: entry.ID, TargetTaskID: target.ID, RecoveryTaskID: recovery.ID, CleanupTaskID: cleanup.ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "A", TaskID: branch.ID}}}
			if err := st.CreateParallelJoin(&join); err != nil {
				t.Fatal(err)
			}
			sf := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
			source := domain.Run{TaskID: entry.ID, ScheduledFor: sf, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
			if err := st.RecordRunAndCreateDeliveries(&source, ""); err != nil {
				t.Fatal(err)
			}
			claims, err := st.ClaimCompletionDeliveries(10)
			if err != nil || len(claims) != 1 {
				t.Fatalf("branch claims=%+v err=%v", claims, err)
			}
			failed := domain.Run{TaskID: branch.ID, ScheduledFor: sf, Outcome: domain.OutcomeFailure, Trigger: domain.TriggerCompletion}
			if err := st.RecordRunAndCreateDeliveries(&failed, claims[0].ID); err != nil {
				t.Fatal(err)
			}
			// Startup recovery returns the pending follow-up to the ordinary admission path.
			runner := &completionRunner{outcomes: map[string]domain.RunOutcome{}, runs: make(chan domain.Run, 8)}
			e, stop := startCompletionEngine(t, st, runner, func(e *Engine) {
				if boundary == "queue" {
					e.running[recovery.ID] = 1
				} else {
					for range cap(e.sem) {
						e.sem <- struct{}{}
					}
				}
			})
			defer stop()
			capacityReleased := false
			releaseCapacity := func() {
				if boundary == "worker" && !capacityReleased {
					capacityReleased = true
					for range cap(e.sem) {
						<-e.sem
					}
				}
			}
			defer releaseCapacity()
			waiting, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{})
			if err != nil || len(waiting) != 1 || waiting[0].RecoverySettlement != domain.JoinSettlementPending || waiting[0].CleanupDeliveryID != "" {
				t.Fatalf("delayed recovery settled early=%+v err=%v", waiting, err)
			}
			if err := st.SetTaskEnabled(recovery.ID, false); err != nil {
				t.Fatal(err)
			}
			if boundary == "queue" {
				e.finish(recovery)
			} else {
				releaseCapacity()
			}
			waitRun(t, runner.runs, cleanup.ID)
			waitFor(t, func() bool {
				v, er := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{})
				return er == nil && len(v) == 1 && v[0].State == domain.JoinInstanceCompleted
			}, "delayed recovery no-start cleanup")
			got, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{})
			if err != nil {
				t.Fatal(err)
			}
			if got[0].RecoverySettlement != domain.JoinSettlementNotStarted || got[0].RecoveryRunID != "" || got[0].RecoveryReason == "" || got[0].CleanupOutcome != domain.OutcomeSuccess {
				t.Fatalf("delayed stage evidence=%+v", got[0])
			}
		})
	}
}
