package gui

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/shruggietech/go-schedule/internal/api/server"
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
	Runtime                 server.RuntimeInfoResponse
	OwnedMachineDataRoot    string
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
	wipe := "No built-in data wipe on this platform"
	if inputs.GOOS == "windows" {
		wipe = "Removed by an explicit data wipe"
	}
	software := "Removed with the application"

	machineScope := func(path string) (storageScope, string) {
		if path == "" {
			return storageScopeMachine, wipe
		}
		if pathWithin(inputs.OwnedMachineDataRoot, path) {
			return storageScopeMachine, wipe
		}
		return storageScopeExternal, "Preserved; outside the application-owned data root"
	}
	dataScope, dataWipe := machineScope(inputs.Runtime.DataDir)
	databaseScope, databaseWipe := machineScope(inputs.Runtime.DatabasePath)
	configScope, configWipe := machineScope(inputs.Runtime.ConfigPath)
	logScope, logWipe := machineScope(inputs.Runtime.LogPath)
	runtimeScope, runtimeWipe := machineScope(inputs.Runtime.LockPath)
	locations := []storageLocation{
		newStorageLocation("Machine data", inputs.Runtime.DataDir, dataScope, preserve, dataWipe, stat),
		newStorageLocation("Task database", inputs.Runtime.DatabasePath, databaseScope, preserve, databaseWipe, stat),
		newStorageLocation("Configuration", inputs.Runtime.ConfigPath, configScope, preserve, configWipe, stat),
		newStorageLocation("Logs", inputs.Runtime.LogPath, logScope, preserve, logWipe, stat),
		newStorageLocation("Runtime state", inputs.Runtime.LockPath, runtimeScope, preserve, runtimeWipe, stat),
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
