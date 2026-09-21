// Package bundle defines the versioned, secret-free portable automation format.
package bundle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// SchemaV1 is the only portable bundle schema accepted by this release.
const SchemaV1 = "go-schedule.bundle/v1"

// Document carries portable automation intent. It deliberately has no daemon
// identity, credentials, raw trigger keys, notification endpoints, run history,
// command, working directory, environment, or run-as identity.
type Document struct {
	Schema     string  `json:"schema"`
	Groups     []Group `json:"groups,omitempty"`
	Tasks      []Task  `json:"tasks,omitempty"`
	Chains     []Chain `json:"chains,omitempty"`
	Exclusions []Issue `json:"exclusions,omitempty"`
}

type Group struct {
	PortableID       string `json:"portable_id"`
	Name             string `json:"name"`
	ParentPortableID string `json:"parent_portable_id,omitempty"`
	Enabled          bool   `json:"enabled"`
}

// Task is only the portable scheduling intent. Execution inputs are never part
// of a v1 document and their omission is disclosed by an exclusion.
type Task struct {
	PortableID        string `json:"portable_id"`
	Name              string `json:"name"`
	GroupPortableID   string `json:"group_portable_id,omitempty"`
	Enabled           bool   `json:"enabled"`
	Timezone          string `json:"timezone"`
	Schedule          string `json:"schedule,omitempty"`
	ScheduleSyntax    string `json:"schedule_syntax,omitempty"`
	OverlapPolicy     string `json:"overlap_policy"`
	CatchupPolicy     string `json:"catchup_policy"`
	MissingDatePolicy string `json:"missing_date_policy"`
	TimeBasis         string `json:"time_basis"`
	DSTGapPolicy      string `json:"dst_gap_policy"`
	DSTOverlapPolicy  string `json:"dst_overlap_policy"`
}

type Chain struct {
	PortableID   string `json:"portable_id"`
	SourceTaskID string `json:"source_task_portable_id"`
	TargetTaskID string `json:"target_task_portable_id"`
	OnOutcome    string `json:"on_outcome"`
}

// Issue names an intentional exclusion, validation problem, or plan finding.
type Issue struct {
	Kind     string `json:"kind"`
	Identity string `json:"identity,omitempty"`
	Message  string `json:"message"`
}

// Action is a non-destructive operation described by a target-bound plan.
type Action string

const (
	ActionCreate     Action = "create"
	ActionUpdate     Action = "update"
	ActionUnchanged  Action = "unchanged"
	ActionConflict   Action = "conflict"
	ActionInvalid    Action = "invalid"
	ActionTargetOnly Action = "target_only"
	ActionApplied    Action = "applied"
	ActionFailed     Action = "failed"
	ActionUncertain  Action = "uncertain"
)

type Item struct {
	Kind       string `json:"kind"`
	PortableID string `json:"portable_id"`
	Name       string `json:"name"`
	Action     Action `json:"action"`
	Message    string `json:"message,omitempty"`
}

// Plan records a deterministic, reviewed comparison of one document and one
// selected daemon snapshot. Apply must echo every binding value.
type Plan struct {
	ID                string `json:"id"`
	BundleDigest      string `json:"bundle_digest"`
	TargetDaemonID    string `json:"target_daemon_id"`
	TargetFingerprint string `json:"target_fingerprint"`
	Items             []Item `json:"items"`
}

// Canonicalize normalizes a document and rejects unsafe or structurally invalid
// content before returning its stable JSON bytes and SHA-256 digest.
func Canonicalize(doc Document) (Document, []byte, string, []Issue) {
	doc.Schema = strings.TrimSpace(doc.Schema)
	issues := Validate(doc)
	sort.Slice(doc.Groups, func(i, j int) bool { return doc.Groups[i].PortableID < doc.Groups[j].PortableID })
	sort.Slice(doc.Tasks, func(i, j int) bool { return doc.Tasks[i].PortableID < doc.Tasks[j].PortableID })
	sort.Slice(doc.Chains, func(i, j int) bool { return doc.Chains[i].PortableID < doc.Chains[j].PortableID })
	sort.Slice(doc.Exclusions, func(i, j int) bool { return issueKey(doc.Exclusions[i]) < issueKey(doc.Exclusions[j]) })
	bytes, err := json.Marshal(doc)
	if err != nil {
		issues = append(issues, Issue{Kind: "document", Message: "bundle cannot be encoded"})
		return doc, nil, "", issues
	}
	sum := sha256.Sum256(bytes)
	return doc, bytes, hex.EncodeToString(sum[:]), issues
}

