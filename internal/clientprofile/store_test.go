package clientprofile

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStoreLifecycleExcludesBearerCanary(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profiles.json")
	store := NewStore(path)
	profile := validProfile(t)
	if _, err := store.Add(profile); err != nil {
		t.Fatal(err)
	}
	if err := store.SetActive(profile.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Rename(profile.ID, "Renamed"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "bearer-canary") {
		t.Fatal("profile file contains bearer canary")
	}
	collection, err := store.Load()
	if err != nil || collection.ActiveDesktopProfileID != profile.ID || collection.Profiles[0].Label != "Renamed" {
		t.Fatalf("collection = %#v, err = %v", collection, err)
	}
	if _, err := store.Remove(profile.ID); err != nil {
		t.Fatal(err)
	}
	collection, err = store.Load()
	if err != nil || collection.ActiveDesktopProfileID != "" || len(collection.Profiles) != 0 {
		t.Fatalf("collection = %#v, err = %v", collection, err)
	}
}

func TestStoreRejectsCorruptAndFutureDocuments(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profiles.json")
	for _, value := range []string{"{", `{"version":2,"profiles":[]}`} {
		if err := os.WriteFile(path, []byte(value), 0o600); err != nil {
			t.Fatal(err)
		}
		_, err := NewStore(path).Load()
		if err == nil {
			t.Fatalf("accepted %q", value)
		}
	}
}

func TestStoreBusyLeavesDocumentIntact(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profiles.json")
	store := NewStore(path)
	if _, err := store.Add(validProfile(t)); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(path)
	if err := os.WriteFile(path+".lock", []byte("busy"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := store.Rename("profile-1", "Other")
	if !errors.Is(err, ErrBusy) {
		t.Fatalf("err = %v", err)
	}
	after, _ := os.ReadFile(path)
	if string(after) != string(before) {
		t.Fatal("busy write changed document")
	}
}

func TestRenameChangesOnlyPresentationFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profiles.json")
	store := NewStore(path)
	profile := validProfile(t)
	profile.CredentialID = "repaired-credential"
	stored, err := store.Add(profile)
	if err != nil {
		t.Fatal(err)
	}
	renamed, err := store.Rename(profile.ID, "Renamed")
	if err != nil {
		t.Fatal(err)
	}
	if renamed.Label != "Renamed" || renamed.CertificateFingerprint != stored.CertificateFingerprint || renamed.CredentialID != "repaired-credential" {
		t.Fatalf("renamed profile = %#v", renamed)
	}
}

func TestRemoveWithUsesCurrentProfileAndKeepsMetadataOnCredentialFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profiles.json")
	store := NewStore(path)
	profile, err := store.Add(validProfile(t))
	if err != nil {
		t.Fatal(err)
	}
	profile.CredentialID = "replacement-credential"
	if _, err := store.Replace(profile); err != nil {
		t.Fatal(err)
	}
	credentialErr := errors.New("keyring unavailable")
	_, err = store.RemoveWith(profile.ID, func(current Profile) error {
		if current.CredentialID != "replacement-credential" {
			t.Fatalf("callback profile = %#v", current)
		}
		return credentialErr
	})
	if !errors.Is(err, credentialErr) {
		t.Fatalf("err = %v", err)
	}
	collection, err := store.Load()
	if err != nil || len(collection.Profiles) != 1 || collection.Profiles[0].CredentialID != "replacement-credential" {
		t.Fatalf("collection=%#v err=%v", collection, err)
	}
}
