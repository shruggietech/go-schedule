package store

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func s145Fixture(t *testing.T, st *Store) (domain.ParallelJoin, domain.Task) {
	t.Helper()
	j, entry := joinFixture(t, st)
	recovery, cleanup := chainTask(t, st, "recover"), chainTask(t, st, "cleanup")
	j.RecoveryTaskID, j.CleanupTaskID = recovery.ID, cleanup.ID
	j.RecoveryOn = nil
	if err := st.UpdateParallelJoin(&j); err != nil {
		t.Fatal(err)
	}
	return j, entry
}
func s145Claim(t *testing.T, st *Store, role domain.CompletionDeliveryRole, count int) []domain.CompletionDelivery {
	t.Helper()
	ds, err := st.ClaimCompletionDeliveries(100)
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) != count {
		t.Fatalf("deliveries=%+v want role=%s count=%d", ds, role, count)
	}
	for _, d := range ds {
		if d.Role != role {
			t.Fatalf("role=%s want=%s", d.Role, role)
		}
	}
	return ds
}
func s145Branch(t *testing.T, j domain.ParallelJoinInstance, delivery string) domain.ParallelJoinInvocation {
	t.Helper()
	for _, b := range j.Branches {
		if b.DeliveryID == delivery {
			return b
		}
	}
	t.Fatalf("missing delivery %s", delivery)
	return domain.ParallelJoinInvocation{}
}
func s145Instance(t *testing.T, st *Store) domain.ParallelJoinInstance {
	t.Helper()
	is := joinInstances(t, st)
	if len(is) != 1 {
		t.Fatalf("instances=%v", is)
	}
	return is[0]
}
func TestS145FullConfigurationAndEmptyConditions(t *testing.T) {
	st := openMem(t)
	j, _ := s145Fixture(t, st)
	got, err := st.GetParallelJoinDefinitionPageContext(t.Context(), j.ID, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if got.RecoveryTaskID != j.RecoveryTaskID || got.CleanupTaskID != j.CleanupTaskID || !reflect.DeepEqual(got.RecoveryOn, []domain.JoinRecoveryCondition{domain.JoinRecoveryOnFailure}) {
		t.Fatalf("definition=%+v", got)
	}
	got.RecoveryOn = []domain.JoinRecoveryCondition{}
	if err := st.UpdateParallelJoin(&got); err != nil {
		t.Fatal(err)
	}
	got, err = st.GetParallelJoin(got.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.RecoveryOn == nil || len(got.RecoveryOn) != 0 {
		t.Fatalf("explicit empty conditions changed: %#v", got.RecoveryOn)
	}
}
func TestS145RecoveryAndCleanupWaitForExactStages(t *testing.T) {
	st := openMem(t)
	j, entry := s145Fixture(t, st)
	joinEntry(t, st, entry)
	ds := s145Claim(t, st, domain.DeliveryRoleJoinBranch, 2)
	first := joinFinish(t, st, ds[0], domain.OutcomeFailure)
	s145Claim(t, st, domain.DeliveryRoleJoinRecovery, 0)
	waiting := s145Instance(t, st)
	if waiting.SelectedRoute != domain.JoinRoutePending || s145Branch(t, waiting, ds[0].ID).RunID != first.ID || s145Branch(t, waiting, ds[0].ID).Settlement != domain.JoinSettlementFailure {
		t.Fatalf("waiting=%+v", waiting)
	}
	joinFinish(t, st, ds[1], domain.OutcomeSuccess)
	recovery := s145Claim(t, st, domain.DeliveryRoleJoinRecovery, 1)[0]
	if recovery.TargetTaskID != j.RecoveryTaskID {
		t.Fatalf("recovery=%+v", recovery)
	}
	ready, _, err := st.CompletionDeliveryReady(recovery.ID)
	if err != nil || !ready {
		t.Fatalf("ready=%v err=%v", ready, err)
	}
	s145Claim(t, st, domain.DeliveryRoleJoinCleanup, 0)
	run := joinFinish(t, st, recovery, domain.OutcomeFailure)
	cleanup := s145Claim(t, st, domain.DeliveryRoleJoinCleanup, 1)[0]
	if cleanup.TargetTaskID != j.CleanupTaskID {
		t.Fatalf("cleanup=%+v", cleanup)
	}
	middle := s145Instance(t, st)
	if middle.SelectedRoute != domain.JoinRouteRecovery || middle.RecoveryRunID != run.ID || middle.TargetDeliveryID != "" || middle.RecoverySettlement != domain.JoinSettlementFailure {
		t.Fatalf("middle=%+v", middle)
	}
	finalRun := joinFinish(t, st, cleanup, domain.OutcomeSuccess)
	final := s145Instance(t, st)
	if final.State != domain.JoinInstanceCompleted || final.CleanupRunID != finalRun.ID || final.CleanupSettlement != domain.JoinSettlementSuccess || final.CompletedAt == nil {
		t.Fatalf("final=%+v", final)
	}
	if err := st.RecordRunAndCreateDeliveries(&finalRun, cleanup.ID); err != nil {
		t.Fatal(err)
	}
	s145Claim(t, st, domain.DeliveryRoleJoinCleanup, 0)
}
func TestS145SuccessThenCleanupAfterNoStart(t *testing.T) {
	st := openMem(t)
	_, entry := s145Fixture(t, st)
	joinEntry(t, st, entry)
	ds := s145Claim(t, st, domain.DeliveryRoleJoinBranch, 2)
	for _, d := range ds {
		joinFinish(t, st, d, domain.OutcomeSuccess)
	}
	target := s145Claim(t, st, domain.DeliveryRoleJoinTarget, 1)[0]
	if err := st.ResolveCompletionDelivery(target.ID, "target disabled before native admission"); err != nil {
		t.Fatal(err)
	}
	cleanup := s145Claim(t, st, domain.DeliveryRoleJoinCleanup, 1)[0]
	got := s145Instance(t, st)
	if got.SelectedRoute != domain.JoinRouteSuccess || got.TargetSettlement != domain.JoinSettlementNotStarted || got.TargetRunID != "" || got.TargetReason != "target disabled before native admission" {
		t.Fatalf("instance=%+v", got)
	}
	if err := st.ResolveCompletionDelivery(cleanup.ID, "cleanup group disabled"); err != nil {
		t.Fatal(err)
	}
	got = s145Instance(t, st)
	if got.CleanupSettlement != domain.JoinSettlementNotStarted || got.CleanupRunID != "" || got.CleanupReason != "cleanup group disabled" || got.State != domain.JoinInstanceCompleted {
		t.Fatalf("final=%+v", got)
	}
}
func TestS145NoStartAndActualStartFailureRemainDistinctEvidence(t *testing.T) {
	for _, native := range []bool{false, true} {
		t.Run(map[bool]string{false: "no_run", true: "native_start_failure"}[native], func(t *testing.T) {
			st := openMem(t)
			j, entry := s145Fixture(t, st)
			j.RecoveryOn = []domain.JoinRecoveryCondition{domain.JoinRecoveryOnNotStarted}
			if err := st.UpdateParallelJoin(&j); err != nil {
				t.Fatal(err)
			}
			joinEntry(t, st, entry)
			ds := s145Claim(t, st, domain.DeliveryRoleJoinBranch, 2)
			var recorded domain.Run
			if native {
				recorded = domain.Run{TaskID: ds[0].TargetTaskID, Outcome: domain.OutcomeFailure, StartFailed: true, ScheduledFor: time.Now().UTC()}
				if err := st.RecordRunAndCreateDeliveries(&recorded, ds[0].ID); err != nil {
					t.Fatal(err)
				}
			} else if err := st.ResolveCompletionDelivery(ds[0].ID, "command unavailable before launch"); err != nil {
				t.Fatal(err)
			}
			s145Claim(t, st, domain.DeliveryRoleJoinRecovery, 0)
			joinFinish(t, st, ds[1], domain.OutcomeSuccess)
			s145Claim(t, st, domain.DeliveryRoleJoinRecovery, 1)
			got := s145Instance(t, st)
			b := s145Branch(t, got, ds[0].ID)
			if b.Settlement != domain.JoinSettlementNotStarted || b.SettledAt == nil || b.RunID != recorded.ID {
				t.Fatalf("branch=%+v recorded=%+v", b, recorded)
			}
			if native && b.Outcome != domain.OutcomeFailure {
				t.Fatalf("actual outcome lost: %+v", b)
			}
		})
	}
}
func TestS145NoSelectedFollowupStillCleans(t *testing.T) {
	st := openMem(t)
	j, entry := s145Fixture(t, st)
	j.RecoveryOn = []domain.JoinRecoveryCondition{}
	if err := st.UpdateParallelJoin(&j); err != nil {
		t.Fatal(err)
	}
	joinEntry(t, st, entry)
	ds := s145Claim(t, st, domain.DeliveryRoleJoinBranch, 2)
	joinFinish(t, st, ds[0], domain.OutcomeFailure)
	joinFinish(t, st, ds[1], domain.OutcomeSuccess)
	s145Claim(t, st, domain.DeliveryRoleJoinCleanup, 1)
	got := s145Instance(t, st)
	if got.SelectedRoute != domain.JoinRouteNone || got.RecoveryDeliveryID != "" || got.TargetDeliveryID != "" {
		t.Fatalf("instance=%+v", got)
	}
}
func TestS145FiniteStageReuseAndActualCascadeCycles(t *testing.T) {
	st := openMem(t)
	j, entry := s145Fixture(t, st)
	j.CleanupTaskID = j.TargetTaskID
	if err := st.UpdateParallelJoin(&j); err != nil {
		t.Fatalf("finite success/cleanup reuse rejected: %v", err)
	}
	c := domain.CompletionChain{SourceTaskID: j.CleanupTaskID, TargetTaskID: entry.ID, OnOutcome: domain.CompletionOnAny}
	if err := st.CreateCompletionChain(&c); err == nil {
		t.Fatal("recursive cleanup->entry cascade accepted")
	}
	j.CleanupTaskID = j.Branches[0].TaskID
	if err := st.UpdateParallelJoin(&j); err != nil {
		t.Fatalf("finite branch/cleanup reuse rejected: %v", err)
	}
}
func TestS145SnapshotReopenAndInvalidationNeverReviveCleanup(t *testing.T) {
	path := filepath.Join(t.TempDir(), "current.db")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	j, entry := s145Fixture(t, st)
	joinEntry(t, st, entry)
	ds := s145Claim(t, st, domain.DeliveryRoleJoinBranch, 2)
	joinFinish(t, st, ds[0], domain.OutcomeFailure)
	newRecovery := chainTask(t, st, "future recovery")
	j.RecoveryTaskID = newRecovery.ID
	if err := st.UpdateParallelJoin(&j); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	st, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if _, err := st.RecoverCompletionDeliveries(); err != nil {
		t.Fatal(err)
	}
	ds = s145Claim(t, st, domain.DeliveryRoleJoinBranch, 1)
	joinFinish(t, st, ds[0], domain.OutcomeSuccess)
	recovery := s145Claim(t, st, domain.DeliveryRoleJoinRecovery, 1)[0]
	if recovery.TargetTaskID == newRecovery.ID {
		t.Fatal("future edit changed immutable selected recovery")
	}
	if err := st.DeleteParallelJoin(j.ID); err != nil {
		t.Fatal(err)
	}
	joinFinish(t, st, recovery, domain.OutcomeSuccess)
	s145Claim(t, st, domain.DeliveryRoleJoinCleanup, 0)
	got := s145Instance(t, st)
	if got.State != domain.JoinInstanceResolved || got.RecoverySettlement != domain.JoinSettlementSuccess {
		t.Fatalf("instance=%+v", got)
	}
}

func TestS145RecoveryCauseSelection(t *testing.T) {
	cases := []struct {
		name        string
		outcome     domain.RunOutcome
		startFailed bool
		condition   domain.JoinRecoveryCondition
		settlement  domain.JoinSettlement
	}{
		{"failure", domain.OutcomeFailure, false, domain.JoinRecoveryOnFailure, domain.JoinSettlementFailure},
		{"canceled", domain.OutcomeCanceled, false, domain.JoinRecoveryOnCanceled, domain.JoinSettlementCanceled},
		{"skipped", domain.OutcomeSkipped, false, domain.JoinRecoveryOnSkipped, domain.JoinSettlementSkipped},
		{"cancellation_cleanup_failure", domain.OutcomeCancellationFailed, false, domain.JoinRecoveryOnFailure, domain.JoinSettlementFailure},
		{"native_start_failure", domain.OutcomeFailure, true, domain.JoinRecoveryOnNotStarted, domain.JoinSettlementNotStarted},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			st := openMem(t)
			j, entry := s145Fixture(t, st)
			j.RecoveryOn = []domain.JoinRecoveryCondition{tc.condition}
			if err := st.UpdateParallelJoin(&j); err != nil {
				t.Fatal(err)
			}
			joinEntry(t, st, entry)
			ds := s145Claim(t, st, domain.DeliveryRoleJoinBranch, 2)
			r := domain.Run{TaskID: ds[0].TargetTaskID, Outcome: tc.outcome, StartFailed: tc.startFailed, ScheduledFor: time.Now().UTC()}
			if err := st.RecordRunAndCreateDeliveries(&r, ds[0].ID); err != nil {
				t.Fatal(err)
			}
			s145Claim(t, st, domain.DeliveryRoleJoinRecovery, 0)
			joinFinish(t, st, ds[1], domain.OutcomeSuccess)
			s145Claim(t, st, domain.DeliveryRoleJoinRecovery, 1)
			got := s145Instance(t, st)
			b := s145Branch(t, got, ds[0].ID)
			if b.Settlement != tc.settlement || b.Outcome != tc.outcome || b.RunID != r.ID {
				t.Fatalf("evidence=%+v", b)
			}
		})
	}
}
func TestS145ConcurrentOverlappingEntriesKeepIndependentRecovery(t *testing.T) {
	st := openMem(t)
	_, entry := s145Fixture(t, st)
	a, b := joinEntry(t, st, entry), joinEntry(t, st, entry)
	ds := s145Claim(t, st, domain.DeliveryRoleJoinBranch, 4)
	var wg sync.WaitGroup
	errs := make(chan error, len(ds))
	for _, d := range ds {
		wg.Add(1)
		go func(d domain.CompletionDelivery) {
			defer wg.Done()
			r := domain.Run{TaskID: d.TargetTaskID, Outcome: domain.OutcomeFailure, ScheduledFor: time.Now().UTC()}
			errs <- st.RecordRunAndCreateDeliveries(&r, d.ID)
		}(d)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	rs := s145Claim(t, st, domain.DeliveryRoleJoinRecovery, 2)
	if rs[0].JoinInstanceID == rs[1].JoinInstanceID || rs[0].SourceRunID == rs[1].SourceRunID {
		t.Fatalf("cross-instance recovery=%+v", rs)
	}
	seen := map[string]bool{a.ID: false, b.ID: false}
	for _, r := range rs {
		seen[r.SourceRunID] = true
		joinFinish(t, st, r, domain.OutcomeSuccess)
	}
	if !seen[a.ID] || !seen[b.ID] {
		t.Fatalf("source identity=%v", seen)
	}
	cs := s145Claim(t, st, domain.DeliveryRoleJoinCleanup, 2)
	if cs[0].JoinInstanceID == cs[1].JoinInstanceID {
		t.Fatalf("cleanup=%+v", cs)
	}
}
func TestS145FiniteTaskReuseExecutesDistinctStages(t *testing.T) {
	st := openMem(t)
	j, entry := s145Fixture(t, st)
	j.TargetTaskID = j.Branches[0].TaskID
	j.RecoveryTaskID = j.TargetTaskID
	j.CleanupTaskID = j.TargetTaskID
	if err := st.UpdateParallelJoin(&j); err != nil {
		t.Fatal(err)
	}
	joinEntry(t, st, entry)
	ds := s145Claim(t, st, domain.DeliveryRoleJoinBranch, 2)
	var ids []string
	for _, d := range ds {
		ids = append(ids, joinFinish(t, st, d, domain.OutcomeSuccess).ID)
	}
	target := s145Claim(t, st, domain.DeliveryRoleJoinTarget, 1)[0]
	ids = append(ids, joinFinish(t, st, target, domain.OutcomeSuccess).ID)
	cleanup := s145Claim(t, st, domain.DeliveryRoleJoinCleanup, 1)[0]
	ids = append(ids, joinFinish(t, st, cleanup, domain.OutcomeSuccess).ID)
	unique := map[string]bool{}
	for _, id := range ids {
		unique[id] = true
	}
	if len(unique) != 4 {
		t.Fatalf("reused native identity=%v", ids)
	}
}
func TestS145RecoveryConfigurationValidationAndAtomicReplacement(t *testing.T) {
	st := openMem(t)
	j, _ := s145Fixture(t, st)
	before, err := st.GetParallelJoin(j.ID)
	if err != nil {
		t.Fatal(err)
	}
	bad := j
	bad.RecoveryOn = []domain.JoinRecoveryCondition{"other"}
	if err := st.UpdateParallelJoin(&bad); !errors.Is(err, ErrInvalidJoin) {
		t.Fatalf("unknown cause=%v", err)
	}
	bad = j
	bad.RecoveryTaskID = ""
	if err := st.UpdateParallelJoin(&bad); !errors.Is(err, ErrInvalidJoin) {
		t.Fatalf("unbound causes=%v", err)
	}
	for _, field := range []string{"recovery", "cleanup"} {
		bad = j
		if field == "recovery" {
			bad.RecoveryTaskID = "absent"
		} else {
			bad.CleanupTaskID = "absent"
		}
		if err := st.UpdateParallelJoin(&bad); !errors.Is(err, ErrInvalidJoin) {
			t.Fatalf("missing %s=%v", field, err)
		}
	}
	bad = j
	bad.CleanupTaskID = j.EntryTaskID
	if err := st.ReplaceAutomationRelationships(nil, []domain.ParallelJoin{bad}); !errors.Is(err, ErrChainCycle) {
		t.Fatalf("atomic cycle=%v", err)
	}
	after, err := st.GetParallelJoin(j.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("failed update changed definition: before=%+v after=%+v", before, after)
	}
	stale := before.UpdatedAt.Add(-time.Second)
	if err := st.UpdateParallelJoinIfUnchanged(&j, &stale); !errors.Is(err, ErrStaleJoin) {
		t.Fatalf("stale=%v", err)
	}
	j.RecoveryOn = []domain.JoinRecoveryCondition{domain.JoinRecoveryOnSkipped, domain.JoinRecoveryOnFailure, domain.JoinRecoveryOnSkipped}
	if err := st.UpdateParallelJoin(&j); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(j.RecoveryOn, []domain.JoinRecoveryCondition{domain.JoinRecoveryOnFailure, domain.JoinRecoveryOnSkipped}) {
		t.Fatalf("canonical conditions=%v", j.RecoveryOn)
	}
}
func TestS145NewDestinationDeletionKeepsLateNativeEvidence(t *testing.T) {
	for _, field := range []string{"recovery", "cleanup"} {
		t.Run(field, func(t *testing.T) {
			st := openMem(t)
			j, entry := s145Fixture(t, st)
			joinEntry(t, st, entry)
			ds := s145Claim(t, st, domain.DeliveryRoleJoinBranch, 2)
			id := j.RecoveryTaskID
			if field == "cleanup" {
				id = j.CleanupTaskID
			}
			if err := st.DeleteTask(id); err != nil {
				t.Fatal(err)
			}
			for _, d := range ds {
				joinFinish(t, st, d, domain.OutcomeFailure)
			}
			s145Claim(t, st, domain.DeliveryRoleJoinRecovery, 0)
			got := s145Instance(t, st)
			if got.State != domain.JoinInstanceResolved || got.RecoveryDeliveryID != "" || got.CleanupDeliveryID != "" {
				t.Fatalf("invalidated=%+v", got)
			}
			for _, b := range got.Branches {
				if b.Settlement != domain.JoinSettlementFailure || b.RunID == "" {
					t.Fatalf("late=%+v", b)
				}
			}
			if _, err := st.GetParallelJoin(j.ID); !errors.Is(err, ErrNotFound) {
				t.Fatalf("invalid definition survived: %v", err)
			}
		})
	}
}
func TestS145SnapshotActivationPersistsUntilCleanupSettles(t *testing.T) {
	st := openMem(t)
	j, entry := s145Fixture(t, st)
	for _, id := range []string{j.RecoveryTaskID, j.CleanupTaskID} {
		task, err := st.GetTask(id)
		if err != nil {
			t.Fatal(err)
		}
		task.ScheduleID = ""
		task.Enabled = true
		if err := st.UpdateTask(&task); err != nil {
			t.Fatal(err)
		}
	}
	joinEntry(t, st, entry)
	ds := s145Claim(t, st, domain.DeliveryRoleJoinBranch, 2)
	oldRecovery, oldCleanup := j.RecoveryTaskID, j.CleanupTaskID
	replacementR, replacementC := chainTask(t, st, "new recovery"), chainTask(t, st, "new cleanup")
	j.RecoveryTaskID, j.CleanupTaskID = replacementR.ID, replacementC.ID
	if err := st.UpdateParallelJoin(&j); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{oldRecovery, oldCleanup} {
		has, err := st.TaskHasJoinSource(id)
		if err != nil || !has {
			t.Fatalf("old snapshot source=%s %v %v", id, has, err)
		}
		task, err := st.GetTask(id)
		if err != nil || !task.Enabled {
			t.Fatalf("old source prematurely disabled %s %+v %v", id, task, err)
		}
	}
	for _, d := range ds {
		joinFinish(t, st, d, domain.OutcomeFailure)
	}
	r := s145Claim(t, st, domain.DeliveryRoleJoinRecovery, 1)[0]
	if r.TargetTaskID != oldRecovery {
		t.Fatal("snapshot recovery changed")
	}
	joinFinish(t, st, r, domain.OutcomeSuccess)
	c := s145Claim(t, st, domain.DeliveryRoleJoinCleanup, 1)[0]
	if c.TargetTaskID != oldCleanup {
		t.Fatal("snapshot cleanup changed")
	}
	joinFinish(t, st, c, domain.OutcomeSuccess)
	for _, id := range []string{oldRecovery, oldCleanup} {
		has, err := st.TaskHasJoinSource(id)
		if err != nil || has {
			t.Fatalf("finished old source=%s %v %v", id, has, err)
		}
	}
}
func TestS145NewObservationTextRemainsBounded(t *testing.T) {
	st := openMem(t)
	j, entry := s145Fixture(t, st)
	for _, id := range []string{j.RecoveryTaskID, j.CleanupTaskID} {
		task, err := st.GetTask(id)
		if err != nil {
			t.Fatal(err)
		}
		task.Name = strings.Repeat("界", 10000)
		if err := st.UpdateTask(&task); err != nil {
			t.Fatal(err)
		}
	}
	j, err := st.GetParallelJoinPage(j.ID, 0, 1, 12)
	if err != nil {
		t.Fatal(err)
	}
	if len(j.RecoveryTaskName) > 12 || len(j.CleanupTaskName) > 12 || len(j.TruncatedFields) < 2 {
		t.Fatalf("unbounded definition=%+v", j)
	}
	joinEntry(t, st, entry)
	ds := s145Claim(t, st, domain.DeliveryRoleJoinBranch, 2)
	for _, d := range ds {
		joinFinish(t, st, d, domain.OutcomeFailure)
	}
	r := s145Claim(t, st, domain.DeliveryRoleJoinRecovery, 1)[0]
	reason := strings.Repeat("界", 10000)
	if err := st.ResolveCompletionDelivery(r.ID, reason); err != nil {
		t.Fatal(err)
	}
	c := s145Claim(t, st, domain.DeliveryRoleJoinCleanup, 1)[0]
	if err := st.ResolveCompletionDelivery(c.ID, reason); err != nil {
		t.Fatal(err)
	}
	full := s145Instance(t, st)
	bounded, err := st.GetParallelJoinInstancePage(full.ID, 0, 1, 12)
	if err != nil {
		t.Fatal(err)
	}
	if full.RecoveryReason != reason || full.CleanupReason != reason || len(bounded.RecoveryReason) > 12 || len(bounded.CleanupReason) > 12 || len(bounded.RecoveryTaskName) > 12 || len(bounded.CleanupTaskName) > 12 {
		t.Fatalf("bounded=%+v", bounded)
	}
}
func TestS145RecoveryCreationRollbackIsAtomic(t *testing.T) {
	st := openMem(t)
	_, entry := s145Fixture(t, st)
	joinEntry(t, st, entry)
	ds := s145Claim(t, st, domain.DeliveryRoleJoinBranch, 2)
	joinFinish(t, st, ds[0], domain.OutcomeFailure)
	if _, err := st.db.Exec("CREATE TRIGGER reject_recovery BEFORE INSERT ON completion_deliveries WHEN NEW.role='join_recovery' BEGIN SELECT RAISE(ABORT,'recovery insert fault'); END"); err != nil {
		t.Fatal(err)
	}
	r := domain.Run{TaskID: ds[1].TargetTaskID, Outcome: domain.OutcomeSuccess, ScheduledFor: time.Now().UTC()}
	if err := st.RecordRunAndCreateDeliveries(&r, ds[1].ID); err == nil {
		t.Fatal("injected creation fault ignored")
	}
	got := s145Instance(t, st)
	if got.SelectedRoute != domain.JoinRoutePending || s145Branch(t, got, ds[1].ID).Settlement != domain.JoinSettlementPending {
		t.Fatalf("partial decision persisted=%+v", got)
	}
	if _, err := st.db.Exec("DROP TRIGGER reject_recovery"); err != nil {
		t.Fatal(err)
	}
	if err := st.RecordRunAndCreateDeliveries(&r, ds[1].ID); err != nil {
		t.Fatal(err)
	}
	s145Claim(t, st, domain.DeliveryRoleJoinRecovery, 1)
}

func TestS145LegacyBlockedReasonRetainsExactCause(t *testing.T) {
	st := openMem(t)
	_, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	ds := s145Claim(t, st, domain.DeliveryRoleJoinBranch, 2)
	joinFinish(t, st, ds[0], domain.OutcomeFailure)
	before := s145Instance(t, st).Reason
	joinFinish(t, st, ds[1], domain.OutcomeSuccess)
	after := s145Instance(t, st)
	if before == "" || after.Reason != before {
		t.Fatalf("legacy blocked reason changed: before=%q after=%q", before, after.Reason)
	}
}
