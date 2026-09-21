package agentaccess

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

const operationTimeout = 2 * time.Second
const guideURL = "https://shruggietech.github.io/go-schedule/mcp/"
const recentActionLimit = 25

// Service serializes Agent Access operations and keeps credentials behind the native boundary.
type Service struct {
	mu      sync.Mutex
	backend Backend
	native  Native
	now     func() time.Time
}

// NewService creates an Agent Access service.
func NewService(backend Backend, native Native) *Service {
	return &Service{backend: backend, native: native, now: func() time.Time { return time.Now().UTC() }}
}

// Workspace returns the authoritative secret-free MCP grant and transport state.
func (s *Service) Workspace(ctx context.Context) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.workspaceResult(ctx, "load_agent_access", "Agent access loaded.")
}

// Enable starts one named localhost client and copies its credential once.
func (s *Service) Enable(ctx context.Context, draft EnableDraft) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	if strings.TrimSpace(draft.ClientName) == "" {
		return Result{Action: "enable_agent_access", Outcome: "rejected", Message: "Enter a client name."}
	}
	callCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	result, err := s.backend.EnableMCPHTTP(callCtx, server.MCPHTTPEnableRequest{Port: draft.Port, AllowedOrigins: draft.AllowedOrigins, ClientName: draft.ClientName, Permission: domain.Capability(draft.Permission), RequireConfirmation: draft.RequireConfirmation})
	cancel()
	if err != nil {
		return Result{Action: "enable_agent_access", Outcome: "rejected", Message: "Localhost access could not be enabled. Check the name, port, and origins, then try again."}
	}
	return s.handoff(ctx, "enable_agent_access", "Localhost access enabled. The credential was copied once.", result)
}

// Rotate replaces and copies the active localhost credential while preserving listener policy.
func (s *Service) Rotate(ctx context.Context) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	callCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	result, err := s.backend.RotateMCPHTTPCredential(callCtx)
	cancel()
	if err != nil {
		return Result{Action: "rotate_agent_access", Outcome: "rejected", Message: "The credential could not be rotated. Refresh Agent Access, then try again."}
	}
	return s.handoff(ctx, "rotate_agent_access", "Credential rotated and copied once.", result)
}

// Revoke invalidates the localhost credential and closes that optional listener only.
func (s *Service) Revoke(ctx context.Context) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	callCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	_, err := s.backend.DisableMCPHTTP(callCtx)
	cancel()
	if err != nil {
		return unavailable("revoke_agent_access")
	}
	return s.workspaceResult(ctx, "revoke_agent_access", "Localhost access revoked and listener stopped.")
}

// CreateGrant creates a remote MCP pairing and copies its one-time enrollment bundle natively.
func (s *Service) CreateGrant(ctx context.Context, draft GrantDraft) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	name, err := domain.NormalizeActorDisplayName(draft.ClientName)
	capability := domain.Capability(draft.Capability)
	if err != nil || (capability != domain.CapabilityObserve && capability != domain.CapabilityOperate && capability != domain.CapabilityManage) {
		return Result{Action: "create_agent_grant", Outcome: "rejected", Message: "Enter a client name and choose Observe, Operate, or Manage."}
	}
	expiresAt, err := s.expiration(draft.Duration, true)
	if err != nil {
		return Result{Action: "create_agent_grant", Outcome: "rejected", Message: err.Error()}
	}
	callCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	pairing, err := s.backend.CreatePairing(callCtx, server.PairingCreateRequest{DisplayName: name, Kind: domain.ActorKindMCP, Capability: capability, ExpiresAt: expiresAt})
	cancel()
	if err != nil {
		return Result{Action: "create_agent_grant", Outcome: "rejected", Message: "The agent grant could not be created. Refresh Agent Access, then try again."}
	}
	grantExpiry := "non-expiring"
	if pairing.GrantExpiresAt != nil {
		grantExpiry = pairing.GrantExpiresAt.UTC().Format(time.RFC3339)
	}
	bundle := fmt.Sprintf("GO_SCHEDULE_MCP_ENROLLMENT_V1\ndaemon_id=%s\npairing_id=%s\nphrase=%s\ndisplay_name=%s\nkind=%s\ncapability=%s\ngrant_expires_at=%s\n", pairing.DaemonID, pairing.ID, pairing.Phrase, name, domain.ActorKindMCP, pairing.Capability, grantExpiry)
	pairing.Phrase = ""
	copyCtx, copyCancel := context.WithTimeout(ctx, operationTimeout)
	copyErr := context.Canceled
	if s.native != nil {
		copyErr = s.native.ClipboardSetText(copyCtx, bundle)
	}
	copyCancel()
	bundle = ""
	if copyErr != nil {
		rollbackCtx, rollbackCancel := context.WithTimeout(context.Background(), operationTimeout)
		_, rollbackErr := s.backend.CancelPairing(rollbackCtx, pairing.ID)
		rollbackCancel()
		if rollbackErr != nil {
			return Result{Action: "create_agent_grant", Outcome: "unavailable", Message: "The enrollment bundle could not be copied and cancellation could not be confirmed. Cancel the pending pairing from a protected local shell."}
		}
		return Result{Action: "create_agent_grant", Outcome: "unavailable", Message: "The enrollment bundle could not be copied, so the pending grant was cancelled."}
	}
	return s.workspaceResult(ctx, "create_agent_grant", "Agent enrollment created. The one-time bundle was copied and expires in 10 minutes.")
}

