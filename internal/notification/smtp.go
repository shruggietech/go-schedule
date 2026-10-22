package notification

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"slices"
	"strings"
	"time"

	"github.com/emersion/go-sasl"
	smtp "github.com/emersion/go-smtp"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// NewSenderWithSMTPRoots injects explicit trust roots for owned SMTP environments.
// A nil pool uses system trust; certificate and hostname verification remain enabled.
func NewSenderWithSMTPRoots(roots *x509.CertPool) *Sender {
	sender := NewSender(nil)
	sender.smtpRoots = roots
	return sender
}

// SendSMTP performs one bounded SMTP attempt, checkpointing the DATA subset before transmission.
func (s *Sender) SendSMTP(ctx context.Context, d domain.NotificationDelivery, checkpoint func([]domain.NotificationRecipient) error) domain.SMTPAttemptResult {
	result := domain.SMTPAttemptResult{Recipients: slices.Clone(d.Recipients)}
	if d.SMTP == nil {
		return smtpFail(result, nil, "configuration", fmt.Errorf("missing settings"), false)
	}
	config, err := ValidateSMTP(*d.SMTP)
	if err != nil {
		return smtpFail(result, nil, "configuration", err, false)
	}
	if len(result.Recipients) == 0 {
		result.Recipients = make([]domain.NotificationRecipient, len(config.Recipients))
		for i, address := range config.Recipients {
			result.Recipients[i] = domain.NotificationRecipient{Address: address, State: "pending"}
		}
	}
	eligible := []int{}
	for i, r := range result.Recipients {
		if r.State == "pending" || r.State == "temporary_failure" || r.State == "uncertain" {
			eligible = append(eligible, i)
		}
	}
	if len(eligible) == 0 {
		return smtpSummary(result)
	}
	auth, err := ParseSMTPAuth(d.Authorization)
	if err != nil {
		return smtpFail(result, eligible, "credentials", err, false)
	}
	message, err := RenderMail(d, config)
	if err != nil {
		return smtpFail(result, eligible, "message", err, false)
	}
	attempt, cancel := context.WithTimeout(ctx, time.Duration(config.TimeoutSeconds)*time.Second)
	defer cancel()
	fail := func(out domain.SMTPAttemptResult, indices []int, phase string, err error, retry bool) domain.SMTPAttemptResult {
		if cause := attempt.Err(); cause != nil {
			err = cause
		}
		return smtpFail(out, indices, phase, err, retry)
	}
	raw, err := (&net.Dialer{}).DialContext(attempt, "tcp", net.JoinHostPort(config.Host, fmt.Sprint(config.Port)))
	if err != nil {
		return fail(result, eligible, "connection", err, true)
	}
	stop := context.AfterFunc(attempt, func() { _ = raw.Close() })
	defer stop()
	defer raw.Close()
	tlsConfig := &tls.Config{ServerName: config.Host, MinVersion: tls.VersionTLS12, RootCAs: s.smtpRoots}
	var client *smtp.Client
	switch config.TLSMode {
	case "tls":
		encrypted := tls.Client(raw, tlsConfig)
		if err := encrypted.HandshakeContext(attempt); err != nil {
			return fail(result, eligible, "TLS", err, smtpTransient(err))
		}
		client = smtp.NewClient(encrypted)
	case "starttls":
		client, err = smtp.NewClientStartTLS(raw, tlsConfig)
		if err != nil {
			return fail(result, eligible, "TLS", err, smtpTransient(err))
		}
	default:
		client = smtp.NewClient(raw)
	}
	defer client.Close()
	client.CommandTimeout = time.Duration(config.TimeoutSeconds) * time.Second
	client.SubmissionTimeout = client.CommandTimeout
	if err := client.Hello("localhost"); err != nil {
		return fail(result, eligible, "greeting", err, smtpTransient(err))
	}
	if auth.Username != "" {
		supported, mechanisms := client.Extension("AUTH")
		if !supported {
			return fail(result, eligible, "authentication", fmt.Errorf("AUTH unavailable"), false)
		}
		available := strings.Fields(strings.ToUpper(mechanisms))
		var mechanism sasl.Client
		switch {
		case slices.Contains(available, "PLAIN"):
			mechanism = sasl.NewPlainClient("", auth.Username, auth.Password)
		case slices.Contains(available, "LOGIN"):
			mechanism = sasl.NewLoginClient(auth.Username, auth.Password)
		default:
			return fail(result, eligible, "authentication", fmt.Errorf("supported AUTH unavailable"), false)
		}
		if err := client.Auth(mechanism); err != nil {
			return fail(result, eligible, "authentication", err, smtpTransient(err))
		}
	}
	if err := client.Mail(config.From, nil); err != nil {
		return fail(result, eligible, "sender", err, smtpTransient(err))
	}
	accepted := []int{}
	for _, i := range eligible {
		if err := client.Rcpt(result.Recipients[i].Address, nil); err != nil {
			var code *smtp.SMTPError
			if !errors.As(err, &code) {
				return fail(result, eligible, "recipient", err, true)
			}
			failed := fail(domain.SMTPAttemptResult{Recipients: result.Recipients}, []int{i}, "recipient", err, code.Temporary())
			result.Recipients = failed.Recipients
			continue
		}
		accepted = append(accepted, i)
	}
	if len(accepted) == 0 {
		return smtpSummary(result)
	}
	for _, i := range accepted {
		result.Recipients[i].State = "submitting"
		result.Recipients[i].Code = 0
		result.Recipients[i].Diagnostic = ""
	}
	if checkpoint == nil {
		return fail(result, accepted, "checkpoint", fmt.Errorf("missing checkpoint"), false)
	}
	if err := checkpoint(slices.Clone(result.Recipients)); err != nil {
		return fail(result, accepted, "checkpoint", err, true)
	}
	writer, err := client.Data()
	if err != nil {
		return fail(result, accepted, "message submission", err, smtpTransient(err))
	}
	if _, err := writer.Write(message); err != nil {
		return smtpUncertain(result, accepted)
	}
	if err := writer.Close(); err != nil {
		var code *smtp.SMTPError
		if errors.As(err, &code) {
			return fail(result, accepted, "message submission", err, code.Temporary())
		}
		return smtpUncertain(result, accepted)
	}
	for _, i := range accepted {
		result.Recipients[i].State = "accepted"
		result.Recipients[i].Code = 250
		result.Recipients[i].Diagnostic = "Mail server accepted the message."
	}
	// The final DATA acknowledgement is authoritative; QUIT is only cleanup.
	_ = client.Quit()
	return smtpSummary(result)
}
func smtpTransient(err error) bool {
	var code *smtp.SMTPError
	if errors.As(err, &code) {
		return code.Temporary()
	}
	var certificate *tls.CertificateVerificationError
	if errors.As(err, &certificate) {
		return false
	}
	var invalid x509.CertificateInvalidError
	if errors.As(err, &invalid) {
		return false
	}
	var unknown x509.UnknownAuthorityError
	return !errors.As(err, &unknown)
}
func smtpFail(result domain.SMTPAttemptResult, indices []int, phase string, err error, retry bool) domain.SMTPAttemptResult {
	var verification *tls.CertificateVerificationError
	if errors.As(err, &verification) {
		phase = "TLS"
	}
	if phase == "configuration" && indices == nil {
		for i := range result.Recipients {
			indices = append(indices, i)
		}
	}
	code := 0
	var smtpErr *smtp.SMTPError
	if errors.As(err, &smtpErr) {
		code = smtpErr.Code
	}
	diagnostic := "SMTP " + phase + " failed."
	switch phase {
	case "configuration":
		diagnostic += " Review the SMTP channel settings."
	case "credentials", "authentication":
		diagnostic += " Check the username, password and server authentication support."
	case "TLS":
		diagnostic += " Check the encryption mode, port, hostname and trusted server certificate; no downgrade was attempted."
	case "connection", "greeting":
		diagnostic += " Check the server address, port and connectivity."
	case "recipient":
		diagnostic += " Check the recipient address or ask the mail administrator about rejection."
	case "sender":
		diagnostic += " Check the sender address and relay permissions."
	case "checkpoint":
		diagnostic += " Recipient progress could not be recorded; inspect daemon storage before retrying."
	default:
		diagnostic += " Check the mail service and channel configuration."
	}
	if code != 0 {
		diagnostic += fmt.Sprintf(" Server status: %d.", code)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		diagnostic += " The attempt timed out."
	}
	if errors.Is(err, context.Canceled) {
		diagnostic += " The attempt was canceled."
	}
	state := "permanent_failure"
	if retry {
		state = "temporary_failure"
	}
	for _, i := range indices {
		if result.Recipients[i].State == "accepted" || result.Recipients[i].State == "permanent_failure" {
			continue
		}
		result.Recipients[i].State = state
		result.Recipients[i].Code = code
		result.Recipients[i].Diagnostic = diagnostic
	}
	result.Status = code
	result.Diagnostic = diagnostic
	return result
}
func smtpUncertain(result domain.SMTPAttemptResult, indices []int) domain.SMTPAttemptResult {
	for _, i := range indices {
		result.Recipients[i].State = "uncertain"
		result.Recipients[i].Code = 0
		result.Recipients[i].Diagnostic = "The connection ended before submission was confirmed; retry may duplicate email."
	}
	result.Diagnostic = "Email submission is uncertain; a bounded retry may duplicate email. Check delivery history and the mail server."
	return result
}
func smtpSummary(result domain.SMTPAttemptResult) domain.SMTPAttemptResult {
	accepted, failed, pending := 0, 0, 0
	for _, r := range result.Recipients {
		switch r.State {
		case "accepted":
			accepted++
		case "permanent_failure":
			failed++
		default:
			pending++
		}
		if r.Code > result.Status {
			result.Status = r.Code
		}
	}
	result.Diagnostic = fmt.Sprintf("Mail server accepted %d recipient(s); %d permanently rejected; %d awaiting retry.", accepted, failed, pending)
	return result
}
