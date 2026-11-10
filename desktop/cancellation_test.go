package main

import (
	"context"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/desktop/connection"
	"github.com/shruggietech/go-schedule/v2/desktop/operations"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

type cancellationAppBackend struct {
	appBackend
	permissions  []string
	capabilities []string
}

func (b cancellationAppBackend) Health(context.Context) (connection.Health, error) {
	return connection.Health{ID: "owned-daemon", DisplayName: "Owned", Version: "v2.0.3", Permissions: b.permissions, Capabilities: b.capabilities}, nil
}

type cancellationAppEmitter struct{ ready chan struct{} }

func (e cancellationAppEmitter) Emit(_ context.Context, _ string, value any) {
	if event, ok := value.(connection.Event); ok && event.Snapshot != nil && event.Snapshot.State == connection.StateConnected {
		select {
		case e.ready <- struct{}{}:
		default:
		}
	}
}

type cancellationFacadeBackend struct {
	facadeOperationsBackend
	calls int
	id    string
}

func (b *cancellationFacadeBackend) CancelRun(_ context.Context, id string) (domain.RunCancellation, error) {
	b.calls++
	b.id = id
	return domain.RunCancellation{RunID: id, Status: "requested"}, nil
}
func TestCancellationFacadeBindsCurrentIdentityAuthorityAndCapability(t *testing.T) {
	for _, tc := range []struct {
		name                      string
		permissions, capabilities []string
		expected                  string
		accepted                  bool
	}{{"current", []string{"manage"}, []string{"run-cancellation"}, "owned-daemon", true}, {"wrong target", []string{"manage"}, []string{"run-cancellation"}, "other-daemon", false}, {"observer", []string{"read"}, []string{"run-cancellation"}, "owned-daemon", false}, {"unsupported", []string{"manage"}, nil, "owned-daemon", false}} {
		t.Run(tc.name, func(t *testing.T) {
			backend := &cancellationFacadeBackend{}
			ready := make(chan struct{}, 1)
			app := newApp(cancellationAppBackend{permissions: tc.permissions, capabilities: tc.capabilities}, cancellationAppEmitter{ready}, nil, appServices{operations: operations.NewService(backend)})
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			app.startup(ctx)
			defer app.shutdown(context.Background())
			select {
			case <-ready:
			case <-time.After(time.Second):
				t.Fatal("connection did not become current")
			}
			r := app.CancelRun("run-1", tc.expected)
			if (r.Outcome == "accepted") != tc.accepted {
				t.Fatalf("result=%+v", r)
			}
			if tc.accepted {
				if backend.calls != 1 || backend.id != "run-1" || r.Cancellation == nil || r.Cancellation.Status != "requested" {
					t.Fatalf("r=%+v calls=%d id=%s", r, backend.calls, backend.id)
				}
			} else if backend.calls != 0 {
				t.Fatalf("denied request reached backend: %d", backend.calls)
			}
		})
	}
}
