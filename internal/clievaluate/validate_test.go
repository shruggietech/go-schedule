package clievaluate

import (
	"strings"
	"testing"
)

func reportFixture(t *testing.T) (*ArtifactStore, Report, Manifest, Corpus) {
	t.Helper()
	store, err := OpenArtifacts(t.TempDir(), 1<<20, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Error(err)
		}
	})
	out, err := store.Write("env/health-out.txt", "text/plain", []byte("daemon ok\n"))
	if err != nil {
		t.Fatal(err)
	}
	stderr, err := store.Write("env/health-err.txt", "text/plain", nil)
	if err != nil {
		t.Fatal(err)
	}
	spec := CaseSpec{ID: "health.read", CommandPath: "health", ScenarioClass: "read", ExecutionKind: "packaged_cli", Platforms: []string{"windows", "linux"}, Argv: []string{"health"}, ExpectedExit: 0, Checks: []CheckSpec{{Name: "daemon health", Expected: "true"}}, Coverage: []string{"health"}}
	m := Manifest{SchemaVersion: Version, Kind: "go-schedule.cli-qa-manifest", Commands: []Command{{Path: "health", Runnable: true}}, Cases: []CaseSpec{spec}}
	md, err := Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	at := "2026-10-21T12:00:00Z"
	zero := 0
	result := CaseResult{ID: spec.ID, EnvironmentID: "env", CommandPath: spec.CommandPath, ScenarioClass: spec.ScenarioClass, ExecutionKind: spec.ExecutionKind, Status: "blocked", ExecutionStatus: "pass", ReasonCode: "pending_review", Argv: spec.Argv, ExpectedExit: 0, ExpectedChecks: spec.Checks, Observed: Observation{ExitCode: &zero, Checks: []Check{{Name: "daemon health", Expected: "true", Actual: "true", Passed: true, Evidence: []string{out.Path}}}}, StdoutArtifact: &out.Path, StderrArtifact: &stderr.Path, StartedAt: at, FinishedAt: at}
	c := Corpus{SchemaVersion: Version, Kind: "go-schedule.cli-qa-corpus", RunID: "run-1"}
	AddMessage(&c, result, "stdout", "daemon ok\n")
	SortCorpus(&c)
	cd, err := Encode(c)
	if err != nil {
		t.Fatal(err)
	}
	r := Report{SchemaVersion: Version, Kind: "go-schedule.cli-qa-report", RunID: "run-1", StartedAt: at, FinishedAt: at, Phase: "pending_review", CompletionState: "completed", CorpusSHA256: Digest(cd), Subject: Subject{Repository: "shruggietech/go-schedule", Commit: strings.Repeat("a", 40), Version: "dev"}, Runner: Runner{Version: Version, ManifestSHA256: Digest(md), Budget: DefaultBudget()}, Environments: []Environment{{ID: "env", OS: "windows", Arch: "amd64", Status: "ready"}}, Cases: []CaseResult{result}, Artifacts: []Artifact{out, stderr}, Cleanup: []Cleanup{{EnvironmentID: "env", Actions: []Check{{Name: "owned tree stopped", Expected: "true", Actual: "true", Passed: true, Evidence: []string{out.Path}}}}}}
	r.Summary = DeriveSummary(r)
	return store, r, m, c
}

func TestReportRejectsForgedPassAndMissingEvidence(t *testing.T) {
	store, r, m, c := reportFixture(t)
	if err := ValidateReport(store, r, m, c); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*Report){func(r *Report) { r.Summary.Passed++ }, func(r *Report) { r.Cases = nil }, func(r *Report) { r.Cases[0].Status = "pass" }, func(r *Report) { r.Cases[0].Observed.Checks[0].Actual = "false" }, func(r *Report) { r.Artifacts[0].SHA256 = strings.Repeat("0", 64) }, func(r *Report) { r.Phase = "finalized"; r.Cases[0].Status = "pass"; r.Summary = DeriveSummary(*r) }} {
		var clone Report
		data, err := Encode(r)
		if err != nil {
			t.Fatal(err)
		}
		if err := Decode(strings.NewReader(string(data)), &clone); err != nil {
			t.Fatal(err)
		}
		mutate(&clone)
		if err := ValidateReport(store, clone, m, c); err == nil {
			t.Fatal("contradictory report accepted")
		}
	}
}

func TestReportRejectsMissingExtraAndStaleReviews(t *testing.T) {
	_, r, _, c := reportFixture(t)
	if err := ApplyReviews(&r, c, ReviewSet{SchemaVersion: Version, Kind: "go-schedule.cli-qa-reviews", RunID: r.RunID, CorpusSHA256: r.CorpusSHA256}); err == nil {
		t.Fatal("missing review accepted")
	}
	m := c.Messages[0]
	review := Review{ID: "review-1", MessageID: m.ID, Cases: m.Cases, Channel: m.Channel, Excerpt: m.Text, Verdict: "pass"}
	for _, name := range rubricNames {
		review.Rubric = append(review.Rubric, Rubric{Name: name, Verdict: "pass", Reason: "Specific fixture review rationale"})
	}
	set := ReviewSet{SchemaVersion: Version, Kind: "go-schedule.cli-qa-reviews", RunID: r.RunID, CorpusSHA256: r.CorpusSHA256, Reviews: []Review{review}}
	stale := set
	stale.CorpusSHA256 = strings.Repeat("0", 64)
	if err := ApplyReviews(&r, c, stale); err == nil {
		t.Fatal("stale review accepted")
	}
	set.Reviews = append(set.Reviews, review)
	if err := ApplyReviews(&r, c, set); err == nil {
		t.Fatal("duplicate review accepted")
	}
	set.Reviews = set.Reviews[:1]
	if err := ApplyReviews(&r, c, set); err != nil {
		t.Fatal(err)
	}
	if r.Summary.Outcome != "pass" {
		t.Fatalf("reviewed outcome=%s", r.Summary.Outcome)
	}
}

func TestCorpusMustComeFromCapturedEvidence(t *testing.T) {
	store, r, m, c := reportFixture(t)
	c.Messages[0].Text = "invented favorable message"
	c.Messages[0].ID = "msg-" + Digest([]byte(c.Messages[0].CommandPath+"\x00"+c.Messages[0].ScenarioClass+"\x00"+c.Messages[0].Channel+"\x00"+c.Messages[0].Text))
	data, err := Encode(c)
	if err != nil {
		t.Fatal(err)
	}
	r.CorpusSHA256 = Digest(data)
	if err := ValidateReport(store, r, m, c); err == nil {
		t.Fatal("invented message accepted despite unrelated captured output")
	}
}

func TestInterruptedActiveCasePreservesChannelsWithoutInventingProductFailure(t *testing.T) {
	store, r, m, c := reportFixture(t)
	r.CompletionState = "interrupted"
	r.Cases[0].ExecutionStatus = "blocked"
	r.Cases[0].Status = "blocked"
	r.Cases[0].ReasonCode = "run_interrupted"
	r.Summary = DeriveSummary(r)
	if err := ValidateReport(store, r, m, c); err != nil {
		t.Fatal(err)
	}
	if r.Summary.ExecutionFailed != 0 || r.Summary.Blocked != 1 {
		t.Fatal("cancellation was classified as a product failure")
	}
	r.CompletionState = "completed"
	r.Summary = DeriveSummary(r)
	if err := ValidateReport(store, r, m, c); err == nil {
		t.Fatal("interrupted active case accepted in a completed run")
	}
}
