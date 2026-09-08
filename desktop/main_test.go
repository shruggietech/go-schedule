package main

import (
	"context"
	"errors"
	"testing"
)

func TestEnsureBundledDaemonReusesReachableDaemon(t *testing.T) {
	spawned := false
	err := ensureBundledDaemon(func(context.Context) error { return nil }, func() error {
		spawned = true
		return nil
	})
	if err != nil {
		t.Fatalf("ensure bundled daemon: %v", err)
	}
	if spawned {
		t.Fatal("spawned a daemon when the existing endpoint was reachable")
	}
}

func TestEnsureBundledDaemonReportsSpawnFailure(t *testing.T) {
	want := errors.New("missing bundled daemon")
	err := ensureBundledDaemon(func(context.Context) error { return errors.New("unreachable") }, func() error { return want })
	if !errors.Is(err, want) {
		t.Fatalf("ensure bundled daemon error = %v, want %v", err, want)
	}
}
