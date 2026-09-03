package gui

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/shruggietech/go-schedule/internal/config"
)

type storageScope string

const (
	storageScopeMachine  storageScope = "Machine"
	storageScopeUser     storageScope = "Current user"
	storageScopeRuntime  storageScope = "Running application"
	storageScopeExternal storageScope = "Configured outside application data"
)

type storageExistence string

const (
	storagePresent storageExistence = "Present"
	storageAbsent  storageExistence = "Not present"
	storageUnknown storageExistence = "Unable to inspect"
)

type storageLocation struct {
	Category            string
	Path                string
	Available           bool
	Scope               storageScope
	Existence           storageExistence
	SoftwareOnlyRemoval string
	ExplicitDataWipe    string
}

type storageLocationInputs struct {
	Config                  config.Config
	PreferencesRoot         string
	ExecutablePath          string
	GOOS                    string
	MaintenanceEvidencePath string
	Stat                    func(string) (os.FileInfo, error)
}

func resolveStorageLocations(inputs storageLocationInputs) []storageLocation {
	stat := inputs.Stat
	if stat == nil {
		stat = os.Stat
	}
	preserve := "Preserved"
	wipe := "Removed by an explicit data wipe"
	software := "Removed with the application"

	logScope := storageScopeMachine
	logSoftwareOnly, logWipe := preserve, wipe
	if !pathWithin(inputs.Config.DataDir, inputs.Config.LogPath()) {
		logScope = storageScopeExternal
		logSoftwareOnly = "Preserved"
		logWipe = "Preserved; outside the application-owned data root"
	}
	locations := []storageLocation{
		newStorageLocation("Machine data", inputs.Config.DataDir, storageScopeMachine, preserve, wipe, stat),
		newStorageLocation("Task database", inputs.Config.DBPath(), storageScopeMachine, preserve, wipe, stat),
		newStorageLocation("Configuration", filepath.Join(inputs.Config.DataDir, "config.json"), storageScopeMachine, preserve, wipe, stat),
		newStorageLocation("Logs", inputs.Config.LogPath(), logScope, logSoftwareOnly, logWipe, stat),
		newStorageLocation("Runtime state", filepath.Join(inputs.Config.DataDir, "goschedd.lock"), storageScopeMachine, preserve, wipe, stat),
		newStorageLocation("Desktop application data", inputs.PreferencesRoot, storageScopeUser, preserve, wipe, stat),
		newStorageLocation("Desktop preferences", filepath.Join(inputs.PreferencesRoot, "preferences.json"), storageScopeUser, preserve, wipe, stat),
	}

	executableDir := ""
	if filepath.IsAbs(inputs.ExecutablePath) {
		executableDir = filepath.Dir(inputs.ExecutablePath)
	}
	locations = append(locations,
		newStorageLocation(
			"Executable directory",
			executableDir,
			storageScopeRuntime,
			"Installer-owned files are removed; development binaries are unaffected",
			"Installer-owned files are removed; development binaries are unaffected",
			stat,
		),
	)

	if executableDir != "" {
		docs := newStorageLocation("Installed documentation", filepath.Join(executableDir, "docs"), storageScopeRuntime, software, software, stat)
		if docs.Existence == storagePresent {
			locations = append(locations, docs)
		}
	}

	if inputs.GOOS == "windows" {
		locations = append(locations, newStorageLocation(
			"Maintenance evidence",
			inputs.MaintenanceEvidencePath,
			storageScopeMachine,
			"Not affected by software-only removal",
			"Removed after a complete wipe; retained only to report incomplete cleanup",
			stat,
		))
	}

	return locations
}

func pathWithin(root, candidate string) bool {
	if !filepath.IsAbs(root) || !filepath.IsAbs(candidate) {
		return false
	}
	relative, err := filepath.Rel(filepath.Clean(root), filepath.Clean(candidate))
	if err != nil {
		return false
	}
	return relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func newStorageLocation(category, path string, scope storageScope, softwareOnly, wipe string, stat func(string) (os.FileInfo, error)) storageLocation {
	location := storageLocation{
		Category:            category,
		Scope:               scope,
		Existence:           storageUnknown,
		SoftwareOnlyRemoval: softwareOnly,
		ExplicitDataWipe:    wipe,
	}
	if path == "" || !filepath.IsAbs(path) {
		return location
	}
	location.Path = filepath.Clean(path)
	location.Available = true
	location.Existence = inspectStoragePath(location.Path, stat)
	return location
}

func inspectStoragePath(path string, stat func(string) (os.FileInfo, error)) storageExistence {
	_, err := stat(path)
	switch {
	case err == nil:
		return storagePresent
	case errors.Is(err, fs.ErrNotExist):
		return storageAbsent
	default:
		return storageUnknown
	}
}

func storageCategories(locations []storageLocation) []string {
	categories := make([]string, len(locations))
	for i, location := range locations {
		categories[i] = location.Category
	}
	return categories
}
