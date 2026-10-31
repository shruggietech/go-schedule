package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	apis "github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/enrollment"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

func TestRemoteNotificationCapabilitiesAndRedaction(t *testing.T) {
	for _, capability := range []domain.Capability{domain.CapabilityObserve, domain.CapabilityOperate, domain.CapabilityManage, domain.CapabilityEnroll} {
		t.Run(string(capability), func(t *testing.T) {
			st, err := store.Open(":memory:")
			if err != nil {
				t.Fatal(err)
			}
			defer st.Close()
			local, err := st.LocalActor()
			if err != nil {
				t.Fatal(err)
			}
			api := apis.New(st, nil, nil, nil, "", nil)
			api.SetActorResolver(ActorID)
			service := enrollment.New(st)
			pairing, err := service.Create("mail fixture", domain.ActorKindJSON, capability)
			if err != nil {
				t.Fatal(err)
			}
			issued, err := service.Exchange(pairing.ID, pairing.Phrase, pairing.DaemonID, pairing.DisplayName, pairing.Kind, pairing.Capability)
			if err != nil {
				t.Fatal(err)
			}
			handler := NewHandler(api.Handler(), service, local.ID)
			config := domain.SMTPConfig{Host: "mail.example.test", Port: 587, TLSMode: "starttls", From: "sender@example.test", Recipients: []string{"ops@example.test"}, TimeoutSeconds: 10}
			channel := domain.NotificationChannel{Name: "existing", Kind: domain.NotificationChannelSMTP, SMTP: &config, Authorization: "{\"username\":\"private-user\",\"password\":\"private-password\"}", Enabled: true}
			if err := st.CreateNotificationChannel(&channel); err != nil {
				t.Fatal(err)
			}
			changedConfig := config
			changedConfig.Host = "other.example.test"
			for _, tc := range []struct {
				method, path string
				body         any
				required     domain.Capability
				status       int
			}{
				{http.MethodGet, "/api/v1/notification-channels", nil, domain.CapabilityObserve, 200},
				{http.MethodGet, "/api/v1/notification-channels/" + channel.ID, nil, domain.CapabilityObserve, 200},
				{http.MethodGet, "/api/v1/notification-deliveries", nil, domain.CapabilityObserve, 200},
				{http.MethodPost, "/api/v1/notification-channels", apis.NotificationChannelCreateRequest{Name: "new", Kind: domain.NotificationChannelSMTP, SMTP: &config}, domain.CapabilityManage, 201},
				{http.MethodPost, "/api/v1/notification-channels/" + channel.ID + "/disable", nil, domain.CapabilityManage, 200},
				{http.MethodPost, "/api/v1/notification-channels/" + channel.ID + "/enable", nil, domain.CapabilityManage, 200},
				{http.MethodPost, "/api/v1/notification-channels/" + channel.ID + "/test", nil, domain.CapabilityOperate, 202},
				{http.MethodPatch, "/api/v1/notification-channels/" + channel.ID, apis.NotificationChannelUpdateRequest{SMTP: &changedConfig}, domain.CapabilityEnroll, 200},
				{http.MethodPost, "/api/v1/notification-channels/" + channel.ID + "/rotate", apis.NotificationChannelRotateRequest{SMTPAuth: &domain.SMTPAuth{}}, domain.CapabilityEnroll, 200},
			} {
				var body io.Reader
				if tc.body != nil {
					encoded, err := json.Marshal(tc.body)
					if err != nil {
						t.Fatal(err)
					}
					body = bytes.NewReader(encoded)
				}
				req := httptest.NewRequest(tc.method, tc.path, body)
				req.Header.Set("Authorization", "Bearer "+issued.Token)
				if tc.body != nil {
					req.Header.Set("Content-Type", "application/json")
				}
				response := httptest.NewRecorder()
				handler.ServeHTTP(response, req)
				want := tc.status
				if !capability.Allows(tc.required) {
					want = http.StatusForbidden
				}
				if response.Code != want {
					t.Fatalf("%s %s status=%d want=%d body=%s", tc.method, tc.path, response.Code, want, response.Body.String())
				}
				if strings.Contains(response.Body.String(), "private-user") || strings.Contains(response.Body.String(), "private-password") {
					t.Fatal("remote credential disclosure")
				}
			}
		})
	}
}

