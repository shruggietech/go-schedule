package main

import (
	"context"
	"github.com/shruggietech/go-schedule/v2/desktop/automation"
	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"net/http"
	"testing"
	"time"
)

type automationFacadeBackend struct {
	automation.Backend
	calls    int
	expected string
}

func (b *automationFacadeBackend) DeleteChain(ctx context.Context, id string) error {
	b.calls++
	c := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b.expected = r.Header.Get(server.ExpectedDaemonHeader)
		w.WriteHeader(http.StatusNoContent)
	}))
	return c.DeleteChain(ctx, id)
}
func (b *automationFacadeBackend) ListTaskDetails(context.Context, string, string) ([]server.TaskResponse, error) {
	return nil, nil
}
func (b *automationFacadeBackend) ListChains(context.Context) ([]domain.CompletionChain, error) {
	return nil, nil
}
func (b *automationFacadeBackend) ListTriggers(context.Context) ([]server.TriggerResponse, error) {
	return nil, nil
}
func (b *automationFacadeBackend) ListTriggerSets(context.Context) ([]server.TriggerSetResponse, error) {
	return nil, nil
}
func (b *automationFacadeBackend) ListFilesystemWatchers(context.Context) ([]server.FilesystemWatcherResponse, error) {
	return nil, nil
}
func (b *automationFacadeBackend) ListJoins(context.Context) ([]domain.ParallelJoin, error) {
	return nil, nil
}
func (b *automationFacadeBackend) ListJoinInstances(context.Context, domain.ParallelJoinInstanceFilter) ([]domain.ParallelJoinInstance, error) {
	return nil, nil
}
func TestAutomationFacadePinsCurrentIdentityAndManageAuthority(t *testing.T) {
	for _, tc := range []struct {
		name, expected string
		permissions    []string
		allowed        bool
	}{{"manager", "owned-daemon", []string{"manage"}, true}, {"other", "other", []string{"manage"}, false}, {"observer", "owned-daemon", []string{"read"}, false}, {"empty", "", []string{"manage"}, false}} {
		t.Run(tc.name, func(t *testing.T) {
			b := &automationFacadeBackend{}
			ready := make(chan struct{}, 1)
			app := newApp(cancellationAppBackend{permissions: tc.permissions, capabilities: []string{"parallel-joins"}}, cancellationAppEmitter{ready}, nil, appServices{automation: automation.NewService(b)})
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			app.startup(ctx)
			defer app.shutdown(context.Background())
			select {
			case <-ready:
			case <-time.After(time.Second):
				t.Fatal("no current target")
			}
			r := app.DeleteChain("chain", tc.expected)
			if (r.Outcome == "accepted") != tc.allowed {
				t.Fatalf("result=%+v", r)
			}
			if tc.allowed {
				if b.calls != 1 || b.expected != "owned-daemon" || r.DaemonID != "owned-daemon" {
					t.Fatalf("calls=%d header=%s result=%+v", b.calls, b.expected, r)
				}
			} else if b.calls != 0 {
				t.Fatal("denied mutation reached backend")
			}
		})
	}
}
