//go:build windows

package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/desktopcontrol"
)

func TestTrayActionsUseLocalSCMState(t *testing.T) {
	t.Parallel()
	cases := []struct {
		snapshot desktopcontrol.Snapshot
		pending  string
		want     []string
	}{
		{snapshot: desktopcontrol.Snapshot{State: "stopped", SCMState: "stopped"}, want: []string{"start"}},
		{snapshot: desktopcontrol.Snapshot{State: "running", SCMState: "running"}, want: []string{"stop", "restart"}},
		{snapshot: desktopcontrol.Snapshot{State: "unreachable", SCMState: "running"}, want: []string{"stop", "restart"}},
		{snapshot: desktopcontrol.Snapshot{State: "starting", SCMState: "starting"}},
		{snapshot: desktopcontrol.Snapshot{State: "not_installed", SCMState: "not_installed"}},
		{snapshot: desktopcontrol.Snapshot{State: "stopped", SCMState: "stopped"}, pending: "start"},
	}
	for _, tc := range cases {
		if got := availableActions(tc.snapshot, tc.pending); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("actions for %+v, pending %q = %v; want %v", tc.snapshot, tc.pending, got, tc.want)
		}
	}
}

func TestTrayTooltipIdentifiesThisComputer(t *testing.T) {
	t.Parallel()
	got := tooltipText(desktopcontrol.Snapshot{State: "not_installed"}, "")
	if !strings.Contains(got, "This computer") || !strings.Contains(got, "not installed") {
		t.Fatalf("tooltip = %q", got)
	}
}
