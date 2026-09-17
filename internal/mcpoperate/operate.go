// Package mcpoperate adds narrowly scoped task controls to an explicitly
// authorized MCP server.
package mcpoperate

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

const (
	SchemaVersion = "1"
	Permission    = "operate"
	cacheLimit    = 256
	cacheLifetime = 10 * time.Minute
	callTimeout   = 10 * time.Second
)

type operationClient interface {
	SetTaskEnabled(context.Context, string, bool) error
	RunNow(context.Context, string) error
	VerifyAccess(context.Context) (domain.Capability, error)
}

type Input struct {
	DaemonID  string `json:"daemon_id" jsonschema:"the exact daemon installation UUID"`
	TaskID    string `json:"task_id" jsonschema:"the exact existing task identifier"`
	RequestID string `json:"request_id" jsonschema:"a caller-generated UUID reused only for retries of this exact operation"`
}

type Result struct {
	SchemaVersion string `json:"schema_version"`
	Permission    string `json:"permission"`
	Operation     string `json:"operation"`
	DaemonID      string `json:"daemon_id"`
	TaskID        string `json:"task_id"`
	RequestID     string `json:"request_id"`
	Outcome       string `json:"outcome"`
	Message       string `json:"message"`
}

const (
	OutcomeAccepted  = "accepted"
	OutcomeRejected  = "rejected"
	OutcomeDenied    = "denied"
	OutcomeUncertain = "uncertain"
)

type cached struct {
	fingerprint string
	result      Result
	expires     time.Time
}

// Executor owns retry deduplication for one runtime MCP session.
type Executor struct {
	mu     sync.Mutex
	client operationClient
	now    func() time.Time
	cache  map[string]cached
}

func New(client operationClient) *Executor {
	return &Executor{client: client, now: time.Now, cache: make(map[string]cached)}
}

// AddTools exposes exactly the three issue #178 operations.
func AddTools(mcpServer *mcp.Server, executor *Executor) {
	closed := false
	destructive := true
	annotations := func(title string, idempotent bool) *mcp.ToolAnnotations {
		return &mcp.ToolAnnotations{Title: title, ReadOnlyHint: false, DestructiveHint: &destructive, IdempotentHint: idempotent, OpenWorldHint: &closed}
	}
	mcp.AddTool(mcpServer, &mcp.Tool{Name: "tasks_run_now", Title: "Run an existing task now", Description: "Request one immediate run on the exact named daemon and task. Reuse request_id only when retrying this exact call.", Annotations: annotations("Run task now", false)}, executor.handler("tasks.run_now"))
	mcp.AddTool(mcpServer, &mcp.Tool{Name: "tasks_enable", Title: "Enable an existing task", Description: "Enable one existing task without changing its definition.", Annotations: annotations("Enable task", true)}, executor.handler("tasks.enable"))
	mcp.AddTool(mcpServer, &mcp.Tool{Name: "tasks_disable", Title: "Disable an existing task", Description: "Disable one existing task without changing or deleting its definition.", Annotations: annotations("Disable task", true)}, executor.handler("tasks.disable"))
}

func (e *Executor) handler(operation string) mcp.ToolHandlerFor[Input, Result] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input Input) (*mcp.CallToolResult, Result, error) {
		result := e.execute(ctx, operation, input)
		return &mcp.CallToolResult{IsError: result.Outcome != OutcomeAccepted}, result, nil
	}
}

