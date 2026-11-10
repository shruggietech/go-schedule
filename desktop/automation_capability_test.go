package main

import (
	"context"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/desktop/automation"
	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

type capabilityAutomationBackend struct {
	automationFacadeBackend
	joinReads, instanceReads int
	unsupported              bool
}

func (b *capabilityAutomationBackend) ListTaskDetails(context.Context, string, string) ([]server.TaskResponse, error) {
	return []server.TaskResponse{{Task: domain.Task{ID: "task", Name: "Preserved task"}}}, nil
}
func (b *capabilityAutomationBackend) ListChains(context.Context) ([]domain.CompletionChain, error) {
	return []domain.CompletionChain{{ID: "chain"}}, nil
}
func (b *capabilityAutomationBackend) ListTriggers(context.Context) ([]server.TriggerResponse, error) {
	return []server.TriggerResponse{{ID: "trigger"}}, nil
}
func (b *capabilityAutomationBackend) ListTriggerSets(context.Context) ([]server.TriggerSetResponse, error) {
	return []server.TriggerSetResponse{{ID: "set"}}, nil
}
func (b *capabilityAutomationBackend) ListFilesystemWatchers(context.Context) ([]server.FilesystemWatcherResponse, error) {
	return []server.FilesystemWatcherResponse{{ID: "watcher"}}, nil
}
func (b *capabilityAutomationBackend) ListJoins(context.Context) ([]domain.ParallelJoin, error) {
	b.joinReads++
	if b.unsupported {
		return nil, &client.StatusError{Code: server.CodeNotFound, Message: "unsupported route"}
	}
	return []domain.ParallelJoin{{ID: "join"}}, nil
}
func (b *capabilityAutomationBackend) ListJoinInstances(context.Context, domain.ParallelJoinInstanceFilter) ([]domain.ParallelJoinInstance, error) {
	b.instanceReads++
	return []domain.ParallelJoinInstance{{ID: "instance"}}, nil
}
func TestAutomationWorkspaceUsesAdvertisedJoinCapability(t *testing.T) {
	for _, tc := range []struct {
		name                    string
		advertised, unsupported bool
	}{{"compatible scheduler without joins", false, true}, {"join-capable scheduler", true, false}, {"advertised joins failure remains visible", true, true}} {
		t.Run(tc.name, func(t *testing.T) {
			b := &capabilityAutomationBackend{unsupported: tc.unsupported}
			ready := make(chan struct{}, 1)
			var capabilities []string
			if tc.advertised {
				capabilities = []string{"parallel-joins"}
			}
			app := newApp(cancellationAppBackend{permissions: []string{"manage"}, capabilities: capabilities}, cancellationAppEmitter{ready}, nil, appServices{automation: automation.NewService(b)})
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			app.startup(ctx)
			defer app.shutdown(context.Background())
			select {
			case <-ready:
			case <-time.After(time.Second):
				t.Fatal("no current scheduler")
			}
			r := app.AutomationWorkspace("owned-daemon")
			if tc.advertised && tc.unsupported {
				if r.Outcome == "accepted" || r.Workspace != nil || b.joinReads != 1 {
					t.Fatalf("advertised failure was hidden: result=%+v joinReads=%d", r, b.joinReads)
				}
				return
			}
			if r.Outcome != "accepted" || r.Workspace == nil {
				t.Fatalf("workspace=%+v", r)
			}
			w := r.Workspace
			if len(w.Tasks) != 1 || len(w.Chains) != 1 || len(w.Triggers) != 1 || len(w.TriggerSets) != 1 || len(w.Watchers) != 1 {
				t.Fatalf("existing families were lost: %+v", w)
			}
			if !tc.advertised {
				if b.joinReads != 0 || b.instanceReads != 0 {
					t.Fatalf("unsupported reads: joins=%d instances=%d", b.joinReads, b.instanceReads)
				}
				changed := app.DeleteChain("chain", "owned-daemon")
				if changed.Outcome != "accepted" || changed.Workspace == nil || b.joinReads != 0 || b.instanceReads != 0 {
					t.Fatalf("mutation refresh read unsupported joins: result=%+v joins=%d instances=%d", changed, b.joinReads, b.instanceReads)
				}
			} else if len(w.Joins) != 1 || len(w.JoinInstances) != 1 || b.joinReads != 1 || b.instanceReads != 1 {
				t.Fatalf("full advertised workspace missing: %+v joinReads=%d instances=%d", w, b.joinReads, b.instanceReads)
			}
		})
	}
}
