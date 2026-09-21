// Package remotemcp exposes the scheduler MCP server as an explicitly enabled
// OAuth-protected HTTPS resource.
package remotemcp

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/modelcontextprotocol/go-sdk/oauthex"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/enrollment"
	"github.com/shruggietech/go-schedule/internal/mcpmanage"
	"github.com/shruggietech/go-schedule/internal/mcpobserve"
	"github.com/shruggietech/go-schedule/internal/mcpoperate"
)

const maxTokenRequestBytes = 16 << 10

type clientFactory func(actorID string) *client.Client

type grant struct {
	digest          [sha256.Size]byte
	sourceDigest    [sha256.Size]byte
	credentialID    string
	actorID         string
	capability      domain.Capability
	actorCapability domain.Capability
	daemonID        string
	resource        string
	scopes          []string
	expires         time.Time
	server          *mcp.Server
}

// Handler owns memory-only OAuth access grants and the remote MCP transport.
type Handler struct {
	mu          sync.Mutex
	config      config.RemoteMCPConfig
	resource    *url.URL
	issuer      string
	metadataURL string
	enrollment  *enrollment.Service
	client      clientFactory
	allowActor  func(string) bool
	version     string
	grants      map[[sha256.Size]byte]*grant
	now         func() time.Time
	random      io.Reader
	protected   http.Handler
	metadata    http.Handler
}

// New constructs the remote MCP endpoint set. The caller must mount it only on
// the existing hardened remote HTTPS listener.
func New(cfg config.RemoteMCPConfig, version string, service *enrollment.Service, factory clientFactory, allowActor func(string) bool) (*Handler, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	resource, _ := url.Parse(cfg.ResourceURL)
	issuer := resource.Scheme + "://" + resource.Host
	h := &Handler{config: cfg, resource: resource, issuer: issuer, metadataURL: issuer + "/.well-known/oauth-protected-resource/mcp", enrollment: service, client: factory, allowActor: allowActor, version: version, grants: make(map[[sha256.Size]byte]*grant), now: func() time.Time { return time.Now().UTC() }, random: rand.Reader}
	h.metadata = auth.ProtectedResourceMetadataHandler(&oauthex.ProtectedResourceMetadata{Resource: cfg.ResourceURL, AuthorizationServers: []string{issuer}, ScopesSupported: supportedScopes(), BearerMethodsSupported: []string{"header"}, ResourceName: "go-schedule remote MCP"})
	stream := mcp.NewStreamableHTTPHandler(h.serverForRequest, &mcp.StreamableHTTPOptions{Stateless: true, MaxRequestBodyBytes: 1 << 20, PropagateRequestCancellation: true, DisableLocalhostProtection: true})
	limitedStream := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info := auth.TokenInfoFromContext(r.Context())
		if info == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		entry, _ := info.Extra["grant"].(*grant)
		if entry == nil || h.allowActor != nil && !h.allowActor(entry.credentialID) {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		stream.ServeHTTP(w, r)
	})
	h.protected = auth.RequireBearerToken(h.verify, &auth.RequireBearerTokenOptions{Scopes: []string{"mcp:observe"}, ResourceMetadataURL: h.metadataURL})(limitedStream)
	return h, nil
}

// Paths returns the exact public paths owned by this handler.
func (h *Handler) Paths() []string {
	return []string{h.resource.Path, "/oauth/token", "/.well-known/oauth-authorization-server", "/.well-known/oauth-protected-resource/mcp"}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if !h.validHost(r) {
		http.Error(w, "invalid host", http.StatusBadRequest)
		return
	}
	switch r.URL.Path {
	case "/.well-known/oauth-protected-resource/mcp":
		h.metadata.ServeHTTP(w, r)
	case "/.well-known/oauth-authorization-server":
		h.serveAuthorizationMetadata(w, r)
	case "/oauth/token":
		h.serveToken(w, r)
	case h.resource.Path:
		if r.URL.Query().Has("access_token") || r.Header.Get("Cookie") != "" || len(r.Header.Values("Authorization")) != 1 {
			h.protected.ServeHTTP(w, r.Clone(context.WithValue(r.Context(), rejectedCredentialKey{}, true)))
			return
		}
		h.protected.ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) validHost(r *http.Request) bool {
	return strings.EqualFold(r.Host, h.resource.Host) && r.URL.Scheme == "" && r.URL.Host == ""
}

func (h *Handler) serveAuthorizationMetadata(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"issuer": h.issuer, "token_endpoint": h.issuer + "/oauth/token", "grant_types_supported": []string{"client_credentials"}, "token_endpoint_auth_methods_supported": []string{"client_secret_basic"}, "scopes_supported": supportedScopes()})
}