func (e *Executor) execute(ctx context.Context, operation string, input Input) Result {
	base := Result{SchemaVersion: SchemaVersion, Permission: Permission, Operation: operation, DaemonID: input.DaemonID, TaskID: input.TaskID, RequestID: input.RequestID}
	if message := validate(input); message != "" {
		base.Outcome, base.Message = OutcomeRejected, message
		return base
	}
	fingerprint := operation + "\x00" + input.DaemonID + "\x00" + input.TaskID
	e.mu.Lock()
	defer e.mu.Unlock()
	e.prune()
	callCtx, cancel := context.WithTimeout(client.ExpectDaemon(ctx, input.DaemonID), callTimeout)
	defer cancel()
	if previous, ok := e.cache[input.RequestID]; ok {
		capability, accessErr := e.client.VerifyAccess(callCtx)
		if accessErr == nil && capability.Allows(domain.CapabilityOperate) {
			if previous.fingerprint == fingerprint {
				return previous.result
			}
			base.Outcome, base.Message = OutcomeRejected, "request_id was already used for a different operation or target"
			return base
		}
		if accessErr != nil && !isAuthorizationDenial(accessErr) {
			base.Outcome, base.Message = classify(operation, accessErr)
			return base
		}
		base.Outcome, base.Message = classify(operation, e.invoke(callCtx, operation, input.TaskID))
		return base
	}
	base.Outcome, base.Message = classify(operation, e.invoke(callCtx, operation, input.TaskID))
	if len(e.cache) >= cacheLimit {
		e.evictOldest()
	}
	e.cache[input.RequestID] = cached{fingerprint: fingerprint, result: base, expires: e.now().Add(cacheLifetime)}
	return base
}

func (e *Executor) invoke(ctx context.Context, operation, taskID string) error {
	switch operation {
	case "tasks.run_now":
		return e.client.RunNow(ctx, taskID)
	case "tasks.enable":
		return e.client.SetTaskEnabled(ctx, taskID, true)
	case "tasks.disable":
		return e.client.SetTaskEnabled(ctx, taskID, false)
	default:
		return errors.New("unsupported operation")
	}
}

func isAuthorizationDenial(err error) bool {
	var status *client.StatusError
	return errors.As(err, &status) && status.Code == server.CodeForbidden
}

func validate(input Input) string {
	if _, err := uuid.Parse(input.DaemonID); err != nil {
		return "daemon_id must be a UUID"
	}
	if _, err := uuid.Parse(input.RequestID); err != nil {
		return "request_id must be a UUID"
	}
	if input.TaskID == "" || len(input.TaskID) > 256 || strings.IndexFunc(input.TaskID, unicode.IsControl) >= 0 {
		return "task_id must contain 1 through 256 non-control characters"
	}
	return ""
}

func classify(operation string, err error) (string, string) {
	if err == nil {
		switch operation {
		case "tasks.run_now":
			return OutcomeAccepted, "Run request accepted."
		case "tasks.enable":
			return OutcomeAccepted, "Task enabled."
		default:
			return OutcomeAccepted, "Task disabled."
		}
	}
	var status *client.StatusError
	if errors.As(err, &status) {
		if status.Code == server.CodeForbidden {
			return OutcomeDenied, "The MCP session is not authorized for this operation."
		}
		if status.Code == server.CodeAuditUnavailable || status.Code == server.CodeInternal {
			return OutcomeUncertain, "The daemon could not prove the final operation outcome. Inspect current task state and activity before retrying."
		}
		return OutcomeRejected, safeMessage(status.Message)
	}
	var uncertain *client.MutationUncertainError
	var connection *client.ConnectionError
	if errors.As(err, &uncertain) || errors.As(err, &connection) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return OutcomeUncertain, "The daemon response was not authoritative. Inspect current task state and activity before retrying."
	}
	return OutcomeUncertain, "The daemon could not prove the final operation outcome. Inspect current task state and activity before retrying."
}

func safeMessage(message string) string {
	message = strings.TrimSpace(message)
	if message == "" || len(message) > 512 || strings.IndexFunc(message, unicode.IsControl) >= 0 {
		return "The daemon rejected the operation."
	}
	return message
}

func (e *Executor) prune() {
	now := e.now()
	for id, value := range e.cache {
		if !value.expires.After(now) {
			delete(e.cache, id)
		}
	}
}

func (e *Executor) evictOldest() {
	var oldestID string
	var oldest time.Time
	for id, value := range e.cache {
		if oldestID == "" || value.expires.Before(oldest) {
			oldestID, oldest = id, value.expires
		}
	}
	delete(e.cache, oldestID)
}

func (r Result) String() string {
	return fmt.Sprintf("%s %s for task %s on daemon %s", r.Outcome, r.Operation, r.TaskID, r.DaemonID)
}
