package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestRequestStoreOperationsCancelQueuedAcquisition(t *testing.T) {
	st := openMem(t)
	actor, _ := st.LocalActor()
	identity, _ := st.DaemonIdentity()
	group := domain.Group{Name: "Preserved group", Enabled: true}
	if err := st.CreateGroup(&group); err != nil {
		t.Fatal(err)
	}
	event := domain.AuditEvent{ID: uuid.NewString(), ActorID: actor.ID, DaemonID: identity.InstallationID, Operation: "groups.delete", TargetKind: "group", TargetID: group.ID, CorrelationID: uuid.NewString(), OccurredAt: time.Now()}
	if err := st.BeginAudit(event); err != nil {
		t.Fatal(err)
	}
	operations := []struct {
		name string
		run  func(context.Context) error
	}{
		{"local_actor", func(c context.Context) error { _, e := st.LocalActorContext(c); return e }},
		{"actor", func(c context.Context) error { _, e := st.GetActorContext(c, actor.ID); return e }},
		{"identity", func(c context.Context) error { _, e := st.DaemonIdentityContext(c); return e }},
		{"tasks", func(c context.Context) error { _, e := st.ListTasksContext(c, "", ""); return e }},
		{"observations", func(c context.Context) error {
			_, e := st.ListTaskObservationsContext(c, "", "", false, 0, 10, 1024)
			return e
		}},
		{"schedule", func(c context.Context) error { _, e := st.GetScheduleContext(c, "owned"); return e }},
		{"groups", func(c context.Context) error { _, e := st.ListGroupsContext(c); return e }},
		{"credentials", func(c context.Context) error { _, e := st.ListCredentialsContext(c); return e }},
		{"delete_group", func(c context.Context) error { return st.DeleteGroupContext(c, group.ID) }},
		{"intent", func(c context.Context) error { return st.BeginAuditContext(c, event) }},
		{"denial", func(c context.Context) error { return st.RecordDeniedAuditContext(c, event) }},
		{"completion", func(c context.Context) error {
			return st.CompleteAuditContext(c, event.ID, domain.AuditResultSucceeded)
		}},
	}
	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			held, err := st.db.Conn(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			defer held.Close()
			baseline := st.db.Stats().WaitCount
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			done := make(chan error, 1)
			go func() { done <- operation.run(ctx) }()
			timeout := time.NewTimer(5 * time.Second)
			defer timeout.Stop()
			tick := time.NewTicker(time.Millisecond)
			defer tick.Stop()
			for st.db.Stats().WaitCount == baseline {
				select {
				case <-tick.C:
				case <-timeout.C:
					t.Fatal("operation did not enter pool wait")
				}
			}
			cancel()
			select {
			case err := <-done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("lost cancellation: %v", err)
				}
			case <-time.After(time.Second):
				t.Fatal("canceled store operation did not return before pool release")
			}
		})
	}
	if _, err := st.GetGroup(group.ID); err != nil {
		t.Fatalf("canceled cleanup mutated group: %v", err)
	}
	events, err := st.ListAudit(domain.AuditQuery{Limit: 10})
	if err != nil || len(events) != 1 || events[0].Result != domain.AuditResultUncertain {
		t.Fatalf("lost uncertain intent: %+v %v", events, err)
	}
}
