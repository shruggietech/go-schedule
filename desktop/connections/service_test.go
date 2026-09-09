package connections

import (
	"context"
	"encoding/pem"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
)

type fakeProfileStore struct {
	collection               clientprofile.Collection
	renamed, active, removed string
}

func (f *fakeProfileStore) Load() (clientprofile.Collection, error) { return f.collection, nil }
func (f *fakeProfileStore) Rename(id, label string) (clientprofile.Profile, error) {
	f.renamed = id + ":" + label
	return f.collection.Profiles[0], nil
}
func (f *fakeProfileStore) SetActive(id string) error {
	f.active = id
	f.collection.ActiveDesktopProfileID = id
	return nil
}
func (f *fakeProfileStore) Remove(id string) (clientprofile.Profile, error) {
	f.removed = id
	removed := f.collection.Profiles[0]
	f.collection.Profiles = nil
	return removed, nil
}

type fakeSecretStore struct {
	deleted   string
	deleteErr error
	token     string
}

func (f *fakeSecretStore) Load(string, string) (string, error) {
	if f.token == "" {
		return "", errors.New("missing")
	}
	return f.token, nil
}
func (f *fakeSecretStore) Delete(daemon, credential string) error {
	f.deleted = daemon + ":" + credential
	return f.deleteErr
}

type fakeManager struct {
	target   connection.Target
	switches int
}

func (f *fakeManager) Switch(_ connection.Backend, target connection.Target) bool {
	f.target = target
	f.switches++
	return true
}

func profileFixture() clientprofile.Profile {
	return clientprofile.Profile{ID: "profile-id", Label: "Workshop", Endpoint: "https://example.test", DaemonID: "daemon-identity", CredentialID: "credential-id", CertificateFingerprint: "abcdef", Capability: "observe", Platform: "linux", ProductVersion: "1.4.0", LastSuccessfulAt: time.Date(2026, 9, 9, 1, 2, 3, 0, time.UTC)}
}

func TestWorkspaceExcludesTrustMaterialAndDisambiguatesProfiles(t *testing.T) {
	profiles := &fakeProfileStore{collection: clientprofile.Collection{Version: clientprofile.CurrentVersion, ActiveDesktopProfileID: "profile-id", Profiles: []clientprofile.Profile{profileFixture()}}}
	service := New(profiles, &fakeSecretStore{}, client.New("local"), client.NewSwitchable(client.New("local")), &fakeManager{})
	result := service.Workspace()
	if result.Outcome != "accepted" || result.Workspace == nil || len(result.Workspace.Profiles) != 1 {
		t.Fatalf("result=%+v", result)
	}
	profile := result.Workspace.Profiles[0]
	if profile.ShortDaemonID != "daemon-i" || !profile.Active || profile.LastSuccessfulAt != "2026-09-09T01:02:03Z" {
		t.Fatalf("profile=%+v", profile)
	}
}

func TestRemovingActiveProfileSwitchesLocalBeforeDeletingCredential(t *testing.T) {
	profiles := &fakeProfileStore{collection: clientprofile.Collection{Version: clientprofile.CurrentVersion, ActiveDesktopProfileID: "profile-id", Profiles: []clientprofile.Profile{profileFixture()}}}
	secrets := &fakeSecretStore{}
	manager := &fakeManager{}
	local := client.New("local")
	result := New(profiles, secrets, local, client.NewSwitchable(local), manager).Remove("profile-id")
	if result.Outcome != "accepted" || profiles.active != "" || secrets.deleted != "daemon-identity:credential-id" || profiles.removed != "profile-id" || manager.target.Kind != "local" {
		t.Fatalf("result=%+v profiles=%+v secrets=%+v manager=%+v", result, profiles, secrets, manager)
	}
}

func TestRemovalKeepsMetadataWhenCredentialDeletionFails(t *testing.T) {
	profiles := &fakeProfileStore{collection: clientprofile.Collection{Version: clientprofile.CurrentVersion, Profiles: []clientprofile.Profile{profileFixture()}}}
	secrets := &fakeSecretStore{deleteErr: errors.New("keyring unavailable")}
	result := New(profiles, secrets, client.New("local"), client.NewSwitchable(client.New("local")), &fakeManager{}).Remove("profile-id")
	if result.Outcome != "rejected" || profiles.removed != "" {
		t.Fatalf("result=%+v removed=%q", result, profiles.removed)
	}
}

func TestRestorePreservesSelectedRemoteBeforeConnectionAttempt(t *testing.T) {
	tlsServer := httptest.NewUnstartedServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	tlsServer.StartTLS()
	defer tlsServer.Close()
	profile := profileFixture()
	profile.Endpoint = tlsServer.URL
	profile.CertificatePEM = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: tlsServer.Certificate().Raw}))
	profiles := &fakeProfileStore{collection: clientprofile.Collection{Version: clientprofile.CurrentVersion, ActiveDesktopProfileID: profile.ID, Profiles: []clientprofile.Profile{profile}}}
	manager := &fakeManager{}
	local := client.New("local")
	router := client.NewSwitchable(local)
	result := New(profiles, &fakeSecretStore{token: "bearer-canary"}, local, router, manager).RestoreSelection(context.Background())
	if result.Outcome != "accepted" || manager.target.ProfileID != profile.ID || manager.target.Kind != "remote" || !router.Remote() {
		t.Fatalf("result=%+v manager=%+v remote=%v", result, manager, router.Remote())
	}
}
