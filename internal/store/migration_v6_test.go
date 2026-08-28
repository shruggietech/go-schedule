package store

import (
	"database/sql"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func openAtV5(t *testing.T, path string) {
	t.Helper()
	openAtV4(t, path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	for _, m := range migrations {
		if m.version != 5 {
			continue
		}
		if _, err := db.Exec(m.stmts); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(`INSERT INTO schema_version(version) VALUES (5)`); err != nil {
			t.Fatal(err)
		}
	}
}

func TestMigrationV6DefaultsEmptyAndPersistsAdjustment(t *testing.T) {
	path := t.TempDir() + "/v5.db"
	openAtV5(t, path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO schedules(id,kind,rrule,human_summary,expression) VALUES(?,?,?,?,?)`,
		"legacy", "recurring", "FREQ=MONTHLY;BYMONTHDAY=15", "legacy", "legacy phrase")
	_ = db.Close()
	if err != nil {
		t.Fatal(err)
	}

	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	legacy, err := st.GetSchedule("legacy")
	if err != nil {
		t.Fatal(err)
	}
	if legacy.CalendarAdjustment != "" || legacy.RRULE != "FREQ=MONTHLY;BYMONTHDAY=15" {
		t.Fatalf("legacy schedule changed: %+v", legacy)
	}

	anchor := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	marked := domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=MONTHLY;BYMONTHDAY=15", Anchor: &anchor, CalendarAdjustment: domain.CalendarAdjustmentNearestWeekday}
	if err := st.CreateSchedule(&marked); err != nil {
		t.Fatal(err)
	}
	id := marked.ID
	_ = st.Close()

	st, err = Open(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer func() { _ = st.Close() }()
	got, err := st.GetSchedule(id)
	if err != nil {
		t.Fatal(err)
	}
	if got.CalendarAdjustment != domain.CalendarAdjustmentNearestWeekday {
		t.Fatalf("adjustment after reopen = %q", got.CalendarAdjustment)
	}
	var count int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM schema_version WHERE version=6`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("v6 rows=%d err=%v", count, err)
	}
}
