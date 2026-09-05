// Package gui implements the go-schedule desktop GUI with Fyne. Its widget
// construction is cgo-free (Fyne's headless test driver renders without OpenGL),
// so the UI is unit-tested here; only the real windowed application entry point
// (cmd/gosched-gui) imports the GL driver and requires cgo.
//
// The GUI talks to the daemon exclusively through the Backend interface (the API
// client implements it), so it operates on the same state as the CLI and is
// fully testable with a fake backend.
package gui

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/gui/viewmodel"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/events"
	"github.com/shruggietech/go-schedule/internal/platform"
	"github.com/shruggietech/go-schedule/internal/winuninstall"
)

// Backend is everything the GUI needs from the daemon. The API client satisfies
// it; tests inject a fake.
type Backend interface {
	viewmodel.API // ListTasks, ListGroups, ListAlerts

	CreateTask(ctx context.Context, req server.TaskCreateRequest) (server.TaskResponse, error)
	// GetTask returns the full task detail including its schedule. The cached
	// task list carries no schedule, so the editor fetches detail on open to
	// show what a task is actually set to.
	GetTask(ctx context.Context, id string) (server.TaskResponse, error)
	UpdateTask(ctx context.Context, id string, req server.TaskUpdateRequest) (server.TaskResponse, error)
	DeleteTask(ctx context.Context, id string) error
	SetTaskEnabled(ctx context.Context, id string, enabled bool) error
	RunNow(ctx context.Context, id string) error
	GetRun(ctx context.Context, id string) (domain.Run, error)
	Preview(ctx context.Context, req server.PreviewRequest) (server.PreviewResponse, error)
	CreateChain(ctx context.Context, req server.ChainCreateRequest) (domain.CompletionChain, error)
	UpdateChain(ctx context.Context, id string, req server.ChainUpdateRequest) (domain.CompletionChain, error)
	DeleteChain(ctx context.Context, id string) error

	CreateGroup(ctx context.Context, req server.GroupCreateRequest) (domain.Group, error)
	SetGroupEnabled(ctx context.Context, id string, enabled bool) error
	DeleteGroup(ctx context.Context, id string) error

	AckAlert(ctx context.Context, id string) error
	GetCalendar(ctx context.Context, from, to time.Time) (server.CalendarResponse, error)
	RuntimeInfo(ctx context.Context) (server.RuntimeInfoResponse, error)
	StreamEvents(ctx context.Context, onEvent func(events.Event)) error
}

// App is the GUI application.
type App struct {
	fyne             fyne.App
	win              fyne.Window
	clipboard        fyne.Clipboard
	backend          Backend
	model            *viewmodel.Model
	appearance       appearancePreferences
	storageLocations []storageLocation
	storageInputs    storageLocationInputs
	options          *optionsView
	taskList         *widget.List
	taskTable        *structuredList
	scheduleTable    *structuredList
	activityTable    *structuredList
	triggerTable     *structuredList

	navigation *navigationShell
	refreshers []func()

	connection         *connectionState
	connectionCard     *widget.Card
	connectionTitle    *widget.Label
	connectionGuidance *widget.Label
	connectionDetail   *widget.Label
	retryButton        *widget.Button
	runCtx             context.Context
	runCancel          context.CancelFunc
	retrySignal        chan struct{}
	stopOnce           sync.Once
	closeOnce          sync.Once
	closeWindow        func()
	editorMu           sync.Mutex
	taskEditorOpen     bool
	activeTaskDialog   dialog.Dialog
}

