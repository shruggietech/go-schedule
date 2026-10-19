//go:build !windows

package processgroup

import (
	"errors"
	"syscall"
	"testing"
	"time"
)

func TestTerminateDrainsEveryGroupMemberBeforeSuccess(t *testing.T) {
	probes, kills := 0, 0
	err := terminate(12345, time.Second, func(pid int, signal syscall.Signal) error {
		if pid != -12345 {
			t.Fatalf("unowned signal target %d", pid)
		}
		if signal == syscall.SIGKILL {
			kills++
			return nil
		}
		if signal != 0 {
			t.Fatalf("unexpected signal %v", signal)
		}
		probes++
		if probes < 4 {
			return nil
		}
		return syscall.ESRCH
	})
	if err != nil || kills != 1 || probes != 4 {
		t.Fatalf("err=%v kills=%d probes=%d", err, kills, probes)
	}
}
func TestTerminateRetainsFailedCleanupAndRejectsUnownedGroups(t *testing.T) {
	for _, id := range []int{-1, 0, 1} {
		if err := terminate(id, time.Second, func(int, syscall.Signal) error { t.Fatal("signaled an unowned group"); return nil }); err == nil {
			t.Fatalf("accepted group %d", id)
		}
	}
	if err := terminate(12345, 0, func(int, syscall.Signal) error { return nil }); err == nil {
		t.Fatal("live group reported clean at deadline")
	}
	if err := terminate(12345, time.Second, func(int, syscall.Signal) error { return syscall.EPERM }); !errors.Is(err, syscall.EPERM) {
		t.Fatal(err)
	}
	if err := terminate(12345, time.Second, func(int, syscall.Signal) error { return syscall.ESRCH }); err != nil {
		t.Fatal(err)
	}
}

func TestTerminateRetriesPostSignalPermissionProbeWithoutFalseSuccess(t *testing.T) {
	probes := 0
	err := terminate(12345, time.Second, func(_ int, signal syscall.Signal) error {
		if signal == syscall.SIGKILL {
			return nil
		}
		probes++
		if probes == 1 {
			return syscall.EPERM
		}
		if probes == 2 {
			return nil
		}
		return syscall.ESRCH
	})
	if err != nil || probes != 3 {
		t.Fatalf("post-signal probe ended cleanup prematurely: %v, probes=%d", err, probes)
	}
	if err := terminate(12345, 0, func(_ int, signal syscall.Signal) error {
		if signal == syscall.SIGKILL {
			return nil
		}
		return syscall.EPERM
	}); err == nil {
		t.Fatal("unconfirmed group reported clean")
	}
}
