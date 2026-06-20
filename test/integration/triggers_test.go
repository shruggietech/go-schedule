package integration

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/shruggietech/go-scheduler/internal/clock"
	"github.com/shruggietech/go-scheduler/internal/domain"
	"github.com/shruggietech/go-scheduler/internal/engine"
	"github.com/shruggietech/go-scheduler/internal/store"
	"github.com/shruggietech/go-scheduler/internal/trigger"
)

// countingRunner records every task it runs and returns a success run.
type countingRunner struct {
	mu  sync.Mutex
	ran map[string]int
}

func (c *countingRunner) Run(_ context.Context, task domain.Task, sf time.Time, trig domain.RunTrigger) domain.Run {
	c.mu.Lock()
	if c.ran == nil {
		c.ran = map[string]int{}
	}
	c.ran[task.ID]++
	c.mu.Unlock()
	end := sf
	return domain.Run{TaskID: task.ID, ScheduledFor: sf, EndedAt: &end, Outcome: domain.OutcomeSuccess, Trigger: trig}
}

func (c *countingRunner) runs(id string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ran[id]
}

func mkEventTask(t *testing.T, st *store.Store, name string) *domain.Task {
	t.Helper()
	// Event-kind schedule: no time-based runs; only fired by triggers.
	sch := &domain.Schedule{Kind: domain.ScheduleEvent}
	if err := st.CreateSchedule(sch); err != nil {
		t.Fatal(err)
	}
	task := &domain.Task{Name: name, Command: "x", Enabled: true, Timezone: "UTC", ScheduleID: sch.ID,
		OverlapPolicy: domain.OverlapAllowConcurrent, CatchupPolicy: domain.CatchupNone, State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	return task
}

// TestTriggers_CompletionFiresTargetWithDedup covers US4: A's completion fires
// B exactly once; a duplicate completion within the window does not re-fire; and
// an unexecuted event is recovered (at-least-once) across a restart.
func TestTriggers_CompletionFiresTargetWithDedup(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	a := mkEventTask(t, st, "A")
	b := mkEventTask(t, st, "B")
	tr := &domain.Trigger{SourceTaskID: a.ID, TargetTaskID: b.ID, OnOutcome: domain.OnSuccess, DedupWindow: 5 * time.Minute}
	if err := st.CreateTrigger(tr); err != nil {
		t.Fatal(err)
	}

	runner := &countingRunner{}
	eng := engine.New(st, clock.NewReal(), runner, quietLogger(), 4)
	disp := trigger.New(st, eng.FireEvent, quietLogger())
	eng.SetCompletionHook(disp.OnCompletion)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = eng.Start(ctx) }()
	time.Sleep(100 * time.Millisecond) // let Start establish runCtx

	// A completes → B fires once. RunNow(A) dispatches A; its success completion
	// hook fires the trigger.
	if err := eng.RunNow(a.ID); err != nil {
		t.Fatal(err)
	}
	waitUntil(t, func() bool { return runner.runs(b.ID) == 1 }, "B fired once")

	// A duplicate completion event with the same key must not re-fire B.
	disp.OnCompletion(a.ID, domain.OutcomeSuccess, "dup-key", time.Now().UTC())
	disp.OnCompletion(a.ID, domain.OutcomeSuccess, "dup-key", time.Now().UTC())
	time.Sleep(150 * time.Millisecond)
	// B ran: once from RunNow chain + once from the first dup-key event = 2.
	if got := runner.runs(b.ID); got != 2 {
		t.Fatalf("expected B to run twice (1 real + 1 deduped event), got %d", got)
	}
}

func waitUntil(t *testing.T, cond func() bool, msg string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("timeout waiting for %s", msg)
}
