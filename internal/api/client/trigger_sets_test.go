package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func TestTypedTriggerSetClientLifecyclePathsAndShapes(t *testing.T) {
	var requests []string
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests = append(requests, request.Method+" "+request.URL.Path)
		status := http.StatusOK
		body := `{"id":"set-1","name":"Agents","target_task_id":"task-1","member_count":2,"enabled_count":2,"members":[]}`
		if request.Method == http.MethodPost && request.URL.Path == "/v1/trigger-sets" {
			status = http.StatusCreated
			body = `{"trigger_set":` + body + `,"members":[{"position":1,"trigger_id":"t1","key":"k1","command":"gosched trigger fire k1"}]}`
		} else if request.Method == http.MethodGet && request.URL.Path == "/v1/trigger-sets" {
			body = `{"trigger_sets":[` + body + `]}`
		} else if strings.HasSuffix(request.URL.Path, "/reveal") || strings.HasSuffix(request.URL.Path, "/rotate") {
			body = `{"trigger_set":` + body + `,"members":[{"position":1,"trigger_id":"t1","key":"k1","command":"gosched trigger fire k1"}]}`
		} else if request.Method == http.MethodDelete {
			status, body = http.StatusNoContent, ""
		}
		return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})
	c := &Client{http: &http.Client{Transport: transport}}
	ctx := context.Background()
	if created, err := c.CreateTriggerSet(ctx, server.TriggerSetCreateRequest{Name: "Agents", TargetTaskID: "task-1", Count: 2}); err != nil || len(created.Members) != 1 {
		t.Fatalf("CreateTriggerSet=%+v err=%v", created, err)
	}
	if listed, err := c.ListTriggerSets(ctx); err != nil || len(listed) != 1 || listed[0].ID != "set-1" {
		t.Fatalf("ListTriggerSets=%+v err=%v", listed, err)
	}
	if _, err := c.GetTriggerSet(ctx, "set-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.RevealTriggerSet(ctx, "set-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.RetargetTriggerSet(ctx, "set-1", "task-2"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.SetTriggerSetEnabled(ctx, "set-1", true); err != nil {
		t.Fatal(err)
	}
	if _, err := c.SetTriggerSetEnabled(ctx, "set-1", false); err != nil {
		t.Fatal(err)
	}
	if _, err := c.RotateTriggerSet(ctx, "set-1"); err != nil {
		t.Fatal(err)
	}
	if err := c.DeleteTriggerSet(ctx, "set-1"); err != nil {
		t.Fatal(err)
	}
	want := []string{"POST /v1/trigger-sets", "GET /v1/trigger-sets", "GET /v1/trigger-sets/set-1", "POST /v1/trigger-sets/set-1/reveal", "PATCH /v1/trigger-sets/set-1", "POST /v1/trigger-sets/set-1/enable", "POST /v1/trigger-sets/set-1/disable", "POST /v1/trigger-sets/set-1/rotate", "DELETE /v1/trigger-sets/set-1"}
	if strings.Join(requests, "|") != strings.Join(want, "|") {
		t.Fatalf("requests=%v want=%v", requests, want)
	}
}
