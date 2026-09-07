package main

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/events"
)

type fakeDaemon struct {
	health      server.HealthResponse
	healthErr   error
	tasks       []domain.Task
	listErr     error
	events      []events.Event
	streamStart chan struct{}
	mu          sync.Mutex
	streams     int
}

func (f *fakeDaemon) Health(context.Context) (server.HealthResponse, error) {
	return f.health, f.healthErr
}

func (f *fakeDaemon) ListTasks(context.Context, string, string) ([]domain.Task, error) {
	return f.tasks, f.listErr
}

func (f *fakeDaemon) StreamEvents(ctx context.Context, onEvent func(events.Event)) error {
	f.mu.Lock()
	f.streams++
	f.mu.Unlock()
	if f.streamStart != nil {
		close(f.streamStart)
	}
	for _, event := range f.events {
		onEvent(event)
	}
	<-ctx.Done()
	return ctx.Err()
}

type fakeNative struct {
	err   error
	calls int
}

func (f *fakeNative) ShowAbout(context.Context) error {
	f.calls++
	return f.err
}

type emittedEvent struct {
	name string
	data any
}

type fakeEmitter struct {
	events chan emittedEvent
}

func (f *fakeEmitter) Emit(_ context.Context, name string, data any) {
	f.events <- emittedEvent{name: name, data: data}
}

func fixedNow() time.Time {
	return time.Date(2026, time.September, 7, 5, 30, 0, 0, time.UTC)
}

func TestSnapshotMapsHealthAndTasks(t *testing.T) {
	daemon := &fakeDaemon{
		health: server.HealthResponse{Status: "ok", Version: "v1.1.1"},
		tasks: []domain.Task{{
			ID: "backup", Name: "Nightly backup", Enabled: true,
			State: domain.TaskActive, ScheduleID: "schedule-1",
		}},
	}
	app := newProofApp(daemon, &fakeNative{}, &fakeEmitter{events: make(chan emittedEvent, 1)}, fixedNow)

	snapshot := app.Snapshot()

	if snapshot.Target.DisplayName != "This computer" || snapshot.Target.Connection != connectionConnected {
		t.Fatalf("target = %+v", snapshot.Target)
	}
	if snapshot.Health.Status != connectionConnected || snapshot.Health.Version != "v1.1.1" {
		t.Fatalf("health = %+v", snapshot.Health)
	}
	if len(snapshot.Tasks) != 1 || snapshot.Tasks[0].ID != "backup" || snapshot.Tasks[0].Schedule != "schedule-1" {
		t.Fatalf("tasks = %+v", snapshot.Tasks)
	}
	if snapshot.GeneratedAt != "2026-09-07T05:30:00Z" {
		t.Fatalf("generated_at = %q", snapshot.GeneratedAt)
	}
}

func TestSnapshotMapsConnectionAndListFailures(t *testing.T) {
	t.Run("health failure", func(t *testing.T) {
		app := newProofApp(&fakeDaemon{healthErr: errors.New("dial secret endpoint")}, &fakeNative{}, &fakeEmitter{events: make(chan emittedEvent, 1)}, fixedNow)
		snapshot := app.Snapshot()
		if snapshot.Target.Connection != connectionDisconnected || snapshot.Health.Status != connectionDisconnected {
			t.Fatalf("snapshot = %+v", snapshot)
		}
		if snapshot.Health.Message == "" || snapshot.Health.Message == "dial secret endpoint" {
			t.Fatalf("unsafe or missing message = %q", snapshot.Health.Message)
		}
	})

	t.Run("list failure", func(t *testing.T) {
		app := newProofApp(&fakeDaemon{health: server.HealthResponse{Status: "ok"}, listErr: errors.New("broken")}, &fakeNative{}, &fakeEmitter{events: make(chan emittedEvent, 1)}, fixedNow)
		snapshot := app.Snapshot()
		if snapshot.Target.Connection != connectionDegraded || snapshot.Health.Status != connectionDegraded {
			t.Fatalf("snapshot = %+v", snapshot)
		}
		if len(snapshot.Tasks) != 0 {
			t.Fatalf("tasks = %+v", snapshot.Tasks)
		}
	})
}

func TestShowAboutMapsNativeActionResult(t *testing.T) {
	native := &fakeNative{}
	app := newProofApp(&fakeDaemon{}, native, &fakeEmitter{events: make(chan emittedEvent, 1)}, fixedNow)
	app.startup(context.Background())
	t.Cleanup(func() { app.shutdown(context.Background()) })

	if result := app.ShowAbout(); result.Outcome != nativeShown || native.calls != 1 {
		t.Fatalf("result = %+v, calls = %d", result, native.calls)
	}

	native.err = errors.New("desktop unavailable")
	if result := app.ShowAbout(); result.Outcome != nativeFailed || result.Message == "desktop unavailable" {
		t.Fatalf("unsafe failure result = %+v", result)
	}
}

func TestEventStreamStartsOnceEmitsAndStops(t *testing.T) {
	started := make(chan struct{})
	daemon := &fakeDaemon{
		streamStart: started,
		events: []events.Event{{Kind: events.KindTask, Task: &events.TaskEvent{
			Verb: events.VerbUpdated, ID: "backup",
		}}},
	}
	emitter := &fakeEmitter{events: make(chan emittedEvent, 2)}
	app := newProofApp(daemon, &fakeNative{}, emitter, fixedNow)

	app.startup(context.Background())
	app.startup(context.Background())

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("event stream did not start")
	}
	select {
	case got := <-emitter.events:
		if got.name != proofEventName {
			t.Fatalf("event name = %q", got.name)
		}
		event, ok := got.data.(ProofEvent)
		if !ok || event.Kind != "task.updated" || event.TaskID != "backup" {
			t.Fatalf("event = %#v", got.data)
		}
	case <-time.After(time.Second):
		t.Fatal("event was not emitted")
	}

	app.shutdown(context.Background())
	daemon.mu.Lock()
	defer daemon.mu.Unlock()
	if daemon.streams != 1 {
		t.Fatalf("stream count = %d", daemon.streams)
	}
}
