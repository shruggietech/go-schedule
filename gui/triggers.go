package gui

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

var triggerColumns = []structuredColumn{
	{Header: "Trigger", Minimum: 130, Weight: 2},
	{Header: "Set", Minimum: 110, Weight: 1},
	{Header: "Position", Minimum: 75, Weight: 1},
	{Header: "Target", Minimum: 150, Weight: 2},
	{Header: "Enabled", Minimum: 75, Weight: 1},
	{Header: "Readiness", Minimum: 110, Weight: 1},
	{Header: "ID", Minimum: 130, Weight: 2},
}

type triggerBackend interface {
	CreateTrigger(context.Context, server.TriggerCreateRequest) (server.TriggerSecretResponse, error)
	UpdateTrigger(context.Context, string, server.TriggerUpdateRequest) (server.TriggerResponse, error)
	DeleteTrigger(context.Context, string) error
	SetTriggerEnabled(context.Context, string, bool) (server.TriggerResponse, error)
	RotateTrigger(context.Context, string) (server.TriggerSecretResponse, error)
	RevealTrigger(context.Context, string) (server.TriggerSecretResponse, error)
	CreateTriggerSet(context.Context, server.TriggerSetCreateRequest) (server.TriggerSetSecretResponse, error)
	ListTriggerSets(context.Context) ([]server.TriggerSetResponse, error)
	RevealTriggerSet(context.Context, string) (server.TriggerSetSecretResponse, error)
	RetargetTriggerSet(context.Context, string, string) (server.TriggerSetResponse, error)
	SetTriggerSetEnabled(context.Context, string, bool) (server.TriggerSetResponse, error)
	RotateTriggerSet(context.Context, string) (server.TriggerSetSecretResponse, error)
	DeleteTriggerSet(context.Context, string) error
}

func triggerRows(items []server.TriggerResponse) []structuredRowModel {
	rows := make([]structuredRowModel, 0, len(items))
	for _, item := range items {
		enabled := "No"
		if item.Enabled {
			enabled = "Yes"
		}
		setName, position := "Standalone", "-"
		if item.SetID != "" {
			setName, position = item.SetName, strconv.Itoa(item.SetPosition)
		}
		cells := []structuredCell{{Text: item.Name}, {Text: setName}, {Text: position}, {Text: item.TargetTaskName, FullText: item.TargetTaskName + " (" + item.TargetTaskID + ")"}, {Text: enabled}, {Text: normalizedWords(item.Readiness, "Unknown", false), FullText: item.Reason}, {Text: item.ID}}
		rows = append(rows, structuredRowModel{Identity: item.ID, Cells: cells, Summary: structuredRowSummary(triggerColumns, cells)})
	}
	return rows
}

