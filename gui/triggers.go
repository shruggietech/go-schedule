package gui

import (
	"context"
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

var triggerColumns = []structuredColumn{
	{Header: "Trigger", Minimum: 130, Weight: 2},
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
}

func triggerRows(items []server.TriggerResponse) []structuredRowModel {
	rows := make([]structuredRowModel, 0, len(items))
	for _, item := range items {
		enabled := "No"
		if item.Enabled {
			enabled = "Yes"
		}
		cells := []structuredCell{{Text: item.Name}, {Text: item.TargetTaskName, FullText: item.TargetTaskName + " (" + item.TargetTaskID + ")"}, {Text: enabled}, {Text: normalizedWords(item.Readiness, "Unknown", false), FullText: item.Reason}, {Text: item.ID}}
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
	copyKey := newToolbarButton("Copy key", theme.ContentCopyIcon(), func() { withSelection(func(item *server.TriggerResponse) { a.revealAndCopyTrigger(item.ID, false) }) })
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
	toolbar := container.NewHBox(newButton, editButton, copyKey, copyCommand, toggle, rotate, remove)
	return container.NewBorder(toolbar, nil, nil, nil, table.root)
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
		a.run(func(ctx context.Context) error {
			_, err := backend.UpdateTrigger(ctx, existing.ID, server.TriggerUpdateRequest{Name: &name.Text, TargetTaskID: &targetID})
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
		value := result.Key
		if command {
			value = result.Command
		}
		fyne.Do(func() {
			a.clipboard.SetContent(value)
			dialog.ShowInformation("Copied", "The trigger value was copied to the clipboard.", a.win)
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
	copyKey := widget.NewButtonWithIcon("Copy key", theme.ContentCopyIcon(), func() { a.clipboard.SetContent(result.Key) })
	copyCommand := widget.NewButtonWithIcon("Copy command", theme.ContentCopyIcon(), func() { a.clipboard.SetContent(result.Command) })
	dialog.ShowCustom(title, "Close", container.NewVBox(widget.NewLabel("Store this key securely. Ordinary trigger views do not display it."), widget.NewLabel("Key"), key, copyKey, widget.NewLabel("Command"), command, copyCommand), a.win)
}
