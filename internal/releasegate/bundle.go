package releasegate

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

const (
	maxBundleEntries = 512
	maxBundleFile    = 128 << 20
	maxBundleBytes   = 512 << 20
)

// ExtractBundle securely extracts a bounded evidence archive to a temporary
// directory. The caller must invoke cleanup when finished.
func ExtractBundle(bundlePath string) (root string, cleanup func(), err error) {
	archive, err := zip.OpenReader(bundlePath)
	if err != nil {
		return "", func() {}, fmt.Errorf("open evidence bundle: %w", err)
	}
	defer archive.Close()
	if len(archive.File) == 0 || len(archive.File) > maxBundleEntries {
		return "", func() {}, fmt.Errorf("evidence bundle contains %d entries; allowed range is 1..%d", len(archive.File), maxBundleEntries)
	}

	root, err = os.MkdirTemp("", "goschedule-release-evidence-")
	if err != nil {
		return "", func() {}, fmt.Errorf("create evidence workspace: %w", err)
	}
	cleanup = func() { _ = os.RemoveAll(root) }
	succeeded := false
	defer func() {
		if !succeeded {
			cleanup()
		}
	}()

	seen := make(map[string]bool, len(archive.File))
	var total uint64
	for _, entry := range archive.File {
		if reason := unsafeRelativePath(entry.Name); reason != "" {
			return "", cleanup, fmt.Errorf("unsafe archive entry %q: %s", entry.Name, reason)
		}
		if seen[entry.Name] {
			return "", cleanup, fmt.Errorf("duplicate archive entry %q", entry.Name)
		}
		seen[entry.Name] = true
		if entry.FileInfo().IsDir() || !entry.Mode().IsRegular() {
			return "", cleanup, fmt.Errorf("archive entry %q is not a regular file", entry.Name)
		}
		if entry.UncompressedSize64 > maxBundleFile {
			return "", cleanup, fmt.Errorf("archive entry %q exceeds %d bytes", entry.Name, maxBundleFile)
		}
		total += entry.UncompressedSize64
		if total > maxBundleBytes {
			return "", cleanup, fmt.Errorf("evidence bundle exceeds %d uncompressed bytes", maxBundleBytes)
		}

		target := filepath.Join(root, filepath.FromSlash(entry.Name))
		if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
			return "", cleanup, fmt.Errorf("create archive directory for %q: %w", entry.Name, err)
		}
		input, err := entry.Open()
		if err != nil {
			return "", cleanup, fmt.Errorf("open archive entry %q: %w", entry.Name, err)
		}
		output, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err != nil {
			input.Close()
			return "", cleanup, fmt.Errorf("create extracted entry %q: %w", entry.Name, err)
		}
		written, copyErr := io.Copy(output, io.LimitReader(input, maxBundleFile+1))
		closeOutputErr := output.Close()
		closeInputErr := input.Close()
		if copyErr != nil || closeOutputErr != nil || closeInputErr != nil {
			return "", cleanup, fmt.Errorf("extract archive entry %q: copy=%v close-output=%v close-input=%v", entry.Name, copyErr, closeOutputErr, closeInputErr)
		}
		if written != int64(entry.UncompressedSize64) {
			return "", cleanup, fmt.Errorf("archive entry %q wrote %d bytes; header declares %d", entry.Name, written, entry.UncompressedSize64)
		}
	}
	if !seen["evidence.json"] {
		return "", cleanup, fmt.Errorf("evidence bundle is missing evidence.json")
	}
	succeeded = true
	return root, cleanup, nil
}

// ValidateBundleContents rejects unlisted files from a canonical archive.
func ValidateBundleContents(root string, evidence Evidence) []string {
	expected := map[string]bool{"evidence.json": true}
	for _, attachment := range evidence.Attachments {
		expected[attachment.Path] = true
	}
	var failures []string
	observed := make(map[string]bool)
	err := filepath.WalkDir(root, func(name string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, name)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		observed[relative] = true
		if !expected[relative] {
			failures = append(failures, fmt.Sprintf("bundle contains unlisted file %q", relative))
		}
		return nil
	})
	if err != nil {
		failures = append(failures, fmt.Sprintf("walk extracted bundle: %v", err))
	}
	for name := range expected {
		if !observed[name] {
			failures = append(failures, fmt.Sprintf("bundle is missing listed file %q", name))
		}
	}
	sort.Strings(failures)
	return failures
}
