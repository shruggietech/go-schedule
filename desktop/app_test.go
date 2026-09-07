package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/desktop/operations"
	"github.com/shruggietech/go-schedule/desktop/taskgroup"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type appBackend struct{}

func (appBackend) Health(context.Context) (connection.Health, error) {
	return connection.Health{Version: "1.2.0", Capabilities: []string{"tasks"}, Permissions: []string{"read"}}, nil
}
func (appBackend) StreamEvents(ctx context.Context, publish func(connection.DomainEvent)) error {
	<-ctx.Done()
	return ctx.Err()
}

type appEmitter struct {
	mu     sync.Mutex
	events []connection.Event
}

func (e *appEmitter) Emit(_ context.Context, name string, value any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if name == desktopEventName {
		e.events = append(e.events, value.(connection.Event))
	}
}

type appNative struct{ quit bool }

func (n *appNative) Quit(context.Context) { n.quit = true }

type facadeTaskBackend struct{ taskgroup.Backend }

type facadeOperationsBackend struct{ operations.Backend }

func (facadeOperationsBackend) GetCalendar(_ context.Context, from, to time.Time) (server.CalendarResponse, error) {
	return server.CalendarResponse{From: from, To: to, Occurrences: []server.Occurrence{}}, nil
}
func (facadeOperationsBackend) ListRuns(context.Context, string, int) ([]domain.Run, error) {
	return []domain.Run{}, nil
}
func (facadeOperationsBackend) ListLogs(context.Context, string, int) (server.LogsResponse, error) {
	return server.LogsResponse{Logs: []domain.LogRecord{}}, nil
}
func (facadeOperationsBackend) ListAlerts(context.Context, bool) ([]domain.Alert, error) {
	return []domain.Alert{}, nil
}
func (facadeOperationsBackend) AckAlert(context.Context, string) error { return nil }

func (facadeTaskBackend) ListTaskDetails(context.Context, string, string) ([]server.TaskResponse, error) {
	return []server.TaskResponse{}, nil
}

func (facadeTaskBackend) ListGroups(context.Context) ([]domain.Group, error) {
	return []domain.Group{}, nil
}

func TestAppFacadeStartsSnapshotsRetriesAndQuits(t *testing.T) {
	emitter := &appEmitter{}
	native := &appNative{}
	app := newApp(appBackend{}, emitter, native)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.startup(ctx)
	deadline := time.After(time.Second)
	for app.Snapshot().State != connection.StateConnected {
		select {
		case <-deadline:
			t.Fatal("connection did not start")
		default:
		}
	}
	snapshot := app.Snapshot()
	if snapshot.Target.DisplayName != "This computer" || snapshot.Target.Version != "1.2.0" {
		t.Fatalf("snapshot=%+v", snapshot)
	}
	if result := app.RetryConnection(); result.Outcome != "accepted" {
		t.Fatalf("retry=%+v", result)
	}
	if result := app.Quit(); result.Outcome != "accepted" || !native.quit {
		t.Fatalf("quit=%+v native=%+v", result, native)
	}
	app.shutdown(context.Background())
	if result := app.RetryConnection(); result.Outcome != "rejected" {
		t.Fatalf("post-close retry=%+v", result)
	}
	for _, event := range emitter.events {
		if event.Message == "" || event.Generation == 0 {
			t.Fatalf("unsafe event=%+v", event)
		}
	}
}

func TestAppFacadeExposesSafeTaskWorkspace(t *testing.T) {
	app := newApp(appBackend{}, nil, nil, appServices{tasks: taskgroup.NewService(facadeTaskBackend{})})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.startup(ctx)
	result := app.Workspace()
	if result.Outcome != "accepted" || result.Workspace == nil || result.Workspace.Tasks == nil || result.Workspace.Groups == nil {
		t.Fatalf("workspace=%+v", result)
	}
	app.shutdown(context.Background())
}

func TestAppFacadeExposesScheduleActivityAndAcknowledgement(t *testing.T) {
	service := operations.NewService(facadeOperationsBackend{})
	app := newApp(appBackend{}, nil, nil, appServices{operations: service})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.startup(ctx)
	if result := app.ScheduleWindow(7); result.Outcome != "accepted" || result.Schedule == nil {
		t.Fatalf("schedule=%+v", result)
	}
	if result := app.ActivityWorkspace(); result.Outcome != "accepted" || result.Activity == nil {
		t.Fatalf("activity=%+v", result)
	}
	if result := app.AcknowledgeAlert("alert-1"); result.Outcome != "accepted" || result.Activity == nil {
		t.Fatalf("acknowledge=%+v", result)
	}
	app.shutdown(context.Background())
}
