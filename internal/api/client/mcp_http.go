package client

import (
	"context"
	"net/http"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

// MCPHTTPStatus returns non-secret runtime listener state.
func (c *Client) MCPHTTPStatus(ctx context.Context) (server.MCPHTTPStatusResponse, error) {
	var out server.MCPHTTPStatusResponse
	err := c.do(ctx, http.MethodGet, "/v1/mcp/http", nil, &out)
	return out, err
}

// EnableMCPHTTP starts the runtime listener and returns its credential once.
func (c *Client) EnableMCPHTTP(ctx context.Context, req server.MCPHTTPEnableRequest) (server.MCPHTTPCredentialResponse, error) {
	var out server.MCPHTTPCredentialResponse
	err := c.do(ctx, http.MethodPost, "/v1/mcp/http/enable", req, &out)
	return out, err
}

// RotateMCPHTTPCredential invalidates the current credential and returns its replacement once.
func (c *Client) RotateMCPHTTPCredential(ctx context.Context) (server.MCPHTTPCredentialResponse, error) {
	var out server.MCPHTTPCredentialResponse
	err := c.do(ctx, http.MethodPost, "/v1/mcp/http/rotate", nil, &out)
	return out, err
}

// DisableMCPHTTP revokes the credential and stops the runtime listener.
func (c *Client) DisableMCPHTTP(ctx context.Context) (server.MCPHTTPStatusResponse, error) {
	var out server.MCPHTTPStatusResponse
	err := c.do(ctx, http.MethodPost, "/v1/mcp/http/disable", nil, &out)
	return out, err
}
