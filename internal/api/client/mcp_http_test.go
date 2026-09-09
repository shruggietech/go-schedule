package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func TestMCPHTTPClientUsesVersionedLifecyclePaths(t *testing.T) {
	want := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/v1/mcp/http"},
		{http.MethodPost, "/v1/mcp/http/enable"},
		{http.MethodPost, "/v1/mcp/http/rotate"},
		{http.MethodPost, "/v1/mcp/http/disable"},
	}
	call := 0
	client := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method != want[call].method || request.URL.Path != want[call].path {
			t.Fatalf("request %d=%s %s", call, request.Method, request.URL.Path)
		}
		if call == 1 {
			body, _ := io.ReadAll(request.Body)
			if !strings.Contains(string(body), `"port":43123`) || !strings.Contains(string(body), `"allowed_origins"`) {
				t.Fatalf("enable body=%s", body)
			}
		}
		call++
		body := `{"enabled":true,"endpoint":"http://127.0.0.1:43123/mcp","allowed_origins":[],"credential":"secret"}`
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}}
	ctx := context.Background()
	if _, err := client.MCPHTTPStatus(ctx); err != nil {
		t.Fatal(err)
	}
	if result, err := client.EnableMCPHTTP(ctx, server.MCPHTTPEnableRequest{Port: 43123, AllowedOrigins: []string{}}); err != nil || result.Credential != "secret" {
		t.Fatalf("enable=%+v err=%v", result, err)
	}
	if _, err := client.RotateMCPHTTPCredential(ctx); err != nil {
		t.Fatal(err)
	}
	if status, err := client.DisableMCPHTTP(ctx); err != nil || !status.Enabled {
		t.Fatalf("disable=%+v err=%v", status, err)
	}
	if call != len(want) {
		t.Fatalf("calls=%d", call)
	}
}
