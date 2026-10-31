package store

import (
	"database/sql"
	"testing"
)

type taskSnapshot struct {
	id, name, command, timezone, scheduleID string
	overlap, catchup, state                 string
	enabled                                 int
}

func readTasksRaw(t *testing.T, path string) map[string]taskSnapshot {
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
	out := map[string]taskSnapshot{}
	for rows.Next() {
		var x taskSnapshot
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

func seedCurrentTasks(t *testing.T, path string) {
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
	for _, x := range []taskSnapshot{
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

func TestCurrentSchemaDefaultsToSkipAndPreservesTasks(t *testing.T) {
	path := t.TempDir() + "/current.db"
	createCurrentStore(t, path)
	seedCurrentTasks(t, path)

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
			t.Errorf("%s: missing_date_policy = %q after bootstrap, want %q", id, policy, "skip")
		}
	}

	// (b) The value the store hands back is the default, not an empty string,
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
		t.Fatalf("schema version = %d after reopen, want >= 5", version)
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
			t.Errorf("task %s disappeared across reopen", id)
			continue
		}
		if a != b {
			t.Errorf("task %s mutated across reopen:\n before: %+v\n after:  %+v", id, b, a)
		}
	}

	scheds := readSchedulesRaw(t, path)
	s, ok := scheds["s-1"]
	if !ok {
		t.Fatal("seeded schedule disappeared across reopen")
	}
	if s.rrule.String != "FREQ=MONTHLY;BYMONTHDAY=31;BYHOUR=9;BYMINUTE=0;BYSECOND=0" {
		t.Errorf("schedule rrule changed across reopen: %q", s.rrule.String)
	}
	if s.anchor.String != "2026-01-31T09:00:00Z" {
		t.Errorf("schedule anchor changed across reopen: %q", s.anchor.String)
	}
}

func TestCurrentSchemaTasksIdempotentReopen(t *testing.T) {
	path := t.TempDir() + "/reopen5.db"
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
		t.Fatalf("bootstrap v5 recorded %d times, want exactly 1", rows)
	}
}
