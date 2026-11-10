package engine

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/clock"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/runcontrol"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

type cancellationCall struct {
	ctx      context.Context
	finish   chan domain.Run
	reported chan struct{}
}
type cancellationRunner struct {
	started chan *cancellationCall
	stop    chan struct{}
	now     time.Time
	hold    bool
}

func (r *cancellationRunner) Run(ctx context.Context, task domain.Task, at time.Time, trigger domain.RunTrigger) domain.Run {
	call := &cancellationCall{ctx: ctx, finish: make(chan domain.Run, 1), reported: make(chan struct{})}
	r.started <- call
	var result domain.Run
	select {
	case result = <-call.finish:
	case <-r.stop:
		result.Outcome = domain.OutcomeFailure
	case <-ctx.Done():
		if r.hold {
			runcontrol.ReportFailure(ctx, errors.New("private-cleanup-detail"))
			close(call.reported)
			select {
			case result = <-call.finish:
			case <-r.stop:
				result.Outcome = domain.OutcomeCancellationFailed
			}
		} else {
			result.Outcome = domain.OutcomeCanceled
			result.Output = "partial"
			code := -1
			result.ExitCode = &code
		}
	}
	result.TaskID = task.ID
	result.ScheduledFor = at
	result.StartedAt = &r.now
	result.EndedAt = &r.now
	result.Trigger = trigger
	return result
}
func newCancellationEngine(t *testing.T, workers int, hold bool) (*Engine, *store.Store, *cancellationRunner, <-chan domain.Run) {
	t.Helper()
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)
	runner := &cancellationRunner{started: make(chan *cancellationCall, 8), stop: make(chan struct{}), now: now, hold: hold}
	e := New(st, clock.NewFake(now), runner, testLogger(), workers)
	e.runCtx = context.Background()
	completed := make(chan domain.Run, 8)
	e.SetOnRun(func(run domain.Run) { completed <- run })
	t.Cleanup(func() {
		close(runner.stop)
		e.runWG.Wait()
		if err := st.Close(); err != nil {
			t.Error(err)
		}
	})
	return e, st, runner, completed
}
func awaitCancellationCall(t *testing.T, runner *cancellationRunner) *cancellationCall {
	t.Helper()
	select {
	case call := <-runner.started:
		return call
	case <-time.After(2 * time.Second):
		t.Fatal("execution did not start")
		return nil
	}
}
func awaitCancellationResult(t *testing.T, completed <-chan domain.Run, id string) domain.Run {
	t.Helper()
	watch := time.NewTimer(2 * time.Second)
	defer watch.Stop()
	for {
		select {
		case got := <-completed:
			if got.ID == id {
				return got
			}
			if got.Outcome != domain.OutcomeQueued && got.Outcome != domain.OutcomeSkipped && got.Outcome != domain.OutcomeCaughtUp {
				t.Fatalf("result=%+v want %s", got, id)
			}
		case <-watch.C:
			t.Fatal("terminal execution was not recorded")
			return domain.Run{}
		}
	}
}

func awaitReleasedWorker(t *testing.T, e *Engine) {
	t.Helper()
	select {
	case e.sem <- struct{}{}:
		<-e.sem
	case <-time.After(2 * time.Second):
		t.Fatal("execution did not release its worker")
	}
}

