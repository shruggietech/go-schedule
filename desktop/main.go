package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/shruggietech/go-schedule/desktop/automation"
	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/desktop/notifications"
	"github.com/shruggietech/go-schedule/desktop/operations"
	"github.com/shruggietech/go-schedule/desktop/settings"
	"github.com/shruggietech/go-schedule/desktop/taskgroup"
	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/autostart"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/ipc"
)

//go:embed all:frontend/dist
var assets embed.FS

type wailsEmitter struct{}

func (wailsEmitter) Emit(ctx context.Context, name string, value any) {
	wailsRuntime.EventsEmit(ctx, name, value)
}

type wailsNative struct{}

func (wailsNative) Quit(ctx context.Context) { wailsRuntime.Quit(ctx) }
func (wailsNative) ClipboardSetText(ctx context.Context, value string) error {
	return wailsRuntime.ClipboardSetText(ctx, value)
}
func (wailsNative) BrowserOpenURL(ctx context.Context, value string) error {
	wailsRuntime.BrowserOpenURL(ctx, value)
	return nil
}

func ensureBundledDaemon(ping func(context.Context) error, spawn func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return autostart.EnsureRunning(ctx, ping, spawn)
}

func main() {
	cfg, err := config.Load("")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	daemon := client.New(ipc.Endpoint(cfg))
	_ = ensureBundledDaemon(func(ctx context.Context) error {
		_, err := daemon.Health(ctx)
		return err
	}, autostart.SpawnDaemon)
	backend := connection.NewLocalBackend(daemon)
	native := wailsNative{}
	app := newApp(backend, wailsEmitter{}, native, appServices{tasks: taskgroup.NewService(taskgroup.NewLocalBackend(daemon)), automation: automation.NewService(automation.NewLocalBackend(daemon)), operations: operations.NewService(operations.NewLocalBackend(daemon)), notifications: notifications.NewService(notifications.NewLocalBackend(daemon)), settings: settings.NewService(settings.NewLocalBackend(daemon), native)})
	if err := wails.Run(&options.App{
		Title: "go-schedule", Width: 1440, Height: 900, MinWidth: 900, MinHeight: 650,
		BackgroundColour: &options.RGBA{R: 245, G: 247, B: 250, A: 1},
		AssetServer:      &assetserver.Options{Assets: assets},
		OnStartup:        app.startup, OnShutdown: app.shutdown, Bind: []any{app},
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
