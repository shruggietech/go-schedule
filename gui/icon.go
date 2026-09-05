package gui

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

// iconBytes is the full-resolution application icon: the go-schedule brand mark
// (a terminal prompt over a five-field cron schedule) as a transparent PNG, used
// for large surfaces like the macOS dock.
//
//go:embed assets/icon.png
var iconBytes []byte

// windowIconBytes is the brand's reduced mark (the prompt chevron and cursor
// line only) on a Night tile, the same purpose-rendered artwork the favicons
// use. Fyne scales the window icon down to ~16px for the title bar; the reduced
// mark stays legible where the full five-cell mark would not, so the title bar
// and taskbar use this tile instead of downscaling the full-resolution source.
//
//go:embed assets/icon-window.png
var windowIconBytes []byte

// appIcon is the full-resolution mark, used for the application-level icon
// (macOS dock and other large surfaces).
var appIcon = fyne.NewStaticResource("icon.png", iconBytes)

// windowIcon is the small crisp mark used for window decorations (title bar and
// the Windows taskbar/alt-tab), where the image is shown at 16–64px.
var windowIcon = fyne.NewStaticResource("icon-window.png", windowIconBytes)
