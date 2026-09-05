package store

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestFilesystemWatcherLifecycleAndNormalization(t *testing.T) {
	st := openMem(t)
	first := &domain.Task{Name: "first", Command: "echo", State: domain.TaskActive}
	second := &domain.Task{Name: "second", Command: "echo", State: domain.TaskActive}
	if err := st.CreateTask(first); err != nil {
		t.Fatal(err)
	}
	if err := st.CreateTask(second); err != nil {
		t.Fatal(err)
	}
	watcher := &domain.FilesystemWatcher{Name: "ready", Kind: domain.WatcherDirectory, Path: "relative", TargetTaskID: first.ID, Enabled: true}
	if err := st.CreateFilesystemWatcher(watcher); err != nil {
		t.Fatal(err)
	}
	if !filepath.IsAbs(watcher.Path) || watcher.Pattern != "*" || watcher.Debounce != DefaultWatcherDebounce || watcher.Stability != DefaultWatcherStability {
		t.Fatalf("normalized watcher=%+v", watcher)
	}
	loaded, err := st.GetFilesystemWatcher(watcher.ID)
	if err != nil || loaded.TargetTaskName != "first" || loaded.Path != watcher.Path {
		t.Fatalf("loaded watcher=%+v err=%v", loaded, err)
	}
	loaded.Name = "renamed"
	loaded.TargetTaskID = second.ID
	loaded.Pattern = "*.ready"
	loaded.Recursive = true
	if err := st.UpdateFilesystemWatcher(&loaded); err != nil {
		t.Fatal(err)
	}
	updated, _ := st.GetFilesystemWatcher(watcher.ID)
	if updated.Name != "renamed" || updated.TargetTaskID != second.ID || updated.Pattern != "*.ready" || !updated.Recursive {
		t.Fatalf("updated watcher=%+v", updated)
	}
	items, err := st.ListFilesystemWatchers()
	if err != nil || len(items) != 1 || items[0].ID != watcher.ID {
		t.Fatalf("watchers=%+v err=%v", items, err)
	}
	if err := st.SetFilesystemWatcherEnabled(watcher.ID, false); err != nil {
		t.Fatal(err)
	}
	updated, _ = st.GetFilesystemWatcher(watcher.ID)
	if updated.Enabled {
		t.Fatal("watcher remained enabled")
	}
	if err := st.DeleteFilesystemWatcher(watcher.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetFilesystemWatcher(watcher.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("deleted watcher error=%v", err)
	}
}

func TestFilesystemWatcherValidationAndTaskCascade(t *testing.T) {
	st := openMem(t)
	invalid := []domain.FilesystemWatcher{
		{},
		{Name: "x", Kind: "other", Path: ".", TargetTaskID: "missing"},
		{Name: "x", Kind: domain.WatcherDirectory, Path: ".", Pattern: "bad/path", TargetTaskID: "missing"},
		{Name: "x", Kind: domain.WatcherDirectory, Path: ".", Pattern: "[", TargetTaskID: "missing"},
		{Name: "x", Kind: domain.WatcherDirectory, Path: ".", TargetTaskID: "missing", Debounce: time.Millisecond, Stability: time.Second},
	}
	for _, watcher := range invalid {
		if err := st.CreateFilesystemWatcher(&watcher); !errors.Is(err, ErrInvalidWatcher) {
			t.Fatalf("watcher=%+v error=%v", watcher, err)
		}
	}
	task := &domain.Task{Name: "target", Command: "echo", State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	watcher := &domain.FilesystemWatcher{Name: "exact", Kind: domain.WatcherFile, Path: filepath.Join(t.TempDir(), "target.txt"), Pattern: "ignored", Recursive: true, TargetTaskID: task.ID, Enabled: true}
	if err := st.CreateFilesystemWatcher(watcher); err != nil {
		t.Fatal(err)
	}
	if watcher.Pattern != "" || watcher.Recursive {
		t.Fatalf("file watcher retained directory options: %+v", watcher)
	}
	if err := st.DeleteTask(task.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetFilesystemWatcher(watcher.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cascaded watcher error=%v", err)
	}
	if err := st.SetFilesystemWatcherEnabled("missing", true); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing state error=%v", err)
	}
	if err := st.DeleteFilesystemWatcher("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing delete error=%v", err)
	}
}

func TestEnabledFilesystemWatcherIsAutomaticSource(t *testing.T) {
	st := openMem(t)
	task := &domain.Task{Name: "target", Command: "echo", State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	watcher := &domain.FilesystemWatcher{Name: "source", Kind: domain.WatcherDirectory, Path: t.TempDir(), TargetTaskID: task.ID, Enabled: true}
	if err := st.CreateFilesystemWatcher(watcher); err != nil {
		t.Fatal(err)
	}
	if has, err := st.TaskHasEnabledWatcher(task.ID); err != nil || !has {
		t.Fatalf("enabled watcher=%t err=%v", has, err)
	}
	if err := st.SetTaskEnabled(task.ID, true); err != nil {
		t.Fatalf("enable watcher-backed task: %v", err)
	}
	if err := st.SetFilesystemWatcherEnabled(watcher.ID, false); err != nil {
		t.Fatal(err)
	}
	loaded, _ := st.GetTask(task.ID)
	if loaded.Enabled {
		t.Fatal("task remained enabled after final watcher source was disabled")
	}
}

func TestFilesystemWatcherAndRunProvenanceRoundTrip(t *testing.T) {
	st := openMem(t)
	task := &domain.Task{Name: "target", Command: "echo", State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	run := &domain.Run{TaskID: task.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerFilesystem, SourceWatcherID: "watcher-id"}
	if err := st.CreateRun(run); err != nil {
		t.Fatal(err)
	}
	loaded, err := st.GetRun(run.ID)
	if err != nil || loaded.Trigger != domain.TriggerFilesystem || loaded.SourceWatcherID != "watcher-id" {
		t.Fatalf("run=%+v err=%v", loaded, err)
	}
}

func TestFilesystemWatcherLifecycleBudgetAtOneHundred(t *testing.T) {
	st := openMem(t)
	task := &domain.Task{Name: "target", Command: "echo", State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	started := time.Now()
	for i := 0; i < 100; i++ {
		watcher := &domain.FilesystemWatcher{Name: "watcher", Kind: domain.WatcherDirectory, Path: t.TempDir(), TargetTaskID: task.ID, Enabled: i%2 == 0}
		if err := st.CreateFilesystemWatcher(watcher); err != nil {
			t.Fatal(err)
		}
	}
	items, err := st.ListFilesystemWatchers()
	if err != nil || len(items) != 100 {
		t.Fatalf("watchers=%d err=%v", len(items), err)
	}
	if elapsed := time.Since(started); elapsed >= time.Second {
		t.Fatalf("100 watcher lifecycle took %v", elapsed)
	}
}
