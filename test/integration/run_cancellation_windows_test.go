//go:build windows

package integration

import "golang.org/x/sys/windows"

func s143ProcessExited(pid int) bool {
	handle, err := windows.OpenProcess(windows.SYNCHRONIZE|windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return err == windows.ERROR_INVALID_PARAMETER
	}
	status, waitErr := windows.WaitForSingleObject(handle, 0)
	closeErr := windows.CloseHandle(handle)
	return waitErr == nil && closeErr == nil && status == windows.WAIT_OBJECT_0
}
