package clievaluate

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/clifixture"
	"github.com/shruggietech/go-schedule/internal/ipc"
	"github.com/shruggietech/go-schedule/internal/store"
)

func TestOwnedMetadataUsesSelectedCompletionAllowance(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("source location unavailable")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	binary := filepath.Join(t.TempDir(), "gosched")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	goPath, err := exec.LookPath("go")
	if err != nil {
		t.Fatal(err)
	}
	buildCtx, cancelBuild := context.WithTimeout(context.Background(), 10*time.Minute)
	built := captureCommand(buildCtx, clifixture.Command{Executable: goPath, Arguments: []string{"build", "-o", binary, "./cmd/gosched"}, Directory: root, Environment: quietEnvironment()}, nil, DefaultBudget().ChannelBytes)
	cancelBuild()
	if built.Code != 0 || built.Failure != nil || built.Cleanup != nil {
		t.Fatalf("CLI build: exit=%d failure=%v cleanup=%v stderr=%s", built.Code, built.Failure, built.Cleanup, built.Err)
	}
	bytes, err := os.ReadFile(binary)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("native metadata CLI sha256=%s runtime=%s platform=%s budget=%+v", Digest(bytes), runtime.Version(), runtime.GOOS, DefaultBudget())
	identityCtx, cancelIdentity := context.WithTimeout(context.Background(), 120*time.Second)
	source, err := identifySource(identityCtx, root)
	cancelIdentity()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("source commit=%s dirty=%v sha256=%s", source.Commit, source.Dirty, source.SourceSHA256)
	s, delay := ownedMetadataSession(t, binary)
	_, steps, err := BuildManifest()
	if err != nil {
		t.Fatal(err)
	}
	selected := map[string]bool{"help.audit": true, "help.profile.default.reset": true, "help.profile.default.set": true}
	count := 0
	var example step
	for _, definition := range steps {
		if !selected[strings.TrimSuffix(definition.Spec.ID, ".human")] {
			continue
		}
		example = definition
		count++
		t.Run(definition.Spec.ID, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.budget.CaseTimeoutMS)*time.Millisecond)
			defer cancel()
			result := s.execute(ctx, definition)
			assertMetadataCompletion(t, result)
			t.Logf("case=%s native_exit=%d state=%s", definition.Spec.ID, result.Code, result.State)
		})
	}
	if count != 3 {
		t.Fatalf("metadata case matrix=%d, want three original metadata definitions", count)
	}
	// This diagnostic deliberately consumes four seconds in each of the four
	// ownership requests. It exercises real OS IPC and the old combined cap,
	// without changing a production per-request deadline or bypassing a guard.
	t.Run("completion_beyond_former_case_cap", func(t *testing.T) {
		delay.Store(int64(4 * time.Second))
		defer delay.Store(0)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.budget.CaseTimeoutMS)*time.Millisecond)
		defer cancel()
		started := time.Now()
		result := s.execute(ctx, example)
		elapsed := time.Since(started)
		assertMetadataCompletion(t, result)
		if elapsed < 16*time.Second {
			t.Fatalf("controlled four-request preparation was not exercised: %s", elapsed)
		}
		t.Logf("controlled beyond-old15s case elapsed=%s native_exit=%d state=%s", elapsed, result.Code, result.State)
	})
	t.Run("explicit_short_budget_still_fails", func(t *testing.T) {
		delay.Store(int64(4 * time.Second))
		defer delay.Store(0)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		result := s.execute(ctx, example)
		var state struct {
			Diagnostics struct {
				CommandRan bool            `json:"command_ran"`
				Requests   []requestTiming `json:"requests"`
			} `json:"diagnostics"`
		}
		if err := json.Unmarshal(result.State, &state); err != nil {
			t.Fatal(err)
		}
		if result.Blocked != "ownership_mismatch" || !strings.Contains(result.Detail, "context deadline exceeded") || state.Diagnostics.CommandRan || len(state.Diagnostics.Requests) != 1 {
			t.Fatalf("explicit limit lost its original guard failure: %+v state=%s", result, result.State)
		}
		t.Logf("explicit100ms command_ran=false original_diagnostic=%s state=%s", result.Detail, result.State)
	})
}

