package engine

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/clock"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
	"github.com/shruggietech/go-schedule/internal/store"
)

type startupRunner struct {
	mu   sync.Mutex
	runs []domain.Run
	ch   chan domain.Run
}

func (r *startupRunner) Run(_ context.Context, task domain.Task, at time.Time, trigger domain.RunTrigger) domain.Run {
	run := domain.Run{TaskID: task.ID, ScheduledFor: at, Outcome: domain.OutcomeSuccess, Trigger: trigger}
	r.mu.Lock()
	r.runs = append(r.runs, run)
	r.mu.Unlock()
	r.ch <- run
	return run
}

func createStartupTask(t *testing.T, st *store.Store, groupID string, enabled bool) domain.Task {
	t.Helper()
	sch := schedule.NewStartup("@reboot")
	if err := st.CreateSchedule(&sch); err != nil {
		t.Fatal(err)
	}
	task := domain.Task{Name: "startup", Command: "x", GroupID: groupID, Enabled: enabled, Timezone: "UTC", ScheduleID: sch.ID, OverlapPolicy: domain.OverlapQueueOne, CatchupPolicy: domain.CatchupOne, State: domain.TaskActive}
	if err := st.CreateTask(&task); err != nil {
		t.Fatal(err)
	}
	return task
}

func runEngineUntilStartup(t *testing.T, st *store.Store, clk *clock.FakeClock, runner *startupRunner, expect bool) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	e := New(st, clk, runner, testLogger(), 2)
	done := make(chan error, 1)
	go func() { done <- e.Start(ctx) }()
	if expect {
		select {
		case run := <-runner.ch:
			if run.Trigger != domain.TriggerStartup || !run.ScheduledFor.Equal(clk.Now()) {
				t.Fatalf("startup run = %+v", run)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("startup run did not dispatch")
		}
	}
	e.Reload()
	select {
	case run := <-runner.ch:
		t.Fatalf("reload dispatched startup: %+v", run)
	case <-time.After(100 * time.Millisecond):
	}
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("engine did not stop")
	}
}

func TestStartupRunsOncePerEngineStartAndNeverOnReload(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	createStartupTask(t, st, "", true)
	clk := clock.NewFake(time.Date(2026, 8, 30, 12, 0, 0, 0, time.UTC))
	runner := &startupRunner{ch: make(chan domain.Run, 4)}
	runEngineUntilStartup(t, st, clk, runner, true)
	runEngineUntilStartup(t, st, clk, runner, true)
	runs, err := st.ListRuns("", 0)
	if err != nil || len(runs) != 2 {
		t.Fatalf("persisted runs = %d, err=%v", len(runs), err)
	}
}

func TestStartupEligibilityExcludesDisabledTaskAndGroup(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	group := domain.Group{Name: "off", Enabled: false}
	if err := st.CreateGroup(&group); err != nil {
		t.Fatal(err)
	}
	createStartupTask(t, st, "", false)
	createStartupTask(t, st, group.ID, true)
	completed := createStartupTask(t, st, "", true)
	if err := st.SetTaskState(completed.ID, domain.TaskCompleted); err != nil {
		t.Fatal(err)
	}
	unknown := domain.Schedule{Kind: domain.ScheduleEvent, TriggerID: "future_event", HumanSummary: "Future event"}
	if err := st.CreateSchedule(&unknown); err != nil {
		t.Fatal(err)
	}
	unknownTask := domain.Task{Name: "unknown", Command: "x", Enabled: true, Timezone: "UTC", ScheduleID: unknown.ID, State: domain.TaskActive}
	if err := st.CreateTask(&unknownTask); err != nil {
		t.Fatal(err)
	}
	clk := clock.NewFake(time.Now().UTC())
	runner := &startupRunner{ch: make(chan domain.Run, 2)}
	runEngineUntilStartup(t, st, clk, runner, false)
	if runs, _ := st.ListRuns("", 0); len(runs) != 0 {
		t.Fatalf("ineligible startup runs = %+v", runs)
	}
}
