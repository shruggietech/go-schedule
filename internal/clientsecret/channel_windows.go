//go:build windows

package clientsecret

import (
	"context"
	"net"
	"strings"

	"github.com/Microsoft/go-winio"
	"github.com/google/uuid"
	"golang.org/x/sys/windows"
)

func listenReadChannel() (net.Listener, string, func(), error) {
	token := windows.GetCurrentProcessToken()
	user, err := token.GetTokenUser()
	if err != nil {
		return nil, "", nil, err
	}
	endpoint := `\\.\pipe\gosched-secret-` + uuid.NewString()
	listener, err := winio.ListenPipe(endpoint, &winio.PipeConfig{SecurityDescriptor: "D:P(A;;GA;;;" + user.User.Sid.String() + ")"})
	if err != nil {
		return nil, "", nil, err
	}
	return listener, endpoint, func() { _ = listener.Close() }, nil
}
func dialReadChannel(ctx context.Context, endpoint string) (net.Conn, error) {
	if !strings.HasPrefix(endpoint, `\\.\pipe\gosched-secret-`) {
		return nil, errReadUnavailable
	}
	return winio.DialPipeContext(ctx, endpoint)
}
func readPeerPID(connection net.Conn, server bool) (int, error) {
	file, ok := connection.(interface{ Fd() uintptr })
	if !ok {
		return 0, errReadUnavailable
	}
	var pid uint32
	var err error
	if server {
		err = windows.GetNamedPipeClientProcessId(windows.Handle(file.Fd()), &pid)
	} else {
		err = windows.GetNamedPipeServerProcessId(windows.Handle(file.Fd()), &pid)
	}
	return int(pid), err
}
func readProcessExecutable(pid int) (string, error) {
	process, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return "", err
	}
	defer func() { _ = windows.CloseHandle(process) }()
	buffer := make([]uint16, 32768)
	size := uint32(len(buffer))
	if err := windows.QueryFullProcessImageName(process, 0, &buffer[0], &size); err != nil {
		return "", err
	}
	return windows.UTF16ToString(buffer[:size]), nil
}
