package connection

import (
	"context"
	"errors"
	"testing"
	"time"
)

func nextState(t *testing.T, events <-chan Event, state State) Event {
	t.Helper()
	for {
		select {
		case event := <-events:
			if event.Snapshot != nil && event.Snapshot.State == state {
				return event
			}
		case <-time.After(time.Second):
			t.Fatalf("state %s not published", state)
		}
	}
}

func firstStream(t *testing.T, backend *backendFake) (chan DomainEvent, chan error) {
	t.Helper()
	deadline := time.After(time.Second)
	for {
		backend.mu.Lock()
		if len(backend.streams) > 0 {
			eventsCh, errorsCh := backend.streams[0], backend.streamErrors[0]
			backend.mu.Unlock()
			return eventsCh, errorsCh
		}
		backend.mu.Unlock()
		select {
		case <-deadline:
			t.Fatal("event stream did not start")
		default:
		}
	}
}

func TestManagerConnectsPublishesEventsAndRejectsStaleGeneration(t *testing.T) {
	backend := &backendFake{results: []backendResult{{health: Health{Version: "1.0.0", Capabilities: []string{"tasks"}, Permissions: []string{"read"}}}}}
	observer := observerFake{events: make(chan Event, 16)}
	manager := newManager(backend, observer, timerScheduler{}, func() time.Time { return time.Unix(1, 0) })
	ctx, cancel := context.WithCancel(context.Background())
	manager.Start(ctx)
	connected := nextState(t, observer.events, StateConnected)
	if connected.Snapshot.Target.Version != "1.0.0" || connected.Generation != 1 {
		t.Fatalf("event=%+v", connected)
	}
	stream, _ := firstStream(t, backend)
	stream <- DomainEvent{Kind: "task.updated", EntityID: "task-1"}
	select {
	case event := <-observer.events:
		if event.Kind != "task.updated" || event.Generation != 1 {
			t.Fatalf("event=%+v", event)
		}
	case <-time.After(time.Second):
		t.Fatal("domain event missing")
	}
	cancel()
	stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second)
	defer stopCancel()
	if err := manager.Stop(stopCtx); err != nil {
		t.Fatal(err)
	}
	stream <- DomainEvent{Kind: "task.updated", EntityID: "stale"}
}

