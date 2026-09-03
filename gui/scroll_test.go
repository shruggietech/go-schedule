package gui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func TestScaleScrollDeltaPreservesPrecisionAndScalesDiscreteWheel(t *testing.T) {
	for _, tc := range []struct {
		name        string
		delta       fyne.Delta
		sensitivity float64
		want        fyne.Delta
	}{
		{name: "precision positive", delta: fyne.NewDelta(2, 12.5), sensitivity: 4, want: fyne.NewDelta(2, 12.5)},
		{name: "precision negative", delta: fyne.NewDelta(-3, -24.9), sensitivity: 4, want: fyne.NewDelta(-3, -24.9)},
		{name: "positive detent", delta: fyne.NewDelta(3, 25), sensitivity: 2, want: fyne.NewDelta(3, 50)},
		{name: "negative detent", delta: fyne.NewDelta(-4, -25), sensitivity: 3, want: fyne.NewDelta(-4, -75)},
		{name: "maximum", delta: fyne.NewDelta(0, -100), sensitivity: 4, want: fyne.NewDelta(0, -400)},
		{name: "invalid preference", delta: fyne.NewDelta(0, -25), sensitivity: 99, want: fyne.NewDelta(0, -50)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := scaleScrollDelta(tc.delta, tc.sensitivity); got != tc.want {
				t.Fatalf("scaleScrollDelta(%+v, %v) = %+v, want %+v", tc.delta, tc.sensitivity, got, tc.want)
			}
		})
	}
}

func TestSensitiveVScrollDelegatesExactlyOnce(t *testing.T) {
	content := canvas.NewRectangle(nil)
	content.SetMinSize(fyne.NewSize(100, 800))
	sensitivity := 2.0
	scroll := newSensitiveVScroll(content, func() float64 { return sensitivity })
	scroll.Resize(fyne.NewSize(100, 100))

	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -25)})
	if got := scroll.Offset.Y; got != 50 {
		t.Fatalf("discrete wheel offset = %v, want 50", got)
	}

	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -10)})
	if got := scroll.Offset.Y; got != 60 {
		t.Fatalf("precision offset = %v, want 60", got)
	}

	sensitivity = 4
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -25)})
	if got := scroll.Offset.Y; got != 160 {
		t.Fatalf("live sensitivity offset = %v, want 160", got)
	}
}

func TestAppOwnedVerticalScrollUsesLivePreference(t *testing.T) {
	ui := &App{appearance: defaultAppearancePreferences()}
	content := canvas.NewRectangle(nil)
	content.SetMinSize(fyne.NewSize(100, 800))
	scroll := ui.newVScroll(content)
	scroll.Resize(fyne.NewSize(100, 100))

	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -25)})
	if got := scroll.Offset.Y; got != 50 {
		t.Fatalf("default offset = %v, want 50", got)
	}
	ui.appearance.ScrollSensitivity = 3
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -25)})
	if got := scroll.Offset.Y; got != 125 {
		t.Fatalf("updated offset = %v, want 125", got)
	}
}
