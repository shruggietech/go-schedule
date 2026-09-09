// Package connection owns the desktop application's transport-neutral daemon lifecycle.
package connection

import (
	"context"
	"time"
)

// State is the complete user-visible connection state vocabulary.
type State string

const (
	StateConnecting   State = "connecting"
	StateConnected    State = "connected"
	StateDegraded     State = "degraded"
	StateRecovering   State = "recovering"
	StateUnavailable  State = "unavailable"
	StateAccessDenied State = "access_denied"
	StateIncompatible State = "incompatible"
	StateTimedOut     State = "timed_out"
)

// Target is the stable, non-sensitive identity shown by feature screens.
type Target struct {
	ID           string   `json:"id"`
	ProfileID    string   `json:"profileId,omitempty"`
	Kind         string   `json:"kind"`
	DisplayName  string   `json:"displayName"`
	Endpoint     string   `json:"endpoint,omitempty"`
	Fingerprint  string   `json:"fingerprint,omitempty"`
	Platform     string   `json:"platform"`
	Architecture string   `json:"architecture,omitempty"`
	Version      string   `json:"version,omitempty"`
	Capabilities []string `json:"capabilities"`
	Permissions  []string `json:"permissions"`
}

// Snapshot is an immutable generation-stamped view of the active connection.
type Snapshot struct {
	Generation       uint64 `json:"generation"`
	Revision         uint64 `json:"revision"`
	State            State  `json:"state"`
	Target           Target `json:"target"`
	Message          string `json:"message"`
	Action           string `json:"action,omitempty"`
	LastSuccessfulAt string `json:"lastSuccessfulAt,omitempty"`
}

// Event is the only payload published to the frontend event channel.
type Event struct {
	ID         string    `json:"id"`
	Kind       string    `json:"kind"`
	Message    string    `json:"message"`
	Generation uint64    `json:"generation"`
	OccurredAt string    `json:"occurredAt"`
	EntityID   string    `json:"entityId,omitempty"`
	Snapshot   *Snapshot `json:"snapshot,omitempty"`
}

// Health is the transport-independent result of daemon negotiation.
type Health struct {
	ID           string
	DisplayName  string
	Platform     string
	Architecture string
	Version      string
	Capabilities []string
	Permissions  []string
}

// DomainEvent is the bounded event information accepted from a backend.
type DomainEvent struct {
	Kind     string
	EntityID string
}

// Backend supplies daemon operations without exposing its transport.
type Backend interface {
	Health(context.Context) (Health, error)
	StreamEvents(context.Context, func(DomainEvent)) error
}

// Scheduler provides cancelable retry waits and a deterministic test seam.
type Scheduler interface {
	After(context.Context, time.Duration) <-chan struct{}
}

// Observer receives safe bridge events.
type Observer interface {
	Publish(Event)
}

// Failure carries a safe state and guidance; its Cause is never serialized.
type Failure struct {
	State   State
	Message string
	Action  string
	Cause   error
}

// LocalTarget returns the stable desktop identity for the bundled IPC daemon.
func LocalTarget() Target { return localTarget() }

// RemoteTarget returns the safe preflight identity for one persisted remote profile.
func RemoteTarget(profileID, daemonID, label, endpoint, fingerprint, platform, architecture, version string) Target {
	return Target{ID: daemonID, ProfileID: profileID, Kind: "remote", DisplayName: label, Endpoint: endpoint, Fingerprint: fingerprint, Platform: platform, Architecture: architecture, Version: version}
}

func (f *Failure) Error() string { return f.Message }
func (f *Failure) Unwrap() error { return f.Cause }
