---
title: Local API
nav_order: 5.5
---

# Local API

**Audience:** client and integration authors\
**Applies to:** the current unreleased local API contract\
**Transport:** local Unix socket or Windows named pipe, never a public TCP port

The CLI and desktop app use the same versioned JSON API hosted by `goschedd`. Errors use `{"error":{"code":"...","field":"...","message":"..."}}`.

## External triggers

An ordinary trigger representation contains `id`, `name`, optional `set_id`, optional `set_name`, optional `set_position`, `target_task_id`, `target_task_name`, `enabled`, `readiness`, `reason`, `created_at`, and `updated_at`. It never contains the raw key. A set member cannot be retargeted individually; use the set endpoint so every member retains the shared target invariant.

| Method | Path | Result |
| --- | --- | --- |
| `GET` | `/v1/triggers` | `{"triggers": [...]}` with redacted records |
| `POST` | `/v1/triggers` | Create from name, target task, and optional enabled state; `201` with the key |
| `GET` | `/v1/triggers/{id}` | One redacted trigger |
| `PATCH` | `/v1/triggers/{id}` | Update name or target task |
| `DELETE` | `/v1/triggers/{id}` | Delete; `204` |
| `POST` | `/v1/triggers/{id}/enable` | Enable the trigger |
| `POST` | `/v1/triggers/{id}/disable` | Disable the trigger |
| `POST` | `/v1/triggers/{id}/rotate` | Atomically replace and return the key |
| `POST` | `/v1/triggers/{id}/reveal` | Explicitly return the current key |
| `POST` | `/v1/triggers/fire` | Accept `{"key":"gst_..."}` and submit one run request; `202` |

Fire failures use stable codes: `trigger_unknown`, `trigger_disabled`, `trigger_target_missing`, `trigger_command_incomplete`, `trigger_task_inactive`, `trigger_task_disabled`, `trigger_group_blocked`, or `trigger_dispatch_unavailable`. Error responses never echo the submitted key.

## Filesystem watchers

A filesystem watcher representation contains `id`, `name`, `kind`, `path`, optional `pattern`, `recursive`, string durations `debounce` and `stability`, target identity, enabled state, runtime `health`, readiness, reason, and timestamps. `kind` is `file` for one exact path or `directory` for basename-glob selection. A directory pattern defaults to `*`; recursive selection excludes linked directories. Durations accept Go duration syntax such as `250ms`, `2s`, or `1m` and must be from 25 milliseconds through one hour.

| Method | Path | Behavior |
|---|---|---|
| `GET` | `/v1/filesystem-watchers` | List definitions joined with current runtime health |
| `POST` | `/v1/filesystem-watchers` | Create a definition and reload observation; `201` |
| `GET` | `/v1/filesystem-watchers/{id}` | Show one definition and current health |
| `PATCH` | `/v1/filesystem-watchers/{id}` | Atomically update selection, timing, target, enabled state, or name and reload observation |
| `DELETE` | `/v1/filesystem-watchers/{id}` | Delete, cancel pending candidates, and reload observation; `204` |
| `POST` | `/v1/filesystem-watchers/{id}/enable` | Enable and reload observation |
| `POST` | `/v1/filesystem-watchers/{id}/disable` | Disable, cancel pending candidates, and reload observation |

Watcher lifecycle and health events use `kind: "filesystem_watcher"` with identity, name, verb, and optional health only. They omit configured and matched paths. Filesystem-originated runs have `trigger: "filesystem_watcher"` and `source_watcher_id`; they never retain the matched path.

### Trigger Sets

