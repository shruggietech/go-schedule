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
