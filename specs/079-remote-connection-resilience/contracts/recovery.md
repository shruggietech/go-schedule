# Remote Recovery Contract

## Connection snapshot projection

The desktop bridge extends its existing snapshot with:

```json
{
  "stale": true,
  "retryAttempt": 2,
  "nextRetryAt": "2026-09-10T01:02:03Z",
  "recovery": "automatic"
}
```

- `stale` is true after a previously connected remote target loses current connectivity and remains true until validated reconnection.
- `retryAttempt` is positive only while an automatic retry is scheduled or running.
- `nextRetryAt` is present only while a retry wait is scheduled.
- `recovery` is `none`, `automatic`, or `manual`.
- Existing fields and local IPC projections remain backward compatible.

## Failure categories

| Category | Automatic retry | Required guidance |
| --- | --- | --- |
| unreachable | Yes | Waiting for endpoint or network recovery, with manual retry available |
| timed out | Yes | Waiting for endpoint or network recovery, with manual retry available |
| stream interrupted | Yes | Last data is stale while live updates recover |
| unauthorized or revoked | No | Repair the credential or ask an administrator to restore access |
| forbidden | No | Ask an administrator for the required grant |
| incompatible | No | Update the client or daemon |
| certificate changed | No | Inspect and explicitly repair trust for the same daemon identity |
| identity changed | No | Stop and create or select the intended target; do not overwrite the pin |

## Mutation contract

Every deliberate daemon-backed mutation produces at most one outbound request. If a remote transport error occurs before an authoritative response is received, the result outcome is `uncertain` and its message states that the operation may have completed. The desktop preserves the user's working state, disables further remote mutation while disconnected, and refreshes the relevant authoritative workspace after recovery.

No client layer automatically retries POST, PUT, PATCH, or DELETE requests. No request body, credential, pairing phrase, or raw transport error is included in the bridge result.

## Refresh-before-stream contract

Each successful recovered connection creates a new generation and publishes a connected snapshot. Mounted feature stores reload their bounded authoritative workspace for that generation. Domain invalidation events from an older generation are rejected. Events from the recovered generation may prompt further reads only after the connected snapshot has made the feature available.
