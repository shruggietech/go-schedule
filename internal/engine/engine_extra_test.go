package engine

import (
	"context"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

// instantRunner completes immediately with success.
type instantRunner struct{}

func (instantRunner) Run(_ context.Context, task domain.Task, sf time.Time, trig domain.RunTrigger) domain.Run {
	now := time.Now()
	return domain.Run{TaskID: task.ID, ScheduledFor: sf, EndedAt: &now, Outcome: domain.OutcomeSuccess, Trigger: trig}
}

func waitFor(t *testing.T, cond func() bool, msg string) {
	t.Helper()
	for i := 0; i < 50; i++ {
		if cond() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timeout waiting for %s", msg)
}

// TestEngine_RunNowFiresHooks covers RunNow, SetOnRun, and the run-recording path.
func TestEngine_RunNowFiresHooks(t *testing.T) {
	st, _ := store.Open(":memory:")
	defer st.Close()
	e := newEngine(st, instantRunner{})

	runs := make(chan domain.Run, 1)
	e.SetOnRun(func(r domain.Run) { runs <- r })

	task := setupTask(t, st, domain.OverlapAllowConcurrent)
	if err := e.RunNow(task.ID); err != nil {
		t.Fatalf("RunNow: %v", err)
	}
	select {
	case r := <-runs:
		if r.Trigger != domain.TriggerManual {
			t.Fatalf("trigger = %q, want manual", r.Trigger)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("onRun hook did not fire")
	}

	// Unknown task is an error.
	if err := e.RunNow("missing"); err == nil {
		t.Fatal("RunNow on missing task should error")
	}
}

// TestEngine_SetOnAlertFires covers SetOnAlert via an overlap alert.
func TestEngine_SetOnAlertFires(t *testing.T) {
	st, _ := store.Open(":memory:")
	defer st.Close()
	r := &blockingRunner{started: make(chan struct{}, 1), release: make(chan struct{})}
	e := newEngine(st, r)

	alerts := make(chan domain.Alert, 1)
	e.SetOnAlert(func(a domain.Alert) { alerts <- a })

	task := setupTask(t, st, domain.OverlapQueueOne)
	now := time.Now().UTC()
	e.dispatch(task, now, domain.TriggerSchedule)
	recv(t, r.started, "first run start")
	e.dispatch(task, now.Add(time.Minute), domain.TriggerSchedule) // queues -> overlap alert

	select {
	case a := <-alerts:
		if a.Kind != domain.AlertOverlapQueued {
			t.Fatalf("alert kind = %q", a.Kind)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("onAlert hook did not fire")
	}
	r.release <- struct{}{}
	recv(t, r.started, "queued run start")
	r.release <- struct{}{}
}

// TestEngine_CompleteOneOff covers the one-off completion path.
func TestEngine_CompleteOneOff(t *testing.T) {
	st, _ := store.Open(":memory:")
	defer st.Close()
	e := newEngine(st, instantRunner{})

	at := time.Now().Add(time.Hour).UTC()
	sch := &domain.Schedule{Kind: domain.ScheduleOneOff, RunAt: &at}
	if err := st.CreateSchedule(sch); err != nil {
		t.Fatal(err)
	}
	task := &domain.Task{
		Name: "once", Command: "x", Enabled: true, Timezone: "UTC", ScheduleID: sch.ID,
		OverlapPolicy: domain.OverlapAllowConcurrent, CatchupPolicy: domain.CatchupNone, State: domain.TaskActive,
	}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	// Seed engine state as if the task were scheduled, then complete it.
	e.tasks[task.ID] = taskCtx{task: *task, sch: *sch}
	e.next[task.ID] = at
	e.completeOneOff(task.ID)

	waitFor(t, func() bool {
		got, err := st.GetTask(task.ID)
		return err == nil && got.State == domain.TaskCompleted
	}, "one-off marked completed")
	if _, ok := e.next[task.ID]; ok {
		t.Fatal("completed one-off should be removed from scheduling")
	}
}
