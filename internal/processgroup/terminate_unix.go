//go:build !windows

// Package processgroup terminates and drains explicitly owned Unix process groups.
package processgroup

import (
	"errors"
	"fmt"
	"syscall"
	"time"
)

// Terminate sends SIGKILL to an owned group and waits until no member remains.
// A cleanup timeout is an error; callers must retain state that may still be in use.
func Terminate(pgid int, timeout time.Duration) error {
	return terminate(pgid, timeout, syscall.Kill)
}

func terminate(pgid int, timeout time.Duration, signal func(int, syscall.Signal) error) error {
	if pgid <= 1 {
		return errors.New("owned process group ID must be greater than one")
	}
	if err := signal(-pgid, syscall.SIGKILL); err != nil {
		if errors.Is(err, syscall.ESRCH) {
			return nil
		}
		return fmt.Errorf("terminate owned process group: %w", err)
	}
	deadline := time.Now().Add(timeout)
	for {
		err := signal(-pgid, 0)
		if errors.Is(err, syscall.ESRCH) {
			return nil
		}
		// An accepted termination can be followed by an ambiguous permission probe.
		// Keep waiting for ESRCH; permission ambiguity is never successful cleanup.
		if err != nil && !errors.Is(err, syscall.EPERM) {
			return fmt.Errorf("inspect owned process group during cleanup: %w", err)
		}
		if !time.Now().Before(deadline) {
			if err != nil {
				return fmt.Errorf("owned process group exit could not be confirmed at cleanup deadline: %w", err)
			}
			return errors.New("owned process group still has members at cleanup deadline")
		}
		time.Sleep(min(10*time.Millisecond, time.Until(deadline)))
	}
}