func assertMetadataCompletion(t *testing.T, result execution) {
	t.Helper()
	var state struct {
		MetadataOnly   bool `json:"metadata_only"`
		ProcessCleanup bool `json:"process_cleanup"`
		Diagnostics    struct {
			CommandRan bool            `json:"command_ran"`
			Requests   []requestTiming `json:"requests"`
		} `json:"diagnostics"`
	}
	if err := json.Unmarshal(result.State, &state); err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 || result.Blocked != "" || result.Detail != "" || result.Cleanup != nil || !strings.Contains(string(result.Out), "Usage:") || !state.MetadataOnly || !state.ProcessCleanup || !state.Diagnostics.CommandRan || len(state.Diagnostics.Requests) != 4 {
		t.Fatalf("metadata result=%+v state=%s", result, result.State)
	}
	for _, passed := range result.Checks {
		if !passed {
			t.Fatalf("metadata assertion failed: %s", result.State)
		}
	}
	for _, request := range state.Diagnostics.Requests {
		if request.Outcome != "completed" {
			t.Fatalf("ownership request did not complete: %+v", request)
		}
	}
}

func ownedMetadataSession(t *testing.T, binary string) (*session, *atomic.Int64) {
	t.Helper()
	base := os.TempDir()
	if runtime.GOOS == "darwin" {
		base = "/tmp"
	}
	parent, err := os.MkdirTemp(base, "qm-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(parent); err != nil {
			t.Error(err)
		}
	})
	w, err := clifixture.NewWorkspace(parent)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := w.CleanupRuntime(); err != nil {
			t.Error(err)
		}
	})
	database, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := database.Close(); err != nil {
			t.Error(err)
		}
	})
	listener, _, err := ipc.Listen(w.Config)
	if err != nil {
		t.Fatal(err)
	}
	info := server.RuntimeInfoResponse{IPCPath: w.IPCPath, ProcessID: os.Getpid(), DataDir: w.Config.DataDir, DatabasePath: w.Config.DBPath(), ConfigPath: w.ConfigPath, LogPath: w.Config.LogPath(), LockPath: filepath.Join(w.Config.DataDir, "goschedd.lock")}
	api := server.NewWithRuntimeInfo(database, nil, nil, nil, w.Config.LogPath(), info, nil)
	delay := &atomic.Int64{}
	handler := http.HandlerFunc(func(out http.ResponseWriter, request *http.Request) {
		if duration := time.Duration(delay.Load()); duration > 0 && (request.URL.Path == "/v1/manifest" || request.URL.Path == "/v1/runtime-info") {
			timer := time.NewTimer(duration)
			defer timer.Stop()
			select {
			case <-timer.C:
			case <-request.Context().Done():
				return
			}
		}
		api.Handler().ServeHTTP(out, request)
	})
	httpServer := &http.Server{Handler: handler, ReadHeaderTimeout: time.Second}
	done := make(chan error, 1)
	go func() { done <- httpServer.Serve(listener) }()
	t.Cleanup(func() {
		if err := httpServer.Close(); err != nil {
			t.Error(err)
		}
		if err := <-done; !errors.Is(err, http.ErrServerClosed) {
			t.Error(err)
		}
	})
	s := &session{w: w, api: client.New(w.IPCPath), profiles: clientprofile.NewStore(filepath.Join(w.Runtime, "client", "profiles.json")), recorder: &clifixture.Recorder{}, redactor: NewRedactor(), vars: map[string]string{"ipc": w.IPCPath}, binaries: map[string]string{"cli": binary}, budget: DefaultBudget()}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	if err := w.Bind(ctx, s.api, os.Getpid()); err != nil {
		t.Fatal(err)
	}
	return s, delay
}
