package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/ipc"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

type wailsEmitter struct{}

func (wailsEmitter) Emit(ctx context.Context, name string, data any) {
	runtime.EventsEmit(ctx, name, data)
}

type wailsNative struct{}

func (wailsNative) ShowAbout(ctx context.Context) error {
	_, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   "About go-schedule",
		Message: "go-schedule Wails foundation proof",
	})
	return err
}

func main() {
	cfg, err := config.Load("")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	app := newProofApp(client.New(ipc.Endpoint(cfg)), wailsNative{}, wailsEmitter{}, time.Now)
	if err := wails.Run(&options.App{
		Title:            "go-schedule",
		Width:            1440,
		Height:           900,
		MinWidth:         900,
		MinHeight:        650,
		BackgroundColour: &options.RGBA{R: 245, G: 247, B: 250, A: 1},
		AssetServer:      &assetserver.Options{Assets: assets},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind:             []any{app},
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
