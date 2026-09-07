package main

import (
	"context"
	"time"

	"github.com/shruggietech/go-schedule/desktop/connection"
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
	emitter eventEmitter
	native  nativeRuntime
	ctx     context.Context
}

func newApp(backend connection.Backend, emitter eventEmitter, native nativeRuntime) *App {
	app := &App{emitter: emitter, native: native}
	app.manager = connection.NewManager(backend, appObserver{app: app})
	return app
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
