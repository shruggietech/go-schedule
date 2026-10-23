package clievaluate

import (
	"bytes"
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
