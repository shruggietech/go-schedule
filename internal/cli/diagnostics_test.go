package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/clientsecret"
	"github.com/shruggietech/go-schedule/internal/clifixture"
)

func TestContextualLookupDiagnosticsIdentifyRequestedObject(t *testing.T) {
	for _, tc := range []struct {
		args []string
		kind string
	}{
		{[]string{"task", "show", "missing-object"}, "task"},
		{[]string{"group", "edit", "missing-object", "--name", "replacement"}, "group"},
		{[]string{"chain", "show", "missing-object"}, "chain"},
		{[]string{"watcher", "show", "missing-object"}, "watcher"},
		{[]string{"trigger", "show", "missing-object"}, "trigger"},
		{[]string{"runs", "get", "missing-object"}, "run"},
	} {
		t.Run(strings.Join(tc.args[:2], " "), func(t *testing.T) {
			var out, diagnostic bytes.Buffer
			api := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = io.WriteString(w, `{"error":{"code":"not_found","message":"not found"}}`)
			}))
			code := ExecuteWithDependencies(context.Background(), tc.args, Dependencies{Client: api, Out: &out, Err: &diagnostic})
			if code != 1 || out.Len() != 0 || !strings.Contains(diagnostic.String(), tc.kind) || !strings.Contains(diagnostic.String(), "missing-object") || !strings.Contains(diagnostic.String(), "--help") {
				t.Fatalf("exit=%d out=%q diagnostic=%q", code, out.String(), diagnostic.String())
			}
		})
	}
}

func TestContextualScheduleValidationExplainsFieldAndCorrection(t *testing.T) {
	for _, tc := range []struct{ field, flag, choice string }{
		{"missing_date_policy", "--missing-date", "last_valid"},
		{"time_basis", "--time-basis", "wall_clock"},
		{"dst_gap_policy", "--dst-gap", "next_valid"},
		{"dst_overlap_policy", "--dst-overlap", "both"},
		{"timezone", "--tz", "America/New_York"},
		{"schedule", "schedule expression", "cron"},
	} {
		t.Run(tc.field, func(t *testing.T) {
			var out, diagnostic bytes.Buffer
			api := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"code": server.CodeValidation, "field": tc.field, "message": "invalid policy"}})
			}))
			code := ExecuteWithDependencies(context.Background(), []string{"schedule", "preview", "every day at 09:00"}, Dependencies{Client: api, Out: &out, Err: &diagnostic})
			if code != 2 || out.Len() != 0 || !strings.Contains(diagnostic.String(), tc.flag) || !strings.Contains(diagnostic.String(), tc.choice) {
				t.Fatalf("exit=%d out=%q diagnostic=%q", code, out.String(), diagnostic.String())
			}
		})
	}
}

func TestCertificateTargetDiagnosticsNameInputAndPreserveExitCategory(t *testing.T) {
	file := filepath.Join(t.TempDir(), "invalid-trust.pem")
	const secret = "owned-certificate-or-token-do-not-print"
	if err := os.WriteFile(file, []byte(secret), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		file       string
		code       int
		correction string
	}{
		{file, 2, "PEM"},
		{file + ".missing", 1, "read"},
	} {
		t.Run(tc.correction, func(t *testing.T) {
			recorder := &clifixture.Recorder{}
			secrets := clientsecret.NewWithBackend(recorder)
			if err := secrets.Save("daemon", "credential", secret); err != nil {
				t.Fatal(err)
			}
			var out, diagnostic bytes.Buffer
			code := ExecuteWithDependencies(context.Background(), []string{"--daemon-endpoint", "https://example.test", "--daemon-id", "daemon", "--credential-id", "credential", "--certificate-file", tc.file, "health"}, Dependencies{Secrets: secrets, Out: &out, Err: &diagnostic})
			if code != tc.code || out.Len() != 0 || !strings.Contains(diagnostic.String(), "--certificate-file") || !strings.Contains(diagnostic.String(), filepath.Base(tc.file)) || !strings.Contains(diagnostic.String(), tc.correction) || strings.Contains(diagnostic.String(), secret) {
				t.Fatalf("exit=%d out=%q diagnostic=%q", code, out.String(), diagnostic.String())
			}
		})
	}
}

func TestContextualDiagnosticsPreserveErrorIdentityAndClassification(t *testing.T) {
	root := newRoot()
	cmd, _, err := root.Find([]string{"task", "show"})
	if err != nil {
		t.Fatal(err)
	}
	status := &client.StatusError{Code: server.CodeValidation, Field: "timezone", Message: "invalid"}
	sentinel := errors.New("original cause")
	original := fmt.Errorf("%w: %w", status, sentinel)
	wrapped := contextualCommandError(cmd, []string{"owned-id"}, original)
	var observed *client.StatusError
	if !errors.Is(wrapped, original) || !errors.Is(wrapped, sentinel) || !errors.As(wrapped, &observed) || observed != status {
		t.Fatal("context lost the original structured error chain")
	}
	if code := handleExecuteError(io.Discard, wrapped); code != 2 {
		t.Fatalf("validation exit=%d", code)
	}
	if contextualCommandError(cmd, nil, wrapped) != wrapped {
		t.Fatal("error context was applied twice")
	}
	uncertain := &client.MutationUncertainError{Operation: "POST /owned", Cause: status}
	if contextualCommandError(cmd, nil, uncertain) != uncertain {
		t.Fatal("uncertain mutation guidance was replaced")
	}
	if contextualCommandError(cmd, nil, sentinel) != sentinel || contextualCommandError(cmd, nil, nil) != nil {
		t.Fatal("unrelated errors were changed")
	}
}

