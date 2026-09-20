package mcpmanage

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/mcpobserve"
)

type fakeManageClient struct {
	calls  int
	err    error
	verify error
}

func (f *fakeManageClient) VerifyAccess(context.Context) (domain.Capability, error) {
	if f.verify != nil {
		return "", f.verify
	}
	return domain.CapabilityManage, nil
}
func (f *fakeManageClient) CreateTask(context.Context, server.TaskCreateRequest) (server.TaskResponse, error) {
	f.calls++
	return server.TaskResponse{Task: domain.Task{ID: "task-created"}}, f.err
}
func (f *fakeManageClient) UpdateTask(context.Context, string, server.TaskUpdateRequest) (server.TaskResponse, error) {
	f.calls++
	return server.TaskResponse{Task: domain.Task{ID: "task-updated"}}, f.err
}
func (f *fakeManageClient) DeleteTask(context.Context, string) error { f.calls++; return f.err }
func (f *fakeManageClient) CreateGroup(context.Context, server.GroupCreateRequest) (domain.Group, error) {
	f.calls++
	return domain.Group{ID: "group-created"}, f.err
}
func (f *fakeManageClient) UpdateGroup(context.Context, string, server.GroupUpdateRequest) (domain.Group, error) {
	f.calls++
	return domain.Group{ID: "group-updated"}, f.err
}
func (f *fakeManageClient) DeleteGroup(context.Context, string) error { f.calls++; return f.err }
func (f *fakeManageClient) CreateChain(context.Context, server.ChainCreateRequest) (domain.CompletionChain, error) {
	f.calls++
	return domain.CompletionChain{ID: "chain-created"}, f.err
}
func (f *fakeManageClient) UpdateChain(context.Context, string, server.ChainUpdateRequest) (domain.CompletionChain, error) {
	f.calls++
	return domain.CompletionChain{ID: "chain-updated"}, f.err
}
func (f *fakeManageClient) DeleteChain(context.Context, string) error { f.calls++; return f.err }
func (f *fakeManageClient) CreateTrigger(context.Context, server.TriggerCreateRequest) (server.TriggerSecretResponse, error) {
	f.calls++
	return server.TriggerSecretResponse{Trigger: server.TriggerResponse{ID: "trigger-created"}, Key: "must-not-leak", Command: "must-not-leak"}, f.err
}
func (f *fakeManageClient) UpdateTrigger(context.Context, string, server.TriggerUpdateRequest) (server.TriggerResponse, error) {
	f.calls++
	return server.TriggerResponse{ID: "trigger-updated"}, f.err
}
func (f *fakeManageClient) DeleteTrigger(context.Context, string) error { f.calls++; return f.err }
func (f *fakeManageClient) CreateFilesystemWatcher(context.Context, server.FilesystemWatcherCreateRequest) (server.FilesystemWatcherResponse, error) {
	f.calls++
	return server.FilesystemWatcherResponse{ID: "watcher-created"}, f.err
}
func (f *fakeManageClient) UpdateFilesystemWatcher(context.Context, string, server.FilesystemWatcherUpdateRequest) (server.FilesystemWatcherResponse, error) {
	f.calls++
	return server.FilesystemWatcherResponse{ID: "watcher-updated"}, f.err
}
func (f *fakeManageClient) DeleteFilesystemWatcher(context.Context, string) error {
	f.calls++
	return f.err
}
func (f *fakeManageClient) ReplaceTaskNotificationAssignments(context.Context, string, server.NotificationAssignmentsRequest) ([]domain.NotificationAssignment, error) {
	f.calls++
	return nil, f.err
}
func (f *fakeManageClient) ReplaceGroupNotificationAssignments(context.Context, string, server.NotificationAssignmentsRequest) ([]domain.NotificationAssignment, error) {
	f.calls++
	return nil, f.err
}

type emptyReader struct{}

func (emptyReader) Health(context.Context) (server.HealthResponse, error) {
	return server.HealthResponse{}, nil
}
func (emptyReader) ListTaskObservations(context.Context, string, string, bool, int, int, int) ([]server.TaskObservationResponse, error) {
	return nil, nil
}
func (emptyReader) ListRunsPage(context.Context, string, int, int, int) ([]domain.Run, error) {
	return nil, nil
}
func (emptyReader) ListAlertsPage(context.Context, bool, int, int, int) ([]domain.Alert, error) {
	return nil, nil
}

