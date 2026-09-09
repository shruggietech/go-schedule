// Package authorization defines the transport-independent operation catalog.
package authorization

import (
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

type AuditClass string

const (
	AuditNone           AuditClass = "none"
	AuditPrivilegedRead AuditClass = "privileged_read"
	AuditMutation       AuditClass = "mutation"
)

type Operation struct {
	Method     string
	Pattern    string
	ID         string
	Capability domain.Capability
	TargetKind string
	Audit      AuditClass
}

func op(method, pattern, id string, capability domain.Capability, target string, audit AuditClass) Operation {
	return Operation{Method: method, Pattern: pattern, ID: id, Capability: capability, TargetKind: target, Audit: audit}
}

var catalog = []Operation{
	op("GET", "/v1/health", "health.read", domain.CapabilityObserve, "daemon", AuditNone),
	op("GET", "/v1/manifest", "manifest.read", domain.CapabilityObserve, "daemon", AuditNone),
	op("PATCH", "/v1/manifest", "manifest.rename", domain.CapabilityEnroll, "daemon", AuditMutation),
	op("POST", "/v1/manifest/reset", "manifest.identity.reset", domain.CapabilityEnroll, "daemon", AuditMutation),
	op("GET", "/v1/runtime-info", "runtime.read", domain.CapabilityManage, "daemon", AuditPrivilegedRead),
	op("GET", "/v1/mcp/http", "mcp.http.read", domain.CapabilityObserve, "mcp_listener", AuditNone),
	op("POST", "/v1/mcp/http/enable", "mcp.http.enable", domain.CapabilityEnroll, "mcp_listener", AuditMutation),
	op("POST", "/v1/mcp/http/rotate", "mcp.http.rotate", domain.CapabilityEnroll, "mcp_listener", AuditMutation),
	op("POST", "/v1/mcp/http/disable", "mcp.http.disable", domain.CapabilityEnroll, "mcp_listener", AuditMutation),
	op("GET", "/v1/tasks", "tasks.list", domain.CapabilityObserve, "task", AuditNone),
	op("POST", "/v1/tasks", "tasks.create", domain.CapabilityManage, "task", AuditMutation),
	op("GET", "/v1/tasks/{id}", "tasks.read", domain.CapabilityObserve, "task", AuditNone),
	op("PATCH", "/v1/tasks/{id}", "tasks.update", domain.CapabilityManage, "task", AuditMutation),
	op("DELETE", "/v1/tasks/{id}", "tasks.delete", domain.CapabilityManage, "task", AuditMutation),
	op("POST", "/v1/tasks/{id}/enable", "tasks.enable", domain.CapabilityManage, "task", AuditMutation),
	op("POST", "/v1/tasks/{id}/disable", "tasks.disable", domain.CapabilityManage, "task", AuditMutation),
	op("POST", "/v1/tasks/{id}/run-now", "tasks.run_now", domain.CapabilityOperate, "task", AuditMutation),
	op("GET", "/v1/triggers", "triggers.list", domain.CapabilityObserve, "trigger", AuditNone),
	op("POST", "/v1/triggers", "triggers.create", domain.CapabilityManage, "trigger", AuditMutation),
	op("POST", "/v1/triggers/fire", "triggers.fire", domain.CapabilityOperate, "trigger", AuditMutation),
	op("GET", "/v1/triggers/{id}", "triggers.read", domain.CapabilityObserve, "trigger", AuditNone),
	op("PATCH", "/v1/triggers/{id}", "triggers.update", domain.CapabilityManage, "trigger", AuditMutation),
	op("DELETE", "/v1/triggers/{id}", "triggers.delete", domain.CapabilityManage, "trigger", AuditMutation),
	op("POST", "/v1/triggers/{id}/enable", "triggers.enable", domain.CapabilityManage, "trigger", AuditMutation),
	op("POST", "/v1/triggers/{id}/disable", "triggers.disable", domain.CapabilityManage, "trigger", AuditMutation),
	op("POST", "/v1/triggers/{id}/rotate", "triggers.rotate", domain.CapabilityEnroll, "trigger_secret", AuditMutation),
	op("POST", "/v1/triggers/{id}/reveal", "triggers.reveal", domain.CapabilityEnroll, "trigger_secret", AuditPrivilegedRead),
	op("GET", "/v1/filesystem-watchers", "watchers.list", domain.CapabilityObserve, "filesystem_watcher", AuditNone),
	op("POST", "/v1/filesystem-watchers", "watchers.create", domain.CapabilityManage, "filesystem_watcher", AuditMutation),
	op("GET", "/v1/filesystem-watchers/{id}", "watchers.read", domain.CapabilityObserve, "filesystem_watcher", AuditNone),
	op("PATCH", "/v1/filesystem-watchers/{id}", "watchers.update", domain.CapabilityManage, "filesystem_watcher", AuditMutation),
	op("DELETE", "/v1/filesystem-watchers/{id}", "watchers.delete", domain.CapabilityManage, "filesystem_watcher", AuditMutation),
	op("POST", "/v1/filesystem-watchers/{id}/enable", "watchers.enable", domain.CapabilityManage, "filesystem_watcher", AuditMutation),
	op("POST", "/v1/filesystem-watchers/{id}/disable", "watchers.disable", domain.CapabilityManage, "filesystem_watcher", AuditMutation),
	op("GET", "/v1/trigger-sets", "trigger_sets.list", domain.CapabilityObserve, "trigger_set", AuditNone),
	op("POST", "/v1/trigger-sets", "trigger_sets.create", domain.CapabilityManage, "trigger_set", AuditMutation),
	op("GET", "/v1/trigger-sets/{id}", "trigger_sets.read", domain.CapabilityObserve, "trigger_set", AuditNone),
	op("PATCH", "/v1/trigger-sets/{id}", "trigger_sets.update", domain.CapabilityManage, "trigger_set", AuditMutation),
	op("DELETE", "/v1/trigger-sets/{id}", "trigger_sets.delete", domain.CapabilityManage, "trigger_set", AuditMutation),
	op("POST", "/v1/trigger-sets/{id}/enable", "trigger_sets.enable", domain.CapabilityManage, "trigger_set", AuditMutation),
	op("POST", "/v1/trigger-sets/{id}/disable", "trigger_sets.disable", domain.CapabilityManage, "trigger_set", AuditMutation),
	op("POST", "/v1/trigger-sets/{id}/rotate", "trigger_sets.rotate", domain.CapabilityEnroll, "trigger_secret", AuditMutation),
	op("POST", "/v1/trigger-sets/{id}/reveal", "trigger_sets.reveal", domain.CapabilityEnroll, "trigger_secret", AuditPrivilegedRead),
	op("GET", "/v1/chains", "chains.list", domain.CapabilityObserve, "chain", AuditNone),
	op("POST", "/v1/chains", "chains.create", domain.CapabilityManage, "chain", AuditMutation),
	op("GET", "/v1/chains/{id}", "chains.read", domain.CapabilityObserve, "chain", AuditNone),
	op("PATCH", "/v1/chains/{id}", "chains.update", domain.CapabilityManage, "chain", AuditMutation),
	op("DELETE", "/v1/chains/{id}", "chains.delete", domain.CapabilityManage, "chain", AuditMutation),
	op("GET", "/v1/groups", "groups.list", domain.CapabilityObserve, "group", AuditNone),
	op("POST", "/v1/groups", "groups.create", domain.CapabilityManage, "group", AuditMutation),
	op("GET", "/v1/groups/{id}", "groups.read", domain.CapabilityObserve, "group", AuditNone),
	op("PATCH", "/v1/groups/{id}", "groups.update", domain.CapabilityManage, "group", AuditMutation),
	op("DELETE", "/v1/groups/{id}", "groups.delete", domain.CapabilityManage, "group", AuditMutation),
	op("POST", "/v1/groups/{id}/enable", "groups.enable", domain.CapabilityManage, "group", AuditMutation),
	op("POST", "/v1/groups/{id}/disable", "groups.disable", domain.CapabilityManage, "group", AuditMutation),
	op("POST", "/v1/schedules/preview", "schedules.preview", domain.CapabilityObserve, "schedule", AuditNone),
	op("GET", "/v1/runs", "runs.list", domain.CapabilityObserve, "run", AuditNone),
	op("GET", "/v1/runs/active", "runs.active.list", domain.CapabilityObserve, "run", AuditNone),
	op("GET", "/v1/runs/{id}", "runs.read", domain.CapabilityObserve, "run", AuditNone),
	op("GET", "/v1/notification-channels", "notifications.channels.list", domain.CapabilityObserve, "notification_channel", AuditNone),
	op("POST", "/v1/notification-channels", "notifications.channels.create", domain.CapabilityManage, "notification_channel", AuditMutation),
	op("GET", "/v1/notification-channels/{id}", "notifications.channels.read", domain.CapabilityObserve, "notification_channel", AuditNone),
	op("PATCH", "/v1/notification-channels/{id}", "notifications.channels.update", domain.CapabilityManage, "notification_channel", AuditMutation),
	op("DELETE", "/v1/notification-channels/{id}", "notifications.channels.delete", domain.CapabilityManage, "notification_channel", AuditMutation),
	op("POST", "/v1/notification-channels/{id}/enable", "notifications.channels.enable", domain.CapabilityManage, "notification_channel", AuditMutation),
	op("POST", "/v1/notification-channels/{id}/disable", "notifications.channels.disable", domain.CapabilityManage, "notification_channel", AuditMutation),
	op("POST", "/v1/notification-channels/{id}/rotate", "notifications.channels.rotate", domain.CapabilityEnroll, "notification_secret", AuditMutation),
	op("POST", "/v1/notification-channels/{id}/test", "notifications.channels.test", domain.CapabilityOperate, "notification_channel", AuditMutation),
	op("GET", "/v1/notification-deliveries", "notifications.deliveries.list", domain.CapabilityObserve, "notification_delivery", AuditNone),
	op("GET", "/v1/tasks/{id}/notifications", "tasks.notifications.read", domain.CapabilityObserve, "notification_assignment", AuditNone),
	op("PUT", "/v1/tasks/{id}/notifications", "tasks.notifications.update", domain.CapabilityManage, "notification_assignment", AuditMutation),
	op("GET", "/v1/tasks/{id}/notifications/effective", "tasks.notifications.effective.read", domain.CapabilityObserve, "notification_assignment", AuditNone),
	op("GET", "/v1/groups/{id}/notifications", "groups.notifications.read", domain.CapabilityObserve, "notification_assignment", AuditNone),
	op("PUT", "/v1/groups/{id}/notifications", "groups.notifications.update", domain.CapabilityManage, "notification_assignment", AuditMutation),
	op("GET", "/v1/alerts", "alerts.list", domain.CapabilityObserve, "alert", AuditNone),
	op("POST", "/v1/alerts/{id}/ack", "alerts.acknowledge", domain.CapabilityOperate, "alert", AuditMutation),
	op("GET", "/v1/logs", "logs.list", domain.CapabilityObserve, "log", AuditNone),
	op("GET", "/v1/calendar", "calendar.read", domain.CapabilityObserve, "calendar", AuditNone),
	op("GET", "/v1/events", "events.stream", domain.CapabilityObserve, "event_stream", AuditNone),
	op("GET", "/v1/access/actors", "actors.list", domain.CapabilityEnroll, "actor", AuditPrivilegedRead),
	op("POST", "/v1/access/actors", "actors.create", domain.CapabilityEnroll, "actor", AuditMutation),
	op("PATCH", "/v1/access/actors/{id}", "actors.update", domain.CapabilityEnroll, "actor", AuditMutation),
	op("POST", "/v1/access/actors/{id}/revoke", "actors.revoke", domain.CapabilityEnroll, "actor", AuditMutation),
	op("GET", "/v1/access/pairings", "pairings.list", domain.CapabilityEnroll, "pairing", AuditPrivilegedRead),
	op("POST", "/v1/access/pairings", "pairings.create", domain.CapabilityEnroll, "pairing", AuditMutation),
	op("POST", "/v1/access/pairings/{id}/cancel", "pairings.cancel", domain.CapabilityEnroll, "pairing", AuditMutation),
	op("GET", "/v1/access/credentials", "credentials.list", domain.CapabilityEnroll, "credential", AuditPrivilegedRead),
	op("POST", "/v1/access/credentials/{id}/rotate", "credentials.rotate", domain.CapabilityEnroll, "credential", AuditMutation),
	op("POST", "/v1/access/credentials/{id}/revoke", "credentials.revoke", domain.CapabilityEnroll, "credential", AuditMutation),
	op("GET", "/v1/audit", "audit.list", domain.CapabilityEnroll, "audit", AuditPrivilegedRead),
	op("GET", "/v1/audit/export", "audit.export", domain.CapabilityEnroll, "audit", AuditPrivilegedRead),
}

func Catalog() []Operation { return append([]Operation(nil), catalog...) }

func Lookup(method, path string) (Operation, bool) {
	for _, operation := range catalog {
		if operation.Method == method && operation.Pattern == path {
			return operation, true
		}
	}
	for _, operation := range catalog {
		if operation.Method == method && match(operation.Pattern, path) {
			return operation, true
		}
	}
	return Operation{}, false
}

func LookupID(id string) (Operation, bool) {
	for _, operation := range catalog {
		if operation.ID == id {
			return operation, true
		}
	}
	return Operation{}, false
}

// Allows applies the shared fail-closed actor and operation checks.
func Allows(actor domain.Actor, operationID string, now time.Time) bool {
	operation, ok := LookupID(operationID)
	return ok && actor.Kind.Valid() && actor.State.Valid() && actor.Capability.Valid() && actor.ActiveAt(now) && actor.Capability.Allows(operation.Capability)
}

func match(pattern, path string) bool {
	left, right := strings.Split(strings.Trim(pattern, "/"), "/"), strings.Split(strings.Trim(path, "/"), "/")
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if strings.HasPrefix(left[i], "{") && strings.HasSuffix(left[i], "}") {
			continue
		}
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func TargetID(operation Operation, path string) string {
	patternParts, pathParts := strings.Split(strings.Trim(operation.Pattern, "/"), "/"), strings.Split(strings.Trim(path, "/"), "/")
	for i := range patternParts {
		if patternParts[i] == "{id}" && i < len(pathParts) {
			if len(pathParts[i]) <= 256 {
				return pathParts[i]
			}
			return ""
		}
	}
	return ""
}
