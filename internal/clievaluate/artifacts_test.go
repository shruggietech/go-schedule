package clievaluate

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

func TestArtifactReplaceSameSanitizedBytesKeepsFileIdentity(t *testing.T) {
	dir := t.TempDir()
	r := NewRedactor()
	s, err := OpenArtifacts(dir, 4096, r)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			t.Error(err)
		}
	}()
	original := []byte("neutral evidence with later-secret-value")
	if err := s.Replace("checkpoint.json", original); err != nil {
		t.Fatal(err)
	}
	checkpoint := filepath.Join(dir, "checkpoint.json")
	oldTime := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	if err := os.Chtimes(checkpoint, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(filepath.Join(dir, "checkpoint.json"))
	if err != nil {
		t.Fatal(err)
	}
	used := s.used
	if err := s.Replace("checkpoint.json", original); err != nil {
		t.Fatal(err)
	}
	after, err := os.Stat(filepath.Join(dir, "checkpoint.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(before, after) || !before.ModTime().Equal(after.ModTime()) {
		t.Fatal("unchanged sanitized bytes were physically replaced")
	}
	if s.used != used {
		t.Fatal("unchanged replacement changed accounting")
	}
	r.Register("later-secret-value")
	if err := s.Replace("checkpoint.json", original); err != nil {
		t.Fatal(err)
	}
	clean, err := os.ReadFile(filepath.Join(dir, "checkpoint.json"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(clean), "later-secret-value") || !strings.Contains(string(clean), "[REDACTED]") {
		t.Fatal("new redaction was skipped")
	}
	if s.used != int64(len(clean)) {
		t.Fatal("redaction replacement accounting differs")
	}
}

func BenchmarkArtifactReplaceUnchanged(b *testing.B) {
	s, err := OpenArtifacts(b.TempDir(), 4096, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			b.Error(err)
		}
	}()
	value := bytes.Repeat([]byte("e"), 1024)
	if err := s.Replace("checkpoint", value); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if err := s.Replace("checkpoint", value); err != nil {
			b.Fatal(err)
		}
	}
}

func TestArtifactFinalizationCancellationPreservesPriorEvidence(t *testing.T) {
	dir := t.TempDir()
	s, err := OpenArtifacts(dir, 1024, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.Replace("prior", []byte("retained")); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s.ctx = ctx
	if err := s.Replace("prior", []byte("overwrite")); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled replace: %v", err)
	}
	if _, err := s.readRegular("prior", 1024); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled read: %v", err)
	}
	s.ctx = nil
	data, err := s.readRegular("prior", 1024)
	if err != nil || string(data) != "retained" {
		t.Fatalf("prior evidence changed: %q %v", data, err)
	}
	if _, err := openArtifactsContext(ctx, dir, 1024, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled traversal: %v", err)
	}
}
