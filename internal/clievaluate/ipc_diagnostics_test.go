package clievaluate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	var holdCredentials atomic.Bool
	var creates, deletes atomic.Int32
	heldStarted := make(chan struct{}, 1)
	heldCanceled := make(chan struct{}, 1)
	handler := http.HandlerFunc(func(out http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/v1/groups" && request.Method == http.MethodPost {
			creates.Add(1)
		}
		if strings.HasPrefix(request.URL.Path, "/v1/groups/") && request.Method == http.MethodDelete {
			deletes.Add(1)
		}
		if holdCredentials.Load() && request.URL.Path == "/v1/access/credentials" {
			heldStarted <- struct{}{}
			<-request.Context().Done()
			heldCanceled <- struct{}{}
			return
		}
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
	t.Run("held_request_cancellation", func(t *testing.T) {
		holdCredentials.Store(true)
		defer holdCredentials.Store(false)
		deadline, stop := context.WithTimeout(context.Background(), 15*time.Second)
		defer stop()
		ctx, cancel := context.WithCancel(deadline)
		defer cancel()
		cancelDone := make(chan struct{})
		go func() {
			defer close(cancelDone)
			select {
			case <-heldStarted:
				cancel()
			case <-deadline.Done():
			}
		}()
		d := &executionDiagnostics{}
		partial, err := s.snapshot(d.phase(ctx, "held_snapshot"))
		<-cancelDone
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("original cancellation lost: %v", err)
		}
		select {
		case <-heldCanceled:
		case <-deadline.Done():
			t.Fatal("server did not receive canceled request context")
		}
		if _, ok := fieldAt(partial, "tasks."+task.ID); !ok {
			t.Fatal("cancellation discarded preceding task collection")
		}
		found := false
		for _, r := range d.requests {
			if r.Route == "/v1/access/credentials" && r.Outcome == "canceled" {
				found = true
			}
		}
		if !found {
			t.Fatalf("request cancellation not attributed: %+v", d.requests)
		}
		holdCredentials.Store(false)
		fresh, err := s.snapshot(deadline)
		if err != nil {
			t.Fatalf("fresh snapshot after cancellation: %v", err)
		}
		if _, ok := fieldAt(fresh, "tasks."+task.ID); !ok {
			t.Fatal("fresh complete snapshot lost task")
		}
	})
	t.Run("mixed_reads_and_mutations", func(t *testing.T) {
		var workers sync.WaitGroup
		for worker := 0; worker < 2; worker++ {
			workers.Add(1)
			go func(worker int) {
				defer workers.Done()
				for i := 0; i < 8; i++ {
					ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
					group, err := s.api.CreateGroup(ctx, server.GroupCreateRequest{Name: fmt.Sprintf("owned-mixed-%d-%d", worker, i)})
					if err == nil {
						_, err = s.api.ListGroups(ctx)
					}
					if err == nil {
						err = s.api.DeleteGroup(ctx, group.ID)
					}
					cancel()
					if err != nil {
						t.Errorf("mixed case %d/%d: %v", worker, i, err)
						return
					}
				}
			}(worker)
		}
		workers.Wait()
		if creates.Load() != 16 || deletes.Load() != 16 {
			t.Errorf("mutations repeated or missing: create=%d delete=%d", creates.Load(), deletes.Load())
		}
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		groups, err := s.api.ListGroups(ctx)
		if err != nil || len(groups) != 0 {
			t.Fatalf("mixed mutation state: groups=%+v error=%v", groups, err)
		}
	})

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
