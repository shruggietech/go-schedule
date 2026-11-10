//go:build windows

package executor

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
	"unicode/utf16"

	"golang.org/x/sys/windows"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestWindowsOwnedCommandResolvesRelativeExecutableInWorkingDirectory(t *testing.T) {
	dir := t.TempDir()
	binary, err := os.ReadFile(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(dir, "owned-helper.exe"), binary, 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(dir, "release"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	task := cancellationTask("child", dir)
	task.Command = ".\\owned-helper"
	task.WorkingDir = dir
	run := New(0).Run(context.Background(), task, time.Now().UTC(), domain.TriggerManual)
	if run.Outcome != domain.OutcomeSuccess || !strings.Contains(run.Output, "partial-output") {
		t.Fatalf("relative command=%+v", run)
	}
}

func TestWindowsEnvironmentPreservesDriveEntriesAndLastCaseInsensitiveOverride(t *testing.T) {
	encoded, err := windowsEnvironment([]string{"PAtH=first", "PATH=selected", "=C:=C:\\fixture", "GOSCHED_UNICODE=héllo 世界"})
	if err != nil {
		t.Fatal(err)
	}
	values := strings.Split(string(utf16.Decode(encoded)), "\x00")
	got := strings.Join(values, "|")
	if strings.Contains(got, "first") || !strings.Contains(got, "PATH=selected") || !strings.Contains(got, "=C:=C:\\fixture") || !strings.Contains(got, "héllo 世界") {
		t.Fatalf("normalized environment=%q", got)
	}
	for _, invalid := range []string{"invalid", "NAME=bad\x00value"} {
		if _, err := windowsEnvironment([]string{invalid}); err == nil {
			t.Fatalf("accepted invalid environment %q", invalid)
		}
	}
}

func TestWindowsOwnedCommandPreservesExplicitBreakawayChild(t *testing.T) {
	dir := t.TempDir()
	task := cancellationTask("windows-breakaway", dir)
	task.Args = []string{"-test.run=^TestWindowsBreakawayHelper$"}
	run := New(0).Run(context.Background(), task, time.Now().UTC(), domain.TriggerManual)
	if run.Outcome != domain.OutcomeSuccess {
		t.Fatalf("explicit breakaway command failed: %+v", run)
	}
	waitCancellationFile(t, filepath.Join(dir, "ready"))
	if err := os.WriteFile(filepath.Join(dir, "release"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	waitCancellationFile(t, filepath.Join(dir, "finished"))
}

func TestWindowsBreakawayHelper(t *testing.T) {
	if os.Getenv("GOSCHED_EXECUTOR_CANCEL_MODE") != "windows-breakaway" {
		return
	}
	child := exec.Command(os.Args[0], "-test.run=^TestExecutorCancellationHelper$")
	child.Env = append(os.Environ(), "GOSCHED_EXECUTOR_CANCEL_MODE=child")
	child.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: windows.CREATE_NO_WINDOW | windows.CREATE_BREAKAWAY_FROM_JOB}
	if err := child.Start(); err != nil {
		t.Fatal(err)
	}
	if err := child.Process.Release(); err != nil {
		t.Fatal(err)
	}
	os.Exit(0)
}
