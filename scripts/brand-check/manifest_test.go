package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fixture struct {
	root     string
	manifest manifestFile
	mapFile  consumerFile
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	root := t.TempDir()
	f := &fixture{
		root:     root,
		manifest: manifestFile{Name: "fixture", Version: "1.0.0"},
		mapFile:  consumerFile{Version: 1},
	}
	for _, file := range []struct {
		path string
		data string
	}{
		{"README.md", "# Brand\n"},
		{"REPOSITORY.md", "# Repository integration\n"},
		{"VERIFY.md", "Status: **PASS**\n"},
		{"brand-guide.pdf", "%PDF-fixture\n"},
		{"logos/mark.svg", `<svg xmlns="http://www.w3.org/2000/svg"><title>Mark</title><path d="M0 0h1v1z"/></svg>` + "\n"},
	} {
		f.write(t, filepath.Join("brand", file.path), []byte(file.data))
		if file.path != "REPOSITORY.md" && file.path != "VERIFY.md" {
			f.addArtifact(t, file.path)
		}
	}
	f.flush(t)
	return f
}

func (f *fixture) write(t *testing.T, path string, data []byte) {
	t.Helper()
	full := filepath.Join(f.root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func (f *fixture) addArtifact(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(f.root, "brand", filepath.FromSlash(path)))
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(data)
	f.manifest.Files = append(f.manifest.Files, artifact{Path: path, Bytes: int64(len(data)), SHA256: hex.EncodeToString(digest[:])})
}

func (f *fixture) flush(t *testing.T) {
	t.Helper()
	manifestData, err := json.Marshal(f.manifest)
	if err != nil {
		t.Fatal(err)
	}
	mapData, err := json.Marshal(f.mapFile)
	if err != nil {
		t.Fatal(err)
	}
	f.write(t, "brand/manifest.json", manifestData)
	f.write(t, "brand/repository-consumers.json", mapData)
}

func requireFailure(t *testing.T, failures []string, want string) {
	t.Helper()
	if !strings.Contains(strings.Join(failures, "\n"), want) {
		t.Fatalf("failures %q do not contain %q", failures, want)
	}
}

func TestCheckRepositoryValid(t *testing.T) {
	f := newFixture(t)
	stats, failures := checkRepository(f.root)
	if len(failures) != 0 {
		t.Fatalf("unexpected failures: %v", failures)
	}
	if stats.Artifacts != 3 || stats.SVGs != 1 {
		t.Fatalf("stats = %+v, want 3 artifacts and 1 SVG", stats)
	}
}

func TestCheckRepositoryManifestFailures(t *testing.T) {
	tests := []struct {
		name string
		edit func(*testing.T, *fixture)
		want string
	}{
		{"missing", func(t *testing.T, f *fixture) { os.Remove(filepath.Join(f.root, "brand", "README.md")) }, "missing artifact"},
		{"hash", func(t *testing.T, f *fixture) { f.write(t, "brand/README.md", []byte("changed\n")) }, "SHA-256"},
		{"malformed", func(t *testing.T, f *fixture) { f.write(t, "brand/manifest.json", []byte("{")) }, "parse"},
		{"escape", func(t *testing.T, f *fixture) { f.manifest.Files[0].Path = "../escape"; f.flush(t) }, "unsafe artifact path"},
		{"drive escape", func(t *testing.T, f *fixture) { f.manifest.Files[0].Path = "C:/escape"; f.flush(t) }, "unsafe artifact path"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := newFixture(t)
			test.edit(t, f)
			_, failures := checkRepository(f.root)
			requireFailure(t, failures, test.want)
		})
	}
}
