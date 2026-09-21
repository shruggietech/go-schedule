package bundles

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/bundle"
)

const callTimeout = 5 * time.Second

type Service struct{ backend Backend }

func NewService(backend Backend) *Service { return &Service{backend: backend} }

func (s *Service) Export(ctx context.Context) Result {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	doc, err := s.backend.ExportBundle(c)
	if err != nil {
		return failure("export_bundle", err)
	}
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return failure("export_bundle", err)
	}
	return Result{Action: "export_bundle", Outcome: "accepted", Message: "Portable intent exported without execution inputs or secrets.", Document: string(data)}
}

func (s *Service) Validate(ctx context.Context, raw string) Result {
	doc, err := parse(raw)
	if err != nil {
		return rejected("validate_bundle", "Enter valid bundle JSON.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	result, err := s.backend.ValidateBundle(c, doc)
	if err != nil {
		return failure("validate_bundle", err)
	}
	if !result.Valid {
		return Result{Action: "validate_bundle", Outcome: "rejected", Message: "Bundle needs attention before it can be previewed.", Issues: result.Issues}
	}
	return Result{Action: "validate_bundle", Outcome: "accepted", Message: "Bundle is valid for review.", Issues: result.Issues}
}

func (s *Service) Preview(ctx context.Context, raw string) Result {
	doc, err := parse(raw)
	if err != nil {
		return rejected("preview_bundle", "Enter valid bundle JSON.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	plan, err := s.backend.PreviewBundle(c, doc)
	if err != nil {
		return failure("preview_bundle", err)
	}
	return Result{Action: "preview_bundle", Outcome: "accepted", Message: "Review the target-bound plan before applying it.", Plan: &plan}
}

func (s *Service) Compare(ctx context.Context, raw string) Result {
	doc, err := parse(raw)
	if err != nil {
		return rejected("compare_bundle", "Enter valid bundle JSON.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	plan, err := s.backend.CompareBundle(c, doc)
	if err != nil {
		return failure("compare_bundle", err)
	}
	return Result{Action: "compare_bundle", Outcome: "accepted", Message: "Target drift was compared without mutation.", Plan: &plan, CompareOnly: true}
}

func (s *Service) Apply(ctx context.Context, plan bundle.Plan) Result {
	if strings.TrimSpace(plan.ID) == "" || strings.TrimSpace(plan.TargetDaemonID) == "" {
		return rejected("apply_bundle", "Preview a bundle before applying it.")
	}
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	result, err := s.backend.ApplyBundle(c, plan)
	if err != nil {
		return failure("apply_bundle", err)
	}
	return Result{Action: "apply_bundle", Outcome: "accepted", Message: "Bundle apply completed. Review every item outcome.", Plan: &result.Plan, Items: result.Outcomes}
}

func parse(raw string) (bundle.Document, error) {
	return bundle.Decode([]byte(raw))
}

func rejected(action, message string) Result {
	return Result{Action: action, Outcome: "rejected", Message: message}
}
func failure(action string, err error) Result {
	return Result{Action: action, Outcome: "unavailable", Message: err.Error()}
}
