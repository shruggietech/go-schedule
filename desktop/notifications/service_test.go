package notifications

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type fakeBackend struct {
	channels        []domain.NotificationChannel
	tasks           []server.TaskResponse
	groups          []domain.Group
	deliveries      []domain.NotificationDelivery
	direct          []domain.NotificationAssignment
	effective       domain.EffectiveNotificationPolicy
	err             error
	created         server.NotificationChannelCreateRequest
	updated         server.NotificationChannelUpdateRequest
	rotated         *string
	replaced        server.NotificationAssignmentsRequest
	directByScope   map[string][]domain.NotificationAssignment
	effectiveByTask map[string]domain.EffectiveNotificationPolicy
	coverageErrors  map[string]error
	coverageDelay   time.Duration
	coverageActive  atomic.Int32
	coverageMaximum atomic.Int32
}

func (f *fakeBackend) waitForCoverage(ctx context.Context) error {
	if f.coverageDelay == 0 {
		return nil
	}
	active := f.coverageActive.Add(1)
	defer f.coverageActive.Add(-1)
	for maximum := f.coverageMaximum.Load(); active > maximum && !f.coverageMaximum.CompareAndSwap(maximum, active); maximum = f.coverageMaximum.Load() {
	}
	timer := time.NewTimer(f.coverageDelay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (f *fakeBackend) ListNotificationChannels(context.Context) ([]domain.NotificationChannel, error) {
	return f.channels, f.err
}
func (f *fakeBackend) CreateNotificationChannel(_ context.Context, r server.NotificationChannelCreateRequest) (domain.NotificationChannel, error) {
	f.created = r
	return domain.NotificationChannel{ID: "new"}, f.err
}
func (f *fakeBackend) UpdateNotificationChannel(_ context.Context, _ string, r server.NotificationChannelUpdateRequest) (domain.NotificationChannel, error) {
	f.updated = r
	return domain.NotificationChannel{}, f.err
}
func (f *fakeBackend) SetNotificationChannelEnabled(context.Context, string, bool) (domain.NotificationChannel, error) {
	return domain.NotificationChannel{}, f.err
}
func (f *fakeBackend) RotateNotificationChannelAuthorization(_ context.Context, _ string, v string) (domain.NotificationChannel, error) {
	f.rotated = &v
	return domain.NotificationChannel{}, f.err
}
func (f *fakeBackend) TestNotificationChannel(context.Context, string) (domain.NotificationDelivery, error) {
	return domain.NotificationDelivery{}, f.err
}
func (f *fakeBackend) DeleteNotificationChannel(context.Context, string) error { return f.err }
func (f *fakeBackend) ListNotificationDeliveries(_ context.Context, filter domain.NotificationDeliveryFilter) ([]domain.NotificationDelivery, error) {
	if filter.Limit != 200 {
		return nil, errors.New("bad limit")
	}
	return f.deliveries, f.err
}
func (f *fakeBackend) ListTaskDetails(context.Context, string, string) ([]server.TaskResponse, error) {
	return f.tasks, f.err
}
func (f *fakeBackend) ListGroups(context.Context) ([]domain.Group, error) { return f.groups, f.err }
func (f *fakeBackend) ListTaskNotificationAssignments(_ context.Context, id string) ([]domain.NotificationAssignment, error) {
	if err := f.coverageErrors["task:"+id]; err != nil {
		return nil, err
	}
	if values, ok := f.directByScope["task:"+id]; ok {
		return values, nil
	}
	return f.direct, f.err
}
func (f *fakeBackend) ListGroupNotificationAssignments(ctx context.Context, id string) ([]domain.NotificationAssignment, error) {
	if err := f.waitForCoverage(ctx); err != nil {
		return nil, err
	}
	if err := f.coverageErrors["group:"+id]; err != nil {
		return nil, err
	}
	if values, ok := f.directByScope["group:"+id]; ok {
		return values, nil
	}
	return f.direct, f.err
}
func (f *fakeBackend) EffectiveTaskNotificationPolicy(ctx context.Context, id string) (domain.EffectiveNotificationPolicy, error) {
	if err := f.waitForCoverage(ctx); err != nil {
		return domain.EffectiveNotificationPolicy{}, err
	}
	if err := f.coverageErrors["task:"+id]; err != nil {
		return domain.EffectiveNotificationPolicy{}, err
	}
	if value, ok := f.effectiveByTask[id]; ok {
		return value, nil
	}
	return f.effective, f.err
}
func (f *fakeBackend) ReplaceTaskNotificationAssignments(_ context.Context, _ string, r server.NotificationAssignmentsRequest) ([]domain.NotificationAssignment, error) {
	f.replaced = r
	return nil, f.err
}
func (f *fakeBackend) ReplaceGroupNotificationAssignments(_ context.Context, _ string, r server.NotificationAssignmentsRequest) ([]domain.NotificationAssignment, error) {
	f.replaced = r
	return nil, f.err
}

func fixture() *fakeBackend {
	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	claimed, completed := now.Add(time.Minute), now.Add(2*time.Minute)
	return &fakeBackend{
		channels: []domain.NotificationChannel{{ID: "c1", Name: "Ops", Kind: domain.NotificationChannelWebhook, Endpoint: "https://secret.example/hook", EndpointSummary: "https://secret.example/...", Authorization: "Bearer secret", HasAuthorization: true, Enabled: true, UpdatedAt: now}},
		groups:   []domain.Group{{ID: "g1", Name: "Parent", Enabled: true}, {ID: "g2", Name: "Child", ParentID: "g1", Enabled: true}},
		tasks:    []server.TaskResponse{{Task: domain.Task{ID: "t1", Name: "Backup", GroupID: "g2"}}},
		deliveries: []domain.NotificationDelivery{
			{ID: "d1", ChannelID: "c1", ChannelName: "Ops", DestinationSummary: "https://secret.example/...", Endpoint: "https://secret.example/hook", Authorization: "Bearer secret", EventKind: domain.NotificationEventTest, State: domain.NotificationDeliveryPending, Attempts: 0, CreatedAt: now},
			{ID: "d2", ChannelID: "c1", ChannelName: "Ops", EventKind: domain.NotificationEventRunCompleted, State: domain.NotificationDeliveryPending, Attempts: 1, NextAttemptAt: now.Add(time.Minute), CreatedAt: now.Add(time.Second)},
			{ID: "d3", ChannelID: "c1", ChannelName: "Ops", EventKind: domain.NotificationEventRunCompleted, State: domain.NotificationDeliveryClaimed, Attempts: 1, ClaimedAt: &claimed, CreatedAt: now.Add(2 * time.Second)},
			{ID: "d4", ChannelID: "c1", ChannelName: "Ops", EventKind: domain.NotificationEventRunCompleted, State: domain.NotificationDeliverySucceeded, Attempts: 1, CompletedAt: &completed, CreatedAt: now.Add(3 * time.Second)},
			{ID: "d5", ChannelID: "c1", ChannelName: "Ops", EventKind: domain.NotificationEventRunCompleted, State: domain.NotificationDeliveryFailed, Attempts: 5, LastStatus: 503, LastError: "receiver unavailable", CompletedAt: &completed, CreatedAt: now.Add(4 * time.Second)},
		},
		direct:    []domain.NotificationAssignment{{ChannelID: "c1", ScopeType: domain.NotificationScopeTask, ScopeID: "t1", OnFailure: true}},
		effective: domain.EffectiveNotificationPolicy{TaskID: "t1", SourceScopeType: domain.NotificationScopeGroup, SourceScopeID: "g1", Assignments: []domain.NotificationAssignment{{ChannelID: "c1", OnFailure: true}}},
	}
}

func TestWorkspaceIsCompleteOrderedAndSecretFree(t *testing.T) {
	f := fixture()
	s := NewService(f)
	s.now = func() time.Time { return time.Date(2026, 9, 7, 13, 0, 0, 0, time.UTC) }
	r := s.Workspace(context.Background())
	if r.Outcome != "accepted" || r.Workspace == nil {
		t.Fatalf("result=%+v", r)
	}
	if got := r.Workspace.Groups[1].Context; got != "Parent / Child" {
		t.Fatalf("context=%q", got)
	}
	if got := r.Workspace.Tasks[0].Context; got != "Parent / Child" {
		t.Fatalf("task context=%q", got)
	}
	states := []string{}
	for _, d := range r.Workspace.Deliveries {
		states = append(states, d.State)
	}
	if !reflect.DeepEqual(states, []string{"failed", "successful", "sending", "retrying", "queued"}) {
		t.Fatalf("states=%v", states)
	}
	if r.Workspace.Deliveries[4].Kind != "test" || r.Workspace.Deliveries[3].Kind != "task_outcome" {
		t.Fatalf("kinds=%+v", r.Workspace.Deliveries)
	}
	encoded, err := json.Marshal(r.Workspace)
	if err != nil {
		t.Fatal(err)
	}
	dump := strings.ToLower(string(encoded))
	if strings.Contains(dump, "bearer secret") || strings.Contains(dump, "/hook") {
		t.Fatalf("secret leaked: %s", dump)
	}
}

func TestWorkspaceRejectsPartialFailure(t *testing.T) {
	f := fixture()
	f.err = errors.New("private daemon detail")
	r := NewService(f).Workspace(context.Background())
	if r.Outcome != "unavailable" || r.Workspace != nil || strings.Contains(r.Message, "private") {
		t.Fatalf("result=%+v", r)
	}
}

func TestWorkspaceIncludesOrderedSecretFreeConfiguredCoverage(t *testing.T) {
	f := fixture()
	f.channels = append(f.channels, domain.NotificationChannel{ID: "c2", Name: "Disabled", Kind: domain.NotificationChannelWebhook, Endpoint: "https://private.example/hook", Authorization: "Bearer private", Enabled: false})
	f.directByScope = map[string][]domain.NotificationAssignment{
		"group:g1": {{ChannelID: "c1", ScopeType: domain.NotificationScopeGroup, ScopeID: "g1", OnFailure: true}},
		"group:g2": {},
	}
	f.effectiveByTask = map[string]domain.EffectiveNotificationPolicy{
		"t1": {TaskID: "t1", SourceScopeType: domain.NotificationScopeGroup, SourceScopeID: "g1", Assignments: []domain.NotificationAssignment{{ChannelID: "c1", OnFailure: true}, {ChannelID: "c2", OnSuccess: true}}},
	}
	r := NewService(f).Workspace(context.Background())
	if r.Outcome != "accepted" || r.Workspace == nil || !r.Workspace.CoverageComplete {
		t.Fatalf("result=%+v", r)
	}
	want := []ConfiguredScope{
		{Type: "group", ID: "g1", Name: "Parent", Context: "Parent", SourceType: "group", SourceName: "Parent", OnFailure: true, DestinationCount: 1, EnabledDestinationCount: 1, EnabledFailureCount: 1},
		{Type: "task", ID: "t1", Name: "Backup", Context: "Parent / Child", SourceType: "group", SourceName: "Parent", OnSuccess: true, OnFailure: true, DestinationCount: 2, EnabledDestinationCount: 1, EnabledFailureCount: 1},
	}
	if !reflect.DeepEqual(r.Workspace.Coverage, want) {
		t.Fatalf("coverage=%+v", r.Workspace.Coverage)
	}
	encoded, err := json.Marshal(r.Workspace.Coverage)
	if err != nil {
		t.Fatal(err)
	}
	dump := strings.ToLower(string(encoded))
	if strings.Contains(dump, "private.example") || strings.Contains(dump, "bearer") {
		t.Fatalf("secret leaked: %s", dump)
	}
}

func TestWorkspacePreservesCoreSnapshotWhenCoverageIsIncomplete(t *testing.T) {
	f := fixture()
	f.coverageErrors = map[string]error{"task:t1": errors.New("private policy detail")}
	r := NewService(f).Workspace(context.Background())
	if r.Outcome != "accepted" || r.Workspace == nil || r.Workspace.CoverageComplete {
		t.Fatalf("result=%+v", r)
	}
	if len(r.Workspace.Channels) != 1 || len(r.Workspace.Deliveries) != 5 {
		t.Fatalf("workspace=%+v", r.Workspace)
	}
	encoded, _ := json.Marshal(r.Workspace)
	if strings.Contains(strings.ToLower(string(encoded)), "private policy") {
		t.Fatalf("private error leaked: %s", encoded)
	}
}

func TestWorkspaceBoundsCoverageLookupConcurrency(t *testing.T) {
	f := fixture()
	f.coverageDelay = 25 * time.Millisecond
	f.tasks = append(f.tasks,
		server.TaskResponse{Task: domain.Task{ID: "t2", Name: "Second"}},
		server.TaskResponse{Task: domain.Task{ID: "t3", Name: "Third"}},
		server.TaskResponse{Task: domain.Task{ID: "t4", Name: "Fourth"}},
	)
	r := NewService(f).Workspace(context.Background())
	if r.Outcome != "accepted" || r.Workspace == nil || !r.Workspace.CoverageComplete {
		t.Fatalf("result=%+v", r)
	}
	if maximum := f.coverageMaximum.Load(); maximum < 2 || maximum > 8 {
		t.Fatalf("coverage concurrency=%d, want 2..8", maximum)
	}
}

func TestSaveChannelValidatesReplacementAndClearsAuthorization(t *testing.T) {
	f := fixture()
	s := NewService(f)
	if r := s.SaveChannel(context.Background(), ChannelDraft{IsNew: true, Name: " ", Endpoint: "http://bad"}); r.Outcome != "rejected" || r.Field != "name" {
		t.Fatalf("name=%+v", r)
	}
	if r := s.SaveChannel(context.Background(), ChannelDraft{ID: "c1", Name: "Ops", ReplaceEndpoint: true}); r.Outcome != "rejected" || r.Field != "endpoint" {
		t.Fatalf("endpoint=%+v", r)
	}
	r := s.SaveChannel(context.Background(), ChannelDraft{ID: "c1", Name: "Repaired", ReplaceAuthorization: true, Authorization: ""})
	if r.Outcome != "accepted" || f.rotated == nil || *f.rotated != "" || f.updated.Name == nil || *f.updated.Name != "Repaired" {
		t.Fatalf("result=%+v updated=%+v rotated=%v", r, f.updated, f.rotated)
	}
}

func TestPolicyUsesAuthoritativeEffectiveSourceAndValidatesAssignments(t *testing.T) {
	f := fixture()
	s := NewService(f)
	r := s.Policy(context.Background(), "task", "t1")
	if r.Outcome != "accepted" || r.Policy == nil || r.Policy.EffectiveSourceName != "Parent" || r.Policy.EffectiveSourceType != "group" {
		t.Fatalf("policy=%+v", r)
	}
	bad := s.SavePolicy(context.Background(), PolicyDraft{ScopeType: "task", ScopeID: "t1", Assignments: []Assignment{{ChannelID: "c1"}}})
	if bad.Outcome != "rejected" || bad.Field != "assignments" {
		t.Fatalf("bad=%+v", bad)
	}
	good := s.SavePolicy(context.Background(), PolicyDraft{ScopeType: "task", ScopeID: "t1", Assignments: []Assignment{{ChannelID: "c1", OnFailure: true}}})
	if good.Outcome != "accepted" || len(f.replaced.Assignments) != 1 || !f.replaced.Assignments[0].OnFailure {
		t.Fatalf("good=%+v request=%+v", good, f.replaced)
	}
}
