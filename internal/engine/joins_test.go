package engine

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

func TestS144ParallelBlockOwnsExactSameTaskBranches(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	entry := setupTask(t, st, domain.OverlapQueueOne)
	branch := setupTask(t, st, domain.OverlapAllowConcurrent)
	target := setupTask(t, st, domain.OverlapQueueOne)
	join := domain.ParallelJoin{Name: "parallel", EntryTaskID: entry.ID, TargetTaskID: target.ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "first", TaskID: branch.ID}, {ID: "b", Name: "second", TaskID: branch.ID}}}
	if err := st.CreateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	runner := &completionRunner{outcomes: map[string]domain.RunOutcome{}, runs: make(chan domain.Run, 8)}
	e, stop := startCompletionEngine(t, st, runner)
	defer stop()
	e.dispatch(entry, time.Date(2026, 11, 10, 0, 0, 0, 0, time.UTC), domain.TriggerManual)
	waitFor(t, func() bool {
		instances, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{DefinitionID: join.ID})
		return err == nil && len(instances) == 1 && instances[0].State == domain.JoinInstanceCompleted
	}, "exact join terminal decision")
	e.runWG.Wait()
	instances, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{DefinitionID: join.ID})
	if err != nil {
		t.Fatal(err)
	}
	if len(instances) != 1 || instances[0].State != domain.JoinInstanceCompleted || len(instances[0].Branches) != 2 {
		t.Fatalf("instances=%+v", instances)
	}
	a, b := instances[0].Branches[0], instances[0].Branches[1]
	if a.RunID == "" || b.RunID == "" || a.RunID == b.RunID || a.DeliveryID == b.DeliveryID {
		t.Fatalf("branch identities=%+v %+v", a, b)
	}
	runs, err := st.ListRuns(target.ID, 0)
	if err != nil || len(runs) != 1 || runs[0].JoinInstanceID != instances[0].ID || runs[0].JoinRole != "target" || runs[0].DeliveryID == "" {
		t.Fatalf("target=%+v err=%v", runs, err)
	}
}

func TestS144QueuedJoinMarkersAndCollapseNeverContribute(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	entry := setupTask(t, st, domain.OverlapQueueOne)
	branch := setupTask(t, st, domain.OverlapQueueOne)
	target := setupTask(t, st, domain.OverlapQueueOne)
	join := domain.ParallelJoin{Name: "block", EntryTaskID: entry.ID, TargetTaskID: target.ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "a", TaskID: branch.ID}, {ID: "b", Name: "b", TaskID: branch.ID}}}
	if err := st.CreateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	source := domain.Run{TaskID: entry.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&source, ""); err != nil {
		t.Fatal(err)
	}
	runner := &completionRunner{outcomes: map[string]domain.RunOutcome{}, runs: make(chan domain.Run, 4)}
	e := newEngine(st, runner)
	e.running[branch.ID] = 1
	var resolvedInstance string
	e.SetOnJoinInstanceChanged(func(id string) { resolvedInstance = id })
	e.drainCompletionDeliveries()
	instances, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{DefinitionID: join.ID})
	if err != nil {
		t.Fatal(err)
	}
	if len(instances) != 1 || instances[0].State != domain.JoinInstanceBlocked || instances[0].TargetDeliveryID != "" {
		t.Fatalf("collapsed=%+v", instances)
	}
	if resolvedInstance != instances[0].ID {
		t.Fatalf("resolution event=%q want=%q", resolvedInstance, instances[0].ID)
	}
	for _, b := range instances[0].Branches {
		if b.Outcome == domain.OutcomeSuccess {
			t.Fatal("queue marker satisfied branch")
		}
	}
	e.finish(branch)
	e.runWG.Wait()
	runs, err := st.ListRuns(target.ID, 0)
	if err != nil || len(runs) != 0 {
		t.Fatalf("blocked target=%+v err=%v", runs, err)
	}
}
func TestS144CanceledBranchBlocksOnlyItsInstance(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	entry := setupTask(t, st, domain.OverlapQueueOne)
	branch := setupTask(t, st, domain.OverlapQueueOne)
	target := setupTask(t, st, domain.OverlapQueueOne)
	join := domain.ParallelJoin{Name: "cancel", EntryTaskID: entry.ID, TargetTaskID: target.ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "a", TaskID: branch.ID}}}
	if err := st.CreateParallelJoin(&join); err != nil {
		t.Fatal(err)
	}
	runner := &completionRunner{outcomes: map[string]domain.RunOutcome{branch.ID: domain.OutcomeCanceled}, runs: make(chan domain.Run, 4)}
	e, stop := startCompletionEngine(t, st, runner)
	defer stop()
	e.dispatch(entry, time.Now().UTC(), domain.TriggerManual)
	waitFor(t, func() bool {
		instances, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{DefinitionID: join.ID})
		return err == nil && len(instances) == 1 && instances[0].State == domain.JoinInstanceBlocked && instances[0].Branches[0].Outcome == domain.OutcomeCanceled
	}, "exact join terminal decision")
	e.runWG.Wait()
	instances, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{DefinitionID: join.ID})
	if err != nil {
		t.Fatal(err)
	}
	if len(instances) != 1 || instances[0].State != domain.JoinInstanceBlocked || instances[0].Branches[0].Outcome != domain.OutcomeCanceled {
		t.Fatalf("instances=%+v", instances)
	}
	task, err := st.GetTask(branch.ID)
	if err != nil || !task.Enabled {
		t.Fatalf("canceled branch disabled task=%+v err=%v", task, err)
	}
}
func TestS144CompletionWorkerAdmissionRechecksTask(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	source := setupTask(t, st, domain.OverlapQueueOne)
	target := setupTask(t, st, domain.OverlapQueueOne)
	addChain(t, st, source, target, domain.CompletionOnSuccess)
	native := domain.Run{TaskID: source.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&native, ""); err != nil {
		t.Fatal(err)
	}
	runner := &completionRunner{outcomes: map[string]domain.RunOutcome{}, runs: make(chan domain.Run, 1)}
	e := newEngine(st, runner)
	for i := 0; i < cap(e.sem); i++ {
		e.sem <- struct{}{}
	}
	e.drainCompletionDeliveries()
	if err := st.SetTaskEnabled(target.ID, false); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < cap(e.sem); i++ {
		<-e.sem
	}
	e.runWG.Wait()
	select {
	case run := <-runner.runs:
		t.Fatalf("disabled worker-delayed execution=%+v", run)
	default:
	}
	deliveries, err := st.ListCompletionDeliveries()
	if err != nil || len(deliveries) != 1 || deliveries[0].State != domain.DeliveryResolved {
		t.Fatalf("deliveries=%+v err=%v", deliveries, err)
	}
}

