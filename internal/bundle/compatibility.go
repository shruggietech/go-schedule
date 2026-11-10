package bundle

import "sort"

// TargetProfile contains the daemon-advertised facts needed to assess one bundle operation.
type TargetProfile struct {
	Platform     string
	Capabilities []string
}

// TargetCompatibility derives requirements from bundle contents and reports missing target support in stable order.
func TargetCompatibility(doc Document, target TargetProfile) []Issue {
	kinds := make([]string, 0, 7)
	if len(doc.Joins) > 0 {
		kinds = append(kinds, "join")
		for _, join := range doc.Joins {
			if joinUsesRecovery(join) {
				kinds = append(kinds, "join_recovery")
				break
			}
		}
	}
	if len(doc.Chains) > 0 {
		kinds = append(kinds, "chain")
	}
	if len(doc.Groups) > 0 {
		kinds = append(kinds, "group")
	}
	if len(doc.SMTPChannels) > 0 {
		kinds = append(kinds, "smtp_channel")
	}
	if len(doc.NotificationPolicies) > 0 {
		kinds = append(kinds, "notification_policy")
	}
	if len(doc.Tasks) > 0 {
		kinds = append(kinds, "task")
	}
	if len(doc.ExternalTriggers) > 0 || len(doc.TriggerSets) > 0 {
		kinds = append(kinds, "external_trigger")
	}
	if len(doc.Watchers) > 0 {
		kinds = append(kinds, "watcher")
	}
	return targetCompatibilityKinds(kinds, target, nil)
}

// PlanTargetCompatibility checks the reviewed item families before a client forwards apply.
func PlanTargetCompatibility(plan Plan, target TargetProfile) []Issue {
	kinds := make([]string, 0, len(plan.Items))
	for _, item := range plan.Items {
		if item.Action != ActionTargetOnly {
			kinds = append(kinds, item.Kind)
		}
	}
	return targetCompatibilityKinds(kinds, target, plan.RequiredCapabilities)
}

func targetCompatibilityKinds(kinds []string, target TargetProfile, extra []string) []Issue {
	requiredSet := map[string]bool{"bundles": true}
	for _, capability := range extra {
		requiredSet[capability] = true
	}
	hasWatcher := false
	for _, kind := range kinds {
		switch kind {
		case "join_recovery":
			requiredSet["parallel-joins"] = true
			requiredSet["parallel-join-recovery"] = true
		case "join":
			requiredSet["parallel-joins"] = true
		case "chain":
			requiredSet["chains"] = true
		case "group":
			requiredSet["groups"] = true
		case "smtp_channel":
			requiredSet["notifications"] = true
			requiredSet["smtp_notifications"] = true
		case "notification_policy":
			requiredSet["notifications"] = true
		case "task":
			requiredSet["schedule"] = true
			requiredSet["tasks"] = true
		case "external_trigger", "trigger_set":
			requiredSet["triggers"] = true
		case "watcher":
			requiredSet["watchers"] = true
			hasWatcher = true
		}
	}
	required := make([]string, 0, len(requiredSet))
	for capability := range requiredSet {
		required = append(required, capability)
	}
	sort.Strings(required)
	available := make(map[string]bool, len(target.Capabilities))
	for _, capability := range target.Capabilities {
		available[capability] = true
	}
	issues := make([]Issue, 0)
	for _, capability := range required {
		if !available[capability] {
			issues = append(issues, Issue{Kind: "target_capability", Identity: capability, Message: "selected daemon does not advertise " + capability + "; select or upgrade a capable daemon"})
		}
	}
	if hasWatcher && target.Platform != "windows" && target.Platform != "linux" && target.Platform != "macos" {
		issues = append(issues, Issue{Kind: "target_platform", Identity: "watchers", Message: "selected daemon has an unknown or unsupported watcher platform; select a Windows, Linux, or macOS daemon"})
	}
	return issues
}

// RequiredCapabilities retains content-derived support requirements in a reviewed plan.
func RequiredCapabilities(doc Document) []string {
	issues := TargetCompatibility(doc, TargetProfile{})
	out := make([]string, 0, len(issues))
	for _, issue := range issues {
		if issue.Kind == "target_capability" {
			out = append(out, issue.Identity)
		}
	}
	return out
}
