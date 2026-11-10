package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/bundle"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

// TestS145OwnedTaskProcess is selected only by an owned disposable native task.
func TestS145OwnedTaskProcess(t *testing.T) {
	role := os.Getenv("GO_SCHEDULE_S145_ROLE")
	if role == "" {
		return
	}
	fmt.Printf("owned recovery fixture started role=%s\n", role)
	ready := filepath.Join(os.Getenv("GO_SCHEDULE_S145_READY"), fmt.Sprintf("%s-%d.ready", role, os.Getpid()))
	if err := os.WriteFile(ready, []byte("native fixture ready\n"), 0600); err != nil {
		t.Fatal(err)
	}
	gate := os.Getenv("GO_SCHEDULE_S145_GATE")
	timer := time.NewTimer(60 * time.Second)
	defer timer.Stop()
	tick := time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()
	for gate != "" {
		if _, err := os.Stat(gate); err == nil {
			break
		}
		select {
		case <-timer.C:
			t.Fatal("owned native fixture gate timed out")
		case <-tick.C:
		}
	}
	fmt.Printf("owned recovery fixture finished role=%s\n", role)
	if os.Getenv("GO_SCHEDULE_S145_RESULT") == "failure" {
		t.Fatal("expected synthetic native task failure")
	}
}

func TestPackagedParallelRecovery(t *testing.T) {
	h := newOwnedPackagedCLI(t, 8)
	w, api, ctx, call := h.workspace, h.api, h.ctx, h.call
	cases := []struct {
		name            string
		cause           string
		entries         int
		repeated        bool
		reuseCleanup    bool
		reuseRecovery   bool
		overlapSkip     bool
		recoveryFails   bool
		disableRecovery bool
		emptyConditions bool
	}{
		{name: "overlapping_success_and_finite_role_reuse", entries: 2, repeated: true, reuseCleanup: true},
		{name: "failure_waits_for_sibling_and_failed_recovery", cause: "failure", entries: 1, recoveryFails: true},
		{name: "cancel_waits_for_sibling", cause: "canceled", entries: 1},
		{name: "no_start_followup_wakes_cleanup", cause: "failure", entries: 1, disableRecovery: true},
		{name: "explicit_empty_policy_still_cleans_up", cause: "failure", entries: 1, emptyConditions: true},
		{name: "success_cleanup_reuse_with_skip", entries: 1, reuseCleanup: true, overlapSkip: true},
		{name: "recovery_cleanup_reuse_with_skip", cause: "failure", entries: 1, reuseRecovery: true, overlapSkip: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := filepath.Join(w.Runtime, "tasks", tc.name)
			ready := filepath.Join(root, "ready")
			if err := os.MkdirAll(ready, 0700); err != nil {
				t.Fatal(err)
			}
			gateA, gateB, gateFollow := filepath.Join(root, "a"), filepath.Join(root, "b"), filepath.Join(root, "follow")
			release := func(path string) {
				t.Helper()
				if err := os.WriteFile(path, []byte("owned release\n"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			task := func(label, role, result, gate string) string {
				t.Helper()
				policy := domain.OverlapAllowConcurrent
				if tc.overlapSkip && (role == "publish" || role == "recover") {
					policy = domain.OverlapSkip
				}
				v, err := api.CreateTask(ctx, server.TaskCreateRequest{Name: tc.name + " " + label, Command: os.Args[0], Args: []string{"-test.run=^TestS145OwnedTaskProcess$"}, Env: map[string]string{"GO_SCHEDULE_S145_ROLE": role, "GO_SCHEDULE_S145_RESULT": result, "GO_SCHEDULE_S145_GATE": gate, "GO_SCHEDULE_S145_READY": ready}, WorkingDir: w.Runtime, OverlapPolicy: string(policy)})
				if err != nil {
					t.Fatal(err)
				}
				return v.Task.ID
			}
			entry := task("Prepare", "entry", "success", "")
			aResult := "success"
			if tc.cause == "failure" {
				aResult = "failure"
			}
			a := task("Export A", "a", aResult, gateA)
			b := task("Export B", "b", "success", gateB)
			publish := task("Publish", "publish", "success", gateFollow)
			recoverResult := "success"
			if tc.recoveryFails {
				recoverResult = "failure"
			}
			recovery := task("Recover", "recover", recoverResult, gateFollow)
			cleanup := task("Cleanup", "cleanup", "success", "")
			if tc.reuseCleanup {
				cleanup = publish
			}
			if tc.reuseRecovery {
				cleanup = recovery
			}
			branches := []domain.ParallelJoinBranch{{ID: "a", Name: "Export A", TaskID: a}, {ID: "b", Name: "Export B", TaskID: b}}
			if tc.repeated {
				branches[1].TaskID = a
			}
			branchBytes, err := json.Marshal(branches)
			if err != nil {
				t.Fatal(err)
			}
			branchFile := filepath.Join(root, "branches.json")
			if err := os.WriteFile(branchFile, branchBytes, 0600); err != nil {
				t.Fatal(err)
			}
			conditions := "failure,canceled,skipped,not_started"
			if tc.emptyConditions {
				conditions = "none"
			}
			var definition domain.ParallelJoin
			raw := call(false, "join", "create", "--name", tc.name, "--entry", entry, "--target", publish, "--branches-file", branchFile, "--recovery", recovery, "--recovery-on", conditions, "--cleanup", cleanup)
			if err := json.Unmarshal(raw, &definition); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := api.DeleteJoin(ctx, definition.ID); err != nil {
					t.Error(err)
				}
			}()
			destinations := []string{a, publish, recovery, cleanup}
			if !tc.repeated {
				destinations = append(destinations, b)
			}
			for _, id := range destinations {
				if err := api.SetTaskEnabled(ctx, id, true); err != nil {
					t.Fatal(err)
				}
			}
			if tc.emptyConditions && (definition.RecoveryOn == nil || len(definition.RecoveryOn) != 0) {
				t.Fatalf("explicit empty policy lost: %+v", definition)
			}
			instances := func() []domain.ParallelJoinInstance {
				t.Helper()
				v, err := api.ListJoinInstances(ctx, domain.ParallelJoinInstanceFilter{DefinitionID: definition.ID, Limit: 100})
				if err != nil {
					t.Fatal(err)
				}
				return v
			}
			for range tc.entries {
				call(false, "task", "run-now", entry)
			}
			waitS143(t, ctx, func() bool {
				v := instances()
				if len(v) != tc.entries {
					return false
				}
				for _, inst := range v {
					for _, br := range inst.Branches {
						if br.RunID == "" || br.Outcome != "" {
							return false
						}
					}
				}
				names, err := filepath.Glob(filepath.Join(ready, "a-*.ready"))
				if err != nil {
					t.Fatal(err)
				}
				want := tc.entries
				if tc.repeated {
					want *= 2
				}
				return len(names) == want
			}, "all exact native branches are owned and active")
			if tc.cause == "canceled" {
				v := instances()[0]
				call(false, "runs", "cancel", v.Branches[0].RunID)
			} else {
				release(gateA)
			}
			if !tc.repeated {
				waitS143(t, ctx, func() bool {
					v := instances()[0]
					for _, br := range v.Branches {
						if br.BranchID == "a" {
							return br.Settlement != ""
						}
					}
					return false
				}, "first exact branch has actual terminal settlement")
				for _, v := range instances() {
					if v.SelectedRoute != "" || v.TargetDeliveryID != "" || v.RecoveryDeliveryID != "" || v.CleanupDeliveryID != "" {
						t.Fatalf("premature route while sibling active: %+v", v)
					}
				}
				if tc.disableRecovery {
					if err := api.SetTaskEnabled(ctx, recovery, false); err != nil {
						t.Fatal(err)
					}
				}
				release(gateB)
			}
			if !tc.disableRecovery && !tc.emptyConditions {
				waitS143(t, ctx, func() bool {
					for _, v := range instances() {
						if v.CleanupDeliveryID != "" {
							t.Fatalf("cleanup raced selected native follow-up: %+v", v)
						}
						if v.TargetRunID == "" && v.RecoveryRunID == "" {
							return false
						}
					}
					return true
				}, "selected exact follow-up is active before cleanup")
				release(gateFollow)
			}
			var complete []domain.ParallelJoinInstance
			waitS143(t, ctx, func() bool {
				complete = instances()
				if len(complete) != tc.entries {
					return false
				}
				for _, v := range complete {
					if v.State != domain.JoinInstanceCompleted || v.CleanupOutcome != domain.OutcomeSuccess {
						return false
					}
				}
				return true
			}, "all configured cleanup stages finish")
			seen := map[string]bool{}
			for _, v := range complete {
				want := domain.JoinRouteSuccess
				if tc.cause != "" {
					want = domain.JoinRouteRecovery
				}
				if tc.emptyConditions {
					want = domain.JoinRouteNone
				}
				if v.SelectedRoute != want {
					t.Fatalf("route=%s want=%s instance=%+v", v.SelectedRoute, want, v)
				}
				if v.EntryRunID == "" || seen[v.EntryRunID] {
					t.Fatalf("entry identity reused: %+v", v)
				}
				seen[v.EntryRunID] = true
				for _, br := range v.Branches {
					run, err := api.GetRun(ctx, br.RunID)
					if err != nil || run.JoinInstanceID != v.ID || run.DeliveryID != br.DeliveryID || run.JoinBranchID != br.ID || run.SourceRunID != v.EntryRunID || run.JoinRole != "branch" || seen[run.ID] {
						t.Fatalf("exact branch=%+v run=%+v err=%v", br, run, err)
					}
					seen[run.ID] = true
					if br.Settlement == "" || br.SettledAt == nil {
						t.Fatalf("missing actual branch settlement: %+v", br)
					}
					if tc.cause == "canceled" && br.BranchID == "a" && br.Outcome != domain.OutcomeCanceled {
						t.Fatalf("cancellation outcome lost: %+v", br)
					}
				}
				if want != domain.JoinRouteSuccess && v.TargetDeliveryID != "" {
					t.Fatalf("unsuccessful aggregate admitted Publish: %+v", v)
				}
				if want != domain.JoinRouteRecovery && v.RecoveryDeliveryID != "" {
					t.Fatalf("unselected recovery admitted: %+v", v)
				}
				if tc.disableRecovery && (v.RecoveryRunID != "" || v.RecoverySettlement != domain.JoinSettlementNotStarted || v.RecoveryReason == "") {
					t.Fatalf("no-start evidence fabricated/lost: %+v", v)
				}
				if tc.recoveryFails && v.RecoveryOutcome != domain.OutcomeFailure {
					t.Fatalf("actual recovery failure lost: %+v", v)
				}
				for _, stage := range []struct{ role, id, delivery string }{{"target", v.TargetRunID, v.TargetDeliveryID}, {"recovery", v.RecoveryRunID, v.RecoveryDeliveryID}, {"cleanup", v.CleanupRunID, v.CleanupDeliveryID}} {
					if stage.id == "" {
						continue
					}
					run, err := api.GetRun(ctx, stage.id)
					if err != nil || run.JoinInstanceID != v.ID || run.DeliveryID != stage.delivery || run.SourceRunID != v.EntryRunID || run.JoinRole != stage.role || seen[run.ID] {
						t.Fatalf("exact stage=%+v run=%+v err=%v", stage, run, err)
					}
					seen[run.ID] = true
					if stage.role != "cleanup" && run.EndedAt != nil && v.CleanupStartedAt != nil && v.CleanupStartedAt.Before(*run.EndedAt) {
						t.Fatalf("cleanup preceded native follow-up end: %+v", v)
					}
				}
				var exact domain.ParallelJoinInstance
				if err := json.Unmarshal(call(false, "join", "instance", v.ID), &exact); err != nil || exact.CleanupRunID != v.CleanupRunID || exact.SelectedRoute != want {
					t.Fatalf("JSON exact=%+v err=%v", exact, err)
				}
				human := string(call(true, "join", "instance", v.ID))
				if !strings.Contains(human, v.CleanupRunID) {
					t.Fatalf("human inspection omitted cleanup run: %s", human)
				}
			}
			active, err := api.ListActiveRuns(ctx)
			if err != nil || len(active) != 0 {
				t.Fatalf("owned work not drained: %+v %v", active, err)
			}
			if tc.emptyConditions {
				file := filepath.Join(root, "routing-bundle.json")
				call(false, "bundle", "export", "--output", file)
				data, err := os.ReadFile(file)
				if err != nil {
					t.Fatal(err)
				}
				doc, err := bundle.Decode(data)
				if err != nil || doc.Schema != bundle.SchemaV4 || len(doc.Joins) != 1 {
					t.Fatalf("portable routing=%+v err=%v", doc, err)
				}
				var decoded map[string]any
				if err := json.Unmarshal(data, &decoded); err != nil {
					t.Fatal(err)
				}
				joins, ok := decoded["joins"].([]any)
				if !ok || len(joins) != 1 {
					t.Fatalf("portable join envelope=%+v", decoded["joins"])
				}
				config, ok := joins[0].(map[string]any)
				if !ok {
					t.Fatalf("portable join record=%+v", joins[0])
				}
				causes, ok := config["recovery_on"].([]any)
				if !ok || len(causes) != 0 {
					t.Fatalf("portable explicit[] lost: %+v", config)
				}
				for _, key := range []string{"entry_task_portable_id", "target_task_portable_id", "recovery_task_portable_id", "cleanup_task_portable_id"} {
					if config[key] == "" || config[key] == entry || config[key] == publish || config[key] == recovery || config[key] == cleanup {
						t.Fatalf("unmapped portable %s=%v", key, config[key])
					}
				}
				for _, forbidden := range []string{"selected_route", "entry_run_id", "delivery_id", "cleanup_run_id", "GO_SCHEDULE_S145_GATE"} {
					if strings.Contains(string(data), forbidden) {
						t.Fatalf("portable runtime field %s", forbidden)
					}
				}
				call(false, "bundle", "validate", file)
			}
		})
	}
	if _, err := api.Health(ctx); err != nil {
		t.Fatal(err)
	}
	receipt := map[string]any{"packaged_cli": true, "scenarios": len(cases), "overlapping_entries": 2, "finite_role_reuse": true, "success_cleanup_skip_reuse": true, "recovery_cleanup_skip_reuse": true, "repeated_task_slots": true, "failure_waits_for_sibling": true, "cancellation_waits_for_sibling": true, "failed_recovery_then_cleanup": true, "no_start_followup_wakes_cleanup": true, "explicit_empty_policy": true, "portable_v4_config_only": true}
	data, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(w.Evidence, "parallel-recovery-result.json"), append(data, '\n'), 0600); err != nil {
		t.Fatal(err)
	}
	t.Logf("retained owned parallel recovery evidence: %s", w.Evidence)
}
