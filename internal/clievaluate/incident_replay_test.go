package clievaluate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/clientprofile"
	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
	"github.com/shruggietech/go-schedule/v2/internal/incident"
	"github.com/shruggietech/go-schedule/v2/internal/ipc"
	"github.com/shruggietech/go-schedule/v2/internal/mcphttp"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

// TestWindowsIncidentReplay is an explicit diagnostic recipe, not a replacement
// for the complete source-bound CLI evaluation and its reviewed report.
func TestWindowsIncidentReplay(t *testing.T) {
	if os.Getenv("GOSCHEDULE_INCIDENT_REPLAY") != "1" {
		t.Skip("opt-in owned incident investigation")
	}
	if runtime.GOOS != "windows" && runtime.GOOS != "linux" {
		t.Skip("catalog execution supports Windows and Linux")
	}
	options, err := parseReplayOptions(os.Getenv("GOSCHEDULE_REPLAY_CASE_TIMEOUT"), os.Getenv("GOSCHEDULE_REPLAY_INCIDENTS"))
	if err != nil {
		t.Fatal(err)
	}
	logIncidentReplay(t, nil, map[string]any{"kind": "go-schedule.incident-replay-controller", "version": 1, "controller": "compiled-current-test", "case_timeout_ms": options.caseTimeout.Milliseconds(), "packaged_incidents_requested": options.incidents, "service_execution": "current-recording-command-harness", "component_server": "compiled-current-test", "packaged_server": "selected-build-receipt"})
	binaries := replayBinaries(t)
	_, definitions, err := BuildManifest()
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]step{}
	for _, definition := range definitions {
		byID[definition.Spec.ID] = definition
	}
	for _, mode := range []string{"component_memory", "component_file_wal", "packaged_file_wal"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var s *session
			var closeOwned func() error
			if mode == "packaged_file_wal" {
				setup, stop := context.WithTimeout(ctx, 120*time.Second)
				budget := DefaultBudget()
				budget.CaseTimeoutMS = options.caseTimeout.Milliseconds()
				s, err = newSession(setup, binaries, budget, NewRedactor(), options.incidents, nil)
				stop()
				if err != nil {
					t.Fatal(err)
				}
				closeOwned = s.close
			} else {
				s, closeOwned = replayComponent(t, mode, binaries)
			}
			s.budget.CaseTimeoutMS = options.caseTimeout.Milliseconds()
			t.Cleanup(func() {
				cleanupErr := closeOwned()
				_, statErr := os.Stat(s.w.Runtime)
				remaining := !os.IsNotExist(statErr)
				logIncidentReplay(t, nil, map[string]any{"kind": "go-schedule.incident-replay-cleanup", "version": 1, "test": t.Name(), "mode": mode, "runtime": s.w.Runtime, "retained": s.w.Evidence, "cleanup_error": errorText(cleanupErr), "runtime_remaining": remaining})
				if cleanupErr != nil || remaining {
					t.Errorf("owned replay cleanup: %v, runtime remaining=%t", cleanupErr, remaining)
				}
			})
			// These unchanged catalog prerequisites and targets preserve all output,
			// state, MCP protocol/confirmation and old-credential assertions.
			for _, id := range []string{"group.create", "group.child", "service.restart", "mcp.http.enable", "mcp.http.rotate", "mcp.bad-port", "group.remove-child"} {
				definition, ok := byID[id]
				if !ok {
					t.Fatalf("original catalog case %s is missing", id)
				}
				t.Run(id, func(t *testing.T) {
					caseCtx, stop := context.WithTimeout(ctx, options.caseTimeout)
					result := s.execute(caseCtx, definition)
					stop()
					passed := result.Blocked == "" && result.Code == definition.Spec.ExpectedExit && result.Cleanup == nil && len(result.Checks) == len(definition.Spec.Checks)
					for _, check := range result.Checks {
						passed = passed && check
					}
					logIncidentReplay(t, s.redactor, map[string]any{"kind": "go-schedule.incident-replay-observation", "version": 1, "test": t.Name(), "mode": mode, "case": id, "execution_kind": definition.Spec.ExecutionKind, "expected_exit": definition.Spec.ExpectedExit, "native_exit": result.Code, "checks": result.Checks, "blocked": result.Blocked, "detail": result.Detail, "cleanup_error": errorText(result.Cleanup), "passed": passed, "state": string(result.State), "stdout": string(result.Out), "stderr": string(result.Err), "private_evidence": s.w.Evidence})
					if !passed {
						t.Fatalf("original case %s failed without retry; retain the incident_replay observation", id)
					}
				})
				if t.Failed() {
					return
				}
			}
		})
		if t.Failed() {
			return
		}
	}
}

type replayOptions struct {
	caseTimeout time.Duration
	incidents   bool
}

func parseReplayOptions(timeout, incidents string) (replayOptions, error) {
	options := replayOptions{caseTimeout: time.Duration(DefaultBudget().CaseTimeoutMS) * time.Millisecond, incidents: true}
	if timeout != "" {
		parsed, err := time.ParseDuration(timeout)
		if err != nil || parsed < time.Millisecond || parsed%time.Millisecond != 0 {
			return options, fmt.Errorf("GOSCHEDULE_REPLAY_CASE_TIMEOUT must be a duration in whole milliseconds of at least 1ms: %q", timeout)
		}
		options.caseTimeout = parsed
	}
	switch incidents {
	case "", "1":
	case "0":
		options.incidents = false
	default:
		return options, fmt.Errorf("GOSCHEDULE_REPLAY_INCIDENTS must be 0 or 1: %q", incidents)
	}
	return options, nil
}

func replayBinaries(t *testing.T) map[string]string {
	t.Helper()
	directory := os.Getenv("GOSCHEDULE_REPLAY_BIN")
	if directory == "" {
		t.Fatal("GOSCHEDULE_REPLAY_BIN must name existing source-bound prepared inputs")
	}
	data, err := os.ReadFile(filepath.Join(directory, "build-receipt.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateSchema("build-receipt", data); err != nil {
		t.Fatal(err)
	}
	var receipt BuildReceipt
	if err := Decode(bytes.NewReader(data), &receipt); err != nil {
		t.Fatal(err)
	}
	actual, err := preparedBinaryDigests(directory)
	if err != nil {
		t.Fatal(err)
	}
	subject := Subject{Repository: receipt.Repository, Commit: receipt.Commit, Dirty: receipt.Dirty, SourceSHA256: receipt.SourceSHA256, Binaries: actual}
	if err := validateBuildReceipt(receipt, subject); err != nil {
		t.Fatal(err)
	}
	logIncidentReplay(t, nil, receipt)
	binaries := map[string]string{}
	for _, binary := range actual {
		if binary.OS == runtime.GOOS {
			binaries[binary.Role] = filepath.Join(directory, filepath.FromSlash(binary.Path))
		}
	}
	return binaries
}

func replayComponent(t *testing.T, mode string, binaries map[string]string) (*session, func() error) {
	t.Helper()
	w, err := clifixture.NewWorkspace("")
	if err != nil {
		t.Fatal(err)
	}
	// Register early cleanup for a setup failure. Successful setup closes through
	// the explicit ordered function below before removing the owned runtime.
	t.Cleanup(func() {
		if err := w.CleanupRuntime(); err != nil {
			t.Error(err)
		}
	})
	database := ":memory:"
	if mode == "component_file_wal" {
		database = w.Config.DBPath()
	}
	st, err := store.Open(database)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	listener, _, err := ipc.Listen(w.Config)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	info := server.RuntimeInfoResponse{IPCPath: w.IPCPath, ProcessID: os.Getpid(), DataDir: w.Config.DataDir, DatabasePath: w.Config.DBPath(), ConfigPath: w.ConfigPath, LogPath: w.Config.LogPath(), LockPath: filepath.Join(w.Config.DataDir, "goschedd.lock")}
	api := server.NewWithRuntimeInfo(st, nil, nil, nil, w.Config.LogPath(), info, nil)
	local := client.New(w.IPCPath)
	manager := mcphttp.NewWithSessions(local, Version, nil)
	api.SetMCPHTTPManager(manager)
	tracker := &incident.Tracker{}
	channel, err := incident.NewChannel(w.Evidence)
	if err != nil {
		t.Fatal(err)
	}
	diagnosticCtx, stopDiagnostics := context.WithCancel(context.Background())
	diagnosticDone := make(chan error, 1)
	go func() { diagnosticDone <- incident.Serve(diagnosticCtx, channel.Dir, tracker, st.PoolStats) }()
	httpServer := &http.Server{Handler: api.HandlerWithDiagnostics(tracker), ReadHeaderTimeout: time.Second}
	served := make(chan error, 1)
	go func() { served <- httpServer.Serve(listener) }()
	closeOwned := func() error {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()
		result := manager.Shutdown(cleanupCtx)
		stopDiagnostics()
		result = errors.Join(result, <-diagnosticDone)
		if err := httpServer.Shutdown(cleanupCtx); err != nil {
			result = errors.Join(result, err, httpServer.Close())
		}
		if err := <-served; !errors.Is(err, http.ErrServerClosed) {
			result = errors.Join(result, err)
		}
		return errors.Join(result, st.Close(), w.CleanupRuntime())
	}
	// The test host is the observed component process, never a process to stop.
	s := &session{w: w, api: local, profiles: clientprofile.NewStore(filepath.Join(w.Runtime, "client", "profiles.json")), recorder: &clifixture.Recorder{}, redactor: NewRedactor(), vars: map[string]string{"ipc": w.IPCPath}, binaries: binaries, budget: DefaultBudget(), incidents: channel, daemon: &clifixture.Process{PID: os.Getpid()}}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	if err := w.Bind(ctx, local, os.Getpid()); err != nil {
		if cleanupErr := closeOwned(); cleanupErr != nil {
			t.Error(cleanupErr)
		}
		t.Fatal(err)
	}
	address, err := availableAddress()
	if err != nil {
		if cleanupErr := closeOwned(); cleanupErr != nil {
			t.Error(cleanupErr)
		}
		t.Fatal(err)
	}
	_, s.vars["mcp_port"], err = net.SplitHostPort(address)
	if err != nil {
		if cleanupErr := closeOwned(); cleanupErr != nil {
			t.Error(cleanupErr)
		}
		t.Fatal(err)
	}
	logIncidentReplay(t, nil, map[string]any{"kind": "go-schedule.incident-replay-storage", "mode": mode, "actual_database": database, "logical_owned_database_path": info.DatabasePath, "pid": os.Getpid(), "evidence": w.Evidence})
	return s, closeOwned
}

func logIncidentReplay(t *testing.T, redactor *Redactor, value any) {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if redactor != nil {
		data = redactor.Sanitize(data)
	}
	t.Logf("incident_replay=%s", data)
}
