# Data Model: Actionable Notification Conditions

## NotificationAssignment additions

| Field | Contract |
|---|---|
| `failure_threshold` | 1 through 100 when `on_failure` is true; existing assignments migrate to 1 |
| `on_failure_to_start` | Immediate problem condition for structured process-start failure |
| `duration_threshold_seconds` | 0 disabled; otherwise 1 through 2,592,000 seconds |
| `on_recovery` | Emit one recovery transition after an active configured problem clears |
| `reminder_interval_seconds` | 0 disabled; otherwise 60 through 2,592,000 seconds |
| `quiet_period_seconds` | 0 disabled; otherwise 60 through 2,592,000 seconds |

At least one of success, failure, failure-to-start, or duration exceeded must be enabled. Recovery and reminders require at least one problem condition.

## NotificationConditionState

| Field | Contract |
|---|---|
| `task_id`, `channel_id` | Composite primary key with cascading task and channel lifecycle |
| `policy_fingerprint` | Deterministic digest of effective assignment identity and settings |
| `consecutive_failures` | Current streak under this fingerprint |
| `active_condition` | Empty, `failure_to_start`, `consecutive_failure`, or `duration_exceeded` |
| `active_since` | First UTC observation of the active primary problem |
| `last_notified_at` | Last created condition delivery time |
| `last_evaluated_run_id` | Last source run applied to this state |
| `updated_at` | Last UTC evaluation time |

Policy fingerprint mismatch resets all state before evaluating the current run. Run identity provides duplicate protection in addition to the run primary key.

## Run addition

`start_failed` is a structured boolean set only when operating-system process creation fails. It is persisted with the run and returned through existing run APIs without exposing the raw error or command.

## NotificationDelivery additions

| Field | Contract |
|---|---|
| `condition_kind` | Empty for transport tests, otherwise `success`, `failure`, `failure_to_start`, `duration_exceeded`, `recovery`, or `daemon_health` |
| `condition_summary` | Generated bounded safe explanation, maximum 512 bytes |
| `event_kind` | Existing `run.completed` and `test`, plus `daemon.health` |

## DaemonHealthState

| Field | Contract |
|---|---|
| `channel_id` | Primary key and cascading channel reference |
| `next_due_at` | Persisted UTC eligibility time |
| `last_delivery_id` | Most recently created heartbeat delivery |
| `updated_at` | Last scheduling update |

Channels add `health_interval_seconds`, where zero disables heartbeats and 60 through 86,400 enables them.

## Task condition transitions

```text
inactive + problem              -> active + first problem delivery
active + same problem early     -> active + no delivery
active + same problem due       -> active + reminder delivery
active + different problem      -> active(new) + first problem delivery
active + no problem             -> inactive + optional recovery delivery
inactive + routine success      -> inactive + optional success delivery
policy fingerprint changed      -> reset, then evaluate current run
```

## Heartbeat transitions

```text
disabled                        -> no state and no delivery
enabled with no state           -> create due heartbeat, set next due
enabled before next due         -> no delivery
enabled at or after next due    -> create one heartbeat if none is non-terminal, advance next due
restart                         -> preserve next due and pending work
```

Migration v20 adds assignment, channel, run, condition-state, and heartbeat-state fields. It rebuilds `notification_deliveries` to expand the checked event vocabulary and copies every existing row unchanged with empty condition metadata.
