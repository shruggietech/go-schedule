package clievaluate

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/testsmtp"
)

func (s *session) startSMTP() error {
	username, password, rotated := uuid.NewString(), uuid.NewString(), uuid.NewString()
	s.redactor.Register(username, password, rotated)
	s.vars["smtp_username"], s.vars["smtp_password"], s.vars["smtp_rotated"] = username, password, rotated
	receiver, err := testsmtp.Start(testsmtp.Options{Username: username, Password: password})
	if err != nil {
		return err
	}
	s.smtp = receiver
	host, port, err := net.SplitHostPort(receiver.Addr())
	if err != nil {
		return err
	}
	s.vars["smtp_host"], s.vars["smtp_port"] = host, port
	for _, v := range []struct{ variable, password string }{{"smtp_auth_file", password}, {"smtp_rotated_file", rotated}} {
		file := filepath.Join(s.w.Runtime, "tasks", v.variable+".json")
		data, err := json.Marshal(domain.SMTPAuth{Username: username, Password: v.password})
		if err != nil {
			return err
		}
		if err := s.writeOwned(file, data); err != nil {
			return err
		}
		s.vars[v.variable] = file
	}
	return nil
}
func (s *session) smtpBefore(ctx context.Context, definition step) error {
	if s.smtp == nil {
		return fmt.Errorf("owned SMTP fixture is unavailable")
	}
	s.smtpBeforeCount = len(s.smtp.Messages())
	switch definition.Spec.ID {
	case "smtp.rotate":
		s.smtp.SetCredentials(s.vars["smtp_username"], s.vars["smtp_rotated"])
	case "smtp.clear":
		s.smtp.SetCredentials("", "")
	case "smtp.task-delivery":
		_, err := s.api.ReplaceTaskNotificationAssignments(ctx, s.vars["base_task"], server.NotificationAssignmentsRequest{Assignments: []server.NotificationAssignmentInput{{ChannelID: s.vars["smtp_channel"], OnSuccess: true}, {ChannelID: s.vars["second_channel"], OnSuccess: true}}})
		return err
	}
	return nil
}
func (s *session) smtpPost(ctx context.Context, definition step, result execution) (bool, any, error) {
	if definition.Spec.ExpectedExit != 0 {
		return result.Code == definition.Spec.ExpectedExit, map[string]any{"negative_input_preserved": result.Code == definition.Spec.ExpectedExit}, nil
	}
	if definition.Spec.ID == "smtp.remove" {
		_, err := s.api.GetNotificationChannel(ctx, s.vars["smtp_channel"])
		return err != nil, map[string]bool{"channel_absent": err != nil}, nil
	}
	if definition.Hook == "smtp-delivery" || definition.Hook == "run-smtp" {
		var observed domain.NotificationDelivery
		err := waitUntil(ctx, func() bool {
			deliveries, err := s.api.ListNotificationDeliveries(ctx, notificationFilter(s.vars["smtp_channel"]))
			if err != nil {
				return false
			}
			for _, delivery := range deliveries {
				match := strings.Contains(string(result.Out), delivery.ID)
				if definition.Hook == "run-smtp" {
					match = delivery.EventKind == domain.NotificationEventRunCompleted && !s.priorRuns[delivery.RunID]
				}
				if match && (delivery.State == domain.NotificationDeliverySucceeded || delivery.State == domain.NotificationDeliveryFailed) {
					observed = delivery
					return true
				}
			}
			return false
		})
		if err != nil {
			return false, observed, err
		}
		accepted := observed.State == domain.NotificationDeliverySucceeded && len(observed.Recipients) == 2
		for _, r := range observed.Recipients {
			accepted = accepted && r.State == "accepted"
		}
		messages := s.smtp.Messages()
		matched := false
		for _, m := range messages[s.smtpBeforeCount:] {
			matched = matched || strings.Contains(string(m.Data), observed.ID) && len(m.Recipients) == 2
		}
		if definition.Hook == "run-smtp" {
			taskOK, run, runErr := s.observeRun(ctx, "run-success")
			if runErr != nil {
				return false, run, runErr
			}
			accepted = accepted && taskOK && strings.Contains(string(result.Out), "triggered run for task "+s.vars["base_task"])
			webhookOK := waitUntil(ctx, func() bool {
				deliveries, err := s.api.ListNotificationDeliveries(ctx, notificationFilter(s.vars["second_channel"]))
				if err != nil {
					return false
				}
				for _, delivery := range deliveries {
					if delivery.RunID == observed.RunID && delivery.State == domain.NotificationDeliverySucceeded {
						return true
					}
				}
				return false
			})
			accepted = accepted && webhookOK == nil
		}
		return accepted && matched, map[string]any{"delivery": observed, "receiver_message_count": len(messages) - s.smtpBeforeCount, "correlated_message": matched}, nil
	}
	channel, err := s.api.GetNotificationChannel(ctx, s.vars["smtp_channel"])
	if err != nil {
		return false, channel, err
	}
	good := channel.Kind == domain.NotificationChannelSMTP && channel.SMTP != nil
	switch definition.Spec.ID {
	case "smtp.create":
		good = good && !channel.Enabled && channel.HasAuthorization && channel.SMTP.HTML
	case "smtp.update":
		good = good && channel.SMTP.TimeoutSeconds == 3 && !channel.SMTP.HTML && len(channel.SMTP.Recipients) == 2
	case "smtp.enable":
		good = good && channel.Enabled
	case "smtp.disable":
		good = good && !channel.Enabled
	case "smtp.rotate":
		good = good && channel.HasAuthorization
	case "smtp.clear":
		good = good && !channel.HasAuthorization
	}
	return good, channel, nil
}
