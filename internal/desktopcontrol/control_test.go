package desktopcontrol

import (
	"context"
	"errors"
	"testing"

	"github.com/shruggietech/go-schedule/internal/service"
)

func TestObserveRequiresSCMAndFreshHealth(t *testing.T) {
	t.Parallel()
	tests := []struct {
		scm     service.State
		health  error
		want    string
		healthN int
	}{
		{scm: service.StateRunning, want: "running", healthN: 1},
		{scm: service.StateRunning, health: errors.New("pipe unavailable"), want: "unreachable", healthN: 1},
		{scm: service.StateStarting, want: "starting"},
		{scm: service.StateStopping, want: "stopping"},
		{scm: service.StateStopped, want: "stopped"},
		{scm: service.StateNotInstalled, want: "not_installed"},
	}
	for _, tc := range tests {
		t.Run(string(tc.scm)+tc.want, func(t *testing.T) {
			t.Parallel()
			calls := 0
			m := Monitor{Query: func() (service.State, error) { return tc.scm, nil }, Health: func(context.Context) error {
				calls++
				return tc.health
			}}
			got := m.Observe(context.Background())
			if got.State != tc.want || calls != tc.healthN {
				t.Fatalf("Observe = %q, health calls = %d; want %q and %d", got.State, calls, tc.want, tc.healthN)
			}
		})
	}
}

func TestObserveQueryFailureIsNotStopped(t *testing.T) {
	t.Parallel()
	m := Monitor{Query: func() (service.State, error) { return service.StateUnknown, errors.New("access denied") }}
	if got := m.Observe(context.Background()); got.State != "unknown" {
		t.Fatalf("got %q, want unknown", got.State)
	}
}

func TestRequestActionUsesObservedState(t *testing.T) {
	t.Parallel()
	state := service.StateStopped
	m := Monitor{
		Query:  func() (service.State, error) { return state, nil },
		Health: func(context.Context) error { return nil },
		Execute: func(action string) error {
			if action != "start" {
				t.Fatalf("action = %q", action)
			}
			state = service.StateRunning
			return nil
		},
	}
	result := m.RequestAction(context.Background(), "start")
	if result.Outcome != "accepted" || result.Snapshot.State != "running" {
		t.Fatalf("result = %+v", result)
	}
}

func TestRequestActionDoesNotClaimElevationCancellationAsSuccess(t *testing.T) {
	t.Parallel()
	m := Monitor{
		Query:   func() (service.State, error) { return service.StateStopped, nil },
		Execute: func(string) error { return errElevationCancelled },
	}
	result := m.RequestAction(context.Background(), "start")
	if result.Outcome != "cancelled" || result.Snapshot.State != "stopped" {
		t.Fatalf("result = %+v", result)
	}
}
