package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var retiredPreferenceKeys = []string{"appearance.font", "appearance.scroll_sensitivity"}

func loadPreferences(deps Dependencies) (DesktopPreferences, error) {
	if !filepath.IsAbs(deps.Paths.Preferences) {
		return DesktopPreferences{}, fmt.Errorf("desktop preference path is unavailable")
	}
	data, err := os.ReadFile(deps.Paths.Preferences)
	if err == nil {
		return decodeCurrentPreferences(data)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return DesktopPreferences{}, fmt.Errorf("read desktop preferences: %w", err)
	}
	prefs := DesktopPreferences{Version: CurrentPreferenceVersion, Appearance: AppearanceSystem, Transition: PreferenceTransition{Status: "not_found", Retired: append([]string(nil), retiredPreferenceKeys...)}}
	legacy, legacyErr := os.ReadFile(deps.Paths.LegacyPreferences)
	switch {
	case legacyErr == nil:
		var values map[string]any
		if json.Unmarshal(legacy, &values) != nil {
			prefs.Transition.Status = "invalid"
		} else if mode, ok := values["appearance.mode"].(string); ok && validAppearance(Appearance(mode)) {
			prefs.Appearance = Appearance(mode)
			prefs.Transition.Status = "migrated"
		} else {
			prefs.Transition.Status = "invalid"
		}
	case errors.Is(legacyErr, fs.ErrNotExist):
		prefs.Transition.Status = "not_found"
	default:
		prefs.Transition.Status = "unreadable"
	}
	if err := writePreferences(deps, prefs); err != nil {
		return DesktopPreferences{}, err
	}
	return prefs, nil
}

func decodeCurrentPreferences(data []byte) (DesktopPreferences, error) {
	var prefs DesktopPreferences
	if err := json.Unmarshal(data, &prefs); err != nil {
		return DesktopPreferences{}, fmt.Errorf("decode desktop preferences: %w", err)
	}
	if prefs.Version != CurrentPreferenceVersion {
		return DesktopPreferences{}, fmt.Errorf("desktop preferences use unsupported version %d", prefs.Version)
	}
	if !validAppearance(prefs.Appearance) {
		return DesktopPreferences{}, fmt.Errorf("desktop preferences contain an unsupported appearance")
	}
	prefs.Transition.Retired = append([]string(nil), prefs.Transition.Retired...)
	return prefs, nil
}

func writePreferences(deps Dependencies, prefs DesktopPreferences) (resultErr error) {
	if !filepath.IsAbs(deps.Paths.Preferences) {
		return fmt.Errorf("desktop preference path is unavailable")
	}
	if !validAppearance(prefs.Appearance) {
		return fmt.Errorf("unsupported appearance %q", prefs.Appearance)
	}
	data, err := json.MarshalIndent(prefs, "", "  ")
	if err != nil {
		return fmt.Errorf("encode desktop preferences: %w", err)
	}
	data = append(data, '\n')
	directory := filepath.Dir(deps.Paths.Preferences)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return fmt.Errorf("create desktop preference directory: %w", err)
	}
	if err := os.Chmod(directory, 0o700); err != nil {
		return fmt.Errorf("secure desktop preference directory: %w", err)
	}
	file, err := os.CreateTemp(directory, ".preferences-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary desktop preferences: %w", err)
	}
	temporary := file.Name()
	closed := false
	defer func() {
		if !closed {
			resultErr = errors.Join(resultErr, file.Close())
		}
		if removeErr := os.Remove(temporary); removeErr != nil && !errors.Is(removeErr, fs.ErrNotExist) {
			resultErr = errors.Join(resultErr, removeErr)
		}
	}()
	if err := file.Chmod(0o600); err != nil {
		return fmt.Errorf("secure temporary desktop preferences: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write temporary desktop preferences: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("flush temporary desktop preferences: %w", err)
	}
	if err := file.Close(); err != nil {
		closed = true
		return fmt.Errorf("close temporary desktop preferences: %w", err)
	}
	closed = true
	if err := deps.Rename(temporary, deps.Paths.Preferences); err != nil {
		return fmt.Errorf("replace desktop preferences: %w", err)
	}
	return nil
}
