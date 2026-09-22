package bundles

import (
	"context"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/bundle"
)

type fakeBackend struct {
	document bundle.Document
	plan     bundle.Plan
}

func (f fakeBackend) ExportBundle(context.Context) (bundle.Document, error) { return f.document, nil }
func (f fakeBackend) ValidateBundle(context.Context, bundle.Document) (client.BundleValidationResponse, error) {
	return client.BundleValidationResponse{Valid: true}, nil
}
func (f fakeBackend) PreviewBundle(context.Context, bundle.Document) (bundle.Plan, error) {
	return f.plan, nil
}
func (f fakeBackend) CompareBundle(context.Context, bundle.Document) (bundle.Plan, error) {
	return f.plan, nil
}
func (f fakeBackend) ApplyBundle(context.Context, bundle.Plan) (server.BundleApplyResponse, error) {
	return server.BundleApplyResponse{Plan: f.plan, Outcomes: []bundle.Item{{Kind: "group", Action: bundle.ActionApplied}}}, nil
}

func TestServiceExportsPreviewsAndAppliesReviewedPlan(t *testing.T) {
	backend := fakeBackend{document: bundle.Document{Schema: bundle.SchemaV1}, plan: bundle.Plan{ID: "plan", TargetDaemonID: "target", BundleDigest: "digest", TargetFingerprint: "fingerprint"}}
	service := NewService(backend)
	if result := service.Export(context.Background()); result.Outcome != "accepted" || result.Document == "" {
		t.Fatalf("export=%+v", result)
	}
	if result := service.Preview(context.Background(), `{"schema":"go-schedule.bundle/v1"}`); result.Outcome != "accepted" || result.Plan == nil || result.Plan.ID != "plan" {
		t.Fatalf("preview=%+v", result)
	}
	if result := service.Apply(context.Background(), backend.plan); result.Outcome != "accepted" || len(result.Items) != 1 {
		t.Fatalf("apply=%+v", result)
	}
}

func TestServiceRejectsInvalidDocumentBeforeBackendCall(t *testing.T) {
	service := NewService(fakeBackend{})
	if result := service.Validate(context.Background(), "not JSON"); result.Outcome != "rejected" {
		t.Fatalf("validate=%+v", result)
	}
}
