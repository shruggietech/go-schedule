package gui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	tasklogic "github.com/shruggietech/go-schedule/internal/task"
)

var taskColumns = []structuredColumn{
	{Header: "Task", Minimum: 150, Weight: 3},
	{Header: "Enabled", Minimum: 80, Alignment: fyne.TextAlignCenter},
	{Header: "Effective", Minimum: 150, Weight: 2},
	{Header: "Lifecycle", Minimum: 90, Alignment: fyne.TextAlignCenter},
	{Header: "Time zone", Minimum: 110, Weight: 1},
	{Header: "Group", Minimum: 110, Weight: 2},
}

func taskRowModel(task domain.Task, groups []domain.Group) structuredRowModel {
	return taskRowModelWithChains(task, groups, nil)
}

func taskRowModelWithChains(task domain.Task, groups []domain.Group, chains []domain.CompletionChain) structuredRowModel {
	return taskRowModelWithSources(task, groups, chains, nil)
}

func taskRowModelWithSources(task domain.Task, groups []domain.Group, chains []domain.CompletionChain, triggers []server.TriggerResponse, watcherSources ...[]server.FilesystemWatcherResponse) structuredRowModel {
	name := tasklogic.DisplayName(task)
	enabled := "Disabled"
	enabledImportance := widget.LowImportance
	if task.Enabled {
		enabled = "Enabled"
		enabledImportance = widget.SuccessImportance
	}
	effective := taskEffectiveStateWithSources(task, groups, chains, triggers, watcherSources...)
	lifecycle := normalizedWords(string(task.State), "Unknown", false)
	lifecycleImportance := widget.MediumImportance
	switch task.State {
	case domain.TaskActive:
		lifecycleImportance = widget.HighImportance
	case domain.TaskDisabled:
		lifecycleImportance = widget.LowImportance
	}
	timezone := task.Timezone
	if timezone == "" {
		timezone = "Unknown"
	}
	group := groupLabelForID(task.GroupID, groups)
	groupImportance := widget.MediumImportance
	if group == groupNoneLabel {
		group = "Not assigned"
		groupImportance = widget.LowImportance
	}
	cells := []structuredCell{
		{Text: name},
		{Text: enabled, Importance: enabledImportance},
		effective,
		{Text: lifecycle, Importance: lifecycleImportance},
		{Text: timezone},
		{Text: group, Importance: groupImportance},
	}
	return structuredRowModel{
		Identity: task.ID,
		Cells:    cells,
		Summary:  structuredRowSummary(taskColumns, cells),
	}
}

func taskEffectiveStateWithSources(task domain.Task, groups []domain.Group, chains []domain.CompletionChain, triggers []server.TriggerResponse, watcherSources ...[]server.FilesystemWatcherResponse) structuredCell {
	if task.State != domain.TaskActive {
		state := normalizedWords(string(task.State), "Unknown", false)
		return structuredCell{Text: "Lifecycle: " + state, Importance: widget.LowImportance}
	}
	hasCompletion := false
	for _, chain := range chains {
		if chain.TargetTaskID == task.ID {
			hasCompletion = true
			break
		}
	}
	hasTrigger := false
	for _, trigger := range triggers {
		if trigger.Enabled && trigger.TargetTaskID == task.ID {
			hasTrigger = true
			break
		}
	}
	hasWatcher := false
	if len(watcherSources) > 0 {
		for _, watcher := range watcherSources[0] {
			if watcher.Enabled && watcher.TargetTaskID == task.ID {
				hasWatcher = true
				break
			}
		}
	}
	readiness := tasklogic.EvaluateReadiness(task, hasCompletion, hasTrigger, hasWatcher)
	if !readiness.CommandReady {
		return structuredCell{Text: "Not runnable", Importance: widget.WarningImportance}
	}
	if !readiness.ActivationReady {
		return structuredCell{Text: "Manual only", Importance: widget.MediumImportance}
	}
	if !task.Enabled {
		return structuredCell{Text: "Task disabled", Importance: widget.LowImportance}
	}
	byID := tasklogic.ByID(groups)
	if blocker, ok := tasklogic.NearestDisabledGroup(task.GroupID, byID); ok {
		return structuredCell{
			Text:       "Blocked by " + groupLabelForID(blocker.ID, groups),
			Importance: widget.WarningImportance,
		}
	}
	if !tasklogic.ChainEnabled(task.GroupID, byID) {
		return structuredCell{Text: "Group chain invalid", Importance: widget.WarningImportance}
	}
	return structuredCell{Text: "Runnable", Importance: widget.SuccessImportance}
}

