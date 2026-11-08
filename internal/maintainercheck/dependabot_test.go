package maintainercheck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDependencyExactValues(t *testing.T) {
	original, err := os.ReadFile(filepath.Join("..", "..", ".github", "dependabot.yml"))
	if err != nil {
		t.Fatal(err)
	}
	baseline := t.TempDir()
	put(t, baseline, ".github/dependabot.yml", string(original))
	valid := newChecker(baseline)
	valid.dependabot()
	if len(valid.findings) != 0 {
		t.Fatal(valid.findings)
	}
	for _, tc := range []struct{ name, old, replacement string }{{"label", "- dependencies", "- non-dependencies"}, {"minor", "- minor", "- not-minor"}, {"storage", "- modernc.org/sqlite", "- modernc.org/sqlite-suffix"}} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			put(t, root, ".github/dependabot.yml", strings.ReplaceAll(string(original), tc.old, tc.replacement))
			c := newChecker(root)
			c.dependabot()
			if len(c.findings) == 0 {
				t.Fatal("invalid exact policy value was accepted")
			}
		})
	}
}
