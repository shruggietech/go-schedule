//go:build linux

package clientsecret

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func readUnixPeerPID(fd int) (int, error) {
	credentials, err := unix.GetsockoptUcred(fd, unix.SOL_SOCKET, unix.SO_PEERCRED)
	if err != nil {
		return 0, err
	}
	return int(credentials.Pid), nil
}
func readProcessExecutable(pid int) (string, error) {
	return os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
}
