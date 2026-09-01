//go:build windows

package gui

import "testing"

func TestWorkAreaForPointUsesSelectedMonitor(t *testing.T) {
	point := winPoint{x: 2400, y: 300}
	wantMonitor := uintptr(22)
	monitorCalls := 0
	infoCalls := 0

	gotW, gotH := workAreaForPoint(
		point,
		func(got winPoint) uintptr {
			monitorCalls++
			if got != point {
				t.Fatalf("monitor lookup point = %#v, want %#v", got, point)
			}
			return wantMonitor
		},
		func(got uintptr) (winRect, bool) {
			infoCalls++
			if got != wantMonitor {
				t.Fatalf("monitor info handle = %d, want %d", got, wantMonitor)
			}
			return winRect{left: 1920, top: 0, right: 2720, bottom: 560}, true
		},
	)

	if gotW != 800 || gotH != 560 {
		t.Fatalf("selected monitor work area = %dx%d, want 800x560", gotW, gotH)
	}
	if monitorCalls != 1 || infoCalls != 1 {
		t.Fatalf("lookup calls = monitor %d, info %d; want 1 each", monitorCalls, infoCalls)
	}
}

func TestWorkAreaForPointRejectsLookupFailure(t *testing.T) {
	if gotW, gotH := workAreaForPoint(winPoint{}, func(winPoint) uintptr { return 0 }, func(uintptr) (winRect, bool) {
		t.Fatal("monitor info must not run without a monitor")
		return winRect{}, false
	}); gotW != 0 || gotH != 0 {
		t.Fatalf("failed lookup work area = %dx%d, want 0x0", gotW, gotH)
	}
}
