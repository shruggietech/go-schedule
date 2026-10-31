package userprefs

import (
	"os"
	"path/filepath"
)

// Paths identifies the current desktop preference document.
type Paths struct{ Preferences string }

// DefaultPaths returns the user-specific desktop preference location.
func DefaultPaths() Paths {
	configDir, _ := os.UserConfigDir()
	var paths Paths
	if configDir != "" {
		paths.Preferences = filepath.Join(configDir, "go-schedule", "desktop", "preferences.json")
	}
	return paths
}
