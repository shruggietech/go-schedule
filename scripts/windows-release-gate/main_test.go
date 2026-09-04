package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/releasegate"
)

func TestRunRejectsUnknownCommandAsUsage(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	if code := run([]string{"unknown"}, &stdout, &stderr); code != 2 {
		t.Fatalf("run() = %d, want 2; stderr=%s", code, stderr.String())
	}
}

func TestRunVerifiesCandidateWithoutEvidence(t *testing.T) {
	t.Parallel()

	root := filepath.Join("..", "..", "test", "fixtures", "windows-release-gate", "passing")
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"verify-candidate",
		"--candidate-manifest", filepath.Join(root, "windows-candidate-manifest.json"),
		"--artifact", filepath.Join(root, "go-schedule_v1.0.0_windows_amd64.msi"),
		"--repository", "shruggietech/go-schedule",
		"--tag", "v1.0.0",
		"--commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() = %d, want 0; stderr=%s", code, stderr.String())
	}
}

func TestRunRejectsMissingFilesAsInternalFailure(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := run([]string{
		"validate",
		"--evidence", "missing.json",
		"--artifact", "missing.msi",
	}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("run() = %d, want 2; stderr=%s", code, stderr.String())
	}
}

func TestRunRendersFormalDispositionPacket(t *testing.T) {
	bundle, manifest, artifact := formalFixtureBundle(t, "attended-windows")
	output := filepath.Join(t.TempDir(), "dispositions")
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"render-dispositions",
		"--bundle", bundle,
		"--candidate-manifest", manifest,
		"--artifact", artifact,
		"--repository", "shruggietech/go-schedule",
		"--tag", "v1.0.0",
		"--commit", strings.Repeat("a", 40),
		"--output-dir", output,
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() = %d, want 0; stderr=%s", code, stderr.String())
	}
	entries, err := os.ReadDir(output)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 11 || !strings.Contains(stdout.String(), "10 issue records") {
		t.Fatalf("entries=%d stdout=%q", len(entries), stdout.String())
	}
}

func TestRunRenderDispositionsRequiresEveryOption(t *testing.T) {
	required := []string{"--bundle", "--candidate-manifest", "--artifact", "--repository", "--tag", "--commit", "--output-dir"}
	for _, omitted := range required {
		omitted := omitted
		t.Run(strings.TrimPrefix(omitted, "--"), func(t *testing.T) {
			args := []string{"render-dispositions"}
			for _, option := range required {
				if option != omitted {
					args = append(args, option, "value")
				}
			}
			var stdout, stderr bytes.Buffer
			if code := run(args, &stdout, &stderr); code != 2 {
				t.Fatalf("run() = %d, want 2; stderr=%s", code, stderr.String())
			}
			if !strings.Contains(stderr.String(), "requires") {
				t.Fatalf("stderr = %q, want required-option diagnostic", stderr.String())
			}
		})
	}
}

func TestRunRenderDispositionsRejectsNonFormalEvidenceAndIdentityDrift(t *testing.T) {
	tests := []struct {
		name       string
		class      string
		repository string
		tag        string
		commit     string
		want       string
	}{
		{"fixture class", "automated-fixture", "shruggietech/go-schedule", "v1.0.0", strings.Repeat("a", 40), "evidence_class"},
		{"repository", "attended-windows", "other/repo", "v1.0.0", strings.Repeat("a", 40), "repository"},
		{"tag", "attended-windows", "shruggietech/go-schedule", "v1.0.1", strings.Repeat("a", 40), "tag"},
		{"commit", "attended-windows", "shruggietech/go-schedule", "v1.0.0", strings.Repeat("b", 40), "commit"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bundle, manifest, artifact := formalFixtureBundle(t, tt.class)
			output := filepath.Join(t.TempDir(), "dispositions")
			var stdout, stderr bytes.Buffer
			code := run([]string{
				"render-dispositions", "--bundle", bundle,
				"--candidate-manifest", manifest, "--artifact", artifact,
				"--repository", tt.repository, "--tag", tt.tag,
				"--commit", tt.commit, "--output-dir", output,
			}, &stdout, &stderr)
			if code != 1 || !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("run()=%d stderr=%q, want validation %q", code, stderr.String(), tt.want)
			}
			if _, err := os.Lstat(output); !os.IsNotExist(err) {
				t.Fatalf("output exists after validation failure: %v", err)
			}
		})
	}
}

func TestRunRenderDispositionsRefusesExistingDestination(t *testing.T) {
	bundle, manifest, artifact := formalFixtureBundle(t, "attended-windows")
	output := filepath.Join(t.TempDir(), "dispositions")
	if err := os.Mkdir(output, 0o700); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(output, "keep.txt")
	if err := os.WriteFile(sentinel, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"render-dispositions", "--bundle", bundle,
		"--candidate-manifest", manifest, "--artifact", artifact,
		"--repository", "shruggietech/go-schedule", "--tag", "v1.0.0",
		"--commit", strings.Repeat("a", 40), "--output-dir", output,
	}, &stdout, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), "already exists") {
		t.Fatalf("run()=%d stderr=%q", code, stderr.String())
	}
	if data, err := os.ReadFile(sentinel); err != nil || string(data) != "keep" {
		t.Fatalf("sentinel changed: data=%q err=%v", data, err)
	}
}

func formalFixtureBundle(t *testing.T, evidenceClass string) (string, string, string) {
	t.Helper()
	root := filepath.Join("..", "..", "test", "fixtures", "windows-release-gate", "passing")
	evidenceFile, err := os.Open(filepath.Join(root, "evidence.json"))
	if err != nil {
		t.Fatal(err)
	}
	evidence, err := releasegate.DecodeEvidence(evidenceFile)
	closeErr := evidenceFile.Close()
	if err != nil || closeErr != nil {
		t.Fatalf("decode fixture: %v; close: %v", err, closeErr)
	}
	evidence.EvidenceClass = evidenceClass
	evidenceJSON, err := json.MarshalIndent(evidence, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	evidenceJSON = append(evidenceJSON, '\n')

	bundle := filepath.Join(t.TempDir(), "evidence.zip")
	file, err := os.Create(bundle)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(file)
	writeZipEntry(t, writer, "evidence.json", evidenceJSON)
	for _, attachment := range evidence.Attachments {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(attachment.Path)))
		if err != nil {
			t.Fatal(err)
		}
		writeZipEntry(t, writer, attachment.Path, data)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	return bundle, filepath.Join(root, "windows-candidate-manifest.json"), filepath.Join(root, "go-schedule_v1.0.0_windows_amd64.msi")
}

func writeZipEntry(t *testing.T, writer *zip.Writer, name string, data []byte) {
	t.Helper()
	entry, err := writer.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(entry, bytes.NewReader(data)); err != nil {
		t.Fatal(err)
	}
}
