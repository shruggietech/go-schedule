package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestS144JoinDefinitionReadPreservesAcceptedWritableText(t *testing.T) {
	s := newTestServer(t)
	entry := apiTask(t, s, "entry")
	branch := apiTask(t, s, "branch")
	target := apiTask(t, s, "target")
	name := strings.Repeat("N", 20000)
	branchName := strings.Repeat("B", 18000)
	req := JoinCreateRequest{Name: name, EntryTaskID: entry.ID, TargetTaskID: target.ID, Branches: []domain.ParallelJoinBranch{{ID: "first", Name: branchName, TaskID: branch.ID}, {ID: "second", Name: "second", TaskID: branch.ID}}}
	created := doJSON(t, s, http.MethodPost, "/v1/joins", req)
	var join domain.ParallelJoin
	if created.Code != http.StatusCreated || json.Unmarshal(created.Body.Bytes(), &join) != nil {
		t.Fatalf("accepted definition=%d %s", created.Code, created.Body.String())
	}
	observed := doJSON(t, s, http.MethodGet, "/v1/joins/"+join.ID, nil)
	var bounded domain.ParallelJoin
	if observed.Code != http.StatusOK || json.Unmarshal(observed.Body.Bytes(), &bounded) != nil {
		t.Fatalf("observe=%d %s", observed.Code, observed.Body.String())
	}
	if bounded.Name == name || bounded.Branches[0].Name == branchName || len(bounded.TruncatedFields) == 0 || len(bounded.Branches[0].TruncatedFields) == 0 {
		t.Fatal("ordinary observation lost its bounded text contract")
	}
	definition := doJSON(t, s, http.MethodGet, "/v1/joins/"+join.ID+"/definition?branch_limit=1", nil)
	var complete domain.ParallelJoin
	if definition.Code != http.StatusOK || json.Unmarshal(definition.Body.Bytes(), &complete) != nil {
		t.Fatalf("authoring definition=%d %s", definition.Code, definition.Body.String())
	}
	if complete.Name != name || len(complete.Branches) != 1 || complete.Branches[0].Name != branchName || complete.Branches[0].ID != "first" || complete.BranchCount != 2 || !complete.HasMoreBranches || len(complete.TruncatedFields) != 0 || len(complete.Branches[0].TruncatedFields) != 0 {
		t.Fatalf("authoring text or pagination changed: name length=%d branches=%+v", len(complete.Name), complete.Branches)
	}
	if complete.EntryTaskName != "" || complete.TargetTaskName != "" || complete.Branches[0].TaskName != "" {
		t.Fatal("authoring response includes unnecessary task display fields")
	}
	next := doJSON(t, s, http.MethodGet, "/v1/joins/"+join.ID+"/definition?branch_offset=1&branch_limit=1", nil)
	var second domain.ParallelJoin
	if next.Code != http.StatusOK || json.Unmarshal(next.Body.Bytes(), &second) != nil || len(second.Branches) != 1 || second.Branches[0].ID != "second" || second.HasMoreBranches || !second.UpdatedAt.Equal(complete.UpdatedAt) {
		t.Fatalf("next authoring page=%d %s", next.Code, next.Body.String())
	}
	audits, err := s.store.ListAudit(domain.AuditQuery{Operation: "joins.definition.read", Result: domain.AuditResultSucceeded, Limit: 10})
	if err != nil || len(audits) != 2 {
		t.Fatalf("authoring read audit count=%d err=%v", len(audits), err)
	}
	for _, audit := range audits {
		if audit.TargetID != join.ID {
			t.Fatalf("authoring audit target=%q", audit.TargetID)
		}
	}
	for _, query := range []string{"branch_limit=0", "branch_limit=1001", "branch_offset=-1"} {
		invalid := doJSON(t, s, http.MethodGet, "/v1/joins/"+join.ID+"/definition?"+query, nil)
		if invalid.Code != http.StatusBadRequest {
			t.Fatalf("authoring page bounds %s=%d", query, invalid.Code)
		}
	}
	observer, err := s.store.CreateActor(domain.ActorKindJSON, "observer", domain.CapabilityObserve, nil)
	if err != nil {
		t.Fatal(err)
	}
	s.SetActorResolver(func(*http.Request) (string, error) { return observer.ID, nil })
	if denied := doJSON(t, s, http.MethodGet, "/v1/joins/"+join.ID+"/definition", nil); denied.Code != http.StatusForbidden {
		t.Fatalf("authoring authority=%d %s", denied.Code, denied.Body.String())
	}
	if allowed := doJSON(t, s, http.MethodGet, "/v1/joins/"+join.ID, nil); allowed.Code != http.StatusOK {
		t.Fatalf("bounded observation authority=%d", allowed.Code)
	}
}