func (a *App) buildTriggersTab() fyne.CanvasObject {
	var items []server.TriggerResponse
	selectedID := ""
	current := func() (*server.TriggerResponse, bool) {
		for i := range items {
			if items[i].ID == selectedID {
				return &items[i], true
			}
		}
		return nil, false
	}
	withSelection := func(fn func(*server.TriggerResponse)) {
		if item, ok := current(); ok {
			fn(item)
			return
		}
		dialog.ShowInformation("No selection", "Select a trigger first.", a.win)
	}
	withSelectedSet := func(fn func(server.TriggerSetResponse)) {
		item, ok := current()
		if !ok {
			dialog.ShowInformation("No selection", "Select a Trigger Set member first.", a.win)
			return
		}
		if item.SetID == "" {
			dialog.ShowInformation("Standalone trigger", "The selected trigger is not a Trigger Set member.", a.win)
			return
		}
		for _, set := range a.model.Snapshot().TriggerSets {
			if set.ID == item.SetID {
				fn(set)
				return
			}
		}
		dialog.ShowInformation("Trigger Set unavailable", "Refresh the view and select the member again.", a.win)
	}
	table := newStructuredList(triggerColumns, "Select a trigger to see its complete values.", func(id string) { selectedID = id }, func(id string) {
		selectedID = id
		if item, ok := current(); ok {
			a.showTriggerEditor(item)
		}
	})
	a.triggerTable = table
	refresh := func() {
		items = a.model.Snapshot().Triggers
		table.setRows(triggerRows(items))
	}
	a.registerRefresher(refresh)
	newButton := newToolbarButton("New", theme.ContentAddIcon(), func() { a.showTriggerEditor(nil) })
	editButton := newToolbarButton("Edit", theme.DocumentCreateIcon(), func() { withSelection(a.showTriggerEditor) })
	revealKey := newToolbarButton("Reveal key", theme.VisibilityIcon(), func() { withSelection(func(item *server.TriggerResponse) { a.revealAndCopyTrigger(item.ID, false) }) })
	copyCommand := newToolbarButton("Copy command", theme.ContentCopyIcon(), func() { withSelection(func(item *server.TriggerResponse) { a.revealAndCopyTrigger(item.ID, true) }) })
	toggle := newToolbarButton("Enable or disable", theme.MediaReplayIcon(), func() {
		withSelection(func(item *server.TriggerResponse) {
			backend, ok := a.backend.(triggerBackend)
			if !ok {
				dialog.ShowError(fmt.Errorf("trigger administration is unavailable"), a.win)
				return
			}
			a.run(func(ctx context.Context) error {
				_, err := backend.SetTriggerEnabled(ctx, item.ID, !item.Enabled)
				return err
			})
		})
	})
	rotate := newToolbarButton("Rotate key", theme.ViewRefreshIcon(), func() {
		withSelection(func(item *server.TriggerResponse) {
			dialog.ShowConfirm("Rotate trigger key", "Replace the key for "+item.Name+"? The old key will stop working immediately.", func(ok bool) {
				if ok {
					a.rotateTrigger(item.ID)
				}
			}, a.win)
		})
	})
	remove := newToolbarButton("Delete", theme.DeleteIcon(), func() {
		withSelection(func(item *server.TriggerResponse) {
			dialog.ShowConfirm("Delete trigger", "Delete "+item.Name+"? Its key will stop working immediately.", func(ok bool) {
				if !ok {
					return
				}
				backend, supported := a.backend.(triggerBackend)
				if !supported {
					dialog.ShowError(fmt.Errorf("trigger administration is unavailable"), a.win)
					return
				}
				a.run(func(ctx context.Context) error { return backend.DeleteTrigger(ctx, item.ID) })
			}, a.win)
		})
	})
	toolbar := container.NewHBox(newButton, editButton, revealKey, copyCommand, toggle, rotate, remove)
	newSet := newToolbarButton("New set", theme.ContentAddIcon(), a.showTriggerSetCreator)
	copySet := newToolbarButton("Copy set", theme.ContentCopyIcon(), func() { withSelectedSet(func(set server.TriggerSetResponse) { a.copyTriggerSet(set) }) })
	retargetSet := newToolbarButton("Retarget set", theme.DocumentCreateIcon(), func() { withSelectedSet(a.showTriggerSetRetarget) })
	toggleSet := newToolbarButton("Enable or disable set", theme.MediaReplayIcon(), func() {
		withSelectedSet(func(set server.TriggerSetResponse) {
			a.confirmTriggerSetEnabled(set, set.EnabledCount != set.MemberCount)
		})
	})
	rotateSet := newToolbarButton("Rotate set", theme.ViewRefreshIcon(), func() { withSelectedSet(a.confirmTriggerSetRotate) })
	deleteSet := newToolbarButton("Delete set", theme.DeleteIcon(), func() { withSelectedSet(a.confirmTriggerSetDelete) })
	setToolbar := container.NewHBox(newSet, copySet, retargetSet, toggleSet, rotateSet, deleteSet)
	triggerContent := container.NewBorder(container.NewVBox(toolbar, setToolbar), nil, nil, nil, table.root)
	split := container.NewVSplit(triggerContent, a.buildFilesystemWatcherPanel())
	split.Offset = 0.55
	return split
}

