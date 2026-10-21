package clievaluate

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArtifactRedactionPrecedesPersistenceAndHashing(t *testing.T) {
	r := NewRedactor("fixture-secret-split")
	store, err := OpenArtifacts(t.TempDir(), 1024, r)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Error(err)
		}
	}()
	var capture bytes.Buffer
	capture.WriteString("fixture-secret-")
	capture.WriteString("split")
	capture.WriteString(` {"token":"new-dynamic-secret"}`)
	a, err := store.Write("case/output.txt", "text/plain", capture.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	data, err := store.Read(a)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "fixture-secret") || strings.Contains(string(data), "new-dynamic-secret") {
		t.Fatal("secret persisted")
	}
	if _, err := store.Write(a.Path, "text/plain", nil); err == nil {
		t.Fatal("artifact overwritten")
	}
	if _, err := store.Write("huge", "text/plain", make([]byte, 1025)); err == nil {
		t.Fatal("artifact budget bypassed")
	}
	a.SHA256 = strings.Repeat("0", 64)
	if _, err := store.Read(a); err == nil {
		t.Fatal("tampered hash accepted")
	}
}

func TestArtifactRootRejectsSymlinkEscape(t *testing.T) {
	dir := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret"), []byte("outside"), 0600); err != nil {
		t.Fatal(err)
	}
	store, err := OpenArtifacts(dir, 1024, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Error(err)
		}
	}()
	if err := os.Symlink(outside, filepath.Join(dir, "link")); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Write("link/escape", "text/plain", []byte("no")); err == nil {
		t.Fatal("symlink write escaped root")
	}
}
