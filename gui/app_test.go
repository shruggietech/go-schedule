package gui

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/events"
)

// fakeBackend implements Backend with in-memory data for headless UI tests.
type fakeBackend struct {
	tasks       []domain.Task
	groups      []domain.Group
	chains      []domain.CompletionChain
	alerts      []domain.Alert
	logs        []domain.LogRecord
	logPath     string
	runtimeInfo server.RuntimeInfoResponse
	runtimeErr  error
	created     int
	lastCreate  server.TaskCreateRequest

	// details keyed by task ID; GetTask serves these and records failures.
	details    map[string]server.TaskResponse
	getTaskErr error
	getTaskIDs []string

	// Updates are recorded under mu: App.run dispatches them on a goroutine, so
	// a test reading them races the UI unless both sides synchronize.
	mu          sync.Mutex
	updated     int
	lastUpdate  server.TaskUpdateRequest
	lastUpdID   string
	previews    int
	lastPreview server.PreviewRequest
	acked       []string
}

// lastUpdateCall returns the recorded update count and the most recent request.
func (f *fakeBackend) lastUpdateCall() (int, string, server.TaskUpdateRequest) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.updated, f.lastUpdID, f.lastUpdate
}

func (f *fakeBackend) acknowledgedAlerts() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.acked...)
}

// waitFor polls cond until it holds or the test times out, so tests can await an
// asynchronous backend call without sleeping for a fixed duration.
func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("timed out waiting for the expected backend call")
}

func TestRootContentFitsEffective800x600Launch(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	available := windowSizeFor(800, 600, 1)
	minimum := ui.win.Content().MinSize()
	if minimum.Width > available.Width || minimum.Height > available.Height {
		t.Fatalf("root minimum %v exceeds effective 800x600 launch viewport %v", minimum, available)
	}
}

func (f *fakeBackend) ListTasks(context.Context, string, string) ([]domain.Task, error) {
	return f.tasks, nil
}
func (f *fakeBackend) ListChains(context.Context) ([]domain.CompletionChain, error) {
	return f.chains, nil
}
func (f *fakeBackend) ListGroups(context.Context) ([]domain.Group, error) { return f.groups, nil }
func (f *fakeBackend) ListAlerts(context.Context, bool) ([]domain.Alert, error) {
	return f.alerts, nil
}
func (f *fakeBackend) ListLogs(context.Context, string, int) (server.LogsResponse, error) {
	return server.LogsResponse{Logs: f.logs, LogPath: f.logPath}, nil
}

func (f *fakeBackend) RuntimeInfo(context.Context) (server.RuntimeInfoResponse, error) {
	return f.runtimeInfo, f.runtimeErr
}
func (f *fakeBackend) CreateTask(_ context.Context, req server.TaskCreateRequest) (server.TaskResponse, error) {
	// Under the same mutex as the updates, and for the same reason: App.run
	// dispatches creates on a goroutine too.
	f.mu.Lock()
	defer f.mu.Unlock()
	f.created++
	f.lastCreate = req
	return server.TaskResponse{}, nil
}

// lastCreateCall returns the recorded create count and the most recent request.
func (f *fakeBackend) lastCreateCall() (int, server.TaskCreateRequest) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.created, f.lastCreate
}
func (f *fakeBackend) lastPreviewCall() (int, server.PreviewRequest) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.previews, f.lastPreview
}
func (f *fakeBackend) GetTask(_ context.Context, id string) (server.TaskResponse, error) {
	f.mu.Lock()
	f.getTaskIDs = append(f.getTaskIDs, id)
	f.mu.Unlock()
	if f.getTaskErr != nil {
		return server.TaskResponse{}, f.getTaskErr
	}
	if d, ok := f.details[id]; ok {
		return d, nil
	}
	for _, t := range f.tasks {
		if t.ID == id {
			return server.TaskResponse{Task: t}, nil
		}
	}
	return server.TaskResponse{}, nil
}

