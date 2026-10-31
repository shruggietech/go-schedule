package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/buildinfo"
	"github.com/shruggietech/go-schedule/v2/internal/clientprofile"
	"github.com/shruggietech/go-schedule/v2/internal/clientsecret"
	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
	"github.com/shruggietech/go-schedule/v2/internal/config"
	"github.com/shruggietech/go-schedule/v2/internal/ipc"
)

func TestInvocationRefusalsKeepUsageCategoryBeforeSideEffects(t *testing.T) {
	for _, args := range [][]string{
		{"health", "--unknown-contract-flag"},
		{"audit", "list", "--limit=not-an-integer"},
		{"unknown-contract-command"},
		{"task", "unknown-contract-command"},
		{"task", "show"},
		{"health", "unwanted", "operands"},
		{"service", "start", "unwanted"},
		{"gui", "unwanted"},
		{"cron", "import", "unwanted", "--dry-run"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			calls := 0
			var diagnostics bytes.Buffer
			deps := Dependencies{Out: io.Discard, Err: &diagnostics,
				Client:         client.NewInProcess(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls++ })),
				DaemonPath:     func() (string, error) { calls++; return "owned", nil },
				GUIPath:        func() (string, error) { calls++; return "owned", nil },
				LaunchGUI:      func(string) error { calls++; return nil },
				ServiceControl: func(string, string, []string) (string, error) { calls++; return "owned", nil },
			}
			// A missing remote profile must not be resolved for an invalid invocation.
			argv := append([]string{"--profile", "must-not-resolve"}, args...)
			if code := ExecuteWithDependencies(context.Background(), argv, deps); code != 2 {
				t.Fatalf("exit=%d, want usage 2: %s", code, diagnostics.String())
			}
			if calls != 0 {
				t.Fatalf("invalid invocation made %d side-effect/connection calls", calls)
			}
			if diagnostics.Len() == 0 {
				t.Fatal("missing refusal diagnostic")
			}
		})
	}
}

