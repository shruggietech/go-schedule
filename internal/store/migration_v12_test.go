package store

import (
	"database/sql"
	"testing"
)

func TestMigrationV12AddsExternalTriggersAndRunProvenance(t *testing.T) {
	path := t.TempDir() + "/v11.db"
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	for _, migration := range migrations[:11] {
		if _, err := db.Exec(migration.stmts); err != nil {
			t.Fatalf("apply v%d: %v", migration.version, err)
		}
	}
	if _, err := db.Exec(`CREATE TABLE schema_version(version INTEGER NOT NULL); INSERT INTO schema_version(version) VALUES(11)`); err != nil {
		t.Fatal(err)
	}
	_ = db.Close()
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	for _, column := range []string{"source_trigger_id"} {
		var count int
		if err := st.db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('runs') WHERE name=?`, column).Scan(&count); err != nil || count != 1 {
			t.Fatalf("column %s count=%d err=%v", column, count, err)
		}
	}
	var tableCount int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='external_triggers'`).Scan(&tableCount); err != nil || tableCount != 1 {
		t.Fatalf("external_triggers table count=%d err=%v", tableCount, err)
	}
}
