# Local Notification API Contract

All endpoints use the existing local authenticated IPC transport and error envelope. Ordinary responses never contain `endpoint`, `authorization`, or any secret reference.

## Channels

- `POST /v1/notification-channels` creates a webhook channel from `name`, `endpoint`, optional `authorization`, and optional `enabled`.
- `GET /v1/notification-channels` lists redacted channels.
- `GET /v1/notification-channels/{id}` returns one redacted channel.
- `PATCH /v1/notification-channels/{id}` updates optional `name`, `endpoint`, or `enabled`; omitted secrets remain unchanged.
- `POST /v1/notification-channels/{id}/enable` and `/disable` change the enabled state.
- `POST /v1/notification-channels/{id}/rotate` replaces authorization from required `authorization`; an empty value removes authorization.
- `POST /v1/notification-channels/{id}/test` creates a `test` delivery and returns `202 Accepted` with the redacted delivery.
- `DELETE /v1/notification-channels/{id}` applies the removal contract.

A channel response contains `id`, `name`, `kind`, `endpoint_summary`, `has_authorization`, `enabled`, `created_at`, and `updated_at`.

## Assignments and Effective Policy

- `GET /v1/tasks/{id}/notifications` and `GET /v1/groups/{id}/notifications` list assignments configured directly on that scope.
- `PUT /v1/tasks/{id}/notifications` and `PUT /v1/groups/{id}/notifications` atomically replace that scope's complete assignment list.
- `GET /v1/tasks/{id}/notifications/effective` returns `task_id`, selected `source_scope_type`, selected `source_scope_id`, and the complete effective assignment list.

Each replacement item contains `channel_id`, `on_success`, and `on_failure`. The list may be empty to resume inheritance. Duplicate channel IDs, missing channels, or items with neither condition produce `validation_failed`.

## Deliveries

`GET /v1/notification-deliveries` accepts optional `channel`, `task`, `run`, and `state` filters plus `limit` from 1 through 1,000. The response contains redacted delivery objects ordered newest first. A delivery object never includes its internal endpoint, authorization, or payload storage fields, but includes the safe parsed payload as `event` for diagnostic correlation.

## Error Mapping

- Invalid input or filter: `400 validation_failed`
- Missing channel, task, group, or run: `404 not_found`
- Duplicate assignment or conflicting state change: `409 conflict`
- Persistence or dispatch setup failure: `500 internal`
