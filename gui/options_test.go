package gui

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
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

func containsString(values []string, want string) bool {
	for _, value := range values {
		if strings.Contains(value, want) {
			return true
		}
	}
	return false
}
