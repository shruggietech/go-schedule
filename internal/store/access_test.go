package store

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestActorLifecycleAndBuiltInProtection(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	local, err := st.LocalActor()
	if err != nil {
		t.Fatal(err)
	}
	if !local.Builtin || local.Kind != domain.ActorKindLocalOS || local.Capability != domain.CapabilityEnroll || !local.ActiveAt(time.Now()) {
		t.Fatalf("local actor=%+v", local)
	}
	capability := domain.CapabilityManage
	actor, err := st.CreateActor(domain.ActorKindCLI, "  Operator CLI  ", domain.CapabilityObserve, nil)
	if err != nil {
		t.Fatal(err)
	}
	updated, err := st.UpdateActor(actor.ID, ActorUpdate{Capability: &capability})
	if err != nil || updated.Capability != capability || updated.DisplayName != "Operator CLI" {
		t.Fatalf("updated=%+v err=%v", updated, err)
	}
	revoked, err := st.RevokeActor(actor.ID)
	if err != nil || revoked.State != domain.ActorStateRevoked || revoked.ActiveAt(time.Now()) {
		t.Fatalf("revoked=%+v err=%v", revoked, err)
	}
	active := domain.ActorStateActive
	if _, err := st.UpdateActor(actor.ID, ActorUpdate{State: &active}); !errors.Is(err, domain.ErrInvalidActor) {
		t.Fatalf("revoked reactivation err=%v", err)
	}
	if _, err := st.RevokeActor(local.ID); !errors.Is(err, ErrBuiltinActor) {
		t.Fatalf("built-in revoke err=%v", err)
	}
	actors, err := st.ListActors()
	if err != nil || len(actors) != 2 || actors[0].ID != local.ID {
		t.Fatalf("actors=%+v err=%v", actors, err)
	}
}

func TestAuditLifecycleFilteringAndAgeRetention(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	local, _ := st.LocalActor()
	identity, _ := st.DaemonIdentity()
	old := domain.AuditEvent{ID: uuid.NewString(), ActorID: local.ID, DaemonID: identity.InstallationID, Operation: "tasks.create", TargetKind: "task", CorrelationID: uuid.NewString(), OccurredAt: time.Now().Add(-AuditRetentionAge - time.Hour)}
	if err := st.BeginAudit(old); err != nil {
		t.Fatal(err)
	}
	current := domain.AuditEvent{ID: uuid.NewString(), ActorID: local.ID, DaemonID: identity.InstallationID, Operation: "tasks.update", TargetKind: "task", TargetID: "task-1", CorrelationID: uuid.NewString(), OccurredAt: time.Now()}
	if err := st.BeginAudit(current); err != nil {
		t.Fatal(err)
	}
	if err := st.CompleteAudit(current.ID, domain.AuditResultSucceeded); err != nil {
		t.Fatal(err)
	}
	denied := domain.AuditEvent{ID: uuid.NewString(), ActorID: local.ID, DaemonID: identity.InstallationID, Operation: "actors.create", TargetKind: "actor", CorrelationID: uuid.NewString(), OccurredAt: time.Now().Add(time.Second)}
	if err := st.RecordDeniedAudit(denied); err != nil {
		t.Fatal(err)
	}
	events, err := st.ListAudit(domain.AuditQuery{ActorID: local.ID, Result: domain.AuditResultSucceeded, Limit: 10})
	if err != nil || len(events) != 1 || events[0].ID != current.ID || events[0].CompletedAt == nil {
		t.Fatalf("events=%+v err=%v", events, err)
	}
	all, err := st.ListAudit(domain.AuditQuery{Limit: 10})
	if err != nil || len(all) != 2 || all[0].ID != current.ID || all[1].Result != domain.AuditResultDenied {
		t.Fatalf("all=%+v err=%v", all, err)
	}
}

