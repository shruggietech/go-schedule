// Package systems owns the read-only multi-daemon operational overview.
package systems

import (
	"time"

	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/internal/domain"
)

// Registration identifies one local or saved desktop target without secret material.
type Registration struct {
	Key           string `json:"key"`
	ProfileID     string `json:"profileId,omitempty"`
	Kind          string `json:"kind"`
	Label         string `json:"label"`
	Endpoint      string `json:"endpoint,omitempty"`
	DaemonID      string `json:"daemonId,omitempty"`
	ShortDaemonID string `json:"shortDaemonId,omitempty"`
	Platform      string `json:"platform,omitempty"`
	Architecture  string `json:"architecture,omitempty"`
	Version       string `json:"version,omitempty"`
}

// Failure is safe target-specific guidance for an unsuccessful observation.
type Failure struct {
	State   connection.State `json:"state"`
	Message string           `json:"message"`
	Action  string           `json:"action"`
}

// Observation combines one registration with its current or session-stale summary.
type Observation struct {
	Registration Registration          `json:"registration"`
	State        connection.State      `json:"state"`
	ObservedAt   string                `json:"observedAt,omitempty"`
	Stale        bool                  `json:"stale"`
	Summary      *domain.SystemSummary `json:"summary,omitempty"`
	Failure      *Failure              `json:"failure,omitempty"`
}

// Snapshot is one generation-scoped refresh result.
type Snapshot struct {
	Generation   uint64        `json:"generation"`
	StartedAt    string        `json:"startedAt"`
	CompletedAt  string        `json:"completedAt"`
	Observations []Observation `json:"observations"`
}

type cachedSummary struct {
	summary domain.SystemSummary
	at      time.Time
}