Ordinary Trigger Set representations include stable set identity, name, target, member and enabled counts, ordered redacted members, and timestamps. Create, reveal, and rotate responses additionally contain ordered member keys and complete commands.

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/v1/trigger-sets` | List redacted Trigger Sets and members |
| `POST` | `/v1/trigger-sets` | Atomically create 1 through 99 members and return ordered secrets |
| `GET` | `/v1/trigger-sets/{id}` | Show one redacted Trigger Set |
| `PATCH` | `/v1/trigger-sets/{id}` | Atomically retarget every member |
| `DELETE` | `/v1/trigger-sets/{id}` | Atomically delete the set and every member; `204` |
| `POST` | `/v1/trigger-sets/{id}/enable` | Atomically enable every member |
| `POST` | `/v1/trigger-sets/{id}/disable` | Atomically disable every member |
| `POST` | `/v1/trigger-sets/{id}/rotate` | Atomically rotate every key and return ordered replacement secrets |
| `POST` | `/v1/trigger-sets/{id}/reveal` | Explicitly return current ordered secrets |

## Runtime storage information

`GET /v1/runtime-info` returns the absolute effective paths used by the running daemon: `data_dir`, `database_path`, optional `config_path`, `log_path`, and `lock_path`. Desktop clients use this endpoint for read-only storage disclosure, including when the daemon was launched with a custom configuration path.

## Completion chains

The chain representation contains `id`, `source_task_id`, `source_task_name`, `target_task_id`, `target_task_name`, `on_outcome`, `created_at`, and `updated_at`.

| Method | Path | Result |
| --- | --- | --- |
| `GET` | `/v1/chains` | `{"chains": [...]}` |
| `POST` | `/v1/chains` | Create from source, target, and outcome; `201` |
| `GET` | `/v1/chains/{id}` | One chain |
| `PATCH` | `/v1/chains/{id}` | Update any non-empty subset of mutable fields |
| `DELETE` | `/v1/chains/{id}` | Delete; `204` |

Create example:

```text
{
  "source_task_id": "source-id",
  "target_task_id": "target-id",
  "on_outcome": "success"
}
```

`on_outcome` is `success`, `failure`, or `any`. Missing resources return `not_found`. Invalid outcomes, self-links, duplicates, and direct or indirect cycles return `validation_failed` without partial mutation.

## Correlated history and events

Completion-triggered entries from `GET /v1/runs` have `trigger` set to `completion` and include optional `source_task_id` and `source_run_id`. Externally triggered entries have `trigger` set to `external_trigger` and include `source_trigger_id`. Filesystem-triggered entries have `trigger` set to `filesystem_watcher` and include `source_watcher_id`. Raw keys and matched file paths are never stored in run history.

`GET /v1/events` emits `kind: "chain"` with a `created`, `updated`, or `deleted` verb. Create and update include the current chain; delete carries its stable ID.

Trigger lifecycle events use `kind: "trigger"` and contain only a redacted trigger or its stable deletion ID.

Trigger Set lifecycle events use `kind: "trigger_set"` and contain set identity, name, member count, and verb without member keys. One set-level mutation publishes one event after its transaction commits.

## Webhook notifications

Notification channels are reusable write-only webhook destinations. Ordinary channel responses contain `id`, `name`, `kind`, `endpoint_summary`, `has_authorization`, `enabled`, `created_at`, and `updated_at`. See [Webhook notifications](notifications.md) for the receiver payload, precedence, retry, duplicate, and security contracts.

| Method | Path | Result |
| --- | --- | --- |
| `GET` | `/v1/notification-channels` | `{"notification_channels": [...]}` with redacted metadata |
| `POST` | `/v1/notification-channels` | Create from `name`, `endpoint`, optional `authorization`, and optional `enabled`; `201` |
| `GET` | `/v1/notification-channels/{id}` | One redacted channel |
| `PATCH` | `/v1/notification-channels/{id}` | Update optional `name`, `endpoint`, or `enabled` |
| `DELETE` | `/v1/notification-channels/{id}` | Remove assignments and unfinished work while preserving safe terminal history; `204` |
| `POST` | `/v1/notification-channels/{id}/enable` | Enable new delivery creation |
| `POST` | `/v1/notification-channels/{id}/disable` | Disable new delivery creation |
| `POST` | `/v1/notification-channels/{id}/rotate` | Replace or clear the write-only `authorization` value |
| `POST` | `/v1/notification-channels/{id}/test` | Queue a transport-equivalent test delivery; `202` |
| `GET` | `/v1/tasks/{id}/notifications` | List direct task assignments |
| `PUT` | `/v1/tasks/{id}/notifications` | Atomically replace direct task assignments |
| `GET` | `/v1/tasks/{id}/notifications/effective` | Return selected source scope and complete effective assignments |
| `GET` | `/v1/groups/{id}/notifications` | List direct group assignments |
| `PUT` | `/v1/groups/{id}/notifications` | Atomically replace direct group assignments |
| `GET` | `/v1/notification-deliveries` | List redacted evidence with optional `channel`, `task`, `run`, `state`, and `limit` filters |

An assignment contains `channel_id`, `on_success`, and `on_failure`; at least one outcome must be true and one channel may occur only once per scope. The nearest non-empty task or group scope replaces all more distant assignments. Delivery responses remain distinct from run responses and include a safe parsed `event`, attempts, timestamps, last status, and bounded diagnostic without protected endpoint or authorization fields.