func TestRemoteNotificationTypedClientBodylessActionsAndNormalizedAuthority(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	local, err := st.LocalActor()
	if err != nil {
		t.Fatal(err)
	}
	api := apis.New(st, nil, nil, nil, "", nil)
	api.SetActorResolver(ActorID)
	service := enrollment.New(st)
	pairing, err := service.Create("typed notification client", domain.ActorKindJSON, domain.CapabilityManage)
	if err != nil {
		t.Fatal(err)
	}
	issued, err := service.Exchange(pairing.ID, pairing.Phrase, pairing.DaemonID, pairing.DisplayName, pairing.Kind, pairing.Capability)
	if err != nil {
		t.Fatal(err)
	}
	receiver := httptest.NewTLSServer(NewHandler(api.Handler(), service, local.ID))
	defer receiver.Close()
	cert := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: receiver.Certificate().Raw}))
	remote, err := client.NewRemote(receiver.URL, cert, issued.Token, pairing.DaemonID)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if _, err := remote.VerifyIdentity(ctx); err != nil {
		t.Fatal(err)
	}
	config := domain.SMTPConfig{Host: "mail.example.test", Port: 587, TLSMode: "starttls", From: "sender@example.test", Recipients: []string{"ops@example.test"}, TimeoutSeconds: 10}
	channel := domain.NotificationChannel{Name: "typed client", Kind: domain.NotificationChannelSMTP, SMTP: &config, Authorization: "{\"username\":\"private-user\",\"password\":\"private-password\"}", Enabled: true}
	if err := st.CreateNotificationChannel(&channel); err != nil {
		t.Fatal(err)
	}
	for _, enabled := range []bool{false, true} {
		updated, err := remote.SetNotificationChannelEnabled(ctx, channel.ID, enabled)
		if err != nil || updated.Enabled != enabled {
			t.Fatalf("enabled=%t result=%+v error=%v", enabled, updated, err)
		}
	}
	delivery, err := remote.TestNotificationChannel(ctx, channel.ID)
	if err != nil || delivery.ChannelID != channel.ID {
		t.Fatalf("test result=%+v error=%v", delivery, err)
	}
	claimed, err := st.ClaimNotificationDeliveries(1, time.Now().Add(time.Minute))
	if err != nil || len(claimed) != 1 {
		t.Fatalf("claim %v", err)
	}
	claimed[0].Recipients[0].State = "accepted"
	if err := st.FinishSMTPAttempt(claimed[0], domain.SMTPAttemptResult{Recipients: claimed[0].Recipients, Status: 250, Diagnostic: "The SMTP server accepted the message."}, time.Now(), time.Now()); err != nil {
		t.Fatal(err)
	}
	for _, typed := range []*client.Client{remote, client.NewInProcess(apis.New(st, nil, nil, nil, "", nil).Handler())} {
		deliveries, err := typed.ListNotificationDeliveries(ctx, domain.NotificationDeliveryFilter{ChannelID: channel.ID})
		if err != nil || len(deliveries) != 1 {
			t.Fatalf("typed history %v", err)
		}
		got := deliveries[0]
		if got.Diagnostic != "The SMTP server accepted the message." || got.LastError != "" {
			t.Fatalf("typed detail mislabeled: %+v", got)
		}
	}
	config.Host = " mail.example.test "
	config.Port = 0
	config.Recipients = []string{"new@example.test"}
	updated, err := remote.UpdateNotificationChannel(ctx, channel.ID, apis.NotificationChannelUpdateRequest{SMTP: &config})
	if err != nil || updated.SMTP == nil || updated.SMTP.Host != "mail.example.test" || updated.SMTP.Port != 587 || !updated.HasAuthorization || updated.SMTP.Recipients[0] != "new@example.test" {
		t.Fatalf("normalized update=%+v error=%v", updated, err)
	}
	config.Host = "other.example.test"
	_, err = remote.UpdateNotificationChannel(ctx, channel.ID, apis.NotificationChannelUpdateRequest{SMTP: &config})
	var denied *client.StatusError
	if !errors.As(err, &denied) || denied.Code != apis.CodeForbidden {
		t.Fatalf("retarget error=%v", err)
	}
	stored, err := st.GetNotificationChannel(channel.ID)
	if err != nil || stored.SMTP.Host != "mail.example.test" {
		t.Fatalf("denied retarget changed channel: %+v error=%v", stored, err)
	}
}
