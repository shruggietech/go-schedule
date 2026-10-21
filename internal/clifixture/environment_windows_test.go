//go:build windows

package clifixture

import (
	"slices"
	"testing"
)

func TestWindowsEnvironmentLastOverrideWinsAndIsSorted(t *testing.T) {
	actual, err := normalizedWindowsEnvironment([]string{"PATH=first", "GOMAXPROCS=8", "Path=last", "GOMAXPROCS=2", "=C:=C:\\owned"})
	want := []string{"=C:=C:\\owned", "GOMAXPROCS=2", "Path=last"}
	if err != nil || !slices.Equal(actual, want) {
		t.Fatalf("environment=%q error=%v", actual, err)
	}
	if _, err := normalizedWindowsEnvironment([]string{"GOMAXPROCS=2\x00INJECT=1"}); err == nil {
		t.Fatal("NUL-bearing environment accepted")
	}
}
