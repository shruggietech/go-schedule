//go:build windows

package gui

import "testing"

func TestWorkAreaForPointUsesSelectedMonitor(t *testing.T) {
	point := winPoint{x: 2400, y: 300}
	wantMonitor := uintptr(22)
	monitorCalls := 0
	infoCalls := 0
	dpiCalls := 0

	gotW, gotH, gotScale := workAreaForPoint(
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
		func(got uintptr) (uint32, bool) {
			dpiCalls++
			if got != wantMonitor {
				t.Fatalf("monitor DPI handle = %d, want %d", got, wantMonitor)
			}
			return 192, true
		},
	)

	if gotW != 800 || gotH != 560 {
		t.Fatalf("selected monitor work area = %dx%d, want 800x560", gotW, gotH)
	}
	if gotScale != 2 {
		t.Fatalf("selected monitor scale = %v, want 2", gotScale)
	}
	if monitorCalls != 1 || infoCalls != 1 || dpiCalls != 1 {
		t.Fatalf("lookup calls = monitor %d, info %d, DPI %d; want 1 each", monitorCalls, infoCalls, dpiCalls)
	}
}

func TestWorkAreaForPointRejectsLookupFailure(t *testing.T) {
	if gotW, gotH, gotScale := workAreaForPoint(winPoint{}, func(winPoint) uintptr { return 0 }, func(uintptr) (winRect, bool) {
		t.Fatal("monitor info must not run without a monitor")
		return winRect{}, false
	}, func(uintptr) (uint32, bool) {
		t.Fatal("monitor DPI must not run without a monitor")
		return 0, false
	}); gotW != 0 || gotH != 0 || gotScale != 0 {
		t.Fatalf("failed lookup work area = %dx%d at scale %v, want 0x0 at scale 0", gotW, gotH, gotScale)
	}
}

func TestWorkAreaForPointRetainsAreaWhenDPILookupFails(t *testing.T) {
	gotW, gotH, gotScale := workAreaForPoint(
		winPoint{},
		func(winPoint) uintptr { return 7 },
		func(uintptr) (winRect, bool) {
			return winRect{left: 10, top: 20, right: 810, bottom: 580}, true
		},
		func(uintptr) (uint32, bool) { return 0, false },
	)

	if gotW != 800 || gotH != 560 || gotScale != 0 {
		t.Fatalf("work area with unavailable DPI = %dx%d at scale %v, want 800x560 at scale 0", gotW, gotH, gotScale)
	}
}
