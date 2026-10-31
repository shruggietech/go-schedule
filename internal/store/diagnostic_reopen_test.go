package store

import (
	"database/sql"
	"testing"
)

func createCurrentDiagnosticFixture(t *testing.T, path string) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(currentSchema + "INSERT INTO schema_version(version) VALUES(24);"); err != nil {
		t.Fatal(err)
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

func TestCurrentSchemaPreservesDiagnosticsWithDefaults(t *testing.T) {
	path := t.TempDir() + "/current.db"
	createCurrentDiagnosticFixture(t, path)

	st, err := Open(path)
	if err != nil {
		t.Fatalf("reopen current database: %v", err)
	}
	defer st.Close()
	run, err := st.GetRun("r")
	if err != nil || run.Output != "old output" || run.OutputTruncated {
		t.Fatalf("retained run = %+v, err=%v", run, err)
	}
	alerts, err := st.ListAlerts(false)
	if err != nil || len(alerts) != 1 || alerts[0].RunID != "" || alerts[0].TaskID != "t" {
		t.Fatalf("retained alerts = %+v, err=%v", alerts, err)
	}
}
