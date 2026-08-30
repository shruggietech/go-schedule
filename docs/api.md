---
title: Local API
nav_order: 5.5
---

# Local API

**Audience:** client and integration authors\
**Applies to:** the unreleased task-completion chaining slice\
**Transport:** local Unix socket or Windows named pipe, never a public TCP port

The CLI and desktop app use the same versioned JSON API hosted by `goschedd`.
Errors use `{"error":{"code":"...","field":"...","message":"..."}}`.

## Completion chains

The chain representation contains `id`, `source_task_id`, `source_task_name`,
`target_task_id`, `target_task_name`, `on_outcome`, `created_at`, and
`updated_at`.

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

`on_outcome` is `success`, `failure`, or `any`. Missing resources return
`not_found`. Invalid outcomes, self-links, duplicates, and direct or indirect
cycles return `validation_failed` without partial mutation.

## Correlated history and events

Completion-triggered entries from `GET /v1/runs` have `trigger` set to
`completion` and include optional `source_task_id` and `source_run_id`. The
fields are absent for older and non-completion history.

`GET /v1/events` emits `kind: "chain"` with a `created`, `updated`, or `deleted`
verb. Create and update include the current chain; delete carries its stable ID.