func TestWebhookDestinationAndDaemonTargetAreIndependent(t *testing.T) {
	const destination = "https://webhook.example.test/events"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/manifest") {
			fmt.Fprint(w, `{"installation_id":"owned-daemon"}`)
			return
		}
		var body struct {
			Endpoint string `json:"endpoint"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Endpoint != destination {
			t.Errorf("destination=%q err=%v", body.Endpoint, err)
		}
		fmt.Fprintf(w, `{"id":"channel","name":"Owned","endpoint_summary":%q}`, body.Endpoint)
	})
	remote := httptest.NewTLSServer(handler)
	defer remote.Close()
	certificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: remote.Certificate().Raw})
	certFile := filepath.Join(t.TempDir(), "trusted.pem")
	if err := os.WriteFile(certFile, certificate, 0o600); err != nil {
		t.Fatal(err)
	}
	socketDir := t.TempDir()
	if runtime.GOOS == "darwin" {
		shortDir, err := os.MkdirTemp("/tmp", "gs-cli-")
		if err != nil {
			t.Fatal(err)
		}
		socketDir = shortDir
		t.Cleanup(func() {
			if err := os.RemoveAll(shortDir); err != nil {
				t.Errorf("remove owned IPC directory: %v", err)
			}
		})
	}
	endpoint := filepath.Join(socketDir, "owned.sock")
	if runtime.GOOS == "windows" {
		endpoint = `\\.\pipe\s109-` + uuid.NewString()
	}
	listener, _, err := ipc.Listen(config.Config{IPCPath: endpoint, AdminGroup: ""})
	if err != nil {
		t.Fatal(err)
	}
	local := &http.Server{Handler: handler}
	finished := make(chan error, 1)
	go func() { finished <- local.Serve(listener) }()
	t.Cleanup(func() {
		if err := local.Close(); err != nil {
			t.Error(err)
		}
		if err := <-finished; err != nil && err != http.ErrServerClosed {
			t.Error(err)
		}
	})
	profiles := &invocationProfileMemory{collection: clientprofile.Collection{Version: clientprofile.CurrentVersion, Profiles: []clientprofile.Profile{{ID: "owned", Label: "Owned", Endpoint: remote.URL, DaemonID: "owned-daemon", CredentialID: "owned-credential", CertificatePEM: string(certificate)}}}}
	recorder := &clifixture.Recorder{}
	secrets := clientsecret.NewWithBackend(recorder)
	if err := secrets.Save("owned-daemon", "owned-credential", "owned-token-do-not-print"); err != nil {
		t.Fatal(err)
	}
	for _, selector := range [][]string{
		{"--ipc-path", endpoint},
		{"--profile", "owned"},
		{"--daemon-endpoint", remote.URL, "--daemon-id", "owned-daemon", "--credential-id", "owned-credential", "--certificate-file", certFile},
	} {
		for _, operation := range [][]string{{"add", "Owned"}, {"update", "channel"}} {
			var out, diagnostic bytes.Buffer
			args := append(append([]string{}, selector...), "--json", "notification", "channel")
			args = append(args, operation...)
			args = append(args, "--endpoint", destination)
			code := ExecuteWithDependencies(context.Background(), args, Dependencies{Profiles: profiles, Secrets: secrets, Out: &out, Err: &diagnostic})
			if code != 0 || !json.Valid(out.Bytes()) || !strings.Contains(out.String(), destination) || strings.Contains(out.String()+diagnostic.String(), "owned-token-do-not-print") {
				t.Fatalf("args=%v exit=%d out=%s stderr=%s", args, code, out.String(), diagnostic.String())
			}
		}
	}
}

func TestRemoteEndpointAliasesAreMutuallyExclusive(t *testing.T) {
	var diagnostic bytes.Buffer
	code := ExecuteWithDependencies(context.Background(), []string{"--endpoint=", "--daemon-endpoint=", "health"}, Dependencies{Out: io.Discard, Err: &diagnostic})
	if code != 2 || !strings.Contains(diagnostic.String(), "choose") {
		t.Fatalf("exit=%d stderr=%q", code, diagnostic.String())
	}
}

func TestProfileRefusesExplicitEmptyRemoteSelectorIntent(t *testing.T) {
	for _, flag := range []string{"endpoint", "daemon-endpoint", "daemon-id", "credential-id", "certificate-file"} {
		var diagnostics bytes.Buffer
		code := ExecuteWithDependencies(context.Background(), []string{"--profile=owned", "--" + flag + "=", "health"}, Dependencies{Profiles: &invocationProfileMemory{}, Out: io.Discard, Err: &diagnostics})
		if code != 2 || !strings.Contains(diagnostics.String(), "cannot be combined") {
			t.Fatalf("flag=%s exit=%d stderr=%q", flag, code, diagnostics.String())
		}
	}
}

func TestGeneratedCommandsPreserveUsageAndLocalHelp(t *testing.T) {
	for _, tc := range []struct {
		args []string
		want int
	}{
		{[]string{"completion", "bash", "extra"}, 2},
		{[]string{"help", "unknown-contract-topic"}, 2},
		{[]string{"help", "task", "unknown-contract-topic"}, 2},
		{[]string{"--profile=missing", "help", "task", "show"}, 0},
		{[]string{"--daemon-endpoint=", "completion", "bash"}, 0},
		{[]string{"__complete", "--help"}, 0},
		{[]string{"__completeNoDesc", ""}, 0},
		{[]string{"__complete"}, 2},
	} {
		var out, diagnostic bytes.Buffer
		code := ExecuteWithDependencies(context.Background(), tc.args, Dependencies{Profiles: &invocationProfileMemory{}, Out: &out, Err: &diagnostic})
		if code != tc.want {
			t.Errorf("%v exit=%d want=%d stderr=%q", tc.args, code, tc.want, diagnostic.String())
		}
	}
	for _, args := range [][]string{{"--help"}, {"task"}, {"help", "task"}, {"completion", "bash"}} {
		if code := ExecuteWithDependencies(context.Background(), args, Dependencies{Out: failingOutput{}, Err: io.Discard}); code != 1 {
			t.Errorf("help/output write failure for %v returned %d", args, code)
		}
	}
}

func TestEmptyRemoteSelectorsNeverFallBackToLocal(t *testing.T) {
	for _, flag := range []string{"profile", "endpoint", "daemon-endpoint", "daemon-id", "credential-id", "certificate-file"} {
		for _, value := range []string{"", "   "} {
			calls := 0
			api := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { calls++; fmt.Fprint(w, `{}`) }))
			var diagnostic bytes.Buffer
			code := ExecuteWithDependencies(context.Background(), []string{"--" + flag + "=" + value, "health"}, Dependencies{Client: api, Out: io.Discard, Err: &diagnostic})
			if code != 2 || calls != 0 || !strings.Contains(diagnostic.String(), flag) {
				t.Errorf("flag=%s exit=%d calls=%d stderr=%q", flag, code, calls, diagnostic.String())
			}
		}
	}
}

func TestMalformedBundleIsUsageButFileIOIsRuntime(t *testing.T) {
	bad := filepath.Join(t.TempDir(), "broken-bundle.json")
	if err := os.WriteFile(bad, []byte(`{"version":`), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		file string
		want int
	}{{bad, 2}, {bad + ".missing", 1}} {
		var diagnostics bytes.Buffer
		code := ExecuteWithDependencies(context.Background(), []string{"bundle", "validate", tc.file}, Dependencies{Client: healthFixture("unused"), Out: io.Discard, Err: &diagnostics})
		if code != tc.want || (!strings.Contains(diagnostics.String(), tc.file) && !strings.Contains(diagnostics.String(), strconv.Quote(tc.file))) {
			t.Fatalf("file=%s exit=%d want=%d stderr=%q", tc.file, code, tc.want, diagnostics.String())
		}
	}
}

func TestOperationalFailureDoesNotBecomeUsage(t *testing.T) {
	var diagnostics bytes.Buffer
	code := ExecuteWithDependencies(context.Background(), []string{"health"}, Dependencies{
		Client: client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"error":{"code":"internal","message":"owned runtime failure"}}`)
		})), Out: io.Discard, Err: &diagnostics,
	})
	if code != 1 {
		t.Fatalf("runtime exit=%d stderr=%q", code, diagnostics.String())
	}
}

