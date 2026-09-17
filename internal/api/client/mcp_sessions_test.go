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

func TestMCPSessionClientUsesLifecyclePathsAndMutationHeaders(t *testing.T) {
	var requests []string
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests = append(requests, request.Method+" "+request.URL.Path)
		if request.URL.Path == "/v1/tasks/task-1/run-now" {
			if got := request.Header.Get(server.MCPSessionHeader); got != "session-secret" {
				t.Fatalf("session header=%q", got)
			}
			if got := request.Header.Get(server.ExpectedDaemonHeader); got != "daemon-1" {
				t.Fatalf("daemon header=%q", got)
			}
		}
		body := "{}"
		status := http.StatusNoContent
		if request.URL.Path == "/v1/mcp/sessions" {
			body = `{"id":"session-1","actor_id":"actor-1","capability":"operate","credential":"session-secret"}`
			status = http.StatusCreated
		}
		return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})
	base := &Client{http: &http.Client{Transport: transport}}
	session, err := base.CreateMCPSession(context.Background(), "Codex", domain.CapabilityOperate)
	if err != nil || session.ID != "session-1" || session.Credential != "session-secret" {
		t.Fatalf("session=%+v err=%v", session, err)
	}
	operating := base.WithMCPSession(session.Credential)
	if err := operating.RunNow(ExpectDaemon(context.Background(), "daemon-1"), "task-1"); err != nil {
		t.Fatal(err)
	}
	if err := base.RevokeMCPSession(context.Background(), session.ID); err != nil {
		t.Fatal(err)
	}
	want := []string{"POST /v1/mcp/sessions", "POST /v1/tasks/task-1/run-now", "DELETE /v1/mcp/sessions/session-1"}
	if len(requests) != len(want) {
		t.Fatalf("requests=%v", requests)
	}
	for index := range want {
		if requests[index] != want[index] {
			t.Fatalf("requests=%v", requests)
		}
	}
}

func TestRemoteClientCannotAttachLocalMCPSessionSecret(t *testing.T) {
	remote := &Client{remote: true, mcpSession: "unexpected"}
	if clone := remote.WithMCPSession("new-secret"); clone.mcpSession != "" {
		t.Fatalf("remote session secret=%q", clone.mcpSession)
	}
}
