package gui

import (
	"context"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/task"
)

// The group hierarchy shows groups and their member tasks in one tree, so
// membership is visible and editable without leaving the view (FR-018/FR-020).
// Node IDs are prefixed because a group and a task may share an identifier, and
// the tree needs to tell them apart.
const (
	groupNodePrefix = "g:"
	taskNodePrefix  = "t:"
	// ungroupedNodeID is a synthetic root collecting tasks that belong to no
	// group. It is always present, so the destination for removing a task from
	// its group is visible even before anything is in it (FR-019).
	ungroupedNodeID = "u:"
	ungroupedLabel  = "Ungrouped"
	// taskRowMarker prefixes task rows so they read differently from groups.
	taskRowMarker = "• "
)

func groupNodeID(id string) string { return groupNodePrefix + id }
func taskNodeID(id string) string  { return taskNodePrefix + id }

// groupTreeModel is the tree's data, derived from the current groups and tasks.
// It is a plain value with no widget dependency so the hierarchy's shape can be
// unit-tested directly.
type groupTreeModel struct {
	children map[string][]string // node ID -> ordered child node IDs ("" = roots)
	labels   map[string]string
	groups   map[string]domain.Group
	tasks    map[string]domain.Task
}

// newGroupTreeModel arranges groups into their hierarchy and files every task
// under its group — or under Ungrouped when it has none, or names one that no
// longer resolves (FR-019a). Every task appears exactly once.
func newGroupTreeModel(groups []domain.Group, tasks []domain.Task) *groupTreeModel {
	m := &groupTreeModel{
		children: map[string][]string{},
		labels:   map[string]string{},
		groups:   task.ByID(groups),
		tasks:    map[string]domain.Task{},
	}

	// Groups first, so member tasks are appended after any child groups.
	for _, g := range groups {
		parent := ""
		if g.ParentID != "" {
			if _, ok := m.groups[g.ParentID]; ok {
				parent = groupNodeID(g.ParentID)
			}
		}
		m.children[parent] = append(m.children[parent], groupNodeID(g.ID))
		label := g.Name
		if !g.Enabled {
			label += "  (disabled)"
		}
		m.labels[groupNodeID(g.ID)] = label
	}

	m.children[""] = append(m.children[""], ungroupedNodeID)
	m.labels[ungroupedNodeID] = ungroupedLabel

	sorted := append([]domain.Task(nil), tasks...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })
	for _, t := range sorted {
		m.tasks[t.ID] = t
		parent := ungroupedNodeID
		if _, ok := m.groups[t.GroupID]; ok && t.GroupID != "" {
			parent = groupNodeID(t.GroupID)
		}
		m.children[parent] = append(m.children[parent], taskNodeID(t.ID))
		m.labels[taskNodeID(t.ID)] = taskRowMarker + t.Name + "   [" + string(t.State) + "]   " +
			boolStr(t.Enabled, "enabled", "disabled")
	}
	return m
}

func (m *groupTreeModel) label(nodeID string) string { return m.labels[nodeID] }

func (m *groupTreeModel) isTask(nodeID string) bool {
	return strings.HasPrefix(nodeID, taskNodePrefix)
}

func (m *groupTreeModel) isGroup(nodeID string) bool {
	return strings.HasPrefix(nodeID, groupNodePrefix)
}

func (m *groupTreeModel) taskID(nodeID string) string {
	return strings.TrimPrefix(nodeID, taskNodePrefix)
}

func (m *groupTreeModel) groupID(nodeID string) string {
	return strings.TrimPrefix(nodeID, groupNodePrefix)
}

// group returns the selected group, if the selection is a group node.
func (m *groupTreeModel) group(nodeID string) (domain.Group, bool) {
	if !m.isGroup(nodeID) {
		return domain.Group{}, false
	}
	g, ok := m.groups[m.groupID(nodeID)]
	return g, ok
}

// moveTaskToGroup reassigns a task from a choice label, carrying the three-way
// intent: a real group assigns, "(none)" removes it from its group.
func (a *App) moveTaskToGroup(taskID, label string) {
	groups := a.model.Snapshot().Groups
	id := groupIDForLabel(label, groups)
	a.run(func(ctx context.Context) error {
		_, err := a.backend.UpdateTask(ctx, taskID, server.TaskUpdateRequest{GroupID: &id})
		return err
	})
}

