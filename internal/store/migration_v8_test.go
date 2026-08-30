package store

import (
	"database/sql"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func openAtV7(t *testing.T, path string) {
	t.Helper()
	openAtV6(t, path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	for _, m := range migrations {
		if m.version == 7 {
			if _, err := db.Exec(m.stmts); err != nil {
				t.Fatal(err)
			}
			if _, err := db.Exec(`INSERT INTO schema_version(version) VALUES (7)`); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestMigrationV8DefaultsAndStdinRoundTrip(t *testing.T) {
	path := t.TempDir() + "/v7.db"
	openAtV7(t, path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO schedules(id,kind,rrule,human_summary,expression,calendar_adjustment) VALUES('s','recurring','FREQ=DAILY;BYHOUR=9','daily','','')`)
	if err == nil {
		_, err = db.Exec(`INSERT INTO tasks(id,name,command,timezone,schedule_id,overlap_policy,catchup_policy,missing_date_policy,time_basis,dst_gap_policy,dst_overlap_policy,state,created_at,updated_at) VALUES('legacy','legacy','true','UTC','s','queue_one','one','skip','wall_clock','next_valid','first','active','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z')`)
	}
	_ = db.Close()
	if err != nil {
		t.Fatal(err)
	}

	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	legacy, err := st.GetTask("legacy")
	if err != nil {
		t.Fatal(err)
	}
	if legacy.Stdin != "" || legacy.Command != "true" || legacy.TimeBasis != domain.TimeBasisWallClock {
		t.Fatalf("legacy task changed: %#v", legacy)
	}
	legacy.Stdin = "alpha\nbeta\n"
	if err := st.UpdateTask(&legacy); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	st, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()
	got, err := st.GetTask("legacy")
	if err != nil {
		t.Fatal(err)
	}
	if got.Stdin != legacy.Stdin {
		t.Fatalf("stdin = %q, want %q", got.Stdin, legacy.Stdin)
	}
}
