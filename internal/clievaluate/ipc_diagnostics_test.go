package clievaluate

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/clifixture"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/ipc"
	"github.com/shruggietech/go-schedule/internal/store"
)

// This uses an owned OS IPC endpoint and real API/store, not the installed service.
func TestOwnedIPCSustainedSnapshotsRetainCompleteRequestEvidence(t *testing.T) {
	w, err := clifixture.NewWorkspace("")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := w.CleanupRuntime(); err != nil {
			t.Error(err)
		}
	}()
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task := domain.Task{Name: "owned request target", Command: "owned-fixture", State: domain.TaskActive, Timezone: "UTC"}
	if err := st.CreateTask(&task); err != nil {
		t.Fatal(err)
	}
	listener, _, err := ipc.Listen(w.Config)
	if err != nil {
		t.Fatal(err)
	}
	info := server.RuntimeInfoResponse{IPCPath: w.IPCPath, ProcessID: os.Getpid(), DataDir: w.Config.DataDir, DatabasePath: w.Config.DBPath(), ConfigPath: w.ConfigPath, LogPath: w.Config.LogPath(), LockPath: filepath.Join(w.Config.DataDir, "goschedd.lock")}
	api := server.NewWithRuntimeInfo(st, nil, nil, nil, w.Config.LogPath(), info, nil)
	var inject atomic.Bool
	var credentialRequests atomic.Int32
	handler := http.HandlerFunc(func(out http.ResponseWriter, request *http.Request) {
		if inject.Load() && request.URL.Path == "/v1/access/credentials" && credentialRequests.Add(1) == 2 {
			out.WriteHeader(http.StatusServiceUnavailable)
			_, _ = out.Write([]byte(`{"error":{"code":"internal","message":"owned inspection failure"}}`))
			return
		}
		api.Handler().ServeHTTP(out, request)
	})
	httpServer := &http.Server{Handler: handler, ReadHeaderTimeout: time.Second}
	done := make(chan error, 1)
	go func() { done <- httpServer.Serve(listener) }()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			t.Error(err)
		}
		if err := <-done; !errors.Is(err, http.ErrServerClosed) {
			t.Error(err)
		}
	}()
	s := &session{w: w, api: client.New(w.IPCPath), profiles: clientprofile.NewStore(filepath.Join(w.Runtime, "client", "profiles.json")), recorder: &clifixture.Recorder{}, redactor: NewRedactor(), vars: map[string]string{"ipc": w.IPCPath}, binaries: map[string]string{}, budget: DefaultBudget()}
	if err := w.Bind(context.Background(), s.api, os.Getpid()); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for worker := 0; worker < 2; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 16; i++ {
				d := &executionDiagnostics{}
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				snapshot, err := s.snapshot(d.phase(ctx, "sustained_snapshot"))
				cancel()
				if err != nil {
					t.Errorf("owned IPC snapshot %d: %v (%s)", i, err, d.finish(execution{}).State)
					return
				}
				if _, ok := fieldAt(snapshot, "tasks."+task.ID); !ok {
					t.Error("complete snapshot lost owned task")
					return
				}
				if len(d.requests) != 10 {
					t.Errorf("request observations %d, want 10", len(d.requests))
					return
				}
			}
		}()
	}
	wg.Wait()
	t.Log("32 complete snapshots, 320 typed requests over the owned native IPC endpoint, no retries")
	// A successful recorded host command must not become a failed command when
	// only its later inspection is unavailable. The fixture never calls the host manager.
	_, steps, err := BuildManifest()
	if err != nil {
		t.Fatal(err)
	}
	var definition step
	found := false
	for _, item := range steps {
		if item.Spec.ID == "service.restart" {
			definition = item
			found = true
			break
		}
	}
	if !found {
		t.Fatal("service restart scenario missing")
	}
	inject.Store(true)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result := s.execute(ctx, definition)
	if result.Code != 0 || result.Blocked != "" || len(result.Out) == 0 || len(result.Checks) < 2 || result.Checks[1] {
		t.Fatalf("command/inspection evidence conflated: %+v", result)
	}
	var state map[string]any
	if err := json.Unmarshal(result.State, &state); err != nil {
		t.Fatal(err)
	}
	if state["after_complete"] != false || state["before_complete"] != true {
		t.Fatalf("incomplete snapshot accepted: %s", result.State)
	}
	if value, ok := fieldAt(state, "partial_after.tasks."+task.ID+".name"); !ok || value != "owned request target" {
		t.Fatalf("partial data discarded: %s", result.State)
	}
	if !strings.Contains(string(result.State), "snapshot credentials") || !strings.Contains(string(result.State), "api_rejected") {
		t.Fatalf("failed collection/request missing: %s", result.State)
	}

}
