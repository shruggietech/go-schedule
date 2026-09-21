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
