//go:build windows

package gui

import (
	"syscall"
	"unsafe"
)

type winRect struct{ left, top, right, bottom int32 }

// workAreaPx returns the primary monitor's work area (the screen minus the
// taskbar) in physical pixels, or 0,0 if it can't be determined. It uses the
// standard-library syscall bridge to user32 — no third-party dependency.
func workAreaPx() (int, int) {
	user32 := syscall.NewLazyDLL("user32.dll")
	spi := user32.NewProc("SystemParametersInfoW")
	const spiGetWorkArea = 0x0030
	var r winRect
	ret, _, _ := spi.Call(spiGetWorkArea, 0, uintptr(unsafe.Pointer(&r)), 0)
	if ret == 0 {
		return 0, 0
	}
	return int(r.right - r.left), int(r.bottom - r.top)
}
