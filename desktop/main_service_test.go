package main

import (
	"context"
	"errors"
	"runtime"
	"testing"

	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/desktopcontrol"
	"github.com/shruggietech/go-schedule/internal/service"
)

func TestLocalServiceConfigPathMatchesInstalledDaemonOnWindows(t *testing.T) {
	t.Parallel()
	if got := localConfigPath("windows"); got != config.DefaultPath() {
		t.Fatalf("Windows service config path = %q, want %q", got, config.DefaultPath())
	}
	for _, goos := range []string{"linux", "darwin"} {
		if got := localConfigPath(goos); got != "" {
			t.Fatalf("%s standalone GUI config path = %q, want empty", goos, got)
		}
	}
}

func TestInstalledServiceNeverAutospawns(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "windows" {
		t.Skip("installed Windows service semantics")
	}
	for _, state := range []service.State{service.StateStopped, service.StateStarting, service.StateStopping, service.StateRunning, service.StateUnknown} {
		if shouldAutoSpawnInstalledService(state, nil) {
			t.Fatalf("%s must not auto-spawn a competing daemon", state)
		}
	}
	if shouldAutoSpawnInstalledService(service.StateUnknown, errors.New("access denied")) {
		t.Fatal("failed SCM query must not auto-spawn")
	}
	if !shouldAutoSpawnInstalledService(service.StateNotInstalled, nil) {
		t.Fatal("standalone GUI should retain bundled daemon behavior")
	}
}

func TestControlLocalServiceRejectsUnconfirmedStop(t *testing.T) {
	t.Parallel()
	calls := 0
	monitor := &desktopcontrol.Monitor{
		Query:   func() (service.State, error) { return service.StateRunning, nil },
		Health:  func(context.Context) error { return nil },
		Execute: func(string) error { calls++; return nil },
	}
	app := &App{ctx: context.Background(), localService: monitor}
	result := app.ControlLocalService("stop", false)
	if result.Outcome != "rejected" || calls != 0 {
		t.Fatalf("result=%+v, executions=%d", result, calls)
	}
}
