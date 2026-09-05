package store

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func triggerSetTask(t *testing.T, st *Store, name string) domain.Task {
	t.Helper()
	task := domain.Task{Name: name, Command: "echo", State: domain.TaskActive}
	if err := st.CreateTask(&task); err != nil {
		t.Fatal(err)
	}
	return task
}

func TestTriggerSetCreateBoundsOrderingAndCascade(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task := triggerSetTask(t, st, "target")
	for _, count := range []int{0, 100} {
		if err := st.CreateTriggerSet(&domain.TriggerSet{Name: "invalid", TargetTaskID: task.ID}, count, true); !errors.Is(err, ErrInvalidTriggerSet) {
			t.Fatalf("count %d error=%v", count, err)
		}
	}
	set := domain.TriggerSet{Name: "Callers", TargetTaskID: task.ID}
	if err := st.CreateTriggerSet(&set, 99, true); err != nil {
		t.Fatal(err)
	}
	if len(set.Members) != 99 {
		t.Fatalf("members=%d", len(set.Members))
	}
	keys := map[string]bool{}
	for index, member := range set.Members {
		if member.SetID != set.ID || member.SetPosition != index+1 || member.TargetTaskID != task.ID || keys[member.Key] {
			t.Fatalf("member %d=%+v", index, member)
		}
		keys[member.Key] = true
	}
	loaded, err := st.GetTriggerSet(set.ID)
	if err != nil || len(loaded.Members) != 99 || loaded.Members[98].SetPosition != 99 {
		t.Fatalf("loaded set=%+v err=%v", loaded, err)
	}
	if err := st.DeleteTask(task.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetTriggerSet(set.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("set survived target deletion: %v", err)
	}
}

func TestTriggerSetMemberIsolationAndFinalCleanup(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task := triggerSetTask(t, st, "target")
	set := domain.TriggerSet{Name: "Pair", TargetTaskID: task.ID}
	if err := st.CreateTriggerSet(&set, 2, true); err != nil {
		t.Fatal(err)
	}
	first, second := set.Members[0], set.Members[1]
	first.Name = "renamed"
	if err := st.UpdateExternalTrigger(&first); err != nil {
		t.Fatal(err)
	}
	if err := st.SetExternalTriggerEnabled(first.ID, false); err != nil {
		t.Fatal(err)
	}
	if _, err := st.RotateExternalTrigger(first.ID); err != nil {
		t.Fatal(err)
	}
	unchanged, err := st.GetExternalTrigger(second.ID)
	if err != nil || unchanged.ID != second.ID || unchanged.Name != second.Name || unchanged.Key != second.Key || unchanged.TargetTaskID != second.TargetTaskID || unchanged.SetID != second.SetID || unchanged.SetPosition != second.SetPosition || unchanged.Enabled != second.Enabled || !unchanged.CreatedAt.Equal(second.CreatedAt) || !unchanged.UpdatedAt.Equal(second.UpdatedAt) {
		t.Fatalf("sibling changed=%+v want=%+v err=%v", unchanged, second, err)
	}
	first.TargetTaskID = "other"
	if err := st.UpdateExternalTrigger(&first); !errors.Is(err, ErrTriggerSetMemberTarget) {
		t.Fatalf("member retarget error=%v", err)
	}
	if err := st.DeleteExternalTrigger(first.ID); err != nil {
		t.Fatal(err)
	}
	remaining, err := st.GetExternalTrigger(second.ID)
	if err != nil || remaining.Key != second.Key || remaining.SetPosition != 2 {
		t.Fatalf("remaining=%+v err=%v", remaining, err)
	}
	if err := st.DeleteExternalTrigger(second.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetTriggerSet(set.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("empty set survived: %v", err)
	}
}

func TestTriggerSetBoundaryCountsRemainUniqueAcrossOneHundredTrials(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task := triggerSetTask(t, st, "target")
	allKeys := make(map[string]bool)
	for trial := 0; trial < 100; trial++ {
		for _, count := range []int{1, 99} {
			set := domain.TriggerSet{Name: fmt.Sprintf("Trial %d count %d", trial, count), TargetTaskID: task.ID}
			if err := st.CreateTriggerSet(&set, count, true); err != nil {
				t.Fatalf("trial %d count %d: %v", trial, count, err)
			}
			if len(set.Members) != count {
				t.Fatalf("trial %d count %d members=%d", trial, count, len(set.Members))
			}
			for index, member := range set.Members {
				if member.SetPosition != index+1 || allKeys[member.Key] {
					t.Fatalf("trial %d member %d=%+v", trial, index, member)
				}
				allKeys[member.Key] = true
			}
			if _, err := st.DeleteTriggerSet(set.ID); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestTriggerSetMaximumOperationsMeetNominalBudget(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	first := triggerSetTask(t, st, "first")
	second := triggerSetTask(t, st, "second")
	set := domain.TriggerSet{Name: "Maximum", TargetTaskID: first.ID}
	assertUnderSecond(t, "create", func() error { return st.CreateTriggerSet(&set, 99, true) })
	assertUnderSecond(t, "reveal", func() error { _, err := st.GetTriggerSet(set.ID); return err })
	assertUnderSecond(t, "retarget", func() error { _, err := st.RetargetTriggerSet(set.ID, second.ID); return err })
	assertUnderSecond(t, "disable", func() error { _, err := st.SetTriggerSetEnabled(set.ID, false); return err })
	assertUnderSecond(t, "enable", func() error { _, err := st.SetTriggerSetEnabled(set.ID, true); return err })
	assertUnderSecond(t, "rotate", func() error { _, err := st.RotateTriggerSet(set.ID); return err })
	assertUnderSecond(t, "delete", func() error { _, err := st.DeleteTriggerSet(set.ID); return err })
}

func TestTriggerSetListsDuplicateNamesAndReportsInvalidReferences(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task := triggerSetTask(t, st, "target")
	for index := 0; index < 2; index++ {
		set := domain.TriggerSet{Name: "Duplicate", TargetTaskID: task.ID}
		if err := st.CreateTriggerSet(&set, index+1, true); err != nil {
			t.Fatal(err)
		}
	}
	sets, err := st.ListTriggerSets()
	if err != nil || len(sets) != 2 || sets[0].Name != "Duplicate" || len(sets[1].Members) == 0 {
		t.Fatalf("sets=%+v err=%v", sets, err)
	}
	if err := st.CreateTriggerSet(&domain.TriggerSet{Name: "Missing", TargetTaskID: "missing"}, 1, true); !errors.Is(err, ErrNotFound) {
		t.Fatalf("create missing target error=%v", err)
	}
	if _, err := st.RetargetTriggerSet(sets[0].ID, ""); !errors.Is(err, ErrInvalidTriggerSet) {
		t.Fatalf("empty retarget error=%v", err)
	}
	if _, err := st.RetargetTriggerSet(sets[0].ID, "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing retarget target error=%v", err)
	}
	if _, err := st.RetargetTriggerSet("missing", task.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing retarget set error=%v", err)
	}
	if _, err := st.SetTriggerSetEnabled("missing", false); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing enable set error=%v", err)
	}
	if _, err := st.RotateTriggerSet("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing rotate set error=%v", err)
	}
	if _, err := st.DeleteTriggerSet("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing delete set error=%v", err)
	}
}

func assertUnderSecond(t *testing.T, name string, operation func() error) {
	t.Helper()
	started := time.Now()
	if err := operation(); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	if elapsed := time.Since(started); elapsed >= time.Second {
		t.Fatalf("%s took %s", name, elapsed)
	}
}

func TestTriggerSetBroadLifecycleAndRollback(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	firstTask := triggerSetTask(t, st, "first")
	secondTask := triggerSetTask(t, st, "second")
	set := domain.TriggerSet{Name: "Pair", TargetTaskID: firstTask.ID}
	if err := st.CreateTriggerSet(&set, 2, true); err != nil {
		t.Fatal(err)
	}
	if err := st.SetTaskEnabled(firstTask.ID, true); err != nil {
		t.Fatal(err)
	}
	retargeted, err := st.RetargetTriggerSet(set.ID, secondTask.ID)
	if err != nil {
		t.Fatal(err)
	}
	if retargeted.TargetTaskID != secondTask.ID || len(retargeted.Members) != 2 || retargeted.Members[0].TargetTaskID != secondTask.ID {
		t.Fatalf("retarget snapshot=%+v", retargeted)
	}
	oldTarget, _ := st.GetTask(firstTask.ID)
	if oldTarget.Enabled {
		t.Fatal("old target remained enabled after final sources moved")
	}
	disabled, err := st.SetTriggerSetEnabled(set.ID, false)
	if err != nil {
		t.Fatal(err)
	}
	if disabled.Members[0].Enabled || disabled.Members[1].Enabled {
		t.Fatalf("disabled snapshot=%+v", disabled)
	}
	loaded, _ := st.GetTriggerSet(set.ID)
	for _, member := range loaded.Members {
		if member.Enabled || member.TargetTaskID != secondTask.ID {
			t.Fatalf("member after broad changes=%+v", member)
		}
	}
	if _, err := st.db.Exec(`CREATE TRIGGER fail_second_rotation BEFORE UPDATE OF key ON external_triggers WHEN OLD.set_position=2 BEGIN SELECT RAISE(ABORT,'injected'); END`); err != nil {
		t.Fatal(err)
	}
	before := loaded.Members[0].Key
	if _, err := st.RotateTriggerSet(set.ID); err == nil {
		t.Fatal("rotation unexpectedly succeeded")
	}
	after, _ := st.GetExternalTrigger(loaded.Members[0].ID)
	if after.Key != before {
		t.Fatal("failed rotation persisted a partial first key")
	}
	if _, err := st.db.Exec(`DROP TRIGGER fail_second_rotation`); err != nil {
		t.Fatal(err)
	}
	rotated, err := st.RotateTriggerSet(set.ID)
	if err != nil || rotated.Members[0].Key == before {
		t.Fatalf("rotation=%+v err=%v", rotated, err)
	}
	if _, err := st.DeleteTriggerSet(set.ID); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkTriggerSetMaximumLifecycle(b *testing.B) {
	for iteration := 0; iteration < b.N; iteration++ {
		st, err := Open(":memory:")
		if err != nil {
			b.Fatal(err)
		}
		task := domain.Task{Name: "target", Command: "echo", State: domain.TaskActive}
		if err := st.CreateTask(&task); err != nil {
			b.Fatal(err)
		}
		set := domain.TriggerSet{Name: fmt.Sprintf("Set %d", iteration), TargetTaskID: task.ID}
		if err := st.CreateTriggerSet(&set, 99, true); err != nil {
			b.Fatal(err)
		}
		if _, err := st.RotateTriggerSet(set.ID); err != nil {
			b.Fatal(err)
		}
		if _, err := st.DeleteTriggerSet(set.ID); err != nil {
			b.Fatal(err)
		}
		_ = st.Close()
	}
}
