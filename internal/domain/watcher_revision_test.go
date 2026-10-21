package domain

import (
	"testing"
	"time"
)

func TestFilesystemWatcherRevisionTracksDefinitionNotJoinedNames(t *testing.T) {
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	original := FilesystemWatcher{ID: "watcher", Name: "Inbox", Kind: WatcherDirectory, Path: "inbox", Pattern: "*.csv", TargetTaskID: "task", TargetTaskName: "Display name", Enabled: true, CreatedAt: now, UpdatedAt: now}
	revision := FilesystemWatcherRevision(original)
	for name, edit := range map[string]func(*FilesystemWatcher){
		"identity":  func(w *FilesystemWatcher) { w.ID = "other" },
		"name":      func(w *FilesystemWatcher) { w.Name = "Other" },
		"kind":      func(w *FilesystemWatcher) { w.Kind = WatcherFile },
		"path":      func(w *FilesystemWatcher) { w.Path = "other" },
		"pattern":   func(w *FilesystemWatcher) { w.Pattern = "*.json" },
		"recursive": func(w *FilesystemWatcher) { w.Recursive = true },
		"debounce":  func(w *FilesystemWatcher) { w.Debounce = time.Second },
		"stability": func(w *FilesystemWatcher) { w.Stability = time.Second },
		"target":    func(w *FilesystemWatcher) { w.TargetTaskID = "other" },
		"enabled":   func(w *FilesystemWatcher) { w.Enabled = false },
		"created":   func(w *FilesystemWatcher) { w.CreatedAt = now.Add(time.Second) },
		"updated":   func(w *FilesystemWatcher) { w.UpdatedAt = now.Add(time.Second) },
	} {
		t.Run(name, func(t *testing.T) {
			changed := original
			edit(&changed)
			if FilesystemWatcherRevision(changed) == revision {
				t.Fatal("definition change retained previous observation identity")
			}
		})
	}
	original.TargetTaskName = "Renamed display name"
	original.CreatedAt = now.In(time.FixedZone("other", 3600))
	original.UpdatedAt = now.In(time.FixedZone("other", 3600))
	if got := FilesystemWatcherRevision(original); got != revision {
		t.Fatal("joined names or equivalent timestamp representations changed the durable revision")
	}
}
