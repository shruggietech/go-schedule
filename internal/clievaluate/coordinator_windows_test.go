//go:build windows

package clievaluate

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestUnavailableInputsStillPersistCompleteCrossEnvironmentLedger(t *testing.T) {
	source, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	directory := filepath.Join(t.TempDir(), "capture")
	r, err := Capture(context.Background(), CaptureOptions{Root: directory, Source: source, BinDirectory: t.TempDir(), WSLDistro: "not-launched", WSLUser: "not-launched", Budget: DefaultBudget()})
	if err != nil {
		t.Fatal(err)
	}
	m, _, err := BuildManifest()
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Environments) != 2 || len(r.Cases) != 2*len(m.Cases) || r.CompletionState != "interrupted" || r.Summary.Outcome != "blocked" {
		t.Fatalf("incomplete unavailable-input ledger: %+v", r.Summary)
	}
	for _, c := range r.Cases {
		if c.ExecutionStatus != "blocked" && c.ExecutionStatus != "not_applicable" {
			t.Fatalf("unavailable input executed %s", c.ID)
		}
	}
	validated, err := ValidateRun(directory)
	if err != nil {
		t.Fatal(err)
	}
	if validated.Summary != r.Summary {
		t.Fatal("persisted ledger changed summary")
	}
}
