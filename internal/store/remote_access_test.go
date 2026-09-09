package store

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestRemoteAccessPersistenceLifecycleAndFailureBounds(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	now := time.Now().UTC()
	session := domain.PairingSession{ID: uuid.NewString(), DisplayName: "Stored client", Kind: domain.ActorKindCLI, Capability: domain.CapabilityOperate, AttemptsRemaining: 5, State: domain.PairingActive, CreatedAt: now, ExpiresAt: now.Add(time.Minute)}
	salt := make([]byte, 16)
	verifier := make([]byte, 32)
	if err := st.CreatePairing(session, salt, verifier); err != nil {
		t.Fatal(err)
	}
	if items, err := st.ListPairings(); err != nil || len(items) != 1 {
		t.Fatalf("pairings=%d err=%v", len(items), err)
	}
	if _, err := st.PairingSalt("missing"); !errors.Is(err, ErrEnrollmentRejected) {
		t.Fatalf("salt err=%v", err)
	}
	actor := domain.Actor{ID: uuid.NewString(), State: domain.ActorStateActive, CreatedAt: now, UpdatedAt: now}
	credential := domain.ClientCredential{ID: uuid.NewString(), Fingerprint: "fingerprint", State: domain.CredentialActive, CreatedAt: now, UpdatedAt: now}
	digest := make([]byte, 32)
	actor, credential, err = st.ExchangePairing(session.ID, verifier, actor, credential, digest, now)
	if err != nil {
		t.Fatal(err)
	}
	if actor.DisplayName != session.DisplayName || credential.ActorID != actor.ID {
		t.Fatalf("actor=%+v credential=%+v", actor, credential)
	}
	if items, err := st.ListCredentials(); err != nil || len(items) != 1 {
		t.Fatalf("credentials=%d err=%v", len(items), err)
	}
	if _, _, err := st.AuthenticateCredential(digest, now); err != nil {
		t.Fatal(err)
	}
	if _, err := st.RotateCredential("missing", digest, "new", now); !errors.Is(err, ErrNotFound) {
		t.Fatalf("rotate err=%v", err)
	}
	if _, err := st.RevokeCredential("missing", now); !errors.Is(err, ErrNotFound) {
		t.Fatalf("revoke err=%v", err)
	}
	if _, err := st.CancelPairing("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cancel err=%v", err)
	}
	if _, _, err := st.AuthenticateCredential([]byte("missing"), now); !errors.Is(err, ErrCredentialRejected) {
		t.Fatalf("authenticate err=%v", err)
	}
}

func TestCreatePairingRejectsUnsafePersistenceInput(t *testing.T) {
	st, _ := Open(":memory:")
	defer st.Close()
	if err := st.CreatePairing(domain.PairingSession{}, nil, nil); !errors.Is(err, domain.ErrInvalidActor) {
		t.Fatalf("err=%v", err)
	}
}
