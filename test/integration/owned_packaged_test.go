package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
)

// ownedPackagedCLI binds every command to one disposable daemon and retains cleanup failures.
type ownedPackagedCLI struct {
	t            *testing.T
	workspace    *clifixture.Workspace
	api          *client.Client
	ctx          context.Context
	cli          string
	sequence     int
	cleanupReady bool
}

func newOwnedPackagedCLI(t *testing.T, workers int) *ownedPackagedCLI {
	t.Helper()
	w, err := clifixture.NewWorkspace("")
	if err != nil {
		t.Fatal(err)
	}
	h := &ownedPackagedCLI{t: t, workspace: w, cleanupReady: true}
	t.Cleanup(func() {
		if !h.cleanupReady {
			t.Logf("retained owned runtime: %s", w.Runtime)
			return
		}
		if err := w.CleanupRuntime(); err != nil {
			t.Error(err)
		}
	})
	w.Config.WorkerPoolSize = workers
	config, err := json.Marshal(w.Config)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(w.ConfigPath, config, 0600); err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(w.Runtime, "bin")
	if err := os.Mkdir(bin, 0700); err != nil {
		t.Fatal(err)
	}
	h.cli = preparePackagedCandidate(t, bin, "gosched", "gosched")
	daemonBinary := preparePackagedCandidate(t, bin, "goschedd", "goschedd")
	daemon, err := clifixture.StartProcess(clifixture.Command{Executable: daemonBinary, Arguments: []string{"--config", w.ConfigPath}, Directory: w.Runtime, StdoutPath: filepath.Join(w.Evidence, "daemon-out.jsonl"), StderrPath: filepath.Join(w.Evidence, "daemon-err.txt")})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := daemon.Stop(); err != nil {
			h.cleanupReady = false
			t.Error(err)
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	t.Cleanup(cancel)
	h.ctx, h.api = ctx, client.New(w.IPCPath)
	if err := waitForV13Health(ctx, h.api); err != nil {
		t.Fatal(err)
	}
	if err := w.Bind(ctx, h.api, daemon.PID); err != nil {
		t.Fatal(err)
	}
	return h
}

func (h *ownedPackagedCLI) call(human bool, args ...string) []byte {
	h.t.Helper()
	if err := h.workspace.Guard(h.ctx, h.api); err != nil {
		h.t.Fatal(err)
	}
	h.sequence++
	argv := []string{"--ipc-path", h.workspace.IPCPath}
	if !human {
		argv = append(argv, "--json")
	}
	argv = append(argv, args...)
	out := filepath.Join(h.workspace.Evidence, fmt.Sprintf("%02d-out.txt", h.sequence))
	diagnostic := filepath.Join(h.workspace.Evidence, fmt.Sprintf("%02d-err.txt", h.sequence))
	p, err := clifixture.StartProcess(clifixture.Command{Executable: h.cli, Arguments: argv, Directory: h.workspace.Runtime, StdoutPath: out, StderrPath: diagnostic})
	if err != nil {
		h.t.Fatal(err)
	}
	defer func() {
		if err := p.Stop(); err != nil {
			h.cleanupReady = false
			h.t.Error(err)
		}
	}()
	if err := p.Wait(h.ctx); err != nil {
		data, readErr := os.ReadFile(diagnostic)
		h.t.Fatalf("CLI %v: %v; diagnostic=%s; diagnostic read=%v", args, err, data, readErr)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		h.t.Fatal(err)
	}
	return data
}
