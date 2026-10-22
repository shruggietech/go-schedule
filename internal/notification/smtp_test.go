package notification

import (
	"context"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/testsmtp"
)

func mailFixtureDelivery(t *testing.T, receiver *testsmtp.Server, mode string) domain.NotificationDelivery {
	t.Helper()
	host, portRaw, err := net.SplitHostPort(receiver.Addr())
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(portRaw)
	if err != nil {
		t.Fatal(err)
	}
	config := domain.SMTPConfig{Host: host, Port: port, TLSMode: mode, From: "sender@example.test", Recipients: []string{"good@example.test", "temporary@example.test", "bad@example.test"}, TimeoutSeconds: 2, HTML: true}
	event := domain.WebhookEvent{Event: "run.completed", Task: &domain.WebhookTask{Name: "保守 <script>alert(1)</script> café", GroupName: "Operations"}, Run: &domain.WebhookRun{ID: "run-owned", Outcome: domain.OutcomeFailure}}
	payload, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	recipients := []domain.NotificationRecipient{}
	for _, address := range config.Recipients {
		recipients = append(recipients, domain.NotificationRecipient{Address: address, State: "pending"})
	}
	return domain.NotificationDelivery{ID: "owned-delivery", Kind: domain.NotificationChannelSMTP, SMTP: &config, Recipients: recipients, Payload: payload, CreatedAt: time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC)}
}
func TestSMTPPartialRecipientsAndRetryExclusions(t *testing.T) {
	receiver, err := testsmtp.Start(testsmtp.Options{Reject: map[string]int{"temporary@example.test": 450, "bad@example.test": 550}})
	if err != nil {
		t.Fatal(err)
	}
	defer receiver.Close()
	d := mailFixtureDelivery(t, receiver, "plain")
	checkpoints := 0
	result := NewSender(nil).SendSMTP(context.Background(), d, func(out []domain.NotificationRecipient) error {
		checkpoints++
		if out[0].State != "submitting" || out[1].State != "temporary_failure" || out[2].State != "permanent_failure" {
			t.Fatalf("checkpoint: %+v", out)
		}
		return nil
	})
	if checkpoints != 1 || result.Recipients[0].State != "accepted" || result.Recipients[1].State != "temporary_failure" || result.Recipients[2].State != "permanent_failure" {
		t.Fatalf("result: %+v", result)
	}
	receiver.SetRecipientCode("temporary@example.test", 0)
	d.Recipients = result.Recipients
	result = NewSender(nil).SendSMTP(context.Background(), d, func(_ []domain.NotificationRecipient) error { return nil })
	if result.Recipients[0].State != "accepted" || result.Recipients[1].State != "accepted" || len(receiver.Messages()) != 2 {
		t.Fatalf("retry result: %+v", result)
	}
	rcpts := strings.Join(receiver.Recipients(), ",")
	if strings.Count(rcpts, "good@example.test") != 1 || strings.Count(rcpts, "bad@example.test") != 1 || strings.Count(rcpts, "temporary@example.test") != 2 {
		t.Fatalf("replayed recipients: %s", rcpts)
	}
	msg, err := mail.ReadMessage(strings.NewReader(string(receiver.Messages()[0].Data)))
	if err != nil {
		t.Fatal(err)
	}
	subject, err := new(mime.WordDecoder).DecodeHeader(msg.Header.Get("Subject"))
	if err != nil || !strings.Contains(subject, "café") {
		t.Fatalf("subject: %q %v", subject, err)
	}
	_, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		t.Fatal(err)
	}
	reader := multipart.NewReader(msg.Body, params["boundary"])
	part, err := reader.NextPart()
	if err != nil {
		t.Fatal(err)
	}
	text, err := io.ReadAll(part)
	if err != nil || !strings.Contains(string(text), "保守") {
		t.Fatalf("text: %s %v", text, err)
	}
	part, err = reader.NextPart()
	if err != nil {
		t.Fatal(err)
	}
	html, err := io.ReadAll(part)
	if err != nil || strings.Contains(string(html), "<script>") || !strings.Contains(string(html), "&lt;script&gt;") {
		t.Fatalf("HTML not escaped: %s %v", html, err)
	}
}
func TestSMTPVerifiedTLSAndAuthentication(t *testing.T) {
	for _, mode := range []string{"starttls", "tls"} {
		t.Run(mode, func(t *testing.T) {
			receiver, err := testsmtp.Start(testsmtp.Options{TLSMode: mode, Username: "owned-user", Password: "owned-secret", ReplyText: "owned-secret"})
			if err != nil {
				t.Fatal(err)
			}
			defer receiver.Close()
			d := mailFixtureDelivery(t, receiver, mode)
			d.SMTP.Recipients = []string{"good@example.test"}
			d.Recipients = d.Recipients[:1]
			d.Authorization, err = EncodeSMTPAuth(domain.SMTPAuth{Username: "owned-user", Password: "owned-secret"})
			if err != nil {
				t.Fatal(err)
			}
			checkpoint := func(_ []domain.NotificationRecipient) error { return nil }
			result := NewSenderWithSMTPRoots(receiver.Roots()).SendSMTP(context.Background(), d, checkpoint)
			if result.Recipients[0].State != "accepted" {
				t.Fatalf("verified TLS: %+v", result)
			}
			result = NewSender(nil).SendSMTP(context.Background(), d, checkpoint)
			if result.Recipients[0].State == "accepted" || strings.Contains(result.Diagnostic, "owned-secret") {
				t.Fatalf("untrusted TLS accepted: %+v", result)
			}
			d.Authorization, err = EncodeSMTPAuth(domain.SMTPAuth{Username: "owned-user", Password: "wrong-secret"})
			if err != nil {
				t.Fatal(err)
			}
			result = NewSenderWithSMTPRoots(receiver.Roots()).SendSMTP(context.Background(), d, checkpoint)
			if result.Recipients[0].State != "permanent_failure" || strings.Contains(result.Diagnostic, "owned-secret") {
				t.Fatalf("authentication: %+v", result)
			}
		})
	}
}
func TestSMTPRequiredTLSNeverDowngrades(t *testing.T) {
	receiver, err := testsmtp.Start(testsmtp.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer receiver.Close()
	d := mailFixtureDelivery(t, receiver, "starttls")
	result := NewSender(nil).SendSMTP(context.Background(), d, func(_ []domain.NotificationRecipient) error { return nil })
	if len(receiver.Messages()) != 0 || !strings.Contains(result.Diagnostic, "no downgrade") {
		t.Fatalf("TLS downgrade: %+v", result)
	}
}
func TestSMTPSharedCancellationBoundsGreeting(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	accepted := make(chan net.Conn, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		accepted <- conn
		var one [1]byte
		_, _ = conn.Read(one[:])
		_ = conn.Close()
	}()
	receiver, err := testsmtp.Start(testsmtp.Options{})
	if err != nil {
		t.Fatal(err)
	}
	d := mailFixtureDelivery(t, receiver, "plain")
	if err := receiver.Close(); err != nil {
		t.Fatal(err)
	}
	host, portRaw, _ := net.SplitHostPort(listener.Addr().String())
	port, _ := strconv.Atoi(portRaw)
	d.SMTP.Host = host
	d.SMTP.Port = port
	d.SMTP.TimeoutSeconds = 60
	ctx, cancel := context.WithCancel(context.Background())
	resultCh := make(chan domain.SMTPAttemptResult, 1)
	go func() {
		resultCh <- NewSender(nil).SendSMTP(ctx, d, func(_ []domain.NotificationRecipient) error { return nil })
	}()
	conn := <-accepted
	defer conn.Close()
	cancel()
	select {
	case result := <-resultCh:
		if !strings.Contains(result.Diagnostic, "greeting") {
			t.Fatalf("diagnostic %+v", result)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("SMTP ignored cancellation")
	}
	<-done
}
func TestSMTPPlainTextAndNonTaskMessages(t *testing.T) {
	for _, event := range []domain.WebhookEvent{{Event: "test"}, {Event: "daemon.health", Daemon: domain.WebhookDaemon{Status: "healthy"}}} {
		payload, err := json.Marshal(event)
		if err != nil {
			t.Fatal(err)
		}
		cfg := smtpTestConfig()
		data, err := RenderMail(domain.NotificationDelivery{ID: "message", Payload: payload}, cfg)
		if err != nil {
			t.Fatal(err)
		}
		msg, err := mail.ReadMessage(strings.NewReader(string(data)))
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(quotedprintable.NewReader(msg.Body))
		if err != nil || !strings.Contains(string(body), "Delivery: message") {
			t.Fatalf("message %v %s", err, body)
		}
	}
}

func TestSMTPQuotedMailboxesPreserveEnvelopeAndHeaders(t *testing.T) {
	receiver, err := testsmtp.Start(testsmtp.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer receiver.Close()
	d := mailFixtureDelivery(t, receiver, "plain")
	raw := `"ops,night"@example.test`
	d.SMTP.From = raw
	d.SMTP.Recipients = []string{raw, `"ops team"@example.test`}
	d.Recipients = nil
	cfg, err := ValidateSMTP(*d.SMTP)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.From != raw || cfg.Recipients[0] != raw {
		t.Fatalf("quoting lost %+v", cfg)
	}
	d.SMTP = &cfg
	result := NewSender(nil).SendSMTP(context.Background(), d, func([]domain.NotificationRecipient) error { return nil })
	for _, r := range result.Recipients {
		if r.State != "accepted" {
			t.Fatalf("envelope rejected %+v", result)
		}
	}
	messages := receiver.Messages()
	if len(messages) != 1 {
		t.Fatalf("messages %d", len(messages))
	}
	message, err := mail.ReadMessage(strings.NewReader(string(messages[0].Data)))
	if err != nil {
		t.Fatal(err)
	}
	from, err := message.Header.AddressList("From")
	if err != nil || len(from) != 1 || from[0].Address != "ops,night@example.test" {
		t.Fatalf("sender %+v %v", from, err)
	}
	recipients, err := message.Header.AddressList("To")
	if err != nil || len(recipients) != 2 || recipients[0].Address != "ops,night@example.test" || recipients[1].Address != "ops team@example.test" {
		t.Fatalf("recipients %+v %v", recipients, err)
	}
	if messages[0].From != "ops,night@example.test" || len(messages[0].Recipients) != 2 || messages[0].Recipients[1] != "ops team@example.test" {
		t.Fatalf("altered envelope %+v", messages[0])
	}
}
