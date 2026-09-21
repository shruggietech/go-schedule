package store

import "testing"

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

func TestPortableIdentityResolutionAndBinding(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

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
}