// NewUI builds the GUI against fyneApp (created by the caller with the GL driver)
// and backend. It constructs the window content but does not show it.
func NewUI(fyneApp fyne.App, backend Backend) *App {
	runCtx, runCancel := context.WithCancel(context.Background())
	a := &App{
		fyne:        fyneApp,
		backend:     backend,
		model:       viewmodel.New(backend),
		connection:  newConnectionState(diagnoseAccess),
		runCtx:      runCtx,
		runCancel:   runCancel,
		retrySignal: make(chan struct{}, 1),
	}
	// Apply the validated per-user theme before constructing widgets. Dark and
	// System remain the safe defaults, and the windowed entry point and headless
	// test driver exercise the same palette and font selection.
	a.appearance = loadAppearancePreferences(fyneApp.Preferences())
	applyBrandTheme(fyneApp.Settings(), a.appearance)
	fyneApp.SetIcon(appIcon)
	a.win = fyneApp.NewWindow("go-schedule")
	a.clipboard = fyneApp.Clipboard()
	executablePath, _ := os.Executable()
	preferencesRoot := ""
	if root := fyneApp.Storage().RootURI(); root != nil && root.Scheme() == "file" {
		preferencesRoot = root.Path()
	}
	ownedMachineDataRoot := platform.DataDir()
	maintenanceEvidencePath := ""
	if runtime.GOOS == "windows" {
		maintenanceEvidencePath = winuninstall.CleanupResultPath(filepath.Dir(ownedMachineDataRoot))
	}
	a.storageInputs = storageLocationInputs{
		OwnedMachineDataRoot:    ownedMachineDataRoot,
		PreferencesRoot:         preferencesRoot,
		ExecutablePath:          executablePath,
		GOOS:                    runtime.GOOS,
		MaintenanceEvidencePath: maintenanceEvidencePath,
	}
	a.storageLocations = resolveStorageLocations(a.storageInputs)
	a.win.SetIcon(windowIcon) // crisp small tile for the title bar (see icon.go)
	// Open as a bounded restored window on the launch monitor, respecting its
	// taskbar and display scale. Unknown work area retains the 1280x800 fallback.
	ww, wh, monitorScale := workAreaPx()
	if monitorScale <= 0 {
		monitorScale = a.win.Canvas().Scale()
	}
	a.win.Resize(windowSizeFor(ww, wh, monitorScale))
	a.win.CenterOnScreen()
	a.win.SetContent(a.buildRoot())
	a.closeWindow = func() {
		a.win.SetCloseIntercept(nil)
		a.win.Close()
	}
	a.win.SetCloseIntercept(a.requestClose)
	a.model.OnChange = func() { fyne.Do(a.onModelChange) }
	return a
}

// buildRoot assembles the leading navigation rail and active content.
func (a *App) buildRoot() fyne.CanvasObject {
	a.navigation = newNavigationShell([]navigationDestinationSpec{
		{ID: navigationTasks, Label: "Tasks", Content: a.buildTasksTab(), Section: navigationDefinitions},
		{ID: navigationGroups, Label: "Groups", Content: a.buildGroupsTab(), Section: navigationDefinitions},
		{ID: navigationChains, Label: "Chains", Content: a.buildChainsTab(), Section: navigationDefinitions},
		{ID: navigationTriggers, Label: "Triggers", Content: a.buildTriggersTab(), Section: navigationDefinitions},
		{ID: navigationSchedule, Label: "Schedule", Content: a.buildScheduleTab(), Section: navigationOperations},
		{ID: navigationActivity, Label: activityTabLabel(0), Content: a.buildLogsTab(), Section: navigationOperations},
		{ID: navigationOptions, Label: "Options", Content: a.buildOptionsTab(), Section: navigationOperations},
		{ID: navigationInfo, Label: "Info", Content: a.buildInfoTab(), Section: navigationOperations},
	}, a.requestClose)
	a.connectionTitle = widget.NewLabel("")
	a.connectionTitle.TextStyle = fyne.TextStyle{Bold: true}
	a.connectionGuidance = widget.NewLabel("")
	a.connectionGuidance.Wrapping = fyne.TextWrapWord
	a.connectionDetail = widget.NewLabel("")
	a.connectionDetail.Wrapping = fyne.TextWrapWord
	a.retryButton = widget.NewButton("Retry", a.retryConnection)
	exitButton := widget.NewButton("Exit", a.requestClose)
	a.connectionCard = widget.NewCard("Connection", "", container.NewVBox(
		a.connectionTitle,
		a.connectionGuidance,
		a.connectionDetail,
		container.NewHBox(a.retryButton, exitButton),
	))
	a.connectionCard.Hide()
	return container.NewBorder(a.connectionCard, nil, nil, nil, a.navigation.root)
}

// Run shows the window, kicks off the first data load and the live event stream,
// and blocks until the window is closed.
func (a *App) Run() {
	a.refreshAll()
	go a.streamEvents()
	a.win.Show()
	if err := writeWindowEvidence(a.win.Canvas().Size(), a.win.Canvas().Scale()); err != nil {
		slog.Error("write attended window evidence", "error", err)
	}
	a.fyne.Run()
	a.stop()
}

// refreshAll reloads model state and every tab's local view.
func (a *App) refreshAll() {
	go func() { _ = a.refreshAllOnce() }()
}

func (a *App) refreshAllOnce() error {
	ctx, cancel := context.WithTimeout(a.runCtx, 10*time.Second)
	defer cancel()
	if err := a.model.Refresh(ctx); err != nil {
		fyne.Do(func() {
			a.connection.finishRetry()
			a.showError(err)
			a.renderConnectionIncident()
		})
		return err
	}
	a.refreshStorageLocations(ctx)
	a.connection.clear()
	fyne.Do(a.renderConnectionIncident)
	for _, r := range a.refreshers {
		rr := r
		fyne.Do(rr)
	}
	return nil
}

func (a *App) refreshStorageLocations(ctx context.Context) {
	runtimeInfo, err := a.backend.RuntimeInfo(ctx)
	if err != nil {
		return
	}
	storageInputs := a.storageInputs
	storageInputs.Runtime = runtimeInfo
	storageLocations := resolveStorageLocations(storageInputs)
	fyne.Do(func() {
		a.storageLocations = storageLocations
		if a.options != nil {
			a.options.setStorageLocations(storageLocations, a.clipboard)
		}
	})
}

