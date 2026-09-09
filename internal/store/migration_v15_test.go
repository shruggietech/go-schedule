package store

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestMigrationV15PreservesExistingStateAndAddsNotifications(t *testing.T) {
	path := filepath.Join(t.TempDir(), "v14.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, migration := range migrations[:14] {
		if _, err := db.Exec(migration.stmts); err != nil {
			t.Fatalf("apply v%d: %v", migration.version, err)
		}
	}
	if _, err := db.Exec(`CREATE TABLE schema_version(version INTEGER NOT NULL); INSERT INTO schema_version(version) VALUES(14)`); err != nil {
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
	if err := st.db.QueryRow(`SELECT MAX(version) FROM schema_version`).Scan(&version); err != nil || version != 18 {
		t.Fatalf("schema version=%d err=%v, want 18", version, err)
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