func (f *fakeBackend) requestedTaskDetails() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.getTaskIDs...)
}
func (f *fakeBackend) UpdateTask(_ context.Context, id string, req server.TaskUpdateRequest) (server.TaskResponse, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.updated++
	f.lastUpdID = id
	f.lastUpdate = req
	return server.TaskResponse{}, nil
}
func (f *fakeBackend) DeleteTask(context.Context, string) error           { return nil }
func (f *fakeBackend) SetTaskEnabled(context.Context, string, bool) error { return nil }
func (f *fakeBackend) RunNow(context.Context, string) error               { return nil }
func (f *fakeBackend) CreateChain(context.Context, server.ChainCreateRequest) (domain.CompletionChain, error) {
	return domain.CompletionChain{}, nil
}
func (f *fakeBackend) UpdateChain(context.Context, string, server.ChainUpdateRequest) (domain.CompletionChain, error) {
	return domain.CompletionChain{}, nil
}
func (f *fakeBackend) DeleteChain(context.Context, string) error { return nil }
func (f *fakeBackend) Preview(_ context.Context, req server.PreviewRequest) (server.PreviewResponse, error) {
	f.mu.Lock()
	f.previews++
	f.lastPreview = req
	f.mu.Unlock()
	return server.PreviewResponse{HumanSummary: "Every day at 09:00"}, nil
}
func (f *fakeBackend) CreateGroup(context.Context, server.GroupCreateRequest) (domain.Group, error) {
	return domain.Group{}, nil
}
func (f *fakeBackend) SetGroupEnabled(context.Context, string, bool) error { return nil }
func (f *fakeBackend) DeleteGroup(context.Context, string) error           { return nil }
func (f *fakeBackend) AckAlert(_ context.Context, id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.acked = append(f.acked, id)
	return nil
}
func (f *fakeBackend) GetCalendar(context.Context, time.Time, time.Time) (server.CalendarResponse, error) {
	return server.CalendarResponse{}, nil
}
func (f *fakeBackend) StreamEvents(ctx context.Context, _ func(events.Event)) error {
	<-ctx.Done()
	return ctx.Err()
}

func TestUI_BuildsCompleteNavigation(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{
		tasks:  []domain.Task{{ID: "t1", Name: "nightly", State: domain.TaskActive, Enabled: true, Timezone: "UTC"}},
		groups: []domain.Group{{ID: "g1", Name: "Backups", Enabled: true}},
		alerts: []domain.Alert{{ID: "a1", Kind: domain.AlertRunFailed, Message: "boom"}},
	})

	want := []string{"Tasks", "Groups", "Chains", "Schedule", "Activity", "Options", "Info"}
	if got := ui.navigation.labels(); !reflect.DeepEqual(got, want) {
		t.Fatalf("navigation = %v, want %v", got, want)
	}
}

func TestUI_WindowTitleIsBranded(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	if got := ui.win.Title(); got != "go-schedule" {
		t.Fatalf("window title = %q, want %q", got, "go-schedule")
	}
}

func TestUI_NoRefreshControls(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	// Walk the whole object tree; no button/label should read "Refresh" (FR-023).
	var walk func(o fyne.CanvasObject)
	walk = func(o fyne.CanvasObject) {
		switch w := o.(type) {
		case *cursorButton:
			if strings.Contains(w.Text, "Refresh") {
				t.Errorf("found a Refresh control: %q", w.Text)
			}
		case *widget.Button:
			if strings.Contains(w.Text, "Refresh") {
				t.Errorf("found a Refresh button: %q", w.Text)
			}
		case *fyne.Container:
			for _, c := range w.Objects {
				walk(c)
			}
		}
	}
	walk(ui.win.Content())
}

func TestUI_TaskEditorBuilds(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	// Opening the editor must not panic and the window keeps a canvas.
	ui.showTaskEditor(nil)
	if ui.win.Canvas() == nil {
		t.Fatal("window canvas missing")
	}
}

func TestTaskEditorPreventsStackingAndReleasesOnClose(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	ui.showTaskEditor(nil)
	first := ui.activeTaskDialog
	if first == nil {
		t.Fatal("first task editor was not opened")
	}
	ui.showTaskEditor(nil)
	if ui.activeTaskDialog != first {
		t.Fatal("second activation stacked another task editor")
	}
	first.Hide()
	if ui.activeTaskDialog != nil {
		t.Fatal("dialog close did not release active editor")
	}
	ui.showTaskEditor(nil)
	if ui.activeTaskDialog == nil || ui.activeTaskDialog == first {
		t.Fatal("editor did not reopen after close")
	}
	ui.activeTaskDialog.Hide()
}

func TestCurrentTaskByIDSurvivesReorderAndRejectsStale(t *testing.T) {
	tasks := []domain.Task{{ID: "second", Name: "Second"}, {ID: "first", Name: "First"}}
	got, ok := currentTaskByID(tasks, "first")
	if !ok || got.ID != "first" || got.Name != "First" {
		t.Fatalf("currentTaskByID after reorder = %+v, %v", got, ok)
	}
	if _, ok := currentTaskByID(tasks, "removed"); ok {
		t.Fatal("stale task ID resolved to an unrelated task")
	}
}

