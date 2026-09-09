// Package client is the shared Go client for the daemon's local API. Both the
// CLI and the GUI use it, so they operate on identical state. It dials the IPC
// transport (Unix socket / named pipe) rather than a network port.
package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/ipc"
)

// Client talks to the daemon over the IPC endpoint.
type Client struct {
	http             *http.Client
	endpoint         string
	baseURL          string
	pathPrefix       string
	bearer           string
	expectedDaemonID string
	remote           bool
	verified         atomic.Bool
	selected         atomic.Pointer[Client]
}

// New returns a client bound to the given IPC endpoint (socket path / pipe name).
func New(endpoint string) *Client {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return ipc.DialContext(ctx, endpoint)
		},
	}
	client := &Client{http: &http.Client{Transport: transport}, endpoint: endpoint, baseURL: "http://ipc"}
	client.verified.Store(true)
	return client
}

// NewSwitchable returns a stable client facade whose selected immutable target can change safely.
func NewSwitchable(initial *Client) *Client {
	client := &Client{}
	client.selected.Store(initial)
	return client
}

// Use changes future requests to next without altering requests already constructed for the prior target.
func (c *Client) Use(next *Client) {
	if next == nil {
		return
	}
	c.selected.Store(next)
}

func (c *Client) target() *Client {
	if selected := c.selected.Load(); selected != nil {
		return selected
	}
	return c
}

// baseURL uses a fixed dummy host; the transport ignores it and dials the IPC
// endpoint instead.
// NewRemote returns a client pinned to one HTTPS origin, trust bundle, bearer credential, and daemon identity.
func NewRemote(endpoint, certificatePEM, bearer, expectedDaemonID string) (*Client, error) {
	canonical, err := clientprofile.NormalizeEndpoint(endpoint)
	if err != nil || strings.TrimSpace(bearer) == "" || strings.TrimSpace(expectedDaemonID) == "" {
		return nil, errors.New("remote target configuration is incomplete")
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(certificatePEM)) {
		return nil, errors.New("remote target certificate is invalid")
	}
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS13, RootCAs: pool}}
	return &Client{http: &http.Client{Transport: transport, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}, endpoint: canonical, baseURL: canonical, pathPrefix: "/api", bearer: bearer, expectedDaemonID: strings.TrimSpace(expectedDaemonID), remote: true}, nil
}

// Remote reports whether this client uses authenticated HTTPS.
func (c *Client) Remote() bool { return c.target().remote }

// Endpoint returns the safe selected endpoint.
func (c *Client) Endpoint() string { return c.target().endpoint }

// VerifyIdentity pins a remote client to its expected installation manifest before feature operations.
func (c *Client) VerifyIdentity(ctx context.Context) (server.ManifestResponse, error) {
	target := c.target()
	if target.remote {
		target.verified.Store(false)
	}
	manifest, err := target.Manifest(ctx)
	if err != nil {
		return server.ManifestResponse{}, err
	}
	if target.remote && manifest.InstallationID != target.expectedDaemonID {
		return server.ManifestResponse{}, &StatusError{Code: server.CodeConflict, Message: "remote daemon identity does not match the selected target"}
	}
	target.verified.Store(true)
	return manifest, nil
}

// StatusError is returned for non-2xx API responses, carrying the API error
// envelope's code and field so callers (e.g. the CLI) can map them to exit codes.
type StatusError struct {
	Code    string
	Field   string
	Message string
}

func (e *StatusError) Error() string { return e.Message }

func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	target := c.target()
	if target.remote && !target.verified.Load() && path != "/v1/health" && path != "/v1/manifest" {
		return nil, &StatusError{Code: server.CodeConflict, Message: "remote daemon identity has not been verified"}
	}
	var encoded []byte
	if body != nil {
		var err error
		encoded, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	request, err := http.NewRequestWithContext(ctx, method, target.baseURL+target.pathPrefix+path, strings.NewReader(string(encoded)))
	if err != nil {
		return nil, err
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if target.remote && target.bearer != "" && path != "/v1/health" && path != "/v1/manifest" {
		request.Header.Set("Authorization", "Bearer "+target.bearer)
	}
	return request, nil
}

// Health calls GET /v1/health.
func (c *Client) Health(ctx context.Context) (server.HealthResponse, error) {
	var out server.HealthResponse
	if err := c.get(ctx, "/v1/health", &out); err != nil {
		return server.HealthResponse{}, err
	}
	return out, nil
}

// RuntimeInfo returns the effective storage paths reported by the connected
// daemon instance.
func (c *Client) RuntimeInfo(ctx context.Context) (server.RuntimeInfoResponse, error) {
	var out server.RuntimeInfoResponse
	if err := c.get(ctx, "/v1/runtime-info", &out); err != nil {
		return server.RuntimeInfoResponse{}, err
	}
	return out, nil
}

// get performs a GET and decodes a JSON body, surfacing the API error envelope.
func (c *Client) get(ctx context.Context, path string, out any) error {
	target := c.target()
	req, err := target.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	resp, err := target.http.Do(req)
	if err != nil {
		return NewConnectionError("GET "+path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		var apiErr server.APIError
		if decErr := json.NewDecoder(resp.Body).Decode(&apiErr); decErr == nil && apiErr.Error.Message != "" {
			return &StatusError{Code: apiErr.Error.Code, Field: apiErr.Error.Field, Message: apiErr.Error.Message}
		}
		return fmt.Errorf("api: %s: unexpected status %d", path, resp.StatusCode)
	}
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("api: decode %s: %w", path, err)
		}
	}
	return nil
}
