# Contract: Task Creation and Failed-Run Diagnostics

## Compatibility boundary

All additions are optional JSON fields or a new read endpoint under the existing local `/v1` API. Existing request and response fields retain their meanings.

## Create task

`POST /v1/tasks`

### Additive request field

```json
{
  "enabled": false
}
```

| Input | Result |
| --- | --- |
| field omitted | Task is created enabled, preserving the pre-S051 contract |
| `false` | Task is atomically created disabled |
| `true` | Task is atomically created enabled |

The returned task contains the authoritative persisted `enabled` value. The server emits only the final created state and performs one scheduler reload.

## Get exact run

`GET /v1/runs/{id}`

### Success

Status `200` with one run object:

```json
{
  "id": "run-id",
  "task_id": "task-id",
  "scheduled_for": "2026-09-05T00:00:00Z",
  "outcome": "failure",
  "exit_code": 7,
  "output": "combined process output",
  "output_truncated": false,
  "trigger": "manual"
}
```

### Not found

Status `404` with the existing `not_found` error envelope. The client must not substitute another run.

## Alert response and events

Newly created `run_failed` alerts include:

```json
{
  "task_id": "task-id",
  "run_id": "run-id",
  "kind": "run_failed"
}
```

Other and legacy alerts omit `run_id`. List and live-event representations use the same additive field.

## Output contract

- `output` contains at most the configured byte cap of combined stdout/stderr.
- `output_truncated` is true if any additional bytes were discarded.
- Empty output and unavailable output are distinct presentation states.
- The diagnostic surface does not add arguments, standard input, environment values, or reconstructed command lines.
