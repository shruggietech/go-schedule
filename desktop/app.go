package main

import (
	"context"
	"time"

	"github.com/shruggietech/go-schedule/desktop/agentaccess"
	"github.com/shruggietech/go-schedule/desktop/automation"
	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/desktop/notifications"
	"github.com/shruggietech/go-schedule/desktop/operations"
	"github.com/shruggietech/go-schedule/desktop/remotepairing"
	"github.com/shruggietech/go-schedule/desktop/settings"
	"github.com/shruggietech/go-schedule/desktop/taskgroup"
)

const desktopEventName = "desktop:event"

type eventEmitter interface {
	Emit(context.Context, string, any)
}
type nativeRuntime interface {
	Quit(context.Context)
	ClipboardSetText(context.Context, string) error
	BrowserOpenURL(context.Context, string) error
}

// ActionResult gives the frontend a stable result without exposing native errors.
type ActionResult struct {
	Action  string `json:"action"`
	Outcome string `json:"outcome"`
	Message string `json:"message"`
}

// App is the deliberately small Wails bridge facade.
type App struct {
	manager       *connection.Manager
	tasks         *taskgroup.Service
	automation    *automation.Service
	operations    *operations.Service
	notifications *notifications.Service
	settings      *settings.Service
	agentAccess   *agentaccess.Service
	remotePairing *remotepairing.Service
	emitter       eventEmitter
	native        nativeRuntime
	ctx           context.Context
}

type appServices struct {
	tasks         *taskgroup.Service
	automation    *automation.Service
	operations    *operations.Service
	notifications *notifications.Service
	settings      *settings.Service
	agentAccess   *agentaccess.Service
	remotePairing *remotepairing.Service
}

func newApp(backend connection.Backend, emitter eventEmitter, native nativeRuntime, services ...appServices) *App {
	app := &App{emitter: emitter, native: native}
	if len(services) > 0 {
		app.tasks = services[0].tasks
		app.automation = services[0].automation
		app.operations = services[0].operations
		app.notifications = services[0].notifications
		app.settings = services[0].settings
		app.agentAccess = services[0].agentAccess
		app.remotePairing = services[0].remotePairing
	}
	app.manager = connection.NewManager(backend, appObserver{app: app})
	return app
}

func (a *App) PairRemote(draft remotepairing.Draft) remotepairing.Result {
	if a.remotePairing == nil || a.ctx == nil {
		return remotepairing.Result{Action: "pair_remote", Outcome: "unavailable", Message: "Remote pairing is unavailable."}
	}
	return a.remotePairing.Pair(a.ctx, draft)
}

// AgentAccessWorkspace returns the safe local MCP projection.
func (a *App) AgentAccessWorkspace() agentaccess.Result {
	if a.agentAccess == nil || a.ctx == nil {
		return agentaccess.Result{Action: "load_agent_access", Outcome: "unavailable", Message: "Agent Access is unavailable."}
	}
	return a.agentAccess.Workspace(a.ctx)
}

// EnableAgentAccess enables one named localhost client and copies its credential natively.
func (a *App) EnableAgentAccess(draft agentaccess.EnableDraft) agentaccess.Result {
	if a.agentAccess == nil || a.ctx == nil {
		return agentaccess.Result{Action: "enable_agent_access", Outcome: "unavailable", Message: "Agent Access is unavailable."}
	}
	return a.agentAccess.Enable(a.ctx, draft)
}

// RotateAgentAccess replaces the current localhost credential.
func (a *App) RotateAgentAccess() agentaccess.Result {
	if a.agentAccess == nil || a.ctx == nil {
		return agentaccess.Result{Action: "rotate_agent_access", Outcome: "unavailable", Message: "Agent Access is unavailable."}
	}
	return a.agentAccess.Rotate(a.ctx)
}

// RevokeAgentAccess closes the localhost listener and invalidates its credential.
func (a *App) RevokeAgentAccess() agentaccess.Result {
	if a.agentAccess == nil || a.ctx == nil {
		return agentaccess.Result{Action: "revoke_agent_access", Outcome: "unavailable", Message: "Agent Access is unavailable."}
	}
	return a.agentAccess.Revoke(a.ctx)
}

