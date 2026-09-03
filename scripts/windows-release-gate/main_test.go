package main

import (
	"bytes"
	"path/filepath"
	"testing"
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
