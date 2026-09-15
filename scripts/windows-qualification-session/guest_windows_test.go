package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestGuestNondestructiveFixturesAndHostRefusal(t *testing.T) {
	script, err := filepath.Abs(filepath.Join("..", "..", "test", "windows", "Start-QualificationGuest.ps1"))
	if err != nil {
		t.Fatal(err)
	}
	engine := filepath.Join(os.Getenv("WINDIR"), "System32", "WindowsPowerShell", "v1.0", "powershell.exe")
	out := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, engine, "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-File", script, "-Fixture", "-FixtureOutput", out)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("fixture: %v\n%s", err, b)
	}
	for _, c := range []struct {
		name, status string
		code         *int
	}{{"success", "completed", intPointer(0)}, {"failure", "failed", intPointer(7)}, {"timeout", "timed-out", nil}, {"reboot", "completed-reboot-required", intPointer(3010)}} {
		b, err := os.ReadFile(filepath.Join(out, c.name+"-phase.json"))
		if err != nil {
			t.Fatal(err)
		}
		var r struct {
			Status   string `json:"status"`
			Attended string `json:"attended_status"`
			Code     *int   `json:"exit_code"`
		}
		if err := json.Unmarshal(b, &r); err != nil {
			t.Fatal(err)
		}
		if r.Status != c.status || r.Attended != "unavailable" || (r.Code == nil) != (c.code == nil) {
			t.Fatalf("bad phase %s: %s", c.name, b)
		}
		if c.code != nil && *r.Code != *c.code {
			t.Fatalf("bad exit code: %s", b)
		}
	}
	cmd = exec.CommandContext(ctx, engine, "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-File", script)
	for _, stream := range []string{"stdout", "stderr"} {
		b, err := os.ReadFile(filepath.Join(out, "timeout-"+stream+".log"))
		if err != nil || !strings.Contains(string(b), "timeout-"+stream) {
			t.Fatalf("timeout diagnostics lost: %v %s", err, b)
		}
	}
	if _, err := os.Stat(filepath.Join(out, "timeout-output.json")); err != nil {
		t.Fatal(err)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	b, err := cmd.CombinedOutput()
	if err == nil || !strings.Contains(string(b), "only in the prepared Windows Sandbox") {
		t.Fatalf("host invocation not refused: %v %s", err, b)
	}
}

func intPointer(n int) *int { return &n }

func TestRejectDirectoryJunction(t *testing.T) {
	root := canonicalTempDir(t)
	target := filepath.Join(root, "target")
	junction := filepath.Join(root, "junction")
	if err := os.Mkdir(target, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "input"), []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("cmd.exe", "/d", "/c", "mklink", "/J", junction, target)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("junction fixture: %v %s", err, b)
	}
	if err := rejectLinks(filepath.Join(junction, "input")); err == nil {
		t.Fatal("junction input accepted")
	}
	if err := rejectLinks(junction); err == nil {
		t.Fatal("junction output parent accepted")
	}
}
