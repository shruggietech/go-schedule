package gui

import (
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/events"
)

func treeTasks() []domain.Task {
	return []domain.Task{
		{ID: "t1", Name: "db-dump", Command: "/bin/true", Timezone: "UTC", Enabled: true, State: domain.TaskActive, GroupID: "g-b-nightly"},
		{ID: "t2", Name: "weekly-report", Command: "/bin/true", Timezone: "UTC", Enabled: true, State: domain.TaskActive, GroupID: "g-reports"},
		{ID: "t3", Name: "loose", Command: "/bin/true", Timezone: "UTC", Enabled: true, State: domain.TaskActive},
		{ID: "t4", Name: "orphan", Command: "/bin/true", Timezone: "UTC", Enabled: true, State: domain.TaskActive, GroupID: "vanished"},
	}
}

// collectNodes walks the tree from its roots, returning every reachable node ID.
func collectNodes(tree *groupTreeModel) []string {
	var out []string
	var walk func(id string)
	walk = func(id string) {
		for _, child := range tree.children[id] {
			out = append(out, child)
			walk(child)
		}
	}
	walk("")
	return out
}

// TestGroupTree_ShowsEveryTaskExactlyOnce covers FR-018/FR-019a and SC-005.
func TestGroupTree_ShowsEveryTaskExactlyOnce(t *testing.T) {
	tree := newGroupTreeModel(nestedGroups(), treeTasks())
	nodes := collectNodes(tree)

	counts := map[string]int{}
	for _, n := range nodes {
		counts[n]++
	}
	for _, task := range treeTasks() {
		id := taskNodeID(task.ID)
		if counts[id] != 1 {
			t.Errorf("task %s appears %d times in the tree, want exactly once", task.Name, counts[id])
		}
	}

	// Membership lands under the right parent.
	if !containsStr(tree.children[groupNodeID("g-b-nightly")], taskNodeID("t1")) {
		t.Error("db-dump is not listed under its group Backups / Nightly")
	}
	if !containsStr(tree.children[groupNodeID("g-reports")], taskNodeID("t2")) {
		t.Error("weekly-report is not listed under Reports")
	}
	// Ungrouped and dangling-reference tasks both land in the ungrouped area.
	ungrouped := tree.children[ungroupedNodeID]
	if !containsStr(ungrouped, taskNodeID("t3")) {
		t.Error("an ungrouped task is not shown in the ungrouped area")
	}
	if !containsStr(ungrouped, taskNodeID("t4")) {
		t.Error("a task whose group cannot be resolved must be shown as ungrouped (FR-019a)")
	}
}

// TestGroupTree_UngroupedNodeAlwaysPresent covers FR-019: the destination for
// removing a task from a group must be visible before there is anything in it.
func TestGroupTree_UngroupedNodeAlwaysPresent(t *testing.T) {
	tree := newGroupTreeModel(nestedGroups(), nil)
	if !containsStr(tree.children[""], ungroupedNodeID) {
		t.Fatal("ungrouped node missing when no task is ungrouped")
	}
	if label := tree.label(ungroupedNodeID); !strings.Contains(strings.ToLower(label), "ungrouped") {
		t.Errorf("ungrouped node label = %q, want it clearly labeled", label)
	}
}

// TestGroupTree_NodeIdentitiesDoNotCollide covers the edge case where a group
// and a task share an identifier.
func TestGroupTree_NodeIdentitiesDoNotCollide(t *testing.T) {
	groups := []domain.Group{{ID: "same", Name: "Shared", Enabled: true}}
	tasks := []domain.Task{{ID: "same", Name: "AlsoShared", GroupID: "same", State: domain.TaskActive}}
	tree := newGroupTreeModel(groups, tasks)

	if groupNodeID("same") == taskNodeID("same") {
		t.Fatal("group and task node IDs collide")
	}
	if !tree.isTask(taskNodeID("same")) {
		t.Error("task node not recognized as a task")
	}
	if tree.isTask(groupNodeID("same")) {
		t.Error("group node misidentified as a task")
	}
	if got := tree.taskID(taskNodeID("same")); got != "same" {
		t.Errorf("taskID = %q, want the bare task id", got)
	}
}

