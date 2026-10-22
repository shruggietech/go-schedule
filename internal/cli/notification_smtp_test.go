package cli

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestSMTPCLILifecycleAndSecretInput(t *testing.T) {
	c := isolatedCommandBackend(t)
	out, err := runFixtureCommand(c, "{\"username\":\"fixture-user\",\"password\":\"fixture-password\"}", "notification", "channel", "add", "Email", "--kind", "smtp", "--smtp-host", "mail.example.test", "--smtp-port", "587", "--smtp-tls", "starttls", "--smtp-from", "sender@example.test", "--smtp-to", "one@example.test,two@example.test", "--smtp-timeout", "5s", "--smtp-html", "--smtp-auth-file", "-", "--json")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "fixture-user") || strings.Contains(out, "fixture-password") {
		t.Fatal("secret in CLI output")
	}
	var channel domain.NotificationChannel
	if err := json.Unmarshal([]byte(out), &channel); err != nil {
		t.Fatal(err)
	}
	if channel.SMTP == nil || len(channel.SMTP.Recipients) != 2 || !channel.HasAuthorization {
		t.Fatalf("channel: %s", out)
	}
	if _, err := runFixtureCommand(c, "", "notification", "channel", "update", channel.ID, "--smtp-timeout", "3s", "--smtp-html=false", "--json"); err != nil {
		t.Fatal(err)
	}
	got, err := c.GetNotificationChannel(context.Background(), channel.ID)
	if err != nil || got.SMTP.TimeoutSeconds != 3 || got.SMTP.HTML || got.SMTP.Host != "mail.example.test" {
		t.Fatalf("update: %+v %v", got, err)
	}
	if _, err := runFixtureCommand(c, "{}", "notification", "channel", "rotate", channel.ID, "--smtp-auth-file", "-"); err != nil {
		t.Fatal(err)
	}
	got, err = c.GetNotificationChannel(context.Background(), channel.ID)
	if err != nil || got.HasAuthorization {
		t.Fatalf("clear: %+v %v", got, err)
	}
	for _, args := range [][]string{{"notification", "channel", "update", channel.ID, "--smtp-timeout", "500ms"}, {"notification", "channel", "add", "bad", "--kind", "smtp", "--smtp-host", "mail.example.test", "--smtp-from", "bad\r\nheader", "--smtp-to", "one@example.test"}, {"notification", "channel", "rotate", channel.ID, "--smtp-auth-file", "-", "--authorization", "conflict"}} {
		if _, err := runFixtureCommand(c, "{}", args...); err == nil {
			t.Fatalf("invalid input accepted: %v", args)
		}
	}
}

func TestSMTPCLIQuotedRecipients(t *testing.T) {
	c := isolatedCommandBackend(t)
	raw := `"ops,night"@example.test`
	out, err := runFixtureCommand(c, "", "notification", "channel", "add", "Quoted", "--kind=smtp", "--smtp-host=mail.example.test", "--smtp-tls=plain", "--smtp-from", raw, "--smtp-to", raw+",other@example.test", "--smtp-to", `"ops team"@example.test`, "--json")
	if err != nil {
		t.Fatal(err)
	}
	var channel domain.NotificationChannel
	if err := json.Unmarshal([]byte(out), &channel); err != nil {
		t.Fatal(err)
	}
	if channel.SMTP == nil || channel.SMTP.From != raw || len(channel.SMTP.Recipients) != 3 || channel.SMTP.Recipients[0] != raw || channel.SMTP.Recipients[2] != `"ops team"@example.test` {
		t.Fatalf("quoted configuration %+v", channel.SMTP)
	}
}
