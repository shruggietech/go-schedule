package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/events"
)

func TestS144JoinAPICRUDRevisionAndBounds(t *testing.T) {
	s := newTestServer(t)
	entry := apiTask(t, s, "entry")
	branch := apiTask(t, s, "branch")
	target := apiTask(t, s, "target")
	req := JoinCreateRequest{Name: "block", EntryTaskID: entry.ID, TargetTaskID: target.ID, Branches: []domain.ParallelJoinBranch{{ID: "one", Name: "one", TaskID: branch.ID}, {ID: "two", Name: "two", TaskID: branch.ID}}}
	rec := doJSON(t, s, http.MethodPost, "/v1/joins", req)
	if rec.Code != 201 {
		t.Fatalf("create=%d %s", rec.Code, rec.Body.String())
	}
	var join domain.ParallelJoin
	if err := json.Unmarshal(rec.Body.Bytes(), &join); err != nil {
		t.Fatal(err)
	}
	page := doJSON(t, s, http.MethodGet, "/v1/joins/"+join.ID+"?branch_limit=1", nil)
	var bounded domain.ParallelJoin
	if err := json.Unmarshal(page.Body.Bytes(), &bounded); err != nil {
		t.Fatal(err)
	}
	if page.Code != 200 || len(bounded.Branches) != 1 || bounded.BranchCount != 2 || !bounded.HasMoreBranches {
		t.Fatalf("bounded=%+v status=%d", bounded, page.Code)
	}
	name := "changed"
	stale := join.UpdatedAt.Add(-time.Second)
	bad := doJSON(t, s, http.MethodPatch, "/v1/joins/"+join.ID, JoinUpdateRequest{Name: &name, ExpectedUpdatedAt: &stale})
	if bad.Code != 409 {
		t.Fatalf("stale=%d %s", bad.Code, bad.Body.String())
	}
	good := doJSON(t, s, http.MethodPatch, "/v1/joins/"+join.ID, JoinUpdateRequest{Name: &name, ExpectedUpdatedAt: &join.UpdatedAt})
	if good.Code != 200 {
		t.Fatalf("update=%d %s", good.Code, good.Body.String())
	}
	if r := doJSON(t, s, http.MethodGet, "/v1/join-instances?limit=0", nil); r.Code != 400 {
		t.Fatalf("zero limit=%d", r.Code)
	}
	if r := doJSON(t, s, http.MethodDelete, "/v1/joins/"+join.ID, nil); r.Code != 204 {
		t.Fatalf("delete=%d %s", r.Code, r.Body.String())
	}
}
func TestS144JoinInputRejectsUnknownFields(t *testing.T) {
	s := newTestServer(t)
	r := doJSON(t, s, http.MethodPost, "/v1/joins", map[string]any{"name": "block", "unknown": true})
	if r.Code != 400 {
		t.Fatalf("status=%d %s", r.Code, r.Body.String())
	}
}

func TestS144JoinLiveMergeUsesExactDeliveryAndTerminalWins(t *testing.T) {
	now := time.Now().UTC()
	instance := domain.ParallelJoinInstance{ID: "instance", Branches: []domain.ParallelJoinInvocation{{ID: "invocation", TaskID: "task", DeliveryID: "delivery", DeliveryState: domain.DeliveryClaimed}}}
	active := map[string]domain.Run{"delivery": {ID: "run", DeliveryID: "delivery", JoinInstanceID: "instance", JoinBranchID: "invocation", JoinRole: "branch", TaskID: "task", StartedAt: &now, CancellationStatus: "requested"}}
	mergeJoinActive(&instance, active, 20)
	if instance.Branches[0].RunID != "run" || instance.Branches[0].CancellationStatus != "requested" {
		t.Fatalf("active=%+v", instance)
	}
	instance.Branches[0].DeliveryState = domain.DeliveryCompleted
	instance.Branches[0].Outcome = domain.OutcomeSuccess
	instance.Branches[0].RunID = "committed"
	mergeJoinActive(&instance, active, 20)
	if instance.Branches[0].RunID != "committed" {
		t.Fatal("stale active replaced committed evidence")
	}
}

func TestS144JoinInputRejectsReadOnlyBranchMetadata(t *testing.T) {
	s := newTestServer(t)
	r := doJSON(t, s, http.MethodPost, "/v1/joins", map[string]any{"name": "block", "branches": []any{map[string]any{"id": "slot", "name": "slot", "task_id": "task", "task_name": ""}}})
	if r.Code != 400 {
		t.Fatalf("read-only branch status=%d %s", r.Code, r.Body.String())
	}
}

