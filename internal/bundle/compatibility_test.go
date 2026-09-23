package bundle

import (
	"reflect"
	"strings"
	"testing"
)

func TestTargetCompatibilityDerivesStableRequirements(t *testing.T) {
	doc := Document{Schema: SchemaV2, Tasks: []Task{{PortableID: "task", Name: "Task"}}, ExternalTriggers: []ExternalTrigger{{PortableID: "trigger", Name: "Trigger"}}, TriggerSets: []TriggerSet{{PortableID: "set", Name: "Set"}}, Watchers: []Watcher{{PortableID: "watcher", Name: "Watcher"}}, NotificationPolicies: []NotificationPolicy{{ScopeType: "task", ScopePortableID: "task"}}}
	profile := TargetProfile{Platform: "linux", Capabilities: []string{"tasks", "future-unknown"}}
	first := TargetCompatibility(doc, profile)
	second := TargetCompatibility(doc, profile)
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("compatibility order changed: %+v %+v", first, second)
	}
	want := []string{"bundles", "notifications", "schedule", "triggers", "watchers"}
	got := make([]string, 0, len(first))
	for _, issue := range first {
		if issue.Kind != "target_capability" || !strings.Contains(issue.Message, "selected daemon") {
			t.Fatalf("issue=%+v", issue)
		}
		got = append(got, issue.Identity)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("missing capabilities=%v want=%v", got, want)
	}
}

func TestTargetCompatibilityPlatformAndSchemaParity(t *testing.T) {
	v1 := Document{Schema: SchemaV1, Groups: []Group{{PortableID: "group", Name: "Group"}}}
	profile := TargetProfile{Platform: "unrecognized", Capabilities: []string{"bundles", "groups"}}
	if issues := TargetCompatibility(v1, profile); len(issues) != 0 {
		t.Fatalf("platform-neutral v1 was rejected: %+v", issues)
	}
	for _, platform := range []string{"windows", "linux", "macos"} {
		v2 := Document{Schema: SchemaV2, Watchers: []Watcher{{PortableID: "watcher", Name: "Watcher"}}}
		compatible := TargetProfile{Platform: platform, Capabilities: []string{"bundles", "watchers"}}
		if issues := TargetCompatibility(v2, compatible); len(issues) != 0 {
			t.Fatalf("platform %s rejected: %+v", platform, issues)
		}
		for _, unsupported := range []string{"", "unrecognized"} {
			compatible.Platform = unsupported
			issues := TargetCompatibility(v2, compatible)
			if len(issues) != 1 || issues[0].Kind != "target_platform" || issues[0].Identity != "watchers" {
				t.Fatalf("platform %q issues=%+v", unsupported, issues)
			}
		}
	}
	_, _, before, _ := Canonicalize(v1)
	_ = TargetCompatibility(v1, profile)
	_, _, after, _ := Canonicalize(v1)
	if before != after {
		t.Fatal("target compatibility changed the v1 digest")
	}
}

func TestPlanTargetCompatibilityChecksReviewedFamilies(t *testing.T) {
	plan := Plan{Items: []Item{{Kind: "watcher", Action: ActionCreate}, {Kind: "notification_policy", Action: ActionUpdate}, {Kind: "watcher", Action: ActionTargetOnly}}}
	issues := PlanTargetCompatibility(plan, TargetProfile{Platform: "linux", Capabilities: []string{"bundles"}})
	if len(issues) != 2 || issues[0].Identity != "notifications" || issues[1].Identity != "watchers" {
		t.Fatalf("plan issues=%+v", issues)
	}
}
