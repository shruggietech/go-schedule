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
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/gui/viewmodel"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/events"
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
	Preview(ctx context.Context, req server.PreviewRequest) (server.PreviewResponse, error)
	CreateChain(ctx context.Context, req server.ChainCreateRequest) (domain.CompletionChain, error)
	UpdateChain(ctx context.Context, id string, req server.ChainUpdateRequest) (domain.CompletionChain, error)
	DeleteChain(ctx context.Context, id string) error

	CreateGroup(ctx context.Context, req server.GroupCreateRequest) (domain.Group, error)
	SetGroupEnabled(ctx context.Context, id string, enabled bool) error
	DeleteGroup(ctx context.Context, id string) error

	AckAlert(ctx context.Context, id string) error
	GetCalendar(ctx context.Context, from, to time.Time) (server.CalendarResponse, error)
	StreamEvents(ctx context.Context, onEvent func(events.Event)) error
}

// App is the GUI application.
type App struct {
	fyne    fyne.App
	win     fyne.Window
	backend Backend
	model   *viewmodel.Model

	tabs       *container.AppTabs
	logsTab    *container.TabItem
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
	// Apply the go-schedule brand theme (dark-first palette + brand fonts). Wired
	// here so the windowed entry point and the headless test driver share one
	// theme, and every NewUI-based test exercises the real palette.
	applyBrandTheme(fyneApp.Settings())
	fyneApp.SetIcon(appIcon)
	a.win = fyneApp.NewWindow("go-schedule")
	a.win.SetIcon(windowIcon) // crisp small tile for the title bar (see icon.go)
	// Open at the screen work area (maximized appearance), respecting the taskbar
	// (FR-001). Falls back to a generous size where the work area is unknown.
	ww, wh := workAreaPx()
	a.win.Resize(windowSizeFor(ww, wh, a.win.Canvas().Scale()))
	a.win.CenterOnScreen()
	a.win.SetContent(a.buildRoot())
	a.win.SetCloseIntercept(func() {
		a.stop()
		a.win.SetCloseIntercept(nil)
		a.win.Close()
	})
	a.model.OnChange = func() { fyne.Do(a.onModelChange) }
	return a
}

// buildRoot assembles the tabbed layout.
func (a *App) buildRoot() fyne.CanvasObject {
	a.tabs = container.NewAppTabs(
		container.NewTabItem("Tasks", a.buildTasksTab()),
		container.NewTabItem("Groups", a.buildGroupsTab()),
		container.NewTabItem("Chains", a.buildChainsTab()),
		container.NewTabItem("Schedule", a.buildScheduleTab()),
	)
	a.logsTab = container.NewTabItem(activityTabLabel(0), a.buildLogsTab())
	a.tabs.Append(a.logsTab)
	a.tabs.Append(container.NewTabItem("Info", a.buildInfoTab()))
	a.tabs.SetTabLocation(container.TabLocationLeading)
	a.connectionTitle = widget.NewLabel("")
	a.connectionTitle.TextStyle = fyne.TextStyle{Bold: true}
	a.connectionGuidance = widget.NewLabel("")
	a.connectionGuidance.Wrapping = fyne.TextWrapWord
	a.connectionDetail = widget.NewLabel("")
	a.connectionDetail.Wrapping = fyne.TextWrapWord
	a.retryButton = widget.NewButton("Retry", a.retryConnection)
	exitButton := widget.NewButton("Exit", func() {
		a.stop()
		a.win.SetCloseIntercept(nil)
		a.win.Close()
	})
	a.connectionCard = widget.NewCard("Connection", "", container.NewVBox(
		a.connectionTitle,
		a.connectionGuidance,
		a.connectionDetail,
		container.NewHBox(a.retryButton, exitButton),
	))
	a.connectionCard.Hide()
	return container.NewBorder(a.connectionCard, nil, nil, nil, a.tabs)
}

// Run shows the window, kicks off the first data load and the live event stream,
// and blocks until the window is closed.
func (a *App) Run() {
	a.refreshAll()
	go a.streamEvents()
	a.win.ShowAndRun()
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
		fyne.Do(func() { a.showError(err) })
		return err
	}
	a.connection.clear()
	fyne.Do(a.renderConnectionIncident)
	for _, r := range a.refreshers {
		rr := r
		fyne.Do(rr)
	}
	return nil
}

// streamEvents consumes the SSE stream and folds events into the model,
// reconnecting after a short delay if the stream drops.
func (a *App) streamEvents() {
	delay := time.Duration(0)
	for {
		err := a.backend.StreamEvents(a.runCtx, func(e events.Event) {
			a.model.ApplyEvent(e)
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
		delay = nextReconnectDelay(delay)
		if !waitForReconnect(a.runCtx, delay, a.retrySignal) {
			return
		}
		if err := a.refreshAllOnce(); err == nil {
			delay = 0
		}
	}
}

func (a *App) retryConnection() {
	a.connection.setRetrying()
	a.renderConnectionIncident()
	select {
	case a.retrySignal <- struct{}{}:
	default:
	}
	a.refreshAll()
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
	if a.runCancel != nil {
		a.runCancel()
	}
}

// onModelChange refreshes the Activity badge and every tab's view when state changes.
func (a *App) onModelChange() {
	a.updateLogsBadge()
	for _, r := range a.refreshers {
		r()
	}
}

// updateLogsBadge shows the bounded count of unacknowledged alerts (the
// actionable subset of activity) on the Activity tab.
func (a *App) updateLogsBadge() {
	if a.logsTab == nil {
		return
	}
	a.logsTab.Text = activityTabLabel(a.model.UnacknowledgedAlerts())
	if a.tabs != nil {
		a.tabs.Refresh()
	}
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
	return context.WithTimeout(context.Background(), 10*time.Second)
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
