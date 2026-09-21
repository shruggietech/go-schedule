# Contract: Actionable Notification Conditions

## Assignment input

```json
{
  "channel_id": "channel-id",
  "on_success": false,
  "on_failure": true,
  "failure_threshold": 3,
  "on_failure_to_start": true,
  "duration_threshold_seconds": 900,
  "on_recovery": true,
  "reminder_interval_seconds": 1800,
  "quiet_period_seconds": 600
}
```

Replacing task or group assignments remains atomic. Existing clients that send only `on_success` and `on_failure` retain threshold-one failure behavior.

## Channel heartbeat input

Channel create and update accept `health_interval_seconds`. Zero disables health heartbeats. Values from 60 through 86,400 enable them.

## Condition evidence

Notification deliveries add optional `condition_kind` and `condition_summary`. Run webhook events add this optional object:

```json
{
  "condition": {
    "kind": "consecutive_failure",
    "summary": "Failure threshold reached after 3 consecutive failures.",
    "streak": 3,
    "threshold": 3,
    "reminder": false
  }
}
```

Condition kinds are `success`, `failure`, `consecutive_failure`, `failure_to_start`, `duration_exceeded`, `recovery`, and `daemon_health`.

## Daemon-health webhook

```json
{
  "schema": "go-schedule.webhook.v1",
  "event": "daemon.health",
  "delivery": {
    "id": "delivery-id",
    "created_at": "2026-09-21T12:00:00Z"
  },
  "daemon": {
    "version": "1.4.0",
    "status": "healthy",
    "next_expected_at": "2026-09-21T12:05:00Z"
  },
  "task": null
}
```

The receiver uses delivery ID for idempotency and treats missing expected heartbeats according to its own monitoring policy.
