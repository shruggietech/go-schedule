package taskgroup

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	tasklogic "github.com/shruggietech/go-schedule/internal/task"
)

type fakeBackend struct {
	tasks   []server.TaskResponse
	groups  []domain.Group
	created server.TaskCreateRequest
	updated server.TaskUpdateRequest
	err     error
}

func (f *fakeBackend) ListTaskDetails(context.Context, string, string) ([]server.TaskResponse, error) {
	return f.tasks, f.err
}
func (f *fakeBackend) GetTask(_ context.Context, id string) (server.TaskResponse, error) {
	for _, task := range f.tasks {
		if task.Task.ID == id {
			return task, f.err
		}
	}
	return server.TaskResponse{}, errors.New("missing")
}
func (f *fakeBackend) CreateTask(_ context.Context, request server.TaskCreateRequest) (server.TaskResponse, error) {
	f.created = request
	task := domain.Task{ID: "new", Name: request.Name, Command: request.Command, Args: request.Args, Env: request.Env, Timezone: request.Timezone, UpdatedAt: time.Now().UTC()}
	return server.TaskResponse{Task: task}, f.err
}
func (f *fakeBackend) UpdateTask(_ context.Context, id string, request server.TaskUpdateRequest) (server.TaskResponse, error) {
	f.updated = request
	return f.GetTask(context.Background(), id)
}
func (f *fakeBackend) DeleteTask(context.Context, string) error           { return f.err }
func (f *fakeBackend) SetTaskEnabled(context.Context, string, bool) error { return f.err }
func (f *fakeBackend) RunNow(context.Context, string) error               { return f.err }
func (f *fakeBackend) Preview(context.Context, server.PreviewRequest) (server.PreviewResponse, error) {
	return server.PreviewResponse{HumanSummary: "Every day"}, f.err
}
func (f *fakeBackend) ListGroups(context.Context) ([]domain.Group, error) { return f.groups, f.err }
func (f *fakeBackend) CreateGroup(context.Context, server.GroupCreateRequest) (domain.Group, error) {
	return domain.Group{}, f.err
}
func (f *fakeBackend) UpdateGroup(context.Context, string, server.GroupUpdateRequest) (domain.Group, error) {
	return domain.Group{}, f.err
}
func (f *fakeBackend) SetGroupEnabled(context.Context, string, bool) error { return f.err }
func (f *fakeBackend) DeleteGroup(context.Context, string) error           { return f.err }

func TestWorkspaceBuildsFullPathsCountsAndEffectiveState(t *testing.T) {
	now := time.Now().UTC()
	root := domain.Group{ID: "root", Name: "Ops", Enabled: false, UpdatedAt: now}
	child := domain.Group{ID: "child", Name: "Nightly", ParentID: root.ID, Enabled: true, UpdatedAt: now}
	task := domain.Task{ID: "task", Name: "", GroupID: child.ID, Command: "echo", Enabled: true, State: domain.TaskActive, Timezone: "UTC", UpdatedAt: now}
	backend := &fakeBackend{groups: []domain.Group{root, child}, tasks: []server.TaskResponse{{Task: task, Readiness: tasklogic.Readiness{CommandReady: true, ActivationReady: true}}}}
	result := NewService(backend).Workspace(context.Background())
	if result.Outcome != "accepted" || result.Workspace.Tasks[0].Name != "unnamed" || result.Workspace.Tasks[0].GroupPath != "Ops / Nightly" || result.Workspace.Tasks[0].EffectiveState != "group_disabled" || result.Workspace.Groups[0].DescendantCount != 1 || result.Workspace.Groups[1].TaskCount != 1 {
		t.Fatalf("result=%+v", result)
	}
}

func TestWorkspaceDisambiguatesDuplicateSiblingsAndRejectsBrokenHierarchy(t *testing.T) {
	now := time.Now().UTC()
	groups := []domain.Group{{ID: "alpha-id", Name: "Duplicate", Enabled: true, UpdatedAt: now}, {ID: "bravo-id", Name: "Duplicate", Enabled: true, UpdatedAt: now}, {ID: "orphan", Name: "Child", ParentID: "missing", Enabled: true, UpdatedAt: now}}
	result := NewService(&fakeBackend{groups: groups}).Workspace(context.Background())
	if result.Outcome != "accepted" || result.Workspace.Groups[0].Path == result.Workspace.Groups[1].Path {
		t.Fatalf("groups=%+v", result.Workspace.Groups)
	}
	for _, group := range result.Workspace.Groups {
		if group.ID == "orphan" && (group.EffectiveEnabled || group.EffectiveReason == "") {
			t.Fatalf("orphan=%+v", group)
		}
	}
}

func TestSaveTaskParsesExactCommandAndRejectsStaleAuthority(t *testing.T) {
	now := time.Now().UTC()
	backend := &fakeBackend{tasks: []server.TaskResponse{{Task: domain.Task{ID: "task", Command: "old", UpdatedAt: now}}}}
	service := NewService(backend)
	stale := service.SaveTask(context.Background(), TaskDraft{TaskDetail: TaskDetail{ID: "task", CommandLine: "echo hi"}, OriginalUpdatedAt: now.Add(-time.Second).Format(time.RFC3339Nano)})
	if stale.Outcome != "stale" {
		t.Fatalf("stale=%+v", stale)
	}
	created := service.SaveTask(context.Background(), TaskDraft{IsNew: true, TaskDetail: TaskDetail{CommandLine: `program "two words"`, Timezone: "Local", Mode: "manual"}})
	if created.Outcome != "accepted" || backend.created.Command != "program" || len(backend.created.Args) != 1 || backend.created.Args[0] != "two words" || backend.created.Enabled == nil || *backend.created.Enabled {
		t.Fatalf("created=%+v request=%+v", created, backend.created)
	}
}

