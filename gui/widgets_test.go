package gui

import (
	"testing"

	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestCursorButton_PointerCursorAndTap(t *testing.T) {
	tapped := false
	b := newCursorButton("Go", nil, widget.HighImportance, func() { tapped = true })

	if b.Cursor() != desktop.PointerCursor {
		t.Fatalf("Cursor() = %v, want PointerCursor", b.Cursor())
	}
	b.OnTapped()
	if !tapped {
		t.Fatal("tap handler not invoked")
	}
}

func TestToolbarButton_PointerCursor(t *testing.T) {
	if b := newToolbarButton("New", theme.ContentAddIcon(), func() {}); b.Cursor() != desktop.PointerCursor {
		t.Fatal("newToolbarButton should report the pointer cursor")
	}
	if b := newToolbarButtonPlain("Toggle", func() {}); b.Cursor() != desktop.PointerCursor {
		t.Fatal("newToolbarButtonPlain should report the pointer cursor")
	}
}

func TestCollapsible_StateAndIcon(t *testing.T) {
	content := widget.NewLabel("inner")
	c := newCollapsible("Advanced Settings", content)

	if c.open {
		t.Fatal("collapsible should start collapsed")
	}
	if content.Visible() {
		t.Fatal("content should be hidden when collapsed")
	}
	if c.header.Icon.Name() != theme.NavigateNextIcon().Name() {
		t.Fatalf("collapsed icon = %s, want NavigateNext (▶)", c.header.Icon.Name())
	}
	if c.header.Cursor() != desktop.PointerCursor {
		t.Fatal("collapsible header should report the pointer cursor")
	}

	c.toggle()
	if !c.open || !content.Visible() {
		t.Fatal("toggle should expand and show content")
	}
	if c.header.Icon.Name() != theme.MenuDropDownIcon().Name() {
		t.Fatalf("expanded icon = %s, want MenuDropDown (▼)", c.header.Icon.Name())
	}
}
