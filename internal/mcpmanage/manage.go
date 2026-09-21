// Package mcpmanage exposes bounded automation-definition mutations to an
// explicitly authorized local MCP session.
package mcpmanage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	Permission    = "manage"
	cacheLimit    = 256
	cacheLifetime = 10 * time.Minute
	callTimeout   = 10 * time.Second
)

const (
	OutcomeAccepted  = "accepted"
	OutcomeRejected  = "rejected"
	OutcomeDenied    = "denied"
	OutcomeUncertain = "uncertain"
)

type manageClient interface {
	VerifyAccess(context.Context) (domain.Capability, error)
	CreateTask(context.Context, server.TaskCreateRequest) (server.TaskResponse, error)
	UpdateTask(context.Context, string, server.TaskUpdateRequest) (server.TaskResponse, error)
	DeleteTask(context.Context, string) error
	CreateGroup(context.Context, server.GroupCreateRequest) (domain.Group, error)
	UpdateGroup(context.Context, string, server.GroupUpdateRequest) (domain.Group, error)
	DeleteGroup(context.Context, string) error
	CreateChain(context.Context, server.ChainCreateRequest) (domain.CompletionChain, error)
	UpdateChain(context.Context, string, server.ChainUpdateRequest) (domain.CompletionChain, error)
	DeleteChain(context.Context, string) error
	CreateTrigger(context.Context, server.TriggerCreateRequest) (server.TriggerSecretResponse, error)
	UpdateTrigger(context.Context, string, server.TriggerUpdateRequest) (server.TriggerResponse, error)
	DeleteTrigger(context.Context, string) error
	CreateFilesystemWatcher(context.Context, server.FilesystemWatcherCreateRequest) (server.FilesystemWatcherResponse, error)
	UpdateFilesystemWatcher(context.Context, string, server.FilesystemWatcherUpdateRequest) (server.FilesystemWatcherResponse, error)
	DeleteFilesystemWatcher(context.Context, string) error
	ReplaceTaskNotificationAssignments(context.Context, string, server.NotificationAssignmentsRequest) ([]domain.NotificationAssignment, error)
	ReplaceGroupNotificationAssignments(context.Context, string, server.NotificationAssignmentsRequest) ([]domain.NotificationAssignment, error)
}

type Envelope struct {
	DaemonID  string `json:"daemon_id" jsonschema:"the exact daemon installation UUID"`
	RequestID string `json:"request_id" jsonschema:"a caller-generated UUID reused only for retries of this exact mutation"`
	Confirmed bool   `json:"confirmed,omitempty" jsonschema:"explicit caller confirmation when the connection policy requires it"`
}

type TaskCreate struct {
	Name              string     `json:"name"`
	GroupID           string     `json:"group_id,omitempty"`
	Command           string     `json:"command"`
	Args              []string   `json:"args,omitempty"`
	WorkingDir        string     `json:"working_dir,omitempty"`
	RunAs             string     `json:"run_as,omitempty"`
	Timezone          string     `json:"timezone,omitempty"`
	Schedule          string     `json:"schedule,omitempty"`
	ScheduleSyntax    string     `json:"schedule_syntax,omitempty"`
	At                *time.Time `json:"at,omitempty"`
	OverlapPolicy     string     `json:"overlap_policy,omitempty"`
	CatchupPolicy     string     `json:"catchup_policy,omitempty"`
	MissingDatePolicy string     `json:"missing_date_policy,omitempty"`
	TimeBasis         string     `json:"time_basis,omitempty"`
	DSTGapPolicy      string     `json:"dst_gap_policy,omitempty"`
	DSTOverlapPolicy  string     `json:"dst_overlap_policy,omitempty"`
	Enabled           *bool      `json:"enabled,omitempty"`
}

type TaskUpdate struct {
	Name              string     `json:"name,omitempty"`
	GroupID           *string    `json:"group_id,omitempty"`
	Command           string     `json:"command,omitempty"`
	Args              []string   `json:"args,omitempty"`
	WorkingDir        string     `json:"working_dir,omitempty"`
	RunAs             string     `json:"run_as,omitempty"`
	Timezone          string     `json:"timezone,omitempty"`
	Schedule          string     `json:"schedule,omitempty"`
	ScheduleSyntax    string     `json:"schedule_syntax,omitempty"`
	At                *time.Time `json:"at,omitempty"`
	OverlapPolicy     string     `json:"overlap_policy,omitempty"`
	CatchupPolicy     string     `json:"catchup_policy,omitempty"`
	MissingDatePolicy string     `json:"missing_date_policy,omitempty"`
	TimeBasis         string     `json:"time_basis,omitempty"`
	DSTGapPolicy      string     `json:"dst_gap_policy,omitempty"`
	DSTOverlapPolicy  string     `json:"dst_overlap_policy,omitempty"`
	ClearName         bool       `json:"clear_name,omitempty"`
	ClearCommand      bool       `json:"clear_command,omitempty"`
	ClearWorkingDir   bool       `json:"clear_working_dir,omitempty"`
	ClearRunAs        bool       `json:"clear_run_as,omitempty"`
	ClearSchedule     bool       `json:"clear_schedule,omitempty"`
	ClearSensitive    bool       `json:"clear_sensitive_inputs,omitempty" jsonschema:"clear stored environment and stdin without reading them"`
}

