package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/buildinfo"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/secretstore"
)

const notificationHistoryLimit = 1000

// ErrNotificationChannelDisabled prevents test work on disabled channels.
var ErrNotificationChannelDisabled = errors.New("store: notification channel is disabled")

// CreateNotificationChannel persists a validated write-only destination.
func (s *Store) CreateNotificationChannel(channel *domain.NotificationChannel) error {
	if channel.ID == "" {
		channel.ID = newID()
	}
	now := time.Now().UTC()
	if channel.CreatedAt.IsZero() {
		channel.CreatedAt = now
	}
	channel.UpdatedAt = now
	endpoint, err := secretstore.Protect(channel.Endpoint)
	if err != nil {
		return fmt.Errorf("store: protect notification endpoint: %w", err)
	}
	authorization, err := secretstore.Protect(channel.Authorization)
	if err != nil {
		return fmt.Errorf("store: protect notification authorization: %w", err)
	}
	_, err = s.db.Exec(`INSERT INTO notification_channels(id,name,kind,endpoint,endpoint_summary,authorization,enabled,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)`, channel.ID, channel.Name, string(channel.Kind), endpoint, channel.EndpointSummary, authorization, boolToInt(channel.Enabled), fmtTime(channel.CreatedAt), fmtTime(channel.UpdatedAt))
	if err != nil {
		return fmt.Errorf("store: create notification channel: %w", err)
	}
	channel.HasAuthorization = channel.Authorization != ""
	return nil
}

const notificationChannelSelect = `SELECT id,name,kind,endpoint,endpoint_summary,authorization,enabled,created_at,updated_at FROM notification_channels`

// GetNotificationChannel returns one channel, including protected values for daemon use.
func (s *Store) GetNotificationChannel(id string) (domain.NotificationChannel, error) {
	return scanNotificationChannel(s.db.QueryRow(notificationChannelSelect+` WHERE id=?`, id))
}

// ListNotificationChannels returns all channels without changing protected values.
func (s *Store) ListNotificationChannels() ([]domain.NotificationChannel, error) {
	rows, err := s.db.Query(notificationChannelSelect + ` ORDER BY name,id`)
	if err != nil {
		return nil, fmt.Errorf("store: list notification channels: %w", err)
	}
	defer rows.Close()
	var out []domain.NotificationChannel
	for rows.Next() {
		channel, err := scanNotificationChannel(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, channel)
	}
	return out, rows.Err()
}

func scanNotificationChannel(sc scanner) (domain.NotificationChannel, error) {
	var channel domain.NotificationChannel
	var kind, created, updated string
	var enabled int
	if err := sc.Scan(&channel.ID, &channel.Name, &kind, &channel.Endpoint, &channel.EndpointSummary, &channel.Authorization, &enabled, &created, &updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return channel, ErrNotFound
		}
		return channel, fmt.Errorf("store: scan notification channel: %w", err)
	}
	channel.Kind, channel.Enabled, channel.HasAuthorization = domain.NotificationChannelKind(kind), enabled != 0, channel.Authorization != ""
	var err error
	channel.Endpoint, err = secretstore.Unprotect(channel.Endpoint)
	if err != nil {
		return channel, fmt.Errorf("store: unprotect notification endpoint: %w", err)
	}
	channel.Authorization, err = secretstore.Unprotect(channel.Authorization)
	if err != nil {
		return channel, fmt.Errorf("store: unprotect notification authorization: %w", err)
	}
	channel.HasAuthorization = channel.Authorization != ""
	channel.CreatedAt, _ = parseTime(created)
	channel.UpdatedAt, _ = parseTime(updated)
	return channel, nil
}

// UpdateNotificationChannel replaces the mutable non-authorization fields.
func (s *Store) UpdateNotificationChannel(channel domain.NotificationChannel) error {
	endpoint, err := secretstore.Protect(channel.Endpoint)
	if err != nil {
		return fmt.Errorf("store: protect notification endpoint: %w", err)
	}
	res, err := s.db.Exec(`UPDATE notification_channels SET name=?,endpoint=?,endpoint_summary=?,enabled=?,updated_at=? WHERE id=?`, channel.Name, endpoint, channel.EndpointSummary, boolToInt(channel.Enabled), fmtTime(time.Now().UTC()), channel.ID)
	return affected(res, err, "update notification channel")
}

