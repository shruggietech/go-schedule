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
	if err := e.dispatchTriggeredTask(trigger.TargetTaskID, dispatchOrigin{trigger: domain.TriggerExternal, sourceTriggerID: trigger.ID}); err != nil {
		return trigger.ID, err
	}
	return trigger.ID, nil
}

// FireFilesystemWatcher submits one request from an active watcher definition.
func (e *Engine) FireFilesystemWatcher(id string) error {
	watcher, err := e.store.GetFilesystemWatcher(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return ErrTriggerUnknown
		}
		return fmt.Errorf("fire filesystem watcher: %w", err)
	}
	if !watcher.Enabled {
		return ErrTriggerDisabled
	}
	return e.dispatchTriggeredTask(watcher.TargetTaskID, dispatchOrigin{trigger: domain.TriggerFilesystem, sourceWatcherID: watcher.ID})
}

func (e *Engine) dispatchTriggeredTask(taskID string, origin dispatchOrigin) error {
	task, err := e.store.GetTask(taskID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return ErrTriggerTargetMissing
		}
		return fmt.Errorf("dispatch triggered task: %w", err)
	}
	if strings.TrimSpace(task.Command) == "" {
		return ErrTriggerCommandIncomplete
	}
	if task.State != domain.TaskActive {
		return ErrTriggerTaskInactive
	}
	if !task.Enabled {
		return ErrTriggerTaskDisabled
	}
	groupsEnabled, err := e.store.GroupChainEnabled(task.GroupID)
	if err != nil {
		return fmt.Errorf("dispatch triggered task: %w", err)
	}
	if !groupsEnabled {
		return ErrTriggerGroupBlocked
	}
	e.dispatchWithOrigin(task, e.clk.Now(), origin)
	return nil
}
