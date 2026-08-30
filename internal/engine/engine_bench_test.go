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

// BenchmarkNextRunNearestWeekday measures the bounded monthly adjustment path
// added for nW schedules, including timezone and missing-date handling.
func BenchmarkNextRunNearestWeekday(b *testing.B) {
	anchor := time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)
	sch := domain.Schedule{
		Kind: domain.ScheduleRecurring, RRULE: "FREQ=MONTHLY;BYMONTHDAY=31;BYHOUR=9;BYMINUTE=0;BYSECOND=0", Anchor: &anchor,
		CalendarAdjustment: domain.CalendarAdjustmentNearestWeekday,
	}
	after := time.Date(2036, 4, 1, 0, 0, 0, 0, time.UTC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = schedule.NextRun(sch, "America/New_York", domain.MissingDateNextValid, after)
	}
}

func BenchmarkNextRunCompositeCron(b *testing.B) {
	anchor := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for name, rrule := range map[string]string{
		"broad":  "FREQ=DAILY;INTERVAL=1;BYHOUR=9,10,11,12,13,14,15,16,17;BYMINUTE=0,10,20,30,40,50;BYSECOND=0;BYDAY=MO,WE,FR",
		"sparse": "FREQ=DAILY;INTERVAL=1;BYHOUR=0;BYMINUTE=0;BYSECOND=0;BYMONTH=2;BYMONTHDAY=29",
	} {
		b.Run(name, func(b *testing.B) {
			sch := domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: rrule, Anchor: &anchor}
			after := time.Date(2036, 3, 1, 0, 0, 0, 0, time.UTC)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, _, err := schedule.NextRun(sch, "America/New_York", domain.MissingDateSkip, after); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkNextRunDSTPolicy(b *testing.B) {
	anchor := time.Date(2026, time.October, 31, 13, 0, 0, 0, time.UTC)
	sch, err := schedule.Parse("every 6 hours starting at 09:00", "America/New_York", anchor)
	if err != nil {
		b.Fatal(err)
	}
	after := time.Date(2026, time.November, 1, 4, 0, 0, 0, time.UTC)
	for name, policy := range map[string]domain.SchedulePolicy{
		"wall_clock": {TimeBasis: domain.TimeBasisWallClock, DSTOverlap: domain.DSTOverlapBoth},
		"elapsed":    {TimeBasis: domain.TimeBasisElapsed},
		"utc":        {TimeBasis: domain.TimeBasisUTC},
	} {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, _, err := schedule.NextRunWithPolicy(sch, "America/New_York", policy, after); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