func TestCancelRunTargetsOneConcurrentExecutionAndPreservesOverlap(t *testing.T) {
	e, st, runner, completed := newCancellationEngine(t, 2, false)
	task := setupTask(t, st, domain.OverlapAllowConcurrent)
	now := runner.now
	e.dispatch(task, now, domain.TriggerManual)
	first := awaitCancellationCall(t, runner)
	firstID := e.ActiveRuns()[0].ID
	e.dispatch(task, now, domain.TriggerManual)
	second := awaitCancellationCall(t, runner)
	var secondID string
	for _, active := range e.ActiveRuns() {
		if active.ID != firstID {
			secondID = active.ID
		}
	}
	got, err := e.CancelRun(context.Background(), firstID)
	if err != nil || got.Status != "requested" || got.RunID != firstID || got.TaskID != task.ID {
		t.Fatalf("request=%+v err=%v", got, err)
	}
	again, err := e.CancelRun(context.Background(), firstID)
	if err != nil || (again.Status != "requested" && again.Status != "already_finished") {
		t.Fatalf("repeat=%+v err=%v", again, err)
	}
	ended := awaitCancellationResult(t, completed, firstID)
	awaitReleasedWorker(t, e)
	if ended.Outcome != domain.OutcomeCanceled || ended.Output != "partial" || ended.ExitCode == nil || *ended.ExitCode != -1 {
		t.Fatalf("canceled result=%+v", ended)
	}
	if !errors.Is(context.Cause(first.ctx), runcontrol.ErrRequested) {
		t.Fatalf("cause=%v", context.Cause(first.ctx))
	}
	select {
	case <-second.ctx.Done():
		t.Fatal("another concurrent execution was canceled")
	default:
	}
	active := e.ActiveRuns()
	if len(active) != 1 || active[0].ID != secondID {
		t.Fatalf("remaining=%+v", active)
	}
	task.OverlapPolicy = domain.OverlapSkip
	e.dispatch(task, now, domain.TriggerManual)
	if got := countOutcomes(t, st, task.ID)[domain.OutcomeSkipped]; got != 1 {
		t.Fatalf("other active instance was forgotten: skipped=%d", got)
	}
	second.finish <- domain.Run{Outcome: domain.OutcomeSuccess, Output: "unrelated"}
	awaitCancellationResult(t, completed, secondID)
	e.runWG.Wait()
	if err := e.RunNow(task.ID); err != nil {
		t.Fatal(err)
	}
	future := awaitCancellationCall(t, runner)
	futureID := e.ActiveRuns()[0].ID
	future.finish <- domain.Run{Outcome: domain.OutcomeSuccess}
	awaitCancellationResult(t, completed, futureID)
	e.runWG.Wait()
	out := countOutcomes(t, st, task.ID)
	if out[domain.OutcomeCanceled] != 1 || out[domain.OutcomeSuccess] != 2 {
		t.Fatalf("history=%v", out)
	}
}
func TestCancelRunBeforeRunnerBeginsAndAfterNaturalCompletion(t *testing.T) {
	t.Run("before launch", func(t *testing.T) {
		e, st, runner, completed := newCancellationEngine(t, 1, false)
		task := setupTask(t, st, domain.OverlapQueueOne)
		request := make(chan domain.RunCancellation, 1)
		e.SetOnRunStarted(func(run domain.Run) {
			result, err := e.CancelRun(context.Background(), run.ID)
			if err != nil {
				t.Error(err)
			}
			request <- result
		})
		e.dispatch(task, runner.now, domain.TriggerManual)
		call := awaitCancellationCall(t, runner)
		got := <-request
		if got.Status != "requested" || !errors.Is(context.Cause(call.ctx), runcontrol.ErrRequested) {
			t.Fatalf("request=%+v cause=%v", got, context.Cause(call.ctx))
		}
		if result := awaitCancellationResult(t, completed, got.RunID); result.Outcome != domain.OutcomeCanceled {
			t.Fatalf("result=%+v", result)
		}
	})
	t.Run("observed completion", func(t *testing.T) {
		e, st, runner, _ := newCancellationEngine(t, 1, false)
		task := setupTask(t, st, domain.OverlapQueueOne)
		result := make(chan domain.RunCancellation, 1)
		var call *cancellationCall
		e.SetOnRun(func(run domain.Run) {
			observed := e.ActiveRuns()
			if len(observed) != 1 || observed[0].ID != run.ID || observed[0].Outcome != domain.OutcomeSuccess || observed[0].Output != "" {
				t.Errorf("finishing observation=%+v", observed)
			}
			got, err := e.CancelRun(context.Background(), run.ID)
			if err != nil {
				t.Error(err)
			}
			if call.ctx.Err() != nil {
				t.Error("finishing execution context was canceled")
			}
			result <- got
		})
		e.dispatch(task, runner.now, domain.TriggerManual)
		call = awaitCancellationCall(t, runner)
		call.finish <- domain.Run{Outcome: domain.OutcomeSuccess, Output: "private-completion-output"}
		got := <-result
		if got.Status != "already_finished" || got.Outcome != domain.OutcomeSuccess {
			t.Fatalf("completion=%+v", got)
		}
		e.runWG.Wait()
		again, err := e.CancelRun(context.Background(), got.RunID)
		if err != nil || again.Status != "already_finished" || again.Outcome != domain.OutcomeSuccess {
			t.Fatalf("persisted=%+v err=%v", again, err)
		}
		persisted, err := st.GetRun(got.RunID)
		if err != nil || persisted.Output != "private-completion-output" {
			t.Fatalf("retained history=%+v err=%v", persisted, err)
		}
	})
}
func TestCancelRunFailureFeedbackRetainsActiveOwnershipAndWorker(t *testing.T) {
	e, st, runner, completed := newCancellationEngine(t, 1, true)
	task := setupTask(t, st, domain.OverlapAllowConcurrent)
	e.dispatch(task, runner.now, domain.TriggerManual)
	call := awaitCancellationCall(t, runner)
	id := e.ActiveRuns()[0].ID
	if _, err := e.CancelRun(context.Background(), id); err != nil {
		t.Fatal(err)
	}
	select {
	case <-call.reported:
	case <-time.After(2 * time.Second):
		t.Fatal("no cancellation failure feedback")
	}
	active := e.ActiveRuns()
	if len(active) != 1 || active[0].ID != id || active[0].EndedAt != nil || active[0].Outcome != "" || active[0].CancellationStatus != "failed" || active[0].CancellationError == "" {
		t.Fatalf("unconfirmed active=%+v", active)
	}
	again, err := e.CancelRun(context.Background(), id)
	if err != nil || again.Status != "requested" || again.CancellationStatus != "failed" {
		t.Fatalf("request=%+v err=%v", again, err)
	}
	if len(e.sem) != 1 {
		t.Fatal("unreaped work released capacity")
	}
	e.dispatch(task, runner.now, domain.TriggerManual)
	select {
	case <-runner.started:
		t.Fatal("unreaped work allowed another worker start")
	default:
	}
	zero := 0
	call.finish <- domain.Run{Outcome: domain.OutcomeCancellationFailed, Output: "retained", ExitCode: &zero}
	ended := awaitCancellationResult(t, completed, id)
	if ended.Outcome != domain.OutcomeCancellationFailed || ended.Output != "retained" || ended.ExitCode == nil || *ended.ExitCode != 0 {
		t.Fatalf("actual observed result=%+v", ended)
	}
	next := awaitCancellationCall(t, runner)
	next.finish <- domain.Run{Outcome: domain.OutcomeSuccess}
	e.runWG.Wait()
}
func TestCancelRunRejectsCallerCancellationAndNonExecutingHistory(t *testing.T) {
	e, st, runner, _ := newCancellationEngine(t, 1, false)
	task := setupTask(t, st, domain.OverlapAllowConcurrent)
	for _, outcome := range []domain.RunOutcome{domain.OutcomeQueued, domain.OutcomeSkipped, domain.OutcomeCaughtUp} {
		run := domain.Run{TaskID: task.ID, ScheduledFor: runner.now, Outcome: outcome, Trigger: domain.TriggerManual}
		if err := st.CreateRun(&run); err != nil {
			t.Fatal(err)
		}
		if _, err := e.CancelRun(context.Background(), run.ID); !errors.Is(err, runcontrol.ErrNotActive) {
			t.Fatalf("%s error=%v", outcome, err)
		}
	}
	if _, err := e.CancelRun(context.Background(), "unknown"); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("unknown=%v", err)
	}
	e.dispatch(task, runner.now, domain.TriggerManual)
	call := awaitCancellationCall(t, runner)
	id := e.ActiveRuns()[0].ID
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := e.CancelRun(ctx, id); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled caller=%v", err)
	}
	if call.ctx.Err() != nil || e.ActiveRuns()[0].CancellationStatus != "" {
		t.Fatal("canceled caller mutated the execution")
	}
	call.finish <- domain.Run{Outcome: domain.OutcomeSuccess}
	e.runWG.Wait()
}
func TestQueuedRunWaitsForAllConcurrentExecutions(t *testing.T) {
	e, st, runner, completed := newCancellationEngine(t, 2, false)
	task := setupTask(t, st, domain.OverlapAllowConcurrent)
	e.dispatch(task, runner.now, domain.TriggerManual)
	first := awaitCancellationCall(t, runner)
	id := e.ActiveRuns()[0].ID
	e.dispatch(task, runner.now, domain.TriggerManual)
	second := awaitCancellationCall(t, runner)
	task.OverlapPolicy = domain.OverlapQueueOne
	e.dispatch(task, runner.now, domain.TriggerExternal)
	first.finish <- domain.Run{Outcome: domain.OutcomeSuccess}
	awaitCancellationResult(t, completed, id)
	awaitReleasedWorker(t, e)
	e.mu.Lock()
	count := e.running[task.ID]
	_, queued := e.queued[task.ID]
	e.mu.Unlock()
	if count != 1 || !queued {
		t.Fatalf("remaining reservations=%d queued=%v", count, queued)
	}
	select {
	case <-runner.started:
		t.Fatal("queued execution started before other active instance completed")
	default:
	}
	second.finish <- domain.Run{Outcome: domain.OutcomeSuccess}
	next := awaitCancellationCall(t, runner)
	next.finish <- domain.Run{Outcome: domain.OutcomeSuccess}
	e.runWG.Wait()
	runs, err := st.ListRuns(task.ID, 0)
	if err != nil {
		t.Fatal(err)
	}
	external := 0
	for _, run := range runs {
		if run.Outcome == domain.OutcomeSuccess && run.Trigger == domain.TriggerExternal {
			external++
		}
	}
	if external != 1 {
		t.Fatalf("queued trigger not retained: %+v", runs)
	}
}
