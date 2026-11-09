package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
)

func TestPairingHelpDescribesCompleteInputFlow(t *testing.T) {
	var out, diagnostic bytes.Buffer
	code := ExecuteWithDependencies(context.Background(), []string{"profile", "pair", "--help"}, Dependencies{Client: client.New("unused"), Out: &out, Err: &diagnostic})
	for _, want := range []string{"standard input", "end-of-input", "command history"} {
		if code != 0 || diagnostic.Len() != 0 || !strings.Contains(out.String(), want) {
			t.Fatalf("exit=%d, help lacks %q: %s %s", code, want, out.String(), diagnostic.String())
		}
	}
}

func TestOneOffTimestampCorrectionsAgree(t *testing.T) {
	for _, args := range [][]string{{"task", "add", "owned", "--command", "unused", "--at", "invalid"}, {"task", "edit", "owned-id", "--at", "invalid"}} {
		var out, diagnostic bytes.Buffer
		code := ExecuteWithDependencies(context.Background(), args, Dependencies{Client: client.New("unused"), Out: &out, Err: &diagnostic})
		if code != 2 || out.Len() != 0 || !strings.Contains(diagnostic.String(), "2026-08-04T09:00:00Z") {
			t.Fatalf("args=%q exit=%d diagnostic=%q", args, code, diagnostic.String())
		}
	}
}

func TestGraphRecoveryUsesKnownContractsWithoutServerEcho(t *testing.T) {
	for _, tc := range []struct {
		path, field, message string
		want                 []string
	}{
		{"group edit", "parent_id", "would create a group cycle", []string{"--parent", "descendant", "cycle"}},
		{"chain create", "target_task_id", "chain would create a cycle", []string{"--source", "--target", "cycle"}},
	} {
		cmd, remaining, err := newRoot().Find(strings.Fields(tc.path))
		if err != nil || len(remaining) != 0 {
			t.Fatalf("command lookup: %v, remaining=%q", err, remaining)
		}
		status := &client.StatusError{Code: server.CodeValidation, Field: tc.field, Message: tc.message}
		observed := contextualCommandError(cmd, []string{"owned"}, status).Error()
		for _, want := range tc.want {
			if !strings.Contains(observed, want) {
				t.Errorf("%s lacks %q: %s", tc.path, want, observed)
			}
		}
		status.Message = "private-server-value-do-not-echo"
		if observed := contextualCommandError(cmd, nil, status).Error(); strings.Contains(observed, status.Message) || strings.Contains(observed, "cycle") {
			t.Fatalf("unrecognized server content was echoed or mislabeled: %s", observed)
		}
	}
}
