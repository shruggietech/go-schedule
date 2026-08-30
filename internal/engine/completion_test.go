package engine

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/clock"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

type completionRunner struct {
	mu       sync.Mutex
	outcomes map[string]domain.RunOutcome
	runs     chan domain.Run
}

func (r *completionRunner) Run(_ context.Context, task domain.Task, scheduledFor time.Time, trigger domain.RunTrigger) domain.Run {
	r.mu.Lock()
	outcome := r.outcomes[task.ID]
	r.mu.Unlock()
	if outcome == "" {
		outcome = domain.OutcomeSuccess
	}
	run := domain.Run{TaskID: task.ID, ScheduledFor: scheduledFor, Outcome: outcome, Trigger: trigger}
	r.runs <- run
	return run
}

func addChain(t *testing.T, st *store.Store, source, target domain.Task, outcome domain.CompletionOutcome) domain.CompletionChain {
	t.Helper()
	chain := domain.CompletionChain{SourceTaskID: source.ID, TargetTaskID: target.ID, OnOutcome: outcome}
	if err := st.CreateCompletionChain(&chain); err != nil {
		t.Fatal(err)
	}
	return chain
}

func waitRun(t *testing.T, ch <-chan domain.Run, taskID string) domain.Run {
	t.Helper()
	deadline := time.After(2 * time.Second)
	for {
		select {
		case run := <-ch:
			if run.TaskID == taskID {
				return run
			}
		case <-deadline:
			t.Fatalf("timeout waiting for task %s", taskID)
		}
	}
}

func TestCompletionChainDispatchesMultiHopWithCorrelation(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	a := setupTask(t, st, domain.OverlapQueueOne)
	b := setupTask(t, st, domain.OverlapQueueOne)
	c := setupTask(t, st, domain.OverlapQueueOne)
	addChain(t, st, a, b, domain.CompletionOnSuccess)
	addChain(t, st, b, c, domain.CompletionOnAny)
	runner := &completionRunner{outcomes: map[string]domain.RunOutcome{}, runs: make(chan domain.Run, 8)}
	e := newEngine(st, runner)

	e.dispatch(a, time.Now().UTC(), domain.TriggerManual)
	_ = waitRun(t, runner.runs, a.ID)
	_ = waitRun(t, runner.runs, b.ID)
	_ = waitRun(t, runner.runs, c.ID)
	waitFor(t, func() bool {
		runs, _ := st.ListRuns("", 0)
		return len(runs) == 3
	}, "three persisted chain runs")

	bRuns, _ := st.ListRuns(b.ID, 0)
	cRuns, _ := st.ListRuns(c.ID, 0)
	if len(bRuns) != 1 || bRuns[0].Trigger != domain.TriggerCompletion || bRuns[0].SourceTaskID != a.ID || bRuns[0].SourceRunID == "" {
		t.Fatalf("B correlation = %+v", bRuns)
	}
	if len(cRuns) != 1 || cRuns[0].Trigger != domain.TriggerCompletion || cRuns[0].SourceTaskID != b.ID || cRuns[0].SourceRunID == "" {
		t.Fatalf("C correlation = %+v", cRuns)
	}
	deliveries, _ := st.ListCompletionDeliveries()
	if len(deliveries) != 2 || deliveries[0].State != domain.DeliveryCompleted || deliveries[1].State != domain.DeliveryCompleted {
		t.Fatalf("delivery states = %+v", deliveries)
	}
}

