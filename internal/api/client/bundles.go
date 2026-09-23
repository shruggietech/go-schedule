package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/bundle"
)

// ExportBundle returns the daemon's deterministic, secret-free portable intent.
func (c *Client) ExportBundle(ctx context.Context) (bundle.Document, error) {
	target, profile, _, err := c.bundleTarget(ctx)
	if err != nil {
		return bundle.Document{}, err
	}
	if issues := bundle.TargetCompatibility(bundle.Document{}, profile); len(issues) > 0 {
		return bundle.Document{}, bundleCompatibilityError(issues)
	}
	var out bundle.Document
	err = target.do(ctx, http.MethodGet, "/v1/bundles/export", nil, &out)
	return out, err
}

type BundleValidationResponse struct {
	Bundle bundle.Document `json:"bundle"`
	Digest string          `json:"digest"`
	Issues []bundle.Issue  `json:"issues"`
	Valid  bool            `json:"valid"`
}

func (c *Client) ValidateBundle(ctx context.Context, document bundle.Document) (BundleValidationResponse, error) {
	target, profile, _, err := c.bundleTarget(ctx)
	if err != nil {
		return BundleValidationResponse{}, err
	}
	if issues := bundle.TargetCompatibility(document, profile); len(issues) > 0 {
		canonical, _, digest, validationIssues := bundle.Canonicalize(document)
		return BundleValidationResponse{Bundle: canonical, Digest: digest, Issues: append(validationIssues, issues...), Valid: false}, nil
	}
	var out BundleValidationResponse
	err = target.do(ctx, http.MethodPost, "/v1/bundles/validate", server.BundleRequest{Bundle: document}, &out)
	return out, err
}

// PreviewBundle creates a single-use, target-bound plan. ApplyBundle accepts only this plan while it remains fresh.
func (c *Client) PreviewBundle(ctx context.Context, document bundle.Document) (bundle.Plan, error) {
	return c.PreviewBundleWithPaths(ctx, document, nil)
}

// PreviewBundleWithPaths supplies target-local watcher paths outside the canonical bundle.
func (c *Client) PreviewBundleWithPaths(ctx context.Context, document bundle.Document, paths map[string]string) (bundle.Plan, error) {
	target, profile, daemonID, err := c.bundleTarget(ctx)
	if err != nil {
		return bundle.Plan{}, err
	}
	if issues := bundle.TargetCompatibility(document, profile); len(issues) > 0 {
		return bundle.Plan{}, bundleCompatibilityError(issues)
	}
	var out bundle.Plan
	err = target.do(ctx, http.MethodPost, "/v1/bundles/preview", server.BundleRequest{Bundle: document, WatcherPaths: paths}, &out)
	if err == nil && out.TargetDaemonID != daemonID {
		return bundle.Plan{}, fmt.Errorf("selected daemon changed during preview; preview again")
	}
	return out, err
}

func (c *Client) CompareBundle(ctx context.Context, document bundle.Document) (bundle.Plan, error) {
	target, profile, daemonID, err := c.bundleTarget(ctx)
	if err != nil {
		return bundle.Plan{}, err
	}
	if issues := bundle.TargetCompatibility(document, profile); len(issues) > 0 {
		return bundle.Plan{}, bundleCompatibilityError(issues)
	}
	var out bundle.Plan
	err = target.do(ctx, http.MethodPost, "/v1/bundles/compare", server.BundleRequest{Bundle: document}, &out)
	if err == nil && out.TargetDaemonID != daemonID {
		return bundle.Plan{}, fmt.Errorf("selected daemon changed during comparison; compare again")
	}
	return out, err
}

func (c *Client) ApplyBundle(ctx context.Context, plan bundle.Plan) (server.BundleApplyResponse, error) {
	target, profile, daemonID, err := c.bundleTarget(ctx)
	if err != nil {
		return server.BundleApplyResponse{}, err
	}
	if issues := bundle.PlanTargetCompatibility(plan, profile); len(issues) > 0 {
		return server.BundleApplyResponse{}, bundleCompatibilityError(issues)
	}
	if plan.TargetDaemonID != daemonID {
		return server.BundleApplyResponse{}, fmt.Errorf("selected daemon changed after preview; select %s and preview again", plan.TargetDaemonID)
	}
	var out server.BundleApplyResponse
	req := server.BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint}
	err = target.do(ctx, http.MethodPost, "/v1/bundles/apply", req, &out)
	return out, err
}

func (c *Client) bundleTarget(ctx context.Context) (*Client, bundle.TargetProfile, string, error) {
	target := c.target()
	manifest, err := target.Manifest(ctx)
	if err != nil {
		return nil, bundle.TargetProfile{}, "", fmt.Errorf("discover selected daemon bundle capabilities: %w", err)
	}
	if manifest.InstallationID == "" {
		return nil, bundle.TargetProfile{}, "", fmt.Errorf("selected daemon manifest has no installation identity; reconnect and try again")
	}
	return target, bundle.TargetProfile{Platform: manifest.Platform.OS, Capabilities: manifest.Capabilities}, manifest.InstallationID, nil
}

func bundleCompatibilityError(issues []bundle.Issue) error {
	messages := make([]string, 0, len(issues))
	for _, issue := range issues {
		messages = append(messages, issue.Message)
	}
	return fmt.Errorf("bundle target is incompatible: %s", strings.Join(messages, "; "))
}
