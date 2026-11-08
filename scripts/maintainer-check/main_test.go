package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunUsageAndLifecycleRoot(t *testing.T) {
	for _, args := range [][]string{nil, {"unknown"}, {"lifecycle", ".", "extra"}} {
		var out, diagnostic bytes.Buffer
		if code := run(args, &out, &diagnostic); code != 2 || !strings.Contains(diagnostic.String(), "usage:") {
			t.Fatal(code, diagnostic.String())
		}
	}
	root := filepath.Join(t.TempDir(), "root with spaces")
	if err := os.MkdirAll(filepath.Join(root, "specs", "001-example"), 0o700); err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string]string{
		"specs/README.md":            "| 001-example | Implemented | record |\n",
		"specs/001-example/spec.md":  "**Status**: Implemented\n**Delivery**: commit abc1234\n",
		"specs/001-example/tasks.md": "- [x] T001 Done\n",
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	var out, diagnostic bytes.Buffer
	if code := run([]string{"lifecycle", root}, &out, &diagnostic); code != 0 || !strings.Contains(out.String(), "1 specification(s)") {
		t.Fatal(code, out.String(), diagnostic.String())
	}
	out.Reset()
	diagnostic.Reset()
	if code := run([]string{"lifecycle", filepath.Join(root, "missing")}, &out, &diagnostic); code != 1 || !strings.Contains(diagnostic.String(), "specs") {
		t.Fatal(code, diagnostic.String())
	}
}
