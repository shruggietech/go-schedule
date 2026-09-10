package connection

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
)

type remoteDaemon interface {
	Health(context.Context) (server.HealthResponse, error)
	VerifyIdentity(context.Context) (server.ManifestResponse, error)
	StreamRemoteEvents(context.Context, func(client.RemoteEvent)) error
}

// RemoteBackend adapts one identity-pinned HTTPS client to the desktop connection contract.
type RemoteBackend struct {
	daemon     remoteDaemon
	capability string
}

type unavailableRemoteBackend struct {
	message string
	action  string
}

// NewUnavailableRemoteBackend preserves a selected remote target while failing closed until it is repaired.
func NewUnavailableRemoteBackend(message, action string) Backend {
	return &unavailableRemoteBackend{message: message, action: action}
}

func (*unavailableRemoteBackend) AutoRetry() bool { return false }
func (backend *unavailableRemoteBackend) Health(context.Context) (Health, error) {
	return Health{}, &Failure{State: StateUnavailable, Message: backend.message, Action: backend.action}
}
func (backend *unavailableRemoteBackend) StreamEvents(context.Context, func(DomainEvent)) error {
	return &Failure{State: StateUnavailable, Message: backend.message, Action: backend.action}
}

func (*RemoteBackend) AutoRetry() bool { return false }

// NewRemoteBackend creates a backend for one immutable remote selection.
func NewRemoteBackend(daemon remoteDaemon, capability string) *RemoteBackend {
	return &RemoteBackend{daemon: daemon, capability: capability}
}

func (backend *RemoteBackend) Health(ctx context.Context) (Health, error) {
	health, err := backend.daemon.Health(ctx)
	if err != nil {
		return Health{}, remoteFailure(err)
	}
	if health.Status != "ok" || !compatibleVersion(health.Version) {
		return Health{}, &Failure{State: StateIncompatible, Message: "The selected remote scheduler version is incompatible.", Action: "Update the remote scheduler or select another connection."}
	}
	manifest, err := backend.daemon.VerifyIdentity(ctx)
	if err != nil {
		return Health{}, remoteFailure(err)
	}
	if manifest.ProductVersion != health.Version || len(manifest.RemoteAPIVersions) == 0 {
		return Health{}, &Failure{State: StateIncompatible, Message: "The selected daemon does not advertise a compatible remote API.", Action: "Update the remote scheduler."}
	}
	return Health{ID: manifest.InstallationID, DisplayName: manifest.DisplayName, Platform: manifest.Platform.OS, Architecture: manifest.Platform.Architecture, Version: manifest.ProductVersion, Capabilities: append([]string(nil), manifest.Capabilities...), Permissions: permissionsFor(backend.capability)}, nil
}

func (backend *RemoteBackend) StreamEvents(ctx context.Context, publish func(DomainEvent)) error {
	return backend.daemon.StreamRemoteEvents(ctx, func(event client.RemoteEvent) {
		kind := string(event.Kind)
		if event.Verb != "" {
			kind += "." + event.Verb
		} else {
			kind += ".changed"
		}
		publish(DomainEvent{Kind: kind, EntityID: event.ResourceID})
	})
}

func permissionsFor(capability string) []string {
	switch capability {
	case "manage":
		return []string{"read", "operate", "manage"}
	case "operate":
		return []string{"read", "operate"}
	case "enroll":
		return []string{"read", "operate", "manage", "enroll"}
	default:
		return []string{"read"}
	}
}

func remoteFailure(err error) error {
	var status *client.StatusError
	if errors.As(err, &status) {
		switch status.Code {
		case "authentication_failed", server.CodeForbidden:
			return &Failure{State: StateAccessDenied, Message: "The selected remote credential was rejected or lacks authority.", Action: "Repair the connection or ask an administrator to update its grant.", Cause: err}
		case server.CodeConflict:
			return &Failure{State: StateIncompatible, Message: "The selected remote daemon identity changed.", Action: "Do not continue until the target is verified.", Cause: err}
		}
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return &Failure{State: StateTimedOut, Message: "The selected remote scheduler did not respond in time.", Action: "Check the network and try again.", Cause: err}
	}
	var connectionErr *client.ConnectionError
	if errors.As(err, &connectionErr) {
		return &Failure{State: StateUnavailable, Message: "The selected remote scheduler is unavailable.", Action: "Check the endpoint and network, then try again.", Cause: err}
	}
	return &Failure{State: StateUnavailable, Message: "The selected remote scheduler could not be reached.", Action: "Check the connection profile and try again.", Cause: fmt.Errorf("remote connection: %w", err)}
}

// ShortID returns a bounded identity suffix suitable for disambiguating display names.
func ShortID(value string) string {
	value = strings.TrimSpace(value)
	if len(value) <= 8 {
		return value
	}
	return value[:8]
}