// EditGrant applies only monotonic authority and expiry changes.
func (s *Service) EditGrant(ctx context.Context, draft GrantEditDraft) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	actor, ok := s.actor(ctx, draft.ActorID)
	if !ok || actor.Builtin || actor.Kind != domain.ActorKindMCP || actor.State != domain.ActorStateActive || !actor.ActiveAt(s.now()) {
		return Result{Action: "edit_agent_grant", Outcome: "rejected", Message: "This active MCP grant is no longer available to change."}
	}
	request := server.ActorUpdateRequest{Monotonic: true}
	if draft.Capability != "" {
		capability := domain.Capability(draft.Capability)
		if !narrowerOrEqual(actor.Capability, capability) {
			return Result{Action: "edit_agent_grant", Outcome: "rejected", Message: "Existing access can only be narrowed. Create a new grant for higher authority."}
		}
		if capability != actor.Capability {
			request.Capability = &capability
		}
	}
	if draft.Duration != "" {
		expiresAt, err := s.expiration(draft.Duration, false)
		if err != nil || expiresAt == nil || (actor.ExpiresAt != nil && !expiresAt.Before(*actor.ExpiresAt)) {
			return Result{Action: "edit_agent_grant", Outcome: "rejected", Message: "Choose an expiry that is earlier than the current grant deadline."}
		}
		request.ExpiresAt = expiresAt
	}
	if request.Capability == nil && request.ExpiresAt == nil {
		return Result{Action: "edit_agent_grant", Outcome: "rejected", Message: "Choose a lower authority or an earlier expiry."}
	}
	callCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	_, err := s.backend.UpdateActor(callCtx, actor.ID, request)
	cancel()
	if err != nil {
		return unavailable("edit_agent_grant")
	}
	return s.workspaceResult(ctx, "edit_agent_grant", "Agent grant narrowed. Existing connections use the new boundary on their next request.")
}

// RevokeGrant permanently revokes one MCP actor.
func (s *Service) RevokeGrant(ctx context.Context, actorID string) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	actor, ok := s.actor(ctx, actorID)
	if !ok || actor.Builtin || actor.Kind != domain.ActorKindMCP || actor.State != domain.ActorStateActive || !actor.ActiveAt(s.now()) {
		return Result{Action: "revoke_agent_grant", Outcome: "rejected", Message: "This active MCP grant is no longer available to revoke."}
	}
	callCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	_, err := s.backend.RevokeActor(callCtx, actorID)
	cancel()
	if err != nil {
		return unavailable("revoke_agent_grant")
	}
	return s.workspaceResult(ctx, "revoke_agent_grant", "Agent grant revoked. Existing connections are denied on their next request.")
}

// Actions returns the newest bounded shared audit evidence for one MCP actor.
func (s *Service) Actions(ctx context.Context, actorID string) ActionsResult {
	s.mu.Lock()
	defer s.mu.Unlock()
	actor, ok := s.actor(ctx, actorID)
	if !ok || actor.Builtin || actor.Kind != domain.ActorKindMCP {
		return ActionsResult{Action: "load_agent_actions", Outcome: "rejected", Message: "Agent activity is unavailable for this grant.", Actions: []Action{}}
	}
	callCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	events, err := s.backend.ListAudit(callCtx, domain.AuditQuery{ActorID: actorID, Limit: recentActionLimit})
	cancel()
	if err != nil {
		return ActionsResult{Action: "load_agent_actions", Outcome: "unavailable", Message: "Recent agent actions could not be loaded.", Actions: []Action{}}
	}
	actions := make([]Action, 0, len(events))
	for _, event := range events {
		actions = append(actions, Action{DaemonID: event.DaemonID, Operation: event.Operation, TargetKind: event.TargetKind, TargetID: event.TargetID, Result: string(event.Result), OccurredAt: event.OccurredAt.UTC().Format(time.RFC3339Nano)})
	}
	return ActionsResult{Action: "load_agent_actions", Outcome: "accepted", Message: "Recent agent actions loaded.", Actions: actions}
}

