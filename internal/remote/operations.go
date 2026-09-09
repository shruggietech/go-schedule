package remote

//go:generate go tool oapi-codegen --config ../../api/openapi/remote-v1.cfg.yaml ../../api/openapi/remote-v1.yaml

import (
	"net/http"
	"strings"

	"github.com/shruggietech/go-schedule/internal/authorization"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type RetryClass string

const (
	RetrySafe      RetryClass = "safe"
	RetryReconnect RetryClass = "reconnect"
	RetryNever     RetryClass = "never"
)

type Operation struct {
	Method           string
	RemotePath       string
	LocalPath        string
	ID               string
	Retry            RetryClass
	BodyLimit        int64
	Public           bool
	Capability       domain.Capability
	TargetKind       string
	Audit            authorization.AuditClass
	SecretExclusions []string
}

var operations = []Operation{
	{Method: http.MethodGet, RemotePath: "/api/v1/health", LocalPath: "/v1/health", ID: "health.read", Retry: RetrySafe, Public: true},
	{Method: http.MethodGet, RemotePath: "/api/v1/manifest", LocalPath: "/v1/manifest", ID: "manifest.read", Retry: RetrySafe, Public: true},
	{Method: http.MethodPost, RemotePath: "/api/v1/enroll", ID: "enrollment.exchange", Retry: RetryNever, BodyLimit: 4096, Public: true},
	{Method: http.MethodGet, RemotePath: "/api/v1/tasks", LocalPath: "/v1/tasks", ID: "tasks.list", Retry: RetrySafe},
	{Method: http.MethodPost, RemotePath: "/api/v1/tasks", LocalPath: "/v1/tasks", ID: "tasks.create", Retry: RetryNever, BodyLimit: 256 << 10},
	{Method: http.MethodGet, RemotePath: "/api/v1/tasks/{id}", LocalPath: "/v1/tasks/{id}", ID: "tasks.read", Retry: RetrySafe},
	{Method: http.MethodPatch, RemotePath: "/api/v1/tasks/{id}", LocalPath: "/v1/tasks/{id}", ID: "tasks.update", Retry: RetryNever, BodyLimit: 256 << 10},
	{Method: http.MethodDelete, RemotePath: "/api/v1/tasks/{id}", LocalPath: "/v1/tasks/{id}", ID: "tasks.delete", Retry: RetryNever},
	{Method: http.MethodPost, RemotePath: "/api/v1/tasks/{id}/enable", LocalPath: "/v1/tasks/{id}/enable", ID: "tasks.enable", Retry: RetryNever},
	{Method: http.MethodPost, RemotePath: "/api/v1/tasks/{id}/disable", LocalPath: "/v1/tasks/{id}/disable", ID: "tasks.disable", Retry: RetryNever},
	{Method: http.MethodPost, RemotePath: "/api/v1/tasks/{id}/run-now", LocalPath: "/v1/tasks/{id}/run-now", ID: "tasks.run_now", Retry: RetryNever},
	{Method: http.MethodGet, RemotePath: "/api/v1/groups", LocalPath: "/v1/groups", ID: "groups.list", Retry: RetrySafe},
	{Method: http.MethodGet, RemotePath: "/api/v1/groups/{id}", LocalPath: "/v1/groups/{id}", ID: "groups.read", Retry: RetrySafe},
	{Method: http.MethodGet, RemotePath: "/api/v1/runs", LocalPath: "/v1/runs", ID: "runs.list", Retry: RetrySafe},
	{Method: http.MethodGet, RemotePath: "/api/v1/runs/active", LocalPath: "/v1/runs/active", ID: "runs.active.list", Retry: RetrySafe},
	{Method: http.MethodGet, RemotePath: "/api/v1/runs/{id}", LocalPath: "/v1/runs/{id}", ID: "runs.read", Retry: RetrySafe},
	{Method: http.MethodGet, RemotePath: "/api/v1/alerts", LocalPath: "/v1/alerts", ID: "alerts.list", Retry: RetrySafe},
	{Method: http.MethodPost, RemotePath: "/api/v1/alerts/{id}/ack", LocalPath: "/v1/alerts/{id}/ack", ID: "alerts.acknowledge", Retry: RetryNever},
	{Method: http.MethodGet, RemotePath: "/api/v1/calendar", LocalPath: "/v1/calendar", ID: "calendar.read", Retry: RetrySafe},
	{Method: http.MethodGet, RemotePath: "/api/v1/events", LocalPath: "/v1/events", ID: "events.stream", Retry: RetryReconnect},
	{Method: http.MethodGet, RemotePath: "/api/v1/access/actors", LocalPath: "/v1/access/actors", ID: "actors.list", Retry: RetrySafe},
	{Method: http.MethodGet, RemotePath: "/api/v1/audit", LocalPath: "/v1/audit", ID: "audit.list", Retry: RetrySafe},
}

func init() {
	for i := range operations {
		if local, ok := authorization.LookupID(operations[i].ID); ok {
			operations[i].Capability = local.Capability
			operations[i].TargetKind = local.TargetKind
			operations[i].Audit = local.Audit
		} else {
			operations[i].TargetKind = "pairing"
			operations[i].Audit = authorization.AuditNone
		}
		operations[i].SecretExclusions = []string{"authorization", "credential", "digest", "phrase", "verifier", "task-input"}
	}
}

func Operations() []Operation {
	result := append([]Operation(nil), operations...)
	for i := range result {
		result[i].SecretExclusions = append([]string(nil), result[i].SecretExclusions...)
	}
	return result
}

func lookup(method, path string) (Operation, bool) {
	for _, operation := range operations {
		if operation.Method == method && pathMatch(operation.RemotePath, path) {
			return operation, true
		}
	}
	return Operation{}, false
}

func pathMatch(pattern, path string) bool {
	a, b := strings.Split(strings.Trim(pattern, "/"), "/"), strings.Split(strings.Trim(path, "/"), "/")
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if strings.HasPrefix(a[i], "{") && strings.HasSuffix(a[i], "}") {
			continue
		}
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func localPath(operation Operation, path string) string {
	if !strings.Contains(operation.LocalPath, "{id}") {
		return operation.LocalPath
	}
	remoteParts, actualParts := strings.Split(strings.Trim(operation.RemotePath, "/"), "/"), strings.Split(strings.Trim(path, "/"), "/")
	for i := range remoteParts {
		if remoteParts[i] == "{id}" {
			return strings.Replace(operation.LocalPath, "{id}", actualParts[i], 1)
		}
	}
	return operation.LocalPath
}
