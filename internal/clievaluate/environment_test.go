package clievaluate

import (
	"slices"
	"strings"
	"testing"
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
