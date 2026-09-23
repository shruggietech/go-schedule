//go:build linux

package desktopcontrol

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestLinuxControlCommandUsesFixedUnitAndNoShell(t *testing.T) {
	t.Parallel()
	for _, action := range []string{"start", "stop", "restart"} {
		path, args, err := linuxControlCommand(action, false)
		if err != nil || path != "/usr/bin/pkexec" || !reflect.DeepEqual(args, []string{"/usr/bin/systemctl", action, linuxServiceUnit}) {
			t.Fatalf("%s: path=%q args=%v err=%v", action, path, args, err)
		}
		path, args, err = linuxControlCommand(action, true)
		if err != nil || path != "/usr/bin/systemctl" || !reflect.DeepEqual(args, []string{action, linuxServiceUnit}) {
			t.Fatalf("root %s: path=%q args=%v err=%v", action, path, args, err)
		}
	}
	if _, _, err := linuxControlCommand("stop; rm -rf /", false); err == nil {
		t.Fatal("unrecognized verb accepted")
	}
}

func TestLinuxControlCancelledCallerDoesNotStartHelper(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	err := requestElevation(ctx, "start")
	if !errors.Is(err, errHelperTimedOut) {
		t.Fatalf("cancelled request error = %v, want interrupted helper", err)
	}
	if time.Since(start) > time.Second {
		t.Fatalf("cancelled request waited for helper")
	}
}
