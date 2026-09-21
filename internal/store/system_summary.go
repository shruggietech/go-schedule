package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// SystemSummaryFacts returns bounded counts and representative records for an observation window.
func (s *Store) SystemSummaryFacts(observedAt time.Time) (domain.SystemSummary, error) {
	observedAt = observedAt.UTC()
	since := observedAt.Add(-24 * time.Hour)
	result := domain.SystemSummary{Schema: domain.SystemSummarySchema, ObservedAt: observedAt}

	if err := s.db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE enabled=1 AND state=?`, string(domain.TaskActive)).Scan(&result.ActiveTaskCount); err != nil {
		return domain.SystemSummary{}, fmt.Errorf("store: count active tasks: %w", err)
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM runs WHERE outcome=? AND COALESCE(ended_at,scheduled_for)>=? AND COALESCE(ended_at,scheduled_for)<=?`, string(domain.OutcomeFailure), fmtTime(since), fmtTime(observedAt)).Scan(&result.RecentFailureCount); err != nil {
		return domain.SystemSummary{}, fmt.Errorf("store: count recent failures: %w", err)
	}
	if result.RecentFailureCount > 0 {
		var ref domain.FailureSummary
		var ended string
		err := s.db.QueryRow(`SELECT r.id,r.task_id,t.name,COALESCE(r.ended_at,r.scheduled_for) FROM runs r JOIN tasks t ON t.id=r.task_id WHERE r.outcome=? AND COALESCE(r.ended_at,r.scheduled_for)>=? AND COALESCE(r.ended_at,r.scheduled_for)<=? ORDER BY COALESCE(r.ended_at,r.scheduled_for) DESC,r.id DESC LIMIT 1`, string(domain.OutcomeFailure), fmtTime(since), fmtTime(observedAt)).Scan(&ref.RunID, &ref.TaskID, &ref.TaskName, &ended)
		if err != nil {
			return domain.SystemSummary{}, fmt.Errorf("store: read recent failure: %w", err)
		}
		ref.EndedAt, err = parseTime(ended)
		if err != nil {
			return domain.SystemSummary{}, fmt.Errorf("store: parse recent failure timestamp: %w", err)
		}
		result.RecentFailure = &ref
	}

	if err := s.db.QueryRow(`SELECT COUNT(*) FROM alerts WHERE acknowledged=0`).Scan(&result.UnacknowledgedAlertCount); err != nil {
		return domain.SystemSummary{}, fmt.Errorf("store: count unacknowledged alerts: %w", err)
	}
	if result.UnacknowledgedAlertCount > 0 {
		var ref domain.AlertSummary
		var taskID, runID sql.NullString
		var severity, kind, created string
		err := s.db.QueryRow(`SELECT id,task_id,run_id,severity,kind,created_at FROM alerts WHERE acknowledged=0 ORDER BY created_at DESC,id DESC LIMIT 1`).Scan(&ref.AlertID, &taskID, &runID, &severity, &kind, &created)
		if err != nil {
			return domain.SystemSummary{}, fmt.Errorf("store: read unacknowledged alert: %w", err)
		}
		ref.TaskID, ref.RunID = taskID.String, runID.String
		ref.Severity, ref.Kind = domain.AlertSeverity(severity), domain.AlertKind(kind)
		ref.CreatedAt, err = parseTime(created)
		if err != nil {
			return domain.SystemSummary{}, fmt.Errorf("store: parse unacknowledged alert timestamp: %w", err)
		}
		result.UnacknowledgedAlert = &ref
	}

	states := []string{string(domain.NotificationDeliveryPending), string(domain.NotificationDeliveryClaimed), string(domain.NotificationDeliveryFailed)}
	problemPredicate := `(state=? OR (state IN (?,?) AND attempts>0)) AND created_at>=? AND created_at<=?`
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM notification_deliveries WHERE `+problemPredicate, states[2], states[0], states[1], fmtTime(since), fmtTime(observedAt)).Scan(&result.NotificationProblemCount); err != nil {
		return domain.SystemSummary{}, fmt.Errorf("store: count notification problems: %w", err)
	}
	if result.NotificationProblemCount > 0 {
		var ref domain.NotificationProblemSummary
		var taskID, runID sql.NullString
		var state, created string
		err := s.db.QueryRow(`SELECT id,task_id,run_id,channel_name,state,created_at FROM notification_deliveries WHERE `+problemPredicate+` ORDER BY created_at DESC,id DESC LIMIT 1`, states[2], states[0], states[1], fmtTime(since), fmtTime(observedAt)).Scan(&ref.DeliveryID, &taskID, &runID, &ref.ChannelName, &state, &created)
		if err != nil {
			return domain.SystemSummary{}, fmt.Errorf("store: read notification problem: %w", err)
		}
		ref.TaskID, ref.RunID = taskID.String, runID.String
		ref.State = domain.NotificationDeliveryState(state)
		ref.CreatedAt, err = parseTime(created)
		if err != nil {
			return domain.SystemSummary{}, fmt.Errorf("store: parse notification problem timestamp: %w", err)
		}
		result.NotificationProblem = &ref
	}
	return result, nil
}
