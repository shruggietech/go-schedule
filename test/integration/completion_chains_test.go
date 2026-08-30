package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/clock"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/engine"
	"github.com/shruggietech/go-schedule/internal/store"
)

func integrationChainTask(t *testing.T, st *store.Store, name string, anchor time.Time) domain.Task {
	t.Helper()
	schedule := domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=HOURLY;INTERVAL=1", Anchor: &anchor}
	if err := st.CreateSchedule(&schedule); err != nil {
		t.Fatal(err)
	}
	task := domain.Task{
		Name: name, Command: "echo", Enabled: true, Timezone: "UTC", ScheduleID: schedule.ID,
		OverlapPolicy: domain.OverlapQueueOne, CatchupPolicy: domain.CatchupNone, State: domain.TaskActive,
	}
	if err := st.CreateTask(&task); err != nil {
		t.Fatal(err)
	}
	return task
}

func integrationChain(t *testing.T, st *store.Store, source, target domain.Task) {
	t.Helper()
	chain := domain.CompletionChain{SourceTaskID: source.ID, TargetTaskID: target.ID, OnOutcome: domain.CompletionOnSuccess}
	if err := st.CreateCompletionChain(&chain); err != nil {
		t.Fatal(err)
	}
}

func TestCompletionChainsCascadeAndCoexistWithSchedule(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	base := time.Date(2026, 8, 30, 12, 0, 0, 0, time.UTC)
	a := integrationChainTask(t, st, "a", base)
	b := integrationChainTask(t, st, "b", base)
	c := integrationChainTask(t, st, "c", base)
	integrationChain(t, st, a, b)
	integrationChain(t, st, b, c)

	fakeClock := clock.NewFake(base)
	runs := make(chan domain.Run, 12)
	eng := engine.New(st, fakeClock, recordingRunner{}, quietLogger(), 4)
	eng.SetOnRun(func(run domain.Run) { runs <- run })
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- eng.Start(ctx) }()
	select {
	case <-eng.Ready():
	case <-time.After(2 * time.Second):
		t.Fatal("engine did not become ready")
	}
	if err := eng.RunNow(a.ID); err != nil {
		t.Fatal(err)
	}
	seen := map[string]domain.Run{}
	for len(seen) < 3 {
		run := waitSignal(t, runs, "A-to-B-to-C completion cascade")
		seen[run.TaskID] = run
	}
	if seen[b.ID].Trigger != domain.TriggerCompletion || seen[b.ID].SourceTaskID != a.ID || seen[c.ID].SourceTaskID != b.ID {
		t.Fatalf("cascade correlation = %+v", seen)
	}

	// B keeps its regular schedule after being used as a completion target.
	waitWaiter(t, fakeClock)
	fakeClock.Advance(time.Hour)
	for {
		run := waitSignal(t, runs, "scheduled target run")
		if run.TaskID == b.ID && run.Trigger == domain.TriggerSchedule {
			break
		}
	}
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("engine did not drain")
	}
}

func TestCompletionChainsHundredWayFanOutUsesBoundedWorkers(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	base := time.Date(2026, 8, 30, 12, 0, 0, 0, time.UTC)
	source := integrationChainTask(t, st, "source", base)
	for i := 0; i < 100; i++ {
		target := integrationChainTask(t, st, fmt.Sprintf("target-%03d", i), base)
		integrationChain(t, st, source, target)
	}
	fakeClock := clock.NewFake(base)
	runs := make(chan domain.Run, 120)
	eng := engine.New(st, fakeClock, recordingRunner{}, quietLogger(), 8)
	eng.SetOnRun(func(run domain.Run) { runs <- run })
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- eng.Start(ctx) }()
	<-eng.Ready()
	if err := eng.RunNow(source.ID); err != nil {
		t.Fatal(err)
	}
	completionCount := 0
	deadline := time.After(5 * time.Second)
	for completionCount < 100 {
		select {
		case run := <-runs:
			if run.Trigger == domain.TriggerCompletion {
				completionCount++
			}
		case <-deadline:
			t.Fatalf("completion fan-out count=%d, want 100", completionCount)
		}
	}
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("engine did not drain fan-out")
	}
	deliveries, err := st.ListCompletionDeliveries()
	if err != nil || len(deliveries) != 100 {
		t.Fatalf("fan-out deliveries=%d err=%v", len(deliveries), err)
	}
	for _, delivery := range deliveries {
		if delivery.State != domain.DeliveryCompleted {
			t.Fatalf("incomplete fan-out delivery: %+v", delivery)
		}
	}
}
