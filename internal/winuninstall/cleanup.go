package winuninstall

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

const resultSchema = 1

// TargetKind identifies which closed application-owned storage class a target belongs to.
type TargetKind string

const (
	TargetMachine TargetKind = "machine"
	TargetProfile TargetKind = "profile"
)

// State is the durable terminal or in-progress state of one cleanup attempt.
type State string

const (
	StateRunning       State = "running"
	StateComplete      State = "complete"
	StateRefused       State = "refused"
	StatePartial       State = "partial"
	StateInternalError State = "internal-error"
)

// Outcome records what happened to one declared application-owned root.
type Outcome string

const (
	OutcomePending Outcome = "pending"
	OutcomeReady   Outcome = "ready"
	OutcomeAbsent  Outcome = "absent"
	OutcomeKept    Outcome = "preserved"
	OutcomeRemoved Outcome = "removed"
	OutcomeRefused Outcome = "refused"
	OutcomeFailed  Outcome = "failed"
)

// Target is one internally derived application-owned root eligible for cleanup.
type Target struct {
	Kind     TargetKind `json:"kind"`
	SID      string     `json:"sid,omitempty"`
	Path     string     `json:"path"`
	base     string
	relative string
}

// Entry combines a declared cleanup target with its recorded outcome.
type Entry struct {
	Target
	Outcome Outcome `json:"outcome"`
	Error   string  `json:"error,omitempty"`
}

// Result is the durable, machine-readable record of one cleanup attempt.
type Result struct {
	Schema        int       `json:"schema"`
	TransactionID string    `json:"transaction_id"`
	State         State     `json:"state"`
	StartedAt     time.Time `json:"started_at"`
	CompletedAt   time.Time `json:"completed_at,omitempty"`
	Remaining     int       `json:"remaining"`
	Entries       []Entry   `json:"entries,omitempty"`
	Error         string    `json:"error,omitempty"`
}

// Backend supplies trusted target discovery, safe removal, and durable evidence storage.
type Backend interface {
	Discover() ([]Target, error)
	Preflight(Target) (bool, error)
	Remove(Target) error
	WriteResult(Result) error
	ClearResult() error
}

// Run performs all-target preflight before deletion and records every state transition.
func Run(backend Backend) Result {
	result := Result{
		Schema:        resultSchema,
		TransactionID: transactionID(),
		State:         StateRunning,
		StartedAt:     time.Now().UTC(),
	}
	targets, err := backend.Discover()
	if err != nil {
		result.State = StateInternalError
		result.Error = err.Error()
		result.CompletedAt = time.Now().UTC()
		if writeErr := backend.WriteResult(result); writeErr != nil {
			result.Error = fmt.Sprintf("%s; record discovery failure: %v", result.Error, writeErr)
		}
		return result
	}
	for _, target := range targets {
		result.Entries = append(result.Entries, Entry{Target: target, Outcome: OutcomePending})
	}
	if err := backend.WriteResult(result); err != nil {
		result.State = StateInternalError
		result.Error = fmt.Sprintf("create cleanup result ledger: %v", err)
		result.CompletedAt = time.Now().UTC()
		if retryErr := backend.WriteResult(result); retryErr != nil {
			result.Error = fmt.Sprintf("%s; record terminal cleanup state: %v", result.Error, retryErr)
		}
		return result
	}

	refused := false
	for index, target := range targets {
		exists, err := backend.Preflight(target)
		if err != nil {
			result.Entries[index].Outcome = OutcomeRefused
			result.Entries[index].Error = err.Error()
			refused = true
		} else if !exists {
			result.Entries[index].Outcome = OutcomeAbsent
		} else {
			result.Entries[index].Outcome = OutcomeReady
		}
	}
	if refused {
		result.State = StateRefused
		for index, entry := range result.Entries {
			if entry.Outcome == OutcomeReady {
				result.Entries[index].Outcome = OutcomeKept
			}
			if entry.Outcome != OutcomeAbsent {
				result.Remaining++
			}
		}
		result.CompletedAt = time.Now().UTC()
		if err := backend.WriteResult(result); err != nil {
			result.State = StateInternalError
			result.Error = fmt.Sprintf("record refused cleanup: %v", err)
		}
		return result
	}

	for index, target := range targets {
		if result.Entries[index].Outcome == OutcomeAbsent {
			continue
		}
		if err := backend.Remove(target); err != nil {
			result.Entries[index].Outcome = OutcomeFailed
			result.Entries[index].Error = err.Error()
			result.Remaining++
		} else {
			result.Entries[index].Outcome = OutcomeRemoved
		}
		if err := backend.WriteResult(result); err != nil && result.Error == "" {
			result.Error = fmt.Sprintf("update cleanup result ledger: %v", err)
		}
	}
	result.CompletedAt = time.Now().UTC()
	if result.Remaining > 0 {
		result.State = StatePartial
	} else if result.Error != "" {
		result.State = StateInternalError
	} else {
		result.State = StateComplete
	}
	if err := backend.WriteResult(result); err != nil {
		result.State = StateInternalError
		result.Error = fmt.Sprintf("record cleanup completion: %v", err)
		return result
	}
	if result.State == StateComplete {
		if err := backend.ClearResult(); err != nil {
			result.State = StateInternalError
			result.Error = fmt.Sprintf("clear completed cleanup evidence: %v", err)
			if writeErr := backend.WriteResult(result); writeErr != nil {
				result.Error = fmt.Sprintf("%s; record evidence-clear failure: %v", result.Error, writeErr)
			}
		}
	}
	return result
}

func transactionID() string {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return fmt.Sprintf("time-%d", time.Now().UTC().UnixNano())
	}
	value[6] = value[6]&0x0f | 0x40
	value[8] = value[8]&0x3f | 0x80
	encoded := hex.EncodeToString(value[:])
	return encoded[0:8] + "-" + encoded[8:12] + "-" + encoded[12:16] + "-" + encoded[16:20] + "-" + encoded[20:]
}
