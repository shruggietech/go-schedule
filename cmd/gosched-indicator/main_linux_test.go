//go:build linux

package main

import (
	"testing"

	"github.com/shruggietech/go-schedule/internal/desktopcontrol"
)

func TestStateLabel(t *testing.T) {
	t.Parallel()
	for input, want := range map[string]string{
		"running": "Running", "stopped": "Stopped", "not_installed": "Not installed", "unreachable": "Unreachable", "unknown": "Unknown", "": "Unknown",
	} {
		if got := stateLabel(input); got != want {
			t.Fatalf("stateLabel(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestAvailableActionsNeverMutateMissingOrTransitionalService(t *testing.T) {
	t.Parallel()
	tests := []struct {
		state, manager string
		start, stop    bool
	}{
		{"stopped", "stopped", true, false},
		{"running", "running", false, true},
		{"unreachable", "running", false, true},
		{"starting", "starting", false, false},
		{"not_installed", "not_installed", false, false},
		{"unknown", "unknown", false, false},
	}
	for _, tt := range tests {
		start, stop := availableActions(desktopcontrol.Snapshot{State: tt.state, SCMState: tt.manager})
		if start != tt.start || stop != tt.stop {
			t.Fatalf("%s/%s: start=%t stop=%t; want start=%t stop=%t", tt.state, tt.manager, start, stop, tt.start, tt.stop)
		}
	}
}
