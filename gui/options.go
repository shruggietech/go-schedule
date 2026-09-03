package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const minimumStorageCategoryWidth float32 = 148

type storageRowView struct {
	location storageLocation
	category *widget.Label
	path     *widget.Label
	details  *widget.Label
	copy     *cursorButton
	root     fyne.CanvasObject
	muted    bool
}

type optionsView struct {
	root             fyne.CanvasObject
	mode             *widget.Select
	font             *widget.Select
	sensitivity      *widget.Slider
	sensitivityValue *widget.Label
	reset            *cursorButton
	storageContent   *fyne.Container
	storageHeader    fyne.CanvasObject
	storageRows      []*storageRowView
}

func (a *App) buildOptionsTab() fyne.CanvasObject {
	view := &optionsView{}
	view.mode = newAlternativeSelect(appearanceModeLabels(), a.appearance.Mode.label(), func(label string) {
		next := a.appearance
		next.Mode = appearanceModeForLabel(label)
		a.applyAppearance(next)
	})
	view.font = newAlternativeSelect(fontChoiceLabels(), a.appearance.Font.label(), func(label string) {
		next := a.appearance
		next.Font = fontChoiceForLabel(label)
		a.applyAppearance(next)
	})
	view.sensitivityValue = widget.NewLabel(scrollSensitivityLabel(a.appearance.ScrollSensitivity))
	view.sensitivityValue.Alignment = fyne.TextAlignTrailing
	view.sensitivity = widget.NewSlider(minimumScrollSensitivity, maximumScrollSensitivity)
	view.sensitivity.Step = scrollSensitivityStep
	view.sensitivity.SetValue(a.appearance.ScrollSensitivity)
	view.sensitivity.OnChanged = func(value float64) {
		next := a.appearance
		next.ScrollSensitivity = value
		a.applyAppearance(next)
		view.sensitivityValue.SetText(scrollSensitivityLabel(a.appearance.ScrollSensitivity))
	}
	view.reset = newToolbarButtonPlain("Restore defaults", func() {
		defaults := defaultAppearancePreferences()
		a.applyAppearance(defaults)
		syncAlternativeSelect(view.mode, defaults.Mode.label(), appearanceModeLabels())
		syncAlternativeSelect(view.font, defaults.Font.label(), fontChoiceLabels())
		view.sensitivity.SetValue(defaults.ScrollSensitivity)
		view.sensitivityValue.SetText(scrollSensitivityLabel(defaults.ScrollSensitivity))
	})

	sensitivityControl := container.NewBorder(nil, nil, nil, view.sensitivityValue, view.sensitivity)
	appearance := widget.NewCard("Appearance", "Changes apply immediately and are saved for this user.", container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Color mode", view.mode),
			widget.NewFormItem("Interface font", view.font),
			widget.NewFormItem("Scroll sensitivity", sensitivityControl),
		),
		container.NewHBox(view.reset),
	))

	storageIntro := widget.NewLabel("These locations are read-only. Daemon paths come from the connected daemon; paths outside application-owned defaults are not claimed as wipe targets.")
	storageIntro.Wrapping = fyne.TextWrapWord
	view.storageHeader = newStorageHeader(minimumStorageCategoryWidth)
	view.storageContent = container.NewVBox(storageIntro, view.storageHeader)
	view.setStorageLocations(a.storageLocations, a.clipboard)
	storage := widget.NewCard("Application storage", "Scope and removal behavior for known local application paths.", view.storageContent)

	content := container.NewVBox(
		widget.NewLabelWithStyle("Options", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		appearance,
		storage,
	)
	view.root = a.newVScroll(container.NewPadded(content))
	a.options = view
	return view.root
}

func newAlternativeSelect(all []string, current string, changed func(string)) *widget.Select {
	selection := widget.NewSelect(all, nil)
	selection.SetSelected(current)
	selection.OnChanged = func(value string) {
		if changed != nil {
			changed(value)
		}
		syncAlternativeSelect(selection, value, all)
	}
	syncAlternativeSelect(selection, current, all)
	return selection
}

func syncAlternativeSelect(selection *widget.Select, current string, all []string) {
	selection.Selected = current
	alternatives := make([]string, 0, len(all)-1)
	for _, candidate := range all {
		if candidate != current {
			alternatives = append(alternatives, candidate)
		}
	}
	selection.SetOptions(alternatives)
}

func scrollSensitivityLabel(value float64) string {
	value = normalizeScrollSensitivity(value)
	label := fmt.Sprintf("%gx", value)
	if value == defaultScrollSensitivity {
		return label + " (default)"
	}
	return label
}