type TaskInput struct {
	Envelope
	Action   string      `json:"action" jsonschema:"create, update, or delete"`
	ObjectID string      `json:"object_id,omitempty"`
	Create   *TaskCreate `json:"create,omitempty"`
	Update   *TaskUpdate `json:"update,omitempty"`
}

type GroupInput struct {
	Envelope
	Action   string                     `json:"action" jsonschema:"create, update, or delete"`
	ObjectID string                     `json:"object_id,omitempty"`
	Create   *server.GroupCreateRequest `json:"create,omitempty"`
	Update   *server.GroupUpdateRequest `json:"update,omitempty"`
}

type ChainInput struct {
	Envelope
	Action   string                     `json:"action" jsonschema:"create, update, or delete"`
	ObjectID string                     `json:"object_id,omitempty"`
	Create   *server.ChainCreateRequest `json:"create,omitempty"`
	Update   *server.ChainUpdateRequest `json:"update,omitempty"`
}

type TriggerInput struct {
	Envelope
	Action   string                       `json:"action" jsonschema:"create, update, or delete"`
	ObjectID string                       `json:"object_id,omitempty"`
	Create   *server.TriggerCreateRequest `json:"create,omitempty"`
	Update   *server.TriggerUpdateRequest `json:"update,omitempty"`
}

type WatcherInput struct {
	Envelope
	Action   string                                 `json:"action" jsonschema:"create, update, or delete"`
	ObjectID string                                 `json:"object_id,omitempty"`
	Create   *server.FilesystemWatcherCreateRequest `json:"create,omitempty"`
	Update   *server.FilesystemWatcherUpdateRequest `json:"update,omitempty"`
}

type AssignmentInput struct {
	Envelope
	Scope       string                               `json:"scope" jsonschema:"task or group"`
	ObjectID    string                               `json:"object_id"`
	Assignments []server.NotificationAssignmentInput `json:"assignments"`
}

type Result struct {
	SchemaVersion string `json:"schema_version"`
	Permission    string `json:"permission"`
	Operation     string `json:"operation"`
	DaemonID      string `json:"daemon_id"`
	RequestID     string `json:"request_id"`
	ObjectKind    string `json:"object_kind"`
	ObjectID      string `json:"object_id,omitempty"`
	Outcome       string `json:"outcome"`
	Message       string `json:"message"`
}

type cached struct {
	fingerprint string
	result      Result
	expires     time.Time
}

// Executor owns one session's confirmation policy and retry cache.
type Executor struct {
	mu                  sync.Mutex
	client              manageClient
	requireConfirmation bool
	now                 func() time.Time
	cache               map[string]cached
}

func New(c manageClient, requireConfirmation bool) *Executor {
	return &Executor{client: c, requireConfirmation: requireConfirmation, now: time.Now, cache: make(map[string]cached)}
}

// AddTools adds exactly six bounded Manage tools.
func AddTools(s *mcp.Server, e *Executor) {
	closed, destructive := false, true
	annotations := func(title string) *mcp.ToolAnnotations {
		return &mcp.ToolAnnotations{Title: title, ReadOnlyHint: false, DestructiveHint: &destructive, IdempotentHint: false, OpenWorldHint: &closed}
	}
	mcp.AddTool(s, &mcp.Tool{Name: "tasks_manage", Title: "Manage one task definition", Description: "Create, update, or delete one task through the scheduler API. Environment and stdin values are intentionally unavailable.", Annotations: annotations("Manage task")}, e.taskHandler())
	mcp.AddTool(s, &mcp.Tool{Name: "groups_manage", Title: "Manage one task group", Description: "Create, update, or delete one task group through the scheduler API.", Annotations: annotations("Manage group")}, e.groupHandler())
	mcp.AddTool(s, &mcp.Tool{Name: "chains_manage", Title: "Manage one completion chain", Description: "Create, update, or delete one completion chain through the scheduler API.", Annotations: annotations("Manage completion chain")}, e.chainHandler())
	mcp.AddTool(s, &mcp.Tool{Name: "triggers_manage", Title: "Manage one external trigger", Description: "Create, update, or delete one trigger. Generated trigger credentials are never returned.", Annotations: annotations("Manage trigger")}, e.triggerHandler())
	mcp.AddTool(s, &mcp.Tool{Name: "watchers_manage", Title: "Manage one filesystem watcher", Description: "Create, update, or delete one watcher through the scheduler API.", Annotations: annotations("Manage watcher")}, e.watcherHandler())
	mcp.AddTool(s, &mcp.Tool{Name: "notification_assignments_replace", Title: "Replace notification assignments", Description: "Atomically replace direct notification assignments for one task or group.", Annotations: annotations("Replace notification assignments")}, e.assignmentHandler())
}

type mutation func(context.Context) (string, error)

