package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/bundle"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestS144PortableJoinRoundTripRemapsEveryTaskAndExcludesInstances(t *testing.T) {
	source := newTestServer(t)
	entry := apiTask(t, source, "prepare")
	branch := apiTask(t, source, "export")
	target := apiTask(t, source, "publish")
	join := domain.ParallelJoin{Name: "parallel", EntryTaskID: entry.ID, TargetTaskID: target.ID, Branches: []domain.ParallelJoinBranch{{ID: "first", Name: "first", TaskID: branch.ID}, {ID: "second", Name: "second", TaskID: branch.ID}}}
	if err := source.store.CreateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	native := domain.Run{TaskID: entry.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := source.store.RecordRunAndCreateDeliveries(&native, ""); err != nil {
		t.Fatal(err)
	}
	exported := doJSON(t, source, http.MethodGet, "/v1/bundles/export", nil)
	if exported.Code != 200 {
		t.Fatalf("export=%d %s", exported.Code, exported.Body.String())
	}
	if strings.Contains(exported.Body.String(), native.ID) || strings.Contains(exported.Body.String(), "join_instance") {
		t.Fatalf("runtime leaked=%s", exported.Body.String())
	}
	var doc bundle.Document
	if err := json.Unmarshal(exported.Body.Bytes(), &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Schema != bundle.SchemaV3 || len(doc.Joins) != 1 {
		t.Fatalf("export=%+v", doc)
	}
	dest := newTestServer(t)
	preview := doJSON(t, dest, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc})
	var plan bundle.Plan
	if preview.Code != 200 || json.Unmarshal(preview.Body.Bytes(), &plan) != nil {
		t.Fatalf("preview=%d %s", preview.Code, preview.Body.String())
	}
	applied := doJSON(t, dest, http.MethodPost, "/v1/bundles/apply", BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint})
	if applied.Code != 200 {
		t.Fatalf("apply=%d %s", applied.Code, applied.Body.String())
	}
	joins, err := dest.store.ListParallelJoins()
	if err != nil || len(joins) != 1 {
		t.Fatalf("joins=%+v err=%v outcomes=%s", joins, err, applied.Body.String())
	}
	imported := joins[0]
	if imported.EntryTaskID == entry.ID || imported.TargetTaskID == target.ID || imported.Branches[0].TaskID == branch.ID || imported.Branches[0].TaskID != imported.Branches[1].TaskID || imported.Branches[0].ID != "first" || imported.Branches[1].ID != "second" {
		t.Fatalf("remapped=%+v", imported)
	}
	instances, err := dest.store.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{})
	if err != nil || len(instances) != 0 {
		t.Fatalf("runtime imported=%+v err=%v", instances, err)
	}
}

func TestS144JoinEditInvalidatesRetainedPortablePlan(t *testing.T) {
	s := newTestServer(t)
	entry := apiTask(t, s, "prepare")
	branch := apiTask(t, s, "export")
	target := apiTask(t, s, "publish")
	join := domain.ParallelJoin{Name: "join", EntryTaskID: entry.ID, TargetTaskID: target.ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "a", TaskID: branch.ID}}}
	if err := s.store.CreateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	doc, err := s.bundleDocument()
	if err != nil {
		t.Fatal(err)
	}
	doc.Joins[0].Name = "requested"
	preview := doJSON(t, s, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc})
	var plan bundle.Plan
	if preview.Code != 200 || json.Unmarshal(preview.Body.Bytes(), &plan) != nil {
		t.Fatalf("preview=%s", preview.Body.String())
	}
	join.Name = "concurrent"
	if err := s.store.UpdateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	applied := doJSON(t, s, http.MethodPost, "/v1/bundles/apply", BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint})
	if applied.Code != 409 {
		t.Fatalf("drift=%d %s", applied.Code, applied.Body.String())
	}
}

func TestS144PortableMixedRelationshipReplacementUsesFinalGraph(t *testing.T) {
	s := newTestServer(t)
	a := apiTask(t, s, "A")
	b := apiTask(t, s, "B")
	c := apiTask(t, s, "C")
	d := apiTask(t, s, "D")
	join := domain.ParallelJoin{Name: "join", EntryTaskID: a.ID, TargetTaskID: c.ID, Branches: []domain.ParallelJoinBranch{{ID: "slot", Name: "slot", TaskID: b.ID}}}
	if err := s.store.CreateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	chain := domain.CompletionChain{SourceTaskID: d.ID, TargetTaskID: a.ID, OnOutcome: domain.CompletionOnSuccess}
	if err := s.store.CreateCompletionChain(&chain); err != nil {
		t.Fatal(err)
	}
	doc, err := s.bundleDocument()
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]string{}
	for _, task := range doc.Tasks {
		ids[task.Name] = task.PortableID
	}
	doc.Joins[0].EntryTaskID = ids["B"]
	doc.Joins[0].Branches[0].TaskID = ids["C"]
	doc.Joins[0].TargetTaskID = ids["D"]
	doc.Chains[0].SourceTaskID = ids["C"]
	preview := doJSON(t, s, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc})
	var plan bundle.Plan
	if preview.Code != 200 || json.Unmarshal(preview.Body.Bytes(), &plan) != nil {
		t.Fatalf("preview=%d %s", preview.Code, preview.Body.String())
	}
	for _, item := range plan.Items {
		if item.Action == bundle.ActionInvalid || item.Action == bundle.ActionConflict {
			t.Fatalf("invalid final graph=%+v", plan)
		}
	}
	applied := doJSON(t, s, http.MethodPost, "/v1/bundles/apply", BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint})
	var result BundleApplyResponse
	if applied.Code != 200 || json.Unmarshal(applied.Body.Bytes(), &result) != nil {
		t.Fatalf("apply=%d %s", applied.Code, applied.Body.String())
	}
	for _, item := range result.Outcomes {
		if item.Action == bundle.ActionFailed || item.Action == bundle.ActionUncertain {
			t.Fatalf("non-atomic mixed replacement=%+v", result)
		}
	}
	updated, err := s.store.GetParallelJoin(join.ID)
	if err != nil || updated.EntryTaskID != b.ID || updated.Branches[0].TaskID != c.ID || updated.TargetTaskID != d.ID {
		t.Fatalf("updated join=%+v err=%v", updated, err)
	}
}