func issueKey(issue Issue) string {
	return issue.Kind + "\x00" + issue.Identity + "\x00" + issue.Message
}

// Validate checks the v1 format, durable identities, and cross-object references.
func Validate(doc Document) []Issue {
	issues := []Issue{}
	if doc.Schema != SchemaV1 {
		issues = append(issues, Issue{Kind: "schema", Message: fmt.Sprintf("unsupported bundle schema %q", doc.Schema)})
	}
	seen := map[string]string{}
	add := func(kind, id, name string) {
		if id == "" {
			issues = append(issues, Issue{Kind: kind, Message: "portable_id is required"})
			return
		}
		if prior, ok := seen[id]; ok {
			issues = append(issues, Issue{Kind: kind, Identity: id, Message: "portable_id duplicates " + prior})
			return
		}
		seen[id] = kind
		if strings.TrimSpace(name) == "" {
			issues = append(issues, Issue{Kind: kind, Identity: id, Message: "name is required"})
		}
	}
	groups := map[string]bool{}
	tasks := map[string]bool{}
	for _, group := range doc.Groups {
		add("group", group.PortableID, group.Name)
		groups[group.PortableID] = true
	}
	for _, task := range doc.Tasks {
		add("task", task.PortableID, task.Name)
		tasks[task.PortableID] = true
		if strings.TrimSpace(task.Timezone) == "" {
			issues = append(issues, Issue{Kind: "task", Identity: task.PortableID, Message: "timezone is required"})
		}
		for _, field := range []struct{ name, value string }{{"overlap_policy", task.OverlapPolicy}, {"catchup_policy", task.CatchupPolicy}, {"missing_date_policy", task.MissingDatePolicy}, {"time_basis", task.TimeBasis}, {"dst_gap_policy", task.DSTGapPolicy}, {"dst_overlap_policy", task.DSTOverlapPolicy}} {
			if strings.TrimSpace(field.value) == "" {
				issues = append(issues, Issue{Kind: "task", Identity: task.PortableID, Message: field.name + " is required"})
			}
		}
	}
	for _, group := range doc.Groups {
		if group.ParentPortableID != "" && !groups[group.ParentPortableID] {
			issues = append(issues, Issue{Kind: "group", Identity: group.PortableID, Message: "parent group is absent from bundle"})
		}
	}
	for _, task := range doc.Tasks {
		if task.GroupPortableID != "" && !groups[task.GroupPortableID] {
			issues = append(issues, Issue{Kind: "task", Identity: task.PortableID, Message: "group is absent from bundle"})
		}
		if task.Schedule != "" && task.ScheduleSyntax == "" {
			issues = append(issues, Issue{Kind: "task", Identity: task.PortableID, Message: "schedule syntax is required when schedule is present"})
		}
		if task.ScheduleSyntax != "" && task.ScheduleSyntax != "human" && task.ScheduleSyntax != "cron" {
			issues = append(issues, Issue{Kind: "task", Identity: task.PortableID, Message: "schedule syntax must be human or cron"})
		}
	}
	for _, chain := range doc.Chains {
		add("chain", chain.PortableID, chain.PortableID)
		if !tasks[chain.SourceTaskID] || !tasks[chain.TargetTaskID] {
			issues = append(issues, Issue{Kind: "chain", Identity: chain.PortableID, Message: "chain task reference is absent from bundle"})
		}
		if chain.SourceTaskID == chain.TargetTaskID {
			issues = append(issues, Issue{Kind: "chain", Identity: chain.PortableID, Message: "chain source and target must differ"})
		}
	}
	return issues
}

// Fingerprint returns the stable digest of a target's exported safe snapshot.
func Fingerprint(target Document) string { _, _, digest, _ := Canonicalize(target); return digest }