// SetNotificationChannelEnabled changes whether new work may be created.
func (s *Store) SetNotificationChannelEnabled(id string, enabled bool) error {
	res, err := s.db.Exec(`UPDATE notification_channels SET enabled=?,updated_at=? WHERE id=?`, boolToInt(enabled), fmtTime(time.Now().UTC()), id)
	return affected(res, err, "set notification channel enabled")
}

// RotateNotificationChannelAuthorization replaces or removes the write-only value.
func (s *Store) RotateNotificationChannelAuthorization(id, authorization string) error {
	protected, err := secretstore.Protect(authorization)
	if err != nil {
		return fmt.Errorf("store: protect notification authorization: %w", err)
	}
	res, err := s.db.Exec(`UPDATE notification_channels SET authorization=?,updated_at=? WHERE id=?`, protected, fmtTime(time.Now().UTC()), id)
	return affected(res, err, "rotate notification channel authorization")
}

// DeleteNotificationChannel erases its secret and unfinished work but preserves terminal evidence.
func (s *Store) DeleteNotificationChannel(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin delete notification channel: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.Exec(`DELETE FROM notification_deliveries WHERE channel_id=? AND state IN (?,?)`, id, string(domain.NotificationDeliveryPending), string(domain.NotificationDeliveryClaimed)); err != nil {
		return fmt.Errorf("store: delete unfinished notification deliveries: %w", err)
	}
	res, err := tx.Exec(`DELETE FROM notification_channels WHERE id=?`, id)
	if err != nil {
		return fmt.Errorf("store: delete notification channel: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: count deleted notification channels: %w", err)
	}
	if count == 0 {
		return ErrNotFound
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit delete notification channel: %w", err)
	}
	return nil
}

