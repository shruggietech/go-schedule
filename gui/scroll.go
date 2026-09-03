package gui

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

const discreteWheelDelta = float32(25)

// sensitiveVScroll preserves Fyne's scroll layout, scrollbar, drag, keyboard,
// and offset behavior while applying the user preference to discrete wheel
// events only. The sensitivity function is read per event so Options changes
// affect already-open views.
type sensitiveVScroll struct {
	container.Scroll
	sensitivity func() float64
}

var _ fyne.Scrollable = (*sensitiveVScroll)(nil)

func newSensitiveVScroll(content fyne.CanvasObject, sensitivity func() float64) *sensitiveVScroll {
	scroll := &sensitiveVScroll{sensitivity: sensitivity}
	scroll.Direction = container.ScrollVerticalOnly
	scroll.Content = content
	scroll.ExtendBaseWidget(scroll)
	return scroll
}

func (s *sensitiveVScroll) Scrolled(event *fyne.ScrollEvent) {
	if event == nil {
		return
	}
	sensitivity := defaultScrollSensitivity
	if s.sensitivity != nil {
		sensitivity = s.sensitivity()
	}
	adjusted := *event
	adjusted.Scrolled = scaleScrollDelta(event.Scrolled, sensitivity)
	s.Scroll.Scrolled(&adjusted)
}

func scaleScrollDelta(delta fyne.Delta, sensitivity float64) fyne.Delta {
	result := delta
	if math.Abs(float64(delta.DY)) >= float64(discreteWheelDelta) {
		result.DY *= float32(normalizeScrollSensitivity(sensitivity))
	}
	return result
}

func (a *App) newVScroll(content fyne.CanvasObject) *sensitiveVScroll {
	return newSensitiveVScroll(content, func() float64 {
		return a.appearance.ScrollSensitivity
	})
}
