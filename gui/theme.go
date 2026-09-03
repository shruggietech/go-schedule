package gui

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// This file applies the go-schedule brand design system to the Fyne GUI. The
// brand is dark-first and close to monochrome: neutral Night/Panel/Line rails
// carry structure, Anchor Blue marks selection/focus/links and the next run,
// Interval Mint marks success/recurrence, and Hold Amber and Stop Red stay
// scarce so they keep their signal. See the go-schedule brand kit for the full
// system.
//
// brandTheme embeds theme.DefaultTheme() so Icon() and any size we do not
// override delegate for free. Everything here is pure Go (no cgo), so the
// gui package keeps its cgo-free, headless-testable guarantee.

// Brand palette (dark surfaces). Opaque unless an alpha is set explicitly.
var (
	cNight  = color.NRGBA{R: 0x07, G: 0x10, B: 0x14, A: 0xFF} // primary dark background
	cPanel  = color.NRGBA{R: 0x0D, G: 0x17, B: 0x1C, A: 0xFF} // cards, inputs, code, menus
	cRaised = color.NRGBA{R: 0x13, G: 0x23, B: 0x2A, A: 0xFF} // elevated / selected / hover
	cLine   = color.NRGBA{R: 0x28, G: 0x41, B: 0x4B, A: 0xFF} // borders, separators, rails
	cText   = color.NRGBA{R: 0xF3, G: 0xF7, B: 0xF8, A: 0xFF} // primary text
	cMuted  = color.NRGBA{R: 0x9B, G: 0xAE, B: 0xB6, A: 0xFF} // secondary text / placeholder

	cInterval = color.NRGBA{R: 0x62, G: 0xD9, B: 0xB7, A: 0xFF} // Interval Mint — success/recurrence
	cAnchor   = color.NRGBA{R: 0x58, G: 0xA6, B: 0xFF, A: 0xFF} // Anchor Blue — primary/focus/links
	cHold     = color.NRGBA{R: 0xF2, G: 0xB8, B: 0x4B, A: 0xFF} // Hold Amber — warning (rare)
	cStop     = color.NRGBA{R: 0xE0, G: 0x5F, B: 0x5F, A: 0xFF} // Stop Red — error only

	// Light surfaces retain the same semantic accents while supplying sufficient
	// contrast for body text, selection, inputs, and focus.
	cPaper       = color.NRGBA{R: 0xF7, G: 0xFA, B: 0xFB, A: 0xFF}
	cLightPanel  = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	cLightRaised = color.NRGBA{R: 0xE7, G: 0xF0, B: 0xF4, A: 0xFF}
	cLightLine   = color.NRGBA{R: 0xB9, G: 0xC9, B: 0xD0, A: 0xFF}
	cInk         = color.NRGBA{R: 0x0D, G: 0x17, B: 0x1C, A: 0xFF}
	cLightMuted  = color.NRGBA{R: 0x4E, G: 0x63, B: 0x6C, A: 0xFF}
	cLightAnchor = color.NRGBA{R: 0x0B, G: 0x62, B: 0xB5, A: 0xFF}
)

// Brand geometry: small radii, calm technical surfaces.
const (
	radiusMd = 6 // buttons/inputs
	radiusSm = 4 // chips, small controls, selection
)

// Brand fonts. Three faces cover the interface: Geist for body/UI, Geist Mono
// for commands/schedules/labels (the highest-value face — the brand sets all
// cron fields, timestamps, and daemon output in mono), and Space Grotesk Bold
// for emphasis/headings. fyne.TextStyle carries no "medium" weight and the UI
// uses no italics, so those variants are intentionally not embedded.
//
//go:embed assets/fonts/Geist-Regular.ttf
var fontBody []byte

//go:embed assets/fonts/GeistMono-Regular.ttf
var fontMono []byte

//go:embed assets/fonts/SpaceGrotesk-Bold.ttf
var fontDisplayBold []byte

var (
	resBody        = fyne.NewStaticResource("Geist-Regular.ttf", fontBody)
	resMono        = fyne.NewStaticResource("GeistMono-Regular.ttf", fontMono)
	resDisplayBold = fyne.NewStaticResource("SpaceGrotesk-Bold.ttf", fontDisplayBold)
)

// brandTheme is the go-schedule Fyne theme configured from bounded appearance
// preferences. It remains immutable after construction so a settings refresh
// cannot observe a partially changed palette or font.
type brandTheme struct {
	fyne.Theme
	appearance appearancePreferences
}

// newBrandTheme returns the go-schedule brand theme.
func newBrandTheme() fyne.Theme { return newBrandThemeFor(defaultAppearancePreferences()) }

func newBrandThemeFor(appearance appearancePreferences) *brandTheme {
	return &brandTheme{Theme: theme.DefaultTheme(), appearance: appearance.normalized()}
}

// applyBrandTheme installs the go-schedule theme unless this app already uses
// it. Reapplying a Fyne theme refreshes every window and clears shared font and
// theme caches, which is unnecessary for repeated UI construction and can
// invalidate text rendering that is already in progress.
func applyBrandTheme(settings fyne.Settings, choices ...appearancePreferences) {
	appearance := defaultAppearancePreferences()
	if len(choices) > 0 {
		appearance = choices[0].normalized()
	}
	if installed, ok := settings.Theme().(*brandTheme); ok && installed.appearance == appearance {
		return
	}
	settings.SetTheme(newBrandThemeFor(appearance))
}

