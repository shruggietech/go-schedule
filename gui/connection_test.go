package gui

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"fyne.io/fyne/v2"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type recoveringBackend struct {
	fakeBackend
	mu        sync.Mutex
	failing   bool
	listCalls int
}

func (b *recoveringBackend) ListTasks(context.Context, string, string) ([]domain.Task, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listCalls++
	if b.failing {
		return nil, client.NewConnectionError("GET /v1/tasks", os.ErrPermission)
	}
	return b.tasks, nil
}

func (b *recoveringBackend) setFailing(failing bool) {
	b.mu.Lock()
	b.failing = failing
	b.mu.Unlock()
}

func (b *recoveringBackend) calls() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.listCalls
}

func TestConnectionIncidentDeduplicatesConcurrentSources(t *testing.T) {
	t.Parallel()
	state := newConnectionState(func() accessDiagnosis { return accessDiagnosis{} })
	errorsBySource := []*client.ConnectionError{
		client.NewConnectionError("GET /v1/tasks", os.ErrPermission),
		client.NewConnectionError("GET /v1/calendar", os.ErrPermission),
		client.NewConnectionError("GET /v1/events", os.ErrPermission),
	}
	var wg sync.WaitGroup
	for _, err := range errorsBySource {
		wg.Add(1)
		go func(err *client.ConnectionError) {
			defer wg.Done()
			state.report(err)
		}(err)
	}
	wg.Wait()
	incident, ok := state.snapshot()
	if !ok {
		t.Fatal("no active incident")
	}
	if incident.VisibleCount != 1 {
		t.Fatalf("visible count = %d, want 1", incident.VisibleCount)
	}
	if incident.Kind != client.ConnectionAccessDenied {
		t.Fatalf("kind = %q", incident.Kind)
	}
}

func TestConnectionPanelKeepsRecoveryControlsReachable(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	ui.showError(client.NewConnectionError("GET /v1/tasks", os.ErrPermission))
	ui.showError(client.NewConnectionError("GET /v1/calendar", os.ErrPermission))
	ui.showError(client.NewConnectionError("GET /v1/events", os.ErrPermission))
	incident, active := ui.connection.snapshot()
	if !active || incident.VisibleCount != 1 {
		t.Fatalf("incident active=%v visible=%d, want active single incident", active, incident.VisibleCount)
	}
	if !ui.connectionCard.Visible() || ui.retryButton.Text != "Retry" || ui.navigation == nil {
		t.Fatalf("panel=%v retry=%q navigation=%v", ui.connectionCard.Visible(), ui.retryButton.Text, ui.navigation != nil)
	}
}

func TestSuccessfulRetryClearsIncidentWithoutReinstall(t *testing.T) {
	backend := &recoveringBackend{failing: true}
	ui := NewUI(testApp, backend)
	if err := ui.refreshAllOnce(); err == nil {
		t.Fatal("initial refresh unexpectedly succeeded")
	}
	waitFor(t, func() bool {
		_, active := ui.connection.snapshot()
		return active
	})
	backend.setFailing(false)
	fyne.DoAndWait(ui.retryConnection)
	if got := backend.calls(); got != 1 {
		t.Fatalf("Retry performed %d refreshes directly, want none", got-1)
	}
	select {
	case <-ui.retrySignal:
	default:
		t.Fatal("Retry did not notify the reconnect coordinator")
	}
	if err := ui.refreshAllOnce(); err != nil {
		t.Fatalf("coordinated retry failed: %v", err)
	}
	waitFor(t, func() bool {
		_, active := ui.connection.snapshot()
		return !active
	})
	if got := backend.calls(); got != 2 {
		t.Fatalf("ListTasks calls = %d, want initial attempt plus one coordinated retry", got)
	}
}

func TestFailedRetryReEnablesControlForAPIStatusError(t *testing.T) {
	ui := NewUI(testApp, &statusErrorBackend{fakeBackend: fakeBackend{}})
	ui.connection.report(client.NewConnectionError("GET /v1/tasks", os.ErrPermission))
	ui.connection.setRetrying()
	if err := ui.refreshAllOnce(); err == nil {
		t.Fatal("retry unexpectedly succeeded")
	}
	incident, active := ui.connection.snapshot()
	if !active || incident.Retrying {
		t.Fatalf("incident active=%v retrying=%v, want active retryable incident", active, incident.Retrying)
	}
}

type statusErrorBackend struct {
	fakeBackend
}

func (b *statusErrorBackend) ListTasks(context.Context, string, string) ([]domain.Task, error) {
	return nil, &client.StatusError{Code: "internal", Message: "temporary API failure"}
}

func TestAPIStatusErrorDoesNotCreateTransportIncident(t *testing.T) {
	ui := NewUI(testApp, &fakeBackend{})
	if ui.recordConnectionError(&client.StatusError{Code: "conflict", Message: "conflict"}) {
		t.Fatal("API response error activated connection incident")
	}
}

func TestAccessDeniedGuidanceRequiresVerifiedStaleToken(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		diagnose accessDiagnosis
		contains string
	}{
		{name: "stale token", diagnose: accessDiagnosis{GroupExists: yes, AccountMember: yes, TokenMember: no}, contains: "sign out"},
		{name: "not enrolled", diagnose: accessDiagnosis{GroupExists: yes, AccountMember: no, TokenMember: no}, contains: "not a member"},
		{name: "unknown", diagnose: accessDiagnosis{}, contains: "could not verify"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			state := newConnectionState(func() accessDiagnosis { return tt.diagnose })
			state.report(client.NewConnectionError("pipe", os.ErrPermission))
			incident, _ := state.snapshot()
			if !containsFold(incident.Guidance, tt.contains) {
				t.Fatalf("guidance = %q, want fragment %q", incident.Guidance, tt.contains)
			}
		})
	}
}

func TestConnectionCopyByCategory(t *testing.T) {
	t.Parallel()
	tests := []struct {
		kind  client.ConnectionFailureKind
		title string
	}{
		{client.ConnectionUnavailable, "Daemon unavailable"},
		{client.ConnectionAccessDenied, "Access denied"},
		{client.ConnectionTimeout, "Connection timed out"},
		{client.ConnectionOtherTransport, "Connection interrupted"},
	}
	for _, tt := range tests {
		state := newConnectionState(func() accessDiagnosis { return accessDiagnosis{} })
		state.report(&client.ConnectionError{Kind: tt.kind, Operation: "test", Cause: errors.New("failure")})
		incident, _ := state.snapshot()
		if incident.Title != tt.title {
			t.Fatalf("kind %q title = %q, want %q", tt.kind, incident.Title, tt.title)
		}
	}
}

func TestReconnectBackoffIsBounded(t *testing.T) {
	t.Parallel()
	delay := time.Duration(0)
	want := []time.Duration{2 * time.Second, 4 * time.Second, 8 * time.Second, 16 * time.Second, 30 * time.Second, 30 * time.Second}
	for i, expected := range want {
		delay = nextReconnectDelay(delay)
		if delay != expected {
			t.Fatalf("delay %d = %s, want %s", i, delay, expected)
		}
	}
}

func TestReconnectBackoffResetsOnlyAfterStreamEvent(t *testing.T) {
	t.Parallel()
	if got := reconnectDelayAfterAttempt(2*time.Second, false); got != 4*time.Second {
		t.Fatalf("delay without stream event = %s, want 4s", got)
	}
	if got := reconnectDelayAfterAttempt(16*time.Second, true); got != 2*time.Second {
		t.Fatalf("delay after stream event = %s, want 2s", got)
	}
}
