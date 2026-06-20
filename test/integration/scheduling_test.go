package integration

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/shruggietech/go-scheduler/internal/clock"
	"github.com/shruggietech/go-scheduler/internal/domain"
	"github.com/shruggietech/go-scheduler/internal/engine"
	"github.com/shruggietech/go-scheduler/internal/store"
)

// recordingRunner returns a success Run without doing real work, so the engine's
// scheduling behavior can be tested deterministically under a fake clock.
type recordingRunner struct{}

func (recordingRunner) Run(_ context.Context, task domain.Task, sf time.Time, trig domain.RunTrigger) domain.Run {
	end := sf
	return domain.Run{TaskID: task.ID, ScheduledFor: sf, EndedAt: &end, Outcome: domain.OutcomeSuccess, Trigger: trig}
}

func quietLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func waitWaiter(t *testing.T, c *clock.FakeClock) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if c.Waiters() >= 1 {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatal("engine never armed its timer")
}

func waitSignal(t *testing.T, ch <-chan domain.Run, msg string) domain.Run {
	t.Helper()
	select {
	case r := <-ch:
		return r
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for %s", msg)
		return domain.Run{}
	}
}

// TestEngine_SchedulesRunsAndResumesAfterRestart covers US1's independent test:
// a recurring task fires at the expected times under a fake clock, and after a
// simulated daemon restart the schedule resumes from the same persisted state.
func TestEngine_SchedulesRunsAndResumesAfterRestart(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	base := time.Date(2026, 6, 19, 8, 0, 0, 0, time.UTC)

	// Recurring every 2 hours, anchored at base.
	sch := &domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=HOURLY;INTERVAL=2", Anchor: &base}
	if err := st.CreateSchedule(sch); err != nil {
		t.Fatal(err)
	}
	task := &domain.Task{
		Name: "tick", Command: "x", Enabled: true, Timezone: "UTC", ScheduleID: sch.ID,
		OverlapPolicy: domain.OverlapQueueOne, CatchupPolicy: domain.CatchupNone, State: domain.TaskActive,
	}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}

	// ---- First engine instance ----
	fc := clock.NewFake(base)
	ran := make(chan domain.Run, 8)
	eng := engine.New(st, fc, recordingRunner{}, quietLogger(), 4)
	eng.SetOnRun(func(r domain.Run) { ran <- r })

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { _ = eng.Start(ctx); close(done) }()

	waitWaiter(t, fc)
	fc.Advance(2 * time.Hour) // -> 10:00, first run due
	r1 := waitSignal(t, ran, "first run")
	if want := base.Add(2 * time.Hour); !r1.ScheduledFor.Equal(want) {
		t.Fatalf("first run scheduled_for = %v, want %v", r1.ScheduledFor, want)
	}

	waitWaiter(t, fc)
	fc.Advance(2 * time.Hour) // -> 12:00, second run
	r2 := waitSignal(t, ran, "second run")
	if want := base.Add(4 * time.Hour); !r2.ScheduledFor.Equal(want) {
		t.Fatalf("second run scheduled_for = %v, want %v", r2.ScheduledFor, want)
	}

	cancel()
	<-done

	// ---- Simulated restart: new engine, same store ----
	resumeAt := base.Add(4 * time.Hour)
	fc2 := clock.NewFake(resumeAt)
	ran2 := make(chan domain.Run, 8)
	eng2 := engine.New(st, fc2, recordingRunner{}, quietLogger(), 4)
	eng2.SetOnRun(func(r domain.Run) { ran2 <- r })
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	go func() { _ = eng2.Start(ctx2) }()

	waitWaiter(t, fc2)
	fc2.Advance(2 * time.Hour) // -> 14:00, schedule resumes
	r3 := waitSignal(t, ran2, "run after restart")
	if want := resumeAt.Add(2 * time.Hour); !r3.ScheduledFor.Equal(want) {
		t.Fatalf("post-restart run scheduled_for = %v, want %v", r3.ScheduledFor, want)
	}

	// History persisted across the restart.
	runs, err := st.ListRuns(task.ID, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) < 3 {
		t.Fatalf("expected >=3 recorded runs across restart, got %d", len(runs))
	}
}
