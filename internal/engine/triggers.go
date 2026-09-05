package engine

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

var (
	ErrTriggerUnknown           = errors.New("external trigger is unknown")
	ErrTriggerDisabled          = errors.New("external trigger is disabled")
	ErrTriggerTargetMissing     = errors.New("external trigger target is missing")
	ErrTriggerCommandIncomplete = errors.New("external trigger target has no command")
	ErrTriggerTaskInactive      = errors.New("external trigger target is inactive")
	ErrTriggerTaskDisabled      = errors.New("external trigger target is disabled")
	ErrTriggerGroupBlocked      = errors.New("external trigger target is blocked by a group")
)

// FireExternalTrigger validates one opaque key and submits one request to the
// normal overlap-aware dispatcher.
func (e *Engine) FireExternalTrigger(key string) (string, error) {
	trigger, err := e.store.GetExternalTriggerByKey(key)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return "", ErrTriggerUnknown
		}
		return "", fmt.Errorf("fire external trigger: %w", err)
	}
	if !trigger.Enabled {
		return trigger.ID, ErrTriggerDisabled
	}
	task, err := e.store.GetTask(trigger.TargetTaskID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return trigger.ID, ErrTriggerTargetMissing
		}
		return trigger.ID, fmt.Errorf("fire external trigger: %w", err)
	}
	if strings.TrimSpace(task.Command) == "" {
		return trigger.ID, ErrTriggerCommandIncomplete
	}
	if task.State != domain.TaskActive {
		return trigger.ID, ErrTriggerTaskInactive
	}
	if !task.Enabled {
		return trigger.ID, ErrTriggerTaskDisabled
	}
	groupsEnabled, err := e.store.GroupChainEnabled(task.GroupID)
	if err != nil {
		return trigger.ID, fmt.Errorf("fire external trigger: %w", err)
	}
	if !groupsEnabled {
		return trigger.ID, ErrTriggerGroupBlocked
	}
	e.dispatchWithOrigin(task, e.clk.Now(), dispatchOrigin{trigger: domain.TriggerExternal, sourceTriggerID: trigger.ID})
	return trigger.ID, nil
}
