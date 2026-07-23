package gui

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/events"
)

// fakeBackend implements Backend with in-memory data for headless UI tests.
type fakeBackend struct {
	tasks      []domain.Task
	groups     []domain.Group
	alerts     []domain.Alert
	logs       []domain.LogRecord
	created    int
	lastCreate server.TaskCreateRequest

	// details keyed by task ID; GetTask serves these and records failures.
	details    map[string]server.TaskResponse
	getTaskErr error

	// Updates are recorded under mu: App.run dispatches them on a goroutine, so
	// a test reading them races the UI unless both sides synchronize.
	mu         sync.Mutex
	updated    int
	lastUpdate server.TaskUpdateRequest
	lastUpdID  string
}

// lastUpdateCall returns the recorded update count and the most recent request.
func (f *fakeBackend) lastUpdateCall() (int, string, server.TaskUpdateRequest) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.updated, f.lastUpdID, f.lastUpdate
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

func (f *fakeBackend) ListTasks(context.Context, string, string) ([]domain.Task, error) {
	return f.tasks, nil
}
func (f *fakeBackend) ListGroups(context.Context) ([]domain.Group, error) { return f.groups, nil }
func (f *fakeBackend) ListAlerts(context.Context, bool) ([]domain.Alert, error) {
	return f.alerts, nil
}
func (f *fakeBackend) ListLogs(context.Context, string, int) ([]domain.LogRecord, error) {
	return f.logs, nil
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
func (f *fakeBackend) GetTask(_ context.Context, id string) (server.TaskResponse, error) {
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
func (f *fakeBackend) Preview(context.Context, server.PreviewRequest) (server.PreviewResponse, error) {
	return server.PreviewResponse{HumanSummary: "Every day at 09:00"}, nil
}
func (f *fakeBackend) CreateGroup(context.Context, server.GroupCreateRequest) (domain.Group, error) {
	return domain.Group{}, nil
}
func (f *fakeBackend) SetGroupEnabled(context.Context, string, bool) error { return nil }
func (f *fakeBackend) DeleteGroup(context.Context, string) error           { return nil }
func (f *fakeBackend) AckAlert(context.Context, string) error              { return nil }
func (f *fakeBackend) GetCalendar(context.Context, time.Time, time.Time) (server.CalendarResponse, error) {
	return server.CalendarResponse{}, nil
}
func (f *fakeBackend) StreamEvents(ctx context.Context, _ func(events.Event)) error {
	<-ctx.Done()
	return ctx.Err()
}

func TestUI_BuildsAllTabs(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{
		tasks:  []domain.Task{{ID: "t1", Name: "nightly", State: domain.TaskActive, Enabled: true, Timezone: "UTC"}},
		groups: []domain.Group{{ID: "g1", Name: "Backups", Enabled: true}},
		alerts: []domain.Alert{{ID: "a1", Kind: domain.AlertRunFailed, Message: "boom"}},
	})

	want := []string{"Tasks", "Schedule", "Groups", "Logs"}
	if len(ui.tabs.Items) != len(want) {
		t.Fatalf("want %d tabs, got %d", len(want), len(ui.tabs.Items))
	}
	for i, w := range want {
		if ui.tabs.Items[i].Text != w {
			t.Fatalf("tab %d = %q, want %q", i, ui.tabs.Items[i].Text, w)
		}
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
	for _, tab := range ui.tabs.Items {
		walk(tab.Content)
	}
}

func TestUI_TaskEditorBuilds(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	// Opening the editor must not panic and the window keeps a canvas.
	ui.showTaskEditor(nil)
	if ui.win.Canvas() == nil {
		t.Fatal("window canvas missing")
	}
}

func TestUI_LogsBadgeReflectsUnacked(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	// Drive the badge synchronously: the production OnChange marshals through
	// fyne.Do on another goroutine, which would race with the assertion below.
	ui.model.OnChange = nil
	ui.model.ApplyEvent(events.Event{Kind: events.KindAlert, Alert: &domain.Alert{ID: "x", Acknowledged: false}})
	ui.updateLogsBadge()
	if ui.logsTab.Text != "Logs (1)" {
		t.Fatalf("logs badge = %q, want Logs (1)", ui.logsTab.Text)
	}
}
