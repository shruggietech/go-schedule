package connection

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const attemptTimeout = 2 * time.Second

var retryDelays = [...]time.Duration{250 * time.Millisecond, time.Second, 5 * time.Second}

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
	backend   Backend
	scheduler Scheduler
	observer  Observer
	now       func() time.Time

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
	return newManager(backend, observer, timerScheduler{}, time.Now)
}

func newManager(backend Backend, observer Observer, scheduler Scheduler, now func() time.Time) *Manager {
	target := localTarget()
	return &Manager{backend: backend, target: target, observer: observer, scheduler: scheduler, now: now, retry: make(chan struct{}, 1), done: make(chan struct{}), snapshot: Snapshot{State: StateConnecting, Target: target, Message: connectingMessage(target)}}
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
	m.snapshot = Snapshot{Generation: m.snapshot.Generation, Revision: m.snapshot.Revision + 1, State: StateConnecting, Target: cloneTarget(target), Message: connectingMessage(target)}
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
	retryIndex := 0

connectionLoop:
	for {
		generation++
		backend, target := m.selected()
		state := StateConnecting
		if generation > 1 {
			state = StateRecovering
		}
		m.publishSelectedSnapshot(generation, state, target, connectingMessage(target), "")
		attemptCtx, cancelAttempt := context.WithTimeout(ctx, attemptTimeout)
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
			retryIndex = 0
			continue connectionLoop
		}
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			failure := failureOf(err)
			m.publishSnapshot(generation, failure.State, failure.Message, failure.Action)
			if failure.State == StateAccessDenied || failure.State == StateIncompatible || !autoRetry(backend) {
				if !m.waitManual(ctx) {
					return
				}
				retryIndex = 0
				continue
			}
			delay := retryDelays[min(retryIndex, len(retryDelays)-1)]
			if retryIndex < len(retryDelays)-1 {
				retryIndex++
			}
			proceed, manual := m.waitRetry(ctx, delay)
			if !proceed {
				return
			}
			if manual {
				retryIndex = 0
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
				retryIndex = 0
				continue connectionLoop
			case <-streamActivity:
				retryIndex = 0
			case <-streamDone:
				cancelStream()
				if ctx.Err() != nil {
					return
				}
				m.publishSnapshot(generation, StateDegraded, "Live updates are temporarily unavailable.", "Try again.")
				if !autoRetry(backend) {
					if !m.waitManual(ctx) {
						return
					}
					retryIndex = 0
					continue connectionLoop
				}
				if streamHadActivity.Load() {
					retryIndex = 0
				}
				delay := retryDelays[min(retryIndex, len(retryDelays)-1)]
				if retryIndex < len(retryDelays)-1 {
					retryIndex++
				}
				proceed, manual := m.waitRetry(ctx, delay)
				if !proceed {
					return
				}
				if manual {
					retryIndex = 0
				}
				continue connectionLoop
			}
		}
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
	m.snapshot = Snapshot{Generation: generation, Revision: m.snapshot.Revision + 1, State: StateConnected, Target: target, Message: target.DisplayName + " is available.", LastSuccessfulAt: m.now().UTC().Format(time.RFC3339)}
	snapshot := cloneSnapshot(m.snapshot)
	m.mu.Unlock()
	m.emitSnapshot(snapshot)
}

func (m *Manager) publishSelectedSnapshot(generation uint64, state State, target Target, message, action string) {
	m.mu.Lock()
	if generation < m.snapshot.Generation {
		m.mu.Unlock()
		return
	}
	m.snapshot = Snapshot{Generation: generation, Revision: m.snapshot.Revision + 1, State: state, Target: cloneTarget(target), Message: message, Action: action, LastSuccessfulAt: m.snapshot.LastSuccessfulAt}
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

func (m *Manager) publishSnapshot(generation uint64, state State, message, action string) {
	m.mu.Lock()
	if generation < m.snapshot.Generation {
		m.mu.Unlock()
		return
	}
	target := cloneTarget(m.snapshot.Target)
	m.snapshot = Snapshot{Generation: generation, Revision: m.snapshot.Revision + 1, State: state, Target: target, Message: message, Action: action, LastSuccessfulAt: m.snapshot.LastSuccessfulAt}
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
