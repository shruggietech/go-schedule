package engine

import (
	"fmt"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

// BenchmarkParallelRecoveryStages isolates the two atomic continuation decisions.
// Instance setup and downstream resolution are outside timing.
func BenchmarkParallelRecoveryStages(b *testing.B) {
	for _, stage := range []string{"recovery_after_last_branch", "cleanup_after_followup"} {
		b.Run(stage, func(b *testing.B) {
			st, err := store.Open(":memory:")
			if err != nil {
				b.Fatal(err)
			}
			b.Cleanup(func() {
				if err := st.Close(); err != nil {
					b.Error(err)
				}
			})
			tasks := make([]domain.Task, 5)
			for i := range tasks {
				tasks[i] = domain.Task{Name: fmt.Sprintf("route-bench-%d", i), Command: "unused", Timezone: "UTC", State: domain.TaskActive}
				if err := st.CreateTask(&tasks[i]); err != nil {
					b.Fatal(err)
				}
			}
			join := domain.ParallelJoin{Name: "route-benchmark", EntryTaskID: tasks[0].ID, TargetTaskID: tasks[2].ID, RecoveryTaskID: tasks[3].ID, CleanupTaskID: tasks[4].ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "first", TaskID: tasks[1].ID}, {ID: "b", Name: "second", TaskID: tasks[1].ID}}}
			if err := st.CreateParallelJoin(&join); err != nil {
				b.Fatal(err)
			}
			now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
			record := func(taskID string, outcome domain.RunOutcome, incoming string) error {
				run := domain.Run{TaskID: taskID, ScheduledFor: now, StartedAt: &now, EndedAt: &now, Outcome: outcome, Trigger: domain.TriggerCompletion}
				return st.RecordRunAndCreateDeliveries(&run, incoming)
			}
			claim := func(role domain.CompletionDeliveryRole) domain.CompletionDelivery {
				deliveries, err := st.ClaimCompletionDeliveries(10)
				if err != nil || len(deliveries) != 1 || deliveries[0].Role != role {
					b.Fatalf("exact role %s claims=%+v err=%v", role, deliveries, err)
				}
				return deliveries[0]
			}
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				b.StopTimer()
				source := domain.Run{TaskID: tasks[0].ID, ScheduledFor: now, StartedAt: &now, EndedAt: &now, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
				if err := st.RecordRunAndCreateDeliveries(&source, ""); err != nil {
					b.Fatal(err)
				}
				branches, err := st.ClaimCompletionDeliveries(10)
				if err != nil || len(branches) != 2 {
					b.Fatalf("branches=%+v err=%v", branches, err)
				}
				if err := record(tasks[1].ID, domain.OutcomeFailure, branches[0].ID); err != nil {
					b.Fatal(err)
				}
				if stage == "recovery_after_last_branch" {
					b.StartTimer()
					err = record(tasks[1].ID, domain.OutcomeSuccess, branches[1].ID)
					b.StopTimer()
					if err != nil {
						b.Fatal(err)
					}
					recovery := claim(domain.DeliveryRoleJoinRecovery)
					if err := record(tasks[3].ID, domain.OutcomeFailure, recovery.ID); err != nil {
						b.Fatal(err)
					}
				} else {
					if err := record(tasks[1].ID, domain.OutcomeSuccess, branches[1].ID); err != nil {
						b.Fatal(err)
					}
					recovery := claim(domain.DeliveryRoleJoinRecovery)
					b.StartTimer()
					err = record(tasks[3].ID, domain.OutcomeFailure, recovery.ID)
					b.StopTimer()
					if err != nil {
						b.Fatal(err)
					}
				}
				cleanup := claim(domain.DeliveryRoleJoinCleanup)
				if err := st.ResolveCompletionDelivery(cleanup.ID, "benchmark downstream resolution"); err != nil {
					b.Fatal(err)
				}
				b.StartTimer()
			}
			b.StopTimer()
		})
	}
}
