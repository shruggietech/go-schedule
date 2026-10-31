package store

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestCurrentSchemaPreservesExistingStateAndAddsNotifications(t *testing.T) {
	path := filepath.Join(t.TempDir(), "v14.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec(currentSchema + `INSERT INTO schema_version(version) VALUES(24)`); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	var version int
	if err := st.db.QueryRow(`SELECT MAX(version) FROM schema_version`).Scan(&version); err != nil || version != 24 {
		t.Fatalf("schema version=%d err=%v, want 24", version, err)
	}
	for _, table := range []string{"notification_channels", "notification_assignments", "notification_deliveries"} {
		var name string
		if err := st.db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&name); err != nil {
			t.Fatalf("missing %s: %v", table, err)
		}
	}
	var deliveries int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM notification_deliveries`).Scan(&deliveries); err != nil || deliveries != 0 {
		t.Fatalf("deliveries=%d err=%v, want zero", deliveries, err)
	}
}
