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

	"github.com/shruggietech/go-schedule/internal/api/client"
)

func TestAlertsJSONHasOnlyOneDocumentAndDiagnosticNotice(t *testing.T) {
	var out, diagnostic bytes.Buffer
	api := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, `{"alerts":[]}`) }))
	code := ExecuteWithDependencies(context.Background(), []string{"--json", "alerts"}, Dependencies{Client: api, Out: &out, Err: &diagnostic})
	if code != 0 || !json.Valid(out.Bytes()) || !strings.Contains(diagnostic.String(), "deprecat") {
		t.Fatalf("exit=%d out=%q stderr=%q", code, out.String(), diagnostic.String())
	}
}

func TestCredentialPairingAndServiceHelpExplainsOperations(t *testing.T) {
	for _, path := range [][]string{{"credential", "list"}, {"credential", "rotate"}, {"credential", "revoke"}, {"pairing", "list"}, {"pairing", "cancel"}, {"service", "install"}, {"service", "start"}, {"service", "stop"}, {"service", "restart"}, {"service", "status"}, {"service", "uninstall"}} {
		root := newRoot()
		cmd, _, err := root.Find(path)
		if err != nil {
			t.Fatal(err)
		}
		if cmd.Short == "" || strings.HasPrefix(cmd.Short, "Service:") || cmd.Long == "" {
			t.Errorf("%v has placeholder/missing operation help: short=%q long=%q", path, cmd.Short, cmd.Long)
		}
	}
}

func TestHumanOperationalListsLabelColumnsAndKeepJSON(t *testing.T) {
	for _, tc := range []struct {
		args     []string
		headings []string
		body     string
	}{
		{[]string{"actor", "list"}, []string{"ACTOR ID", "CAPABILITY", "NAME"}, `{"actors":[]}`},
		{[]string{"audit", "list"}, []string{"TIME", "ACTOR ID", "TARGET ID"}, `{"events":[]}`},
		{[]string{"pairing", "list"}, []string{"PAIRING ID", "ATTEMPTS"}, `{"pairings":[]}`},
		{[]string{"notification", "channel", "list"}, []string{"CHANNEL ID", "ENDPOINT"}, `{"notification_channels":[]}`},
		{[]string{"watcher", "list"}, []string{"WATCHER ID", "READINESS", "PATH"}, `{"filesystem_watchers":[]}`},
		{[]string{"profile", "list"}, []string{"PROFILE ID", "DAEMON ID"}, `{}`},
	} {
		for _, asJSON := range []bool{false, true} {
			var out, diagnostic bytes.Buffer
			args := append([]string{}, tc.args...)
			if asJSON {
				args = append([]string{"--json"}, args...)
			}
			api := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, tc.body) }))
			profiles := &invocationProfileMemory{}
			code := ExecuteWithDependencies(context.Background(), args, Dependencies{Client: api, Profiles: profiles, Out: &out, Err: &diagnostic})
			if code != 0 {
				t.Fatalf("args=%v exit=%d stderr=%q", args, code, diagnostic.String())
			}
			if asJSON {
				if !json.Valid(out.Bytes()) {
					t.Errorf("invalid JSON for %v: %q", args, out.String())
				}
				continue
			}
			for _, heading := range tc.headings {
				if !strings.Contains(out.String(), heading) {
					t.Errorf("%v missing %q: %q", args, heading, out.String())
				}
			}
		}
	}
}

func TestListOutputWriteFailureIsReported(t *testing.T) {
	api := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `{"actors":[{"id":"owned","display_name":"Owned"}]}`)
	}))
	if code := ExecuteWithDependencies(context.Background(), []string{"actor", "list"}, Dependencies{Client: api, Out: failingOutput{}, Err: io.Discard}); code != 1 {
		t.Fatalf("lost list writer failure, exit=%d", code)
	}
}

type failingOutput struct{}

func (failingOutput) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }
