package clientprofile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

// DefaultPath resolves the shared interactive-user profile document.
func DefaultPath() string {
	root, err := os.UserConfigDir()
	if err != nil || root == "" {
		return ""
	}
	return filepath.Join(root, "go-schedule", "desktop", "profiles.json")
}

// Store serializes profile mutations within this process and with an exclusive lock file across processes.
type Store struct {
	path string
	now  func() time.Time
	mu   sync.Mutex
}

// NewStore creates a profile store at path, or the platform default when path is empty.
func NewStore(path string) *Store {
	if path == "" {
		path = DefaultPath()
	}
	return &Store{path: path, now: time.Now}
}

// Path returns the resolved profile document path.
func (s *Store) Path() string { return s.path }

// Load reads and validates one complete profile document.
func (s *Store) Load() (Collection, error) {
	if s.path == "" {
		return Collection{}, fmt.Errorf("profile path is unavailable")
	}
	data, err := os.ReadFile(s.path)
	if errors.Is(err, fs.ErrNotExist) {
		return Collection{Version: CurrentVersion, Profiles: []Profile{}}, nil
	}
	if err != nil {
		return Collection{}, fmt.Errorf("read connection profiles: %w", err)
	}
	var collection Collection
	if err := json.Unmarshal(data, &collection); err != nil {
		return Collection{}, fmt.Errorf("decode connection profiles: %w", err)
	}
	return validateCollection(collection)
}

// Add persists one new profile and returns its normalized value.
func (s *Store) Add(profile Profile) (Profile, error) {
	var result Profile
	err := s.update(func(collection *Collection) error {
		if len(collection.Profiles) >= MaxProfiles {
			return ErrInvalid
		}
		for _, candidate := range collection.Profiles {
			if candidate.ID == profile.ID {
				return ErrInvalid
			}
		}
		normalized, err := Normalize(profile)
		if err != nil {
			return err
		}
		collection.Profiles = append(collection.Profiles, normalized)
		result = normalized
		return nil
	})
	return result, err
}

// Replace updates an existing profile while preserving its ID and creation time.
func (s *Store) Replace(profile Profile) (Profile, error) {
	var result Profile
	err := s.update(func(collection *Collection) error {
		for i := range collection.Profiles {
			if collection.Profiles[i].ID == profile.ID {
				profile.CreatedAt = collection.Profiles[i].CreatedAt
				normalized, err := Normalize(profile)
				if err != nil {
					return err
				}
				collection.Profiles[i] = normalized
				result = normalized
				return nil
			}
		}
		return ErrNotFound
	})
	return result, err
}

// Rename changes only one profile's local presentation label.
func (s *Store) Rename(id, label string) (Profile, error) {
	collection, err := s.Load()
	if err != nil {
		return Profile{}, err
	}
	profile, err := collection.Find(id)
	if err != nil {
		return Profile{}, err
	}
	profile.Label = label
	profile.UpdatedAt = s.now().UTC()
	return s.Replace(profile)
}

// SetActive selects one remote profile or This computer when id is empty.
func (s *Store) SetActive(id string) error {
	return s.update(func(collection *Collection) error {
		if id != "" {
			if _, err := collection.Find(id); err != nil {
				return err
			}
		}
		collection.ActiveDesktopProfileID = id
		return nil
	})
}

// Remove deletes profile metadata. Credential deletion must occur before this call.
func (s *Store) Remove(id string) (Profile, error) {
	var removed Profile
	err := s.update(func(collection *Collection) error {
		for i := range collection.Profiles {
			if collection.Profiles[i].ID == id {
				removed = collection.Profiles[i]
				collection.Profiles = slices.Delete(collection.Profiles, i, i+1)
				if collection.ActiveDesktopProfileID == id {
					collection.ActiveDesktopProfileID = ""
				}
				return nil
			}
		}
		return ErrNotFound
	})
	return removed, err
}

func (s *Store) update(mutate func(*Collection) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.path == "" {
		return fmt.Errorf("profile path is unavailable")
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("create profile directory: %w", err)
	}
	lock, err := os.OpenFile(s.path+".lock", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if errors.Is(err, fs.ErrExist) {
		return ErrBusy
	}
	if err != nil {
		return fmt.Errorf("lock connection profiles: %w", err)
	}
	_ = lock.Close()
	defer os.Remove(s.path + ".lock")
	collection, err := s.Load()
	if err != nil {
		return err
	}
	if err := mutate(&collection); err != nil {
		return err
	}
	collection.Version = CurrentVersion
	validated, err := validateCollection(collection)
	if err != nil {
		return err
	}
	return writeAtomic(s.path, validated)
}

func validateCollection(collection Collection) (Collection, error) {
	if collection.Version > CurrentVersion {
		return Collection{}, ErrFuture
	}
	if collection.Version != CurrentVersion || len(collection.Profiles) > MaxProfiles {
		return Collection{}, ErrInvalid
	}
	seen := make(map[string]struct{}, len(collection.Profiles))
	for i := range collection.Profiles {
		profile, err := Normalize(collection.Profiles[i])
		if err != nil {
			return Collection{}, err
		}
		if _, exists := seen[profile.ID]; exists {
			return Collection{}, ErrInvalid
		}
		seen[profile.ID] = struct{}{}
		collection.Profiles[i] = profile
	}
	if collection.ActiveDesktopProfileID != "" {
		if _, exists := seen[collection.ActiveDesktopProfileID]; !exists {
			return Collection{}, ErrInvalid
		}
	}
	return collection, nil
}

func writeAtomic(path string, collection Collection) (resultErr error) {
	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	file, err := os.CreateTemp(filepath.Dir(path), ".profiles-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary profiles: %w", err)
	}
	temporary := file.Name()
	closed := false
	defer func() {
		if !closed {
			resultErr = errors.Join(resultErr, file.Close())
		}
		_ = os.Remove(temporary)
	}()
	if err := file.Chmod(0o600); err != nil {
		return fmt.Errorf("secure temporary profiles: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write temporary profiles: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("flush temporary profiles: %w", err)
	}
	if err := file.Close(); err != nil {
		closed = true
		return fmt.Errorf("close temporary profiles: %w", err)
	}
	closed = true
	if err := os.Rename(temporary, path); err != nil {
		return fmt.Errorf("replace connection profiles: %w", err)
	}
	return nil
}
