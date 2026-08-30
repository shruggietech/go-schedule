package engine

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/clock"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

// blockingRunner signals when a run starts and blocks until released, letting
// tests hold a task in the "running" state deterministically.
type blockingRunner struct {
	started chan struct{}
	release chan struct{}
}

func (r *blockingRunner) Run(_ context.Context, task domain.Task, sf time.Time, trig domain.RunTrigger) domain.Run {
	r.started <- struct{}{}
	<-r.release
	now := time.Now()
	return domain.Run{TaskID: task.ID, ScheduledFor: sf, EndedAt: &now, Outcome: domain.OutcomeSuccess, Trigger: trig}
}

func testLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func setupTask(t *testing.T, st *store.Store, policy domain.OverlapPolicy) domain.Task {
	t.Helper()
	sch := &domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=MINUTELY;INTERVAL=1"}
	if err := st.CreateSchedule(sch); err != nil {
		t.Fatal(err)
	}
	task := &domain.Task{
		Name: "t", Command: "x", Enabled: true, Timezone: "UTC", ScheduleID: sch.ID,
		OverlapPolicy: policy, CatchupPolicy: domain.CatchupNone, State: domain.TaskActive,
	}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	return *task
}

func newEngine(st *store.Store, r Runner) *Engine {
	e := New(st, clock.NewReal(), r, testLogger(), 4)
	e.runCtx = context.Background()
	return e
}

func recv(t *testing.T, ch <-chan struct{}, msg string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for %s", msg)
	}
}

func notRecv(t *testing.T, ch <-chan struct{}, msg string) {
	t.Helper()
	select {
	case <-ch:
		t.Fatalf("unexpected %s", msg)
	case <-time.After(100 * time.Millisecond):
	}
}

func countOutcomes(t *testing.T, st *store.Store, taskID string) map[domain.RunOutcome]int {
	t.Helper()
	runs, err := st.ListRuns(taskID, 0)
	if err != nil {
		t.Fatal(err)
	}
	m := map[domain.RunOutcome]int{}
	for _, r := range runs {
		m[r.Outcome]++
	}
	return m
}

func TestOverlap_QueueOne(t *testing.T) {
	st, _ := store.Open(":memory:")
	defer st.Close()
	r := &blockingRunner{started: make(chan struct{}, 1), release: make(chan struct{})}
	e := newEngine(st, r)
	task := setupTask(t, st, domain.OverlapQueueOne)
	now := time.Now().UTC()

	e.dispatch(task, now, domain.TriggerSchedule)
	recv(t, r.started, "first run start")

	// While running, dispatch twice more: first queues, second is dropped.
	e.dispatch(task, now.Add(time.Minute), domain.TriggerSchedule)
	e.dispatch(task, now.Add(2*time.Minute), domain.TriggerSchedule)
	notRecv(t, r.started, "second start before release")

	// Exactly one alert was raised for the overlap.
	if alerts, _ := st.ListAlerts(true); len(alerts) != 1 || alerts[0].Kind != domain.AlertOverlapQueued {
		t.Fatalf("expected 1 overlap alert, got %+v", alerts)
	}

	// Release the first run; the single queued run should start.
	r.release <- struct{}{}
	recv(t, r.started, "queued run start")
	r.release <- struct{}{}

	// Drain and verify outcomes: 2 successes, exactly 1 queued marker.
	time.Sleep(100 * time.Millisecond)
	out := countOutcomes(t, st, task.ID)
	if out[domain.OutcomeSuccess] != 2 {
		t.Fatalf("want 2 successes, got %d (%v)", out[domain.OutcomeSuccess], out)
	}
	if out[domain.OutcomeQueued] != 1 {
		t.Fatalf("want exactly 1 queued marker, got %d (%v)", out[domain.OutcomeQueued], out)
	}
}

func TestOverlap_QueuedRunRetainsTriggerOrigin(t *testing.T) {
	st, _ := store.Open(":memory:")
	defer st.Close()
	r := &blockingRunner{started: make(chan struct{}, 2), release: make(chan struct{})}
	e := newEngine(st, r)
	task := setupTask(t, st, domain.OverlapQueueOne)
	now := time.Now().UTC()
	e.dispatch(task, now, domain.TriggerStartup)
	recv(t, r.started, "startup run")
	e.dispatch(task, now.Add(time.Second), domain.TriggerManual)
	r.release <- struct{}{}
	recv(t, r.started, "queued manual run")
	r.release <- struct{}{}
	waitFor(t, func() bool {
		runs, _ := st.ListRuns(task.ID, 0)
		for _, run := range runs {
			if run.Outcome == domain.OutcomeSuccess && run.Trigger == domain.TriggerManual {
				return true
			}
		}
		return false
	}, "queued manual trigger persistence")
}

func TestOverlap_Skip(t *testing.T) {
	st, _ := store.Open(":memory:")
	defer st.Close()
	r := &blockingRunner{started: make(chan struct{}, 1), release: make(chan struct{})}
	e := newEngine(st, r)
	task := setupTask(t, st, domain.OverlapSkip)
	now := time.Now().UTC()

	e.dispatch(task, now, domain.TriggerSchedule)
	recv(t, r.started, "first run start")

	e.dispatch(task, now.Add(time.Minute), domain.TriggerSchedule) // should be skipped, not queued
	notRecv(t, r.started, "second start (skip policy)")

	r.release <- struct{}{}
	time.Sleep(100 * time.Millisecond)

	out := countOutcomes(t, st, task.ID)
	if out[domain.OutcomeSkipped] != 1 {
		t.Fatalf("want 1 skipped, got %v", out)
	}
	if out[domain.OutcomeSuccess] != 1 {
		t.Fatalf("want 1 success, got %v", out)
	}
}

