package gui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

// taskRowText renders one task-list row. It names the task's group so
// membership is visible without opening the editor; an ungrouped task simply
// omits that column rather than showing a placeholder.
func taskRowText(t domain.Task, groups []domain.Group) string {
	row := fmt.Sprintf("%s   [%s]   %s   %s",
		t.Name, t.State, boolStr(t.Enabled, "enabled", "disabled"), t.Timezone)
	if label := groupLabelForID(t.GroupID, groups); label != groupNoneLabel {
		row += "   " + label
	}
	return row
}

// taskDetailFor fetches a task's full detail (task + schedule) so the editor can
// show what the task is actually set to. The cached task list carries no
// schedule, so this is the only way to populate the timing fields.
//
// A failed lookup is degraded, never fatal: the caller falls back to the task it
// already holds, with no schedule attached, so an unrelated edit (renaming,
// fixing a command) is not blocked by a transient read failure. The editor then
// leaves the timing fields blank — which on save keeps the stored schedule — and
// says so (FR-009).
func (a *App) taskDetailFor(t domain.Task) *server.TaskResponse {
	ctx, cancel := a.bgCtx()
	defer cancel()
	detail, err := a.backend.GetTask(ctx, t.ID)
	if err != nil {
		return &server.TaskResponse{Task: t}
	}
	return &detail
}

func (a *App) buildTasksTab() fyne.CanvasObject {
	var tasks []domain.Task
	selected := -1

	list := widget.NewList(
		func() int { return len(tasks) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(taskRowText(tasks[i], a.model.Snapshot().Groups))
		},
	)
	list.OnSelected = func(id widget.ListItemID) { selected = id }
	list.OnUnselected = func(widget.ListItemID) { selected = -1 }

	refresh := func() {
		tasks = a.model.Snapshot().Tasks
		list.Refresh()
	}
	a.registerRefresher(refresh)

	cur := func() (domain.Task, bool) {
		if selected < 0 || selected >= len(tasks) {
			return domain.Task{}, false
		}
		return tasks[selected], true
	}
	withSel := func(fn func(t domain.Task)) {
		if t, ok := cur(); ok {
			fn(t)
		} else {
			dialog.ShowInformation("No selection", "Select a task first.", a.win)
		}
	}

	newBtn := newToolbarButton("New", theme.ContentAddIcon(), func() { a.showTaskEditor(nil) })
	editBtn := newToolbarButton("Edit", theme.DocumentCreateIcon(), func() {
		withSel(func(t domain.Task) { a.showTaskEditor(a.taskDetailFor(t)) })
	})
	runBtn := newToolbarButton("Run now", theme.MediaPlayIcon(), func() {
		withSel(func(t domain.Task) { a.run(func(ctx context.Context) error { return a.backend.RunNow(ctx, t.ID) }) })
	})
	toggleBtn := newToolbarButtonPlain("Enable/Disable", func() {
		withSel(func(t domain.Task) {
			a.run(func(ctx context.Context) error { return a.backend.SetTaskEnabled(ctx, t.ID, !t.Enabled) })
		})
	})
	delBtn := newToolbarButton("Delete", theme.DeleteIcon(), func() {
		withSel(func(t domain.Task) {
			dialog.ShowConfirm("Delete task", "Delete "+t.Name+"?", func(ok bool) {
				if ok {
					a.run(func(ctx context.Context) error { return a.backend.DeleteTask(ctx, t.ID) })
				}
			}, a.win)
		})
	})
	// No manual Refresh: the view updates live from the event stream (FR-023).
	toolbar := container.NewHBox(newBtn, editBtn, runBtn, toggleBtn, delBtn)
	return container.NewBorder(toolbar, nil, nil, nil, list)
}
