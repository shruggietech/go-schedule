package brandrelease

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixtureArchive(t *testing.T, extra string, badManifest bool, symlink bool) (Pin, string) {
	t.Helper()
	packageID := "go-schedule-brand-2.0.0-bb2.4.0"
	bundle, err := json.Marshal(map[string]any{
		"package":     map[string]string{"id": packageID, "filename": packageID + ".zip", "brand_version": "2.0.0", "brandbuilder_version": "2.4.0"},
		"publication": map[string]string{"tag": "v2.4.0"},
	})
	if err != nil {
		t.Fatal(err)
	}
	contents := map[string][]byte{"README.md": []byte("# Brand\n"), "enforcement/bundle.json": bundle}
	manifest := Manifest{Name: "go-schedule-brand-kit", Version: "2.0.0"}
	for _, name := range []string{"README.md", "enforcement/bundle.json"} {
		data := contents[name]
		sum := sha256.Sum256(data)
		manifest.Files = append(manifest.Files, Artifact{Path: name, Bytes: int64(len(data)), SHA256: hex.EncodeToString(sum[:])})
	}
	if badManifest {
		manifest.Files[0].SHA256 = strings.Repeat("0", 64)
	}
	contents["manifest.json"], err = json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	var archive bytes.Buffer
	writer := zip.NewWriter(&archive)
	for _, name := range []string{"README.md", "enforcement/bundle.json", "manifest.json"} {
		entry, createErr := writer.Create(name)
		if createErr != nil {
			t.Fatal(createErr)
		}
		if _, writeErr := entry.Write(contents[name]); writeErr != nil {
			t.Fatal(writeErr)
		}
	}
	if extra != "" {
		header := &zip.FileHeader{Name: extra, Method: zip.Store}
		if symlink {
			header.SetMode(os.ModeSymlink | 0o777)
		}
		entry, createErr := writer.CreateHeader(header)
		if createErr != nil {
			t.Fatal(createErr)
		}
		if _, writeErr := entry.Write([]byte("target")); writeErr != nil {
			t.Fatal(writeErr)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	filename := filepath.Join(t.TempDir(), packageID+".zip")
	if err := os.WriteFile(filename, archive.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(archive.Bytes())
	pin := Pin{Package: packageID, BrandVersion: "2.0.0", BrandBuilderVersion: "2.4.0", ReleaseURL: "https://github.com/shruggietech/shruggie-brand/releases/download/v2.4.0/" + packageID + ".zip", SHA256: hex.EncodeToString(sum[:])}
	return pin, filename
}

func TestOpenValidatesArchiveBeforeUse(t *testing.T) {
	pin, filename := fixtureArchive(t, "", false, false)
	kit, err := Open(pin, filename)
	if err != nil {
		t.Fatal(err)
	}
	if len(kit.Manifest.Files) != 2 {
		t.Fatalf("manifest has %d files", len(kit.Manifest.Files))
	}
	badPin := pin
	badPin.SHA256 = strings.Repeat("0", 64)
	if _, err := Open(badPin, filename); err == nil || !strings.Contains(err.Error(), "SHA-256") {
		t.Fatalf("wrong digest error = %v", err)
	}
}

func TestOpenRejectsUnsafeOrChangedInventory(t *testing.T) {
	for _, test := range []struct {
		name        string
		extra       string
		badManifest bool
		symlink     bool
		want        string
	}{
		{"traversal", "../escape", false, false, "unsafe path"},
		{"case duplicate", "readme.md", false, false, "duplicate path"},
		{"symlink", "linked", false, true, "regular file"},
		{"undeclared", "surprise.txt", false, false, "undeclared file"},
		{"bad manifest hash", "", true, false, "wrong size or SHA-256"},
	} {
		t.Run(test.name, func(t *testing.T) {
			pin, filename := fixtureArchive(t, test.extra, test.badManifest, test.symlink)
			_, err := Open(pin, filename)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestVerifyChecksumsRequiresExactPin(t *testing.T) {
	pin, _ := fixtureArchive(t, "", false, false)
	filename := filepath.Join(t.TempDir(), "SHA256SUMS")
	if err := os.WriteFile(filename, []byte(pin.SHA256+"  ./"+pin.Package+".zip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := VerifyChecksums(pin, filename); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filename, []byte(strings.Repeat("0", 64)+"  ./"+pin.Package+".zip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := VerifyChecksums(pin, filename); err == nil {
		t.Fatal("accepted wrong formal release checksum")
	}
}

func TestRenderDocsManifestUsesRelativeIconPaths(t *testing.T) {
	kit := &Kit{Files: map[string][]byte{
		"favicons/site.webmanifest":          []byte(`{"name":"go-schedule","custom":{"retained":true},"icons":[{"src":"/maskable-icon-192x192.png","sizes":"192x192","type":"image/png","purpose":"maskable","extra":"retained"}]}`),
		"favicons/maskable-icon-192x192.png": []byte("png"),
	}}
	output, err := Render(kit, Mapping{Source: "favicons/site.webmanifest", Transform: "docs-site-manifest"})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(output, []byte(`"src": "maskable-icon-192x192.png"`)) {
		t.Fatalf("unexpected docs manifest: %s", output)
	}
	if !bytes.Contains(output, []byte(`"retained": true`)) || !bytes.Contains(output, []byte(`"extra": "retained"`)) {
		t.Fatalf("docs manifest dropped official extension fields: %s", output)
	}
}

func TestValidateConsumersRequiresDocsManifestIcons(t *testing.T) {
	root := t.TempDir()
	kit := &Kit{Pin: Pin{Package: "kit"}, Files: map[string][]byte{
		"favicons/site.webmanifest":          []byte(`{"icons":[{"src":"/maskable-icon-192x192.png"}]}`),
		"favicons/maskable-icon-192x192.png": []byte("png"),
	}}
	consumers := ConsumerMap{Version: 2, Mappings: []Mapping{{Source: "favicons/site.webmanifest", Targets: []string{"docs/assets/favicons/site.webmanifest"}, Purpose: "docs", Transform: "docs-site-manifest"}}}
	if err := ValidateConsumers(root, kit, consumers); err == nil || !strings.Contains(err.Error(), "no matching consumer") {
		t.Fatalf("missing icon consumer error = %v", err)
	}
	consumers.Mappings = append(consumers.Mappings, Mapping{Source: "favicons/maskable-icon-192x192.png", Targets: []string{"docs/assets/favicons/maskable-icon-192x192.png"}, Purpose: "docs"})
	if err := ValidateConsumers(root, kit, consumers); err != nil {
		t.Fatal(err)
	}
}

func TestValidateConsumersRejectsEscapesAndSymlinks(t *testing.T) {
	root := t.TempDir()
	kit := &Kit{Pin: Pin{Package: "kit"}, Files: map[string][]byte{"README.md": []byte("source")}}
	for _, target := range []string{"../escape", "docs/assets/brand/../../escape", "secrets/file.txt"} {
		mapping := ConsumerMap{Version: 2, Mappings: []Mapping{{Source: "README.md", Targets: []string{target}, Purpose: "test"}}}
		if err := ValidateConsumers(root, kit, mapping); err == nil {
			t.Errorf("accepted target %q", target)
		}
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(root, "docs")); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	mapping := ConsumerMap{Version: 2, Mappings: []Mapping{{Source: "README.md", Targets: []string{"docs/assets/brand/logo.svg"}, Purpose: "test"}}}
	if err := ValidateConsumers(root, kit, mapping); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("symlink error = %v", err)
	}
}
