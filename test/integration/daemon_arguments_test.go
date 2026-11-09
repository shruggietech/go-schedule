package integration

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
)

func TestPackagedDaemonRejectsPositionalArgumentsBeforeStartup(t *testing.T) {
	w, err := clifixture.NewWorkspace("")
	if err != nil {
		t.Fatal(err)
	}
	cleanupReady := true
	t.Cleanup(func() {
		if cleanupReady {
			if err := w.CleanupRuntime(); err != nil {
				t.Error(err)
			}
		}
	})
	bin := preparePackagedCandidate(t, w.Runtime, "goschedd", "goschedd")
	out, diagnostic := filepath.Join(w.Evidence, "stdout.txt"), filepath.Join(w.Evidence, "stderr.txt")
	p, err := clifixture.StartProcess(clifixture.Command{Executable: bin, Arguments: []string{"--config", w.ConfigPath, "unexpected-operand"}, Directory: w.Runtime, StdoutPath: out, StderrPath: diagnostic})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	waitErr := p.Wait(ctx)
	if err := p.Stop(); err != nil {
		cleanupReady = false
		t.Fatal(err)
	}
	code, known := p.ExitCode()
	stdout, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	stderr, err := os.ReadFile(diagnostic)
	if err != nil {
		t.Fatal(err)
	}
	_, stateErr := os.Stat(w.Config.DBPath())
	if !known || code != 2 || !strings.Contains(string(stderr), "unexpected argument") || len(stdout) != 0 || !os.IsNotExist(stateErr) {
		t.Fatalf("stray operand was not refused before startup: known=%t exit=%d wait=%v stdout=%q stderr=%q database_stat=%v; private evidence=%s", known, code, waitErr, stdout, stderr, stateErr, w.Evidence)
	}
}
