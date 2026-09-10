package connection

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"
)

const localAttemptTimeout = 2 * time.Second

var localRetryDelays = [...]time.Duration{250 * time.Millisecond, time.Second, 5 * time.Second}
var remoteRetryDelays = [...]time.Duration{500 * time.Millisecond, 2 * time.Second, 10 * time.Second, 30 * time.Second}

type timerScheduler struct{}

func (timerScheduler) After(ctx context.Context, delay time.Duration) <-chan struct{} {
	ready := make(chan struct{})
	var once sync.Once
	closeReady := func() { once.Do(func() { close(ready) }) }
	timer := time.AfterFunc(delay, closeReady)
	context.AfterFunc(ctx, func() { timer.Stop(); closeReady() })
	return ready
}

// Manager owns exactly one connection generation and event stream at a time.
type Manager struct {
	backend    Backend
	scheduler  Scheduler
	observer   Observer
	now        func() time.Time
	retryDelay func(Target, int) time.Duration

	mu              sync.RWMutex
	target          Target
	snapshot        Snapshot
	cancel          context.CancelFunc
	done            chan struct{}
	retry           chan struct{}
	retryQueued     bool
	retryProcessing bool
	started         bool
	stopped         bool
	sequence        atomic.Uint64
}

// NewManager creates an idle manager with an honest initial snapshot.
func NewManager(backend Backend, observer Observer) *Manager {
	manager := newManager(backend, observer, timerScheduler{}, time.Now)
	manager.retryDelay = jitteredRetryDelay
	return manager
}

func newManager(backend Backend, observer Observer, scheduler Scheduler, now func() time.Time) *Manager {
	target := localTarget()
	return &Manager{backend: backend, target: target, observer: observer, scheduler: scheduler, now: now, retryDelay: baseRetryDelay, retry: make(chan struct{}, 1), done: make(chan struct{}), snapshot: Snapshot{State: StateConnecting, Target: target, Message: connectingMessage(target), Recovery: RecoveryNone}}
}

func baseRetryDelay(target Target, attempt int) time.Duration {
	delays := localRetryDelays[:]
	if target.Kind == "remote" {
		delays = remoteRetryDelays[:]
	}
	return delays[min(max(attempt-1, 0), len(delays)-1)]
}

func jitteredRetryDelay(target Target, attempt int) time.Duration {
	base := baseRetryDelay(target, attempt)
	if target.Kind != "remote" {
		return base
	}
	window := base / 4
	if window <= 0 {
		return base
	}
	delay := base - window + time.Duration(rand.Int64N(int64(2*window)+1))
	return min(delay, remoteRetryDelays[len(remoteRetryDelays)-1])
}

func attemptTimeoutFor(backend Backend) time.Duration {
	if policy, ok := backend.(interface{ AttemptTimeout() time.Duration }); ok && policy.AttemptTimeout() > 0 {
		return policy.AttemptTimeout()
	}
	return localAttemptTimeout
}

// Switch cancels the active generation and selects an immutable backend for future work.
func (m *Manager) Switch(backend Backend, target Target) bool {
	if backend == nil || target.ID == "" || target.DisplayName == "" {
		return false
	}
	m.mu.Lock()
	if m.stopped {
		m.mu.Unlock()
		return false
	}
	m.backend = backend
	m.target = cloneTarget(target)
	started := m.started
	m.snapshot = Snapshot{Generation: m.snapshot.Generation, Revision: m.snapshot.Revision + 1, State: StateConnecting, Target: cloneTarget(target), Message: connectingMessage(target), LastSuccessfulAt: target.lastSuccessfulAt, Stale: target.Kind == "remote" && target.lastSuccessfulAt != "", Recovery: RecoveryNone}
	snapshot := cloneSnapshot(m.snapshot)
	if started && !m.retryQueued {
		m.retryQueued = true
		m.retry <- struct{}{}
	}
	m.mu.Unlock()
	m.emitSnapshot(snapshot)
	return true
}