// ReplaceNotificationAssignments atomically replaces one scope's complete policy.
func (s *Store) ReplaceNotificationAssignments(scopeType domain.NotificationScopeType, scopeID string, assignments []domain.NotificationAssignment) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin replace notification assignments: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	column, table := "task_id", "tasks"
	if scopeType == domain.NotificationScopeGroup {
		column, table = "group_id", "groups"
	} else if scopeType != domain.NotificationScopeTask {
		return fmt.Errorf("store: invalid notification scope %q", scopeType)
	}
	var exists int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM `+table+` WHERE id=?`, scopeID).Scan(&exists); err != nil || exists == 0 {
		if err != nil {
			return fmt.Errorf("store: validate notification scope: %w", err)
		}
		return ErrNotFound
	}
	if _, err := tx.Exec(`DELETE FROM notification_assignments WHERE `+column+`=?`, scopeID); err != nil {
		return fmt.Errorf("store: clear notification assignments: %w", err)
	}
	now, seen := time.Now().UTC(), map[string]bool{}
	for i := range assignments {
		a := &assignments[i]
		if a.ChannelID == "" || (!a.OnSuccess && !a.OnFailure) || seen[a.ChannelID] {
			return fmt.Errorf("store: invalid or duplicate notification assignment")
		}
		seen[a.ChannelID] = true
		var channelExists int
		if err := tx.QueryRow(`SELECT COUNT(*) FROM notification_channels WHERE id=?`, a.ChannelID).Scan(&channelExists); err != nil || channelExists == 0 {
			if err != nil {
				return fmt.Errorf("store: validate notification channel: %w", err)
			}
			return ErrNotFound
		}
		a.ID, a.ScopeType, a.ScopeID, a.CreatedAt, a.UpdatedAt = newID(), scopeType, scopeID, now, now
		var taskID, groupID any
		if scopeType == domain.NotificationScopeTask {
			taskID = scopeID
		} else {
			groupID = scopeID
		}
		if _, err := tx.Exec(`INSERT INTO notification_assignments(id,channel_id,task_id,group_id,on_success,on_failure,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?)`, a.ID, a.ChannelID, taskID, groupID, boolToInt(a.OnSuccess), boolToInt(a.OnFailure), fmtTime(now), fmtTime(now)); err != nil {
			return fmt.Errorf("store: create notification assignment: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit notification assignments: %w", err)
	}
	return nil
}

// ListNotificationAssignments returns the policy configured directly on one scope.
func (s *Store) ListNotificationAssignments(scopeType domain.NotificationScopeType, scopeID string) ([]domain.NotificationAssignment, error) {
	column := "task_id"
	if scopeType == domain.NotificationScopeGroup {
		column = "group_id"
	} else if scopeType != domain.NotificationScopeTask {
		return nil, fmt.Errorf("store: invalid notification scope %q", scopeType)
	}
	return listNotificationAssignments(s.db, column, scopeID, scopeType)
}

type notificationQuerier interface {
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
}

func listNotificationAssignments(q notificationQuerier, column, scopeID string, scopeType domain.NotificationScopeType) ([]domain.NotificationAssignment, error) {
	rows, err := q.Query(`SELECT id,channel_id,on_success,on_failure,created_at,updated_at FROM notification_assignments WHERE `+column+`=? ORDER BY created_at,id`, scopeID)
	if err != nil {
		return nil, fmt.Errorf("store: list notification assignments: %w", err)
	}
	defer rows.Close()
	var out []domain.NotificationAssignment
	for rows.Next() {
		var a domain.NotificationAssignment
		var success, failure int
		var created, updated string
		if err := rows.Scan(&a.ID, &a.ChannelID, &success, &failure, &created, &updated); err != nil {
			return nil, fmt.Errorf("store: scan notification assignment: %w", err)
		}
		a.ScopeType, a.ScopeID, a.OnSuccess, a.OnFailure = scopeType, scopeID, success != 0, failure != 0
		a.CreatedAt, _ = parseTime(created)
		a.UpdatedAt, _ = parseTime(updated)
		out = append(out, a)
	}
	return out, rows.Err()
}

// EffectiveNotificationPolicy selects the nearest non-empty policy for a task.
func (s *Store) EffectiveNotificationPolicy(taskID string) (domain.EffectiveNotificationPolicy, error) {
	return effectiveNotificationPolicy(s.db, taskID)
}

func effectiveNotificationPolicy(q notificationQuerier, taskID string) (domain.EffectiveNotificationPolicy, error) {
	policy := domain.EffectiveNotificationPolicy{TaskID: taskID, SourceScopeType: domain.NotificationScopeNone, Assignments: []domain.NotificationAssignment{}}
	var groupID sql.NullString
	if err := q.QueryRow(`SELECT group_id FROM tasks WHERE id=?`, taskID).Scan(&groupID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return policy, ErrNotFound
		}
		return policy, fmt.Errorf("store: read task notification scope: %w", err)
	}
	assignments, err := listNotificationAssignments(q, "task_id", taskID, domain.NotificationScopeTask)
	if err != nil {
		return policy, err
	}
	if len(assignments) > 0 {
		policy.SourceScopeType, policy.SourceScopeID, policy.Assignments = domain.NotificationScopeTask, taskID, assignments
		return policy, nil
	}
	for current := groupID.String; current != ""; {
		assignments, err = listNotificationAssignments(q, "group_id", current, domain.NotificationScopeGroup)
		if err != nil {
			return policy, err
		}
		if len(assignments) > 0 {
			policy.SourceScopeType, policy.SourceScopeID, policy.Assignments = domain.NotificationScopeGroup, current, assignments
			return policy, nil
		}
		var parent sql.NullString
		if err := q.QueryRow(`SELECT parent_id FROM groups WHERE id=?`, current).Scan(&parent); err != nil {
			return policy, fmt.Errorf("store: read notification group ancestry: %w", err)
		}
		current = parent.String
	}
	return policy, nil
}

type notificationExecer interface {
	Exec(string, ...any) (sql.Result, error)
}

func (s *Store) insertNotificationDelivery(q notificationExecer, d domain.NotificationDelivery) error {
	endpoint, err := secretstore.Protect(d.Endpoint)
	if err != nil {
		return fmt.Errorf("store: protect delivery endpoint: %w", err)
	}
	authorization, err := secretstore.Protect(d.Authorization)
	if err != nil {
		return fmt.Errorf("store: protect delivery authorization: %w", err)
	}
	_, err = q.Exec(`INSERT INTO notification_deliveries(id,channel_id,channel_name,destination_summary,endpoint,authorization,event_kind,task_id,run_id,task_name,group_id,group_name,payload,state,attempts,next_attempt_at,created_at,claimed_at,completed_at,last_status,last_error) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, d.ID, nullStr(d.ChannelID), d.ChannelName, d.DestinationSummary, endpoint, authorization, string(d.EventKind), nullStr(d.TaskID), nullStr(d.RunID), d.TaskName, d.GroupID, d.GroupName, []byte(d.Payload), string(d.State), d.Attempts, fmtTime(d.NextAttemptAt), fmtTime(d.CreatedAt), fmtTimePtr(d.ClaimedAt), fmtTimePtr(d.CompletedAt), d.LastStatus, d.LastError)
	if err != nil {
		return fmt.Errorf("store: create notification delivery: %w", err)
	}
	return nil
}

