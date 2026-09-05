// Package engine is the scheduling core: a single timer-driven loop computes
// the next run for each active task, dispatches due runs through a bounded
// worker pool, and records history. Time is read through an injected Clock so
// behavior is deterministic under test. Overlap handling lives in overlap.go.
package engine

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/shruggietech/go-schedule/internal/catchup"
	"github.com/shruggietech/go-schedule/internal/clock"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
	"github.com/shruggietech/go-schedule/internal/store"
	watchruntime "github.com/shruggietech/go-schedule/internal/watcher"
)

// Runner executes a task and returns its Run record. The executor implements it;
// tests inject a fake.
type Runner interface {
	Run(ctx context.Context, task domain.Task, scheduledFor time.Time, trigger domain.RunTrigger) domain.Run
}

// taskCtx caches a task with its schedule for next-run computation.
type taskCtx struct {
	task domain.Task
	sch  domain.Schedule
}

// DispatchLatencyBudget is the documented ceiling for dispatch latency, the
// interval from a run's scheduled time to the start of its execution. Per the
// constitution's Performance principle (IV), the p99 of this latency must stay
// under this budget "under nominal load", and the budget lives next to the
// dispatch code it governs. TestDispatchLatencyP99 asserts against this value.
const DispatchLatencyBudget = 100 * time.Millisecond

// Engine schedules and dispatches task runs.
type Engine struct {
	store  *store.Store
	clk    clock.Clock
	runner Runner
	log    *slog.Logger
	sem    chan struct{} // bounded worker pool

	mu      sync.Mutex
	tasks   map[string]taskCtx   // active tasks by ID
	next    map[string]time.Time // next scheduled run (UTC) by task ID
	running map[string]bool
	queued  map[string]pendingRun // one queued pending run, by task ID

	reload    chan struct{}
	ready     chan struct{}
	readyOnce sync.Once
	runCtx    context.Context
	runWG     sync.WaitGroup // tracks in-flight runs for graceful drain
	onRun     func(domain.Run)
	onAlert   func(domain.Alert)
	watchers  *watchruntime.Manager
}

// New constructs an Engine. workers bounds concurrent task executions.
func New(st *store.Store, clk clock.Clock, runner Runner, log *slog.Logger, workers int) *Engine {
	if workers <= 0 {
		workers = 1
	}
	engine := &Engine{
		store:   st,
		clk:     clk,
		runner:  runner,
		log:     log,
		sem:     make(chan struct{}, workers),
		tasks:   map[string]taskCtx{},
		next:    map[string]time.Time{},
		running: map[string]bool{},
		queued:  map[string]pendingRun{},
		reload:  make(chan struct{}, 1),
		ready:   make(chan struct{}),
	}
	engine.watchers = watchruntime.New(st, engine, clk, log)
	return engine
}

// SetOnRun registers a callback invoked after each run is recorded (used for
// alerts/event streaming and for test synchronization).
func (e *Engine) SetOnRun(f func(domain.Run)) { e.onRun = f }

// SetOnAlert registers a callback invoked after each alert is raised (used to
// stream alerts to GUI clients).
func (e *Engine) SetOnAlert(f func(domain.Alert)) { e.onAlert = f }

// SetOnWatcherHealth registers a callback for watcher runtime health transitions.
func (e *Engine) SetOnWatcherHealth(f func(domain.FilesystemWatcher, domain.WatcherHealth)) {
	e.watchers.SetHealthReporter(f)
}

// ReloadWatchers asks the watcher runtime to rebuild native registrations.
func (e *Engine) ReloadWatchers() { e.watchers.Reload() }

// WatcherHealth returns the current runtime health for a watcher.
func (e *Engine) WatcherHealth(id string) domain.WatcherHealth { return e.watchers.Health(id) }

// Reload asks the loop to recompute schedules from the store (call after tasks
// change). Non-blocking and coalesced.
func (e *Engine) Reload() {
	select {
	case e.reload <- struct{}{}:
	default:
	}
}

