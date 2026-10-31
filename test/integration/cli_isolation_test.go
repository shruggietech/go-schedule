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
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
)

// TestCLIIsolationBenignTask is the disposable task program, not a host action.
func TestCLIIsolationBenignTask(t *testing.T) {
	path := os.Getenv("GO_SCHEDULE_CLI_BENIGN_OUTPUT")
	if path == "" {
		return
	}
	if err := os.WriteFile(path, []byte("owned CLI task executed\n"), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestPackagedCLIUsesOnlyOwnedDaemon(t *testing.T) {
	w, err := clifixture.NewWorkspace("")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("retained CLI evidence: %s", w.Evidence)
	cleanupReady := true
	t.Cleanup(func() {
		if !cleanupReady {
			t.Logf("retained runtime after process cleanup failure: %s", w.Runtime)
			return
		}
		if err := w.CleanupRuntime(); err != nil {
			t.Errorf("owned cleanup: %v", err)
		}
	})
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	binDir := filepath.Join(w.Runtime, "bin")
	if err := os.Mkdir(binDir, 0700); err != nil {
		t.Fatal(err)
	}
	cliBinary := preparePackagedCandidate(t, binDir, "gosched", "gosched")
	daemonBinary := preparePackagedCandidate(t, binDir, "goschedd", "goschedd")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	sequence := 0
	run := func(binary string, args []string) []byte {
		t.Helper()
		sequence++
		prefix := fmt.Sprintf("%02d", sequence)
		stdout, stderr := filepath.Join(w.Evidence, prefix+"-stdout.txt"), filepath.Join(w.Evidence, prefix+"-stderr.txt")
		p, err := clifixture.StartProcess(clifixture.Command{Executable: binary, Arguments: args, Directory: root, StdoutPath: stdout, StderrPath: stderr})
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := p.Stop(); err != nil {
				cleanupReady = false
				t.Errorf("owned command process cleanup: %v", err)
			}
		}()
		if err := p.Wait(ctx); err != nil {
			diagnostics, _ := os.ReadFile(stderr)
			t.Fatalf("owned command %q failed: %v\n%s", args, err, diagnostics)
		}
		output, err := os.ReadFile(stdout)
		if err != nil {
			t.Fatal(err)
		}
		return output
	}
	daemon, err := clifixture.StartProcess(clifixture.Command{Executable: daemonBinary, Arguments: []string{"--config", w.ConfigPath}, Directory: w.Runtime, StdoutPath: filepath.Join(w.Evidence, "daemon-stdout.jsonl"), StderrPath: filepath.Join(w.Evidence, "daemon-stderr.txt")})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := daemon.Stop(); err != nil {
			cleanupReady = false
			t.Errorf("owned daemon process-tree cleanup: %v", err)
		}
		stopCtx, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		_ = daemon.Wait(stopCtx)
	})
	api := client.New(w.IPCPath)
	if err := waitForV13Health(ctx, api); err != nil {
		t.Fatal(err)
	}
	if err := w.Bind(ctx, api, daemon.PID); err != nil {
		t.Fatal(err)
	}
	call := func(mutation bool, args ...string) []byte {
		t.Helper()
		if mutation {
			if err := w.Guard(ctx, api); err != nil {
				t.Fatal(err)
			}
		}
		return run(cliBinary, append([]string{"--ipc-path", w.IPCPath, "--json"}, args...))
	}
	call(false, "health")
	outputFile := filepath.Join(w.Runtime, "tasks", "benign-output.txt")
	if err := w.Contains(outputFile); err != nil {
		t.Fatal(err)
	}
	createdJSON := call(true, "task", "add", "owned benign CLI fixture", "--command", os.Args[0], "--arg=-test.run=^TestCLIIsolationBenignTask$", "--env=GO_SCHEDULE_CLI_BENIGN_OUTPUT="+outputFile)
	var created server.TaskResponse
	if err := json.Unmarshal(createdJSON, &created); err != nil || created.Task.ID == "" {
		t.Fatalf("task response=%s err=%v", createdJSON, err)
	}
	call(true, "task", "run-now", created.Task.ID)
	deadline := time.Now().Add(15 * time.Second)
	for {
		if _, err := os.Stat(outputFile); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("owned benign task did not run")
		}
		time.Sleep(30 * time.Millisecond)
	}
	call(false, "task", "show", created.Task.ID)
	call(false, "runs", "--task", created.Task.ID)
	call(true, "task", "rm", created.Task.ID)
	if err := w.Guard(ctx, client.New(w.IPCPath+"-unowned")); err == nil {
		t.Fatal("mismatched target permitted mutation")
	}
	originalID := w.DaemonID
	w.DaemonID = "unowned-installation"
	if err := w.Guard(ctx, api); err == nil {
		t.Fatal("mismatched daemon identity permitted mutation")
	}
	w.DaemonID = originalID
	data, err := json.MarshalIndent(map[string]any{"evidence_kind": clifixture.PackagedCLI, "run_id": w.ID, "daemon_id": w.DaemonID, "daemon_pid": daemon.PID, "ipc_path": w.IPCPath, "task_id": created.Task.ID, "operations": []string{"health", "task.add", "task.run", "task.show", "runs", "task.rm"}}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(w.Evidence, "packaged-cli.json"), append(data, '\n'), 0600); err != nil {
		t.Fatal(err)
	}
}
