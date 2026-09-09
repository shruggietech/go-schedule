package main

import (
	"context"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/desktop/agentaccess"
	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/desktop/notifications"
	"github.com/shruggietech/go-schedule/desktop/operations"
	"github.com/shruggietech/go-schedule/desktop/settings"
	"github.com/shruggietech/go-schedule/desktop/taskgroup"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type appBackend struct{}

func (appBackend) Health(context.Context) (connection.Health, error) {
	return connection.Health{ID: "daemon-1", DisplayName: "Workshop", Platform: "windows", Architecture: "amd64", Version: "1.2.0", Capabilities: []string{"tasks"}, Permissions: []string{"read"}}, nil
}
func (appBackend) StreamEvents(ctx context.Context, publish func(connection.DomainEvent)) error {
	<-ctx.Done()
	return ctx.Err()
}

type appEmitter struct {
	mu     sync.Mutex
	events []connection.Event
}

func (e *appEmitter) Emit(_ context.Context, name string, value any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if name == desktopEventName {
		e.events = append(e.events, value.(connection.Event))
	}
}

type appNative struct {
	quit   bool
	copied string
	opened string
}

func (n *appNative) Quit(context.Context) { n.quit = true }
func (n *appNative) ClipboardSetText(_ context.Context, value string) error {
	n.copied = value
	return nil
}
func (n *appNative) BrowserOpenURL(_ context.Context, value string) error {
	n.opened = value
	return nil
}

type facadeTaskBackend struct{ taskgroup.Backend }

type facadeOperationsBackend struct{ operations.Backend }

type facadeNotificationsBackend struct{ notifications.Backend }

type facadeSettingsBackend struct{}

type facadeAgentAccessBackend struct{ status server.MCPHTTPStatusResponse }

func (b *facadeAgentAccessBackend) MCPHTTPStatus(context.Context) (server.MCPHTTPStatusResponse, error) {
	return b.status, nil
}
func (b *facadeAgentAccessBackend) EnableMCPHTTP(_ context.Context, req server.MCPHTTPEnableRequest) (server.MCPHTTPCredentialResponse, error) {
	b.status = server.MCPHTTPStatusResponse{Enabled: true, ClientName: req.ClientName, AllowedOrigins: []string{}}
	return server.MCPHTTPCredentialResponse{MCPHTTPStatusResponse: b.status, Credential: "one-time"}, nil
}
func (b *facadeAgentAccessBackend) RotateMCPHTTPCredential(context.Context) (server.MCPHTTPCredentialResponse, error) {
	return server.MCPHTTPCredentialResponse{MCPHTTPStatusResponse: b.status, Credential: "replacement"}, nil
}
func (b *facadeAgentAccessBackend) DisableMCPHTTP(context.Context) (server.MCPHTTPStatusResponse, error) {
	b.status = server.MCPHTTPStatusResponse{AllowedOrigins: []string{}}
	return b.status, nil
}

func (facadeSettingsBackend) RuntimeInfo(context.Context) (server.RuntimeInfoResponse, error) {
	return server.RuntimeInfoResponse{}, nil
}

func (facadeOperationsBackend) GetCalendar(_ context.Context, from, to time.Time) (server.CalendarResponse, error) {
	return server.CalendarResponse{From: from, To: to, Occurrences: []server.Occurrence{}}, nil
}
func (facadeOperationsBackend) ListRuns(context.Context, string, int) ([]domain.Run, error) {
	return []domain.Run{}, nil
}
func (facadeOperationsBackend) ListActiveRuns(context.Context) ([]domain.Run, error) {
	return []domain.Run{}, nil
}
func (facadeOperationsBackend) ListLogs(context.Context, string, int) (server.LogsResponse, error) {
	return server.LogsResponse{Logs: []domain.LogRecord{}}, nil
}
func (facadeOperationsBackend) ListAlertsLimited(context.Context, bool, int) ([]domain.Alert, error) {
	return []domain.Alert{}, nil
}
func (facadeOperationsBackend) AckAlert(context.Context, string) error { return nil }

func (facadeNotificationsBackend) ListNotificationChannels(context.Context) ([]domain.NotificationChannel, error) {
	return []domain.NotificationChannel{}, nil
}
func (facadeNotificationsBackend) ListNotificationDeliveries(context.Context, domain.NotificationDeliveryFilter) ([]domain.NotificationDelivery, error) {
	return []domain.NotificationDelivery{}, nil
}
func (facadeNotificationsBackend) ListTaskDetails(context.Context, string, string) ([]server.TaskResponse, error) {
	return []server.TaskResponse{}, nil
}
func (facadeNotificationsBackend) ListGroups(context.Context) ([]domain.Group, error) {
	return []domain.Group{}, nil
}

func (facadeTaskBackend) ListTaskDetails(context.Context, string, string) ([]server.TaskResponse, error) {
	return []server.TaskResponse{}, nil
}

func (facadeTaskBackend) ListGroups(context.Context) ([]domain.Group, error) {
	return []domain.Group{}, nil
}

func TestAppFacadeStartsSnapshotsRetriesAndQuits(t *testing.T) {
	emitter := &appEmitter{}
	native := &appNative{}
	app := newApp(appBackend{}, emitter, native)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.startup(ctx)
	deadline := time.After(time.Second)
	for app.Snapshot().State != connection.StateConnected {
		select {
		case <-deadline:
			t.Fatal("connection did not start")
		default:
		}
	}
	snapshot := app.Snapshot()
	if snapshot.Target.ID != "daemon-1" || snapshot.Target.DisplayName != "Workshop" || snapshot.Target.Platform != "windows" || snapshot.Target.Architecture != "amd64" || snapshot.Target.Version != "1.2.0" {
		t.Fatalf("snapshot=%+v", snapshot)
	}
	if result := app.RetryConnection(); result.Outcome != "accepted" {
		t.Fatalf("retry=%+v", result)
	}
	if result := app.Quit(); result.Outcome != "accepted" || !native.quit {
		t.Fatalf("quit=%+v native=%+v", result, native)
	}
	app.shutdown(context.Background())
	if result := app.RetryConnection(); result.Outcome != "rejected" {
		t.Fatalf("post-close retry=%+v", result)
	}
	for _, event := range emitter.events {
		if event.Message == "" || event.Generation == 0 {
			t.Fatalf("unsafe event=%+v", event)
		}
	}
}

func TestAppFacadeExposesSafeTaskWorkspace(t *testing.T) {
	app := newApp(appBackend{}, nil, nil, appServices{tasks: taskgroup.NewService(facadeTaskBackend{})})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.startup(ctx)
	result := app.Workspace()
	if result.Outcome != "accepted" || result.Workspace == nil || result.Workspace.Tasks == nil || result.Workspace.Groups == nil {
		t.Fatalf("workspace=%+v", result)
	}
	app.shutdown(context.Background())
}

func TestAppFacadeExposesScheduleActivityAndAcknowledgement(t *testing.T) {
	service := operations.NewService(facadeOperationsBackend{})
	app := newApp(appBackend{}, nil, nil, appServices{operations: service})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.startup(ctx)
	if result := app.ScheduleWindow(7); result.Outcome != "accepted" || result.Schedule == nil {
		t.Fatalf("schedule=%+v", result)
	}
	if result := app.ActivityWorkspace(); result.Outcome != "accepted" || result.Activity == nil {
		t.Fatalf("activity=%+v", result)
	}
	if result := app.AcknowledgeAlert("alert-1"); result.Outcome != "accepted" || result.Activity == nil {
		t.Fatalf("acknowledge=%+v", result)
	}
	app.shutdown(context.Background())
}

func TestAppFacadeExposesSecretFreeNotificationWorkspace(t *testing.T) {
	service := notifications.NewService(facadeNotificationsBackend{})
	app := newApp(appBackend{}, nil, nil, appServices{notifications: service})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.startup(ctx)
	if result := app.NotificationWorkspace(); result.Outcome != "accepted" || result.Workspace == nil {
		t.Fatalf("notifications=%+v", result)
	}
	app.shutdown(context.Background())
}

func TestAppFacadeExposesDesktopSettings(t *testing.T) {
	native := &appNative{}
	deps := settings.DefaultDependencies()
	deps.Paths.Preferences = t.TempDir() + "/preferences.json"
	deps.Paths.LegacyPreferences = t.TempDir() + "/legacy.json"
	deps.Paths.ApplicationData = t.TempDir()
	service := settings.NewServiceWithDependencies(facadeSettingsBackend{}, native, deps)
	app := newApp(appBackend{}, nil, native, appServices{settings: service})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.startup(ctx)
	if result := app.SettingsWorkspace(); result.Outcome != "accepted" || result.Workspace == nil {
		t.Fatalf("settings=%+v", result)
	}
	if result := app.SaveAppearance("dark"); result.Outcome != "accepted" || result.Workspace.Preferences.Appearance != settings.AppearanceDark {
		t.Fatalf("appearance=%+v", result)
	}
	if result := app.RestoreDesktopPreferences(); result.Outcome != "accepted" || result.Workspace.Preferences.Appearance != settings.AppearanceSystem {
		t.Fatalf("restore=%+v", result)
	}
	if result := app.CopyStoragePath("desktop-preferences"); result.Outcome != "accepted" || native.copied != filepath.Clean(deps.Paths.Preferences) {
		t.Fatalf("copy=%+v copied=%q", result, native.copied)
	}
	if result := app.OpenProductLink("documentation"); result.Outcome != "accepted" || native.opened != "https://shruggietech.github.io/go-schedule/" {
		t.Fatalf("open=%+v opened=%q", result, native.opened)
	}
	app.shutdown(context.Background())
}

func TestAppFacadeExposesAgentAccessWithoutCredential(t *testing.T) {
	native := &appNative{}
	backend := &facadeAgentAccessBackend{status: server.MCPHTTPStatusResponse{AllowedOrigins: []string{}}}
	service := agentaccess.NewService(backend, native)
	app := newApp(appBackend{}, nil, native, appServices{agentAccess: service})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.startup(ctx)
	if result := app.AgentAccessWorkspace(); result.Outcome != "accepted" || result.Workspace == nil {
		t.Fatalf("workspace=%+v", result)
	}
	if result := app.EnableAgentAccess(agentaccess.EnableDraft{ClientName: "Codex", Port: 43123}); result.Outcome != "accepted" || native.copied != "one-time" || result.Workspace == nil {
		t.Fatalf("enable=%+v copied=%q", result, native.copied)
	}
	if result := app.RotateAgentAccess(); result.Outcome != "accepted" || native.copied != "replacement" {
		t.Fatalf("rotate=%+v copied=%q", result, native.copied)
	}
	if result := app.RevokeAgentAccess(); result.Outcome != "accepted" || result.Workspace == nil || result.Workspace.HTTP.Enabled {
		t.Fatalf("revoke=%+v", result)
	}
	app.shutdown(context.Background())
}