func (a *App) showTriggerEditor(existing *server.TriggerResponse) {
	backend, ok := a.backend.(triggerBackend)
	if !ok {
		dialog.ShowError(fmt.Errorf("trigger administration is unavailable"), a.win)
		return
	}
	domainTasks := a.model.Snapshot().Tasks
	if len(domainTasks) == 0 {
		dialog.ShowInformation("Task required", "Create a task before adding a trigger.", a.win)
		return
	}
	sort.Slice(domainTasks, func(i, j int) bool { return domainTasks[i].Name < domainTasks[j].Name })
	labels := make([]string, 0, len(domainTasks))
	ids := make(map[string]string, len(domainTasks))
	byID := make(map[string]string, len(domainTasks))
	for _, task := range domainTasks {
		label := task.Name + " (" + task.ID + ")"
		labels = append(labels, label)
		ids[label], byID[task.ID] = task.ID, label
	}
	name := widget.NewEntry()
	target := widget.NewSelect(labels, nil)
	enabled := widget.NewCheck("Allow this trigger to invoke its task", nil)
	title, confirm := "New trigger", "Create"
	if existing == nil {
		target.SetSelected(labels[0])
		enabled.SetChecked(true)
	} else {
		title, confirm = "Edit trigger", "Save"
		name.SetText(existing.Name)
		target.SetSelected(byID[existing.TargetTaskID])
		enabled.SetChecked(existing.Enabled)
	}
	form := dialog.NewForm(title, confirm, "Cancel", []*widget.FormItem{widget.NewFormItem("Name", name), widget.NewFormItem("Target task", target), widget.NewFormItem("Enabled", enabled)}, func(ok bool) {
		if !ok {
			return
		}
		if name.Text == "" || ids[target.Selected] == "" {
			dialog.ShowError(fmt.Errorf("name and target task are required"), a.win)
			return
		}
		if existing == nil {
			a.runTriggerCreate(backend, server.TriggerCreateRequest{Name: name.Text, TargetTaskID: ids[target.Selected], Enabled: &enabled.Checked})
			return
		}
		targetID := ids[target.Selected]
		request := server.TriggerUpdateRequest{Name: &name.Text}
		if targetID != existing.TargetTaskID {
			request.TargetTaskID = &targetID
		}
		a.run(func(ctx context.Context) error {
			_, err := backend.UpdateTrigger(ctx, existing.ID, request)
			if err != nil {
				return err
			}
			if existing.Enabled != enabled.Checked {
				_, err = backend.SetTriggerEnabled(ctx, existing.ID, enabled.Checked)
			}
			return err
		})
	}, a.win)
	form.Resize(fyne.NewSize(560, form.MinSize().Height))
	form.Show()
}

func (a *App) showTriggerSetCreator() {
	backend, ok := a.backend.(triggerBackend)
	if !ok {
		dialog.ShowError(fmt.Errorf("trigger set administration is unavailable"), a.win)
		return
	}
	tasks, labels, ids := sortedTriggerTasks(a.model.Snapshot().Tasks)
	if len(tasks) == 0 {
		dialog.ShowInformation("Task required", "Create a task before adding a Trigger Set.", a.win)
		return
	}
	name := widget.NewEntry()
	target := widget.NewSelect(labels, nil)
	target.SetSelected(labels[0])
	count := widget.NewEntry()
	count.SetText("3")
	enabled := widget.NewCheck("Allow every member to invoke its task", nil)
	enabled.SetChecked(true)
	form := dialog.NewForm("New Trigger Set", "Create", "Cancel", []*widget.FormItem{widget.NewFormItem("Name", name), widget.NewFormItem("Target task", target), widget.NewFormItem("Member count (1-99)", count), widget.NewFormItem("Enabled", enabled)}, func(ok bool) {
		if !ok {
			return
		}
		n, err := strconv.Atoi(strings.TrimSpace(count.Text))
		if strings.TrimSpace(name.Text) == "" || ids[target.Selected] == "" || err != nil || n < 1 || n > 99 {
			dialog.ShowError(fmt.Errorf("name, target task, and a member count from 1 through 99 are required"), a.win)
			return
		}
		go func() {
			ctx, cancel := a.bgCtx()
			defer cancel()
			result, err := backend.CreateTriggerSet(ctx, server.TriggerSetCreateRequest{Name: name.Text, TargetTaskID: ids[target.Selected], Count: n, Enabled: &enabled.Checked})
			if err != nil {
				fyne.Do(func() { a.showError(err) })
				return
			}
			fyne.Do(func() { a.showTriggerSetSecrets("Trigger Set created", result) })
			a.refreshAll()
		}()
	}, a.win)
	form.Resize(fyne.NewSize(600, form.MinSize().Height))
	form.Show()
}

