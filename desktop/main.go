package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/shruggietech/go-schedule/desktop/agentaccess"
	"github.com/shruggietech/go-schedule/desktop/automation"
	"github.com/shruggietech/go-schedule/desktop/bundles"
	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/desktop/connections"
	"github.com/shruggietech/go-schedule/desktop/notifications"
	"github.com/shruggietech/go-schedule/desktop/operations"
	"github.com/shruggietech/go-schedule/desktop/remotepairing"
	"github.com/shruggietech/go-schedule/desktop/search"
	"github.com/shruggietech/go-schedule/desktop/settings"
	"github.com/shruggietech/go-schedule/desktop/systems"
	"github.com/shruggietech/go-schedule/desktop/taskgroup"
	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/autostart"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/clientsecret"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/desktopcontrol"
	"github.com/shruggietech/go-schedule/internal/ipc"
	"github.com/shruggietech/go-schedule/internal/service"
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

func shouldAutoSpawnInstalledService(state service.State, err error) bool {
	if runtime.GOOS != "windows" && runtime.GOOS != "linux" {
		return true
	}
	// A deliberate stop must survive opening the GUI. When SCM cannot be
	// queried, failing closed also avoids spawning a competing daemon.
	return err == nil && state == service.StateNotInstalled
}

func localConfigPath(goos string) string {
	if goos == "windows" || goos == "linux" {
		return config.DefaultPath()
	}
	return ""
}

func main() {
	instance, owner, instanceErr := desktopcontrol.ClaimGUI()
	if instanceErr != nil {
		fmt.Fprintln(os.Stderr, instanceErr)
		os.Exit(1)
	}
	if !owner {
		_, _ = desktopcontrol.SignalGUI()
		return
	}
	defer instance.Close() //nolint:errcheck // process-lifetime mutex handle
	cfg, err := config.Load(localConfigPath(runtime.GOOS))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	localDaemon := client.New(ipc.Endpoint(cfg))
	if state, statusErr := service.QueryState(); shouldAutoSpawnInstalledService(state, statusErr) {
		_ = ensureBundledDaemon(func(ctx context.Context) error {
			_, err := localDaemon.Health(ctx)
			return err
		}, autostart.SpawnDaemon)
	}
	router := client.NewSwitchable(localDaemon)
	backend := connection.NewLocalBackend(localDaemon)
	profileStore := clientprofile.NewStore("")
	secretStore := clientsecret.New()
	native := wailsNative{}
	app := newApp(backend, wailsEmitter{}, native, appServices{tasks: taskgroup.NewService(taskgroup.NewLocalBackend(router)), automation: automation.NewService(automation.NewLocalBackend(router)), bundles: bundles.NewService(bundles.NewLocalBackend(router)), operations: operations.NewService(operations.NewLocalBackend(router)), notifications: notifications.NewService(notifications.NewLocalBackend(router)), settings: settings.NewService(settings.NewLocalBackend(localDaemon), native), agentAccess: agentaccess.NewService(agentaccess.NewLocalBackend(router), native), remotePairing: remotepairing.NewWithStores(secretStore, profileStore)})
	if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
		monitor := desktopcontrol.NewMonitor(func(ctx context.Context) error {
			_, err := localDaemon.Health(ctx)
			return err
		})
		app.localService = &monitor
	}
	app.connections = connections.New(profileStore, secretStore, localDaemon, router, app.manager)
	app.systems = systems.New(profileStore, secretStore, localDaemon)
	app.search = search.New(profileStore, secretStore, localDaemon)
	if err := wails.Run(&options.App{
		Title: "go-schedule", Width: 1440, Height: 900, MinWidth: 800, MinHeight: 600,
		BackgroundColour: &options.RGBA{R: 245, G: 247, B: 250, A: 1},
		Windows:          &windows.Options{Theme: windows.SystemDefault},
		AssetServer:      &assetserver.Options{Assets: assets},
		OnStartup:        app.startup, OnShutdown: app.shutdown, Bind: []any{app},
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
