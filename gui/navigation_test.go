package gui

import (
	"reflect"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestNavigationOrderSelectionAndLabelUpdate(t *testing.T) {
	exited := 0
	shell := newNavigationShell([]navigationDestinationSpec{
		{ID: navigationTasks, Label: "Tasks", Content: canvas.NewRectangle(nil)},
		{ID: navigationGroups, Label: "Groups", Content: canvas.NewRectangle(nil)},
		{ID: navigationChains, Label: "Chains", Content: canvas.NewRectangle(nil)},
		{ID: navigationSchedule, Label: "Schedule", Content: canvas.NewRectangle(nil)},
		{ID: navigationActivity, Label: "Activity", Content: canvas.NewRectangle(nil)},
		{ID: navigationOptions, Label: "Options", Content: canvas.NewRectangle(nil)},
		{ID: navigationInfo, Label: "Info", Content: canvas.NewRectangle(nil)},
	}, func() { exited++ })

	want := []string{"Tasks", "Groups", "Chains", "Schedule", "Activity", "Options", "Info"}
	if got := shell.labels(); !reflect.DeepEqual(got, want) {
		t.Fatalf("navigation labels = %v, want %v", got, want)
	}
	if shell.selected != navigationTasks || shell.destinations[0].button.Importance != widget.HighImportance {
		t.Fatalf("initial selection = %q", shell.selected)
	}
	shell.destinations[5].button.OnTapped()
	if shell.selected != navigationOptions || shell.destinations[5].button.Importance != widget.HighImportance {
		t.Fatalf("selected = %q, want Options", shell.selected)
	}
	if shell.destinations[0].content.Visible() || !shell.destinations[5].content.Visible() {
		t.Fatal("selection did not swap visible content")
	}
	if !shell.updateLabel(navigationActivity, "Activity (99+)") {
		t.Fatal("Activity label update failed")
	}
	if shell.destinations[4].button.Text != "Activity (99+)" {
		t.Fatalf("Activity label = %q", shell.destinations[4].button.Text)
	}
	if shell.exit.Importance != widget.DangerImportance || shell.exit.Icon != theme.LogoutIcon() {
		t.Fatalf("Exit treatment = importance %v icon %v", shell.exit.Importance, shell.exit.Icon)
	}
	shell.exit.OnTapped()
	if exited != 1 || shell.selected != navigationOptions {
		t.Fatalf("Exit changed selection or did not execute: selected=%q count=%d", shell.selected, exited)
	}
}

func TestNavigationSeparatesDefinitionsFromOperations(t *testing.T) {
	shell := newNavigationShell([]navigationDestinationSpec{
		{ID: navigationTasks, Label: "Tasks", Content: canvas.NewRectangle(nil), Section: navigationDefinitions},
		{ID: navigationChains, Label: "Chains", Content: canvas.NewRectangle(nil), Section: navigationDefinitions},
		{ID: navigationSchedule, Label: "Schedule", Content: canvas.NewRectangle(nil), Section: navigationOperations},
	}, func() {})
	if shell.sectionBoundary == nil || !shell.sectionBoundary.Visible() {
		t.Fatal("semantic navigation section boundary is missing")
	}
	if got := shell.labels(); !reflect.DeepEqual(got, []string{"Tasks", "Chains", "Schedule"}) {
		t.Fatalf("navigation labels = %v", got)
	}
}

func TestNavigationRailHasFullHeightContentBoundary(t *testing.T) {
	shell := newNavigationShell([]navigationDestinationSpec{
		{ID: navigationTasks, Label: "Tasks", Content: canvas.NewRectangle(nil)},
		{ID: navigationInfo, Label: "Info", Content: canvas.NewRectangle(nil)},
	}, func() {})
	if shell.boundary == nil {
		t.Fatal("navigation shell has no rail/content boundary")
	}
	shell.root.Resize(fyne.NewSize(800, 600))
	shell.root.Refresh()
	if got := shell.boundary.Size().Height; got < 600-theme.Padding()*2 {
		t.Fatalf("boundary height = %v, want full shell height", got)
	}
	if got := shell.boundary.Position().X; got < shell.rail.MinSize().Width-theme.Padding() {
		t.Fatalf("boundary x = %v, rail width %v", got, shell.rail.MinSize().Width)
	}
}

func TestNavigationRailGeometryAtSupportedSizes(t *testing.T) {
	shell := newNavigationShell([]navigationDestinationSpec{
		{ID: navigationTasks, Label: "Tasks", Content: canvas.NewRectangle(nil)},
		{ID: navigationActivity, Label: "Activity (99+)", Content: canvas.NewRectangle(nil)},
		{ID: navigationOptions, Label: "Options", Content: canvas.NewRectangle(nil)},
		{ID: navigationInfo, Label: "Info", Content: canvas.NewRectangle(nil)},
	}, func() {})

	wantMin := navigationRailMinimumWidth([]string{"Tasks", "Activity (99+)", "Options", "Info"})
	if got := shell.rail.MinSize().Width; got < wantMin || got < minimumNavigationRailWidth {
		t.Fatalf("rail minimum width = %v, want >= %v and >= %v", got, wantMin, minimumNavigationRailWidth)
	}
	for _, size := range []fyne.Size{fyne.NewSize(1280, 800), fyne.NewSize(800, 600)} {
		shell.rail.Resize(fyne.NewSize(shell.rail.MinSize().Width, size.Height))
		shell.rail.Refresh()
		padding := theme.Padding()
		if got := shell.exit.Position().Y + shell.exit.Size().Height; got > size.Height-padding/2 || got < size.Height-padding*2 {
			t.Errorf("size %v Exit bottom = %v, want near %v", size, got, size.Height-padding)
		}
		if got := shell.exit.Position().X + shell.exit.Size().Width; got < shell.rail.Size().Width-padding*2 {
			t.Errorf("size %v Exit right = %v, rail width %v", size, got, shell.rail.Size().Width)
		}
	}
}

func TestNavigationRailWidthAccountsForLabelAndPadding(t *testing.T) {
	short := navigationRailMinimumWidth([]string{"Info"})
	long := navigationRailMinimumWidth([]string{"An intentionally long navigation destination"})
	if short != minimumNavigationRailWidth {
		t.Fatalf("short width = %v, want floor %v", short, minimumNavigationRailWidth)
	}
	if long <= short {
		t.Fatalf("long width = %v, want > %v", long, short)
	}
}