// OpenAgentAccessGuide opens the fixed official setup guide.
func (a *App) OpenAgentAccessGuide() agentaccess.Result {
	if a.agentAccess == nil || a.ctx == nil {
		return agentaccess.Result{Action: "open_agent_access_guide", Outcome: "unavailable", Message: "Agent Access is unavailable."}
	}
	return a.agentAccess.OpenGuide(a.ctx)
}

func (a *App) NotificationWorkspace() notifications.Result {
	if a.notifications == nil || a.ctx == nil {
		return notifications.Result{Action: "load_notifications", Outcome: "unavailable", Message: "Notifications are unavailable."}
	}
	return a.notifications.Workspace(a.ctx)
}
func (a *App) SaveNotificationChannel(d notifications.ChannelDraft) notifications.Result {
	if a.notifications == nil || a.ctx == nil {
		return notifications.Result{Action: "save_notification_channel", Outcome: "unavailable", Message: "Notifications are unavailable."}
	}
	return a.notifications.SaveChannel(a.ctx, d)
}
func (a *App) SetNotificationChannelEnabled(id string, enabled bool) notifications.Result {
	if a.notifications == nil || a.ctx == nil {
		return notifications.Result{Action: "toggle_notification_channel", Outcome: "unavailable", Message: "Notifications are unavailable."}
	}
	return a.notifications.SetChannelEnabled(a.ctx, id, enabled)
}
func (a *App) TestNotificationChannel(id string) notifications.Result {
	if a.notifications == nil || a.ctx == nil {
		return notifications.Result{Action: "test_notification_channel", Outcome: "unavailable", Message: "Notifications are unavailable."}
	}
	return a.notifications.TestChannel(a.ctx, id)
}
func (a *App) DeleteNotificationChannel(id string) notifications.Result {
	if a.notifications == nil || a.ctx == nil {
		return notifications.Result{Action: "delete_notification_channel", Outcome: "unavailable", Message: "Notifications are unavailable."}
	}
	return a.notifications.DeleteChannel(a.ctx, id)
}
func (a *App) NotificationPolicy(scopeType, scopeID string) notifications.Result {
	if a.notifications == nil || a.ctx == nil {
		return notifications.Result{Action: "load_notification_policy", Outcome: "unavailable", Message: "Notifications are unavailable."}
	}
	return a.notifications.Policy(a.ctx, scopeType, scopeID)
}
func (a *App) SaveNotificationPolicy(d notifications.PolicyDraft) notifications.Result {
	if a.notifications == nil || a.ctx == nil {
		return notifications.Result{Action: "save_notification_policy", Outcome: "unavailable", Message: "Notifications are unavailable."}
	}
	return a.notifications.SavePolicy(a.ctx, d)
}

// SettingsWorkspace returns desktop-local settings and storage information.
func (a *App) SettingsWorkspace() settings.Result {
	if a.settings == nil || a.ctx == nil {
		return settings.Result{Action: "load_settings", Outcome: "unavailable", Message: "Desktop settings are unavailable."}
	}
	return a.settings.Workspace(a.ctx)
}

// SaveAppearance persists one validated desktop appearance.
func (a *App) SaveAppearance(value string) settings.Result {
	if a.settings == nil || a.ctx == nil {
		return settings.Result{Action: "save_appearance", Outcome: "unavailable", Message: "Desktop settings are unavailable."}
	}
	return a.settings.SaveAppearance(a.ctx, value)
}

// RestoreDesktopPreferences restores only desktop-local preference defaults.
func (a *App) RestoreDesktopPreferences() settings.Result {
	if a.settings == nil || a.ctx == nil {
		return settings.Result{Action: "restore_preferences", Outcome: "unavailable", Message: "Desktop settings are unavailable."}
	}
	return a.settings.Restore(a.ctx)
}

// CopyStoragePath copies one current backend-resolved storage path.
func (a *App) CopyStoragePath(id string) settings.Result {
	if a.settings == nil || a.ctx == nil {
		return settings.Result{Action: "copy_storage_path", Outcome: "unavailable", Message: "Desktop settings are unavailable."}
	}
	return a.settings.CopyStoragePath(a.ctx, id)
}

// OpenProductLink opens one fixed product destination in the system browser.
func (a *App) OpenProductLink(key string) settings.Result {
	if a.settings == nil || a.ctx == nil {
		return settings.Result{Action: "open_product_link", Outcome: "unavailable", Message: "Desktop settings are unavailable."}
	}
	return a.settings.OpenProductLink(a.ctx, key)
}

