package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return f(request) }

func TestTypedChainClientLifecyclePathsAndShapes(t *testing.T) {
	var requests []string
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests = append(requests, request.Method+" "+request.URL.Path)
		status, body := http.StatusOK, `{"id":"c1","source_task_id":"a","source_task_name":"A","target_task_id":"b","target_task_name":"B","on_outcome":"success"}`
		switch {
		case request.Method == http.MethodPost:
			status = http.StatusCreated
		case request.Method == http.MethodGet && request.URL.Path == "/v1/chains":
			body = `{"chains":[` + body + `]}`
		case request.Method == http.MethodDelete:
			status, body = http.StatusNoContent, ""
		}
		return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})
	client := &Client{http: &http.Client{Transport: transport}}
	ctx := context.Background()
	created, err := client.CreateChain(ctx, server.ChainCreateRequest{SourceTaskID: "a", TargetTaskID: "b", OnOutcome: "success"})
	if err != nil || created.ID != "c1" || created.OnOutcome != domain.CompletionOnSuccess {
		t.Fatalf("CreateChain=%+v err=%v", created, err)
	}
	listed, err := client.ListChains(ctx)
	if err != nil || len(listed) != 1 || listed[0].SourceTaskName != "A" {
		t.Fatalf("ListChains=%+v err=%v", listed, err)
	}
	if _, err := client.GetChain(ctx, "c1"); err != nil {
		t.Fatal(err)
	}
	outcome := "any"
	if _, err := client.UpdateChain(ctx, "c1", server.ChainUpdateRequest{OnOutcome: &outcome}); err != nil {
		t.Fatal(err)
	}
	if err := client.DeleteChain(ctx, "c1"); err != nil {
		t.Fatal(err)
	}
	want := []string{"POST /v1/chains", "GET /v1/chains", "GET /v1/chains/c1", "PATCH /v1/chains/c1", "DELETE /v1/chains/c1"}
	if strings.Join(requests, "|") != strings.Join(want, "|") {
		t.Fatalf("requests=%v want=%v", requests, want)
	}
}

func TestTypedChainClientPreservesValidationEnvelope(t *testing.T) {
	client := &Client{http: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		body := `{"error":{"code":"validation_failed","field":"target_task_id","message":"chain would create a cycle"}}`
		return &http.Response{StatusCode: http.StatusBadRequest, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}}
	_, err := client.CreateChain(context.Background(), server.ChainCreateRequest{})
	status, ok := err.(*StatusError)
	if !ok || status.Code != server.CodeValidation || status.Field != "target_task_id" {
		t.Fatalf("error=%T %+v", err, err)
	}
}