func TestS144JoinMutationPreservesLargeAcceptedDefinitionAndBodyValidation(t *testing.T) {
	s := newTestServer(t)
	entry := apiTask(t, s, "entry")
	branch := apiTask(t, s, "branch")
	target := apiTask(t, s, "target")
	next := apiTask(t, s, "next")
	name := strings.Repeat("N", 96000)
	branchName := strings.Repeat("B", 80000)
	join := domain.ParallelJoin{Name: name, EntryTaskID: entry.ID, TargetTaskID: target.ID, Branches: []domain.ParallelJoinBranch{{ID: "slot", Name: branchName, TaskID: branch.ID}}}
	if err := s.store.CreateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	// The current store already accepts this definition; a complete edit must preserve it.
	branches := []domain.ParallelJoinBranch{{ID: "slot", Name: branchName, TaskID: branch.ID}}
	patched := doJSON(t, s, http.MethodPatch, "/v1/joins/"+join.ID, JoinUpdateRequest{Name: &name, Branches: &branches, TargetTaskID: &next.ID, ExpectedUpdatedAt: &join.UpdatedAt})
	if patched.Code != http.StatusOK {
		t.Fatalf("large accepted edit=%d %s", patched.Code, patched.Body.String())
	}
	persisted, err := s.store.GetParallelJoin(join.ID)
	if err != nil || persisted.Name != name || len(persisted.Branches) != 1 || persisted.Branches[0].Name != branchName || persisted.TargetTaskID != next.ID {
		t.Fatalf("large edit changed text: name length=%d branch count=%d err=%v", len(persisted.Name), len(persisted.Branches), err)
	}
	for _, method := range []string{http.MethodPost, http.MethodPatch} {
		route := "/v1/joins"
		var payload any = JoinCreateRequest{Name: strings.Repeat("X", 2<<20), EntryTaskID: entry.ID, TargetTaskID: next.ID, Branches: branches}
		if method == http.MethodPatch {
			route += "/" + join.ID
			oversize := strings.Repeat("X", 2<<20)
			payload = JoinUpdateRequest{Name: &oversize}
		}
		rejected := doJSON(t, s, method, route, payload)
		if rejected.Code != http.StatusRequestEntityTooLarge {
			t.Fatalf("%s oversize=%d %s", method, rejected.Code, rejected.Body.String())
		}
	}
	streamed := httptest.NewRequest(http.MethodPatch, "/v1/joins/"+join.ID, strings.NewReader(`{"name":"`+strings.Repeat("X", 2<<20)+`"}`))
	streamed.ContentLength = -1
	streamResponse := httptest.NewRecorder()
	s.Handler().ServeHTTP(streamResponse, streamed)
	if streamResponse.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("streamed oversize=%d %s", streamResponse.Code, streamResponse.Body.String())
	}
	for _, raw := range []string{`{"name":"changed","unknown":true}`, `{"name":"changed"} {}`} {
		req := httptest.NewRequest(http.MethodPatch, "/v1/joins/"+join.ID, strings.NewReader(raw))
		response := httptest.NewRecorder()
		s.Handler().ServeHTTP(response, req)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("invalid body=%d %s", response.Code, response.Body.String())
		}
	}
	after, err := s.store.GetParallelJoin(join.ID)
	if err != nil || after.Name != name || after.Branches[0].Name != branchName {
		t.Fatal("rejected body changed accepted definition")
	}
}