func (v *optionsView) setStorageLocations(locations []storageLocation, clipboard fyne.Clipboard) {
	v.storageRows = nil
	categoryWidth := storageCategoryWidth(locations)
	v.storageHeader = newStorageHeader(categoryWidth)
	objects := []fyne.CanvasObject{v.storageContent.Objects[0], v.storageHeader}
	for _, location := range locations {
		row := newStorageRowView(location, clipboard, categoryWidth)
		v.storageRows = append(v.storageRows, row)
		objects = append(objects, row.root)
	}
	v.storageContent.Objects = objects
	v.storageContent.Refresh()
}

func storageCategoryWidth(locations []storageLocation) float32 {
	width := minimumStorageCategoryWidth
	for _, location := range locations {
		candidate := fyne.MeasureText(location.Category, theme.TextSize(), fyne.TextStyle{Bold: true}).Width + theme.Padding()*2
		if candidate > width {
			width = candidate
		}
	}
	return width
}

func newStorageHeader(categoryWidth float32) fyne.CanvasObject {
	category := newStorageText("Category", false, true, false)
	details := newStorageText("Location and removal details", false, true, false)
	action := newStorageText("Action", false, true, false)
	return container.NewVBox(
		container.New(&storageColumnsLayout{categoryWidth: categoryWidth}, category, details, action),
		widget.NewSeparator(),
	)
}

func newStorageRowView(location storageLocation, clipboard fyne.Clipboard, categoryWidths ...float32) *storageRowView {
	categoryWidth := minimumStorageCategoryWidth
	if len(categoryWidths) > 0 {
		categoryWidth = categoryWidths[0]
	}
	muted := !location.Available
	pathText := location.Path
	if muted {
		pathText = "Unavailable in this environment"
	}
	category := newStorageText(location.Category, muted, true, false)
	path := newStorageText(pathText, muted, false, location.Available)
	path.Wrapping = fyne.TextWrapBreak
	details := newStorageText(fmt.Sprintf(
		"Scope: %s  •  Status: %s  •  Software-only removal: %s  •  Explicit data wipe: %s",
		location.Scope,
		location.Existence,
		location.SoftwareOnlyRemoval,
		location.ExplicitDataWipe,
	), muted, false, false)
	details.Wrapping = fyne.TextWrapWord
	copyButton := newToolbarButtonPlain("Copy", func() {
		if location.Available && clipboard != nil {
			clipboard.SetContent(location.Path)
		}
	})
	if muted || clipboard == nil {
		copyButton.Disable()
	}
	middle := container.NewVBox(path, details)
	columns := container.New(&storageColumnsLayout{categoryWidth: categoryWidth}, category, middle, copyButton)
	root := container.NewVBox(columns, widget.NewSeparator())
	return &storageRowView{
		location: location,
		category: category,
		path:     path,
		details:  details,
		copy:     copyButton,
		root:     root,
		muted:    muted,
	}
}

func newStorageText(text string, muted, bold, selectable bool) *widget.Label {
	label := widget.NewLabel(text)
	label.TextStyle.Bold = bold
	label.Selectable = selectable
	if muted {
		label.Importance = widget.LowImportance
	}
	return label
}

type storageColumnsLayout struct{ categoryWidth float32 }

func (l *storageColumnsLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	padding := theme.Padding()
	middle := objects[1].MinSize()
	action := objects[2].MinSize()
	height := maxFloat32(objects[0].MinSize().Height, middle.Height, action.Height) + padding*2
	return fyne.NewSize(l.categoryWidth+middle.Width+action.Width+padding*4, height)
}

func (l *storageColumnsLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	padding := theme.Padding()
	actionSize := objects[2].MinSize()
	categoryWidth := l.categoryWidth
	availableMiddle := size.Width - categoryWidth - actionSize.Width - padding*4
	if availableMiddle < 1 {
		availableMiddle = 1
	}
	contentHeight := size.Height - padding*2
	objects[0].Move(fyne.NewPos(padding, padding))
	objects[0].Resize(fyne.NewSize(categoryWidth, contentHeight))
	objects[1].Move(fyne.NewPos(categoryWidth+padding*2, padding))
	objects[1].Resize(fyne.NewSize(availableMiddle, contentHeight))
	objects[2].Move(fyne.NewPos(size.Width-padding-actionSize.Width, padding))
	objects[2].Resize(fyne.NewSize(actionSize.Width, actionSize.Height))
}

func maxFloat32(values ...float32) float32 {
	max := float32(0)
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max
}

func (a *App) applyAppearance(value appearancePreferences) {
	value = value.normalized()
	a.appearance = value
	saveAppearancePreferences(a.fyne.Preferences(), value)
	applyBrandTheme(a.fyne.Settings(), value)
}
