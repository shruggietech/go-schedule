//go:build darwin

package clientsecret

import (
	"bytes"
	"runtime"
	"unsafe"

	"golang.org/x/sys/unix"
)

func readUnixPeerPID(fd int) (int, error) {
	return unix.GetsockoptInt(fd, unix.SOL_LOCAL, unix.LOCAL_PEERPID)
}
func readProcessExecutable(pid int) (string, error) {
	// Match libproc's proc_pidpath: PROC_INFO_CALL_PIDINFO (2),
	// PROC_PIDPATHINFO (11). The kernel resolves the executable vnode;
	// mutable argv or process titles are never authentication input.
	buffer := make([]byte, 4096)
	_, _, errno := unix.Syscall6(unix.SYS_PROC_INFO, 2, uintptr(pid), 11, 0, uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))
	runtime.KeepAlive(buffer)
	if errno != 0 {
		return "", errReadUnavailable
	}
	end := bytes.IndexByte(buffer, 0)
	if end <= 0 {
		return "", errReadUnavailable
	}
	return string(buffer[:end]), nil
}