func TestDuplicateWatcherBindingsRefuseBeforePreviewRequest(t *testing.T) {
	file := filepath.Join(t.TempDir(), "owned.json")
	if err := os.WriteFile(file, []byte(`{"schema":"go-schedule.bundle/v2"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, second := range []string{"owned=/first", "owned=/second"} {
		requests := 0
		var diagnostics bytes.Buffer
		api := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { requests++; fmt.Fprint(w, `{}`) }))
		code := ExecuteWithDependencies(context.Background(), []string{"bundle", "preview", file, "--watcher-path=owned=/first", "--watcher-path=" + second}, Dependencies{Client: api, Out: io.Discard, Err: &diagnostics})
		if code != 2 || requests != 0 || !strings.Contains(diagnostics.String(), "owned") || !strings.Contains(diagnostics.String(), "duplicate") {
			t.Fatalf("exit=%d requests=%d stderr=%q", code, requests, diagnostics.String())
		}
	}
}

func TestRootVersionRetainsParsedFlagWithoutTargetResolution(t *testing.T) {
	for _, args := range [][]string{{"--version"}, {"-v"}, {"--version=true"}, {"--profile=missing", "--version"}} {
		var out, diagnostic bytes.Buffer
		code := ExecuteWithDependencies(context.Background(), args, Dependencies{Profiles: &invocationProfileMemory{}, Out: &out, Err: &diagnostic})
		want := "gosched version " + buildinfo.Version + "\n"
		if code != 0 || out.String() != want || diagnostic.Len() != 0 {
			t.Errorf("%v exit=%d stdout=%q want=%q stderr=%q", args, code, out.String(), want, diagnostic.String())
		}
	}
}