// Start begins the single-owner loop. Repeated calls are harmless.
func (m *Manager) Start(parent context.Context) {
	m.mu.Lock()
	if m.started || m.stopped {
		m.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(parent)
	m.cancel, m.started = cancel, true
	m.mu.Unlock()
	go m.run(ctx)
}

// Snapshot returns a detached copy of current safe state.
func (m *Manager) Snapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return cloneSnapshot(m.snapshot)
}

// Retry requests one fresh generation and coalesces repeated requests.
func (m *Manager) Retry() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.stopped {
		return false
	}
	if !m.retryQueued {
		m.retryQueued = true
		m.retry <- struct{}{}
	}
	return true
}

// ReconcileMutation blocks further remote work immediately and starts a fresh
// generation whose authoritative reads can settle an uncertain mutation.
func (m *Manager) ReconcileMutation() bool {
	m.mu.Lock()
	if m.stopped || m.target.Kind != "remote" {
		m.mu.Unlock()
		return false
	}
	m.snapshot = Snapshot{Generation: m.snapshot.Generation, Revision: m.snapshot.Revision + 1, State: StateRecovering, Target: cloneTarget(m.snapshot.Target), Message: "Refreshing authoritative remote state before another change.", Action: "Wait for reconnection to complete.", LastSuccessfulAt: m.snapshot.LastSuccessfulAt, Stale: true, Recovery: RecoveryAutomatic}
	snapshot := cloneSnapshot(m.snapshot)
	if !m.retryQueued {
		m.retryQueued = true
		m.retry <- struct{}{}
	}
	m.mu.Unlock()
	m.emitSnapshot(snapshot)
	return true
}

func (m *Manager) beginRetry() {
	m.mu.Lock()
	m.retryProcessing = true
	m.mu.Unlock()
}

func (m *Manager) completeRetry() {
	m.mu.Lock()
	if m.retryProcessing {
		m.retryQueued = false
		m.retryProcessing = false
	}
	m.mu.Unlock()
}

