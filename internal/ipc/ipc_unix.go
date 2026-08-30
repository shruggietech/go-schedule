//go:build !windows

package ipc

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/shruggietech/go-schedule/internal/config"
)

type unixIPCOperations struct {
	lookupGroup func(string) (int, error)
	mkdirAll    func(string, os.FileMode) error
	stat        func(string) (os.FileInfo, error)
	chown       func(string, int, int) error
	chmod       func(string, os.FileMode) error
	remove      func(string) error
	listen      func(string, string) (net.Listener, error)
}

var productionUnixIPCOperations = unixIPCOperations{
	lookupGroup: func(name string) (int, error) {
		group, err := user.LookupGroup(name)
		if err != nil {
			return 0, err
		}
		gid, err := strconv.Atoi(group.Gid)
		if err != nil {
			return 0, fmt.Errorf("parse gid %q: %w", group.Gid, err)
		}
		return gid, nil
	},
	mkdirAll: os.MkdirAll,
	stat:     os.Stat,
	chown:    os.Chown,
	chmod:    os.Chmod,
	remove:   os.Remove,
	listen:   net.Listen,
}

func defaultEndpoint(cfg config.Config) string {
	return filepath.Join(cfg.DataDir, "goschedd.sock")
}

// Listen creates a Unix domain socket and verifies its effective local access
// policy before returning it to the daemon.
func Listen(cfg config.Config) (net.Listener, AccessInfo, error) {
	return listenUnix(cfg, productionUnixIPCOperations)
}

func listenUnix(cfg config.Config, ops unixIPCOperations) (net.Listener, AccessInfo, error) {
	endpoint := Endpoint(cfg)
	parent := filepath.Dir(endpoint)
	managedParent := cfg.IPCPath == ""
	access := AccessInfo{Mode: AccessModeCompatibility}
	wantParentMode := os.FileMode(0o755)
	wantSocketMode := os.FileMode(0o666)
	wantGID := -1

	if cfg.AdminGroup != "" {
		gid, err := ops.lookupGroup(cfg.AdminGroup)
		if err != nil {
			return nil, AccessInfo{}, fmt.Errorf("ipc: resolve admin_group %q: %w", cfg.AdminGroup, err)
		}
		wantGID = gid
		wantParentMode = 0o770
		wantSocketMode = 0o660
		access = AccessInfo{Mode: AccessModeRestricted, AdminGroup: cfg.AdminGroup}
	}

	missingParents := make([]string, 0)
	for candidate := parent; ; candidate = filepath.Dir(candidate) {
		if _, err := ops.stat(candidate); err == nil {
			break
		} else if !os.IsNotExist(err) {
			return nil, AccessInfo{}, fmt.Errorf("ipc: inspect socket dir %s: %w", candidate, err)
		}
		missingParents = append(missingParents, candidate)
		next := filepath.Dir(candidate)
		if next == candidate {
			return nil, AccessInfo{}, fmt.Errorf("ipc: no existing ancestor for socket dir %s", parent)
		}
	}
	createdParent := len(missingParents) > 0
	if createdParent {
		if err := ops.mkdirAll(parent, wantParentMode); err != nil {
			return nil, AccessInfo{}, fmt.Errorf("ipc: create socket dir %s: %w", parent, err)
		}
	}

	pathsToSecure := make([]string, 0, len(missingParents)+1)
	for i := len(missingParents) - 1; i >= 0; i-- {
		pathsToSecure = append(pathsToSecure, missingParents[i])
	}
	if managedParent && !createdParent {
		pathsToSecure = append(pathsToSecure, parent)
	}
	for _, path := range pathsToSecure {
		if wantGID >= 0 {
			if err := ops.chown(path, -1, wantGID); err != nil {
				return nil, AccessInfo{}, fmt.Errorf("ipc: set socket dir group %s: %w", path, err)
			}
		}
		if err := ops.chmod(path, wantParentMode); err != nil {
			return nil, AccessInfo{}, fmt.Errorf("ipc: set socket dir permissions %s: %w", path, err)
		}
		if err := verifyUnixPath(ops, path, wantParentMode, wantGID); err != nil {
			return nil, AccessInfo{}, fmt.Errorf("ipc: verify socket dir %s: %w", path, err)
		}
	}
	if len(pathsToSecure) == 0 && cfg.AdminGroup != "" {
		if err := verifyUnixPath(ops, parent, wantParentMode, wantGID); err != nil {
			return nil, AccessInfo{}, fmt.Errorf("ipc: verify socket dir %s: %w", parent, err)
		}
	}

	if info, err := ops.stat(endpoint); err == nil {
		if info.Mode()&os.ModeSocket == 0 {
			return nil, AccessInfo{}, fmt.Errorf("ipc: stale endpoint %s is not a socket", endpoint)
		}
		if err := ops.remove(endpoint); err != nil {
			return nil, AccessInfo{}, fmt.Errorf("ipc: remove stale socket: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return nil, AccessInfo{}, fmt.Errorf("ipc: inspect stale socket: %w", err)
	}

	l, err := ops.listen("unix", endpoint)
	if err != nil {
		return nil, AccessInfo{}, fmt.Errorf("ipc: listen %s: %w", endpoint, err)
	}
	cleanup := func() {
		_ = l.Close()
		_ = ops.remove(endpoint)
	}
	if wantGID >= 0 {
		if err := ops.chown(endpoint, -1, wantGID); err != nil {
			cleanup()
			return nil, AccessInfo{}, fmt.Errorf("ipc: set socket group %s: %w", endpoint, err)
		}
	}
	if err := ops.chmod(endpoint, wantSocketMode); err != nil {
		cleanup()
		return nil, AccessInfo{}, fmt.Errorf("ipc: set socket permissions %s: %w", endpoint, err)
	}
	if err := verifyUnixPath(ops, endpoint, wantSocketMode, wantGID); err != nil {
		cleanup()
		return nil, AccessInfo{}, fmt.Errorf("ipc: verify socket %s: %w", endpoint, err)
	}
	return l, access, nil
}

func verifyUnixPath(ops unixIPCOperations, path string, wantMode os.FileMode, wantGID int) error {
	info, err := ops.stat(path)
	if err != nil {
		return err
	}
	if got := info.Mode().Perm(); got != wantMode {
		return fmt.Errorf("permissions are %04o, want %04o", got, wantMode)
	}
	if wantGID >= 0 {
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("group ownership is unavailable")
		}
		if got := int(stat.Gid); got != wantGID {
			return fmt.Errorf("group id is %d, want %d", got, wantGID)
		}
	}
	return nil
}

// DialContext connects to the Unix domain socket endpoint.
func DialContext(ctx context.Context, endpoint string) (net.Conn, error) {
	var d net.Dialer
	return d.DialContext(ctx, "unix", endpoint)
}
