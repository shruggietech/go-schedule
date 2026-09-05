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
	"github.com/shruggietech/go-schedule/internal/domain"
)

var watcherColumns = []structuredColumn{{Header: "Watcher", Minimum: 130, Weight: 2}, {Header: "Selection", Minimum: 160, Weight: 2}, {Header: "Target", Minimum: 140, Weight: 2}, {Header: "Enabled", Minimum: 75, Weight: 1}, {Header: "Health", Minimum: 100, Weight: 1}, {Header: "Readiness", Minimum: 110, Weight: 1}}

type watcherBackend interface {
	CreateFilesystemWatcher(context.Context, server.FilesystemWatcherCreateRequest) (server.FilesystemWatcherResponse, error)
	UpdateFilesystemWatcher(context.Context, string, server.FilesystemWatcherUpdateRequest) (server.FilesystemWatcherResponse, error)
	SetFilesystemWatcherEnabled(context.Context, string, bool) (server.FilesystemWatcherResponse, error)
	DeleteFilesystemWatcher(context.Context, string) error
}

func watcherRows(items []server.FilesystemWatcherResponse) []structuredRowModel {
	rows := make([]structuredRowModel, 0, len(items))
	for _, item := range items {
		selection := item.Path
		if item.Kind == domain.WatcherDirectory {
			selection += " (" + item.Pattern + ")"
		}
		enabled := "No"
		if item.Enabled {
			enabled = "Yes"
		}
		cells := []structuredCell{{Text: item.Name}, {Text: selection}, {Text: item.TargetTaskName, FullText: item.TargetTaskID}, {Text: enabled}, {Text: normalizedWords(string(item.Health.State), "Unknown", false), FullText: item.Health.Reason}, {Text: normalizedWords(item.Readiness, "Unknown", false), FullText: item.Reason}}
		rows = append(rows, structuredRowModel{Identity: item.ID, Cells: cells, Summary: structuredRowSummary(watcherColumns, cells)})
	}
	return rows
}

func (a *App) buildFilesystemWatcherPanel() fyne.CanvasObject {
	var items []server.FilesystemWatcherResponse
	selectedID := ""
	current := func() (*server.FilesystemWatcherResponse, bool) {
		for i := range items {
			if items[i].ID == selectedID {
				return &items[i], true
			}
		}
		return nil, false
	}
	withSelection := func(fn func(*server.FilesystemWatcherResponse)) {
		if item, ok := current(); ok {
			fn(item)
			return
		}
		dialog.ShowInformation("No selection", "Select a filesystem watcher first.", a.win)
	}
	table := newStructuredList(watcherColumns, "Select a watcher to see its full path and health detail.", func(id string) { selectedID = id }, func(id string) {
		selectedID = id
		if item, ok := current(); ok {
			a.showFilesystemWatcherEditor(item)
		}
	})
	a.watcherTable = table
	a.registerRefresher(func() { items = a.model.Snapshot().Watchers; table.setRows(watcherRows(items)) })
	newButton := newToolbarButton("New watcher", theme.ContentAddIcon(), func() { a.showFilesystemWatcherEditor(nil) })
	editButton := newToolbarButton("Edit watcher", theme.DocumentCreateIcon(), func() { withSelection(a.showFilesystemWatcherEditor) })
	toggleButton := newToolbarButton("Enable or disable watcher", theme.MediaReplayIcon(), func() {
		withSelection(func(item *server.FilesystemWatcherResponse) {
			backend, ok := a.backend.(watcherBackend)
			if !ok {
				dialog.ShowError(fmt.Errorf("watcher administration is unavailable"), a.win)
				return
			}
			a.run(func(ctx context.Context) error {
				_, err := backend.SetFilesystemWatcherEnabled(ctx, item.ID, !item.Enabled)
				return err
			})
		})
	})
	deleteButton := newToolbarButton("Delete watcher", theme.DeleteIcon(), func() {
		withSelection(func(item *server.FilesystemWatcherResponse) {
			dialog.ShowConfirm("Delete filesystem watcher", "Delete "+item.Name+"?", func(ok bool) {
				if !ok {
					return
				}
				backend, supported := a.backend.(watcherBackend)
				if !supported {
					dialog.ShowError(fmt.Errorf("watcher administration is unavailable"), a.win)
					return
				}
				a.run(func(ctx context.Context) error { return backend.DeleteFilesystemWatcher(ctx, item.ID) })
			}, a.win)
		})
	})
	return container.NewBorder(container.NewVBox(widget.NewLabelWithStyle("Filesystem watchers", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), container.NewHBox(newButton, editButton, toggleButton, deleteButton)), nil, nil, nil, table.root)
}

