package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type storageRowView struct {
	location storageLocation
	path     *widget.Label
	copy     *cursorButton
	root     fyne.CanvasObject
}

type optionsView struct {
	root        fyne.CanvasObject
	mode        *widget.Select
	font        *widget.Select
	reset       *cursorButton
	storageRows []*storageRowView
}

func (a *App) buildOptionsTab() fyne.CanvasObject {
	view := &optionsView{}
	view.mode = widget.NewSelect(appearanceModeLabels(), nil)
	view.mode.SetSelected(a.appearance.Mode.label())
	view.font = widget.NewSelect(fontChoiceLabels(), nil)
	view.font.SetSelected(a.appearance.Font.label())
	view.mode.OnChanged = func(label string) {
		next := a.appearance
		next.Mode = appearanceModeForLabel(label)
		a.applyAppearance(next)
	}
	view.font.OnChanged = func(label string) {
		next := a.appearance
		next.Font = fontChoiceForLabel(label)
		a.applyAppearance(next)
	}
	view.reset = newToolbarButtonPlain("Restore defaults", func() {
		a.applyAppearance(defaultAppearancePreferences())
		view.mode.SetSelected(appearanceDark.label())
		view.font.SetSelected(fontBrand.label())
	})

	appearance := widget.NewCard("Appearance", "Changes apply immediately and are saved for this user.", container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Color mode", view.mode),
			widget.NewFormItem("Interface font", view.font),
		),
		container.NewHBox(view.reset),
	))

	storageItems := []fyne.CanvasObject{
		widget.NewLabel("These locations are read-only. Paths outside the application-owned defaults are not listed or claimed as wipe targets."),
	}
	storageItems[0].(*widget.Label).Wrapping = fyne.TextWrapWord
	for _, location := range a.storageLocations {
		row := newStorageRowView(location, a.clipboard)
		view.storageRows = append(view.storageRows, row)
		storageItems = append(storageItems, row.root)
	}
	storage := widget.NewCard("Application storage", "Scope and removal behavior for known local application paths.", container.NewVBox(storageItems...))

	content := container.NewVBox(
		widget.NewLabelWithStyle("Options", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		appearance,
		storage,
	)
	view.root = container.NewVScroll(container.NewPadded(content))
	a.options = view
	return view.root
}

func newStorageRowView(location storageLocation, clipboard fyne.Clipboard) *storageRowView {
	pathText := location.Path
	if !location.Available {
		pathText = "Unavailable in this environment"
	}
	path := widget.NewLabel(pathText)
	path.Wrapping = fyne.TextWrapBreak
	path.Selectable = location.Available
	copyButton := newToolbarButtonPlain("Copy", func() {
		if location.Available {
			clipboard.SetContent(location.Path)
		}
	})
	if !location.Available {
		copyButton.Disable()
	}
	details := widget.NewLabel(fmt.Sprintf(
		"Scope: %s\nStatus: %s\nSoftware-only removal: %s\nExplicit data wipe: %s",
		location.Scope,
		location.Existence,
		location.SoftwareOnlyRemoval,
		location.ExplicitDataWipe,
	))
	details.Wrapping = fyne.TextWrapWord
	root := widget.NewCard(location.Category, "", container.NewVBox(
		container.NewBorder(nil, nil, nil, copyButton, path),
		details,
	))
	return &storageRowView{location: location, path: path, copy: copyButton, root: root}
}

func (a *App) applyAppearance(value appearancePreferences) {
	value = value.normalized()
	a.appearance = value
	saveAppearancePreferences(a.fyne.Preferences(), value)
	applyBrandTheme(a.fyne.Settings(), value)
}
