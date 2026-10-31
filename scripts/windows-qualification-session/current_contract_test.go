package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCurrentSessionRequiresNoHistoricalInstaller(t *testing.T) {
	m, root := fixture(t)
	var inputs []input
	for _, in := range m.Inputs {
		if in.Role != "baseline" {
			inputs = append(inputs, in)
		}
	}
	m.Inputs = inputs
	out := filepath.Join(root, "current-session")
	if err := prepare(m, out, root); err != nil {
		t.Fatalf("current session needs historical MSI: %v", err)
	}
	if _, err := os.Stat(filepath.Join(out, "fresh.wsb")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(out, "upgrade.wsb")); !os.IsNotExist(err) {
		t.Fatal("historical upgrade session generated")
	}
	if _, ok := filenames["baseline"]; ok {
		t.Fatal("historical baseline role still supported")
	}
}

func TestCurrentSessionRejectsHistoricalInstallerRole(t *testing.T) {
	m, root := fixture(t)
	m.Inputs[1].Role = "baseline"
	if err := prepare(m, filepath.Join(root, "historical-session"), root); err == nil {
		t.Fatal("historical MSI input accepted")
	}
}
