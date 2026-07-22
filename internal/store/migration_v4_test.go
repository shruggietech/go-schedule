package store

import (
	"database/sql"
	"testing"
)

// Migration v4 adds schedules.expression. It is a forward-only, non-destructive
// migration on a safety-critical surface (see CLAUDE.md non-negotiables), so it
// is pinned by an explicit upgrade test rather than only by the store's normal
// round-trip coverage: an already-installed database must survive the upgrade
// with every stored schedule — and therefore every task's timing — untouched.

// v3Schedule is a schedule row as it exists in a pre-v4 database, read back
// through raw SQL so the assertion does not depend on the current Go struct.
type v3Schedule struct {
	id, kind, humanSummary string
	rrule, anchor, runAt   sql.NullString
}

// openAtV3 creates a database carrying only migrations v1..v3 — the shape a
// v0.3.0 installation has on disk — without going through Open (which would
// apply every migration including the one under test).
func openAtV3(t *testing.T, path string) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_version (version INTEGER NOT NULL)`); err != nil {
		t.Fatal(err)
	}
	for _, m := range migrations {
		if m.version > 3 {
			continue
		}
		if _, err := db.Exec(m.stmts); err != nil {
			t.Fatalf("seed migration %d: %v", m.version, err)
		}
		if _, err := db.Exec(`INSERT INTO schema_version(version) VALUES (?)`, m.version); err != nil {
			t.Fatal(err)
		}
	}
}

func readSchedulesRaw(t *testing.T, path string) map[string]v3Schedule {
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
	out := map[string]v3Schedule{}
	for rows.Next() {
		var s v3Schedule
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

// TestMigration_V4PreservesExistingSchedules is the FR-002 / FR-024 gate: a
// pre-v4 database upgrades with the new column present and defaulted, and every
// pre-existing schedule row otherwise byte-identical. If this test fails,
// installed users' task timings have moved and the change must not ship.
func TestMigration_V4PreservesExistingSchedules(t *testing.T) {
	path := t.TempDir() + "/v3.db"
	openAtV3(t, path)

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

	// Upgrade through the real path.
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
			t.Errorf("%s: expression = %q after migration, want empty default", id, expr)
		}
	}

	var version int
	if err := st.db.QueryRow(`SELECT COALESCE(MAX(version),0) FROM schema_version`).Scan(&version); err != nil {
		t.Fatal(err)
	}
	if version < 4 {
		t.Fatalf("schema version = %d after upgrade, want >= 4", version)
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
			t.Errorf("schedule %s disappeared across the migration", id)
			continue
		}
		if a != b {
			t.Errorf("schedule %s mutated across the migration:\n before: %+v\n after:  %+v", id, b, a)
		}
	}
}

// TestMigration_V4IdempotentReopen verifies re-opening an already-upgraded
// database does not re-apply v4 (ALTER TABLE ADD COLUMN is not idempotent on
// its own — the version gate is what makes it safe).
func TestMigration_V4IdempotentReopen(t *testing.T) {
	path := t.TempDir() + "/reopen.db"
	openAtV3(t, path)

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
	if err := st2.db.QueryRow(`SELECT COUNT(*) FROM schema_version WHERE version=4`).Scan(&rows); err != nil {
		t.Fatal(err)
	}
	if rows != 1 {
		t.Fatalf("migration v4 recorded %d times, want exactly 1", rows)
	}
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