// withAlpha returns c with its alpha replaced.
func withAlpha(c color.NRGBA, a uint8) color.NRGBA { c.A = a; return c }

// Color maps Fyne's semantic color roles onto the selected palette. Explicit
// modes force their variant; Follow system honors the driver's current value.
func (t *brandTheme) Color(n fyne.ThemeColorName, requested fyne.ThemeVariant) color.Color {
	variant, _ := t.appearance.Mode.themeVariant(requested)
	if variant == theme.VariantLight {
		return t.lightColor(n, variant)
	}
	switch n {
	case theme.ColorNameBackground:
		return cNight
	case theme.ColorNameForeground:
		return cText
	case theme.ColorNamePrimary, theme.ColorNameFocus, theme.ColorNameHyperlink:
		return cAnchor
	case theme.ColorNameButton, theme.ColorNameInputBackground,
		theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground,
		theme.ColorNameHeaderBackground:
		return cPanel
	case theme.ColorNameSelection:
		return withAlpha(cAnchor, 0x40) // ~25% Anchor fill for selected rows/text
	case theme.ColorNameHover, theme.ColorNamePressed, theme.ColorNameDisabledButton:
		return cRaised
	case theme.ColorNameInputBorder, theme.ColorNameSeparator:
		return cLine
	case theme.ColorNameScrollBar:
		return withAlpha(cLine, 0xB0)
	case theme.ColorNameScrollBarBackground:
		return cNight
	case theme.ColorNamePlaceHolder, theme.ColorNameDisabled:
		return cMuted
	case theme.ColorNameShadow:
		return color.NRGBA{A: 0x66} // quiet black shadow, brand avoids heavy elevation
	case theme.ColorNameSuccess:
		return cInterval
	case theme.ColorNameWarning:
		return cHold
	case theme.ColorNameError:
		return cStop
	// Text drawn ON a filled accent. The accents are light, so ink is Night;
	// the red is dark enough that white text reads better on it.
	case theme.ColorNameForegroundOnPrimary, theme.ColorNameForegroundOnSuccess,
		theme.ColorNameForegroundOnWarning:
		return cNight
	case theme.ColorNameForegroundOnError:
		return cText
	default:
		// Anything unmapped uses the default dark theme so new Fyne color roles
		// still resolve sensibly.
		return t.Theme.Color(n, variant)
	}
}

func (t *brandTheme) lightColor(n fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground, theme.ColorNameScrollBarBackground:
		return cPaper
	case theme.ColorNameForeground:
		return cInk
	case theme.ColorNamePrimary, theme.ColorNameFocus, theme.ColorNameHyperlink:
		return cLightAnchor
	case theme.ColorNameButton, theme.ColorNameInputBackground,
		theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground,
		theme.ColorNameHeaderBackground:
		return cLightPanel
	case theme.ColorNameSelection:
		return withAlpha(cLightAnchor, 0x38)
	case theme.ColorNameHover, theme.ColorNamePressed, theme.ColorNameDisabledButton:
		return cLightRaised
	case theme.ColorNameInputBorder, theme.ColorNameSeparator:
		return cLightLine
	case theme.ColorNameScrollBar:
		return withAlpha(cLightLine, 0xD0)
	case theme.ColorNamePlaceHolder, theme.ColorNameDisabled:
		return cLightMuted
	case theme.ColorNameShadow:
		return color.NRGBA{A: 0x30}
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 0x0C, G: 0x78, B: 0x5D, A: 0xFF}
	case theme.ColorNameWarning:
		return color.NRGBA{R: 0x8A, G: 0x58, B: 0x00, A: 0xFF}
	case theme.ColorNameError:
		return color.NRGBA{R: 0xB1, G: 0x2B, B: 0x34, A: 0xFF}
	case theme.ColorNameForegroundOnPrimary, theme.ColorNameForegroundOnSuccess,
		theme.ColorNameForegroundOnWarning, theme.ColorNameForegroundOnError:
		return cLightPanel
	default:
		return t.Theme.Color(n, variant)
	}
}

// Font selects a brand face per text style. The branch order matches Fyne's own
// default selection (Monospace, then Bold, then Symbol, then regular).
func (t *brandTheme) Font(s fyne.TextStyle) fyne.Resource {
	if s.Symbol {
		return t.Theme.Font(s)
	}
	switch t.appearance.Font {
	case fontSystem:
		return t.Theme.Font(s)
	case fontMonospace:
		return resMono
	}
	switch {
	case s.Monospace:
		return resMono
	case s.Bold:
		return resDisplayBold
	default:
		return resBody
	}
}

// Size applies the brand's small-radius geometry; everything else delegates.
func (t *brandTheme) Size(n fyne.ThemeSizeName) float32 {
	switch n {
	case theme.SizeNameInputRadius:
		return radiusMd
	case theme.SizeNameSelectionRadius, theme.SizeNameScrollBarRadius:
		return radiusSm
	default:
		return t.Theme.Size(n)
	}
}
