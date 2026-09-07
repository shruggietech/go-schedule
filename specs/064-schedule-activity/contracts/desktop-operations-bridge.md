# Contract: Desktop Operations Bridge

## Schedule

### `ScheduleWindow(days) OperationResult`

Loads occurrences from one day before the current time through the requested 1-day, 7-day, or 30-day future window. Invalid window values are rejected. The successful result contains a complete ordered ScheduleSnapshot.

## Activity

### `ActivityWorkspace() OperationResult`

Loads authoritative active executions, up to 200 recent persisted runs, 200 recent daemon logs, and 200 recent alerts. Active identities are assigned before execution, published through the run event stream, and retained by the persisted record at completion. A successful result contains one complete ActivityWorkspace. Any failed constituent read returns unavailable without a partial workspace.

### `AcknowledgeAlert(id) OperationResult`

Acknowledges one exact non-empty alert identifier, then returns a complete refreshed ActivityWorkspace. Duplicate frontend activation is suppressed while pending.

### `AcknowledgeAlerts(ids) OperationResult`

Acknowledges each distinct non-empty identifier from the currently visible filtered view, then returns a complete refreshed ActivityWorkspace. An empty list succeeds without a mutation.

## Shared guarantees

- Every backend call has a bounded timeout.
- Errors use a small safe vocabulary and do not reflect arbitrary daemon content.
- Schedule predictions remain distinct from recorded runs.
- Activity runs, logs, and alerts remain distinct typed collections.
- Exact daemon log-path metadata is forwarded without probing or fallback.
- The frontend preserves only the last complete successful snapshot.
- Relevant domain events are debounced and stale async responses are discarded.
- Mutation controls are disabled while disconnected or pending.
