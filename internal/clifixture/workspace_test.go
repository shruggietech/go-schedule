package clifixture

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/clientsecret"
	"github.com/shruggietech/go-schedule/v2/internal/ipc"
)

func TestWorkspaceIdentityAndContainment(t *testing.T) {
	w, err := NewWorkspace("")
	if err != nil {
		t.Fatal(err)
	}
	info := server.RuntimeInfoResponse{IPCPath: w.IPCPath, ProcessID: 123, DataDir: w.Config.DataDir, DatabasePath: w.Config.DBPath(), ConfigPath: w.ConfigPath, LogPath: w.Config.LogPath(), LockPath: filepath.Join(w.Config.DataDir, "goschedd.lock")}
	installation := "owned-installation"
	api := client.NewInProcess(http.HandlerFunc(func(out http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/v1/manifest" {
			_ = json.NewEncoder(out).Encode(server.ManifestResponse{InstallationID: installation})
			return
		}
		_ = json.NewEncoder(out).Encode(info)
	}))
	// In-process transport is deliberately not accepted as an owned IPC target.
	if err := w.Bind(context.Background(), api, 123); err == nil {
		t.Fatal("in-process endpoint accepted as owned IPC")
	}
	if err := w.Contains(filepath.Join(w.Runtime, "future", "output")); err != nil {
		t.Fatal(err)
	}
	if err := w.Contains(filepath.Dir(w.Root)); err == nil {
		t.Fatal("outside parent accepted")
	}
	link := filepath.Join(w.Runtime, "escape")
	outside := t.TempDir()
	if err := os.Symlink(outside, link); err == nil {
		if err := w.Contains(filepath.Join(link, "output")); err == nil {
			t.Fatal("symlink escape accepted")
		}
	}
	if err := w.CleanupRuntime(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(w.Runtime); !os.IsNotExist(err) {
		t.Fatalf("runtime remains: %v", err)
	}
	if _, err := os.Stat(filepath.Join(w.Evidence, "ownership.json")); err != nil {
		t.Fatalf("evidence lost: %v", err)
	}
}

func TestGuardChecksEveryAuthoritativeOwnedField(t *testing.T) {
	w, err := NewWorkspace("")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = w.CleanupRuntime() })
	info := server.RuntimeInfoResponse{IPCPath: w.IPCPath, ProcessID: 123, DataDir: w.Config.DataDir, DatabasePath: w.Config.DBPath(), ConfigPath: w.ConfigPath, LogPath: w.Config.LogPath(), LockPath: filepath.Join(w.Config.DataDir, "goschedd.lock")}
	original := info
	var mu sync.Mutex
	installation := "owned-installation"
	listener, _, err := ipc.Listen(w.Config)
	if err != nil {
		t.Fatal(err)
	}
	httpServer := &http.Server{Handler: http.HandlerFunc(func(out http.ResponseWriter, req *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if req.URL.Path == "/v1/manifest" {
			_ = json.NewEncoder(out).Encode(server.ManifestResponse{InstallationID: installation})
			return
		}
		_ = json.NewEncoder(out).Encode(info)
	})}
	go func() { _ = httpServer.Serve(listener) }()
	t.Cleanup(func() { _ = httpServer.Close() })
	api := client.New(w.IPCPath)
	if err := w.Bind(context.Background(), api, 123); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"pid", "ipc", "data", "database", "config", "log", "lock", "installation"} {
		mu.Lock()
		info = original
		installation = "owned-installation"
		switch field {
		case "pid":
			info.ProcessID++
		case "ipc":
			info.IPCPath += "-wrong"
		case "data":
			info.DataDir = filepath.Dir(w.Root)
		case "database":
			info.DatabasePath = filepath.Join(w.Evidence, "wrong.db")
		case "config":
			info.ConfigPath = filepath.Join(w.Evidence, "wrong.json")
		case "log":
			info.LogPath = filepath.Join(w.Evidence, "wrong.log")
		case "lock":
			info.LockPath = filepath.Join(w.Evidence, "wrong.lock")
		case "installation":
			installation = "different-installation"
		}
		mu.Unlock()
		if err := w.Guard(context.Background(), api); err == nil {
			t.Errorf("%s mismatch permitted mutation", field)
		}
	}
}

func TestRecorderSecretsAndFailures(t *testing.T) {
	r := &Recorder{}
	store := clientsecret.NewWithBackend(r)
	if err := store.Probe(); err != nil {
		t.Fatal(err)
	}
	if err := store.Save("daemon", "credential", "do-not-print-token"); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load("daemon", "credential")
	if err != nil || got != "do-not-print-token" {
		t.Fatal(err)
	}
	for _, intent := range r.Intents() {
		if intent.Evidence != CommandHarness {
			t.Fatal(intent)
		}
	}
	if err := store.Delete("daemon", "credential"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Load("daemon", "credential"); err == nil {
		t.Fatal("deleted fixture secret remains")
	}
}