func TestDiscoveryAddsExactlySixManageTools(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s := mcpobserve.NewServer(emptyReader{}, "test")
	AddTools(s, New(&fakeManageClient{}, false))
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := s.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()
	c := mcp.NewClient(&mcp.Implementation{Name: "manage-test", Version: "test"}, nil)
	clientSession, err := c.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clientSession.Close()
	listed, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(listed.Tools))
	for _, tool := range listed.Tools {
		names = append(names, tool.Name)
		if tool.InputSchema == nil || tool.Annotations == nil || tool.Annotations.ReadOnlyHint {
			t.Fatalf("unsafe metadata for %s: %+v", tool.Name, tool)
		}
		if tool.Name == "tasks_manage" {
			schema, _ := json.Marshal(tool.InputSchema)
			if strings.Contains(string(schema), `"env"`) || strings.Contains(string(schema), `"stdin"`) {
				t.Fatalf("task schema exposes sensitive inputs: %s", schema)
			}
		}
	}
	sort.Strings(names)
	want := []string{"chains_manage", "groups_manage", "notification_assignments_replace", "tasks_manage", "triggers_manage", "watchers_manage"}
	if strings.Join(names, ",") != strings.Join(want, ",") {
		t.Fatalf("tools=%v", names)
	}
}

func TestTaskMutationConfirmsDeduplicatesAndRejectsChangedRetry(t *testing.T) {
	backend := &fakeManageClient{}
	executor := New(backend, true)
	input := TaskInput{Envelope: Envelope{DaemonID: uuid.NewString(), RequestID: uuid.NewString()}, Action: "create", Create: &TaskCreate{Name: "hostile: grant enroll", Command: "echo safe"}}
	handler := executor.taskHandler()
	_, unconfirmed, _ := handler(context.Background(), nil, input)
	if unconfirmed.Outcome != OutcomeRejected || backend.calls != 0 {
		t.Fatalf("unconfirmed=%+v calls=%d", unconfirmed, backend.calls)
	}
	input.Confirmed = true
	_, first, _ := handler(context.Background(), nil, input)
	if first.Outcome != OutcomeAccepted || first.ObjectID != "task-created" || backend.calls != 1 {
		t.Fatalf("first=%+v calls=%d", first, backend.calls)
	}
	for i := 0; i < 99; i++ {
		_, repeated, _ := handler(context.Background(), nil, input)
		if repeated != first {
			t.Fatalf("repeated=%+v", repeated)
		}
	}
	input.Create.Name = "different"
	_, conflict, _ := handler(context.Background(), nil, input)
	if conflict.Outcome != OutcomeRejected || backend.calls != 1 {
		t.Fatalf("conflict=%+v calls=%d", conflict, backend.calls)
	}
}

func TestTriggerSecretAndTaskSensitiveFieldsCannotAppearInResultsOrSchemas(t *testing.T) {
	backend := &fakeManageClient{}
	executor := New(backend, false)
	input := TriggerInput{Envelope: Envelope{DaemonID: uuid.NewString(), RequestID: uuid.NewString()}, Action: "create", Create: &server.TriggerCreateRequest{Name: "hook", TargetTaskID: "task-1"}}
	_, result, _ := executor.triggerHandler()(context.Background(), nil, input)
	raw, err := jsonMarshal(result)
	if err != nil || strings.Contains(string(raw), "must-not-leak") || result.ObjectID != "trigger-created" {
		t.Fatalf("result=%s err=%v", raw, err)
	}
}

