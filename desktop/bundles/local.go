package bundles

import (
	"context"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/bundle"
)

type Backend interface {
	ExportBundle(context.Context) (bundle.Document, error)
	ValidateBundle(context.Context, bundle.Document) (client.BundleValidationResponse, error)
	PreviewBundle(context.Context, bundle.Document) (bundle.Plan, error)
	CompareBundle(context.Context, bundle.Document) (bundle.Plan, error)
	ApplyBundle(context.Context, bundle.Plan) (server.BundleApplyResponse, error)
}

type LocalBackend struct{ daemon Backend }

func NewLocalBackend(daemon Backend) *LocalBackend { return &LocalBackend{daemon: daemon} }
func (b *LocalBackend) ExportBundle(ctx context.Context) (bundle.Document, error) {
	return b.daemon.ExportBundle(ctx)
}
func (b *LocalBackend) ValidateBundle(ctx context.Context, doc bundle.Document) (client.BundleValidationResponse, error) {
	return b.daemon.ValidateBundle(ctx, doc)
}
func (b *LocalBackend) PreviewBundle(ctx context.Context, doc bundle.Document) (bundle.Plan, error) {
	return b.daemon.PreviewBundle(ctx, doc)
}
func (b *LocalBackend) CompareBundle(ctx context.Context, doc bundle.Document) (bundle.Plan, error) {
	return b.daemon.CompareBundle(ctx, doc)
}
func (b *LocalBackend) ApplyBundle(ctx context.Context, plan bundle.Plan) (server.BundleApplyResponse, error) {
	return b.daemon.ApplyBundle(ctx, plan)
}
