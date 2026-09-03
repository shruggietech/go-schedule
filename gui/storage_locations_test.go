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

	"github.com/shruggietech/go-schedule/internal/config"
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
	cfg := config.Default()
	cfg.DataDir = root
	cfg.LogFilePath = filepath.Join(root, "custom", "events.log")

	var inspected []string
	exists := map[string]bool{
		root:                                     true,
		cfg.DBPath():                             true,
		filepath.Join(prefs, "preferences.json"): true,
		filepath.Dir(exe):                        true,
	}
	locations := resolveStorageLocations(storageLocationInputs{
		Config:          cfg,
		PreferencesRoot: prefs,
		ExecutablePath:  exe,
		GOOS:            "linux",
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
	}
	assertStorageLocation(t, locations, "Machine data", storageScopeMachine, storagePresent)
	assertStorageLocation(t, locations, "Desktop preferences", storageScopeUser, storagePresent)
	assertStorageLocation(t, locations, "Logs", storageScopeMachine, storageAbsent)
}

func TestResolveStorageLocationsWindowsMaintenanceAndUnknownState(t *testing.T) {
	base := filepath.Join(t.TempDir(), "ProgramData")
	cfg := config.Default()
	cfg.DataDir = filepath.Join(base, "goschedule")
	locations := resolveStorageLocations(storageLocationInputs{
		Config:                  cfg,
		PreferencesRoot:         filepath.Join(t.TempDir(), "prefs"),
		ExecutablePath:          filepath.Join(t.TempDir(), "gosched-gui.exe"),
		GOOS:                    "windows",
		MaintenanceEvidencePath: winuninstall.CleanupResultPath(base),
		Stat: func(path string) (os.FileInfo, error) {
			return nil, errors.New("access denied")
		},
	})
	assertStorageLocation(t, locations, "Maintenance evidence", storageScopeMachine, storageUnknown)
}

func TestResolveStorageLocationsRejectsRelativeValues(t *testing.T) {
	cfg := config.Default()
	cfg.DataDir = "relative-data"
	locations := resolveStorageLocations(storageLocationInputs{
		Config:          cfg,
		PreferencesRoot: "relative-preferences",
		ExecutablePath:  "relative-executable",
		GOOS:            "linux",
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
	cfg := config.Default()
	cfg.DataDir = filepath.Join(base, "owned")
	cfg.LogFilePath = filepath.Join(base, "operator", "events.log")
	locations := resolveStorageLocations(storageLocationInputs{
		Config:          cfg,
		PreferencesRoot: filepath.Join(base, "preferences"),
		ExecutablePath:  filepath.Join(base, "gosched-gui.exe"),
		GOOS:            "linux",
		Stat:            func(string) (os.FileInfo, error) { return nil, os.ErrNotExist },
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
