package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/platform"
)

// TestS143OwnedTaskProcess runs only as an explicitly configured disposable task.
func TestS143OwnedTaskProcess(t *testing.T) {
	role := os.Getenv("GO_SCHEDULE_S143_ROLE")
	if role == "" {
		return
	}
	if role == "child" {
		timer := time.NewTimer(90 * time.Second)
		defer timer.Stop()
		<-timer.C
		return
	}
	root := os.Getenv("GO_SCHEDULE_S143_READY")
	if role != "parent" || root == "" {
		t.Fatal("invalid owned helper selection")
	}
	child := exec.Command(os.Args[0], "-test.run=^TestS143OwnedTaskProcess$")
	child.Env = append(os.Environ(), "GO_SCHEDULE_S143_ROLE=child")
	child.Stdout = os.Stdout
	child.Stderr = os.Stderr
	platform.HideConsole(child)
	if err := child.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = child.Process.Kill(); _ = child.Wait() }()
	parentPID, childPID := os.Getpid(), child.Process.Pid
	if _, err := fmt.Fprintf(os.Stdout, "owned cancellation fixture parent=%d child=%d\n", parentPID, childPID); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(map[string]int{"parent": parentPID, "child": childPID})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, fmt.Sprintf("%d.json", parentPID)), data, 0600); err != nil {
		t.Fatal(err)
	}
	timer := time.NewTimer(90 * time.Second)
	defer timer.Stop()
	<-timer.C
}

func waitS143(t *testing.T, ctx context.Context, predicate func() bool, description string) {
	t.Helper()
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	for {
		if predicate() {
			return
		}
		select {
		case <-ctx.Done():
			t.Fatalf("%s: %v", description, ctx.Err())
		case <-ticker.C:
		}
	}
}

