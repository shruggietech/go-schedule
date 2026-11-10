package engine

import (
	"context"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

type noStartPublicationRunner struct {
	sharedID, branchID string
	manualStarted      chan struct{}
	manualRelease      chan struct{}
}

func (r *noStartPublicationRunner) Run(_ context.Context, task domain.Task, sf time.Time, trigger domain.RunTrigger) domain.Run {
	if task.ID == r.sharedID && trigger == domain.TriggerManual {
		r.manualStarted <- struct{}{}
		<-r.manualRelease
	}
	outcome := domain.OutcomeSuccess
	if task.ID == r.branchID {
		outcome = domain.OutcomeFailure
	}
	ended := sf.Add(time.Second)
	return domain.Run{TaskID: task.ID, ScheduledFor: sf, EndedAt: &ended, Outcome: outcome, Trigger: trigger}
}

func TestS145NoStartPublicationRetiresRejectedReservationBeforeCleanup(t *testing.T) {
	for _, boundary := range []string{"worker", "queue", "unreserved_other_active"} {
		t.Run(boundary, func(t *testing.T) {
			st, err := store.Open(":memory:")
			if err != nil {
				t.Fatal(err)
			}
			defer st.Close()
			entry := setupTask(t, st, domain.OverlapAllowConcurrent)
			branch := setupTask(t, st, domain.OverlapAllowConcurrent)
			target := setupTask(t, st, domain.OverlapAllowConcurrent)
			policy := domain.OverlapSkip
			if boundary == "queue" {
				policy = domain.OverlapQueueOne
			}
			shared := setupTask(t, st, policy)
			join := domain.ParallelJoin{Name: "no-start publication " + boundary, EntryTaskID: entry.ID, TargetTaskID: target.ID, RecoveryTaskID: shared.ID, CleanupTaskID: shared.ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "A", TaskID: branch.ID}}}
			if err := st.CreateParallelJoin(&join); err != nil {
				t.Fatal(err)
			}
			sf := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
			if boundary == "worker" {
				original := domain.Run{TaskID: entry.ID, ScheduledFor: sf, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
				if err := st.RecordRunAndCreateDeliveries(&original, ""); err != nil {
					t.Fatal(err)
				}
				claims, err := st.ClaimCompletionDeliveries(10)
				if err != nil || len(claims) != 1 {
					t.Fatalf("branch claims=%+v err=%v", claims, err)
				}
				failed := domain.Run{TaskID: branch.ID, ScheduledFor: sf, Outcome: domain.OutcomeFailure, Trigger: domain.TriggerCompletion}
				if err := st.RecordRunAndCreateDeliveries(&failed, claims[0].ID); err != nil {
					t.Fatal(err)
				}
			}
			runner := &noStartPublicationRunner{sharedID: shared.ID, branchID: branch.ID, manualStarted: make(chan struct{}, 1), manualRelease: make(chan struct{})}
			e, stop := startCompletionEngine(t, st, runner, func(e *Engine) {
				if boundary == "worker" {
					for range cap(e.sem) {
						e.sem <- struct{}{}
					}
				}
			})
			defer stop()
			published := make(chan struct{}, 1)
			publicationRelease := make(chan struct{})
			callbackErrors := make(chan error, 1)
			release := func(ch chan struct{}) {
				select {
				case <-ch:
				default:
					close(ch)
				}
			}
			capacityReleased := boundary != "worker"
			releaseCapacity := func() {
				if !capacityReleased {
					capacityReleased = true
					for range cap(e.sem) {
						<-e.sem
					}
				}
			}
			defer func() { release(runner.manualRelease); release(publicationRelease); releaseCapacity() }()
			e.SetOnJoinInstanceChanged(func(id string) {
				instance, err := st.GetParallelJoinInstance(id)
				if err != nil {
					callbackErrors <- err
					return
				}
				if instance.DefinitionID != join.ID || instance.RecoverySettlement != domain.JoinSettlementNotStarted || instance.CleanupSettlement != domain.JoinSettlementPending {
					return
				}
				// Model a legitimate owner re-enable between the rejected follow-up and
				// a concurrent cleanup drain. No old native execution remains here.
				if err := st.SetTaskEnabled(shared.ID, true); err != nil {
					callbackErrors <- err
					return
				}
				e.drainCompletionDeliveries()
				published <- struct{}{}
				<-publicationRelease
			})

			if boundary != "worker" {
				e.dispatch(shared, e.clk.Now(), domain.TriggerManual)
				recv(t, runner.manualStarted, "real manual predecessor")
			}
			if boundary == "unreserved_other_active" {
				shared.Enabled = false
				if err := st.UpdateTask(&shared); err != nil {
					t.Fatal(err)
				}
				e.dispatch(entry, e.clk.Now(), domain.TriggerManual)
			} else {
				if boundary == "queue" {
					e.dispatch(entry, e.clk.Now(), domain.TriggerManual)
					waitFor(t, func() bool {
						e.mu.Lock()
						defer e.mu.Unlock()
						pending, ok := e.queued[shared.ID]
						return ok && pending.origin.joinRole == "recovery"
					}, "recovery queued behind real native work")
				}
				shared.Enabled = false
				shared.OverlapPolicy = domain.OverlapSkip
				if err := st.UpdateTask(&shared); err != nil {
					t.Fatal(err)
				}
				if boundary == "worker" {
					releaseCapacity()
				} else {
					release(runner.manualRelease)
				}
			}
			select {
			case <-published:
			case err := <-callbackErrors:
				t.Fatal(err)
			case <-time.After(5 * time.Second):
				t.Fatal("rejected reservation did not publish its no-start disposition")
			}
			waitFor(t, func() bool {
				v, er := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{DefinitionID: join.ID})
				return er == nil && len(v) == 1 && v[0].CleanupSettlement != domain.JoinSettlementPending
			}, "re-enabled cleanup native result")
			instances, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{DefinitionID: join.ID})
			if err != nil || len(instances) != 1 {
				t.Fatalf("instances=%+v err=%v", instances, err)
			}
			got := instances[0]
			wantOutcome := domain.OutcomeSuccess
			if boundary == "unreserved_other_active" {
				wantOutcome = domain.OutcomeSkipped
				e.mu.Lock()
				count := e.running[shared.ID]
				e.mu.Unlock()
				active := false
				for _, run := range e.ActiveRuns() {
					active = active || (run.TaskID == shared.ID && run.Trigger == domain.TriggerManual && run.Outcome == "")
				}
				if count != 1 || !active {
					t.Fatalf("unreserved rejection retired genuine work: reservations=%d active=%t", count, active)
				}
			}
			if got.RecoverySettlement != domain.JoinSettlementNotStarted || got.RecoveryRunID != "" || got.CleanupOutcome != wantOutcome || got.CleanupRunID == "" {
				t.Fatalf("never-started predecessor blocked newly eligible cleanup: recovery settlement=%s run=%q cleanup settlement=%s outcome=%s reason=%q", got.RecoverySettlement, got.RecoveryRunID, got.CleanupSettlement, got.CleanupOutcome, got.CleanupReason)
			}
		})
	}
}

