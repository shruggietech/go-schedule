package userprefs

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"
)

func testStore(t *testing.T) *Store {
	t.Helper()
	root := t.TempDir()
	return NewStore(Options{Path: filepath.Join(root, "current", "preferences.json"), LegacyPath: filepath.Join(root, "legacy.json")})
}

func TestLegacyMigrationAndCurrentAuthority(t *testing.T) {
	store := testStore(t)
	if err := os.WriteFile(store.legacyPath, []byte(`{"appearance.mode":"dark","appearance.font":"old"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	prefs, err := store.Load()
	if err != nil || prefs.Appearance != AppearanceDark || prefs.Transition.Status != "migrated" || len(prefs.Transition.Retired) != 2 || prefs.Popups.Enabled {
		t.Fatalf("prefs=%+v err=%v", prefs, err)
	}
	if err := os.WriteFile(store.legacyPath, []byte(`{"appearance.mode":"light"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	again, err := store.Load()
	if err != nil || !reflect.DeepEqual(prefs, again) {
		t.Fatalf("current authority changed: %+v %v", again, err)
	}
}

func TestInvalidPatchDoesNotInitializeOrReplacePreferences(t *testing.T) {
	store := testStore(t)
	invalid := Appearance("sepia")
	if _, err := store.Patch(Patch{Appearance: &invalid}); err == nil {
		t.Fatal("invalid appearance accepted")
	}
	if _, err := os.Stat(store.Path()); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("invalid patch initialized file: %v", err)
	}
	if _, err := store.Load(); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(store.Path())
	empty := []string{}
	if _, err := store.Patch(Patch{PopupConditions: &empty}); err == nil {
		t.Fatal("empty conditions accepted")
	}
	after, _ := os.ReadFile(store.Path())
	if string(before) != string(after) {
		t.Fatal("invalid patch changed current bytes")
	}
}

func TestIndependentStoresPreserveConcurrentScopedUpdates(t *testing.T) {
	store := testStore(t)
	if _, err := store.Load(); err != nil {
		t.Fatal(err)
	}
	other := NewStore(Options{Path: store.Path()})
	dark, enabled := AppearanceDark, true
	start := make(chan struct{})
	var wg sync.WaitGroup
	for _, action := range []func() error{
		func() error { _, err := store.Patch(Patch{Appearance: &dark}); return err },
		func() error { _, err := other.Patch(Patch{PopupEnabled: &enabled}); return err },
	} {
		wg.Add(1)
		go func(run func() error) {
			defer wg.Done()
			<-start
			if err := run(); err != nil {
				t.Error(err)
			}
		}(action)
	}
	close(start)
	wg.Wait()
	prefs, err := store.Load()
	if err != nil || prefs.Appearance != AppearanceDark || !prefs.Popups.Enabled {
		t.Fatalf("lost update: %+v %v", prefs, err)
	}
}

func TestSecondWriterReadsAfterFirstReplacementAndReadersNeedNoLock(t *testing.T) {
	store := testStore(t)
	if _, err := store.Load(); err != nil {
		t.Fatal(err)
	}
	entered, release := make(chan struct{}), make(chan struct{})
	first := NewStore(Options{Path: store.Path(), Rename: func(from, to string) error { close(entered); <-release; return os.Rename(from, to) }})
	dark, enabled := AppearanceDark, true
	firstResult, secondResult := make(chan error, 1), make(chan error, 1)
	go func() { _, err := first.Patch(Patch{Appearance: &dark}); firstResult <- err }()
	<-entered
	// The current file remains readable while an independent writer holds the mutation lock.
	if prefs, err := store.Load(); err != nil || prefs.Appearance != AppearanceSystem {
		close(release)
		t.Fatalf("read existing document under writer lock: %+v %v", prefs, err)
	}
	go func() { _, err := store.Patch(Patch{PopupEnabled: &enabled}); secondResult <- err }()
	select {
	case err := <-secondResult:
		close(release)
		t.Fatalf("second writer bypassed pending replacement: %v", err)
	case <-time.After(25 * time.Millisecond):
	}
	close(release)
	if err := <-firstResult; err != nil {
		t.Fatal(err)
	}
	if err := <-secondResult; err != nil {
		t.Fatal(err)
	}
	prefs, err := store.Load()
	if err != nil || prefs.Appearance != AppearanceDark || !prefs.Popups.Enabled {
		t.Fatalf("second writer used stale state: %+v %v", prefs, err)
	}
}

func TestRelativePathsNeverInitialize(t *testing.T) {
	store := NewStore(Options{Path: "relative/preferences.json"})
	if _, err := store.Load(); err == nil {
		t.Fatal("relative read accepted")
	}
	dark := AppearanceDark
	if _, err := store.Patch(Patch{Appearance: &dark}); err == nil {
		t.Fatal("relative mutation accepted")
	}
	if _, err := store.Reset("all"); err == nil {
		t.Fatal("relative reset accepted")
	}
}

func TestAbandonedLockFileDoesNotBlockFuturePreferenceWrites(t *testing.T) {
	store := testStore(t)
	if _, err := store.Load(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(store.Path()+".lock", []byte("abandoned marker"), 0o600); err != nil {
		t.Fatal(err)
	}
	dark := AppearanceDark
	prefs, err := store.Patch(Patch{Appearance: &dark})
	if err != nil || prefs.Appearance != AppearanceDark {
		t.Fatalf("abandoned lock marker prevented update: %+v %v", prefs, err)
	}
}

func TestScopedResetAndExplicitDaemonClear(t *testing.T) {
	store := testStore(t)
	dark, enabled, daemons := AppearanceDark, true, []string{"daemon-a"}
	if _, err := store.Patch(Patch{Appearance: &dark, PopupEnabled: &enabled, PopupDaemonIDs: &daemons}); err != nil {
		t.Fatal(err)
	}
	empty := []string{}
	prefs, err := store.Patch(Patch{PopupDaemonIDs: &empty})
	if err != nil || prefs.Popups.DaemonIDs == nil || len(prefs.Popups.DaemonIDs) != 0 || !prefs.Popups.Enabled {
		t.Fatalf("clear=%+v %v", prefs, err)
	}
	prefs, err = store.Reset("appearance")
	if err != nil || prefs.Appearance != AppearanceSystem || !prefs.Popups.Enabled {
		t.Fatalf("reset appearance=%+v %v", prefs, err)
	}
	if _, err := store.Patch(Patch{Appearance: &dark}); err != nil {
		t.Fatal(err)
	}
	prefs, err = store.Reset("popups")
	if err != nil || prefs.Appearance != AppearanceDark || prefs.Popups.Enabled {
		t.Fatalf("reset popups=%+v %v", prefs, err)
	}
	if _, err := store.Reset("unknown"); err == nil {
		t.Fatal("unknown reset scope accepted")
	}
}

func TestFailedReplacementPreservesCurrentAndExplicitRestoreRepairsInvalid(t *testing.T) {
	store := testStore(t)
	if _, err := store.Load(); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(store.Path())
	blocked := NewStore(Options{Path: store.Path(), Rename: func(string, string) error { return errors.New("blocked replacement") }})
	dark := AppearanceDark
	if _, err := blocked.Patch(Patch{Appearance: &dark}); err == nil {
		t.Fatal("blocked replacement accepted")
	}
	after, _ := os.ReadFile(store.Path())
	if string(before) != string(after) {
		t.Fatal("failed replacement changed current bytes")
	}
	for _, invalid := range []string{`{`, `{"version":2,"appearance":"system"}`} {
		if err := os.WriteFile(store.Path(), []byte(invalid), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := store.Load(); err == nil {
			t.Fatal("invalid current preferences accepted")
		}
	}
	if err := os.WriteFile(store.Path(), []byte(`{`), 0o600); err != nil {
		t.Fatal(err)
	}
	prefs, err := store.Reset("all")
	if err != nil || prefs.Appearance != AppearanceSystem || prefs.Transition.Status != "not_required" {
		t.Fatalf("restore=%+v %v", prefs, err)
	}
}
