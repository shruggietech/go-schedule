//go:build windows

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/desktopcontrol"
)

func TestTrayUsesInstalledServiceIPCPath(t *testing.T) {
	root := t.TempDir()
	t.Setenv("ProgramData", root)
	path := filepath.Join(root, "goschedule", "config.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	want := `\\.\pipe\goschedd-custom`
	data, err := json.Marshal(map[string]string{"ipc_path": want})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := localServiceEndpoint()
	if err != nil || got != want {
		t.Fatalf("tray endpoint = %q, %v; want %q", got, err, want)
	}
}

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
