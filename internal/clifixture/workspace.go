// Package clifixture owns disposable CLI evaluation state and process trees.
// It never discovers or controls installed services or existing desktop sessions.
package clifixture

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/config"
)

// EvidenceKind identifies which execution boundary produced an observation.
type EvidenceKind string

const (
	PackagedCLI    EvidenceKind = "packaged_cli"
	CommandHarness EvidenceKind = "command_harness"
	EngineFixture  EvidenceKind = "deterministic_engine"
)

// Workspace separates retained evidence from disposable runtime state.
type Workspace struct {
	ownedRoot, ownedID, boundDaemonID                string
	boundPID                                         int
	ID, Root, Runtime, Evidence, ConfigPath, IPCPath string
	Config                                           config.Config
	PID                                              int
	DaemonID                                         string
}

// NewWorkspace creates a new UUID namespace. The parent must already exist.
func NewWorkspace(parent string) (*Workspace, error) {
	if parent == "" {
		parent = os.TempDir()
		// macOS per-user temporary paths can exceed sockaddr_un capacity.
		// A UUID directory under canonical /private/tmp remains private and owned.
		if runtime.GOOS == "darwin" {
			parent = "/tmp"
		}
	}
	resolved, err := filepath.EvalSymlinks(parent)
	if err != nil {
		return nil, err
	}
	resolved, err = filepath.Abs(resolved)
	if err != nil {
		return nil, err
	}
	id := uuid.NewString()
	root := filepath.Join(resolved, "gosched-"+id)
	socketPath := filepath.Join(root, "runtime", "ipc", "daemon.sock")
	limit := 108
	if runtime.GOOS == "darwin" {
		limit = 104
	}
	if runtime.GOOS != "windows" && len([]byte(socketPath)) >= limit {
		return nil, errors.New("owned IPC path is too long; choose a shorter workspace parent directory")
	}
	if err := os.Mkdir(root, 0700); err != nil {
		return nil, err
	}
	w := &Workspace{ID: id, Root: root, Runtime: filepath.Join(root, "runtime"), Evidence: filepath.Join(root, "evidence"), ownedRoot: root, ownedID: id}
	for _, path := range []string{w.Runtime, w.Evidence, filepath.Join(w.Runtime, "state"), filepath.Join(w.Runtime, "ipc"), filepath.Join(w.Runtime, "client"), filepath.Join(w.Runtime, "watchers"), filepath.Join(w.Runtime, "exports"), filepath.Join(w.Runtime, "tasks"), filepath.Join(w.Runtime, "certificates")} {
		if err := os.Mkdir(path, 0700); err != nil {
			return nil, err
		}
	}
	w.Config = config.Default()
	w.Config.DataDir = filepath.Join(w.Runtime, "state")
	w.Config.AdminGroup = ""
	w.Config.WorkerPoolSize = 1
	w.Config.OutputCapBytes = 32768
	w.Config.LogMaxSizeBytes = 262144
	w.Config.LogMaxFiles = 2
	w.Config.LogRingSize = 128
	w.IPCPath = filepath.Join(w.Runtime, "ipc", "daemon.sock")
	if runtime.GOOS == "windows" {
		w.IPCPath = `\\.\pipe\gosched-qa-` + id
	}
	w.Config.IPCPath = w.IPCPath
	w.ConfigPath = filepath.Join(w.Runtime, "config.json")
	data, err := json.MarshalIndent(w.Config, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(w.ConfigPath, append(data, '\n'), 0600); err != nil {
		return nil, err
	}
	manifest, err := json.MarshalIndent(map[string]string{"run_id": id, "root": root, "runtime": w.Runtime, "evidence": w.Evidence, "ipc_path": w.IPCPath}, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(w.Evidence, "ownership.json"), append(manifest, '\n'), 0600); err != nil {
		return nil, err
	}
	return w, nil
}

// Contains rejects traversal and every existing symlink or junction escape.
// Nonexistent leaves are checked through their nearest existing ancestor.
func (w *Workspace) Contains(path string) error {
	if w == nil || w.Root == "" || w.ownedRoot == "" || !samePath(w.Root, w.ownedRoot) || w.ID != w.ownedID || !filepath.IsAbs(path) {
		return errors.New("owned workspace requires an absolute path")
	}
	root, err := filepath.EvalSymlinks(w.Root)
	if err != nil {
		return err
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return err
	}
	if !samePath(root, w.Root) {
		return errors.New("owned workspace root changed")
	}
	clean := filepath.Clean(path)
	rel, err := filepath.Rel(root, clean)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path is outside owned workspace: %s", path)
	}
	candidate := clean
	for {
		_, err := os.Lstat(candidate)
		if err == nil {
			break
		}
		if !os.IsNotExist(err) {
			return err
		}
		next := filepath.Dir(candidate)
		if next == candidate {
			return errors.New("owned path has no existing ancestor")
		}
		candidate = next
	}
	actual, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return err
	}
	if !samePath(actual, candidate) {
		return fmt.Errorf("owned path traverses a link: %s", path)
	}
	return nil
}

