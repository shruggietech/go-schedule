//go:build !windows && !linux

package desktopcontrol

import (
	"context"
	"errors"
)

var errElevationCancelled = errors.New("elevation cancelled")

func requestElevation(context.Context, string) error {
	return errors.New("windows service controls are unavailable on this platform")
}
