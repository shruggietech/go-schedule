// Package userprefs owns the shared headless desktop preference contract.
package userprefs

import (
	"fmt"
	"strings"
)

// CurrentVersion is the durable desktop preference schema accepted by this build.
const CurrentVersion = 1

// Appearance is the complete desktop color-mode vocabulary.
type Appearance string

const (
	// AppearanceSystem follows the operating system's current color mode.
	AppearanceSystem Appearance = "system"
	// AppearanceLight always uses the light color mode.
	AppearanceLight Appearance = "light"
	// AppearanceDark always uses the dark color mode.
	AppearanceDark Appearance = "dark"
)

// Preferences is the shared versioned desktop-local preference document.
type Preferences struct {
	Version    int              `json:"version"`
	Appearance Appearance       `json:"appearance"`
	Popups     PopupPreferences `json:"popups"`
}

// PopupPreferences filters desktop notifications while the desktop is running.
type PopupPreferences struct {
	Enabled    bool     `json:"enabled"`
	Conditions []string `json:"conditions"`
	Severities []string `json:"severities"`
	DaemonIDs  []string `json:"daemonIds"`
}

// DefaultPopups disables delivery and includes every supported condition and source.
func DefaultPopups() PopupPreferences {
	return PopupPreferences{Conditions: []string{"failure", "success", "alert"}, Severities: []string{"info", "warning", "error"}, DaemonIDs: []string{}}
}

// Defaults constructs complete current preferences.
func Defaults() Preferences {
	return Preferences{Version: CurrentVersion, Appearance: AppearanceSystem, Popups: DefaultPopups()}
}

// ValidAppearance reports whether value is one of the supported color modes.
func ValidAppearance(value Appearance) bool {
	return value == AppearanceSystem || value == AppearanceLight || value == AppearanceDark
}

// ValidatePopups checks supported, nonempty conditions and severities plus source IDs.
func ValidatePopups(value PopupPreferences) error {
	if len(value.Conditions) == 0 || len(value.Severities) == 0 {
		return fmt.Errorf("choose at least one popup condition and severity")
	}
	if len(value.Conditions) > 3 || len(value.Severities) > 3 || len(value.DaemonIDs) > 64 {
		return fmt.Errorf("too many popup filters")
	}
	if !uniqueAllowed(value.Conditions, map[string]bool{"failure": true, "success": true, "alert": true}) || !uniqueAllowed(value.Severities, map[string]bool{"info": true, "warning": true, "error": true}) {
		return fmt.Errorf("unsupported or duplicate popup condition or severity; conditions: failure, success, alert; severities: info, warning, error")
	}
	seen := make(map[string]bool, len(value.DaemonIDs))
	for _, id := range value.DaemonIDs {
		if strings.TrimSpace(id) != id || id == "" || len(id) > 128 || seen[id] {
			return fmt.Errorf("popup daemon IDs must be nonempty, unique, at most 128 bytes, and have no surrounding whitespace")
		}
		seen[id] = true
	}
	return nil
}

func uniqueAllowed(values []string, allowed map[string]bool) bool {
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		if !allowed[value] || seen[value] {
			return false
		}
		seen[value] = true
	}
	return true
}

// Patch distinguishes an omitted setting from a false value or cleared source filter.
type Patch struct {
	Appearance      *Appearance
	PopupEnabled    *bool
	PopupConditions *[]string
	PopupSeverities *[]string
	PopupDaemonIDs  *[]string
}

func (p Patch) apply(prefs Preferences) Preferences {
	if p.Appearance != nil {
		prefs.Appearance = *p.Appearance
	}
	if p.PopupEnabled != nil {
		prefs.Popups.Enabled = *p.PopupEnabled
	}
	if p.PopupConditions != nil {
		prefs.Popups.Conditions = append([]string{}, (*p.PopupConditions)...)
	}
	if p.PopupSeverities != nil {
		prefs.Popups.Severities = append([]string{}, (*p.PopupSeverities)...)
	}
	if p.PopupDaemonIDs != nil {
		prefs.Popups.DaemonIDs = append([]string{}, (*p.PopupDaemonIDs)...)
	}
	return prefs
}

// ValidatePatch rejects malformed input before a missing store can initialize.
func ValidatePatch(p Patch) error {
	if p.Appearance == nil && p.PopupEnabled == nil && p.PopupConditions == nil && p.PopupSeverities == nil && p.PopupDaemonIDs == nil {
		return fmt.Errorf("provide at least one preference to change")
	}
	if p.Appearance != nil && !ValidAppearance(*p.Appearance) {
		return fmt.Errorf("appearance must be system, light, or dark")
	}
	return ValidatePopups(p.apply(Defaults()).Popups)
}
