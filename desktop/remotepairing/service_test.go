package remotepairing

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/clientprofile"
)

type fakeSecrets struct {
	probeErr                           error
	daemon, credential, token, deleted string
}

func (f *fakeSecrets) Probe() error { return f.probeErr }
func (f *fakeSecrets) Save(daemon, credential, token string) error {
	f.daemon, f.credential, f.token = daemon, credential, token
	return nil
}
func (f *fakeSecrets) Delete(_ string, credential string) error { f.deleted = credential; return nil }

type fakeProfiles struct {
	collection clientprofile.Collection
	added      clientprofile.Profile
}

func (f *fakeProfiles) Load() (clientprofile.Collection, error) { return f.collection, nil }
func (f *fakeProfiles) Add(profile clientprofile.Profile) (clientprofile.Profile, error) {
	f.added = profile
	f.collection = clientprofile.Collection{Version: clientprofile.CurrentVersion, Profiles: []clientprofile.Profile{profile}}
	return profile, nil
}
func (f *fakeProfiles) Replace(profile clientprofile.Profile) (clientprofile.Profile, error) {
	f.added = profile
	return profile, nil
}

type testError string

func (e testError) Error() string { return string(e) }

func TestPairProbesKeyringVerifiesIdentityAndPersistsSecretFreeProfile(t *testing.T) {
	enrollmentCalls := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/api/v1/enroll":
			enrollmentCalls++
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "credential-id", "daemon_id": "daemon-id", "token": "bearer-value", "actor": map[string]any{"kind": "desktop", "capability": "observe"}})
		case "/api/v1/manifest":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"installation_id":"daemon-id","display_name":"Remote daemon","product_version":"v1.4.0","remote_api_versions":["v1"],"capabilities":["tasks"],"platform":{"os":"linux","architecture":"amd64"}}`))
		default:
			t.Fatalf("path=%s", request.URL.Path)
		}
	}))
	defer server.Close()
	certificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw})
	secrets, profiles := &fakeSecrets{}, &fakeProfiles{}
	result := NewWithStores(secrets, profiles).Pair(context.Background(), Draft{ProfileLabel: "Production", Address: server.URL, DaemonID: "daemon-id", PairingID: "pairing-id", Phrase: "one-time", CertificatePEM: string(certificate), DisplayName: "Desktop", Capability: "observe"})
	if result.Outcome != "accepted" || enrollmentCalls != 1 || secrets.token != "bearer-value" || result.CredentialID != "credential-id" || result.ProfileID == "" {
		t.Fatalf("result=%+v calls=%d secrets=%+v", result, enrollmentCalls, secrets)
	}
	data, _ := json.Marshal(profiles.added)
	if profiles.added.Label != "Production" || profiles.added.CertificateFingerprint == "" || strings.Contains(string(data), "bearer-value") {
		t.Fatalf("profile=%+v", profiles.added)
	}
}

func TestPairDoesNotConsumePhraseWhenNativeStoreUnavailable(t *testing.T) {
	calls := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls++ }))
	defer server.Close()
	result := NewWithStores(&fakeSecrets{probeErr: testError("unavailable")}, &fakeProfiles{}).Pair(context.Background(), Draft{Address: server.URL, DaemonID: "daemon", PairingID: "pair", Phrase: "phrase", CertificatePEM: "invalid"})
	if result.Outcome != "rejected" || calls != 0 {
		t.Fatalf("result=%+v calls=%d", result, calls)
	}
}

func TestPairDoesNotForwardEnrollmentSecretsAcrossRedirects(t *testing.T) {
	receiverCalls := 0
	receiver := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { receiverCalls++ }))
	defer receiver.Close()
	redirectorCalls := 0
	redirector := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		redirectorCalls++
		http.Redirect(w, &http.Request{}, receiver.URL+"/api/v1/enroll", http.StatusTemporaryRedirect)
	}))
	defer redirector.Close()
	certificates := strings.Join([]string{string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: redirector.Certificate().Raw})), string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: receiver.Certificate().Raw}))}, "")
	result := NewWithStores(&fakeSecrets{}, &fakeProfiles{}).Pair(context.Background(), Draft{Address: redirector.URL, DaemonID: "daemon-id", PairingID: "pairing-id", Phrase: "one-time", CertificatePEM: certificates, DisplayName: "Desktop", Capability: "observe"})
	if result.Outcome != "rejected" || redirectorCalls != 1 || receiverCalls != 0 {
		t.Fatalf("result=%+v redirectorCalls=%d receiverCalls=%d", result, redirectorCalls, receiverCalls)
	}
}