func TestS144JoinAuthorityIdentityAuditAndRedaction(t *testing.T) {
	s := newTestServer(t)
	entry := apiTask(t, s, "entry")
	branch := apiTask(t, s, "branch")
	target := apiTask(t, s, "target")
	req := JoinCreateRequest{Name: "block", EntryTaskID: entry.ID, TargetTaskID: target.ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "a", TaskID: branch.ID}}}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	actor, err := s.store.CreateActor(domain.ActorKindJSON, "observer", domain.CapabilityObserve, nil)
	if err != nil {
		t.Fatal(err)
	}
	s.SetActorResolver(func(*http.Request) (string, error) { return actor.ID, nil })
	denied := doJSON(t, s, http.MethodPost, "/v1/joins", req)
	if denied.Code != 403 {
		t.Fatalf("observer write=%d", denied.Code)
	}
	read := doJSON(t, s, http.MethodGet, "/v1/joins", nil)
	if read.Code != 200 {
		t.Fatalf("observer read=%d", read.Code)
	}
	manager, err := s.store.CreateActor(domain.ActorKindJSON, "manager", domain.CapabilityManage, nil)
	if err != nil {
		t.Fatal(err)
	}
	s.SetActorResolver(func(*http.Request) (string, error) { return manager.ID, nil })
	wrong := httptest.NewRequest(http.MethodPost, "/v1/joins", bytes.NewReader(data))
	wrong.Header.Set(ExpectedDaemonHeader, "wrong-daemon")
	out := httptest.NewRecorder()
	s.Handler().ServeHTTP(out, wrong)
	if out.Code != 409 {
		t.Fatalf("wrong daemon=%d %s", out.Code, out.Body.String())
	}
	created := doJSON(t, s, http.MethodPost, "/v1/joins", req)
	if created.Code != 201 {
		t.Fatalf("manager=%d %s", created.Code, created.Body.String())
	}
	var join domain.ParallelJoin
	if err := json.Unmarshal(created.Body.Bytes(), &join); err != nil {
		t.Fatal(err)
	}
	audits, err := s.store.ListAudit(domain.AuditQuery{Operation: "joins.create", ActorID: manager.ID, Result: domain.AuditResultSucceeded, Limit: 10})
	if err != nil || len(audits) != 1 || audits[0].TargetID != join.ID {
		t.Fatalf("audit=%+v err=%v", audits, err)
	}
	instanceSource := domain.Run{TaskID: entry.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := s.store.RecordRunAndCreateDeliveries(&instanceSource, ""); err != nil {
		t.Fatal(err)
	}
	observation := doJSON(t, s, http.MethodGet, "/v1/join-instances?join_id="+join.ID+"&branch_limit=1", nil)
	if observation.Code != 200 {
		t.Fatalf("instances=%d %s", observation.Code, observation.Body.String())
	}
	var response struct {
		Instances []domain.ParallelJoinInstance `json:"instances"`
	}
	if err := json.Unmarshal(observation.Body.Bytes(), &response); err != nil || len(response.Instances) != 1 || response.Instances[0].EntryRunID != instanceSource.ID {
		t.Fatalf("instances=%+v err=%v", response, err)
	}
}

func TestS144ResolvedJoinStillShowsOwnedActiveExecution(t *testing.T) {
	now := time.Now().UTC()
	instance := domain.ParallelJoinInstance{ID: "instance", State: domain.JoinInstanceResolved, Reason: "definition deleted", Branches: []domain.ParallelJoinInvocation{{ID: "invocation", TaskID: "task", DeliveryID: "delivery", DeliveryState: domain.DeliveryResolved, Reason: "definition deleted"}}}
	active := map[string]domain.Run{"delivery": {ID: "owned", TaskID: "task", DeliveryID: "delivery", JoinInstanceID: "instance", JoinBranchID: "invocation", JoinRole: "branch", StartedAt: &now, CancellationStatus: "requested", CancellationError: "cleanup detail"}}
	mergeJoinActive(&instance, active, 100)
	if instance.State != domain.JoinInstanceResolved || instance.Reason != "definition deleted" || instance.Branches[0].Reason != "definition deleted" || instance.Branches[0].RunID != "owned" || instance.Branches[0].CancellationStatus != "requested" {
		t.Fatalf("resolved active=%+v", instance)
	}
	instance.Branches[0].Outcome = domain.OutcomeCanceled
	instance.Branches[0].RunID = "terminal"
	mergeJoinActive(&instance, active, 100)
	if instance.Branches[0].RunID != "terminal" || instance.Branches[0].Outcome != domain.OutcomeCanceled {
		t.Fatalf("terminal replaced=%+v", instance)
	}
}

