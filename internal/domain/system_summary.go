package domain

import "time"

// SystemSummarySchema identifies the additive operational summary contract.
const SystemSummarySchema = "go-schedule.system-summary.v1"

// UpcomingSummary identifies the nearest scheduled task occurrence without exposing execution details.
type UpcomingSummary struct {
	TaskID       string    `json:"task_id"`
	TaskName     string    `json:"task_name"`
	ScheduledFor time.Time `json:"scheduled_for"`
}

// FailureSummary identifies one recent failed run without exposing captured output.
type FailureSummary struct {
	RunID    string    `json:"run_id"`
	TaskID   string    `json:"task_id"`
	TaskName string    `json:"task_name"`
	EndedAt  time.Time `json:"ended_at"`
}

// AlertSummary identifies one unacknowledged alert without exposing its message.
type AlertSummary struct {
	AlertID   string        `json:"alert_id"`
	TaskID    string        `json:"task_id,omitempty"`
	RunID     string        `json:"run_id,omitempty"`
	Severity  AlertSeverity `json:"severity"`
	Kind      AlertKind     `json:"kind"`
	CreatedAt time.Time     `json:"created_at"`
}

// NotificationProblemSummary identifies one failed or retrying delivery without exposing its destination or payload.
type NotificationProblemSummary struct {
	DeliveryID  string                    `json:"delivery_id"`
	TaskID      string                    `json:"task_id,omitempty"`
	RunID       string                    `json:"run_id,omitempty"`
	ChannelName string                    `json:"channel_name"`
	State       NotificationDeliveryState `json:"state"`
	CreatedAt   time.Time                 `json:"created_at"`
}

// SystemSummary is one bounded, read-only operational observation from a daemon.
type SystemSummary struct {
	Schema                   string                      `json:"schema"`
	ObservedAt               time.Time                   `json:"observed_at"`
	ActiveTaskCount          int                         `json:"active_task_count"`
	NextOccurrence           *UpcomingSummary            `json:"next_occurrence,omitempty"`
	RecentFailureCount       int                         `json:"recent_failure_count"`
	RecentFailure            *FailureSummary             `json:"recent_failure,omitempty"`
	UnacknowledgedAlertCount int                         `json:"unacknowledged_alert_count"`
	UnacknowledgedAlert      *AlertSummary               `json:"unacknowledged_alert,omitempty"`
	NotificationProblemCount int                         `json:"notification_problem_count"`
	NotificationProblem      *NotificationProblemSummary `json:"notification_problem,omitempty"`
}
