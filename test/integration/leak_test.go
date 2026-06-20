package integration

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/shruggietech/go-scheduler/internal/clock"
	"github.com/shruggietech/go-scheduler/internal/domain"
	"github.com/shruggietech/go-scheduler/internal/engine"
	"github.com/shruggietech/go-scheduler/internal/store"
)

// TestEngine_NoGoroutineLeak runs many task executions and then shuts the engine
// down, asserting goroutines drain back to a baseline. The engine's run
// goroutines are tracked by a WaitGroup and drained on context cancellation, so
// after Start returns there should be no lingering goroutines.
func TestEngine_NoGoroutineLeak(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	base := time.Date(2026, 6, 19, 8, 0, 0, 0, time.UTC)
	sch := &domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=SECONDLY;INTERVAL=1", Anchor: &base}
	if err := st.CreateSchedule(sch); err != nil {
		t.Fatal(err)
	}
	task := &domain.Task{
		Name: "loop", Command: "x", Enabled: true, Timezone: "UTC", ScheduleID: sch.ID,
		OverlapPolicy: domain.OverlapAllowConcurrent, CatchupPolicy: domain.CatchupNone, State: domain.TaskActive,
	}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}

	// Settle and snapshot a baseline goroutine count.
	runtime.GC()
	time.Sleep(50 * time.Millisecond)
	baseline := runtime.NumGoroutine()

	fc := clock.NewFake(base)
	ran := make(chan domain.Run, 4096)
	eng := engine.New(st, fc, recordingRunner{}, quietLogger(), 8)
	eng.SetOnRun(func(r domain.Run) {
		select {
		case ran <- r:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { _ = eng.Start(ctx); close(done) }()

	// Drive ~500 executions through the worker pool.
	for i := 0; i < 500; i++ {
		waitWaiter(t, fc)
		fc.Advance(time.Second)
		<-ran
	}

	cancel()
	<-done // Start drains in-flight runs before returning

	// Allow scheduler bookkeeping to settle, then compare goroutine counts.
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	after := runtime.NumGoroutine()

	if after > baseline+3 { // small tolerance for runtime/test noise
		t.Fatalf("possible goroutine leak: baseline=%d after=%d (ran 500 tasks)", baseline, after)
	}
}
