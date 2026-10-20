package userprefs

import (
	"os"
	"path/filepath"
	"runtime"
)

// Paths identifies current preferences and their retired Fyne migration source.
type Paths struct{ Preferences, LegacyPreferences string }

// DefaultPaths preserves the existing user-specific desktop and retired Fyne locations.
func DefaultPaths() Paths {
	configDir, _ := os.UserConfigDir()
	homeDir, _ := os.UserHomeDir()
	var paths Paths
	if configDir != "" {
		paths.Preferences = filepath.Join(configDir, "go-schedule", "desktop", "preferences.json")
		paths.LegacyPreferences = filepath.Join(configDir, "fyne", "tech.shruggie.goschedule", "preferences.json")
	}
	if runtime.GOOS == "darwin" && homeDir != "" {
		paths.LegacyPreferences = filepath.Join(homeDir, "Library", "Preferences", "fyne", "tech.shruggie.goschedule", "preferences.json")
	}
	return paths
}
