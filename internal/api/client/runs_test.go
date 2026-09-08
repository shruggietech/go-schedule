package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetRunUsesExactEscapedPathAndDecodesDiagnostics(t *testing.T) {
	var path string
	client := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		path = request.URL.EscapedPath()
		body := `{"id":"run one","task_id":"task","outcome":"failure","output":"partial","output_truncated":true,"trigger":"manual","scheduled_for":"2026-09-05T00:00:00Z"}`
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}}
	run, err := client.GetRun(context.Background(), "run one")
	if err != nil {
		t.Fatal(err)
	}
	if path != "/v1/runs/run%20one" || run.ID != "run one" || !run.OutputTruncated {
		t.Fatalf("path=%q run=%+v", path, run)
	}
}

func TestActiveRunsAndLimitedAlertsUseBoundedPaths(t *testing.T) {
	var paths []string
	client := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		paths = append(paths, request.URL.RequestURI())
		body := `{"runs":[{"id":"active-1","task_id":"task-1","scheduled_for":"2026-09-07T00:00:00Z"}]}`
		if strings.HasPrefix(request.URL.Path, "/v1/alerts") {
			body = `{"alerts":[]}`
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}}
	if runs, err := client.ListActiveRuns(context.Background()); err != nil || len(runs) != 1 {
		t.Fatalf("runs=%+v err=%v", runs, err)
	}
	if _, err := client.ListAlertsLimited(context.Background(), true, 200); err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 || paths[0] != "/v1/runs/active" || paths[1] != "/v1/alerts?limit=200&unacked=true" {
		t.Fatalf("paths=%v", paths)
	}
}

func TestPagedRunsAndAlertsSendBoundaryParameters(t *testing.T) {
	var paths []string
	client := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		paths = append(paths, request.URL.RequestURI())
		collection := "runs"
		if request.URL.Path == "/v1/alerts" {
			collection = "alerts"
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"` + collection + `":[]}`))}, nil
	})}}
	if _, err := client.ListRunsPage(context.Background(), "task", 100, 101, 8192); err != nil {
		t.Fatal(err)
	}
	if _, err := client.ListAlertsPage(context.Background(), true, 200, 101, 2048); err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 || paths[0] != "/v1/runs?limit=101&offset=100&output_limit=8192&task=task" || paths[1] != "/v1/alerts?limit=101&message_limit=2048&offset=200&unacked=true" {
		t.Fatalf("paths=%v", paths)
	}
}
