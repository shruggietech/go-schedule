package enrollment

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

func TestPairingIsSingleUseAndCredentialCanRotateAndRevoke(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	service := New(st)
	service.memory = 64
	pairing, err := service.Create("My laptop", domain.ActorKindDesktop, domain.CapabilityManage)
	if err != nil {
		t.Fatal(err)
	}
	if len(pairing.Phrase) == 0 || time.Until(pairing.ExpiresAt) > PairingLifetime {
		t.Fatal("invalid pairing lifetime")
	}
	issued, err := service.Exchange(pairing.ID, pairing.Phrase, pairing.DaemonID, pairing.DisplayName, pairing.Kind, pairing.Capability)
	if err != nil {
		t.Fatal(err)
	}
	if issued.Token == "" || issued.Actor.DisplayName != "My laptop" {
		t.Fatal("missing issued relationship")
	}
	if _, err := service.Exchange(pairing.ID, pairing.Phrase, pairing.DaemonID, pairing.DisplayName, pairing.Kind, pairing.Capability); !errors.Is(err, ErrRejected) {
		t.Fatalf("replay err=%v", err)
	}
	if _, _, err := service.Authenticate(issued.Token); err != nil {
		t.Fatal(err)
	}
	rotated, err := service.Rotate(issued.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := service.Authenticate(issued.Token); !errors.Is(err, ErrRejected) {
		t.Fatalf("old token err=%v", err)
	}
	if _, _, err := service.Authenticate(rotated.Token); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Revoke(issued.ID); err != nil {
		t.Fatal(err)
	}
	if _, _, err := service.Authenticate(rotated.Token); !errors.Is(err, ErrRejected) {
		t.Fatalf("revoked token err=%v", err)
	}
}

func TestConcurrentPairingExchangeCreatesExactlyOneRelationship(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	service := New(st)
	service.memory = 64
	pairing, err := service.Create("Concurrent client", domain.ActorKindCLI, domain.CapabilityOperate)
	if err != nil {
		t.Fatal(err)
	}
	var wait sync.WaitGroup
	successes := make(chan domain.IssuedCredential, 8)
	for range 8 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			issued, err := service.Exchange(pairing.ID, pairing.Phrase, pairing.DaemonID, pairing.DisplayName, pairing.Kind, pairing.Capability)
			if err == nil {
				successes <- issued
			}
		}()
	}
	wait.Wait()
	close(successes)
	if len(successes) != 1 {
		t.Fatalf("successful exchanges=%d, want 1", len(successes))
	}
	credentials, err := service.Credentials()
	if err != nil || len(credentials) != 1 {
		t.Fatalf("credentials=%d err=%v", len(credentials), err)
	}
}

func TestCancelledExpiredAndWrongDaemonPairingsFailClosed(t *testing.T) {
	st, _ := store.Open(":memory:")
	defer st.Close()
	service := New(st)
	service.memory = 64
	pairing, _ := service.Create("Cancelled", domain.ActorKindDesktop, domain.CapabilityObserve)
	service.Cancel(pairing.ID)
	if _, err := service.Exchange(pairing.ID, pairing.Phrase, pairing.DaemonID, pairing.DisplayName, pairing.Kind, pairing.Capability); !errors.Is(err, ErrRejected) {
		t.Fatal(err)
	}
	pairing, _ = service.Create("Wrong daemon", domain.ActorKindDesktop, domain.CapabilityObserve)
	if _, err := service.Exchange(pairing.ID, pairing.Phrase, "wrong", pairing.DisplayName, pairing.Kind, pairing.Capability); !errors.Is(err, ErrRejected) {
		t.Fatal(err)
	}
	service.now = func() time.Time { return pairing.ExpiresAt.Add(time.Second) }
	if _, err := service.Exchange(pairing.ID, pairing.Phrase, pairing.DaemonID, pairing.DisplayName, pairing.Kind, pairing.Capability); !errors.Is(err, ErrRejected) {
		t.Fatal(err)
	}
}

func TestWrongPhraseConsumesAttempts(t *testing.T) {
	st, _ := store.Open(":memory:")
	defer st.Close()
	service := New(st)
	service.memory = 64
	pairing, _ := service.Create("JSON client", domain.ActorKindJSON, domain.CapabilityObserve)
	wrong := "amber-amber-amber-amber-amber-amber-amber-amber-amber-amber"
	if wrong == pairing.Phrase {
		wrong = "apple-apple-apple-apple-apple-apple-apple-apple-apple-apple"
	}
	for range PairingAttempts {
		if _, err := service.Exchange(pairing.ID, wrong, pairing.DaemonID, pairing.DisplayName, pairing.Kind, pairing.Capability); !errors.Is(err, ErrRejected) {
			t.Fatal(err)
		}
	}
	items, _ := service.List()
	if items[0].State != domain.PairingExhausted || items[0].AttemptsRemaining != 0 {
		t.Fatalf("pairing=%+v", items[0])
	}
}
