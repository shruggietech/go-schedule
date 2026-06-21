# API Contract: Logs endpoint, event stream changes, Triggers removal

Local HTTP/JSON API served over IPC (Unix socket / named pipe). Error envelope unchanged
(`{"error":{code,field,message}}`).

## NEW: `GET /v1/logs`

Returns the most recent log records from the daemon's in-memory ring, newest first.

**Query parameters**:

| Param        | Type   | Default | Notes                                                        |
|--------------|--------|---------|-------------------------------------------------------------|
| `severity`   | string | (all)   | One of `info`,`warning`,`error`; filters server-side        |
| `limit`      | int    | ring N  | Max records to return (capped at ring size)                 |
| `since`      | RFC3339| (none)  | Only records at/after this time                             |

**200 response**:

```json
{
  "logs": [
    {
      "id": "00000000000123",
      "time": "2026-06-20T17:04:11Z",
      "severity": "error",
      "source": "executor",
      "message": "task run failed",
      "task_id": "tsk_abc",
      "run_id": "run_def",
      "attrs": { "exit_code": 1, "error": "exec: \"backup.sh\": file does not exist" }
    }
  ]
}
```

**Notes**: the ring is bounded; older history lives only in the on-disk log file
([log-file.md](log-file.md)). `severity=error` MUST return only error records (US3 acceptance #2).

## CHANGED: `GET /v1/events` (SSE stream)

Existing SSE stream gains new event kinds. Each SSE `data:` line is a JSON `events.Event`:

```jsonc
// existing
{ "kind": "run",   "run":   { /* domain.Run */ } }
{ "kind": "alert", "alert": { /* domain.Alert */ } }
// NEW
{ "kind": "log",   "log":   { /* LogRecord */ } }
{ "kind": "task",  "task":  { "verb": "created|updated|deleted", "task":  { /* domain.Task */ },  "id": "tsk_..." } }
{ "kind": "group", "group": { "verb": "created|updated|deleted", "group": { /* domain.Group */ }, "id": "grp_..." } }
```

**Guarantees**: delivery is best-effort/non-blocking (a slow client drops events, as today).
Clients MUST treat events as idempotent hints and dedupe by id; on reconnect the client does a
full reload before resuming (FR-024). For `deleted`, `task`/`group` object MAY be null and only
`id` is guaranteed.

**Publishing sites** (server-side): `log` from the logbus handler; `task`/`group` from the
corresponding API mutation handlers after a successful store write (create/update/delete/
enable/disable); `run`/`alert` unchanged.

## REMOVED routes (Triggers)

The following are deleted and now return the standard `404 not_found` envelope via the fallback
handler:

- `GET /v1/triggers`
- `POST /v1/triggers`
- `DELETE /v1/triggers/{id}`

The API client loses `CreateTrigger`/`ListTriggers`/`DeleteTrigger`; the GUI `Backend` interface
drops the same three methods.

## UNCHANGED (reference)

`/v1/alerts` and `/v1/alerts/{id}/ack` remain (alerts still power part of the Logs view).
`/v1/calendar` remains and is the data source for the new calendar view. `/v1/tasks`,
`/v1/groups`, `/v1/runs`, `/v1/schedules/preview`, `/v1/health` unchanged.
