package notifications

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	webhook "github.com/shruggietech/go-schedule/internal/notification"
)

const callTimeout = 3 * time.Second
const deliveryLimit = 200

type Service struct {
	backend Backend
	now     func() time.Time
}

func NewService(backend Backend) *Service { return &Service{backend: backend, now: time.Now} }

func (s *Service) Workspace(ctx context.Context) Result {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	channels, err := s.backend.ListNotificationChannels(c)
	if err != nil {
		return failure("load_notifications", err)
	}
	tasks, err := s.backend.ListTaskDetails(c, "", "")
	if err != nil {
		return failure("load_notifications", err)
	}
	groups, err := s.backend.ListGroups(c)
	if err != nil {
		return failure("load_notifications", err)
	}
	deliveries, err := s.backend.ListNotificationDeliveries(c, domain.NotificationDeliveryFilter{Limit: deliveryLimit})
	if err != nil {
		return failure("load_notifications", err)
	}
	workspace := buildWorkspace(channels, tasks, groups, deliveries, s.now())
	workspace.Coverage, workspace.CoverageComplete = s.configuredCoverage(c, &workspace)
	return Result{Action: "load_notifications", Outcome: "accepted", Message: "Notifications are up to date.", Workspace: &workspace}
}

func buildWorkspace(channels []domain.NotificationChannel, tasks []server.TaskResponse, groups []domain.Group, deliveries []domain.NotificationDelivery, loadedAt time.Time) Workspace {
	groupByID := make(map[string]domain.Group, len(groups))
	for _, group := range groups {
		groupByID[group.ID] = group
	}
	groupPath := func(id string) string {
		parts, seen := []string{}, map[string]bool{}
		for id != "" && !seen[id] {
			seen[id] = true
			group, ok := groupByID[id]
			if !ok {
				break
			}
			parts = append([]string{group.Name}, parts...)
			id = group.ParentID
		}
		return strings.Join(parts, " / ")
	}
	w := Workspace{Channels: make([]Channel, 0, len(channels)), Tasks: make([]Scope, 0, len(tasks)), Groups: make([]Scope, 0, len(groups)), Coverage: []ConfiguredScope{}, CoverageComplete: true, Deliveries: make([]Delivery, 0, len(deliveries)), LoadedAt: timestamp(loadedAt)}
	for _, value := range channels {
		w.Channels = append(w.Channels, Channel{ID: value.ID, Name: fallback(value.Name, "Unnamed channel"), Kind: string(value.Kind), EndpointSummary: value.EndpointSummary, HasAuthorization: value.HasAuthorization, Enabled: value.Enabled, UpdatedAt: timestamp(value.UpdatedAt)})
	}
	for _, value := range groups {
		w.Groups = append(w.Groups, Scope{Type: "group", ID: value.ID, Name: fallback(value.Name, "Unnamed group"), Context: groupPath(value.ID)})
	}
	for _, value := range tasks {
		w.Tasks = append(w.Tasks, Scope{Type: "task", ID: value.Task.ID, Name: fallback(value.Task.Name, "Unnamed task"), Context: fallback(groupPath(value.Task.GroupID), "Ungrouped")})
	}
	for _, value := range deliveries {
		w.Deliveries = append(w.Deliveries, mapDelivery(value))
	}
	sort.SliceStable(w.Channels, func(i, j int) bool { return strings.ToLower(w.Channels[i].Name) < strings.ToLower(w.Channels[j].Name) })
	sort.SliceStable(w.Groups, func(i, j int) bool {
		return strings.ToLower(w.Groups[i].Context) < strings.ToLower(w.Groups[j].Context)
	})
	sort.SliceStable(w.Tasks, func(i, j int) bool { return strings.ToLower(w.Tasks[i].Name) < strings.ToLower(w.Tasks[j].Name) })
	sort.SliceStable(w.Deliveries, func(i, j int) bool { return w.Deliveries[i].CreatedAt > w.Deliveries[j].CreatedAt })
	return w
}