func TestOverlap_AllowConcurrent(t *testing.T) {
	st, _ := store.Open(":memory:")
	defer st.Close()
	r := &blockingRunner{started: make(chan struct{}, 2), release: make(chan struct{})}
	e := newEngine(st, r)
	task := setupTask(t, st, domain.OverlapAllowConcurrent)
	now := time.Now().UTC()

	e.dispatch(task, now, domain.TriggerSchedule)
	recv(t, r.started, "first start")
	e.dispatch(task, now.Add(time.Minute), domain.TriggerSchedule)
	recv(t, r.started, "second concurrent start") // both run at once

	r.release <- struct{}{}
	r.release <- struct{}{}
	time.Sleep(100 * time.Millisecond)
	if out := countOutcomes(t, st, task.ID); out[domain.OutcomeSuccess] != 2 {
		t.Fatalf("want 2 concurrent successes, got %v", out)
	}
}

func claimedCompletion(t *testing.T, st *store.Store, source, target domain.Task, chainID, sourceRunID string) domain.CompletionDelivery {
	t.Helper()
	if chainID == "" {
		chain := domain.CompletionChain{SourceTaskID: source.ID, TargetTaskID: target.ID, OnOutcome: domain.CompletionOnSuccess}
		if err := st.CreateCompletionChain(&chain); err != nil {
			t.Fatal(err)
		}
		chainID = chain.ID
	}
	run := domain.Run{ID: sourceRunID, TaskID: source.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&run, ""); err != nil {
		t.Fatal(err)
	}
	claimed, err := st.ClaimCompletionDeliveries(1)
	if err != nil || len(claimed) != 1 || claimed[0].ChainID != chainID {
		t.Fatalf("claimed=%+v err=%v", claimed, err)
	}
	return claimed[0]
}

func TestOverlap_CompletionQueuePreservesCorrelationAndResolvesExtra(t *testing.T) {
	st, _ := store.Open(":memory:")
	defer st.Close()
	runner := &blockingRunner{started: make(chan struct{}, 2), release: make(chan struct{})}
	e := newEngine(st, runner)
	source := setupTask(t, st, domain.OverlapQueueOne)
	target := setupTask(t, st, domain.OverlapQueueOne)
	chain := domain.CompletionChain{SourceTaskID: source.ID, TargetTaskID: target.ID, OnOutcome: domain.CompletionOnSuccess}
	if err := st.CreateCompletionChain(&chain); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	e.dispatch(target, now, domain.TriggerManual)
	recv(t, runner.started, "running target")
	first := claimedCompletion(t, st, source, target, chain.ID, "source-one")
	e.dispatchWithOrigin(target, now, dispatchOrigin{trigger: domain.TriggerCompletion, sourceTaskID: source.ID, sourceRunID: first.SourceRunID, deliveryID: first.ID})
	second := claimedCompletion(t, st, source, target, chain.ID, "source-two")
	e.dispatchWithOrigin(target, now, dispatchOrigin{trigger: domain.TriggerCompletion, sourceTaskID: source.ID, sourceRunID: second.SourceRunID, deliveryID: second.ID})
	runner.release <- struct{}{}
	recv(t, runner.started, "queued completion target")
	runner.release <- struct{}{}
	waitFor(t, func() bool {
		deliveries, _ := st.ListCompletionDeliveries()
		states := map[string]domain.DeliveryState{}
		for _, delivery := range deliveries {
			states[delivery.SourceRunID] = delivery.State
		}
		return states["source-one"] == domain.DeliveryCompleted && states["source-two"] == domain.DeliveryResolved
	}, "queued and collapsed completion delivery states")
	runs, _ := st.ListRuns(target.ID, 0)
	found := false
	for _, run := range runs {
		if run.Outcome == domain.OutcomeSuccess && run.Trigger == domain.TriggerCompletion && run.SourceRunID == "source-one" {
			found = true
		}
	}
	if !found {
		t.Fatalf("queued completion correlation missing from %+v", runs)
	}
}

func TestOverlap_CompletionSkipClosesDeliveryWithoutDownstreamWork(t *testing.T) {
	st, _ := store.Open(":memory:")
	defer st.Close()
	runner := &blockingRunner{started: make(chan struct{}, 1), release: make(chan struct{})}
	e := newEngine(st, runner)
	source := setupTask(t, st, domain.OverlapQueueOne)
	target := setupTask(t, st, domain.OverlapSkip)
	delivery := claimedCompletion(t, st, source, target, "", "source-skip")
	now := time.Now().UTC()
	e.dispatch(target, now, domain.TriggerManual)
	recv(t, runner.started, "running skip target")
	e.dispatchWithOrigin(target, now, dispatchOrigin{trigger: domain.TriggerCompletion, sourceTaskID: source.ID, sourceRunID: delivery.SourceRunID, deliveryID: delivery.ID})
	deliveries, err := st.ListCompletionDeliveries()
	if err != nil || len(deliveries) != 1 || deliveries[0].State != domain.DeliveryCompleted {
		t.Fatalf("skipped delivery=%+v err=%v", deliveries, err)
	}
	runs, _ := st.ListRuns(target.ID, 0)
	if len(runs) != 1 || runs[0].Outcome != domain.OutcomeSkipped || runs[0].SourceRunID != "source-skip" {
		t.Fatalf("skipped correlated history=%+v", runs)
	}
	runner.release <- struct{}{}
	waitFor(t, func() bool {
		runs, _ := st.ListRuns(target.ID, 0)
		return len(runs) == 2
	}, "original target completion")
}
