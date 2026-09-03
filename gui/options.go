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
	root           fyne.CanvasObject
	mode           *widget.Select
	font           *widget.Select
	reset          *cursorButton
	storageContent *fyne.Container
	storageRows    []*storageRowView
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

	storageIntro := widget.NewLabel("These locations are read-only. Daemon paths come from the connected daemon; paths outside application-owned defaults are not claimed as wipe targets.")
	storageIntro.Wrapping = fyne.TextWrapWord
	view.storageContent = container.NewVBox(storageIntro)
	view.setStorageLocations(a.storageLocations, a.clipboard)
	storage := widget.NewCard("Application storage", "Scope and removal behavior for known local application paths.", view.storageContent)

	content := container.NewVBox(
		widget.NewLabelWithStyle("Options", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		appearance,
		storage,
	)
	view.root = container.NewVScroll(container.NewPadded(content))
	a.options = view
	return view.root
}

func (v *optionsView) setStorageLocations(locations []storageLocation, clipboard fyne.Clipboard) {
	v.storageRows = nil
	objects := v.storageContent.Objects[:1]
	for _, location := range locations {
		row := newStorageRowView(location, clipboard)
		v.storageRows = append(v.storageRows, row)
		objects = append(objects, row.root)
	}
	v.storageContent.Objects = objects
	v.storageContent.Refresh()
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
