//go:build !windows

package gui

// workAreaPx returns an unknown work area and scale on non-Windows platforms;
// the caller then falls back to a generous default window size.
func workAreaPx() (int, int, float32) { return 0, 0, 0 }
