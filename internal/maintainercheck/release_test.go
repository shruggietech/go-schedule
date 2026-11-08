package maintainercheck

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseNotesLiteralRoot(t *testing.T) {
	for _, name := range []string{"checkout[1]", "checkout[unclosed"} {
		t.Run(name, func(t *testing.T) {
			root := filepath.Join(t.TempDir(), name)
			put(t, root, ".github/release-notes/v1.2.3.md", "## Highlights\n\n- One.\n- Two.\n- Three.\n- Four.\n\nRead the [full changelog](https://github.com/shruggietech/go-schedule/blob/v1.2.3/CHANGELOG.md#123) for every change.\n")
			c := newChecker(root)
			c.releaseNotes()
			if len(c.findings) != 0 {
				t.Fatal(c.findings)
			}
		})
	}
}

func TestComparisonOrderAndDependencies(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		pass       bool
	}{
		{"same block before", "jobs:\n  publish:\n    steps:\n      - run: |\n          go run ./scripts/comparison-check -tag x -commit \"$TAG_COMMIT\"\n          gh release upload x y\n", true},
		{"continued bypass", "jobs:\n  publish:\n    steps:\n      - run: |\n          go run ./scripts/comparison-check -tag x -commit \"$TAG_COMMIT\" \\\n            || true\n          gh release upload x y\n", false},
		{"same block after", "jobs:\n  publish:\n    steps:\n      - run: |\n          gh release upload x y\n          go run ./scripts/comparison-check -tag x -commit \"$TAG_COMMIT\"\n", false},
		{"unrelated job", "jobs:\n  review:\n    steps:\n      - run: go run ./scripts/comparison-check -tag x -commit \"$TAG_COMMIT\"\n  publish:\n    steps:\n      - run: gh release upload x y\n", false},
		{"dependency", "jobs:\n  review:\n    steps:\n      - run: go run ./scripts/comparison-check -tag x -commit \"$TAG_COMMIT\"\n  publish:\n    needs: review\n    steps:\n      - run: gh release upload x y\n", true},
		{"conditional", "jobs:\n  publish:\n    steps:\n      - run: go run ./scripts/comparison-check -tag x -commit \"$TAG_COMMIT\"\n        if: false\n      - run: gh release upload x y\n", false},
		{"wrong source", "jobs:\n  publish:\n    steps:\n      - run: go run ./scripts/comparison-check -tag x -commit main\n      - run: gh release upload x y\n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			file := ".github/workflows/promote-release.yml"
			put(t, root, file, tc.body)
			c := newChecker(root)
			c.comparison(file, c.yaml(file))
			if (len(c.findings) == 0) != tc.pass {
				t.Fatal(c.findings)
			}
		})
	}
}

func TestChildFailureIsRetained(t *testing.T) {
	c := newChecker(t.TempDir())
	c.child("helper", "executable-that-does-not-exist-s136")
	if !strings.Contains(strings.Join(c.findings, "\n"), "child failed") {
		t.Fatal(c.findings)
	}
}

func BenchmarkWorkflowEvaluation(b *testing.B) {
	root := b.TempDir()
	put(b, root, "workflow.yml", "jobs:\n  test:\n    steps:\n      - uses: actions/checkout@v7\n      - uses: actions/setup-go@v7\n")
	b.ResetTimer()
	for b.Loop() {
		c := newChecker(root)
		c.actions("workflow.yml", c.yaml("workflow.yml"))
		if len(c.findings) > 0 {
			b.Fatal(c.findings)
		}
	}
}