// Ready is closed after Start freezes and dispatches the initial startup-task
// snapshot. Mutation-serving callers can wait on it to ensure newly created or
// enabled startup tasks remain deferred until the next daemon lifecycle.
func (e *Engine) Ready() <-chan struct{} { return e.ready }

// Start runs the scheduling loop until ctx is cancelled, then drains in-flight
// runs. It blocks, must be called once, and is normally run in a goroutine.
func (e *Engine) Start(ctx context.Context) error {
	e.runCtx = ctx
	watchDone := make(chan error, 1)
	go func() { watchDone <- e.watchers.Run(ctx) }()
	select {
	case <-e.watchers.Ready():
	case err := <-watchDone:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
	e.recompute(e.clk.Now())
	if recovered, err := e.store.RecoverCompletionDeliveries(); err != nil {
		e.log.Error("engine: recover completion deliveries", "err", err)
	} else if recovered > 0 {
		e.log.Warn("engine: recovered interrupted completion deliveries", "count", recovered)
	}
	e.drainCompletionDeliveries()
	e.runStartup(e.clk.Now())
	e.runCatchup(e.clk.Now())
	e.readyOnce.Do(func() { close(e.ready) })
	for {
		d, has := e.untilNext(e.clk.Now())
		var wake <-chan time.Time
		var timer *clock.Timer
		if has {
			timer = e.clk.NewTimer(d)
			wake = timer.C
		}
		select {
		case <-ctx.Done():
			if timer != nil {
				timer.Stop()
			}
			e.runWG.Wait()
			<-watchDone
			return ctx.Err()
		case <-e.reload:
			if timer != nil {
				timer.Stop()
			}
			e.recompute(e.clk.Now())
		case now := <-wake:
			e.runDue(now)
		}
	}
}

// runStartup dispatches the eligible startup-event snapshot loaded by the
// initial recompute. Start calls it exactly once; Reload never does.
func (e *Engine) runStartup(now time.Time) {
	e.mu.Lock()
	tasks := make([]domain.Task, 0)
	for _, tc := range e.tasks {
		if schedule.IsStartup(tc.sch) {
			tasks = append(tasks, tc.task)
		}
	}
	e.mu.Unlock()
	for _, task := range tasks {
		e.dispatch(task, now, domain.TriggerStartup)
	}
}

// recompute rebuilds the active task set and their next run times from the store.
func (e *Engine) recompute(now time.Time) {
	tasks, err := e.store.ListTasks("", string(domain.TaskActive))
	if err != nil {
		e.log.Error("engine: list tasks", "err", err)
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.tasks = map[string]taskCtx{}
	newNext := map[string]time.Time{}
	for _, task := range tasks {
		if !task.Enabled {
			continue
		}
		if task.ScheduleID == "" {
			continue
		}
		// A task is ineligible if any ancestor group is disabled (cascade).
		if ok, err := e.store.GroupChainEnabled(task.GroupID); err == nil && !ok {
			continue
		}
		sch, err := e.store.GetSchedule(task.ScheduleID)
		if err != nil {
			e.log.Error("engine: get schedule", "task", task.ID, "err", err)
			continue
		}
		e.tasks[task.ID] = taskCtx{task: task, sch: sch}
		if n, ok, err := schedule.NextRunWithPolicy(sch, task.Timezone, task.SchedulePolicy(), now); err == nil && ok {
			newNext[task.ID] = n
		}
	}
	e.next = newNext
}

// runCatchup performs one catch-up run per eligible task that missed scheduled
// runs during downtime. The catch-up run is recorded at `now` (so a subsequent
// restart does not re-trigger it) and honors the task's overlap policy via
// dispatch. Normal scheduling (computed in recompute) resumes afterward.
func (e *Engine) runCatchup(now time.Time) {
	e.mu.Lock()
	tasks := make([]taskCtx, 0, len(e.tasks))
	for _, tc := range e.tasks {
		tasks = append(tasks, tc)
	}
	e.mu.Unlock()

	for _, tc := range tasks {
		runs, err := e.store.ListRuns(tc.task.ID, 1)
		if err != nil || len(runs) == 0 {
			continue // never run → nothing to catch up
		}
		dec, err := catchup.EvaluateWithPolicy(tc.sch, tc.task.Timezone, runs[0].ScheduledFor, true, tc.task.CatchupPolicy, tc.task.SchedulePolicy(), now)
		if err != nil {
			e.log.Error("engine: catchup evaluate", "task", tc.task.ID, "err", err)
			continue
		}
		if !dec.ShouldCatchUp {
			continue
		}
		e.log.Warn("missed run(s) during downtime; performing one catch-up",
			"task", tc.task.ID, "name", tc.task.Name, "first_missed", dec.FirstMissed)
		e.raiseAlert(tc.task.ID, domain.SeverityWarning, domain.AlertMissedRun,
			"missed run(s) during downtime; running one catch-up")
		e.dispatch(tc.task, now, domain.TriggerCatchup)
	}
}

// untilNext returns the duration until the earliest scheduled run, and whether
// any run is scheduled.
func (e *Engine) untilNext(now time.Time) (time.Duration, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	var earliest time.Time
	has := false
	for _, t := range e.next {
		if !has || t.Before(earliest) {
			earliest, has = t, true
		}
	}
	if !has {
		return 0, false
	}
	if d := earliest.Sub(now); d > 0 {
		return d, true
	}
	return 0, true
}

// runDue dispatches every task whose next run is at or before now and advances
// its schedule.
func (e *Engine) runDue(now time.Time) {
	e.mu.Lock()
	due := make([]string, 0)
	for id, t := range e.next {
		if !t.After(now) {
			due = append(due, id)
		}
	}
	e.mu.Unlock()

	for _, id := range due {
		e.mu.Lock()
		tc, ok := e.tasks[id]
		scheduledFor := e.next[id]
		e.mu.Unlock()
		if !ok {
			continue
		}

		e.dispatch(tc.task, scheduledFor, domain.TriggerSchedule)

		// Advance (or retire one-off) using the schedule.
		if tc.sch.Kind == domain.ScheduleOneOff {
			e.completeOneOff(id)
			continue
		}
		if n, ok, err := schedule.NextRunWithPolicy(tc.sch, tc.task.Timezone, tc.task.SchedulePolicy(), scheduledFor); err == nil && ok {
			e.mu.Lock()
			e.next[id] = n
			e.mu.Unlock()
		} else {
			e.mu.Lock()
			delete(e.next, id)
			e.mu.Unlock()
		}
	}
}

// completeOneOff marks a one-off task completed and removes it from scheduling.
func (e *Engine) completeOneOff(taskID string) {
	e.mu.Lock()
	delete(e.next, taskID)
	delete(e.tasks, taskID)
	e.mu.Unlock()
	if err := e.store.SetTaskState(taskID, domain.TaskCompleted); err != nil {
		e.log.Error("engine: complete one-off", "task", taskID, "err", err)
	}
}

// launch runs a task through the worker pool and records the result.
func (e *Engine) launch(task domain.Task, scheduledFor time.Time, origin dispatchOrigin) {
	e.runWG.Add(1)
	go func() {
		defer e.runWG.Done()
		e.sem <- struct{}{}
		defer func() { <-e.sem }()

		run := e.runner.Run(e.runCtx, task, scheduledFor, origin.trigger)
		run.Trigger = origin.trigger
		run.SourceTaskID = origin.sourceTaskID
		run.SourceRunID = origin.sourceRunID
		run.SourceTriggerID = origin.sourceTriggerID
		run.SourceWatcherID = origin.sourceWatcherID
		e.recordRun(run, origin.deliveryID)
		e.finish(task)
	}()
}

// recordRun persists a run, raises a failure alert when needed, and notifies
// the onRun callback.
func (e *Engine) recordRun(run domain.Run, incomingDeliveryID string) {
	if err := e.store.RecordRunAndCreateDeliveries(&run, incomingDeliveryID); err != nil {
		e.log.Error("engine: record run", "task", run.TaskID, "err", err)
		return
	}
	e.log.Info("engine: recorded run", "task", run.TaskID, "run", run.ID, "trigger", run.Trigger,
		"source_task", run.SourceTaskID, "source_run", run.SourceRunID, "source_trigger", run.SourceTriggerID, "source_watcher", run.SourceWatcherID, "delivery", incomingDeliveryID)
	if run.Outcome == domain.OutcomeFailure {
		e.raiseRunAlert(run.TaskID, run.ID, domain.SeverityError, domain.AlertRunFailed, "task run failed")
	}
	if e.onRun != nil {
		e.onRun(run)
	}
	if e.runCtx.Err() == nil {
		e.drainCompletionDeliveries()
	}
}

// drainCompletionDeliveries claims durable pending work in bounded batches and
// dispatches each target through the same eligibility and overlap machinery as
// every other run origin. It is event-driven by startup and committed runs, so
// completion chaining adds no polling loop or long-lived goroutine.
func (e *Engine) drainCompletionDeliveries() {
	for {
		if e.runCtx.Err() != nil {
			return
		}
		deliveries, err := e.store.ClaimCompletionDeliveries(100)
		if err != nil {
			e.log.Error("engine: claim completion deliveries", "err", err)
			return
		}
		for _, delivery := range deliveries {
			if e.runCtx.Err() != nil {
				return
			}
			if _, err := e.store.GetCompletionChain(delivery.ChainID); err != nil {
				e.resolveCompletionDelivery(delivery.ID, "completion chain no longer exists")
				continue
			}
			task, err := e.store.GetTask(delivery.TargetTaskID)
			if err != nil {
				e.resolveCompletionDelivery(delivery.ID, "target task no longer exists")
				continue
			}
			if !task.Enabled || task.State != domain.TaskActive {
				e.resolveCompletionDelivery(delivery.ID, "target task is disabled or inactive")
				continue
			}
			if enabled, err := e.store.GroupChainEnabled(task.GroupID); err != nil {
				e.resolveCompletionDelivery(delivery.ID, "target group eligibility could not be read")
				continue
			} else if !enabled {
				e.resolveCompletionDelivery(delivery.ID, "target group chain is disabled")
				continue
			}
			e.log.Info("engine: dispatch completion delivery", "delivery", delivery.ID, "chain", delivery.ChainID,
				"source_task", delivery.SourceTaskID, "source_run", delivery.SourceRunID, "target_task", delivery.TargetTaskID,
				"attempt", delivery.Attempts)
			e.dispatchWithOrigin(task, e.clk.Now(), dispatchOrigin{
				trigger:      domain.TriggerCompletion,
				sourceTaskID: delivery.SourceTaskID,
				sourceRunID:  delivery.SourceRunID,
				deliveryID:   delivery.ID,
			})
		}
		if len(deliveries) < 100 {
			return
		}
	}
}

func (e *Engine) resolveCompletionDelivery(id, reason string) {
	if err := e.store.ResolveCompletionDelivery(id, reason); err != nil {
		e.log.Error("engine: resolve completion delivery", "delivery", id, "reason", reason, "err", err)
		return
	}
	e.log.Warn("engine: completion delivery resolved without execution", "delivery", id, "reason", reason)
}

// finish marks a task no longer running and dispatches any queued pending run.
func (e *Engine) finish(task domain.Task) {
	e.mu.Lock()
	e.running[task.ID] = false
	pending, queued := e.queued[task.ID]
	delete(e.queued, task.ID)
	e.mu.Unlock()

	if queued {
		e.mu.Lock()
		e.running[task.ID] = true
		e.mu.Unlock()
		e.launch(task, pending.scheduledFor, pending.origin)
	}
}

// raiseAlert stores an alert and logs it.
func (e *Engine) raiseAlert(taskID string, sev domain.AlertSeverity, kind domain.AlertKind, msg string) {
	e.raiseRunAlert(taskID, "", sev, kind, msg)
}

func (e *Engine) raiseRunAlert(taskID, runID string, sev domain.AlertSeverity, kind domain.AlertKind, msg string) {
	a := domain.Alert{TaskID: taskID, RunID: runID, Severity: sev, Kind: kind, Message: msg}
	if err := e.store.CreateAlert(&a); err != nil {
		e.log.Error("engine: create alert", "err", err)
	}
	if e.onAlert != nil {
		e.onAlert(a)
	}
}
