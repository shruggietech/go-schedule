package engine

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

func TestFireFilesystemWatcherRecordsProvenance(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task := &domain.Task{Name: "target", Command: "echo", State: domain.TaskActive, OverlapPolicy: domain.OverlapAllowConcurrent}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	watcher := &domain.FilesystemWatcher{Name: "drop", Kind: domain.WatcherDirectory, Path: t.TempDir(), Pattern: "*", TargetTaskID: task.ID, Enabled: true}
	if err := st.CreateFilesystemWatcher(watcher); err != nil {
		t.Fatal(err)
	}
	if err := st.SetTaskEnabled(task.ID, true); err != nil {
		t.Fatal(err)
	}
	engine := newEngine(st, instantRunner{})
	if err := engine.FireFilesystemWatcher(watcher.ID); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		runs, err := st.ListRuns(task.ID, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(runs) == 1 {
			if runs[0].Trigger != domain.TriggerFilesystem || runs[0].SourceWatcherID != watcher.ID || runs[0].SourceTriggerID != "" {
				t.Fatalf("run provenance = %+v", runs[0])
			}
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("watcher run was not recorded")
}