func TestS144JoinEventProjectionCarriesOnlyInvalidationIdentity(t *testing.T) {
	for _, event := range []events.Event{{Kind: events.KindJoin, Join: &events.JoinEvent{ID: "definition", Verb: events.VerbUpdated, Join: &domain.ParallelJoin{Name: "private name"}}}, {Kind: events.KindJoinInstance, JoinInstance: &events.JoinInstanceEvent{ID: "instance"}}} {
		projected, visible := remoteEventProjection(event)
		if !visible || projected.ResourceID == "" {
			t.Fatalf("projection=%+v visible=%v", projected, visible)
		}
		data, err := json.Marshal(projected)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(data, []byte("private name")) {
			t.Fatalf("unbounded event leaked=%s", data)
		}
	}
}

type joinObservationScheduler struct {
	runs       []domain.Run
	beforeRead func()
}

func (*joinObservationScheduler) Reload()             {}
func (*joinObservationScheduler) RunNow(string) error { return nil }
func (s *joinObservationScheduler) ActiveRuns() []domain.Run {
	if s.beforeRead != nil {
		f := s.beforeRead
		s.beforeRead = nil
		f()
	}
	return s.runs
}

func TestS144InstanceAPIReadsTerminalAfterActiveSample(t *testing.T) {
	s := newTestServer(t)
	entry := apiTask(t, s, "entry")
	branch := apiTask(t, s, "branch")
	target := apiTask(t, s, "target")
	join := domain.ParallelJoin{Name: "block", EntryTaskID: entry.ID, TargetTaskID: target.ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "a", TaskID: branch.ID}}}
	if err := s.store.CreateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	source := domain.Run{TaskID: entry.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := s.store.RecordRunAndCreateDeliveries(&source, ""); err != nil {
		t.Fatal(err)
	}
	deliveries, err := s.store.ClaimCompletionDeliveries(10)
	if err != nil || len(deliveries) != 1 {
		t.Fatalf("claims=%+v err=%v", deliveries, err)
	}
	delivery := deliveries[0]
	active := domain.Run{ID: "stale-active", TaskID: branch.ID, DeliveryID: delivery.ID, JoinInstanceID: delivery.JoinInstanceID, JoinBranchID: delivery.JoinBranchID, JoinRole: "branch"}
	s.sched = &joinObservationScheduler{runs: []domain.Run{active}, beforeRead: func() {
		terminal := domain.Run{ID: "committed", TaskID: branch.ID, ScheduledFor: source.ScheduledFor, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerCompletion}
		if err := s.store.RecordRunAndCreateDeliveries(&terminal, delivery.ID); err != nil {
			t.Fatal(err)
		}
	}}
	response := doJSON(t, s, http.MethodGet, "/v1/join-instances/"+delivery.JoinInstanceID, nil)
	var instance domain.ParallelJoinInstance
	if response.Code != 200 || json.Unmarshal(response.Body.Bytes(), &instance) != nil {
		t.Fatalf("get=%d %s", response.Code, response.Body.String())
	}
	if instance.Branches[0].RunID != "committed" || instance.Branches[0].Outcome != domain.OutcomeSuccess {
		t.Fatalf("terminal overwritten=%+v", instance)
	}
}

