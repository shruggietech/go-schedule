// Package search owns cross-daemon discovery and target-safe actions.
package search

import (
	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/desktop/systems"
	"github.com/shruggietech/go-schedule/internal/domain"
)

// Request is one explicit cross-daemon query.
type Request struct {
	Query string              `json:"query"`
	Kinds []domain.SearchKind `json:"kinds,omitempty"`
	Limit int                 `json:"limit,omitempty"`
}

// Match binds one daemon result to its immutable source registration.
type Match struct {
	RegistrationKey  string                `json:"registrationKey"`
	ExpectedDaemonID string                `json:"expectedDaemonId"`
	SourceLabel      string                `json:"sourceLabel"`
	SourceShortID    string                `json:"sourceShortId"`
	Result           domain.SearchMatch    `json:"result"`
	AvailableActions []domain.SearchAction `json:"availableActions"`
	DisabledReason   string                `json:"disabledReason,omitempty"`
}

// Observation is one registration's current search outcome.
type Observation struct {
	Registration systems.Registration `json:"registration"`
	State        connection.State     `json:"state"`
	ObservedAt   string               `json:"observedAt,omitempty"`
	Truncated    bool                 `json:"truncated"`
	Matches      []Match              `json:"matches"`
	Failure      *systems.Failure     `json:"failure,omitempty"`
}

// Snapshot is one generation-scoped progressive search result.
type Snapshot struct {
	Generation   uint64        `json:"generation"`
	Query        string        `json:"query"`
	StartedAt    string        `json:"startedAt"`
	CompletedAt  string        `json:"completedAt,omitempty"`
	Complete     bool          `json:"complete"`
	Observations []Observation `json:"observations"`
}

// Selection is one explicitly confirmed object and target identity.
type Selection struct {
	RegistrationKey  string            `json:"registrationKey"`
	ExpectedDaemonID string            `json:"expectedDaemonId"`
	Kind             domain.SearchKind `json:"kind"`
	ObjectID         string            `json:"objectId"`
	DisplayName      string            `json:"displayName"`
}

// ActionIntent applies one compatible action to explicit selections.
type ActionIntent struct {
	Action     domain.SearchAction `json:"action"`
	Selections []Selection         `json:"selections"`
}

// ActionOutcome preserves independent evidence for one object.
type ActionOutcome struct {
	RegistrationKey  string `json:"registrationKey"`
	ExpectedDaemonID string `json:"expectedDaemonId"`
	CurrentDaemonID  string `json:"currentDaemonId,omitempty"`
	Kind             string `json:"kind"`
	ObjectID         string `json:"objectId"`
	DisplayName      string `json:"displayName"`
	Outcome          string `json:"outcome"`
	Message          string `json:"message"`
}

// ActionBatchResult summarizes non-transactional per-object outcomes.
type ActionBatchResult struct {
	Action   domain.SearchAction `json:"action"`
	Outcome  string              `json:"outcome"`
	Message  string              `json:"message"`
	Outcomes []ActionOutcome     `json:"outcomes"`
}
