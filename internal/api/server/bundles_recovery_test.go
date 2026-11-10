package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/bundle"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestS145PortableRecoveryCleanupRoundTripAndConfigurationDrift(t *testing.T) {
	source := newTestServer(t)
	entry := apiTask(t, source, "prepare")
	branch := apiTask(t, source, "export")
	target := apiTask(t, source, "publish")
	recovery := apiTask(t, source, "recover")
	cleanup := apiTask(t, source, "cleanup")
	join := domain.ParallelJoin{Name: "routing", EntryTaskID: entry.ID, TargetTaskID: target.ID, RecoveryTaskID: recovery.ID, RecoveryOn: []domain.JoinRecoveryCondition{}, CleanupTaskID: cleanup.ID, Branches: []domain.ParallelJoinBranch{{ID: "slot", Name: "slot", TaskID: branch.ID}}}
	if err := source.store.CreateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	exported := doJSON(t, source, http.MethodGet, "/v1/bundles/export", nil)
	var doc bundle.Document
	if exported.Code != http.StatusOK || json.Unmarshal(exported.Body.Bytes(), &doc) != nil || doc.Schema != bundle.SchemaV4 || len(doc.Joins) != 1 || doc.Joins[0].RecoveryOn == nil || len(*doc.Joins[0].RecoveryOn) != 0 {
		t.Fatalf("v4 export status=%d document=%+v", exported.Code, doc)
	}
	if strings.Contains(exported.Body.String(), "settlement") || strings.Contains(exported.Body.String(), "selected_route") || strings.Contains(exported.Body.String(), "delivery_id") {
		t.Fatal("runtime state leaked into portable configuration")
	}
	dest := newTestServer(t)
	preview := doJSON(t, dest, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc})
	var plan bundle.Plan
	if preview.Code != http.StatusOK || json.Unmarshal(preview.Body.Bytes(), &plan) != nil {
		t.Fatalf("preview=%d %s", preview.Code, preview.Body.String())
	}
	required := false
	for _, capability := range plan.RequiredCapabilities {
		required = required || capability == "parallel-join-recovery"
	}
	if !required {
		t.Fatalf("review lost required capability: %+v", plan)
	}
	applied := doJSON(t, dest, http.MethodPost, "/v1/bundles/apply", BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint})
	var response BundleApplyResponse
	if applied.Code != http.StatusOK || json.Unmarshal(applied.Body.Bytes(), &response) != nil {
		t.Fatalf("apply=%d %s", applied.Code, applied.Body.String())
	}
	for _, item := range response.Outcomes {
		if item.Action == bundle.ActionFailed || item.Action == bundle.ActionUncertain {
			t.Fatalf("incomplete portable application: %+v", response)
		}
	}
	imported, err := dest.store.ListParallelJoins()
	if err != nil || len(imported) != 1 {
		t.Fatalf("joins=%+v err=%v", imported, err)
	}
	names := map[string]string{}
	tasks, err := dest.store.ListTasks("", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, task := range tasks {
		names[task.ID] = task.Name
	}
	got := imported[0]
	if names[got.EntryTaskID] != "prepare" || names[got.TargetTaskID] != "publish" || names[got.RecoveryTaskID] != "recover" || names[got.CleanupTaskID] != "cleanup" || names[got.Branches[0].TaskID] != "export" || got.RecoveryTaskID == recovery.ID || got.CleanupTaskID == cleanup.ID || got.RecoveryOn == nil || len(got.RecoveryOn) != 0 {
		t.Fatalf("lost or unmapped routing=%+v", got)
	}
	retained := doJSON(t, source, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc})
	if retained.Code != http.StatusOK || json.Unmarshal(retained.Body.Bytes(), &plan) != nil {
		t.Fatalf("retained=%d %s", retained.Code, retained.Body.String())
	}
	join.CleanupTaskID = target.ID
	if err := source.store.UpdateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	stale := doJSON(t, source, http.MethodPost, "/v1/bundles/apply", BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint})
	if stale.Code != http.StatusConflict {
		t.Fatalf("cleanup-only drift accepted: %d %s", stale.Code, stale.Body.String())
	}
}
