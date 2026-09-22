package client

import (
	"context"
	"net/http"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/bundle"
)

// ExportBundle returns the daemon's deterministic, secret-free portable intent.
func (c *Client) ExportBundle(ctx context.Context) (bundle.Document, error) {
	var out bundle.Document
	err := c.do(ctx, http.MethodGet, "/v1/bundles/export", nil, &out)
	return out, err
}

type BundleValidationResponse struct {
	Bundle bundle.Document `json:"bundle"`
	Digest string          `json:"digest"`
	Issues []bundle.Issue  `json:"issues"`
	Valid  bool            `json:"valid"`
}

func (c *Client) ValidateBundle(ctx context.Context, document bundle.Document) (BundleValidationResponse, error) {
	var out BundleValidationResponse
	err := c.do(ctx, http.MethodPost, "/v1/bundles/validate", server.BundleRequest{Bundle: document}, &out)
	return out, err
}

// PreviewBundle creates a single-use, target-bound plan. ApplyBundle accepts only this plan while it remains fresh.
func (c *Client) PreviewBundle(ctx context.Context, document bundle.Document) (bundle.Plan, error) {
	var out bundle.Plan
	err := c.do(ctx, http.MethodPost, "/v1/bundles/preview", server.BundleRequest{Bundle: document}, &out)
	return out, err
}

func (c *Client) CompareBundle(ctx context.Context, document bundle.Document) (bundle.Plan, error) {
	var out bundle.Plan
	err := c.do(ctx, http.MethodPost, "/v1/bundles/compare", server.BundleRequest{Bundle: document}, &out)
	return out, err
}

func (c *Client) ApplyBundle(ctx context.Context, plan bundle.Plan) (server.BundleApplyResponse, error) {
	var out server.BundleApplyResponse
	req := server.BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint}
	err := c.do(ctx, http.MethodPost, "/v1/bundles/apply", req, &out)
	return out, err
}
