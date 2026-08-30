package store

import (
	"database/sql"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func openAtV6(t *testing.T, path string) {
	t.Helper()
	openAtV5(t, path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	for _, m := range migrations {
		if m.version == 6 {
			if _, err := db.Exec(m.stmts); err != nil {
				t.Fatal(err)
			}
			if _, err := db.Exec(`INSERT INTO schema_version(version) VALUES (6)`); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestMigrationV7DefaultsAndPolicyRoundTrip(t *testing.T) {
	path := t.TempDir() + "/v6.db"
	openAtV6(t, path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO schedules(id,kind,rrule,human_summary,expression,calendar_adjustment) VALUES('s','recurring','FREQ=DAILY;BYHOUR=9','daily','','')`)
	if err == nil {
		_, err = db.Exec(`INSERT INTO tasks(id,name,command,timezone,schedule_id,overlap_policy,catchup_policy,missing_date_policy,state,created_at,updated_at) VALUES('legacy','legacy','true','UTC','s','queue_one','one','skip','active','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z')`)
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
	if legacy.TimeBasis != domain.TimeBasisWallClock || legacy.DSTGapPolicy != domain.DSTGapNextValid || legacy.DSTOverlapPolicy != domain.DSTOverlapFirst {
		t.Fatalf("legacy defaults = %#v", legacy.SchedulePolicy())
	}

	legacy.TimeBasis = domain.TimeBasisElapsed
	legacy.DSTGapPolicy = domain.DSTGapSkip
	legacy.DSTOverlapPolicy = domain.DSTOverlapBoth
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
	if got.TimeBasis != domain.TimeBasisElapsed || got.DSTGapPolicy != domain.DSTGapSkip || got.DSTOverlapPolicy != domain.DSTOverlapBoth {
		t.Fatalf("round trip = %#v", got.SchedulePolicy())
	}
}

func TestMigrationV7PersistsElapsedEpoch(t *testing.T) {
	st, err := Open(t.TempDir() + "/epoch.db")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()
	epoch := time.Date(2026, time.March, 7, 14, 0, 0, 0, time.UTC)
	sch := domain.Schedule{
		Kind: domain.ScheduleRecurring, RRULE: "FREQ=DAILY;BYHOUR=9;BYMINUTE=0;BYSECOND=0",
		Anchor: &epoch, ElapsedEpoch: &epoch,
	}
	if err := st.CreateSchedule(&sch); err != nil {
		t.Fatal(err)
	}
	got, err := st.GetSchedule(sch.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ElapsedEpoch == nil || !got.ElapsedEpoch.Equal(epoch) {
		t.Fatalf("elapsed epoch = %v, want %s", got.ElapsedEpoch, epoch)
	}
}