func TestCompletionChainHonorsOutcomeAndResolvesIneligibleTarget(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	source := setupTask(t, st, domain.OverlapQueueOne)
	successTarget := setupTask(t, st, domain.OverlapQueueOne)
	failureTarget := setupTask(t, st, domain.OverlapQueueOne)
	anyTarget := setupTask(t, st, domain.OverlapQueueOne)
	addChain(t, st, source, successTarget, domain.CompletionOnSuccess)
	addChain(t, st, source, failureTarget, domain.CompletionOnFailure)
	addChain(t, st, source, anyTarget, domain.CompletionOnAny)
	if err := st.SetTaskEnabled(anyTarget.ID, false); err != nil {
		t.Fatal(err)
	}
	runner := &completionRunner{
		outcomes: map[string]domain.RunOutcome{source.ID: domain.OutcomeFailure},
		runs:     make(chan domain.Run, 8),
	}
	e := newEngine(st, runner)
	e.dispatch(source, time.Now().UTC(), domain.TriggerManual)
	_ = waitRun(t, runner.runs, source.ID)
	_ = waitRun(t, runner.runs, failureTarget.ID)
	waitFor(t, func() bool {
		deliveries, _ := st.ListCompletionDeliveries()
		return len(deliveries) == 2 && deliveries[0].State != domain.DeliveryClaimed && deliveries[1].State != domain.DeliveryClaimed
	}, "failure and any delivery resolution")
	if runs, _ := st.ListRuns(successTarget.ID, 0); len(runs) != 0 {
		t.Fatalf("success-only target ran after failure: %+v", runs)
	}
	if runs, _ := st.ListRuns(anyTarget.ID, 0); len(runs) != 0 {
		t.Fatalf("disabled any-outcome target ran: %+v", runs)
	}
	deliveries, _ := st.ListCompletionDeliveries()
	states := map[string]domain.DeliveryState{}
	for _, delivery := range deliveries {
		states[delivery.TargetTaskID] = delivery.State
	}
	if states[failureTarget.ID] != domain.DeliveryCompleted || states[anyTarget.ID] != domain.DeliveryResolved {
		t.Fatalf("delivery states by target = %+v", states)
	}
}

func TestCompletionDeliveryReplaysOnEngineStart(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	source := setupTask(t, st, domain.OverlapQueueOne)
	target := setupTask(t, st, domain.OverlapQueueOne)
	addChain(t, st, source, target, domain.CompletionOnSuccess)
	sourceRun := domain.Run{TaskID: source.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&sourceRun, ""); err != nil {
		t.Fatal(err)
	}
	claimed, err := st.ClaimCompletionDeliveries(100)
	if err != nil || len(claimed) != 1 {
		t.Fatalf("pre-crash claim = %+v, err=%v", claimed, err)
	}

	runner := &completionRunner{outcomes: map[string]domain.RunOutcome{}, runs: make(chan domain.Run, 8)}
	e := New(st, clock.NewReal(), runner, testLogger(), 2)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- e.Start(ctx) }()
	_ = waitRun(t, runner.runs, target.ID)
	waitFor(t, func() bool {
		deliveries, _ := st.ListCompletionDeliveries()
		return len(deliveries) == 1 && deliveries[0].State == domain.DeliveryCompleted && deliveries[0].Attempts == 2
	}, "replayed delivery completion")
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("engine did not stop")
	}
}

func TestCompletionDeliveryResolvesWhenChainDeletedBeforeClaim(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	source := setupTask(t, st, domain.OverlapQueueOne)
	target := setupTask(t, st, domain.OverlapQueueOne)
	chain := addChain(t, st, source, target, domain.CompletionOnSuccess)
	sourceRun := domain.Run{TaskID: source.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&sourceRun, ""); err != nil {
		t.Fatal(err)
	}
	if err := st.DeleteCompletionChain(chain.ID); err != nil {
		t.Fatal(err)
	}
	runner := &completionRunner{outcomes: map[string]domain.RunOutcome{}, runs: make(chan domain.Run, 2)}
	e := newEngine(st, runner)
	e.drainCompletionDeliveries()
	deliveries, err := st.ListCompletionDeliveries()
	if err != nil || len(deliveries) != 1 || deliveries[0].State != domain.DeliveryResolved {
		t.Fatalf("deleted-chain delivery = %+v, err=%v", deliveries, err)
	}
	select {
	case run := <-runner.runs:
		t.Fatalf("deleted chain dispatched target: %+v", run)
	default:
	}
}
