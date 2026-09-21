# Contract: System Summary

## Request

```http
GET /v1/system-summary
Authorization: Bearer <credential>
Accept: application/json
```

- Local IPC uses the same path without exposing bearer material.
- Remote credentials require Observe authority. Operate and Manage are not required.
- The operation is read-only and creates no audit mutation event.

## Success response

```json
{
  "schema": "go-schedule.system-summary.v1",
  "observed_at": "2026-09-21T15:00:00Z",
  "active_task_count": 3,
  "next_occurrence": {
    "task_id": "task-1",
    "task_name": "Daily export",
    "scheduled_for": "2026-09-21T16:00:00Z"
  },
  "recent_failure_count": 1,
  "recent_failure": {
    "run_id": "run-8",
    "task_id": "task-2",
    "task_name": "Archive",
    "ended_at": "2026-09-21T12:30:00Z"
  },
  "unacknowledged_alert_count": 1,
  "unacknowledged_alert": {
    "alert_id": "alert-4",
    "task_id": "task-2",
    "run_id": "run-8",
    "severity": "error",
    "kind": "run_failure",
    "created_at": "2026-09-21T12:30:01Z"
  },
  "notification_problem_count": 1,
  "notification_problem": {
    "delivery_id": "delivery-3",
    "task_id": "task-2",
    "run_id": "run-8",
    "channel_name": "Operations webhook",
    "state": "failed",
    "created_at": "2026-09-21T12:30:02Z"
  }
}
```

All timestamps are RFC 3339 UTC values. Optional representative objects are omitted when their count is zero. Unknown additive fields must be ignored by clients.

## Error responses

| Status | Meaning |
|---|---|
| 401 | Credential missing, unknown, revoked, or invalid |
| 403 | Credential lacks Observe authority |
| 404 | Older daemon does not implement the endpoint; desktop classifies the target as incompatible for this overview |
| 405 | Method other than GET |
| 500 | Bounded summary query failed; response contains a safe diagnostic only |

Transport, trust, pinned identity, manifest compatibility, and five-second timeout failures remain client-side typed connection states.

## Window semantics

- `observed_at` is captured once from the daemon clock for the whole response.
- Recent failures and notification problems use `[observed_at - 24h, observed_at]`.
- Upcoming work uses `(observed_at, observed_at + 24h]`.
- Counts include every matching row in the window.
- Each representative is the nearest upcoming or newest recent matching record.

## Security boundary

The response must never include task commands, arguments, working directories, environments, stdin, run output, alert messages, notification endpoints, authorization material, payloads, certificate bodies, or credential identifiers.
