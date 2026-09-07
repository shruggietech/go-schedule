package settings

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

type fakeBackend struct {
	info server.RuntimeInfoResponse
	err  error
}

type deadlineBackend struct{ bounded bool }

func (b *deadlineBackend) RuntimeInfo(ctx context.Context) (server.RuntimeInfoResponse, error) {
	deadline, ok := ctx.Deadline()
	b.bounded = ok && time.Until(deadline) > 0 && time.Until(deadline) <= runtimeInfoTimeout
	return server.RuntimeInfoResponse{}, errors.New("offline")
}

func (f fakeBackend) RuntimeInfo(context.Context) (server.RuntimeInfoResponse, error) {
	return f.info, f.err
}

type fakeNative struct {
	mu      sync.Mutex
	copied  string
	opened  string
	copyErr error
	openErr error
}

func (f *fakeNative) ClipboardSetText(_ context.Context, value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.copied = value
	return f.copyErr
}

func (f *fakeNative) BrowserOpenURL(_ context.Context, value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.opened = value
	return f.openErr
}

func testDependencies(t *testing.T) Dependencies {
	t.Helper()
	root := t.TempDir()
	return Dependencies{
		Paths: Paths{
			Preferences:       filepath.Join(root, "current", "preferences.json"),
			LegacyPreferences: filepath.Join(root, "legacy", "preferences.json"),
			ApplicationData:   filepath.Join(root, "current"),
			MachineData:       filepath.Join(root, "machine"),
			Executable:        filepath.Join(root, "bin", "go-schedule"),
			GOOS:              "linux",
		},
		Stat: os.Stat, Now: func() time.Time { return time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC) }, Rename: os.Rename,
	}
}

func writeFixture(t *testing.T, path, value string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(value), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestPreferenceMigrationMatrixAndEstablishedAuthority(t *testing.T) {
	tests := []struct {
		name, legacy, status string
		appearance           Appearance
		directory            bool
	}{
		{name: "missing", status: "not_found", appearance: AppearanceSystem},
		{name: "valid light", legacy: `{"appearance.mode":"light","appearance.font":"ubuntu","appearance.scroll_sensitivity":4}`, status: "migrated", appearance: AppearanceLight},
		{name: "valid dark", legacy: `{"appearance.mode":"dark"}`, status: "migrated", appearance: AppearanceDark},
		{name: "valid system", legacy: `{"appearance.mode":"system"}`, status: "migrated", appearance: AppearanceSystem},
		{name: "invalid mode", legacy: `{"appearance.mode":"sepia"}`, status: "invalid", appearance: AppearanceSystem},
		{name: "malformed", legacy: `{`, status: "invalid", appearance: AppearanceSystem},
		{name: "unreadable", status: "unreadable", appearance: AppearanceSystem, directory: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := testDependencies(t)
			if tt.directory {
				if err := os.MkdirAll(deps.Paths.LegacyPreferences, 0o700); err != nil {
					t.Fatal(err)
				}
			} else if tt.legacy != "" {
				writeFixture(t, deps.Paths.LegacyPreferences, tt.legacy)
			}
			prefs, err := loadPreferences(deps)
			if err != nil {
				t.Fatal(err)
			}
			if prefs.Appearance != tt.appearance || prefs.Transition.Status != tt.status {
				t.Fatalf("preferences=%+v", prefs)
			}
			if len(prefs.Transition.Retired) != 2 {
				t.Fatalf("retired=%v", prefs.Transition.Retired)
			}
			if tt.directory {
				return
			}
			writeFixture(t, deps.Paths.LegacyPreferences, `{"appearance.mode":"dark"}`)
			again, err := loadPreferences(deps)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(again, prefs) {
				t.Fatalf("legacy overwrote current preferences: got %+v want %+v", again, prefs)
			}
		})
	}
}

func TestCurrentPreferencesRejectUnsupportedState(t *testing.T) {
	deps := testDependencies(t)
	for _, fixture := range []string{`{"version":2,"appearance":"system"}`, `{"version":1,"appearance":"sepia"}`, `{`} {
		writeFixture(t, deps.Paths.Preferences, fixture)
		if _, err := loadPreferences(deps); err == nil {
			t.Fatalf("fixture %q was accepted", fixture)
		}
	}
}

func TestRelativePreferencePathIsRejectedWithoutWriting(t *testing.T) {
	deps := testDependencies(t)
	deps.Paths.Preferences = "relative/preferences.json"
	if _, err := loadPreferences(deps); err == nil {
		t.Fatal("relative preference path was accepted")
	}
}

func TestSaveRestoreAndFailedReplacement(t *testing.T) {
	deps := testDependencies(t)
	service := NewServiceWithDependencies(fakeBackend{err: errors.New("offline")}, &fakeNative{}, deps)
	ctx := context.Background()
	if result := service.SaveAppearance(ctx, "light"); result.Outcome != "accepted" || result.Workspace.Preferences.Appearance != AppearanceLight {
		t.Fatalf("save=%+v", result)
	}
	if result := service.SaveAppearance(ctx, "sepia"); result.Outcome != "rejected" {
		t.Fatalf("invalid=%+v", result)
	}
	if result := service.Restore(ctx); result.Outcome != "accepted" || result.Workspace.Preferences.Appearance != AppearanceSystem || result.Workspace.Preferences.Transition.Status != "not_required" {
		t.Fatalf("restore=%+v", result)
	}
	before, err := os.ReadFile(deps.Paths.Preferences)
	if err != nil {
		t.Fatal(err)
	}
	deps.Rename = func(string, string) error { return errors.New("blocked") }
	failed := NewServiceWithDependencies(fakeBackend{}, &fakeNative{}, deps).SaveAppearance(ctx, "dark")
	if failed.Outcome != "unavailable" {
		t.Fatalf("failed=%+v", failed)
	}
	after, err := os.ReadFile(deps.Paths.Preferences)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Fatal("failed replacement changed current preferences")
	}
}

