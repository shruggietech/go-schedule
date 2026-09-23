//go:build !windows

package desktopcontrol

import "errors"

var errElevationCancelled = errors.New("elevation cancelled")

func requestElevation(string) error {
	return errors.New("windows service controls are unavailable on this platform")
}
