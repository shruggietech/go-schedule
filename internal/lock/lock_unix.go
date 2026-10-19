//go:build !windows

package lock

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

// Lock holds an acquired exclusive flock.
type Lock struct {
	f *os.File
}

// Acquire takes an exclusive, non-blocking lock on path, creating the file if
// needed. It returns an error if another process already holds the lock.
func Acquire(path string) (*Lock, error) {
	held, err := AcquireFile(path, 0o644)
	if errors.Is(err, ErrHeld) {
		return nil, fmt.Errorf("another goschedd instance is already running (lock held: %s)", path)
	}
	return held, err
}

// AcquireFile takes an exclusive nonblocking OS lock, creating its file with mode.
// A leftover file is harmless; handle closure or process exit releases ownership.
func AcquireFile(path string, mode os.FileMode) (*Lock, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, mode)
	if err != nil {
		return nil, fmt.Errorf("lock: open %s: %w", path, err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = f.Close()
		if errors.Is(err, syscall.EWOULDBLOCK) || errors.Is(err, syscall.EAGAIN) {
			return nil, fmt.Errorf("lock %s: %w", path, ErrHeld)
		}
		return nil, fmt.Errorf("lock: acquire %s: %w", path, err)
	}
	return &Lock{f: f}, nil
}

// Release unlocks and closes the lock file.
func (l *Lock) Release() error {
	if l == nil || l.f == nil {
		return nil
	}
	_ = syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
	return l.f.Close()
}
