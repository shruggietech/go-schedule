package logbus

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRotatingWriter_RotatesAndBounds(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "logs", "goschedule.log")
	// Tiny cap so each ~10-byte write rotates; keep 2 files.
	w, err := NewRotatingWriter(path, 10, 2)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	for i := 0; i < 5; i++ {
		if _, err := w.Write([]byte("0123456789\n")); err != nil { // 11 bytes > cap
			t.Fatal(err)
		}
	}
	_ = w.Close()

	// With maxFiles=2 we keep at most the active file + 1 backup.
	entries, _ := os.ReadDir(filepath.Dir(path))
	count := 0
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "goschedule.log") {
			count++
		}
	}
	if count > 2 {
		t.Fatalf("retained %d log files, want <= 2 (bounded retention)", count)
	}
}

func TestRotatingWriter_CreatesDirAndAppends(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "goschedule.log")
	w, err := NewRotatingWriter(path, 1<<20, 5)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = w.Write([]byte("hello\n"))
	_ = w.Close()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if !strings.Contains(string(b), "hello") {
		t.Fatalf("log content = %q", string(b))
	}
}
