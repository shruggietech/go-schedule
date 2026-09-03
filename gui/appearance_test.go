package gui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func TestAppearancePreferencesDefaultsAndValidation(t *testing.T) {
	prefs := testApp.Preferences()
	oldMode := prefs.String(appearanceModePreferenceKey)
	oldFont := prefs.String(appearanceFontPreferenceKey)
	t.Cleanup(func() {
		prefs.SetString(appearanceModePreferenceKey, oldMode)
		prefs.SetString(appearanceFontPreferenceKey, oldFont)
	})

	for _, tc := range []struct {
		name     string
		mode     string
		font     string
		wantMode appearanceMode
		wantFont fontChoice
	}{
		{name: "missing", wantMode: appearanceDark, wantFont: fontBrand},
		{name: "valid", mode: "light", font: "system", wantMode: appearanceLight, wantFont: fontSystem},
		{name: "future values", mode: "sepia", font: "downloaded", wantMode: appearanceDark, wantFont: fontBrand},
		{name: "independent fallback", mode: "system", font: "broken", wantMode: appearanceSystem, wantFont: fontBrand},
	} {
		t.Run(tc.name, func(t *testing.T) {
			prefs.SetString(appearanceModePreferenceKey, tc.mode)
			prefs.SetString(appearanceFontPreferenceKey, tc.font)
			got := loadAppearancePreferences(prefs)
			if got.Mode != tc.wantMode || got.Font != tc.wantFont {
				t.Fatalf("loadAppearancePreferences() = %+v, want mode=%q font=%q", got, tc.wantMode, tc.wantFont)
			}
		})
	}
}

func TestAppearancePreferencesPersistAndReset(t *testing.T) {
	prefs := testApp.Preferences()
	oldMode := prefs.String(appearanceModePreferenceKey)
	oldFont := prefs.String(appearanceFontPreferenceKey)
	t.Cleanup(func() {
		prefs.SetString(appearanceModePreferenceKey, oldMode)
		prefs.SetString(appearanceFontPreferenceKey, oldFont)
	})

	want := appearancePreferences{Mode: appearanceSystem, Font: fontMonospace}
	saveAppearancePreferences(prefs, want)
	if got := loadAppearancePreferences(prefs); got != want {
		t.Fatalf("restored preferences = %+v, want %+v", got, want)
	}

	resetAppearancePreferences(prefs)
	if got := loadAppearancePreferences(prefs); got != defaultAppearancePreferences() {
		t.Fatalf("reset preferences = %+v, want %+v", got, defaultAppearancePreferences())
	}
}

func TestAppearanceChoicesAreBounded(t *testing.T) {
	if got := appearanceModeLabels(); len(got) != 3 {
		t.Fatalf("appearance modes = %v, want exactly 3", got)
	}
	if got := fontChoiceLabels(); len(got) != 3 {
		t.Fatalf("font choices = %v, want exactly 3", got)
	}
	for _, mode := range []appearanceMode{appearanceDark, appearanceLight, appearanceSystem} {
		variant, ok := mode.themeVariant(fyne.ThemeVariant(99))
		if !ok {
			t.Fatalf("mode %q should be valid", mode)
		}
		if mode == appearanceSystem && variant != theme.VariantDark {
			t.Fatalf("unknown system variant = %v, want safe dark fallback", variant)
		}
	}
}