func (a *App) ScheduleWindow(days int) operations.OperationResult {
	if a.operations == nil || a.ctx == nil {
		return operations.OperationResult{Action: "load_schedule", Outcome: "unavailable", Message: "Schedule is unavailable."}
	}
	return a.operations.ScheduleWindow(a.ctx, days)
}

func (a *App) ActivityWorkspace() operations.OperationResult {
	if a.operations == nil || a.ctx == nil {
		return operations.OperationResult{Action: "load_activity", Outcome: "unavailable", Message: "Activity is unavailable."}
	}
	return a.operations.ActivityWorkspace(a.ctx)
}

func (a *App) AcknowledgeAlert(id string) operations.OperationResult {
	if a.operations == nil || a.ctx == nil {
		return operations.OperationResult{Action: "acknowledge_alerts", Outcome: "unavailable", Message: "Activity is unavailable."}
	}
	return a.operations.AcknowledgeAlert(a.ctx, id)
}

func (a *App) AcknowledgeAlerts(ids []string) operations.OperationResult {
	if a.operations == nil || a.ctx == nil {
		return operations.OperationResult{Action: "acknowledge_alerts", Outcome: "unavailable", Message: "Activity is unavailable."}
	}
	return a.operations.AcknowledgeAlerts(a.ctx, ids)
}

