package store

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestMigrationV16PreservesExistingStateAndAddsIdentity(t *testing.T) {
	path := filepath.Join(t.TempDir(), "v15.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	for _, migration := range migrations[:15] {
		if _, err := db.Exec(migration.stmts); err != nil {
			t.Fatalf("apply v%d: %v", migration.version, err)
		}
	}
	if _, err := db.Exec(`CREATE TABLE schema_version(version INTEGER NOT NULL); INSERT INTO schema_version(version) VALUES(15); INSERT INTO groups(id,name,enabled,created_at,updated_at) VALUES('group-1','Preserved',1,'2026-09-09T00:00:00Z','2026-09-09T00:00:00Z'); INSERT INTO schedules(id,kind,human_summary,expression) VALUES('schedule-1','event','At startup','@reboot'); INSERT INTO tasks(id,name,command,enabled,timezone,schedule_id,overlap_policy,catchup_policy,state,created_at,updated_at) VALUES('task-1','Preserved task','echo',1,'UTC','schedule-1','queue_one','one','active','2026-09-09T00:00:00Z','2026-09-09T00:00:00Z')`); err != nil {
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
	var version, groups, tasks, identities int
	if err := st.db.QueryRow(`SELECT MAX(version) FROM schema_version`).Scan(&version); err != nil {
		t.Fatal(err)
	}
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM groups WHERE id='group-1' AND name='Preserved'`).Scan(&groups); err != nil {
		t.Fatal(err)
	}
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE id='task-1' AND name='Preserved task'`).Scan(&tasks); err != nil {
		t.Fatal(err)
	}
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM daemon_identity`).Scan(&identities); err != nil {
		t.Fatal(err)
	}
	if version != 16 || groups != 1 || tasks != 1 || identities != 1 {
		t.Fatalf("version=%d groups=%d tasks=%d identities=%d", version, groups, tasks, identities)
	}
}
