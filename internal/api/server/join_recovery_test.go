package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestS145RecoveryAuthoringDefaultsEmptyPolicyClearingAndFullText(t *testing.T) {
	s := newTestServer(t)
	entry := apiTask(t, s, "entry")
	branch := apiTask(t, s, "branch")
	target := apiTask(t, s, "publish")
	recovery := apiTask(t, s, strings.Repeat("R", 10000))
	cleanup := apiTask(t, s, strings.Repeat("C", 11000))
	name := strings.Repeat("界", 3000)
	created := doJSON(t, s, http.MethodPost, "/v1/joins", JoinCreateRequest{Name: name, EntryTaskID: entry.ID, TargetTaskID: target.ID, RecoveryTaskID: recovery.ID, CleanupTaskID: cleanup.ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "A", TaskID: branch.ID}}})
	var join domain.ParallelJoin
	if created.Code != http.StatusCreated || json.Unmarshal(created.Body.Bytes(), &join) != nil {
		t.Fatalf("create=%d %s", created.Code, created.Body.String())
	}
	if len(join.RecoveryOn) != 1 || join.RecoveryOn[0] != domain.JoinRecoveryOnFailure {
		t.Fatalf("default=%+v", join)
	}
	route := "/v1/joins/" + join.ID
	empty := []domain.JoinRecoveryCondition{}
	changed := doJSON(t, s, http.MethodPatch, route, JoinUpdateRequest{RecoveryOn: &empty, ExpectedUpdatedAt: &join.UpdatedAt})
	if changed.Code != http.StatusOK || json.Unmarshal(changed.Body.Bytes(), &join) != nil || join.RecoveryOn == nil || len(join.RecoveryOn) != 0 {
		t.Fatalf("explicit empty=%d %s", changed.Code, changed.Body.String())
	}
	read := doJSON(t, s, http.MethodGet, route+"/definition", nil)
	var full domain.ParallelJoin
	if read.Code != http.StatusOK || json.Unmarshal(read.Body.Bytes(), &full) != nil || full.Name != name || full.RecoveryTaskID != recovery.ID || full.CleanupTaskID != cleanup.ID || full.RecoveryTaskName != "" || full.CleanupTaskName != "" || full.RecoveryOn == nil || len(full.RecoveryOn) != 0 {
		t.Fatalf("complete authoring status=%d definition=%+v", read.Code, full)
	}
	observed := doJSON(t, s, http.MethodGet, route+"?text_limit=16", nil)
	var bounded domain.ParallelJoin
	if observed.Code != http.StatusOK || json.Unmarshal(observed.Body.Bytes(), &bounded) != nil || len(bounded.RecoveryTaskName) > 16 || len(bounded.CleanupTaskName) > 16 {
		t.Fatalf("bounded=%d %+v", observed.Code, bounded)
	}
	removed := ""
	clear := doJSON(t, s, http.MethodPatch, route, JoinUpdateRequest{RecoveryTaskID: &removed, ExpectedUpdatedAt: &join.UpdatedAt})
	join = domain.ParallelJoin{}
	if clear.Code != http.StatusOK || json.Unmarshal(clear.Body.Bytes(), &join) != nil || join.RecoveryTaskID != "" || len(join.RecoveryOn) != 0 || join.CleanupTaskID != cleanup.ID {
		t.Fatalf("clear recovery status=%d recovery=%q causes=%v cleanup=%q", clear.Code, join.RecoveryTaskID, join.RecoveryOn, join.CleanupTaskID)
	}
	enable := doJSON(t, s, http.MethodPatch, route, JoinUpdateRequest{RecoveryTaskID: &recovery.ID, ExpectedUpdatedAt: &join.UpdatedAt})
	if enable.Code != http.StatusOK || json.Unmarshal(enable.Body.Bytes(), &join) != nil || len(join.RecoveryOn) != 1 || join.RecoveryOn[0] != domain.JoinRecoveryOnFailure {
		t.Fatalf("enable default=%d %s", enable.Code, enable.Body.String())
	}
	stale := doJSON(t, s, http.MethodPatch, route, JoinUpdateRequest{CleanupTaskID: &removed, ExpectedUpdatedAt: &full.UpdatedAt})
	if stale.Code != http.StatusConflict {
		t.Fatalf("stale routing edit=%d %s", stale.Code, stale.Body.String())
	}
	for _, payload := range []map[string]any{{"recovery_on": []string{"success"}}, {"cleanup_task_id": "missing"}, {"recovery_task_id": "", "recovery_on": []string{"failure"}}, {"selected_route": "recovery"}} {
		bad := doJSON(t, s, http.MethodPatch, route, payload)
		if bad.Code != http.StatusBadRequest && bad.Code != http.StatusNotFound {
			t.Fatalf("invalid routing=%d %s", bad.Code, bad.Body.String())
		}
	}
	after, err := s.store.GetParallelJoin(join.ID)
	if err != nil || after.Name != name || after.RecoveryTaskID != recovery.ID || after.CleanupTaskID != cleanup.ID || len(after.RecoveryOn) != 1 {
		t.Fatalf("rejected write changed configuration=%+v err=%v", after, err)
	}
}

