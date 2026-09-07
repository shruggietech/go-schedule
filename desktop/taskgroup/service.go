package taskgroup

import (
	"context"
	"errors"
	"maps"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/commandline"
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
	groups, err := s.backend.ListGroups(c)
	if err != nil {
		return failure("load", err)
	}
	workspace := buildWorkspace(tasks, groups)
	return accepted("load", "Tasks and groups are up to date.", "", &workspace, nil)
}

func (s *Service) Task(ctx context.Context, id string) OperationResult {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	response, err := s.backend.GetTask(c, id)
	if err != nil {
		return failure("load", err)
	}
	detail := detailFrom(response)
	return accepted("load", "Task loaded.", id, nil, &detail)
}

func (s *Service) PreviewTask(ctx context.Context, draft TaskDraft) OperationResult {
	invocation, err := parseDraftCommand(draft.CommandLine)
	if err != nil {
		return rejected("preview", "command", "Enter a valid direct command line.")
	}
	result := OperationResult{Action: "preview", Outcome: "accepted", Message: "Command is valid.", Command: &CommandPreview{Program: invocation.Program, Args: invocation.Args}}
	detail := draft.TaskDetail
	result.Task = &detail
	if draft.Mode == "recurring" && strings.TrimSpace(draft.Schedule) != "" {
		c, cancel := context.WithTimeout(ctx, callTimeout)
		defer cancel()
		preview, err := s.backend.Preview(c, previewRequest(draft))
		if err != nil {
			return failure("preview", err)
		}
		result.Message = preview.HumanSummary
		detail.ScheduleSummary = preview.HumanSummary
		detail.PolicySummary = preview.PolicySummary
		detail.NextRuns = times(preview.NextRuns)
	}
	return result
}

