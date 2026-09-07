package main

import (
	"context"
	"time"

	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/desktop/taskgroup"
)

const desktopEventName = "desktop:event"

type eventEmitter interface {
	Emit(context.Context, string, any)
}
type nativeRuntime interface{ Quit(context.Context) }

// ActionResult gives the frontend a stable result without exposing native errors.
type ActionResult struct {
	Action  string `json:"action"`
	Outcome string `json:"outcome"`
	Message string `json:"message"`
}

// App is the deliberately small Wails bridge facade.
type App struct {
	manager *connection.Manager
	tasks   *taskgroup.Service
	emitter eventEmitter
	native  nativeRuntime
	ctx     context.Context
}

func newApp(backend connection.Backend, emitter eventEmitter, native nativeRuntime, taskServices ...*taskgroup.Service) *App {
	app := &App{emitter: emitter, native: native}
	if len(taskServices) > 0 {
		app.tasks = taskServices[0]
	}
	app.manager = connection.NewManager(backend, appObserver{app: app})
	return app
}

func (a *App) Workspace() taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "load", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.Workspace(a.ctx)
}

func (a *App) Task(id string) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "load", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.Task(a.ctx, id)
}

func (a *App) PreviewTask(draft taskgroup.TaskDraft) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "preview", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.PreviewTask(a.ctx, draft)
}
func (a *App) SaveTask(draft taskgroup.TaskDraft) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "save_task", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.SaveTask(a.ctx, draft)
}
func (a *App) RunTask(id string) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "run_task", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.RunTask(a.ctx, id)
}
func (a *App) SetTaskEnabled(id string, enabled bool) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "toggle_task", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.SetTaskEnabled(a.ctx, id, enabled)
}
func (a *App) DeleteTask(id string) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "delete_task", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.DeleteTask(a.ctx, id)
}
func (a *App) SaveGroup(draft taskgroup.GroupDraft) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "save_group", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.SaveGroup(a.ctx, draft)
}
func (a *App) SetGroupEnabled(id string, enabled bool) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "toggle_group", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.SetGroupEnabled(a.ctx, id, enabled)
}
func (a *App) DeleteGroup(id string) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "delete_group", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.DeleteGroup(a.ctx, id)
}

type appObserver struct{ app *App }

func (o appObserver) Publish(event connection.Event) {
	if o.app.emitter != nil && o.app.ctx != nil {
		o.app.emitter.Emit(o.app.ctx, desktopEventName, event)
	}
}

func (a *App) startup(ctx context.Context) { a.ctx = ctx; a.manager.Start(ctx) }

func (a *App) shutdown(context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = a.manager.Stop(ctx)
}

// Snapshot returns only the transport-neutral connection contract.
func (a *App) Snapshot() connection.Snapshot { return a.manager.Snapshot() }

// RetryConnection requests a fresh, coalesced generation.
func (a *App) RetryConnection() ActionResult {
	if !a.manager.Retry() {
		return ActionResult{Action: "retry", Outcome: "rejected", Message: "The application is closing."}
	}
	return ActionResult{Action: "retry", Outcome: "accepted", Message: "Trying the local scheduler service again."}
}

// Quit performs the native application exit action.
func (a *App) Quit() ActionResult {
	if a.native == nil || a.ctx == nil {
		return ActionResult{Action: "quit", Outcome: "unavailable", Message: "Native exit is unavailable."}
	}
	a.native.Quit(a.ctx)
	return ActionResult{Action: "quit", Outcome: "accepted", Message: "Closing go-schedule."}
}
