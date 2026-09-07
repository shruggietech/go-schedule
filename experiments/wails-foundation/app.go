package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/events"
)

const proofEventName = "proof:event"

type connectionState string

const (
	connectionConnected    connectionState = "connected"
	connectionDisconnected connectionState = "disconnected"
	connectionDegraded     connectionState = "degraded"
)

type nativeOutcome string

const (
	nativeShown  nativeOutcome = "shown"
	nativeFailed nativeOutcome = "failed"
)

type TargetIdentity struct {
	ID          string          `json:"id"`
	DisplayName string          `json:"displayName"`
	Platform    string          `json:"platform"`
	Connection  connectionState `json:"connection"`
	Detail      string          `json:"detail"`
}

type HealthSummary struct {
	Status  connectionState `json:"status"`
	Version string          `json:"version,omitempty"`
	Message string          `json:"message"`
}

type TaskSummary struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Enabled    bool             `json:"enabled"`
	State      domain.TaskState `json:"state"`
	Schedule   string           `json:"schedule"`
	NextRun    string           `json:"nextRun"`
	LastResult string           `json:"lastResult"`
}

type ProofSnapshot struct {
	Target      TargetIdentity `json:"target"`
	Health      HealthSummary  `json:"health"`
	Tasks       []TaskSummary  `json:"tasks"`
	GeneratedAt string         `json:"generatedAt"`
}

type ProofEvent struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Message    string `json:"message"`
	OccurredAt string `json:"occurredAt"`
	TaskID     string `json:"taskId,omitempty"`
}

type NativeActionResult struct {
	Action  string        `json:"action"`
	Outcome nativeOutcome `json:"outcome"`
	Message string        `json:"message"`
}

type daemonClient interface {
	Health(context.Context) (server.HealthResponse, error)
	ListTasks(context.Context, string, string) ([]domain.Task, error)
	StreamEvents(context.Context, func(events.Event)) error
}

type nativeActions interface {
	ShowAbout(context.Context) error
}

type eventEmitter interface {
	Emit(context.Context, string, any)
}

// ProofApp is the Wails-bound application facade. It keeps daemon details and
// raw failures on the Go side of the bridge.
type ProofApp struct {
	daemon  daemonClient
	native  nativeActions
	emitter eventEmitter
	now     func() time.Time

	mu       sync.Mutex
	ctx      context.Context
	cancel   context.CancelFunc
	started  bool
	wg       sync.WaitGroup
	sequence atomic.Uint64
}

func newProofApp(daemon daemonClient, native nativeActions, emitter eventEmitter, now func() time.Time) *ProofApp {
	return &ProofApp{daemon: daemon, native: native, emitter: emitter, now: now}
}

func (a *ProofApp) startup(ctx context.Context) {
	a.mu.Lock()
	if a.started {
		a.mu.Unlock()
		return
	}
	a.ctx, a.cancel = context.WithCancel(ctx)
	a.started = true
	streamCtx := a.ctx
	a.wg.Add(1)
	a.mu.Unlock()

	go func() {
		defer a.wg.Done()
		err := a.daemon.StreamEvents(streamCtx, func(event events.Event) {
			if mapped, ok := a.mapProofEvent(event); ok {
				a.emitter.Emit(streamCtx, proofEventName, mapped)
			}
		})
		if err != nil && !errors.Is(err, context.Canceled) && streamCtx.Err() == nil {
			a.emitter.Emit(streamCtx, proofEventName, ProofEvent{ID: "stream-degraded", Kind: "connection.degraded", Message: "Live updates are temporarily unavailable.", OccurredAt: a.now().UTC().Format(time.RFC3339)})
		}
	}()
}

func (a *ProofApp) shutdown(context.Context) {
	a.mu.Lock()
	cancel := a.cancel
	a.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	a.wg.Wait()
}

func (a *ProofApp) operationContext() (context.Context, context.CancelFunc) {
	a.mu.Lock()
	ctx := a.ctx
	a.mu.Unlock()
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, 2*time.Second)
}

// Snapshot returns a deliberately small, non-sensitive view of daemon state.
func (a *ProofApp) Snapshot() ProofSnapshot {
	ctx, cancel := a.operationContext()
	defer cancel()

	snapshot := ProofSnapshot{
		Target:      TargetIdentity{ID: "local", DisplayName: "This computer", Platform: platformName(), Connection: connectionConnected, Detail: "Local scheduler service"},
		Health:      HealthSummary{Status: connectionConnected, Message: "Scheduler service is available."},
		Tasks:       []TaskSummary{},
		GeneratedAt: a.now().UTC().Format(time.RFC3339),
	}
	health, err := a.daemon.Health(ctx)
	if err != nil {
		snapshot.Target.Connection = connectionDisconnected
		snapshot.Health = HealthSummary{Status: connectionDisconnected, Message: "Scheduler service is unavailable."}
		return snapshot
	}
	snapshot.Health.Version = health.Version

	tasks, err := a.daemon.ListTasks(ctx, "", "")
	if err != nil {
		snapshot.Target.Connection = connectionDegraded
		snapshot.Health.Status = connectionDegraded
		snapshot.Health.Message = "Scheduler service is available, but tasks could not be loaded."
		return snapshot
	}
	for _, task := range tasks {
		snapshot.Tasks = append(snapshot.Tasks, TaskSummary{
			ID: task.ID, Name: task.Name, Enabled: task.Enabled, State: task.State, Schedule: task.ScheduleID, LastResult: "unavailable",
		})
	}
	return snapshot
}

// ShowAbout proves a native desktop interaction without exposing its error.
func (a *ProofApp) ShowAbout() NativeActionResult {
	ctx, cancel := a.operationContext()
	defer cancel()
	if err := a.native.ShowAbout(ctx); err != nil {
		return NativeActionResult{Action: "about-dialog", Outcome: nativeFailed, Message: "The native dialog could not be opened."}
	}
	return NativeActionResult{Action: "about-dialog", Outcome: nativeShown, Message: "About dialog opened."}
}

func (a *ProofApp) mapProofEvent(event events.Event) (ProofEvent, bool) {
	stamp := a.now().UTC().Format(time.RFC3339)
	sequence := a.sequence.Add(1)
	if event.Kind == events.KindTask && event.Task != nil {
		kind := "task." + string(event.Task.Verb)
		return ProofEvent{ID: fmt.Sprintf("%s:%s:%d", kind, event.Task.ID, sequence), Kind: kind, Message: "Task state changed.", OccurredAt: stamp, TaskID: event.Task.ID}, true
	}
	if event.Kind != "" {
		kind := string(event.Kind) + ".changed"
		return ProofEvent{ID: fmt.Sprintf("%s:%d", kind, sequence), Kind: kind, Message: "Scheduler state changed.", OccurredAt: stamp}, true
	}
	return ProofEvent{}, false
}

func platformName() string {
	if runtime.GOOS == "darwin" {
		return "macos"
	}
	return runtime.GOOS
}
