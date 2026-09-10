package connection

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

const remoteStageTimeout = 2 * time.Second
const remoteAttemptTimeout = 7 * time.Second

type remoteDaemon interface {
	Health(context.Context) (server.HealthResponse, error)
	VerifyIdentity(context.Context) (server.ManifestResponse, error)
	VerifyAccess(context.Context) (domain.Capability, error)
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

func (*RemoteBackend) AutoRetry() bool               { return true }
func (*RemoteBackend) AttemptTimeout() time.Duration { return remoteAttemptTimeout }

// NewRemoteBackend creates a backend for one immutable remote selection.
func NewRemoteBackend(daemon remoteDaemon, capability string) *RemoteBackend {
	return &RemoteBackend{daemon: daemon, capability: capability}
}

func (backend *RemoteBackend) Health(ctx context.Context) (Health, error) {
	healthCtx, cancelHealth := context.WithTimeout(ctx, remoteStageTimeout)
	health, err := backend.daemon.Health(healthCtx)
	cancelHealth()
	if err != nil {
		return Health{}, remoteFailure(err)
	}
	if health.Status != "ok" || !compatibleVersion(health.Version) {
		return Health{}, &Failure{State: StateIncompatible, Message: "The selected remote scheduler version is incompatible.", Action: "Update the remote scheduler or select another connection."}
	}
	identityCtx, cancelIdentity := context.WithTimeout(ctx, remoteStageTimeout)
	manifest, err := backend.daemon.VerifyIdentity(identityCtx)
	cancelIdentity()
	if err != nil {
		return Health{}, remoteFailure(err)
	}
	if manifest.ProductVersion != health.Version || len(manifest.RemoteAPIVersions) == 0 {
		return Health{}, &Failure{State: StateIncompatible, Message: "The selected daemon does not advertise a compatible remote API.", Action: "Update the remote scheduler."}
	}
	accessCtx, cancelAccess := context.WithTimeout(ctx, remoteStageTimeout)
	capability, err := backend.daemon.VerifyAccess(accessCtx)
	cancelAccess()
	if err != nil {
		return Health{}, remoteFailure(err)
	}
	expected := domain.Capability(backend.capability)
	if !capability.Valid() || !capability.Allows(expected) {
		return Health{}, &Failure{State: StateForbidden, Message: "The selected remote credential no longer has its expected authority.", Action: "Ask an administrator to restore its grant or repair this connection."}
	}
	return Health{ID: manifest.InstallationID, DisplayName: manifest.DisplayName, Platform: manifest.Platform.OS, Architecture: manifest.Platform.Architecture, Version: manifest.ProductVersion, Capabilities: append([]string(nil), manifest.Capabilities...), Permissions: permissionsFor(string(capability))}, nil
}

func (backend *RemoteBackend) StreamEvents(ctx context.Context, publish func(DomainEvent)) error {
	err := backend.daemon.StreamRemoteEvents(ctx, func(event client.RemoteEvent) {
		kind := string(event.Kind)
		if event.Verb != "" {
			kind += "." + event.Verb
		} else {
			kind += ".changed"
		}
		publish(DomainEvent{Kind: kind, EntityID: event.ResourceID})
	})
	if err != nil && !errors.Is(err, context.Canceled) {
		return remoteFailure(err)
	}
	return err
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
		case "credential_revoked":
			return &Failure{State: StateRevoked, Message: "The selected remote credential was revoked.", Action: "Repair the connection with a new pairing phrase.", Cause: err}
		case "unauthorized", "authentication_failed":
			return &Failure{State: StateUnauthorized, Message: "The selected remote credential was rejected.", Action: "Repair the connection or verify its credential.", Cause: err}
		case server.CodeForbidden:
			return &Failure{State: StateForbidden, Message: "The selected remote credential lacks authority for this operation.", Action: "Ask an administrator to update its grant.", Cause: err}
		case server.CodeConflict:
			return &Failure{State: StateIdentityChanged, Message: "The selected remote daemon identity changed.", Action: "Do not continue until the target is verified.", Cause: err}
		}
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return &Failure{State: StateTimedOut, Message: "The selected remote scheduler did not respond in time.", Action: "Check the network and try again.", Cause: err}
	}
	var connectionErr *client.ConnectionError
	if errors.As(err, &connectionErr) {
		if connectionErr.Kind == client.ConnectionTrustFailure {
			return &Failure{State: StateTrustChanged, Message: "The selected remote certificate is no longer trusted.", Action: "Inspect the certificate change, then repair this connection explicitly.", Cause: err}
		}
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
