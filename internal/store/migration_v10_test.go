package store

import (
	"database/sql"
	"testing"
)

func createV9DiagnosticFixture(t *testing.T, path string, conflictingAlertColumn bool) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE schema_version (version INTEGER NOT NULL)`); err != nil {
		t.Fatal(err)
	}
	for _, migration := range migrations {
		if migration.version > 9 {
			break
		}
		if _, err := db.Exec(migration.stmts); err != nil {
			t.Fatalf("apply historical migration %d: %v", migration.version, err)
		}
		if _, err := db.Exec(`INSERT INTO schema_version(version) VALUES(?)`, migration.version); err != nil {
			t.Fatal(err)
		}
	}
	if conflictingAlertColumn {
		if _, err := db.Exec(`ALTER TABLE alerts ADD COLUMN run_id TEXT`); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec(`INSERT INTO schedules(id,kind,human_summary,expression,calendar_adjustment) VALUES('s','recurring','hourly','','')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO tasks(id,name,command,args_json,working_dir,env_json,stdin,run_as,enabled,timezone,schedule_id,overlap_policy,catchup_policy,missing_date_policy,time_basis,dst_gap_policy,dst_overlap_policy,state,created_at,updated_at) VALUES('t','old','echo','[]','','{}','','',1,'UTC','s','queue_one','one','skip','wall_clock','next_valid','first','active','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO runs(id,task_id,scheduled_for,outcome,output,trigger) VALUES('r','t','2026-01-01T00:00:00Z','failure','old output','manual')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO alerts(id,task_id,severity,kind,message,created_at,acknowledged) VALUES('a','t','error','run_failed','task run failed','2026-01-01T00:00:01Z',0)`); err != nil {
		t.Fatal(err)
	}
}

func TestMigrationV10PreservesV9DiagnosticsWithSafeDefaults(t *testing.T) {
	path := t.TempDir() + "/v9.db"
	createV9DiagnosticFixture(t, path, false)

	st, err := Open(path)
	if err != nil {
		t.Fatalf("open v9 database: %v", err)
	}
	defer st.Close()
	run, err := st.GetRun("r")
	if err != nil || run.Output != "old output" || run.OutputTruncated {
		t.Fatalf("legacy run = %+v, err=%v", run, err)
	}
	alerts, err := st.ListAlerts(false)
	if err != nil || len(alerts) != 1 || alerts[0].RunID != "" || alerts[0].TaskID != "t" {
		t.Fatalf("legacy alerts = %+v, err=%v", alerts, err)
	}
}

func TestMigrationV10FailureRollsBackBothColumnsAndVersion(t *testing.T) {
	path := t.TempDir() + "/v9-conflict.db"
	createV9DiagnosticFixture(t, path, true)

	if st, err := Open(path); err == nil {
		_ = st.Close()
		t.Fatal("conflicting v10 migration unexpectedly succeeded")
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var version int
	if err := db.QueryRow(`SELECT MAX(version) FROM schema_version`).Scan(&version); err != nil || version != 9 {
		t.Fatalf("schema version=%d err=%v, want 9", version, err)
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('runs') WHERE name='output_truncated'`).Scan(&count); err != nil || count != 0 {
		t.Fatalf("rolled-back runs.output_truncated count=%d err=%v", count, err)
	}
}