// CreateTestNotificationDelivery creates durable test work without a task run.
func (s *Store) CreateTestNotificationDelivery(channelID string) (domain.NotificationDelivery, error) {
	channel, err := s.GetNotificationChannel(channelID)
	if err != nil {
		return domain.NotificationDelivery{}, err
	}
	if !channel.Enabled {
		return domain.NotificationDelivery{}, ErrNotificationChannelDisabled
	}
	now := time.Now().UTC()
	delivery := domain.NotificationDelivery{ID: newID(), ChannelID: channel.ID, ChannelName: channel.Name, DestinationSummary: channel.EndpointSummary, Endpoint: channel.Endpoint, Authorization: channel.Authorization, EventKind: domain.NotificationEventTest, State: domain.NotificationDeliveryPending, NextAttemptAt: now, CreatedAt: now}
	event := domain.WebhookEvent{Schema: "go-schedule.webhook.v1", Event: string(domain.NotificationEventTest), Delivery: domain.WebhookDelivery{ID: delivery.ID, CreatedAt: now}, Daemon: domain.WebhookDaemon{Version: buildinfo.Version}, Task: nil}
	delivery.Payload, err = json.Marshal(event)
	if err != nil {
		return delivery, fmt.Errorf("store: encode test notification: %w", err)
	}
	if err := s.insertNotificationDelivery(s.db, delivery); err != nil {
		return delivery, err
	}
	return delivery, nil
}

func (s *Store) createRunNotificationDeliveries(tx *sql.Tx, run domain.Run) error {
	if run.Outcome != domain.OutcomeSuccess && run.Outcome != domain.OutcomeFailure {
		return nil
	}
	policy, err := effectiveNotificationPolicy(tx, run.TaskID)
	if err != nil {
		return err
	}
	var taskName string
	var groupID sql.NullString
	if err := tx.QueryRow(`SELECT name,group_id FROM tasks WHERE id=?`, run.TaskID).Scan(&taskName, &groupID); err != nil {
		return fmt.Errorf("store: snapshot notification task: %w", err)
	}
	groupName := ""
	if groupID.String != "" {
		if err := tx.QueryRow(`SELECT name FROM groups WHERE id=?`, groupID.String).Scan(&groupName); err != nil {
			return fmt.Errorf("store: snapshot notification group: %w", err)
		}
	}
	now := time.Now().UTC()
	for _, assignment := range policy.Assignments {
		if (run.Outcome == domain.OutcomeSuccess && !assignment.OnSuccess) || (run.Outcome == domain.OutcomeFailure && !assignment.OnFailure) {
			continue
		}
		channel, err := scanNotificationChannel(tx.QueryRow(notificationChannelSelect+` WHERE id=? AND enabled=1`, assignment.ChannelID))
		if errors.Is(err, ErrNotFound) {
			continue
		}
		if err != nil {
			return err
		}
		d := domain.NotificationDelivery{ID: newID(), ChannelID: channel.ID, ChannelName: channel.Name, DestinationSummary: channel.EndpointSummary, Endpoint: channel.Endpoint, Authorization: channel.Authorization, EventKind: domain.NotificationEventRunCompleted, TaskID: run.TaskID, RunID: run.ID, TaskName: taskName, GroupID: groupID.String, GroupName: groupName, State: domain.NotificationDeliveryPending, NextAttemptAt: now, CreatedAt: now}
		duration := int64(0)
		if run.StartedAt != nil && run.EndedAt != nil && !run.EndedAt.Before(*run.StartedAt) {
			duration = run.EndedAt.Sub(*run.StartedAt).Milliseconds()
		}
		event := domain.WebhookEvent{Schema: "go-schedule.webhook.v1", Event: string(domain.NotificationEventRunCompleted), Delivery: domain.WebhookDelivery{ID: d.ID, CreatedAt: now}, Daemon: domain.WebhookDaemon{Version: buildinfo.Version}, Task: &domain.WebhookTask{ID: run.TaskID, Name: taskName, GroupID: groupID.String, GroupName: groupName}, Run: &domain.WebhookRun{ID: run.ID, Outcome: run.Outcome, Trigger: run.Trigger, ScheduledFor: run.ScheduledFor, StartedAt: run.StartedAt, EndedAt: run.EndedAt, DurationMS: duration, ExitCode: run.ExitCode}}
		d.Payload, err = json.Marshal(event)
		if err != nil {
			return fmt.Errorf("store: encode run notification: %w", err)
		}
		if err := s.insertNotificationDelivery(tx, d); err != nil {
			return err
		}
	}
	return nil
}

