package incident

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDaemonEvidenceDoesNotNeedTheBusySQLConnection(t *testing.T) {
	channel, err := NewChannel(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	tracker := &Tracker{}
	finish := tracker.Begin("GET", "GET /v1/tasks")
	defer finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- Serve(ctx, channel.Dir, tracker, func() sql.DBStats { return sql.DBStats{InUse: 1, WaitCount: 7} })
	}()
	defer func() {
		cancel()
		if err := <-done; err != nil {
			t.Error(err)
		}
	}()
	requestCtx, stop := context.WithTimeout(context.Background(), 2*time.Second)
	defer stop()
	response, err := channel.Collect(requestCtx, os.Getpid(), "before")
	if err != nil {
		t.Fatal(err)
	}
	if response.Target.PID != os.Getpid() || response.Pool.InUse != 1 || response.Pool.WaitCount != 7 || len(response.Active) != 1 || response.Active[0].Route != "GET /v1/tasks" {
		t.Fatalf("lost live progress: %+v", response)
	}
	if response.Stack.Status != "collected" && response.Stack.Status != "truncated" {
		t.Fatalf("stack unavailable: %+v", response.Stack)
	}
	if err := channel.Validate(response, os.Getpid()+1); err == nil {
		t.Fatal("accepted different daemon")
	}
}

func TestLiveHandlerStallIsDistinctFromSQLPoolWait(t *testing.T) {
	channel, err := NewChannel(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	tracker := &Tracker{}
	started := make(chan struct{})
	release := make(chan struct{})
	handler := tracker.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Stage(r.Context(), "handler")
		close(started)
		<-release
		w.WriteHeader(204)
	}), func(*http.Request) string { return "GET /v1/tasks" })
	requestDone := make(chan struct{})
	go func() {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/v1/tasks?secret=not-recorded", nil))
		close(requestDone)
	}()
	<-started
	defer func() { close(release); <-requestDone }()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Serve(ctx, channel.Dir, tracker, func() sql.DBStats { return sql.DBStats{Idle: 1} }) }()
	defer func() {
		cancel()
		if err := <-done; err != nil {
			t.Error(err)
		}
	}()
	collect, stop := context.WithTimeout(context.Background(), 2*time.Second)
	defer stop()
	snapshot, err := channel.Collect(collect, os.Getpid(), "after")
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Active) != 1 || snapshot.Active[0].Stage != "handler" || snapshot.Pool.InUse != 0 || snapshot.Pool.WaitCount != 0 {
		t.Fatalf("conflated handler and SQL wait: %+v", snapshot)
	}
	if strings.Contains(snapshot.Active[0].Route, "secret") {
		t.Fatal("request query leaked")
	}
}

func TestDaemonRejectsWrongTargetAndOversizedRequest(t *testing.T) {
	for _, kind := range []string{"wrong_pid", "oversized"} {
		t.Run(kind, func(t *testing.T) {
			channel, err := NewChannel(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			data := []byte(strings.Repeat("x", 1025))
			if kind == "wrong_pid" {
				data = []byte(`{"id":"aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa","session":"` + channel.ID + `","pid":-1,"phase":"before"}`)
			}
			if err := channel.Write("1.request.json", data); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := Serve(ctx, channel.Dir, &Tracker{}, func() sql.DBStats { return sql.DBStats{} }); err == nil {
				t.Fatal("accepted malformed target/request")
			}
		})
	}
}

func TestDaemonDiagnosticsRequireTheOwnedWorkspace(t *testing.T) {
	root := filepath.Join(t.TempDir(), "gosched-aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa")
	evidence := filepath.Join(root, "evidence")
	state := filepath.Join(root, "runtime", "state")
	for _, dir := range []string{evidence, state} {
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatal(err)
		}
	}
	// Canonicalize Windows short temporary aliases just as Workspace does.
	var err error
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	evidence = filepath.Join(root, "evidence")
	state = filepath.Join(root, "runtime", "state")
	manifest, _ := json.Marshal(map[string]string{"run_id": "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", "root": root, "runtime": filepath.Join(root, "runtime"), "evidence": evidence})
	if err := os.WriteFile(filepath.Join(evidence, "ownership.json"), manifest, 0600); err != nil {
		t.Fatal(err)
	}
	channel, err := NewChannel(evidence)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateWorkspace(channel.Dir, state); err != nil {
		t.Fatal(err)
	}
	if err := ValidateWorkspace(channel.Dir, filepath.Join(t.TempDir(), "installed-state")); err == nil {
		t.Fatal("accepted a different state directory")
	}
	if err := os.WriteFile(filepath.Join(evidence, "ownership.json"), []byte(`{"run_id":"changed"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := ValidateWorkspace(channel.Dir, state); err == nil {
		t.Fatal("accepted changed ownership identity")
	}
}

func TestDaemonCollectionIsBoundedAbsentAndCanceled(t *testing.T) {
	channel, err := NewChannel(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := channel.Collect(ctx, os.Getpid(), "before"); err == nil {
		t.Fatal("accepted canceled collection")
	}
	ctx, stop := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer stop()
	if _, err := channel.Collect(ctx, os.Getpid(), "before"); err == nil {
		t.Fatal("missing responder reported success")
	}
	tracker := &Tracker{}
	var ends []func()
	for range 70 {
		ends = append(ends, tracker.Begin("SECRET-METHOD", "GET /v1/tasks/{id}"))
	}
	active, omitted := tracker.Snapshot()
	if len(active) != 64 || omitted != 6 {
		t.Fatalf("unbounded tracker: %d/%d", len(active), omitted)
	}
	for _, a := range active {
		if strings.Contains(a.Method, "SECRET") {
			t.Fatal("arbitrary method leaked")
		}
	}
	for _, end := range ends {
		end()
	}
	active, _ = tracker.Snapshot()
	if len(active) != 0 {
		t.Fatal("completed request retained")
	}
}
