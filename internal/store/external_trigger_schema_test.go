package store

import (
	"database/sql"
	"testing"
)

func TestCurrentSchemaAddsExternalTriggersAndRunProvenance(t *testing.T) {
	path := t.TempDir() + "/current.db"
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := db.Exec(currentSchema + `INSERT INTO schema_version(version) VALUES(24)`); err != nil {
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
