//go:build windows

package winuninstall

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveDoesNotFollowDescendantReparsePoint(t *testing.T) {
	base := t.TempDir()
	owned := filepath.Join(base, "owned")
	outside := t.TempDir()
	sentinel := filepath.Join(outside, "sentinel.txt")
	if err := os.Mkdir(owned, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sentinel, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(owned, "redirect")); err != nil {
		if errors.Is(err, os.ErrPermission) {
			t.Skipf("directory symlink privilege unavailable: %v", err)
		}
		t.Fatalf("create directory symlink: %v", err)
	}
	backend := &windowsBackend{}
	if err := backend.Remove(Target{Path: owned, base: base, relative: "owned"}); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(sentinel); err != nil || string(data) != "keep" {
		t.Fatalf("outside sentinel changed: data=%q err=%v", data, err)
	}
	if _, err := os.Stat(owned); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("owned root remains: %v", err)
	}
}

func TestPreflightRejectsReparseRoot(t *testing.T) {
	base := t.TempDir()
	outside := t.TempDir()
	owned := filepath.Join(base, "owned")
	if err := os.Symlink(outside, owned); err != nil {
		if errors.Is(err, os.ErrPermission) {
			t.Skipf("directory symlink privilege unavailable: %v", err)
		}
		t.Fatalf("create directory symlink: %v", err)
	}
	backend := &windowsBackend{}
	if _, err := backend.Preflight(Target{Path: owned, base: base, relative: "owned"}); err == nil {
		t.Fatal("preflight accepted a reparse-point root")
	}
}

func TestPreflightRejectsCandidateOutsideDeclaredRoot(t *testing.T) {
	base := t.TempDir()
	outside := filepath.Join(t.TempDir(), "owned")
	if err := os.Mkdir(outside, 0o700); err != nil {
		t.Fatal(err)
	}
	backend := &windowsBackend{}
	if _, err := backend.Preflight(Target{Path: outside, base: base, relative: "owned"}); err == nil {
		t.Fatal("preflight accepted a candidate outside its declared root")
	}
}
