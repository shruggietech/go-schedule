package engine

import (
	"fmt"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// RunNow triggers an immediate manual run of the task, honoring its overlap
// policy. It satisfies the API's Scheduler interface.
func (e *Engine) RunNow(taskID string) error {
	task, err := e.store.GetTask(taskID)
	if err != nil {
		return fmt.Errorf("run-now: %w", err)
	}
	e.dispatch(task, e.clk.Now(), domain.TriggerManual)
	return nil
}
