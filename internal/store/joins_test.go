package store

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func joinFixture(t *testing.T, st *Store) (domain.ParallelJoin, domain.Task) {
	t.Helper()
	entry := chainTask(t, st, "prepare")
	branch := chainTask(t, st, "export")
	target := chainTask(t, st, "publish")
	j := domain.ParallelJoin{Name: "release", EntryTaskID: entry.ID, TargetTaskID: target.ID, Branches: []domain.ParallelJoinBranch{{Name: "first", TaskID: branch.ID}, {Name: "second", TaskID: branch.ID}}}
	if err := st.CreateParallelJoin(&j); err != nil {
		t.Fatal(err)
	}
	return j, entry
}
func joinEntry(t *testing.T, st *Store, entry domain.Task) domain.Run {
	t.Helper()
	r := domain.Run{TaskID: entry.ID, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual, ScheduledFor: time.Now().UTC()}
	if err := st.RecordRunAndCreateDeliveries(&r, ""); err != nil {
		t.Fatal(err)
	}
	return r
}
func joinFinish(t *testing.T, st *Store, d domain.CompletionDelivery, outcome domain.RunOutcome) domain.Run {
	t.Helper()
	r := domain.Run{TaskID: d.TargetTaskID, Outcome: outcome, Trigger: domain.TriggerCompletion, ScheduledFor: time.Now().UTC()}
	if err := st.RecordRunAndCreateDeliveries(&r, d.ID); err != nil {
		t.Fatal(err)
	}
	return r
}
func joinInstances(t *testing.T, st *Store) []domain.ParallelJoinInstance {
	t.Helper()
	all, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{})
	if err != nil {
		t.Fatal(err)
	}
	return all
}

