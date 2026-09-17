package mcpsession

import (
	"errors"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

type fakeStore struct {
	actor   domain.Actor
	revoked string
}

func (s *fakeStore) CreateActor(kind domain.ActorKind, name string, capability domain.Capability, expires *time.Time) (domain.Actor, error) {
	s.actor = domain.Actor{ID: "actor-1", Kind: kind, DisplayName: name, Capability: capability, State: domain.ActorStateActive, ExpiresAt: expires}
	return s.actor, nil
}

func (s *fakeStore) RevokeActor(id string) (domain.Actor, error) {
	if id != s.actor.ID {
		return domain.Actor{}, store.ErrNotFound
	}
	s.revoked = id
	s.actor.State = domain.ActorStateRevoked
	return s.actor, nil
}

func TestRegistryCreatesResolvesAndRetainsAttributionAfterRevocation(t *testing.T) {
	st := &fakeStore{}
	registry := New(st)
	now := time.Date(2026, 9, 17, 12, 0, 0, 0, time.UTC)
	registry.now = func() time.Time { return now }
	session, secret, err := registry.Create("Codex", domain.CapabilityOperate)
	if err != nil || secret == "" || session.ActorID != "actor-1" || session.CreatedAt != now || session.ExpiresAt != now.Add(Lifetime) {
		t.Fatalf("session=%+v secret=%t err=%v", session, secret != "", err)
	}
	if actorID, err := registry.Resolve(secret); err != nil || actorID != "actor-1" {
		t.Fatalf("resolve=%q err=%v", actorID, err)
	}
	if _, err := registry.Revoke(session.ID); err != nil || st.revoked != "actor-1" {
		t.Fatalf("revoke=%q err=%v", st.revoked, err)
	}
	if actorID, err := registry.Resolve(secret); err != nil || actorID != "actor-1" {
		t.Fatalf("revoked attribution=%q err=%v", actorID, err)
	}
	if _, err := registry.Resolve("wrong"); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("wrong secret err=%v", err)
	}
}

func TestRegistryRejectsManageAndUnknownSessions(t *testing.T) {
	registry := New(&fakeStore{})
	if _, _, err := registry.Create("Codex", domain.CapabilityManage); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("manage err=%v", err)
	}
	if _, err := registry.Revoke("missing"); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("missing err=%v", err)
	}
}