// Preview compares one validated bundle to one selected target snapshot without
// mutation. Target-only state is visible drift, never an inferred deletion.
func Preview(id, daemonID string, source, target Document) Plan {
	canonical, _, digest, issues := Canonicalize(source)
	plan := Plan{ID: id, BundleDigest: digest, TargetDaemonID: daemonID, TargetFingerprint: Fingerprint(target), Items: make([]Item, 0)}
	if len(issues) > 0 {
		for _, issue := range issues {
			plan.Items = append(plan.Items, Item{Kind: issue.Kind, PortableID: issue.Identity, Action: ActionInvalid, Message: issue.Message})
		}
		return plan
	}
	targetGroups := mapByID(target.Groups, func(v Group) string { return v.PortableID })
	for _, v := range canonical.Groups {
		plan.Items = append(plan.Items, compareGroup(v, targetGroups[v.PortableID]))
	}
	targetTasks := mapByID(target.Tasks, func(v Task) string { return v.PortableID })
	for _, v := range canonical.Tasks {
		plan.Items = append(plan.Items, compareTask(v, targetTasks[v.PortableID]))
	}
	targetChains := mapByID(target.Chains, func(v Chain) string { return v.PortableID })
	for _, v := range canonical.Chains {
		plan.Items = append(plan.Items, compareChain(v, targetChains[v.PortableID]))
	}
	for _, group := range target.Groups {
		if _, found := groupByID(canonical.Groups, group.PortableID); !found {
			plan.Items = append(plan.Items, Item{Kind: "group", PortableID: group.PortableID, Name: group.Name, Action: ActionTargetOnly, Message: "exists only on target"})
		}
	}
	for _, task := range target.Tasks {
		if _, found := taskByID(canonical.Tasks, task.PortableID); !found {
			plan.Items = append(plan.Items, Item{Kind: "task", PortableID: task.PortableID, Name: task.Name, Action: ActionTargetOnly, Message: "exists only on target"})
		}
	}
	for _, chain := range target.Chains {
		if _, found := chainByID(canonical.Chains, chain.PortableID); !found {
			plan.Items = append(plan.Items, Item{Kind: "chain", PortableID: chain.PortableID, Name: chain.PortableID, Action: ActionTargetOnly, Message: "exists only on target"})
		}
	}
	sort.Slice(plan.Items, func(i, j int) bool { return planKey(plan.Items[i]) < planKey(plan.Items[j]) })
	return plan
}

func planKey(item Item) string {
	return item.Kind + "\x00" + item.PortableID + "\x00" + string(item.Action)
}
func mapByID[T any](values []T, id func(T) string) map[string]T {
	out := make(map[string]T, len(values))
	for _, value := range values {
		out[id(value)] = value
	}
	return out
}
func groupByID(values []Group, id string) (Group, bool) {
	for _, value := range values {
		if value.PortableID == id {
			return value, true
		}
	}
	return Group{}, false
}
func taskByID(values []Task, id string) (Task, bool) {
	for _, value := range values {
		if value.PortableID == id {
			return value, true
		}
	}
	return Task{}, false
}
func chainByID(values []Chain, id string) (Chain, bool) {
	for _, value := range values {
		if value.PortableID == id {
			return value, true
		}
	}
	return Chain{}, false
}
func compareGroup(source Group, target Group) Item {
	if target.PortableID == "" {
		return Item{Kind: "group", PortableID: source.PortableID, Name: source.Name, Action: ActionCreate}
	}
	if source == target {
		return Item{Kind: "group", PortableID: source.PortableID, Name: source.Name, Action: ActionUnchanged}
	}
	return Item{Kind: "group", PortableID: source.PortableID, Name: source.Name, Action: ActionUpdate}
}
func compareTask(source Task, target Task) Item {
	if target.PortableID == "" {
		return Item{Kind: "task", PortableID: source.PortableID, Name: source.Name, Action: ActionCreate}
	}
	if source == target {
		return Item{Kind: "task", PortableID: source.PortableID, Name: source.Name, Action: ActionUnchanged}
	}
	return Item{Kind: "task", PortableID: source.PortableID, Name: source.Name, Action: ActionUpdate}
}
func compareChain(source Chain, target Chain) Item {
	if target.PortableID == "" {
		return Item{Kind: "chain", PortableID: source.PortableID, Name: source.PortableID, Action: ActionCreate}
	}
	if source == target {
		return Item{Kind: "chain", PortableID: source.PortableID, Name: source.PortableID, Action: ActionUnchanged}
	}
	return Item{Kind: "chain", PortableID: source.PortableID, Name: source.PortableID, Action: ActionUpdate}
}
