package clievaluate

import (
	"context"
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/notification"
	"github.com/shruggietech/go-schedule/internal/store"
	"github.com/shruggietech/go-schedule/internal/testsmtp"
)

func smtpTransportPolicy(ctx context.Context, redactor *Redactor) execution {
	user, password := uuid.NewString(), uuid.NewString()
	redactor.Register(user, password)
	observations := []NamedValue{}
	good := true
	for _, mode := range []string{"plain", "starttls", "tls"} {
		receiver, err := testsmtp.Start(testsmtp.Options{TLSMode: mode, Username: user, Password: password, ReplyText: password, Reject: map[string]int{"temporary@example.test": 450, "rejected@example.test": 550}})
		if err != nil {
			return policyResult(observations, false, err)
		}
		host, portRaw, err := net.SplitHostPort(receiver.Addr())
		if err != nil {
			_ = receiver.Close()
			return policyResult(observations, false, err)
		}
		port, err := strconv.Atoi(portRaw)
		if err != nil {
			_ = receiver.Close()
			return policyResult(observations, false, err)
		}
		config := domain.SMTPConfig{Host: host, Port: port, TLSMode: mode, From: "sender@example.test", Recipients: []string{"accepted@example.test", "temporary@example.test", "rejected@example.test"}, TimeoutSeconds: 2, HTML: true}
		auth, err := notification.EncodeSMTPAuth(domain.SMTPAuth{Username: user, Password: password})
		if err != nil {
			_ = receiver.Close()
			return policyResult(observations, false, err)
		}
		payload, err := json.Marshal(domain.WebhookEvent{Event: "run.completed", Task: &domain.WebhookTask{Name: "QA 世界 <script>unsafe markup</script>"}, Run: &domain.WebhookRun{Outcome: domain.OutcomeFailure}})
		if err != nil {
			_ = receiver.Close()
			return policyResult(observations, false, err)
		}
		delivery := domain.NotificationDelivery{ID: "owned-smtp", Kind: domain.NotificationChannelSMTP, SMTP: &config, Authorization: auth, Payload: payload, CreatedAt: time.Now().UTC()}
		checkpoints := 0
		checkpoint := func(_ []domain.NotificationRecipient) error { checkpoints++; return nil }
		sender := notification.NewSenderWithSMTPRoots(receiver.Roots())
		result := sender.SendSMTP(ctx, delivery, checkpoint)
		caseOK := len(result.Recipients) == 3 && result.Recipients[0].State == "accepted" && result.Recipients[1].State == "temporary_failure" && result.Recipients[2].State == "permanent_failure" && checkpoints == 1 && !strings.Contains(result.Diagnostic, password)
		if caseOK {
			receiver.SetRecipientCode("temporary@example.test", 0)
			delivery.Recipients = result.Recipients
			retried := sender.SendSMTP(ctx, delivery, checkpoint)
			rcpts := strings.Join(receiver.Recipients(), ",")
			caseOK = len(retried.Recipients) == 3 && retried.Recipients[0].State == "accepted" && retried.Recipients[1].State == "accepted" && strings.Count(rcpts, "accepted@example.test") == 1 && strings.Count(rcpts, "rejected@example.test") == 1 && strings.Count(rcpts, "temporary@example.test") == 2
		}
		if mode != "plain" {
			rejected := notification.NewSender(nil).SendSMTP(ctx, delivery, checkpoint)
			caseOK = caseOK && len(rejected.Recipients) == 3 && rejected.Recipients[1].State != "accepted"
		}
		good = good && caseOK
		observations = append(observations, NamedValue{Name: "SMTP " + mode + " TLS/authentication/partial recipients/retry exclusion", Value: strconv.FormatBool(caseOK)})
		if err := receiver.Close(); err != nil {
			return policyResult(observations, false, err)
		}
	}
	persisted, err := smtpRecoveryPolicy()
	good = good && persisted
	observations = append(observations, NamedValue{Name: "accepted recipient retained and interrupted submission identified after restart", Value: strconv.FormatBool(persisted)})
	result := policyResult(observations, good, err)
	result.Logs = mustObservation(observations)
	return result
}
func smtpRecoveryPolicy() (bool, error) {
	st, err := store.Open(":memory:")
	if err != nil {
		return false, err
	}
	defer st.Close()
	config := &domain.SMTPConfig{Host: "127.0.0.1", Port: 25, TLSMode: "plain", From: "sender@example.test", Recipients: []string{"accepted@example.test", "uncertain@example.test"}, TimeoutSeconds: 1}
	c := domain.NotificationChannel{Name: "SMTP recovery", Kind: domain.NotificationChannelSMTP, SMTP: config, Enabled: true}
	if err := st.CreateNotificationChannel(&c); err != nil {
		return false, err
	}
	if _, err := st.CreateTestNotificationDelivery(c.ID); err != nil {
		return false, err
	}
	claimed, err := st.ClaimNotificationDeliveries(1, time.Now().Add(time.Minute))
	if err != nil || len(claimed) != 1 {
		return false, err
	}
	recipients := claimed[0].Recipients
	recipients[0].State = "accepted"
	recipients[1].State = "submitting"
	if err := st.CheckpointSMTPRecipients(claimed[0].ID, recipients); err != nil {
		return false, err
	}
	if _, err := st.RecoverNotificationDeliveries(time.Now()); err != nil {
		return false, err
	}
	reclaimed, err := st.ClaimNotificationDeliveries(1, time.Now().Add(time.Minute))
	if err != nil || len(reclaimed) != 1 {
		return false, err
	}
	return len(reclaimed[0].Recipients) == 2 && reclaimed[0].Recipients[0].State == "accepted" && reclaimed[0].Recipients[1].State == "uncertain", nil
}
