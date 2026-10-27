package incident

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

func TestIncidentArtifactsAreExclusiveBoundedAndAttributable(t *testing.T) {
	parent := canonicalTempDir(t)
	capture, err := New(parent)
	if err != nil {
		t.Fatal(err)
	}
	if err := capture.Write("sample.txt", []byte("owned")); err != nil {
		t.Fatal(err)
	}
	if err := capture.Write("sample.txt", []byte("replacement")); err == nil {
		t.Fatal("overwrote original artifact")
	}
	if err := capture.Publish("published.json", []byte(`{"complete":true}`)); err != nil {
		t.Fatal(err)
	}
	if err := capture.Publish("published.json", []byte(`{"replaced":true}`)); err == nil {
		t.Fatal("replaced published artifact")
	}
	artifact := capture.Inspect("sample.txt", 5)
	if artifact.Status != "collected" || artifact.Bytes != 5 || len(artifact.SHA256) != 64 {
		t.Fatalf("artifact: %+v", artifact)
	}
	if got := capture.Inspect("sample.txt", 4); got.Status != "oversized" {
		t.Fatalf("size limit: %+v", got)
	}
	if err := capture.DiscardOversized("sample.txt", 4); err != nil {
		t.Fatal(err)
	}
	if got := capture.Inspect("sample.txt", 5); got.Status != "absent" {
		t.Fatal("retained oversized producer output")
	}
	if got := capture.Inspect("missing.txt", 10); got.Status != "absent" {
		t.Fatalf("missing: %+v", got)
	}
	if err := capture.Write("../outside", []byte("no")); err == nil {
		t.Fatal("accepted traversal")
	}
	if got := capture.Inspect("../outside", 10); got.Status != "unavailable" {
		t.Fatalf("escaped artifact: %+v", got)
	}
	record := Record{Kind: "go-schedule.incident", Version: 1, ID: capture.ID, Category: "script", Phase: "command", Target: Target{PID: 123}, Original: Original{NativeExit: 37}, Collection: Collection{Status: "absent"}, Artifacts: []Artifact{artifact}}
	if err := capture.Save(record); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(capture.Dir, "incident.json"))
	if err != nil {
		t.Fatal(err)
	}
	var decoded Record
	if err := json.Unmarshal(data, &decoded); err != nil || decoded.Original.NativeExit != 37 || decoded.Target.PID != 123 {
		t.Fatalf("independent original outcome: %s %v", data, err)
	}
}

func TestIncidentRejectsLinkedParentAndArtifact(t *testing.T) {
	parent, outside := canonicalTempDir(t), canonicalTempDir(t)
	link := filepath.Join(parent, "linked")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("native link capability unavailable: %v", err)
	}
	if _, err := New(link); err == nil {
		t.Fatal("accepted linked parent")
	}
	capture, err := New(parent)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outside, "private"), []byte("not owned"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "private"), filepath.Join(capture.Dir, "sample")); err != nil {
		t.Fatal(err)
	}
	if a := capture.Inspect("sample", 100); a.Status != "unavailable" {
		t.Fatalf("followed linked artifact: %+v", a)
	}
}

func TestConcurrentIncidentsHaveDistinctOwnedDirectories(t *testing.T) {
	parent := canonicalTempDir(t)
	var wg sync.WaitGroup
	ids := make(chan string, 6)
	for range 6 {
		wg.Go(func() {
			c, err := New(parent)
			if err != nil {
				t.Error(err)
				return
			}
			ids <- c.ID
		})
	}
	wg.Wait()
	close(ids)
	seen := map[string]bool{}
	for id := range ids {
		if seen[id] {
			t.Fatal("incident identity collision")
		}
		seen[id] = true
	}
	if len(seen) != 6 {
		t.Fatalf("created %d incidents", len(seen))
	}
}

func TestIncidentRejectsChangedDirectoryIdentity(t *testing.T) {
	capture, err := New(canonicalTempDir(t))
	if err != nil {
		t.Fatal(err)
	}
	other := filepath.Join(canonicalTempDir(t), filepath.Base(capture.Dir))
	if err := os.Mkdir(other, 0700); err != nil {
		t.Fatal(err)
	}
	capture.Dir = other
	if err := capture.Write("unexpected.txt", []byte("not the original capture")); err == nil {
		t.Fatal("accepted a changed capture directory")
	}
	if _, err := capture.CoreCLREnvironment(os.Environ()); err == nil {
		t.Fatal("configured crash output outside original capture identity")
	}
	if _, err := os.Stat(filepath.Join(other, "unexpected.txt")); !os.IsNotExist(err) {
		t.Fatal("wrote outside the original capture identity")
	}
}

func canonicalTempDir(t *testing.T) string {
	t.Helper()
	parent, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return parent
}

func TestIncidentTempFixtureResolvesPlatformAlias(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX temporary-root alias")
	}
	root := canonicalTempDir(t)
	actual, alias := filepath.Join(root, "actual"), filepath.Join(root, "system-temp")
	if err := os.Mkdir(actual, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(actual, alias); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TMPDIR", alias)
	t.Run("owned_fixture", func(t *testing.T) {
		parent := canonicalTempDir(t)
		capture, err := New(parent)
		if err != nil {
			t.Fatal(err)
		}
		if err := capture.Write("owned.txt", []byte("owned")); err != nil {
			t.Fatal(err)
		}
		callerLink := filepath.Join(parent, "caller-link")
		if err := os.Symlink(actual, callerLink); err != nil {
			t.Fatal(err)
		}
		if _, err := New(callerLink); err == nil {
			t.Fatal("accepted caller-created link")
		}
	})
}