func TestEditTaskByIDUsesCurrentModelAndDegradedDetail(t *testing.T) {
	backend := &fakeBackend{
		tasks:      []domain.Task{{ID: "wanted", Name: "Wanted"}},
		getTaskErr: errors.New("temporary detail failure"),
	}
	ui := NewUI(testApp, backend)
	ui.model.OnChange = nil
	if err := ui.model.Refresh(context.Background()); err != nil {
		t.Fatal(err)
	}
	if ui.editTaskByID("stale") {
		t.Fatal("stale identity opened an editor")
	}
	if !ui.editTaskByID("wanted") {
		t.Fatal("current identity did not open an editor")
	}
	if ui.activeTaskDialog == nil {
		t.Fatal("degraded detail lookup did not preserve the edit fallback")
	}
	ui.activeTaskDialog.Hide()
}

func TestTaskListSingleDoubleRefreshAndStaleIdentity(t *testing.T) {
	backend := &fakeBackend{tasks: []domain.Task{
		{ID: "first", Name: "First", Timezone: "UTC"},
		{ID: "second", Name: "Second", Timezone: "UTC"},
	}}
	ui := NewUI(testApp, backend)
	ui.model.OnChange = nil
	if err := ui.model.Refresh(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, refresh := range ui.refreshers {
		refresh()
	}

	row := ui.taskList.CreateItem().(*structuredRow)
	ui.taskList.UpdateItem(0, row)
	test.Tap(row)
	if ui.activeTaskDialog != nil || len(backend.requestedTaskDetails()) != 0 {
		t.Fatal("single tap opened or fetched an editor")
	}
	test.DoubleTap(row)
	if got := backend.requestedTaskDetails(); !reflect.DeepEqual(got, []string{"first"}) {
		t.Fatalf("double activation fetched %v, want first", got)
	}
	ui.activeTaskDialog.Hide()

	backend.tasks = []domain.Task{
		{ID: "second", Name: "Second", Timezone: "UTC"},
		{ID: "first", Name: "First", Timezone: "UTC"},
	}
	if err := ui.model.Refresh(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, refresh := range ui.refreshers {
		refresh()
	}
	ui.taskList.UpdateItem(1, row)
	test.DoubleTap(row)
	if got := backend.requestedTaskDetails(); !reflect.DeepEqual(got, []string{"first", "first"}) {
		t.Fatalf("reordered row fetched %v, want stable first identity", got)
	}
	ui.activeTaskDialog.Hide()

	backend.tasks = []domain.Task{{ID: "second", Name: "Second", Timezone: "UTC"}}
	if err := ui.model.Refresh(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, refresh := range ui.refreshers {
		refresh()
	}
	test.DoubleTap(row)
	if got := backend.requestedTaskDetails(); !reflect.DeepEqual(got, []string{"first", "first"}) {
		t.Fatalf("stale row fetched details: %v", got)
	}
}

func TestTaskTableHasFixedHeadersAndFullValueDisclosure(t *testing.T) {
	backend := &fakeBackend{tasks: []domain.Task{{
		ID: "long", Name: "A very long Unicode task 備份", Enabled: true,
		State: domain.TaskActive, Timezone: "America/Argentina/Buenos_Aires",
	}}}
	ui := NewUI(testApp, backend)
	ui.model.OnChange = nil
	if err := ui.model.Refresh(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, refresh := range ui.refreshers {
		refresh()
	}
	if ui.taskTable == nil || ui.taskTable.header == nil || ui.taskTable.list != ui.taskList {
		t.Fatal("Tasks did not expose the shared fixed-header table")
	}
	wantHeaders := []string{"Task", "Enabled", "Lifecycle", "Time zone", "Group"}
	for index, want := range wantHeaders {
		if got := ui.taskTable.header.labels[index].Text; got != want {
			t.Errorf("header %d=%q, want %q", index, got, want)
		}
	}
	row := ui.taskList.CreateItem().(*structuredRow)
	ui.taskList.UpdateItem(0, row)
	for index, label := range row.labels {
		if label.Truncation != fyne.TextTruncateEllipsis {
			t.Errorf("cell %d truncation=%v", index, label.Truncation)
		}
	}
	ui.taskList.Select(0)
	for _, want := range []string{"Task: A very long Unicode task 備份", "Enabled: Enabled", "Lifecycle: Active", "Time zone: America/Argentina/Buenos_Aires", "Group: Not assigned"} {
		if !strings.Contains(ui.taskTable.disclosure.Text, want) {
			t.Errorf("disclosure %q missing %q", ui.taskTable.disclosure.Text, want)
		}
	}
}

func TestTaskListRefreshReconcilesVisibleSelectionByStableID(t *testing.T) {
	backend := &fakeBackend{tasks: []domain.Task{
		{ID: "first", Name: "First", Timezone: "UTC"},
		{ID: "second", Name: "Second", Timezone: "UTC"},
	}}
	ui := NewUI(testApp, backend)
	ui.model.OnChange = nil
	if err := ui.model.Refresh(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, refresh := range ui.refreshers {
		refresh()
	}
	ui.taskList.Select(0)

	var selected, unselected []widget.ListItemID
	originalSelected := ui.taskList.OnSelected
	originalUnselected := ui.taskList.OnUnselected
	ui.taskList.OnSelected = func(id widget.ListItemID) {
		selected = append(selected, id)
		originalSelected(id)
	}
	ui.taskList.OnUnselected = func(id widget.ListItemID) {
		unselected = append(unselected, id)
		originalUnselected(id)
	}

	backend.tasks = []domain.Task{
		{ID: "second", Name: "Second", Timezone: "UTC"},
		{ID: "first", Name: "First", Timezone: "UTC"},
	}
	if err := ui.model.Refresh(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, refresh := range ui.refreshers {
		refresh()
	}
	if !reflect.DeepEqual(unselected, []widget.ListItemID{0}) {
		t.Fatalf("unselected rows = %v, want [0]", unselected)
	}
	if !reflect.DeepEqual(selected, []widget.ListItemID{1}) {
		t.Fatalf("selected rows = %v, want [1]", selected)
	}

	backend.tasks = []domain.Task{{ID: "second", Name: "Second", Timezone: "UTC"}}
	if err := ui.model.Refresh(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, refresh := range ui.refreshers {
		refresh()
	}
	if !reflect.DeepEqual(unselected, []widget.ListItemID{0, 1}) {
		t.Fatalf("unselected rows after removal = %v, want [0 1]", unselected)
	}
	if !reflect.DeepEqual(selected, []widget.ListItemID{1}) {
		t.Fatalf("removed identity was reselected: %v", selected)
	}
}

func TestActivityTabLabel(t *testing.T) {
	tests := []struct {
		count int
		want  string
	}{
		{count: -1, want: "Activity"},
		{count: 0, want: "Activity"},
		{count: 1, want: "Activity (1)"},
		{count: 99, want: "Activity (99)"},
		{count: 100, want: "Activity (99+)"},
		{count: 12_345, want: "Activity (99+)"},
	}

	for _, tt := range tests {
		if got := activityTabLabel(tt.count); got != tt.want {
			t.Errorf("activityTabLabel(%d) = %q, want %q", tt.count, got, tt.want)
		}
	}
}

func TestEditorOwnershipGuard(t *testing.T) {
	ui := &App{}
	if !ui.claimTaskEditor() {
		t.Fatal("first editor claim should succeed")
	}
	if ui.claimTaskEditor() {
		t.Fatal("second editor claim should be rejected")
	}
	ui.releaseTaskEditor()
	if !ui.claimTaskEditor() {
		t.Fatal("claim should succeed after close release")
	}
}

func TestShutdownCoordinatorRunsCancellationAndCloseOnce(t *testing.T) {
	ui := &App{}
	var cancelCount atomic.Int32
	var closeCount atomic.Int32
	ui.runCancel = func() { cancelCount.Add(1) }
	ui.closeWindow = func() { closeCount.Add(1) }

	var wg sync.WaitGroup
	for range 32 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ui.requestClose()
			ui.stop()
		}()
	}
	wg.Wait()
	if got := cancelCount.Load(); got != 1 {
		t.Fatalf("cancellation count = %d, want 1", got)
	}
	if got := closeCount.Load(); got != 1 {
		t.Fatalf("close count = %d, want 1", got)
	}
}

func TestShutdownCancelsBackendContexts(t *testing.T) {
	runCtx, cancel := context.WithCancel(context.Background())
	ui := &App{runCtx: runCtx, runCancel: cancel, closeWindow: func() {}}
	backendCtx, backendCancel := ui.bgCtx()
	defer backendCancel()
	ui.requestClose()
	select {
	case <-backendCtx.Done():
	default:
		t.Fatal("orderly shutdown did not cancel an in-flight backend context")
	}
}

func TestUI_ActivityBadgeReflectsUnacked(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	if got := ui.navigation.labels(); !reflect.DeepEqual(got, []string{"Tasks", "Groups", "Chains", "Schedule", "Activity", "Options", "Info"}) {
		t.Fatalf("initial navigation = %v", got)
	}
	// Drive the badge synchronously: the production OnChange marshals through
	// fyne.Do on another goroutine, which would race with the assertion below.
	ui.model.OnChange = nil
	ui.model.ApplyEvent(events.Event{Kind: events.KindAlert, Alert: &domain.Alert{ID: "x", Acknowledged: false}})
	ui.updateLogsBadge()
	if got := ui.navigation.label(navigationActivity); got != "Activity (1)" {
		t.Fatalf("activity badge = %q, want Activity (1)", got)
	}
	if got := ui.navigation.labels(); !reflect.DeepEqual(got, []string{"Tasks", "Groups", "Chains", "Schedule", "Activity (1)", "Options", "Info"}) {
		t.Fatal("Activity badge update moved navigation destinations")
	}
}
