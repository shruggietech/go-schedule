package mcphttp

import (
	"context"
	"errors"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/mcpoperate"
	"github.com/shruggietech/go-schedule/internal/mcpsession"
)

func TestCanceledEnableDiscardsUncommittedPreparation(t *testing.T) {
	manager := New(fakeReader{}, "test", nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	revoked := 0
	manager.sessionCreate = func(context.Context, string, domain.Capability) (server.MCPSessionCredentialResponse, error) {
		cancel()
		return server.MCPSessionCredentialResponse{Session: mcpsession.Session{ID: "uncommitted-session", ActorID: "owned-actor"}, Credential: "owned-private-token"}, nil
	}
	manager.sessionRevoke = func(cleanup context.Context, id string) error {
		if cleanup.Err() != nil || id != "uncommitted-session" {
			t.Errorf("uncommitted cleanup lost its independent context or identity: %v %s", cleanup.Err(), id)
		}
		revoked++
		return nil
	}
	manager.operateFor = func(string) *mcpoperate.Executor { return mcpoperate.New(&fakeOperator{}) }
	t.Cleanup(func() {
		if err := manager.Shutdown(context.Background()); err != nil {
			t.Error(err)
		}
	})
	port := freePort(t)
	_, err := manager.Enable(ctx, server.MCPHTTPEnableRequest{Port: port, Permission: domain.CapabilityOperate})
	if !errors.Is(err, context.Canceled) || manager.Status().Enabled || revoked != 1 {
		t.Errorf("canceled preparation was published or retained: error=%v enabled=%t revocations=%d", err, manager.Status().Enabled, revoked)
	}
	listener, err := net.Listen("tcp4", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)))
	if err != nil {
		t.Errorf("uncommitted listener was retained: %v", err)
	} else if err := listener.Close(); err != nil {
		t.Error(err)
	}
}

func TestDisableCompletesCommittedCleanupAfterCallerCancellation(t *testing.T) {
	manager := New(fakeReader{}, "test", nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	revoked := 0
	manager.sessionCreate = func(context.Context, string, domain.Capability) (server.MCPSessionCredentialResponse, error) {
		return server.MCPSessionCredentialResponse{Session: mcpsession.Session{ID: "committed-session", ActorID: "owned-actor"}, Credential: "owned-private-token"}, nil
	}
	manager.sessionRevoke = func(cleanup context.Context, id string) error {
		cancel()
		if cleanup.Err() != nil || id != "committed-session" {
			t.Errorf("committed cleanup inherited caller cancellation: %v %s", cleanup.Err(), id)
		}
		revoked++
		return nil
	}
	manager.operateFor = func(string) *mcpoperate.Executor { return mcpoperate.New(&fakeOperator{}) }
	if _, err := manager.Enable(context.Background(), server.MCPHTTPEnableRequest{Port: freePort(t), Permission: domain.CapabilityOperate}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := manager.Shutdown(context.Background()); err != nil {
			t.Error(err)
		}
	})
	status, err := manager.Disable(ctx)
	if err != nil || status.Enabled || revoked != 1 || ctx.Err() == nil {
		t.Fatalf("committed revocation was undone or lost cleanup: enabled=%t revocations=%d error=%v", status.Enabled, revoked, err)
	}
}

func TestCanceledLifecycleDoesNotChangeListenerOrCredential(t *testing.T) {
	for _, operation := range []string{"enable", "rotate", "disable"} {
		t.Run(operation, func(t *testing.T) {
			manager := New(fakeReader{}, "test", nil)
			if operation != "enable" {
				var err error
				_, err = manager.Enable(context.Background(), server.MCPHTTPEnableRequest{Port: freePort(t)})
				if err != nil {
					t.Fatal(err)
				}
			}
			t.Cleanup(func() {
				if err := manager.Shutdown(context.Background()); err != nil {
					t.Error(err)
				}
			})
			before := manager.Status()
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			var err error
			switch operation {
			case "enable":
				_, err = manager.Enable(ctx, server.MCPHTTPEnableRequest{Port: freePort(t)})
			case "rotate":
				_, err = manager.Rotate(ctx)
			case "disable":
				_, err = manager.Disable(ctx)
			}
			after := manager.Status()
			if !errors.Is(err, context.Canceled) || before.Enabled != after.Enabled || before.Endpoint != after.Endpoint || before.CredentialFingerprint != after.CredentialFingerprint {
				t.Fatalf("canceled %s changed lifecycle: error=%v before=%+v after=%+v", operation, err, before, after)
			}
		})
	}
}

func TestQueuedRotationCancellationReturnsBeforeLifecycleRelease(t *testing.T) {
	manager := New(fakeReader{}, "test", nil)
	entered, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	unblock := func() { once.Do(func() { close(release) }) }
	manager.sessionCreate = func(context.Context, string, domain.Capability) (server.MCPSessionCredentialResponse, error) {
		close(entered)
		<-release
		return server.MCPSessionCredentialResponse{Session: mcpsession.Session{ID: "owned-session", ActorID: "owned-actor", Capability: domain.CapabilityOperate}, Credential: "owned-private-token"}, nil
	}
	manager.sessionRevoke = func(context.Context, string) error { return nil }
	manager.operateFor = func(string) *mcpoperate.Executor { return mcpoperate.New(&fakeOperator{}) }
	t.Cleanup(func() {
		unblock()
		if err := manager.Shutdown(context.Background()); err != nil {
			t.Error(err)
		}
	})
	type enableResult struct {
		status server.MCPHTTPCredentialResponse
		err    error
	}
	enabled := make(chan enableResult, 1)
	port := freePort(t)
	go func() {
		status, err := manager.Enable(context.Background(), server.MCPHTTPEnableRequest{Port: port, Permission: domain.CapabilityOperate})
		enabled <- enableResult{status, err}
	}()
	deadline, stop := context.WithTimeout(context.Background(), 30*time.Second)
	defer stop()
	select {
	case <-entered:
	case <-deadline.Done():
		t.Fatal("owned enable did not reach held session creation")
	}
	base, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx := &observedLifecycleContext{Context: base, selected: make(chan struct{})}
	rotated := make(chan error, 1)
	go func() {
		_, err := manager.Rotate(ctx)
		rotated <- err
	}()
	select {
	case <-ctx.selected:
	case <-deadline.Done():
		t.Error("queued operation did not observe caller cancellation")
	}
	cancel()
	var rotationErr error
	returned := false
	select {
	case rotationErr = <-rotated:
		returned = true
		if !errors.Is(rotationErr, context.Canceled) {
			t.Errorf("queued rotation did not preserve cancellation: %v", rotationErr)
		}
	case <-deadline.Done():
		t.Error("canceled rotation waited for the unrelated held lifecycle operation")
	}
	unblock()
	first := <-enabled
	if !returned {
		rotationErr = <-rotated
	}
	if first.err != nil {
		t.Fatal(first.err)
	}
	if manager.Status().CredentialFingerprint != first.status.CredentialFingerprint {
		t.Errorf("canceled queued rotation replaced the newly enabled credential: %v", rotationErr)
	}
}

// Done acknowledges that lifecycle admission reached its cancelable select.
// The held Enable keeps admission unavailable until this test releases it.
type observedLifecycleContext struct {
	context.Context
	selected chan struct{}
	once     sync.Once
}

func (c *observedLifecycleContext) Done() <-chan struct{} {
	c.once.Do(func() { close(c.selected) })
	return c.Context.Done()
}