// Stop cancels owned work and waits within the caller's deadline.
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	if m.stopped {
		done := m.done
		m.mu.Unlock()
		select {
		case <-done:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	m.stopped = true
	cancel, started := m.cancel, m.started
	if !started {
		close(m.done)
	}
	m.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	select {
	case <-m.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *Manager) run(ctx context.Context) {
	defer close(m.done)
	var generation uint64
	retryAttempt := 0

connectionLoop:
	for {
		generation++
		backend, target := m.selected()
		state := StateConnecting
		if generation > 1 {
			state = StateRecovering
		}
		m.publishSelectedSnapshot(generation, state, target, connectingMessage(target), "", retryAttempt)
		attemptCtx, cancelAttempt := context.WithTimeout(ctx, attemptTimeoutFor(backend))
		type healthResult struct {
			health Health
			err    error
		}
		result := make(chan healthResult, 1)
		go func() {
			health, err := backend.Health(attemptCtx)
			result <- healthResult{health: health, err: err}
		}()
		m.completeRetry()
		var health Health
		var err error
		select {
		case completed := <-result:
			health, err = completed.health, completed.err
			cancelAttempt()
		case <-ctx.Done():
			cancelAttempt()
			<-result
			return
		case <-m.retry:
			m.beginRetry()
			cancelAttempt()
			<-result
			retryAttempt = 0
			continue connectionLoop
		}
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			failure := failureOf(err)
			if terminalState(failure.State) || !autoRetry(backend) {
				m.publishManualFailure(generation, failure)
				if !m.waitManual(ctx) {
					return
				}
				retryAttempt = 0
				continue
			}
			retryAttempt++
			delay := m.retryDelay(target, retryAttempt)
			m.publishAutomaticFailure(generation, failure.State, failure.Message, failure.Action, retryAttempt, delay)
			proceed, manual := m.waitRetry(ctx, delay)
			if !proceed {
				return
			}
			if manual {
				retryAttempt = 0
			}
			continue
		}

		m.publishConnected(generation, target, health)
		streamCtx, cancelStream := context.WithCancel(ctx)
		streamDone := make(chan error, 1)
		streamActivity := make(chan struct{}, 1)
		var streamHadActivity atomic.Bool
		go func(active uint64) {
			streamDone <- backend.StreamEvents(streamCtx, func(event DomainEvent) {
				if m.publishDomain(active, event) {
					streamHadActivity.Store(true)
					select {
					case streamActivity <- struct{}{}:
					default:
					}
				}
			})
		}(generation)
		for {
			select {
			case <-ctx.Done():
				cancelStream()
				<-streamDone
				return
			case <-m.retry:
				m.beginRetry()
				cancelStream()
				<-streamDone
				retryAttempt = 0
				continue connectionLoop
			case <-streamActivity:
				retryAttempt = 0
			case streamErr := <-streamDone:
				cancelStream()
				if ctx.Err() != nil {
					return
				}
				failure := failureOf(streamErr)
				if terminalState(failure.State) || !autoRetry(backend) {
					m.publishManualFailure(generation, failure)
					if !m.waitManual(ctx) {
						return
					}
					retryAttempt = 0
					continue connectionLoop
				}
				if streamHadActivity.Load() {
					retryAttempt = 0
				}
				retryAttempt++
				delay := m.retryDelay(target, retryAttempt)
				m.publishAutomaticFailure(generation, StateDegraded, "Live updates are temporarily unavailable.", "Retry now.", retryAttempt, delay)
				proceed, manual := m.waitRetry(ctx, delay)
				if !proceed {
					return
				}
				if manual {
					retryAttempt = 0
				}
				continue connectionLoop
			}
		}
	}
}

func terminalState(state State) bool {
	switch state {
	case StateAccessDenied, StateUnauthorized, StateRevoked, StateForbidden, StateIncompatible, StateTrustChanged, StateIdentityChanged:
		return true
	default:
		return false
	}
}

func autoRetry(backend Backend) bool {
	policy, ok := backend.(interface{ AutoRetry() bool })
	return !ok || policy.AutoRetry()
}

func (m *Manager) waitManual(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	case <-m.retry:
		m.beginRetry()
		return true
	}
}

func (m *Manager) waitRetry(ctx context.Context, delay time.Duration) (bool, bool) {
	waitCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	select {
	case <-ctx.Done():
		return false, false
	case <-m.retry:
		m.beginRetry()
		return true, true
	case <-m.scheduler.After(waitCtx, delay):
		return true, false
	}
}

func (m *Manager) publishConnected(generation uint64, target Target, health Health) {
	m.mu.Lock()
	if generation < m.snapshot.Generation {
		m.mu.Unlock()
		return
	}
	target.ID = health.ID
	target.DisplayName = health.DisplayName
	target.Platform = health.Platform
	target.Architecture = health.Architecture
	target.Version = health.Version
	target.Capabilities = append([]string(nil), health.Capabilities...)
	target.Permissions = append([]string(nil), health.Permissions...)
	contact := m.now().UTC().Format(time.RFC3339)
	target.lastSuccessfulAt = contact
	m.snapshot = Snapshot{Generation: generation, Revision: m.snapshot.Revision + 1, State: StateConnected, Target: target, Message: target.DisplayName + " is available.", LastSuccessfulAt: contact, Recovery: RecoveryNone}
	snapshot := cloneSnapshot(m.snapshot)
	m.mu.Unlock()
	m.emitSnapshot(snapshot)
}

func (m *Manager) publishSelectedSnapshot(generation uint64, state State, target Target, message, action string, retryAttempt int) {
	m.mu.Lock()
	if generation < m.snapshot.Generation {
		m.mu.Unlock()
		return
	}
	lastSuccessfulAt := m.snapshot.LastSuccessfulAt
	if lastSuccessfulAt == "" {
		lastSuccessfulAt = target.lastSuccessfulAt
	}
	recovery := RecoveryNone
	if retryAttempt > 0 {
		recovery = RecoveryAutomatic
	}
	m.snapshot = Snapshot{Generation: generation, Revision: m.snapshot.Revision + 1, State: state, Target: cloneTarget(target), Message: message, Action: action, LastSuccessfulAt: lastSuccessfulAt, Stale: target.Kind == "remote" && lastSuccessfulAt != "", RetryAttempt: retryAttempt, Recovery: recovery}
	snapshot := cloneSnapshot(m.snapshot)
	m.mu.Unlock()
	m.emitSnapshot(snapshot)
}