func (s *Service) configuredCoverage(ctx context.Context, workspace *Workspace) ([]ConfiguredScope, bool) {
	enabled := make(map[string]bool, len(workspace.Channels))
	for _, channel := range workspace.Channels {
		enabled[channel.ID] = channel.Enabled
	}
	coverage := make([]ConfiguredScope, 0, len(workspace.Tasks)+len(workspace.Groups))
	complete := true
	appendSummary := func(scope Scope, sourceType, sourceID string, values []domain.NotificationAssignment) {
		if len(values) == 0 {
			return
		}
		summary := ConfiguredScope{Type: scope.Type, ID: scope.ID, Name: scope.Name, Context: scope.Context, SourceType: sourceType, SourceName: sourceName(workspace, domain.NotificationScopeType(sourceType), sourceID, scope.Name)}
		seen := map[string]bool{}
		for _, value := range values {
			summary.OnSuccess = summary.OnSuccess || value.OnSuccess
			summary.OnFailure = summary.OnFailure || value.OnFailure
			if seen[value.ChannelID] {
				continue
			}
			seen[value.ChannelID] = true
			summary.DestinationCount++
			if enabled[value.ChannelID] {
				summary.EnabledDestinationCount++
			}
		}
		coverage = append(coverage, summary)
	}
	for _, scope := range workspace.Groups {
		values, err := s.backend.ListGroupNotificationAssignments(ctx, scope.ID)
		if err != nil {
			complete = false
			continue
		}
		appendSummary(scope, string(domain.NotificationScopeGroup), scope.ID, values)
	}
	for _, scope := range workspace.Tasks {
		policy, err := s.backend.EffectiveTaskNotificationPolicy(ctx, scope.ID)
		if err != nil {
			complete = false
			continue
		}
		appendSummary(scope, string(policy.SourceScopeType), policy.SourceScopeID, policy.Assignments)
	}
	sort.SliceStable(coverage, func(i, j int) bool {
		if coverage[i].Type != coverage[j].Type {
			return coverage[i].Type < coverage[j].Type
		}
		return strings.ToLower(coverage[i].Context+coverage[i].Name) < strings.ToLower(coverage[j].Context+coverage[j].Name)
	})
	return coverage, complete
}

func mapDelivery(value domain.NotificationDelivery) Delivery {
	kind := "task_outcome"
	if value.EventKind == domain.NotificationEventTest {
		kind = "test"
	}
	state := string(value.State)
	switch value.State {
	case domain.NotificationDeliveryPending:
		if value.Attempts > 0 {
			state = "retrying"
		} else {
			state = "queued"
		}
	case domain.NotificationDeliveryClaimed:
		state = "sending"
	case domain.NotificationDeliverySucceeded:
		state = "successful"
	case domain.NotificationDeliveryFailed:
		state = "failed"
	}
	return Delivery{ID: value.ID, ChannelID: value.ChannelID, ChannelName: fallback(value.ChannelName, "Removed channel"), DestinationSummary: value.DestinationSummary, Kind: kind, TaskID: value.TaskID, TaskName: value.TaskName, GroupID: value.GroupID, GroupName: value.GroupName, RunID: value.RunID, State: state, Attempts: value.Attempts, NextAttemptAt: timestamp(value.NextAttemptAt), CreatedAt: timestamp(value.CreatedAt), ClaimedAt: optionalTimestamp(value.ClaimedAt), CompletedAt: optionalTimestamp(value.CompletedAt), LastStatus: value.LastStatus, LastError: value.LastError}
}

