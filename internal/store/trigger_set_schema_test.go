package store

import (
	"database/sql"
	"testing"
)

func TestCurrentSchemaAddsTriggerSetsAndPreservesStandaloneTriggers(t *testing.T) {
	path := t.TempDir() + "/current.db"
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := db.Exec(currentSchema + `INSERT INTO schema_version(version) VALUES(24)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO schedules(id,kind,human_summary,expression) VALUES('schedule','once','',''); INSERT INTO tasks(id,name,command,args_json,working_dir,env_json,run_as,enabled,timezone,schedule_id,overlap_policy,catchup_policy,state,created_at,updated_at) VALUES('task','target','echo','[]','','{}','',1,'UTC','schedule','queue_one','one','active','2026-09-05T00:00:00Z','2026-09-05T00:00:00Z'); INSERT INTO external_triggers(id,name,key,target_task_id,enabled,created_at,updated_at) VALUES('trigger','standalone','gst_secret','task',1,'2026-09-05T00:00:00Z','2026-09-05T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	_ = db.Close()
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	var tableCount int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='external_trigger_sets'`).Scan(&tableCount); err != nil || tableCount != 1 {
		t.Fatalf("trigger set table count=%d err=%v", tableCount, err)
	}
	for _, column := range []string{"set_id", "set_position"} {
		var count int
		if err := st.db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('external_triggers') WHERE name=?`, column).Scan(&count); err != nil || count != 1 {
			t.Fatalf("column %s count=%d err=%v", column, count, err)
		}
	}
	trigger, err := st.GetExternalTrigger("trigger")
	if err != nil || trigger.SetID != "" || trigger.SetPosition != 0 || trigger.Key != "gst_secret" {
		t.Fatalf("preserved standalone trigger=%+v err=%v", trigger, err)
	}
}
