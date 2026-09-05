// Package domain defines the core entities of the scheduler (Group, Task,
// Schedule, Run, Alert) and their enumerations. Entities are plain data with
// UTC timestamps; persistence lives in internal/store and behavior in the
// engine. Keeping the types in one low-level package avoids import cycles
// between the task, schedule, store, and engine packages.
package domain

import "time"

// ---- Enumerations -------------------------------------------------------

// TaskState is the lifecycle state of a Task.
type TaskState string

const (
	TaskActive    TaskState = "active"
	TaskCompleted TaskState = "completed" // one-off ran, or recurrence exhausted
	TaskDisabled  TaskState = "disabled"
)

// OverlapPolicy controls what happens when a task is still running at its next
// trigger time.
type OverlapPolicy string

const (
	OverlapQueueOne        OverlapPolicy = "queue_one" // default: queue exactly one pending run
	OverlapSkip            OverlapPolicy = "skip"      // skip the new trigger
	OverlapAllowConcurrent OverlapPolicy = "allow_concurrent"
)

// CatchupPolicy controls catch-up behavior after downtime.
type CatchupPolicy string

const (
	CatchupOne  CatchupPolicy = "one"  // default: one catch-up run if any were missed
	CatchupNone CatchupPolicy = "none" // never catch up
)

// MissingDatePolicy controls what a schedule does in a period that has no
// matching date: February in a non-leap year for a rule on the 29th, a 30-day
// month for a rule on the 31st, a month with only four Fridays for a rule on the
// fifth. It is consulted only by rules that address a date or an ordinal
// weekday; for interval and plain weekday rules it is inert.
type MissingDatePolicy string

const (
	// MissingDateSkip is the default and the behavior of every schedule created
	// before the policy existed: the period simply produces no run.
	MissingDateSkip MissingDatePolicy = "skip"
	// MissingDateLastValid falls back to the last date in the period that does
	// exist, Feb 29 → Feb 28, the 31st → the 30th, the 5th Friday → the last
	// Friday.
	MissingDateLastValid MissingDatePolicy = "last_valid"
	// MissingDateNextValid rolls forward into the following period, Feb 29 →
	// Mar 1, without displacing that period's own occurrence.
	MissingDateNextValid MissingDatePolicy = "next_valid"
)

// TimeBasis controls which clock gives recurrence fields their meaning.
type TimeBasis string

const (
	TimeBasisWallClock TimeBasis = "wall_clock"
	TimeBasisElapsed   TimeBasis = "elapsed"
	TimeBasisUTC       TimeBasis = "utc"
)

// DSTGapPolicy controls a wall-clock occurrence whose local reading does not
// exist because the UTC offset moves forward.
type DSTGapPolicy string

const (
	DSTGapNextValid DSTGapPolicy = "next_valid"
	DSTGapSkip      DSTGapPolicy = "skip"
)

// DSTOverlapPolicy controls a wall-clock occurrence whose local reading maps
// to two concrete instants because the UTC offset moves backward.
type DSTOverlapPolicy string

const (
	DSTOverlapFirst DSTOverlapPolicy = "first"
	DSTOverlapBoth  DSTOverlapPolicy = "both"
	DSTOverlapLast  DSTOverlapPolicy = "last"
)

// SchedulePolicy is the complete task-level calendar anomaly contract passed
// to every occurrence-producing path.
type SchedulePolicy struct {
	MissingDate MissingDatePolicy
	TimeBasis   TimeBasis
	DSTGap      DSTGapPolicy
	DSTOverlap  DSTOverlapPolicy
}

// Effective returns p with compatibility defaults substituted for zero values.
func (p SchedulePolicy) Effective() SchedulePolicy {
	if p.MissingDate == "" {
		p.MissingDate = MissingDateSkip
	}
	if p.TimeBasis == "" {
		p.TimeBasis = TimeBasisWallClock
	}
	if p.DSTGap == "" {
		p.DSTGap = DSTGapNextValid
	}
	if p.DSTOverlap == "" {
		p.DSTOverlap = DSTOverlapFirst
	}
	return p
}

// ScheduleKind distinguishes the timing model of a Schedule.
type ScheduleKind string

const (
	ScheduleOneOff    ScheduleKind = "one_off"
	ScheduleRecurring ScheduleKind = "recurring"
	ScheduleEvent     ScheduleKind = "event"
)

// StartupEventID is the stable identity of the daemon-start schedule event.
// Other event sources remain outside the current scheduling contract.
const StartupEventID = "scheduler_startup"

// CalendarAdjustment identifies a bounded calendar operation applied after an
// RRULE selects its intended date. The empty value means the RRULE is complete.
type CalendarAdjustment string

const (
	// CalendarAdjustmentNearestWeekday moves a monthly numbered date to the
	// nearest Monday-through-Friday date without crossing its resolved month.
	CalendarAdjustmentNearestWeekday CalendarAdjustment = "nearest_weekday"
)

