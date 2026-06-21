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
