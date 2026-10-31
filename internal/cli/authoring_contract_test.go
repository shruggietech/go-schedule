package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
)

func executeAuthoring(t *testing.T, argv []string, input string) (map[string]any, string, int, error) {
	t.Helper()
	var body map[string]any
	var path string
	var out, stderr bytes.Buffer
	calls := 0
	c := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		path = r.URL.Path
		if r.Body != nil {
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil && err != io.EOF {
				t.Errorf("decode request: %v", err)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		var response any = server.TaskResponse{}
		if r.URL.Path == "/v1/schedules/preview" {
			response = server.PreviewResponse{HumanSummary: "Every weekday", SourceSyntax: "human", PolicySummary: "wall clock"}
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	root := newRootWithDependencies(Dependencies{Client: c, In: strings.NewReader(input), Out: &out, Err: &stderr})
	root.SetArgs(argv)
	err := root.Execute()
	return body, path, calls, err
}

func TestTaskCreationTransportsCompleteIntent(t *testing.T) {
	body, path, calls, err := executeAuthoring(t, []string{"task", "add", "fixture", "--command", "echo", "--arg", "", "--env", "EXAMPLE=value", "--cwd", "fixture", "--stdin-file", "-", "--run-as", "fixture-user", "--disabled", "--schedule", "0 9 * * 1-5", "--syntax", "cron", "--json"}, "Unicode input: \u03bb\n")
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 || path != "/v1/tasks" || body["enabled"] != false || body["stdin"] != "Unicode input: \u03bb\n" || body["run_as"] != "fixture-user" || body["schedule_syntax"] != "cron" {
		t.Fatalf("path=%s calls=%d body=%v", path, calls, body)
	}
	if !reflect.DeepEqual(body["args"], []any{""}) {
		t.Fatalf("one empty argument lost: %v", body["args"])
	}
}

func TestTaskEditingPreservesExplicitClearsAndOmission(t *testing.T) {
	body, _, _, err := executeAuthoring(t, []string{"task", "edit", "fixture", "--clear-args", "--clear-env", "--clear-cwd", "--clear-run-as", "--clear-stdin", "--json"}, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"clear_env", "clear_working_dir", "clear_run_as"} {
		if body[key] != true {
			t.Errorf("missing %s: %v", key, body)
		}
	}
	if body["stdin"] != "" || !reflect.DeepEqual(body["args"], []any{}) {
		t.Fatalf("clear stdin/args lost: %v", body)
	}
	body, _, _, err = executeAuthoring(t, []string{"task", "edit", "fixture", "--name", "renamed", "--json"}, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"clear_env", "clear_working_dir", "clear_run_as", "stdin", "env"} {
		if _, exists := body[key]; exists {
			t.Errorf("omitted %s changed: %v", key, body)
		}
	}
	if body["args"] != nil {
		t.Fatalf("omitted args must retain values: %v", body)
	}
}

func TestTaskAuthoringRejectsConflictingIntentBeforeRequest(t *testing.T) {
	for _, flags := range [][]string{
		{"--clear-args", "--arg", ""},
		{"--clear-env", "--env", "EXAMPLE=value"},
		{"--clear-cwd", "--cwd", ""},
		{"--clear-run-as", "--run-as", "user"},
		{"--clear-stdin", "--stdin", ""},
		{"--stdin", "text", "--stdin-file", "-"},
		{"--syntax", "cron"},
		{"--schedule", "every day", "--syntax", "invalid"},
	} {
		t.Run(strings.Join(flags, " "), func(t *testing.T) {
			_, _, calls, err := executeAuthoring(t, append([]string{"task", "edit", "fixture"}, flags...), "")
			if err == nil || calls != 0 {
				t.Fatalf("err=%v calls=%d", err, calls)
			}
		})
	}
}

func TestGroupAuthoringTransportsDisabledRenameAndRootIntent(t *testing.T) {
	body, path, _, err := executeAuthoring(t, []string{"group", "add", "fixture", "--disabled", "--parent", "parent", "--json"}, "")
	if err != nil || path != "/v1/groups" || body["enabled"] != false || body["parent_id"] != "parent" {
		t.Fatalf("path=%s body=%v err=%v", path, body, err)
	}
	body, path, _, err = executeAuthoring(t, []string{"group", "edit", "fixture", "--name", "renamed", "--parent", "", "--json"}, "")
	if err != nil || path != "/v1/groups/fixture" || body["name"] != "renamed" || body["parent_id"] != "" {
		t.Fatalf("path=%s body=%v err=%v", path, body, err)
	}
}

func TestSchedulePreviewTransportsCompleteNonMutatingPolicy(t *testing.T) {
	body, path, calls, err := executeAuthoring(t, []string{"schedule", "preview", "every month on the 31st at 09:00", "--syntax", "human", "--tz", "America/New_York", "--missing-date", "last_valid", "--time-basis", "wall_clock", "--dst-gap", "skip", "--dst-overlap", "both", "--json"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 || path != "/v1/schedules/preview" || body["schedule_syntax"] != "human" || body["timezone"] != "America/New_York" || body["missing_date_policy"] != "last_valid" || body["time_basis"] != "wall_clock" || body["dst_gap_policy"] != "skip" || body["dst_overlap_policy"] != "both" {
		t.Fatalf("calls=%d path=%s body=%v", calls, path, body)
	}
}
