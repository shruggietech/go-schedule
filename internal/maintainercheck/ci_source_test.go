package maintainercheck

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestHostedCheckoutUsesRecordedCandidate(t *testing.T) {
	const selected = "${{ github.event.pull_request.head.sha || github.sha }}"
	c := newChecker(filepath.Join("..", ".."))
	for _, file := range []string{".github/workflows/ci.yml", ".github/workflows/codeql.yml"} {
		n := c.yaml(file)
		if len(c.findings) != 0 {
			t.Fatal(c.findings)
		}
		count := 0
		for _, step := range steps(n) {
			if value(field(step, "uses")) == "actions/checkout@v7" {
				count++
				if value(field(step, "with", "ref")) != selected {
					t.Fatalf("%s checkout must select the candidate head", file)
				}
			}
		}
		if count == 0 {
			t.Fatalf("%s has no checkout", file)
		}
	}
	ci := c.yaml(".github/workflows/ci.yml")
	if value(field(ci, "env", "GOSCHED_CI_SOURCE_SHA")) != selected {
		t.Fatal("CI artifact source must bind the selected checkout")
	}
	text := c.read(".github/workflows/ci.yml")
	if strings.Contains(text, "CI build from commit $env:GITHUB_SHA") || strings.Count(text, "CI build from commit $env:GOSCHED_CI_SOURCE_SHA") != 2 {
		t.Fatal("qualification origin must name the actual candidate checkout, not the event merge SHA")
	}
}
