package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type navigationID string

const (
	navigationTasks    navigationID = "tasks"
	navigationGroups   navigationID = "groups"
	navigationChains   navigationID = "chains"
	navigationSchedule navigationID = "schedule"
	navigationActivity navigationID = "activity"
	navigationOptions  navigationID = "options"
	navigationInfo     navigationID = "info"
)

const minimumNavigationRailWidth float32 = 168

type navigationDestinationSpec struct {
	ID      navigationID
	Label   string
	Content fyne.CanvasObject
}

type navigationDestination struct {
	id      navigationID
	button  *cursorButton
	content fyne.CanvasObject
}

type navigationShell struct {
	root         fyne.CanvasObject
	rail         *fyne.Container
	content      *fyne.Container
	boundary     *widget.Separator
	exit         *cursorButton
	destinations []*navigationDestination
	selected     navigationID
}

func newNavigationShell(specs []navigationDestinationSpec, onExit func()) *navigationShell {
	shell := &navigationShell{}
	buttons := make([]fyne.CanvasObject, 0, len(specs))
	contents := make([]fyne.CanvasObject, 0, len(specs))
	labels := make([]string, 0, len(specs))
	for _, spec := range specs {
		destination := &navigationDestination{id: spec.ID, content: spec.Content}
		destination.button = newCursorButton(spec.Label, nil, widget.LowImportance, func() {
			shell.selectDestination(destination.id)
		})
		destination.button.Alignment = widget.ButtonAlignLeading
		shell.destinations = append(shell.destinations, destination)
		buttons = append(buttons, destination.button)
		contents = append(contents, destination.content)
		labels = append(labels, spec.Label)
	}

	shell.content = container.NewStack(contents...)
	shell.exit = newCursorButton("Exit", theme.LogoutIcon(), widget.DangerImportance, onExit)
	shell.exit.Alignment = widget.ButtonAlignTrailing
	top := container.NewVBox(buttons...)
	shell.rail = container.New(&navigationRailLayout{minimumWidth: navigationRailMinimumWidth(labels)}, top, widget.NewSeparator(), shell.exit)
	shell.boundary = widget.NewSeparator()
	railWithBoundary := container.NewBorder(nil, nil, shell.rail, shell.boundary)
	shell.root = container.NewBorder(nil, nil, railWithBoundary, nil, shell.content)
	if len(shell.destinations) > 0 {
		shell.selectDestination(shell.destinations[0].id)
	}
	return shell
}

func (s *navigationShell) selectDestination(id navigationID) bool {
	found := false
	for _, destination := range s.destinations {
		selected := destination.id == id
		if selected {
			found = true
			destination.content.Show()
			destination.button.Importance = widget.HighImportance
		} else {
			destination.content.Hide()
			destination.button.Importance = widget.LowImportance
		}
		destination.button.Refresh()
	}
	if found {
		s.selected = id
		s.content.Refresh()
	}
	return found
}

func (s *navigationShell) updateLabel(id navigationID, label string) bool {
	for _, destination := range s.destinations {
		if destination.id == id {
			destination.button.SetText(label)
			return true
		}
	}
	return false
}

func (s *navigationShell) label(id navigationID) string {
	for _, destination := range s.destinations {
		if destination.id == id {
			return destination.button.Text
		}
	}
	return ""
}

func (s *navigationShell) contentFor(id navigationID) fyne.CanvasObject {
	for _, destination := range s.destinations {
		if destination.id == id {
			return destination.content
		}
	}
	return nil
}

func (s *navigationShell) labels() []string {
	labels := make([]string, len(s.destinations))
	for i, destination := range s.destinations {
		labels[i] = destination.button.Text
	}
	return labels
}

func navigationRailMinimumWidth(labels []string) float32 {
	width := minimumNavigationRailWidth
	for _, label := range labels {
		measured := fyne.MeasureText(label, theme.TextSize(), fyne.TextStyle{Bold: true}).Width + theme.Padding()*4
		if measured > width {
			width = measured
		}
	}
	return width
}

// navigationRailLayout keeps destinations at the top and Exit at the
// bottom-right, with a separator and symmetric padding independent of height.
type navigationRailLayout struct{ minimumWidth float32 }

func (l *navigationRailLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	padding := theme.Padding()
	top := objects[0].MinSize()
	separator := objects[1].MinSize()
	exit := objects[2].MinSize()
	width := l.minimumWidth
	if candidate := top.Width + padding*2; candidate > width {
		width = candidate
	}
	if candidate := exit.Width + padding*2; candidate > width {
		width = candidate
	}
	return fyne.NewSize(width, top.Height+separator.Height+exit.Height+padding*4)
}

func (l *navigationRailLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	padding := theme.Padding()
	top := objects[0]
	separator := objects[1]
	exit := objects[2]
	exitSize := exit.MinSize()
	separatorSize := separator.MinSize()
	separatorY := size.Height - exitSize.Height - padding*3 - separatorSize.Height
	if separatorY < padding {
		separatorY = padding
	}
	top.Move(fyne.NewPos(padding, padding))
	top.Resize(fyne.NewSize(size.Width-padding*2, separatorY-padding))
	separator.Move(fyne.NewPos(padding, separatorY))
	separator.Resize(fyne.NewSize(size.Width-padding*2, separatorSize.Height))
	exit.Move(fyne.NewPos(size.Width-padding-exitSize.Width, size.Height-padding-exitSize.Height))
	exit.Resize(exitSize)
}
