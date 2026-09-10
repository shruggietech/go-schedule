package connection

import (
	"context"
	"errors"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
)

type remoteDaemonStub struct {
	health   server.HealthResponse
	manifest server.ManifestResponse
	err      error
	events   []client.RemoteEvent
}

func (stub remoteDaemonStub) Health(context.Context) (server.HealthResponse, error) {
	return stub.health, stub.err
}
func (stub remoteDaemonStub) VerifyIdentity(context.Context) (server.ManifestResponse, error) {
	return stub.manifest, stub.err
}
func (stub remoteDaemonStub) StreamRemoteEvents(_ context.Context, publish func(client.RemoteEvent)) error {
	for _, event := range stub.events {
		publish(event)
	}
	return errors.New("closed")
}

func TestRemoteBackendNegotiatesIdentityAuthorityAndEvents(t *testing.T) {
	stub := remoteDaemonStub{health: server.HealthResponse{Status: "ok", Version: "v1.4.0"}, manifest: server.ManifestResponse{InstallationID: "daemon-1", DisplayName: "Remote", ProductVersion: "v1.4.0", RemoteAPIVersions: []string{"v1"}, Capabilities: []string{"tasks"}, Platform: server.ManifestPlatform{OS: "linux", Architecture: "amd64"}}, events: []client.RemoteEvent{{Kind: "task", ResourceID: "task-1", Verb: "updated"}}}
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

func TestRemoteBackendClassifiesIdentityAndAuthorityFailures(t *testing.T) {
	for _, test := range []struct {
		code  string
		state State
	}{{"authentication_failed", StateAccessDenied}, {server.CodeForbidden, StateAccessDenied}, {server.CodeConflict, StateIncompatible}} {
		backend := NewRemoteBackend(remoteDaemonStub{err: &client.StatusError{Code: test.code, Message: "secret-canary"}}, "observe")
		_, err := backend.Health(context.Background())
		var failure *Failure
		if !errors.As(err, &failure) || failure.State != test.state || failure.Message == "secret-canary" {
			t.Fatalf("code %s failure = %#v", test.code, failure)
		}
	}
}
