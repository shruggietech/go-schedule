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

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/mcpobserve"
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
	opMu    sync.Mutex
	mu      sync.RWMutex
	reader  observeReader
	version string
	log     *slog.Logger
	status  server.MCPHTTPStatusResponse
	digest  [sha256.Size]byte
	http    *http.Server
	listen  net.Listener
	now     func() time.Time
}

// New constructs a disabled manager.
func New(reader observeReader, version string, log *slog.Logger) *Manager {
	if log == nil {
		log = slog.New(slog.DiscardHandler)
	}
	return &Manager{reader: reader, version: version, log: log, status: disabledStatus(), now: time.Now}
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
func (m *Manager) Enable(_ context.Context, req server.MCPHTTPEnableRequest) (server.MCPHTTPCredentialResponse, error) {
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
	credential, digest, fingerprint, err := newCredential()
	if err != nil {
		return server.MCPHTTPCredentialResponse{}, fmt.Errorf("generate localhost MCP credential: %w", err)
	}
	host := net.JoinHostPort("127.0.0.1", strconv.Itoa(req.Port))
	listener, err := net.Listen("tcp4", host)
	if err != nil {
		return server.MCPHTTPCredentialResponse{}, conflictError("cannot enable localhost MCP because the selected port is unavailable")
	}
	enabledAt := m.now().UTC()
	status := server.MCPHTTPStatusResponse{
		Enabled:               true,
		Endpoint:              "http://" + host + "/mcp",
		AllowedOrigins:        origins,
		CredentialFingerprint: fingerprint,
		EnabledAt:             &enabledAt,
		ClientName:            clientName,
	}
	httpServer := &http.Server{Handler: streamableHandler(m), ReadHeaderTimeout: 5 * time.Second}
	m.mu.Lock()
	m.status = status
	m.digest = digest
	m.http = httpServer
	m.listen = listener
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
		m.clearLocked()
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
	m.clearLocked()
	m.mu.Unlock()
	if listener != nil {
		_ = listener.Close()
	}
	shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, net.ErrClosed) {
		_ = httpServer.Close()
		return disabledStatus(), fmt.Errorf("stop localhost MCP listener: %w", err)
	}
	m.log.Info("localhost MCP disabled")
	return disabledStatus(), nil
}

// Shutdown revokes and stops the listener during daemon termination.
func (m *Manager) Shutdown(ctx context.Context) error {
	_, err := m.Disable(ctx)
	return err
}

func (m *Manager) clearLocked() {
	m.status = disabledStatus()
	m.digest = [sha256.Size]byte{}
	m.http = nil
	m.listen = nil
}

func (m *Manager) newObserveServer() *mcp.Server {
	return mcpobserve.NewServer(m.reader, m.version)
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
