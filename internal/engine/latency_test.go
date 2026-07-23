package engine

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/clock"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

// timingRunner records the dispatch latency of each run — the interval from the
// run's scheduled time to the moment its execution starts — and does no other
// work, so the measurement reflects scheduling overhead alone (queue, goroutine
// hand-off, and semaphore acquisition), not command execution. It mirrors the
// no-op runner used by BenchmarkDispatch but timestamps the start.
type timingRunner struct {
	latencies chan<- time.Duration
}

func (r timingRunner) Run(_ context.Context, task domain.Task, sf time.Time, trig domain.RunTrigger) domain.Run {
	start := time.Now()
	r.latencies <- start.Sub(sf)
	end := start
	return domain.Run{
		TaskID: task.ID, ScheduledFor: sf, StartedAt: &start, EndedAt: &end,
		Outcome: domain.OutcomeSuccess, Trigger: trig,
	}
}

// TestDispatchLatencyP99 enforces the constitution's Performance budget
// (Principle IV): the p99 of dispatch latency must stay under
// DispatchLatencyBudget. It dispatches a fixed number of runs serially — waiting
// for each to complete before dispatching the next, so every sample measures
// pure dispatch overhead under nominal load rather than self-induced queue
// contention — collects the per-dispatch latencies, and asserts the p99 against
// the budget. The assertion carries several orders of magnitude of headroom (the
// real overhead is microseconds against a 100ms ceiling), so it is stable on
// loaded CI hardware and does not depend on real sleeps for correctness.
func TestDispatchLatencyP99(t *testing.T) {
	const n = 2000

	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	sch := &domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=MINUTELY;INTERVAL=1"}
	if err := st.CreateSchedule(sch); err != nil {
		t.Fatal(err)
	}
	task := domain.Task{
		Name: "latency", Command: "x", Enabled: true, Timezone: "UTC", ScheduleID: sch.ID,
		OverlapPolicy: domain.OverlapAllowConcurrent, CatchupPolicy: domain.CatchupNone, State: domain.TaskActive,
	}
	if err := st.CreateTask(&task); err != nil {
		t.Fatal(err)
	}

	latencies := make(chan time.Duration, n)
	done := make(chan struct{}, n)
	eng := New(st, clock.NewReal(), timingRunner{latencies: latencies}, testLogger(), 8)
	eng.runCtx = context.Background()
	eng.SetOnRun(func(domain.Run) { done <- struct{}{} })

	for i := 0; i < n; i++ {
		eng.dispatch(task, time.Now(), domain.TriggerManual)
		<-done
	}

	samples := make([]time.Duration, n)
	for i := 0; i < n; i++ {
		samples[i] = <-latencies
	}
	sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })

	p99 := samples[int(0.99*float64(n))]
	t.Logf("dispatch latency over %d samples: p50=%v p99=%v max=%v (budget %v)",
		n, samples[n/2], p99, samples[n-1], DispatchLatencyBudget)

	if p99 >= DispatchLatencyBudget {
		t.Fatalf("p99 dispatch latency %v exceeds budget %v", p99, DispatchLatencyBudget)
	}
}