// TestGroupTree_TaskRowsAreDistinguishable covers FR-018.
func TestGroupTree_TaskRowsAreDistinguishable(t *testing.T) {
	tree := newGroupTreeModel(nestedGroups(), treeTasks())
	groupLabel := tree.label(groupNodeID("g-reports"))
	taskLabel := tree.label(taskNodeID("t2"))

	if !strings.Contains(taskLabel, "weekly-report") {
		t.Errorf("task label = %q, want it to name the task", taskLabel)
	}
	if groupLabel == taskLabel {
		t.Error("group and task rows render identically")
	}
	if !strings.HasPrefix(taskLabel, taskRowMarker) {
		t.Errorf("task label = %q, want the task marker prefix so rows read differently from groups", taskLabel)
	}
}

// TestGroupsTab_MoveToGroupIssuesUpdate covers FR-020, including the clear case.
func TestGroupsTab_MoveToGroupIssuesUpdate(t *testing.T) {
	fb := &fakeBackend{groups: nestedGroups(), tasks: treeTasks()}
	ui := NewUI(testApp, fb)
	if err := ui.model.Refresh(t.Context()); err != nil {
		t.Fatal(err)
	}

	// Move an assigned task to another group.
	ui.moveTaskToGroup("t1", groupLabelForID("g-reports", fb.groups))
	waitFor(t, func() bool { n, _, _ := fb.lastUpdateCall(); return n == 1 })
	_, id, req := fb.lastUpdateCall()
	if id != "t1" {
		t.Errorf("updated task %q, want t1", id)
	}
	if req.GroupID == nil || *req.GroupID != "g-reports" {
		t.Errorf("GroupID = %v, want pointer to g-reports", req.GroupID)
	}

	// Move it out of all groups.
	ui.moveTaskToGroup("t1", groupNoneLabel)
	waitFor(t, func() bool { n, _, _ := fb.lastUpdateCall(); return n == 2 })
	if _, _, req = fb.lastUpdateCall(); req.GroupID == nil || *req.GroupID != "" {
		t.Errorf("GroupID = %v, want pointer to empty string (ungroup)", req.GroupID)
	}
}

// TestGroupTree_ReflectsTaskEvents covers FR-022/SC-006: membership changes show
// up without a manual refresh. The Groups tab now reads Tasks, a dependency it
// did not previously have, so the live-update assumption is pinned here.
func TestGroupTree_ReflectsTaskEvents(t *testing.T) {
	fb := &fakeBackend{groups: nestedGroups(), tasks: treeTasks()}
	ui := NewUI(testApp, fb)
	if err := ui.model.Refresh(t.Context()); err != nil {
		t.Fatal(err)
	}

	moved := treeTasks()[0]
	moved.GroupID = "g-reports"
	ui.model.ApplyEvent(events.Event{
		Kind: events.KindTask,
		Task: &events.TaskEvent{Verb: events.VerbUpdated, ID: moved.ID, Task: &moved},
	})

	snap := ui.model.Snapshot()
	tree := newGroupTreeModel(snap.Groups, snap.Tasks)
	if containsStr(tree.children[groupNodeID("g-b-nightly")], taskNodeID("t1")) {
		t.Error("task still shown under its old group after a task event")
	}
	if !containsStr(tree.children[groupNodeID("g-reports")], taskNodeID("t1")) {
		t.Error("task not shown under its new group after a task event")
	}

	// Deletion removes it from the tree entirely.
	ui.model.ApplyEvent(events.Event{
		Kind: events.KindTask,
		Task: &events.TaskEvent{Verb: events.VerbDeleted, ID: "t1"},
	})
	snap = ui.model.Snapshot()
	tree = newGroupTreeModel(snap.Groups, snap.Tasks)
	if containsStr(collectNodes(tree), taskNodeID("t1")) {
		t.Error("deleted task still present in the tree")
	}
}

func containsStr(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
