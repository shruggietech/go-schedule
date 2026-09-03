package gui

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func TestOptionsAppearanceAppliesPersistsAndResets(t *testing.T) {
	prefs := testApp.Preferences()
	oldMode := prefs.String(appearanceModePreferenceKey)
	oldFont := prefs.String(appearanceFontPreferenceKey)
	t.Cleanup(func() {
		prefs.SetString(appearanceModePreferenceKey, oldMode)
		prefs.SetString(appearanceFontPreferenceKey, oldFont)
		applyBrandTheme(testApp.Settings(), loadAppearancePreferences(prefs))
	})
	resetAppearancePreferences(prefs)

	ui := NewUI(testApp, &fakeBackend{})
	if ui.options == nil {
		t.Fatal("Options view was not constructed")
	}
	ui.options.mode.SetSelected("Light")
	ui.options.font.SetSelected("System")

	want := appearancePreferences{Mode: appearanceLight, Font: fontSystem}
	if got := loadAppearancePreferences(prefs); got != want {
		t.Fatalf("persisted appearance = %+v, want %+v", got, want)
	}
	installed, ok := testApp.Settings().Theme().(*brandTheme)
	if !ok || installed.appearance != want {
		t.Fatalf("installed theme = %#v, want appearance %+v", installed, want)
	}

	ui.options.reset.OnTapped()
	if got := loadAppearancePreferences(prefs); got != defaultAppearancePreferences() {
		t.Fatalf("reset appearance = %+v", got)
	}
	if ui.options.mode.Selected != "Dark" || ui.options.font.Selected != "Brand" {
		t.Fatalf("reset controls = mode %q font %q", ui.options.mode.Selected, ui.options.font.Selected)
	}
}

func TestOptionsEveryAppearanceCombination(t *testing.T) {
	prefs := testApp.Preferences()
	old := loadAppearancePreferences(prefs)
	t.Cleanup(func() {
		saveAppearancePreferences(prefs, old)
		applyBrandTheme(testApp.Settings(), old)
	})
	ui := NewUI(testApp, &fakeBackend{})
	for _, mode := range []appearanceMode{appearanceDark, appearanceLight, appearanceSystem} {
		for _, font := range []fontChoice{fontBrand, fontSystem, fontMonospace} {
			want := appearancePreferences{Mode: mode, Font: font}
			ui.applyAppearance(want)
			installed, ok := testApp.Settings().Theme().(*brandTheme)
			if !ok || installed.appearance != want {
				t.Fatalf("apply %+v installed %#v", want, installed)
			}
		}
	}
}

func TestStorageRowsAreSelectableAndCopyExactPath(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	if len(ui.options.storageRows) == 0 {
		t.Fatal("Options has no storage rows")
	}
	for _, row := range ui.options.storageRows {
		if row.location.Available {
			if !row.path.Selectable {
				t.Errorf("%s path is not selectable", row.location.Category)
			}
			row.copy.OnTapped()
			if got := ui.clipboard.Content(); got != row.location.Path {
				t.Fatalf("copied %q, want exact %q", got, row.location.Path)
			}
			continue
		}
		if !row.copy.Disabled() {
			t.Errorf("%s unavailable row copy is enabled", row.location.Category)
		}
		if row.path.Text == "" {
			t.Errorf("%s unavailable row has no explanation", row.location.Category)
		}
	}
}

func TestOptionsStorageRowsExposeLifecycleAndScope(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	var labels []string
	walkInfoObjects(ui.options.root, func(object fyne.CanvasObject) {
		if label, ok := object.(*widget.Label); ok {
			labels = append(labels, label.Text)
		}
	})
	for _, row := range ui.options.storageRows {
		for _, want := range []string{string(row.location.Scope), string(row.location.Existence), row.location.SoftwareOnlyRemoval, row.location.ExplicitDataWipe} {
			if !containsString(labels, want) {
				t.Errorf("storage row %s missing visible value %q", row.location.Category, want)
			}
		}
	}
}

func TestOptionsRefreshUsesConnectedDaemonRuntimePaths(t *testing.T) {
	base := t.TempDir()
	custom := filepath.Join(base, "custom-daemon")
	backend := &fakeBackend{runtimeInfo: server.RuntimeInfoResponse{
		DataDir:      custom,
		DatabasePath: filepath.Join(custom, "custom.db"),
		ConfigPath:   filepath.Join(base, "daemon.json"),
		LogPath:      filepath.Join(base, "logs", "custom.log"),
		LockPath:     filepath.Join(custom, "custom.lock"),
	}}
	ui := NewUI(testApp, backend)
	ui.refreshStorageLocations(t.Context())

	want := map[string]string{
		"Machine data":  backend.runtimeInfo.DataDir,
		"Task database": backend.runtimeInfo.DatabasePath,
		"Configuration": backend.runtimeInfo.ConfigPath,
		"Logs":          backend.runtimeInfo.LogPath,
		"Runtime state": backend.runtimeInfo.LockPath,
	}
	for _, row := range ui.options.storageRows {
		if path, ok := want[row.location.Category]; ok {
			if row.location.Path != path {
				t.Errorf("%s path = %q, want %q", row.location.Category, row.location.Path, path)
			}
			delete(want, row.location.Category)
		}
	}
	if len(want) != 0 {
		t.Fatalf("missing daemon runtime rows: %v", want)
	}

	// A second refresh must replace, not append to, the daemon inventory.
	backend.runtimeInfo.DataDir = filepath.Join(base, "moved")
	backend.runtimeInfo.DatabasePath = filepath.Join(backend.runtimeInfo.DataDir, "custom.db")
	ui.refreshStorageLocations(t.Context())
	for _, row := range ui.options.storageRows {
		if row.location.Category == "Machine data" && row.location.Path != backend.runtimeInfo.DataDir {
			t.Fatalf("refreshed machine data = %q, want %q", row.location.Path, backend.runtimeInfo.DataDir)
		}
	}

	backend.runtimeErr = errors.New("runtime metadata temporarily unavailable")
	previous := backend.runtimeInfo.DataDir
	backend.runtimeInfo.DataDir = filepath.Join(base, "must-not-replace-known-state")
	ui.refreshStorageLocations(t.Context())
	for _, row := range ui.options.storageRows {
		if row.location.Category == "Machine data" && row.location.Path != previous {
			t.Fatalf("failed metadata refresh replaced known path with %q", row.location.Path)
		}
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if strings.Contains(value, want) {
			return true
		}
	}
	return false
}
