package client

import (
	"context"
	"net/http"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

// Manifest returns daemon-owned identity and compatibility facts.
func (c *Client) Manifest(ctx context.Context) (server.ManifestResponse, error) {
	var out server.ManifestResponse
	err := c.do(ctx, http.MethodGet, "/v1/manifest", nil, &out)
	return out, err
}

// RenameDaemon updates the daemon's operator-facing display name.
func (c *Client) RenameDaemon(ctx context.Context, displayName string) (server.ManifestResponse, error) {
	var out server.ManifestResponse
	body := struct {
		DisplayName string `json:"display_name"`
	}{DisplayName: displayName}
	err := c.do(ctx, http.MethodPatch, "/v1/manifest", body, &out)
	return out, err
}

// ResetDaemonIdentity replaces the installation identity after exact acknowledgement.
func (c *Client) ResetDaemonIdentity(ctx context.Context, confirmation string) (server.ManifestResponse, error) {
	var out server.ManifestResponse
	body := struct {
		ConfirmInstallationID string `json:"confirm_installation_id"`
	}{ConfirmInstallationID: confirmation}
	err := c.do(ctx, http.MethodPost, "/v1/manifest/reset", body, &out)
	return out, err
}
