//go:build !windows

package clientsecret

import (
	"context"
	"net"
	"os"
	"path/filepath"
)

func listenReadChannel() (net.Listener, string, func(), error) {
	directory, err := os.MkdirTemp("/tmp", "gosched-secret-")
	if err != nil {
		return nil, "", nil, err
	}
	endpoint := filepath.Join(directory, "s")
	listener, err := net.Listen("unix", endpoint)
	if err != nil {
		_ = os.Remove(directory)
		return nil, "", nil, err
	}
	if err := os.Chmod(endpoint, 0600); err != nil {
		_ = listener.Close()
		_ = os.Remove(directory)
		return nil, "", nil, err
	}
	return listener, endpoint, func() { _ = listener.Close(); _ = os.Remove(directory) }, nil
}
func dialReadChannel(ctx context.Context, endpoint string) (net.Conn, error) {
	if !filepath.IsAbs(endpoint) {
		return nil, errReadUnavailable
	}
	var dialer net.Dialer
	return dialer.DialContext(ctx, "unix", endpoint)
}
func readPeerPID(connection net.Conn, _ bool) (int, error) {
	unixConnection, ok := connection.(*net.UnixConn)
	if !ok {
		return 0, errReadUnavailable
	}
	raw, err := unixConnection.SyscallConn()
	if err != nil {
		return 0, err
	}
	var pid int
	var queryErr error
	if err := raw.Control(func(fd uintptr) { pid, queryErr = readUnixPeerPID(int(fd)) }); err != nil {
		return 0, err
	}
	return pid, queryErr
}