func (a *App) AutomationWorkspace() automation.OperationResult {
	if a.automation == nil || a.ctx == nil {
		return automation.OperationResult{Action: "load", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.Workspace(a.ctx)
}
func (a *App) SaveChain(d automation.ChainDraft) automation.OperationResult {
	if a.automation == nil || a.ctx == nil {
		return automation.OperationResult{Action: "save_chain", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.SaveChain(a.ctx, d)
}
func (a *App) DeleteChain(id string) automation.OperationResult {
	if a.automation == nil || a.ctx == nil {
		return automation.OperationResult{Action: "delete_chain", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.DeleteChain(a.ctx, id)
}
func (a *App) SaveTrigger(d automation.TriggerDraft) automation.SecretResult {
	if a.automation == nil || a.ctx == nil {
		return automation.SecretResult{Action: "save_trigger", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.SaveTrigger(a.ctx, d)
}
func (a *App) SetTriggerEnabled(id string, enabled bool) automation.OperationResult {
	if a.automation == nil || a.ctx == nil {
		return automation.OperationResult{Action: "toggle_trigger", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.SetTriggerEnabled(a.ctx, id, enabled)
}
func (a *App) RevealTrigger(id string) automation.SecretResult {
	if a.automation == nil || a.ctx == nil {
		return automation.SecretResult{Action: "reveal_trigger", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.RevealTrigger(a.ctx, id)
}
func (a *App) RotateTrigger(id string) automation.SecretResult {
	if a.automation == nil || a.ctx == nil {
		return automation.SecretResult{Action: "rotate_trigger", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.RotateTrigger(a.ctx, id)
}
func (a *App) FireTrigger(id string) automation.OperationResult {
	if a.automation == nil || a.ctx == nil {
		return automation.OperationResult{Action: "fire_trigger", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.FireTrigger(a.ctx, id)
}
func (a *App) DeleteTrigger(id string) automation.OperationResult {
	if a.automation == nil || a.ctx == nil {
		return automation.OperationResult{Action: "delete_trigger", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.DeleteTrigger(a.ctx, id)
}
func (a *App) CreateTriggerSet(d automation.TriggerSetDraft) automation.SecretResult {
	if a.automation == nil || a.ctx == nil {
		return automation.SecretResult{Action: "create_trigger_set", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.CreateTriggerSet(a.ctx, d)
}
func (a *App) RetargetTriggerSet(id, target, updated string, overwrite bool) automation.OperationResult {
	if a.automation == nil || a.ctx == nil {
		return automation.OperationResult{Action: "retarget_trigger_set", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.RetargetTriggerSet(a.ctx, id, target, updated, overwrite)
}
func (a *App) SetTriggerSetEnabled(id string, enabled bool) automation.OperationResult {
	if a.automation == nil || a.ctx == nil {
		return automation.OperationResult{Action: "toggle_trigger_set", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.SetTriggerSetEnabled(a.ctx, id, enabled)
}
func (a *App) RevealTriggerSet(id string) automation.SecretResult {
	if a.automation == nil || a.ctx == nil {
		return automation.SecretResult{Action: "reveal_trigger_set", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.RevealTriggerSet(a.ctx, id)
}
func (a *App) RotateTriggerSet(id string) automation.SecretResult {
	if a.automation == nil || a.ctx == nil {
		return automation.SecretResult{Action: "rotate_trigger_set", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.RotateTriggerSet(a.ctx, id)
}
func (a *App) DeleteTriggerSet(id string) automation.OperationResult {
	if a.automation == nil || a.ctx == nil {
		return automation.OperationResult{Action: "delete_trigger_set", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.DeleteTriggerSet(a.ctx, id)
}
func (a *App) SaveFilesystemWatcher(d automation.WatcherDraft) automation.OperationResult {
	if a.automation == nil || a.ctx == nil {
		return automation.OperationResult{Action: "save_watcher", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.SaveWatcher(a.ctx, d)
}
func (a *App) SetFilesystemWatcherEnabled(id string, enabled bool) automation.OperationResult {
	if a.automation == nil || a.ctx == nil {
		return automation.OperationResult{Action: "toggle_watcher", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.SetWatcherEnabled(a.ctx, id, enabled)
}
func (a *App) DeleteFilesystemWatcher(id string) automation.OperationResult {
	if a.automation == nil || a.ctx == nil {
		return automation.OperationResult{Action: "delete_watcher", Outcome: "unavailable", Message: "Automation sources are unavailable."}
	}
	return a.automation.DeleteWatcher(a.ctx, id)
}

func (a *App) Workspace() taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "load", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.Workspace(a.ctx)
}

func (a *App) Task(id string) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "load", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.Task(a.ctx, id)
}

func (a *App) PreviewTask(draft taskgroup.TaskDraft) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "preview", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.PreviewTask(a.ctx, draft)
}
func (a *App) SaveTask(draft taskgroup.TaskDraft) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "save_task", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.SaveTask(a.ctx, draft)
}
func (a *App) RunTask(id string) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "run_task", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.RunTask(a.ctx, id)
}
func (a *App) SetTaskEnabled(id string, enabled bool) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "toggle_task", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.SetTaskEnabled(a.ctx, id, enabled)
}
func (a *App) DeleteTask(id string) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "delete_task", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.DeleteTask(a.ctx, id)
}
func (a *App) SaveGroup(draft taskgroup.GroupDraft) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "save_group", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.SaveGroup(a.ctx, draft)
}
func (a *App) SetGroupEnabled(id string, enabled bool) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "toggle_group", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.SetGroupEnabled(a.ctx, id, enabled)
}
func (a *App) DeleteGroup(id string) taskgroup.OperationResult {
	if a.tasks == nil || a.ctx == nil {
		return taskgroup.OperationResult{Action: "delete_group", Outcome: "unavailable", Message: "Task authoring is unavailable."}
	}
	return a.tasks.DeleteGroup(a.ctx, id)
}

type appObserver struct{ app *App }

func (o appObserver) Publish(event connection.Event) {
	if o.app.emitter != nil && o.app.ctx != nil {
		o.app.emitter.Emit(o.app.ctx, desktopEventName, event)
	}
}

func (a *App) startup(ctx context.Context) { a.ctx = ctx; a.manager.Start(ctx) }

func (a *App) shutdown(context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = a.manager.Stop(ctx)
}

// Snapshot returns only the transport-neutral connection contract.
func (a *App) Snapshot() connection.Snapshot { return a.manager.Snapshot() }

// RetryConnection requests a fresh, coalesced generation.
func (a *App) RetryConnection() ActionResult {
	if !a.manager.Retry() {
		return ActionResult{Action: "retry", Outcome: "rejected", Message: "The application is closing."}
	}
	return ActionResult{Action: "retry", Outcome: "accepted", Message: "Trying the local scheduler service again."}
}

// Quit performs the native application exit action.
func (a *App) Quit() ActionResult {
	if a.native == nil || a.ctx == nil {
		return ActionResult{Action: "quit", Outcome: "unavailable", Message: "Native exit is unavailable."}
	}
	a.native.Quit(a.ctx)
	return ActionResult{Action: "quit", Outcome: "accepted", Message: "Closing go-schedule."}
}
