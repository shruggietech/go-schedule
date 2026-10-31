package userprefs

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/lock"
	"github.com/shruggietech/go-schedule/v2/internal/platform"
)

func TestPreferenceProcessExitReleasesOwnership(t *testing.T) {
	store := testStore(t)
	if _, err := store.Load(); err != nil {
		t.Fatal(err)
	}
	ready := filepath.Join(t.TempDir(), "ready")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestPreferenceLockChild$", "--", store.Path()+".lock", ready)
	command.Env = append(os.Environ(), "GOSCHEDULE_PREFS_LOCK_CHILD=1")
	command.Stdin = strings.NewReader("")
	var diagnostics bytes.Buffer
	command.Stdout = &diagnostics
	command.Stderr = &diagnostics
	platform.HideConsole(command)
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	waited := false
	defer func() {
		if !waited {
			_ = command.Process.Kill()
			_ = command.Wait()
		}
	}()
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		if ctx.Err() != nil {
			t.Fatalf("child did not acquire isolated lock: %v", ctx.Err())
		}
		time.Sleep(5 * time.Millisecond)
	}
	// Killing the exact owned child exercises release without its normal cleanup path.
	if err := command.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	if err := command.Wait(); err == nil {
		t.Fatal("owned child unexpectedly exited successfully")
	}
	waited = true
	dark := AppearanceDark
	if prefs, err := store.Patch(Patch{Appearance: &dark}); err != nil || prefs.Appearance != AppearanceDark {
		t.Fatalf("process exit retained ownership: %+v %v", prefs, err)
	}
}

func TestPreferenceLockChild(t *testing.T) {
	if os.Getenv("GOSCHEDULE_PREFS_LOCK_CHILD") != "1" {
		return
	}
	held, err := lock.AcquireFile(os.Args[len(os.Args)-2], 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer held.Release()
	if err := os.WriteFile(os.Args[len(os.Args)-1], []byte("ready"), 0o600); err != nil {
		t.Fatal(err)
	}
	for {
		time.Sleep(time.Second)
	}
}
