package gui

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/winuninstall"
)

type storageFileInfo struct{ name string }

func (i storageFileInfo) Name() string     { return i.name }
func (storageFileInfo) Size() int64        { return 0 }
func (storageFileInfo) Mode() fs.FileMode  { return 0 }
func (storageFileInfo) ModTime() time.Time { return time.Time{} }
func (storageFileInfo) IsDir() bool        { return false }
func (storageFileInfo) Sys() any           { return nil }

func TestResolveStorageLocationsClassifiesDeclaredPaths(t *testing.T) {
	root := filepath.Join(t.TempDir(), "machine")
	prefs := filepath.Join(t.TempDir(), "preferences")
	exe := filepath.Join(t.TempDir(), "install", "gosched-gui.exe")
	runtimeInfo := server.RuntimeInfoResponse{
		DataDir:      root,
		DatabasePath: filepath.Join(root, "goschedule.db"),
		ConfigPath:   filepath.Join(root, "config.json"),
		LogPath:      filepath.Join(root, "custom", "events.log"),
		LockPath:     filepath.Join(root, "goschedd.lock"),
	}

	var inspected []string
	exists := map[string]bool{
		root:                                     true,
		runtimeInfo.DatabasePath:                 true,
		filepath.Join(prefs, "preferences.json"): true,
		filepath.Dir(exe):                        true,
	}
	locations := resolveStorageLocations(storageLocationInputs{
		Runtime:              runtimeInfo,
		OwnedMachineDataRoot: root,
		PreferencesRoot:      prefs,
		ExecutablePath:       exe,
		GOOS:                 "linux",
		Stat: func(path string) (os.FileInfo, error) {
			inspected = append(inspected, path)
			if exists[path] {
				return storageFileInfo{name: filepath.Base(path)}, nil
			}
			return nil, os.ErrNotExist
		},
	})

	wantCategories := []string{
		"Machine data", "Task database", "Configuration", "Logs", "Runtime state",
		"Desktop application data", "Desktop preferences", "Executable directory",
	}
	if got := storageCategories(locations); !reflect.DeepEqual(got, wantCategories) {
		t.Fatalf("categories = %v, want %v", got, wantCategories)
	}
	if len(inspected) != len(locations)+1 { // absent installed docs is probed, then omitted
		t.Fatalf("inspected %d paths for %d rows, want one exact probe per row plus docs", len(inspected), len(locations))
	}
	for _, location := range locations {
		if !filepath.IsAbs(location.Path) {
			t.Errorf("%s path %q is not absolute", location.Category, location.Path)
		}
		if location.Scope == "" || location.SoftwareOnlyRemoval == "" || location.ExplicitDataWipe == "" {
			t.Errorf("%s lacks ownership/lifecycle copy: %+v", location.Category, location)
		}
		if strings.Contains(location.ExplicitDataWipe, "Removed by") {
			t.Errorf("%s claims a non-Windows data wipe: %+v", location.Category, location)
		}
	}
	assertStorageLocation(t, locations, "Machine data", storageScopeMachine, storagePresent)
	assertStorageLocation(t, locations, "Desktop preferences", storageScopeUser, storagePresent)
	assertStorageLocation(t, locations, "Logs", storageScopeMachine, storageAbsent)
}

func TestResolveStorageLocationsWindowsMaintenanceAndUnknownState(t *testing.T) {
	base := filepath.Join(t.TempDir(), "ProgramData")
	dataDir := filepath.Join(base, "goschedule")
	locations := resolveStorageLocations(storageLocationInputs{
		Runtime: server.RuntimeInfoResponse{
			DataDir:      dataDir,
			DatabasePath: filepath.Join(dataDir, "goschedule.db"),
			ConfigPath:   filepath.Join(dataDir, "config.json"),
			LogPath:      filepath.Join(dataDir, "logs", "goschedule.log"),
			LockPath:     filepath.Join(dataDir, "goschedd.lock"),
		},
		OwnedMachineDataRoot:    dataDir,
		PreferencesRoot:         filepath.Join(t.TempDir(), "prefs"),
		ExecutablePath:          filepath.Join(t.TempDir(), "gosched-gui.exe"),
		GOOS:                    "windows",
		MaintenanceEvidencePath: winuninstall.CleanupResultPath(base),
		Stat: func(path string) (os.FileInfo, error) {
			return nil, errors.New("access denied")
		},
	})
	assertStorageLocation(t, locations, "Maintenance evidence", storageScopeMachine, storageUnknown)
	for _, location := range locations {
		if location.Category == "Machine data" && location.ExplicitDataWipe != "Removed by an explicit data wipe" {
			t.Fatalf("Windows machine data wipe text = %q", location.ExplicitDataWipe)
		}
	}
}

