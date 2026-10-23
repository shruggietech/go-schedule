package clievaluate

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTimingPolicyHonorsValidatedBudgetAndBoundedReserve(t *testing.T) {
	for _, total := range []int64{1800000, 7200000} {
		budget := DefaultBudget()
		budget.TotalTimeoutMS = total
		policy, err := TimingForBudget(budget)
		if err != nil {
			t.Fatal(err)
		}
		if policy.TotalTimeoutMS != total || policy.WorkerShutdownMS <= 0 || policy.ReportTimeoutMS <= 0 || policy.OuterTimeoutMS != total+policy.WorkerShutdownMS+policy.ReportTimeoutMS+15000 {
			t.Fatalf("invalid policy: %+v", policy)
		}
		if policy.OuterTimeoutMS > total+int64((2*time.Minute)/time.Millisecond) {
			t.Fatalf("unbounded reserve: %+v", policy)
		}
		if budget.CaseTimeoutMS != 15000 {
			t.Fatal("case limit changed")
		}
	}
	invalid := DefaultBudget()
	invalid.TotalTimeoutMS = 7200001
	if _, err := TimingForBudget(invalid); err == nil {
		t.Fatal("invalid budget accepted")
	}
}

func TestSelectedBudgetResumesOriginalAndRejectsIncompleteOverrides(t *testing.T) {
	dir := t.TempDir()
	budget := DefaultBudget()
	budget.TotalTimeoutMS = 7200000
	data, err := Encode(captureIdentity{Budget: budget})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "identity.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
	got, err := SelectedBudget("", dir)
	if err != nil || got != budget {
		t.Fatalf("resume budget %+v %v", got, err)
	}
	input := filepath.Join(dir, "budget.json")
	if err := os.WriteFile(input, []byte(`{"total_timeout_ms":7200000}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := SelectedBudget(input, dir); err == nil {
		t.Fatal("partial override silently inherited defaults")
	}
}

func TestFinalizationSharesHardReserveAfterExecutionDeadline(t *testing.T) {
	budget := DefaultBudget()
	policy, err := TimingForBudget(budget)
	if err != nil {
		t.Fatal(err)
	}
	expired := time.Now().Add(-50 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), expired)
	defer cancel()
	report, stop := finalizationContext(ctx, policy)
	defer stop()
	deadline, _ := report.Deadline()
	hard := expired.Add(105 * time.Second)
	if deadline.After(hard) || report.Err() != nil {
		t.Fatalf("report reserve is expired or reset: %s %v", deadline, report.Err())
	}
}

func TestWorkerFinalizationUsesRemainingExecutionBudget(t *testing.T) {
	execution, stop := context.WithTimeout(context.Background(), time.Second)
	defer stop()
	final, cancel := workerFinalizationContext(execution, 20*time.Millisecond)
	defer cancel()
	select {
	case <-final.Done():
		t.Fatal("normal finalization was cut short by an unrelated grace cap")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestWorkerFinalizationBoundsEarlyCancellationAndJoinsWatcher(t *testing.T) {
	execution, stop := context.WithTimeout(context.Background(), time.Hour)
	final, cancel := workerFinalizationContext(execution, 100*time.Millisecond)
	defer cancel()
	stop()
	if final.Err() != nil {
		t.Fatal("execution cancellation removed the independent reporting reserve")
	}
	select {
	case <-final.Done():
	case <-time.After(time.Second):
		t.Fatal("early cancellation retained the original hour-long deadline")
	}
}

func TestWorkerFinalizationPreCanceledAndUnboundedContextsHaveFiniteReserve(t *testing.T) {
	canceled, stop := context.WithCancel(context.Background())
	stop()
	for _, execution := range []context.Context{canceled, context.Background()} {
		final, cancel := workerFinalizationContext(execution, 100*time.Millisecond)
		if final.Err() != nil {
			cancel()
			t.Fatal("reporting reserve was already canceled")
		}
		select {
		case <-final.Done():
		case <-time.After(time.Second):
			cancel()
			t.Fatal("reporting reserve was unbounded")
		}
		cancel()
	}
}

func TestWorkerFinalizationRetainsOnlyGraceBeyondExecutionDeadline(t *testing.T) {
	execution, stop := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer stop()
	final, cancel := workerFinalizationContext(execution, 200*time.Millisecond)
	defer cancel()
	<-execution.Done()
	if final.Err() != nil {
		t.Fatal("execution deadline erased bounded reporting reserve")
	}
	select {
	case <-final.Done():
	case <-time.After(time.Second):
		t.Fatal("reporting outlived its finite grace")
	}
}