func (s *Service) SaveTask(ctx context.Context, draft TaskDraft) OperationResult {
	invocation, err := parseDraftCommand(draft.CommandLine)
	if err != nil {
		return rejected("save_task", "command", "Enter a valid direct command line.")
	}
	env, field := environment(draft.Environment)
	if field != "" {
		return rejected("save_task", field, "Environment keys must be nonempty and unique.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	var response server.TaskResponse
	if draft.IsNew || draft.ID == "" {
		req, err := createRequest(draft, invocation, env)
		if err != nil {
			return rejected("save_task", "at", "Enter a valid future date and time.")
		}
		response, err = s.backend.CreateTask(c, req)
	} else {
		current, loadErr := s.backend.GetTask(c, draft.ID)
		if loadErr != nil {
			return failure("save_task", loadErr)
		}
		if !draft.OverwriteStale && draft.OriginalUpdatedAt != "" && current.Task.UpdatedAt.Format(time.RFC3339Nano) != draft.OriginalUpdatedAt {
			return OperationResult{Action: "save_task", Outcome: "stale", Message: "This task changed after you opened it. Reload it or explicitly overwrite the newer version.", EntityID: draft.ID}
		}
		req, err := updateRequest(draft, current, invocation, env)
		if err != nil {
			return rejected("save_task", "at", "Enter a valid future date and time.")
		}
		response, err = s.backend.UpdateTask(c, draft.ID, req)
		if err == nil && response.Task.Enabled != draft.Enabled {
			if err = s.backend.SetTaskEnabled(c, draft.ID, draft.Enabled); err == nil {
				response, err = s.backend.GetTask(c, draft.ID)
			}
		}
	}
	if err != nil {
		return failure("save_task", err)
	}
	detail := detailFrom(response)
	workspaceResult := s.Workspace(ctx)
	return accepted("save_task", "Task saved.", response.Task.ID, workspaceResult.Workspace, &detail)
}

func (s *Service) RunTask(ctx context.Context, id string) OperationResult {
	return s.mutateTask(ctx, "run_task", id, func(c context.Context) error { return s.backend.RunNow(c, id) })
}
func (s *Service) SetTaskEnabled(ctx context.Context, id string, enabled bool) OperationResult {
	return s.mutateTask(ctx, "toggle_task", id, func(c context.Context) error { return s.backend.SetTaskEnabled(c, id, enabled) })
}
func (s *Service) DeleteTask(ctx context.Context, id string) OperationResult {
	return s.mutateTask(ctx, "delete_task", id, func(c context.Context) error { return s.backend.DeleteTask(c, id) })
}

func (s *Service) mutateTask(ctx context.Context, action, id string, mutate func(context.Context) error) OperationResult {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	if err := mutate(c); err != nil {
		return failure(action, err)
	}
	workspace := s.Workspace(ctx)
	return accepted(action, "Task action accepted.", id, workspace.Workspace, nil)
}

func (s *Service) SaveGroup(ctx context.Context, draft GroupDraft) OperationResult {
	name := strings.TrimSpace(draft.Name)
	if name == "" {
		return rejected("save_group", "name", "Enter a group name.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	entityID := draft.ID
	if !draft.IsNew && draft.ID != "" {
		groups, err := s.backend.ListGroups(c)
		if err != nil {
			return failure("save_group", err)
		}
		var current *domain.Group
		for i := range groups {
			group := groups[i]
			if group.ID == draft.ID && !draft.OverwriteStale && draft.OriginalUpdatedAt != "" && group.UpdatedAt.Format(time.RFC3339Nano) != draft.OriginalUpdatedAt {
				return OperationResult{Action: "save_group", Outcome: "stale", Message: "This group changed after you opened it. Reload it or explicitly overwrite the newer version.", EntityID: draft.ID}
			}
			if group.ID == draft.ID {
				current = &groups[i]
			}
		}
		if current == nil {
			return rejected("save_group", "", "The selected group no longer exists.")
		}
		parent := draft.ParentID
		request := server.GroupUpdateRequest{}
		if name != current.Name {
			request.Name = name
		}
		if draft.ParentID != current.ParentID {
			request.Parent = &parent
		}
		group, err := s.backend.UpdateGroup(c, draft.ID, request)
		if err != nil {
			return failure("save_group", err)
		}
		if group.Enabled != draft.Enabled {
			if err = s.backend.SetGroupEnabled(c, group.ID, draft.Enabled); err != nil {
				return failure("save_group", err)
			}
		}
	} else {
		group, err := s.backend.CreateGroup(c, server.GroupCreateRequest{Name: name, ParentID: draft.ParentID, Enabled: &draft.Enabled})
		if err != nil {
			return failure("save_group", err)
		}
		entityID = group.ID
	}
	workspace := s.Workspace(ctx)
	return accepted("save_group", "Group saved.", entityID, workspace.Workspace, nil)
}

func (s *Service) SetGroupEnabled(ctx context.Context, id string, enabled bool) OperationResult {
	return s.mutateGroup(ctx, "toggle_group", id, func(c context.Context) error { return s.backend.SetGroupEnabled(c, id, enabled) })
}
func (s *Service) DeleteGroup(ctx context.Context, id string) OperationResult {
	return s.mutateGroup(ctx, "delete_group", id, func(c context.Context) error { return s.backend.DeleteGroup(c, id) })
}
func (s *Service) mutateGroup(ctx context.Context, action, id string, mutate func(context.Context) error) OperationResult {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	if err := mutate(c); err != nil {
		return failure(action, err)
	}
	workspace := s.Workspace(ctx)
	return accepted(action, "Group action accepted.", id, workspace.Workspace, nil)
}

func buildWorkspace(details []server.TaskResponse, groups []domain.Group) Workspace {
	byID := tasklogic.ByID(groups)
	paths := groupPaths(groups)
	groupSummaries := make([]GroupSummary, 0, len(groups))
	for _, group := range groups {
		descendants := tasklogic.DescendantIDs(group.ID, groups)
		validHierarchy := validGroupChain(group.ID, byID)
		summary := GroupSummary{ID: group.ID, Name: group.Name, ParentID: group.ParentID, Path: paths[group.ID], Depth: groupDepth(group.ID, byID), DeclaredEnabled: group.Enabled, EffectiveEnabled: validHierarchy && tasklogic.ChainEnabled(group.ID, byID), DescendantCount: len(descendants), UpdatedAt: group.UpdatedAt.Format(time.RFC3339Nano)}
		for _, candidate := range groups {
			if candidate.ParentID == group.ID {
				summary.ChildCount++
			}
		}
		if !summary.EffectiveEnabled {
			if !validHierarchy {
				summary.EffectiveReason = "The group hierarchy is invalid."
			} else if blocker, ok := tasklogic.NearestDisabledGroup(group.ID, byID); ok {
				summary.EffectiveReason = paths[blocker.ID] + " is disabled."
			} else {
				summary.EffectiveReason = "The group hierarchy is invalid."
			}
		}
		groupSummaries = append(groupSummaries, summary)
	}
	tasks := make([]TaskSummary, 0, len(details))
	for _, detail := range details {
		value := detail.Task
		summary := TaskSummary{ID: value.ID, Name: tasklogic.DisplayName(value), GroupID: value.GroupID, GroupPath: paths[value.GroupID], CommandConfigured: detail.Readiness.CommandReady, DeclaredEnabled: value.Enabled, Lifecycle: string(value.State), Timezone: value.Timezone, PolicySummary: detail.PolicySummary, UpdatedAt: value.UpdatedAt.Format(time.RFC3339Nano), NextRuns: times(detail.NextRuns)}
		if summary.GroupPath == "" {
			summary.GroupPath = "Not assigned"
		}
		if detail.Schedule != nil {
			summary.ScheduleSummary = detail.Schedule.HumanSummary
		} else {
			summary.ScheduleSummary = "Manual only"
		}
		summary.EffectiveState, summary.EffectiveReason = effectiveTask(detail, byID, paths)
		tasks = append(tasks, summary)
		for i := range groupSummaries {
			if groupSummaries[i].ID == value.GroupID {
				groupSummaries[i].TaskCount++
			}
		}
	}
	sort.SliceStable(groupSummaries, func(i, j int) bool { return groupSummaries[i].Path < groupSummaries[j].Path })
	return Workspace{Tasks: tasks, Groups: groupSummaries, LoadedAt: time.Now().UTC().Format(time.RFC3339Nano)}
}

func groupPaths(groups []domain.Group) map[string]string {
	byID := tasklogic.ByID(groups)
	out := map[string]string{}
	for _, group := range groups {
		names, seen := []string{}, map[string]bool{}
		for id := group.ID; id != ""; {
			if seen[id] {
				names = []string{"Invalid hierarchy"}
				break
			}
			seen[id] = true
			value, ok := byID[id]
			if !ok {
				names = append([]string{"Missing group"}, names...)
				break
			}
			names = append([]string{groupLabel(value, groups)}, names...)
			id = value.ParentID
		}
		out[group.ID] = strings.Join(names, " / ")
	}
	return out
}
func groupLabel(group domain.Group, groups []domain.Group) string {
	duplicates := 0
	for _, candidate := range groups {
		if candidate.ParentID == group.ParentID && candidate.Name == group.Name {
			duplicates++
		}
	}
	if duplicates < 2 {
		return group.Name
	}
	id := group.ID
	if len(id) > 8 {
		id = id[:8]
	}
	return group.Name + " [" + id + "]"
}
func validGroupChain(id string, groups map[string]domain.Group) bool {
	seen := map[string]bool{}
	for id != "" {
		if seen[id] {
			return false
		}
		seen[id] = true
		group, ok := groups[id]
		if !ok {
			return false
		}
		id = group.ParentID
	}
	return true
}
func groupDepth(id string, groups map[string]domain.Group) int {
	depth, seen := 0, map[string]bool{}
	for {
		group, ok := groups[id]
		if !ok || group.ParentID == "" || seen[id] {
			return depth
		}
		seen[id] = true
		depth++
		id = group.ParentID
	}
}
func effectiveTask(detail server.TaskResponse, groups map[string]domain.Group, paths map[string]string) (string, string) {
	task := detail.Task
	if task.GroupID != "" {
		if _, ok := groups[task.GroupID]; !ok {
			return "invalid_group", "The assigned group no longer exists."
		}
		if !validGroupChain(task.GroupID, groups) {
			return "invalid_group", "The group hierarchy is invalid."
		}
		if blocker, ok := tasklogic.NearestDisabledGroup(task.GroupID, groups); ok {
			return "group_disabled", paths[blocker.ID] + " is disabled."
		}
	}
	switch {
	case task.State != domain.TaskActive:
		return "terminal", detail.Readiness.Reason
	case !detail.Readiness.CommandReady:
		return "not_runnable", detail.Readiness.Reason
	case !detail.Readiness.ActivationReady:
		return "manual_only", detail.Readiness.Reason
	case !task.Enabled:
		return "task_disabled", "Automatic activation is disabled."
	default:
		return "runnable", "Ready to run."
	}
}
func detailFrom(response server.TaskResponse) TaskDetail {
	task := response.Task
	command, _ := commandline.Format(task.Command, task.Args)
	detail := TaskDetail{ID: task.ID, Name: task.Name, GroupID: task.GroupID, CommandLine: command, WorkingDir: task.WorkingDir, Stdin: task.Stdin, RunAs: task.RunAs, Enabled: task.Enabled, Timezone: task.Timezone, Mode: "manual", OverlapPolicy: string(task.OverlapPolicy), CatchupPolicy: string(task.CatchupPolicy), MissingDatePolicy: string(task.MissingDatePolicy), TimeBasis: string(task.TimeBasis), DSTGapPolicy: string(task.DSTGapPolicy), DSTOverlapPolicy: string(task.DSTOverlapPolicy), PolicySummary: response.PolicySummary, Readiness: string(response.Readiness.Status), UpdatedAt: task.UpdatedAt.Format(time.RFC3339Nano), NextRuns: times(response.NextRuns)}
	keys := make([]string, 0, len(task.Env))
	for key := range task.Env {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		detail.Environment = append(detail.Environment, EnvironmentRow{Key: key, Value: task.Env[key]})
	}
	if response.Schedule != nil {
		detail.ScheduleSummary = response.Schedule.HumanSummary
		if response.Schedule.Kind == domain.ScheduleOneOff {
			detail.Mode = "one_off"
			if response.Schedule.RunAt != nil {
				detail.At = response.Schedule.RunAt.Format(time.RFC3339)
			}
		} else {
			detail.Mode = "recurring"
			detail.Schedule = response.Schedule.Expression
			detail.ScheduleSyntax = response.Schedule.SourceSyntax
		}
	}
	return detail
}
func environment(rows []EnvironmentRow) (map[string]string, string) {
	values := map[string]string{}
	for _, row := range rows {
		key := strings.TrimSpace(row.Key)
		if key == "" {
			return nil, "environment"
		}
		if _, exists := values[key]; exists {
			return nil, "environment"
		}
		values[key] = row.Value
	}
	return values, ""
}

func parseDraftCommand(value string) (commandline.Invocation, error) {
	if strings.TrimSpace(value) == "" {
		return commandline.Invocation{}, nil
	}
	return commandline.Parse(value)
}
func previewRequest(d TaskDraft) server.PreviewRequest {
	return server.PreviewRequest{Schedule: d.Schedule, ScheduleSyntax: d.ScheduleSyntax, Timezone: d.Timezone, MissingDatePolicy: d.MissingDatePolicy, TimeBasis: d.TimeBasis, DSTGapPolicy: d.DSTGapPolicy, DSTOverlapPolicy: d.DSTOverlapPolicy}
}
func createRequest(d TaskDraft, i commandline.Invocation, env map[string]string) (server.TaskCreateRequest, error) {
	req := server.TaskCreateRequest{Name: d.Name, GroupID: d.GroupID, Command: i.Program, Args: i.Args, WorkingDir: d.WorkingDir, Env: env, Stdin: d.Stdin, RunAs: d.RunAs, Timezone: d.Timezone, OverlapPolicy: d.OverlapPolicy, CatchupPolicy: d.CatchupPolicy, MissingDatePolicy: d.MissingDatePolicy, TimeBasis: d.TimeBasis, DSTGapPolicy: d.DSTGapPolicy, DSTOverlapPolicy: d.DSTOverlapPolicy, Enabled: &d.Enabled}
	if d.Mode == "recurring" {
		req.Schedule, req.ScheduleSyntax = d.Schedule, d.ScheduleSyntax
	} else if d.Mode == "one_off" {
		value, err := time.Parse(time.RFC3339, d.At)
		if err != nil {
			return req, err
		}
		req.At = &value
	}
	return req, nil
}
func updateRequest(d TaskDraft, current server.TaskResponse, i commandline.Invocation, env map[string]string) (server.TaskUpdateRequest, error) {
	req := server.TaskUpdateRequest{}
	if d.Name == "" && current.Task.Name != "" {
		req.ClearName = true
	} else if d.Name != current.Task.Name {
		req.Name = d.Name
	}
	if i.Program == "" && current.Task.Command != "" {
		req.ClearCommand = true
	} else if i.Program != current.Task.Command {
		req.Command = i.Program
	}
	if !slices.Equal(i.Args, current.Task.Args) {
		req.Args = i.Args
	}
	if !maps.Equal(env, current.Task.Env) {
		req.Env = env
	}
	if d.Stdin != current.Task.Stdin {
		req.Stdin = &d.Stdin
	}
	if d.GroupID != current.Task.GroupID {
		req.GroupID = &d.GroupID
	}
	if d.WorkingDir == "" {
		req.ClearWorkingDir = current.Task.WorkingDir != ""
	} else if d.WorkingDir != current.Task.WorkingDir {
		req.WorkingDir = d.WorkingDir
	}
	if d.RunAs == "" {
		req.ClearRunAs = current.Task.RunAs != ""
	} else if d.RunAs != current.Task.RunAs {
		req.RunAs = d.RunAs
	}
	if d.Timezone != current.Task.Timezone {
		req.Timezone = d.Timezone
	}
	if d.OverlapPolicy != string(current.Task.OverlapPolicy) {
		req.OverlapPolicy = d.OverlapPolicy
	}
	if d.CatchupPolicy != string(current.Task.CatchupPolicy) {
		req.CatchupPolicy = d.CatchupPolicy
	}
	if d.MissingDatePolicy != string(current.Task.MissingDatePolicy) {
		req.MissingDatePolicy = d.MissingDatePolicy
	}
	if d.TimeBasis != string(current.Task.TimeBasis) {
		req.TimeBasis = d.TimeBasis
	}
	if d.DSTGapPolicy != string(current.Task.DSTGapPolicy) {
		req.DSTGapPolicy = d.DSTGapPolicy
	}
	if d.DSTOverlapPolicy != string(current.Task.DSTOverlapPolicy) {
		req.DSTOverlapPolicy = d.DSTOverlapPolicy
	}
	if d.Mode == "manual" {
		req.ClearSchedule = current.Schedule != nil
	} else if d.Mode == "one_off" {
		value, err := time.Parse(time.RFC3339, d.At)
		if err != nil {
			return req, err
		}
		if current.Schedule == nil || current.Schedule.RunAt == nil || !current.Schedule.RunAt.Equal(value) {
			req.At = &value
		}
	} else if current.Schedule == nil || d.Schedule != current.Schedule.Expression || d.ScheduleSyntax != current.Schedule.SourceSyntax || d.Timezone != current.Task.Timezone {
		req.Schedule, req.ScheduleSyntax = d.Schedule, d.ScheduleSyntax
	}
	return req, nil
}
func times(values []time.Time) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, value.Format(time.RFC3339))
	}
	return out
}
func accepted(action, message, id string, workspace *Workspace, task *TaskDetail) OperationResult {
	return OperationResult{Action: action, Outcome: "accepted", Message: message, EntityID: id, Workspace: workspace, Task: task}
}
func rejected(action, field, message string) OperationResult {
	return OperationResult{Action: action, Outcome: "rejected", Message: message, Field: field}
}
func failure(action string, err error) OperationResult {
	var status *client.StatusError
	if errors.As(err, &status) {
		field := status.Field
		known := map[string]bool{"name": true, "group_id": true, "command": true, "working_dir": true, "environment": true, "stdin": true, "run_as": true, "timezone": true, "schedule": true, "schedule_syntax": true, "at": true, "overlap_policy": true, "catchup_policy": true, "missing_date_policy": true, "time_basis": true, "dst_gap_policy": true, "dst_overlap_policy": true, "enabled": true, "parent_id": true}
		if !known[field] {
			field = ""
		}
		message := "The request was rejected. Review the highlighted field and try again."
		if status.Code == server.CodeNotFound {
			message = "The selected item no longer exists."
		}
		return OperationResult{Action: action, Outcome: "rejected", Message: message, Field: field}
	}
	return OperationResult{Action: action, Outcome: "unavailable", Message: "The local scheduler service is unavailable. Try again."}
}
