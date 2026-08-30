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

const (
	chainOutcomeSuccess = "Only when the source succeeds"
	chainOutcomeFailure = "Only when the source fails"
	chainOutcomeAny     = "Whenever the source finishes"
)

func chainRowText(chain domain.CompletionChain) string {
	return fmt.Sprintf("%s  →  %s   [%s]", chain.SourceTaskName, chain.TargetTaskName, chainOutcomeLabel(chain.OnOutcome))
}

func chainOutcomeLabel(outcome domain.CompletionOutcome) string {
	switch outcome {
	case domain.CompletionOnSuccess:
		return "on success"
	case domain.CompletionOnFailure:
		return "on failure"
	default:
		return "on any completion"
	}
}

func chainOutcomeChoice(outcome domain.CompletionOutcome) string {
	switch outcome {
	case domain.CompletionOnFailure:
		return chainOutcomeFailure
	case domain.CompletionOnAny:
		return chainOutcomeAny
	default:
		return chainOutcomeSuccess
	}
}

func chainOutcomeForChoice(choice string) string {
	switch choice {
	case chainOutcomeFailure:
		return string(domain.CompletionOnFailure)
	case chainOutcomeAny:
		return string(domain.CompletionOnAny)
	default:
		return string(domain.CompletionOnSuccess)
	}
}

func chainTaskLabels(tasks []domain.Task) ([]string, map[string]string, map[string]string) {
	sorted := append([]domain.Task(nil), tasks...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Name == sorted[j].Name {
			return sorted[i].ID < sorted[j].ID
		}
		return sorted[i].Name < sorted[j].Name
	})
	labels := make([]string, 0, len(sorted))
	ids := make(map[string]string, len(sorted))
	byID := make(map[string]string, len(sorted))
	for _, task := range sorted {
		label := task.Name + " (" + task.ID + ")"
		labels = append(labels, label)
		ids[label] = task.ID
		byID[task.ID] = label
	}
	return labels, ids, byID
}

func (a *App) buildChainsTab() fyne.CanvasObject {
	var chains []domain.CompletionChain
	selected := -1
	list := widget.NewList(
		func() int { return len(chains) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(chainRowText(chains[id]))
		},
	)
	list.OnSelected = func(id widget.ListItemID) { selected = id }
	list.OnUnselected = func(widget.ListItemID) { selected = -1 }
	empty := widget.NewLabel("No completion chains yet. Create one to run a target task after a source task finishes.")
	empty.Alignment = fyne.TextAlignCenter
	refresh := func() {
		chains = a.model.Snapshot().Chains
		if selected >= len(chains) {
			selected = -1
			list.UnselectAll()
		}
		if len(chains) == 0 {
			list.Hide()
			empty.Show()
		} else {
			empty.Hide()
			list.Show()
		}
		list.Refresh()
	}
	a.registerRefresher(refresh)

	current := func() (*domain.CompletionChain, bool) {
		if selected < 0 || selected >= len(chains) {
			return nil, false
		}
		chain := chains[selected]
		return &chain, true
	}
	withSelection := func(fn func(*domain.CompletionChain)) {
		if chain, ok := current(); ok {
			fn(chain)
			return
		}
		dialog.ShowInformation("No selection", "Select a completion chain first.", a.win)
	}

	newButton := newToolbarButton("New", theme.ContentAddIcon(), func() { a.showChainEditor(nil) })
	editButton := newToolbarButton("Edit", theme.DocumentCreateIcon(), func() {
		withSelection(a.showChainEditor)
	})
	deleteButton := newToolbarButton("Delete", theme.DeleteIcon(), func() {
		withSelection(func(chain *domain.CompletionChain) {
			dialog.ShowConfirm("Delete completion chain", "Delete "+chain.SourceTaskName+" → "+chain.TargetTaskName+"?", func(ok bool) {
				if ok {
					a.run(func(ctx context.Context) error { return a.backend.DeleteChain(ctx, chain.ID) })
				}
			}, a.win)
		})
	})
	toolbar := container.NewHBox(newButton, editButton, deleteButton)
	return container.NewBorder(toolbar, nil, nil, nil, container.NewStack(list, empty))
}

func (a *App) showChainEditor(existing *domain.CompletionChain) {
	tasks := a.model.Snapshot().Tasks
	labels, ids, byID := chainTaskLabels(tasks)
	if len(labels) < 2 {
		dialog.ShowInformation("Two tasks required", "Create at least two tasks before adding a completion chain.", a.win)
		return
	}
	source := widget.NewSelect(labels, nil)
	target := widget.NewSelect(labels, nil)
	outcome := widget.NewSelect([]string{chainOutcomeSuccess, chainOutcomeFailure, chainOutcomeAny}, nil)
	if existing == nil {
		source.SetSelected(labels[0])
		target.SetSelected(labels[1])
		outcome.SetSelected(chainOutcomeSuccess)
	} else {
		source.SetSelected(byID[existing.SourceTaskID])
		target.SetSelected(byID[existing.TargetTaskID])
		outcome.SetSelected(chainOutcomeChoice(existing.OnOutcome))
	}
	title, confirm := "New completion chain", "Create"
	if existing != nil {
		title, confirm = "Edit completion chain", "Save"
	}
	items := []*widget.FormItem{
		widget.NewFormItem("After task", source),
		widget.NewFormItem("Run task", target),
		widget.NewFormItem("Condition", outcome),
	}
	dialog.NewForm(title, confirm, "Cancel", items, func(ok bool) {
		if !ok {
			return
		}
		sourceID, targetID := ids[source.Selected], ids[target.Selected]
		if sourceID == "" || targetID == "" || sourceID == targetID {
			dialog.ShowError(fmt.Errorf("source and target must be different tasks"), a.win)
			return
		}
		on := chainOutcomeForChoice(outcome.Selected)
		if existing == nil {
			a.run(func(ctx context.Context) error {
				_, err := a.backend.CreateChain(ctx, server.ChainCreateRequest{SourceTaskID: sourceID, TargetTaskID: targetID, OnOutcome: on})
				return err
			})
			return
		}
		a.run(func(ctx context.Context) error {
			_, err := a.backend.UpdateChain(ctx, existing.ID, server.ChainUpdateRequest{SourceTaskID: &sourceID, TargetTaskID: &targetID, OnOutcome: &on})
			return err
		})
	}, a.win).Show()
}
