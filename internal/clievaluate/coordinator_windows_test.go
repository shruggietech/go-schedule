//go:build windows

package clievaluate

import (
	"context"
	"encoding/json"
	"fmt"
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

func TestCanceledCaptureRetainsTerminalLedgerOrActionableRecovery(t *testing.T) {
	for _, quota := range []int64{64 << 20, 1 << 20} {
		t.Run(fmt.Sprint(quota), func(t *testing.T) {
			source, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			root := filepath.Join(t.TempDir(), "canceled")
			budget := DefaultBudget()
			budget.ArtifactBytes = quota
			options := CaptureOptions{Root: root, Source: source, BinDirectory: t.TempDir(), WSLDistro: "never-launched", WSLUser: "never-launched", Budget: budget}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			report, captureErr := Capture(ctx, options)
			if quota == 64<<20 {
				if captureErr != nil || report.CompletionState != "interrupted" || report.Phase != "pending_review" || report.Summary.Outcome != "blocked" {
					t.Fatalf("canceled ledger: %+v %v", report.Summary, captureErr)
				}
				if _, err := ValidateRun(root); err != nil {
					t.Fatal(err)
				}
			} else {
				if captureErr == nil {
					t.Fatal("report that exceeded its artifact quota advertised completion")
				}
				data, err := os.ReadFile(filepath.Join(root, "recovery.json"))
				if err != nil {
					t.Fatal(err)
				}
				var recovery struct {
					State    string         `json:"state"`
					Inputs   CaptureOptions `json:"original_inputs"`
					Recovery string         `json:"recovery"`
					Detail   string         `json:"detail"`
				}
				if err := json.Unmarshal(data, &recovery); err != nil {
					t.Fatal(err)
				}
				if recovery.State != "recovery_required" || recovery.Inputs != options || recovery.Recovery == "" || recovery.Detail == "" {
					t.Fatalf("recovery record lost original inputs: %s", data)
				}
			}
			if _, err := os.Stat(filepath.Join(root, "attempt-001")); !os.IsNotExist(err) {
				t.Fatalf("canceled preflight started a worker: %v", err)
			}
		})
	}
}