func (e *Executor) execute(ctx context.Context, envelope Envelope, operation, kind, objectID, fingerprint, successMessage string, invoke mutation) Result {
	result := Result{SchemaVersion: SchemaVersion, Permission: Permission, Operation: operation, DaemonID: envelope.DaemonID, RequestID: envelope.RequestID, ObjectKind: kind, ObjectID: objectID}
	if message := validateEnvelope(envelope, e.requireConfirmation); message != "" {
		result.Outcome, result.Message = OutcomeRejected, message
		return result
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.prune()
	callCtx, cancel := context.WithTimeout(client.ExpectDaemon(ctx, envelope.DaemonID), callTimeout)
	defer cancel()
	if prior, ok := e.cache[envelope.RequestID]; ok {
		capability, err := e.client.VerifyAccess(callCtx)
		if err == nil && capability.Allows(domain.CapabilityManage) {
			if prior.fingerprint == fingerprint {
				return prior.result
			}
			result.Outcome, result.Message = OutcomeRejected, "request_id was already used for a different mutation"
			return result
		}
		if err != nil && !isAuthorizationDenial(err) {
			result.Outcome, result.Message = classify(err)
			return result
		}
	}
	resolvedID, err := invoke(callCtx)
	if resolvedID != "" {
		result.ObjectID = resolvedID
	}
	result.Outcome, result.Message = classify(err)
	if err == nil && successMessage != "" {
		result.Message = successMessage
	}
	if len(e.cache) >= cacheLimit {
		e.evictOldest()
	}
	e.cache[envelope.RequestID] = cached{fingerprint: fingerprint, result: result, expires: e.now().Add(cacheLifetime)}
	return result
}

func validateEnvelope(input Envelope, requireConfirmation bool) string {
	if _, err := uuid.Parse(input.DaemonID); err != nil {
		return "daemon_id must be a UUID"
	}
	if _, err := uuid.Parse(input.RequestID); err != nil {
		return "request_id must be a UUID"
	}
	if requireConfirmation && !input.Confirmed {
		return "this Manage session requires confirmed=true for every mutation"
	}
	return ""
}

func validateAction(action, objectID string, create, update bool) string {
	if objectID != "" && !validObjectID(objectID) {
		return "object_id must contain no more than 256 non-control characters"
	}
	switch action {
	case "create":
		if objectID != "" || !create || update {
			return "create requires only a create definition and no object_id"
		}
	case "update":
		if objectID == "" || create || !update {
			return "update requires object_id and only an update definition"
		}
	case "delete":
		if objectID == "" || create || update {
			return "delete requires object_id and no definition"
		}
	default:
		return "action must be create, update, or delete"
	}
	return ""
}

func validObjectID(objectID string) bool {
	return objectID != "" && len(objectID) <= 256 && strings.IndexFunc(objectID, unicode.IsControl) < 0
}

func fingerprint(value any) string {
	raw, _ := json.Marshal(value)
	digest := sha256.Sum256(raw)
	return hex.EncodeToString(digest[:])
}

func rejected(envelope Envelope, operation, kind, objectID, message string) Result {
	return Result{SchemaVersion: SchemaVersion, Permission: Permission, Operation: operation, DaemonID: envelope.DaemonID, RequestID: envelope.RequestID, ObjectKind: kind, ObjectID: objectID, Outcome: OutcomeRejected, Message: message}
}

func classify(err error) (string, string) {
	if err == nil {
		return OutcomeAccepted, "Definition mutation accepted."
	}
	var status *client.StatusError
	if errors.As(err, &status) {
		if status.Code == server.CodeForbidden {
			return OutcomeDenied, "The MCP session is not authorized for this operation."
		}
		if status.Code == server.CodeAuditUnavailable || status.Code == server.CodeInternal {
			return OutcomeUncertain, "The daemon could not prove the final mutation outcome. Inspect current state before retrying."
		}
		return OutcomeRejected, safeMessage(status.Message)
	}
	var uncertain *client.MutationUncertainError
	var connection *client.ConnectionError
	if errors.As(err, &uncertain) || errors.As(err, &connection) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return OutcomeUncertain, "The daemon response was not authoritative. Inspect current state before retrying."
	}
	return OutcomeUncertain, "The daemon could not prove the final mutation outcome. Inspect current state before retrying."
}

func safeMessage(message string) string {
	message = strings.TrimSpace(message)
	if message == "" || len(message) > 512 || strings.IndexFunc(message, unicode.IsControl) >= 0 {
		return "The daemon rejected the mutation."
	}
	return message
}

func isAuthorizationDenial(err error) bool {
	var status *client.StatusError
	return errors.As(err, &status) && status.Code == server.CodeForbidden
}

func (e *Executor) prune() {
	now := e.now()
	for id, entry := range e.cache {
		if !entry.expires.After(now) {
			delete(e.cache, id)
		}
	}
}

func (e *Executor) evictOldest() {
	var id string
	var expiry time.Time
	for candidate, entry := range e.cache {
		if id == "" || entry.expires.Before(expiry) {
			id, expiry = candidate, entry.expires
		}
	}
	delete(e.cache, id)
}