// OpenGuide opens the fixed official MCP guide.
func (s *Service) OpenGuide(ctx context.Context) Result {
	if s.native == nil || s.native.BrowserOpenURL(ctx, guideURL) != nil {
		return unavailable("open_agent_access_guide")
	}
	return Result{Action: "open_agent_access_guide", Outcome: "accepted", Message: "Agent Access guide opened."}
}

func (s *Service) workspaceResult(ctx context.Context, action, message string) Result {
	workspace, err := s.workspace(ctx)
	if err != nil {
		return unavailable(action)
	}
	return Result{Action: action, Outcome: "accepted", Message: message, Workspace: &workspace}
}

func (s *Service) workspace(ctx context.Context) (Workspace, error) {
	callCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	manifest, err := s.backend.Manifest(callCtx)
	if err != nil {
		return Workspace{}, err
	}
	status, err := s.backend.MCPHTTPStatus(callCtx)
	if err != nil {
		return Workspace{}, err
	}
	actors, err := s.backend.ListActors(callCtx)
	if err != nil {
		return Workspace{}, err
	}
	credentials, err := s.backend.ListCredentials(callCtx)
	if err != nil {
		return Workspace{}, err
	}
	credentialByActor := make(map[string]domain.ClientCredential, len(credentials))
	for _, credential := range credentials {
		credentialByActor[credential.ActorID] = credential
	}
	grants := make([]Grant, 0)
	stdioActive := false
	now := s.now()
	for _, actor := range actors {
		if actor.Builtin || actor.Kind != domain.ActorKindMCP {
			continue
		}
		transport := "stdio"
		lastUsed := ""
		fingerprint := ""
		if credential, ok := credentialByActor[actor.ID]; ok {
			transport = "remote_https"
			fingerprint = credential.Fingerprint
			if credential.LastUsedAt != nil {
				lastUsed = credential.LastUsedAt.UTC().Format(time.RFC3339Nano)
			}
		} else if actor.ID == status.ActorID {
			transport = "localhost_http"
			if status.LastAccessedAt != nil {
				lastUsed = status.LastAccessedAt.UTC().Format(time.RFC3339Nano)
			}
		}
		if transport == "stdio" {
			events, auditErr := s.backend.ListAudit(callCtx, domain.AuditQuery{ActorID: actor.ID, Limit: 1})
			if auditErr == nil && len(events) > 0 {
				lastUsed = events[0].OccurredAt.UTC().Format(time.RFC3339Nano)
			}
		}
		state := actor.State
		if state == domain.ActorStateActive && !actor.ActiveAt(now) {
			state = domain.ActorStateExpired
		}
		if transport == "stdio" && state == domain.ActorStateActive {
			stdioActive = true
		}
		expires := ""
		if actor.ExpiresAt != nil {
			expires = actor.ExpiresAt.UTC().Format(time.RFC3339Nano)
		}
		grants = append(grants, Grant{ID: actor.ID, ClientName: actor.DisplayName, DaemonID: manifest.InstallationID, DaemonName: manifest.DisplayName, Capability: string(actor.Capability), CapabilityDescription: authorityDescription(actor.Capability), Transport: transport, CreatedAt: actor.CreatedAt.UTC().Format(time.RFC3339Nano), LastUsedAt: lastUsed, ExpiresAt: expires, State: string(state), CredentialFingerprint: fingerprint})
	}
	sort.SliceStable(grants, func(i, j int) bool {
		if grants[i].State != grants[j].State {
			return grants[i].State == string(domain.ActorStateActive)
		}
		return grants[i].CreatedAt > grants[j].CreatedAt
	})
	httpStatus := projectHTTP(status)
	remoteEnabled := slices.Contains(manifest.Capabilities, "remote-mcp")
	mcpState := "off"
	if status.Enabled || remoteEnabled || stdioActive {
		mcpState = "active"
	}
	return Workspace{
		Daemon:           Daemon{ID: manifest.InstallationID, Name: manifest.DisplayName},
		MCPState:         mcpState,
		StdioDescription: "Available on demand when an MCP host launches `gosched mcp serve`. Stdio opens no network listener and does not enable HTTP access.",
		HTTP:             httpStatus,
		Authorities:      []Authority{{Name: "Observe", Status: "available", Description: authorityDescription(domain.CapabilityObserve)}, {Name: "Operate", Status: "available", Description: authorityDescription(domain.CapabilityOperate)}, {Name: "Manage", Status: "available", Description: authorityDescription(domain.CapabilityManage)}},
		Transports: []Transport{
			{ID: "stdio", Name: "Stdio", State: stdioState(stdioActive), Description: "Available only while an MCP host explicitly launches it. No listening port is opened."},
			{ID: "localhost_http", Name: "Localhost HTTP", State: enabledState(status.Enabled), Description: "Optional loopback listener controlled separately below."},
			{ID: "remote_https", Name: "Remote HTTPS", State: enabledState(remoteEnabled), Description: "Separately configured authenticated listener. Grants do not enable or disable it."},
		},
		Grants: grants,
	}, nil
}

