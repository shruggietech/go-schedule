package search

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/domain"
)

// Execute revalidates every source and object immediately before dispatch.
func (s *Service) Execute(ctx context.Context, intent ActionIntent) ActionBatchResult {
	if !mutableAction(intent.Action) || len(intent.Selections) == 0 {
		return ActionBatchResult{Action: intent.Action, Outcome: "rejected", Message: "Choose one supported action and at least one compatible result.", Outcomes: []ActionOutcome{}}
	}
	for _, selection := range intent.Selections {
		if selection.RegistrationKey == "" || selection.ExpectedDaemonID == "" || selection.ObjectID == "" || !compatible(intent.Action, selection.Kind) {
			return ActionBatchResult{Action: intent.Action, Outcome: "rejected", Message: "Every selected result must identify one compatible object and daemon.", Outcomes: []ActionOutcome{}}
		}
	}
	groups := map[string][]Selection{}
	for _, selection := range intent.Selections {
		groups[selection.RegistrationKey] = append(groups[selection.RegistrationKey], selection)
	}
	results := make(chan []ActionOutcome, len(groups))
	var workers sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrentTargets)
	for key, selections := range groups {
		workers.Add(1)
		go func(key string, selections []Selection) {
			defer workers.Done()
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				results <- rejectAll(selections, "The action was canceled before this scheduler was reached.", "")
				return
			}
			results <- s.executeTarget(ctx, intent.Action, key, selections)
		}(key, selections)
	}
	go func() { workers.Wait(); close(results) }()
	outcomes := make([]ActionOutcome, 0, len(intent.Selections))
	for values := range results {
		outcomes = append(outcomes, values...)
	}
	accepted := 0
	uncertain := 0
	for _, outcome := range outcomes {
		if outcome.Outcome == "accepted" {
			accepted++
		} else if outcome.Outcome == "uncertain" {
			uncertain++
		}
	}
	summary := "rejected"
	if accepted == len(outcomes) {
		summary = "accepted"
	} else if accepted > 0 || uncertain > 0 {
		summary = "partial"
	}
	return ActionBatchResult{Action: intent.Action, Outcome: summary, Message: fmt.Sprintf("%d of %d requested objects accepted. Review the target-specific outcomes.", accepted, len(outcomes)), Outcomes: outcomes}
}

func (s *Service) executeTarget(ctx context.Context, action domain.SearchAction, key string, selections []Selection) []ActionOutcome {
	value, err := s.targetFor(key)
	if err != nil {
		return rejectAll(selections, "The saved connection is unavailable.", "")
	}
	health, err := value.backend.Health(ctx)
	if err != nil {
		return rejectAll(selections, "The scheduler identity or authority could not be revalidated.", "")
	}
	if !canOperate(health.Permissions) {
		return rejectAll(selections, "The current credential does not have Operate authority.", health.ID)
	}
	results := make([]ActionOutcome, 0, len(selections))
	for _, selection := range selections {
		base := ActionOutcome{RegistrationKey: selection.RegistrationKey, ExpectedDaemonID: selection.ExpectedDaemonID, CurrentDaemonID: health.ID, Kind: string(selection.Kind), ObjectID: selection.ObjectID, DisplayName: selection.DisplayName, Outcome: "rejected"}
		if health.ID != selection.ExpectedDaemonID {
			base.Message = "The scheduler identity changed. No action was performed."
			results = append(results, base)
			continue
		}
		if err := revalidateAndDispatch(ctx, value.client, action, selection); err != nil {
			var uncertain *client.MutationUncertainError
			if errors.As(err, &uncertain) {
				base.Outcome = "uncertain"
				base.Message = "The remote action may have completed. Refresh before deciding whether to try again."
			} else {
				base.Message = actionFailureMessage(err)
			}
			results = append(results, base)
			continue
		}
		base.Outcome, base.Message = "accepted", "The action was accepted by "+value.registration.Label+"."
		results = append(results, base)
	}
	return results
}

func canOperate(permissions []string) bool {
	for _, permission := range permissions {
		if permission == "operate" || permission == "manage" || permission == "enroll" {
			return true
		}
	}
	return false
}

func revalidateAndDispatch(ctx context.Context, daemon daemonClient, action domain.SearchAction, selection Selection) error {
	switch action {
	case domain.SearchActionEnable, domain.SearchActionDisable, domain.SearchActionRunNow:
		task, err := daemon.GetTask(ctx, selection.ObjectID)
		if err != nil {
			return err
		}
		if task.Task.ID != selection.ObjectID {
			return errors.New("task identity mismatch")
		}
		if action == domain.SearchActionEnable {
			if task.Task.Enabled {
				return errors.New("task is already enabled")
			}
			return daemon.SetTaskEnabled(ctx, selection.ObjectID, true)
		}
		if action == domain.SearchActionDisable {
			if !task.Task.Enabled {
				return errors.New("task is already disabled")
			}
			return daemon.SetTaskEnabled(ctx, selection.ObjectID, false)
		}
		return daemon.RunNow(ctx, selection.ObjectID)
	case domain.SearchActionAcknowledge:
		const pageSize = 200
		for offset := 0; ; offset += pageSize {
			alerts, err := daemon.ListAlertsPage(ctx, true, offset, pageSize, 0)
			if err != nil {
				return err
			}
			for _, alert := range alerts {
				if alert.ID == selection.ObjectID {
					return daemon.AckAlert(ctx, selection.ObjectID)
				}
			}
			if len(alerts) < pageSize {
				break
			}
		}
		return errors.New("alert no longer exists")
	default:
		return errors.New("unsupported action")
	}
}

func mutableAction(action domain.SearchAction) bool {
	return action == domain.SearchActionAcknowledge || action == domain.SearchActionEnable || action == domain.SearchActionDisable || action == domain.SearchActionRunNow
}
func compatible(action domain.SearchAction, kind domain.SearchKind) bool {
	if action == domain.SearchActionAcknowledge {
		return kind == domain.SearchKindAlert
	}
	return kind == domain.SearchKindTask
}
func rejectAll(selections []Selection, message, current string) []ActionOutcome {
	result := make([]ActionOutcome, 0, len(selections))
	for _, selection := range selections {
		result = append(result, ActionOutcome{RegistrationKey: selection.RegistrationKey, ExpectedDaemonID: selection.ExpectedDaemonID, CurrentDaemonID: current, Kind: string(selection.Kind), ObjectID: selection.ObjectID, DisplayName: selection.DisplayName, Outcome: "rejected", Message: message})
	}
	return result
}
func actionFailureMessage(err error) string {
	var status *client.StatusError
	if errors.As(err, &status) && status.Code == "not_found" {
		return "The selected object no longer exists. No action was performed."
	}
	return "The object changed or the action could not be confirmed. Refresh before retrying."
}
