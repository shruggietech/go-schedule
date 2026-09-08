package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNotificationSecretsAreExcludedFromJSON(t *testing.T) {
	secretURL := "https://user:pass@example.test/hook?token=secret"
	authorization := "Bearer very-secret"
	values := []any{
		NotificationChannel{ID: "channel", Endpoint: secretURL, EndpointSummary: "https://example.test", Authorization: authorization},
		NotificationDelivery{ID: "delivery", Endpoint: secretURL, DestinationSummary: "https://example.test", Authorization: authorization, Payload: json.RawMessage(`{"safe":true}`)},
	}
	for _, value := range values {
		encoded, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		for _, secret := range []string{secretURL, authorization, "user:pass", "token=secret"} {
			if strings.Contains(string(encoded), secret) {
				t.Fatalf("JSON disclosed %q: %s", secret, encoded)
			}
		}
	}
}
