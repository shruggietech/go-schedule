package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/brandrelease"
)

func TestInvalidArchiveLeavesInstalledKitUntouched(t *testing.T) {
	root := t.TempDir()
	brandDir := filepath.Join(root, "brand")
	if err := os.MkdirAll(brandDir, 0o755); err != nil {
		t.Fatal(err)
	}
	pin := brandrelease.Pin{Package: "go-schedule-brand-2.0.0-bb2.4.0", BrandVersion: "2.0.0", BrandBuilderVersion: "2.4.0", ReleaseURL: "https://github.com/shruggietech/shruggie-brand/releases/download/v2.4.0/go-schedule-brand-2.0.0-bb2.4.0.zip", SHA256: strings.Repeat("0", 64)}
	pinData, err := json.Marshal(pin)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(brandDir, "source.json"), pinData, 0o644); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(brandDir, "README.md")
	if err := os.WriteFile(sentinel, []byte("installed identity"), 0o644); err != nil {
		t.Fatal(err)
	}
	archive := filepath.Join(t.TempDir(), pin.Package+".zip")
	if err := os.WriteFile(archive, []byte("invalid archive"), 0o644); err != nil {
		t.Fatal(err)
	}
	checksums := filepath.Join(t.TempDir(), "SHA256SUMS")
	if err := os.WriteFile(checksums, []byte(pin.SHA256+"  ./"+pin.Package+".zip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := run(root, archive, checksums, "unused.json"); err == nil {
		t.Fatal("accepted invalid archive")
	}
	got, err := os.ReadFile(sentinel)
	if err != nil || string(got) != "installed identity" {
		t.Fatalf("installed kit changed: %q, %v", got, err)
	}
}
