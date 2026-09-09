package store

import "testing"

func TestMigrationV17AddsOneLocalActorAndAuditStorage(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	var version, actors, audit int
	if err := st.db.QueryRow(`SELECT MAX(version) FROM schema_version`).Scan(&version); err != nil {
		t.Fatal(err)
	}
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM actors WHERE builtin=1`).Scan(&actors); err != nil {
		t.Fatal(err)
	}
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM audit_events`).Scan(&audit); err != nil {
		t.Fatal(err)
	}
	if version != 17 || actors != 1 || audit != 0 {
		t.Fatalf("version=%d actors=%d audit=%d", version, actors, audit)
	}
}
