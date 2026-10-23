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