// RecoverNotificationDeliveries makes interrupted claims replayable without resetting attempts.
func (s *Store) RecoverNotificationDeliveries(now time.Time) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("store: begin recover notification deliveries: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	res, err := tx.Exec(`UPDATE notification_deliveries SET state=CASE WHEN attempts>=3 THEN ? ELSE ? END,next_attempt_at=?,claimed_at=NULL,completed_at=CASE WHEN attempts>=3 THEN ? ELSE NULL END,endpoint=CASE WHEN attempts>=3 THEN '' ELSE endpoint END,authorization=CASE WHEN attempts>=3 THEN '' ELSE authorization END,last_error=CASE WHEN attempts>=3 THEN 'retry budget exhausted during recovery' ELSE 'recovered after daemon restart' END WHERE state=?`, string(domain.NotificationDeliveryFailed), string(domain.NotificationDeliveryPending), fmtTime(now), fmtTime(now), string(domain.NotificationDeliveryClaimed))
	if err != nil {
		return 0, fmt.Errorf("store: recover notification deliveries: %w", err)
	}
	if err := pruneNotificationDeliveries(tx); err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("store: count recovered notification deliveries: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("store: commit recovered notification deliveries: %w", err)
	}
	return count, nil
}

// ClaimNotificationDeliveries claims eligible work and increments attempt counts.
func (s *Store) ClaimNotificationDeliveries(limit int, now time.Time) ([]domain.NotificationDelivery, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("store: begin claim notification deliveries: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	rows, err := tx.Query(notificationDeliverySelect+` WHERE state=? AND attempts<3 AND next_attempt_at<=? ORDER BY next_attempt_at,created_at,id LIMIT ?`, string(domain.NotificationDeliveryPending), fmtTime(now), limit)
	if err != nil {
		return nil, fmt.Errorf("store: list pending notification deliveries: %w", err)
	}
	var out []domain.NotificationDelivery
	for rows.Next() {
		d, err := scanNotificationDelivery(rows)
		if err != nil {
			_ = rows.Close()
			return nil, err
		}
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("store: iterate pending notification deliveries: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("store: close pending notification deliveries: %w", err)
	}
	for i := range out {
		res, err := tx.Exec(`UPDATE notification_deliveries SET state=?,attempts=attempts+1,claimed_at=? WHERE id=? AND state=?`, string(domain.NotificationDeliveryClaimed), fmtTime(now), out[i].ID, string(domain.NotificationDeliveryPending))
		if err != nil {
			return nil, fmt.Errorf("store: claim notification delivery: %w", err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("store: count claimed notification delivery: %w", err)
		}
		if count != 1 {
			return nil, fmt.Errorf("store: claim notification delivery %s: %w", out[i].ID, ErrNotFound)
		}
		out[i].State, out[i].Attempts, out[i].ClaimedAt = domain.NotificationDeliveryClaimed, out[i].Attempts+1, &now
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("store: commit notification delivery claims: %w", err)
	}
	return out, nil
}

// RetryNotificationDelivery returns a failed attempt to pending within its budget.
func (s *Store) RetryNotificationDelivery(id string, next time.Time, status int, diagnostic string) error {
	res, err := s.db.Exec(`UPDATE notification_deliveries SET state=?,next_attempt_at=?,claimed_at=NULL,last_status=?,last_error=? WHERE id=? AND state=? AND attempts<3`, string(domain.NotificationDeliveryPending), fmtTime(next), status, boundedDiagnostic(diagnostic), id, string(domain.NotificationDeliveryClaimed))
	return affected(res, err, "retry notification delivery")
}

// CompleteNotificationDelivery records terminal evidence and erases protected snapshots.
func (s *Store) CompleteNotificationDelivery(id string, succeeded bool, status int, diagnostic string, now time.Time) error {
	state := domain.NotificationDeliveryFailed
	if succeeded {
		state = domain.NotificationDeliverySucceeded
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin complete notification delivery: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	res, err := tx.Exec(`UPDATE notification_deliveries SET state=?,completed_at=?,claimed_at=NULL,last_status=?,last_error=?,endpoint='',authorization='' WHERE id=? AND state=?`, string(state), fmtTime(now), status, boundedDiagnostic(diagnostic), id, string(domain.NotificationDeliveryClaimed))
	if err != nil {
		return fmt.Errorf("store: complete notification delivery: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: count completed notification delivery: %w", err)
	}
	if count != 1 {
		return ErrNotFound
	}
	if err := pruneNotificationDeliveries(tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit notification delivery: %w", err)
	}
	return nil
}

func pruneNotificationDeliveries(tx *sql.Tx) error {
	if _, err := tx.Exec(`DELETE FROM notification_deliveries WHERE id IN (SELECT id FROM notification_deliveries WHERE state IN (?,?) ORDER BY completed_at DESC,id DESC LIMIT -1 OFFSET ?)`, string(domain.NotificationDeliverySucceeded), string(domain.NotificationDeliveryFailed), notificationHistoryLimit); err != nil {
		return fmt.Errorf("store: prune notification deliveries: %w", err)
	}
	return nil
}

const notificationDeliverySelect = `SELECT id,channel_id,channel_name,destination_summary,endpoint,authorization,event_kind,task_id,run_id,task_name,group_id,group_name,payload,state,attempts,next_attempt_at,created_at,claimed_at,completed_at,last_status,last_error FROM notification_deliveries`

func scanNotificationDelivery(sc scanner) (domain.NotificationDelivery, error) {
	var d domain.NotificationDelivery
	var channelID, taskID, runID sql.NullString
	var eventKind, state, next, created string
	var claimed, completed sql.NullString
	var payload []byte
	if err := sc.Scan(&d.ID, &channelID, &d.ChannelName, &d.DestinationSummary, &d.Endpoint, &d.Authorization, &eventKind, &taskID, &runID, &d.TaskName, &d.GroupID, &d.GroupName, &payload, &state, &d.Attempts, &next, &created, &claimed, &completed, &d.LastStatus, &d.LastError); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return d, ErrNotFound
		}
		return d, fmt.Errorf("store: scan notification delivery: %w", err)
	}
	d.ChannelID, d.TaskID, d.RunID, d.EventKind, d.State, d.Payload = channelID.String, taskID.String, runID.String, domain.NotificationEventKind(eventKind), domain.NotificationDeliveryState(state), payload
	var err error
	d.Endpoint, err = secretstore.Unprotect(d.Endpoint)
	if err != nil {
		return d, fmt.Errorf("store: unprotect delivery endpoint: %w", err)
	}
	d.Authorization, err = secretstore.Unprotect(d.Authorization)
	if err != nil {
		return d, fmt.Errorf("store: unprotect delivery authorization: %w", err)
	}
	d.NextAttemptAt, _ = parseTime(next)
	d.CreatedAt, _ = parseTime(created)
	d.ClaimedAt, _ = parseTimePtr(claimed)
	d.CompletedAt, _ = parseTimePtr(completed)
	return d, nil
}

// ListNotificationDeliveries returns bounded redacted history and pending work.
func (s *Store) ListNotificationDeliveries(filter domain.NotificationDeliveryFilter) ([]domain.NotificationDelivery, error) {
	if filter.Limit <= 0 {
		filter.Limit = 100
	}
	if filter.Limit > 1000 {
		filter.Limit = 1000
	}
	conditions, args := []string{"1=1"}, []any{}
	values := []struct{ value, column string }{{filter.ChannelID, "channel_id"}, {filter.TaskID, "task_id"}, {filter.RunID, "run_id"}}
	for _, item := range values {
		if item.value != "" {
			conditions, args = append(conditions, item.column+`=?`), append(args, item.value)
		}
	}
	if filter.State != "" {
		conditions, args = append(conditions, `state=?`), append(args, string(filter.State))
	}
	args = append(args, filter.Limit)
	rows, err := s.db.Query(notificationDeliverySelect+` WHERE `+strings.Join(conditions, ` AND `)+` ORDER BY created_at DESC,id DESC LIMIT ?`, args...)
	if err != nil {
		return nil, fmt.Errorf("store: list notification deliveries: %w", err)
	}
	defer rows.Close()
	var out []domain.NotificationDelivery
	for rows.Next() {
		d, err := scanNotificationDelivery(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func boundedDiagnostic(value string) string {
	value = strings.TrimSpace(value)
	if len(value) > 512 {
		return value[:512]
	}
	return value
}
