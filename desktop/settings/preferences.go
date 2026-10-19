package settings

import (
	"fmt"
	"path/filepath"

	"github.com/shruggietech/go-schedule/internal/userprefs"
)

func preferenceStore(deps Dependencies) (*userprefs.Store, error) {
	if !filepath.IsAbs(deps.Paths.Preferences) {
		return nil, fmt.Errorf("desktop preference path is unavailable")
	}
	return userprefs.NewStore(userprefs.Options{Path: deps.Paths.Preferences, LegacyPath: deps.Paths.LegacyPreferences, Rename: deps.Rename}), nil
}

func loadPreferences(deps Dependencies) (DesktopPreferences, error) {
	store, err := preferenceStore(deps)
	if err != nil {
		return DesktopPreferences{}, err
	}
	return store.Load()
}
