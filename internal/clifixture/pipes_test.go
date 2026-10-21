package clifixture

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestOwnedPipedHelper(t *testing.T) {
	if os.Args[len(os.Args)-1] != "gosched-qa-pipe" {
		return
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		os.Exit(3)
	}
	fmt.Printf("echo:%s", data)
	fmt.Fprint(os.Stderr, "fixture diagnostic")
	os.Exit(7)
}

func TestOwnedPipesCaptureExactOutputAndExit(t *testing.T) {
	root, err := filepath.Abs(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	p, pipes, err := StartPiped(Command{Executable: os.Args[0], Arguments: []string{"-test.run=^TestOwnedPipedHelper$", "--", "gosched-qa-pipe"}, Directory: root})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := p.Stop(); err != nil {
			t.Error(err)
		}
		if err := pipes.Close(); err != nil {
			t.Error(err)
		}
	}()
	if _, err := io.WriteString(pipes.Stdin, "hello 世界"); err != nil {
		t.Fatal(err)
	}
	if err := pipes.Stdin.Close(); err != nil {
		t.Fatal(err)
	}
	output, err := io.ReadAll(pipes.Stdout)
	if err != nil {
		t.Fatal(err)
	}
	diagnostic, err := io.ReadAll(pipes.Stderr)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := p.Wait(ctx); err == nil {
		t.Fatal("nonzero child reported success")
	}
	code, done := p.ExitCode()
	if !done || code != 7 || string(output) != "echo:hello 世界" || !strings.Contains(string(diagnostic), "fixture diagnostic") {
		t.Fatalf("code=%d done=%v stdout=%q stderr=%q", code, done, output, diagnostic)
	}
}
