package clievaluate

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func frozenFixture(t *testing.T) (*ArtifactStore, Report, Manifest, Corpus) {
	t.Helper()
	store, r, m, c := reportFixture(t)
	manifest, _ := Encode(m)
	corpus, _ := Encode(c)
	identity, _ := Encode(frozenMetadata(r))
	id, env := r.Cases[0].ID, r.Environments[0].ID
	events := []byte{}
	for i, kind := range []string{"case_started", "case_finished"} {
		e := Event{Version, "go-schedule.cli-qa-event", r.RunID, int64(i + 1), r.StartedAt, kind, &env, &id, "fixture transition"}
		data, err := Encode(e)
		if err != nil {
			t.Fatal(err)
		}
		data, err = compactJSON(data)
		if err != nil {
			t.Fatal(err)
		}
		events = append(events, data...)
		events = append(events, '\n')
	}
	for _, document := range []struct {
		name, media string
		data        []byte
	}{{"manifest.json", "application/json", manifest}, {"corpus.json", "application/json", corpus}, {"events.jsonl", "application/x-ndjson", events}, {"capture-identity.json", "application/json", identity}} {
		a, err := store.Write(document.name, document.media, document.data)
		if err != nil {
			t.Fatal(err)
		}
		r.Artifacts = append(r.Artifacts, a)
	}
	frozen := EvidenceManifest{Version, "go-schedule.cli-qa-evidence-manifest", r.RunID, r.CorpusSHA256, r.Artifacts}
	data, err := Encode(frozen)
	if err != nil {
		t.Fatal(err)
	}
	r.Runner.EvidenceSHA256 = Digest(data)
	if err := store.Replace("evidence-manifest.json", data); err != nil {
		t.Fatal(err)
	}
	if err := saveReport(store, r); err != nil {
		t.Fatal(err)
	}
	return store, r, m, c
}

func TestFrozenIdentityRejectsRewrittenSourceAndResourceClaims(t *testing.T) {
	store, report, _, _ := frozenFixture(t)
	for _, change := range []func(*Report){func(r *Report) { r.Subject.Commit = "different" }, func(r *Report) { r.Runner.Budget.CPUPercent++ }, func(r *Report) { r.Environments[0].Version = "invented" }, func(r *Report) { r.CompletionState = "interrupted" }} {
		data, err := Encode(report)
		if err != nil {
			t.Fatal(err)
		}
		var altered Report
		if err := Decode(bytes.NewReader(data), &altered); err != nil {
			t.Fatal(err)
		}
		change(&altered)
		if err := validateFrozen(store, altered); err == nil {
			t.Fatal("rewritten immutable capture identity accepted")
		}
	}
}

func TestFrozenRunRoundTripAndReadOnlyValidation(t *testing.T) {
	store, r, _, c := frozenFixture(t)
	root := store.root.Name()
	before, err := os.ReadFile(filepath.Join(root, "report.json"))
	if err != nil {
		t.Fatal(err)
	}
	validated, err := ValidateRun(root)
	if err != nil || validated.RunID != r.RunID {
		t.Fatalf("validate: %v", err)
	}
	after, err := os.ReadFile(filepath.Join(root, "report.json"))
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("validate changed the canonical report")
	}
	message := c.Messages[0]
	review := Review{ID: "review-1", MessageID: message.ID, Cases: message.Cases, Channel: message.Channel, Excerpt: message.Text, Verdict: "pass"}
	for _, name := range rubricNames {
		review.Rubric = append(review.Rubric, Rubric{Name: name, Verdict: "pass", Reason: "The fixture says daemon ok and contains no secret or unhelpful correction."})
	}
	set := ReviewSet{Version, "go-schedule.cli-qa-reviews", r.RunID, r.CorpusSHA256, []Review{review}, nil}
	data, err := Encode(set)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "reviews.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	final, err := ProcessReport(root, path)
	if err != nil {
		t.Fatal(err)
	}
	if final.Phase != "finalized" || final.Summary.Outcome != "pass" {
		t.Fatalf("final outcome=%+v", final.Summary)
	}
	if _, err := ValidateRun(root); err != nil {
		t.Fatal(err)
	}
}

func TestJournalRejectsTruncationReorderingAndUnknownCases(t *testing.T) {
	store, r, _, _ := frozenFixture(t)
	data, err := store.readRegular("events.jsonl", documentLimit)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateEvents(data, r); err != nil {
		t.Fatal(err)
	}
	for _, bad := range [][]byte{bytes.TrimSuffix(data, []byte("\n")), bytes.Replace(data, []byte(`"sequence":2`), []byte(`"sequence":3`), 1), bytes.ReplaceAll(data, []byte("health.read"), []byte("unknown-case"))} {
		if err := validateEvents(bad, r); err == nil {
			t.Fatal("corrupt journal accepted")
		}
	}
}

func TestFrozenChannelsReconcileLateSecretAndRejectOmittedOutput(t *testing.T) {
	store, r, m, c := reportFixture(t)
	store.Redactor.Register("daemon ok")
	frozen, err := freezeChannels(store, &r)
	if err != nil {
		t.Fatal(err)
	}
	if frozen.Messages[0].Text != "[REDACTED]\n" {
		t.Fatal("late fixture registry did not reconcile retained evidence")
	}
	data, err := Encode(frozen)
	if err != nil {
		t.Fatal(err)
	}
	r.CorpusSHA256 = Digest(data)
	if err := ValidateReport(store, r, m, frozen); err != nil {
		t.Fatal(err)
	}
	_ = c
	extra, err := store.Write("env/extra.txt", "text/plain", []byte("second diagnostic"))
	if err != nil {
		t.Fatal(err)
	}
	r.Artifacts = append(r.Artifacts, extra)
	r.Cases[0].StderrArtifact = &extra.Path
	if err := ValidateReport(store, r, m, frozen); err == nil {
		t.Fatal("nonempty channel omitted from review corpus")
	}
}