func sortedTriggerTasks(tasks []domain.Task) ([]domain.Task, []string, map[string]string) {
	tasks = append([]domain.Task(nil), tasks...)
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].Name < tasks[j].Name })
	labels := make([]string, 0, len(tasks))
	ids := make(map[string]string, len(tasks))
	for _, task := range tasks {
		label := task.Name + " (" + task.ID + ")"
		labels = append(labels, label)
		ids[label] = task.ID
	}
	return tasks, labels, ids
}

func (a *App) copyTriggerSet(set server.TriggerSetResponse) {
	backend, ok := a.backend.(triggerBackend)
	if !ok {
		return
	}
	go func() {
		ctx, cancel := a.bgCtx()
		defer cancel()
		result, err := backend.RevealTriggerSet(ctx, set.ID)
		if err != nil {
			fyne.Do(func() { a.showError(err) })
			return
		}
		fyne.Do(func() {
			a.clipboard.SetContent(triggerSetCommands(result))
			dialog.ShowInformation("Copied", fmt.Sprintf("Copied %d commands from %s.", len(result.Members), set.Name), a.win)
		})
	}()
}

func triggerSetCommands(result server.TriggerSetSecretResponse) string {
	var builder strings.Builder
	for _, member := range result.Members {
		builder.WriteString(member.Command)
		builder.WriteByte('\n')
	}
	return builder.String()
}

func (a *App) showTriggerSetRetarget(set server.TriggerSetResponse) {
	backend, ok := a.backend.(triggerBackend)
	if !ok {
		return
	}
	_, labels, ids := sortedTriggerTasks(a.model.Snapshot().Tasks)
	if len(labels) == 0 {
		return
	}
	target := widget.NewSelect(labels, nil)
	target.SetSelected(labels[0])
	form := dialog.NewForm("Retarget "+set.Name, "Retarget", "Cancel", []*widget.FormItem{widget.NewFormItem("New target for "+strconv.Itoa(set.MemberCount)+" members", target)}, func(ok bool) {
		if !ok {
			return
		}
		a.run(func(ctx context.Context) error {
			_, err := backend.RetargetTriggerSet(ctx, set.ID, ids[target.Selected])
			return err
		})
	}, a.win)
	form.Resize(fyne.NewSize(600, form.MinSize().Height))
	form.Show()
}

func (a *App) confirmTriggerSetEnabled(set server.TriggerSetResponse, enabled bool) {
	verb := "disable"
	if enabled {
		verb = "enable"
	}
	titleVerb := strings.ToUpper(verb[:1]) + verb[1:]
	dialog.ShowConfirm(titleVerb+" Trigger Set", fmt.Sprintf("%s all %d members of %s?", titleVerb, set.MemberCount, set.Name), func(ok bool) {
		if !ok {
			return
		}
		backend, supported := a.backend.(triggerBackend)
		if !supported {
			return
		}
		a.run(func(ctx context.Context) error {
			_, err := backend.SetTriggerSetEnabled(ctx, set.ID, enabled)
			return err
		})
	}, a.win)
}

func (a *App) confirmTriggerSetRotate(set server.TriggerSetResponse) {
	dialog.ShowConfirm("Rotate Trigger Set keys", fmt.Sprintf("Replace the keys for all %d members of %s? Every old key will stop working immediately.", set.MemberCount, set.Name), func(ok bool) {
		if !ok {
			return
		}
		backend, supported := a.backend.(triggerBackend)
		if !supported {
			return
		}
		go func() {
			ctx, cancel := a.bgCtx()
			defer cancel()
			result, err := backend.RotateTriggerSet(ctx, set.ID)
			if err != nil {
				fyne.Do(func() { a.showError(err) })
				return
			}
			fyne.Do(func() { a.showTriggerSetSecrets("Trigger Set keys rotated", result) })
			a.refreshAll()
		}()
	}, a.win)
}

