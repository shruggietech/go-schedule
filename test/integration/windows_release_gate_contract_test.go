package integration_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/releasegate"
)

func TestReleaseWorkflowStagesEveryUploadAsDraft(t *testing.T) {
	t.Parallel()

	workflow := releaseGateFile(t, ".github", "workflows", "release.yml")
	uploads := strings.Count(workflow, "uses: softprops/action-gh-release@v3")
	if uploads == 0 || strings.Count(workflow, "draft: true") != uploads {
		t.Fatalf("release uploads=%d draft declarations=%d", uploads, strings.Count(workflow, "draft: true"))
	}
	for _, forbidden := range []string{"--draft=false", "SHA256SUMS.txt"} {
		if strings.Contains(workflow, forbidden) {
			t.Fatalf("release staging contains forbidden %q", forbidden)
		}
	}
	for _, required := range []string{
		"windows-candidate-manifest.json",
		"-RunId \"${GITHUB_RUN_ID}\"",
		"-RunAttempt \"${GITHUB_RUN_ATTEMPT}\"",
	} {
		if !strings.Contains(workflow, required) {
			t.Fatalf("release staging missing %q", required)
		}
	}
}

func TestPromotionOrdersExactGateChecksumsAndPublication(t *testing.T) {
	t.Parallel()

	workflow := releaseGateFile(t, ".github", "workflows", "promote-release.yml")
	gate := strings.Index(workflow, "windows-release-gate verify-bundle")
	checksums := strings.Index(workflow, "LC_ALL=C sha256sum *")
	promote := strings.Index(workflow, "-F draft=false")
	if gate < 0 || checksums <= gate || promote <= checksums {
		t.Fatalf("unsafe order gate=%d checksums=%d promote=%d", gate, checksums, promote)
	}
	for _, required := range []string{
		"actions: read",
		"jq -r '.draft'",
		"expected-assets.txt",
		"git ls-remote origin",
		"windows-release-gate verify-candidate",
		"windows-candidate-manifest.json",
		"windows-attended-evidence.zip",
		"--candidate-manifest \"$MANIFEST\"",
		".github/workflows/release.yml",
		".conclusion')\" = success",
	} {
		if !strings.Contains(workflow, required) {
			t.Fatalf("promotion missing %q", required)
		}
	}
	if strings.Contains(workflow, "go build") || strings.Contains(workflow, "wix build") {
		t.Fatal("promotion rebuilds a tested release artifact")
	}
}

func TestAttendedCollectorUsesCanonicalScenariosAndHiddenChildren(t *testing.T) {
	t.Parallel()

	script := releaseGateFile(t, "test", "windows", "Invoke-ReleaseCandidateAttended.ps1")
	for _, id := range releasegate.RequiredScenarioIDs() {
		if !strings.Contains(script, "'"+id+"'") {
			t.Fatalf("collector scenario %q is missing", id)
		}
	}
	for _, required := range []string{
		"UseShellExecute = $false",
		"CreateNoWindow = $true",
		"RedirectStandardOutput = $true",
		"RedirectStandardError = $true",
		"RedirectStandardInput = $true",
		"ReadToEndAsync()",
		"Expected exactly one visible top-level window",
		"OpenProcessToken",
		"ProcessUserSid",
		"Fullscreen",
		"FyneEvidencePath",
		"Refusing to overwrite",
	} {
		if !strings.Contains(script, required) {
			t.Fatalf("attended collector missing %q", required)
		}
	}
}

func TestMSIInspectorCanWriteExactCandidateManifest(t *testing.T) {
	t.Parallel()

	script := releaseGateFile(t, "test", "windows", "inspect-installer.ps1")
	for _, required := range []string{
		"CandidateManifestPath",
		"WHERE ``Property``='ProductCode'",
		"product_code = $productCode.ToUpperInvariant()",
		"run_attempt = $RunAttempt",
		"FileMode]::CreateNew",
		"UTF8Encoding]::new($false)",
	} {
		if !strings.Contains(script, required) {
			t.Fatalf("MSI inspector missing %q", required)
		}
	}
}

func releaseGateFile(t *testing.T, parts ...string) string {
	t.Helper()
	_, current, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(current), "..", ".."))
	data, err := os.ReadFile(filepath.Join(append([]string{root}, parts...)...))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
