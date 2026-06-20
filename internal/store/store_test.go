package store

import (
	"testing"
	"time"

	"github.com/shruggietech/go-scheduler/internal/domain"
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

func TestTriggers_CRUDAndDedup(t *testing.T) {
	st := openMem(t)
	// Need two tasks for FK.
	mk := func(name string) string {
		sch := &domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=MINUTELY;INTERVAL=1"}
		_ = st.CreateSchedule(sch)
		task := &domain.Task{Name: name, Command: "x", Timezone: "UTC", ScheduleID: sch.ID, State: domain.TaskActive}
		_ = st.CreateTask(task)
		return task.ID
	}
	a, b := mk("A"), mk("B")

	tr := &domain.Trigger{SourceTaskID: a, TargetTaskID: b, DedupWindow: time.Minute}
	if err := st.CreateTrigger(tr); err != nil {
		t.Fatal(err)
	}
	if tr.OnOutcome != domain.OnSuccess {
		t.Fatalf("default outcome should be success, got %q", tr.OnOutcome)
	}

	if got, _ := st.ListTriggers(); len(got) != 1 {
		t.Fatalf("ListTriggers = %d, want 1", len(got))
	}
	if got, _ := st.ListTriggersBySource(a); len(got) != 1 {
		t.Fatalf("ListTriggersBySource = %d, want 1", len(got))
	}
	if g, err := st.GetTrigger(tr.ID); err != nil || g.ID != tr.ID {
		t.Fatalf("GetTrigger: %v", err)
	}
	if _, err := st.GetTrigger("missing"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	now := time.Now().UTC()
	if ok, _ := st.ClaimEvent(tr.ID, "k", time.Minute, now); !ok {
		t.Fatal("first claim should succeed")
	}
	if ok, _ := st.ClaimEvent(tr.ID, "k", time.Minute, now.Add(30*time.Second)); ok {
		t.Fatal("claim within window should be deduped")
	}
	if ok, _ := st.ClaimEvent(tr.ID, "k", time.Minute, now.Add(2*time.Minute)); !ok {
		t.Fatal("claim after window should succeed (new event)")
	}
	if err := st.MarkExecuted(tr.ID, "k"); err != nil {
		t.Fatal(err)
	}
	if p, _ := st.PendingClaims(); len(p) != 0 {
		t.Fatalf("no pending claims expected after MarkExecuted, got %d", len(p))
	}

	if err := st.DeleteTrigger(tr.ID); err != nil {
		t.Fatal(err)
	}
	if err := st.DeleteTrigger("missing"); err != ErrNotFound {
		t.Fatalf("delete missing should be ErrNotFound, got %v", err)
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
