// Package brandrelease validates a pinned official brand archive and its repository consumers.
package brandrelease

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// Pin identifies the one formal release accepted by this repository.
type Pin struct {
	Package             string `json:"package"`
	BrandVersion        string `json:"brand_version"`
	BrandBuilderVersion string `json:"brandbuilder_version"`
	ReleaseURL          string `json:"release_url"`
	SHA256              string `json:"sha256"`
}

// Artifact is one file declared by the official kit manifest.
type Artifact struct {
	Path   string `json:"path"`
	Bytes  int64  `json:"bytes"`
	SHA256 string `json:"sha256"`
}

// Manifest is the official kit's file inventory.
type Manifest struct {
	Name    string     `json:"name"`
	Version string     `json:"version"`
	Files   []Artifact `json:"files"`
}

// Mapping names one official source and its repository consumers.
type Mapping struct {
	Source    string   `json:"source"`
	Targets   []string `json:"targets"`
	Purpose   string   `json:"purpose"`
	Transform string   `json:"transform,omitempty"`
}

// ConsumerMap is the repository-owned source-to-target contract.
type ConsumerMap struct {
	Version  int       `json:"version"`
	Mappings []Mapping `json:"mappings"`
}

// Kit is a validated archive held in memory for import and offline checks.
type Kit struct {
	Pin      Pin
	Manifest Manifest
	Files    map[string][]byte
	Archive  []byte
}

// ReadPin loads a checked release pin.
func ReadPin(filename string) (Pin, error) {
	var pin Pin
	data, err := os.ReadFile(filename)
	if err != nil {
		return pin, fmt.Errorf("read brand source pin: %w", err)
	}
	if err := json.Unmarshal(data, &pin); err != nil {
		return pin, fmt.Errorf("parse brand source pin: %w", err)
	}
	if pin.Package == "" || pin.BrandVersion == "" || pin.BrandBuilderVersion == "" || !strings.HasPrefix(pin.ReleaseURL, "https://github.com/shruggietech/shruggie-brand/releases/download/") || len(pin.SHA256) != 64 {
		return pin, fmt.Errorf("brand source pin has incomplete release identity")
	}
	if err := SafePath(pin.Package + ".zip"); err != nil || strings.Contains(pin.Package, "/") || strings.Contains(pin.Package, "..") {
		return pin, fmt.Errorf("brand source pin has unsafe package name %q", pin.Package)
	}
	if _, err := hex.DecodeString(pin.SHA256); err != nil {
		return pin, fmt.Errorf("brand source pin has invalid SHA-256: %w", err)
	}
	if !strings.HasSuffix(pin.ReleaseURL, "/"+pin.Package+".zip") {
		return pin, fmt.Errorf("brand source URL does not match package %q", pin.Package)
	}
	return pin, nil
}

// SafePath rejects archive and mapping paths that could escape the checkout.
func SafePath(value string) error {
	if value == "" || strings.ContainsAny(value, "\\:\x00") || strings.HasPrefix(value, "/") || path.Clean(value) != value || value == "." || value == ".." || strings.HasPrefix(value, "../") {
		return fmt.Errorf("unsafe path %q", value)
	}
	return nil
}

