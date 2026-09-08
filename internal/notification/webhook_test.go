package notification

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestValidateEndpointAndSafeSummary(t *testing.T) {
	for _, raw := range []string{"ftp://example.test/hook", "http://example.test/hook", "https:///missing", "https://example.test/hook#fragment"} {
		if _, err := ValidateEndpoint(raw); err == nil {
			t.Fatalf("accepted %q", raw)
		}
	}
	summary, err := ValidateEndpoint("https://user:pass@example.test:8443/hook?token=secret")
	if err != nil || summary != "https://example.test:8443" {
		t.Fatalf("summary=%q err=%v", summary, err)
	}
	if _, err := ValidateEndpoint("http://127.0.0.1:8080/hook"); err != nil {
		t.Fatal(err)
	}
}

func TestSenderUsesStableHeadersAndRefusesRedirects(t *testing.T) {
	var gotDelivery, gotEvent, gotAuth string
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotDelivery, gotEvent, gotAuth = r.Header.Get("X-Go-Schedule-Delivery"), r.Header.Get("X-Go-Schedule-Event"), r.Header.Get("Authorization")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer receiver.Close()
	delivery := domain.NotificationDelivery{ID: "delivery-1", Endpoint: receiver.URL, Authorization: "Bearer secret", EventKind: domain.NotificationEventTest, Payload: json.RawMessage(`{"schema":"go-schedule.webhook.v1"}`)}
	status, diagnostic, err := NewSender(nil).Send(context.Background(), delivery)
	if err != nil || status != http.StatusNoContent || diagnostic != "" || gotDelivery != delivery.ID || gotEvent != "test" || gotAuth != "Bearer secret" {
		t.Fatalf("status=%d diagnostic=%q headers=%q/%q/%q err=%v", status, diagnostic, gotDelivery, gotEvent, gotAuth, err)
	}

	targetCalled := false
	target := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { targetCalled = true }))
	defer target.Close()
	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusTemporaryRedirect)
	}))
	defer redirect.Close()
	delivery.Endpoint = redirect.URL
	status, diagnostic, err = NewSender(nil).Send(context.Background(), delivery)
	if err == nil || status != http.StatusTemporaryRedirect || !strings.Contains(diagnostic, "307") || targetCalled {
		t.Fatalf("redirect status=%d diagnostic=%q called=%t err=%v", status, diagnostic, targetCalled, err)
	}
}
