package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// cursorButton is a widget.Button that shows the pointer (hand) cursor on hover.
// Fyne's stock button keeps the default arrow cursor, giving no visual hint that
// it is clickable; implementing desktop.Cursorable opts into the link-style
// pointer so buttons read as interactive (FR-021).
type cursorButton struct {
	widget.Button
}

// newCursorButton builds a pointer-cursor button with the given label and tap
// handler. importance controls the visual emphasis (e.g. widget.HighImportance
// for the primary Save action).
func newCursorButton(label string, icon fyne.Resource, importance widget.Importance, tapped func()) *cursorButton {
	b := &cursorButton{}
	b.ExtendBaseWidget(b)
	b.Text = label
	b.Icon = icon
	b.Importance = importance
	b.OnTapped = tapped
	return b
}

// Cursor implements desktop.Cursorable.
func (b *cursorButton) Cursor() desktop.Cursor { return desktop.PointerCursor }

// newToolbarButton builds a pointer-cursor button for app toolbars/dialogs with
// the default (medium) emphasis, matching widget.NewButtonWithIcon but with the
// hand cursor on hover (FR-013).
func newToolbarButton(label string, icon fyne.Resource, tapped func()) *cursorButton {
	return newCursorButton(label, icon, widget.MediumImportance, tapped)
}

// newToolbarButtonPlain is newToolbarButton without an icon.
func newToolbarButtonPlain(label string, tapped func()) *cursorButton {
	return newCursorButton(label, nil, widget.MediumImportance, tapped)
}

// collapsible is a disclosure section whose header arrow points right (▶) when
// collapsed and down (▼) when expanded — the standard convention. Fyne's
// widget.Accordion hardcodes the opposite icons, so this small widget is used
// instead (FR-009). Its header is a cursorButton, so it also satisfies the
// app-wide pointer-cursor rule (FR-013).
type collapsible struct {
	widget.BaseWidget
	header  *cursorButton
	content fyne.CanvasObject
	box     *fyne.Container
	open    bool
}

// newCollapsible builds a collapsed disclosure section with the given title and
// content.
func newCollapsible(title string, content fyne.CanvasObject) *collapsible {
	c := &collapsible{content: content}
	c.ExtendBaseWidget(c)
	c.header = newCursorButton(title, theme.NavigateNextIcon(), widget.LowImportance, c.toggle)
	c.header.Alignment = widget.ButtonAlignLeading
	c.header.IconPlacement = widget.ButtonIconLeadingText
	content.Hide()
	c.box = container.NewVBox(c.header, content)
	return c
}

func (c *collapsible) toggle() { c.SetOpen(!c.open) }

// SetOpen expands or collapses the section, updating the arrow and content
// visibility.
func (c *collapsible) SetOpen(open bool) {
	c.open = open
	if open {
		c.header.SetIcon(theme.MenuDropDownIcon())
		c.content.Show()
	} else {
		c.header.SetIcon(theme.NavigateNextIcon())
		c.content.Hide()
	}
	c.Refresh()
}

// CreateRenderer implements fyne.Widget.
func (c *collapsible) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.box)
}
