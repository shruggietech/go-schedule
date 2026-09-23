//go:build !windows && !linux

package desktopcontrol

import "context"

// Instance is unused on platforms without the Windows tray.
type Instance struct{}

func ClaimGUI() (*Instance, bool, error)                { return &Instance{}, true, nil }
func ClaimTray() (*Instance, bool, error)               { return &Instance{}, true, nil }
func (i *Instance) Close() error                        { return nil }
func SignalGUI() (bool, error)                          { return false, nil }
func ListenGUI(context.Context, func()) (func(), error) { return func() {}, nil }
