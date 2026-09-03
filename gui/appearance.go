package gui

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const (
	appearanceModePreferenceKey    = "appearance.mode"
	appearanceFontPreferenceKey    = "appearance.font"
	scrollSensitivityPreferenceKey = "appearance.scroll_sensitivity"

	minimumScrollSensitivity = 1.0
	maximumScrollSensitivity = 4.0
	scrollSensitivityStep    = 0.5
	defaultScrollSensitivity = 2.0
)

type appearanceMode string

const (
	appearanceDark   appearanceMode = "dark"
	appearanceLight  appearanceMode = "light"
	appearanceSystem appearanceMode = "system"
)

type fontChoice string

const (
	fontBrand     fontChoice = "brand"
	fontSystem    fontChoice = "system"
	fontInter     fontChoice = "inter"
	fontUbuntu    fontChoice = "ubuntu"
	fontMonospace fontChoice = "monospace"
)

type appearancePreferences struct {
	Mode              appearanceMode
	Font              fontChoice
	ScrollSensitivity float64
}

func defaultAppearancePreferences() appearancePreferences {
	return appearancePreferences{
		Mode:              appearanceDark,
		Font:              fontSystem,
		ScrollSensitivity: defaultScrollSensitivity,
	}
}

func loadAppearancePreferences(preferences fyne.Preferences) appearancePreferences {
	return appearancePreferences{
		Mode:              normalizeAppearanceMode(preferences.String(appearanceModePreferenceKey)),
		Font:              normalizeFontChoice(preferences.String(appearanceFontPreferenceKey)),
		ScrollSensitivity: normalizeScrollSensitivity(preferences.FloatWithFallback(scrollSensitivityPreferenceKey, defaultScrollSensitivity)),
	}
}

func saveAppearancePreferences(preferences fyne.Preferences, value appearancePreferences) {
	value = value.normalized()
	preferences.SetString(appearanceModePreferenceKey, string(value.Mode))
	preferences.SetString(appearanceFontPreferenceKey, string(value.Font))
	preferences.SetFloat(scrollSensitivityPreferenceKey, value.ScrollSensitivity)
}

func resetAppearancePreferences(preferences fyne.Preferences) {
	saveAppearancePreferences(preferences, defaultAppearancePreferences())
}

func (p appearancePreferences) normalized() appearancePreferences {
	return appearancePreferences{
		Mode:              normalizeAppearanceMode(string(p.Mode)),
		Font:              normalizeFontChoice(string(p.Font)),
		ScrollSensitivity: normalizeScrollSensitivity(p.ScrollSensitivity),
	}
}

func (p appearancePreferences) sameTheme(other appearancePreferences) bool {
	p = p.normalized()
	other = other.normalized()
	return p.Mode == other.Mode && p.Font == other.Font
}

func normalizeAppearanceMode(value string) appearanceMode {
	switch appearanceMode(value) {
	case appearanceDark, appearanceLight, appearanceSystem:
		return appearanceMode(value)
	default:
		return appearanceDark
	}
}

func normalizeFontChoice(value string) fontChoice {
	switch fontChoice(value) {
	case fontBrand, fontSystem, fontInter, fontUbuntu, fontMonospace:
		return fontChoice(value)
	default:
		return fontSystem
	}
}

func normalizeScrollSensitivity(value float64) float64 {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < minimumScrollSensitivity || value > maximumScrollSensitivity {
		return defaultScrollSensitivity
	}
	steps := math.Round((value - minimumScrollSensitivity) / scrollSensitivityStep)
	return minimumScrollSensitivity + steps*scrollSensitivityStep
}

func appearanceModeLabels() []string {
	return []string{"Dark", "Light", "Follow system"}
}

func fontChoiceLabels() []string {
	return []string{"System", "Geist (brand)", "Inter", "Ubuntu", "Monospace"}
}

func appearanceModeForLabel(label string) appearanceMode {
	switch label {
	case "Light":
		return appearanceLight
	case "Follow system":
		return appearanceSystem
	default:
		return appearanceDark
	}
}

func (m appearanceMode) label() string {
	switch normalizeAppearanceMode(string(m)) {
	case appearanceLight:
		return "Light"
	case appearanceSystem:
		return "Follow system"
	default:
		return "Dark"
	}
}

func fontChoiceForLabel(label string) fontChoice {
	switch label {
	case "Geist (brand)":
		return fontBrand
	case "Inter":
		return fontInter
	case "Ubuntu":
		return fontUbuntu
	case "Monospace":
		return fontMonospace
	default:
		return fontSystem
	}
}

func (f fontChoice) label() string {
	switch normalizeFontChoice(string(f)) {
	case fontBrand:
		return "Geist (brand)"
	case fontInter:
		return "Inter"
	case fontUbuntu:
		return "Ubuntu"
	case fontMonospace:
		return "Monospace"
	default:
		return "System"
	}
}

// themeVariant returns the effective variant and whether the mode is valid.
// Follow system preserves the variant selected by the Fyne driver.
func (m appearanceMode) themeVariant(system fyne.ThemeVariant) (fyne.ThemeVariant, bool) {
	switch m {
	case appearanceDark:
		return theme.VariantDark, true
	case appearanceLight:
		return theme.VariantLight, true
	case appearanceSystem:
		if system == theme.VariantLight || system == theme.VariantDark {
			return system, true
		}
		return theme.VariantDark, true
	default:
		return theme.VariantDark, false
	}
}
