package agentaccess

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type fakeBackend struct {
	status      server.MCPHTTPStatusResponse
	manifest    server.ManifestResponse
	actors      []domain.Actor
	credentials []domain.ClientCredential
	events      []domain.AuditEvent
	pairing     domain.PairingSecret
	secret      string
	err         error
	disableErr  error
	disabled    int
	lastUpdate  server.ActorUpdateRequest
}

func (b *fakeBackend) Manifest(context.Context) (server.ManifestResponse, error) {
	if b.manifest.InstallationID == "" {
		b.manifest = server.ManifestResponse{InstallationID: "daemon-1", DisplayName: "This computer"}
	}
	return b.manifest, b.err
}
func (b *fakeBackend) ListActors(context.Context) ([]domain.Actor, error) {
	return append([]domain.Actor(nil), b.actors...), b.err
}
func (b *fakeBackend) ListCredentials(context.Context) ([]domain.ClientCredential, error) {
	return append([]domain.ClientCredential(nil), b.credentials...), b.err
}
func (b *fakeBackend) ListAudit(_ context.Context, query domain.AuditQuery) ([]domain.AuditEvent, error) {
	var result []domain.AuditEvent
	for _, event := range b.events {
		if query.ActorID == "" || query.ActorID == event.ActorID {
			result = append(result, event)
		}
	}
	if query.Limit > 0 && len(result) > query.Limit {
		result = result[:query.Limit]
	}
	return result, b.err
}
func (b *fakeBackend) CreatePairing(_ context.Context, request server.PairingCreateRequest) (domain.PairingSecret, error) {
	if b.pairing.ID == "" {
		b.pairing = domain.PairingSecret{PairingSession: domain.PairingSession{ID: "pairing-1", DisplayName: request.DisplayName, Kind: request.Kind, Capability: request.Capability, ExpiresAt: time.Now().Add(10 * time.Minute), GrantExpiresAt: request.ExpiresAt}, Phrase: "one-time-phrase", DaemonID: "daemon-1"}
	}
	return b.pairing, b.err
}
func (b *fakeBackend) CancelPairing(context.Context, string) (domain.PairingSession, error) {
	return b.pairing.PairingSession, b.disableErr
}
func (b *fakeBackend) UpdateActor(_ context.Context, id string, request server.ActorUpdateRequest) (domain.Actor, error) {
	b.lastUpdate = request
	for i := range b.actors {
		if b.actors[i].ID != id {
			continue
		}
		if request.Capability != nil {
			b.actors[i].Capability = *request.Capability
		}
		if request.ExpiresAt != nil {
			b.actors[i].ExpiresAt = request.ExpiresAt
		}
		return b.actors[i], b.err
	}
	return domain.Actor{}, errors.New("missing")
}
func (b *fakeBackend) RevokeActor(_ context.Context, id string) (domain.Actor, error) {
	for i := range b.actors {
		if b.actors[i].ID == id {
			b.actors[i].State = domain.ActorStateRevoked
			return b.actors[i], b.err
		}
	}
	return domain.Actor{}, errors.New("missing")
}

func (b *fakeBackend) MCPHTTPStatus(context.Context) (server.MCPHTTPStatusResponse, error) {
	return b.status, b.err
}
func (b *fakeBackend) EnableMCPHTTP(_ context.Context, req server.MCPHTTPEnableRequest) (server.MCPHTTPCredentialResponse, error) {
	b.status.Enabled = true
	b.status.ClientName = strings.TrimSpace(req.ClientName)
	b.status.Permission = req.Permission
	b.status.RequireConfirmation = req.RequireConfirmation
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
	backend := &fakeBackend{status: server.MCPHTTPStatusResponse{Enabled: true, Endpoint: "http://127.0.0.1:43123/mcp", AllowedOrigins: []string{}, ClientName: "Codex", Permission: domain.CapabilityOperate, RequestCount: 3, LastAccessedAt: &now}}
	result := NewService(backend, &fakeNative{}).Workspace(context.Background())
	if result.Outcome != "accepted" || result.Workspace == nil || len(result.Workspace.Authorities) != 3 || result.Workspace.Authorities[1].Status != "available" || result.Workspace.Authorities[2].Status != "available" || result.Workspace.HTTP.Permission != "operate" || result.Workspace.HTTP.RequestCount != 3 {
		t.Fatalf("workspace = %+v", result)
	}
	if strings.Contains(result.Message, "Bearer") {
		t.Fatal("workspace exposed credential language")
	}
}