// RunOutcome is the result of a single execution.
type RunOutcome string

const (
	OutcomeSuccess  RunOutcome = "success"
	OutcomeFailure  RunOutcome = "failure"
	OutcomeSkipped  RunOutcome = "skipped"
	OutcomeCaughtUp RunOutcome = "caught_up"
	OutcomeQueued   RunOutcome = "queued"
)

// RunTrigger records what caused a run.
type RunTrigger string

const (
	TriggerSchedule   RunTrigger = "schedule"
	TriggerEvent      RunTrigger = "event"
	TriggerStartup    RunTrigger = "startup"
	TriggerCatchup    RunTrigger = "catchup"
	TriggerManual     RunTrigger = "manual"
	TriggerCompletion RunTrigger = "completion"
	TriggerExternal   RunTrigger = "external_trigger"
)

// CompletionOutcome selects which terminal source outcomes activate a chain.
type CompletionOutcome string

const (
	CompletionOnSuccess CompletionOutcome = "success"
	CompletionOnFailure CompletionOutcome = "failure"
	CompletionOnAny     CompletionOutcome = "any"
)

// DeliveryState is the durable lifecycle of one completion delivery.
type DeliveryState string

const (
	DeliveryPending   DeliveryState = "pending"
	DeliveryClaimed   DeliveryState = "claimed"
	DeliveryCompleted DeliveryState = "completed"
	DeliveryResolved  DeliveryState = "resolved"
)

// AlertSeverity and AlertKind classify surfaced conditions.
type AlertSeverity string

const (
	SeverityInfo    AlertSeverity = "info"
	SeverityWarning AlertSeverity = "warning"
	SeverityError   AlertSeverity = "error"
)

type AlertKind string

const (
	AlertOverlapQueued AlertKind = "overlap_queued"
	AlertRunFailed     AlertKind = "run_failed"
	AlertMissedRun     AlertKind = "missed_run"
	AlertService       AlertKind = "service"
)

// ---- Entities -----------------------------------------------------------

// Group is a named container forming a nested hierarchy. ParentID is empty for
// top-level groups. Disabling cascades to descendants and their tasks.
type Group struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ParentID  string    `json:"parent_id,omitempty"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Task is a unit of work with a schedule, timezone, and execution policies.
type Task struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	GroupID       string            `json:"group_id,omitempty"`
	Command       string            `json:"command"`
	Args          []string          `json:"args,omitempty"`
	WorkingDir    string            `json:"working_dir,omitempty"`
	Env           map[string]string `json:"env,omitempty"`
	Stdin         string            `json:"stdin,omitempty"`
	RunAs         string            `json:"run_as,omitempty"`
	Enabled       bool              `json:"enabled"`
	Timezone      string            `json:"timezone"`
	ScheduleID    string            `json:"schedule_id"`
	OverlapPolicy OverlapPolicy     `json:"overlap_policy"`
	CatchupPolicy CatchupPolicy     `json:"catchup_policy"`
	// MissingDatePolicy lives on the task rather than the schedule because
	// replacing a task's schedule phrase creates a new schedule row: a
	// schedule-borne policy would silently reset to the default on an unrelated
	// edit, changing run times without the operator asking.
	MissingDatePolicy MissingDatePolicy `json:"missing_date_policy"`
	TimeBasis         TimeBasis         `json:"time_basis"`
	DSTGapPolicy      DSTGapPolicy      `json:"dst_gap_policy"`
	DSTOverlapPolicy  DSTOverlapPolicy  `json:"dst_overlap_policy"`
	State             TaskState         `json:"state"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// SchedulePolicy returns the task's effective recurrence policy set.
func (t Task) SchedulePolicy() SchedulePolicy {
	return (SchedulePolicy{
		MissingDate: t.MissingDatePolicy,
		TimeBasis:   t.TimeBasis,
		DSTGap:      t.DSTGapPolicy,
		DSTOverlap:  t.DSTOverlapPolicy,
	}).Effective()
}

