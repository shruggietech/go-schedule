//go:build windows

package clievaluate

import (
	"context"
	"os"
	"testing"
)

func TestResumeRefusesLiveOwnedLauncherWithoutTouchingIt(t *testing.T) {
	store, _, _, _ := reportFixture(t)
	identity, err := processIdentity(os.Getpid())
	if err != nil || identity == "" {
		t.Fatalf("live identity: %q %v", identity, err)
	}
	data, err := Encode(workerLease{PID: os.Getpid(), Creation: identity, RunID: "owned"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Write("attempt-001/windows/lease.json", "application/json", data); err != nil {
		t.Fatal(err)
	}
	if err := confirmPreviousWorkersStopped(context.Background(), store, captureIdentity{RunID: "owned", Attempt: 1}); err == nil {
		t.Fatal("live owned launcher accepted")
	}
	if actual, err := processIdentity(os.Getpid()); err != nil || actual != identity {
		t.Fatalf("probe changed process identity: %q %v", actual, err)
	}
	data, err = Encode(workerLease{PID: os.Getpid(), Creation: "different creation", RunID: "owned"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Replace("attempt-001/windows/lease.json", data); err != nil {
		t.Fatal(err)
	}
	if err := confirmPreviousWorkersStopped(context.Background(), store, captureIdentity{RunID: "owned", Attempt: 1}); err != nil {
		t.Fatalf("reused PID blocked: %v", err)
	}
}
