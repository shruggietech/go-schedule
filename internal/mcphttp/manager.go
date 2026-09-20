package mcphttp

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/mcpmanage"
	"github.com/shruggietech/go-schedule/internal/mcpobserve"
	"github.com/shruggietech/go-schedule/internal/mcpoperate"
)

const (
	shutdownTimeout    = 5 * time.Second
	defaultClientName  = "Local MCP client"
	maxClientNameBytes = 64
)

type observeReader interface {
	Health(context.Context) (server.HealthResponse, error)
	ListTaskObservations(context.Context, string, string, bool, int, int, int) ([]server.TaskObservationResponse, error)
	ListRunsPage(context.Context, string, int, int, int) ([]domain.Run, error)
	ListAlertsPage(context.Context, bool, int, int, int) ([]domain.Alert, error)
}

// Manager owns one runtime-only listener and its credential digest.
type Manager struct {
	opMu          sync.Mutex
	mu            sync.RWMutex
	reader        observeReader
	version       string
	log           *slog.Logger
	status        server.MCPHTTPStatusResponse
	digest        [sha256.Size]byte
	http          *http.Server
	listen        net.Listener
	now           func() time.Time
	sessionCreate func(context.Context, string, domain.Capability) (server.MCPSessionCredentialResponse, error)
	sessionRevoke func(context.Context, string) error
	operateFor    func(string) *mcpoperate.Executor
	manageFor     func(string, bool) *mcpmanage.Executor
	sessionID     string
	operate       *mcpoperate.Executor
	manage        *mcpmanage.Executor
}

// New constructs a disabled manager.
func New(reader observeReader, version string, log *slog.Logger) *Manager {
	if log == nil {
		log = slog.New(slog.DiscardHandler)
	}
	return &Manager{reader: reader, version: version, log: log, status: disabledStatus(), now: time.Now}
}

// NewWithSessions enables explicit Operate and Manage sessions in addition to the
// Observe-only default.
func NewWithSessions(local *client.Client, version string, log *slog.Logger) *Manager {
	manager := New(local, version, log)
	manager.sessionCreate = local.CreateMCPSession
	manager.sessionRevoke = local.RevokeMCPSession
	manager.operateFor = func(secret string) *mcpoperate.Executor { return mcpoperate.New(local.WithMCPSession(secret)) }
	manager.manageFor = func(secret string, confirm bool) *mcpmanage.Executor {
		return mcpmanage.New(local.WithMCPSession(secret), confirm)
	}
	return manager
}

func disabledStatus() server.MCPHTTPStatusResponse {
	return server.MCPHTTPStatusResponse{AllowedOrigins: []string{}}
}

// Status returns a non-secret immutable snapshot.
func (m *Manager) Status() server.MCPHTTPStatusResponse {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return cloneStatus(m.status)
}