func TestResolveStorageLocationsRejectsRelativeValues(t *testing.T) {
	locations := resolveStorageLocations(storageLocationInputs{
		Runtime: server.RuntimeInfoResponse{
			DataDir:      "relative-data",
			DatabasePath: "relative-data/db",
			ConfigPath:   "relative-config",
			LogPath:      "relative-log",
			LockPath:     "relative-lock",
		},
		OwnedMachineDataRoot: "relative-owned-root",
		PreferencesRoot:      "relative-preferences",
		ExecutablePath:       "relative-executable",
		GOOS:                 "linux",
		Stat: func(string) (os.FileInfo, error) {
			t.Fatal("relative unavailable paths must not be inspected")
			return nil, nil
		},
	})
	if len(locations) == 0 {
		t.Fatal("expected honest unavailable rows")
	}
	for _, location := range locations {
		if location.Available || location.Path != "" || location.Existence != storageUnknown {
			t.Errorf("relative path should be unavailable: %+v", location)
		}
	}
}

func TestStorageExistenceDoesNotTraverse(t *testing.T) {
	want := filepath.Join(t.TempDir(), "one", "exact", "path")
	calls := 0
	got := inspectStoragePath(want, func(path string) (os.FileInfo, error) {
		calls++
		if path != want {
			t.Fatalf("stat path = %q, want exact %q", path, want)
		}
		return nil, os.ErrNotExist
	})
	if got != storageAbsent || calls != 1 {
		t.Fatalf("inspectStoragePath() = %q after %d calls, want absent after 1", got, calls)
	}
}

func TestResolveStorageLocationsDoesNotClaimExternalLogAsWipeOwned(t *testing.T) {
	base := t.TempDir()
	owned := filepath.Join(base, "owned")
	configured := filepath.Join(base, "operator")
	locations := resolveStorageLocations(storageLocationInputs{
		Runtime: server.RuntimeInfoResponse{
			DataDir:      configured,
			DatabasePath: filepath.Join(configured, "goschedule.db"),
			ConfigPath:   filepath.Join(configured, "config.json"),
			LogPath:      filepath.Join(base, "operator-logs", "events.log"),
			LockPath:     filepath.Join(configured, "goschedd.lock"),
		},
		OwnedMachineDataRoot: owned,
		PreferencesRoot:      filepath.Join(base, "preferences"),
		ExecutablePath:       filepath.Join(base, "gosched-gui.exe"),
		GOOS:                 "linux",
		Stat:                 func(string) (os.FileInfo, error) { return nil, os.ErrNotExist },
	})
	for _, location := range locations {
		if location.Category != "Logs" {
			continue
		}
		if location.Scope != storageScopeExternal || !strings.Contains(location.ExplicitDataWipe, "Preserved") {
			t.Fatalf("external log incorrectly claimed as wipe-owned: %+v", location)
		}
		return
	}
	t.Fatal("missing Logs row")
}

func assertStorageLocation(t *testing.T, locations []storageLocation, category string, scope storageScope, existence storageExistence) {
	t.Helper()
	for _, location := range locations {
		if location.Category == category {
			if location.Scope != scope || location.Existence != existence {
				t.Fatalf("%s = %+v, want scope=%q existence=%q", category, location, scope, existence)
			}
			return
		}
	}
	t.Fatalf("missing storage location %q", category)
}
