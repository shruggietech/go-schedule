package engine

import (
	"context"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/clock"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
	"github.com/shruggietech/go-schedule/internal/store"
)

// noopRunner returns a success run without doing work, isolating scheduling
// overhead from command execution time.
type noopRunner struct{}

func (noopRunner) Run(_ context.Context, task domain.Task, sf time.Time, trig domain.RunTrigger) domain.Run {
	end := sf
	return domain.Run{TaskID: task.ID, ScheduledFor: sf, EndedAt: &end, Outcome: domain.OutcomeSuccess, Trigger: trig}
}

// BenchmarkDispatch measures the per-run scheduling overhead (dispatch through
// the worker pool to a recorded run), excluding command execution. The
// Performance principle's budget is p99 dispatch latency < 100ms; this overhead
// should be orders of magnitude smaller (microseconds).
func BenchmarkDispatch(b *testing.B) {
	st, err := store.Open(":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer st.Close()

	sch := &domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=MINUTELY;INTERVAL=1"}
	_ = st.CreateSchedule(sch)
	task := domain.Task{
		Name: "bench", Command: "x", Enabled: true, Timezone: "UTC", ScheduleID: sch.ID,
		OverlapPolicy: domain.OverlapAllowConcurrent, CatchupPolicy: domain.CatchupNone, State: domain.TaskActive,
	}
	if err := st.CreateTask(&task); err != nil {
		b.Fatal(err)
	}

	done := make(chan struct{}, 256)
	eng := New(st, clock.NewReal(), noopRunner{}, testLogger(), 8)
	eng.runCtx = context.Background()
	eng.SetOnRun(func(domain.Run) { done <- struct{}{} })

	now := time.Now().UTC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eng.dispatch(task, now, domain.TriggerManual)
		<-done
	}
	b.StopTimer()
}

// BenchmarkNextRun measures the per-task next-run computation, the hot path when
// the engine recomputes schedules for many tasks.
func BenchmarkNextRun(b *testing.B) {
	anchor := time.Date(2026, 6, 19, 8, 0, 0, 0, time.UTC)
	sch := domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=MONTHLY;BYDAY=+3WE;BYHOUR=14;BYMINUTE=0;BYSECOND=0", Anchor: &anchor}
	after := anchor
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = schedule.NextRun(sch, "America/New_York", domain.MissingDateSkip, after)
	}
}
