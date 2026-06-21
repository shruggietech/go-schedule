package config

import (
	"log/slog"
	"os"
	"path/filepath"
)

// legacyDirName is the pre-rebrand data-directory base name (go-scheduler era).
const legacyDirName = "goscheduler"

// currentDirName is the data-directory base name after the go-schedule rebrand.
const currentDirName = "goschedule"

// legacyDBName / currentDBName are the SQLite filenames before and after the
// rebrand.
const (
	legacyDBName  = "goscheduler.db"
	currentDBName = "goschedule.db"
)

// MigrateLegacyPaths performs a one-time, best-effort move of a pre-rebrand data
// directory (".../goscheduler") and database file ("goscheduler.db") onto the new
// "goschedule" names. It runs before the daemon creates its data directory.
//
// It is intentionally non-fatal: any failure is logged and the daemon proceeds
// with the new (possibly empty) location rather than blocking startup. Existing
// data is never deleted — on a failed move the legacy directory is left intact
// and its location is logged so the operator can recover it.
//
// The migration only applies to the default location (a DataDir whose base name
// is "goschedule"); a custom DataDir is left untouched.
func MigrateLegacyPaths(cfg Config, log *slog.Logger) {
	newDir := cfg.DataDir
	if filepath.Base(newDir) != currentDirName {
		return // custom data dir; nothing to migrate
	}
	legacyDir := filepath.Join(filepath.Dir(newDir), legacyDirName)

	// Already migrated or a fresh install: the new dir exists, so do nothing.
	if _, err := os.Stat(newDir); err == nil {
		return
	}
	// No legacy data to bring forward: fresh install.
	if _, err := os.Stat(legacyDir); err != nil {
		return
	}

	if err := os.Rename(legacyDir, newDir); err != nil {
		log.Warn("could not migrate legacy data directory; using a fresh location",
			"legacy", legacyDir, "new", newDir, "error", err)
		return
	}
	log.Info("migrated legacy data directory", "from", legacyDir, "to", newDir)

	// Rename the database (and its WAL/SHM sidecars) inside the moved directory.
	for _, suffix := range []string{"", "-wal", "-shm"} {
		oldPath := filepath.Join(newDir, legacyDBName+suffix)
		newPath := filepath.Join(newDir, currentDBName+suffix)
		if _, err := os.Stat(oldPath); err != nil {
			continue
		}
		if err := os.Rename(oldPath, newPath); err != nil {
			log.Warn("could not rename legacy database file",
				"from", oldPath, "to", newPath, "error", err)
		}
	}
}
