package automation

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	tasklogic "github.com/shruggietech/go-schedule/internal/task"
)

const callTimeout = 3 * time.Second

type Service struct{ backend Backend }

func NewService(backend Backend) *Service { return &Service{backend: backend} }

func (s *Service) Workspace(ctx context.Context) OperationResult {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	tasks, err := s.backend.ListTaskDetails(c, "", "")
	if err != nil {
		return failure("load", err)
	}
	chains, err := s.backend.ListChains(c)
	if err != nil {
		return failure("load", err)
	}
	triggers, err := s.backend.ListTriggers(c)
	if err != nil {
		return failure("load", err)
	}
	sets, err := s.backend.ListTriggerSets(c)
	if err != nil {
		return failure("load", err)
	}
	watchers, err := s.backend.ListFilesystemWatchers(c)
	if err != nil {
		return failure("load", err)
	}
	workspace := buildWorkspace(tasks, chains, triggers, sets, watchers)
	return accepted("load", "Automation sources are up to date.", "", &workspace)
}

func buildWorkspace(tasks []server.TaskResponse, chains []domain.CompletionChain, triggers []server.TriggerResponse, sets []server.TriggerSetResponse, watchers []server.FilesystemWatcherResponse) Workspace {
	w := Workspace{Tasks: make([]TaskChoice, 0, len(tasks)), Chains: make([]ChainSummary, 0, len(chains)), Triggers: []TriggerSummary{}, TriggerSets: make([]TriggerSetSummary, 0, len(sets)), Watchers: make([]WatcherSummary, 0, len(watchers)), LoadedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	for _, value := range tasks {
		name := tasklogic.DisplayName(value.Task)
		w.Tasks = append(w.Tasks, TaskChoice{ID: value.Task.ID, Name: name, Readiness: string(value.Readiness.Status), Reason: value.Readiness.Reason})
	}
	for _, value := range chains {
		readiness, reason := "ready", "Runs the target after the selected source outcome."
		if value.SourceTaskName == "" {
			readiness, reason = "source_missing", "Source task is missing."
		} else if value.TargetTaskName == "" {
			readiness, reason = "target_missing", "Target task is missing."
		}
		w.Chains = append(w.Chains, ChainSummary{ID: value.ID, SourceTaskID: value.SourceTaskID, SourceTaskName: fallback(value.SourceTaskName, "Missing task"), TargetTaskID: value.TargetTaskID, TargetTaskName: fallback(value.TargetTaskName, "Missing task"), OnOutcome: string(value.OnOutcome), Readiness: readiness, Reason: reason, UpdatedAt: value.UpdatedAt.Format(time.RFC3339Nano)})
	}
	for _, value := range triggers {
		if value.SetID != "" {
			continue
		}
		w.Triggers = append(w.Triggers, TriggerSummary{ID: value.ID, Name: value.Name, TargetTaskID: value.TargetTaskID, TargetTaskName: fallback(value.TargetTaskName, "Missing task"), Readiness: value.Readiness, Reason: value.Reason, UpdatedAt: value.UpdatedAt.Format(time.RFC3339Nano), Enabled: value.Enabled})
	}
	for _, value := range sets {
		summary := TriggerSetSummary{ID: value.ID, Name: value.Name, TargetTaskID: value.TargetTaskID, TargetTaskName: fallback(value.TargetTaskName, "Missing task"), MemberCount: value.MemberCount, EnabledCount: value.EnabledCount, Readiness: "ready", Reason: "All enabled members are ready.", UpdatedAt: value.UpdatedAt.Format(time.RFC3339Nano), Members: make([]TriggerSetMember, 0, len(value.Members))}
		for _, member := range value.Members {
			summary.Members = append(summary.Members, TriggerSetMember{ID: member.ID, Name: member.Name, Position: member.SetPosition, Enabled: member.Enabled, Readiness: member.Readiness, Reason: member.Reason})
			if member.Readiness != "ready" && summary.Readiness == "ready" {
				summary.Readiness, summary.Reason = member.Readiness, member.Reason
			}
		}
		sort.SliceStable(summary.Members, func(i, j int) bool { return summary.Members[i].Position < summary.Members[j].Position })
		w.TriggerSets = append(w.TriggerSets, summary)
	}
	for _, value := range watchers {
		w.Watchers = append(w.Watchers, WatcherSummary{ID: value.ID, Name: value.Name, Kind: string(value.Kind), Path: value.Path, Pattern: value.Pattern, Recursive: value.Recursive, Debounce: value.Debounce, Stability: value.Stability, TargetTaskID: value.TargetTaskID, TargetTaskName: fallback(value.TargetTaskName, "Missing task"), Enabled: value.Enabled, Health: string(value.Health.State), HealthReason: value.Health.Reason, Readiness: value.Readiness, Reason: value.Reason, UpdatedAt: value.UpdatedAt.Format(time.RFC3339Nano)})
	}
	sort.SliceStable(w.Tasks, func(i, j int) bool { return strings.ToLower(w.Tasks[i].Name) < strings.ToLower(w.Tasks[j].Name) })
	return w
}

func (s *Service) SaveChain(ctx context.Context, draft ChainDraft) OperationResult {
	if draft.SourceTaskID == "" {
		return rejected("save_chain", "sourceTaskId", "Select a source task.")
	}
	if draft.TargetTaskID == "" {
		return rejected("save_chain", "targetTaskId", "Select a target task.")
	}
	if draft.SourceTaskID == draft.TargetTaskID {
		return rejected("save_chain", "targetTaskId", "Source and target tasks must be different.")
	}
	if draft.OnOutcome != "success" && draft.OnOutcome != "failure" && draft.OnOutcome != "any" {
		return rejected("save_chain", "onOutcome", "Select a completion outcome.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	entityID := draft.ID
	if draft.IsNew {
		value, err := s.backend.CreateChain(c, server.ChainCreateRequest{SourceTaskID: draft.SourceTaskID, TargetTaskID: draft.TargetTaskID, OnOutcome: draft.OnOutcome})
		if err != nil {
			return failure("save_chain", err)
		}
		entityID = value.ID
	} else {
		current, err := s.backend.GetChain(c, draft.ID)
		if err != nil {
			return failure("save_chain", err)
		}
		if stale(current.UpdatedAt, draft.OriginalUpdatedAt, draft.OverwriteStale) {
			return staleResult("save_chain", draft.ID)
		}
		source, target, outcome := draft.SourceTaskID, draft.TargetTaskID, draft.OnOutcome
		if _, err = s.backend.UpdateChain(c, draft.ID, server.ChainUpdateRequest{SourceTaskID: &source, TargetTaskID: &target, OnOutcome: &outcome}); err != nil {
			return failure("save_chain", err)
		}
	}
	return s.refreshed(ctx, "save_chain", "Completion chain saved.", entityID)
}

func (s *Service) DeleteChain(ctx context.Context, id string) OperationResult {
	return s.mutate(ctx, "delete_chain", id, func(c context.Context) error { return s.backend.DeleteChain(c, id) })
}

func (s *Service) SaveTrigger(ctx context.Context, draft TriggerDraft) SecretResult {
	name := strings.TrimSpace(draft.Name)
	if name == "" {
		return secretRejected("save_trigger", "name", "Enter a trigger name.")
	}
	if draft.TargetTaskID == "" {
		return secretRejected("save_trigger", "targetTaskId", "Select a target task.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	if draft.IsNew {
		value, err := s.backend.CreateTrigger(c, server.TriggerCreateRequest{Name: name, TargetTaskID: draft.TargetTaskID, Enabled: &draft.Enabled})
		if err != nil {
			return secretFailure("save_trigger", err)
		}
		return triggerSecret("save_trigger", "Trigger created. Save this key now.", value)
	}
	current, err := s.backend.GetTrigger(c, draft.ID)
	if err != nil {
		return secretFailure("save_trigger", err)
	}
	if stale(current.UpdatedAt, draft.OriginalUpdatedAt, draft.OverwriteStale) {
		return SecretResult{Action: "save_trigger", Outcome: "stale", Message: "This trigger changed after you opened it. Reload it or explicitly overwrite the newer version.", EntityID: draft.ID}
	}
	if _, err = s.backend.UpdateTrigger(c, draft.ID, server.TriggerUpdateRequest{Name: &name, TargetTaskID: &draft.TargetTaskID}); err != nil {
		return secretFailure("save_trigger", err)
	}
	if current.Enabled != draft.Enabled {
		if _, err = s.backend.SetTriggerEnabled(c, draft.ID, draft.Enabled); err != nil {
			return secretFailure("save_trigger", err)
		}
	}
	return SecretResult{Action: "save_trigger", Outcome: "accepted", Message: "Trigger saved.", EntityID: draft.ID}
}

func (s *Service) SetTriggerEnabled(ctx context.Context, id string, enabled bool) OperationResult {
	return s.mutate(ctx, "toggle_trigger", id, func(c context.Context) error { _, err := s.backend.SetTriggerEnabled(c, id, enabled); return err })
}
func (s *Service) RevealTrigger(ctx context.Context, id string) SecretResult {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	value, err := s.backend.RevealTrigger(c, id)
	if err != nil {
		return secretFailure("reveal_trigger", err)
	}
	return triggerSecret("reveal_trigger", "Trigger key revealed.", value)
}
func (s *Service) RotateTrigger(ctx context.Context, id string) SecretResult {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	value, err := s.backend.RotateTrigger(c, id)
	if err != nil {
		return secretFailure("rotate_trigger", err)
	}
	return triggerSecret("rotate_trigger", "Trigger key rotated. Save the replacement now.", value)
}
func (s *Service) FireTrigger(ctx context.Context, id string) OperationResult {
	return s.mutate(ctx, "fire_trigger", id, func(c context.Context) error {
		value, err := s.backend.RevealTrigger(c, id)
		if err != nil {
			return err
		}
		return s.backend.FireTrigger(c, value.Key)
	})
}
func (s *Service) DeleteTrigger(ctx context.Context, id string) OperationResult {
	return s.mutate(ctx, "delete_trigger", id, func(c context.Context) error { return s.backend.DeleteTrigger(c, id) })
}

func (s *Service) CreateTriggerSet(ctx context.Context, draft TriggerSetDraft) SecretResult {
	name := strings.TrimSpace(draft.Name)
	if name == "" {
		return secretRejected("create_trigger_set", "name", "Enter a Trigger Set name.")
	}
	if draft.TargetTaskID == "" {
		return secretRejected("create_trigger_set", "targetTaskId", "Select a target task.")
	}
	if draft.Count < 1 || draft.Count > 99 {
		return secretRejected("create_trigger_set", "count", "Enter a member count from 1 through 99.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	value, err := s.backend.CreateTriggerSet(c, server.TriggerSetCreateRequest{Name: name, TargetTaskID: draft.TargetTaskID, Count: draft.Count, Enabled: &draft.Enabled})
	if err != nil {
		return secretFailure("create_trigger_set", err)
	}
	return setSecret("create_trigger_set", "Trigger Set created. Save these keys now.", value)
}

func (s *Service) RetargetTriggerSet(ctx context.Context, id, target, original string, overwrite bool) OperationResult {
	if target == "" {
		return rejected("retarget_trigger_set", "targetTaskId", "Select a target task.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	current, err := s.backend.GetTriggerSet(c, id)
	if err != nil {
		return failure("retarget_trigger_set", err)
	}
	if stale(current.UpdatedAt, original, overwrite) {
		return staleResult("retarget_trigger_set", id)
	}
	if _, err = s.backend.RetargetTriggerSet(c, id, target); err != nil {
		return failure("retarget_trigger_set", err)
	}
	return s.refreshed(ctx, "retarget_trigger_set", "Trigger Set retargeted.", id)
}
func (s *Service) SetTriggerSetEnabled(ctx context.Context, id string, enabled bool) OperationResult {
	return s.mutate(ctx, "toggle_trigger_set", id, func(c context.Context) error { _, err := s.backend.SetTriggerSetEnabled(c, id, enabled); return err })
}
func (s *Service) RevealTriggerSet(ctx context.Context, id string) SecretResult {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	value, err := s.backend.RevealTriggerSet(c, id)
	if err != nil {
		return secretFailure("reveal_trigger_set", err)
	}
	return setSecret("reveal_trigger_set", "Trigger Set keys revealed.", value)
}
func (s *Service) RotateTriggerSet(ctx context.Context, id string) SecretResult {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	value, err := s.backend.RotateTriggerSet(c, id)
	if err != nil {
		return secretFailure("rotate_trigger_set", err)
	}
	return setSecret("rotate_trigger_set", "Trigger Set keys rotated. Save the replacements now.", value)
}
func (s *Service) DeleteTriggerSet(ctx context.Context, id string) OperationResult {
	return s.mutate(ctx, "delete_trigger_set", id, func(c context.Context) error { return s.backend.DeleteTriggerSet(c, id) })
}

func (s *Service) SaveWatcher(ctx context.Context, draft WatcherDraft) OperationResult {
	name := strings.TrimSpace(draft.Name)
	if name == "" {
		return rejected("save_watcher", "name", "Enter a watcher name.")
	}
	if draft.Kind != "file" && draft.Kind != "directory" {
		return rejected("save_watcher", "kind", "Select a watcher kind.")
	}
	if strings.TrimSpace(draft.Path) == "" {
		return rejected("save_watcher", "path", "Enter a path.")
	}
	if draft.TargetTaskID == "" {
		return rejected("save_watcher", "targetTaskId", "Select a target task.")
	}
	if _, err := time.ParseDuration(draft.Debounce); err != nil {
		return rejected("save_watcher", "debounce", "Enter a valid debounce duration.")
	}
	if _, err := time.ParseDuration(draft.Stability); err != nil {
		return rejected("save_watcher", "stability", "Enter a valid stability duration.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	entityID := draft.ID
	if draft.IsNew {
		value, err := s.backend.CreateFilesystemWatcher(c, server.FilesystemWatcherCreateRequest{Name: name, Kind: domain.WatcherKind(draft.Kind), Path: draft.Path, Pattern: draft.Pattern, Recursive: draft.Recursive, Debounce: draft.Debounce, Stability: draft.Stability, TargetTaskID: draft.TargetTaskID, Enabled: &draft.Enabled})
		if err != nil {
			return failure("save_watcher", err)
		}
		entityID = value.ID
	} else {
		current, err := s.backend.GetFilesystemWatcher(c, draft.ID)
		if err != nil {
			return failure("save_watcher", err)
		}
		if stale(current.UpdatedAt, draft.OriginalUpdatedAt, draft.OverwriteStale) {
			return staleResult("save_watcher", draft.ID)
		}
		kind, path, pattern, recursive, debounce, stability, target, enabled := domain.WatcherKind(draft.Kind), draft.Path, draft.Pattern, draft.Recursive, draft.Debounce, draft.Stability, draft.TargetTaskID, draft.Enabled
		_, err = s.backend.UpdateFilesystemWatcher(c, draft.ID, server.FilesystemWatcherUpdateRequest{Name: &name, Kind: &kind, Path: &path, Pattern: &pattern, Recursive: &recursive, Debounce: &debounce, Stability: &stability, TargetTaskID: &target, Enabled: &enabled})
		if err != nil {
			return failure("save_watcher", err)
		}
	}
	return s.refreshed(ctx, "save_watcher", "Filesystem watcher saved.", entityID)
}
func (s *Service) SetWatcherEnabled(ctx context.Context, id string, enabled bool) OperationResult {
	return s.mutate(ctx, "toggle_watcher", id, func(c context.Context) error {
		_, err := s.backend.SetFilesystemWatcherEnabled(c, id, enabled)
		return err
	})
}
func (s *Service) DeleteWatcher(ctx context.Context, id string) OperationResult {
	return s.mutate(ctx, "delete_watcher", id, func(c context.Context) error { return s.backend.DeleteFilesystemWatcher(c, id) })
}

func (s *Service) mutate(ctx context.Context, action, id string, fn func(context.Context) error) OperationResult {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	if err := fn(c); err != nil {
		return failure(action, err)
	}
	return s.refreshed(ctx, action, "Automation source action accepted.", id)
}
func (s *Service) refreshed(ctx context.Context, action, message, id string) OperationResult {
	loaded := s.Workspace(ctx)
	if loaded.Outcome != "accepted" {
		return accepted(action, message+" Refresh the workspace to see the latest state.", id, nil)
	}
	return accepted(action, message, id, loaded.Workspace)
}
func fallback(value, fallbackValue string) string {
	if strings.TrimSpace(value) == "" {
		return fallbackValue
	}
	return value
}
func stale(updated time.Time, original string, overwrite bool) bool {
	return !overwrite && original != "" && updated.Format(time.RFC3339Nano) != original
}
func accepted(action, message, id string, workspace *Workspace) OperationResult {
	return OperationResult{Action: action, Outcome: "accepted", Message: message, EntityID: id, Workspace: workspace}
}
func rejected(action, field, message string) OperationResult {
	return OperationResult{Action: action, Outcome: "rejected", Message: message, Field: field}
}
func staleResult(action, id string) OperationResult {
	return OperationResult{Action: action, Outcome: "stale", Message: "This item changed after you opened it. Reload it or explicitly overwrite the newer version.", EntityID: id}
}
func triggerSecret(action, message string, value server.TriggerSecretResponse) SecretResult {
	return SecretResult{Action: action, Outcome: "accepted", Message: message, EntityID: value.Trigger.ID, Title: value.Trigger.Name, Secrets: []Secret{{Label: value.Trigger.Name, Key: value.Key, Command: value.Command}}}
}
func setSecret(action, message string, value server.TriggerSetSecretResponse) SecretResult {
	out := SecretResult{Action: action, Outcome: "accepted", Message: message, EntityID: value.TriggerSet.ID, Title: value.TriggerSet.Name, Secrets: make([]Secret, 0, len(value.Members))}
	for _, member := range value.Members {
		out.Secrets = append(out.Secrets, Secret{Label: fmt.Sprintf("Member %d", member.Position), Key: member.Key, Command: member.Command})
	}
	return out
}
func secretRejected(action, field, message string) SecretResult {
	return SecretResult{Action: action, Outcome: "rejected", Message: message, Field: field}
}
func failure(action string, err error) OperationResult {
	outcome, field, message := safeError(err)
	return OperationResult{Action: action, Outcome: outcome, Field: field, Message: message}
}
func secretFailure(action string, err error) SecretResult {
	outcome, field, message := safeError(err)
	return SecretResult{Action: action, Outcome: outcome, Field: field, Message: message}
}
func safeError(err error) (string, string, string) {
	var status *client.StatusError
	if errors.As(err, &status) {
		fields := map[string]string{"name": "name", "source_task_id": "sourceTaskId", "target_task_id": "targetTaskId", "on_outcome": "onOutcome", "kind": "kind", "path": "path", "pattern": "pattern", "debounce": "debounce", "stability": "stability", "count": "count"}
		field := fields[status.Field]
		if status.Code == server.CodeNotFound {
			return "rejected", field, "The selected automation source no longer exists."
		}
		return "rejected", field, "The request was rejected. Review the highlighted field and try again."
	}
	return "unavailable", "", "The local scheduler service is unavailable. Try again."
}
