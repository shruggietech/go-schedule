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

// workAreaPx returns the work area and effective display scale for the monitor
// nearest the launch pointer. The dimensions and scale come from the same
// HMONITOR so mixed-DPI displays cannot be combined accidentally.
func workAreaPx() (int, int, float32) {
	user32 := syscall.NewLazyDLL("user32.dll")
	shcore := syscall.NewLazyDLL("shcore.dll")
	getCursorPos := user32.NewProc("GetCursorPos")
	monitorFromPoint := user32.NewProc("MonitorFromPoint")
	getMonitorInfo := user32.NewProc("GetMonitorInfoW")
	getDPIForMonitor := shcore.NewProc("GetDpiForMonitor")
	dpiLookupAvailable := getDPIForMonitor.Find() == nil

	var point winPoint
	ret, _, _ := getCursorPos.Call(uintptr(unsafe.Pointer(&point)))
	if ret == 0 {
		return 0, 0, 0
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
		func(monitor uintptr) (uint32, bool) {
			if !dpiLookupAvailable {
				return 0, false
			}
			const monitorDPITypeEffective = 0
			var dpiX, dpiY uint32
			result, _, _ := getDPIForMonitor.Call(
				monitor,
				monitorDPITypeEffective,
				uintptr(unsafe.Pointer(&dpiX)),
				uintptr(unsafe.Pointer(&dpiY)),
			)
			return dpiX, result == 0 && dpiX > 0
		},
	)
}

func workAreaForPoint(
	point winPoint,
	monitorForPoint func(winPoint) uintptr,
	infoForMonitor func(uintptr) (winRect, bool),
	dpiForMonitor func(uintptr) (uint32, bool),
) (int, int, float32) {
	monitor := monitorForPoint(point)
	if monitor == 0 {
		return 0, 0, 0
	}
	work, ok := infoForMonitor(monitor)
	if !ok || work.right <= work.left || work.bottom <= work.top {
		return 0, 0, 0
	}
	dpi, dpiOK := dpiForMonitor(monitor)
	if !dpiOK {
		return int(work.right - work.left), int(work.bottom - work.top), 0
	}
	const defaultWindowsDPI = 96
	return int(work.right - work.left), int(work.bottom - work.top), float32(dpi) / defaultWindowsDPI
}