// Enable validates, binds, issues a credential, and atomically publishes the listener.
func (m *Manager) Enable(ctx context.Context, req server.MCPHTTPEnableRequest) (server.MCPHTTPCredentialResponse, error) {
	m.opMu.Lock()
	defer m.opMu.Unlock()
	if req.Port < 1 || req.Port > 65535 {
		return server.MCPHTTPCredentialResponse{}, validationError("port", "port must be between 1 and 65535")
	}
	origins, err := normalizeOrigins(req.AllowedOrigins)
	if err != nil {
		return server.MCPHTTPCredentialResponse{}, err
	}
	clientName, err := normalizeClientName(req.ClientName)
	if err != nil {
		return server.MCPHTTPCredentialResponse{}, err
	}
	m.mu.RLock()
	enabled := m.status.Enabled
	m.mu.RUnlock()
	if enabled {
		return server.MCPHTTPCredentialResponse{}, conflictError("localhost MCP is already enabled; disable it before changing port or origins")
	}
	permission := req.Permission
	if permission == "" {
		permission = domain.CapabilityObserve
	}
	if permission != domain.CapabilityObserve && permission != domain.CapabilityOperate && permission != domain.CapabilityManage {
		return server.MCPHTTPCredentialResponse{}, validationError("permission", "permission must be observe, operate, or manage")
	}
	if req.RequireConfirmation && permission != domain.CapabilityManage {
		return server.MCPHTTPCredentialResponse{}, validationError("require_confirmation", "confirmation policy requires Manage permission")
	}
	if permission.Allows(domain.CapabilityOperate) && (m.sessionCreate == nil || m.sessionRevoke == nil || m.operateFor == nil) {
		return server.MCPHTTPCredentialResponse{}, conflictError("MCP mutation sessions are unavailable")
	}
	if permission == domain.CapabilityManage && m.manageFor == nil {
		return server.MCPHTTPCredentialResponse{}, conflictError("MCP Manage sessions are unavailable")
	}
	credential, digest, fingerprint, err := newCredential()
	if err != nil {
		return server.MCPHTTPCredentialResponse{}, fmt.Errorf("generate localhost MCP credential: %w", err)
	}
	host := net.JoinHostPort("127.0.0.1", strconv.Itoa(req.Port))
	listener, err := net.Listen("tcp4", host)
	if err != nil {
		return server.MCPHTTPCredentialResponse{}, conflictError("cannot enable localhost MCP because the selected port is unavailable")
	}
	var sessionID, actorID string
	var operate *mcpoperate.Executor
	var manage *mcpmanage.Executor
	if permission.Allows(domain.CapabilityOperate) {
		session, sessionErr := m.sessionCreate(ctx, clientName, permission)
		if sessionErr != nil {
			_ = listener.Close()
			return server.MCPHTTPCredentialResponse{}, fmt.Errorf("create MCP mutation session: %w", sessionErr)
		}
		sessionID, actorID = session.ID, session.ActorID
		operate = m.operateFor(session.Credential)
		if permission == domain.CapabilityManage {
			manage = m.manageFor(session.Credential, req.RequireConfirmation)
		}
	}
	enabledAt := m.now().UTC()
	status := server.MCPHTTPStatusResponse{
		Enabled:               true,
		Endpoint:              "http://" + host + "/mcp",
		AllowedOrigins:        origins,
		CredentialFingerprint: fingerprint,
		EnabledAt:             &enabledAt,
		ClientName:            clientName,
		Permission:            permission,
		ActorID:               actorID,
		RequireConfirmation:   req.RequireConfirmation,
	}
	httpServer := &http.Server{Handler: streamableHandler(m), ReadHeaderTimeout: 5 * time.Second}
	m.mu.Lock()
	m.status = status
	m.digest = digest
	m.http = httpServer
	m.listen = listener
	m.sessionID = sessionID
	m.operate = operate
	m.manage = manage
	m.mu.Unlock()
	go m.serve(httpServer, listener)
	m.log.Info("localhost MCP enabled", "endpoint", status.Endpoint, "credential_fingerprint", fingerprint)
	return server.MCPHTTPCredentialResponse{MCPHTTPStatusResponse: cloneStatus(status), Credential: credential}, nil
}

func (m *Manager) serve(httpServer *http.Server, listener net.Listener) {
	err := httpServer.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) && !errors.Is(err, net.ErrClosed) {
		m.log.Error("localhost MCP listener stopped unexpectedly", "error", err)
	}
	m.mu.Lock()
	if m.http == httpServer {
		sessionID := m.clearLocked()
		m.mu.Unlock()
		m.revokeRuntimeSession(sessionID)
		return
	}
	m.mu.Unlock()
}

