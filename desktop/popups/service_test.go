package popups

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/desktop/settings"
	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type fakePreferences struct{ value settings.PopupPreferences }

func (f *fakePreferences) PopupPreferences() (settings.PopupPreferences, error) { return f.value, nil }

type fakePresenter struct {
	mu      sync.Mutex
	notices []Notice
	denied  bool
}

func (*fakePresenter) Available() bool    { return true }
func (f *fakePresenter) Authorized() bool { return !f.denied }
func (f *fakePresenter) Present(value Notice) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.notices = append(f.notices, value)
	return nil
}
func (f *fakePresenter) count() int { f.mu.Lock(); defer f.mu.Unlock(); return len(f.notices) }

type fakeClient struct {
	runs   map[string]domain.Run
	alerts []domain.Alert
	stream func(context.Context, func(client.RemoteEvent)) error
}

func (f *fakeClient) StreamRemoteEvents(ctx context.Context, send func(client.RemoteEvent)) error {
	if f.stream != nil {
		return f.stream(ctx, send)
	}
	return nil
}
func (f *fakeClient) GetRun(_ context.Context, id string) (domain.Run, error) { return f.runs[id], nil }
func (*fakeClient) GetTask(context.Context, string) (server.TaskResponse, error) {
	return server.TaskResponse{}, nil
}
func (f *fakeClient) ListAlertsLimited(context.Context, bool, int) ([]domain.Alert, error) {
	return f.alerts, nil
}

func TestDeduplicationUsesDaemonAndExactRecord(t *testing.T) {
	end := time.Now()
	client := &fakeClient{runs: map[string]domain.Run{"same": {ID: "same", TaskID: "task", Outcome: domain.OutcomeFailure, EndedAt: &end}}}
	prefs := &fakePreferences{value: settings.DefaultPopupPreferences()}
	prefs.value.Enabled = true
	presenter := &fakePresenter{}
	service := New(nil, prefs, presenter)
	first := Source{DaemonID: "daemon-a", Label: "Office", Client: client}
	second := Source{DaemonID: "daemon-b", Label: "Office", Client: client}
	event := clientEvent("run", "same")
	service.process(context.Background(), first, event)
	service.Refresh() // Reconnection must not clear the bounded identity cache.
	service.process(context.Background(), first, event)
	service.process(context.Background(), second, event)
	if presenter.count() != 2 {
		t.Fatalf("wanted one per daemon, got %d", presenter.count())
	}
	if presenter.notices[0].Title == presenter.notices[1].Title {
		t.Fatal("daemon attribution was lost")
	}
	if got, ok := service.Intent(presenter.notices[1].ID); !ok || got.DaemonID != "daemon-b" || got.RecordID != "same" {
		t.Fatalf("unsafe intent: %+v %v", got, ok)
	}
}

func TestMuteFiltersAndNonterminalRun(t *testing.T) {
	end := time.Now()
	client := &fakeClient{runs: map[string]domain.Run{"pending": {ID: "pending", Outcome: domain.OutcomeSuccess}, "done": {ID: "done", Outcome: domain.OutcomeSuccess, EndedAt: &end}}}
	prefs := &fakePreferences{value: settings.DefaultPopupPreferences()}
	presenter := &fakePresenter{}
	service := New(nil, prefs, presenter)
	source := Source{DaemonID: "daemon-a", Label: "A", Client: client}
	service.process(context.Background(), source, clientEvent("run", "done"))
	prefs.value.Enabled = true
	service.process(context.Background(), source, clientEvent("run", "pending"))
	prefs.value.Conditions = []string{"failure", "alert"}
	service.process(context.Background(), source, clientEvent("run", "done"))
	prefs.value.Conditions = []string{"success"}
	prefs.value.DaemonIDs = []string{"daemon-b"}
	service.process(context.Background(), source, clientEvent("run", "done"))
	if presenter.count() != 0 {
		t.Fatalf("excluded events produced %d popups", presenter.count())
	}
	prefs.value.DaemonIDs = nil
	prefs.value.Severities = []string{"error"} // Alert severity does not hide a selected successful-run condition.
	service.process(context.Background(), source, clientEvent("run", "done"))
	if presenter.count() != 1 {
		t.Fatalf("matching event count=%d", presenter.count())
	}
	prefs.value.Enabled = false
	service.process(context.Background(), source, clientEvent("run", "new"))
	if presenter.count() != 1 {
		t.Fatal("muted event was presented")
	}
}

func TestAlertSeverityAndCancelledContext(t *testing.T) {
	prefs := &fakePreferences{value: settings.DefaultPopupPreferences()}
	prefs.value.Enabled = true
	prefs.value.Severities = []string{"error"}
	presenter := &fakePresenter{}
	service := New(nil, prefs, presenter)
	source := Source{DaemonID: "daemon-a", Label: "A", Client: &fakeClient{alerts: []domain.Alert{{ID: "warning", Severity: domain.SeverityWarning, Kind: domain.AlertService}, {ID: "error", Severity: domain.SeverityError, Kind: domain.AlertService}}}}
	service.process(context.Background(), source, clientEvent("alert", "warning"))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	service.process(ctx, source, clientEvent("alert", "error"))
	if presenter.count() != 0 {
		t.Fatal("excluded or cancelled alert was presented")
	}
	service.process(context.Background(), source, clientEvent("alert", "error"))
	if presenter.count() != 1 {
		t.Fatalf("matching alert count=%d", presenter.count())
	}
}

func clientEvent(kind, id string) client.RemoteEvent {
	return client.RemoteEvent{Kind: kind, ResourceID: id}
}

type fakeHealth struct{ id string }

func (f fakeHealth) Health(context.Context) (connection.Health, error) {
	return connection.Health{ID: f.id}, nil
}

func TestStreamStopsOnMute(t *testing.T) {
	prefs := &fakePreferences{value: settings.DefaultPopupPreferences()}
	prefs.value.Enabled = true
	started, exited := make(chan struct{}, 1), make(chan struct{}, 1)
	events := &fakeClient{stream: func(ctx context.Context, _ func(client.RemoteEvent)) error {
		started <- struct{}{}
		<-ctx.Done()
		exited <- struct{}{}
		return ctx.Err()
	}}
	service := New(func() []Source {
		return []Source{{DaemonID: "daemon-a", Label: "A", Client: events, Backend: fakeHealth{id: "daemon-a"}}}
	}, prefs, &fakePresenter{})
	service.Start(context.Background())
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("stream did not start")
	}
	prefs.value.Enabled = false
	service.Refresh()
	select {
	case <-exited:
	case <-time.After(time.Second):
		t.Fatal("stream did not stop after mute")
	}
	service.Stop()
}

func TestDeniedPermissionDoesNotStartObservation(t *testing.T) {
	prefs := &fakePreferences{value: settings.DefaultPopupPreferences()}
	prefs.value.Enabled = true
	loaded := false
	service := New(func() []Source { loaded = true; return nil }, prefs, &fakePresenter{denied: true})
	service.Start(context.Background())
	if loaded {
		t.Fatal("denied notifications still opened daemon sources")
	}
	service.Stop()
}

func TestPopupTextIsSingleLineAndBounded(t *testing.T) {
	if got := compact("  One\nTwo\tThree  ", 20); got != "One Two Three" {
		t.Fatalf("unsafe popup text %q", got)
	}
	if got := compact("abcdefghijkl", 6); got != "abcde…" {
		t.Fatalf("unbounded popup text %q", got)
	}
}
