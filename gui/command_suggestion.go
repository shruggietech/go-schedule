package gui

import (
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type suggestionEntry struct {
	widget.Entry
	suggestion string
}

func (e *suggestionEntry) AcceptsTab() bool {
	if e.Text != "" || e.suggestion == "" {
		return false
	}
	if app := fyne.CurrentApp(); app != nil {
		if driver, ok := app.Driver().(desktop.Driver); ok && driver.CurrentKeyModifiers()&fyne.KeyModifierShift != 0 {
			return false
		}
	}
	return true
}

func newSuggestionEntry(suggestion string) *suggestionEntry {
	entry := &suggestionEntry{suggestion: suggestion}
	entry.ExtendBaseWidget(entry)
	entry.MultiLine = true
	return entry
}

func (e *suggestionEntry) TypedKey(event *fyne.KeyEvent) {
	if event.Name == fyne.KeyTab && e.Text == "" && e.suggestion != "" {
		e.SetText(e.suggestion)
		e.CursorRow = 0
		e.CursorColumn = utf8.RuneCountInString(e.suggestion)
		e.Refresh()
		return
	}
	e.Entry.TypedKey(event)
}
