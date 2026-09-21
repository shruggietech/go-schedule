package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/mcpsession"
)

type expectedDaemonKey struct{}

// ExpectDaemon binds one mutation request to the exact installation identity
// supplied by its caller. The server validates this value inside the audited
// mutation boundary.
func ExpectDaemon(ctx context.Context, daemonID string) context.Context {
	return context.WithValue(ctx, expectedDaemonKey{}, daemonID)
}

func expectedDaemon(ctx context.Context) string {
	value, _ := ctx.Value(expectedDaemonKey{}).(string)
	return value
}

func (c *Client) CreateMCPSession(ctx context.Context, clientName string, capability domain.Capability) (server.MCPSessionCredentialResponse, error) {
	var out server.MCPSessionCredentialResponse
	err := c.do(ctx, http.MethodPost, "/v1/mcp/sessions", server.MCPSessionCreateRequest{ClientName: clientName, Capability: capability}, &out)
	return out, err
}

func (c *Client) ListMCPSessions(ctx context.Context) ([]mcpsession.Session, error) {
	var out server.MCPSessionListResponse
	err := c.do(ctx, http.MethodGet, "/v1/mcp/sessions", nil, &out)
	return out.Sessions, err
}

func (c *Client) RevokeMCPSession(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/mcp/sessions/"+url.PathEscape(id), nil, nil)
}

// WithMCPSession returns an immutable local target that authenticates as the
// runtime MCP actor. It deliberately cannot add a session header to a remote
// target.
func (c *Client) WithMCPSession(secret string) *Client {
	target := c.target()
	clone := &Client{
		http:             target.http,
		endpoint:         target.endpoint,
		baseURL:          target.baseURL,
		pathPrefix:       target.pathPrefix,
		bearer:           target.bearer,
		expectedDaemonID: target.expectedDaemonID,
		remote:           target.remote,
		identityTimeout:  target.identityTimeout,
	}
	clone.verified.Store(target.verified.Load())
	if !target.remote {
		clone.mcpSession = secret
	}
	return clone
}
