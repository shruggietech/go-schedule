//go:build !windows

package clifixture

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLongSocketParentFailsBeforeCreatingOwnedState(t *testing.T) {
	parent := filepath.Join(t.TempDir(), strings.Repeat("long-parent-", 12))
	if err := os.Mkdir(parent, 0700); err != nil {
		t.Fatal(err)
	}
	if _, err := NewWorkspace(parent); err == nil || !strings.Contains(err.Error(), "shorter workspace parent") {
		t.Fatalf("invalid socket path accepted: %v", err)
	}
	entries, err := os.ReadDir(parent)
	if err != nil || len(entries) != 0 {
		t.Fatalf("failed placement left owned state: %v %v", entries, err)
	}
}
