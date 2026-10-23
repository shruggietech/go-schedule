package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func encodeSMTPConfig(c *domain.SMTPConfig) string {
	if c == nil {
		return ""
	}
	// SMTPConfig contains only JSON-supported scalar values and slices.
	data, _ := json.Marshal(c)
	return string(data)
}
func decodeSMTPConfig(raw string) (*domain.SMTPConfig, error) {
	if raw == "" {
		return nil, nil
	}
	var c domain.SMTPConfig
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		return nil, fmt.Errorf("store: decode SMTP settings")
	}
	return &c, nil
}
func initialSMTPRecipients(c *domain.SMTPConfig) []domain.NotificationRecipient {
	if c == nil {
		return nil
	}
	out := make([]domain.NotificationRecipient, len(c.Recipients))
	for i, address := range c.Recipients {
		out[i] = domain.NotificationRecipient{Address: address, State: "pending"}
	}
	return out
}
func encodeSMTPRecipients(values []domain.NotificationRecipient) string {
	if values == nil {
		return "[]"
	}
	data, _ := json.Marshal(values)
	return string(data)
}

// CheckpointSMTPRecipients records the recipients about to enter DATA before transmission.
func (s *Store) CheckpointSMTPRecipients(id string, values []domain.NotificationRecipient) error {
	if len(values) == 0 || len(values) > 100 {
		return fmt.Errorf("store: invalid SMTP recipient checkpoint")
	}
	res, err := s.db.Exec("UPDATE notification_deliveries SET recipients=? WHERE id=? AND kind='smtp' AND state='claimed'", encodeSMTPRecipients(values), id)
	return affected(res, err, "checkpoint SMTP recipients")
}

// FinishSMTPAttempt atomically retains recipient results and schedules or completes work.
func (s *Store) FinishSMTPAttempt(d domain.NotificationDelivery, result domain.SMTPAttemptResult, next, now time.Time) error {
	if len(result.Recipients) == 0 || len(result.Recipients) > 100 {
		return fmt.Errorf("store: invalid SMTP recipient results")
	}
	succeeded, retry := true, false
	for _, recipient := range result.Recipients {
		if recipient.State != "accepted" {
			succeeded = false
		}
		if recipient.State == "pending" || recipient.State == "temporary_failure" || recipient.State == "uncertain" {
			retry = true
		}
	}
	state := domain.NotificationDeliveryFailed
	var completed any = fmtTime(now)
	retain := false
	if succeeded {
		state = domain.NotificationDeliverySucceeded
	} else if retry && d.Attempts < 3 {
		state = domain.NotificationDeliveryPending
		completed = nil
	}
	if state == domain.NotificationDeliveryPending {
		retain = true
	}
	diagnostic := boundedDiagnostic(result.Diagnostic)
	lastError := diagnostic
	if succeeded {
		lastError = ""
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin SMTP result: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	res, err := tx.Exec("UPDATE notification_deliveries SET state=?,recipients=?,next_attempt_at=?,claimed_at=NULL,completed_at=?,last_status=?,last_error=?,diagnostic=?,authorization=CASE WHEN ? THEN authorization ELSE '' END,endpoint=CASE WHEN ? THEN endpoint ELSE '' END WHERE id=? AND kind='smtp' AND state='claimed' AND attempts=?", string(state), encodeSMTPRecipients(result.Recipients), fmtTime(next), completed, result.Status, lastError, diagnostic, retain, retain, d.ID, d.Attempts)
	if err := affected(res, err, "record SMTP result"); err != nil {
		return err
	}
	if err := pruneNotificationDeliveries(tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit SMTP result: %w", err)
	}
	return nil
}

func recoverSMTPRecipients(tx *sql.Tx) error {
	rows, err := tx.Query("SELECT id,recipients FROM notification_deliveries WHERE kind='smtp' AND state='claimed'")
	if err != nil {
		return fmt.Errorf("store: read interrupted SMTP recipients: %w", err)
	}
	type update struct{ id, raw string }
	var updates []update
	for rows.Next() {
		var id, raw string
		if err := rows.Scan(&id, &raw); err != nil {
			_ = rows.Close()
			return fmt.Errorf("store: scan interrupted SMTP recipients: %w", err)
		}
		var recipients []domain.NotificationRecipient
		if err := json.Unmarshal([]byte(raw), &recipients); err != nil {
			_ = rows.Close()
			return fmt.Errorf("store: decode interrupted SMTP recipients")
		}
		for i := range recipients {
			if recipients[i].State == "submitting" {
				recipients[i].State = "uncertain"
				recipients[i].Diagnostic = "Submission was interrupted before its confirmation was saved; retry may duplicate email."
			}
		}
		updates = append(updates, update{id, encodeSMTPRecipients(recipients)})
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	for _, u := range updates {
		if _, err := tx.Exec("UPDATE notification_deliveries SET recipients=? WHERE id=?", u.raw, u.id); err != nil {
			return fmt.Errorf("store: retain interrupted SMTP recipients: %w", err)
		}
	}
	return nil
}

func deliveryTransportKind(kind domain.NotificationChannelKind) string {
	if kind == "" {
		return string(domain.NotificationChannelWebhook)
	}
	return string(kind)
}
