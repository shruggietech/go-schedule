# Data Model: Dependable Webhook Notifications

## NotificationChannel

Represents one reusable webhook destination.

| Field | Contract |
|---|---|
| `id` | Stable UUID primary key |
| `name` | Required trimmed display name |
| `kind` | `webhook` in S067 |
| `endpoint` | Required write-only validated URL, never serialized |
| `endpoint_summary` | Derived scheme and host without user information, path, query, or fragment |
| `authorization` | Optional write-only header value, never serialized |
| `has_authorization` | Derived response boolean |
| `enabled` | Controls creation of new deliveries |
| `created_at`, `updated_at` | UTC RFC 3339 timestamps |

## NotificationAssignment

Connects a channel to exactly one task or group scope.

| Field | Contract |
|---|---|
| `id` | Stable UUID primary key |
| `channel_id` | Existing notification channel |
| `scope_type` | `task` or `group` |
| `scope_id` | Existing task or group identity |
| `on_success` | Include successful terminal runs |
| `on_failure` | Include failed terminal runs |
| `created_at`, `updated_at` | UTC RFC 3339 timestamps |

The pair `(scope_type, scope_id, channel_id)` is unique. At least one condition must be true. Deleting the source scope or channel deletes the assignment.

## EffectiveNotificationPolicy

Computed, not persisted.

| Field | Contract |
|---|---|
| `task_id` | Task being evaluated |
| `source_scope_type` | `task`, `group`, or `none` |
| `source_scope_id` | Selected task or nearest ancestor group |
| `assignments` | Complete assignments from only the selected scope |

## NotificationDelivery

Represents durable work and redacted evidence distinct from a run.

| Field | Contract |
|---|---|
| `id` | Stable UUID and receiver idempotency key |
| `channel_id` | Nullable live-channel reference |
| `channel_name` | Immutable safe snapshot |
| `destination_summary` | Immutable safe scheme-and-host snapshot |
| `endpoint` | Write-only attempt snapshot, cleared at terminal state |
| `authorization` | Write-only attempt snapshot, cleared at terminal state |
| `event_kind` | `run.completed` or `test` |
| `task_id`, `run_id` | Nullable correlations; present for run completion |
| `task_name`, `group_id`, `group_name` | Immutable safe source snapshots |
| `payload` | Immutable JSON snapshot containing no protected fields |
| `state` | `pending`, `claimed`, `succeeded`, or `failed` |
| `attempts` | Total attempts already started, maximum three |
| `next_attempt_at` | Earliest UTC eligibility for pending retry |
| `created_at`, `claimed_at`, `completed_at` | UTC lifecycle timestamps |
| `last_status` | Last HTTP response status, zero for transport failures |
| `last_error` | Sanitized bounded diagnostic, never a URL or authorization value |

## State Transitions

```text
pending -> claimed -> succeeded
pending -> claimed -> pending
pending -> claimed -> failed
claimed -> pending (daemon restart, attempts preserved)
pending or claimed -> deleted (channel removal only)
```

The claimed-to-pending retry transition increments no counter because claiming already incremented `attempts`. A claim at attempt three must transition to succeeded or failed, never pending. Every terminal transition clears `endpoint` and `authorization`, then prunes terminal history beyond 1,000 newest records.

## Migration and Removal

Migration v15 creates all three tables and indexes without rewriting existing groups, tasks, runs, completion chains, triggers, watchers, or alerts. Existing databases therefore begin with no notification behavior. Channel removal is transactional: assignments and non-terminal deliveries are deleted, terminal delivery `channel_id` values become null, the channel row holding its secret is deleted, and safe terminal snapshots remain.
