package store

import (
	"database/sql"
	"testing"
)

type scheduleSnapshot struct {
	id, kind, humanSummary string
	rrule, anchor, runAt   sql.NullString
}

func readSchedulesRaw(t *testing.T, path string) map[string]scheduleSnapshot {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	rows, err := db.Query(`SELECT id,kind,rrule,anchor,run_at,human_summary FROM schedules ORDER BY id`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	out := map[string]scheduleSnapshot{}
	for rows.Next() {
		var s scheduleSnapshot
		if err := rows.Scan(&s.id, &s.kind, &s.rrule, &s.anchor, &s.runAt, &s.humanSummary); err != nil {
			t.Fatal(err)
		}
		out[s.id] = s
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestCurrentSchemaPreservesExistingSchedules(t *testing.T) {
	path := t.TempDir() + "/current.db"
	createCurrentStore(t, path)

	// Seed one row of each kind, covering the columns that carry timing.
	seed := func() {
		db, err := sql.Open("sqlite", path)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = db.Close() }()
		stmts := []struct {
			id, kind, rrule, anchor, runAt, summary string
		}{
			{"s-recur", "recurring", "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR;BYHOUR=9;BYMINUTE=0;BYSECOND=0", "2026-01-01T09:00:00Z", "", "Every weekday at 09:00"},
			{"s-interval", "recurring", "FREQ=MINUTELY;INTERVAL=15", "2026-01-01T00:00:00Z", "", "Every 15 minutes"},
			{"s-oneoff", "one_off", "", "", "2026-12-31T23:59:00Z", "Once at 2026-12-31 23:59 UTC"},
		}
		for _, s := range stmts {
			_, err := db.Exec(
				`INSERT INTO schedules(id,kind,rrule,anchor,run_at,human_summary) VALUES(?,?,?,?,?,?)`,
				s.id, s.kind, nilIfEmpty(s.rrule), nilIfEmpty(s.anchor), nilIfEmpty(s.runAt), s.summary,
			)
			if err != nil {
				t.Fatalf("seed %s: %v", s.id, err)
			}
		}
	}
	seed()
	before := readSchedulesRaw(t, path)
	if len(before) != 3 {
		t.Fatalf("seeded %d schedules, want 3", len(before))
	}

	st, err := Open(path)
	if err != nil {
		t.Fatalf("upgrade from v3 failed: %v", err)
	}
	closed := false
	t.Cleanup(func() {
		if !closed {
			_ = st.Close()
		}
	})

	// (a) The new column exists and defaults to the empty string.
	for id := range before {
		var expr string
		if err := st.db.QueryRow(`SELECT expression FROM schedules WHERE id=?`, id).Scan(&expr); err != nil {
			t.Fatalf("expression column missing after v4 for %s: %v", id, err)
		}
		if expr != "" {
			t.Errorf("%s: expression = %q after bootstrap, want empty default", id, expr)
		}
	}

	var version int
	if err := st.db.QueryRow(`SELECT COALESCE(MAX(version),0) FROM schema_version`).Scan(&version); err != nil {
		t.Fatal(err)
	}
	if version < 4 {
		t.Fatalf("schema version = %d after reopen, want >= 4", version)
	}
	_ = st.Close()
	closed = true

	// (b) Every other column of every pre-existing row is unchanged.
	after := readSchedulesRaw(t, path)
	if len(after) != len(before) {
		t.Fatalf("row count changed: %d -> %d", len(before), len(after))
	}
	for id, b := range before {
		a, ok := after[id]
		if !ok {
			t.Errorf("schedule %s disappeared across reopen", id)
			continue
		}
		if a != b {
			t.Errorf("schedule %s mutated across reopen:\n before: %+v\n after:  %+v", id, b, a)
		}
	}
}

func TestCurrentSchemaSchedulesIdempotentReopen(t *testing.T) {
	path := t.TempDir() + "/reopen.db"
	createCurrentStore(t, path)

	st1, err := Open(path)
	if err != nil {
		t.Fatalf("first upgrade: %v", err)
	}
	_ = st1.Close()

	st2, err := Open(path)
	if err != nil {
		t.Fatalf("reopen of upgraded db should be a clean no-op, got: %v", err)
	}
	defer func() { _ = st2.Close() }()

	var rows int
	if err := st2.db.QueryRow(`SELECT COUNT(*) FROM schema_version WHERE version=24`).Scan(&rows); err != nil {
		t.Fatal(err)
	}
	if rows != 1 {
		t.Fatalf("bootstrap v4 recorded %d times, want exactly 1", rows)
	}
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
