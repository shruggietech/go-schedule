package domain

import (
	"errors"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// Capability is a closed, monotonic authority level.
type Capability string

const (
	CapabilityObserve Capability = "observe"
	CapabilityOperate Capability = "operate"
	CapabilityManage  Capability = "manage"
	CapabilityEnroll  Capability = "enroll"
)

var ErrInvalidActor = errors.New("invalid actor")

func (c Capability) Valid() bool { return c.rank() != 0 }

func (c Capability) Allows(required Capability) bool {
	return c.rank() != 0 && required.rank() != 0 && c.rank() >= required.rank()
}

func (c Capability) rank() int {
	switch c {
	case CapabilityObserve:
		return 1
	case CapabilityOperate:
		return 2
	case CapabilityManage:
		return 3
	case CapabilityEnroll:
		return 4
	default:
		return 0
	}
}

type ActorKind string

const (
	ActorKindLocalOS ActorKind = "local_os"
	ActorKindDesktop ActorKind = "desktop"
	ActorKindCLI     ActorKind = "cli"
	ActorKindJSON    ActorKind = "json"
	ActorKindMCP     ActorKind = "mcp"
)

func (k ActorKind) Valid() bool {
	switch k {
	case ActorKindLocalOS, ActorKindDesktop, ActorKindCLI, ActorKindJSON, ActorKindMCP:
		return true
	default:
		return false
	}
}

type ActorState string

const (
	ActorStateActive  ActorState = "active"
	ActorStateExpired ActorState = "expired"
	ActorStateRevoked ActorState = "revoked"
)

func (s ActorState) Valid() bool {
	return s == ActorStateActive || s == ActorStateExpired || s == ActorStateRevoked
}

type Actor struct {
	ID          string     `json:"id"`
	Kind        ActorKind  `json:"kind"`
	DisplayName string     `json:"display_name"`
	Capability  Capability `json:"capability"`
	State       ActorState `json:"state"`
	Builtin     bool       `json:"builtin"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

func NormalizeActorDisplayName(value string) (string, error) {
	if !utf8.ValidString(value) {
		return "", ErrInvalidActor
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return "", ErrInvalidActor
		}
	}
	value = strings.TrimSpace(value)
	if utf8.RuneCountInString(value) < 1 || utf8.RuneCountInString(value) > 80 {
		return "", ErrInvalidActor
	}
	return value, nil
}

func (a Actor) ActiveAt(now time.Time) bool {
	return a.State == ActorStateActive && (a.ExpiresAt == nil || a.ExpiresAt.After(now))
}

type AuditResult string

const (
	AuditResultUncertain AuditResult = "uncertain"
	AuditResultSucceeded AuditResult = "succeeded"
	AuditResultFailed    AuditResult = "failed"
	AuditResultDenied    AuditResult = "denied"
)

func (r AuditResult) Valid() bool {
	return r == AuditResultUncertain || r == AuditResultSucceeded || r == AuditResultFailed || r == AuditResultDenied
}

type AuditEvent struct {
	ID            string      `json:"id"`
	ActorID       string      `json:"actor_id,omitempty"`
	DaemonID      string      `json:"daemon_id"`
	Operation     string      `json:"operation"`
	TargetKind    string      `json:"target_kind"`
	TargetID      string      `json:"target_id,omitempty"`
	Result        AuditResult `json:"result"`
	CorrelationID string      `json:"correlation_id"`
	OccurredAt    time.Time   `json:"occurred_at"`
	CompletedAt   *time.Time  `json:"completed_at,omitempty"`
}

type AuditQuery struct {
	ActorID   string
	Operation string
	Result    AuditResult
	Since     *time.Time
	Until     *time.Time
	Limit     int
}