func TestS144DeletedJoinAPIStillExposesOwnedActiveRun(t *testing.T) {
	s := newTestServer(t)
	entry := apiTask(t, s, "entry")
	branch := apiTask(t, s, "branch")
	target := apiTask(t, s, "target")
	join := domain.ParallelJoin{Name: "block", EntryTaskID: entry.ID, TargetTaskID: target.ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "a", TaskID: branch.ID}}}
	if err := s.store.CreateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	source := domain.Run{TaskID: entry.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := s.store.RecordRunAndCreateDeliveries(&source, ""); err != nil {
		t.Fatal(err)
	}
	deliveries, err := s.store.ClaimCompletionDeliveries(10)
	if err != nil || len(deliveries) != 1 {
		t.Fatalf("claims=%+v err=%v", deliveries, err)
	}
	delivery := deliveries[0]
	now := time.Now().UTC()
	s.sched = &joinObservationScheduler{runs: []domain.Run{{ID: "owned", TaskID: branch.ID, DeliveryID: delivery.ID, JoinInstanceID: delivery.JoinInstanceID, JoinBranchID: delivery.JoinBranchID, JoinRole: "branch", StartedAt: &now, CancellationStatus: "requested"}}}
	if err := s.store.DeleteParallelJoin(join.ID); err != nil {
		t.Fatal(err)
	}
	response := doJSON(t, s, http.MethodGet, "/v1/join-instances/"+delivery.JoinInstanceID, nil)
	var instance domain.ParallelJoinInstance
	if response.Code != 200 || json.Unmarshal(response.Body.Bytes(), &instance) != nil {
		t.Fatalf("get=%d %s", response.Code, response.Body.String())
	}
	if instance.State != domain.JoinInstanceResolved || instance.Branches[0].RunID != "owned" || instance.Branches[0].CancellationStatus != "requested" || instance.Reason == "" {
		t.Fatalf("deleted active=%+v", instance)
	}
}

func TestS144JoinAPIBranchReplacementRemovesAddsAndReorders(t *testing.T) {
	s := newTestServer(t)
	entry := apiTask(t, s, "entry")
	branch := apiTask(t, s, "branch")
	target := apiTask(t, s, "target")
	first := domain.ParallelJoinBranch{ID: "first", Name: "first", TaskID: branch.ID}
	second := domain.ParallelJoinBranch{ID: "second", Name: "second", TaskID: branch.ID}
	third := domain.ParallelJoinBranch{ID: "third", Name: "third", TaskID: branch.ID}
	created := doJSON(t, s, http.MethodPost, "/v1/joins", JoinCreateRequest{Name: "block", EntryTaskID: entry.ID, TargetTaskID: target.ID, Branches: []domain.ParallelJoinBranch{first, second}})
	var join domain.ParallelJoin
	if created.Code != http.StatusCreated || json.Unmarshal(created.Body.Bytes(), &join) != nil {
		t.Fatalf("create=%d %s", created.Code, created.Body.String())
	}
	initialRevision := join.UpdatedAt
	for _, step := range []struct {
		name     string
		branches []domain.ParallelJoinBranch
	}{
		{"remove", []domain.ParallelJoinBranch{first}},
		{"add", []domain.ParallelJoinBranch{first, third}},
		{"reorder", []domain.ParallelJoinBranch{third, first}},
	} {
		updated := doJSON(t, s, http.MethodPatch, "/v1/joins/"+join.ID, JoinUpdateRequest{Branches: &step.branches, ExpectedUpdatedAt: &join.UpdatedAt})
		if updated.Code != http.StatusOK {
			t.Fatalf("%s status=%d body=%s", step.name, updated.Code, updated.Body.String())
		}
		if err := json.Unmarshal(updated.Body.Bytes(), &join); err != nil {
			t.Fatal(err)
		}
		if len(join.Branches) != len(step.branches) || join.BranchCount != len(step.branches) || join.BranchOffset != 0 || join.HasMoreBranches || len(join.TruncatedFields) != 0 {
			t.Fatalf("%s observation metadata=%+v", step.name, join)
		}
		persisted, err := s.store.GetParallelJoin(join.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(persisted.Branches) != len(step.branches) {
			t.Fatalf("%s persisted=%+v", step.name, persisted)
		}
		for i, want := range step.branches {
			if join.Branches[i].ID != want.ID || persisted.Branches[i].ID != want.ID || persisted.Branches[i].TaskID != want.TaskID {
				t.Fatalf("%s slot %d identity changed: response=%+v persisted=%+v", step.name, i, join, persisted)
			}
		}
	}
	staleBranches := []domain.ParallelJoinBranch{first, second}
	stale := doJSON(t, s, http.MethodPatch, "/v1/joins/"+join.ID, JoinUpdateRequest{Branches: &staleBranches, ExpectedUpdatedAt: &initialRevision})
	if stale.Code != http.StatusConflict {
		t.Fatalf("stale replacement=%d %s", stale.Code, stale.Body.String())
	}
	persisted, err := s.store.GetParallelJoin(join.ID)
	if err != nil || len(persisted.Branches) != 2 || persisted.Branches[0].ID != third.ID || persisted.Branches[1].ID != first.ID {
		t.Fatalf("stale replacement changed slots: %+v err=%v", persisted, err)
	}
}