func TestSaveTaskClearsOptionalFieldsAndMapsFailuresSafely(t *testing.T) {
	now := time.Now().UTC()
	backend := &fakeBackend{tasks: []server.TaskResponse{{Task: domain.Task{ID: "task", Command: "echo", WorkingDir: "secret-path", RunAs: "secret-user", UpdatedAt: now}}}}
	result := NewService(backend).SaveTask(context.Background(), TaskDraft{TaskDetail: TaskDetail{ID: "task", CommandLine: "echo"}, OriginalUpdatedAt: now.Format(time.RFC3339Nano), OverwriteStale: true})
	if result.Outcome != "accepted" || !backend.updated.ClearWorkingDir || !backend.updated.ClearRunAs {
		t.Fatalf("result=%+v request=%+v", result, backend.updated)
	}
	backend.err = errors.New("backend path C:/private command password")
	failure := NewService(backend).Workspace(context.Background())
	if failure.Outcome != "unavailable" || failure.Message != "The local scheduler service is unavailable. Try again." {
		t.Fatalf("failure=%+v", failure)
	}
}

func TestTaskDetailAndCommandPreviewUseEmptyCollections(t *testing.T) {
	now := time.Now().UTC()
	backend := &fakeBackend{tasks: []server.TaskResponse{{Task: domain.Task{ID: "task", Command: "/usr/bin/sw_vers", Timezone: "UTC", UpdatedAt: now}}}}
	service := NewService(backend)
	detail := service.Task(context.Background(), "task")
	if detail.Task == nil || detail.Task.Environment == nil {
		t.Fatalf("task detail must serialize an empty environment array: %+v", detail)
	}
	preview := service.PreviewTask(context.Background(), TaskDraft{TaskDetail: TaskDetail{CommandLine: "/usr/bin/sw_vers", Mode: "manual"}})
	if preview.Command == nil || preview.Command.Args == nil || len(preview.Command.Args) != 0 {
		t.Fatalf("command preview must serialize an empty argument array: %+v", preview)
	}
}

func TestOneOffPreviewUsesTaskTimezoneAndRejectsInvalidTimes(t *testing.T) {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(&fakeBackend{})
	future := time.Now().In(location).Add(24 * time.Hour).Truncate(time.Minute)
	draft := TaskDraft{TaskDetail: TaskDetail{CommandLine: "echo", Mode: "one_off", Timezone: "America/New_York", At: future.Format("2006-01-02T15:04")}}
	result := service.PreviewTask(context.Background(), draft)
	if result.Outcome != "accepted" || result.Task == nil || len(result.Task.NextRuns) != 1 || result.Task.NextRuns[0] != future.Format(time.RFC3339) || result.Task.ScheduleSummary == "" || result.Task.PolicySummary == "" {
		t.Fatalf("preview=%+v", result)
	}
	past := draft
	past.At = time.Now().In(location).Add(-time.Hour).Format("2006-01-02T15:04")
	if rejected := service.PreviewTask(context.Background(), past); rejected.Outcome != "rejected" || rejected.Field != "at" {
		t.Fatalf("past preview=%+v", rejected)
	}
	invalid := draft
	invalid.At = "not-a-time"
	if rejected := service.PreviewTask(context.Background(), invalid); rejected.Outcome != "rejected" || rejected.Field != "at" {
		t.Fatalf("invalid preview=%+v", rejected)
	}
}

func TestRecurringDraftRequiresScheduleForPreviewAndSave(t *testing.T) {
	service := NewService(&fakeBackend{})
	draft := TaskDraft{IsNew: true, TaskDetail: TaskDetail{CommandLine: "echo", Mode: "recurring", Schedule: "   "}}
	for _, result := range []OperationResult{service.PreviewTask(context.Background(), draft), service.SaveTask(context.Background(), draft)} {
		if result.Outcome != "rejected" || result.Field != "schedule" {
			t.Fatalf("result=%+v", result)
		}
	}
}

func TestTaskDetailPreservesOneOffInstant(t *testing.T) {
	runAt := time.Date(2030, time.January, 15, 14, 0, 0, 0, time.UTC)
	detail := detailFrom(server.TaskResponse{Task: domain.Task{Timezone: "America/New_York"}, Schedule: &domain.Schedule{Kind: domain.ScheduleOneOff, RunAt: &runAt}})
	if detail.At != "2030-01-15T14:00:00Z" {
		t.Fatalf("at=%q", detail.At)
	}
}

func TestOneHundredRepeatedMutationsAndCancellationsStayBounded(t *testing.T) {
	service := NewService(&fakeBackend{})
	for range 100 {
		if result := service.RunTask(context.Background(), "task"); result.Outcome != "accepted" {
			t.Fatalf("mutation=%+v", result)
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if result := service.Workspace(ctx); result.Outcome != "accepted" {
			t.Fatalf("cancellation=%+v", result)
		}
	}
}
