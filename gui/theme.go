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

// brandTheme is the go-schedule Fyne theme. It is dark-first: it ignores the
// requested variant and always returns the dark palette, so the app presents
// the brand consistently regardless of the OS light/dark setting.
type brandTheme struct{ fyne.Theme }

// newBrandTheme returns the go-schedule brand theme.
func newBrandTheme() fyne.Theme { return &brandTheme{Theme: theme.DefaultTheme()} }

// withAlpha returns c with its alpha replaced.
func withAlpha(c color.NRGBA, a uint8) color.NRGBA { c.A = a; return c }

// Color maps Fyne's semantic color roles onto the brand palette. The variant is
// ignored on purpose (dark-first, see the type doc).
func (t *brandTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
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
		return t.Theme.Color(n, theme.VariantDark)
	}
}

// Font selects a brand face per text style. The branch order matches Fyne's own
// default selection (Monospace, then Bold, then Symbol, then regular).
func (t *brandTheme) Font(s fyne.TextStyle) fyne.Resource {
	switch {
	case s.Monospace:
		return resMono
	case s.Bold:
		return resDisplayBold
	case s.Symbol:
		return t.Theme.Font(s) // keep the default symbol font
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
