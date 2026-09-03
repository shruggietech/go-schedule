package gui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/buildinfo"
)

func (a *App) buildInfoTab() fyne.CanvasObject {
	return buildInfoContentWithScroll(buildinfo.Version, func(content fyne.CanvasObject) fyne.CanvasObject {
		return a.newVScroll(content)
	})
}

func buildInfoContent(version string) fyne.CanvasObject {
	return buildInfoContentWithScroll(version, func(content fyne.CanvasObject) fyne.CanvasObject {
		return newSensitiveVScroll(content, func() float64 { return defaultScrollSensitivity })
	})
}

func buildInfoContentWithScroll(version string, wrap func(fyne.CanvasObject) fyne.CanvasObject) fyne.CanvasObject {
	mark := canvas.NewImageFromResource(appIcon)
	mark.FillMode = canvas.ImageFillContain
	mark.SetMinSize(fyne.NewSize(180, 180))

	title := widget.NewLabelWithStyle("go-schedule", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	versionLabel := widget.NewLabel("Version " + version)
	versionLabel.Alignment = fyne.TextAlignCenter
	attribution := widget.NewLabel("Built and maintained by ShruggieTech")
	attribution.Alignment = fyne.TextAlignCenter

	content := container.NewVBox(
		container.NewCenter(mark),
		title,
		versionLabel,
		widget.NewSeparator(),
		attribution,
		centeredInfoLink("ShruggieTech", &url.URL{Scheme: "https", Host: "shruggie.tech"}),
		centeredInfoLink("Source repository", &url.URL{Scheme: "https", Host: "github.com", Path: "/shruggietech/go-schedule"}),
		centeredInfoLink("Documentation", &url.URL{Scheme: "https", Host: "shruggietech.github.io", Path: "/go-schedule/"}),
	)
	return wrap(container.NewPadded(content))
}

func centeredInfoLink(text string, destination *url.URL) fyne.CanvasObject {
	link := widget.NewHyperlink(text, destination)
	link.Alignment = fyne.TextAlignCenter
	return link
}
