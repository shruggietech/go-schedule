package gui

import (
	"math"
	"reflect"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func TestAppearancePreferencesDefaultsAndValidation(t *testing.T) {
	prefs := testApp.Preferences()
	oldMode := prefs.String(appearanceModePreferenceKey)
	oldFont := prefs.String(appearanceFontPreferenceKey)
	oldScroll := prefs.FloatWithFallback(scrollSensitivityPreferenceKey, defaultScrollSensitivity)
	t.Cleanup(func() {
		prefs.SetString(appearanceModePreferenceKey, oldMode)
		prefs.SetString(appearanceFontPreferenceKey, oldFont)
		prefs.SetFloat(scrollSensitivityPreferenceKey, oldScroll)
	})

	for _, tc := range []struct {
		name       string
		mode       string
		font       string
		scroll     float64
		wantMode   appearanceMode
		wantFont   fontChoice
		wantScroll float64
	}{
		{name: "missing", wantMode: appearanceDark, wantFont: fontSystem, wantScroll: defaultScrollSensitivity},
		{name: "valid", mode: "light", font: "ubuntu", scroll: 3.5, wantMode: appearanceLight, wantFont: fontUbuntu, wantScroll: 3.5},
		{name: "legacy explicit brand", mode: "dark", font: "brand", scroll: 2, wantMode: appearanceDark, wantFont: fontBrand, wantScroll: 2},
		{name: "future values", mode: "sepia", font: "downloaded", scroll: 99, wantMode: appearanceDark, wantFont: fontSystem, wantScroll: defaultScrollSensitivity},
		{name: "independent fallback", mode: "system", font: "broken", scroll: 0.5, wantMode: appearanceSystem, wantFont: fontSystem, wantScroll: defaultScrollSensitivity},
	} {
		t.Run(tc.name, func(t *testing.T) {
			prefs.SetString(appearanceModePreferenceKey, tc.mode)
			prefs.SetString(appearanceFontPreferenceKey, tc.font)
			if tc.scroll == 0 {
				prefs.RemoveValue(scrollSensitivityPreferenceKey)
			} else {
				prefs.SetFloat(scrollSensitivityPreferenceKey, tc.scroll)
			}
			got := loadAppearancePreferences(prefs)
			if got.Mode != tc.wantMode || got.Font != tc.wantFont || got.ScrollSensitivity != tc.wantScroll {
				t.Fatalf("loadAppearancePreferences() = %+v, want mode=%q font=%q scroll=%v", got, tc.wantMode, tc.wantFont, tc.wantScroll)
			}
		})
	}
}

func TestAppearancePreferencesPersistAndReset(t *testing.T) {
	prefs := testApp.Preferences()
	oldMode := prefs.String(appearanceModePreferenceKey)
	oldFont := prefs.String(appearanceFontPreferenceKey)
	oldScroll := prefs.FloatWithFallback(scrollSensitivityPreferenceKey, defaultScrollSensitivity)
	t.Cleanup(func() {
		prefs.SetString(appearanceModePreferenceKey, oldMode)
		prefs.SetString(appearanceFontPreferenceKey, oldFont)
		prefs.SetFloat(scrollSensitivityPreferenceKey, oldScroll)
	})

	want := appearancePreferences{Mode: appearanceSystem, Font: fontMonospace, ScrollSensitivity: 3.5}
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
	wantFonts := []string{"System", "Geist (brand)", "Inter", "Ubuntu", "Monospace"}
	if got := fontChoiceLabels(); !reflect.DeepEqual(got, wantFonts) {
		t.Fatalf("font choices = %v, want %v", got, wantFonts)
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

func TestScrollSensitivityNormalization(t *testing.T) {
	for _, tc := range []struct {
		value float64
		want  float64
	}{
		{value: -1, want: defaultScrollSensitivity},
		{value: 0, want: defaultScrollSensitivity},
		{value: 1, want: 1},
		{value: 1.24, want: 1},
		{value: 1.26, want: 1.5},
		{value: 3.74, want: 3.5},
		{value: 3.76, want: 4},
		{value: 5, want: defaultScrollSensitivity},
		{value: math.NaN(), want: defaultScrollSensitivity},
	} {
		if got := normalizeScrollSensitivity(tc.value); got != tc.want {
			t.Errorf("normalizeScrollSensitivity(%v) = %v, want %v", tc.value, got, tc.want)
		}
	}
}
