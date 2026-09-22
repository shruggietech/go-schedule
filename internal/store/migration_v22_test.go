package store

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestPortableIdentityIsStableAndIndependentFromObjectID(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	first, err := st.PortableID("task", "daemon-local-task")
	if err != nil {
		t.Fatal(err)
	}
	second, err := st.PortableID("task", "daemon-local-task")
	if err != nil {
		t.Fatal(err)
	}
	if first == "" || first != second || first == "daemon-local-task" {
		t.Fatalf("portable identities %q %q", first, second)
	}
	var version int
	if err := st.db.QueryRow(`SELECT MAX(version) FROM schema_version`).Scan(&version); err != nil || version != 22 {
		t.Fatalf("version=%d err=%v", version, err)
	}
}

func TestPortableSourceIdentitiesResolveOnlyLiveRecords(t *testing.T) {
	st := openMem(t)
	task := domain.Task{Name: "target", Command: "echo", State: domain.TaskActive}
	if err := st.CreateTask(&task); err != nil {
		t.Fatal(err)
	}
	trigger := domain.ExternalTrigger{Name: "call", TargetTaskID: task.ID}
	if err := st.CreateExternalTrigger(&trigger); err != nil {
		t.Fatal(err)
	}
	set := domain.TriggerSet{Name: "set", TargetTaskID: task.ID}
	if err := st.CreateTriggerSet(&set, 2, false); err != nil {
		t.Fatal(err)
	}
	watcher := domain.FilesystemWatcher{Name: "files", Kind: domain.WatcherDirectory, Path: t.TempDir(), TargetTaskID: task.ID, Debounce: 250 * time.Millisecond, Stability: 500 * time.Millisecond}
	if err := st.CreateFilesystemWatcher(&watcher); err != nil {
		t.Fatal(err)
	}
	var triggerPortableID string
	for _, value := range []struct{ kind, id string }{{"external_trigger", trigger.ID}, {"trigger_set", set.ID}, {"watcher", watcher.ID}} {
		portableID, err := st.PortableID(value.kind, value.id)
		if err != nil {
			t.Fatal(err)
		}
		if value.kind == "external_trigger" {
			triggerPortableID = portableID
		}
		resolved, err := st.ObjectIDForPortableID(value.kind, portableID)
		if err != nil || resolved != value.id {
			t.Fatalf("%s resolved %q err=%v", value.kind, resolved, err)
		}
	}
	if err := st.DeleteExternalTrigger(trigger.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := st.ObjectIDForPortableID("external_trigger", triggerPortableID); err != ErrNotFound {
		t.Fatalf("stale trigger identity error=%v", err)
	}
	if _, err := st.ObjectIDForPortableID("external_trigger", "missing"); err != ErrNotFound {
		t.Fatalf("missing identity error=%v", err)
	}
	if _, err := st.ObjectIDForPortableID("unknown", "missing"); err != ErrNotFound {
		t.Fatalf("unknown absent identity error=%v", err)
	}
	if err := st.BindPortableID("unknown", "not-an-object", "unsupported-id"); err != nil {
		t.Fatal(err)
	}
	if _, err := st.ObjectIDForPortableID("unknown", "unsupported-id"); err == nil {
		t.Fatal("unsupported object kind resolved")
	}
}

func TestPortableIdentityResolutionAndBinding(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	for _, id := range []string{"task-1", "task-2"} {
		task := domain.Task{ID: id, Name: id, Command: "echo", Timezone: "UTC", State: domain.TaskActive}
		if err := st.CreateTask(&task); err != nil {
			t.Fatal(err)
		}
	}

	if _, err := st.PortableID("", "task-1"); err == nil {
		t.Fatal("PortableID accepted an empty kind")
	}
	if _, err := st.PortableID("task", ""); err == nil {
		t.Fatal("PortableID accepted an empty object ID")
	}
	if _, err := st.ObjectIDForPortableID("task", "missing"); err != ErrNotFound {
		t.Fatalf("missing portable identity error = %v, want %v", err, ErrNotFound)
	}
	if err := st.BindPortableID("task", "task-1", "portable-task-1"); err != nil {
		t.Fatal(err)
	}
	objectID, err := st.ObjectIDForPortableID("task", "portable-task-1")
	if err != nil {
		t.Fatal(err)
	}
	if objectID != "task-1" {
		t.Fatalf("object ID = %q, want task-1", objectID)
	}
	if err := st.BindPortableID("", "task-2", "portable-task-2"); err == nil {
		t.Fatal("BindPortableID accepted an empty kind")
	}
	if err := st.BindPortableID("task", "task-1", "portable-task-2"); err == nil {
		t.Fatal("BindPortableID accepted a conflicting object binding")
	}
	if _, err := st.db.Exec(`INSERT INTO portable_identities(object_kind,object_id,portable_id,created_at) VALUES(?,?,?,?)`, "task", "deleted-task", "stale-portable-task", fmtTime(st.now())); err != nil {
		t.Fatal(err)
	}
	if _, err := st.ObjectIDForPortableID("task", "stale-portable-task"); err != ErrNotFound {
		t.Fatalf("stale portable identity error = %v, want %v", err, ErrNotFound)
	}
	if err := st.BindPortableID("task", "task-2", "stale-portable-task"); err != nil {
		t.Fatalf("stale portable identity was not removed: %v", err)
	}
}
