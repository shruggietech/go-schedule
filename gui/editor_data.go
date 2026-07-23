package gui

import (
	"strings"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/task"
)

// This file holds presentation-only data for the task editor: the human-readable
// labels shown for overlap/catch-up policies (the stored wire values are
// unchanged), the curated timezone suggestions, and the command-line preview
// builder. Keeping them here keeps editor.go focused on widget wiring.

// overlapChoice / catchupChoice pair a friendly label with its stored value.
type policyChoice[T ~string] struct {
	label string
	value T
}

var overlapChoices = []policyChoice[domain.OverlapPolicy]{
	{"Queue one run", domain.OverlapQueueOne},
	{"Skip this run", domain.OverlapSkip},
	{"Allow concurrent runs", domain.OverlapAllowConcurrent},
}

var catchupChoices = []policyChoice[domain.CatchupPolicy]{
	{"Run once to catch up", domain.CatchupOne},
	{"Skip missed runs", domain.CatchupNone},
}

// missingDateChoices covers what a schedule does in a period with no matching
// date — February for a rule on the 29th, a 30-day month for a rule on the 31st.
// The labels name the outcome rather than the stored value, because "last_valid"
// tells an operator nothing about what their task will do.
var missingDateChoices = []policyChoice[domain.MissingDatePolicy]{
	{"Skip that period", domain.MissingDateSkip},
	{"Use the last valid date", domain.MissingDateLastValid},
	{"Roll into the next period", domain.MissingDateNextValid},
}

// overlapLabels / catchupLabels are the ordered display strings for the selects.
func overlapLabels() []string     { return labelsOf(overlapChoices) }
func catchupLabels() []string     { return labelsOf(catchupChoices) }
func missingDateLabels() []string { return labelsOf(missingDateChoices) }

func labelsOf[T ~string](cs []policyChoice[T]) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.label
	}
	return out
}

// overlapValue maps a display label back to its stored value, falling back to the
// default (first choice) for any unknown label so the UI never crashes on legacy
// or empty input.
func overlapValue(label string) domain.OverlapPolicy {
	for _, c := range overlapChoices {
		if c.label == label {
			return c.value
		}
	}
	return overlapChoices[0].value
}

func catchupValue(label string) domain.CatchupPolicy {
	for _, c := range catchupChoices {
		if c.label == label {
			return c.value
		}
	}
	return catchupChoices[0].value
}

func missingDateValue(label string) domain.MissingDatePolicy {
	for _, c := range missingDateChoices {
		if c.label == label {
			return c.value
		}
	}
	return missingDateChoices[0].value
}

// overlapLabel maps a stored value back to its display label (default label for
// unknown values).
func overlapLabel(v domain.OverlapPolicy) string {
	for _, c := range overlapChoices {
		if c.value == v {
			return c.label
		}
	}
	return overlapChoices[0].label
}

func catchupLabel(v domain.CatchupPolicy) string {
	for _, c := range catchupChoices {
		if c.value == v {
			return c.label
		}
	}
	return catchupChoices[0].label
}

func missingDateLabel(v domain.MissingDatePolicy) string {
	for _, c := range missingDateChoices {
		if c.value == v {
			return c.label
		}
	}
	return missingDateChoices[0].label
}

// --- group choices -------------------------------------------------------

// groupNoneLabel is the choice meaning "this task belongs to no group". It is
// always offered, so removing a task from its group is reachable wherever a
// group can be chosen.
const groupNoneLabel = "(none)"

// groupChoiceLabels builds the ordered choice list for every group picker in the
// app: groupNoneLabel first, then one entry per group rendered as its full path
// through the hierarchy ("Backups / Nightly"). Paths rather than bare names,
// because two groups at different levels may share a name and a bare list would
// make them indistinguishable.
//
// This is the single source for group choices — the task editor and the Groups
// tab's move action both use it, so the two paths to the same operation cannot
// drift apart.
func groupChoiceLabels(groups []domain.Group) []string {
	out := []string{groupNoneLabel}
	var walk func(nodes []*task.TreeNode, prefix string)
	walk = func(nodes []*task.TreeNode, prefix string) {
		for _, n := range nodes {
			path := n.Group.Name
			if prefix != "" {
				path = prefix + groupPathSep + n.Group.Name
			}
			out = append(out, path)
			walk(n.Children, path)
		}
	}
	walk(task.BuildForest(groups), "")
	return out
}

const groupPathSep = " / "

// groupLabelForID renders a group ID as its choice label. An empty or
// unresolvable ID reads as groupNoneLabel: a dangling reference is presented as
// "ungrouped" rather than raised as an error (FR-019a).
func groupLabelForID(id string, groups []domain.Group) string {
	if id == "" {
		return groupNoneLabel
	}
	byID := task.ByID(groups)
	g, ok := byID[id]
	if !ok {
		return groupNoneLabel
	}
	path := g.Name
	// Walk up to the root, guarding against a malformed cycle by bounding the
	// climb to the number of groups.
	for i := 0; i < len(groups) && g.ParentID != ""; i++ {
		parent, ok := byID[g.ParentID]
		if !ok {
			break
		}
		path = parent.Name + groupPathSep + path
		g = parent
	}
	return path
}

// groupIDForLabel maps a choice label back to its group ID; groupNoneLabel and
// any unknown label yield the empty ID.
func groupIDForLabel(label string, groups []domain.Group) string {
	if label == "" || label == groupNoneLabel {
		return ""
	}
	for _, g := range groups {
		if groupLabelForID(g.ID, groups) == label {
			return g.ID
		}
	}
	return ""
}

// commonZones seeds the timezone SelectEntry. It is a curated, ordered subset of
// the IANA database for quick selection; any other valid IANA name typed by the
// user is still accepted (validated via timezone.Resolve).
var commonZones = []string{
	"Local", "UTC",
	"America/New_York", "America/Chicago", "America/Denver", "America/Los_Angeles",
	"America/Sao_Paulo",
	"Europe/London", "Europe/Paris", "Europe/Berlin", "Europe/Moscow",
	"Asia/Kolkata", "Asia/Shanghai", "Asia/Tokyo",
	"Australia/Sydney", "Pacific/Auckland",
}

// commandLinePreview renders the resolved command line for display only: the
// command followed by each argument, with whitespace-bearing tokens quoted for
// readability. Execution still receives the raw argument slice — this never
// re-parses or shell-splits.
func commandLinePreview(command string, args []string) string {
	command = strings.TrimSpace(command)
	if command == "" {
		return ""
	}
	parts := make([]string, 0, len(args)+1)
	parts = append(parts, quoteForDisplay(command))
	for _, a := range args {
		parts = append(parts, quoteForDisplay(a))
	}
	return strings.Join(parts, " ")
}

func quoteForDisplay(s string) string {
	if s == "" || strings.ContainsAny(s, " \t") {
		return `"` + s + `"`
	}
	return s
}
