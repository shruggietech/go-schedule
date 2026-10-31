package store

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func rotationFixture(t *testing.T) (*Store, domain.Actor, domain.ClientCredential, []byte) {
	t.Helper()
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	now := time.Now().UTC()
	session := domain.PairingSession{ID: uuid.NewString(), DisplayName: "Rotation client", Kind: domain.ActorKindCLI, Capability: domain.CapabilityManage, AttemptsRemaining: 5, State: domain.PairingActive, CreatedAt: now, ExpiresAt: now.Add(time.Minute)}
	verifier := make([]byte, 32)
	if err := st.CreatePairing(session, make([]byte, 16), verifier); err != nil {
		t.Fatal(err)
	}
	actor := domain.Actor{ID: uuid.NewString(), State: domain.ActorStateActive, CreatedAt: now, UpdatedAt: now}
	credential := domain.ClientCredential{ID: uuid.NewString(), Fingerprint: "original", State: domain.CredentialActive, CreatedAt: now, UpdatedAt: now}
	digest := bytes.Repeat([]byte{1}, 32)
	actor, credential, err = st.ExchangePairing(session.ID, verifier, actor, credential, digest, now)
	if err != nil {
		t.Fatal(err)
	}
	return st, actor, credential, digest
}

func TestCredentialRotationRollsBackWhenActorContextCannotBeRead(t *testing.T) {
	st, actor, credential, oldDigest := rotationFixture(t)
	if _, err := st.db.Exec("UPDATE actors SET created_at='invalid' WHERE id=?", actor.ID); err != nil {
		t.Fatal(err)
	}
	if _, _, err := st.RotateCredential(credential.ID, bytes.Repeat([]byte{2}, 32), "replacement", time.Now().UTC()); err == nil {
		t.Error("rotation succeeded with unreadable actor context")
	}
	assertRotationUnchanged(t, st, credential.ID, oldDigest)
	if _, err := st.db.Exec("UPDATE actors SET created_at=? WHERE id=?", fmtTime(actor.CreatedAt), actor.ID); err != nil {
		t.Fatal(err)
	}
	if _, _, err := st.AuthenticateCredential(oldDigest, time.Now().UTC()); err != nil {
		t.Fatal("failed context assembly invalidated existing token")
	}
}

func TestCredentialRotationRollsBackWhenReplacementResponseCannotBeRead(t *testing.T) {
	st, _, credential, oldDigest := rotationFixture(t)
	if _, err := st.db.Exec("CREATE TRIGGER corrupt_rotated_timestamp AFTER UPDATE OF digest ON client_credentials BEGIN UPDATE client_credentials SET updated_at='invalid' WHERE id=NEW.id; END"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := st.RotateCredential(credential.ID, bytes.Repeat([]byte{2}, 32), "replacement", time.Now().UTC()); err == nil {
		t.Error("rotation succeeded with unreadable replacement response")
	}
	assertRotationUnchanged(t, st, credential.ID, oldDigest)
	if _, _, err := st.AuthenticateCredential(oldDigest, time.Now().UTC()); err != nil {
		t.Fatal("response failure invalidated existing token")
	}
}

func assertRotationUnchanged(t *testing.T, st *Store, id string, expectedDigest []byte) {
	t.Helper()
	var digest []byte
	var fingerprint string
	if err := st.db.QueryRow("SELECT digest,fingerprint FROM client_credentials WHERE id=?", id).Scan(&digest, &fingerprint); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(digest, expectedDigest) || fingerprint != "original" {
		t.Error("failed rotation changed the stored credential")
	}
}

func TestCredentialRotationDoesNotChangeExistingEligibility(t *testing.T) {
	for _, item := range []struct {
		name                            string
		actorState                      domain.ActorState
		actorExpired, credentialExpired bool
		wantState                       domain.ActorState
	}{
		{name: "revoked actor", actorState: domain.ActorStateRevoked, wantState: domain.ActorStateRevoked},
		{name: "expired actor", actorState: domain.ActorStateActive, actorExpired: true, wantState: domain.ActorStateExpired},
		{name: "expired credential", actorState: domain.ActorStateActive, credentialExpired: true, wantState: domain.ActorStateActive},
	} {
		t.Run(item.name, func(t *testing.T) {
			st, originalActor, credential, _ := rotationFixture(t)
			if _, err := st.db.Exec("UPDATE actors SET state=? WHERE id=?", item.actorState, originalActor.ID); err != nil {
				t.Fatal(err)
			}
			expired := fmtTime(time.Now().UTC().Add(-time.Hour))
			if item.actorExpired {
				if _, err := st.db.Exec("UPDATE actors SET expires_at=? WHERE id=?", expired, originalActor.ID); err != nil {
					t.Fatal(err)
				}
			}
			if item.credentialExpired {
				if _, err := st.db.Exec("UPDATE client_credentials SET expires_at=? WHERE id=?", expired, credential.ID); err != nil {
					t.Fatal(err)
				}
			}
			replacement := bytes.Repeat([]byte{2}, 32)
			rotated, actor, err := st.RotateCredential(credential.ID, replacement, "replacement", time.Now().UTC())
			if err != nil {
				t.Fatalf("rotation changed eligibility: %v", err)
			}
			if actor.ID != originalActor.ID || actor.State != item.wantState || rotated.State != domain.CredentialActive || rotated.LastUsedAt != nil {
				t.Fatal("incorrect actor context or credential state")
			}
			if item.credentialExpired && (rotated.ExpiresAt == nil || !rotated.ExpiresAt.Before(time.Now())) {
				t.Fatal("rotation changed credential expiry")
			}
			var stored []byte
			if err := st.db.QueryRow("SELECT digest FROM client_credentials WHERE id=?", credential.ID).Scan(&stored); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(stored, replacement) {
				t.Fatal("rotation did not update the credential")
			}
		})
	}
}
