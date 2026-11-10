package bundle

import (
	"strings"
	"testing"
)

func TestS144PortableJoinsRequireV3AndParticipateInDrift(t *testing.T) {
	doc := Document{Schema: SchemaV3, Tasks: []Task{{PortableID: "e", Name: "entry"}, {PortableID: "b", Name: "branch"}, {PortableID: "t", Name: "target"}}, Joins: []Join{{PortableID: "j", Name: "join", EntryTaskID: "e", TargetTaskID: "t", Branches: []JoinBranch{{ID: "a", Name: "first", TaskID: "b"}, {ID: "b", Name: "second", TaskID: "b"}}}}}
	for i := range doc.Tasks {
		doc.Tasks[i].Timezone = "UTC"
		doc.Tasks[i].OverlapPolicy = "queue_one"
		doc.Tasks[i].CatchupPolicy = "one"
		doc.Tasks[i].MissingDatePolicy = "skip"
		doc.Tasks[i].TimeBasis = "wall_clock"
		doc.Tasks[i].DSTGapPolicy = "next_valid"
		doc.Tasks[i].DSTOverlapPolicy = "first"
	}
	_, data, digest, issues := Canonicalize(doc)
	if len(issues) > 0 {
		t.Fatalf("issues=%+v", issues)
	}
	if !strings.Contains(string(data), "task_portable_id") {
		t.Fatalf("portable=%s", data)
	}
	doc.Joins[0].Branches[1].Name = "different"
	if Fingerprint(doc) == digest {
		t.Fatal("join changes omitted from drift")
	}
	doc.Schema = SchemaV2
	if len(Validate(doc)) == 0 {
		t.Fatal("older schema accepted joins")
	}
}
func TestS144PortableMixedGraphRejectsCycle(t *testing.T) {
	doc := Document{Schema: SchemaV3, Joins: []Join{{PortableID: "j", EntryTaskID: "e", TargetTaskID: "t", Branches: []JoinBranch{{ID: "a", Name: "a", TaskID: "b"}}}}, Chains: []Chain{{PortableID: "c", SourceTaskID: "t", TargetTaskID: "e", OnOutcome: "success"}}}
	if len(ValidateProjectedChains(doc, Document{})) == 0 {
		t.Fatal("mixed cycle accepted")
	}
}
