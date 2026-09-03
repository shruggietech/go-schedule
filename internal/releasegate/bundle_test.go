package releasegate

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractBundleRejectsTraversalAndDuplicateEntries(t *testing.T) {
	t.Parallel()

	for _, entries := range [][]string{{"../escape"}, {"evidence.json", "evidence.json"}} {
		bundle := filepath.Join(t.TempDir(), "evidence.zip")
		file, err := os.Create(bundle)
		if err != nil {
			t.Fatal(err)
		}
		writer := zip.NewWriter(file)
		for _, name := range entries {
			entry, createErr := writer.Create(name)
			if createErr != nil {
				t.Fatal(createErr)
			}
			if _, writeErr := entry.Write([]byte("fixture")); writeErr != nil {
				t.Fatal(writeErr)
			}
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}

		if _, cleanup, err := ExtractBundle(bundle); err == nil {
			cleanup()
			t.Fatalf("ExtractBundle(%v) unexpectedly succeeded", entries)
		}
	}
}

func TestValidateBundleContentsRejectsUnlistedFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	for _, name := range []string{"evidence.json", "extra.txt"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("fixture"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	failures := ValidateBundleContents(root, Evidence{})
	if !containsFailure(failures, "unlisted file") {
		t.Fatalf("ValidateBundleContents() failures = %v", failures)
	}
}

func TestExtractBundleRejectsEntryCountAndNonRegularFiles(t *testing.T) {
	t.Parallel()

	t.Run("entry count", func(t *testing.T) {
		bundle := filepath.Join(t.TempDir(), "many.zip")
		file, err := os.Create(bundle)
		if err != nil {
			t.Fatal(err)
		}
		writer := zip.NewWriter(file)
		for i := 0; i <= maxBundleEntries; i++ {
			entry, createErr := writer.Create(fmt.Sprintf("entry-%03d", i))
			if createErr != nil {
				t.Fatal(createErr)
			}
			if _, writeErr := entry.Write([]byte("x")); writeErr != nil {
				t.Fatal(writeErr)
			}
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
		if _, cleanup, err := ExtractBundle(bundle); err == nil {
			cleanup()
			t.Fatal("ExtractBundle() unexpectedly accepted too many entries")
		}
	})

	t.Run("symbolic link", func(t *testing.T) {
		bundle := filepath.Join(t.TempDir(), "link.zip")
		file, err := os.Create(bundle)
		if err != nil {
			t.Fatal(err)
		}
		writer := zip.NewWriter(file)
		header := &zip.FileHeader{Name: "evidence.json"}
		header.SetMode(os.ModeSymlink | 0o777)
		entry, err := writer.CreateHeader(header)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write([]byte("target")); err != nil {
			t.Fatal(err)
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
		if _, cleanup, err := ExtractBundle(bundle); err == nil {
			cleanup()
			t.Fatal("ExtractBundle() unexpectedly accepted a symbolic link")
		}
	})
}
