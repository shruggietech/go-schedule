package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/clifixture"
)

func healthFixture(version string) *client.Client {
	return client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, `{"version":%q,"ok":true}`, version) }))
}

type invocationProfileMemory struct{ collection clientprofile.Collection }

func (p *invocationProfileMemory) Load() (clientprofile.Collection, error) { return p.collection, nil }
func (p *invocationProfileMemory) Add(value clientprofile.Profile) (clientprofile.Profile, error) {
	p.collection.Profiles = append(p.collection.Profiles, value)
	return value, nil
}
func (p *invocationProfileMemory) Rename(id, label string) (clientprofile.Profile, error) {
	return clientprofile.Profile{}, errors.New("unused rename adapter")
}
func (p *invocationProfileMemory) SetActive(id string) error {
	p.collection.ActiveDesktopProfileID = id
	return nil
}
func (p *invocationProfileMemory) RemoveWith(id string, before func(clientprofile.Profile) error) (clientprofile.Profile, error) {
	for i, value := range p.collection.Profiles {
		if value.ID == id {
			if err := before(value); err != nil {
				return clientprofile.Profile{}, err
			}
			p.collection.Profiles = append(p.collection.Profiles[:i], p.collection.Profiles[i+1:]...)
			return value, nil
		}
	}
	return clientprofile.Profile{}, errors.New("fixture profile does not exist")
}

func TestActualProfileHandlerUsesRecordingSecretsAndScopedMetadata(t *testing.T) {
	profile := clientprofile.Profile{ID: "owned-profile", Label: "Owned", DaemonID: "owned-daemon", CredentialID: "owned-credential"}
	profiles := &invocationProfileMemory{collection: clientprofile.Collection{Version: clientprofile.CurrentVersion, Profiles: []clientprofile.Profile{profile}}}
	recorder := &clifixture.Recorder{}
	secrets := recorder.SecretStore()
	if err := secrets.Save(profile.DaemonID, profile.CredentialID, "fixture-secret-do-not-print"); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if code := ExecuteWithDependencies(context.Background(), []string{"profile", "remove", profile.ID, "--confirm", profile.ID}, Dependencies{Profiles: profiles, Secrets: secrets, Out: &output, Err: io.Discard}); code != 0 {
		t.Fatalf("exit=%d output=%q", code, output.String())
	}
	if len(profiles.collection.Profiles) != 0 {
		t.Fatal("owned profile metadata remains")
	}
	if _, err := secrets.Load(profile.DaemonID, profile.CredentialID); err == nil {
		t.Fatal("owned secret remains")
	}
	encoded, err := json.Marshal(recorder.Intents())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "fixture-secret") || bytes.Contains(encoded, []byte("fixture-secret")) {
		t.Fatal("recording output disclosed fixture secret")
	}
	found := false
	for _, intent := range recorder.Intents() {
		if intent.Operation == "secret.delete" {
			found = true
		}
		if intent.Evidence != clifixture.CommandHarness {
			t.Fatal(intent)
		}
	}
	if !found {
		t.Fatal("real handler did not request credential deletion")
	}
}

func TestInvocationRootsDoNotShareTargetJSONOrIO(t *testing.T) {
	var first, second bytes.Buffer
	a := newRootWithDependencies(Dependencies{Client: healthFixture("first"), Out: &first, Err: io.Discard})
	b := newRootWithDependencies(Dependencies{Client: healthFixture("second"), Out: &second, Err: io.Discard})
	a.SetArgs([]string{"--json", "health"})
	b.SetArgs([]string{"health"})
	var wg sync.WaitGroup
	for _, root := range []*cobra.Command{a, b} {
		wg.Add(1)
		go func(root *cobra.Command) {
			defer wg.Done()
			if err := root.Execute(); err != nil {
				t.Error(err)
			}
		}(root)
	}
	wg.Wait()
	if !strings.Contains(first.String(), `"version": "first"`) || !strings.Contains(second.String(), "daemon ok (version second)") {
		t.Fatalf("first=%q second=%q", first.String(), second.String())
	}
	if stateFor(a).client == stateFor(b).client || stateFor(a).deps.Secrets == stateFor(b).deps.Secrets {
		t.Fatal("roots share mutable invocation state")
	}
}

func TestExplicitIPCNeverFallsBackToInjectedDefault(t *testing.T) {
	endpoint := filepath.Join(t.TempDir(), "missing.sock")
	if runtime.GOOS == "windows" {
		endpoint = `\\.\pipe\s107-missing-invocation`
	}
	for _, args := range [][]string{{"--ipc-path", "relative", "health"}, {"--ipc-path", "", "health"}, {"--ipc-path", endpoint, "--profile", "production", "health"}, {"--ipc-path", endpoint, "health"}} {
		var output bytes.Buffer
		root := newRootWithDependencies(Dependencies{Client: healthFixture("default-must-not-be-contacted"), Out: &output, Err: io.Discard})
		root.SetArgs(args)
		if err := root.Execute(); err == nil {
			t.Fatalf("selector succeeded: %q", args)
		}
		if output.Len() != 0 {
			t.Fatalf("default contacted for %q: %s", args, output.String())
		}
	}
}

