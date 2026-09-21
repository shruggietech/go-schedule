package remotemcp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/internal/api/client"
	apis "github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/enrollment"
	"github.com/shruggietech/go-schedule/internal/remote"
	"github.com/shruggietech/go-schedule/internal/store"
)

const testResource = "https://scheduler.example/mcp"

func TestMetadataTokenExchangeAndMCPInitialization(t *testing.T) {
	service, cleanup := testEnrollment(t)
	defer cleanup()
	issued := issueCredential(t, service, domain.ActorKindMCP, domain.CapabilityManage)
	handler := newTestHandler(t, service)

	metadata := request(handler, http.MethodGet, "/.well-known/oauth-protected-resource/mcp", "", "", "")
	if metadata.Code != http.StatusOK {
		t.Fatalf("metadata status=%d body=%s", metadata.Code, metadata.Body.String())
	}
	var document map[string]any
	if err := json.Unmarshal(metadata.Body.Bytes(), &document); err != nil || document["resource"] != testResource {
		t.Fatalf("metadata=%s err=%v", metadata.Body.String(), err)
	}
	authorizationMetadata := request(handler, http.MethodGet, "/.well-known/oauth-authorization-server", "", "", "")
	if authorizationMetadata.Code != http.StatusOK {
		t.Fatalf("authorization metadata status=%d body=%s", authorizationMetadata.Code, authorizationMetadata.Body.String())
	}
	var authorizationDocument map[string]any
	if err := json.Unmarshal(authorizationMetadata.Body.Bytes(), &authorizationDocument); err != nil {
		t.Fatal(err)
	}
	responseTypes, ok := authorizationDocument["response_types_supported"].([]any)
	if !ok || len(responseTypes) != 0 {
		t.Fatalf("response_types_supported=%#v", authorizationDocument["response_types_supported"])
	}

	token := exchange(t, handler, issued.ID, issued.Token, "mcp:manage")
	requestBody := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1"}}}`
	response := request(handler, http.MethodPost, "/mcp", requestBody, "Bearer "+token, "application/json")
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"serverInfo"`) {
		t.Fatalf("initialize status=%d body=%s", response.Code, response.Body.String())
	}
}

func TestOfficialSDKNegotiatesManageSurface(t *testing.T) {
	service, cleanup := testEnrollment(t)
	defer cleanup()
	issued := issueCredential(t, service, domain.ActorKindMCP, domain.CapabilityManage)
	handler := newTestHandler(t, service)
	token := exchange(t, handler, issued.ID, issued.Token, "mcp:manage")

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = "scheduler.example"
		handler.ServeHTTP(w, r)
	}))
	defer server.Close()
	httpClient := server.Client()
	httpClient.Transport = bearerTransport{base: httpClient.Transport, token: token}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	session, err := mcp.NewClient(&mcp.Implementation{Name: "S093 test", Version: "1"}, nil).Connect(ctx, &mcp.StreamableClientTransport{Endpoint: server.URL + "/mcp", HTTPClient: httpClient, DisableStandaloneSSE: true}, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer session.Close()
	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	if len(tools.Tools) != 9 {
		t.Fatalf("manage tools=%d, want 9", len(tools.Tools))
	}
}

func TestAuthoritySurfacesAreMonotonic(t *testing.T) {
	service, cleanup := testEnrollment(t)
	defer cleanup()
	handler := newTestHandler(t, service)
	for _, test := range []struct {
		capability domain.Capability
		tools      int
	}{
		{domain.CapabilityObserve, 0},
		{domain.CapabilityOperate, 3},
		{domain.CapabilityManage, 9},
	} {
		t.Run(string(test.capability), func(t *testing.T) {
			server := handler.newServer(&grant{actorID: "actor", capability: test.capability})
			serverTransport, clientTransport := mcp.NewInMemoryTransports()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			serverSession, err := server.Connect(ctx, serverTransport, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer serverSession.Close()
			clientSession, err := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil).Connect(ctx, clientTransport, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer clientSession.Close()
			result, err := clientSession.ListTools(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Tools) != test.tools {
				t.Fatalf("tools=%d, want %d", len(result.Tools), test.tools)
			}
		})
	}
}

func TestRequestedScopesAreAnUnorderedAuthoritySet(t *testing.T) {
	scopes, capability, ok := requestedScopes("mcp:manage mcp:operate", domain.CapabilityManage)
	if !ok || capability != domain.CapabilityManage || strings.Join(scopes, " ") != "mcp:observe mcp:operate mcp:manage" {
		t.Fatalf("scopes=%v capability=%q ok=%v", scopes, capability, ok)
	}
	if _, _, ok := requestedScopes("mcp:manage mcp:operate", domain.CapabilityOperate); ok {
		t.Fatal("an Operate actor must not gain Manage through scope ordering")
	}
}

func TestTokenRefreshReusesMutationDeduplicationServer(t *testing.T) {
	service, cleanup := testEnrollment(t)
	defer cleanup()
	issued := issueCredential(t, service, domain.ActorKindMCP, domain.CapabilityManage)
	handler := newTestHandler(t, service)
	firstToken := exchange(t, handler, issued.ID, issued.Token, "mcp:manage")
	secondToken := exchange(t, handler, issued.ID, issued.Token, "mcp:manage")
	firstRequest := httptest.NewRequest(http.MethodPost, testResource, nil)
	firstRequest.URL.Scheme, firstRequest.URL.Host = "", ""
	firstInfo, err := handler.verify(context.Background(), firstToken, firstRequest)
	if err != nil {
		t.Fatal(err)
	}
	secondInfo, err := handler.verify(context.Background(), secondToken, firstRequest)
	if err != nil {
		t.Fatal(err)
	}
	firstGrant := firstInfo.Extra["grant"].(*grant)
	secondGrant := secondInfo.Extra["grant"].(*grant)
	if firstGrant.server != secondGrant.server {
		t.Fatal("token refresh must retain the credential-scoped mutation deduplication server")
	}
}

func TestLiveGrantKeepsDeduplicationServerPastCacheLifetime(t *testing.T) {
	service, cleanup := testEnrollment(t)
	defer cleanup()
	issued := issueCredential(t, service, domain.ActorKindMCP, domain.CapabilityManage)
	handler := newTestHandler(t, service)
	handler.config.AccessTokenLifetimeSeconds = 20 * 60
	now := time.Date(2026, time.September, 20, 12, 0, 0, 0, time.UTC)
	handler.now = func() time.Time { return now }
	firstToken := exchange(t, handler, issued.ID, issued.Token, "mcp:manage")
	now = now.Add(serverCacheLifetime + time.Minute)
	secondToken := exchange(t, handler, issued.ID, issued.Token, "mcp:manage")
	request := httptest.NewRequest(http.MethodPost, testResource, nil)
	request.URL.Scheme, request.URL.Host = "", ""
	firstInfo, err := handler.verify(context.Background(), firstToken, request)
	if err != nil {
		t.Fatal(err)
	}
	secondInfo, err := handler.verify(context.Background(), secondToken, request)
	if err != nil {
		t.Fatal(err)
	}
	firstGrant := firstInfo.Extra["grant"].(*grant)
	secondGrant := secondInfo.Extra["grant"].(*grant)
	if firstGrant.server != secondGrant.server {
		t.Fatal("a live grant must keep its credential-scoped deduplication server")
	}
}

type bearerTransport struct {
	base  http.RoundTripper
	token string
}

func (t bearerTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	clone := request.Clone(request.Context())
	clone.Header = request.Header.Clone()
	clone.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(clone)
}

func TestTokenExchangeFailsClosedAndRevocationIsImmediate(t *testing.T) {
	service, cleanup := testEnrollment(t)
	defer cleanup()
	mcpCredential := issueCredential(t, service, domain.ActorKindMCP, domain.CapabilityOperate)
	desktopCredential := issueCredential(t, service, domain.ActorKindDesktop, domain.CapabilityManage)
	pairing, err := service.Create("unexchanged MCP", domain.ActorKindMCP, domain.CapabilityObserve)
	if err != nil {
		t.Fatal(err)
	}
	handler := newTestHandler(t, service)

	for _, test := range []struct {
		name, clientID, secret, resource, scope string
		want                                    int
	}{
		{"wrong client", "different", mcpCredential.Token, testResource, "mcp:observe", http.StatusUnauthorized},
		{"wrong secret", mcpCredential.ID, "not-a-token", testResource, "mcp:observe", http.StatusUnauthorized},
		{"pairing phrase", pairing.ID, pairing.Phrase, testResource, "mcp:observe", http.StatusUnauthorized},
		{"wrong resource", mcpCredential.ID, mcpCredential.Token, "https://other.example/mcp", "mcp:observe", http.StatusBadRequest},
		{"excessive scope", mcpCredential.ID, mcpCredential.Token, testResource, "mcp:manage", http.StatusBadRequest},
		{"non MCP actor", desktopCredential.ID, desktopCredential.Token, testResource, "mcp:observe", http.StatusUnauthorized},
	} {
		t.Run(test.name, func(t *testing.T) {
			form := url.Values{"grant_type": {"client_credentials"}, "resource": {test.resource}, "scope": {test.scope}}
			response := tokenRequest(handler, test.clientID, test.secret, form)
			if response.Code != test.want {
				t.Fatalf("status=%d want=%d body=%s", response.Code, test.want, response.Body.String())
			}
		})
	}

	token := exchange(t, handler, mcpCredential.ID, mcpCredential.Token, "mcp:operate")
	request := httptest.NewRequest(http.MethodPost, testResource, nil)
	request.URL.Scheme, request.URL.Host = "", ""
	if _, err := handler.verify(context.Background(), token, request); err != nil {
		t.Fatalf("verify active token: %v", err)
	}
	rotated, err := service.Rotate(mcpCredential.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := handler.verify(context.Background(), token, request); !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatalf("verify after rotation=%v, want invalid token", err)
	}
	rotatedToken := exchange(t, handler, rotated.ID, rotated.Token, "mcp:operate")
	if _, err := service.Revoke(rotated.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := handler.verify(context.Background(), rotatedToken, request); !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatalf("verify after revocation=%v, want invalid token", err)
	}
}

func TestAccessGrantExpiresAtConfiguredBoundary(t *testing.T) {
	service, cleanup := testEnrollment(t)
	defer cleanup()
	issued := issueCredential(t, service, domain.ActorKindMCP, domain.CapabilityObserve)
	handler := newTestHandler(t, service)
	now := time.Date(2026, 9, 20, 12, 0, 0, 0, time.UTC)
	handler.now = func() time.Time { return now }
	token := exchange(t, handler, issued.ID, issued.Token, "mcp:observe")
	request := httptest.NewRequest(http.MethodPost, testResource, nil)
	request.URL.Scheme, request.URL.Host = "", ""
	now = now.Add(61 * time.Second)
	if _, err := handler.verify(context.Background(), token, request); !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatalf("verify expired grant=%v, want invalid token", err)
	}
}

func TestAccessGrantRejectsChangedActorAuthorityOnNextRequest(t *testing.T) {
	database, err := store.Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	service := enrollment.New(database)
	issued := issueCredential(t, service, domain.ActorKindMCP, domain.CapabilityOperate)
	handler := newTestHandler(t, service)
	token := exchange(t, handler, issued.ID, issued.Token, "mcp:operate")
	request := httptest.NewRequest(http.MethodPost, testResource, nil)
	request.URL.Scheme, request.URL.Host = "", ""
	observe := domain.CapabilityObserve
	if _, err := database.UpdateActor(issued.Actor.ID, store.ActorUpdate{Capability: &observe}); err != nil {
		t.Fatal(err)
	}
	if _, err := handler.verify(context.Background(), token, request); !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatalf("verify after narrowing=%v, want invalid token", err)
	}
}

func TestAccessTokenRejectsQueryCookieAndDuplicateAuthorization(t *testing.T) {
	service, cleanup := testEnrollment(t)
	defer cleanup()
	issued := issueCredential(t, service, domain.ActorKindMCP, domain.CapabilityObserve)
	handler := newTestHandler(t, service)
	token := exchange(t, handler, issued.ID, issued.Token, "")

	for _, mutate := range []func(*http.Request){
		func(r *http.Request) { r.URL.RawQuery = "access_token=" + token },
		func(r *http.Request) { r.Header.Set("Cookie", "access_token="+token) },
		func(r *http.Request) { r.Header.Add("Authorization", "Bearer "+token) },
	} {
		request := httptest.NewRequest(http.MethodPost, testResource, strings.NewReader("{}"))
		request.URL.Scheme, request.URL.Host = "", ""
		request.Header.Set("Authorization", "Bearer "+token)
		request.Header.Set("Content-Type", "application/json")
		mutate(request)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusUnauthorized {
			t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
		}
	}
}

func TestActorBoundClientPreservesPersistentAuthority(t *testing.T) {
	database, err := store.Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	service := enrollment.New(database)
	issued := issueCredential(t, service, domain.ActorKindMCP, domain.CapabilityOperate)
	api := apis.New(database, nil, nil, nil, "", nil)
	api.SetActorResolver(remote.ActorID)
	bound := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.Handler().ServeHTTP(w, r.WithContext(remote.WithActorID(r.Context(), issued.Actor.ID)))
	}))
	capability, err := bound.VerifyAccess(context.Background())
	if err != nil || capability != domain.CapabilityOperate {
		t.Fatalf("capability=%q err=%v", capability, err)
	}
}

func TestDaemonIdentityResetInvalidatesGrant(t *testing.T) {
	database, err := store.Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	service := enrollment.New(database)
	issued := issueCredential(t, service, domain.ActorKindMCP, domain.CapabilityObserve)
	handler := newTestHandler(t, service)
	token := exchange(t, handler, issued.ID, issued.Token, "mcp:observe")
	identity, err := database.DaemonIdentity()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := database.ResetDaemonIdentity(identity.InstallationID); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, testResource, nil)
	request.URL.Scheme, request.URL.Host = "", ""
	if _, err := handler.verify(context.Background(), token, request); !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatalf("verify after identity reset=%v, want invalid token", err)
	}
}

func TestSupportedDeploymentResourceHosts(t *testing.T) {
	service, cleanup := testEnrollment(t)
	defer cleanup()
	for _, resourceURL := range []string{
		"https://10.0.0.2:8443/mcp",
		"https://scheduler.private.example/mcp",
		"https://scheduler.public.example/mcp",
	} {
		t.Run(resourceURL, func(t *testing.T) {
			handler, err := New(config.RemoteMCPConfig{Enabled: true, ResourceURL: resourceURL}, "test", service, func(string) *client.Client {
				return client.NewInProcess(http.NotFoundHandler())
			}, func(string) bool { return true })
			if err != nil {
				t.Fatal(err)
			}
			resource, _ := url.Parse(resourceURL)
			request := httptest.NewRequest(http.MethodGet, resourceURL[:len(resourceURL)-len("/mcp")]+"/.well-known/oauth-protected-resource/mcp", nil)
			request.URL.Scheme, request.URL.Host = "", ""
			request.Host = resource.Host
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != http.StatusOK {
				t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
			}
			request.Host = "wrong.example"
			blocked := httptest.NewRecorder()
			handler.ServeHTTP(blocked, request)
			if blocked.Code != http.StatusBadRequest {
				t.Fatalf("wrong-host status=%d", blocked.Code)
			}
		})
	}
}

func newTestHandler(t *testing.T, service *enrollment.Service) *Handler {
	t.Helper()
	handler, err := New(config.RemoteMCPConfig{Enabled: true, ResourceURL: testResource, AccessTokenLifetimeSeconds: 60}, "test", service, func(string) *client.Client {
		return client.NewInProcess(http.NotFoundHandler())
	}, func(string) bool { return true })
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func testEnrollment(t *testing.T) (*enrollment.Service, func()) {
	t.Helper()
	database, err := store.Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatal(err)
	}
	return enrollment.New(database), func() { _ = database.Close() }
}

func issueCredential(t *testing.T, service *enrollment.Service, kind domain.ActorKind, capability domain.Capability) domain.IssuedCredential {
	t.Helper()
	secret, err := service.Create("test client", kind, capability)
	if err != nil {
		t.Fatal(err)
	}
	issued, err := service.Exchange(secret.ID, secret.Phrase, secret.DaemonID, secret.DisplayName, secret.Kind, secret.Capability)
	if err != nil {
		t.Fatal(err)
	}
	return issued
}

func exchange(t *testing.T, handler *Handler, clientID, secret, scope string) string {
	t.Helper()
	form := url.Values{"grant_type": {"client_credentials"}, "resource": {testResource}}
	if scope != "" {
		form.Set("scope", scope)
	}
	response := tokenRequest(handler, clientID, secret, form)
	if response.Code != http.StatusOK {
		t.Fatalf("token status=%d body=%s", response.Code, response.Body.String())
	}
	var body struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil || body.AccessToken == "" {
		t.Fatalf("token body=%s err=%v", response.Body.String(), err)
	}
	return body.AccessToken
}

func tokenRequest(handler http.Handler, clientID, secret string, form url.Values) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, "https://scheduler.example/oauth/token", strings.NewReader(form.Encode()))
	request.URL.Scheme, request.URL.Host = "", ""
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(clientID, secret)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func request(handler http.Handler, method, path, body, authorization, contentType string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, "https://scheduler.example"+path, strings.NewReader(body))
	request.URL.Scheme, request.URL.Host = "", ""
	if authorization != "" {
		request.Header.Set("Authorization", authorization)
	}
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	request.Header.Set("Accept", "application/json, text/event-stream")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}
