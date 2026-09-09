package agentaccess

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

const operationTimeout = 2 * time.Second
const guideURL = "https://shruggietech.github.io/go-schedule/mcp/"

// Service serializes Agent Access operations and keeps credentials behind the native boundary.
type Service struct {
	mu      sync.Mutex
	backend Backend
	native  Native
}

// NewService creates an Agent Access service.
func NewService(backend Backend, native Native) *Service {
	return &Service{backend: backend, native: native}
}

// Workspace returns authoritative non-secret local MCP state.
func (s *Service) Workspace(ctx context.Context) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	status, err := s.status(ctx)
	if err != nil {
		return unavailable("load_agent_access")
	}
	workspace := project(status)
	return Result{Action: "load_agent_access", Outcome: "accepted", Message: "Agent access loaded.", Workspace: &workspace}
}

// Enable starts one named localhost client and copies its credential once.
func (s *Service) Enable(ctx context.Context, draft EnableDraft) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	if strings.TrimSpace(draft.ClientName) == "" {
		return Result{Action: "enable_agent_access", Outcome: "rejected", Message: "Enter a client name."}
	}
	callCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	result, err := s.backend.EnableMCPHTTP(callCtx, server.MCPHTTPEnableRequest{Port: draft.Port, AllowedOrigins: draft.AllowedOrigins, ClientName: draft.ClientName})
	cancel()
	if err != nil {
		return Result{Action: "enable_agent_access", Outcome: "rejected", Message: "Localhost access could not be enabled. Check the name, port, and origins, then try again."}
	}
	return s.handoff(ctx, "enable_agent_access", "Localhost access enabled. The credential was copied once.", result)
}

// Rotate replaces and copies the active credential while preserving listener policy.
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

// Revoke invalidates the credential and closes the optional listener.
func (s *Service) Revoke(ctx context.Context) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	callCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	status, err := s.backend.DisableMCPHTTP(callCtx)
	cancel()
	if err != nil {
		return unavailable("revoke_agent_access")
	}
	workspace := project(status)
	return Result{Action: "revoke_agent_access", Outcome: "accepted", Message: "Localhost access revoked and listener stopped.", Workspace: &workspace}
}

// OpenGuide opens the fixed official MCP guide.
func (s *Service) OpenGuide(ctx context.Context) Result {
	if s.native == nil || s.native.BrowserOpenURL(ctx, guideURL) != nil {
		return unavailable("open_agent_access_guide")
	}
	return Result{Action: "open_agent_access_guide", Outcome: "accepted", Message: "Agent Access guide opened."}
}

func (s *Service) status(ctx context.Context) (server.MCPHTTPStatusResponse, error) {
	callCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	return s.backend.MCPHTTPStatus(callCtx)
}

func (s *Service) handoff(ctx context.Context, action, message string, result server.MCPHTTPCredentialResponse) Result {
	credential := result.Credential
	result.Credential = ""
	if s.native == nil || s.native.ClipboardSetText(ctx, credential) != nil {
		rollbackCtx, cancel := context.WithTimeout(context.Background(), operationTimeout)
		_, _ = s.backend.DisableMCPHTTP(rollbackCtx)
		cancel()
		return Result{Action: action, Outcome: "unavailable", Message: "The credential could not be copied, so localhost access was disabled."}
	}
	workspace := project(result.MCPHTTPStatusResponse)
	return Result{Action: action, Outcome: "accepted", Message: message, Workspace: &workspace}
}

func project(status server.MCPHTTPStatusResponse) Workspace {
	httpStatus := HTTPStatus{Enabled: status.Enabled, Endpoint: status.Endpoint, AllowedOrigins: append([]string{}, status.AllowedOrigins...), CredentialFingerprint: status.CredentialFingerprint, ClientName: status.ClientName, RequestCount: status.RequestCount}
	if status.EnabledAt != nil {
		httpStatus.EnabledAt = status.EnabledAt.UTC().Format(time.RFC3339Nano)
	}
	if status.LastAccessedAt != nil {
		httpStatus.LastAccessedAt = status.LastAccessedAt.UTC().Format(time.RFC3339Nano)
	}
	return Workspace{
		StdioDescription: "Available on demand when an MCP host launches `gosched mcp serve`. Stdio opens no network listener and inherits this user's local daemon access.",
		HTTP:             httpStatus,
		Authorities:      []Authority{{Name: "Observe", Status: "available", Description: "Read bounded scheduler health, task, run, alert, and summary data."}, {Name: "Operate", Status: "future", Description: "Unavailable. Agents cannot run or change scheduled work."}, {Name: "Manage", Status: "future", Description: "Unavailable. Agents cannot change configuration or access controls."}},
	}
}

func unavailable(action string) Result {
	return Result{Action: action, Outcome: "unavailable", Message: "Agent Access is unavailable. Check the local scheduler connection, then try again."}
}
