package clievaluate

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestWSLUsesIndependentOwnedLimitsWithoutShell(t *testing.T) {
	b := DefaultBudget()
	id := "cc06c2ec-016c-4f9b-90d5-b72283763aa7"
	args, err := WSLWorkerArguments(id, "/tmp/owned space", "/tmp/owned space/worker", []string{"worker", "literal $() value"}, b)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"--property=CPUQuota=25%", "--property=MemoryMax=1073741824", "--property=TasksMax=128", "--property=KillMode=control-group", "--no-ask-password", "--expand-environment=no", "literal $() value"} {
		if !slices.Contains(args, want) {
			t.Fatalf("missing worker control %q", want)
		}
	}
	if strings.Contains(strings.Join(args, "\n"), ".wslconfig") {
		t.Fatal("global WSL configuration touched")
	}
	if _, err := WSLWorkerArguments("../../other", "/tmp", "/worker", nil, b); err == nil {
		t.Fatal("unowned unit accepted")
	}
	b.CPUPercent = 0
	if _, err := WSLWorkerArguments(id, "/tmp", "/worker", nil, b); err == nil {
		t.Fatal("unbounded CPU accepted")
	}
}

func TestWSLUsesSharedShutdownAllowance(t *testing.T) {
	budget := DefaultBudget()
	args, err := WSLWorkerArguments("cc06c2ec-016c-4f9b-90d5-b72283763aa7", "/tmp/owned", "/tmp/owned/worker", nil, budget)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(args, "--property=RuntimeMaxSec=1845s") {
		t.Fatalf("unit hard limit disagrees with the 45-second worker allowance: %v", args)
	}
}

func TestWSLHardDeadlineUsesRemainingCoordinatorTimeWithoutReset(t *testing.T) {
	budget := DefaultBudget()
	budget.TotalTimeoutMS = 7200000
	for _, remaining := range []time.Duration{0, time.Millisecond, 20 * time.Second, 2 * time.Hour} {
		args, err := wslWorkerArguments("cc06c2ec-016c-4f9b-90d5-b72283763aa7", "/tmp/owned", "/tmp/owned/worker", nil, budget, remaining)
		if err != nil {
			t.Fatal(err)
		}
		expected := remaining + 45*time.Second
		seconds := (expected + time.Second - 1) / time.Second
		want := fmt.Sprintf("--property=RuntimeMaxSec=%ds", seconds)
		if !slices.Contains(args, want) || !slices.Contains(args, "--property=CPUQuota=25%") || !slices.Contains(args, "--property=TasksMax=128") {
			t.Fatalf("remaining budget or controls changed: %v", args)
		}
	}
	for _, bad := range []time.Duration{-time.Millisecond, 2*time.Hour + time.Millisecond} {
		if _, err := wslWorkerArguments("cc06c2ec-016c-4f9b-90d5-b72283763aa7", "/tmp/owned", "/worker", nil, budget, bad); err == nil {
			t.Fatal("invalid remaining budget accepted")
		}
	}
}
