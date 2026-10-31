package clifixture

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/platform"
)

func TestOwnedProcessTree(t *testing.T) {
	if mode := os.Getenv("GO_SCHEDULE_OWNED_PROCESS_HELPER"); mode != "" {
		path := os.Getenv("GO_SCHEDULE_OWNED_PROCESS_OUTPUT")
		if mode == "child" {
			for i := 0; ; i++ {
				_ = os.WriteFile(path, []byte(fmt.Sprint(i)), 0600)
				time.Sleep(30 * time.Millisecond)
			}
		}
		child := exec.Command(os.Args[0], "-test.run=^TestOwnedProcessTree$")
		platform.HideConsole(child)
		child.Env = append(os.Environ(), "GO_SCHEDULE_OWNED_PROCESS_HELPER=child")
		if err := child.Start(); err != nil {
			os.Exit(2)
		}
		_ = child.Wait()
		os.Exit(0)
	}
	dir := t.TempDir()
	heartbeat := filepath.Join(dir, "heartbeat")
	p, err := StartProcess(Command{Executable: os.Args[0], Arguments: []string{"-test.run=^TestOwnedProcessTree$"}, Environment: append(os.Environ(), "GO_SCHEDULE_OWNED_PROCESS_HELPER=parent", "GO_SCHEDULE_OWNED_PROCESS_OUTPUT="+heartbeat), Directory: dir, StdoutPath: filepath.Join(dir, "stdout"), StderrPath: filepath.Join(dir, "stderr")})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = p.Stop() })
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(heartbeat); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("owned grandchild did not become ready")
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err := p.Stop(); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = p.Wait(ctx)
	time.Sleep(80 * time.Millisecond)
	first, err := os.ReadFile(heartbeat)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(120 * time.Millisecond)
	second, err := os.ReadFile(heartbeat)
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Fatalf("owned grandchild survived cleanup: %q -> %q", first, second)
	}
	if err := p.Stop(); err != nil {
		t.Fatalf("repeated cleanup: %v", err)
	}
}

func TestProcessRequiresExplicitAbsolutePaths(t *testing.T) {
	if _, err := StartProcess(Command{Executable: "relative"}); err == nil {
		t.Fatal("implicit process location accepted")
	}
}
