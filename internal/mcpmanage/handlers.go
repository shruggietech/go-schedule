package mcpmanage

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func toolResult(result Result) (*mcp.CallToolResult, Result, error) {
	return &mcp.CallToolResult{IsError: result.Outcome != OutcomeAccepted}, result, nil
}

func (e *Executor) taskHandler() mcp.ToolHandlerFor[TaskInput, Result] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input TaskInput) (*mcp.CallToolResult, Result, error) {
		operation := "tasks." + input.Action
		if message := validateAction(input.Action, input.ObjectID, input.Create != nil, input.Update != nil); message != "" {
			return toolResult(rejected(input.Envelope, operation, "task", input.ObjectID, message))
		}
		invoke := func(ctx context.Context) (string, error) {
			switch input.Action {
			case "create":
				response, err := e.client.CreateTask(ctx, input.Create.request())
				return response.Task.ID, err
			case "update":
				response, err := e.client.UpdateTask(ctx, input.ObjectID, input.Update.request())
				return response.Task.ID, err
			default:
				return input.ObjectID, e.client.DeleteTask(ctx, input.ObjectID)
			}
		}
		return toolResult(e.execute(ctx, input.Envelope, operation, "task", input.ObjectID, fingerprint(input), "Task definition mutation accepted.", invoke))
	}
}

func (e *Executor) groupHandler() mcp.ToolHandlerFor[GroupInput, Result] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input GroupInput) (*mcp.CallToolResult, Result, error) {
		operation := "groups." + input.Action
		if message := validateAction(input.Action, input.ObjectID, input.Create != nil, input.Update != nil); message != "" {
			return toolResult(rejected(input.Envelope, operation, "group", input.ObjectID, message))
		}
		invoke := func(ctx context.Context) (string, error) {
			switch input.Action {
			case "create":
				value, err := e.client.CreateGroup(ctx, *input.Create)
				return value.ID, err
			case "update":
				value, err := e.client.UpdateGroup(ctx, input.ObjectID, *input.Update)
				return value.ID, err
			default:
				return input.ObjectID, e.client.DeleteGroup(ctx, input.ObjectID)
			}
		}
		return toolResult(e.execute(ctx, input.Envelope, operation, "group", input.ObjectID, fingerprint(input), "Group definition mutation accepted.", invoke))
	}
}

func (e *Executor) chainHandler() mcp.ToolHandlerFor[ChainInput, Result] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input ChainInput) (*mcp.CallToolResult, Result, error) {
		operation := "chains." + input.Action
		if message := validateAction(input.Action, input.ObjectID, input.Create != nil, input.Update != nil); message != "" {
			return toolResult(rejected(input.Envelope, operation, "chain", input.ObjectID, message))
		}
		invoke := func(ctx context.Context) (string, error) {
			switch input.Action {
			case "create":
				value, err := e.client.CreateChain(ctx, *input.Create)
				return value.ID, err
			case "update":
				value, err := e.client.UpdateChain(ctx, input.ObjectID, *input.Update)
				return value.ID, err
			default:
				return input.ObjectID, e.client.DeleteChain(ctx, input.ObjectID)
			}
		}
		return toolResult(e.execute(ctx, input.Envelope, operation, "chain", input.ObjectID, fingerprint(input), "Completion-chain mutation accepted.", invoke))
	}
}

func (e *Executor) triggerHandler() mcp.ToolHandlerFor[TriggerInput, Result] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input TriggerInput) (*mcp.CallToolResult, Result, error) {
		operation := "triggers." + input.Action
		if message := validateAction(input.Action, input.ObjectID, input.Create != nil, input.Update != nil); message != "" {
			return toolResult(rejected(input.Envelope, operation, "trigger", input.ObjectID, message))
		}
		invoke := func(ctx context.Context) (string, error) {
			switch input.Action {
			case "create":
				value, err := e.client.CreateTrigger(ctx, *input.Create)
				return value.Trigger.ID, err
			case "update":
				value, err := e.client.UpdateTrigger(ctx, input.ObjectID, *input.Update)
				return value.ID, err
			default:
				return input.ObjectID, e.client.DeleteTrigger(ctx, input.ObjectID)
			}
		}
		message := "Trigger definition mutation accepted."
		if input.Action == "create" {
			message = "Trigger created. Its generated credential was withheld from MCP output; use an Enroll-authorized human client for credential handling."
		}
		return toolResult(e.execute(ctx, input.Envelope, operation, "trigger", input.ObjectID, fingerprint(input), message, invoke))
	}
}

