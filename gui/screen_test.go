package gui

import "testing"

func TestWindowSizeFor(t *testing.T) {
	// Unknown work area → generous default.
	if got := windowSizeFor(0, 0, 1); got != defaultWindowSize {
		t.Fatalf("unknown work area = %v, want default %v", got, defaultWindowSize)
	}

	// Scale 1: work area minus the chrome margins.
	got := windowSizeFor(1920, 1040, 1)
	if got.Width != 1918 || got.Height != 1008 {
		t.Fatalf("size at scale 1 = %v, want {1918, 1008}", got)
	}

	// HiDPI: pixels are divided by the scale before the margins.
	got = windowSizeFor(3840, 2080, 2)
	if got.Width != 1918 || got.Height != 1008 {
		t.Fatalf("size at scale 2 = %v, want {1918, 1008}", got)
	}

	// Tiny work area clamps to the minimum.
	got = windowSizeFor(400, 300, 1)
	if got.Width != 640 || got.Height != 480 {
		t.Fatalf("tiny work area = %v, want clamp {640, 480}", got)
	}
}
