package clievaluate

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestSourceContentChangesRejectAnAlreadyDirtyResume(t *testing.T) {
	root := t.TempDir()
	paths := []string{"tracked.go", "untracked.go"}
	for _, name := range paths {
		if err := os.WriteFile(filepath.Join(root, name), []byte("initial"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	initial, err := digestSourceFiles(context.Background(), root, paths)
	if err != nil {
		t.Fatal(err)
	}
	prior := captureIdentity{Subject: Subject{Commit: "same-commit", Dirty: true, SourceSHA256: initial}}
	for _, name := range paths {
		if err := os.WriteFile(filepath.Join(root, name), []byte("changed"), 0600); err != nil {
			t.Fatal(err)
		}
		digest, err := digestSourceFiles(context.Background(), root, paths)
		if err != nil {
			t.Fatal(err)
		}
		current := prior
		current.Subject.SourceSHA256 = digest
		if err := compareResumeInputs(prior, current); err == nil {
			t.Fatal("changed bytes in an already dirty checkout were accepted")
		}
		if err := os.WriteFile(filepath.Join(root, name), []byte("initial"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Remove(filepath.Join(root, paths[0])); err != nil {
		t.Fatal(err)
	}
	deleted, err := digestSourceFiles(context.Background(), root, paths)
	if err != nil || deleted == initial {
		t.Fatal("deleted tracked path lost its source identity", err)
	}
	if err := os.WriteFile(filepath.Join(root, "ignored-runtime"), []byte("changed ignored data"), 0600); err != nil {
		t.Fatal(err)
	}
	again, err := digestSourceFiles(context.Background(), root, []string{paths[1], paths[0], paths[1]})
	if err != nil || again != deleted {
		t.Fatal("path order, duplicate entries or an unlisted ignored runtime changed source identity", err)
	}
}