// Schedule is the timing definition for a task. Exactly one of (RRULE+Anchor),
// RunAt, or TriggerID is populated, matching Kind. All times are UTC.
//
// Three fields carry timing text and are easy to confuse:
//   - RRULE plus an optional CalendarAdjustment are authoritative timing input.
//   - HumanSummary is what the system says back to the user ("Every weekday at
//     09:00"). Display only, and deliberately not re-parseable.
//   - Expression is the accepted source ("weekdays at 09:00" or
//     "0 9 * * 1-5"). It exists so a client can put the source back into the
//     field it came from, and is inert with respect to execution.
type Schedule struct {
	ID    string       `json:"id"`
	Kind  ScheduleKind `json:"kind"`
	RRULE string       `json:"rrule,omitempty"`
	// CalendarAdjustment carries the one supported calendar operation that a
	// single RFC 5545 RRULE cannot represent. It is persisted and authoritative;
	// Expression remains inert.
	CalendarAdjustment CalendarAdjustment `json:"calendar_adjustment,omitempty"`
	Anchor             *time.Time         `json:"anchor,omitempty"`
	// ElapsedEpoch is the persisted absolute phase for fixed-duration schedules.
	// It is independent of Timezone, which is presentation-only in elapsed mode.
	ElapsedEpoch *time.Time `json:"elapsed_epoch,omitempty"`
	RunAt        *time.Time `json:"run_at,omitempty"`
	TriggerID    string     `json:"trigger_id,omitempty"`
	HumanSummary string     `json:"human_summary"`
	// Expression is the human or cron source this schedule was parsed from,
	// suitable for re-submission. Empty for one-offs and legacy rows. It is never
	// an execution input: nothing on the scheduling path may read it.
	Expression string `json:"expression,omitempty"`
	// SourceSyntax identifies Expression for API clients. It is derived at the
	// response boundary and is deliberately not persisted.
	SourceSyntax string `json:"source_syntax,omitempty"`
}

// Run is a single execution record (append-only history).
type Run struct {
	ID           string     `json:"id"`
	TaskID       string     `json:"task_id"`
	ScheduledFor time.Time  `json:"scheduled_for"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	EndedAt      *time.Time `json:"ended_at,omitempty"`
	Outcome      RunOutcome `json:"outcome"`
	ExitCode     *int       `json:"exit_code,omitempty"`
	Output       string     `json:"output,omitempty"`
	// OutputTruncated reports that the configured capture cap discarded one or
	// more output bytes. It is metadata so the retained output stays within cap.
	OutputTruncated bool       `json:"output_truncated,omitempty"`
	Trigger         RunTrigger `json:"trigger"`
	SourceTaskID    string     `json:"source_task_id,omitempty"`
	SourceRunID     string     `json:"source_run_id,omitempty"`
	SourceTriggerID string     `json:"source_trigger_id,omitempty"`
}

// ExternalTrigger maps one opaque local key to one task invocation.
// Key is deliberately excluded from ordinary JSON serialization.
type ExternalTrigger struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Key            string    `json:"-"`
	TargetTaskID   string    `json:"target_task_id"`
	TargetTaskName string    `json:"target_task_name,omitempty"`
	Enabled        bool      `json:"enabled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CompletionChain connects a source task's terminal outcome to a target task.
// It supplements both tasks' normal schedules.
type CompletionChain struct {
	ID             string            `json:"id"`
	SourceTaskID   string            `json:"source_task_id"`
	SourceTaskName string            `json:"source_task_name,omitempty"`
	TargetTaskID   string            `json:"target_task_id"`
	TargetTaskName string            `json:"target_task_name,omitempty"`
	OnOutcome      CompletionOutcome `json:"on_outcome"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// CompletionDelivery is the durable processing record for one chain and one
// immutable source run.
type CompletionDelivery struct {
	ID           string        `json:"id"`
	ChainID      string        `json:"chain_id"`
	SourceTaskID string        `json:"source_task_id"`
	TargetTaskID string        `json:"target_task_id"`
	SourceRunID  string        `json:"source_run_id"`
	State        DeliveryState `json:"state"`
	Attempts     int           `json:"attempts"`
	CreatedAt    time.Time     `json:"created_at"`
	ClaimedAt    *time.Time    `json:"claimed_at,omitempty"`
	CompletedAt  *time.Time    `json:"completed_at,omitempty"`
	TargetRunID  string        `json:"target_run_id,omitempty"`
	Resolution   string        `json:"resolution,omitempty"`
}

// LogRecord is a single structured log entry surfaced in the GUI Logs view and
// persisted to the on-disk log file. It is produced by the daemon's slog handler
// (see internal/logbus); it is not stored in SQLite. Attrs holds the structured
// slog attributes that form the "cause/context" detail.
type LogRecord struct {
	ID       string         `json:"id"`
	Time     time.Time      `json:"time"`
	Severity AlertSeverity  `json:"severity"` // info | warning | error
	Source   string         `json:"source,omitempty"`
	Message  string         `json:"message"`
	TaskID   string         `json:"task_id,omitempty"`
	RunID    string         `json:"run_id,omitempty"`
	Attrs    map[string]any `json:"attrs,omitempty"`
}

// Alert is a surfaced condition shown in the GUI and reflected in logs.
type Alert struct {
	ID     string `json:"id"`
	TaskID string `json:"task_id,omitempty"`
	// RunID correlates a run failure to the exact persisted run. Other and
	// legacy alerts leave it empty.
	RunID        string        `json:"run_id,omitempty"`
	Severity     AlertSeverity `json:"severity"`
	Kind         AlertKind     `json:"kind"`
	Message      string        `json:"message"`
	CreatedAt    time.Time     `json:"created_at"`
	Acknowledged bool          `json:"acknowledged"`
}
