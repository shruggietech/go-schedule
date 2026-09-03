package gui

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func TestOptionsAppearanceAppliesPersistsAndResets(t *testing.T) {
	prefs := testApp.Preferences()
	oldMode := prefs.String(appearanceModePreferenceKey)
	oldFont := prefs.String(appearanceFontPreferenceKey)
	oldScroll := prefs.FloatWithFallback(scrollSensitivityPreferenceKey, defaultScrollSensitivity)
	t.Cleanup(func() {
		prefs.SetString(appearanceModePreferenceKey, oldMode)
		prefs.SetString(appearanceFontPreferenceKey, oldFont)
		prefs.SetFloat(scrollSensitivityPreferenceKey, oldScroll)
		applyBrandTheme(testApp.Settings(), loadAppearancePreferences(prefs))
	})
	resetAppearancePreferences(prefs)

	ui := NewUI(testApp, &fakeBackend{})
	if ui.options == nil {
		t.Fatal("Options view was not constructed")
	}
	ui.options.mode.SetSelected("Light")
	ui.options.font.SetSelected("System")

	ui.options.sensitivity.SetValue(3.5)
	want := appearancePreferences{Mode: appearanceLight, Font: fontSystem, ScrollSensitivity: 3.5}
	if got := loadAppearancePreferences(prefs); got != want {
		t.Fatalf("persisted appearance = %+v, want %+v", got, want)
	}
	installed, ok := testApp.Settings().Theme().(*brandTheme)
	if !ok || !installed.appearance.sameTheme(want) {
		t.Fatalf("installed theme = %#v, want theme choices %+v", installed, want)
	}

	ui.options.reset.OnTapped()
	if got := loadAppearancePreferences(prefs); got != defaultAppearancePreferences() {
		t.Fatalf("reset appearance = %+v", got)
	}
	if ui.options.mode.Selected != "Dark" || ui.options.font.Selected != "System" || ui.options.sensitivity.Value != defaultScrollSensitivity {
		t.Fatalf("reset controls = mode %q font %q scroll %v", ui.options.mode.Selected, ui.options.font.Selected, ui.options.sensitivity.Value)
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
		for _, font := range []fontChoice{fontSystem, fontBrand, fontInter, fontUbuntu, fontMonospace} {
			want := appearancePreferences{Mode: mode, Font: font, ScrollSensitivity: defaultScrollSensitivity}
			ui.applyAppearance(want)
			installed, ok := testApp.Settings().Theme().(*brandTheme)
			if !ok || installed.appearance != want {
				t.Fatalf("apply %+v installed %#v", want, installed)
			}
		}
	}
}

func TestOptionsSelectorsOfferOnlyAlternatives(t *testing.T) {
	prefs := testApp.Preferences()
	old := loadAppearancePreferences(prefs)
	t.Cleanup(func() {
		saveAppearancePreferences(prefs, old)
		applyBrandTheme(testApp.Settings(), old)
	})
	saveAppearancePreferences(prefs, defaultAppearancePreferences())

	ui := NewUI(testApp, &fakeBackend{})
	if got, want := ui.options.mode.Options, []string{"Light", "Follow system"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("mode alternatives = %v, want %v", got, want)
	}
	if got, want := ui.options.font.Options, []string{"Geist (brand)", "Inter", "Ubuntu", "Monospace"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("font alternatives = %v, want %v", got, want)
	}

	ui.options.mode.SetSelected("Light")
	if ui.options.mode.Selected != "Light" {
		t.Fatalf("selected mode = %q, want Light", ui.options.mode.Selected)
	}
	if got, want := ui.options.mode.Options, []string{"Dark", "Follow system"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("updated mode alternatives = %v, want %v", got, want)
	}
	ui.options.font.SetSelected("Ubuntu")
	if ui.options.font.Selected != "Ubuntu" {
		t.Fatalf("selected font = %q, want Ubuntu", ui.options.font.Selected)
	}
	if containsString(ui.options.font.Options, "Ubuntu") {
		t.Fatalf("font alternatives still contain current value: %v", ui.options.font.Options)
	}
}

