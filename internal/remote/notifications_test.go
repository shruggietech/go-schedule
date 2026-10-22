package remote

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apis "github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/enrollment"
	"github.com/shruggietech/go-schedule/internal/store"
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
				{http.MethodPost, "/api/v1/notification-channels/" + channel.ID + "/test", nil, domain.CapabilityOperate, 202},
				{http.MethodPatch, "/api/v1/notification-channels/" + channel.ID, apis.NotificationChannelUpdateRequest{SMTP: &changedConfig}, domain.CapabilityEnroll, 200},
				{http.MethodPost, "/api/v1/notification-channels/" + channel.ID + "/rotate", apis.NotificationChannelRotateRequest{SMTPAuth: &domain.SMTPAuth{}}, domain.CapabilityEnroll, 200},
			} {
				body, err := json.Marshal(tc.body)
				if err != nil {
					t.Fatal(err)
				}
				req := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(body))
				req.Header.Set("Authorization", "Bearer "+issued.Token)
				if tc.method != http.MethodGet {
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
