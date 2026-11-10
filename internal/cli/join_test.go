package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestJoinCommandsExposeCompleteLifecycle(t *testing.T) {
	r := newJoinCmd()
	for _, name := range []string{"create", "list", "show", "update", "rm", "instances", "instance"} {
		c, _, e := r.Find([]string{name})
		if e != nil || c == r {
			t.Errorf("missing %s", name)
		}
	}
}
func TestJoinCreateRequiresReferencesBeforeContactingDaemon(t *testing.T) {
	c := joinCreate()
	c.SetArgs([]string{"--name", "Example", "--entry", "entry", "--target", "target"})
	if err := c.Execute(); !errors.Is(err, errUsage) || !strings.Contains(err.Error(), "branches-file") {
		t.Fatalf("err=%v", err)
	}
}

func TestJoinBranchesFileAndStdinPreserveDistinctSameTaskInvocations(t *testing.T) {
	data := `[{"id":"a","name":"Export A","task_id":"same"},{"id":"b","name":"Export B","task_id":"same"}]`
	file := filepath.Join(t.TempDir(), "branches.json")
	if err := os.WriteFile(file, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	for _, source := range []string{"-", file} {
		body, p, calls, err := executeAuthoring(t, []string{"join", "create", "--name", "Exports", "--entry", "entry", "--target", "target", "--branches-file", source, "--json"}, data)
		if err != nil || p != "/v1/joins" || calls != 1 {
			t.Fatalf("path=%s calls=%d err=%v", p, calls, err)
		}
		branches := body["branches"].([]any)
		if len(branches) != 2 || branches[0].(map[string]any)["task_id"] != branches[1].(map[string]any)["task_id"] || branches[0].(map[string]any)["id"] == branches[1].(map[string]any)["id"] {
			t.Fatalf("branches=%+v", branches)
		}
	}
}
func TestJoinBranchesRejectReadOnlyMetadataAndDuplicateIdentity(t *testing.T) {
	for _, input := range []string{`[{"id":"a","name":"A","task_id":"task","task_name":"readonly"}]`, `[{"id":"same","name":"A","task_id":"task"},{"id":"same","name":"B","task_id":"task"}]`} {
		_, _, calls, err := executeAuthoring(t, []string{"join", "create", "--name", "Exports", "--entry", "entry", "--target", "target", "--branches-file", "-"}, input)
		if err == nil || calls != 0 {
			t.Fatalf("err=%v calls=%d", err, calls)
		}
	}
}

func TestJoinShowReportsBoundedHumanTextWithoutChangingJSON(t *testing.T) {
	for _, branchOnly := range []bool{false, true} {
		t.Run(fmt.Sprintf("branch_only_%v", branchOnly), func(t *testing.T) {
			v := domain.ParallelJoin{ID: "join", Name: strings.Repeat("長", 1365), EntryTaskID: "entry", TargetTaskID: "target", BranchCount: 1, Branches: []domain.ParallelJoinBranch{{ID: "slot", Name: "Branch", TaskID: "task"}}}
			if branchOnly {
				v.Branches[0].TruncatedFields = []string{"name"}
			} else {
				v.TruncatedFields = []string{"name"}
			}
			c := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v1/joins/join" {
					t.Errorf("observation used authoring path: %s", r.URL.Path)
				}
				if err := json.NewEncoder(w).Encode(v); err != nil {
					t.Error(err)
				}
			}))
			for _, asJSON := range []bool{false, true} {
				var out, diagnostic bytes.Buffer
				args := []string{"join", "show", "join"}
				if asJSON {
					args = append(args, "--json")
				}
				root := newRootWithDependencies(Dependencies{Client: c, Out: &out, Err: &diagnostic})
				root.SetArgs(args)
				if err := root.Execute(); err != nil {
					t.Fatal(err)
				}
				if asJSON {
					var result domain.ParallelJoin
					if err := json.Unmarshal(out.Bytes(), &result); err != nil {
						t.Fatalf("JSON output gained human prose: %v", err)
					}
					if !joinObservationTruncated(result) {
						t.Error("JSON omitted truncation metadata")
					}
				} else if !strings.Contains(out.String(), "some display text is truncated") {
					t.Error("human observation omitted actual truncation notice")
				}
			}
		})
	}
}
