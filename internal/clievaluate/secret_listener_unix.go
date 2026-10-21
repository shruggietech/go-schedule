//go:build !windows

package clievaluate

import (
	"net"
	"os"
	"path/filepath"
)

func secretListener(root string) (net.Listener, string, error) {
	endpoint := filepath.Join(root, "ipc", "secret.sock")
	listener, err := net.Listen("unix", endpoint)
	if err != nil {
		return nil, "", err
	}
	if err := os.Chmod(endpoint, 0600); err != nil {
		_ = listener.Close()
		return nil, "", err
	}
	return listener, endpoint, nil
}