func (a *App) confirmTriggerSetDelete(set server.TriggerSetResponse) {
	dialog.ShowConfirm("Delete Trigger Set", fmt.Sprintf("Delete %s and all %d members? Every key will stop working immediately.", set.Name, set.MemberCount), func(ok bool) {
		if !ok {
			return
		}
		backend, supported := a.backend.(triggerBackend)
		if !supported {
			return
		}
		a.run(func(ctx context.Context) error { return backend.DeleteTriggerSet(ctx, set.ID) })
	}, a.win)
}

func (a *App) showTriggerSetSecrets(title string, result server.TriggerSetSecretResponse) {
	commands := triggerSetCommands(result)
	content := widget.NewMultiLineEntry()
	content.SetText(commands)
	content.Disable()
	content.SetMinRowsVisible(8)
	copyStatus := widget.NewLabel("")
	copyStatus.Importance = widget.SuccessImportance
	copyButton := widget.NewButtonWithIcon("Copy all commands", theme.ContentCopyIcon(), func() {
		a.clipboard.SetContent(commands)
		copyStatus.SetText(fmt.Sprintf("Copied %d commands to the clipboard.", len(result.Members)))
	})
	dialog.ShowCustom(title, "Close", container.NewVBox(widget.NewLabel(fmt.Sprintf("%d ordered commands for %s. Store them securely.", len(result.Members), result.TriggerSet.Name)), content, copyButton, copyStatus), a.win)
}

func (a *App) runTriggerCreate(backend triggerBackend, request server.TriggerCreateRequest) {
	go func() {
		ctx, cancel := a.bgCtx()
		defer cancel()
		result, err := backend.CreateTrigger(ctx, request)
		if err != nil {
			fyne.Do(func() { a.showError(err) })
			return
		}
		fyne.Do(func() { a.showTriggerSecret("Trigger created", result) })
		a.refreshAll()
	}()
}

func (a *App) revealAndCopyTrigger(id string, command bool) {
	backend, ok := a.backend.(triggerBackend)
	if !ok {
		dialog.ShowError(fmt.Errorf("trigger administration is unavailable"), a.win)
		return
	}
	go func() {
		ctx, cancel := a.bgCtx()
		defer cancel()
		result, err := backend.RevealTrigger(ctx, id)
		if err != nil {
			fyne.Do(func() { a.showError(err) })
			return
		}
		fyne.Do(func() {
			if !command {
				a.showTriggerSecret("Trigger key", result)
				return
			}
			a.clipboard.SetContent(result.Command)
			dialog.ShowInformation("Copied", "The trigger command was copied to the clipboard.", a.win)
		})
	}()
}

func (a *App) rotateTrigger(id string) {
	backend, ok := a.backend.(triggerBackend)
	if !ok {
		return
	}
	go func() {
		ctx, cancel := a.bgCtx()
		defer cancel()
		result, err := backend.RotateTrigger(ctx, id)
		if err != nil {
			fyne.Do(func() { a.showError(err) })
			return
		}
		fyne.Do(func() { a.showTriggerSecret("Trigger key rotated", result) })
		a.refreshAll()
	}()
}

func (a *App) showTriggerSecret(title string, result server.TriggerSecretResponse) {
	key := widget.NewEntry()
	key.SetText(result.Key)
	key.Disable()
	command := widget.NewEntry()
	command.SetText(result.Command)
	command.Disable()
	copyStatus := widget.NewLabel("")
	copyStatus.Importance = widget.SuccessImportance
	copyKey := widget.NewButtonWithIcon("Copy key", theme.ContentCopyIcon(), func() {
		a.clipboard.SetContent(result.Key)
		copyStatus.SetText("Key copied to the clipboard.")
	})
	copyCommand := widget.NewButtonWithIcon("Copy command", theme.ContentCopyIcon(), func() {
		a.clipboard.SetContent(result.Command)
		copyStatus.SetText("Command copied to the clipboard.")
	})
	dialog.ShowCustom(title, "Close", container.NewVBox(widget.NewLabel("Store this key securely. Ordinary trigger views do not display it."), widget.NewLabel("Key"), key, copyKey, widget.NewLabel("Command"), command, copyCommand, copyStatus), a.win)
}
