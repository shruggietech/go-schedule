package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/desktop/taskgroup"
	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/ipc"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

type wailsEmitter struct{}

func (wailsEmitter) Emit(ctx context.Context, name string, value any) {
	wailsRuntime.EventsEmit(ctx, name, value)
}

type wailsNative struct{}

func (wailsNative) Quit(ctx context.Context) { wailsRuntime.Quit(ctx) }

func main() {
	cfg, err := config.Load("")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	daemon := client.New(ipc.Endpoint(cfg))
	backend := connection.NewLocalBackend(daemon)
	app := newApp(backend, wailsEmitter{}, wailsNative{}, taskgroup.NewService(taskgroup.NewLocalBackend(daemon)))
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
