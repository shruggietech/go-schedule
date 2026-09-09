package authorization

import (
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestCatalogIsUniqueValidAndFailClosed(t *testing.T) {
	seenRoutes, seenIDs := map[string]bool{}, map[string]bool{}
	for _, operation := range Catalog() {
		route := operation.Method + " " + operation.Pattern
		if seenRoutes[route] || seenIDs[operation.ID] {
			t.Fatalf("duplicate catalog entry %s %s", route, operation.ID)
		}
		seenRoutes[route], seenIDs[operation.ID] = true, true
		if !operation.Capability.Valid() || operation.TargetKind == "" || (operation.Audit != AuditNone && operation.Audit != AuditMutation && operation.Audit != AuditPrivilegedRead) {
			t.Fatalf("invalid operation: %+v", operation)
		}
	}
	if _, ok := Lookup("POST", "/v1/unknown"); ok {
		t.Fatal("unknown route must fail closed")
	}
	if _, ok := LookupID("unknown.operation"); ok {
		t.Fatal("unknown operation ID must fail closed")
	}
}

func TestStaticRouteWinsOverParameterizedRoute(t *testing.T) {
	operation, ok := Lookup("GET", "/v1/runs/active")
	if !ok || operation.ID != "runs.active.list" || operation.Capability != domain.CapabilityObserve {
		t.Fatalf("operation=%+v ok=%v", operation, ok)
	}
}

func TestAllowsRejectsUnknownExpiredAndRevokedActors(t *testing.T) {
	now := time.Now()
	actor := domain.Actor{Kind: domain.ActorKindCLI, Capability: domain.CapabilityOperate, State: domain.ActorStateActive}
	if !Allows(actor, "tasks.run_now", now) || Allows(actor, "tasks.create", now) || Allows(actor, "unknown.operation", now) {
		t.Fatal("capability mapping is not fail closed")
	}
	past := now.Add(-time.Second)
	actor.ExpiresAt = &past
	if Allows(actor, "health.read", now) {
		t.Fatal("expired actor was authorized")
	}
	actor.ExpiresAt, actor.State = nil, domain.ActorStateRevoked
	if Allows(actor, "health.read", now) {
		t.Fatal("revoked actor was authorized")
	}
}

func TestTargetIDIsBounded(t *testing.T) {
	operation, _ := Lookup("GET", "/v1/tasks/task-1")
	if got := TargetID(operation, "/v1/tasks/task-1"); got != "task-1" {
		t.Fatalf("target=%q", got)
	}
	if got := TargetID(operation, "/v1/tasks/"+strings.Repeat("x", 257)); got != "" {
		t.Fatalf("oversized target=%q", got)
	}
}