func TestIPCSelectorConflictsPreserveExplicitEmptyIntent(t *testing.T) {
	endpoint := filepath.Join(t.TempDir(), "owned.sock")
	if runtime.GOOS == "windows" {
		endpoint = `\\.\pipe\owned-s107`
	}
	for _, flag := range []string{"profile", "endpoint", "daemon-id", "credential-id", "certificate-file"} {
		var diagnostics bytes.Buffer
		code := ExecuteWithDependencies(context.Background(), []string{"--ipc-path", endpoint, "--" + flag + "=", "health"}, Dependencies{Client: healthFixture("unused"), Out: io.Discard, Err: &diagnostics})
		if code != 2 || !strings.Contains(diagnostics.String(), "cannot be combined") {
			t.Fatalf("flag=%s code=%d stderr=%q", flag, code, diagnostics.String())
		}
	}
}

func TestGeneratedHelpAndCompletionDoNotResolveRemoteTargets(t *testing.T) {
	for _, args := range [][]string{{"--profile", "missing", "--help"}, {"--profile", "missing", "completion", "bash"}} {
		var output bytes.Buffer
		if code := ExecuteWithDependencies(context.Background(), args, Dependencies{Profiles: clientprofile.NewStore(filepath.Join(t.TempDir(), "empty.json")), Out: &output, Err: io.Discard}); code != 0 || output.Len() == 0 {
			t.Fatalf("args=%v code=%d output=%q", args, code, output.String())
		}
	}
}

func TestMCPInvocationCancellationClosesOwnedBlockingStreams(t *testing.T) {
	input, producer := io.Pipe()
	output, consumer := io.Pipe()
	defer producer.Close()
	defer output.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	finished := make(chan int, 1)
	go func() {
		finished <- ExecuteWithDependencies(ctx, []string{"mcp", "serve"}, Dependencies{Client: healthFixture("unused"), In: input, Out: consumer, Err: io.Discard})
	}()
	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		_ = input.Close()
		_ = consumer.Close()
		t.Fatal("owned MCP invocation did not stop after cancellation")
	}
	if _, err := producer.Write([]byte("closed")); err == nil {
		t.Fatal("owned input pipe remained open")
	}
}

func TestInvocationRequestUsesCanceledParent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	root := newRootWithDependencies(Dependencies{Client: healthFixture("unused"), Out: io.Discard, Err: io.Discard})
	root.SetArgs([]string{"health"})
	// Assert the exact handler context even for an in-process transport.
	root.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		cmd.SetContext(context.WithValue(cmd.Context(), invocationKey{}, stateFor(root)))
		got, stop := reqCtx(cmd)
		defer stop()
		return got.Err()
	}
	if err := root.ExecuteContext(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v", err)
	}
}

func TestRecordingHostAdaptersRunActualHandlers(t *testing.T) {
	var actions []string
	for _, action := range []string{"install", "start", "stop", "restart", "status", "uninstall"} {
		deps := Dependencies{Out: io.Discard, Err: io.Discard, DaemonPath: func() (string, error) { return "owned-goschedd", nil }, ServiceControl: func(got, exe string, args []string) (string, error) {
			actions = append(actions, got)
			if exe != "owned-goschedd" {
				t.Fatal(exe)
			}
			return "recorded", nil
		}}
		if code := ExecuteWithDependencies(context.Background(), []string{"service", action}, deps); code != 0 {
			t.Fatalf("%s exit=%d", action, code)
		}
	}
	launched := ""
	if code := ExecuteWithDependencies(context.Background(), []string{"gui"}, Dependencies{Out: io.Discard, Err: io.Discard, GUIPath: func() (string, error) { return "owned-gui", nil }, LaunchGUI: func(path string) error { launched = path; return nil }}); code != 0 || launched != "owned-gui" {
		t.Fatalf("launch=%q code=%d", launched, code)
	}
	if len(actions) != 6 {
		t.Fatalf("actions=%v", actions)
	}
	var diagnostics bytes.Buffer
	if code := ExecuteWithDependencies(context.Background(), []string{"service", "start"}, Dependencies{Out: io.Discard, Err: &diagnostics, DaemonPath: func() (string, error) { return "owned", nil }, ServiceControl: func(string, string, []string) (string, error) { return "", errors.New("recorded service unavailable") }}); code != 1 || !strings.Contains(diagnostics.String(), "recorded service unavailable") {
		t.Fatalf("code=%d stderr=%q", code, diagnostics.String())
	}
}

func TestLocalProfileCommandUsesOnlyScopedStore(t *testing.T) {
	store := clientprofile.NewStore(filepath.Join(t.TempDir(), "profiles.json"))
	var output bytes.Buffer
	if code := ExecuteWithDependencies(context.Background(), []string{"--profile", "missing", "--json", "profile", "list"}, Dependencies{Profiles: store, Out: &output, Err: io.Discard, Now: func() time.Time { return time.Unix(1, 0) }}); code != 0 || !strings.Contains(output.String(), `"profiles": []`) {
		t.Fatalf("code=%d output=%q", code, output.String())
	}
}
