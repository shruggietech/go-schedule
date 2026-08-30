package gui

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// nrgba coerces a themed color back to NRGBA for comparison. Fyne returns the
// exact color.Color we hand it, so this is a straight type assertion in practice
// but stays robust if Fyne ever wraps the value.
func nrgba(t *testing.T, c color.Color) color.NRGBA {
	t.Helper()
	if n, ok := c.(color.NRGBA); ok {
		return n
	}
	r, g, b, a := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

func TestBrandThemeColors(t *testing.T) {
	th := newBrandTheme()
	cases := []struct {
		name fyne.ThemeColorName
		want color.NRGBA
	}{
		{theme.ColorNameBackground, cNight},
		{theme.ColorNameForeground, cText},
		{theme.ColorNamePrimary, cAnchor},
		{theme.ColorNameFocus, cAnchor},
		{theme.ColorNameHyperlink, cAnchor},
		{theme.ColorNameSuccess, cInterval},
		{theme.ColorNameWarning, cHold},
		{theme.ColorNameError, cStop},
		{theme.ColorNameInputBackground, cPanel},
		{theme.ColorNameInputBorder, cLine},
		{theme.ColorNamePlaceHolder, cMuted},
		// Contrast rule: ink on light accents is Night; text on red stays light.
		{theme.ColorNameForegroundOnPrimary, cNight},
		{theme.ColorNameForegroundOnSuccess, cNight},
		{theme.ColorNameForegroundOnWarning, cNight},
		{theme.ColorNameForegroundOnError, cText},
	}
	for _, tc := range cases {
		got := nrgba(t, th.Color(tc.name, theme.VariantDark))
		if got != tc.want {
			t.Errorf("Color(%s) = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// The theme is dark-first: a Light variant request still returns the dark color.
func TestBrandThemeIgnoresVariant(t *testing.T) {
	th := newBrandTheme()
	if got := nrgba(t, th.Color(theme.ColorNameBackground, theme.VariantLight)); got != cNight {
		t.Errorf("Background under VariantLight = %v, want Night %v", got, cNight)
	}
}

func TestBrandThemeFonts(t *testing.T) {
	th := newBrandTheme()
	cases := []struct {
		style fyne.TextStyle
		want  string
	}{
		{fyne.TextStyle{Monospace: true}, "GeistMono-Regular.ttf"},
		{fyne.TextStyle{Bold: true}, "SpaceGrotesk-Bold.ttf"},
		{fyne.TextStyle{}, "Geist-Regular.ttf"},
	}
	for _, tc := range cases {
		if got := th.Font(tc.style).Name(); got != tc.want {
			t.Errorf("Font(%+v) = %q, want %q", tc.style, got, tc.want)
		}
	}
}

func TestBrandThemeRadii(t *testing.T) {
	th := newBrandTheme()
	if got := th.Size(theme.SizeNameInputRadius); got != radiusMd {
		t.Errorf("InputRadius = %v, want %v", got, radiusMd)
	}
	if got := th.Size(theme.SizeNameSelectionRadius); got != radiusSm {
		t.Errorf("SelectionRadius = %v, want %v", got, radiusSm)
	}
	// A delegated size stays non-zero (sanity that the embedded default answers).
	if got := th.Size(theme.SizeNameText); got <= 0 {
		t.Errorf("delegated text size = %v, want > 0", got)
	}
}

func TestApplyBrandThemeIsIdempotent(t *testing.T) {
	applyBrandTheme(testApp.Settings())
	installed := testApp.Settings().Theme()

	applyBrandTheme(testApp.Settings())

	if got := testApp.Settings().Theme(); got != installed {
		t.Fatal("applyBrandTheme replaced an already-installed brand theme")
	}
}
