package gui

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
)

const (
	initialReconnectDelay = 2 * time.Second
	maximumReconnectDelay = 30 * time.Second
)

type evidenceValue uint8

const (
	unknown evidenceValue = iota
	no
	yes
)

type accessDiagnosis struct {
	Service       string
	GroupExists   evidenceValue
	AccountMember evidenceValue
	TokenMember   evidenceValue
	Detail        string
}

type connectionIncident struct {
	Kind         client.ConnectionFailureKind
	Title        string
	Guidance     string
	Detail       string
	Retrying     bool
	Revision     uint64
	VisibleCount uint64
}

type connectionState struct {
	mu       sync.RWMutex
	active   bool
	incident connectionIncident
	diagnose func() accessDiagnosis
}

func newConnectionState(diagnose func() accessDiagnosis) *connectionState {
	return &connectionState{diagnose: diagnose}
}

func (s *connectionState) report(err *client.ConnectionError) connectionIncident {
	s.mu.Lock()
	if s.active && s.incident.Kind == err.Kind {
		s.incident.Detail = err.Error()
		s.incident.Retrying = false
		s.incident.Revision++
		incident := s.incident
		s.mu.Unlock()
		return incident
	}
	s.mu.Unlock()

	diagnosis := accessDiagnosis{}
	if err.Kind == client.ConnectionAccessDenied && s.diagnose != nil {
		diagnosis = s.diagnose()
	}
	title, guidance := connectionCopy(err.Kind, diagnosis)

	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.active {
		s.incident.VisibleCount = 1
	}
	s.active = true
	s.incident.Kind = err.Kind
	s.incident.Title = title
	s.incident.Guidance = guidance
	s.incident.Detail = err.Error()
	s.incident.Retrying = false
	s.incident.Revision++
	return s.incident
}

func (s *connectionState) setRetrying() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.active {
		return
	}
	s.incident.Retrying = true
	s.incident.Revision++
}

func (s *connectionState) finishRetry() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.active || !s.incident.Retrying {
		return
	}
	s.incident.Retrying = false
	s.incident.Revision++
}

func (s *connectionState) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.active = false
	s.incident = connectionIncident{}
}

func (s *connectionState) snapshot() (connectionIncident, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.incident, s.active
}

func connectionCopy(kind client.ConnectionFailureKind, diagnosis accessDiagnosis) (string, string) {
	switch kind {
	case client.ConnectionAccessDenied:
		switch {
		case diagnosis.GroupExists == yes && diagnosis.AccountMember == yes && diagnosis.TokenMember == no:
			return "Access denied", "Your account belongs to goschedadmin, but this login session does not include it. Sign out of Windows and back in, then launch go-schedule again."
		case diagnosis.AccountMember == no:
			return "Access denied", "Your account is not a member of the local goschedadmin group. Ask an administrator to verify the installation and group membership."
		default:
			return "Access denied", "The daemon rejected this Windows session. go-schedule could not verify whether goschedadmin membership or a stale login token caused it. Check the service and group diagnostics below."
		}
	case client.ConnectionUnavailable:
		return "Daemon unavailable", "Check that the goschedd service is installed and running, then try again."
	case client.ConnectionTimeout:
		return "Connection timed out", "The daemon did not respond before the request deadline. Check the service and try again."
	default:
		return "Connection interrupted", "The local daemon connection failed. Check the service and try again."
	}
}

func nextReconnectDelay(current time.Duration) time.Duration {
	if current <= 0 {
		return initialReconnectDelay
	}
	next := current * 2
	if next > maximumReconnectDelay {
		return maximumReconnectDelay
	}
	return next
}

func reconnectDelayAfterAttempt(current time.Duration, streamRecovered bool) time.Duration {
	if streamRecovered {
		current = 0
	}
	return nextReconnectDelay(current)
}

func waitForReconnect(ctx context.Context, delay time.Duration, retry <-chan struct{}) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-retry:
		return true
	case <-timer.C:
		return true
	}
}

func asConnectionError(err error) (*client.ConnectionError, bool) {
	var connectionErr *client.ConnectionError
	ok := errors.As(err, &connectionErr)
	return connectionErr, ok
}

func containsFold(value, fragment string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(fragment))
}
