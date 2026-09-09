package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type accessClientFake struct {
	actors []domain.Actor
	events []domain.AuditEvent
	actor  domain.Actor
	export []byte
}

func (f *accessClientFake) ListActors(context.Context) ([]domain.Actor, error) { return f.actors, nil }
func (f *accessClientFake) CreateActor(_ context.Context, request server.ActorCreateRequest) (domain.Actor, error) {
	f.actor = domain.Actor{ID: "actor-1", DisplayName: request.DisplayName, Kind: request.Kind, Capability: request.Capability, State: domain.ActorStateActive}
	return f.actor, nil
}
func (f *accessClientFake) UpdateActor(_ context.Context, id string, _ server.ActorUpdateRequest) (domain.Actor, error) {
	f.actor.ID = id
	return f.actor, nil
}
func (f *accessClientFake) RevokeActor(_ context.Context, id string) (domain.Actor, error) {
	f.actor.ID, f.actor.State = id, domain.ActorStateRevoked
	return f.actor, nil
}
func (f *accessClientFake) ListAudit(context.Context, domain.AuditQuery) ([]domain.AuditEvent, error) {
	return f.events, nil
}
func (f *accessClientFake) ExportAudit(context.Context, domain.AuditQuery) ([]byte, error) {
	return f.export, nil
}

func TestActorCommandsCoverLifecycle(t *testing.T) {
	fake := &accessClientFake{actors: []domain.Actor{{ID: "actor-1", DisplayName: "CLI", Kind: domain.ActorKindCLI, Capability: domain.CapabilityObserve, State: domain.ActorStateActive}}}
	for _, test := range []struct {
		args []string
		want string
	}{
		{[]string{"list"}, "actor-1"},
		{[]string{"create", "Desktop", "--kind", "desktop", "--capability", "operate"}, "Created actor Desktop"},
		{[]string{"update", "actor-1", "--name", "Release"}, "Updated actor"},
		{[]string{"revoke", "actor-1"}, "Revoked actor"},
	} {
		var output bytes.Buffer
		command := newActorCmdWithClient(fake)
		command.SetOut(&output)
		command.SetArgs(test.args)
		if err := command.Execute(); err != nil {
			t.Fatalf("%v: %v", test.args, err)
		}
		if !strings.Contains(output.String(), test.want) {
			t.Fatalf("%v output=%q", test.args, output.String())
		}
	}
}

func TestAuditCommandsListAndExport(t *testing.T) {
	fake := &accessClientFake{events: []domain.AuditEvent{{ID: "event-1", Operation: "tasks.create", Result: domain.AuditResultSucceeded, OccurredAt: time.Unix(0, 0).UTC()}}, export: []byte("{\"id\":\"event-1\"}\n")}
	for _, test := range []struct {
		args []string
		want string
	}{
		{[]string{"list", "--result", "succeeded"}, "tasks.create"},
		{[]string{"export"}, "event-1"},
	} {
		var output bytes.Buffer
		command := newAuditCmdWithClient(fake)
		command.SetOut(&output)
		command.SetArgs(test.args)
		if err := command.Execute(); err != nil {
			t.Fatalf("%v: %v", test.args, err)
		}
		if !strings.Contains(output.String(), test.want) {
			t.Fatalf("%v output=%q", test.args, output.String())
		}
	}
}