// taskDetailFor fetches a task's full detail (task + schedule) so the editor can
// show what the task is actually set to. The cached task list carries no
// schedule, so this is the only way to populate the timing fields.
//
// A failed lookup is degraded, never fatal: the caller falls back to the task it
// already holds, with no schedule attached, so an unrelated edit (renaming,
// fixing a command) is not blocked by a transient read failure. The editor then
// leaves the timing fields blank, which on save keeps the stored schedule, and
// says so (FR-009).
func (a *App) taskDetailFor(task domain.Task) *server.TaskResponse {
	ctx, cancel := a.bgCtx()
	defer cancel()
	detail, err := a.backend.GetTask(ctx, task.ID)
	if err != nil {
		return &server.TaskResponse{Task: task}
	}
	return &detail
}

func (a *App) buildTasksTab() fyne.CanvasObject {
	var tasks []domain.Task
	table := newStructuredList(
		taskColumns,
		"Select a task to see its complete values.",
		nil,
		func(id string) { a.editTaskByID(id) },
	)
	a.taskTable = table
	a.taskList = table.list

	refresh := func() {
		snapshot := a.model.Snapshot()
		tasks = snapshot.Tasks
		rows := make([]structuredRowModel, len(tasks))
		for index, task := range tasks {
			rows[index] = taskRowModelWithSources(task, snapshot.Groups, snapshot.Chains, snapshot.Triggers, snapshot.Watchers)
		}
		table.setRows(rows)
	}
	a.registerRefresher(refresh)

	cur := func() (domain.Task, bool) {
		return currentTaskByID(tasks, table.selectedIdentity)
	}
	withSel := func(fn func(task domain.Task)) {
		if task, ok := cur(); ok {
			fn(task)
		} else {
			dialog.ShowInformation("No selection", "Select a task first.", a.win)
		}
	}

	newBtn := newToolbarButton("New", theme.ContentAddIcon(), func() { a.showTaskEditor(nil) })
	editBtn := newToolbarButton("Edit", theme.DocumentCreateIcon(), func() {
		withSel(func(task domain.Task) { a.showTaskEditor(a.taskDetailFor(task)) })
	})
	runBtn := newToolbarButton("Run now", theme.MediaPlayIcon(), func() {
		withSel(func(task domain.Task) {
			a.run(func(ctx context.Context) error { return a.backend.RunNow(ctx, task.ID) })
		})
	})
	toggleBtn := newToolbarButtonPlain("Enable/Disable", func() {
		withSel(func(task domain.Task) {
			a.run(func(ctx context.Context) error { return a.backend.SetTaskEnabled(ctx, task.ID, !task.Enabled) })
		})
	})
	delBtn := newToolbarButton("Delete", theme.DeleteIcon(), func() {
		withSel(func(task domain.Task) {
			dialog.ShowConfirm("Delete task", "Delete "+tasklogic.DisplayName(task)+"?", func(ok bool) {
				if ok {
					a.run(func(ctx context.Context) error { return a.backend.DeleteTask(ctx, task.ID) })
				}
			}, a.win)
		})
	})
	// No manual Refresh: the view updates live from the event stream (FR-023).
	toolbar := container.NewHBox(newBtn, editBtn, runBtn, toggleBtn, delBtn)
	return container.NewBorder(toolbar, nil, nil, nil, table.root)
}

func currentTaskByID(tasks []domain.Task, id string) (domain.Task, bool) {
	if id == "" {
		return domain.Task{}, false
	}
	for _, task := range tasks {
		if task.ID == id {
			return task, true
		}
	}
	return domain.Task{}, false
}

func (a *App) editTaskByID(id string) bool {
	task, ok := currentTaskByID(a.model.Snapshot().Tasks, id)
	if !ok {
		return false
	}
	a.showTaskEditor(a.taskDetailFor(task))
	return true
}
