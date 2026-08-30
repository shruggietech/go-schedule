//go:build !windows

package ipc

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/shruggietech/go-schedule/internal/config"
)

func currentGroup(t *testing.T) (string, int) {
	t.Helper()
	current, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	group, err := user.LookupGroupId(current.Gid)
	if err != nil {
		t.Fatal(err)
	}
	gid, err := strconv.Atoi(group.Gid)
	if err != nil {
		t.Fatal(err)
	}
	return group.Name, gid
}

func pathGID(t *testing.T, path string) int {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	return int(info.Sys().(*syscall.Stat_t).Gid)
}

func TestListenUnixRestrictedPermissions(t *testing.T) {
	group, gid := currentGroup(t)
	cfg := config.Default()
	cfg.DataDir = filepath.Join(t.TempDir(), "data")
	cfg.AdminGroup = group
	ln, access, err := Listen(cfg)
	if err != nil {
		t.Fatalf("Listen() error = %v", err)
	}
	defer ln.Close()
	endpoint := Endpoint(cfg)
	for path, wantMode := range map[string]os.FileMode{cfg.DataDir: 0o770, endpoint: 0o660} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got != wantMode {
			t.Errorf("%s mode = %04o, want %04o", path, got, wantMode)
		}
		if got := pathGID(t, path); got != gid {
			t.Errorf("%s gid = %d, want %d", path, got, gid)
		}
	}
	if access != (AccessInfo{Mode: AccessModeRestricted, AdminGroup: group}) {
		t.Fatalf("access = %+v", access)
	}
}

func TestListenUnixMissingGroupFailsClosed(t *testing.T) {
	cfg := config.Default()
	cfg.DataDir = filepath.Join(t.TempDir(), "data")
	cfg.AdminGroup = "group-that-does-not-exist-for-goschedule-tests"
	_, _, err := Listen(cfg)
	if err == nil || !strings.Contains(err.Error(), "admin_group") || !strings.Contains(err.Error(), cfg.AdminGroup) {
		t.Fatalf("error = %v, want actionable admin_group failure", err)
	}
	if _, statErr := os.Stat(Endpoint(cfg)); !os.IsNotExist(statErr) {
		t.Fatalf("endpoint exists after lookup failure: %v", statErr)
	}
}

func TestListenUnixDoesNotRewriteUnsafeCustomParent(t *testing.T) {
	group, _ := currentGroup(t)
	parent := filepath.Join(t.TempDir(), "custom")
	if err := os.Mkdir(parent, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(parent, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := config.Default()
	cfg.IPCPath = filepath.Join(parent, "goschedd.sock")
	cfg.AdminGroup = group
	_, _, err := Listen(cfg)
	if err == nil || !strings.Contains(err.Error(), "verify socket dir") {
		t.Fatalf("Listen() error = %v, want unsafe custom directory rejection", err)
	}
	info, statErr := os.Stat(parent)
	if statErr != nil {
		t.Fatal(statErr)
	}
	if got := info.Mode().Perm(); got != 0o755 {
		t.Fatalf("custom parent mode changed to %04o", got)
	}
}

func TestListenUnixCompatibilityMode(t *testing.T) {
	cfg := config.Default()
	cfg.DataDir = filepath.Join(t.TempDir(), "data")
	cfg.AdminGroup = ""
	ln, access, err := Listen(cfg)
	if err != nil {
		t.Fatalf("Listen() error = %v", err)
	}
	defer ln.Close()
	for path, wantMode := range map[string]os.FileMode{cfg.DataDir: 0o755, Endpoint(cfg): 0o666} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got != wantMode {
			t.Errorf("%s mode = %04o, want %04o", path, got, wantMode)
		}
	}
	if access != (AccessInfo{Mode: AccessModeCompatibility}) {
		t.Fatalf("access = %+v", access)
	}
}

func TestListenUnixSecuresEveryCreatedCustomParent(t *testing.T) {
	group, gid := currentGroup(t)
	root := t.TempDir()
	first := filepath.Join(root, "first")
	parent := filepath.Join(first, "second")
	cfg := config.Default()
	cfg.IPCPath = filepath.Join(parent, "goschedd.sock")
	cfg.AdminGroup = group
	ln, _, err := Listen(cfg)
	if err != nil {
		t.Fatalf("Listen() error = %v", err)
	}
	defer ln.Close()
	for _, path := range []string{first, parent} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got != 0o770 {
			t.Errorf("%s mode = %04o, want 0770", path, got)
		}
		if got := pathGID(t, path); got != gid {
			t.Errorf("%s gid = %d, want %d", path, got, gid)
		}
	}
}

func TestListenUnixRemovesStaleSocket(t *testing.T) {
	group, _ := currentGroup(t)
	cfg := config.Default()
	cfg.DataDir = filepath.Join(t.TempDir(), "data")
	cfg.AdminGroup = group
	first, _, err := Listen(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, _, err := Listen(cfg)
	if err != nil {
		t.Fatalf("Listen() with stale socket error = %v", err)
	}
	second.Close()
}

func TestListenUnixCleansUpAfterPermissionFailure(t *testing.T) {
	group, _ := currentGroup(t)
	cfg := config.Default()
	cfg.DataDir = filepath.Join(t.TempDir(), "data")
	cfg.AdminGroup = group
	endpoint := Endpoint(cfg)
	ops := productionUnixIPCOperations
	realChmod := ops.chmod
	ops.chmod = func(path string, mode os.FileMode) error {
		if path == endpoint {
			return errors.New("injected chmod failure")
		}
		return realChmod(path, mode)
	}
	_, _, err := listenUnix(cfg, ops)
	if err == nil || !strings.Contains(err.Error(), "injected chmod failure") {
		t.Fatalf("listenUnix() error = %v", err)
	}
	if _, statErr := os.Stat(endpoint); !os.IsNotExist(statErr) {
		t.Fatalf("endpoint remains after failure: %v", statErr)
	}
}