func TestPackagedCLICancelsExactRunAndOwnedDescendants(t *testing.T) {
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
		} else {
			t.Logf("owned runtime retained after cleanup failure: %s", w.Runtime)
		}
	})
	w.Config.WorkerPoolSize = 2
	data, err := json.Marshal(w.Config)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(w.ConfigPath, data, 0600); err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(w.Runtime, "bin")
	if err := os.Mkdir(bin, 0700); err != nil {
		t.Fatal(err)
	}
	cliBinary := preparePackagedCandidate(t, bin, "gosched", "gosched")
	daemonBinary := preparePackagedCandidate(t, bin, "goschedd", "goschedd")
	daemon, err := clifixture.StartProcess(clifixture.Command{Executable: daemonBinary, Arguments: []string{"--config", w.ConfigPath}, Directory: w.Runtime, StdoutPath: filepath.Join(w.Evidence, "daemon-out.jsonl"), StderrPath: filepath.Join(w.Evidence, "daemon-err.txt")})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := daemon.Stop(); err != nil {
			cleanupReady = false
			t.Error(err)
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	api := client.New(w.IPCPath)
	if err := waitForV13Health(ctx, api); err != nil {
		t.Fatal(err)
	}
	if err := w.Bind(ctx, api, daemon.PID); err != nil {
		t.Fatal(err)
	}
	sequence := 0
	call := func(jsonOutput bool, args ...string) []byte {
		t.Helper()
		if err := w.Guard(ctx, api); err != nil {
			t.Fatal(err)
		}
		sequence++
		argv := []string{"--ipc-path", w.IPCPath}
		if jsonOutput {
			argv = append(argv, "--json")
		}
		argv = append(argv, args...)
		out := filepath.Join(w.Evidence, fmt.Sprintf("%02d-out.txt", sequence))
		diagnostic := filepath.Join(w.Evidence, fmt.Sprintf("%02d-err.txt", sequence))
		p, err := clifixture.StartProcess(clifixture.Command{Executable: cliBinary, Arguments: argv, Directory: w.Runtime, StdoutPath: out, StderrPath: diagnostic})
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := p.Stop(); err != nil {
				cleanupReady = false
				t.Error(err)
			}
		}()
		if err := p.Wait(ctx); err != nil {
			d, _ := os.ReadFile(diagnostic)
			t.Fatalf("CLI %v: %v %s", args, err, d)
		}
		data, err := os.ReadFile(out)
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	ready := filepath.Join(w.Runtime, "tasks", "cancel-ready")
	if err := os.Mkdir(ready, 0700); err != nil {
		t.Fatal(err)
	}
	created, err := api.CreateTask(ctx, server.TaskCreateRequest{Name: "owned cancellation fixture", Command: os.Args[0], Args: []string{"-test.run=^TestS143OwnedTaskProcess$"}, Env: map[string]string{"GO_SCHEDULE_S143_ROLE": "parent", "GO_SCHEDULE_S143_READY": ready}, WorkingDir: w.Runtime, OverlapPolicy: string(domain.OverlapAllowConcurrent)})
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		call(true, "task", "run-now", created.Task.ID)
	}
	var active []domain.Run
	waitS143(t, ctx, func() bool {
		entries, e := os.ReadDir(ready)
		if e != nil || len(entries) != 2 {
			return false
		}
		active, e = api.ListActiveRuns(ctx)
		return e == nil && len(active) == 2
	}, "two owned process trees become active")
	selected, sibling := active[0].ID, active[1].ID
	output := call(true, "runs", "cancel", selected)
	var ack domain.RunCancellation
	if err := json.Unmarshal(output, &ack); err != nil || ack.RunID != selected || ack.Status != "requested" {
		t.Fatalf("ack=%s err=%v", output, err)
	}
	var final domain.Run
	waitS143(t, ctx, func() bool { var e error; final, e = api.GetRun(ctx, selected); return e == nil }, "selected cancellation reaches persisted history")
	if final.Outcome != domain.OutcomeCanceled || final.ExitCode == nil || !strings.Contains(final.Output, "owned cancellation fixture") {
		t.Fatalf("canceled native result/output=%+v", final)
	}
	var parent, child int
	if _, err := fmt.Sscanf(strings.TrimSpace(final.Output), "owned cancellation fixture parent=%d child=%d", &parent, &child); err != nil || parent <= 0 || child <= 0 {
		t.Fatalf("owned process attribution output=%q err=%v", final.Output, err)
	}
	waitS143(t, ctx, func() bool { return s143ProcessExited(parent) && s143ProcessExited(child) }, "selected root and descendant exit independently")
	still, err := api.ListActiveRuns(ctx)
	if err != nil || len(still) != 1 || still[0].ID != sibling {
		t.Fatalf("sibling was disturbed: active=%+v err=%v", still, err)
	}
	for range 2 {
		raw := call(true, "runs", "cancel", selected)
		var again domain.RunCancellation
		if err := json.Unmarshal(raw, &again); err != nil || again.Status != "already_finished" || again.Outcome != domain.OutcomeCanceled {
			t.Fatalf("duplicate ack=%s err=%v", raw, err)
		}
	}
	human := string(call(false, "runs", "cancel", sibling))
	if !strings.Contains(strings.ToLower(human), "requested") || !strings.Contains(human, sibling) {
		t.Fatalf("human ack=%q", human)
	}
	waitS143(t, ctx, func() bool { runs, e := api.ListActiveRuns(ctx); return e == nil && len(runs) == 0 }, "remaining owned execution finishes")
	runs, err := api.ListRuns(ctx, created.Task.ID, 0)
	if err != nil || len(runs) != 2 {
		t.Fatalf("terminal history=%+v err=%v", runs, err)
	}
	for _, r := range runs {
		if r.Outcome != domain.OutcomeCanceled {
			t.Fatalf("terminal result=%+v", r)
		}
	}
	if _, err := api.Health(ctx); err != nil {
		t.Fatalf("daemon stopped after cancellation: %v", err)
	}
	if err := os.WriteFile(filepath.Join(w.Evidence, "cancellation-result.json"), []byte(`{"packaged_cli":true,"selected_tree_stopped":true,"sibling_untouched":true,"history_records":2,"duplicate_history":false}`), 0600); err != nil {
		t.Fatal(err)
	}
	t.Logf("retained owned cancellation evidence: %s", w.Evidence)
}
