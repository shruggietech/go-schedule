package server

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func apiTask(t *testing.T, s *Server, name string) domain.Task {
	t.Helper()
	rec := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: name, Command: "echo", Schedule: "every hour", Timezone: "UTC",
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create task status=%d body=%s", rec.Code, rec.Body.String())
	}
	var response TaskResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	return response.Task
}

func TestChainsAPILifecycleAndResponseShape(t *testing.T) {
	s := newTestServer(t)
	source := apiTask(t, s, "source")
	target := apiTask(t, s, "target")
	created := doJSON(t, s, http.MethodPost, "/v1/chains", ChainCreateRequest{
		SourceTaskID: source.ID, TargetTaskID: target.ID, OnOutcome: "success",
	})
	if created.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", created.Code, created.Body.String())
	}
	var chain domain.CompletionChain
	if err := json.Unmarshal(created.Body.Bytes(), &chain); err != nil {
		t.Fatal(err)
	}
	if chain.ID == "" || chain.SourceTaskName != "source" || chain.TargetTaskName != "target" || chain.OnOutcome != domain.CompletionOnSuccess {
		t.Fatalf("created chain = %+v", chain)
	}
	list := doJSON(t, s, http.MethodGet, "/v1/chains", nil)
	if list.Code != http.StatusOK || !json.Valid(list.Body.Bytes()) {
		t.Fatalf("list status=%d body=%s", list.Code, list.Body.String())
	}
	outcome := "failure"
	updated := doJSON(t, s, http.MethodPatch, "/v1/chains/"+chain.ID, ChainUpdateRequest{OnOutcome: &outcome})
	if updated.Code != http.StatusOK {
		t.Fatalf("update status=%d body=%s", updated.Code, updated.Body.String())
	}
	if err := json.Unmarshal(updated.Body.Bytes(), &chain); err != nil || chain.OnOutcome != domain.CompletionOnFailure {
		t.Fatalf("updated chain=%+v err=%v", chain, err)
	}
	if deleted := doJSON(t, s, http.MethodDelete, "/v1/chains/"+chain.ID, nil); deleted.Code != http.StatusNoContent {
		t.Fatalf("delete status=%d body=%s", deleted.Code, deleted.Body.String())
	}
	if got := doJSON(t, s, http.MethodGet, "/v1/chains/"+chain.ID, nil); got.Code != http.StatusNotFound {
		t.Fatalf("get deleted status=%d", got.Code)
	}
}

func TestChainsAPIValidationLeavesStateUnchanged(t *testing.T) {
	s := newTestServer(t)
	a := apiTask(t, s, "a")
	b := apiTask(t, s, "b")
	c := apiTask(t, s, "c")
	create := func(source, target domain.Task) domain.CompletionChain {
		rec := doJSON(t, s, http.MethodPost, "/v1/chains", ChainCreateRequest{SourceTaskID: source.ID, TargetTaskID: target.ID, OnOutcome: "success"})
		if rec.Code != http.StatusCreated {
			t.Fatalf("create status=%d body=%s", rec.Code, rec.Body.String())
		}
		var chain domain.CompletionChain
		_ = json.Unmarshal(rec.Body.Bytes(), &chain)
		return chain
	}
	ab := create(a, b)
	_ = create(b, c)
	for name, request := range map[string]ChainCreateRequest{
		"self":      {SourceTaskID: a.ID, TargetTaskID: a.ID, OnOutcome: "success"},
		"outcome":   {SourceTaskID: a.ID, TargetTaskID: c.ID, OnOutcome: "eventually"},
		"duplicate": {SourceTaskID: a.ID, TargetTaskID: b.ID, OnOutcome: "success"},
		"cycle":     {SourceTaskID: c.ID, TargetTaskID: a.ID, OnOutcome: "any"},
	} {
		t.Run(name, func(t *testing.T) {
			rec := doJSON(t, s, http.MethodPost, "/v1/chains", request)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
			}
		})
	}
	empty := doJSON(t, s, http.MethodPatch, "/v1/chains/"+ab.ID, ChainUpdateRequest{})
	if empty.Code != http.StatusBadRequest {
		t.Fatalf("empty patch status=%d body=%s", empty.Code, empty.Body.String())
	}
	got := doJSON(t, s, http.MethodGet, "/v1/chains/"+ab.ID, nil)
	var unchanged domain.CompletionChain
	_ = json.Unmarshal(got.Body.Bytes(), &unchanged)
	if unchanged.SourceTaskID != a.ID || unchanged.TargetTaskID != b.ID || unchanged.OnOutcome != domain.CompletionOnSuccess {
		t.Fatalf("invalid requests mutated chain: %+v", unchanged)
	}
}