func (s *Service) SaveChannel(ctx context.Context, draft ChannelDraft) Result {
	name := strings.TrimSpace(draft.Name)
	if name == "" {
		return rejected("save_notification_channel", "name", "Enter a channel name.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	entityID := draft.ID
	if draft.IsNew || entityID == "" {
		endpoint := strings.TrimSpace(draft.Endpoint)
		if _, err := webhook.ValidateEndpoint(endpoint); err != nil {
			return rejected("save_notification_channel", "endpoint", "Enter a valid HTTPS webhook endpoint.")
		}
		created, err := s.backend.CreateNotificationChannel(c, server.NotificationChannelCreateRequest{Name: name, Endpoint: endpoint, Authorization: draft.Authorization})
		if err != nil {
			return failure("save_notification_channel", err)
		}
		entityID = created.ID
	} else {
		req := server.NotificationChannelUpdateRequest{Name: &name}
		if draft.ReplaceEndpoint {
			endpoint := strings.TrimSpace(draft.Endpoint)
			if _, err := webhook.ValidateEndpoint(endpoint); err != nil {
				return rejected("save_notification_channel", "endpoint", "Enter a valid HTTPS replacement endpoint.")
			}
			req.Endpoint = &endpoint
		}
		if _, err := s.backend.UpdateNotificationChannel(c, entityID, req); err != nil {
			return failure("save_notification_channel", err)
		}
		if draft.ReplaceAuthorization {
			if _, err := s.backend.RotateNotificationChannelAuthorization(c, entityID, draft.Authorization); err != nil {
				return failure("save_notification_channel", err)
			}
		}
	}
	return s.refreshed(ctx, "save_notification_channel", "Channel saved.", entityID)
}

func (s *Service) SetChannelEnabled(ctx context.Context, id string, enabled bool) Result {
	message := "Channel disabled."
	if enabled {
		message = "Channel enabled."
	}
	return s.mutateChannel(ctx, "toggle_notification_channel", message, id, func(c context.Context) error {
		_, err := s.backend.SetNotificationChannelEnabled(c, id, enabled)
		return err
	})
}
func (s *Service) TestChannel(ctx context.Context, id string) Result {
	return s.mutateChannel(ctx, "test_notification_channel", "Test notification queued.", id, func(c context.Context) error { _, err := s.backend.TestNotificationChannel(c, id); return err })
}
func (s *Service) DeleteChannel(ctx context.Context, id string) Result {
	return s.mutateChannel(ctx, "delete_notification_channel", "Channel removed.", id, func(c context.Context) error { return s.backend.DeleteNotificationChannel(c, id) })
}

func (s *Service) mutateChannel(ctx context.Context, action, message, id string, mutate func(context.Context) error) Result {
	id = strings.TrimSpace(id)
	if id == "" {
		return rejected(action, "id", "Choose a channel.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	err := mutate(c)
	cancel()
	if err != nil {
		return failure(action, err)
	}
	return s.refreshed(ctx, action, message, id)
}
func (s *Service) refreshed(ctx context.Context, action, message, id string) Result {
	loaded := s.Workspace(ctx)
	if loaded.Outcome != "accepted" {
		return Result{Action: action, Outcome: "accepted", Message: message + " Refresh Notifications to see the latest state.", EntityID: id}
	}
	loaded.Action = action
	loaded.Message = message
	loaded.EntityID = id
	return loaded
}

func (s *Service) Policy(ctx context.Context, scopeType, scopeID string) Result {
	workspace, err := s.policyScopes(ctx)
	if err != nil {
		return failure("load_notification_policy", err)
	}
	scope, ok := findScope(workspace, scopeType, scopeID)
	if !ok {
		return rejected("load_notification_policy", "scopeId", "Choose an available task or group.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	var direct []domain.NotificationAssignment
	err = nil
	policy := Policy{Scope: scope, DirectAssignments: []Assignment{}, EffectiveAssignments: []Assignment{}, EffectiveSourceType: "none", EffectiveSourceName: "No policy"}
	if scopeType == "task" {
		direct, err = s.backend.ListTaskNotificationAssignments(c, scopeID)
		if err == nil {
			var effective domain.EffectiveNotificationPolicy
			effective, err = s.backend.EffectiveTaskNotificationPolicy(c, scopeID)
			if err == nil {
				policy.EffectiveSourceType = string(effective.SourceScopeType)
				policy.EffectiveSourceID = effective.SourceScopeID
				policy.EffectiveSourceName = sourceName(workspace, effective.SourceScopeType, effective.SourceScopeID, scope.Name)
				policy.EffectiveAssignments = mapAssignments(effective.Assignments)
			}
		}
	} else if scopeType == "group" {
		direct, err = s.backend.ListGroupNotificationAssignments(c, scopeID)
		policy.EffectiveSourceType = "group"
		policy.EffectiveSourceID = scopeID
		policy.EffectiveSourceName = scope.Name
	} else {
		return rejected("load_notification_policy", "scopeType", "Choose a task or group.")
	}
	if err != nil {
		return failure("load_notification_policy", err)
	}
	policy.DirectAssignments = mapAssignments(direct)
	if scopeType == "group" {
		policy.EffectiveAssignments = policy.DirectAssignments
	}
	return Result{Action: "load_notification_policy", Outcome: "accepted", Message: "Notification policy loaded.", EntityID: scopeID, Policy: &policy}
}

func (s *Service) policyScopes(ctx context.Context) (*Workspace, error) {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	tasks, err := s.backend.ListTaskDetails(c, "", "")
	if err != nil {
		return nil, err
	}
	groups, err := s.backend.ListGroups(c)
	if err != nil {
		return nil, err
	}
	workspace := buildWorkspace(nil, tasks, groups, nil, s.now())
	return &workspace, nil
}

func (s *Service) SavePolicy(ctx context.Context, draft PolicyDraft) Result {
	if draft.ScopeType != "task" && draft.ScopeType != "group" {
		return rejected("save_notification_policy", "scopeType", "Choose a task or group.")
	}
	if strings.TrimSpace(draft.ScopeID) == "" {
		return rejected("save_notification_policy", "scopeId", "Choose an available task or group.")
	}
	seen := map[string]bool{}
	req := server.NotificationAssignmentsRequest{Assignments: make([]server.NotificationAssignmentInput, 0, len(draft.Assignments))}
	for _, item := range draft.Assignments {
		id := strings.TrimSpace(item.ChannelID)
		if id == "" || seen[id] || (!item.OnSuccess && !item.OnFailure) {
			return rejected("save_notification_policy", "assignments", "Each selected channel must be unique and notify on failure, success, or both.")
		}
		seen[id] = true
		req.Assignments = append(req.Assignments, server.NotificationAssignmentInput{ChannelID: id, OnSuccess: item.OnSuccess, OnFailure: item.OnFailure})
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	var err error
	if draft.ScopeType == "task" {
		_, err = s.backend.ReplaceTaskNotificationAssignments(c, draft.ScopeID, req)
	} else {
		_, err = s.backend.ReplaceGroupNotificationAssignments(c, draft.ScopeID, req)
	}
	cancel()
	if err != nil {
		return failure("save_notification_policy", err)
	}
	result := s.Policy(ctx, draft.ScopeType, draft.ScopeID)
	if result.Outcome != "accepted" {
		return Result{Action: "save_notification_policy", Outcome: "accepted", Message: "Notification policy saved. Reload it to see the latest state.", EntityID: draft.ScopeID}
	}
	result.Action = "save_notification_policy"
	result.Message = "Notification policy saved."
	return result
}

func findScope(workspace *Workspace, scopeType, id string) (Scope, bool) {
	values := workspace.Tasks
	if scopeType == "group" {
		values = workspace.Groups
	} else {
		values = workspace.Tasks
	}
	for _, value := range values {
		if value.ID == id {
			return value, true
		}
	}
	return Scope{}, false
}
func sourceName(workspace *Workspace, kind domain.NotificationScopeType, id, taskName string) string {
	if kind == domain.NotificationScopeNone {
		return "No policy"
	}
	if kind == domain.NotificationScopeTask {
		return taskName
	}
	scope, ok := findScope(workspace, "group", id)
	if ok {
		return scope.Name
	}
	return "Configured group"
}
func mapAssignments(values []domain.NotificationAssignment) []Assignment {
	out := make([]Assignment, 0, len(values))
	for _, value := range values {
		out = append(out, Assignment{ChannelID: value.ChannelID, OnSuccess: value.OnSuccess, OnFailure: value.OnFailure})
	}
	return out
}
func timestamp(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}
func optionalTimestamp(value *time.Time) string {
	if value == nil {
		return ""
	}
	return timestamp(*value)
}
func fallback(value, otherwise string) string {
	if strings.TrimSpace(value) == "" {
		return otherwise
	}
	return value
}
func rejected(action, field, message string) Result {
	return Result{Action: action, Outcome: "rejected", Field: field, Message: message}
}
func failure(action string, err error) Result {
	var status *client.StatusError
	if errors.As(err, &status) {
		switch status.Code {
		case server.CodeValidation:
			return Result{Action: action, Outcome: "rejected", Field: status.Field, Message: "Review the highlighted notification details and try again."}
		case server.CodeConflict:
			return Result{Action: action, Outcome: "conflict", Field: status.Field, Message: "Notification state changed or conflicts with this request. Refresh and try again."}
		case server.CodeNotFound:
			return Result{Action: action, Outcome: "stale", Message: "That notification resource no longer exists. Refresh Notifications."}
		}
	}
	return Result{Action: action, Outcome: "unavailable", Message: "The local scheduler service is unavailable. Try again."}
}
