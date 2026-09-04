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
		"/attempts/${RUN_ATTEMPT}/jobs?per_page=100",
		"Build & stage GUI (windows-latest)",
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
		"evidence_class = 'attended-windows'",
		"attachments/tasks",
	} {
		if !strings.Contains(script, required) {
			t.Fatalf("attended collector missing %q", required)
		}
	}
	for _, metric := range []string{"prior_config_available", "prior_logs_available"} {
		if strings.Count(script, metric) != 2 {
			t.Fatalf("attended collector must template %q for both reinstall outcomes", metric)
		}
	}
	for _, required := range []string{
		"New-ScenarioTemplate",
		"$scenarioId.StartsWith('desktop.')",
		"minimum_text_contrast",
		"touchpad_unavailable_reason",
		"schedule_row_count",
		"activity_row_count",
		"horizontal_scrollbar_present",
	} {
		if !strings.Contains(script, required) {
			t.Fatalf("S047 collector template contract missing %q", required)
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

func TestWindowsReleaseGateRendersOfflineFailClosedDispositionPacket(t *testing.T) {
	t.Parallel()

	command := releaseGateFile(t, "scripts", "windows-release-gate", "main.go")
	for _, required := range []string{
		`case "render-dispositions":`,
		`releasegate.Validate(evidence, root, *artifactPath, *options)`,
		`releasegate.ValidateBundleContents(root, evidence)`,
		`releasegate.ValidateCandidateManifest(evidence.Candidate, manifest)`,
		`releasegate.WriteDispositionPacket(*outputDir, evidence)`,
	} {
		if !strings.Contains(command, required) {
			t.Fatalf("render-dispositions command missing %q", required)
		}
	}

	renderer := releaseGateFile(t, "internal", "releasegate", "disposition.go")
	for _, required := range []string{
		`Issue: 96`, `Issue: 98`, `Issue: 101`, `Issue: 104`, `Issue: 105`,
		`Issue: 106`, `Issue: 109`, `Issue: 111`, `Issue: 112`, `Issue: 113`,
		`fmt.Sprintf("issue-%03d.md"`, "packet.json",
		"os.MkdirTemp", "os.Rename", "os.Lstat",
	} {
		if !strings.Contains(renderer, required) {
			t.Fatalf("disposition renderer missing %q", required)
		}
	}
	for _, forbidden := range []string{"net/http", "api.github.com", "gh issue", "gh release"} {
		if strings.Contains(command, forbidden) || strings.Contains(renderer, forbidden) {
			t.Fatalf("offline disposition path contains forbidden %q", forbidden)
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
