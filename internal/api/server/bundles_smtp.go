package server

import (
	"fmt"

	"github.com/shruggietech/go-schedule/internal/bundle"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/notification"
)

func smtpSameAuthority(a, b *domain.SMTPConfig) bool {
	return a != nil && b != nil && a.Host == b.Host && a.Port == b.Port && a.TLSMode == b.TLSMode
}
func (s *Server) bindBundleSMTP(source bundle.SMTPChannel, bindings map[string]bundleChannelBinding) (bundle.Action, string, error) {
	channels, err := s.store.ListNotificationChannels()
	if err != nil {
		return "", "", err
	}
	var matches []domain.NotificationChannel
	for _, channel := range channels {
		if channel.Name == source.Name {
			matches = append(matches, channel)
		}
	}
	if len(matches) != 1 || matches[0].Kind != domain.NotificationChannelSMTP {
		return bundle.ActionConflict, "Provision exactly one target-local SMTP channel with this name, then preview again.", nil
	}
	c := matches[0]
	if c.HasAuthorization != source.RequiresAuthorization {
		return bundle.ActionConflict, "Provision or explicitly clear this target-local channel's credentials before preview.", nil
	}
	if c.HasAuthorization && !smtpSameAuthority(c.SMTP, &source.Settings) {
		return bundle.ActionConflict, "Configure the target-local mail server and credentials before preview; bundles do not redirect stored credentials.", nil
	}
	bindings[source.Name] = bundleChannelBinding{ID: c.ID, UpdatedAt: c.UpdatedAt}
	return "", "", nil
}
func (s *Server) validateBundleChannels(bindings map[string]bundleChannelBinding) error {
	channels, err := s.store.ListNotificationChannels()
	if err != nil {
		return err
	}
	counts := map[string]int{}
	for _, channel := range channels {
		counts[channel.Name]++
	}
	for name, binding := range bindings {
		channel, err := s.store.GetNotificationChannel(binding.ID)
		if err != nil || counts[name] != 1 || channel.Name != name || !channel.UpdatedAt.Equal(binding.UpdatedAt) {
			return fmt.Errorf("reviewed notification channel changed after preview; preview again")
		}
	}
	return nil
}
func (s *Server) applyBundleSMTP(source bundle.SMTPChannel, bindings map[string]bundleChannelBinding) bundle.Item {
	item := bundle.Item{Kind: "smtp_channel", PortableID: source.Name, Name: source.Name}
	binding, ok := bindings[source.Name]
	if !ok {
		item.Action, item.Message = bundle.ActionFailed, "reviewed SMTP binding is missing"
		return item
	}
	c, err := s.store.GetNotificationChannel(binding.ID)
	if err != nil || !c.UpdatedAt.Equal(binding.UpdatedAt) || c.Name != source.Name || c.Kind != domain.NotificationChannelSMTP || c.HasAuthorization != source.RequiresAuthorization || (c.HasAuthorization && !smtpSameAuthority(c.SMTP, &source.Settings)) {
		item.Action, item.Message = bundle.ActionFailed, "reviewed SMTP channel changed; preview again"
		return item
	}
	config, err := notification.ValidateSMTP(source.Settings)
	if err != nil {
		item.Action, item.Message = bundle.ActionFailed, "SMTP settings are invalid"
		return item
	}
	c.SMTP = &config
	c.EndpointSummary = notification.SMTPSummary(config)
	c.Enabled = source.Enabled
	c.HealthIntervalSeconds = source.HealthIntervalSeconds
	if err := s.store.UpdateNotificationChannel(c); err != nil {
		item.Action, item.Message = bundle.ActionFailed, "SMTP channel settings could not be saved"
		return item
	}
	updated, err := s.store.GetNotificationChannel(c.ID)
	if err != nil {
		item.Action, item.Message = bundle.ActionUncertain, "SMTP settings were saved but the new revision could not be read"
		return item
	}
	bindings[source.Name] = bundleChannelBinding{ID: c.ID, UpdatedAt: updated.UpdatedAt}
	item.Action, item.Message = bundle.ActionApplied, "Updated nonsecret mail settings; target-local credentials retained."
	return item
}
