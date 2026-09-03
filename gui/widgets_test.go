package gui

import (
	"testing"

	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
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

func TestTaskRowDoubleActivationUsesBoundIdentity(t *testing.T) {
	var selected []string
	var activated []string
	row := newTaskRow(
		func(id string) { selected = append(selected, id) },
		func(id string) { activated = append(activated, id) },
	)
	row.bind("task-a", "Task A")
	if row.Text != "Task A" {
		t.Fatalf("row text = %q", row.Text)
	}
	test.Tap(row)
	if len(selected) != 1 || selected[0] != "task-a" || len(activated) != 0 {
		t.Fatalf("single tap selected=%v activated=%v", selected, activated)
	}
	test.DoubleTap(row)
	row.bind("task-b", "Task B")
	test.DoubleTap(row)
	if len(activated) != 2 || activated[0] != "task-a" || activated[1] != "task-b" {
		t.Fatalf("activated IDs = %v", activated)
	}
}

func TestTaskRowIgnoresUnboundIdentity(t *testing.T) {
	selected, activated := false, false
	row := newTaskRow(func(string) { selected = true }, func(string) { activated = true })
	test.Tap(row)
	test.DoubleTap(row)
	if selected || activated {
		t.Fatal("unbound row selected or activated a task")
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