func TestWorkspaceProjectsTransportAwareSecretFreeGrants(t *testing.T) {
	now := time.Date(2026, 9, 20, 12, 0, 0, 0, time.UTC)
	expires := now.Add(24 * time.Hour)
	backend := &fakeBackend{
		manifest: server.ManifestResponse{InstallationID: "daemon-1", DisplayName: "Workshop", Capabilities: []string{"remote-mcp"}},
		status:   server.MCPHTTPStatusResponse{Enabled: true, ActorID: "local-http", ClientName: "Local agent", Permission: domain.CapabilityOperate, AllowedOrigins: []string{}},
		actors: []domain.Actor{
			{ID: "remote", Kind: domain.ActorKindMCP, DisplayName: "Remote agent", Capability: domain.CapabilityManage, State: domain.ActorStateActive, CreatedAt: now, ExpiresAt: &expires},
			{ID: "local-http", Kind: domain.ActorKindMCP, DisplayName: "Local agent", Capability: domain.CapabilityOperate, State: domain.ActorStateActive, CreatedAt: now},
			{ID: "cli", Kind: domain.ActorKindCLI, DisplayName: "Not an agent", Capability: domain.CapabilityObserve, State: domain.ActorStateActive, CreatedAt: now},
		},
		credentials: []domain.ClientCredential{{ID: "credential", ActorID: "remote", Fingerprint: "safe-fingerprint", State: domain.CredentialActive, CreatedAt: now, UpdatedAt: now}},
	}
	service := NewService(backend, &fakeNative{})
	service.now = func() time.Time { return now }
	result := service.Workspace(context.Background())
	if result.Outcome != "accepted" || result.Workspace == nil || result.Workspace.MCPState != "active" || len(result.Workspace.Grants) != 2 {
		t.Fatalf("result=%+v", result)
	}
	if result.Workspace.Grants[0].Transport != "remote_https" || result.Workspace.Grants[0].DaemonID != "daemon-1" || result.Workspace.Grants[0].CredentialFingerprint != "safe-fingerprint" {
		t.Fatalf("grants=%+v", result.Workspace.Grants)
	}
	if result.Workspace.Transports[2].State != "active" {
		t.Fatalf("transports=%+v", result.Workspace.Transports)
	}
	if strings.Contains(fmt.Sprintf("%+v", result.Workspace), "one-time-phrase") {
		t.Fatal("workspace contains protected enrollment material")
	}
}

func TestRemoteGrantDoesNotImplyRemoteListenerEnablement(t *testing.T) {
	now := time.Date(2026, 9, 20, 12, 0, 0, 0, time.UTC)
	backend := &fakeBackend{
		actors:      []domain.Actor{{ID: "remote", Kind: domain.ActorKindMCP, DisplayName: "Remote agent", Capability: domain.CapabilityObserve, State: domain.ActorStateActive, CreatedAt: now}},
		credentials: []domain.ClientCredential{{ID: "credential", ActorID: "remote", Fingerprint: "safe-fingerprint", State: domain.CredentialActive, CreatedAt: now, UpdatedAt: now}},
		status:      server.MCPHTTPStatusResponse{AllowedOrigins: []string{}},
	}
	service := NewService(backend, &fakeNative{})
	service.now = func() time.Time { return now }
	result := service.Workspace(context.Background())
	if result.Workspace == nil || result.Workspace.MCPState != "off" || result.Workspace.Transports[2].State != "off" || len(result.Workspace.Grants) != 1 {
		t.Fatalf("workspace=%+v", result.Workspace)
	}
}

