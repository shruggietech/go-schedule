# Data Model: Production Wails Shell and Local Connection

## Connection State

Closed values: `connecting`, `connected`, `degraded`, `recovering`, `unavailable`, `access_denied`, `incompatible`, and `timed_out`.

| Field | Type | Rule |
| --- | --- | --- |
| state | enum | One closed connection state |
| title | string | Approved short operator-facing label |
| message | string | Sanitized explanation with no raw transport value |
| action | enum | `none`, `retry`, `open_service_help`, `refresh_login`, or `update_daemon` |
| retryable | boolean | True only when manual retry is meaningful |

## Desktop Target

| Field | Type | Rule |
| --- | --- | --- |
| id | string | Stable adapter identity; local target is `local` |
| displayName | string | Local target is `This computer` |
| platform | enum | `windows`, `macos`, or `linux` |
| daemonVersion | string | Sanitized version returned by health negotiation |
| capabilities | ordered string list | Closed manifest of feature families supported by the adapter |
| permissions | ordered string list | Closed manifest of effective desktop authority |

## Connection Snapshot

| Field | Type | Rule |
| --- | --- | --- |
| generation | unsigned integer | Increases for every connection attempt |
| revision | unsigned integer | Increases for every committed state transition so asynchronous same-generation snapshots remain ordered |
| target | Desktop Target | Always present, including unavailable state |
| connection | Connection State | Current safe status and action |
| lastSuccessfulContact | optional RFC 3339 timestamp | Updated only after compatible health succeeds |

Snapshots are immutable values. Consumers order first by generation and then by revision, and never merge an older value over a newer one.

## Desktop Event

| Field | Type | Rule |
| --- | --- | --- |
| id | string | Stable kind, generation, and monotonic sequence |
| generation | unsigned integer | Must equal the manager's active generation when emitted |
| kind | enum | `connection.changed` or sanitized daemon-domain family plus `.changed` |
| message | string | Approved bounded status text |
| occurredAt | RFC 3339 timestamp | UTC |
| taskId | optional string | Stable task identity only for task events |

## Retry Schedule

| Attempt class | Delay |
| --- | --- |
| first transient retry | 250 milliseconds |
| second transient retry | 1 second |
| subsequent transient retries | 5 seconds maximum |

Manual retry cancels the active attempt or wait and immediately creates a new generation. A successfully received stream event resets the retry schedule.

## Shell Route

| Field | Type | Rule |
| --- | --- | --- |
| id | enum | `tasks`, `automations`, `schedule`, `activity`, or `settings` |
| label | string | Unique accessible navigation label |
| description | string | Concise page purpose |
| delivery | string | Honest placeholder naming the owning follow-up issue until implemented |

## Appearance Preference

Closed values: `system`, `light`, `dark`. S061 keeps the value in the current application session. System mode follows the operating-system color preference through media queries.

## Component Contract

Each reusable component defines semantic element or role, required accessible name, keyboard behavior, focus behavior, state props, safe content boundary, and one direct behavior or accessibility test. Components do not fetch connection data or call Wails directly.

## State Transitions

```text
start -> connecting
connecting -> connected | unavailable | access_denied | incompatible | timed_out
connected -> degraded | recovering | shutdown
degraded -> recovering | connected | shutdown
unavailable | timed_out -> recovering | shutdown
recovering -> connected | unavailable | access_denied | incompatible | timed_out | degraded
access_denied | incompatible -> connecting by explicit action | shutdown
any active state -> shutdown
```

Only the manager loop commits transitions. Results and events from older generations are discarded.
