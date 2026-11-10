package engine

import (
	"context"
	"fmt"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/runcontrol"
)

// runCancellation belongs to one stable execution identity, never a PID.
// engine.mu guards its request and observed-result transitions.
type runCancellation struct {
	cancel    context.CancelCauseFunc
	requested bool
	result    *domain.Run
}

// CancelRun requests cancellation of one active execution without waiting for termination.
// Once the runner has returned, its actual result wins even while persistence is pending.
func (e *Engine) CancelRun(ctx context.Context, id string) (domain.RunCancellation, error) {
	if err := ctx.Err(); err != nil {
		return domain.RunCancellation{}, err
	}
	e.mu.Lock()
	if control, ok := e.cancellations[id]; ok {
		if control.result != nil {
			result := finishedCancellation(*control.result)
			e.mu.Unlock()
			return result, nil
		}
		if err := ctx.Err(); err != nil {
			e.mu.Unlock()
			return domain.RunCancellation{}, err
		}
		active := e.active[id]
		if !control.requested {
			control.requested = true
			active.CancellationStatus = "requested"
			e.active[id] = active
		}
		result := domain.RunCancellation{RunID: id, TaskID: active.TaskID, Status: "requested", CancellationStatus: active.CancellationStatus, CancellationError: active.CancellationError}
		cancel := control.cancel
		e.mu.Unlock()
		cancel(runcontrol.ErrRequested)
		return result, nil
	}
	e.mu.Unlock()
	run, err := e.store.GetRunResultContext(ctx, id)
	if err != nil {
		return domain.RunCancellation{}, err
	}
	switch run.Outcome {
	case domain.OutcomeSuccess, domain.OutcomeFailure, domain.OutcomeCanceled, domain.OutcomeCancellationFailed:
		return finishedCancellation(run), nil
	}
	return domain.RunCancellation{}, fmt.Errorf("cancel run: %w", runcontrol.ErrNotActive)
}

func finishedCancellation(run domain.Run) domain.RunCancellation {
	return domain.RunCancellation{RunID: run.ID, TaskID: run.TaskID, Status: "already_finished", Outcome: run.Outcome, CancellationStatus: run.CancellationStatus, CancellationError: run.CancellationError}
}

func (e *Engine) cancellationFailed(id string, err error) {
	e.mu.Lock()
	control, exists := e.cancellations[id]
	if !exists || control.result != nil || !control.requested {
		e.mu.Unlock()
		return
	}
	active := e.active[id]
	active.CancellationStatus = "failed"
	active.CancellationError = "Cancellation cleanup could not be confirmed. Refresh this run and inspect the daemon log."
	e.active[id] = active
	e.mu.Unlock()
	e.log.Error("engine: cancellation cleanup unconfirmed", "run", id, "task", active.TaskID, "err", err)
}
