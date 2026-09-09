package mcphttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type fakeReader struct{}

func (fakeReader) Health(context.Context) (server.HealthResponse, error) {
	return server.HealthResponse{Status: "ok", Version: "test"}, nil
}
func (fakeReader) ListTaskObservations(context.Context, string, string, bool, int, int, int) ([]server.TaskObservationResponse, error) {
	return []server.TaskObservationResponse{}, nil
}
func (fakeReader) ListRunsPage(context.Context, string, int, int, int) ([]domain.Run, error) {
	return []domain.Run{}, nil
}
func (fakeReader) ListAlertsPage(context.Context, bool, int, int, int) ([]domain.Alert, error) {
	return []domain.Alert{}, nil
}

type cancelReader struct {
	fakeReader
	started  chan struct{}
	canceled chan struct{}
}

func (r *cancelReader) Health(ctx context.Context) (server.HealthResponse, error) {
	close(r.started)
	<-ctx.Done()
	close(r.canceled)
	return server.HealthResponse{}, ctx.Err()
}

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := ln.Close(); err != nil {
		t.Fatal(err)
	}
	return port
}

func enableManager(t *testing.T, origins ...string) (*Manager, server.MCPHTTPCredentialResponse) {
	t.Helper()
	manager := New(fakeReader{}, "test", nil)
	result, err := manager.Enable(context.Background(), server.MCPHTTPEnableRequest{Port: freePort(t), AllowedOrigins: origins})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.Shutdown(context.Background()) })
	return manager, result
}

