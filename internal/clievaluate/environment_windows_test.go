//go:build windows

package clievaluate

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
)

func TestWindowsLimitProbeHelper(t *testing.T) {
	if os.Args[len(os.Args)-1] != "s108-limit-probe" {
		return
	}
	limits, err := CurrentLimits("", DefaultBudget())
	if err != nil {
		os.Exit(2)
	}
	if err := json.NewEncoder(os.Stdout).Encode(limits); err != nil {
		os.Exit(3)
	}
	os.Exit(0)
}

func TestWindowsWorkerVerifiesAppliedHardLimits(t *testing.T) {
	b := DefaultBudget()
	p, pipes, err := clifixture.StartPiped(clifixture.Command{Executable: os.Args[0], Arguments: []string{"-test.run=^TestWindowsLimitProbeHelper$", "--", "s108-limit-probe"}, Directory: t.TempDir(), Limits: &clifixture.ProcessLimits{CPUPercent: b.CPUPercent, MemoryBytes: uint64(b.MemoryBytes), MaxProcesses: uint32(b.WindowsProcesses)}})
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
	if err := pipes.Stdin.Close(); err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(pipes.Stdout)
	if err != nil {
		t.Fatal(err)
	}
	diagnostic, err := io.ReadAll(pipes.Stderr)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := p.Wait(ctx); err != nil {
		t.Fatalf("controlled worker: %v %s", err, diagnostic)
	}
	var applied []NamedValue
	if err := json.Unmarshal(out, &applied); err != nil || len(applied) < 5 {
		t.Fatalf("actual controls=%s err=%v", out, err)
	}
}
