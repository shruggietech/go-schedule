package gui

import "fyne.io/fyne/v2"

// defaultWindowSize is used when the screen work area can't be determined.
var defaultWindowSize = fyne.NewSize(1280, 800)

// windowSizeFor converts a screen work area (in physical pixels, excluding the
// taskbar) to Fyne units for the main window, leaving a small margin for the
// window chrome so the frame stays within the work area. When the work area is
// unknown (workW/workH <= 0, e.g. non-Windows) it returns a generous default.
func windowSizeFor(workW, workH int, scale float32) fyne.Size {
	if workW <= 0 || workH <= 0 || scale <= 0 {
		return defaultWindowSize
	}
	w := float32(workW)/scale - 2
	h := float32(workH)/scale - 32 // leave room for the title bar
	if w < 640 {
		w = 640
	}
	if h < 480 {
		h = 480
	}
	return fyne.NewSize(w, h)
}
