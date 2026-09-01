package gui

import (
	"math"
	"testing"

	"fyne.io/fyne/v2"
)

func TestWindowSizeFor(t *testing.T) {
	tests := []struct {
		name         string
		workW, workH int
		scale        float32
		want         fyne.Size
	}{
		{name: "unknown work area", want: fyne.NewSize(1280, 800)},
		{name: "800x600", workW: 800, workH: 600, scale: 1, want: fyne.NewSize(720, 540)},
		{name: "1024x768", workW: 1024, workH: 768, scale: 1, want: fyne.NewSize(921.6, 691.2)},
		{name: "1366x768", workW: 1366, workH: 768, scale: 1, want: fyne.NewSize(1229.4, 691.2)},
		{name: "1920x1080", workW: 1920, workH: 1080, scale: 1, want: fyne.NewSize(1280, 800)},
		{name: "200 percent scaling", workW: 3840, workH: 2160, scale: 2, want: fyne.NewSize(1280, 800)},
		{name: "smaller than old minimum", workW: 400, workH: 300, scale: 1, want: fyne.NewSize(360, 270)},
		{name: "invalid scale", workW: 1920, workH: 1080, scale: 0, want: fyne.NewSize(1280, 800)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windowSizeFor(tt.workW, tt.workH, tt.scale)
			if !sizesClose(got, tt.want) {
				t.Fatalf("windowSizeFor(%d, %d, %g) = %v, want %v", tt.workW, tt.workH, tt.scale, got, tt.want)
			}
		})
	}
}

func sizesClose(got, want fyne.Size) bool {
	const tolerance = 0.01
	return math.Abs(float64(got.Width-want.Width)) < tolerance &&
		math.Abs(float64(got.Height-want.Height)) < tolerance
}