func TestS144QueuedCompletionRechecksAncestorAndCommand(t *testing.T) {
	for _, change := range []string{"ancestor_disabled", "command_removed"} {
		t.Run(change, func(t *testing.T) {
			st, err := store.Open(":memory:")
			if err != nil {
				t.Fatal(err)
			}
			defer st.Close()
			source := setupTask(t, st, domain.OverlapQueueOne)
			target := setupTask(t, st, domain.OverlapQueueOne)
			parent := domain.Group{Name: "parent", Enabled: true}
			if err := st.CreateGroup(&parent); err != nil {
				t.Fatal(err)
			}
			child := domain.Group{Name: "child", ParentID: parent.ID, Enabled: true}
			if err := st.CreateGroup(&child); err != nil {
				t.Fatal(err)
			}
			target.GroupID = child.ID
			if err := st.UpdateTask(&target); err != nil {
				t.Fatal(err)
			}
			addChain(t, st, source, target, domain.CompletionOnSuccess)
			sourceRun := domain.Run{TaskID: source.ID, ScheduledFor: time.Date(2026, 11, 10, 0, 0, 0, 0, time.UTC), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
			if err := st.RecordRunAndCreateDeliveries(&sourceRun, ""); err != nil {
				t.Fatal(err)
			}
			runner := &completionRunner{outcomes: map[string]domain.RunOutcome{}, runs: make(chan domain.Run, 1)}
			e := newEngine(st, runner)
			e.running[target.ID] = 1
			e.drainCompletionDeliveries()
			if change == "ancestor_disabled" {
				if err := st.SetGroupEnabled(parent.ID, false); err != nil {
					t.Fatal(err)
				}
			} else {
				target.Command = ""
				if err := st.UpdateTask(&target); err != nil {
					t.Fatal(err)
				}
			}
			e.finish(target)
			e.runWG.Wait()
			select {
			case run := <-runner.runs:
				t.Fatalf("invalidated execution=%+v", run)
			default:
			}
			deliveries, err := st.ListCompletionDeliveries()
			if err != nil || len(deliveries) != 1 || deliveries[0].State != domain.DeliveryResolved || deliveries[0].Resolution == "" {
				t.Fatalf("delivery=%+v err=%v", deliveries, err)
			}
		})
	}
}
