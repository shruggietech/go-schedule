//go:build windows

package gui

import (
	"syscall"
	"unsafe"
)

type winPoint struct{ x, y int32 }
type winRect struct{ left, top, right, bottom int32 }
type monitorInfo struct {
	cbSize  uint32
	monitor winRect
	work    winRect
	flags   uint32
}

// workAreaPx returns the monitor work area nearest the launch pointer in
// physical pixels, or 0,0 if it cannot be determined. Selecting by pointer
// avoids silently sizing a launch for a different, larger primary monitor.
func workAreaPx() (int, int) {
	user32 := syscall.NewLazyDLL("user32.dll")
	getCursorPos := user32.NewProc("GetCursorPos")
	monitorFromPoint := user32.NewProc("MonitorFromPoint")
	getMonitorInfo := user32.NewProc("GetMonitorInfoW")

	var point winPoint
	ret, _, _ := getCursorPos.Call(uintptr(unsafe.Pointer(&point)))
	if ret == 0 {
		return 0, 0
	}

	return workAreaForPoint(
		point,
		func(p winPoint) uintptr {
			const monitorDefaultToNearest = 2
			packed := uintptr(uint32(p.x)) | uintptr(uint64(uint32(p.y))<<32)
			monitor, _, _ := monitorFromPoint.Call(packed, monitorDefaultToNearest)
			return monitor
		},
		func(monitor uintptr) (winRect, bool) {
			info := monitorInfo{cbSize: uint32(unsafe.Sizeof(monitorInfo{}))}
			ok, _, _ := getMonitorInfo.Call(monitor, uintptr(unsafe.Pointer(&info)))
			return info.work, ok != 0
		},
	)
}

func workAreaForPoint(
	point winPoint,
	monitorForPoint func(winPoint) uintptr,
	infoForMonitor func(uintptr) (winRect, bool),
) (int, int) {
	monitor := monitorForPoint(point)
	if monitor == 0 {
		return 0, 0
	}
	work, ok := infoForMonitor(monitor)
	if !ok || work.right <= work.left || work.bottom <= work.top {
		return 0, 0
	}
	return int(work.right - work.left), int(work.bottom - work.top)
}