func samePath(a, b string) bool {
	if runtime.GOOS == "windows" {
		return strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
	}
	return filepath.Clean(a) == filepath.Clean(b)
}

// Bind records the expected identity once the run-owned child is ready.
func (w *Workspace) Bind(ctx context.Context, api *client.Client, pid int) error {
	if pid <= 0 {
		return errors.New("owned daemon PID is missing")
	}
	if w.boundPID != 0 {
		return errors.New("owned daemon identity is already bound")
	}
	if api == nil || api.Remote() || api.Endpoint() != w.IPCPath {
		return errors.New("target is not this run's owned local daemon")
	}
	if err := w.Contains(w.ConfigPath); err != nil {
		return err
	}
	manifest, err := api.Manifest(ctx)
	if err != nil {
		return err
	}
	if manifest.InstallationID == "" {
		return errors.New("owned daemon installation identity is missing")
	}
	candidate := *w
	candidate.PID = pid
	candidate.DaemonID = manifest.InstallationID
	candidate.boundPID = pid
	candidate.boundDaemonID = manifest.InstallationID
	if err := candidate.Guard(ctx, api); err != nil {
		return err
	}
	*w = candidate
	return nil
}

// Guard must precede each fixture mutation. A mismatch cannot fall back.
func (w *Workspace) Guard(ctx context.Context, api *client.Client) error {
	if w.PID <= 0 || w.PID != w.boundPID || w.DaemonID == "" || w.DaemonID != w.boundDaemonID || api == nil || api.Remote() || api.Endpoint() != w.IPCPath {
		return errors.New("target is not this run's owned local daemon")
	}
	manifest, err := api.Manifest(ctx)
	if err != nil {
		return err
	}
	if manifest.InstallationID != w.DaemonID {
		return errors.New("owned daemon installation identity changed")
	}
	info, err := api.RuntimeInfo(ctx)
	if err != nil {
		return err
	}
	if info.ProcessID != w.PID || info.IPCPath != w.IPCPath {
		return errors.New("owned daemon process or IPC identity changed")
	}
	for _, pair := range [][2]string{{info.DataDir, w.Config.DataDir}, {info.DatabasePath, w.Config.DBPath()}, {info.ConfigPath, w.ConfigPath}, {info.LogPath, w.Config.LogPath()}, {info.LockPath, filepath.Join(w.Config.DataDir, "goschedd.lock")}} {
		if !samePath(pair[0], pair[1]) {
			return fmt.Errorf("owned daemon state path changed: %s", pair[0])
		}
		if err := w.Contains(pair[0]); err != nil {
			return err
		}
	}
	return nil
}

// CleanupRuntime removes only disposable state; evidence and ownership survive.
func (w *Workspace) CleanupRuntime() error {
	if err := w.Contains(w.Runtime); err != nil {
		return err
	}
	if !samePath(w.Runtime, filepath.Join(w.Root, "runtime")) {
		return errors.New("runtime cleanup root changed")
	}
	return os.RemoveAll(w.Runtime)
}
