package clievaluate

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestIncrementalCheckpointRecoversCompletedCasesAndRejectsCorruption(t *testing.T) {
	store, report, _, _ := reportFixture(t)
	if err := store.root.MkdirAll("worker", 0700); err != nil {
		t.Fatal(err)
	}
	completed := report.Cases[0]
	report.Phase = "capturing"
	report.CompletionState = "running"
	report.Cases[0].ExecutionStatus = "blocked"
	report.Cases[0].Observed = Observation{}
	report.Cases[0].ReasonCode = "not_started"
	report.Cases[0].StdoutArtifact = nil
	report.Cases[0].StderrArtifact = nil
	artifacts := append([]Artifact{}, report.Artifacts...)
	report.Artifacts = nil
	data, err := Encode(report)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Replace("worker/report.json", data); err != nil {
		t.Fatal(err)
	}
	if err := persistCase(store, "worker", completed, artifacts); err != nil {
		t.Fatal(err)
	}
	recovered, err := loadWorkerCheckpoint(store, "worker")
	if err != nil {
		t.Fatal(err)
	}
	if recovered.Cases[0].ExecutionStatus != "pass" || len(recovered.Artifacts) != len(artifacts) {
		t.Fatal("completed case lost between full checkpoints")
	}
	unchanged, err := store.readRegular("worker/report.json", documentLimit)
	if err != nil || !bytes.Equal(unchanged, data) {
		t.Fatal("incremental update rewrote the complete ledger")
	}
	if err := store.Replace("worker/progress/health.read.json", []byte(`{"case":`)); err != nil {
		t.Fatal(err)
	}
	if _, err := loadWorkerCheckpoint(store, "worker"); err == nil {
		t.Fatal("truncated case checkpoint accepted")
	}
}

func TestJournalAppendPreservesPrefixAndEnforcesBudgetAndRedaction(t *testing.T) {
	store, _, _, _ := reportFixture(t)
	store.Redactor.Register("owned-secret")
	for _, data := range []string{"first\n", "owned-secret\n"} {
		if err := store.AppendJournal("journal.jsonl", []byte(data)); err != nil {
			t.Fatal(err)
		}
	}
	data, err := store.readRegular("journal.jsonl", 1024)
	if err != nil || string(data) != "first\n[REDACTED]\n" {
		t.Fatalf("journal=%q err=%v", data, err)
	}
	store.limit = store.used
	if err := store.AppendJournal("journal.jsonl", []byte("overflow\n")); err == nil {
		t.Fatal("journal exceeded budget")
	}
	if err := store.AppendJournal("../outside", nil); err == nil {
		t.Fatal("journal escaped artifact root")
	}
}

func TestFrozenProgressUpdatesLateRedactedArtifactIdentities(t *testing.T) {
	store, r, _, _ := reportFixture(t)
	if err := store.root.MkdirAll("worker", 0700); err != nil {
		t.Fatal(err)
	}
	if err := persistCase(store, "worker", r.Cases[0], r.Artifacts); err != nil {
		t.Fatal(err)
	}
	store.Redactor.Register("daemon ok")
	if _, err := freezeChannels(store, &r); err != nil {
		t.Fatal(err)
	}
	if err := freezeProgress(store, "worker", r); err != nil {
		t.Fatal(err)
	}
	data, err := store.readRegular("worker/progress/"+r.Cases[0].ID+".json", documentLimit)
	if err != nil {
		t.Fatal(err)
	}
	var saved caseCheckpoint
	if err := Decode(bytes.NewReader(data), &saved); err != nil {
		t.Fatal(err)
	}
	for _, a := range saved.Artifacts {
		if _, err := store.Read(a); err != nil {
			t.Fatalf("frozen checkpoint retained stale identity: %v", err)
		}
	}
}

func TestLatestWorkerCheckpointKeepsPostCaseCleanupEvidence(t *testing.T) {
	store, r, _, _ := reportFixture(t)
	if err := store.root.MkdirAll("worker", 0700); err != nil {
		t.Fatal(err)
	}
	if err := persistCase(store, "worker", r.Cases[0], r.Artifacts); err != nil {
		t.Fatal(err)
	}
	finalLog, err := store.Write("worker/final-daemon.log", "text/plain", []byte("owned daemon stopped\n"))
	if err != nil {
		t.Fatal(err)
	}
	r.Artifacts = append(r.Artifacts, finalLog)
	r.Cases[0].LogArtifacts = append(r.Cases[0].LogArtifacts, finalLog.Path)
	r.Phase = "capturing"
	r.CompletionState = "running"
	r.Cleanup = []Cleanup{{EnvironmentID: r.Cases[0].EnvironmentID, Remaining: []string{}, Actions: []Check{{Name: "owned runtime removed", Expected: "true", Actual: "true", Passed: true, Evidence: []string{finalLog.Path}}}}}
	data, err := Encode(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Replace("worker/report.json", data); err != nil {
		t.Fatal(err)
	}
	recovered, err := loadWorkerCheckpoint(store, "worker")
	if err != nil {
		t.Fatal(err)
	}
	if len(recovered.Cases[0].LogArtifacts) != 1 || recovered.Cases[0].LogArtifacts[0] != finalLog.Path {
		t.Fatal("older per-case progress erased the later cleanup log")
	}
	if len(recovered.Cleanup[0].Remaining) != 0 || !recovered.Cleanup[0].Actions[0].Passed {
		t.Fatal("latest cleanup result was lost")
	}
}

func TestWorkerFreezeFailureRetainsLatestCleanupCheckpoint(t *testing.T) {
	store, r, _, _ := reportFixture(t)
	if err := store.root.MkdirAll("worker", 0700); err != nil {
		t.Fatal(err)
	}
	r.Phase = "capturing"
	r.CompletionState = "running"
	r.Cleanup = []Cleanup{{EnvironmentID: r.Cases[0].EnvironmentID, Remaining: []string{}, Actions: []Check{{Name: "owned runtime removed", Expected: "true", Actual: "true", Passed: true, Evidence: []string{r.Artifacts[0].Path}}}}}
	if err := store.Replace(r.Artifacts[0].Path, []byte("tampered")); err != nil {
		t.Fatal(err)
	}
	_, err := freezeWorkerEvidence(context.Background(), store, "worker", &r)
	if err == nil || !strings.Contains(err.Error(), "freeze worker channels") {
		t.Fatalf("stage-specific failure missing: %v", err)
	}
	var saved Report
	if err := store.ReadDocument("worker/report.json", "report", &saved); err != nil {
		t.Fatal(err)
	}
	if saved.Phase != "capturing" || saved.CompletionState != "running" {
		t.Fatal("failed freeze claimed terminal completion")
	}
	if len(saved.Cleanup[0].Remaining) != 0 || !saved.Cleanup[0].Actions[0].Passed {
		t.Fatal("completed cleanup was not durable before failed freezing")
	}
	if saved.Cases[0].Observed.ExitCode == nil || *saved.Cases[0].Observed.ExitCode != 0 {
		t.Fatal("completed command outcome was lost")
	}
}
