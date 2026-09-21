# Contract: Desktop Search and Actions

## Search bridge

```text
SearchAcrossSystems(request: SearchRequest) -> SearchSnapshot
event: search:event -> SearchSnapshot
```

The initial snapshot contains one searching observation per current registration. Progressive snapshots replace observations only within the same generation. The final snapshot excludes registrations removed while the search was active.

## Action bridge

```text
ExecuteSearchAction(intent: ActionIntent) -> ActionBatchResult
```

The desktop validates that:

- the action is one of acknowledge, enable, disable, or run-now;
- every selection supports the same action and has a registration key, expected daemon identity, kind, and object identity;
- the current profile still exists and resolves its stored credential;
- health and identity checks produce the expected daemon identity and current Operate authority;
- the exact task or alert still exists and remains compatible with the action.

## Batch result

```json
{
  "action": "run_now",
  "outcome": "partial",
  "message": "2 of 3 requested objects accepted. Review the target-specific outcomes.",
  "outcomes": [
    {
      "registrationKey": "local",
      "expectedDaemonId": "daemon-a",
      "kind": "task",
      "objectId": "task-1",
      "displayName": "Daily export",
      "outcome": "accepted",
      "message": "Run request accepted by This computer.",
      "currentDaemonId": "daemon-a"
    },
    {
      "registrationKey": "profile-b",
      "expectedDaemonId": "daemon-b",
      "kind": "task",
      "objectId": "task-9",
      "displayName": "Daily export",
      "outcome": "rejected",
      "message": "The selected task no longer exists on Remote B.",
      "currentDaemonId": "daemon-b"
    }
  ]
}
```

Top-level outcomes are accepted, partial, rejected, or unavailable. They summarize only the included per-object outcomes and never imply transactionality or rollback.

## Retry and uncertainty

- Enable, disable, and acknowledgement use their existing idempotent endpoints.
- Run-now uses the existing remote request-identity behavior.
- An uncertain result remains uncertain until a refreshed search or Activity view provides current evidence. The desktop does not retry it automatically.
