package maintainercheck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func put(t testing.TB, root, name, data string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestActionFields(t *testing.T) {
	for _, tc := range []struct{ name, text, want string }{
		{"comment", "# uses: example/unknown@v1\njobs: {}\n", ""},
		{"script", "jobs:\n  test:\n    steps:\n      - run: 'echo uses: example/unknown@v1'\n", ""},
		{"step", "jobs:\n  test:\n    steps:\n      - uses: example/unknown@v1\n", "unaudited action reference example/unknown@v1"},
		{"reusable", "jobs:\n  test:\n    uses: example/unknown@v1\n", "unaudited action reference example/unknown@v1"},
		{"obsolete", "jobs:\n  test:\n    steps:\n      - uses: actions/checkout@v4\n", "expected actions/checkout@v7"},
		{"anchor", "approved: &approved actions/checkout@v7\njobs:\n  test:\n    steps:\n      - uses: *approved\n", ""},
		{"cycle", "jobs: &jobs {self: *jobs}\n", "cyclic YAML alias"},
		{"duplicate", "jobs:\n  test: {}\n  test: {}\n", "duplicate"},
		{"malformed", "jobs: [\n", "YAML"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			put(t, root, "workflow.yml", tc.text)
			c := newChecker(root)
			n := c.yaml("workflow.yml")
			c.actions("workflow.yml", n)
			got := strings.Join(c.findings, "\n")
			if tc.want == "" && got != "" {
				t.Fatal(got)
			}
			if tc.want != "" && !strings.Contains(got, tc.want) {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
		})
	}
}

func TestEvaluationObservesLaterMutation(t *testing.T) {
	root := t.TempDir()
	put(t, root, "workflow.yml", "jobs: {}\n")
	c := newChecker(root)
	c.actions("workflow.yml", c.yaml("workflow.yml"))
	if len(c.findings) != 0 {
		t.Fatal(c.findings)
	}
	put(t, root, "workflow.yml", "jobs:\n  test:\n    uses: example/unknown@v1\n")
	c = newChecker(root)
	c.actions("workflow.yml", c.yaml("workflow.yml"))
	if len(c.findings) == 0 {
		t.Fatal("stale cached evaluation")
	}
}
