# Data Model: All Systems Operational Overview

## SystemRegistration

One independently saved or local target.

| Field | Type | Rules |
|---|---|---|
| key | string | `local` or the immutable saved profile ID; unique within a snapshot |
| kind | enum | `local` or `remote` |
| label | string | User-facing label; not an identity key |
| endpointSummary | string | Safe local designation or HTTPS origin; no credential material |
| profileID | string | Present only for saved profiles |
| daemonID | string | Last known pinned identity when available |
| platform | string | Last known or currently observed platform |
| architecture | string | Last known or currently observed architecture |
| version | string | Last known or currently observed daemon version |

## OperationalSummary

One immutable, daemon-produced observation.

| Field | Type | Rules |
|---|---|---|
| schema | string | `go-schedule.system-summary.v1` |
| observedAt | RFC 3339 timestamp | Produced from the daemon clock |
| activeTaskCount | integer | Non-negative count of enabled active tasks |
| nextOccurrence | optional UpcomingReference | Nearest occurrence in `(observedAt, observedAt + 24h]` |
| recentFailureCount | integer | Failures in `[observedAt - 24h, observedAt]` |
| recentFailure | optional FailureReference | Newest matching failed run |
| unacknowledgedAlertCount | integer | Current unacknowledged alert count |
| unacknowledgedAlert | optional AlertReference | Newest unacknowledged alert |
| notificationProblemCount | integer | Failed or retrying deliveries created in the previous 24 hours |
| notificationProblem | optional NotificationReference | Newest matching delivery |

The response contains no task commands, arguments, environment, stdin, run output, alert message body, delivery endpoint, authorization, event payload, certificate, or credential.

## Representative references

### UpcomingReference

| Field | Type | Rules |
|---|---|---|
| taskID | string | Stable task identifier |
| taskName | string | Safe display label |
| scheduledFor | RFC 3339 timestamp | Occurs inside the next 24 hours |

### FailureReference

| Field | Type | Rules |
|---|---|---|
| runID | string | Stable run identifier |
| taskID | string | Stable task identifier |
| taskName | string | Safe display label |
| endedAt | RFC 3339 timestamp | Inside the previous 24 hours |

### AlertReference

| Field | Type | Rules |
|---|---|---|
| alertID | string | Stable alert identifier |
| taskID | string | Optional task identifier |
| runID | string | Optional run identifier |
| severity | string | Existing alert severity |
| kind | string | Existing alert kind |
| createdAt | RFC 3339 timestamp | Observation metadata only |

### NotificationReference

| Field | Type | Rules |
|---|---|---|
| deliveryID | string | Stable delivery identifier |
| taskID | string | Optional task identifier |
| runID | string | Optional run identifier |
| channelName | string | Safe channel display name |
| state | string | `pending`, `claimed`, or `failed` when it represents a retry or failure |
| createdAt | RFC 3339 timestamp | Inside the previous 24 hours |

## SystemObservation

One desktop-owned result for a registration and refresh generation.

| Field | Type | Rules |
|---|---|---|
| generation | integer | Monotonically increasing within the service process |
| registration | SystemRegistration | Snapshot identity used for this refresh |
| state | connection state | Existing typed connection taxonomy plus `not_contacted` and `stale` display classification |
| summary | optional OperationalSummary | Current on success, cached only within the process on later failure |
| stale | boolean | True only when a failed refresh retains prior session data |
| failure | optional safe failure | State-specific title, explanation, and next action |

### State transitions

```text
not_contacted -> refreshing -> current
not_contacted -> refreshing -> failure_state
current -> refreshing -> current
current -> refreshing -> stale + failure_state
stale -> refreshing -> current
stale -> refreshing -> stale + failure_state
any_generation -> superseded -> discarded
removed_registration -> late_result -> discarded
```

## OverviewSnapshot

| Field | Type | Rules |
|---|---|---|
| generation | integer | Identifies the completed refresh |
| startedAt | RFC 3339 timestamp | Desktop clock |
| completedAt | RFC 3339 timestamp | Desktop clock |
| observations | SystemObservation array | One entry per current registration, stable key ordering |

## DrilldownIntent

| Field | Type | Rules |
|---|---|---|
| registrationKey | string | `local` or exact profile ID |
| destination | enum | `tasks`, `schedule`, `activity`, or `notifications` |
| taskID | optional string | Source task context |
| recordID | optional string | Run, alert, or delivery context |

The intent is transient. Connection selection succeeds before the route changes, and a failed selection preserves the overview snapshot.