func TestS144ExactDuplicateBranchesAndIdempotentEntry(t *testing.T) {
	st := openMem(t)
	j, entry := joinFixture(t, st)
	r := joinEntry(t, st, entry)
	if err := st.RecordRunAndCreateDeliveries(&r, ""); err != nil {
		t.Fatal(err)
	}
	instances := joinInstances(t, st)
	if len(instances) != 1 || len(instances[0].Branches) != 2 {
		t.Fatalf("instances=%+v", instances)
	}
	a, b := instances[0].Branches[0], instances[0].Branches[1]
	if a.ID == b.ID || a.DeliveryID == b.DeliveryID || a.TaskID != b.TaskID || a.BranchID != j.Branches[0].ID {
		t.Fatalf("identities %v %v", a, b)
	}
	deliveries, err := st.ClaimCompletionDeliveries(10)
	if err != nil || len(deliveries) != 2 {
		t.Fatalf("claim=%+v %v", deliveries, err)
	}
	first := joinFinish(t, st, deliveries[0], domain.OutcomeSuccess)
	if first.DeliveryID != deliveries[0].ID || first.JoinInstanceID != instances[0].ID || first.JoinBranchID != deliveries[0].JoinBranchID || first.JoinRole != "branch" {
		t.Fatalf("run correlation=%+v", first)
	}
	if all, err := st.ClaimCompletionDeliveries(10); err != nil || len(all) != 0 {
		t.Fatalf("premature target=%+v %v", all, err)
	}
	second := joinFinish(t, st, deliveries[1], domain.OutcomeSuccess)
	if err := st.RecordRunAndCreateDeliveries(&second, deliveries[1].ID); err != nil {
		t.Fatal(err)
	}
	targets, err := st.ClaimCompletionDeliveries(10)
	if err != nil || len(targets) != 1 || targets[0].Role != domain.DeliveryRoleJoinTarget || targets[0].SourceRunID != r.ID {
		t.Fatalf("target=%+v %v", targets, err)
	}
	joinFinish(t, st, targets[0], domain.OutcomeFailure)
	got, _ := st.GetParallelJoinInstance(instances[0].ID)
	if got.State != domain.JoinInstanceCompleted || got.TargetOutcome != domain.OutcomeFailure {
		t.Fatalf("instance=%+v", got)
	}
	conflict := r
	conflict.Outcome = domain.OutcomeFailure
	if err := st.RecordRunAndCreateDeliveries(&conflict, ""); err == nil {
		t.Fatal("conflicting run identity accepted")
	}
}
func TestS144InterleavingAndConcurrentLastBranch(t *testing.T) {
	st := openMem(t)
	_, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	joinEntry(t, st, entry)
	ds, err := st.ClaimCompletionDeliveries(10)
	if err != nil || len(ds) != 4 {
		t.Fatalf("deliveries=%+v %v", ds, err)
	}
	grouped := map[string][]domain.CompletionDelivery{}
	for _, d := range ds {
		grouped[d.JoinInstanceID] = append(grouped[d.JoinInstanceID], d)
	}
	var pairs [][]domain.CompletionDelivery
	for _, pair := range grouped {
		pairs = append(pairs, pair)
	}
	joinFinish(t, st, pairs[0][0], domain.OutcomeSuccess)
	joinFinish(t, st, pairs[1][1], domain.OutcomeSuccess)
	if ds, _ := st.ClaimCompletionDeliveries(10); len(ds) != 0 {
		t.Fatal("cross-instance target")
	}
	var wg sync.WaitGroup
	errc := make(chan error, 2)
	for _, pair := range pairs {
		d := pair[0]
		if d.ID == pairs[0][0].ID {
			d = pair[1]
		}
		wg.Add(1)
		go func(d domain.CompletionDelivery) {
			defer wg.Done()
			r := domain.Run{TaskID: d.TargetTaskID, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerCompletion, ScheduledFor: time.Now().UTC()}
			errc <- st.RecordRunAndCreateDeliveries(&r, d.ID)
		}(d)
	}
	wg.Wait()
	close(errc)
	for err := range errc {
		if err != nil {
			t.Fatal(err)
		}
	}
	if ds, _ := st.ClaimCompletionDeliveries(10); len(ds) != 2 {
		t.Fatalf("targets=%+v", ds)
	}
}
func TestS144TerminalBranchesNeverReleaseTarget(t *testing.T) {
	for _, outcome := range []domain.RunOutcome{domain.OutcomeFailure, domain.OutcomeSkipped, domain.OutcomeCanceled, domain.OutcomeCancellationFailed} {
		t.Run(string(outcome), func(t *testing.T) {
			st := openMem(t)
			_, entry := joinFixture(t, st)
			joinEntry(t, st, entry)
			ds, _ := st.ClaimCompletionDeliveries(10)
			joinFinish(t, st, ds[0], outcome)
			joinFinish(t, st, ds[1], domain.OutcomeSuccess)
			got := joinInstances(t, st)[0]
			if got.State != domain.JoinInstanceBlocked || got.Reason == "" {
				t.Fatalf("instance=%+v", got)
			}
			if ds, _ := st.ClaimCompletionDeliveries(10); len(ds) != 0 {
				t.Fatal("blocked target launched")
			}
		})
	}
}
func TestS144SnapshotEditDeleteAndLateResult(t *testing.T) {
	st := openMem(t)
	j, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	before := joinInstances(t, st)[0]
	ds, _ := st.ClaimCompletionDeliveries(10)
	old := j.UpdatedAt
	j.Name = "new"
	j.Branches[0].Name = "renamed"
	j.TargetTaskID = chainTask(t, st, "new target").ID
	if err := st.UpdateParallelJoinIfUnchanged(&j, &old); err != nil {
		t.Fatal(err)
	}
	if err := st.UpdateParallelJoinIfUnchanged(&j, &old); !errors.Is(err, ErrStaleJoin) {
		t.Fatalf("stale=%v", err)
	}
	got, _ := st.GetParallelJoinInstance(before.ID)
	if got.DefinitionName != before.DefinitionName || got.TargetTaskID != before.TargetTaskID || got.Branches[0].Name != before.Branches[0].Name {
		t.Fatalf("retargeted=%+v", got)
	}
	if err := st.DeleteParallelJoin(j.ID); err != nil {
		t.Fatal(err)
	}
	ready, reason, err := st.CompletionDeliveryReady(ds[0].ID)
	if err != nil || ready || reason == "" {
		t.Fatalf("ready=%v %s %v", ready, reason, err)
	}
	joinFinish(t, st, ds[0], domain.OutcomeSuccess)
	got, _ = st.GetParallelJoinInstance(before.ID)
	if got.State != domain.JoinInstanceResolved || !joinHasResult(got, ds[0].JoinBranchID, domain.OutcomeSuccess) {
		t.Fatalf("late evidence=%+v", got)
	}
	if all, _ := st.ClaimCompletionDeliveries(10); len(all) != 0 {
		t.Fatal("invalidated work replayed")
	}
}
func TestS144CurrentReopenTaskDeletionAndPaging(t *testing.T) {
	path := filepath.Join(t.TempDir(), "joins.db")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	j, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	ds, _ := st.ClaimCompletionDeliveries(10)
	joinFinish(t, st, ds[0], domain.OutcomeSuccess)
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	st, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := st.Close(); err != nil {
			t.Error(err)
		}
	}()
	if n, err := st.RecoverCompletionDeliveries(); err != nil || n != 1 {
		t.Fatalf("recover=%d %v", n, err)
	}
	got, err := st.GetParallelJoinPage(j.ID, 1, 1, 2)
	if err != nil || got.BranchCount != 2 || got.BranchOffset != 1 || len(got.Branches) != 1 || len(got.TruncatedFields) == 0 {
		t.Fatalf("page=%+v %v", got, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := st.GetParallelJoinPageContext(ctx, j.ID, 0, 1, 10); !errors.Is(err, context.Canceled) {
		t.Fatalf("context=%v", err)
	}
	inst := joinInstances(t, st)[0]
	if err := st.DeleteTask(j.Branches[1].TaskID); err != nil {
		t.Fatal(err)
	}
	after, _ := st.GetParallelJoinInstance(inst.ID)
	if after.State != domain.JoinInstanceResolved || !joinHasResult(after, ds[0].JoinBranchID, domain.OutcomeSuccess) || after.Branches[0].TaskID == "" {
		t.Fatalf("deleted attribution=%+v", after)
	}
}
func TestS144CombinedGraphAndAtomicReplacement(t *testing.T) {
	st := openMem(t)
	j, _ := joinFixture(t, st)
	c := domain.CompletionChain{SourceTaskID: j.TargetTaskID, TargetTaskID: j.EntryTaskID, OnOutcome: domain.CompletionOnSuccess}
	if err := st.CreateCompletionChain(&c); !errors.Is(err, ErrChainCycle) {
		t.Fatalf("cross cycle=%v", err)
	}
	j.Branches[1].Name = j.Branches[0].Name
	if err := st.UpdateParallelJoin(&j); err == nil {
		t.Fatal("duplicate branch name accepted")
	}
	j.Branches[1].Name = "second"
	oldTarget := j.TargetTaskID
	next := chainTask(t, st, "next")
	c = domain.CompletionChain{SourceTaskID: oldTarget, TargetTaskID: next.ID, OnOutcome: domain.CompletionOnSuccess}
	if err := st.CreateCompletionChain(&c); err != nil {
		t.Fatal(err)
	}
	c.SourceTaskID = next.ID
	c.TargetTaskID = j.EntryTaskID
	j.TargetTaskID = next.ID
	if err := st.ReplaceAutomationRelationships([]domain.CompletionChain{c}, []domain.ParallelJoin{j}); !errors.Is(err, ErrChainCycle) {
		t.Fatalf("final mixed cycle=%v", err)
	}
	got, _ := st.GetParallelJoin(j.ID)
	if got.TargetTaskID != oldTarget {
		t.Fatal("failed replacement partially committed")
	}
}

func joinHasResult(j domain.ParallelJoinInstance, id string, outcome domain.RunOutcome) bool {
	for _, b := range j.Branches {
		if b.ID == id {
			return b.Outcome == outcome
		}
	}
	return false
}

func TestS144SnapshotReadinessAndQueueResolution(t *testing.T) {
	st := openMem(t)
	entry := chainTask(t, st, "entry")
	oldBranch := domain.Task{Name: "old branch", Command: "echo", Timezone: "UTC", State: domain.TaskActive}
	oldTarget := domain.Task{Name: "old target", Command: "echo", Timezone: "UTC", State: domain.TaskActive}
	if err := st.CreateTask(&oldBranch); err != nil {
		t.Fatal(err)
	}
	if err := st.CreateTask(&oldTarget); err != nil {
		t.Fatal(err)
	}
	j := domain.ParallelJoin{Name: "snapshot", EntryTaskID: entry.ID, TargetTaskID: oldTarget.ID, Branches: []domain.ParallelJoinBranch{{Name: "old", TaskID: oldBranch.ID}}}
	if err := st.CreateParallelJoin(&j); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{oldBranch.ID, oldTarget.ID} {
		if err := st.SetTaskEnabled(id, true); err != nil {
			t.Fatal(err)
		}
	}
	joinEntry(t, st, entry)
	ds, _ := st.ClaimCompletionDeliveries(10)
	newBranch := chainTask(t, st, "new branch")
	newTarget := chainTask(t, st, "new target")
	j.Branches = []domain.ParallelJoinBranch{{Name: "replacement", TaskID: newBranch.ID}}
	j.BranchCount = 1
	j.TargetTaskID = newTarget.ID
	if err := st.UpdateParallelJoin(&j); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{oldBranch.ID, oldTarget.ID} {
		task, err := st.GetTask(id)
		if err != nil || !task.Enabled {
			t.Fatalf("snapshot source disabled %s: %+v %v", id, task, err)
		}
	}
	if ok, reason, err := st.CompletionDeliveryReady(ds[0].ID); err != nil || !ok {
		t.Fatalf("edited snapshot invalid: %v %s %v", ok, reason, err)
	}
	if err := st.ResolveCompletionDelivery(ds[0].ID, "collapsed by queue-one overlap policy"); err != nil {
		t.Fatal(err)
	}
	if err := st.ResolveCompletionDelivery(ds[0].ID, "repeat"); err != nil {
		t.Fatal(err)
	}
	got := joinInstances(t, st)[0]
	if got.State != domain.JoinInstanceBlocked || got.Branches[0].DeliveryState != domain.DeliveryResolved || got.Reason == "" {
		t.Fatalf("collapsed=%+v", got)
	}
	for _, id := range []string{oldBranch.ID, oldTarget.ID} {
		task, _ := st.GetTask(id)
		if task.Enabled {
			t.Fatalf("finished obsolete source stayed enabled: %s", id)
		}
	}
	obs, err := st.ListTaskObservations("", "", false, 0, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, o := range obs {
		if o.Task.ID == newBranch.ID {
			found = o.HasJoin
		}
	}
	if !found {
		t.Fatal("join source absent from bounded readiness observation")
	}
}
func TestS144SameInstanceConcurrentLastContributions(t *testing.T) {
	st := openMem(t)
	_, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	ds, _ := st.ClaimCompletionDeliveries(10)
	var wg sync.WaitGroup
	errs := make(chan error, len(ds))
	for _, d := range ds {
		wg.Add(1)
		go func(d domain.CompletionDelivery) {
			defer wg.Done()
			r := domain.Run{TaskID: d.TargetTaskID, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerCompletion, ScheduledFor: time.Now().UTC()}
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
	targets, err := st.ClaimCompletionDeliveries(10)
	if err != nil || len(targets) != 1 {
		t.Fatalf("unique target=%+v %v", targets, err)
	}
	if err := st.ResolveCompletionDelivery(targets[0].ID, "target disabled"); err != nil {
		t.Fatal(err)
	}
	got := joinInstances(t, st)[0]
	if got.State != domain.JoinInstanceResolved {
		t.Fatalf("target resolution=%+v", got)
	}
}
func TestS144MixedReplacementSupportsCreatesAndValidGraphSwap(t *testing.T) {
	st := openMem(t)
	j, _ := joinFixture(t, st)
	next := chainTask(t, st, "next")
	c := domain.CompletionChain{SourceTaskID: j.TargetTaskID, TargetTaskID: next.ID, OnOutcome: domain.CompletionOnSuccess}
	if err := st.CreateCompletionChain(&c); err != nil {
		t.Fatal(err)
	}
	oldEntry := j.EntryTaskID
	j.EntryTaskID = next.ID
	c.SourceTaskID = oldEntry
	c.TargetTaskID = next.ID
	if err := st.ReplaceAutomationRelationships([]domain.CompletionChain{c}, []domain.ParallelJoin{j}); err != nil {
		t.Fatalf("valid final graph swap=%v", err)
	}
	added := domain.ParallelJoin{ID: "portable-created-join", Name: "created", EntryTaskID: oldEntry, TargetTaskID: j.TargetTaskID, Branches: []domain.ParallelJoinBranch{{ID: "stable-slot", Name: "slot", TaskID: next.ID}}}
	if err := st.ReplaceAutomationRelationships(nil, []domain.ParallelJoin{added}); err != nil {
		t.Fatal(err)
	}
	got, err := st.GetParallelJoin(added.ID)
	if err != nil || got.Branches[0].ID != "stable-slot" {
		t.Fatalf("upsert=%+v %v", got, err)
	}
	if err := st.BindPortableID("join", added.ID, "portable-identity"); err != nil {
		t.Fatal(err)
	}
	if id, err := st.ObjectIDForPortableID("join", "portable-identity"); err != nil || id != added.ID {
		t.Fatalf("portable=%s %v", id, err)
	}
}
func TestS144PagedReadsCancelWhilePoolOccupied(t *testing.T) {
	st := openMem(t)
	j, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	instance := joinInstances(t, st)[0]
	conn, err := st.db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := st.ListParallelJoinsPageContext(ctx, 0, 1, 5, 0, 1); !errors.Is(err, context.Canceled) {
		t.Fatalf("list=%v", err)
	}
	if _, err := st.GetParallelJoinPageContext(ctx, j.ID, 0, 1, 5); !errors.Is(err, context.Canceled) {
		t.Fatalf("exact=%v", err)
	}
	if _, err := st.GetParallelJoinInstancePageContext(ctx, instance.ID, 0, 1, 5); !errors.Is(err, context.Canceled) {
		t.Fatalf("instance=%v", err)
	}
	if _, err := st.ListParallelJoinInstancesContext(ctx, domain.ParallelJoinInstanceFilter{Limit: 1}); !errors.Is(err, context.Canceled) {
		t.Fatalf("instances=%v", err)
	}
}
func TestS144LateDeletedTaskNativeResultAndExactBinding(t *testing.T) {
	st := openMem(t)
	j, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	ds, _ := st.ClaimCompletionDeliveries(10)
	instance := joinInstances(t, st)[0]
	wrong := domain.Run{TaskID: entry.ID, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerCompletion, ScheduledFor: time.Now().UTC()}
	if err := st.RecordRunAndCreateDeliveries(&wrong, ds[0].ID); err == nil {
		t.Fatal("wrong task satisfied delivery")
	}
	if err := st.DeleteTask(j.Branches[0].TaskID); err != nil {
		t.Fatal(err)
	}
	r := joinFinish(t, st, ds[0], domain.OutcomeSuccess)
	got, _ := st.GetParallelJoinInstance(instance.ID)
	if got.State != domain.JoinInstanceResolved || !joinHasResult(got, ds[0].JoinBranchID, domain.OutcomeSuccess) {
		t.Fatalf("deleted late evidence=%+v", got)
	}
	if err := st.RecordRunAndCreateDeliveries(&r, ds[0].ID); err != nil {
		t.Fatal(err)
	}
}

func TestS144DefinitionValidationAndMissingReads(t *testing.T) {
	st := openMem(t)
	j, _ := joinFixture(t, st)
	cases := []struct {
		name   string
		change func(*domain.ParallelJoin)
	}{{"empty name", func(v *domain.ParallelJoin) { v.Name = "  " }}, {"empty branches", func(v *domain.ParallelJoin) {
		v.Branches = nil
		v.BranchCount = 0
	}}, {"duplicate branch ID", func(v *domain.ParallelJoin) { v.Branches[1].ID = v.Branches[0].ID }}, {"empty branch name", func(v *domain.ParallelJoin) { v.Branches[0].Name = " " }}, {"missing entry", func(v *domain.ParallelJoin) { v.EntryTaskID = "missing" }}, {"missing target", func(v *domain.ParallelJoin) { v.TargetTaskID = "missing" }}, {"missing branch", func(v *domain.ParallelJoin) { v.Branches[0].TaskID = "missing" }}, {"partial definition", func(v *domain.ParallelJoin) { v.HasMoreBranches = true }}, {"offset definition", func(v *domain.ParallelJoin) { v.BranchOffset = 1 }}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			copy := j
			copy.ID = ""
			copy.Branches = append([]domain.ParallelJoinBranch(nil), j.Branches...)
			tc.change(&copy)
			if err := st.CreateParallelJoin(&copy); !errors.Is(err, ErrInvalidJoin) {
				t.Fatalf("validation=%v", err)
			}
		})
	}
	cycle := j
	cycle.ID = ""
	cycle.Branches = append([]domain.ParallelJoinBranch(nil), j.Branches...)
	cycle.Branches[0].TaskID = j.EntryTaskID
	if err := st.CreateParallelJoin(&cycle); !errors.Is(err, ErrChainCycle) {
		t.Fatalf("self cycle=%v", err)
	}
	if _, err := st.GetParallelJoin("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing definition=%v", err)
	}
	if _, err := st.GetParallelJoinInstance("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing instance=%v", err)
	}
	if err := st.DeleteParallelJoin("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing delete=%v", err)
	}
	copy := j
	copy.ID = "missing"
	if err := st.UpdateParallelJoin(&copy); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing update=%v", err)
	}
	if _, err := st.GetCompletionDelivery("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing delivery=%v", err)
	}
	if ready, reason, err := st.CompletionDeliveryReady("missing"); err != nil || ready || reason == "" {
		t.Fatalf("missing readiness=%v %s %v", ready, reason, err)
	}
	if err := st.ResolveCompletionDelivery("missing", "absent"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing resolution=%v", err)
	}
}
func TestS144ListsFiltersAndUTF8BranchPagination(t *testing.T) {
	st := openMem(t)
	j, entry := joinFixture(t, st)
	j.Name = "日曜日 workflow"
	j.Branches[0].Name = "一番目 branch"
	if err := st.UpdateParallelJoin(&j); err != nil {
		t.Fatal(err)
	}
	r := joinEntry(t, st, entry)
	joinEntry(t, st, entry)
	joins, err := st.ListParallelJoins()
	if err != nil || len(joins) != 1 || len(joins[0].Branches) != 2 {
		t.Fatalf("definitions=%+v %v", joins, err)
	}
	page, err := st.ListParallelJoinsPage(0, 1, 4, 0, 1)
	if err != nil || len(page) != 1 || len(page[0].Branches) != 1 || !page[0].HasMoreBranches || page[0].BranchCount != 2 || len(page[0].TruncatedFields) == 0 || len(page[0].Branches[0].TruncatedFields) == 0 {
		t.Fatalf("bounded definitions=%+v %v", page, err)
	}
	if !utf8.ValidString(page[0].Name) || !utf8.ValidString(page[0].Branches[0].Name) || len(page[0].Name) > 4 {
		t.Fatal("bounded UTF8 name corrupted")
	}
	if page, err := st.ListParallelJoinsPage(1, 1, 1, 0, 1); err != nil || len(page) != 0 {
		t.Fatalf("definition offset=%+v %v", page, err)
	}
	instances, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{DefinitionID: j.ID, EntryRunID: r.ID, State: domain.JoinInstanceWaiting, Limit: 1, TextLimit: 4, BranchLimit: 1})
	if err != nil || len(instances) != 1 || !instances[0].HasMoreBranches || len(instances[0].Branches) != 1 || len(instances[0].TruncatedFields) == 0 {
		t.Fatalf("filtered instances=%+v %v", instances, err)
	}
	exact, err := st.GetParallelJoinInstancePage(instances[0].ID, 1, 1, 4)
	if err != nil || exact.BranchOffset != 1 || exact.HasMoreBranches || len(exact.Branches) != 1 {
		t.Fatalf("exact page=%+v %v", exact, err)
	}
	if instances, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{Offset: 1, Limit: 1}); err != nil || len(instances) != 1 {
		t.Fatalf("instance offset=%+v %v", instances, err)
	}
	if source, err := st.TaskHasJoinSource(j.Branches[0].TaskID); err != nil || !source {
		t.Fatalf("source=%v %v", source, err)
	}
}
func TestS144InputIdentityAndPendingDeliveryCannotPretendCompletion(t *testing.T) {
	st := openMem(t)
	_, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	inst := joinInstances(t, st)[0]
	d, err := st.GetCompletionDelivery(inst.Branches[0].DeliveryID)
	if err != nil {
		t.Fatal(err)
	}
	pending := domain.Run{TaskID: d.TargetTaskID, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerCompletion, ScheduledFor: time.Now().UTC()}
	if err := st.RecordRunAndCreateDeliveries(&pending, d.ID); err == nil {
		t.Fatal("unclaimed delivery accepted")
	}
	ds, _ := st.ClaimCompletionDeliveries(10)
	queued := domain.Run{TaskID: ds[0].TargetTaskID, Outcome: domain.OutcomeQueued, Trigger: domain.TriggerCompletion, ScheduledFor: time.Now().UTC()}
	if err := st.RecordRunAndCreateDeliveries(&queued, ds[0].ID); err == nil {
		t.Fatal("queued marker consumed required work")
	}
	r := joinFinish(t, st, ds[0], domain.OutcomeSuccess)
	conflicting := r
	conflicting.Outcome = domain.OutcomeFailure
	if err := st.RecordRunAndCreateDeliveries(&conflicting, ds[0].ID); err == nil {
		t.Fatal("conflicting delivery result accepted")
	}
	if ready, reason, err := st.CompletionDeliveryReady(ds[0].ID); err != nil || ready || reason == "" {
		t.Fatalf("terminal readiness=%v %s %v", ready, reason, err)
	}
	if err := st.RecordRunAndCreateDeliveries(&domain.Run{TaskID: entry.ID, Outcome: domain.OutcomeSuccess}, "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing incoming=%v", err)
	}
}
func TestS144ClosedStorageReturnsErrors(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	j := domain.ParallelJoin{Name: "closed"}
	checks := []func() error{func() error { return st.CreateParallelJoin(&j) }, func() error { return st.UpdateParallelJoin(&j) }, func() error { return st.DeleteParallelJoin("x") }, func() error {
		_, err := st.GetParallelJoin("x")
		return err
	}, func() error {
		_, err := st.ListParallelJoins()
		return err
	}, func() error {
		_, err := st.GetParallelJoinInstance("x")
		return err
	}, func() error {
		_, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{})
		return err
	}, func() error {
		_, err := st.TaskHasJoinSource("x")
		return err
	}, func() error {
		_, _, err := st.CompletionDeliveryReady("x")
		return err
	}, func() error { return st.ResolveCompletionDelivery("x", "closed") }, func() error { return st.RecordRunAndCreateDeliveries(&domain.Run{TaskID: "x"}, "") }, func() error { return st.ReplaceAutomationRelationships(nil, []domain.ParallelJoin{j}) }}
	for i, check := range checks {
		if err := check(); err == nil {
			t.Fatalf("closed operation %d succeeded", i)
		}
	}
}