func TestEnableCopiesCredentialAndClipboardFailureRevokes(t *testing.T) {
	backend := &fakeBackend{status: server.MCPHTTPStatusResponse{AllowedOrigins: []string{}}, secret: "top-secret"}
	native := &fakeNative{}
	service := NewService(backend, native)
	result := service.Enable(context.Background(), EnableDraft{ClientName: "Codex", Port: 43123, Permission: "operate"})
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

func TestEnableForwardsManageConfirmationPolicy(t *testing.T) {
	backend := &fakeBackend{status: server.MCPHTTPStatusResponse{AllowedOrigins: []string{}}, secret: "one-time"}
	result := NewService(backend, &fakeNative{}).Enable(context.Background(), EnableDraft{ClientName: "Codex", Port: 43123, Permission: "manage", RequireConfirmation: true})
	if result.Outcome != "accepted" || backend.status.Permission != domain.CapabilityManage || !backend.status.RequireConfirmation {
		t.Fatalf("result=%+v status=%+v", result, backend.status)
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

func TestCreateGrantCopiesBundleWithoutReturningPhraseAndRollsBackFailure(t *testing.T) {
	now := time.Date(2026, 9, 20, 12, 0, 0, 0, time.UTC)
	backend := &fakeBackend{status: server.MCPHTTPStatusResponse{AllowedOrigins: []string{}}}
	native := &fakeNative{}
	service := NewService(backend, native)
	service.now = func() time.Time { return now }
	result := service.CreateGrant(context.Background(), GrantDraft{ClientName: "Release agent", Capability: "manage", Duration: "7d"})
	if result.Outcome != "accepted" || !strings.Contains(native.copied, "GO_SCHEDULE_MCP_ENROLLMENT_V1") || !strings.Contains(native.copied, "phrase=one-time-phrase") || strings.Contains(fmt.Sprintf("%+v", result), "one-time-phrase") {
		t.Fatalf("result=%+v copied=%q", result, native.copied)
	}
	for _, field := range []string{"display_name=Release agent", "kind=mcp", "capability=manage"} {
		if !strings.Contains(native.copied, field) {
			t.Fatalf("enrollment bundle missing %q: %q", field, native.copied)
		}
	}
	if backend.pairing.GrantExpiresAt == nil || !backend.pairing.GrantExpiresAt.Equal(now.Add(7*24*time.Hour)) {
		t.Fatalf("pairing=%+v", backend.pairing)
	}
	backend = &fakeBackend{status: server.MCPHTTPStatusResponse{AllowedOrigins: []string{}}}
	result = NewService(backend, &fakeNative{copyErr: errors.New("denied")}).CreateGrant(context.Background(), GrantDraft{ClientName: "Agent", Capability: "observe", Duration: "1h"})
	if result.Outcome != "unavailable" || !strings.Contains(result.Message, "cancelled") {
		t.Fatalf("rollback=%+v", result)
	}
}

func TestGrantEditsAreMonotonicAndActionsAreBounded(t *testing.T) {
	now := time.Date(2026, 9, 20, 12, 0, 0, 0, time.UTC)
	expires := now.Add(30 * 24 * time.Hour)
	actor := domain.Actor{ID: "actor-1", Kind: domain.ActorKindMCP, DisplayName: "Agent", Capability: domain.CapabilityManage, State: domain.ActorStateActive, CreatedAt: now, ExpiresAt: &expires}
	events := make([]domain.AuditEvent, 30)
	for i := range events {
		events[i] = domain.AuditEvent{ActorID: actor.ID, DaemonID: "daemon-1", Operation: "tasks.update", TargetKind: "task", Result: domain.AuditResultSucceeded, OccurredAt: now.Add(time.Duration(i) * time.Minute)}
	}
	backend := &fakeBackend{status: server.MCPHTTPStatusResponse{AllowedOrigins: []string{}}, actors: []domain.Actor{actor}, events: events}
	service := NewService(backend, &fakeNative{})
	service.now = func() time.Time { return now }
	if result := service.EditGrant(context.Background(), GrantEditDraft{ActorID: actor.ID, Capability: "operate", Duration: "7d"}); result.Outcome != "accepted" {
		t.Fatalf("edit=%+v", result)
	}
	if backend.actors[0].Capability != domain.CapabilityOperate || backend.actors[0].ExpiresAt == nil || !backend.actors[0].ExpiresAt.Equal(now.Add(7*24*time.Hour)) {
		t.Fatalf("actor=%+v", backend.actors[0])
	}
	if !backend.lastUpdate.Monotonic {
		t.Fatal("desktop grant edit did not request atomic monotonic enforcement")
	}
	if result := service.EditGrant(context.Background(), GrantEditDraft{ActorID: actor.ID, Capability: "manage"}); result.Outcome != "rejected" {
		t.Fatalf("widen=%+v", result)
	}
	actions := service.Actions(context.Background(), actor.ID)
	if actions.Outcome != "accepted" || len(actions.Actions) != recentActionLimit {
		t.Fatalf("actions=%+v", actions)
	}
	if result := service.RevokeGrant(context.Background(), actor.ID); result.Outcome != "accepted" || backend.actors[0].State != domain.ActorStateRevoked {
		t.Fatalf("revoke=%+v actor=%+v", result, backend.actors[0])
	}
}

func TestWorkspaceUsesNewestAuditAsStdioLastUse(t *testing.T) {
	now := time.Date(2026, 9, 20, 12, 0, 0, 0, time.UTC)
	actor := domain.Actor{ID: "stdio-agent", Kind: domain.ActorKindMCP, DisplayName: "Stdio agent", Capability: domain.CapabilityObserve, State: domain.ActorStateActive, CreatedAt: now.Add(-time.Hour)}
	backend := &fakeBackend{status: server.MCPHTTPStatusResponse{AllowedOrigins: []string{}}, actors: []domain.Actor{actor}, events: []domain.AuditEvent{
		{ActorID: actor.ID, OccurredAt: now},
		{ActorID: actor.ID, OccurredAt: now.Add(-time.Minute)},
	}}
	service := NewService(backend, &fakeNative{})
	service.now = func() time.Time { return now }
	result := service.Workspace(context.Background())
	if result.Workspace == nil || len(result.Workspace.Grants) != 1 || result.Workspace.Grants[0].Transport != "stdio" || result.Workspace.Grants[0].LastUsedAt != now.Format(time.RFC3339Nano) {
		t.Fatalf("workspace=%+v", result.Workspace)
	}
}
