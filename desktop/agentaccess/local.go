package agentaccess

import (
	"context"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/mcpsession"
)

// Backend supplies the authoritative daemon-owned MCP lifecycle.
type Backend interface {
	Manifest(context.Context) (server.ManifestResponse, error)
	ListActors(context.Context) ([]domain.Actor, error)
	ListCredentials(context.Context) ([]domain.ClientCredential, error)
	ListMCPSessions(context.Context) ([]mcpsession.Session, error)
	ListAudit(context.Context, domain.AuditQuery) ([]domain.AuditEvent, error)
	CreatePairing(context.Context, server.PairingCreateRequest) (domain.PairingSecret, error)
	CancelPairing(context.Context, string) (domain.PairingSession, error)
	UpdateActor(context.Context, string, server.ActorUpdateRequest) (domain.Actor, error)
	RevokeActor(context.Context, string) (domain.Actor, error)
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
type LocalBackend struct{ daemon *client.Client }

// NewLocalBackend creates an Agent Access adapter without opening a listener.
func NewLocalBackend(daemon *client.Client) *LocalBackend { return &LocalBackend{daemon: daemon} }

func (b *LocalBackend) Manifest(ctx context.Context) (server.ManifestResponse, error) {
	return b.daemon.Manifest(ctx)
}
func (b *LocalBackend) ListActors(ctx context.Context) ([]domain.Actor, error) {
	return b.daemon.ListActors(ctx)
}
func (b *LocalBackend) ListCredentials(ctx context.Context) ([]domain.ClientCredential, error) {
	return b.daemon.ListCredentials(ctx)
}
func (b *LocalBackend) ListMCPSessions(ctx context.Context) ([]mcpsession.Session, error) {
	return b.daemon.ListMCPSessions(ctx)
}
func (b *LocalBackend) ListAudit(ctx context.Context, query domain.AuditQuery) ([]domain.AuditEvent, error) {
	return b.daemon.ListAudit(ctx, query)
}
func (b *LocalBackend) CreatePairing(ctx context.Context, request server.PairingCreateRequest) (domain.PairingSecret, error) {
	return b.daemon.CreatePairing(ctx, request)
}
func (b *LocalBackend) CancelPairing(ctx context.Context, id string) (domain.PairingSession, error) {
	return b.daemon.CancelPairing(ctx, id)
}
func (b *LocalBackend) UpdateActor(ctx context.Context, id string, request server.ActorUpdateRequest) (domain.Actor, error) {
	return b.daemon.UpdateActor(ctx, id, request)
}
func (b *LocalBackend) RevokeActor(ctx context.Context, id string) (domain.Actor, error) {
	return b.daemon.RevokeActor(ctx, id)
}

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
