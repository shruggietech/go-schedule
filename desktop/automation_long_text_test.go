package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/desktop/automation"
	"github.com/shruggietech/go-schedule/v2/desktop/connection"
	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

func TestJoinEditorTargetOnlyEditPreservesAcceptedUnicodeNames(t *testing.T) {
	for _, repetitions := range []int{700, 2000, 9000} {
		t.Run(fmt.Sprintf("unicode_bytes_%d", repetitions*9), func(t *testing.T) {
			st, err := store.Open(":memory:")
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = st.Close() })
			api := server.New(st, nil, nil, nil, "", slog.New(slog.DiscardHandler))
			var patches [][]byte
			c := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPatch && strings.HasPrefix(r.URL.Path, "/v1/joins/") {
					raw, e := io.ReadAll(r.Body)
					if e != nil {
						t.Error(e)
					}
					patches = append(patches, raw)
					r.Body = io.NopCloser(bytes.NewReader(raw))
				}
				api.Handler().ServeHTTP(w, r)
			}))
			ctx := context.Background()
			tasks := make([]domain.Task, 4)
			for i := range tasks {
				v, e := c.CreateTask(ctx, server.TaskCreateRequest{Name: fmt.Sprintf("Task %d", i), Command: "echo", Schedule: "every hour", Timezone: "UTC"})
				if e != nil {
					t.Fatal(e)
				}
				tasks[i] = v.Task
			}
			longName := strings.Repeat("長λ🙂", repetitions)
			branches := make([]domain.ParallelJoinBranch, 101)
			for i := range branches {
				branches[i] = domain.ParallelJoinBranch{ID: fmt.Sprintf("slot-%03d", i), Name: fmt.Sprintf("Branch %03d", i), TaskID: tasks[1].ID}
			}
			branches[100].Name = "Last " + longName
			created := domain.ParallelJoin{Name: "Join " + longName, EntryTaskID: tasks[0].ID, TargetTaskID: tasks[2].ID, Branches: branches}
			if repetitions > 2000 {
				if err = st.CreateParallelJoin(&created); err != nil {
					t.Fatal(err)
				}
			} else {
				created, err = c.CreateJoin(ctx, server.JoinCreateRequest{Name: created.Name, EntryTaskID: created.EntryTaskID, TargetTaskID: created.TargetTaskID, Branches: branches})
				if err != nil {
					t.Fatal(err)
				}
			}
			manifest, e := c.Manifest(ctx)
			if e != nil {
				t.Fatal(e)
			}
			ready := make(chan struct{}, 1)
			app := newApp(authoringAppHealth{daemonID: manifest.InstallationID}, cancellationAppEmitter{ready}, nil, appServices{automation: automation.NewService(automation.NewLocalBackend(c))})
			appCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			app.startup(appCtx)
			defer app.shutdown(context.Background())
			select {
			case <-ready:
			case <-time.After(time.Second):
				t.Fatal("authoring daemon did not become current")
			}
			expected := manifest.InstallationID
			workspace := app.AutomationWorkspace(expected)
			if workspace.Outcome != "accepted" || workspace.Workspace == nil || len(workspace.Workspace.Joins) != 1 {
				t.Fatalf("workspace outcome=%s", workspace.Outcome)
			}
			bounded := workspace.Workspace.Joins[0]
			if bounded.Name == created.Name || bounded.Branches[100].Name == branches[100].Name {
				t.Fatal("workspace observation unexpectedly returned full authoring text")
			}
			authoring := app.JoinDefinition(created.ID, expected)
			if authoring.Outcome != "accepted" || authoring.Join == nil {
				t.Fatalf("authoring read outcome=%s message=%s", authoring.Outcome, authoring.Message)
			}
			observed := *authoring.Join
			// The editor opens only after the full authoring read, including all paged slots.
			draft := automation.JoinDraft{ExpectedDaemonID: expected, ID: observed.ID, Name: observed.Name, EntryTaskID: observed.EntryTaskID, TargetTaskID: tasks[3].ID, Branches: observed.Branches, OriginalUpdatedAt: observed.UpdatedAt.Format(time.RFC3339Nano)}
			saved := app.SaveJoin(draft)
			if saved.Outcome != "accepted" {
				t.Fatalf("save outcome=%s message=%s", saved.Outcome, saved.Message)
			}
			if len(patches) != 1 {
				t.Fatalf("target-only patch count=%d", len(patches))
			}
			var targetPatch map[string]json.RawMessage
			if err = json.Unmarshal(patches[0], &targetPatch); err != nil {
				t.Fatal(err)
			}
			if len(targetPatch) != 2 || targetPatch["target_task_id"] == nil || targetPatch["expected_updated_at"] == nil {
				t.Fatalf("target-only patch included unchanged fields: %v", mapKeys(targetPatch))
			}
			persisted, err := st.GetParallelJoin(created.ID)
			if err != nil {
				t.Fatal(err)
			}
			if persisted.Name != created.Name {
				t.Errorf("target-only edit changed join name: got %d bytes want %d", len(persisted.Name), len(created.Name))
			}
			if len(persisted.Branches) != 101 || persisted.Branches[100].Name != branches[100].Name {
				t.Errorf("target-only edit changed paged branch name: got %d bytes want %d", len(persisted.Branches[100].Name), len(branches[100].Name))
			}
			if persisted.TargetTaskID != tasks[3].ID {
				t.Error("target edit was not applied")
			}
			draft.TargetTaskID = tasks[2].ID
			stale := app.SaveJoin(draft)
			if stale.Outcome != "stale" {
				t.Fatalf("old editor revision outcome=%s", stale.Outcome)
			}
			afterStale, err := st.GetParallelJoin(created.ID)
			if err != nil {
				t.Fatal(err)
			}
			if afterStale.TargetTaskID != tasks[3].ID {
				t.Error("stale editor retargeted the definition")
			}
			if len(patches) != 1 {
				t.Error("stale preflight submitted a patch")
			}
			replacementRead := app.JoinDefinition(created.ID, expected)
			if replacementRead.Outcome != "accepted" || replacementRead.Join == nil {
				t.Fatalf("replacement read=%s", replacementRead.Outcome)
			}
			replacement := *replacementRead.Join
			replacement.Branches[0].Name = "Renamed first branch"
			replaced := app.SaveJoin(automation.JoinDraft{ExpectedDaemonID: expected, ID: replacement.ID, Name: replacement.Name, EntryTaskID: replacement.EntryTaskID, TargetTaskID: replacement.TargetTaskID, Branches: replacement.Branches, OriginalUpdatedAt: replacement.UpdatedAt.Format(time.RFC3339Nano)})
			if replaced.Outcome != "accepted" {
				t.Fatalf("replacement outcome=%s", replaced.Outcome)
			}
			final, err := st.GetParallelJoin(created.ID)
			if err != nil {
				t.Fatal(err)
			}
			if final.Name != created.Name || len(final.Branches) != 101 || final.Branches[0].Name != "Renamed first branch" || final.Branches[100].Name != branches[100].Name {
				t.Error("branch replacement changed untouched full names or branch identities")
			}
			var branchPatch map[string]json.RawMessage
			if len(patches) != 2 {
				t.Fatalf("replacement patch count=%d", len(patches))
			}
			if err = json.Unmarshal(patches[1], &branchPatch); err != nil {
				t.Fatal(err)
			}
			if len(branchPatch) != 2 || branchPatch["branches"] == nil || branchPatch["expected_updated_at"] == nil {
				t.Fatalf("branch replacement included unchanged fields: %v", mapKeys(branchPatch))
			}
			noOp := app.SaveJoin(automation.JoinDraft{ExpectedDaemonID: expected, ID: final.ID, Name: final.Name, EntryTaskID: final.EntryTaskID, TargetTaskID: final.TargetTaskID, Branches: final.Branches, OriginalUpdatedAt: final.UpdatedAt.Format(time.RFC3339Nano)})
			if noOp.Outcome != "accepted" || !strings.Contains(noOp.Message, "No changes") || len(patches) != 2 {
				t.Fatalf("no-op outcome=%s patch count=%d", noOp.Outcome, len(patches))
			}
		})
	}
}

func mapKeys(v map[string]json.RawMessage) []string {
	out := make([]string, 0, len(v))
	for k := range v {
		out = append(out, k)
	}
	return out
}

type authoringAppHealth struct {
	appBackend
	daemonID string
}

func (b authoringAppHealth) Health(context.Context) (connection.Health, error) {
	return connection.Health{ID: b.daemonID, DisplayName: "Owned authoring daemon", Version: "v2.0.3", Permissions: []string{"manage"}, Capabilities: []string{"parallel-joins"}}, nil
}
