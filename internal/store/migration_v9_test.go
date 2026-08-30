package store

import (
	"database/sql"
	"testing"
)

func TestMigrationV9UpgradesV8HistoryWithoutLegacyTriggerRevival(t *testing.T) {
	path := t.TempDir() + "/v8.db"
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE schema_version (version INTEGER NOT NULL)`); err != nil {
		t.Fatal(err)
	}
	for _, migration := range migrations {
		if migration.version > 8 {
			break
		}
		if _, err := db.Exec(migration.stmts); err != nil {
			t.Fatalf("apply historical migration %d: %v", migration.version, err)
		}
		if _, err := db.Exec(`INSERT INTO schema_version(version) VALUES(?)`, migration.version); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec(`INSERT INTO schedules(id,kind,human_summary,expression,calendar_adjustment) VALUES('s','recurring','hourly','','')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO tasks(id,name,command,args_json,working_dir,env_json,stdin,run_as,enabled,timezone,schedule_id,overlap_policy,catchup_policy,missing_date_policy,time_basis,dst_gap_policy,dst_overlap_policy,state,created_at,updated_at) VALUES('t','old','echo','[]','','{}','','',1,'UTC','s','queue_one','one','skip','wall_clock','next_valid','first','active','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO runs(id,task_id,scheduled_for,outcome,output,trigger) VALUES('r','t','2026-01-01T00:00:00Z','success','','schedule')`); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	st, err := Open(path)
	if err != nil {
		t.Fatalf("open v8 database: %v", err)
	}
	defer st.Close()
	runs, err := st.ListRuns("t", 0)
	if err != nil || len(runs) != 1 || runs[0].SourceTaskID != "" || runs[0].SourceRunID != "" {
		t.Fatalf("legacy run = %+v, err=%v", runs, err)
	}
	for _, table := range []string{"completion_chains", "completion_deliveries"} {
		var count int
		if err := st.db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&count); err != nil || count != 1 {
			t.Fatalf("table %q count=%d err=%v", table, count, err)
		}
	}
	for _, table := range []string{"triggers", "dedup_ledger"} {
		var count int
		if err := st.db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&count); err != nil || count != 0 {
			t.Fatalf("legacy table %q count=%d err=%v", table, count, err)
		}
	}
}

func TestMigrationV9FailureRollsBackSchemaAndVersion(t *testing.T) {
	path := t.TempDir() + "/v8-conflict.db"
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE schema_version (version INTEGER NOT NULL)`); err != nil {
		t.Fatal(err)
	}
	for _, migration := range migrations {
		if migration.version > 8 {
			break
		}
		if _, err := db.Exec(migration.stmts); err != nil {
			t.Fatalf("apply historical migration %d: %v", migration.version, err)
		}
		if _, err := db.Exec(`INSERT INTO schema_version(version) VALUES(?)`, migration.version); err != nil {
			t.Fatal(err)
		}
	}
	// Simulate an incompatible local schema edit so v9 fails after beginning.
	if _, err := db.Exec(`ALTER TABLE runs ADD COLUMN source_task_id TEXT`); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	if st, err := Open(path); err == nil {
		_ = st.Close()
		t.Fatal("conflicting v9 migration unexpectedly succeeded")
	}
	db, err = sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var version int
	if err := db.QueryRow(`SELECT MAX(version) FROM schema_version`).Scan(&version); err != nil || version != 8 {
		t.Fatalf("schema version=%d err=%v, want 8", version, err)
	}
	for _, table := range []string{"completion_chains", "completion_deliveries"} {
		var count int
		if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&count); err != nil || count != 0 {
			t.Fatalf("rolled-back table %q count=%d err=%v", table, count, err)
		}
	}
}
