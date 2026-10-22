package notification

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func smtpTestConfig() domain.SMTPConfig {
	return domain.SMTPConfig{Host: "mail.example.test", Port: 587, TLSMode: "starttls", From: "scheduler@example.test", Recipients: []string{"ops@example.test"}, TimeoutSeconds: 10}
}

func TestSMTPConfigValidation(t *testing.T) {
	good := smtpTestConfig()
	normalized, err := ValidateSMTP(good)
	if err != nil || normalized.Host != good.Host {
		t.Fatalf("valid config: %v %v", normalized, err)
	}
	cases := []struct {
		name   string
		change func(*domain.SMTPConfig)
		field  string
	}{
		{"host", func(c *domain.SMTPConfig) { c.Host = "smtp://bad" }, "host"},
		{"port", func(c *domain.SMTPConfig) { c.Port = 65536 }, "port"},
		{"tls", func(c *domain.SMTPConfig) { c.TLSMode = "opportunistic" }, "tls_mode"},
		{"from injection", func(c *domain.SMTPConfig) { c.From = "a@example.test\r\nBcc:bad@example.test" }, "from"},
		{"missing recipients", func(c *domain.SMTPConfig) { c.Recipients = nil }, "recipients"},
		{"recipient injection", func(c *domain.SMTPConfig) { c.Recipients = []string{"a@example.test\x00bad"} }, "recipients"},
		{"unicode envelope", func(c *domain.SMTPConfig) { c.Recipients = []string{"é@example.test"} }, "recipients"},
		{"recipient bound", func(c *domain.SMTPConfig) { c.Recipients = make([]string, 101) }, "recipients"},
		{"timeout", func(c *domain.SMTPConfig) { c.TimeoutSeconds = 61 }, "timeout_seconds"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := smtpTestConfig()
			tc.change(&c)
			_, err := ValidateSMTP(c)
			if err == nil || !strings.Contains(err.Error(), tc.field) {
				t.Fatalf("expected field %s, got %v", tc.field, err)
			}
		})
	}
	good.Recipients = []string{"ops@example.test", "ops@example.test"}
	got, err := ValidateSMTP(good)
	if err != nil || len(got.Recipients) != 1 {
		t.Fatalf("deduplicate: %v %v", got, err)
	}
}

func TestSMTPAuthenticationAndPublicRedaction(t *testing.T) {
	auth := domain.SMTPAuth{Username: "operator", Password: "sensitive-fixture"}
	raw, err := EncodeSMTPAuth(auth)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ParseSMTPAuth(raw)
	if err != nil || got != auth {
		t.Fatalf("auth round trip: %v", err)
	}
	for _, raw := range []string{"not-json", "{\"username\":\"a\",\"password\":\"b\",\"extra\":\"c\"}", "{\"username\":\"a\",\"password\":\"\"}"} {
		if _, err := ParseSMTPAuth(raw); err == nil {
			t.Fatalf("accepted invalid auth")
		}
	}
	config := smtpTestConfig()
	channel := domain.NotificationChannel{Kind: domain.NotificationChannelSMTP, SMTP: &config, Authorization: raw, HasAuthorization: true}
	data, err := json.Marshal(channel)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), auth.Password) || strings.Contains(string(data), auth.Username) {
		t.Fatal("credential leaked")
	}
}
