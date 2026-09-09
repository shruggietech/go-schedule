package agentaccess

import (
	"context"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

// Backend supplies the authoritative daemon-owned MCP lifecycle.
type Backend interface {
	MCPHTTPStatus(context.Context) (server.MCPHTTPStatusResponse, error)
	EnableMCPHTTP(context.Context, server.MCPHTTPEnableRequest) (server.MCPHTTPCredentialResponse, error)
	RotateMCPHTTPCredential(context.Context) (server.MCPHTTPCredentialResponse, error)
	DisableMCPHTTP(context.Context) (server.MCPHTTPStatusResponse, error)
}

// Native supplies the bounded operating-system actions used by Agent Access.
type Native interface {
	ClipboardSetText(context.Context, string) error
	BrowserOpenURL(context.Context, string) error
}

// LocalBackend adapts the protected local daemon client.
type LocalBackend struct{ daemon Backend }

// NewLocalBackend creates an Agent Access adapter without opening a listener.
func NewLocalBackend(daemon Backend) *LocalBackend { return &LocalBackend{daemon: daemon} }

func (b *LocalBackend) MCPHTTPStatus(ctx context.Context) (server.MCPHTTPStatusResponse, error) {
	return b.daemon.MCPHTTPStatus(ctx)
}
func (b *LocalBackend) EnableMCPHTTP(ctx context.Context, req server.MCPHTTPEnableRequest) (server.MCPHTTPCredentialResponse, error) {
	return b.daemon.EnableMCPHTTP(ctx, req)
}
func (b *LocalBackend) RotateMCPHTTPCredential(ctx context.Context) (server.MCPHTTPCredentialResponse, error) {
	return b.daemon.RotateMCPHTTPCredential(ctx)
}
func (b *LocalBackend) DisableMCPHTTP(ctx context.Context) (server.MCPHTTPStatusResponse, error) {
	return b.daemon.DisableMCPHTTP(ctx)
}
