package releasegate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCurrentEvidenceContractRejectsHistoricalMatrices(t *testing.T) {
	for _, id := range RequiredScenarioIDs() {
		if id == "setup.upgrade" || id == "window.retained-profile" {
			t.Errorf("historical scenario remains required: %s", id)
		}
	}
	root, artifact, e := passingEvidence(t)
	e.SchemaVersion = 2
	var observations []Observation
	for _, o := range e.Observations {
		if o.ID != "setup.upgrade" && o.ID != "window.retained-profile" {
			observations = append(observations, o)
		}
	}
	e.Observations = observations
	var environments []Environment
	for _, env := range e.Environments {
		if !strings.HasPrefix(env.ProfileState, "retained-") {
			environments = append(environments, env)
		}
	}
	e.Environments = environments
	if failures := Validate(e, root, artifact, ExpectedIdentity{}); len(failures) != 0 {
		t.Fatalf("current packet requires historical evidence: %v", failures)
	}
	e.SchemaVersion = 1
	if failures := Validate(e, root, artifact, ExpectedIdentity{}); !strings.Contains(strings.Join(failures, " "), "schema_version") {
		t.Fatal("historical schema1 accepted as current qualification")
	}
}

func TestCurrentSyntheticPacketIsCompleteAndExplicitlyNonNative(t *testing.T) {
	root := filepath.Join("testdata", "current")
	file, err := os.Open(filepath.Join(root, "evidence.json"))
	if err != nil {
		t.Fatal(err)
	}
	e, err := DecodeEvidence(file)
	file.Close()
	if err != nil {
		t.Fatal(err)
	}
	artifact := filepath.Join(root, e.Candidate.Filename)
	if failures := ValidateFixture(e, root, artifact, ExpectedIdentity{}); len(failures) != 0 {
		t.Fatalf("current synthetic fixture invalid: %v", failures)
	}
	if failures := Validate(e, root, artifact, ExpectedIdentity{}); !containsFailure(failures, "evidence_class") {
		t.Fatal("synthetic packet accepted as native qualification")
	}
}

func TestCurrentEvidenceRejectsRetiredObservationAndProfile(t *testing.T) {
	for _, id := range []string{"setup.upgrade", "window.retained-profile"} {
		t.Run(id, func(t *testing.T) {
			root, artifact, e := passingEvidence(t)
			o := e.Observations[0]
			o.ID = id
			e.Observations = append(e.Observations, o)
			if failures := Validate(e, root, artifact, ExpectedIdentity{}); !containsFailure(failures, "is not a required scenario") {
				t.Fatalf("retired observation accepted: %v", failures)
			}
		})
	}
	root, artifact, e := passingEvidence(t)
	e.Environments[0].ProfileState = "retained-v1.1.1"
	if failures := Validate(e, root, artifact, ExpectedIdentity{}); !containsFailure(failures, "profile_state") {
		t.Fatalf("historical profile accepted: %v", failures)
	}
}