func (m *Manager) selected() (Backend, Target) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.backend, cloneTarget(m.target)
}

func connectingMessage(target Target) string {
	if target.Kind == "remote" {
		return "Connecting to " + target.DisplayName + "."
	}
	return "Connecting to the local scheduler service."
}

func (m *Manager) publishManualFailure(generation uint64, failure *Failure) {
	m.mu.Lock()
	if generation < m.snapshot.Generation {
		m.mu.Unlock()
		return
	}
	target := cloneTarget(m.snapshot.Target)
	m.snapshot = Snapshot{Generation: generation, Revision: m.snapshot.Revision + 1, State: failure.State, Target: target, Message: failure.Message, Action: failure.Action, LastSuccessfulAt: m.snapshot.LastSuccessfulAt, Stale: target.Kind == "remote" && m.snapshot.LastSuccessfulAt != "", Recovery: RecoveryManual}
	snapshot := cloneSnapshot(m.snapshot)
	m.mu.Unlock()
	m.emitSnapshot(snapshot)
}

func (m *Manager) publishAutomaticFailure(generation uint64, state State, message, action string, attempt int, delay time.Duration) {
	m.mu.Lock()
	if generation < m.snapshot.Generation {
		m.mu.Unlock()
		return
	}
	target := cloneTarget(m.snapshot.Target)
	m.snapshot = Snapshot{Generation: generation, Revision: m.snapshot.Revision + 1, State: state, Target: target, Message: message, Action: action, LastSuccessfulAt: m.snapshot.LastSuccessfulAt, Stale: target.Kind == "remote" && m.snapshot.LastSuccessfulAt != "", RetryAttempt: attempt, NextRetryAt: m.now().UTC().Add(delay).Format(time.RFC3339Nano), Recovery: RecoveryAutomatic}
	snapshot := cloneSnapshot(m.snapshot)
	m.mu.Unlock()
	m.emitSnapshot(snapshot)
}

func (m *Manager) emitSnapshot(snapshot Snapshot) {
	if m.observer == nil {
		return
	}
	m.observer.Publish(Event{ID: fmt.Sprintf("connection:%d:%s", snapshot.Generation, snapshot.State), Kind: "connection.changed", Message: snapshot.Message, Generation: snapshot.Generation, OccurredAt: m.now().UTC().Format(time.RFC3339), Snapshot: &snapshot})
}

func (m *Manager) publishDomain(generation uint64, event DomainEvent) bool {
	m.mu.RLock()
	current := m.snapshot.Generation
	connected := m.snapshot.State == StateConnected
	m.mu.RUnlock()
	if generation != current || !connected {
		return false
	}
	if m.observer == nil {
		return true
	}
	sequence := m.sequence.Add(1)
	m.observer.Publish(Event{ID: fmt.Sprintf("%s:%s:%d", event.Kind, event.EntityID, sequence), Kind: event.Kind, Message: eventMessage(event.Kind), Generation: generation, OccurredAt: m.now().UTC().Format(time.RFC3339), EntityID: event.EntityID})
	return true
}

func failureOf(err error) *Failure {
	var failure *Failure
	if errors.As(err, &failure) {
		return failure
	}
	return &Failure{State: StateUnavailable, Message: "The local scheduler service is unavailable.", Action: "Check the service, then try again.", Cause: err}
}

func cloneTarget(value Target) Target {
	value.Capabilities = append([]string(nil), value.Capabilities...)
	value.Permissions = append([]string(nil), value.Permissions...)
	return value
}
func cloneSnapshot(value Snapshot) Snapshot { value.Target = cloneTarget(value.Target); return value }
