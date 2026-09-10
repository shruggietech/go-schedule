package connection

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type remoteDaemonStub struct {
	health     server.HealthResponse
	manifest   server.ManifestResponse
	capability domain.Capability
	err        error
	events     []client.RemoteEvent
}

type delayedRemoteDaemon struct{ remoteDaemonStub }

func waitStage(ctx context.Context) error {
	select {
	case <-time.After(700 * time.Millisecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (stub delayedRemoteDaemon) Health(ctx context.Context) (server.HealthResponse, error) {
	if err := waitStage(ctx); err != nil {
		return server.HealthResponse{}, err
	}
	return stub.remoteDaemonStub.Health(ctx)
}
func (stub delayedRemoteDaemon) VerifyIdentity(ctx context.Context) (server.ManifestResponse, error) {
	if err := waitStage(ctx); err != nil {
		return server.ManifestResponse{}, err
	}
	return stub.remoteDaemonStub.VerifyIdentity(ctx)
}
func (stub delayedRemoteDaemon) VerifyAccess(ctx context.Context) (domain.Capability, error) {
	if err := waitStage(ctx); err != nil {
		return "", err
	}
	return stub.remoteDaemonStub.VerifyAccess(ctx)
}

func (stub remoteDaemonStub) Health(context.Context) (server.HealthResponse, error) {
	return stub.health, stub.err
}
func (stub remoteDaemonStub) VerifyIdentity(context.Context) (server.ManifestResponse, error) {
	return stub.manifest, stub.err
}
func (stub remoteDaemonStub) VerifyAccess(context.Context) (domain.Capability, error) {
	capability := stub.capability
	if capability == "" {
		capability = domain.CapabilityObserve
	}
	return capability, stub.err
}
func (stub remoteDaemonStub) StreamRemoteEvents(_ context.Context, publish func(client.RemoteEvent)) error {
	for _, event := range stub.events {
		publish(event)
	}
	return errors.New("closed")
}

func TestRemoteBackendNegotiatesIdentityAuthorityAndEvents(t *testing.T) {
	stub := remoteDaemonStub{health: server.HealthResponse{Status: "ok", Version: "v1.4.0"}, manifest: server.ManifestResponse{InstallationID: "daemon-1", DisplayName: "Remote", ProductVersion: "v1.4.0", RemoteAPIVersions: []string{"v1"}, Capabilities: []string{"tasks"}, Platform: server.ManifestPlatform{OS: "linux", Architecture: "amd64"}}, capability: domain.CapabilityOperate, events: []client.RemoteEvent{{Kind: "task", ResourceID: "task-1", Verb: "updated"}}}
	backend := NewRemoteBackend(stub, "operate")
	health, err := backend.Health(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if health.ID != "daemon-1" || len(health.Permissions) != 2 || health.Permissions[1] != "operate" {
		t.Fatalf("health = %#v", health)
	}
	var event DomainEvent
	_ = backend.StreamEvents(context.Background(), func(value DomainEvent) { event = value })
	if event.Kind != "task.updated" || event.EntityID != "task-1" {
		t.Fatalf("event = %#v", event)
	}
}

func TestRemoteBackendRejectsAuthorityDowngrade(t *testing.T) {
	stub := remoteDaemonStub{health: server.HealthResponse{Status: "ok", Version: "v1.4.0"}, manifest: server.ManifestResponse{InstallationID: "daemon-1", DisplayName: "Remote", ProductVersion: "v1.4.0", RemoteAPIVersions: []string{"v1"}, Platform: server.ManifestPlatform{OS: "linux"}}, capability: domain.CapabilityObserve}
	_, err := NewRemoteBackend(stub, "manage").Health(context.Background())
	var failure *Failure
	if !errors.As(err, &failure) || failure.State != StateForbidden {
		t.Fatalf("failure = %#v", failure)
	}
}

func TestRemoteBackendGivesEachNegotiationStageItsOwnTimeout(t *testing.T) {
	stub := remoteDaemonStub{health: server.HealthResponse{Status: "ok", Version: "v1.4.0"}, manifest: server.ManifestResponse{InstallationID: "daemon-1", DisplayName: "Remote", ProductVersion: "v1.4.0", RemoteAPIVersions: []string{"v1"}, Platform: server.ManifestPlatform{OS: "linux"}}}
	backend := NewRemoteBackend(delayedRemoteDaemon{remoteDaemonStub: stub}, "observe")
	started := time.Now()
	if _, err := backend.Health(context.Background()); err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(started); elapsed < 2*time.Second || elapsed >= remoteAttemptTimeout {
		t.Fatalf("negotiation elapsed=%s", elapsed)
	}
	if got := attemptTimeoutFor(backend); got != remoteAttemptTimeout {
		t.Fatalf("attempt timeout=%s", got)
	}
}

func TestRemoteBackendClassifiesIdentityAndAuthorityFailures(t *testing.T) {
	for _, test := range []struct {
		code  string
		state State
	}{{"authentication_failed", StateUnauthorized}, {"credential_revoked", StateRevoked}, {server.CodeForbidden, StateForbidden}, {server.CodeConflict, StateIdentityChanged}} {
		backend := NewRemoteBackend(remoteDaemonStub{err: &client.StatusError{Code: test.code, Message: "secret-canary"}}, "observe")
		_, err := backend.Health(context.Background())
		var failure *Failure
		if !errors.As(err, &failure) || failure.State != test.state || failure.Message == "secret-canary" {
			t.Fatalf("code %s failure = %#v", test.code, failure)
		}
	}
}
