# Contract: External Triggers

## Local API

- `POST /v1/triggers` creates a trigger and returns its raw key once.
- `GET /v1/triggers` lists redacted triggers with target and readiness information.
- `GET /v1/triggers/{id}` returns one redacted trigger.
- `PATCH /v1/triggers/{id}` changes the name or target.
- `DELETE /v1/triggers/{id}` deletes the trigger.
- `POST /v1/triggers/{id}/enable` enables the trigger.
- `POST /v1/triggers/{id}/disable` disables the trigger.
- `POST /v1/triggers/{id}/rotate` replaces the key and returns the new raw key once.
- `POST /v1/triggers/{id}/reveal` explicitly returns the current raw key.
- `POST /v1/triggers/fire` accepts `{"key":"gst_..."}` and dispatches the target task once.

Ordinary trigger resources never contain the raw key. Sensitive responses include the key only for create, rotate, or reveal.

## Fire errors

| HTTP status | Code | Meaning |
|---|---|---|
| 404 | `trigger_unknown` | No trigger matches the supplied key |
| 409 | `trigger_disabled` | The trigger is disabled |
| 409 | `trigger_target_missing` | The target task no longer exists |
| 409 | `trigger_command_incomplete` | The target has no runnable command |
| 409 | `trigger_task_inactive` | The target is not active |
| 409 | `trigger_task_disabled` | The target is disabled |
| 409 | `trigger_group_blocked` | An ancestor group disables execution |
| 503 | `trigger_dispatch_unavailable` | The scheduler cannot accept the run |

Errors never echo the supplied key.

## CLI

The `gosched trigger` command supports `create`, `list`, `show`, `update`, `enable`, `disable`, `rotate`, `rm`, and `fire`. `show` is redacted unless `--reveal-key` is explicit. The external invocation is `gosched trigger fire <key>`.

## Desktop GUI

The Definitions navigation section adds Triggers after Chains. The table displays Trigger, Target, Enabled, Readiness, and ID. Users can create, edit, reveal or copy a key, copy a complete fire command, enable, disable, rotate, and delete triggers. Rotation and deletion require confirmation, and lifecycle changes refresh through the existing event stream.
