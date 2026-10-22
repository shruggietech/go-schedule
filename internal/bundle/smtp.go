package bundle

import (
	"strings"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/notification"
)

// SMTPChannel carries nonsecret mail settings bound to one target-local channel by name.
type SMTPChannel struct {
	Name                  string            `json:"name"`
	Settings              domain.SMTPConfig `json:"settings"`
	Enabled               bool              `json:"enabled"`
	HealthIntervalSeconds int64             `json:"health_interval_seconds"`
	RequiresAuthorization bool              `json:"requires_authorization"`
}

func validateSMTPChannels(doc Document) []Issue {
	var issues []Issue
	seen := map[string]bool{}
	for _, channel := range doc.SMTPChannels {
		if strings.TrimSpace(channel.Name) == "" || len(channel.Name) > 256 || seen[channel.Name] {
			issues = append(issues, Issue{Kind: "smtp_channel", Identity: channel.Name, Message: "SMTP channel names must be nonempty, unique and at most 256 bytes"})
		}
		seen[channel.Name] = true
		if _, err := notification.ValidateSMTP(channel.Settings); err != nil {
			issues = append(issues, Issue{Kind: "smtp_channel", Identity: channel.Name, Message: err.Error()})
		}
		if channel.HealthIntervalSeconds != 0 && (channel.HealthIntervalSeconds < 60 || channel.HealthIntervalSeconds > 86400) {
			issues = append(issues, Issue{Kind: "smtp_channel", Identity: channel.Name, Message: "SMTP health interval must be 0 or 60-86400 seconds"})
		}
	}
	return issues
}