func TestS144MixedPortableCreatesAndDuplicateInput(t *testing.T) {
	st := openMem(t)
	j, _ := joinFixture(t, st)
	a := chainTask(t, st, "portable a")
	b := chainTask(t, st, "portable b")
	chain := domain.CompletionChain{ID: "portable-chain", SourceTaskID: a.ID, TargetTaskID: b.ID, OnOutcome: domain.CompletionOnAny}
	if err := st.ReplaceAutomationRelationships([]domain.CompletionChain{chain}, nil); err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceAutomationRelationships([]domain.CompletionChain{chain, chain}, nil); !errors.Is(err, ErrInvalidJoin) {
		t.Fatalf("duplicate chain mutation=%v", err)
	}
	if err := st.ReplaceAutomationRelationships(nil, []domain.ParallelJoin{j, j}); !errors.Is(err, ErrInvalidJoin) {
		t.Fatalf("duplicate join mutation=%v", err)
	}
	bad := j
	bad.Name = ""
	if err := st.ReplaceParallelJoins([]domain.ParallelJoin{bad}); !errors.Is(err, ErrInvalidJoin) {
		t.Fatalf("invalid replacement=%v", err)
	}
	bad = j
	bad.TargetTaskID = "missing"
	if err := st.ReplaceParallelJoins([]domain.ParallelJoin{bad}); !errors.Is(err, ErrInvalidJoin) {
		t.Fatalf("invalid replacement reference=%v", err)
	}
	if err := st.ReplaceParallelJoins([]domain.ParallelJoin{j}); err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceAutomationRelationships(nil, nil); err != nil {
		t.Fatal(err)
	}
}
func TestS144StoreDecisionsRollbackOnDurableWriteFailure(t *testing.T) {
	st := openMem(t)
	_, entry := joinFixture(t, st)
	if _, err := st.db.Exec("CREATE TRIGGER fail_invocation BEFORE INSERT ON parallel_join_invocations BEGIN SELECT RAISE(ABORT,'owned invocation fixture failure'); END"); err != nil {
		t.Fatal(err)
	}
	failed := domain.Run{TaskID: entry.ID, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual, ScheduledFor: time.Now().UTC()}
	if err := st.RecordRunAndCreateDeliveries(&failed, ""); err == nil {
		t.Fatal("failed invocation creation reported committed")
	}
	if _, err := st.GetRun(failed.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("partial entry run committed: %v", err)
	}
	if len(joinInstances(t, st)) != 0 {
		t.Fatal("partial instance committed")
	}
	if _, err := st.db.Exec("DROP TRIGGER fail_invocation"); err != nil {
		t.Fatal(err)
	}
	joinEntry(t, st, entry)
	ds, _ := st.ClaimCompletionDeliveries(10)
	joinFinish(t, st, ds[0], domain.OutcomeSuccess)
	if _, err := st.db.Exec("CREATE TRIGGER fail_target BEFORE INSERT ON completion_deliveries WHEN NEW.role='join_target' BEGIN SELECT RAISE(ABORT,'target fixture failure'); END"); err != nil {
		t.Fatal(err)
	}
	failed = domain.Run{TaskID: ds[1].TargetTaskID, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerCompletion, ScheduledFor: time.Now().UTC()}
	if err := st.RecordRunAndCreateDeliveries(&failed, ds[1].ID); err == nil {
		t.Fatal("failed target decision reported committed")
	}
	d, err := st.GetCompletionDelivery(ds[1].ID)
	if err != nil || d.State != domain.DeliveryClaimed {
		t.Fatalf("failed decision consumed incoming=%+v %v", d, err)
	}
	if _, err := st.GetRun(failed.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("partial branch result committed=%v", err)
	}
	if _, err := st.db.Exec("DROP TRIGGER fail_target"); err != nil {
		t.Fatal(err)
	}
	if err := st.RecordRunAndCreateDeliveries(&failed, ds[1].ID); err != nil {
		t.Fatal(err)
	}
	targets, _ := st.ClaimCompletionDeliveries(10)
	if len(targets) != 1 {
		t.Fatalf("retry targets=%+v", targets)
	}
	if ready, _, err := st.CompletionDeliveryReady(targets[0].ID); err != nil || !ready {
		t.Fatalf("target readiness=%v %v", ready, err)
	}
}
func TestS144DefinitionDeletionRollsBackInvalidationOnFailure(t *testing.T) {
	st := openMem(t)
	j, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	ds, _ := st.ClaimCompletionDeliveries(10)
	if _, err := st.db.Exec("CREATE TRIGGER fail_delete BEFORE DELETE ON parallel_joins BEGIN SELECT RAISE(ABORT,'definition delete fixture failure'); END"); err != nil {
		t.Fatal(err)
	}
	if err := st.DeleteParallelJoin(j.ID); err == nil {
		t.Fatal("delete reported success despite storage failure")
	}
	got := joinInstances(t, st)[0]
	if got.State != domain.JoinInstanceWaiting {
		t.Fatalf("partial invalidation=%+v", got)
	}
	if ready, _, err := st.CompletionDeliveryReady(ds[0].ID); err != nil || !ready {
		t.Fatalf("partial delivery invalidation=%v %v", ready, err)
	}
}
func TestS144PendingTargetInvalidationAndLateNativeEvidence(t *testing.T) {
	st := openMem(t)
	j, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	ds, _ := st.ClaimCompletionDeliveries(10)
	for _, d := range ds {
		joinFinish(t, st, d, domain.OutcomeSuccess)
	}
	targets, _ := st.ClaimCompletionDeliveries(10)
	if err := st.DeleteTask(j.TargetTaskID); err != nil {
		t.Fatal(err)
	}
	r := joinFinish(t, st, targets[0], domain.OutcomeCanceled)
	got := joinInstances(t, st)[0]
	if got.State != domain.JoinInstanceResolved || got.TargetRunID != r.ID || got.TargetOutcome != domain.OutcomeCanceled {
		t.Fatalf("late invalidated target=%+v", got)
	}
	if err := st.RecordRunAndCreateDeliveries(&r, targets[0].ID); err != nil {
		t.Fatal(err)
	}
	r.Outcome = domain.OutcomeSuccess
	if err := st.RecordRunAndCreateDeliveries(&r, targets[0].ID); err == nil {
		t.Fatal("conflicting deleted native outcome accepted")
	}
}
func TestS144MalformedObservationTimestampsReturnErrors(t *testing.T) {
	st := openMem(t)
	j, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	inst := joinInstances(t, st)[0]
	for _, column := range []string{"created_at", "updated_at"} {
		if _, err := st.db.Exec("UPDATE parallel_joins SET "+column+"='invalid' WHERE id=?", j.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := st.GetParallelJoin(j.ID); err == nil {
			t.Fatalf("malformed definition %s silently accepted", column)
		}
		if _, err := st.db.Exec("UPDATE parallel_joins SET "+column+"=? WHERE id=?", fmtTime(time.Now().UTC()), j.ID); err != nil {
			t.Fatal(err)
		}
	}
	for _, column := range []string{"created_at", "updated_at", "completed_at", "target_started_at", "target_ended_at"} {
		if _, err := st.db.Exec("UPDATE parallel_join_instances SET "+column+"='invalid' WHERE id=?", inst.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := st.GetParallelJoinInstance(inst.ID); err == nil {
			t.Fatalf("malformed instance %s silently accepted", column)
		}
		if _, err := st.db.Exec("UPDATE parallel_join_instances SET "+column+"=? WHERE id=?", fmtTime(time.Now().UTC()), inst.ID); err != nil {
			t.Fatal(err)
		}
	}
	for _, column := range []string{"started_at", "ended_at"} {
		if _, err := st.db.Exec("UPDATE parallel_join_invocations SET "+column+"='invalid' WHERE id=?", inst.Branches[0].ID); err != nil {
			t.Fatal(err)
		}
		if _, err := st.GetParallelJoinInstance(inst.ID); err == nil {
			t.Fatalf("malformed invocation %s silently accepted", column)
		}
		if _, err := st.db.Exec("UPDATE parallel_join_invocations SET "+column+"=? WHERE id=?", fmtTime(time.Now().UTC()), inst.Branches[0].ID); err != nil {
			t.Fatal(err)
		}
	}
}

func TestS144ReadinessRejectsBrokenOwnershipAndRetainsOrdinaryChains(t *testing.T) {
	st := openMem(t)
	_, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	ds, _ := st.ClaimCompletionDeliveries(10)
	d := ds[0]
	if _, err := st.db.Exec("UPDATE completion_deliveries SET source_run_id='another-entry' WHERE id=?", d.ID); err != nil {
		t.Fatal(err)
	}
	if ready, reason, err := st.CompletionDeliveryReady(d.ID); err != nil || ready || reason == "" {
		t.Fatalf("mismatched source=%v %s %v", ready, reason, err)
	}
	if _, err := st.db.Exec("UPDATE completion_deliveries SET source_run_id=? WHERE id=?", d.SourceRunID, d.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := st.db.Exec("UPDATE parallel_join_invocations SET task_id='wrong-task' WHERE id=?", d.JoinBranchID); err != nil {
		t.Fatal(err)
	}
	if ready, _, err := st.CompletionDeliveryReady(d.ID); err != nil || ready {
		t.Fatalf("mismatched invocation=%v %v", ready, err)
	}
	a := chainTask(t, st, "ordinary source")
	b := chainTask(t, st, "ordinary target")
	c := domain.CompletionChain{SourceTaskID: a.ID, TargetTaskID: b.ID, OnOutcome: domain.CompletionOnSuccess}
	if err := st.CreateCompletionChain(&c); err != nil {
		t.Fatal(err)
	}
	joinEntry(t, st, a)
	ordinary, _ := st.ClaimCompletionDeliveries(10)
	if len(ordinary) != 1 {
		t.Fatalf("ordinary=%+v", ordinary)
	}
	if ready, _, err := st.CompletionDeliveryReady(ordinary[0].ID); err != nil || !ready {
		t.Fatalf("ordinary ready=%v %v", ready, err)
	}
	if err := st.DeleteCompletionChain(c.ID); err != nil {
		t.Fatal(err)
	}
	if ready, reason, err := st.CompletionDeliveryReady(ordinary[0].ID); err != nil || ready || reason == "" {
		t.Fatalf("deleted chain=%v %s %v", ready, reason, err)
	}
}
func TestS144IdempotentNativeTimesAndMonotonicRevisions(t *testing.T) {
	st := openMem(t)
	j, entry := joinFixture(t, st)
	frozen := j.UpdatedAt.Add(-time.Hour)
	st.now = func() time.Time { return frozen }
	old := j.UpdatedAt
	j.Name = "edited under fixed clock"
	if err := st.UpdateParallelJoinIfUnchanged(&j, &old); err != nil {
		t.Fatal(err)
	}
	if !j.UpdatedAt.After(old) {
		t.Fatal("revision failed to advance")
	}
	if err := st.ReplaceParallelJoins([]domain.ParallelJoin{j}); err != nil {
		t.Fatal(err)
	}
	start := time.Date(2026, 11, 10, 8, 30, 0, 123456789, time.FixedZone("fixture", -5*60*60))
	end := start.Add(time.Second)
	r := domain.Run{TaskID: entry.ID, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual, ScheduledFor: start, StartedAt: &start, EndedAt: &end, Output: "fixture output", OutputTruncated: true}
	if err := st.RecordRunAndCreateDeliveries(&r, ""); err != nil {
		t.Fatal(err)
	}
	if err := st.RecordRunAndCreateDeliveries(&r, ""); err != nil {
		t.Fatalf("exact normalized timestamps were not idempotent: %v", err)
	}
}

func TestS144AdmissionProjectionIsFreshAndExact(t *testing.T) {
	st := openMem(t)
	j, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	ds, _ := st.ClaimCompletionDeliveries(10)
	d := ds[0]
	actual, task, ready, reason, err := st.GetCompletionAdmission(d.ID)
	if err != nil || !ready || actual.ID != d.ID || actual.JoinBranchID != d.JoinBranchID || task.ID != d.TargetTaskID {
		t.Fatalf("admission=%+v %+v %v %s %v", actual, task, ready, reason, err)
	}
	if err := st.SetTaskEnabled(task.ID, false); err != nil {
		t.Fatal(err)
	}
	if _, _, ready, reason, err := st.GetCompletionAdmission(d.ID); err != nil || ready || reason == "" {
		t.Fatalf("disabled admission=%v %s %v", ready, reason, err)
	}
	if err := st.SetTaskEnabled(task.ID, true); err != nil {
		t.Fatal(err)
	}
	group := domain.Group{Name: "disabled parent", Enabled: false}
	if err := st.CreateGroup(&group); err != nil {
		t.Fatal(err)
	}
	task.GroupID = group.ID
	if err := st.UpdateTask(&task); err != nil {
		t.Fatal(err)
	}
	if _, _, ready, reason, err := st.GetCompletionAdmission(d.ID); err != nil || ready || reason == "" {
		t.Fatalf("ancestor admission=%v %s %v", ready, reason, err)
	}
	if err := st.SetGroupEnabled(group.ID, true); err != nil {
		t.Fatal(err)
	}
	if _, _, ready, reason, err := st.GetCompletionAdmission(d.ID); err != nil || !ready {
		t.Fatalf("reenabled admission=%v %s %v", ready, reason, err)
	}
	if err := st.DeleteParallelJoin(j.ID); err != nil {
		t.Fatal(err)
	}
	actual, _, ready, reason, err = st.GetCompletionAdmission(d.ID)
	if err != nil || ready || actual.ID != d.ID || reason == "" {
		t.Fatalf("invalidated admission=%+v %v %s %v", actual, ready, reason, err)
	}
	if _, _, ready, reason, err := st.GetCompletionAdmission("missing"); err != nil || ready || reason == "" {
		t.Fatalf("missing admission=%v %s %v", ready, reason, err)
	}
}

func TestS144AdmissionProjectionPreservesGroupCycleAndWhitespaceRefusal(t *testing.T) {
	st := openMem(t)
	_, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	ds, _ := st.ClaimCompletionDeliveries(10)
	d := ds[0]
	task, err := st.GetTask(d.TargetTaskID)
	if err != nil {
		t.Fatal(err)
	}
	group := domain.Group{Name: "cycle fixture", Enabled: true}
	if err := st.CreateGroup(&group); err != nil {
		t.Fatal(err)
	}
	task.GroupID = group.ID
	if err := st.UpdateTask(&task); err != nil {
		t.Fatal(err)
	}
	if _, err := st.db.Exec("UPDATE groups SET parent_id=id WHERE id=?", group.ID); err != nil {
		t.Fatal(err)
	}
	if _, _, ready, reason, err := st.GetCompletionAdmission(d.ID); err != nil || ready || reason == "" {
		t.Fatalf("cyclic group admission=%v %s %v", ready, reason, err)
	}
	if _, err := st.db.Exec("UPDATE groups SET parent_id=NULL WHERE id=?", group.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := st.db.Exec("UPDATE tasks SET command=?,enabled=1 WHERE id=?", "\t\n\u00a0", task.ID); err != nil {
		t.Fatal(err)
	}
	if _, _, ready, reason, err := st.GetCompletionAdmission(d.ID); err != nil || ready || reason != "target task has no command" {
		t.Fatalf("whitespace command admission=%v %s %v", ready, reason, err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	if _, _, _, _, err := st.GetCompletionAdmission(d.ID); err == nil {
		t.Fatal("closed admission reported certainty")
	}
}

func TestS144PreparedStatementsOwnLifetimeAndFreshAdmission(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if st.completionAdmission == nil || st.runInsert == nil || st.automationMatches == nil || st.chainDeliveryInsert == nil {
		t.Fatal("hot-path statements were not prepared with the store")
	}
	j, entry := joinFixture(t, st)
	joinEntry(t, st, entry)
	ds, _ := st.ClaimCompletionDeliveries(10)
	d := ds[0]
	if _, _, ready, _, err := st.GetCompletionAdmission(d.ID); err != nil || !ready {
		t.Fatalf("prepared admission=%v %v", ready, err)
	}
	if err := st.DeleteParallelJoin(j.ID); err != nil {
		t.Fatal(err)
	}
	if _, _, ready, _, err := st.GetCompletionAdmission(d.ID); err != nil || ready {
		t.Fatalf("prepared statement reused stale result=%v %v", ready, err)
	}
	saved := st.completionAdmission
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := saved.QueryRow(d.ID).Scan(&n); err == nil {
		t.Fatal("prepared statement outlived store close")
	}
	if err := st.Close(); err != nil {
		t.Fatalf("repeated close=%v", err)
	}
}