// streamEvents consumes the SSE stream and folds events into the model,
// reconnecting after a short delay if the stream drops.
func (a *App) streamEvents() {
	delay := time.Duration(0)
	for {
		var streamRecovered atomic.Bool
		err := a.backend.StreamEvents(a.runCtx, func(e events.Event) {
			streamRecovered.Store(true)
			a.model.ApplyEvent(e)
			if e.Kind == events.KindTrigger || e.Kind == events.KindTriggerSet || e.Kind == events.KindTask || e.Kind == events.KindGroup {
				a.refreshAll()
			}
		})
		if a.runCtx.Err() != nil {
			return
		}
		if err == nil {
			return
		}
		if !a.reportConnectionError(err) {
			fyne.Do(func() { a.showError(err) })
			return
		}
		delay = reconnectDelayAfterAttempt(delay, streamRecovered.Load())
		if !waitForReconnect(a.runCtx, delay, a.retrySignal) {
			return
		}
		_ = a.refreshAllOnce()
	}
}

func (a *App) retryConnection() {
	a.connection.setRetrying()
	a.renderConnectionIncident()
	select {
	case a.retrySignal <- struct{}{}:
	default:
	}
}

func (a *App) reportConnectionError(err error) bool {
	if !a.recordConnectionError(err) {
		return false
	}
	fyne.Do(a.renderConnectionIncident)
	return true
}

func (a *App) recordConnectionError(err error) bool {
	connectionErr, ok := asConnectionError(err)
	if !ok {
		return false
	}
	a.connection.report(connectionErr)
	return true
}

func (a *App) renderConnectionIncident() {
	if a.connectionCard == nil {
		return
	}
	incident, active := a.connection.snapshot()
	if !active {
		a.connectionCard.Hide()
		return
	}
	a.connectionTitle.SetText(incident.Title)
	a.connectionGuidance.SetText(incident.Guidance)
	a.connectionDetail.SetText(incident.Detail)
	if incident.Retrying {
		a.retryButton.SetText("Retrying…")
		a.retryButton.Disable()
	} else {
		a.retryButton.SetText("Retry")
		a.retryButton.Enable()
	}
	a.connectionCard.Show()
}

func (a *App) stop() {
	a.stopOnce.Do(func() {
		if a.runCancel != nil {
			a.runCancel()
		}
	})
}

func (a *App) requestClose() {
	a.closeOnce.Do(func() {
		a.stop()
		if a.closeWindow != nil {
			a.closeWindow()
		}
	})
}

func (a *App) claimTaskEditor() bool {
	a.editorMu.Lock()
	defer a.editorMu.Unlock()
	if a.taskEditorOpen {
		return false
	}
	a.taskEditorOpen = true
	return true
}

func (a *App) releaseTaskEditor() {
	a.editorMu.Lock()
	a.taskEditorOpen = false
	a.activeTaskDialog = nil
	a.editorMu.Unlock()
}

func (a *App) setActiveTaskDialog(dialog dialog.Dialog) {
	a.editorMu.Lock()
	a.activeTaskDialog = dialog
	a.editorMu.Unlock()
}

// onModelChange refreshes the Activity badge and every destination view when state changes.
func (a *App) onModelChange() {
	a.updateLogsBadge()
	for _, r := range a.refreshers {
		r()
	}
}

// updateLogsBadge shows the bounded count of unacknowledged alerts (the
// actionable subset of activity) on the Activity destination.
func (a *App) updateLogsBadge() {
	if a.navigation == nil {
		return
	}
	a.navigation.updateLabel(navigationActivity, activityTabLabel(a.model.UnacknowledgedAlerts()))
}

func activityTabLabel(unacknowledged int) string {
	if unacknowledged <= 0 {
		return "Activity"
	}
	if unacknowledged > 99 {
		return "Activity (99+)"
	}
	return "Activity (" + itoa(unacknowledged) + ")"
}

func (a *App) registerRefresher(f func()) { a.refreshers = append(a.refreshers, f) }

// bgCtx returns a short-lived context for a backend call.
func (a *App) bgCtx() (context.Context, context.CancelFunc) {
	base := a.runCtx
	if base == nil {
		base = context.Background()
	}
	return context.WithTimeout(base, 10*time.Second)
}

// run executes a backend mutation in the background and refreshes on success.
func (a *App) run(fn func(ctx context.Context) error) {
	go func() {
		ctx, cancel := a.bgCtx()
		defer cancel()
		if err := fn(ctx); err != nil {
			fyne.Do(func() { a.showError(err) })
			return
		}
		a.refreshAll()
	}()
}
