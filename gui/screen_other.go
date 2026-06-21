//go:build !windows

package gui

// workAreaPx returns 0,0 on non-Windows platforms; the caller then falls back to
// a generous default window size.
func workAreaPx() (int, int) { return 0, 0 }