func (h *Handler) serveToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	contentType, _, contentTypeErr := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if contentTypeErr != nil || contentType != "application/x-www-form-urlencoded" || r.URL.RawQuery != "" {
		oauthError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxTokenRequestBytes)
	if err := r.ParseForm(); err != nil || !onlyFormFields(r.PostForm, "grant_type", "resource", "scope") {
		oauthError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	clientID, secret, ok := r.BasicAuth()
	if !ok || len(r.Header.Values("Authorization")) != 1 || clientID == "" || secret == "" {
		oauthError(w, http.StatusUnauthorized, "invalid_client")
		return
	}
	if one(r.PostForm, "grant_type") != "client_credentials" {
		oauthError(w, http.StatusBadRequest, "unsupported_grant_type")
		return
	}
	if one(r.PostForm, "resource") != h.config.ResourceURL {
		oauthError(w, http.StatusBadRequest, "invalid_target")
		return
	}
	credential, actor, err := h.enrollment.Authenticate(secret)
	if err != nil || credential.ID != clientID || actor.Kind != domain.ActorKindMCP {
		oauthError(w, http.StatusUnauthorized, "invalid_client")
		return
	}
	scopes, capability, ok := requestedScopes(one(r.PostForm, "scope"), actor.Capability)
	if !ok {
		oauthError(w, http.StatusBadRequest, "invalid_scope")
		return
	}
	rawSecret, _ := base64.RawURLEncoding.DecodeString(secret)
	sourceDigest := sha256.Sum256(rawSecret)
	accessRaw := make([]byte, 32)
	if _, err := io.ReadFull(h.random, accessRaw); err != nil {
		oauthError(w, http.StatusInternalServerError, "server_error")
		return
	}
	accessToken := base64.RawURLEncoding.EncodeToString(accessRaw)
	digest := sha256.Sum256(accessRaw)
	identity, err := h.enrollment.Identity()
	if err != nil {
		oauthError(w, http.StatusInternalServerError, "server_error")
		return
	}
	entry := &grant{digest: digest, sourceDigest: sourceDigest, credentialID: credential.ID, actorID: actor.ID, capability: capability, actorCapability: actor.Capability, daemonID: identity.InstallationID, resource: h.config.ResourceURL, scopes: scopes, expires: h.now().Add(h.config.AccessTokenLifetime())}
	entry.server = h.newServer(entry)
	h.mu.Lock()
	h.pruneLocked()
	if len(h.grants) >= 1024 {
		h.evictOldestLocked()
	}
	h.grants[digest] = entry
	h.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"access_token": accessToken, "token_type": "Bearer", "expires_in": int(h.config.AccessTokenLifetime().Seconds()), "scope": strings.Join(scopes, " ")})
}

type rejectedCredentialKey struct{}

func (h *Handler) verify(_ context.Context, token string, r *http.Request) (*auth.TokenInfo, error) {
	if rejected, _ := r.Context().Value(rejectedCredentialKey{}).(bool); rejected || !h.validHost(r) || r.URL.Path != h.resource.Path || strings.ContainsAny(token, " \t\r\n") {
		return nil, auth.ErrInvalidToken
	}
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(raw) != 32 {
		return nil, auth.ErrInvalidToken
	}
	digest := sha256.Sum256(raw)
	h.mu.Lock()
	h.pruneLocked()
	entry := h.grants[digest]
	h.mu.Unlock()
	if entry == nil || entry.resource != h.config.ResourceURL || !entry.expires.After(h.now()) {
		return nil, auth.ErrInvalidToken
	}
	credential, actor, err := h.enrollment.AuthenticateDigest(entry.sourceDigest[:])
	if err != nil || credential.ID != entry.credentialID || actor.ID != entry.actorID || actor.Kind != domain.ActorKindMCP || actor.Capability != entry.actorCapability || !actor.Capability.Allows(entry.capability) {
		return nil, auth.ErrInvalidToken
	}
	identity, identityErr := h.enrollment.Identity()
	if identityErr != nil || identity.InstallationID != entry.daemonID {
		return nil, auth.ErrInvalidToken
	}
	return &auth.TokenInfo{Scopes: append([]string(nil), entry.scopes...), Expiration: entry.expires, UserID: entry.actorID, Extra: map[string]any{"grant": entry}}, nil
}

