package store

import (
	"errors"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func chainTask(t *testing.T, st *Store, name string) domain.Task {
	t.Helper()
	schedule := domain.Schedule{
		Kind:         domain.ScheduleEvent,
		TriggerID:    domain.StartupEventID,
		Expression:   "@reboot",
		HumanSummary: "At scheduler startup",
	}
	if err := st.CreateSchedule(&schedule); err != nil {
		t.Fatal(err)
	}
	task := domain.Task{
		Name:          name,
		Command:       "echo",
		Enabled:       true,
		Timezone:      "UTC",
		ScheduleID:    schedule.ID,
		OverlapPolicy: domain.OverlapQueueOne,
		CatchupPolicy: domain.CatchupOne,
		State:         domain.TaskActive,
	}
	if err := st.CreateTask(&task); err != nil {
		t.Fatal(err)
	}
	return task
}

func TestRecordRunAndCreateDeliveriesPersistsOutputTruncation(t *testing.T) {
	st := openMem(t)
	task := chainTask(t, st, "source")
	run := domain.Run{
		TaskID: task.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeFailure,
		Output: "bounded output", OutputTruncated: true, Trigger: domain.TriggerManual,
	}
	if err := st.RecordRunAndCreateDeliveries(&run, ""); err != nil {
		t.Fatal(err)
	}
	got, err := st.GetRun(run.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !got.OutputTruncated || got.Output != run.Output {
		t.Fatalf("stored run diagnostics=%+v, want truncation and retained output", got)
	}
}

func TestCompletionChainsCRUDAndCycleValidation(t *testing.T) {
	st := openMem(t)
	a := chainTask(t, st, "a")
	b := chainTask(t, st, "b")
	c := chainTask(t, st, "c")

	ab := domain.CompletionChain{SourceTaskID: a.ID, TargetTaskID: b.ID, OnOutcome: domain.CompletionOnSuccess}
	if err := st.CreateCompletionChain(&ab); err != nil {
		t.Fatal(err)
	}
	if err := st.CreateCompletionChain(&domain.CompletionChain{SourceTaskID: a.ID, TargetTaskID: b.ID, OnOutcome: domain.CompletionOnSuccess}); !errors.Is(err, ErrDuplicateChain) {
		t.Fatalf("duplicate error = %v, want ErrDuplicateChain", err)
	}
	bc := domain.CompletionChain{SourceTaskID: b.ID, TargetTaskID: c.ID, OnOutcome: domain.CompletionOnAny}
	if err := st.CreateCompletionChain(&bc); err != nil {
		t.Fatal(err)
	}
	if err := st.CreateCompletionChain(&domain.CompletionChain{SourceTaskID: c.ID, TargetTaskID: a.ID, OnOutcome: domain.CompletionOnFailure}); !errors.Is(err, ErrChainCycle) {
		t.Fatalf("indirect cycle error = %v, want ErrChainCycle", err)
	}

	got, err := st.GetCompletionChain(ab.ID)
	if err != nil || got.SourceTaskName != "a" || got.TargetTaskName != "b" {
		t.Fatalf("GetCompletionChain = %+v, err=%v", got, err)
	}
	got.OnOutcome = domain.CompletionOnFailure
	if err := st.UpdateCompletionChain(&got); err != nil {
		t.Fatal(err)
	}
	all, err := st.ListCompletionChains()
	if err != nil || len(all) != 2 {
		t.Fatalf("ListCompletionChains len=%d err=%v", len(all), err)
	}
	if err := st.DeleteCompletionChain(ab.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetCompletionChain(ab.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("deleted chain error = %v, want ErrNotFound", err)
	}
}

func TestCompletionChainsRejectIndirectCycleAcrossHundredTasks(t *testing.T) {
	st := openMem(t)
	tasks := make([]domain.Task, 100)
	for i := range tasks {
		tasks[i] = chainTask(t, st, "node-"+time.Unix(int64(i), 0).UTC().Format("150405"))
		if i == 0 {
			continue
		}
		chain := domain.CompletionChain{SourceTaskID: tasks[i-1].ID, TargetTaskID: tasks[i].ID, OnOutcome: domain.CompletionOnAny}
		if err := st.CreateCompletionChain(&chain); err != nil {
			t.Fatalf("create edge %d: %v", i, err)
		}
	}
	closing := domain.CompletionChain{SourceTaskID: tasks[99].ID, TargetTaskID: tasks[0].ID, OnOutcome: domain.CompletionOnSuccess}
	if err := st.CreateCompletionChain(&closing); !errors.Is(err, ErrChainCycle) {
		t.Fatalf("100-task closing edge error = %v, want ErrChainCycle", err)
	}
}

func TestCompletionDeliveryAtomicLifecycleAndCorrelation(t *testing.T) {
	st := openMem(t)
	source := chainTask(t, st, "source")
	target := chainTask(t, st, "target")
	chain := domain.CompletionChain{SourceTaskID: source.ID, TargetTaskID: target.ID, OnOutcome: domain.CompletionOnSuccess}
	if err := st.CreateCompletionChain(&chain); err != nil {
		t.Fatal(err)
	}

	sourceRun := domain.Run{
		ID:           "source-run",
		TaskID:       source.ID,
		ScheduledFor: time.Now().UTC(),
		Outcome:      domain.OutcomeSuccess,
		Trigger:      domain.TriggerManual,
	}
	if err := st.RecordRunAndCreateDeliveries(&sourceRun, ""); err != nil {
		t.Fatal(err)
	}
	deliveries, err := st.ListCompletionDeliveries()
	if err != nil || len(deliveries) != 1 || deliveries[0].State != domain.DeliveryPending {
		t.Fatalf("pending deliveries = %+v, err=%v", deliveries, err)
	}
	claimed, err := st.ClaimCompletionDeliveries(100)
	if err != nil || len(claimed) != 1 || claimed[0].Attempts != 1 {
		t.Fatalf("claimed deliveries = %+v, err=%v", claimed, err)
	}

	targetRun := domain.Run{
		ID:           "target-run",
		TaskID:       target.ID,
		ScheduledFor: time.Now().UTC(),
		Outcome:      domain.OutcomeSuccess,
		Trigger:      domain.TriggerCompletion,
		SourceTaskID: source.ID,
		SourceRunID:  sourceRun.ID,
	}
	if err := st.RecordRunAndCreateDeliveries(&targetRun, claimed[0].ID); err != nil {
		t.Fatal(err)
	}
	deliveries, err = st.ListCompletionDeliveries()
	if err != nil || len(deliveries) != 1 || deliveries[0].State != domain.DeliveryCompleted || deliveries[0].TargetRunID != targetRun.ID {
		t.Fatalf("completed deliveries = %+v, err=%v", deliveries, err)
	}
	runs, err := st.ListRuns(target.ID, 0)
	if err != nil || len(runs) != 1 || runs[0].SourceTaskID != source.ID || runs[0].SourceRunID != sourceRun.ID {
		t.Fatalf("correlated target runs = %+v, err=%v", runs, err)
	}
}

func TestDeletingChainPreservesPendingDeliveryForTerminalResolution(t *testing.T) {
	st := openMem(t)
	source := chainTask(t, st, "source")
	target := chainTask(t, st, "target")
	chain := domain.CompletionChain{SourceTaskID: source.ID, TargetTaskID: target.ID, OnOutcome: domain.CompletionOnSuccess}
	if err := st.CreateCompletionChain(&chain); err != nil {
		t.Fatal(err)
	}
	run := domain.Run{TaskID: source.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&run, ""); err != nil {
		t.Fatal(err)
	}
	if err := st.DeleteCompletionChain(chain.ID); err != nil {
		t.Fatal(err)
	}
	deliveries, err := st.ListCompletionDeliveries()
	if err != nil || len(deliveries) != 1 || deliveries[0].State != domain.DeliveryPending {
		t.Fatalf("delivery after chain deletion = %+v, err=%v", deliveries, err)
	}
}

func TestCompletionDeliveryRecoveryIsBoundedAndCompletedDoesNotReplay(t *testing.T) {
	st := openMem(t)
	source := chainTask(t, st, "source")
	target := chainTask(t, st, "target")
	chain := domain.CompletionChain{SourceTaskID: source.ID, TargetTaskID: target.ID, OnOutcome: domain.CompletionOnFailure}
	if err := st.CreateCompletionChain(&chain); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 101; i++ {
		run := domain.Run{TaskID: source.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeFailure, Trigger: domain.TriggerSchedule}
		if err := st.RecordRunAndCreateDeliveries(&run, ""); err != nil {
			t.Fatal(err)
		}
	}
	first, err := st.ClaimCompletionDeliveries(100)
	if err != nil || len(first) != 100 {
		t.Fatalf("first claim len=%d err=%v", len(first), err)
	}
	if recovered, err := st.RecoverCompletionDeliveries(); err != nil || recovered != 100 {
		t.Fatalf("recovered=%d err=%v", recovered, err)
	}
	replayed, err := st.ClaimCompletionDeliveries(100)
	if err != nil || len(replayed) != 100 {
		t.Fatalf("replay claim len=%d err=%v", len(replayed), err)
	}
	if err := st.ResolveCompletionDelivery(replayed[0].ID, "test terminal resolution"); err != nil {
		t.Fatal(err)
	}
	if recovered, err := st.RecoverCompletionDeliveries(); err != nil || recovered != 99 {
		t.Fatalf("second recovered=%d err=%v", recovered, err)
	}
}

func TestCompletedDeliveryDoesNotReplayAcrossHundredRecoveryPasses(t *testing.T) {
	st := openMem(t)
	source := chainTask(t, st, "source")
	target := chainTask(t, st, "target")
	chain := domain.CompletionChain{SourceTaskID: source.ID, TargetTaskID: target.ID, OnOutcome: domain.CompletionOnSuccess}
	if err := st.CreateCompletionChain(&chain); err != nil {
		t.Fatal(err)
	}
	sourceRun := domain.Run{TaskID: source.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&sourceRun, ""); err != nil {
		t.Fatal(err)
	}
	claimed, err := st.ClaimCompletionDeliveries(100)
	if err != nil || len(claimed) != 1 {
		t.Fatalf("claim=%+v err=%v", claimed, err)
	}
	targetRun := domain.Run{TaskID: target.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerCompletion, SourceTaskID: source.ID, SourceRunID: sourceRun.ID}
	if err := st.RecordRunAndCreateDeliveries(&targetRun, claimed[0].ID); err != nil {
		t.Fatal(err)
	}
	for restart := 0; restart < 100; restart++ {
		if recovered, err := st.RecoverCompletionDeliveries(); err != nil || recovered != 0 {
			t.Fatalf("restart %d recovered=%d err=%v", restart, recovered, err)
		}
		if replayed, err := st.ClaimCompletionDeliveries(100); err != nil || len(replayed) != 0 {
			t.Fatalf("restart %d replayed=%+v err=%v", restart, replayed, err)
		}
	}
}

func TestDuplicateSourceRunInsertionCannotDuplicateDelivery(t *testing.T) {
	st := openMem(t)
	source := chainTask(t, st, "source")
	target := chainTask(t, st, "target")
	chain := domain.CompletionChain{SourceTaskID: source.ID, TargetTaskID: target.ID, OnOutcome: domain.CompletionOnSuccess}
	if err := st.CreateCompletionChain(&chain); err != nil {
		t.Fatal(err)
	}
	run := domain.Run{ID: "immutable-source-run", TaskID: source.ID, ScheduledFor: time.Now().UTC(), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&run, ""); err != nil {
		t.Fatal(err)
	}
	for attempt := 0; attempt < 100; attempt++ {
		if err := st.RecordRunAndCreateDeliveries(&run, ""); err == nil {
			t.Fatalf("duplicate source run insertion %d unexpectedly succeeded", attempt)
		}
	}
	deliveries, err := st.ListCompletionDeliveries()
	if err != nil || len(deliveries) != 1 {
		t.Fatalf("deliveries=%+v err=%v", deliveries, err)
	}
}

func TestMigrationV9CreatesCompletionSchemaWithoutLegacyTriggers(t *testing.T) {
	st := openMem(t)
	for _, table := range []string{"completion_chains", "completion_deliveries"} {
		var count int
		if err := st.db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&count); err != nil || count != 1 {
			t.Fatalf("table %q count=%d err=%v", table, count, err)
		}
	}
	for _, table := range []string{"triggers", "dedup_ledger"} {
		var count int
		if err := st.db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&count); err != nil || count != 0 {
			t.Fatalf("legacy table %q count=%d err=%v", table, count, err)
		}
	}
}
