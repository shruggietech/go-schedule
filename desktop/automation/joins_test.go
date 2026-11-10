package automation

import (
	"context"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

type s144UncertainBackend struct{ Backend }

func (s144UncertainBackend) CreateChain(context.Context, server.ChainCreateRequest) (domain.CompletionChain, error) {
	return domain.CompletionChain{}, &client.MutationUncertainError{Operation: "POST /v1/chains", Cause: context.DeadlineExceeded}
}
func (s144UncertainBackend) DeleteChain(context.Context, string) error {
	return &client.MutationUncertainError{Operation: "DELETE /v1/chains/fixture", Cause: context.DeadlineExceeded}
}
func (s144UncertainBackend) CreateTrigger(context.Context, server.TriggerCreateRequest) (server.TriggerSecretResponse, error) {
	return server.TriggerSecretResponse{}, &client.MutationUncertainError{Operation: "POST /v1/triggers", Cause: context.DeadlineExceeded}
}

func TestS144BaselineAutomationUncertainty(t *testing.T) {
	svc := NewService(s144UncertainBackend{})
	t.Run("chain_create", func(t *testing.T) {
		r := svc.SaveChain(context.Background(), ChainDraft{SourceTaskID: "entry", TargetTaskID: "target", OnOutcome: "success", IsNew: true})
		if r.Outcome != "uncertain" {
			t.Fatalf("accepted-or-unknown mutation must retain uncertainty; got outcome=%q message=%q", r.Outcome, r.Message)
		}
	})
	t.Run("chain_delete", func(t *testing.T) {
		r := svc.DeleteChain(context.Background(), "fixture")
		if r.Outcome != "uncertain" {
			t.Fatalf("accepted-or-unknown mutation must retain uncertainty; got outcome=%q message=%q", r.Outcome, r.Message)
		}
	})
	t.Run("trigger_create", func(t *testing.T) {
		r := svc.SaveTrigger(context.Background(), TriggerDraft{Name: "fixture", TargetTaskID: "target", IsNew: true})
		if r.Outcome != "uncertain" {
			t.Fatalf("accepted-or-unknown secret mutation must retain uncertainty; got outcome=%q message=%q", r.Outcome, r.Message)
		}
	})
}
