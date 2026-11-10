package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestJoinExactBranchesNeverUseAmbientTaskHistory(t *testing.T) {
	calls := 0
	c := &Client{http: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.Path != "/v1/join-instances/exact" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		b := domain.ParallelJoinInstance{ID: "exact", EntryRunID: "entry-run", BranchCount: 2, Branches: []domain.ParallelJoinInvocation{{ID: "inv-a", BranchID: "a", TaskID: "same", RunID: "run-a"}}, HasMoreBranches: true}
		if r.URL.Query().Get("branch_offset") == "1" {
			b.Branches = []domain.ParallelJoinInvocation{{ID: "inv-b", BranchID: "b", TaskID: "same", RunID: "run-b"}}
			b.HasMoreBranches = false
			b.BranchOffset = 1
		}
		raw, _ := json.Marshal(b)
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(string(raw)))}, nil
	})}}
	got, err := c.GetJoinInstance(context.Background(), "exact")
	if err != nil || calls != 2 || len(got.Branches) != 2 || got.Branches[0].RunID == got.Branches[1].RunID {
		t.Fatalf("got=%+v calls=%d err=%v", got, calls, err)
	}
}

func TestJoinMutationOmitsObservedMetadata(t *testing.T) {
	c := NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(raw), "task_name") || strings.Contains(string(raw), "truncated_fields") {
			t.Fatalf("read-only metadata submitted: %s", raw)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"id":"join"}`)
	}))
	_, err := c.CreateJoin(context.Background(), server.JoinCreateRequest{Name: "X", Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "A", TaskID: "task", TaskName: "Observed", TruncatedFields: []string{"name"}}}})
	if err != nil {
		t.Fatal(err)
	}
}
