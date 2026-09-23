//go:build linux

package desktopcontrol

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/unix"
)

func sessionKey() string {
	identity := os.Getenv("XDG_SESSION_ID")
	if identity == "" {
		identity = os.Getenv("WAYLAND_DISPLAY")
	}
	if identity == "" {
		identity = os.Getenv("DISPLAY")
	}
	if identity == "" {
		identity = "default"
	}
	sum := sha256.Sum256([]byte(identity))
	return hex.EncodeToString(sum[:8])
}

func guiSocketName() string { return "gui-open-" + sessionKey() + ".sock" }

// Instance owns a Linux session-local advisory lock until Close.
type Instance struct{ file *os.File }

func sessionDir() (string, error) {
	base := os.Getenv("XDG_RUNTIME_DIR")
	if base == "" {
		var err error
		base, err = os.UserCacheDir()
		if err != nil {
			return "", fmt.Errorf("locate user cache directory: %w", err)
		}
	}
	dir := filepath.Join(base, "go-schedule")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("create desktop session directory: %w", err)
	}
	info, err := os.Lstat(dir)
	if err != nil {
		return "", fmt.Errorf("inspect desktop session directory: %w", err)
	}
	if !info.IsDir() || info.Mode().Perm()&0077 != 0 {
		return "", fmt.Errorf("desktop session directory %s must be a private directory", dir)
	}
	return dir, nil
}

// ClaimGUI returns false if this user session already owns the GUI lock.
func ClaimGUI() (*Instance, bool, error) { return claimLinux("gui-" + sessionKey() + ".lock") }

// ClaimTray returns false if this user session already owns the indicator lock.
func ClaimTray() (*Instance, bool, error) { return claimLinux("indicator-" + sessionKey() + ".lock") }

func claimLinux(name string) (*Instance, bool, error) {
	dir, err := sessionDir()
	if err != nil {
		return nil, false, err
	}
	file, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, false, fmt.Errorf("open desktop lock: %w", err)
	}
	if err := unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		_ = file.Close()
		if errors.Is(err, unix.EWOULDBLOCK) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("claim desktop lock: %w", err)
	}
	return &Instance{file: file}, true, nil
}

// Close releases the session-local lock.
func (i *Instance) Close() error {
	if i == nil || i.file == nil {
		return nil
	}
	return i.file.Close()
}

// SignalGUI asks the existing GUI to show itself using a same-user socket.
func SignalGUI() (bool, error) {
	dir, err := sessionDir()
	if err != nil {
		return false, err
	}
	conn, err := net.DialTimeout("unix", filepath.Join(dir, guiSocketName()), 500*time.Millisecond)
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, unix.ECONNREFUSED) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("signal existing GUI: %w", err)
	}
	defer conn.Close() //nolint:errcheck // short-lived activation connection
	if _, err := conn.Write([]byte("open\n")); err != nil {
		return false, fmt.Errorf("write GUI activation: %w", err)
	}
	return true, nil
}

// ListenGUI starts an activation listener owned by the GUI lifecycle.
func ListenGUI(ctx context.Context, show func()) (func(), error) {
	dir, err := sessionDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, guiSocketName())
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("remove stale GUI socket: %w", err)
	}
	listener, err := net.Listen("unix", path)
	if err != nil {
		return nil, fmt.Errorf("listen for GUI activation: %w", err)
	}
	if err := os.Chmod(path, 0600); err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("protect GUI activation socket: %w", err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			var buf [5]byte
			_ = conn.SetReadDeadline(time.Now().Add(time.Second))
			_, readErr := io.ReadFull(conn, buf[:])
			_ = conn.Close()
			if readErr == nil && string(buf[:]) == "open\n" && ctx.Err() == nil {
				show()
			}
		}
	}()
	return func() {
		_ = listener.Close()
		<-done
		_ = os.Remove(path)
	}, nil
}
