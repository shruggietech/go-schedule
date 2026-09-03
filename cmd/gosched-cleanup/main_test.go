package main

import (
	"testing"

	"github.com/shruggietech/go-schedule/internal/winuninstall"
)

func TestRunAcceptsOnlyWipeWithoutPathArguments(t *testing.T) {
	called := false
	wipe := func() winuninstall.Result {
		called = true
		return winuninstall.Result{State: winuninstall.StateComplete}
	}
	if code := run([]string{"wipe"}, wipe); code != exitSuccess || !called {
		t.Fatalf("wipe code = %d, called = %v", code, called)
	}
	for _, args := range [][]string{nil, {"wipe", `C:\outside`}, {"report"}, {"--path", `C:\outside`}} {
		called = false
		if code := run(args, wipe); code != exitUsage || called {
			t.Errorf("run(%q) = %d, called = %v", args, code, called)
		}
	}
}

func TestRunMapsIncompleteCleanupToDistinctExit(t *testing.T) {
	for _, state := range []winuninstall.State{winuninstall.StateRefused, winuninstall.StatePartial} {
		code := run([]string{"wipe"}, func() winuninstall.Result { return winuninstall.Result{State: state} })
		if code != exitIncomplete {
			t.Errorf("state %q exit = %d, want %d", state, code, exitIncomplete)
		}
	}
	code := run([]string{"wipe"}, func() winuninstall.Result { return winuninstall.Result{State: winuninstall.StateInternalError} })
	if code != exitInternalError {
		t.Fatalf("internal-error exit = %d, want %d", code, exitInternalError)
	}
}
