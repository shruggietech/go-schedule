package bundle

import (
	"encoding/json"
	"strings"
	"testing"
)

func recoveryDocument(t *testing.T) Document {
	t.Helper()
	doc, err := Decode([]byte(`{"schema":"go-schedule.bundle/v4","tasks":[{"portable_id":"e","name":"entry"},{"portable_id":"b","name":"branch"},{"portable_id":"t","name":"target"},{"portable_id":"r","name":"recovery"},{"portable_id":"c","name":"cleanup"}],"joins":[{"portable_id":"j","name":"join","entry_task_portable_id":"e","target_task_portable_id":"t","recovery_task_portable_id":"r","recovery_on":[],"cleanup_task_portable_id":"c","branches":[{"id":"a","name":"a","task_portable_id":"b"}]}]}`))
	if err != nil {
		t.Fatalf("v4 routing configuration rejected: %v", err)
	}
	for i := range doc.Tasks {
		doc.Tasks[i].Timezone = "UTC"
		doc.Tasks[i].OverlapPolicy = "queue_one"
		doc.Tasks[i].CatchupPolicy = "one"
		doc.Tasks[i].MissingDatePolicy = "skip"
		doc.Tasks[i].TimeBasis = "wall_clock"
		doc.Tasks[i].DSTGapPolicy = "next_valid"
		doc.Tasks[i].DSTOverlapPolicy = "first"
	}
	return doc
}

func TestS145PortableV4RecoveryWireAndCapabilities(t *testing.T) {
	doc := recoveryDocument(t)
	canonical, data, _, issues := Canonicalize(doc)
	if len(issues) != 0 {
		t.Fatalf("v4 invalid: %+v", issues)
	}
	if !strings.Contains(string(data), `"recovery_on":[]`) {
		t.Fatalf("explicit empty causes lost: %s", data)
	}
	profile := TargetProfile{Platform: "windows", Capabilities: []string{"bundles", "tasks", "schedule", "parallel-joins"}}
	if issues := TargetCompatibility(canonical, profile); len(issues) != 1 || issues[0].Identity != "parallel-join-recovery" {
		t.Fatalf("new configuration requires specific capability: %+v", issues)
	}
	plan := Preview("p", "d", canonical, Document{Schema: SchemaV3})
	if issues := PlanTargetCompatibility(plan, profile); len(issues) != 1 || issues[0].Identity != "parallel-join-recovery" {
		t.Fatalf("reviewed plan lost capability requirement: %+v", issues)
	}
	profile.Capabilities = append(profile.Capabilities, "parallel-join-recovery")
	if issues := TargetCompatibility(canonical, profile); len(issues) != 0 {
		t.Fatalf("capable target rejected: %+v", issues)
	}
	var wire map[string]any
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatal(err)
	}
	wire["schema"] = SchemaV3
	legacy, err := json.Marshal(wire)
	if err != nil {
		t.Fatal(err)
	}
	older, err := Decode(legacy)
	if err != nil {
		t.Fatal(err)
	}
	if len(Validate(older)) == 0 {
		t.Fatal("new options silently accepted as v3")
	}
}

func TestS145PortableRecoveryCanonicalSetDriftAndActualCausalGraph(t *testing.T) {
	doc := recoveryDocument(t)
	causes := []string{"skipped", "failure", "failure"}
	doc.Joins[0].RecoveryOn = &causes
	canonical, _, digest, issues := Canonicalize(doc)
	if len(issues) > 0 || canonical.Joins[0].RecoveryOn == nil || strings.Join(*canonical.Joins[0].RecoveryOn, ",") != "failure,skipped" {
		t.Fatalf("canonical cause set=%+v issues=%+v", canonical.Joins, issues)
	}
	empty := []string{}
	canonical.Joins[0].RecoveryOn = &empty
	if Fingerprint(canonical) == digest {
		t.Fatal("empty condition selection omitted from drift")
	}
	canonical.Joins[0].CleanupTaskID = canonical.Joins[0].TargetTaskID
	if issues := Validate(canonical); len(issues) > 0 {
		t.Fatalf("finite stage reuse rejected: %+v", issues)
	}
	canonical.Chains = []Chain{{PortableID: "back", SourceTaskID: canonical.Joins[0].CleanupTaskID, TargetTaskID: canonical.Joins[0].EntryTaskID, OnOutcome: "success"}}
	if len(ValidateProjectedChains(canonical, Document{})) == 0 {
		t.Fatal("cleanup-to-entry automatic cycle accepted")
	}
	canonical.Chains = nil
	unknown := []string{"success"}
	canonical.Joins[0].RecoveryOn = &unknown
	if len(Validate(canonical)) == 0 {
		t.Fatal("unknown recovery cause accepted")
	}
}

func TestS145BaseV3WireOmitsUnconfiguredNewFields(t *testing.T) {
	doc := Document{Schema: SchemaV3, Joins: []Join{{PortableID: "j", Name: "join", EntryTaskID: "e", TargetTaskID: "t", Branches: []JoinBranch{{ID: "a", Name: "A", TaskID: "b"}}}}}
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "recovery_") || strings.Contains(string(data), "cleanup_") {
		t.Fatalf("base v3 wire changed: %s", data)
	}
}
