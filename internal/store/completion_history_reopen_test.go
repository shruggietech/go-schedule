package store

import (
	"database/sql"
	"testing"
)

func TestCurrentSchemaPreservesRunHistoryWithoutRetiredTriggers(t *testing.T) {
	path := t.TempDir() + "/current.db"
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(currentSchema + "INSERT INTO schema_version(version) VALUES(24);"); err != nil {
		t.Fatal(err)
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
		t.Fatalf("reopen current database: %v", err)
	}
	defer st.Close()
	runs, err := st.ListRuns("t", 0)
	if err != nil || len(runs) != 1 || runs[0].SourceTaskID != "" || runs[0].SourceRunID != "" {
		t.Fatalf("retained run = %+v, err=%v", runs, err)
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