func TestAuditCountRetentionKeepsNewestTenThousand(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	identity, _ := st.DaemonIdentity()
	now := time.Now().UTC()
	if _, err := st.db.Exec(`WITH RECURSIVE n(x) AS (VALUES(1) UNION ALL SELECT x+1 FROM n WHERE x<=10000) INSERT INTO audit_events(id,daemon_id,operation,target_kind,result,correlation_id,occurred_at) SELECT printf('old-%05d',x),?,'tasks.create','task','denied',printf('correlation-%05d',x),? FROM n`, identity.InstallationID, fmtTime(now.Add(-time.Minute))); err != nil {
		t.Fatal(err)
	}
	event := domain.AuditEvent{ID: uuid.NewString(), DaemonID: identity.InstallationID, Operation: "tasks.create", TargetKind: "task", CorrelationID: uuid.NewString(), OccurredAt: now}
	if err := st.BeginAudit(event); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM audit_events`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != AuditRetentionCount {
		t.Fatalf("audit count=%d, want %d", count, AuditRetentionCount)
	}
}

func TestActorValidationUpdatesAndAuditQueryBoundaries(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if _, err := st.GetActor(uuid.NewString()); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing actor err=%v", err)
	}
	for _, test := range []struct {
		kind       domain.ActorKind
		name       string
		capability domain.Capability
		expires    *time.Time
	}{
		{domain.ActorKindLocalOS, "duplicate local", domain.CapabilityEnroll, nil},
		{domain.ActorKindCLI, "", domain.CapabilityObserve, nil},
		{domain.ActorKindCLI, "client", domain.Capability("owner"), nil},
		{domain.ActorKindCLI, "client", domain.CapabilityObserve, timePointer(time.Now().Add(-time.Minute))},
	} {
		if _, err := st.CreateActor(test.kind, test.name, test.capability, test.expires); !errors.Is(err, domain.ErrInvalidActor) {
			t.Fatalf("invalid create err=%v", err)
		}
	}
	future := time.Now().Add(time.Hour).UTC()
	actor, err := st.CreateActor(domain.ActorKindDesktop, "Desktop", domain.CapabilityOperate, &future)
	if err != nil {
		t.Fatal(err)
	}
	name := "Renamed"
	state := domain.ActorStateExpired
	clear := (*time.Time)(nil)
	updated, err := st.UpdateActor(actor.ID, ActorUpdate{DisplayName: &name, State: &state, ExpiresAt: &clear})
	if err != nil || updated.DisplayName != name || updated.State != state || updated.ExpiresAt != nil {
		t.Fatalf("updated=%+v err=%v", updated, err)
	}
	invalidCapability := domain.Capability("owner")
	if _, err := st.UpdateActor(actor.ID, ActorUpdate{Capability: &invalidCapability}); !errors.Is(err, domain.ErrInvalidActor) {
		t.Fatalf("invalid capability err=%v", err)
	}
	invalidState := domain.ActorState("pending")
	if _, err := st.UpdateActor(actor.ID, ActorUpdate{State: &invalidState}); !errors.Is(err, domain.ErrInvalidActor) {
		t.Fatalf("invalid state err=%v", err)
	}
	identity, _ := st.DaemonIdentity()
	if err := st.BeginAudit(domain.AuditEvent{DaemonID: identity.InstallationID, Operation: "tasks.read", TargetKind: "task"}); err != nil {
		t.Fatal(err)
	}
	if err := st.RecordDeniedAudit(domain.AuditEvent{DaemonID: identity.InstallationID, Operation: "tasks.create", TargetKind: "task"}); err != nil {
		t.Fatal(err)
	}
	since, until := time.Now().Add(-time.Hour), time.Now().Add(time.Hour)
	if _, err := st.ListAudit(domain.AuditQuery{Operation: "tasks.read", Since: &since, Until: &until, Limit: 10}); err != nil {
		t.Fatal(err)
	}
	for _, query := range []domain.AuditQuery{{Limit: -1}, {Limit: 1001}, {Result: domain.AuditResult("ok")}, {Since: &until, Until: &since}} {
		if _, err := st.ListAudit(query); err == nil {
			t.Fatalf("accepted invalid query %+v", query)
		}
	}
	if err := st.CompleteAudit("missing", domain.AuditResultSucceeded); err == nil {
		t.Fatal("completed missing event")
	}
	if err := st.CompleteAudit("missing", domain.AuditResultDenied); err == nil {
		t.Fatal("accepted invalid completion result")
	}
	for _, event := range []domain.AuditEvent{
		{ID: "bad", DaemonID: uuid.NewString(), Operation: "x", TargetKind: "x", Result: domain.AuditResultDenied, CorrelationID: uuid.NewString(), OccurredAt: time.Now()},
		{ID: uuid.NewString(), DaemonID: "bad", Operation: "x", TargetKind: "x", Result: domain.AuditResultDenied, CorrelationID: uuid.NewString(), OccurredAt: time.Now()},
		{ID: uuid.NewString(), DaemonID: uuid.NewString(), ActorID: "bad", Operation: "x", TargetKind: "x", Result: domain.AuditResultDenied, CorrelationID: uuid.NewString(), OccurredAt: time.Now()},
		{ID: uuid.NewString(), DaemonID: uuid.NewString(), Operation: "", TargetKind: "", Result: domain.AuditResult("bad"), CorrelationID: uuid.NewString()},
	} {
		if err := validateAuditEvent(event); err == nil {
			t.Fatalf("accepted invalid event %+v", event)
		}
	}
}

func timePointer(value time.Time) *time.Time { return &value }
