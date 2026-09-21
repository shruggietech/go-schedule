# Contract: Daemon Search

## Request

```http
GET /v1/search?q=archive&kind=task&kind=failure&limit=50
Authorization: Bearer <credential>
Accept: application/json
```

- Local IPC uses the same path without exposing bearer material.
- `q` is required after trimming, accepts at most 200 UTF-8 characters, and matches case-insensitively.
- `kind` may repeat and accepts task, group, failure, schedule, or alert. Omission searches all kinds.
- `limit` defaults to 50 and accepts one through 50.
- Observe authority is sufficient. The operation is read-only and creates no mutation audit event.

## Success response

```json
{
  "schema": "go-schedule.daemon-search.v1",
  "query": "archive",
  "observed_at": "2026-09-21T18:00:00Z",
  "truncated": false,
  "results": [
    {
      "kind": "task",
      "object_id": "task-2",
      "task_id": "task-2",
      "name": "Archive",
      "context": "Active, enabled, recurring",
      "enabled": true,
      "action_hints": ["disable", "run_now"]
    },
    {
      "kind": "failure",
      "object_id": "run-8",
      "task_id": "task-2",
      "name": "Archive",
      "context": "Failed run",
      "occurred_at": "2026-09-21T12:30:00Z",
      "action_hints": []
    }
  ]
}
```

- Results sort by kind, case-folded name, occurred time descending when present, and object ID.
- The daemon queries up to `limit + 1` safe matches to calculate `truncated`, then returns at most `limit`.
- Schedule matches are the nearest upcoming occurrence for matching active scheduled tasks in the next 30 days.
- Failure matches include failed runs only. Alert matches include unacknowledged alerts only.
- Unknown additive fields and action hints are ignored by older clients.

## Error responses

| Status | Meaning |
|---|---|
| 400 | Query, kind, or limit is invalid |
| 401 | Credential is missing, invalid, unknown, or revoked |
| 403 | Credential lacks Observe authority |
| 404 | Older daemon does not implement search |
| 405 | Method other than GET |
| 500 | A bounded search query or occurrence calculation failed |

## Security boundary

The response never includes task commands, arguments, working directories, environments, stdin, run output, alert messages, credentials, credential identifiers, notification payloads, certificate bodies, or secret-bearing endpoints.
