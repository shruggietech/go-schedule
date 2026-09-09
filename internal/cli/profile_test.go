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
