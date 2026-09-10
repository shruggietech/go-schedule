package cli

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/clientprofile"
)

func TestResolveTargetPreservesLocalDefaultAndRejectsConflicts(t *testing.T) {
	local, description, err := resolveTarget(context.Background(), remoteTargetFlags{})
	if err != nil || local.Remote() || description != "" {
		t.Fatalf("local=%v description=%q err=%v", local, description, err)
	}
	if _, _, err := resolveTarget(context.Background(), remoteTargetFlags{Profile: "prod", Endpoint: "https://example.test"}); err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("conflict err=%v", err)
	}
	if _, _, err := resolveTarget(context.Background(), remoteTargetFlags{Endpoint: "https://example.test"}); err == nil || !strings.Contains(err.Error(), "required together") {
		t.Fatalf("incomplete err=%v", err)
	}
}

func TestSafeProfileJSONExcludesCertificateAndBearerFields(t *testing.T) {
	encoded, err := json.Marshal(safeProfileOf(clientprofile.Profile{ID: "id", Label: "Production", CertificatePEM: "certificate-canary"}))
	if err != nil {
		t.Fatal(err)
	}
	text := string(encoded)
	if strings.Contains(text, "certificate-canary") || strings.Contains(text, "certificate_pem") || strings.Contains(text, "token") || strings.Contains(text, "phrase") {
		t.Fatalf("unsafe profile JSON: %s", text)
	}
}

type renameStoreFake struct {
	collection clientprofile.Collection
	id         string
	label      string
}

func (f *renameStoreFake) Load() (clientprofile.Collection, error) { return f.collection, nil }
func (f *renameStoreFake) Rename(id, label string) (clientprofile.Profile, error) {
	f.id, f.label = id, label
	return clientprofile.Profile{ID: id, Label: label}, nil
}

func TestRenameProfileResolvesUnambiguousLabel(t *testing.T) {
	store := &renameStoreFake{collection: clientprofile.Collection{Version: clientprofile.CurrentVersion, Profiles: []clientprofile.Profile{{ID: "profile-id", Label: "Production"}}}}
	profile, err := renameProfile(store, "Production", "Primary")
	if err != nil || store.id != "profile-id" || store.label != "Primary" || profile.Label != "Primary" {
		t.Fatalf("profile=%+v store=%+v err=%v", profile, store, err)
	}
}
