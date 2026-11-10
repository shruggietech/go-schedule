package engine

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

// overlapPublicationRunner holds the first owned follow-up until the test has
// arranged any genuine concurrent or queued reservation. The returned native
// result and its later persistence callback are separate deterministic barriers.
type overlapPublicationRunner struct {
	sharedID, branchID             string
	failBranch                     bool
	followupGate, manualGate       chan struct{}
	followupStarted, manualStarted chan struct{}
	completions                    atomic.Int32
}

func (r *overlapPublicationRunner) Run(_ context.Context, task domain.Task, sf time.Time, trigger domain.RunTrigger) domain.Run {
	if task.ID == r.sharedID {
		if trigger == domain.TriggerManual {
			r.manualStarted <- struct{}{}
			<-r.manualGate
		} else if r.completions.Add(1) == 1 {
			r.followupStarted <- struct{}{}
			<-r.followupGate
		}
	}
	outcome := domain.OutcomeSuccess
	if r.failBranch && task.ID == r.branchID {
		outcome = domain.OutcomeFailure
	}
	ended := sf.Add(time.Second)
	return domain.Run{TaskID: task.ID, ScheduledFor: sf, EndedAt: &ended, Outcome: outcome, Trigger: trigger}
}

func TestS145CleanupPublicationRetiresOnlyCompletedReservation(t *testing.T) {
	for _, mode := range []string{"success_skip", "recovery_skip", "genuine_other_active", "queued_handoff"} {
		t.Run(mode, func(t *testing.T) {
			st, err := store.Open(":memory:")
			if err != nil {
				t.Fatal(err)
			}
			defer st.Close()
			entry := setupTask(t, st, domain.OverlapAllowConcurrent)
			branch := setupTask(t, st, domain.OverlapAllowConcurrent)
			policy := domain.OverlapSkip
			if mode == "genuine_other_active" {
				policy = domain.OverlapAllowConcurrent
			}
			if mode == "queued_handoff" {
				policy = domain.OverlapQueueOne
			}
			shared := setupTask(t, st, policy)
			join := domain.ParallelJoin{Name: mode, EntryTaskID: entry.ID, TargetTaskID: shared.ID, CleanupTaskID: shared.ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "A", TaskID: branch.ID}}}
			followupRole := "target"
			if mode == "recovery_skip" {
				unused := setupTask(t, st, domain.OverlapAllowConcurrent)
				join.TargetTaskID = unused.ID
				join.RecoveryTaskID = shared.ID
				followupRole = "recovery"
			}
			if err := st.CreateParallelJoin(&join); err != nil {
				t.Fatal(err)
			}
			runner := &overlapPublicationRunner{sharedID: shared.ID, branchID: branch.ID, failBranch: mode == "recovery_skip", followupGate: make(chan struct{}), manualGate: make(chan struct{}), followupStarted: make(chan struct{}, 1), manualStarted: make(chan struct{}, 1)}
			e, stop := startCompletionEngine(t, st, runner)
			defer stop()
			releaseCommit := make(chan struct{})
			release := func(ch chan struct{}) {
				select {
				case <-ch:
				default:
					close(ch)
				}
			}
			defer func() { release(runner.followupGate); release(runner.manualGate); release(releaseCommit) }()
			committed := make(chan domain.Run, 1)
			e.SetOnRun(func(run domain.Run) {
				if run.JoinRole == followupRole && run.Outcome != domain.OutcomeQueued {
					committed <- run
					<-releaseCommit
				}
			})
			if mode == "genuine_other_active" {
				e.dispatch(shared, e.clk.Now(), domain.TriggerManual)
				recv(t, runner.manualStarted, "genuine other run")
			}
			e.dispatch(entry, e.clk.Now(), domain.TriggerManual)
			recv(t, runner.followupStarted, "selected native follow-up")
			if mode == "queued_handoff" {
				e.dispatch(shared, e.clk.Now(), domain.TriggerManual)
			}
			if mode == "genuine_other_active" {
				shared.OverlapPolicy = domain.OverlapSkip
				if err := st.UpdateTask(&shared); err != nil {
					t.Fatal(err)
				}
			}
			release(runner.followupGate)
			var predecessor domain.Run
			select {
			case predecessor = <-committed:
			case <-time.After(5 * time.Second):
				t.Fatal("follow-up completion did not reach persistence barrier")
			}
			// The callback holds the predecessor goroutine after its completion commit.
			// A pre-existing or concurrent drain can see cleanup before that callback
			// returns; correctness cannot depend on moving the new wake after it.
			e.drainCompletionDeliveries()
			decision, err := e.CancelRun(context.Background(), predecessor.ID)
			if err != nil || decision.Status != "already_finished" || decision.Outcome != domain.OutcomeSuccess {
				t.Fatalf("terminal cancellation observation=%+v err=%v", decision, err)
			}
			instances, err := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{DefinitionID: join.ID})
			if err != nil || len(instances) != 1 {
				t.Fatalf("instances=%+v err=%v", instances, err)
			}
			if mode == "queued_handoff" {
				if instances[0].CleanupSettlement != domain.JoinSettlementPending || instances[0].CleanupOutcome != "" || instances[0].CleanupReason != "" {
					t.Fatalf("cleanup collapsed behind completed predecessor instead of queued handoff: settlement=%s outcome=%s reason=%q", instances[0].CleanupSettlement, instances[0].CleanupOutcome, instances[0].CleanupReason)
				}
				recv(t, runner.manualStarted, "queued reservation handoff")
				e.mu.Lock()
				count := e.running[shared.ID]
				e.mu.Unlock()
				if count != 1 {
					t.Fatalf("queued handoff running count=%d want=1", count)
				}
				release(releaseCommit)
				release(runner.manualGate)
			}
			waitFor(t, func() bool {
				v, er := st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{DefinitionID: join.ID})
				return er == nil && len(v) == 1 && v[0].CleanupSettlement != domain.JoinSettlementPending
			}, "cleanup terminal evidence")
			instances, err = st.ListParallelJoinInstances(domain.ParallelJoinInstanceFilter{DefinitionID: join.ID})
			if err != nil {
				t.Fatal(err)
			}
			got := instances[0]
			if mode == "genuine_other_active" {
				e.mu.Lock()
				count := e.running[shared.ID]
				e.mu.Unlock()
				if got.CleanupOutcome != domain.OutcomeSkipped || runner.completions.Load() != 1 || count != 1 {
					t.Fatalf("genuine overlap lost: outcome=%s native completions=%d reservations=%d", got.CleanupOutcome, runner.completions.Load(), count)
				}
			} else if got.CleanupOutcome != domain.OutcomeSuccess || got.CleanupRunID == "" || got.CleanupRunID == predecessor.ID || runner.completions.Load() != 2 {
				t.Fatalf("completed predecessor blocked reused cleanup: settlement=%s outcome=%s native completions=%d", got.CleanupSettlement, got.CleanupOutcome, runner.completions.Load())
			}
		})
	}
}