// buildGroupsTab renders the group hierarchy with its member tasks, plus
// enable/disable, add/delete, and moving a task between groups
// (FR-018 – FR-021).
func (a *App) buildGroupsTab() fyne.CanvasObject {
	model := newGroupTreeModel(nil, nil)

	tree := widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID { return model.children[string(id)] },
		func(id widget.TreeNodeID) bool { return len(model.children[string(id)]) > 0 },
		func(bool) fyne.CanvasObject { return widget.NewLabel("template") },
		func(id widget.TreeNodeID, _ bool, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(model.label(string(id)))
		},
	)
	selected := ""
	tree.OnSelected = func(id widget.TreeNodeID) { selected = string(id) }
	tree.OnUnselected = func(widget.TreeNodeID) { selected = "" }

	refresh := func() {
		snap := a.model.Snapshot()
		model = newGroupTreeModel(snap.Groups, snap.Tasks)
		tree.Refresh()
	}
	a.registerRefresher(refresh)

	// withGroup / withTask keep each action bound to the kind of node it applies
	// to, so a group action never fires against a selected task (FR-021).
	withGroup := func(fn func(g domain.Group)) {
		if g, ok := model.group(selected); ok {
			fn(g)
			return
		}
		dialog.ShowInformation("Select a group", "Select a group first.", a.win)
	}
	withTask := func(fn func(id string)) {
		if model.isTask(selected) {
			fn(model.taskID(selected))
			return
		}
		dialog.ShowInformation("Select a task", "Select a task first.", a.win)
	}

	addBtn := newToolbarButton("New Group", theme.ContentAddIcon(), func() {
		nameEntry := widget.NewEntry()
		parent := "" // a selected group becomes the parent
		parentNote := "top-level"
		if g, ok := model.group(selected); ok {
			parent = g.ID
			parentNote = "under " + g.Name
		}
		items := []*widget.FormItem{
			widget.NewFormItem("Name", nameEntry),
			widget.NewFormItem("Parent", widget.NewLabel(parentNote)),
		}
		dialog.NewForm("New Group", "Create", "Cancel", items, func(ok bool) {
			if !ok || nameEntry.Text == "" {
				return
			}
			a.run(func(ctx context.Context) error {
				_, err := a.backend.CreateGroup(ctx, server.GroupCreateRequest{Name: nameEntry.Text, ParentID: parent})
				return err
			})
		}, a.win).Show()
	})
	toggleBtn := newToolbarButtonPlain("Enable/Disable", func() {
		withGroup(func(g domain.Group) {
			a.run(func(ctx context.Context) error { return a.backend.SetGroupEnabled(ctx, g.ID, !g.Enabled) })
		})
	})
	moveBtn := newToolbarButtonPlain("Move to group…", func() {
		withTask(func(id string) {
			groups := a.model.Snapshot().Groups
			sel := widget.NewSelect(groupChoiceLabels(groups), nil)
			sel.SetSelected(groupLabelForID(model.tasks[id].GroupID, groups))
			items := []*widget.FormItem{widget.NewFormItem("Group", sel)}
			dialog.NewForm("Move to group", "Move", "Cancel", items, func(ok bool) {
				if ok {
					a.moveTaskToGroup(id, sel.Selected)
				}
			}, a.win).Show()
		})
	})
	delBtn := newToolbarButton("Delete", theme.DeleteIcon(), func() {
		withGroup(func(g domain.Group) {
			dialog.ShowConfirm("Delete group", "Delete "+g.Name+" (children cascade; its tasks become ungrouped)?", func(yes bool) {
				if yes {
					a.run(func(ctx context.Context) error { return a.backend.DeleteGroup(ctx, g.ID) })
				}
			}, a.win)
		})
	})
	// No manual Refresh: the view updates live from the event stream (FR-023).
	toolbar := container.NewHBox(addBtn, toggleBtn, moveBtn, delBtn)
	return container.NewBorder(toolbar, nil, nil, nil, tree)
}
