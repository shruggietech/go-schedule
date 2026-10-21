package clievaluate

import (
	"bytes"
	"strings"
	"testing"
)

func TestResumeRejectsEveryChangedFrozenInput(t *testing.T) {
	prior := captureIdentity{RunID: "prior", Subject: Subject{Repository: "repo", Commit: strings.Repeat("a", 40), Version: "dev", Binaries: []Binary{{Role: "cli", OS: "windows", Arch: "amd64", Path: "C:\\owned.exe", SHA256: strings.Repeat("b", 64)}}}, Budget: DefaultBudget(), Distro: "fixture", User: "owned", ManifestSHA256: strings.Repeat("c", 64), Attempt: 1}
	for _, change := range []func(*captureIdentity){func(v *captureIdentity) { v.Subject.Commit = strings.Repeat("d", 40) }, func(v *captureIdentity) { v.Subject.Dirty = true }, func(v *captureIdentity) { v.Subject.Binaries[0].SHA256 = strings.Repeat("e", 64) }, func(v *captureIdentity) { v.Budget.CPUPercent++ }, func(v *captureIdentity) { v.Distro = "other" }, func(v *captureIdentity) { v.User = "other" }, func(v *captureIdentity) { v.ManifestSHA256 = strings.Repeat("f", 64) }} {
		data, err := Encode(prior)
		if err != nil {
			t.Fatal(err)
		}
		var current captureIdentity
		if err := Decode(bytes.NewReader(data), &current); err != nil {
			t.Fatal(err)
		}
		if err := compareResumeInputs(prior, current); err != nil {
			t.Fatal(err)
		}
		change(&current)
		if err := compareResumeInputs(prior, current); err == nil {
			t.Fatal("changed input was accepted")
		}
	}
}

func TestInterruptedCheckpointRetainsTruthfulAggregate(t *testing.T) {
	_, r, _, _ := reportFixture(t)
	r.Phase = "finalized"
	r.CompletionState = "interrupted"
	r.Cases[0].Status = "pass"
	if summary := DeriveSummary(r); summary.Outcome != "blocked" || summary.ExecutionPassed != 1 {
		t.Fatalf("interrupted result lost execution truth: %+v", summary)
	}
}

func TestResumeRequiresEstablishedStoppedUnit(t *testing.T) {
	for _, test := range []struct {
		code     int
		state    string
		accepted bool
	}{{0, "inactive", true}, {0, "failed", true}, {1, "", false}, {1, "inactive", false}, {0, "", false}, {0, "active", false}, {0, "deactivating", false}} {
		if got := validateStoppedUnit(test.code, test.state) == nil; got != test.accepted {
			t.Errorf("code=%d state=%q accepted=%v", test.code, test.state, got)
		}
	}
}

func TestResumePreservesPriorCleanupObligationsAndOwnershipEvidence(t *testing.T) {
	prior := Report{Cleanup: []Cleanup{{EnvironmentID: "windows", Remaining: []string{"owned/runtime"}, Actions: []Check{{Name: "stopped", Evidence: []string{"old/ownership.json"}}}}}, Artifacts: []Artifact{{Path: "old/ownership.json"}, {Path: "old/stdout.txt"}}}
	cleanup, artifacts := resumeCleanup(prior)
	current := mergeCleanup(cleanup["windows"], Cleanup{EnvironmentID: "windows", Retained: []string{"new/evidence"}})
	if len(current.Remaining) != 1 || current.Remaining[0] != "owned/runtime" || len(artifacts) != 1 || artifacts[0].Path != "old/ownership.json" {
		t.Fatal("resume discarded prior cleanup obligations")
	}
}
