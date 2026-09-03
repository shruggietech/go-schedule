package gui

import (
	"image/color"
	"math"
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

func TestBrandThemeTextContrastAndFocusVisibility(t *testing.T) {
	for _, mode := range []appearanceMode{appearanceDark, appearanceLight, appearanceSystem} {
		for _, variant := range []fyne.ThemeVariant{theme.VariantDark, theme.VariantLight} {
			th := newBrandThemeFor(appearancePreferences{Mode: mode, Font: fontBrand})
			background := th.Color(theme.ColorNameBackground, variant)
			for _, name := range []fyne.ThemeColorName{theme.ColorNameForeground, theme.ColorNameHyperlink} {
				if ratio := contrastRatio(th.Color(name, variant), background); ratio < 4.5 {
					t.Errorf("mode=%s variant=%v %s contrast = %.2f, want >= 4.5", mode, variant, name, ratio)
				}
			}
			for _, surface := range []fyne.ThemeColorName{theme.ColorNameBackground, theme.ColorNameInputBackground} {
				if ratio := contrastRatio(th.Color(theme.ColorNameError, variant), th.Color(surface, variant)); ratio < 3 {
					t.Errorf("mode=%s variant=%v standalone error on %s contrast = %.2f, want >= 3", mode, variant, surface, ratio)
				}
			}
			focus := blendThemeColor(background, th.Color(theme.ColorNameFocus, variant))
			if ratio := contrastRatio(focus, background); ratio < 3 {
				t.Errorf("mode=%s variant=%v focus contrast = %.2f, want >= 3", mode, variant, ratio)
			}
		}
	}
}

func TestStructuredRowSurfacesPreserveTextAndStateContrast(t *testing.T) {
	for _, mode := range []appearanceMode{appearanceDark, appearanceLight, appearanceSystem} {
		for _, variant := range []fyne.ThemeVariant{theme.VariantDark, theme.VariantLight} {
			th := newBrandThemeFor(appearancePreferences{Mode: mode, Font: fontSystem})
			background := th.Color(theme.ColorNameBackground, variant)
			alternate := blendThemeColor(background, th.Color(colorNameAlternateRow, variant))
			for name, surface := range map[string]color.Color{"base": background, "alternate": alternate} {
				if ratio := contrastRatio(th.Color(theme.ColorNameForeground, variant), surface); ratio < 4.5 {
					t.Errorf("mode=%s variant=%v %s text contrast=%.2f, want >=4.5", mode, variant, name, ratio)
				}
				for stateName, state := range map[string]fyne.ThemeColorName{
					"hover": theme.ColorNameHover, "focus": theme.ColorNameFocus, "selection": theme.ColorNameSelection,
				} {
					stateSurface := blendThemeColor(surface, th.Color(state, variant))
					if ratio := contrastRatio(th.Color(theme.ColorNameForeground, variant), stateSurface); ratio < 4.5 {
						t.Errorf("mode=%s variant=%v %s+%s text contrast=%.2f, want >=4.5", mode, variant, name, stateName, ratio)
					}
					if ratio := contrastRatio(stateSurface, surface); ratio < 1.1 {
						t.Errorf("mode=%s variant=%v %s+%s surface contrast=%.2f, want >=1.1", mode, variant, name, stateName, ratio)
					}
				}
			}
			for _, semantic := range []fyne.ThemeColorName{theme.ColorNamePrimary, theme.ColorNameSuccess, theme.ColorNameWarning, theme.ColorNameError} {
				if ratio := contrastRatio(th.Color(semantic, variant), background); ratio < 3 {
					t.Errorf("mode=%s variant=%v semantic=%s contrast=%.2f, want >=3", mode, variant, semantic, ratio)
				}
			}
		}
	}
}

func contrastRatio(a, b color.Color) float64 {
	la, lb := relativeLuminance(a), relativeLuminance(b)
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

func relativeLuminance(value color.Color) float64 {
	r, g, b, _ := value.RGBA()
	linear := func(component uint32) float64 {
		c := float64(component) / 65535
		if c <= 0.04045 {
			return c / 12.92
		}
		return math.Pow((c+0.055)/1.055, 2.4)
	}
	return 0.2126*linear(r) + 0.7152*linear(g) + 0.0722*linear(b)
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
		{theme.ColorNameFocus, withAlpha(cAnchor, 0x90)},
		{theme.ColorNameHyperlink, cAnchor},
		{theme.ColorNameSuccess, cInterval},
		{theme.ColorNameWarning, cHold},
		{theme.ColorNameError, cStop},
		{theme.ColorNameInputBackground, cPanel},
		{theme.ColorNameInputBorder, cLine},
		{theme.ColorNamePlaceHolder, cMuted},
		// Contrast rule: ink on the light accents is Night.
		{theme.ColorNameForegroundOnPrimary, cNight},
		{theme.ColorNameForegroundOnSuccess, cNight},
		{theme.ColorNameForegroundOnWarning, cNight},
		{theme.ColorNameForegroundOnError, cNight},
	}
	for _, tc := range cases {
		got := nrgba(t, th.Color(tc.name, theme.VariantDark))
		if got != tc.want {
			t.Errorf("Color(%s) = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestBrandThemeAppearanceVariants(t *testing.T) {
	for _, tc := range []struct {
		name     string
		mode     appearanceMode
		variant  fyne.ThemeVariant
		wantBack color.NRGBA
		wantText color.NRGBA
	}{
		{name: "dark ignores light request", mode: appearanceDark, variant: theme.VariantLight, wantBack: cNight, wantText: cText},
		{name: "light ignores dark request", mode: appearanceLight, variant: theme.VariantDark, wantBack: cPaper, wantText: cInk},
		{name: "system dark", mode: appearanceSystem, variant: theme.VariantDark, wantBack: cNight, wantText: cText},
		{name: "system light", mode: appearanceSystem, variant: theme.VariantLight, wantBack: cPaper, wantText: cInk},
	} {
		t.Run(tc.name, func(t *testing.T) {
			th := newBrandThemeFor(appearancePreferences{Mode: tc.mode, Font: fontBrand})
			if got := nrgba(t, th.Color(theme.ColorNameBackground, tc.variant)); got != tc.wantBack {
				t.Errorf("background = %v, want %v", got, tc.wantBack)
			}
			if got := nrgba(t, th.Color(theme.ColorNameForeground, tc.variant)); got != tc.wantText {
				t.Errorf("foreground = %v, want %v", got, tc.wantText)
			}
		})
	}
}

func TestBrandThemeFonts(t *testing.T) {
	th := newBrandThemeFor(appearancePreferences{Mode: appearanceDark, Font: fontBrand})
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

func TestBrandThemeFontChoices(t *testing.T) {
	regular := fyne.TextStyle{}
	bold := fyne.TextStyle{Bold: true}
	symbol := fyne.TextStyle{Symbol: true}
	defaults := theme.DefaultTheme()

	systemTheme := newBrandThemeFor(appearancePreferences{Mode: appearanceDark, Font: fontSystem})
	if got, want := systemTheme.Font(regular).Name(), defaults.Font(regular).Name(); got != want {
		t.Fatalf("system regular font = %q, want %q", got, want)
	}
	if got, want := systemTheme.Font(bold).Name(), defaults.Font(bold).Name(); got != want {
		t.Fatalf("system bold font = %q, want %q", got, want)
	}

	monoTheme := newBrandThemeFor(appearancePreferences{Mode: appearanceLight, Font: fontMonospace})
	if got := monoTheme.Font(regular).Name(); got != "GeistMono-Regular.ttf" {
		t.Fatalf("monospace regular font = %q", got)
	}
	if got := monoTheme.Font(bold).Name(); got != "GeistMono-Regular.ttf" {
		t.Fatalf("monospace bold font = %q", got)
	}
	if got, want := monoTheme.Font(symbol).Name(), defaults.Font(symbol).Name(); got != want {
		t.Fatalf("symbol font = %q, want delegated %q", got, want)
	}
}

func TestBrandThemeCuratedFontResources(t *testing.T) {
	regular := fyne.TextStyle{}
	bold := fyne.TextStyle{Bold: true}
	monospace := fyne.TextStyle{Monospace: true}
	symbol := fyne.TextStyle{Symbol: true}
	defaults := theme.DefaultTheme()

	for _, tc := range []struct {
		choice      fontChoice
		wantRegular string
		wantBold    string
		wantMono    string
	}{
		{choice: fontBrand, wantRegular: "Geist-Regular.ttf", wantBold: "SpaceGrotesk-Bold.ttf", wantMono: "GeistMono-Regular.ttf"},
		{choice: fontInter, wantRegular: "Inter-Regular.ttf", wantBold: "Inter-Bold.ttf", wantMono: "GeistMono-Regular.ttf"},
		{choice: fontUbuntu, wantRegular: "UbuntuSans-Regular.ttf", wantBold: "UbuntuSans-Bold.ttf", wantMono: "GeistMono-Regular.ttf"},
		{choice: fontMonospace, wantRegular: "GeistMono-Regular.ttf", wantBold: "GeistMono-Regular.ttf", wantMono: "GeistMono-Regular.ttf"},
	} {
		t.Run(string(tc.choice), func(t *testing.T) {
			th := newBrandThemeFor(appearancePreferences{Mode: appearanceDark, Font: tc.choice})
			if got := th.Font(regular).Name(); got != tc.wantRegular {
				t.Errorf("regular = %q, want %q", got, tc.wantRegular)
			}
			if got := th.Font(bold).Name(); got != tc.wantBold {
				t.Errorf("bold = %q, want %q", got, tc.wantBold)
			}
			if got := th.Font(monospace).Name(); got != tc.wantMono {
				t.Errorf("monospace = %q, want %q", got, tc.wantMono)
			}
			if got, want := th.Font(symbol).Name(), defaults.Font(symbol).Name(); got != want {
				t.Errorf("symbol = %q, want %q", got, want)
			}
		})
	}
}

func TestBrandThemeButtonStateCompositesRemainReadable(t *testing.T) {
	for _, mode := range []appearanceMode{appearanceDark, appearanceLight} {
		th := newBrandThemeFor(appearancePreferences{Mode: mode, Font: fontSystem})
		variant, _ := mode.themeVariant(theme.VariantDark)
		states := []fyne.ThemeColorName{theme.ColorNameHover, theme.ColorNamePressed, theme.ColorNameFocus}
		for _, state := range states {
			if overlay := nrgba(t, th.Color(state, variant)); overlay.A == 0xff {
				t.Errorf("mode=%s state=%s overlay is opaque and would erase the importance background", mode, state)
			}
		}

		pairs := []struct {
			name string
			fg   fyne.ThemeColorName
			bg   fyne.ThemeColorName
		}{
			{name: "ordinary", fg: theme.ColorNameForeground, bg: theme.ColorNameButton},
			{name: "primary", fg: theme.ColorNameForegroundOnPrimary, bg: theme.ColorNamePrimary},
			{name: "danger", fg: theme.ColorNameForegroundOnError, bg: theme.ColorNameError},
			{name: "success", fg: theme.ColorNameForegroundOnSuccess, bg: theme.ColorNameSuccess},
			{name: "warning", fg: theme.ColorNameForegroundOnWarning, bg: theme.ColorNameWarning},
		}
		for _, pair := range pairs {
			if ratio := contrastRatio(th.Color(pair.fg, variant), th.Color(pair.bg, variant)); ratio < 4.5 {
				t.Errorf("mode=%s %s rest text contrast = %.2f, want >= 4.5", mode, pair.name, ratio)
			}
			for _, state := range states {
				background := blendThemeColor(th.Color(pair.bg, variant), th.Color(state, variant))
				if ratio := contrastRatio(th.Color(pair.fg, variant), background); ratio < 4.5 {
					t.Errorf("mode=%s %s+%s text contrast = %.2f, want >= 4.5", mode, pair.name, state, ratio)
				}
			}
		}
		if ratio := contrastRatio(th.Color(theme.ColorNameDisabled, variant), th.Color(theme.ColorNameDisabledButton, variant)); ratio < 4.5 {
			t.Errorf("mode=%s disabled text contrast = %.2f, want >= 4.5", mode, ratio)
		}
		focusOrdinary := blendThemeColor(th.Color(theme.ColorNameButton, variant), th.Color(theme.ColorNameFocus, variant))
		if ratio := contrastRatio(focusOrdinary, th.Color(theme.ColorNameBackground, variant)); ratio < 3 {
			t.Errorf("mode=%s ordinary focus surface contrast = %.2f, want >= 3", mode, ratio)
		}
	}
}

func blendThemeColor(under, over color.Color) color.Color {
	dstR, dstG, dstB, dstA := under.RGBA()
	srcR, srcG, srcB, srcA := over.RGBA()
	blend := func(src, dst, alpha uint32) uint16 {
		return uint16((src + dst - (dst * alpha / 0xffff)) & 0xffff)
	}
	return color.RGBA64{
		R: blend(srcR, dstR, srcA),
		G: blend(srcG, dstG, srcA),
		B: blend(srcB, dstB, srcA),
		A: blend(srcA, dstA, srcA),
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