func (a *App) showFilesystemWatcherEditor(existing *server.FilesystemWatcherResponse) {
	backend, ok := a.backend.(watcherBackend)
	if !ok {
		dialog.ShowError(fmt.Errorf("watcher administration is unavailable"), a.win)
		return
	}
	tasks := append([]domain.Task(nil), a.model.Snapshot().Tasks...)
	if len(tasks) == 0 {
		dialog.ShowInformation("Task required", "Create a task before adding a filesystem watcher.", a.win)
		return
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].Name < tasks[j].Name })
	labels := make([]string, 0, len(tasks))
	ids := map[string]string{}
	byID := map[string]string{}
	for _, task := range tasks {
		label := task.Name + " (" + task.ID + ")"
		labels = append(labels, label)
		ids[label] = task.ID
		byID[task.ID] = label
	}
	name := widget.NewEntry()
	kind := widget.NewSelect([]string{string(domain.WatcherFile), string(domain.WatcherDirectory)}, nil)
	path := widget.NewEntry()
	pattern := widget.NewEntry()
	recursive := widget.NewCheck("Include descendant directories", nil)
	debounce := widget.NewEntry()
	stability := widget.NewEntry()
	target := widget.NewSelect(labels, nil)
	enabled := widget.NewCheck("Observe changes and invoke the target task", nil)
	title, confirm := "New filesystem watcher", "Create"
	if existing == nil {
		kind.SetSelected(string(domain.WatcherFile))
		target.SetSelected(labels[0])
		debounce.SetText("250ms")
		stability.SetText("500ms")
		enabled.SetChecked(true)
	} else {
		title, confirm = "Edit filesystem watcher", "Save"
		name.SetText(existing.Name)
		kind.SetSelected(string(existing.Kind))
		path.SetText(existing.Path)
		pattern.SetText(existing.Pattern)
		recursive.SetChecked(existing.Recursive)
		debounce.SetText(existing.Debounce)
		stability.SetText(existing.Stability)
		target.SetSelected(byID[existing.TargetTaskID])
		enabled.SetChecked(existing.Enabled)
	}
	formItems := []*widget.FormItem{widget.NewFormItem("Name", name), widget.NewFormItem("Kind", kind), widget.NewFormItem("Path", path), widget.NewFormItem("File-name pattern", pattern), widget.NewFormItem("Recursive", recursive), widget.NewFormItem("Debounce", debounce), widget.NewFormItem("Stability", stability), widget.NewFormItem("Target task", target), widget.NewFormItem("Enabled", enabled)}
	form := dialog.NewForm(title, confirm, "Cancel", formItems, func(ok bool) {
		if !ok {
			return
		}
		if name.Text == "" || path.Text == "" || ids[target.Selected] == "" {
			dialog.ShowError(fmt.Errorf("name, path, and target task are required"), a.win)
			return
		}
		request := server.FilesystemWatcherCreateRequest{Name: name.Text, Kind: domain.WatcherKind(kind.Selected), Path: path.Text, Pattern: pattern.Text, Recursive: recursive.Checked, Debounce: debounce.Text, Stability: stability.Text, TargetTaskID: ids[target.Selected], Enabled: &enabled.Checked}
		a.run(func(ctx context.Context) error {
			if existing == nil {
				_, err := backend.CreateFilesystemWatcher(ctx, request)
				return err
			}
			update := server.FilesystemWatcherUpdateRequest{Name: &request.Name, Kind: &request.Kind, Path: &request.Path, Pattern: &request.Pattern, Recursive: &request.Recursive, Debounce: &request.Debounce, Stability: &request.Stability, TargetTaskID: &request.TargetTaskID, Enabled: request.Enabled}
			_, err := backend.UpdateFilesystemWatcher(ctx, existing.ID, update)
			return err
		})
	}, a.win)
	form.Resize(fyne.NewSize(680, form.MinSize().Height))
	form.Show()
}
