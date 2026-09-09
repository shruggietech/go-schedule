package agentaccess

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

type fakeBackend struct {
	status     server.MCPHTTPStatusResponse
	secret     string
	err        error
	disableErr error
	disabled   int
}

func (b *fakeBackend) MCPHTTPStatus(context.Context) (server.MCPHTTPStatusResponse, error) {
	return b.status, b.err
}
func (b *fakeBackend) EnableMCPHTTP(_ context.Context, req server.MCPHTTPEnableRequest) (server.MCPHTTPCredentialResponse, error) {
	b.status.Enabled = true
	b.status.ClientName = strings.TrimSpace(req.ClientName)
	return server.MCPHTTPCredentialResponse{MCPHTTPStatusResponse: b.status, Credential: b.secret}, b.err
}
func (b *fakeBackend) RotateMCPHTTPCredential(context.Context) (server.MCPHTTPCredentialResponse, error) {
	return server.MCPHTTPCredentialResponse{MCPHTTPStatusResponse: b.status, Credential: b.secret}, b.err
}
func (b *fakeBackend) DisableMCPHTTP(context.Context) (server.MCPHTTPStatusResponse, error) {
	b.disabled++
	if b.disableErr != nil {
		return b.status, b.disableErr
	}
	b.status = server.MCPHTTPStatusResponse{AllowedOrigins: []string{}}
	return b.status, nil
}

type fakeNative struct {
	copied         string
	opened         string
	copyErr        error
	waitForContext bool
}

func (n *fakeNative) ClipboardSetText(ctx context.Context, value string) error {
	n.copied = value
	if n.waitForContext {
		<-ctx.Done()
		return ctx.Err()
	}
	return n.copyErr
}
func (n *fakeNative) BrowserOpenURL(_ context.Context, value string) error {
	n.opened = value
	return nil
}

func TestWorkspaceProjectsAuthorityAndEvidenceWithoutSecret(t *testing.T) {
	now := time.Date(2026, 9, 9, 12, 0, 0, 0, time.UTC)
	backend := &fakeBackend{status: server.MCPHTTPStatusResponse{Enabled: true, Endpoint: "http://127.0.0.1:43123/mcp", AllowedOrigins: []string{}, ClientName: "Codex", RequestCount: 3, LastAccessedAt: &now}}
	result := NewService(backend, &fakeNative{}).Workspace(context.Background())
	if result.Outcome != "accepted" || result.Workspace == nil || len(result.Workspace.Authorities) != 3 || result.Workspace.HTTP.RequestCount != 3 {
		t.Fatalf("workspace = %+v", result)
	}
	if strings.Contains(result.Message, "Bearer") {
		t.Fatal("workspace exposed credential language")
	}
}

func TestEnableCopiesCredentialAndClipboardFailureRevokes(t *testing.T) {
	backend := &fakeBackend{status: server.MCPHTTPStatusResponse{AllowedOrigins: []string{}}, secret: "top-secret"}
	native := &fakeNative{}
	service := NewService(backend, native)
	result := service.Enable(context.Background(), EnableDraft{ClientName: "Codex", Port: 43123})
	if result.Outcome != "accepted" || native.copied != "top-secret" || result.Workspace == nil {
		t.Fatalf("enable = %+v copied=%q", result, native.copied)
	}
	failing := &fakeNative{copyErr: errors.New("denied")}
	result = NewService(backend, failing).Rotate(context.Background())
	if result.Outcome != "unavailable" || result.Workspace == nil || backend.disabled != 1 || backend.status.Enabled {
		t.Fatalf("rollback = %+v disabled=%d", result, backend.disabled)
	}
	backend.status.Enabled = true
	backend.disableErr = errors.New("timeout")
	result = NewService(backend, failing).Rotate(context.Background())
	if result.Workspace != nil || !strings.Contains(result.Message, "could not be confirmed") || !backend.status.Enabled {
		t.Fatalf("unconfirmed rollback = %+v status=%+v", result, backend.status)
	}
}

func TestEnableValidatesNameAndGuideIsFixed(t *testing.T) {
	backend := &fakeBackend{}
	native := &fakeNative{}
	service := NewService(backend, native)
	if result := service.Enable(context.Background(), EnableDraft{ClientName: "  ", Port: 43123}); result.Outcome != "rejected" {
		t.Fatalf("blank name = %+v", result)
	}
	if result := service.OpenGuide(context.Background()); result.Outcome != "accepted" || native.opened != guideURL {
		t.Fatalf("guide = %+v url=%q", result, native.opened)
	}
}

func TestClipboardHandoffCancellationStillRollsBack(t *testing.T) {
	backend := &fakeBackend{status: server.MCPHTTPStatusResponse{Enabled: true, AllowedOrigins: []string{}}, secret: "one-time"}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result := NewService(backend, &fakeNative{waitForContext: true}).Rotate(ctx)
	if result.Workspace == nil || result.Workspace.HTTP.Enabled || backend.disabled != 1 {
		t.Fatalf("deadline rollback = %+v disabled=%d", result, backend.disabled)
	}
}
