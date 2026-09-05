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

An ordinary trigger representation contains `id`, `name`, `target_task_id`, `target_task_name`, `enabled`, `readiness`, `reason`, `created_at`, and `updated_at`. It never contains the raw key.

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

Completion-triggered entries from `GET /v1/runs` have `trigger` set to `completion` and include optional `source_task_id` and `source_run_id`. Externally triggered entries have `trigger` set to `external_trigger` and include `source_trigger_id`. Raw keys are never stored in run history.

`GET /v1/events` emits `kind: "chain"` with a `created`, `updated`, or `deleted` verb. Create and update include the current chain; delete carries its stable ID.

Trigger lifecycle events use `kind: "trigger"` and contain only a redacted trigger or its stable deletion ID.
