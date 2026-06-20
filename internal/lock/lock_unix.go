//go:build !windows

package lock

import (
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
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, fmt.Errorf("lock: open %s: %w", path, err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("another goschedd instance is already running (lock held: %s)", path)
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
