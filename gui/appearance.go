package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const (
	appearanceModePreferenceKey = "appearance.mode"
	appearanceFontPreferenceKey = "appearance.font"
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
	fontMonospace fontChoice = "monospace"
)

type appearancePreferences struct {
	Mode appearanceMode
	Font fontChoice
}

func defaultAppearancePreferences() appearancePreferences {
	return appearancePreferences{Mode: appearanceDark, Font: fontBrand}
}

func loadAppearancePreferences(preferences fyne.Preferences) appearancePreferences {
	return appearancePreferences{
		Mode: normalizeAppearanceMode(preferences.String(appearanceModePreferenceKey)),
		Font: normalizeFontChoice(preferences.String(appearanceFontPreferenceKey)),
	}
}

func saveAppearancePreferences(preferences fyne.Preferences, value appearancePreferences) {
	value = value.normalized()
	preferences.SetString(appearanceModePreferenceKey, string(value.Mode))
	preferences.SetString(appearanceFontPreferenceKey, string(value.Font))
}

func resetAppearancePreferences(preferences fyne.Preferences) {
	saveAppearancePreferences(preferences, defaultAppearancePreferences())
}

func (p appearancePreferences) normalized() appearancePreferences {
	return appearancePreferences{
		Mode: normalizeAppearanceMode(string(p.Mode)),
		Font: normalizeFontChoice(string(p.Font)),
	}
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
	case fontBrand, fontSystem, fontMonospace:
		return fontChoice(value)
	default:
		return fontBrand
	}
}

func appearanceModeLabels() []string {
	return []string{"Dark", "Light", "Follow system"}
}

func fontChoiceLabels() []string {
	return []string{"Brand", "System", "Monospace"}
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
	case "System":
		return fontSystem
	case "Monospace":
		return fontMonospace
	default:
		return fontBrand
	}
}

func (f fontChoice) label() string {
	switch normalizeFontChoice(string(f)) {
	case fontSystem:
		return "System"
	case fontMonospace:
		return "Monospace"
	default:
		return "Brand"
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
