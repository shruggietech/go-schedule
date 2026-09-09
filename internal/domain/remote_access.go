package domain

import "time"

type PairingState string

const (
	PairingActive    PairingState = "active"
	PairingConsumed  PairingState = "consumed"
	PairingCancelled PairingState = "cancelled"
	PairingExpired   PairingState = "expired"
	PairingExhausted PairingState = "exhausted"
)

func (s PairingState) Valid() bool {
	switch s {
	case PairingActive, PairingConsumed, PairingCancelled, PairingExpired, PairingExhausted:
		return true
	default:
		return false
	}
}

type PairingSession struct {
	ID                string       `json:"id"`
	DisplayName       string       `json:"display_name"`
	Kind              ActorKind    `json:"kind"`
	Capability        Capability   `json:"capability"`
	AttemptsRemaining int          `json:"attempts_remaining"`
	State             PairingState `json:"state"`
	CreatedAt         time.Time    `json:"created_at"`
	ExpiresAt         time.Time    `json:"expires_at"`
	CompletedAt       *time.Time   `json:"completed_at,omitempty"`
}

type PairingSecret struct {
	PairingSession
	Phrase   string `json:"phrase"`
	DaemonID string `json:"daemon_id"`
}

type CredentialState string

const (
	CredentialActive  CredentialState = "active"
	CredentialRevoked CredentialState = "revoked"
)

func (s CredentialState) Valid() bool { return s == CredentialActive || s == CredentialRevoked }

type ClientCredential struct {
	ID          string          `json:"id"`
	ActorID     string          `json:"actor_id"`
	Fingerprint string          `json:"fingerprint"`
	State       CredentialState `json:"state"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	LastUsedAt  *time.Time      `json:"last_used_at,omitempty"`
	ExpiresAt   *time.Time      `json:"expires_at,omitempty"`
	RevokedAt   *time.Time      `json:"revoked_at,omitempty"`
}

type IssuedCredential struct {
	ClientCredential
	Actor    Actor  `json:"actor"`
	Token    string `json:"token"`
	DaemonID string `json:"daemon_id"`
}
