package config

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// TestMigrateLegacyPaths_MovesLegacyDirAndDB verifies a pre-rebrand directory and
// database are moved onto the new names when the new location does not yet exist.
func TestMigrateLegacyPaths_MovesLegacyDirAndDB(t *testing.T) {
	base := t.TempDir()
	legacyDir := filepath.Join(base, legacyDirName)
	if err := os.MkdirAll(legacyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Seed the legacy DB plus a WAL sidecar to confirm both are renamed.
	writeFile(t, filepath.Join(legacyDir, legacyDBName), "db")
	writeFile(t, filepath.Join(legacyDir, legacyDBName+"-wal"), "wal")

	newDir := filepath.Join(base, currentDirName)
	MigrateLegacyPaths(Config{DataDir: newDir}, discardLogger())

	if _, err := os.Stat(legacyDir); !os.IsNotExist(err) {
		t.Errorf("legacy dir should have been moved away, stat err = %v", err)
	}
	if got := readFile(t, filepath.Join(newDir, currentDBName)); got != "db" {
		t.Errorf("renamed db content = %q, want %q", got, "db")
	}
	if got := readFile(t, filepath.Join(newDir, currentDBName+"-wal")); got != "wal" {
		t.Errorf("renamed wal content = %q, want %q", got, "wal")
	}
}

// TestMigrateLegacyPaths_NoopWhenNewExists verifies an existing new directory is
// left untouched (already migrated or a fresh install).
func TestMigrateLegacyPaths_NoopWhenNewExists(t *testing.T) {
	base := t.TempDir()
	legacyDir := filepath.Join(base, legacyDirName)
	newDir := filepath.Join(base, currentDirName)
	if err := os.MkdirAll(legacyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(newDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(legacyDir, legacyDBName), "old")

	MigrateLegacyPaths(Config{DataDir: newDir}, discardLogger())

	// Legacy data must be preserved (not deleted) and the new dir untouched.
	if _, err := os.Stat(filepath.Join(legacyDir, legacyDBName)); err != nil {
		t.Errorf("legacy data must be preserved when new dir exists: %v", err)
	}
	if _, err := os.Stat(filepath.Join(newDir, currentDBName)); !os.IsNotExist(err) {
		t.Errorf("new dir should not have gained a migrated db, stat err = %v", err)
	}
}

// TestMigrateLegacyPaths_NoopFreshInstall verifies no error/effect when neither
// the legacy nor the new directory exists.
func TestMigrateLegacyPaths_NoopFreshInstall(t *testing.T) {
	base := t.TempDir()
	newDir := filepath.Join(base, currentDirName)
	MigrateLegacyPaths(Config{DataDir: newDir}, discardLogger())
	if _, err := os.Stat(newDir); !os.IsNotExist(err) {
		t.Errorf("fresh install should not create the new dir here, stat err = %v", err)
	}
}

// TestMigrateLegacyPaths_SkipsCustomDataDir verifies a non-default DataDir (base
// name != goschedule) is never migrated.
func TestMigrateLegacyPaths_SkipsCustomDataDir(t *testing.T) {
	base := t.TempDir()
	legacyDir := filepath.Join(base, legacyDirName)
	if err := os.MkdirAll(legacyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	custom := filepath.Join(base, "my-custom-dir")
	MigrateLegacyPaths(Config{DataDir: custom}, discardLogger())

	if _, err := os.Stat(legacyDir); err != nil {
		t.Errorf("legacy dir must be untouched for a custom data dir: %v", err)
	}
	if _, err := os.Stat(custom); !os.IsNotExist(err) {
		t.Errorf("custom dir should not be created by migration, stat err = %v", err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
