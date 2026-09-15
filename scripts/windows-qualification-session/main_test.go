package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixture(t *testing.T) (manifest, string) {
	t.Helper()
	root := canonicalTempDir(t)
	m := manifest{Schema: 1, Repository: "shruggietech/go-schedule", Tag: "v1.4.0", Commit: strings.Repeat("a", 40), RunID: 123, RunAttempt: 1}
	for _, name := range []string{"candidate", "baseline", "powershell", "webview2"} {
		data := []byte(name)
		path := filepath.Join(root, name)
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		h := sha256.Sum256(data)
		sources := map[string]string{"candidate": "https://github.com/shruggietech/go-schedule/actions/runs/123", "baseline": "https://github.com/shruggietech/go-schedule/releases/download/v1.1.1/go-schedule_v1.1.1_windows_amd64.msi", "powershell": "https://github.com/PowerShell/PowerShell/releases/download/v7.5.0/PowerShell-7.5.0-win-x64.zip", "webview2": "https://developer.microsoft.com/en-us/microsoft-edge/webview2/"}
		m.Inputs = append(m.Inputs, input{Role: name, Path: path, Bytes: int64(len(data)), SHA256: hex.EncodeToString(h[:]), Source: sources[name]})
	}
	for _, role := range []string{"bootstrap", "collector"} {
		if err := os.WriteFile(filepath.Join(root, filenames[role]), []byte(role+" fixture"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	return m, root
}

func canonicalTempDir(t *testing.T) string {
	t.Helper()
	path, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func TestPrepareAllowsSiblingOutput(t *testing.T) {
	m, root := fixture(t)
	if err := prepare(m, filepath.Join(root, "new-session"), root); err != nil {
		t.Fatal(err)
	}
}

func TestPrepareRejectsInputAliases(t *testing.T) {
	for _, hardlink := range []bool{false, true} {
		m, root := fixture(t)
		alias := m.Inputs[0].Path
		if hardlink {
			alias = filepath.Join(root, "candidate-hardlink")
			if err := os.Link(m.Inputs[0].Path, alias); err != nil {
				t.Fatal(err)
			}
		}
		m.Inputs[1].Path = alias
		m.Inputs[1].Bytes = m.Inputs[0].Bytes
		m.Inputs[1].SHA256 = m.Inputs[0].SHA256
		if err := prepare(m, filepath.Join(root, "new-session"), root); err == nil || !strings.Contains(err.Error(), "alias") {
			t.Fatalf("input alias accepted: %v", err)
		}
	}
}

func TestPrepareValidation(t *testing.T) {
	for _, mode := range []string{"tamper", "size", "duplicate", "commit", "occupied", "overlap"} {
		t.Run(mode, func(t *testing.T) {
			m, root := fixture(t)
			out := filepath.Join(canonicalTempDir(t), "new")
			switch mode {
			case "tamper":
				m.Inputs[0].SHA256 = strings.Repeat("0", 64)
			case "size":
				m.Inputs[0].Bytes++
			case "duplicate":
				m.Inputs[1].Role = "candidate"
			case "commit":
				m.Commit = "main"
			case "occupied":
				if err := os.Mkdir(out, 0700); err != nil {
					t.Fatal(err)
				}
			case "overlap":
				out = root
			}
			if err := prepare(m, out, root); err == nil {
				t.Fatal("invalid preparation accepted")
			}
		})
	}
}

func TestConfigurationEscapesHostPaths(t *testing.T) {
	s := configuration(`C:\inputs & [x]`, `C:\outputs <x>`, "fresh")
	var v struct{ XMLName xml.Name }
	if err := xml.Unmarshal([]byte(s), &v); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(s, "&amp;") || !strings.Contains(s, "&lt;") {
		t.Fatal("host paths not escaped")
	}
	if !strings.Contains(s, "<ReadOnly>true</ReadOnly>") || !strings.Contains(s, "-WindowStyle Hidden") {
		t.Fatal("unsafe configuration")
	}
}

func TestPreparationExportsNoAttestation(t *testing.T) {
	m, root := fixture(t)
	if err := os.WriteFile(filepath.Join(root, "Start-QualificationGuest.ps1"), []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "Invoke-ReleaseCandidateAttended.ps1"), []byte("collector fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(canonicalTempDir(t), "package")
	if err := prepare(m, out, root); err != nil {
		t.Fatal(err)
	}
	for _, scenario := range []string{"fresh", "upgrade"} {
		if _, err := os.Stat(filepath.Join(out, scenario+".wsb")); err != nil {
			t.Fatal(err)
		}
		entries, err := os.ReadDir(filepath.Join(out, scenario+"-exports"))
		if err != nil || len(entries) != 0 {
			t.Fatal("export must start empty", err)
		}
	}
	b, err := os.ReadFile(filepath.Join(out, "inputs", "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var got manifest
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Inputs) != 6 || strings.Contains(string(b), "attested_at") {
		t.Fatal("bad packaged manifest")
	}
}