func TestConcurrentPreferenceMutationsRemainValid(t *testing.T) {
	deps := testDependencies(t)
	service := NewServiceWithDependencies(fakeBackend{}, &fakeNative{}, deps)
	var wg sync.WaitGroup
	for _, mode := range []string{"light", "dark", "system", "light"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if got := service.SaveAppearance(context.Background(), mode); got.Outcome != "accepted" {
				t.Errorf("save=%+v", got)
			}
		}()
	}
	wg.Wait()
	prefs, err := loadPreferences(deps)
	if err != nil {
		t.Fatal(err)
	}
	if !validAppearance(prefs.Appearance) {
		t.Fatalf("appearance=%q", prefs.Appearance)
	}
}

func TestWorkspacePreservesStorageTruthOnlineAndOffline(t *testing.T) {
	deps := testDependencies(t)
	inside := filepath.Join(deps.Paths.MachineData, "goschedule.db")
	outside := filepath.Join(filepath.Dir(deps.Paths.MachineData), "external", "events.log")
	info := server.RuntimeInfoResponse{DataDir: deps.Paths.MachineData, DatabasePath: inside, ConfigPath: filepath.Join(deps.Paths.MachineData, "config.json"), LogPath: outside, LockPath: filepath.Join(deps.Paths.MachineData, "lock")}
	result := NewServiceWithDependencies(fakeBackend{info: info}, &fakeNative{}, deps).Workspace(context.Background())
	if result.Outcome != "accepted" || !result.Workspace.DaemonAvailable {
		t.Fatalf("workspace=%+v", result)
	}
	byID := map[string]StorageRecord{}
	for _, record := range result.Workspace.Storage {
		byID[record.ID] = record
	}
	if byID["task-database"].Path != inside || byID["task-database"].Owner != "go-schedule" {
		t.Fatalf("database=%+v", byID["task-database"])
	}
	if byID["logs"].Owner != "external" || byID["logs"].Scope != "external" {
		t.Fatalf("logs=%+v", byID["logs"])
	}
	offline := NewServiceWithDependencies(fakeBackend{err: errors.New("offline")}, &fakeNative{}, deps).Workspace(context.Background())
	if offline.Outcome != "accepted" || offline.Workspace.DaemonAvailable || offline.Workspace.Storage[0].Path != "" || offline.Workspace.Storage[0].Copyable {
		t.Fatalf("offline=%+v", offline)
	}
}

func TestRuntimeInfoReadIsBounded(t *testing.T) {
	deps := testDependencies(t)
	backend := &deadlineBackend{}
	result := NewServiceWithDependencies(backend, &fakeNative{}, deps).Workspace(context.Background())
	if result.Outcome != "accepted" || !backend.bounded {
		t.Fatalf("result=%+v bounded=%v", result, backend.bounded)
	}
}

func TestNativeActionsResolveAllowlistedIdentifiers(t *testing.T) {
	deps := testDependencies(t)
	native := &fakeNative{}
	service := NewServiceWithDependencies(fakeBackend{err: errors.New("offline")}, native, deps)
	ctx := context.Background()
	if got := service.CopyStoragePath(ctx, "desktop-preferences"); got.Outcome != "accepted" {
		t.Fatalf("copy=%+v", got)
	}
	if native.copied != deps.Paths.Preferences {
		t.Fatalf("copied=%q", native.copied)
	}
	native.copied = ""
	if got := service.CopyStoragePath(ctx, "../../arbitrary"); got.Outcome != "rejected" || native.copied != "" {
		t.Fatalf("unknown copy=%+v copied=%q", got, native.copied)
	}
	if got := service.OpenProductLink(ctx, "source"); got.Outcome != "accepted" || native.opened != "https://github.com/shruggietech/go-schedule" {
		t.Fatalf("open=%+v destination=%q", got, native.opened)
	}
	native.opened = ""
	if got := service.OpenProductLink(ctx, "javascript:alert(1)"); got.Outcome != "rejected" || native.opened != "" {
		t.Fatalf("unknown open=%+v destination=%q", got, native.opened)
	}
}

func TestNativeActionFailuresAreSafe(t *testing.T) {
	deps := testDependencies(t)
	native := &fakeNative{copyErr: errors.New("clipboard detail"), openErr: errors.New("browser detail")}
	service := NewServiceWithDependencies(fakeBackend{}, native, deps)
	if got := service.CopyStoragePath(context.Background(), "desktop-preferences"); got.Outcome != "unavailable" || got.Message != "The storage path could not be copied." {
		t.Fatalf("copy=%+v", got)
	}
	if got := service.OpenProductLink(context.Background(), "documentation"); got.Outcome != "unavailable" || got.Message != "The product link could not be opened." {
		t.Fatalf("open=%+v", got)
	}
}
