package engine

import (
	"time"

	"github.com/shruggietech/go-scheduler/internal/domain"
)

// dispatch applies the task's overlap policy and either runs the task now,
// queues a single pending run, or skips it.
//
//   - queue_one (default): if the task is running, queue exactly one pending run,
//     log a warning, and raise an overlap alert; additional triggers while one is
//     already queued are dropped.
//   - skip: if running, record a skipped run and drop the trigger.
//   - allow_concurrent: always run, even if an instance is already running.
func (e *Engine) dispatch(task domain.Task, scheduledFor time.Time, trigger domain.RunTrigger) {
	e.mu.Lock()
	if e.running[task.ID] {
		switch task.OverlapPolicy {
		case domain.OverlapSkip:
			e.mu.Unlock()
			e.recordRun(domain.Run{
				TaskID: task.ID, ScheduledFor: scheduledFor,
				Outcome: domain.OutcomeSkipped, Trigger: trigger,
			})
			return

		case domain.OverlapAllowConcurrent:
			// fall through to run concurrently

		default: // OverlapQueueOne (and unset)
			if _, already := e.queued[task.ID]; already {
				e.mu.Unlock()
				return // one already pending; drop extra triggers
			}
			e.queued[task.ID] = scheduledFor
			e.mu.Unlock()
			e.log.Warn("task still running at next trigger; queued one pending run",
				"task", task.ID, "name", task.Name, "scheduled_for", scheduledFor)
			e.raiseAlert(task.ID, domain.SeverityWarning, domain.AlertOverlapQueued,
				"previous run still in progress; queued one pending run")
			e.recordRun(domain.Run{
				TaskID: task.ID, ScheduledFor: scheduledFor,
				Outcome: domain.OutcomeQueued, Trigger: trigger,
			})
			return
		}
	}
	e.running[task.ID] = true
	e.mu.Unlock()
	e.launch(task, scheduledFor, trigger)
}