// Open validates the exact archive digest, its complete file inventory, and bundle identity.
func Open(pin Pin, archivePath string) (*Kit, error) {
	if filepath.Base(archivePath) != pin.Package+".zip" {
		return nil, fmt.Errorf("archive filename does not match pinned package %q", pin.Package)
	}
	archive, err := os.ReadFile(archivePath)
	if err != nil {
		return nil, fmt.Errorf("read brand archive: %w", err)
	}
	digest := sha256.Sum256(archive)
	if got := hex.EncodeToString(digest[:]); !strings.EqualFold(got, pin.SHA256) {
		return nil, fmt.Errorf("brand archive SHA-256 %s differs from pin %s", got, pin.SHA256)
	}
	reader, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		return nil, fmt.Errorf("open brand ZIP: %w", err)
	}
	files := make(map[string][]byte, len(reader.File))
	seen := make(map[string]bool, len(reader.File))
	var total int64
	for _, entry := range reader.File {
		name := entry.Name
		if err := SafePath(name); err != nil {
			return nil, fmt.Errorf("archive: %w", err)
		}
		key := strings.ToLower(name)
		if seen[key] {
			return nil, fmt.Errorf("archive has duplicate path %q", name)
		}
		seen[key] = true
		if !entry.Mode().IsRegular() || entry.UncompressedSize64 > 32<<20 {
			return nil, fmt.Errorf("archive path %q is not a bounded regular file", name)
		}
		total += int64(entry.UncompressedSize64)
		if total > 128<<20 {
			return nil, fmt.Errorf("archive expanded size exceeds 128 MiB")
		}
		stream, err := entry.Open()
		if err != nil {
			return nil, fmt.Errorf("open archive file %q: %w", name, err)
		}
		data, readErr := io.ReadAll(io.LimitReader(stream, int64(entry.UncompressedSize64)+1))
		closeErr := stream.Close()
		if readErr != nil || closeErr != nil {
			return nil, fmt.Errorf("read archive file %q: %v; close: %v", name, readErr, closeErr)
		}
		if uint64(len(data)) != entry.UncompressedSize64 {
			return nil, fmt.Errorf("archive file %q has unexpected size", name)
		}
		files[name] = data
	}
	var manifest Manifest
	if err := json.Unmarshal(files["manifest.json"], &manifest); err != nil {
		return nil, fmt.Errorf("parse brand manifest: %w", err)
	}
	if manifest.Name != "go-schedule-brand-kit" || manifest.Version != pin.BrandVersion || len(manifest.Files) == 0 {
		return nil, fmt.Errorf("brand manifest identity differs from pin")
	}
	declared := make(map[string]bool, len(manifest.Files))
	for _, item := range manifest.Files {
		if err := SafePath(item.Path); err != nil {
			return nil, fmt.Errorf("manifest: %w", err)
		}
		key := strings.ToLower(item.Path)
		if declared[key] {
			return nil, fmt.Errorf("manifest has duplicate path %q", item.Path)
		}
		declared[key] = true
		data, ok := files[item.Path]
		if !ok {
			return nil, fmt.Errorf("manifest file %q is missing from archive", item.Path)
		}
		sum := sha256.Sum256(data)
		if int64(len(data)) != item.Bytes || !strings.EqualFold(hex.EncodeToString(sum[:]), item.SHA256) {
			return nil, fmt.Errorf("manifest file %q has wrong size or SHA-256", item.Path)
		}
		if isText(item.Path) && (!utf8.Valid(data) || bytes.HasPrefix(data, []byte("\xef\xbb\xbf")) || bytes.Contains(data, []byte("\xef\xbf\xbd"))) {
			return nil, fmt.Errorf("manifest file %q is not clean UTF-8", item.Path)
		}
	}
	for name := range files {
		if !declared[strings.ToLower(name)] && name != "manifest.json" && name != "LICENSE" && name != "LICENSE-BRAND.md" && name != "NOTICE" {
			return nil, fmt.Errorf("archive has undeclared file %q", name)
		}
	}
	var bundle struct {
		Package struct {
			ID                  string `json:"id"`
			Filename            string `json:"filename"`
			BrandVersion        string `json:"brand_version"`
			BrandBuilderVersion string `json:"brandbuilder_version"`
		} `json:"package"`
		Publication struct {
			Tag string `json:"tag"`
		} `json:"publication"`
	}
	if err := json.Unmarshal(files["enforcement/bundle.json"], &bundle); err != nil {
		return nil, fmt.Errorf("parse brand bundle: %w", err)
	}
	if bundle.Package.ID != pin.Package || bundle.Package.Filename != pin.Package+".zip" || bundle.Package.BrandVersion != pin.BrandVersion || bundle.Package.BrandBuilderVersion != pin.BrandBuilderVersion || !strings.Contains(pin.ReleaseURL, "/"+bundle.Publication.Tag+"/") {
		return nil, fmt.Errorf("brand bundle identity differs from pin")
	}
	return &Kit{Pin: pin, Manifest: manifest, Files: files, Archive: archive}, nil
}

func isText(name string) bool {
	switch strings.ToLower(path.Ext(name)) {
	case ".css", ".html", ".js", ".json", ".jsx", ".md", ".py", ".rs", ".svg", ".toml", ".ts", ".tsx", ".txt", ".xml", ".yaml", ".yml":
		return true
	}
	return false
}

// VerifyChecksums requires the formal release inventory to agree with the pin.
func VerifyChecksums(pin Pin, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read release checksums: %w", err)
	}
	var matched bool
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 || strings.TrimPrefix(fields[1], "./") != pin.Package+".zip" {
			continue
		}
		if matched || !strings.EqualFold(fields[0], pin.SHA256) {
			return fmt.Errorf("release checksum entry for %s is duplicate or differs from pin", pin.Package)
		}
		matched = true
	}
	if !matched {
		return fmt.Errorf("release checksums omit %s.zip", pin.Package)
	}
	return nil
}
