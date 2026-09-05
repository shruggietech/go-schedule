package store

import (
	"database/sql"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestMigrationV14AddsFilesystemWatchersAndPreservesV13Data(t *testing.T) {
	path := t.TempDir() + "/v13.db"
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	for _, migration := range migrations[:13] {
		if _, err := db.Exec(migration.stmts); err != nil {
			t.Fatalf("apply v%d: %v", migration.version, err)
		}
	}
	if _, err := db.Exec(`CREATE TABLE schema_version(version INTEGER NOT NULL); INSERT INTO schema_version(version) VALUES(13)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO schedules(id,kind,human_summary,expression) VALUES('schedule','once','',''); INSERT INTO tasks(id,name,command,args_json,working_dir,env_json,run_as,enabled,timezone,schedule_id,overlap_policy,catchup_policy,state,created_at,updated_at) VALUES('task','target','echo','[]','','{}','',1,'UTC','schedule','queue_one','one','active','2026-09-05T00:00:00Z','2026-09-05T00:00:00Z'); INSERT INTO runs(id,task_id,scheduled_for,outcome,output,trigger) VALUES('run','task','2026-09-05T00:00:00Z','success','','manual')`); err != nil {
		t.Fatal(err)
	}
	_ = db.Close()

	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	var tableCount, columnCount int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='filesystem_watchers'`).Scan(&tableCount); err != nil || tableCount != 1 {
		t.Fatalf("watcher table count=%d err=%v", tableCount, err)
	}
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('runs') WHERE name='source_watcher_id'`).Scan(&columnCount); err != nil || columnCount != 1 {
		t.Fatalf("source watcher column count=%d err=%v", columnCount, err)
	}
	run, err := st.GetRun("run")
	if err != nil || run.TaskID != "task" || run.SourceWatcherID != "" || run.Trigger != domain.TriggerManual {
		t.Fatalf("preserved run=%+v err=%v", run, err)
	}
	var version int
	if err := st.db.QueryRow(`SELECT MAX(version) FROM schema_version`).Scan(&version); err != nil || version != 14 {
		t.Fatalf("schema version=%d err=%v", version, err)
	}
}
