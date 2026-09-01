//go:build windows

package gui

import (
	"os"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/ipc"
)

// TestNativeWindowsConnectionRecovery is an opt-in walkthrough probe. It is
// skipped in ordinary CI because it requires the installed service and the
// interactive account state produced by the MSI.
func TestNativeWindowsConnectionRecovery(t *testing.T) {
	expect := os.Getenv("GOSCHED_NATIVE_EXPECT")
	if expect == "" {
		t.Skip("set GOSCHED_NATIVE_EXPECT=stale or connected on an installed Windows host")
	}
	if expect != "stale" && expect != "connected" {
		t.Fatalf("GOSCHED_NATIVE_EXPECT = %q, want stale or connected", expect)
	}

	diagnosis := diagnoseAccess()
	if diagnosis.Service != "running" || diagnosis.GroupExists != yes || diagnosis.AccountMember != yes {
		t.Fatalf("native prerequisites: service=%q group=%v account=%v detail=%q",
			diagnosis.Service, diagnosis.GroupExists, diagnosis.AccountMember, diagnosis.Detail)
	}

	ui := NewUI(testApp, client.New(ipc.Endpoint(config.Default())))
	err := ui.refreshAllOnce()
	if expect == "stale" {
		if diagnosis.TokenMember != no {
			t.Fatalf("token membership = %v, want absent stale token", diagnosis.TokenMember)
		}
		if err == nil {
			t.Fatal("stale token unexpectedly reached the daemon")
		}
		waitFor(t, func() bool {
			incident, active := ui.connection.snapshot()
			return active && incident.Kind == client.ConnectionAccessDenied &&
				containsFold(incident.Guidance, "sign out")
		})
		return
	}

	if diagnosis.TokenMember != yes {
		t.Fatalf("token membership = %v, want refreshed group token", diagnosis.TokenMember)
	}
	if err != nil {
		t.Fatalf("refreshed token could not load the GUI model: %v", err)
	}
	if _, active := ui.connection.snapshot(); active {
		t.Fatal("connection incident remained active after native recovery")
	}
	if ui.tabs == nil || len(ui.tabs.Items) == 0 {
		t.Fatal("normal GUI tabs are unavailable after native recovery")
	}
}
