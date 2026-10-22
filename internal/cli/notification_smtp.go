package cli

import (
	"fmt"
	"io"
	"net/mail"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/notification"
)

type smtpFlags struct {
	config     domain.SMTPConfig
	timeout    time.Duration
	authFile   string
	recipients []string
}

func (f *smtpFlags) bind(cmd *cobra.Command, auth bool) {
	cmd.Flags().StringVar(&f.config.Host, "smtp-host", "", "SMTP server hostname or IP address")
	cmd.Flags().IntVar(&f.config.Port, "smtp-port", 0, "SMTP port (default follows the encryption mode)")
	cmd.Flags().StringVar(&f.config.TLSMode, "smtp-tls", "starttls", "SMTP encryption: starttls (required), tls (implicit), or plain (unencrypted)")
	cmd.Flags().StringVar(&f.config.From, "smtp-from", "", "sender email address")
	cmd.Flags().StringArrayVar(&f.recipients, "smtp-to", nil, "recipient email address (repeatable or comma-separated, at most 100)")
	cmd.Flags().DurationVar(&f.timeout, "smtp-timeout", 10*time.Second, "shared SMTP attempt timeout (whole seconds, 1s-60s)")
	cmd.Flags().BoolVar(&f.config.HTML, "smtp-html", false, "include simple escaped HTML alongside plain text")
	if auth {
		cmd.Flags().StringVar(&f.authFile, "smtp-auth-file", "", "write-only SMTP username/password JSON file; - reads standard input")
	}
}
func smtpSettingsChanged(cmd *cobra.Command) bool {
	for _, name := range []string{"smtp-host", "smtp-port", "smtp-tls", "smtp-from", "smtp-to", "smtp-timeout", "smtp-html"} {
		if cmd.Flags().Changed(name) {
			return true
		}
	}
	return false
}
func (f *smtpFlags) settings(cmd *cobra.Command, current *domain.SMTPConfig) (*domain.SMTPConfig, error) {
	c := f.config
	if current != nil {
		c = *current
		if cmd.Flags().Changed("smtp-host") {
			c.Host = f.config.Host
		}
		if cmd.Flags().Changed("smtp-port") {
			c.Port = f.config.Port
		}
		if cmd.Flags().Changed("smtp-tls") {
			c.TLSMode = f.config.TLSMode
		}
		if cmd.Flags().Changed("smtp-from") {
			c.From = f.config.From
		}
		if cmd.Flags().Changed("smtp-html") {
			c.HTML = f.config.HTML
		}
	}
	if current == nil || cmd.Flags().Changed("smtp-to") {
		c.Recipients = nil
		for _, raw := range f.recipients {
			if len(raw) > 32000 || strings.IndexFunc(raw, func(r rune) bool { return r < 32 || r >= 127 }) >= 0 {
				return nil, fmtUsage("--smtp-to requires ASCII mailboxes without control characters")
			}
			addresses, err := mail.ParseAddressList(raw)
			if err != nil {
				return nil, fmtUsage("--smtp-to requires valid email addresses, separated by commas or repeated flags")
			}
			for _, address := range addresses {
				c.Recipients = append(c.Recipients, address.String())
			}
		}
	}
	if current == nil || cmd.Flags().Changed("smtp-timeout") {
		if f.timeout%time.Second != 0 {
			return nil, fmtUsage("--smtp-timeout requires whole seconds between 1s and 60s")
		}
		c.TimeoutSeconds = int(f.timeout / time.Second)
		if c.TimeoutSeconds < 1 || c.TimeoutSeconds > 60 {
			return nil, fmtUsage("--smtp-timeout must be between 1s and 60s")
		}
	}
	normalized, err := notification.ValidateSMTP(c)
	if err != nil {
		return nil, fmtUsage(err.Error())
	}
	return &normalized, nil
}
func readSMTPAuth(cmd *cobra.Command, filename string) (string, error) {
	reader := cmd.InOrStdin()
	if filename != "-" {
		file, err := os.Open(filename)
		if err != nil {
			return "", fmtUsage("cannot read --smtp-auth-file; check the path and file permissions")
		}
		defer file.Close()
		reader = file
	}
	data, err := io.ReadAll(io.LimitReader(reader, 8193))
	if err != nil {
		return "", fmtUsage("cannot read --smtp-auth-file")
	}
	if !utf8.Valid(data) {
		return "", fmtUsage("--smtp-auth-file must contain UTF-8 JSON")
	}
	if len(data) > 8192 {
		return "", fmtUsage("--smtp-auth-file must be at most 8192 bytes")
	}
	auth, err := notification.ParseSMTPAuth(string(data))
	if err != nil {
		return "", fmtUsage(err.Error())
	}
	value, err := notification.EncodeSMTPAuth(auth)
	if err != nil {
		return "", fmt.Errorf("encode SMTP credentials: %w", err)
	}
	return value, nil
}
