package gui

import "fyne.io/fyne/v2"

// defaultWindowSize is used when the screen work area can't be determined.
var defaultWindowSize = fyne.NewSize(1280, 800)

// windowSizeFor converts a monitor work area (in physical pixels, excluding the
// taskbar) to Fyne units for the main window. A fresh launch prefers 1280x800
// but is capped independently to 90 percent of each available logical
// dimension. When the work area is unknown it returns the preferred default.
func windowSizeFor(workW, workH int, scale float32) fyne.Size {
	if workW <= 0 || workH <= 0 || scale <= 0 {
		return defaultWindowSize
	}
	const workAreaFraction = float32(0.9)
	w := float32(workW) / scale * workAreaFraction
	h := float32(workH) / scale * workAreaFraction
	if w > defaultWindowSize.Width {
		w = defaultWindowSize.Width
	}
	if h > defaultWindowSize.Height {
		h = defaultWindowSize.Height
	}
	return fyne.NewSize(w, h)
}
