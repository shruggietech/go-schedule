package integration_test

import (
	"strings"
	"testing"
)

func TestWaiverPromotionBindsTrustedHelperBeforeTagCheckout(t *testing.T) {
	t.Parallel()
	workflow := releaseGateFile(t, ".github", "workflows", "promote-release-waived.yml")
	trusted := strings.Index(workflow, "ref: $"+"{{ github.sha }}")
	build := strings.Index(workflow, "go build -o \"$RUNNER_TEMP/release-promotion\" ./scripts/release-promotion")
	tagged := strings.Index(workflow, "ref: refs/tags/")
	if trusted < 0 || build <= trusted || tagged <= build {
		t.Fatalf("trusted helper must be built before tagged checkout: trusted=%d build=%d tag=%d", trusted, build, tagged)
	}
	for _, forbidden := range []string{"/releases/tags/", "gh release download", "gh release upload", "--clobber"} {
		if strings.Contains(workflow, forbidden) {
			t.Fatalf("standing-waiver workflow independently resolves or overwrites release: %s", forbidden)
		}
	}
	comparison := strings.Index(workflow, "go run ./scripts/comparison-check -tag")
	ensure := strings.Index(workflow, "\"$RUNNER_TEMP/release-promotion\" ensure-checksum")
	tagCheck := strings.Index(workflow, "test \"$REMOTE_TAG\" = \"$TAG_COMMIT\"")
	publish := strings.Index(workflow, "\"$RUNNER_TEMP/release-promotion\" publish")
	if comparison < 0 || ensure <= comparison || tagCheck <= ensure || publish <= tagCheck {
		t.Fatalf("unsafe mutation ordering comparison=%d checksum=%d tag=%d publish=%d", comparison, ensure, tagCheck, publish)
	}
	for _, required := range []string{"actions: read", "contents: write", "--state \"$RUNNER_TEMP/release-selection.json\"", "windows-release-gate verify-candidate", "expected-assets.txt", "/attempts/$" + "{RUN_ATTEMPT}/jobs?per_page=100", "Build & stage GUI (windows-latest)", "sha256sum -c SHA256SUMS.txt", "test \"$(wc -l < SHA256SUMS.txt)\" -eq 8", "test \"$(find . -maxdepth 1 -type f | wc -l)\" -eq 9"} {
		if !strings.Contains(workflow, required) {
			t.Fatalf("standing-waiver workflow lost %q", required)
		}
	}
}
