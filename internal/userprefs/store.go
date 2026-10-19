package userprefs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/shruggietech/go-schedule/internal/lock"
)

// ErrBusy reports a preference writer that did not finish within the bounded wait.
var ErrBusy = errors.New("desktop preference store is busy; retry after the other writer finishes")

// Options identifies preference files and an optional atomic-replacement seam.
type Options struct {
	Path       string
	LegacyPath string
	Rename     func(string, string) error
}

// Store coordinates scoped updates across desktop and CLI processes.
type Store struct {
	path, legacyPath string
	rename           func(string, string) error
	mu               sync.Mutex
}

// NewStore uses production paths only when no explicit current path was supplied.
func NewStore(options Options) *Store {
	if options.Path == "" {
		paths := DefaultPaths()
		options.Path = paths.Preferences
		if options.LegacyPath == "" {
			options.LegacyPath = paths.LegacyPreferences
		}
	}
	if options.Rename == nil {
		options.Rename = os.Rename
	}
	return &Store{path: options.Path, legacyPath: options.LegacyPath, rename: options.Rename}
}

// Path returns the current preference document's resolved location.
func (s *Store) Path() string { return s.path }

// Load reads the current document, initializing or migrating only when absent.
func (s *Store) Load() (Preferences, error) {
	if !filepath.IsAbs(s.path) {
		return Preferences{}, fmt.Errorf("desktop preference path is unavailable or not absolute")
	}
	// Existing documents are atomically replaced, so readers need no write access or lock.
	data, err := os.ReadFile(s.path)
	if err == nil {
		return Decode(data)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return Preferences{}, fmt.Errorf("read desktop preferences: %w", err)
	}
	var prefs Preferences
	err = s.locked(func() error { var err error; prefs, err = s.load(); return err })
	return prefs, err
}

// Patch atomically changes only supplied fields after validating the complete patch.
func (s *Store) Patch(patch Patch) (Preferences, error) {
	if err := ValidatePatch(patch); err != nil {
		return Preferences{}, err
	}
	var prefs Preferences
	err := s.locked(func() error {
		current, err := s.load()
		if err != nil {
			return err
		}
		prefs = patch.apply(current)
		return s.write(prefs)
	})
	return prefs, err
}

// Reset restores appearance, popups, or all desktop preferences to their defaults.
// The explicit all scope also repairs an invalid current preference document.
func (s *Store) Reset(scope string) (Preferences, error) {
	if scope != "all" && scope != "appearance" && scope != "popups" {
		return Preferences{}, fmt.Errorf("reset scope must be appearance, popups, or all")
	}
	var prefs Preferences
	err := s.locked(func() error {
		if scope == "all" {
			prefs = Defaults("not_required")
		} else {
			var err error
			prefs, err = s.load()
			if err != nil {
				return err
			}
			if scope == "appearance" {
				prefs.Appearance = AppearanceSystem
			} else {
				prefs.Popups = DefaultPopups()
			}
		}
		return s.write(prefs)
	})
	return prefs, err
}

func (s *Store) locked(action func() error) (resultErr error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !filepath.IsAbs(s.path) {
		return fmt.Errorf("desktop preference path is unavailable or not absolute")
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("create desktop preference directory: %w", err)
	}
	deadline := time.Now().Add(time.Second)
	for {
		held, err := lock.AcquireFile(s.path+".lock", 0o600)
		if errors.Is(err, lock.ErrHeld) {
			if time.Now().After(deadline) {
				return ErrBusy
			}
			time.Sleep(5 * time.Millisecond)
			continue
		}
		if err != nil {
			return fmt.Errorf("lock desktop preferences: %w", err)
		}
		// Never unlink an OS lock file: doing so could let later writers lock a different inode.
		defer func() { resultErr = errors.Join(resultErr, held.Release()) }()
		return action()
	}
}

func (s *Store) load() (Preferences, error) {
	data, err := os.ReadFile(s.path)
	if err == nil {
		return Decode(data)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return Preferences{}, fmt.Errorf("read desktop preferences: %w", err)
	}
	prefs := Defaults("not_found")
	legacy, legacyErr := os.ReadFile(s.legacyPath)
	switch {
	case legacyErr == nil:
		var values map[string]any
		if json.Unmarshal(legacy, &values) != nil {
			prefs.Transition.Status = "invalid"
		} else if mode, ok := values["appearance.mode"].(string); ok && ValidAppearance(Appearance(mode)) {
			prefs.Appearance = Appearance(mode)
			prefs.Transition.Status = "migrated"
		} else {
			prefs.Transition.Status = "invalid"
		}
	case errors.Is(legacyErr, fs.ErrNotExist) || s.legacyPath == "":
		prefs.Transition.Status = "not_found"
	default:
		prefs.Transition.Status = "unreadable"
	}
	if err := s.write(prefs); err != nil {
		return Preferences{}, err
	}
	return prefs, nil
}

// Decode validates an existing schema-v1 document and fills legacy popup defaults.
func Decode(data []byte) (Preferences, error) {
	var prefs Preferences
	if err := json.Unmarshal(data, &prefs); err != nil {
		return Preferences{}, fmt.Errorf("decode desktop preferences: %w", err)
	}
	if prefs.Version != CurrentVersion {
		return Preferences{}, fmt.Errorf("desktop preferences use unsupported version %d", prefs.Version)
	}
	if !ValidAppearance(prefs.Appearance) {
		return Preferences{}, fmt.Errorf("desktop preferences contain an unsupported appearance")
	}
	prefs.Transition.Retired = append([]string(nil), prefs.Transition.Retired...)
	if prefs.Popups.Conditions == nil && prefs.Popups.Severities == nil && prefs.Popups.DaemonIDs == nil {
		prefs.Popups = DefaultPopups()
	}
	if err := ValidatePopups(prefs.Popups); err != nil {
		return Preferences{}, err
	}
	return prefs, nil
}

func (s *Store) write(prefs Preferences) (resultErr error) {
	if prefs.Version != CurrentVersion || !ValidAppearance(prefs.Appearance) {
		return fmt.Errorf("unsupported desktop preference schema or appearance")
	}
	if err := ValidatePopups(prefs.Popups); err != nil {
		return err
	}
	data, err := json.MarshalIndent(prefs, "", "  ")
	if err != nil {
		return fmt.Errorf("encode desktop preferences: %w", err)
	}
	data = append(data, '\n')
	directory := filepath.Dir(s.path)
	if err := os.Chmod(directory, 0o700); err != nil {
		return fmt.Errorf("secure desktop preference directory: %w", err)
	}
	file, err := os.CreateTemp(directory, ".preferences-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary desktop preferences: %w", err)
	}
	temporary, closed := file.Name(), false
	defer func() {
		if !closed {
			resultErr = errors.Join(resultErr, file.Close())
		}
		if err := os.Remove(temporary); err != nil && !errors.Is(err, fs.ErrNotExist) {
			resultErr = errors.Join(resultErr, err)
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
	if err := s.rename(temporary, s.path); err != nil {
		return fmt.Errorf("replace desktop preferences: %w", err)
	}
	return nil
}
