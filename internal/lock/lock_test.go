package lock

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenericFileLockAndDaemonCompatibility(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shared.lock")
	held, err := AcquireFile(path, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if second, err := AcquireFile(path, 0o600); !errors.Is(err, ErrHeld) {
		if second != nil {
			_ = second.Release()
		}
		_ = held.Release()
		t.Fatalf("expected generic contention, got %v", err)
	}
	if second, err := Acquire(path); err == nil || !strings.Contains(err.Error(), "another goschedd instance is already running") {
		if second != nil {
			_ = second.Release()
		}
		_ = held.Release()
		t.Fatalf("daemon contention contract changed: %v", err)
	}
	if err := held.Release(); err != nil {
		t.Fatal(err)
	}
	// The persistent file has no ownership once its handle has closed.
	again, err := AcquireFile(path, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := again.Release(); err != nil {
		t.Fatal(err)
	}
}