// Rotate atomically replaces the active digest without restarting the listener.
func (m *Manager) Rotate(_ context.Context) (server.MCPHTTPCredentialResponse, error) {
	m.opMu.Lock()
	defer m.opMu.Unlock()
	credential, digest, fingerprint, err := newCredential()
	if err != nil {
		return server.MCPHTTPCredentialResponse{}, fmt.Errorf("generate localhost MCP credential: %w", err)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.status.Enabled {
		return server.MCPHTTPCredentialResponse{}, conflictError("localhost MCP is disabled")
	}
	m.digest = digest
	m.status.CredentialFingerprint = fingerprint
	m.status.LastAccessedAt = nil
	m.status.RequestCount = 0
	status := cloneStatus(m.status)
	m.log.Info("localhost MCP credential rotated", "credential_fingerprint", fingerprint)
	return server.MCPHTTPCredentialResponse{MCPHTTPStatusResponse: status, Credential: credential}, nil
}

// Disable revokes access first and then performs bounded graceful shutdown.
func (m *Manager) Disable(ctx context.Context) (server.MCPHTTPStatusResponse, error) {
	m.opMu.Lock()
	defer m.opMu.Unlock()
	m.mu.Lock()
	if !m.status.Enabled {
		m.mu.Unlock()
		return disabledStatus(), nil
	}
	httpServer := m.http
	listener := m.listen
	sessionID := m.clearLocked()
	m.mu.Unlock()
	if listener != nil {
		_ = listener.Close()
	}
	shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, net.ErrClosed) {
		closeErr := httpServer.Close()
		m.revokeRuntimeSession(sessionID)
		if closeErr != nil && !errors.Is(closeErr, net.ErrClosed) && !errors.Is(closeErr, http.ErrServerClosed) {
			return disabledStatus(), fmt.Errorf("force stop localhost MCP listener after graceful shutdown failed: %w", closeErr)
		}
		m.log.Warn("localhost MCP required forced shutdown", "error", err)
		return disabledStatus(), nil
	}
	m.revokeRuntimeSession(sessionID)
	m.log.Info("localhost MCP disabled")
	return disabledStatus(), nil
}

// Shutdown revokes and stops the listener during daemon termination.
func (m *Manager) Shutdown(ctx context.Context) error {
	_, err := m.Disable(ctx)
	return err
}

func (m *Manager) clearLocked() string {
	sessionID := m.sessionID
	m.status = disabledStatus()
	m.digest = [sha256.Size]byte{}
	m.http = nil
	m.listen = nil
	m.sessionID = ""
	m.operate = nil
	m.manage = nil
	return sessionID
}

func (m *Manager) newMCPServer() *mcp.Server {
	server := mcpobserve.NewServer(m.reader, m.version)
	m.mu.RLock()
	operate := m.operate
	manage := m.manage
	permission := m.status.Permission
	m.mu.RUnlock()
	if permission.Allows(domain.CapabilityOperate) && operate != nil {
		mcpoperate.AddTools(server, operate)
	}
	if permission == domain.CapabilityManage && manage != nil {
		mcpmanage.AddTools(server, manage)
	}
	return server
}

func (m *Manager) revokeRuntimeSession(id string) {
	if id == "" || m.sessionRevoke == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := m.sessionRevoke(ctx, id); err != nil {
		m.log.Error("revoke localhost MCP session", "error", err)
	}
}

func newCredential() (string, [sha256.Size]byte, string, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", [sha256.Size]byte{}, "", err
	}
	credential := base64.RawURLEncoding.EncodeToString(raw[:])
	digest := sha256.Sum256([]byte(credential))
	fingerprint := base64.RawURLEncoding.EncodeToString(digest[:9])
	return credential, digest, fingerprint, nil
}

func cloneStatus(status server.MCPHTTPStatusResponse) server.MCPHTTPStatusResponse {
	status.AllowedOrigins = append([]string{}, status.AllowedOrigins...)
	status.EnabledAt = cloneTime(status.EnabledAt)
	status.LastAccessedAt = cloneTime(status.LastAccessedAt)
	return status
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func normalizeClientName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultClientName, nil
	}
	if !utf8.ValidString(value) || len(value) > maxClientNameBytes {
		return "", validationError("client_name", "client name must be valid UTF-8 and no more than 64 bytes")
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return "", validationError("client_name", "client name must not contain control characters")
		}
	}
	return value, nil
}

func validationError(field, message string) error {
	return &server.MCPHTTPControlError{Code: server.CodeValidation, Field: field, Message: message}
}

func conflictError(message string) error {
	return &server.MCPHTTPControlError{Code: server.CodeConflict, Message: message}
}
