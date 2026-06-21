package store

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func openMem(t *testing.T) *Store {
	t.Helper()
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

// TestMigration_V3RemovesTriggers verifies migration v3 drops the triggers
// feature tables and the store opens cleanly (the v1/v2 migrations create the
// tables; v3 drops them, so a freshly opened DB must not have them).
func TestMigration_V3RemovesTriggers(t *testing.T) {
	st := openMem(t)

	for _, tbl := range []string{"triggers", "dedup_ledger"} {
		var name string
		err := st.db.QueryRow(
			`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, tbl,
		).Scan(&name)
		if err == nil {
			t.Errorf("table %q should have been dropped by migration v3", tbl)
		}
	}

	var version int
	if err := st.db.QueryRow(`SELECT COALESCE(MAX(version),0) FROM schema_version`).Scan(&version); err != nil {
		t.Fatal(err)
	}
	if version < 3 {
		t.Fatalf("schema version = %d, want >= 3", version)
	}
}

// TestMigration_V3IdempotentReopen verifies re-opening an already-migrated
// database is a clean no-op (the DROP ... IF EXISTS guards in v3 do not error,
// and v3 is not re-applied).
func TestMigration_V3IdempotentReopen(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/goschedule.db"
	st1, err := Open(path)
	if err != nil {
		t.Fatalf("first open: %v", err)
	}
	_ = st1.Close()
	st2, err := Open(path)
	if err != nil {
		t.Fatalf("reopen of migrated db should be a clean no-op: %v", err)
	}
	_ = st2.Close()
}

func TestGroups_RenameReparentTreeAndChain(t *testing.T) {
	st := openMem(t)
	root := &domain.Group{Name: "root", Enabled: true}
	_ = st.CreateGroup(root)
	child := &domain.Group{Name: "child", ParentID: root.ID, Enabled: true}
	_ = st.CreateGroup(child)
	loose := &domain.Group{Name: "loose", Enabled: true}
	_ = st.CreateGroup(loose)

	if err := st.RenameGroup(root.ID, "ROOT"); err != nil {
		t.Fatal(err)
	}
	if g, _ := st.GetGroup(root.ID); g.Name != "ROOT" {
		t.Fatalf("rename failed: %q", g.Name)
	}

	// Reparent loose under child (valid).
	if err := st.SetGroupParent(loose.ID, child.ID); err != nil {
		t.Fatalf("valid reparent failed: %v", err)
	}
	// Cycle: root under loose (loose is now a descendant of root).
	if err := st.SetGroupParent(root.ID, loose.ID); err != ErrCycle {
		t.Fatalf("expected ErrCycle, got %v", err)
	}

	tree, err := st.GroupTree()
	if err != nil {
		t.Fatal(err)
	}
	if len(tree) != 1 || tree[0].Group.ID != root.ID {
		t.Fatalf("expected single root in tree, got %d", len(tree))
	}

	// Chain enabled: disable root → child ineligible.
	_ = st.SetGroupEnabled(root.ID, false)
	if ok, _ := st.GroupChainEnabled(child.ID); ok {
		t.Fatal("child chain should be disabled when root is disabled")
	}
	if ok, _ := st.GroupChainEnabled(""); !ok {
		t.Fatal("empty group is always enabled")
	}
}

func TestSchedule_GetNotFoundAndRunsLimit(t *testing.T) {
	st := openMem(t)
	if _, err := st.GetSchedule("nope"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	sch := &domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=MINUTELY;INTERVAL=1"}
	_ = st.CreateSchedule(sch)
	task := &domain.Task{Name: "t", Command: "x", Timezone: "UTC", ScheduleID: sch.ID, State: domain.TaskActive}
	_ = st.CreateTask(task)
	for i := 0; i < 5; i++ {
		when := time.Now().UTC().Add(time.Duration(i) * time.Minute)
		_ = st.CreateRun(&domain.Run{TaskID: task.ID, ScheduledFor: when, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerSchedule})
	}
	runs, err := st.ListRuns(task.ID, 2)
	if err != nil || len(runs) != 2 {
		t.Fatalf("ListRuns limit: got %d err %v", len(runs), err)
	}
	all, _ := st.ListRuns("", 0)
	if len(all) != 5 {
		t.Fatalf("ListRuns all: got %d, want 5", len(all))
	}
}

func TestAlerts_AckFlow(t *testing.T) {
	st := openMem(t)
	a := &domain.Alert{Severity: domain.SeverityWarning, Kind: domain.AlertOverlapQueued, Message: "x"}
	_ = st.CreateAlert(a)
	if got, _ := st.ListAlerts(true); len(got) != 1 {
		t.Fatalf("unacked = %d, want 1", len(got))
	}
	if err := st.AckAlert(a.ID); err != nil {
		t.Fatal(err)
	}
	if got, _ := st.ListAlerts(true); len(got) != 0 {
		t.Fatal("should be no unacked after ack")
	}
	if err := st.AckAlert("missing"); err != ErrNotFound {
		t.Fatalf("ack missing should be ErrNotFound, got %v", err)
	}
}