func TestContextualDiagnosticsNeverReplayTriggerKeyOrServerPayload(t *testing.T) {
	const secret = "owned-trigger-secret-do-not-print"
	for _, code := range []string{server.CodeNotFound, server.CodeValidation} {
		var out, diagnostic bytes.Buffer
		api := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"code": code, "field": secret + "\n" + secret, "message": secret}})
		}))
		exit := ExecuteWithDependencies(context.Background(), []string{"trigger", "fire", secret}, Dependencies{Client: api, Out: &out, Err: &diagnostic})
		if exit == 0 || out.Len() != 0 || diagnostic.Len() == 0 || strings.Contains(diagnostic.String(), secret) {
			t.Fatalf("diagnostic leaked input/server payload or lost error: exit=%d", exit)
		}
	}
}

func TestContextualLookupEscapesIdentifierControlCharacters(t *testing.T) {
	root := newRoot()
	cmd, _, err := root.Find([]string{"task", "show"})
	if err != nil {
		t.Fatal(err)
	}
	err = contextualCommandError(cmd, []string{"owned\n\x1b[31m"}, &client.StatusError{Code: server.CodeNotFound, Message: "not found"})
	if strings.ContainsAny(err.Error(), "\n\x1b") || !strings.Contains(err.Error(), `owned\n\x1b[31m`) {
		t.Fatal("identifier was not safely quoted")
	}
}

func TestContextualBundleDiagnosticsPreserveBindingIdentityWithoutPath(t *testing.T) {
	file := filepath.Join(t.TempDir(), "owned-bundle.json")
	if err := os.WriteFile(file, []byte(`{"schema":"go-schedule.bundle/v2"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	const privatePath = "private-bound-path-do-not-print"
	for _, reason := range []string{"unknown", "duplicate"} {
		var out, diagnostic bytes.Buffer
		api := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			message := `watcher path binding "unknown-watcher" is not a watcher in this bundle; use a watcher portable_id from the bundle or remove the binding`
			if reason == "duplicate" {
				message = `duplicate watcher path binding "unknown-watcher"; supply each watcher portable_id exactly once`
			}
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"code": server.CodeValidation, "field": "watcher_paths", "message": message}})
		}))
		code := ExecuteWithDependencies(context.Background(), []string{"bundle", "preview", file, "--watcher-path=unknown-watcher=" + privatePath}, Dependencies{Client: api, Out: &out, Err: &diagnostic})
		if code != 2 || out.Len() != 0 || !strings.Contains(diagnostic.String(), "--watcher-path") || !strings.Contains(diagnostic.String(), "unknown-watcher") || !strings.Contains(diagnostic.String(), "portable_id") || strings.Contains(diagnostic.String(), privatePath) {
			t.Fatalf("exit=%d diagnostic=%q", code, diagnostic.String())
		}
	}
}

func TestContextualBundleValidationIdentifiesDocumentAndNextAction(t *testing.T) {
	file := filepath.Join(t.TempDir(), "owned-bundle.json")
	if err := os.WriteFile(file, []byte(`{"schema":"go-schedule.bundle/v2"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, diagnostic bytes.Buffer
	api := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"code": server.CodeValidation, "field": "bundle", "message": "bundle is invalid"}})
	}))
	code := ExecuteWithDependencies(context.Background(), []string{"bundle", "preview", file}, Dependencies{Client: api, Out: &out, Err: &diagnostic})
	if code != 2 || out.Len() != 0 || !strings.Contains(diagnostic.String(), filepath.Base(file)) || !strings.Contains(diagnostic.String(), "bundle validate") {
		t.Fatalf("exit=%d diagnostic=%q", code, diagnostic.String())
	}
}

func TestMissingProfileDiagnosticPreservesIdentityAndCause(t *testing.T) {
	profiles := &invocationProfileMemory{}
	_, _, err := resolveTargetWith(context.Background(), remoteTargetFlags{Profile: "missing-profile"}, makeInvocation(Dependencies{Profiles: profiles}).deps)
	if !errors.Is(err, errUsage) || !errors.Is(err, clientprofile.ErrNotFound) || !strings.Contains(err.Error(), "missing-profile") || !strings.Contains(err.Error(), "profile list") {
		t.Fatalf("missing profile context/cause: %v", err)
	}
	var out, diagnostic bytes.Buffer
	code := ExecuteWithDependencies(context.Background(), []string{"profile", "show", "missing-profile"}, Dependencies{Profiles: profiles, Out: &out, Err: &diagnostic})
	if code != 2 || out.Len() != 0 || !strings.Contains(diagnostic.String(), "missing-profile") || !strings.Contains(diagnostic.String(), "profile list") {
		t.Fatalf("profile lookup exit=%d diagnostic=%q", code, diagnostic.String())
	}
}
