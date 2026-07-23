package store

import (
	"database/sql"
	"testing"
)

// Migration v5 adds tasks.missing_date_policy. Like v4 it is forward-only and
// non-destructive on a safety-critical surface (CLAUDE.md non-negotiables), and
// it carries a stronger obligation than v4 did: v4's new column was inert, while
// this one is read on the scheduling path. A wrong default would silently move
// installed users' run times, so the default is pinned here explicitly.

// v4Task is a task row as it exists in a pre-v5 database, read through raw SQL
// so the assertion does not depend on the current Go struct.
type v4Task struct {
	id, name, command, timezone, scheduleID string
	overlap, catchup, state                 string
	enabled                                 int
}

// openAtV4 creates a database carrying only migrations v1..v4 — the shape a
// v0.6.0 installation has on disk — without going through Open, which would
// apply the migration under test.
func openAtV4(t *testing.T, path string) {
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
		if m.version > 4 {
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

func readTasksRaw(t *testing.T, path string) map[string]v4Task {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	rows, err := db.Query(`SELECT id,name,command,timezone,schedule_id,overlap_policy,catchup_policy,state,enabled FROM tasks ORDER BY id`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	out := map[string]v4Task{}
	for rows.Next() {
		var x v4Task
		if err := rows.Scan(&x.id, &x.name, &x.command, &x.timezone, &x.scheduleID,
			&x.overlap, &x.catchup, &x.state, &x.enabled); err != nil {
			t.Fatal(err)
		}
		out[x.id] = x
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}

// seedV4Tasks writes one schedule and two tasks in the pre-v5 column shape.
func seedV4Tasks(t *testing.T, path string) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	if _, err := db.Exec(
		`INSERT INTO schedules(id,kind,rrule,anchor,human_summary) VALUES(?,?,?,?,?)`,
		"s-1", "recurring", "FREQ=MONTHLY;BYMONTHDAY=31;BYHOUR=9;BYMINUTE=0;BYSECOND=0",
		"2026-01-31T09:00:00Z", "The 31st of every month at 09:00",
	); err != nil {
		t.Fatal(err)
	}
	for _, x := range []v4Task{
		{id: "t-1", name: "month end", command: "close", timezone: "UTC", scheduleID: "s-1",
			overlap: "queue_one", catchup: "one", state: "active", enabled: 1},
		{id: "t-2", name: "disabled one", command: "probe", timezone: "America/New_York", scheduleID: "s-1",
			overlap: "skip", catchup: "none", state: "disabled", enabled: 0},
	} {
		if _, err := db.Exec(
			`INSERT INTO tasks(id,name,command,timezone,schedule_id,overlap_policy,catchup_policy,state,enabled,created_at,updated_at)
			 VALUES(?,?,?,?,?,?,?,?,?,?,?)`,
			x.id, x.name, x.command, x.timezone, x.scheduleID, x.overlap, x.catchup, x.state, x.enabled,
			"2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z",
		); err != nil {
			t.Fatalf("seed %s: %v", x.id, err)
		}
	}
}

// TestMigration_V5DefaultsToSkipAndPreservesTasks is the FR-020 / FR-026 / SC-006
// gate: a pre-v5 database upgrades with the new column present and defaulted to
// "skip" — the behavior those tasks already had — and every pre-existing task row
// otherwise byte-identical. If this fails, installed users' run times have moved
// and the change must not ship.
func TestMigration_V5DefaultsToSkipAndPreservesTasks(t *testing.T) {
	path := t.TempDir() + "/v4.db"
	openAtV4(t, path)
	seedV4Tasks(t, path)

	before := readTasksRaw(t, path)
	if len(before) != 2 {
		t.Fatalf("seeded %d tasks, want 2", len(before))
	}

	st, err := Open(path)
	if err != nil {
		t.Fatalf("upgrade from v4 failed: %v", err)
	}
	closed := false
	t.Cleanup(func() {
		if !closed {
			_ = st.Close()
		}
	})

	// (a) The new column exists and every pre-existing task reads "skip".
	for id := range before {
		var policy string
		if err := st.db.QueryRow(`SELECT missing_date_policy FROM tasks WHERE id=?`, id).Scan(&policy); err != nil {
			t.Fatalf("missing_date_policy column missing after v5 for %s: %v", id, err)
		}
		if policy != "skip" {
			t.Errorf("%s: missing_date_policy = %q after migration, want %q", id, policy, "skip")
		}
	}

	// (b) The value the store hands back is the default, not an empty string —
	// an empty string would compare unequal to every known policy.
	task, err := st.GetTask("t-1")
	if err != nil {
		t.Fatal(err)
	}
	if task.MissingDatePolicy != "skip" {
		t.Errorf("GetTask MissingDatePolicy = %q, want %q", task.MissingDatePolicy, "skip")
	}

	var version int
	if err := st.db.QueryRow(`SELECT COALESCE(MAX(version),0) FROM schema_version`).Scan(&version); err != nil {
		t.Fatal(err)
	}
	if version < 5 {
		t.Fatalf("schema version = %d after upgrade, want >= 5", version)
	}
	_ = st.Close()
	closed = true

	// (c) Every other column of every pre-existing row is unchanged.
	after := readTasksRaw(t, path)
	if len(after) != len(before) {
		t.Fatalf("row count changed: %d -> %d", len(before), len(after))
	}
	for id, b := range before {
		a, ok := after[id]
		if !ok {
			t.Errorf("task %s disappeared across the migration", id)
			continue
		}
		if a != b {
			t.Errorf("task %s mutated across the migration:\n before: %+v\n after:  %+v", id, b, a)
		}
	}

	// (d) The schedule's timing columns are untouched — this migration must not
	// have reached the schedules table at all.
	scheds := readSchedulesRaw(t, path)
	s, ok := scheds["s-1"]
	if !ok {
		t.Fatal("seeded schedule disappeared across the migration")
	}
	if s.rrule.String != "FREQ=MONTHLY;BYMONTHDAY=31;BYHOUR=9;BYMINUTE=0;BYSECOND=0" {
		t.Errorf("schedule rrule changed across the migration: %q", s.rrule.String)
	}
	if s.anchor.String != "2026-01-31T09:00:00Z" {
		t.Errorf("schedule anchor changed across the migration: %q", s.anchor.String)
	}
}

// TestMigration_V5IdempotentReopen verifies re-opening an already-upgraded
// database does not re-apply v5 (ALTER TABLE ADD COLUMN is not idempotent on its
// own — the version gate is what makes it safe).
func TestMigration_V5IdempotentReopen(t *testing.T) {
	path := t.TempDir() + "/reopen5.db"
	openAtV4(t, path)

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
	if err := st2.db.QueryRow(`SELECT COUNT(*) FROM schema_version WHERE version=5`).Scan(&rows); err != nil {
		t.Fatal(err)
	}
	if rows != 1 {
		t.Fatalf("migration v5 recorded %d times, want exactly 1", rows)
	}
}