func TestManagerLifecycleIsRuntimeOnlyAndCredentialsRotate(t *testing.T) {
	manager := New(fakeReader{}, "test", nil)
	if status := manager.Status(); status.Enabled || status.AllowedOrigins == nil {
		t.Fatalf("initial status = %+v", status)
	}
	port := freePort(t)
	first, err := manager.Enable(context.Background(), server.MCPHTTPEnableRequest{Port: port, AllowedOrigins: []string{"HTTP://127.0.0.1:3000/", "http://127.0.0.1:3000"}})
	if err != nil {
		t.Fatal(err)
	}
	if !first.Enabled || first.Credential == "" || first.CredentialFingerprint == "" || len(first.AllowedOrigins) != 1 {
		t.Fatalf("enable = %+v", first)
	}
	if _, err := manager.Enable(context.Background(), server.MCPHTTPEnableRequest{Port: freePort(t)}); err == nil {
		t.Fatal("second enable succeeded")
	}
	rotated, err := manager.Rotate(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if rotated.Credential == first.Credential || rotated.Endpoint != first.Endpoint || rotated.CredentialFingerprint == first.CredentialFingerprint {
		t.Fatalf("rotation did not replace only credential: first=%+v rotated=%+v", first, rotated)
	}
	status := manager.Status()
	encoded, _ := json.Marshal(status)
	if strings.Contains(string(encoded), first.Credential) || strings.Contains(string(encoded), rotated.Credential) {
		t.Fatalf("status leaked credential: %s", encoded)
	}
	if _, err := manager.Disable(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := manager.Disable(context.Background()); err != nil {
		t.Fatal(err)
	}
	if manager.Status().Enabled {
		t.Fatal("manager still enabled")
	}
	if _, err := manager.Rotate(context.Background()); err == nil {
		t.Fatal("rotation succeeded while disabled")
	}
	if conn, err := net.DialTimeout("tcp4", net.JoinHostPort("127.0.0.1", stringPort(port)), 100*time.Millisecond); err == nil {
		_ = conn.Close()
		t.Fatal("listener remained open")
	}
}

func TestEnableValidationAndPortConflictAreAtomic(t *testing.T) {
	manager := New(fakeReader{}, "test", nil)
	invalid := []server.MCPHTTPEnableRequest{
		{Port: 0},
		{Port: 65536},
		{Port: 1234, AllowedOrigins: []string{"null"}},
		{Port: 1234, AllowedOrigins: []string{"http://localhost:3000"}},
		{Port: 1234, AllowedOrigins: []string{"http://127.0.0.1"}},
		{Port: 1234, AllowedOrigins: []string{"http://127.0.0.1:3000/path"}},
		{Port: 1234, AllowedOrigins: []string{"ftp://127.0.0.1:3000"}},
	}
	for _, req := range invalid {
		if _, err := manager.Enable(context.Background(), req); err == nil {
			t.Fatalf("invalid request succeeded: %+v", req)
		}
		if manager.Status().Enabled {
			t.Fatal("invalid request published enabled state")
		}
	}
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	if _, err := manager.Enable(context.Background(), server.MCPHTTPEnableRequest{Port: ln.Addr().(*net.TCPAddr).Port}); err == nil {
		t.Fatal("occupied port succeeded")
	}
	if manager.Status().Enabled {
		t.Fatal("port conflict published enabled state")
	}
}

func TestSecurityMatricesRejectBeforeMCP(t *testing.T) {
	manager, result := enableManager(t, "http://127.0.0.1:3000")
	host := strings.TrimSuffix(strings.TrimPrefix(result.Endpoint, "http://"), "/mcp")
	initBody := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1"}}}`
	tests := []struct {
		name    string
		host    string
		origins []string
		auth    []string
		forward string
		status  int
	}{
		{name: "valid native", host: host, auth: []string{"Bearer " + result.Credential}, status: http.StatusOK},
		{name: "valid browser", host: host, origins: []string{"http://127.0.0.1:3000"}, auth: []string{"Bearer " + result.Credential}, status: http.StatusOK},
		{name: "localhost host", host: "localhost:" + strings.Split(host, ":")[1], auth: []string{"Bearer " + result.Credential}, status: http.StatusForbidden},
		{name: "alternate loopback", host: "127.0.0.2:" + strings.Split(host, ":")[1], auth: []string{"Bearer " + result.Credential}, status: http.StatusForbidden},
		{name: "host trailing dot", host: "127.0.0.1.:" + strings.Split(host, ":")[1], auth: []string{"Bearer " + result.Credential}, status: http.StatusForbidden},
		{name: "host without port", host: "127.0.0.1", auth: []string{"Bearer " + result.Credential}, status: http.StatusForbidden},
		{name: "forwarded host ignored", host: host, auth: []string{"Bearer " + result.Credential}, forward: "evil.test", status: http.StatusOK},
		{name: "forwarded host cannot repair host", host: "evil.test", auth: []string{"Bearer " + result.Credential}, forward: host, status: http.StatusForbidden},
		{name: "hostile origin", host: host, origins: []string{"https://evil.test"}, auth: []string{"Bearer " + result.Credential}, status: http.StatusForbidden},
		{name: "origin wrong port", host: host, origins: []string{"http://127.0.0.1:3001"}, auth: []string{"Bearer " + result.Credential}, status: http.StatusForbidden},
		{name: "duplicate origin", host: host, origins: []string{"http://127.0.0.1:3000", "http://127.0.0.1:3000"}, auth: []string{"Bearer " + result.Credential}, status: http.StatusForbidden},
		{name: "missing auth", host: host, status: http.StatusUnauthorized},
		{name: "basic auth", host: host, auth: []string{"Basic " + result.Credential}, status: http.StatusUnauthorized},
		{name: "wrong auth", host: host, auth: []string{"Bearer wrong"}, status: http.StatusUnauthorized},
		{name: "suffixed auth", host: host, auth: []string{"Bearer " + result.Credential + "x"}, status: http.StatusUnauthorized},
		{name: "duplicate auth", host: host, auth: []string{"Bearer " + result.Credential, "Bearer " + result.Credential}, status: http.StatusUnauthorized},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, result.Endpoint, strings.NewReader(initBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Host = tc.host
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json, text/event-stream")
			for _, value := range tc.origins {
				req.Header.Add("Origin", value)
			}
			for _, value := range tc.auth {
				req.Header.Add("Authorization", value)
			}
			if tc.forward != "" {
				req.Header.Set("X-Forwarded-Host", tc.forward)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode != tc.status {
				t.Fatalf("status=%d body=%s", resp.StatusCode, body)
			}
			if resp.Header.Get("Access-Control-Allow-Origin") != "" || strings.Contains(string(body), result.Credential) {
				t.Fatalf("unsafe response headers=%v body=%q", resp.Header, body)
			}
		})
	}
	old := result.Credential
	rotated, err := manager.Rotate(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := rawInitialize(t, result.Endpoint, old, "2025-11-25"); got != http.StatusUnauthorized {
		t.Fatalf("stale credential status=%d", got)
	}
	if got := rawInitialize(t, result.Endpoint, rotated.Credential, "2025-11-25"); got != http.StatusOK {
		t.Fatalf("rotated credential status=%d", got)
	}
}

func TestOfficialSDKHTTPDiscoveryAndSupportedRevisions(t *testing.T) {
	manager, result := enableManager(t)
	client := &http.Client{Transport: authTransport{credential: result.Credential, base: http.DefaultTransport}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	observeClient := mcp.NewClient(&mcp.Implementation{Name: "http-test", Version: "test"}, &mcp.ClientOptions{Capabilities: &mcp.ClientCapabilities{}})
	session, err := observeClient.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: result.Endpoint, HTTPClient: client, DisableStandaloneSSE: true, MaxRetries: -1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	if session.InitializeResult().Capabilities.Tools != nil {
		t.Fatal("HTTP advertised tools")
	}
	resources, err := session.ListResources(ctx, nil)
	if err != nil || len(resources.Resources) != 5 {
		t.Fatalf("resources=%+v err=%v", resources, err)
	}
	if _, err := session.ReadResource(ctx, &mcp.ReadResourceParams{URI: resources.Resources[0].URI}); err != nil {
		t.Fatal(err)
	}
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := manager.newObserveServer().Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()
	stdioClient := mcp.NewClient(&mcp.Implementation{Name: "stdio-parity", Version: "test"}, &mcp.ClientOptions{Capabilities: &mcp.ClientCapabilities{}})
	stdioSession, err := stdioClient.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer stdioSession.Close()
	stdioResources, err := stdioSession.ListResources(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	httpDiscovery, _ := json.Marshal(resources)
	stdioDiscovery, _ := json.Marshal(stdioResources)
	if !bytes.Equal(httpDiscovery, stdioDiscovery) {
		t.Fatalf("transport discovery differs: http=%s stdio=%s", httpDiscovery, stdioDiscovery)
	}
	httpHealth, err := session.ReadResource(ctx, &mcp.ReadResourceParams{URI: resources.Resources[0].URI})
	if err != nil {
		t.Fatal(err)
	}
	stdioHealth, err := stdioSession.ReadResource(ctx, &mcp.ReadResourceParams{URI: resources.Resources[0].URI})
	if err != nil {
		t.Fatal(err)
	}
	if normalizedResource(t, httpHealth.Contents[0].Text) != normalizedResource(t, stdioHealth.Contents[0].Text) {
		t.Fatalf("transport resource contracts differ: http=%s stdio=%s", httpHealth.Contents[0].Text, stdioHealth.Contents[0].Text)
	}
	for _, revision := range []string{"2025-11-25", "2026-07-28"} {
		if got := rawInitialize(t, result.Endpoint, result.Credential, revision); got != http.StatusOK {
			t.Fatalf("revision %s status=%d", revision, got)
		}
	}
}

func normalizedResource(t *testing.T, text string) string {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal([]byte(text), &value); err != nil {
		t.Fatal(err)
	}
	delete(value, "generated_at")
	normalized, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(normalized)
}

func TestRequestBodyBoundAndConcurrentLifecycle(t *testing.T) {
	manager, result := enableManager(t)
	req, err := http.NewRequest(http.MethodPost, result.Endpoint, bytes.NewReader(make([]byte, maxRequestBodyBytes+1)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+result.Credential)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized status=%d", resp.StatusCode)
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = manager.Status()
			_, _ = manager.Rotate(context.Background())
		}()
	}
	wg.Wait()
}

func TestHTTPResourceCancellationReachesObserveRead(t *testing.T) {
	reader := &cancelReader{started: make(chan struct{}), canceled: make(chan struct{})}
	manager := New(reader, "test", nil)
	result, err := manager.Enable(context.Background(), server.MCPHTTPEnableRequest{Port: freePort(t)})
	if err != nil {
		t.Fatal(err)
	}
	defer manager.Shutdown(context.Background())
	client := &http.Client{Transport: authTransport{credential: result.Credential, base: http.DefaultTransport}}
	observeClient := mcp.NewClient(&mcp.Implementation{Name: "cancel-test", Version: "test"}, nil)
	session, err := observeClient.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: result.Endpoint, HTTPClient: client, DisableStandaloneSSE: true, MaxRetries: -1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	readCtx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := session.ReadResource(readCtx, &mcp.ReadResourceParams{URI: "goschedule://daemon/health"})
		done <- err
	}()
	select {
	case <-reader.started:
	case <-time.After(time.Second):
		t.Fatal("Observe read did not start")
	}
	cancel()
	select {
	case <-reader.canceled:
	case <-time.After(time.Second):
		t.Fatal("request cancellation did not reach Observe read")
	}
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("canceled read succeeded")
		}
	case <-time.After(time.Second):
		t.Fatal("canceled read did not return")
	}
}

type authTransport struct {
	credential string
	base       http.RoundTripper
}

func (t authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	clone.Header.Set("Authorization", "Bearer "+t.credential)
	return t.base.RoundTrip(clone)
}

func rawInitialize(t *testing.T, endpoint, credential, revision string) int {
	t.Helper()
	body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"` + revision + `","capabilities":{},"clientInfo":{"name":"test","version":"1"}}}`
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+credential)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode
}

func stringPort(port int) string {
	return strconv.Itoa(port)
}
