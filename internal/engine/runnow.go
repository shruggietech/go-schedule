package engine

import (
	"fmt"
	"strings"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

// RunNow triggers an immediate manual run of the task, honoring its overlap
// policy. It satisfies the API's Scheduler interface.
func (e *Engine) RunNow(taskID string) error {
	task, err := e.store.GetTask(taskID)
	if err != nil {
		return fmt.Errorf("run-now: %w", err)
	}
	if strings.TrimSpace(task.Command) == "" {
		return fmt.Errorf("run-now: %w", store.ErrTaskNotRunnable)
	}
	e.dispatch(task, e.clk.Now(), domain.TriggerManual)
	return nil
}
