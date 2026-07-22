package gui

import (
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// nestedGroups has two groups named "Nightly" at different levels, which is the
// case a bare name list cannot distinguish (FR-013).
func nestedGroups() []domain.Group {
	return []domain.Group{
		{ID: "g-backups", Name: "Backups", Enabled: true},
		{ID: "g-b-nightly", Name: "Nightly", ParentID: "g-backups", Enabled: true},
		{ID: "g-reports", Name: "Reports", Enabled: true},
		{ID: "g-r-nightly", Name: "Nightly", ParentID: "g-reports", Enabled: true},
	}
}

// TestGroupChoices_PathsAndRoundTrip covers FR-012/FR-013.
func TestGroupChoices_PathsAndRoundTrip(t *testing.T) {
	groups := nestedGroups()
	labels := groupChoiceLabels(groups)

	if len(labels) == 0 || labels[0] != groupNoneLabel {
		t.Fatalf("first choice = %q, want %q so removing a task from its group is always offered",
			labels, groupNoneLabel)
	}
	if len(labels) != len(groups)+1 {
		t.Fatalf("got %d choices for %d groups, want one per group plus %q: %v",
			len(labels), len(groups), groupNoneLabel, labels)
	}

	// The two same-named groups must be distinguishable.
	seen := map[string]bool{}
	for _, l := range labels {
		if seen[l] {
			t.Errorf("duplicate choice label %q: same-named groups at different levels are indistinguishable", l)
		}
		seen[l] = true
	}

	// Every group round-trips label -> id -> label.
	for _, g := range groups {
		label := groupLabelForID(g.ID, groups)
		if label == groupNoneLabel {
			t.Errorf("group %s (%s) rendered as %q", g.ID, g.Name, label)
			continue
		}
		if got := groupIDForLabel(label, groups); got != g.ID {
			t.Errorf("round-trip for %s: label %q -> id %q", g.ID, label, got)
		}
	}

	// "(none)" maps to the empty id, and an unresolvable id reads as "(none)".
	if got := groupIDForLabel(groupNoneLabel, groups); got != "" {
		t.Errorf("%q -> %q, want empty id", groupNoneLabel, got)
	}
	if got := groupLabelForID("ghost", groups); got != groupNoneLabel {
		t.Errorf("unresolvable group id -> %q, want %q (FR-019a)", got, groupNoneLabel)
	}
	if got := groupLabelForID("", groups); got != groupNoneLabel {
		t.Errorf("empty group id -> %q, want %q", got, groupNoneLabel)
	}
}

// TestEditor_GroupPrefillAndSubmit covers FR-012 and the FR-019a fallback.
func TestEditor_GroupPrefillAndSubmit(t *testing.T) {
	detail := recurringDetail("weekdays at 09:00")
	detail.Task.GroupID = "g-b-nightly"

	fb := &fakeBackend{groups: nestedGroups()}
	ui := NewUI(testApp, fb)
	if err := ui.model.Refresh(t.Context()); err != nil {
		t.Fatal(err)
	}
	e := newTaskEditor(ui, detail)
	e.previewSync = true
	e.build()

	want := groupLabelForID("g-b-nightly", fb.groups)
	if e.group.Selected != want {
		t.Errorf("group prefill = %q, want %q", e.group.Selected, want)
	}
	if e.isDirty() {
		t.Error("prefilled group must be part of the dirty baseline")
	}

	// Reassigning reaches the submitted request.
	e.group.SetSelected(groupLabelForID("g-reports", fb.groups))
	if got := e.buildForm().groupID; got == nil || *got != "g-reports" {
		t.Errorf("buildForm groupID = %v, want pointer to g-reports", got)
	}
}

// TestEditor_UnresolvableGroupPrefillsAsNone covers FR-019a.
func TestEditor_UnresolvableGroupPrefillsAsNone(t *testing.T) {
	detail := recurringDetail("weekdays at 09:00")
	detail.Task.GroupID = "vanished"

	fb := &fakeBackend{groups: nestedGroups()}
	ui := NewUI(testApp, fb)
	if err := ui.model.Refresh(t.Context()); err != nil {
		t.Fatal(err)
	}
	e := newTaskEditor(ui, detail)
	e.previewSync = true
	e.build()

	if e.group.Selected != groupNoneLabel {
		t.Errorf("unresolvable group prefilled as %q, want %q with no error", e.group.Selected, groupNoneLabel)
	}
}

// TestEditor_NoneMeansClearOnlyWhenItWasGrouped covers FR-014 / User Story 3.
// Choosing "(none)" for a task that was in a group is a deliberate removal; for
// a task that was never grouped it is a no-op and must not emit a write.
func TestEditor_NoneMeansClearOnlyWhenItWasGrouped(t *testing.T) {
	fb := &fakeBackend{groups: nestedGroups()}
	ui := NewUI(testApp, fb)
	if err := ui.model.Refresh(t.Context()); err != nil {
		t.Fatal(err)
	}

	grouped := recurringDetail("weekdays at 09:00")
	grouped.Task.GroupID = "g-reports"
	e := newTaskEditor(ui, grouped)
	e.previewSync = true
	e.build()
	e.group.SetSelected(groupNoneLabel)
	got := e.buildForm().groupID
	if got == nil || *got != "" {
		t.Errorf("clearing a grouped task: groupID = %v, want pointer to empty string", got)
	}

	ungrouped := recurringDetail("weekdays at 09:00")
	e2 := newTaskEditor(ui, ungrouped)
	e2.previewSync = true
	e2.build()
	if got := e2.buildForm().groupID; got != nil {
		t.Errorf("untouched ungrouped task: groupID = %q, want nil (no pointless write)", *got)
	}
}

// TestTasksTab_RowShowsGroup covers FR-017.
func TestTasksTab_RowShowsGroup(t *testing.T) {
	groups := nestedGroups()
	tasks := []domain.Task{
		{ID: "t1", Name: "db-dump", Command: "/bin/true", Timezone: "UTC", Enabled: true, State: domain.TaskActive, GroupID: "g-b-nightly"},
		{ID: "t2", Name: "loose", Command: "/bin/true", Timezone: "UTC", Enabled: true, State: domain.TaskActive},
	}

	if got := taskRowText(tasks[0], groups); !containsAll(got, "db-dump", "Nightly") {
		t.Errorf("grouped task row = %q, want it to name the task and its group", got)
	}
	if got := taskRowText(tasks[1], groups); !containsAll(got, "loose") {
		t.Errorf("ungrouped task row = %q, want it to name the task", got)
	}
}

func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}
