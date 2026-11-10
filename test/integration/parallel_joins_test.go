package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/bundle"
	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

// TestS144OwnedTaskProcess is an explicitly selected, bounded native task fixture.
func TestS144OwnedTaskProcess(t *testing.T) {
	role := os.Getenv("GO_SCHEDULE_S144_ROLE")
	if role == "" {
		return
	}
	gate := os.Getenv("GO_SCHEDULE_S144_GATE")
	if role == "branch" {
		deadline := time.NewTimer(60 * time.Second)
		defer deadline.Stop()
		tick := time.NewTicker(10 * time.Millisecond)
		defer tick.Stop()
		for {
			if _, err := os.Stat(gate); err == nil {
				break
			}
			select {
			case <-deadline.C:
				t.Fatal("owned branch gate timed out")
			case <-tick.C:
			}
		}
	}
	fmt.Printf("owned parallel fixture %s\n", role)
}

func TestPackagedCLIParallelJoinsCorrelateNativeRuns(t *testing.T) {
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
			t.Logf("retained owned runtime: %s", w.Runtime)
		}
	})
	w.Config.WorkerPoolSize = 4
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
	seq := 0
	call := func(human bool, args ...string) []byte {
		t.Helper()
		if err := w.Guard(ctx, api); err != nil {
			t.Fatal(err)
		}
		seq++
		argv := []string{"--ipc-path", w.IPCPath}
		if !human {
			argv = append(argv, "--json")
		}
		argv = append(argv, args...)
		out, diagnostic := filepath.Join(w.Evidence, fmt.Sprintf("%02d-out.txt", seq)), filepath.Join(w.Evidence, fmt.Sprintf("%02d-err.txt", seq))
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
	task := func(name, role, gate string) string {
		t.Helper()
		created, err := api.CreateTask(ctx, server.TaskCreateRequest{Name: name, Command: os.Args[0], Args: []string{"-test.run=^TestS144OwnedTaskProcess$"}, Env: map[string]string{"GO_SCHEDULE_S144_ROLE": role, "GO_SCHEDULE_S144_GATE": gate}, WorkingDir: w.Runtime, OverlapPolicy: string(domain.OverlapAllowConcurrent)})
		if err != nil {
			t.Fatal(err)
		}
		return created.Task.ID
	}
	gateA, gateB := filepath.Join(w.Runtime, "tasks", "allow-a"), filepath.Join(w.Runtime, "tasks", "allow-b")
	entry, a, b, target := task("Prepare", "entry", ""), task("Export A", "branch", gateA), task("Export B", "branch", gateB), task("Publish", "target", "")
	branchesFile := filepath.Join(w.Runtime, "branches.json")
	create := func(branches []domain.ParallelJoinBranch) domain.ParallelJoin {
		t.Helper()
		data, err := json.Marshal(branches)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(branchesFile, data, 0600); err != nil {
			t.Fatal(err)
		}
		var join domain.ParallelJoin
		data = call(false, "join", "create", "--name=Owned exports", "--entry", entry, "--target", target, "--branches-file", branchesFile)
		if err := json.Unmarshal(data, &join); err != nil {
			t.Fatal(err)
		}
		for _, branch := range branches {
			if err := api.SetTaskEnabled(ctx, branch.TaskID, true); err != nil {
				t.Fatal(err)
			}
		}
		if err := api.SetTaskEnabled(ctx, target, true); err != nil {
			t.Fatal(err)
		}
		return join
	}
	instances := func(id string) []domain.ParallelJoinInstance {
		t.Helper()
		values, err := api.ListJoinInstances(ctx, domain.ParallelJoinInstanceFilter{DefinitionID: id, Limit: 100})
		if err != nil {
			t.Fatal(err)
		}
		return values
	}
	targets := func() []domain.Run {
		t.Helper()
		runs, err := api.ListRuns(ctx, target, 0)
		if err != nil {
			t.Fatal(err)
		}
		return runs
	}
	release := func(path string) {
		t.Helper()
		if err := os.WriteFile(path, []byte("owned release\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	join := create([]domain.ParallelJoinBranch{{ID: "a", Name: "Export A", TaskID: a}, {ID: "b", Name: "Export B", TaskID: b}})
	for range 2 {
		call(false, "task", "run-now", entry)
	}
	var initial []domain.ParallelJoinInstance
	waitS143(t, ctx, func() bool {
		initial = instances(join.ID)
		if len(initial) != 2 {
			return false
		}
		active, e := api.ListActiveRuns(ctx)
		return e == nil && len(active) == 4
	}, "two entry instances own four native branches")
	if len(targets()) != 0 {
		t.Fatal("target ran before branches completed")
	}
	release(gateA)
	waitS143(t, ctx, func() bool {
		for _, v := range instances(join.ID) {
			success := 0
			for _, b := range v.Branches {
				if b.Outcome == domain.OutcomeSuccess {
					success++
				}
			}
			if success != 1 || v.State != domain.JoinInstanceWaiting || v.TargetDeliveryID != "" {
				return false
			}
		}
		return true
	}, "one exact branch per instance remains outstanding")
	if len(targets()) != 0 {
		t.Fatal("one branch released target prematurely")
	}
	release(gateB)
	var complete []domain.ParallelJoinInstance
	waitS143(t, ctx, func() bool {
		complete = instances(join.ID)
		for _, v := range complete {
			if v.State != domain.JoinInstanceCompleted || v.TargetOutcome != domain.OutcomeSuccess {
				return false
			}
		}
		return len(complete) == 2 && len(targets()) == 2
	}, "each entry releases exactly one target")
	seen := map[string]bool{}
	for _, v := range complete {
		if seen[v.EntryRunID] || v.EntryRunID == "" || len(v.Branches) != 2 {
			t.Fatalf("entry correlation=%+v", v)
		}
		seen[v.EntryRunID] = true
		for _, b := range v.Branches {
			run, err := api.GetRun(ctx, b.RunID)
			if err != nil || run.JoinInstanceID != v.ID || run.JoinBranchID != b.ID || run.DeliveryID != b.DeliveryID || run.SourceRunID != v.EntryRunID || run.JoinRole != "branch" || run.Outcome != domain.OutcomeSuccess || !strings.Contains(run.Output, "owned parallel fixture branch") || seen[run.ID] {
				t.Fatalf("exact contribution=%+v run=%+v err=%v", b, run, err)
			}
			seen[run.ID] = true
		}
		var exact domain.ParallelJoinInstance
		if err := json.Unmarshal(call(false, "join", "instance", v.ID), &exact); err != nil || exact.ID != v.ID || exact.TargetRunID != v.TargetRunID {
			t.Fatalf("exact CLI=%+v err=%v", exact, err)
		}
		if text := string(call(true, "join", "instance", v.ID)); !strings.Contains(text, v.EntryRunID) || !strings.Contains(text, v.TargetRunID) {
			t.Fatalf("human exact evidence=%q", text)
		}
	}
	call(false, "join", "rm", join.ID)
	repeated := create([]domain.ParallelJoinBranch{{ID: "first", Name: "First export", TaskID: a}, {ID: "second", Name: "Second export", TaskID: a}})
	call(false, "task", "run-now", entry)
	var duplicate domain.ParallelJoinInstance
	waitS143(t, ctx, func() bool {
		v := instances(repeated.ID)
		if len(v) != 1 {
			return false
		}
		duplicate = v[0]
		return duplicate.State == domain.JoinInstanceCompleted && duplicate.TargetOutcome == domain.OutcomeSuccess
	}, "same-task branch slots finish independently")
	if len(duplicate.Branches) != 2 || duplicate.Branches[0].ID == duplicate.Branches[1].ID || duplicate.Branches[0].RunID == duplicate.Branches[1].RunID || len(targets()) != 3 {
		t.Fatalf("duplicate task slots=%+v", duplicate)
	}
	call(false, "join", "rm", repeated.ID)
	if err := os.Remove(gateA); err != nil {
		t.Fatal(err)
	}
	blocked := create([]domain.ParallelJoinBranch{{ID: "cancel", Name: "Cancel export", TaskID: a}, {ID: "sibling", Name: "Independent export", TaskID: b}})
	call(false, "task", "run-now", entry)
	var canceled domain.ParallelJoinInstance
	waitS143(t, ctx, func() bool {
		v := instances(blocked.ID)
		if len(v) != 1 {
			return false
		}
		canceled = v[0]
		for _, br := range canceled.Branches {
			if br.TaskID == a && br.RunID != "" && br.Outcome == "" {
				return true
			}
		}
		return false
	}, "exact cancelable branch is active")
	for _, br := range canceled.Branches {
		if br.TaskID == a {
			call(false, "runs", "cancel", br.RunID)
		}
	}
	waitS143(t, ctx, func() bool {
		v := instances(blocked.ID)[0]
		if v.State != domain.JoinInstanceBlocked {
			return false
		}
		for _, br := range v.Branches {
			if br.TaskID == a && br.Outcome != domain.OutcomeCanceled {
				return false
			}
			if br.TaskID == b && br.Outcome != domain.OutcomeSuccess {
				return false
			}
		}
		return true
	}, "canceled instance blocks while sibling retains its native success")
	if len(targets()) != 3 {
		t.Fatal("canceled prerequisite released a target")
	}
	if _, err := api.Health(ctx); err != nil {
		t.Fatal(err)
	}
	portableFile := filepath.Join(w.Runtime, "portable-joins.json")
	call(false, "bundle", "export", "--output", portableFile)
	portableBytes, err := os.ReadFile(portableFile)
	if err != nil {
		t.Fatal(err)
	}
	portable, err := bundle.Decode(portableBytes)
	if err != nil || portable.Schema != bundle.SchemaV3 || len(portable.Joins) != 1 || len(portable.Joins[0].Branches) != 2 {
		t.Fatalf("portable joins=%+v err=%v", portable.Joins, err)
	}
	if portable.Joins[0].EntryTaskID == entry || portable.Joins[0].TargetTaskID == target {
		t.Fatal("portable definitions retained native task identities")
	}
	for _, forbidden := range []string{"join_instance_id", "entry_run_id", "delivery_id", "target_run_id", "GO_SCHEDULE_S144_GATE"} {
		if strings.Contains(string(portableBytes), forbidden) {
			t.Fatalf("portable document exposed runtime field %s", forbidden)
		}
	}
	call(false, "bundle", "validate", portableFile)
	receipt := map[string]any{"packaged_cli": true, "overlapping_entries": 2, "unique_targets": 3, "duplicate_task_slots": true, "exact_native_correlation": true, "canceled_instance_blocked": true, "sibling_retained": true, "portable_v3_config_only": true}
	data, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(w.Evidence, "parallel-join-result.json"), append(data, '\n'), 0600); err != nil {
		t.Fatal(err)
	}
	t.Logf("retained owned parallel join evidence: %s", w.Evidence)
}
