package notification

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

var smtpHostLabel = regexp.MustCompile("^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$")

// SMTPConfigError names an invalid authoring field without echoing its value.
type SMTPConfigError struct{ Field, Message string }

func (e *SMTPConfigError) Error() string { return "smtp." + e.Field + ": " + e.Message }
func smtpInvalid(field, message string) error {
	return &SMTPConfigError{Field: field, Message: message}
}

// ValidateSMTP normalizes nonsecret settings and rejects command/header injection.
func ValidateSMTP(c domain.SMTPConfig) (domain.SMTPConfig, error) {
	c.Host = strings.TrimSpace(c.Host)
	if !validSMTPHost(c.Host) {
		return c, smtpInvalid("host", "use a DNS hostname or IP address without a URL scheme or port")
	}
	switch c.TLSMode {
	case "starttls", "tls", "plain":
	default:
		return c, smtpInvalid("tls_mode", "choose starttls (required), tls (implicit TLS), or plain (unencrypted relay)")
	}
	if c.Port == 0 {
		switch c.TLSMode {
		case "tls":
			c.Port = 465
		case "plain":
			c.Port = 25
		default:
			c.Port = 587
		}
	}
	if c.Port < 1 || c.Port > 65535 {
		return c, smtpInvalid("port", "use a port between 1 and 65535")
	}
	var err error
	c.From, err = mailbox(c.From)
	if err != nil {
		return c, smtpInvalid("from", "use one valid ASCII email address without control characters")
	}
	if len(c.Recipients) < 1 || len(c.Recipients) > 100 {
		return c, smtpInvalid("recipients", "provide 1 to 100 email addresses")
	}
	recipients := make([]string, 0, len(c.Recipients))
	seen := map[string]bool{}
	for _, raw := range c.Recipients {
		addr, err := mailbox(raw)
		if err != nil {
			return c, smtpInvalid("recipients", "use valid ASCII email addresses without control characters")
		}
		if !seen[addr] {
			seen[addr] = true
			recipients = append(recipients, addr)
		}
	}
	c.Recipients = recipients
	if c.TimeoutSeconds == 0 {
		c.TimeoutSeconds = 10
	}
	if c.TimeoutSeconds < 1 || c.TimeoutSeconds > 60 {
		return c, smtpInvalid("timeout_seconds", "use 1 to 60 seconds")
	}
	return c, nil
}

func validSMTPHost(host string) bool {
	if len(host) == 0 || len(host) > 253 || strings.ContainsAny(host, "\r\n\x00%") {
		return false
	}
	if net.ParseIP(host) != nil {
		return true
	}
	for _, label := range strings.Split(strings.TrimSuffix(host, "."), ".") {
		if len(label) > 63 || !smtpHostLabel.MatchString(label) {
			return false
		}
	}
	return true
}
func mailbox(raw string) (string, error) {
	if len(raw) > 320 || strings.IndexFunc(raw, func(r rune) bool { return r < 32 || r >= 127 }) >= 0 {
		return "", fmt.Errorf("invalid mailbox")
	}
	addr, err := mail.ParseAddress(strings.TrimSpace(raw))
	if err != nil || len(addr.Address) > 254 || !strings.Contains(addr.Address, "@") {
		return "", fmt.Errorf("invalid mailbox")
	}
	canonical := (&mail.Address{Address: addr.Address}).String()
	return strings.TrimSuffix(strings.TrimPrefix(canonical, "<"), ">"), nil
}

// EncodeSMTPAuth validates write-only credentials for the existing protected slot.
func EncodeSMTPAuth(auth domain.SMTPAuth) (string, error) {
	if auth.Username == "" && auth.Password == "" {
		return "", nil
	}
	if auth.Username == "" || auth.Password == "" || len(auth.Username) > 1024 || len(auth.Password) > 4096 || strings.IndexFunc(auth.Username+auth.Password, unicode.IsControl) >= 0 {
		return "", smtpInvalid("auth", "provide both username and password without control characters, or clear both")
	}
	raw, err := json.Marshal(auth)
	if err != nil {
		return "", smtpInvalid("auth", "could not encode credentials")
	}
	return string(raw), nil
}

// ParseSMTPAuth decodes protected credentials without exposing malformed input.
func ParseSMTPAuth(raw string) (domain.SMTPAuth, error) {
	var auth domain.SMTPAuth
	if raw == "" {
		return auth, nil
	}
	if len(raw) > 8192 {
		return auth, smtpInvalid("auth", "credential input is too large")
	}
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&auth); err != nil {
		return domain.SMTPAuth{}, smtpInvalid("auth", "provide a JSON object with username and password")
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return domain.SMTPAuth{}, smtpInvalid("auth", "provide one JSON credential object")
	}
	if _, err := EncodeSMTPAuth(auth); err != nil {
		return domain.SMTPAuth{}, err
	}
	return auth, nil
}

// SMTPSummary identifies the destination without publishing credentials.
func SMTPSummary(c domain.SMTPConfig) string {
	return "smtp+" + c.TLSMode + "://" + net.JoinHostPort(c.Host, fmt.Sprint(c.Port))
}
