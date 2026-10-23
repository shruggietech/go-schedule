package clievaluate

import (
	"context"
	"fmt"
	"path/filepath"
	"time"
)

// TimingPolicy separates the selected execution budget from bounded cleanup and reporting.
type TimingPolicy struct {
	TotalTimeoutMS   int64 `json:"total_timeout_ms"`
	WorkerShutdownMS int64 `json:"worker_shutdown_ms"`
	ReportTimeoutMS  int64 `json:"report_timeout_ms"`
	OuterTimeoutMS   int64 `json:"outer_timeout_ms"`
}

// TimingForBudget is the single validated policy used by the coordinator and Windows wrapper.
func TimingForBudget(budget Budget) (TimingPolicy, error) {
	if err := ValidateBudget(budget); err != nil {
		return TimingPolicy{}, err
	}
	return TimingPolicy{budget.TotalTimeoutMS, 45000, 60000, budget.TotalTimeoutMS + 120000}, nil
}

// SelectedBudget uses the original recorded budget for a resume without an explicit override.
func SelectedBudget(input, resumeRoot string) (Budget, error) {
	budget := DefaultBudget()
	if input != "" {
		budget = Budget{}
		if err := ReadInput(input, &budget); err != nil {
			return Budget{}, err
		}
	} else if resumeRoot != "" {
		var identity captureIdentity
		if err := ReadInput(filepath.Join(resumeRoot, "identity.json"), &identity); err != nil {
			return Budget{}, fmt.Errorf("read original resume budget: %w", err)
		}
		budget = identity.Budget
	}
	return budget, ValidateBudget(budget)
}

func finalizationContext(execution context.Context, policy TimingPolicy) (context.Context, context.CancelFunc) {
	deadline := time.Now().Add(time.Duration(policy.ReportTimeoutMS) * time.Millisecond)
	if executionDeadline, ok := execution.Deadline(); ok {
		hard := executionDeadline.Add(time.Duration(policy.WorkerShutdownMS+policy.ReportTimeoutMS) * time.Millisecond)
		if hard.Before(deadline) {
			deadline = hard
		}
	}
	return context.WithDeadline(context.WithoutCancel(execution), deadline)
}
