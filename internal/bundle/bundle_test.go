package bundle

import (
	"strings"
	"testing"
)

func TestCanonicalizeIsDeterministicAndRejectsUnresolvedReferences(t *testing.T) {
	doc := Document{Schema: SchemaV1, Groups: []Group{{PortableID: "g2", Name: "Two"}, {PortableID: "g1", Name: "One"}}, Tasks: []Task{{PortableID: "t1", Name: "Task", GroupPortableID: "g1", Timezone: "UTC", Schedule: "weekdays at 09:00", ScheduleSyntax: "human", OverlapPolicy: "queue_one", CatchupPolicy: "one", MissingDatePolicy: "skip", TimeBasis: "wall_clock", DSTGapPolicy: "next_valid", DSTOverlapPolicy: "first"}}}
	_, first, digest, issues := Canonicalize(doc)
	if len(issues) != 0 || digest == "" {
		t.Fatalf("issues=%v digest=%q", issues, digest)
	}
	doc.Groups[0], doc.Groups[1] = doc.Groups[1], doc.Groups[0]
	_, second, secondDigest, issues := Canonicalize(doc)
	if len(issues) != 0 || string(first) != string(second) || digest != secondDigest {
		t.Fatalf("canonical output changed: %s %s", first, second)
	}
	_, _, _, issues = Canonicalize(Document{Schema: SchemaV1, Tasks: []Task{{PortableID: "t", Name: "Task", GroupPortableID: "missing"}}})
	foundGroup := false
	for _, issue := range issues {
		foundGroup = foundGroup || strings.Contains(issue.Message, "group")
	}
	if len(issues) == 0 || !foundGroup {
		t.Fatalf("issues=%v", issues)
	}
}

func TestPreviewReportsTargetOnlyWithoutDelete(t *testing.T) {
	source := Document{Schema: SchemaV1, Groups: []Group{{PortableID: "source", Name: "Source"}}}
	target := Document{Schema: SchemaV1, Groups: []Group{{PortableID: "target", Name: "Target"}}}
	plan := Preview("plan", "daemon", source, target)
	if len(plan.Items) != 2 || plan.Items[0].Action != ActionCreate || plan.Items[1].Action != ActionTargetOnly {
		t.Fatalf("items=%+v", plan.Items)
	}
	for _, item := range plan.Items {
		if item.Action == "delete" {
			t.Fatal("preview inferred delete")
		}
	}
}

func TestDecodeRejectsUnknownAndTrailingBundleContent(t *testing.T) {
	for _, raw := range []string{
		`{"schema":"go-schedule.bundle/v1","command":"unsafe"}`,
		`{"schema":"go-schedule.bundle/v1"} {}`,
		`{"schema":"go-schedule.bundle/v1","tasks":[{"portable_id":"task","name":"task","timezone":"UTC","overlap_policy":"queue_one","catchup_policy":"one","missing_date_policy":"skip","time_basis":"wall_clock","dst_gap_policy":"next_valid","dst_overlap_policy":"first","command":"unsafe"}]}`,
	} {
		if _, err := Decode([]byte(raw)); err == nil {
			t.Fatalf("Decode accepted %s", raw)
		}
	}
}

func TestPreviewReportsSameNameAsIdentityConflict(t *testing.T) {
	source := Document{Schema: SchemaV1, Groups: []Group{{PortableID: "source", Name: "Shared"}}, Tasks: []Task{{PortableID: "source-task", Name: "Shared task", Timezone: "UTC", OverlapPolicy: "queue_one", CatchupPolicy: "one", MissingDatePolicy: "skip", TimeBasis: "wall_clock", DSTGapPolicy: "next_valid", DSTOverlapPolicy: "first"}}}
	target := Document{Schema: SchemaV1, Groups: []Group{{PortableID: "target", Name: "Shared"}}, Tasks: []Task{{PortableID: "target-task", Name: "Shared task", Timezone: "UTC", OverlapPolicy: "queue_one", CatchupPolicy: "one", MissingDatePolicy: "skip", TimeBasis: "wall_clock", DSTGapPolicy: "next_valid", DSTOverlapPolicy: "first"}}}
	plan := Preview("plan", "daemon", source, target)
	for _, item := range plan.Items {
		if (item.Kind == "group" || item.Kind == "task") && (item.PortableID == "source" || item.PortableID == "source-task") {
			if item.Action != ActionConflict {
				t.Fatalf("item %+v was not a name conflict", item)
			}
		}
	}
}

func TestValidateRejectsInvalidAndCyclicCompletionChains(t *testing.T) {
	doc := Document{Schema: SchemaV1, Tasks: []Task{
		{PortableID: "a", Name: "A", Timezone: "UTC", OverlapPolicy: "queue_one", CatchupPolicy: "one", MissingDatePolicy: "skip", TimeBasis: "wall_clock", DSTGapPolicy: "next_valid", DSTOverlapPolicy: "first"},
		{PortableID: "b", Name: "B", Timezone: "UTC", OverlapPolicy: "queue_one", CatchupPolicy: "one", MissingDatePolicy: "skip", TimeBasis: "wall_clock", DSTGapPolicy: "next_valid", DSTOverlapPolicy: "first"},
	}, Chains: []Chain{{PortableID: "a-b", SourceTaskID: "a", TargetTaskID: "b", OnOutcome: "unsupported"}, {PortableID: "b-a", SourceTaskID: "b", TargetTaskID: "a", OnOutcome: "any"}}}
	issues := Validate(doc)
	if len(issues) < 2 {
		t.Fatalf("issues=%+v, want outcome and cycle findings", issues)
	}
	if issues := ValidateProjectedChains(Document{Schema: SchemaV1, Chains: []Chain{{PortableID: "b-a", SourceTaskID: "b", TargetTaskID: "a", OnOutcome: "any"}}}, Document{Schema: SchemaV1, Chains: []Chain{{PortableID: "a-b", SourceTaskID: "a", TargetTaskID: "b", OnOutcome: "any"}}}); len(issues) == 0 {
		t.Fatal("projected chain cycle was accepted")
	}
}