func (s *Service) actor(ctx context.Context, id string) (domain.Actor, bool) {
	if id == "" {
		return domain.Actor{}, false
	}
	callCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	actors, err := s.backend.ListActors(callCtx)
	if err != nil {
		return domain.Actor{}, false
	}
	for _, actor := range actors {
		if actor.ID == id {
			return actor, true
		}
	}
	return domain.Actor{}, false
}

func (s *Service) expiration(value string, allowNonExpiring bool) (*time.Time, error) {
	if value == "non-expiring" && allowNonExpiring {
		return nil, nil
	}
	durations := map[string]time.Duration{"1h": time.Hour, "24h": 24 * time.Hour, "7d": 7 * 24 * time.Hour, "30d": 30 * 24 * time.Hour}
	duration, ok := durations[value]
	if !ok {
		return nil, fmt.Errorf("Choose 1 hour, 24 hours, 7 days, 30 days, or deliberate non-expiring access")
	}
	expires := s.now().Add(duration).UTC()
	return &expires, nil
}

func (s *Service) handoff(ctx context.Context, action, message string, result server.MCPHTTPCredentialResponse) Result {
	credential := result.Credential
	result.Credential = ""
	copyCtx, copyCancel := context.WithTimeout(ctx, operationTimeout)
	copyErr := context.Canceled
	if s.native != nil {
		copyErr = s.native.ClipboardSetText(copyCtx, credential)
	}
	copyCancel()
	credential = ""
	if copyErr != nil {
		rollbackCtx, cancel := context.WithTimeout(context.Background(), operationTimeout)
		_, err := s.backend.DisableMCPHTTP(rollbackCtx)
		cancel()
		if err != nil {
			return Result{Action: action, Outcome: "unavailable", Message: "The credential could not be copied and revocation could not be confirmed. Run `gosched mcp http disable` now, then refresh Agent Access."}
		}
		result := s.workspaceResult(ctx, action, "The credential could not be copied, so localhost access was disabled.")
		result.Outcome = "unavailable"
		return result
	}
	return s.workspaceResult(ctx, action, message)
}

func projectHTTP(status server.MCPHTTPStatusResponse) HTTPStatus {
	result := HTTPStatus{Enabled: status.Enabled, Endpoint: status.Endpoint, AllowedOrigins: append([]string{}, status.AllowedOrigins...), CredentialFingerprint: status.CredentialFingerprint, ClientName: status.ClientName, Permission: string(status.Permission), RequestCount: status.RequestCount, RequireConfirmation: status.RequireConfirmation}
	if status.EnabledAt != nil {
		result.EnabledAt = status.EnabledAt.UTC().Format(time.RFC3339Nano)
	}
	if status.LastAccessedAt != nil {
		result.LastAccessedAt = status.LastAccessedAt.UTC().Format(time.RFC3339Nano)
	}
	return result
}

func authorityDescription(capability domain.Capability) string {
	switch capability {
	case domain.CapabilityObserve:
		return "Read bounded scheduler health, task, run, alert, and summary data."
	case domain.CapabilityOperate:
		return "Run, enable, or disable an existing task using exact daemon, task, and request identifiers."
	case domain.CapabilityManage:
		return "High-impact access to create, update, or delete bounded automation definitions. Access-control and credential tools remain excluded."
	default:
		return "Unknown authority."
	}
}

func narrowerOrEqual(current, next domain.Capability) bool {
	rank := map[domain.Capability]int{domain.CapabilityObserve: 1, domain.CapabilityOperate: 2, domain.CapabilityManage: 3}
	return rank[next] > 0 && rank[next] <= rank[current]
}

func enabledState(enabled bool) string {
	if enabled {
		return "active"
	}
	return "off"
}

func stdioState(active bool) string {
	if active {
		return "active"
	}
	return "available_on_demand"
}

func unavailable(action string) Result {
	return Result{Action: action, Outcome: "unavailable", Message: "Agent Access is unavailable. Check the local scheduler connection, then try again."}
}