func TestS145AdmissionFailureRetiresOnlyOwnedReservation(t *testing.T) {
	for _, mode := range []string{"owned_read_error", "unreserved_read_error", "owned_shutdown", "unreserved_shutdown"} {
		t.Run(mode, func(t *testing.T) {
			st, err := store.Open(":memory:")
			if err != nil {
				t.Fatal(err)
			}
			defer st.Close()
			source := setupTask(t, st, domain.OverlapAllowConcurrent)
			target := setupTask(t, st, domain.OverlapSkip)
			delivery := claimedCompletion(t, st, source, target, "", "admission-source")
			runner := &completionRunner{outcomes: map[string]domain.RunOutcome{}, runs: make(chan domain.Run, 1)}
			e := newEngine(st, runner)
			e.running[target.ID] = 1
			var reservation *domain.Task
			owned := mode == "owned_read_error" || mode == "owned_shutdown"
			if owned {
				reservation = &target
			}
			readError := mode == "owned_read_error" || mode == "unreserved_read_error"
			if readError {
				if err := st.Close(); err != nil {
					t.Fatal(err)
				}
			} else {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				e.runCtx = ctx
			}
			if _, ok := e.eligibleCompletionTask(delivery.ID, reservation); ok {
				t.Fatal("failed admission accepted work")
			}
			e.mu.Lock()
			count := e.running[target.ID]
			e.mu.Unlock()
			want := 1
			if owned {
				want = 0
			}
			if count != want {
				t.Fatalf("reservation count=%d want=%d", count, want)
			}
			select {
			case run := <-runner.runs:
				t.Fatalf("failed admission executed native work: %+v", run)
			default:
			}
			if !readError {
				got, err := st.GetCompletionDelivery(delivery.ID)
				if err != nil || got.State != domain.DeliveryClaimed {
					t.Fatalf("shutdown destroyed restartable claim: %+v err=%v", got, err)
				}
			}
		})
	}
}
