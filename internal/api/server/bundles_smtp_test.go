package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/bundle"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestSMTPBundleDefinitionsProvisioningDriftAndStaleCredentials(t *testing.T) {
	s := newTestServer(t)
	config := domain.SMTPConfig{Host: "mail.example.test", Port: 587, TLSMode: "starttls", From: "sender@example.test", Recipients: []string{"ops@example.test"}, TimeoutSeconds: 10}
	source := bundle.SMTPChannel{Name: "Mail", Settings: config, Enabled: true, RequiresAuthorization: true}
	doc := bundle.Document{Schema: bundle.SchemaV2, SMTPChannels: []bundle.SMTPChannel{source}}
	preview := func() bundle.Plan {
		t.Helper()
		rec := doJSON(t, s, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc})
		if rec.Code != 200 {
			t.Fatalf("preview: %d %s", rec.Code, rec.Body.String())
		}
		var plan bundle.Plan
		if err := json.Unmarshal(rec.Body.Bytes(), &plan); err != nil {
			t.Fatal(err)
		}
		return plan
	}
	apply := func(p bundle.Plan) int {
		t.Helper()
		return doJSON(t, s, http.MethodPost, "/v1/bundles/apply", BundleApplyRequest{PlanID: p.ID, BundleDigest: p.BundleDigest, TargetDaemonID: p.TargetDaemonID, TargetFingerprint: p.TargetFingerprint}).Code
	}
	plan := preview()
	if len(plan.Items) != 1 || plan.Items[0].Action != bundle.ActionConflict || apply(plan) != http.StatusConflict {
		t.Fatalf("missing provisioning: %+v", plan)
	}
	channel := domain.NotificationChannel{Name: "Mail", Kind: domain.NotificationChannelSMTP, SMTP: &config, Authorization: "protected-fixture", Enabled: false}
	if err := s.store.CreateNotificationChannel(&channel); err != nil {
		t.Fatal(err)
	}
	plan = preview()
	if plan.Items[0].Action != bundle.ActionUpdate {
		t.Fatalf("drift: %+v", plan)
	}
	if err := s.store.RotateNotificationChannelAuthorization(channel.ID, "rotated-fixture"); err != nil {
		t.Fatal(err)
	}
	if apply(plan) != http.StatusConflict {
		t.Fatal("stale credential binding accepted")
	}
	current, err := s.store.GetNotificationChannel(channel.ID)
	if err != nil || current.Enabled {
		t.Fatal("stale apply changed channel")
	}
	plan = preview()
	if apply(plan) != 200 {
		t.Fatal("valid apply failed")
	}
	current, err = s.store.GetNotificationChannel(channel.ID)
	if err != nil || !current.Enabled || current.Authorization != "rotated-fixture" {
		t.Fatalf("apply lost local credential: %v", err)
	}
	plan = preview()
	if plan.Items[0].Action != bundle.ActionUnchanged {
		t.Fatalf("round trip: %+v", plan)
	}
	exported := doJSON(t, s, http.MethodGet, "/v1/bundles/export", nil)
	if exported.Code != 200 || !strings.Contains(exported.Body.String(), "\"smtp_channels\"") || strings.Contains(exported.Body.String(), "rotated-fixture") {
		t.Fatalf("export: %d %s", exported.Code, exported.Body.String())
	}
	doc.SMTPChannels[0].Settings.Host = "other.example.test"
	plan = preview()
	if plan.Items[0].Action != bundle.ActionConflict || apply(plan) != http.StatusConflict {
		t.Fatal("bundle redirected existing credentials")
	}
}
