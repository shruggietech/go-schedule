package remotepairing

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeSecrets struct {
	probeErr                  error
	daemon, credential, token string
}

func (f *fakeSecrets) Probe() error { return f.probeErr }
func (f *fakeSecrets) Save(daemon, credential, token string) error {
	f.daemon, f.credential, f.token = daemon, credential, token
	return nil
}

type testError string

func (e testError) Error() string { return string(e) }

func TestPairProbesKeyringBeforeExchangeAndStoresOnlyNativeSecret(t *testing.T) {
	calls := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Path != "/api/v1/enroll" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{"id": "credential-id", "daemon_id": "daemon-id", "token": "bearer-value", "actor": map[string]any{}})
	}))
	defer server.Close()
	certificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw})
	secrets := &fakeSecrets{}
	result := NewWithStore(secrets).Pair(context.Background(), Draft{Address: server.URL, DaemonID: "daemon-id", PairingID: "pairing-id", Phrase: "one-time", CertificatePEM: string(certificate), DisplayName: "Desktop", Capability: "observe"})
	if result.Outcome != "accepted" || calls != 1 || secrets.token != "bearer-value" || result.CredentialID != "credential-id" {
		t.Fatalf("result=%+v calls=%d secrets=%+v", result, calls, secrets)
	}
}

func TestPairDoesNotConsumePhraseWhenNativeStoreUnavailable(t *testing.T) {
	calls := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls++ }))
	defer server.Close()
	result := NewWithStore(&fakeSecrets{probeErr: testError("unavailable")}).Pair(context.Background(), Draft{Address: server.URL, DaemonID: "daemon", PairingID: "pair", Phrase: "phrase", CertificatePEM: "invalid"})
	if result.Outcome != "rejected" || calls != 0 {
		t.Fatalf("result=%+v calls=%d", result, calls)
	}
}
