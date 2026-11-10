package engine

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

func TestS144QueuedCompletionRechecksEligibility(t *testing.T) {
	for _, invalidation := range []string{"task_disabled", "chain_deleted"} {
		t.Run(invalidation, func(t *testing.T) {
			st, err := store.Open(":memory:")
			if err != nil {
				t.Fatal(err)
			}
			defer st.Close()
			source := setupTask(t, st, domain.OverlapQueueOne)
			target := setupTask(t, st, domain.OverlapQueueOne)
			chain := addChain(t, st, source, target, domain.CompletionOnSuccess)
			runner := &completionRunner{outcomes: map[string]domain.RunOutcome{}, runs: make(chan domain.Run, 1)}
			e := newEngine(st, runner)
			e.running[target.ID] = 1 // Deterministically hold the ordinary overlap reservation.
			sourceRun := domain.Run{TaskID: source.ID, ScheduledFor: time.Date(2026, 11, 10, 0, 0, 0, 0, time.UTC), Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
			if err := st.RecordRunAndCreateDeliveries(&sourceRun, ""); err != nil {
				t.Fatal(err)
			}
			e.drainCompletionDeliveries()
			pending, ok := e.queued[target.ID]
			if !ok || pending.origin.deliveryID == "" {
				t.Fatal("completion delivery was not queued")
			}
			if invalidation == "task_disabled" {
				if err := st.SetTaskEnabled(target.ID, false); err != nil {
					t.Fatal(err)
				}
			} else if err := st.DeleteCompletionChain(chain.ID); err != nil {
				t.Fatal(err)
			}
			e.finish(target) // Release the original reservation only after invalidation committed.
			e.runWG.Wait()
			deliveries, err := st.ListCompletionDeliveries()
			if err != nil {
				t.Fatal(err)
			}
			select {
			case run := <-runner.runs:
				t.Fatalf("invalidated queued target executed: trigger=%s source_run=%s delivery_state=%s; want no native execution and explicit resolved delivery", run.Trigger, run.SourceRunID, deliveries[0].State)
			default:
			}
			if len(deliveries) != 1 || deliveries[0].State != domain.DeliveryResolved {
				t.Fatalf("delivery states=%+v; want resolved", deliveries)
			}
		})
	}
}
