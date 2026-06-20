//go:build windows

package lock

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

// Lock holds an acquired exclusive Windows file lock.
type Lock struct {
	f *os.File
}

// Acquire takes an exclusive, fail-immediately lock on path, creating the file
// if needed. It returns an error if another process already holds the lock.
func Acquire(path string) (*Lock, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, fmt.Errorf("lock: open %s: %w", path, err)
	}
	ol := new(windows.Overlapped)
	err = windows.LockFileEx(
		windows.Handle(f.Fd()),
		windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY,
		0, 1, 0, ol,
	)
	if err != nil {
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
	ol := new(windows.Overlapped)
	_ = windows.UnlockFileEx(windows.Handle(l.f.Fd()), 0, 1, 0, ol)
	return l.f.Close()
}
