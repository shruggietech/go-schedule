package engine

import (
	"fmt"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

// BenchmarkParallelJoinLastBranchDecision isolates the atomic final contribution
// and unique target admission. Setup and downstream cleanup are outside timing.
func BenchmarkParallelJoinLastBranchDecision(b *testing.B) {
	st, err := store.Open(":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer st.Close()
	tasks := make([]domain.Task, 3)
	for i := range tasks {
		tasks[i] = domain.Task{Name: fmt.Sprintf("task-%d", i), Command: "unused", Timezone: "UTC", State: domain.TaskActive}
		if err := st.CreateTask(&tasks[i]); err != nil {
			b.Fatal(err)
		}
	}
	join := domain.ParallelJoin{Name: "benchmark", EntryTaskID: tasks[0].ID, TargetTaskID: tasks[2].ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "first", TaskID: tasks[1].ID}, {ID: "b", Name: "second", TaskID: tasks[1].ID}}}
	if err := st.CreateParallelJoin(&join); err != nil {
		b.Fatal(err)
	}
	now := time.Date(2026, 11, 10, 0, 0, 0, 0, time.UTC)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		source := domain.Run{TaskID: tasks[0].ID, ScheduledFor: now, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
		if err := st.RecordRunAndCreateDeliveries(&source, ""); err != nil {
			b.Fatal(err)
		}
		deliveries, err := st.ClaimCompletionDeliveries(10)
		if err != nil || len(deliveries) != 2 {
			b.Fatalf("branch claims=%d err=%v", len(deliveries), err)
		}
		first := domain.Run{TaskID: tasks[1].ID, ScheduledFor: now, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerCompletion}
		if err := st.RecordRunAndCreateDeliveries(&first, deliveries[0].ID); err != nil {
			b.Fatal(err)
		}
		last := domain.Run{TaskID: tasks[1].ID, ScheduledFor: now, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerCompletion}
		b.StartTimer()
		err = st.RecordRunAndCreateDeliveries(&last, deliveries[1].ID)
		b.StopTimer()
		if err != nil {
			b.Fatal(err)
		}
		targets, err := st.ClaimCompletionDeliveries(10)
		if err != nil || len(targets) != 1 {
			b.Fatalf("target claims=%d err=%v", len(targets), err)
		}
		if err := st.ResolveCompletionDelivery(targets[0].ID, "benchmark cleanup"); err != nil {
			b.Fatal(err)
		}
		b.StartTimer()
	}
}