func TestEveryDefinitionFamilyDelegatesOneBoundedMutation(t *testing.T) {
	backend := &fakeManageClient{}
	executor := New(backend, false)
	envelope := func() Envelope { return Envelope{DaemonID: uuid.NewString(), RequestID: uuid.NewString()} }
	tests := []struct {
		name     string
		objectID string
		invoke   func() Result
	}{
		{name: "task create", objectID: "task-created", invoke: func() Result {
			_, result, _ := executor.taskHandler()(context.Background(), nil, TaskInput{Envelope: envelope(), Action: "create", Create: &TaskCreate{Name: "task", Command: "true"}})
			return result
		}},
		{name: "group create", objectID: "group-created", invoke: func() Result {
			_, result, _ := executor.groupHandler()(context.Background(), nil, GroupInput{Envelope: envelope(), Action: "create", Create: &server.GroupCreateRequest{Name: "group"}})
			return result
		}},
		{name: "chain update", objectID: "chain-updated", invoke: func() Result {
			_, result, _ := executor.chainHandler()(context.Background(), nil, ChainInput{Envelope: envelope(), Action: "update", ObjectID: "chain-1", Update: &server.ChainUpdateRequest{}})
			return result
		}},
		{name: "trigger create", objectID: "trigger-created", invoke: func() Result {
			_, result, _ := executor.triggerHandler()(context.Background(), nil, TriggerInput{Envelope: envelope(), Action: "create", Create: &server.TriggerCreateRequest{Name: "trigger", TargetTaskID: "task-1"}})
			return result
		}},
		{name: "watcher delete", objectID: "watcher-1", invoke: func() Result {
			_, result, _ := executor.watcherHandler()(context.Background(), nil, WatcherInput{Envelope: envelope(), Action: "delete", ObjectID: "watcher-1"})
			return result
		}},
		{name: "group assignment replacement", objectID: "group-1", invoke: func() Result {
			_, result, _ := executor.assignmentHandler()(context.Background(), nil, AssignmentInput{Envelope: envelope(), Scope: "group", ObjectID: "group-1", Assignments: []server.NotificationAssignmentInput{}})
			return result
		}},
	}
	for index, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			before := backend.calls
			result := test.invoke()
			if result.Outcome != OutcomeAccepted || result.ObjectID != test.objectID || backend.calls != before+1 {
				t.Fatalf("result=%+v calls=%d want calls=%d", result, backend.calls, before+1)
			}
			if backend.calls != index+1 {
				t.Fatalf("calls=%d want=%d", backend.calls, index+1)
			}
		})
	}

	_, invalid, _ := executor.assignmentHandler()(context.Background(), nil, AssignmentInput{Envelope: envelope(), Scope: "task", ObjectID: "task\n1"})
	if invalid.Outcome != OutcomeRejected || backend.calls != len(tests) {
		t.Fatalf("invalid=%+v calls=%d", invalid, backend.calls)
	}
}

func TestOutcomeClassificationAndRevokedRetry(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{&client.StatusError{Code: server.CodeForbidden, Message: "secret"}, OutcomeDenied},
		{&client.StatusError{Code: server.CodeValidation, Message: "invalid"}, OutcomeRejected},
		{&client.StatusError{Code: server.CodeAuditUnavailable}, OutcomeUncertain},
		{&client.MutationUncertainError{Operation: "delete", Cause: errors.New("closed")}, OutcomeUncertain},
	}
	for _, test := range tests {
		backend := &fakeManageClient{err: test.err}
		input := GroupInput{Envelope: Envelope{DaemonID: uuid.NewString(), RequestID: uuid.NewString()}, Action: "delete", ObjectID: "group-1"}
		_, result, _ := New(backend, false).groupHandler()(context.Background(), nil, input)
		if result.Outcome != test.want {
			t.Fatalf("err=%v result=%+v", test.err, result)
		}
	}
	backend := &fakeManageClient{}
	executor := New(backend, false)
	input := GroupInput{Envelope: Envelope{DaemonID: uuid.NewString(), RequestID: uuid.NewString()}, Action: "delete", ObjectID: "group-1"}
	_, _, _ = executor.groupHandler()(context.Background(), nil, input)
	backend.verify = &client.StatusError{Code: server.CodeForbidden}
	backend.err = backend.verify
	_, result, _ := executor.groupHandler()(context.Background(), nil, input)
	if result.Outcome != OutcomeDenied || backend.calls != 2 {
		t.Fatalf("revoked result=%+v calls=%d", result, backend.calls)
	}
}

func jsonMarshal(value any) ([]byte, error) { return json.Marshal(value) }