func (h *Handler) serverForRequest(r *http.Request) *mcp.Server {
	info := auth.TokenInfoFromContext(r.Context())
	if info == nil {
		return nil
	}
	entry, _ := info.Extra["grant"].(*grant)
	if entry == nil {
		return nil
	}
	return entry.server
}

func (h *Handler) newServer(entry *grant) *mcp.Server {
	api := h.client(entry.actorID)
	server := mcpobserve.NewServer(api, h.version)
	if entry.capability.Allows(domain.CapabilityOperate) {
		mcpoperate.AddTools(server, mcpoperate.New(api))
	}
	if entry.capability.Allows(domain.CapabilityManage) {
		mcpmanage.AddTools(server, mcpmanage.New(api, false))
	}
	return server
}

func requestedScopes(value string, allowed domain.Capability) ([]string, domain.Capability, bool) {
	if value == "" {
		value = "mcp:observe"
	}
	wanted := strings.Fields(value)
	if len(wanted) == 0 {
		return nil, "", false
	}
	seen := make(map[string]bool)
	capability := domain.CapabilityObserve
	for _, scope := range wanted {
		if seen[scope] {
			return nil, "", false
		}
		seen[scope] = true
		switch scope {
		case "mcp:observe":
		case "mcp:operate":
			capability = domain.CapabilityOperate
		case "mcp:manage":
			capability = domain.CapabilityManage
		default:
			return nil, "", false
		}
	}
	if !allowed.Allows(capability) {
		return nil, "", false
	}
	// Grants carry all lower scopes so ordinary scope checks remain monotonic.
	result := []string{"mcp:observe"}
	if capability.Allows(domain.CapabilityOperate) {
		result = append(result, "mcp:operate")
	}
	if capability.Allows(domain.CapabilityManage) {
		result = append(result, "mcp:manage")
	}
	return result, capability, true
}

func supportedScopes() []string { return []string{"mcp:observe", "mcp:operate", "mcp:manage"} }

func onlyFormFields(values url.Values, allowed ...string) bool {
	allowedSet := make(map[string]bool, len(allowed))
	for _, name := range allowed {
		allowedSet[name] = true
	}
	for name, entries := range values {
		if !allowedSet[name] || len(entries) != 1 {
			return false
		}
	}
	return true
}

func one(values url.Values, name string) string {
	entries := values[name]
	if len(entries) != 1 {
		return ""
	}
	return entries[0]
}

func (h *Handler) pruneLocked() {
	now := h.now()
	for digest, entry := range h.grants {
		if !entry.expires.After(now) {
			delete(h.grants, digest)
		}
	}
}

func (h *Handler) evictOldestLocked() {
	type candidate struct {
		digest [sha256.Size]byte
		expiry time.Time
	}
	values := make([]candidate, 0, len(h.grants))
	for digest, entry := range h.grants {
		values = append(values, candidate{digest, entry.expires})
	}
	sort.Slice(values, func(i, j int) bool { return values[i].expiry.Before(values[j].expiry) })
	if len(values) > 0 {
		delete(h.grants, values[0].digest)
	}
}

func oauthError(w http.ResponseWriter, status int, code string) {
	if status == http.StatusUnauthorized {
		w.Header().Set("WWW-Authenticate", `Basic realm="go-schedule MCP"`)
	}
	writeJSON(w, status, map[string]string{"error": code})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