func TestS145RecoveryAndCleanupInstanceFiltersAndNoStartEvidence(t *testing.T) {
	s := newTestServer(t)
	entry := apiTask(t, s, "entry")
	branch := apiTask(t, s, "branch")
	target := apiTask(t, s, "publish")
	recovery := apiTask(t, s, "recover")
	cleanup := apiTask(t, s, "cleanup")
	join := domain.ParallelJoin{Name: "routing", EntryTaskID: entry.ID, TargetTaskID: target.ID, RecoveryTaskID: recovery.ID, RecoveryOn: []domain.JoinRecoveryCondition{domain.JoinRecoveryOnNotStarted}, CleanupTaskID: cleanup.ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "A", TaskID: branch.ID}}}
	if err := s.store.CreateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	source := domain.Run{TaskID: entry.ID, ScheduledFor: time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := s.store.RecordRunAndCreateDeliveries(&source, ""); err != nil {
		t.Fatal(err)
	}
	claimed, err := s.store.ClaimCompletionDeliveries(100)
	if err != nil || len(claimed) != 1 {
		t.Fatalf("claims=%+v err=%v", claimed, err)
	}
	if err := s.store.ResolveCompletionDelivery(claimed[0].ID, strings.Repeat("no branch ", 30)); err != nil {
		t.Fatal(err)
	}
	response := doJSON(t, s, http.MethodGet, "/v1/join-instances?state=recovery_pending&text_limit=16", nil)
	var list struct {
		Instances []domain.ParallelJoinInstance `json:"instances"`
	}
	if response.Code != http.StatusOK || json.Unmarshal(response.Body.Bytes(), &list) != nil || len(list.Instances) != 1 || list.Instances[0].SelectedRoute != domain.JoinRouteRecovery || list.Instances[0].Branches[0].Settlement != domain.JoinSettlementNotStarted {
		t.Fatalf("recovery filter=%d %s", response.Code, response.Body.String())
	}
	claimed, err = s.store.ClaimCompletionDeliveries(100)
	if err != nil || len(claimed) != 1 || claimed[0].Role != domain.DeliveryRoleJoinRecovery {
		t.Fatalf("recovery claims=%+v err=%v", claimed, err)
	}
	if err := s.store.ResolveCompletionDelivery(claimed[0].ID, strings.Repeat("no recovery ", 30)); err != nil {
		t.Fatal(err)
	}
	response = doJSON(t, s, http.MethodGet, "/v1/join-instances?state=cleanup_pending&text_limit=16", nil)
	if response.Code != http.StatusOK || json.Unmarshal(response.Body.Bytes(), &list) != nil || len(list.Instances) != 1 || list.Instances[0].RecoveryRunID != "" || list.Instances[0].RecoverySettlement != domain.JoinSettlementNotStarted || list.Instances[0].RecoverySettledAt == nil || len(list.Instances[0].RecoveryReason) > 16 {
		t.Fatalf("cleanup filter=%d %s", response.Code, response.Body.String())
	}
	if response := doJSON(t, s, http.MethodGet, "/v1/join-instances?state=wrong", nil); response.Code != http.StatusBadRequest {
		t.Fatalf("unknown state=%d", response.Code)
	}
}
