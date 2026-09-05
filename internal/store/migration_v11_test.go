package store

import (
	"database/sql"
	"testing"
)

func TestMigrationV11PreservesConfiguredTasksAndAllowsNoSchedule(t *testing.T) {
	path := t.TempDir() + "/v10.db"
	createV9DiagnosticFixture(t, path, false)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(migrations[9].stmts); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO schema_version(version) VALUES(10)`); err != nil {
		t.Fatal(err)
	}
	_ = db.Close()

	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	configured, err := st.GetTask("t")
	if err != nil || configured.ScheduleID != "s" {
		t.Fatalf("configured task = %+v, err=%v", configured, err)
	}
	if _, err := st.db.Exec(`INSERT INTO tasks(id,name,command,schedule_id,created_at,updated_at) VALUES('draft','','',NULL,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z')`); err != nil {
		t.Fatalf("insert unscheduled task: %v", err)
	}
	if _, err := st.GetTask("draft"); err != nil {
		t.Fatalf("read unscheduled task: %v", err)
	}
	var violations int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM pragma_foreign_key_check`).Scan(&violations); err != nil || violations != 0 {
		t.Fatalf("foreign key violations=%d err=%v", violations, err)
	}
}
