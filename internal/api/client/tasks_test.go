package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestListTaskDetailsRequestsOptInShape(t *testing.T) {
	var query string
	client := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		query = request.URL.RawQuery
		body := `{"tasks":[{"task":{"id":"task-1"},"schedule":null,"readiness":{"status":"not_runnable"},"policy_summary":"","next_runs":[]}]}`
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}}
	tasks, err := client.ListTaskDetails(context.Background(), "group one", "active")
	if err != nil || len(tasks) != 1 || tasks[0].Task.ID != "task-1" || query != "details=true&group=group+one&state=active" {
		t.Fatalf("query=%q tasks=%+v err=%v", query, tasks, err)
	}
}

func TestListTaskObservationsRequestsBoundedProjection(t *testing.T) {
	var query string
	client := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		query = request.URL.RawQuery
		body := `{"tasks":[{"id":"task-1","name":"task","enabled":true,"state":"active","timezone":"UTC","readiness":"ready","has_schedule":true,"next_runs":[],"updated_at":"2026-09-08T00:00:00Z"}]}`
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}}
	tasks, err := client.ListTaskObservations(context.Background(), "group one", "active", true, 100, 101, 2048)
	wantQuery := "group=group+one&limit=101&observation=true&offset=100&scheduled=true&state=active&text_limit=2048"
	if err != nil || len(tasks) != 1 || tasks[0].ID != "task-1" || query != wantQuery {
		t.Fatalf("query=%q tasks=%+v err=%v", query, tasks, err)
	}
}