func TestManagerExactRetryCadenceAndManualInterruption(t *testing.T) {
	failure := &Failure{State: StateUnavailable, Message: "Unavailable.", Action: "Retry."}
	backend := &backendFake{results: []backendResult{{err: failure}, {err: failure}, {err: failure}, {health: Health{Version: "1.0.0"}}}}
	observer := observerFake{events: make(chan Event, 32)}
	scheduler := schedulerFake{delays: make(chan time.Duration, 4), releases: make(chan func(), 4)}
	manager := newManager(backend, observer, scheduler, time.Now)
	manager.Start(context.Background())
	wants := []time.Duration{250 * time.Millisecond, time.Second, 5 * time.Second}
	for index, want := range wants {
		nextState(t, observer.events, StateUnavailable)
		if got := <-scheduler.delays; got != want {
			t.Fatalf("delay %d=%s want %s", index, got, want)
		}
		release := <-scheduler.releases
		if index == 1 {
			if !manager.Retry() || !manager.Retry() {
				t.Fatal("manual retry rejected")
			}
		} else {
			release()
		}
	}
	nextState(t, observer.events, StateConnected)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := manager.Stop(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestManagerTerminalFailureWaitsForManualRetry(t *testing.T) {
	backend := &backendFake{results: []backendResult{{err: &Failure{State: StateAccessDenied, Message: "Denied.", Action: "Fix permissions."}}, {health: Health{Version: "1.0.0"}}}}
	observer := observerFake{events: make(chan Event, 8)}
	manager := newManager(backend, observer, timerScheduler{}, time.Now)
	manager.Start(context.Background())
	nextState(t, observer.events, StateAccessDenied)
	if !manager.Retry() {
		t.Fatal("retry rejected")
	}
	nextState(t, observer.events, StateConnected)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = manager.Stop(ctx)
}

func TestManagerPublishesEveryConnectionState(t *testing.T) {
	transient := &Failure{State: StateTimedOut, Message: "Timed out.", Action: "Retry."}
	backend := &backendFake{results: []backendResult{{err: transient}, {health: Health{Version: "1.0.0"}}}}
	observer := observerFake{events: make(chan Event, 12)}
	scheduler := schedulerFake{delays: make(chan time.Duration, 1), releases: make(chan func(), 1)}
	manager := newManager(backend, observer, scheduler, time.Now)
	manager.Start(context.Background())
	nextState(t, observer.events, StateConnecting)
	nextState(t, observer.events, StateTimedOut)
	<-scheduler.delays
	(<-scheduler.releases)()
	nextState(t, observer.events, StateRecovering)
	nextState(t, observer.events, StateConnected)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := manager.Stop(ctx); err != nil {
		t.Fatal(err)
	}

	for _, state := range []State{StateUnavailable, StateAccessDenied, StateIncompatible} {
		t.Run(string(state), func(t *testing.T) {
			failureBackend := &backendFake{results: []backendResult{{err: &Failure{State: state, Message: "Safe failure.", Action: "Try again."}}}}
			failureObserver := observerFake{events: make(chan Event, 4)}
			failureManager := newManager(failureBackend, failureObserver, timerScheduler{}, time.Now)
			failureManager.Start(context.Background())
			nextState(t, failureObserver.events, state)
			stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second)
			defer stopCancel()
			if err := failureManager.Stop(stopCtx); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSnapshotIsDetachedFromManagerState(t *testing.T) {
	manager := newManager(&backendFake{}, nil, timerScheduler{}, time.Now)
	manager.snapshot.Target.Capabilities = []string{"tasks"}
	first := manager.Snapshot()
	first.Target.Capabilities[0] = "modified"
	if got := manager.Snapshot().Target.Capabilities[0]; got != "tasks" {
		t.Fatalf("manager snapshot was mutated through caller: %q", got)
	}
}

func TestManagerRejectsDomainEventFromStaleGeneration(t *testing.T) {
	observer := observerFake{events: make(chan Event, 1)}
	manager := newManager(&backendFake{}, observer, timerScheduler{}, time.Now)
	manager.snapshot = Snapshot{Generation: 2, State: StateConnected, Target: localTarget()}
	manager.publishDomain(1, DomainEvent{Kind: "task.updated", EntityID: "stale"})
	select {
	case event := <-observer.events:
		t.Fatalf("stale event was published: %+v", event)
	default:
	}
}

func TestManagerDegradesAndRecoversAfterEventFailure(t *testing.T) {
	backend := &backendFake{results: []backendResult{{health: Health{Version: "1.0.0"}}}}
	observer := observerFake{events: make(chan Event, 16)}
	scheduler := schedulerFake{delays: make(chan time.Duration, 2), releases: make(chan func(), 2)}
	manager := newManager(backend, observer, scheduler, time.Now)
	manager.Start(context.Background())
	nextState(t, observer.events, StateConnected)
	_, errorsCh := firstStream(t, backend)
	errorsCh <- errors.New("private endpoint")
	event := nextState(t, observer.events, StateDegraded)
	if event.Message != "Live updates are temporarily unavailable." {
		t.Fatalf("event=%+v", event)
	}
	if got := <-scheduler.delays; got != 250*time.Millisecond {
		t.Fatalf("delay=%s", got)
	}
	(<-scheduler.releases)()
	nextState(t, observer.events, StateConnected)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = manager.Stop(ctx)
}

func TestManagerLifecycleOneHundredCycles(t *testing.T) {
	for index := 0; index < 100; index++ {
		backend := &backendFake{results: []backendResult{{health: Health{Version: "1.0.0"}}}}
		observer := observerFake{events: make(chan Event, 4)}
		manager := newManager(backend, observer, timerScheduler{}, time.Now)
		manager.Start(context.Background())
		nextState(t, observer.events, StateConnected)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		if err := manager.Stop(ctx); err != nil {
			t.Fatalf("cycle %d: %v", index, err)
		}
		cancel()
		if manager.Retry() {
			t.Fatalf("cycle %d accepted retry after close", index)
		}
	}
}