func (e *Executor) watcherHandler() mcp.ToolHandlerFor[WatcherInput, Result] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input WatcherInput) (*mcp.CallToolResult, Result, error) {
		operation := "watchers." + input.Action
		if message := validateAction(input.Action, input.ObjectID, input.Create != nil, input.Update != nil); message != "" {
			return toolResult(rejected(input.Envelope, operation, "watcher", input.ObjectID, message))
		}
		invoke := func(ctx context.Context) (string, error) {
			switch input.Action {
			case "create":
				value, err := e.client.CreateFilesystemWatcher(ctx, *input.Create)
				return value.ID, err
			case "update":
				value, err := e.client.UpdateFilesystemWatcher(ctx, input.ObjectID, *input.Update)
				return value.ID, err
			default:
				return input.ObjectID, e.client.DeleteFilesystemWatcher(ctx, input.ObjectID)
			}
		}
		return toolResult(e.execute(ctx, input.Envelope, operation, "watcher", input.ObjectID, fingerprint(input), "Filesystem-watcher mutation accepted.", invoke))
	}
}

func (e *Executor) assignmentHandler() mcp.ToolHandlerFor[AssignmentInput, Result] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input AssignmentInput) (*mcp.CallToolResult, Result, error) {
		operation := "notification_assignments.replace"
		if input.Scope != "task" && input.Scope != "group" {
			return toolResult(rejected(input.Envelope, operation, "notification_assignment", input.ObjectID, "scope must be task or group"))
		}
		if !validObjectID(input.ObjectID) {
			return toolResult(rejected(input.Envelope, operation, "notification_assignment", input.ObjectID, "object_id must contain no more than 256 non-control characters"))
		}
		invoke := func(ctx context.Context) (string, error) {
			request := server.NotificationAssignmentsRequest{Assignments: input.Assignments}
			if input.Scope == "task" {
				_, err := e.client.ReplaceTaskNotificationAssignments(ctx, input.ObjectID, request)
				return input.ObjectID, err
			}
			_, err := e.client.ReplaceGroupNotificationAssignments(ctx, input.ObjectID, request)
			return input.ObjectID, err
		}
		return toolResult(e.execute(ctx, input.Envelope, operation, input.Scope+"_notification_assignment", input.ObjectID, fingerprint(input), "Notification assignments replaced atomically.", invoke))
	}
}

func (input TaskCreate) request() server.TaskCreateRequest {
	return server.TaskCreateRequest{Name: input.Name, GroupID: input.GroupID, Command: input.Command, Args: input.Args, WorkingDir: input.WorkingDir, RunAs: input.RunAs, Timezone: input.Timezone, Schedule: input.Schedule, ScheduleSyntax: input.ScheduleSyntax, At: input.At, OverlapPolicy: input.OverlapPolicy, CatchupPolicy: input.CatchupPolicy, MissingDatePolicy: input.MissingDatePolicy, TimeBasis: input.TimeBasis, DSTGapPolicy: input.DSTGapPolicy, DSTOverlapPolicy: input.DSTOverlapPolicy, Enabled: input.Enabled}
}

func (input TaskUpdate) request() server.TaskUpdateRequest {
	request := server.TaskUpdateRequest{Name: input.Name, GroupID: input.GroupID, Command: input.Command, Args: input.Args, WorkingDir: input.WorkingDir, RunAs: input.RunAs, Timezone: input.Timezone, Schedule: input.Schedule, ScheduleSyntax: input.ScheduleSyntax, At: input.At, OverlapPolicy: input.OverlapPolicy, CatchupPolicy: input.CatchupPolicy, MissingDatePolicy: input.MissingDatePolicy, TimeBasis: input.TimeBasis, DSTGapPolicy: input.DSTGapPolicy, DSTOverlapPolicy: input.DSTOverlapPolicy, ClearName: input.ClearName, ClearCommand: input.ClearCommand, ClearWorkingDir: input.ClearWorkingDir, ClearRunAs: input.ClearRunAs, ClearSchedule: input.ClearSchedule}
	if input.ClearSensitive {
		request.Env = map[string]string{}
		empty := ""
		request.Stdin = &empty
	}
	return request
}
