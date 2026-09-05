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
