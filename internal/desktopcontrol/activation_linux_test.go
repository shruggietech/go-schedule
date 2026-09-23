//go:build linux

package desktopcontrol

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestLinuxSessionSingletonAndGUIActivation(t *testing.T) {
	t.Setenv("XDG_RUNTIME_DIR", t.TempDir())
	t.Setenv("XDG_SESSION_ID", "session-one")
	first, owner, err := ClaimGUI()
	if err != nil || !owner {
		t.Fatalf("first GUI claim: owner=%t err=%v", owner, err)
	}
	defer first.Close() //nolint:errcheck // test cleanup
	second, owner, err := ClaimGUI()
	if err != nil || owner || second != nil {
		t.Fatalf("second GUI claim: instance=%v owner=%t err=%v", second, owner, err)
	}
	t.Setenv("XDG_SESSION_ID", "session-two")
	otherSession, owner, err := ClaimGUI()
	if err != nil || !owner {
		t.Fatalf("other session GUI claim: owner=%t err=%v", owner, err)
	}
	defer otherSession.Close() //nolint:errcheck // test cleanup
	t.Setenv("XDG_SESSION_ID", "session-one")
	activated := make(chan struct{}, 1)
	stop, err := ListenGUI(context.Background(), func() { activated <- struct{}{} })
	if err != nil {
		t.Fatal(err)
	}
	defer stop()
	found, err := SignalGUI()
	if err != nil || !found {
		t.Fatalf("signal GUI: found=%t err=%v", found, err)
	}
	select {
	case <-activated:
	case <-time.After(time.Second):
		t.Fatal("GUI activation was not delivered")
	}
	tray, owner, err := ClaimTray()
	if err != nil || !owner {
		t.Fatalf("indicator claim: owner=%t err=%v", owner, err)
	}
	defer tray.Close() //nolint:errcheck // test cleanup
	if info, err := os.Stat(os.Getenv("XDG_RUNTIME_DIR") + "/go-schedule"); err != nil || info.Mode().Perm() != 0700 {
		t.Fatalf("session directory: %v, %v", info, err)
	}
}