func TestOptionsFontChangeReflowsStorageColumns(t *testing.T) {
	prefs := testApp.Preferences()
	old := loadAppearancePreferences(prefs)
	t.Cleanup(func() {
		saveAppearancePreferences(prefs, old)
		applyBrandTheme(testApp.Settings(), old)
	})
	saveAppearancePreferences(prefs, defaultAppearancePreferences())

	ui := NewUI(testApp, &fakeBackend{})
	previousHeader := ui.options.storageHeader
	ui.options.font.SetSelected("Monospace")

	if ui.options.storageHeader == previousHeader {
		t.Fatal("font change did not rebuild the measured storage header")
	}
	if got, want := ui.options.storageCategoryWidth, storageCategoryWidth(ui.storageLocations); got != want {
		t.Fatalf("category width after font change = %v, want freshly measured %v", got, want)
	}
}

func TestOptionsScrollSensitivityControl(t *testing.T) {
	prefs := testApp.Preferences()
	old := loadAppearancePreferences(prefs)
	t.Cleanup(func() {
		saveAppearancePreferences(prefs, old)
		applyBrandTheme(testApp.Settings(), old)
	})
	saveAppearancePreferences(prefs, defaultAppearancePreferences())

	ui := NewUI(testApp, &fakeBackend{})
	slider := ui.options.sensitivity
	if slider.Min != minimumScrollSensitivity || slider.Max != maximumScrollSensitivity || slider.Step != scrollSensitivityStep {
		t.Fatalf("slider bounds = %v..%v step %v", slider.Min, slider.Max, slider.Step)
	}
	if slider.Value != defaultScrollSensitivity || ui.options.sensitivityValue.Text != "2x (default)" {
		t.Fatalf("default slider = %v label=%q", slider.Value, ui.options.sensitivityValue.Text)
	}
	slider.SetValue(3.5)
	if got := loadAppearancePreferences(prefs).ScrollSensitivity; got != 3.5 {
		t.Fatalf("persisted sensitivity = %v, want 3.5", got)
	}
	if ui.options.sensitivityValue.Text != "3.5x" {
		t.Fatalf("sensitivity label = %q, want 3.5x", ui.options.sensitivityValue.Text)
	}
}

func TestOptionsStorageRowsUseCompactAlignedPresentation(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	if ui.options.storageHeader == nil {
		t.Fatal("storage table has no header")
	}
	if len(ui.options.storageRows) == 0 {
		t.Fatal("storage table has no rows")
	}
	for _, row := range ui.options.storageRows {
		if _, isCard := row.root.(*widget.Card); isCard {
			t.Fatalf("storage row %q still uses card presentation", row.location.Category)
		}
		if row.category == nil || row.details == nil || row.copy == nil {
			t.Fatalf("storage row %q is missing aligned columns", row.location.Category)
		}
		if !row.location.Available && !row.muted {
			t.Fatalf("unavailable row %q is not wholly muted", row.location.Category)
		}
	}
}

func TestStorageRowWrapsWithinSmallViewportWithoutHorizontalScroll(t *testing.T) {
	location := storageLocation{
		Category:            "Machine data",
		Path:                `C:\very-long-storage-root\one-unbroken-segment-that-must-wrap-without-growing-the-options-view\another-unbroken-segment\go-schedule`,
		Scope:               storageScopeMachine,
		Existence:           storagePresent,
		SoftwareOnlyRemoval: "Preserved",
		ExplicitDataWipe:    "Removed only after explicit confirmation",
		Available:           true,
	}
	row := newStorageRowView(location, testApp.Clipboard())
	row.root.Resize(fyne.NewSize(600, 140))
	row.root.Refresh()

	if got := row.root.MinSize().Width; got > 600 {
		t.Fatalf("row minimum width = %v, want <= 600 without horizontal scrolling", got)
	}
	if right := row.copy.Position().X + row.copy.Size().Width; right > 600 || right < 600-theme.Padding()*2 {
		t.Fatalf("Copy right edge = %v, want within far-right inset of 600", right)
	}
	if row.path.Size().Width <= 0 || row.path.Size().Width >= 600 {
		t.Fatalf("path width = %v, want flexible middle column", row.path.Size().Width)
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
		switch label := object.(type) {
		case *widget.Label:
			labels = append(labels, label.Text)
		case *widget.RichText:
			labels = append(labels, label.String())
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
