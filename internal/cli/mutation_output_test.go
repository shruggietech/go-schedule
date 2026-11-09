package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
)

type acknowledgementCase struct {
	args   []string
	human  string
	fields map[string]any
}

func acknowledgementCases() []acknowledgementCase {
	cases := []acknowledgementCase{
		{[]string{"alerts", "ack", "owned-id"}, "acknowledged alert owned-id\n", map[string]any{"id": "owned-id", "status": "acknowledged"}},
		{[]string{"notification", "channel", "rm", "owned-id"}, "deleted notification channel owned-id\n", map[string]any{"id": "owned-id", "status": "deleted"}},
		{[]string{"task", "enable", "owned-id"}, "task owned-id enabled\n", map[string]any{"id": "owned-id", "enabled": true}},
		{[]string{"task", "disable", "owned-id"}, "task owned-id disabled\n", map[string]any{"id": "owned-id", "enabled": false}},
		{[]string{"task", "rm", "owned-id"}, "deleted task owned-id\n", map[string]any{"id": "owned-id", "status": "deleted"}},
		{[]string{"task", "run-now", "owned-id"}, "triggered run for task owned-id\n", map[string]any{"id": "owned-id", "status": "accepted"}},
		{[]string{"group", "enable", "owned-id"}, "group owned-id enabled\n", map[string]any{"id": "owned-id", "enabled": true}},
		{[]string{"group", "disable", "owned-id"}, "group owned-id disabled\n", map[string]any{"id": "owned-id", "enabled": false}},
		{[]string{"group", "rm", "owned-id"}, "deleted group owned-id\n", map[string]any{"id": "owned-id", "status": "deleted"}},
		{[]string{"chain", "rm", "owned-id"}, "deleted chain owned-id\n", map[string]any{"id": "owned-id", "status": "deleted"}},
		{[]string{"gui"}, "launched gosched-gui\n", map[string]any{"status": "launched"}},
	}
	for _, action := range []string{"install", "uninstall", "start", "stop", "restart", "status"} {
		cases = append(cases, acknowledgementCase{[]string{"service", action}, action + ": owned service result\n", map[string]any{"action": action, "status": "owned service result"}})
	}
	return cases
}

func acknowledgementDependencies(out, diagnostic io.Writer, calls *int) Dependencies {
	api := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { *calls++; w.WriteHeader(http.StatusNoContent) }))
	return Dependencies{Client: api, Out: out, Err: diagnostic, DaemonPath: func() (string, error) { return "owned-daemon", nil }, ServiceControl: func(string, string, []string) (string, error) { *calls++; return "owned service result", nil }, GUIPath: func() (string, error) { return "owned-gui", nil }, LaunchGUI: func(string) error { *calls++; return nil }}
}

func TestMutationAcknowledgementsRespectOutputMode(t *testing.T) {
	for _, tc := range acknowledgementCases() {
		for _, asJSON := range []bool{false, true} {
			t.Run(strings.Join(tc.args, " ")+fmt.Sprint("/json=", asJSON), func(t *testing.T) {
				var out, diagnostic bytes.Buffer
				calls := 0
				args := append([]string{}, tc.args...)
				if asJSON {
					args = append([]string{"--json"}, args...)
				}
				code := ExecuteWithDependencies(context.Background(), args, acknowledgementDependencies(&out, &diagnostic, &calls))
				if code != 0 || calls != 1 {
					t.Fatalf("exit=%d calls=%d stderr=%q", code, calls, diagnostic.String())
				}
				if !asJSON {
					if out.String() != tc.human {
						t.Fatalf("human output=%q want %q", out.String(), tc.human)
					}
					return
				}
				var got map[string]any
				if err := json.Unmarshal(out.Bytes(), &got); err != nil {
					t.Fatalf("not one JSON acknowledgement: %q (%v)", out.String(), err)
				}
				for key, want := range tc.fields {
					if got[key] != want {
						t.Errorf("%s=%v want %v", key, got[key], want)
					}
				}
			})
		}
	}
}

func TestMutationAcknowledgementsReportOutputFailure(t *testing.T) {
	for _, tc := range acknowledgementCases() {
		for _, asJSON := range []bool{false, true} {
			t.Run(strings.Join(tc.args, " ")+fmt.Sprint("/json=", asJSON), func(t *testing.T) {
				calls := 0
				args := append([]string{}, tc.args...)
				if asJSON {
					args = append([]string{"--json"}, args...)
				}
				code := ExecuteWithDependencies(context.Background(), args, acknowledgementDependencies(failingOutput{}, io.Discard, &calls))
				if code != 1 || calls != 1 {
					t.Fatalf("lost output failure or replayed operation: exit=%d calls=%d", code, calls)
				}
			})
		}
	}
}
