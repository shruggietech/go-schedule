//go:build windows

package desktopcontrol

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

var errElevationCancelled = errors.New("elevation cancelled")

const seeMaskNoCloseProcess = 0x00000040

type shellExecuteInfo struct {
	size       uint32
	mask       uint32
	hwnd       windows.Handle
	verb       *uint16
	file       *uint16
	parameters *uint16
	directory  *uint16
	show       int32
	instance   windows.Handle
	idList     uintptr
	class      *uint16
	classKey   windows.Handle
	hotKey     uint32
	icon       windows.Handle
	process    windows.Handle
}

func requestElevation(action string) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locate executable: %w", err)
	}
	helper := filepath.Join(filepath.Dir(exe), "gosched-tray.exe")
	if _, err := os.Stat(helper); err != nil {
		return fmt.Errorf("service-control companion is unavailable at %s: %w", helper, err)
	}
	verb, err := windows.UTF16PtrFromString("runas")
	if err != nil {
		return err
	}
	file, err := windows.UTF16PtrFromString(helper)
	if err != nil {
		return err
	}
	params, err := windows.UTF16PtrFromString("--elevated-service-action " + action)
	if err != nil {
		return err
	}
	info := shellExecuteInfo{size: uint32(unsafe.Sizeof(shellExecuteInfo{})), mask: seeMaskNoCloseProcess, verb: verb, file: file, parameters: params}
	proc := windows.NewLazySystemDLL("shell32.dll").NewProc("ShellExecuteExW")
	ok, _, callErr := proc.Call(uintptr(unsafe.Pointer(&info)))
	runtime.KeepAlive(verb)
	runtime.KeepAlive(file)
	runtime.KeepAlive(params)
	if ok == 0 {
		if errors.Is(callErr, windows.ERROR_CANCELLED) {
			return errElevationCancelled
		}
		return fmt.Errorf("windows elevation prompt: %w", callErr)
	}
	defer windows.CloseHandle(info.process) //nolint:errcheck // process handle cleanup
	event, err := windows.WaitForSingleObject(info.process, 35_000)
	if err != nil {
		return fmt.Errorf("wait for service helper: %w", err)
	}
	if event != windows.WAIT_OBJECT_0 {
		return errors.New("service helper did not finish within 35 seconds")
	}
	var exitCode uint32
	if err := windows.GetExitCodeProcess(info.process, &exitCode); err != nil {
		return fmt.Errorf("read service helper outcome: %w", err)
	}
	if exitCode != 0 {
		return fmt.Errorf("windows service manager rejected %s (helper exit %d)", action, exitCode)
	}
	return nil
}
