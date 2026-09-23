//go:build windows

package desktopcontrol

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sys/windows"
)

const (
	guiMutexName  = "Local\\go-schedule-gui-v1"
	guiEventName  = "Local\\go-schedule-gui-open-v1"
	trayMutexName = "Local\\go-schedule-tray-v1"
)

// Instance owns a session-local mutex until Close.
type Instance struct{ handle windows.Handle }

// ClaimGUI returns false when the current session already has a GUI.
func ClaimGUI() (*Instance, bool, error) { return claim(guiMutexName) }

// ClaimTray returns false when the current session already has a companion.
func ClaimTray() (*Instance, bool, error) { return claim(trayMutexName) }

func claim(name string) (*Instance, bool, error) {
	ptr, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, false, err
	}
	h, err := windows.CreateMutex(nil, false, ptr)
	if errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
		if h != 0 {
			windows.CloseHandle(h) //nolint:errcheck // duplicate handle
		}
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("claim %s: %w", name, err)
	}
	return &Instance{handle: h}, true, nil
}

// Close releases an owned session instance.
func (i *Instance) Close() error {
	if i == nil || i.handle == 0 {
		return nil
	}
	return windows.CloseHandle(i.handle)
}

// SignalGUI asks a running GUI in the current logon session to show itself.
// It returns false if no GUI listener exists.
func SignalGUI() (bool, error) {
	ptr, err := windows.UTF16PtrFromString(guiEventName)
	if err != nil {
		return false, err
	}
	h, err := windows.OpenEvent(windows.EVENT_MODIFY_STATE, false, ptr)
	if errors.Is(err, windows.ERROR_FILE_NOT_FOUND) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("open GUI activation: %w", err)
	}
	defer windows.CloseHandle(h) //nolint:errcheck // short-lived event handle
	return true, windows.SetEvent(h)
}

// ListenGUI starts a cancellable activation listener. The callback must be
// safe to invoke from the listener goroutine.
func ListenGUI(ctx context.Context, show func()) (func(), error) {
	ptr, err := windows.UTF16PtrFromString(guiEventName)
	if err != nil {
		return nil, err
	}
	h, err := windows.CreateEvent(nil, 0, 0, ptr)
	if err != nil && !errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
		return nil, fmt.Errorf("create GUI activation: %w", err)
	}
	done := make(chan struct{})
	listenCtx, cancel := context.WithCancel(ctx)
	go func() {
		defer close(done)
		for {
			event, waitErr := windows.WaitForSingleObject(h, 1000)
			if waitErr != nil || listenCtx.Err() != nil {
				return
			}
			if event == windows.WAIT_OBJECT_0 {
				show()
			}
		}
	}()
	return func() {
		cancel()
		_ = windows.SetEvent(h)
		<-done
		_ = windows.CloseHandle(h) // listener shutdown; callback has exited
	}, nil
}
