package domain

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestRemoteAccessStatesAndSafeMetadata(t *testing.T) {
	for _, state := range []PairingState{PairingActive, PairingConsumed, PairingCancelled, PairingExpired, PairingExhausted} {
		if !state.Valid() {
			t.Fatalf("invalid pairing state %q", state)
		}
	}
	if PairingState("unknown").Valid() || !CredentialActive.Valid() || !CredentialRevoked.Valid() || CredentialState("unknown").Valid() {
		t.Fatal("state validation is not fail closed")
	}
	credential := ClientCredential{ID: "credential", ActorID: "actor", Fingerprint: "safe", State: CredentialActive}
	encoded, _ := json.Marshal(credential)
	for _, protected := range []string{"token", "digest", "verifier", "phrase"} {
		if strings.Contains(string(encoded), protected) {
			t.Fatalf("metadata contains %q: %s", protected, encoded)
		}
	}
}

func TestPairingCarriesSeparateEnrollmentAndGrantDeadlines(t *testing.T) {
	now := time.Now().UTC()
	grantExpires := now.Add(7 * 24 * time.Hour)
	pairing := PairingSession{CreatedAt: now, ExpiresAt: now.Add(10 * time.Minute), GrantExpiresAt: &grantExpires}
	if pairing.GrantExpiresAt == nil || !pairing.GrantExpiresAt.Equal(grantExpires) || !pairing.ExpiresAt.Before(*pairing.GrantExpiresAt) {
		t.Fatalf("pairing=%+v", pairing)
	}
}
