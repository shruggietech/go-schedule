//go:build windows

package clievaluate

import (
	"errors"
	"strconv"

	"golang.org/x/sys/windows"
)

func processIdentity(pid int) (string, error) {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		if errors.Is(err, windows.ERROR_INVALID_PARAMETER) {
			return "", nil
		}
		return "", err
	}
	defer func() { _ = windows.CloseHandle(handle) }()
	var created, exit, kernel, user windows.Filetime
	if err := windows.GetProcessTimes(handle, &created, &exit, &kernel, &user); err != nil {
		return "", err
	}
	var code uint32
	if err := windows.GetExitCodeProcess(handle, &code); err != nil {
		return "", err
	}
	if code != 259 {
		return "", nil
	}
	return strconv.FormatUint(uint64(created.HighDateTime)<<32|uint64(created.LowDateTime), 10), nil
}
