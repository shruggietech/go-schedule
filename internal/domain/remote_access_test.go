package domain

import (
	"encoding/json"
	"strings"
	"testing"
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
