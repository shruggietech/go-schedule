package domain

import (
	"encoding/json"
	"time"
)

// NotificationChannelKind identifies an outbound notification transport.
type NotificationChannelKind string

const (
	// NotificationChannelWebhook is the versioned JSON-over-HTTP transport.
	NotificationChannelWebhook NotificationChannelKind = "webhook"
)

// NotificationScopeType identifies where a channel assignment is configured.
type NotificationScopeType string

const (
	NotificationScopeTask  NotificationScopeType = "task"
	NotificationScopeGroup NotificationScopeType = "group"
	NotificationScopeNone  NotificationScopeType = "none"
)

// NotificationDeliveryState is the durable lifecycle of outbound work.
type NotificationDeliveryState string

const (
	NotificationDeliveryPending   NotificationDeliveryState = "pending"
	NotificationDeliveryClaimed   NotificationDeliveryState = "claimed"
	NotificationDeliverySucceeded NotificationDeliveryState = "succeeded"
	NotificationDeliveryFailed    NotificationDeliveryState = "failed"
)

// NotificationEventKind distinguishes real run notifications from channel tests.
type NotificationEventKind string

const (
	NotificationEventRunCompleted NotificationEventKind = "run.completed"
	NotificationEventTest         NotificationEventKind = "test"
	NotificationEventDaemonHealth NotificationEventKind = "daemon.health"
)

// NotificationConditionKind explains why a delivery was created.
type NotificationConditionKind string

const (
	NotificationConditionSuccess            NotificationConditionKind = "success"
	NotificationConditionFailure            NotificationConditionKind = "failure"
	NotificationConditionFailureToStart     NotificationConditionKind = "failure_to_start"
	NotificationConditionDurationExceeded   NotificationConditionKind = "duration_exceeded"
	NotificationConditionRecovery           NotificationConditionKind = "recovery"
	NotificationConditionDaemonHealth       NotificationConditionKind = "daemon_health"
	NotificationConditionConsecutiveFailure NotificationConditionKind = "consecutive_failure"
)

// NotificationChannel is a reusable destination. Endpoint and Authorization
// are write-only protected values and are intentionally excluded from JSON.
type NotificationChannel struct {
	ID                    string                  `json:"id"`
	Name                  string                  `json:"name"`
	Kind                  NotificationChannelKind `json:"kind"`
	Endpoint              string                  `json:"-"`
	EndpointSummary       string                  `json:"endpoint_summary"`
	Authorization         string                  `json:"-"`
	HasAuthorization      bool                    `json:"has_authorization"`
	Enabled               bool                    `json:"enabled"`
	HealthIntervalSeconds int64                   `json:"health_interval_seconds"`
	CreatedAt             time.Time               `json:"created_at"`
	UpdatedAt             time.Time               `json:"updated_at"`
}

// NotificationAssignment binds a channel and outcome conditions to one scope.
type NotificationAssignment struct {
	ID                       string                `json:"id"`
	ChannelID                string                `json:"channel_id"`
	ScopeType                NotificationScopeType `json:"scope_type"`
	ScopeID                  string                `json:"scope_id"`
	OnSuccess                bool                  `json:"on_success"`
	OnFailure                bool                  `json:"on_failure"`
	FailureThreshold         int                   `json:"failure_threshold"`
	OnFailureToStart         bool                  `json:"on_failure_to_start"`
	DurationThresholdSeconds int64                 `json:"duration_threshold_seconds"`
	OnRecovery               bool                  `json:"on_recovery"`
	ReminderIntervalSeconds  int64                 `json:"reminder_interval_seconds"`
	QuietPeriodSeconds       int64                 `json:"quiet_period_seconds"`
	CreatedAt                time.Time             `json:"created_at"`
	UpdatedAt                time.Time             `json:"updated_at"`
}

// EffectiveNotificationPolicy explains the one scope selected for a task.
type EffectiveNotificationPolicy struct {
	TaskID          string                   `json:"task_id"`
	SourceScopeType NotificationScopeType    `json:"source_scope_type"`
	SourceScopeID   string                   `json:"source_scope_id,omitempty"`
	Assignments     []NotificationAssignment `json:"assignments"`
}

// NotificationDelivery is durable outbound work and its redacted evidence.
// Endpoint and Authorization exist only while work can still be attempted.
type NotificationDelivery struct {
	ID                 string                    `json:"id"`
	ChannelID          string                    `json:"channel_id,omitempty"`
	ChannelName        string                    `json:"channel_name"`
	DestinationSummary string                    `json:"destination_summary"`
	Endpoint           string                    `json:"-"`
	Authorization      string                    `json:"-"`
	EventKind          NotificationEventKind     `json:"event_kind"`
	TaskID             string                    `json:"task_id,omitempty"`
	RunID              string                    `json:"run_id,omitempty"`
	TaskName           string                    `json:"task_name,omitempty"`
	GroupID            string                    `json:"group_id,omitempty"`
	GroupName          string                    `json:"group_name,omitempty"`
	Payload            json.RawMessage           `json:"event"`
	State              NotificationDeliveryState `json:"state"`
	Attempts           int                       `json:"attempts"`
	NextAttemptAt      time.Time                 `json:"next_attempt_at"`
	CreatedAt          time.Time                 `json:"created_at"`
	ClaimedAt          *time.Time                `json:"claimed_at,omitempty"`
	CompletedAt        *time.Time                `json:"completed_at,omitempty"`
	LastStatus         int                       `json:"last_status,omitempty"`
	LastError          string                    `json:"last_error,omitempty"`
	ConditionKind      NotificationConditionKind `json:"condition_kind,omitempty"`
	ConditionSummary   string                    `json:"condition_summary,omitempty"`
}

// NotificationDeliveryFilter bounds and narrows delivery history queries.
type NotificationDeliveryFilter struct {
	ChannelID string
	TaskID    string
	RunID     string
	State     NotificationDeliveryState
	Limit     int
}

// WebhookEvent is the stable receiver-facing payload.
type WebhookEvent struct {
	Schema    string            `json:"schema"`
	Event     string            `json:"event"`
	Delivery  WebhookDelivery   `json:"delivery"`
	Daemon    WebhookDaemon     `json:"daemon"`
	Task      *WebhookTask      `json:"task"`
	Run       *WebhookRun       `json:"run,omitempty"`
	Condition *WebhookCondition `json:"condition,omitempty"`
}

type WebhookDelivery struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type WebhookDaemon struct {
	Version        string     `json:"version"`
	Status         string     `json:"status,omitempty"`
	NextExpectedAt *time.Time `json:"next_expected_at,omitempty"`
}

// WebhookCondition is safe, bounded evidence explaining delivery creation.
type WebhookCondition struct {
	Kind      NotificationConditionKind `json:"kind"`
	Summary   string                    `json:"summary"`
	Streak    int                       `json:"streak,omitempty"`
	Threshold int                       `json:"threshold,omitempty"`
	Reminder  bool                      `json:"reminder,omitempty"`
}

type WebhookTask struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	GroupID   string `json:"group_id,omitempty"`
	GroupName string `json:"group_name,omitempty"`
}

type WebhookRun struct {
	ID           string     `json:"id"`
	Outcome      RunOutcome `json:"outcome"`
	Trigger      RunTrigger `json:"trigger"`
	ScheduledFor time.Time  `json:"scheduled_for"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	EndedAt      *time.Time `json:"ended_at,omitempty"`
	DurationMS   int64      `json:"duration_ms,omitempty"`
	ExitCode     *int       `json:"exit_code,omitempty"`
}
