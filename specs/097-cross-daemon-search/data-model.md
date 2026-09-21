# Data Model: Cross-Daemon Search and Target-Safe Actions

## SearchRequest

| Field | Type | Rules |
|---|---|---|
| query | string | Trimmed, non-empty, case-insensitive match input; maximum 200 characters |
| kinds | result-kind set | Optional subset of task, group, failure, schedule, alert |
| limit | integer | One through 50; defaults to 50 per daemon |
| generation | integer | Desktop-owned monotonically increasing search generation |

## DaemonSearchResponse

| Field | Type | Rules |
|---|---|---|
| schema | string | `go-schedule.daemon-search.v1` |
| query | string | Normalized query used by the daemon |
| observedAt | RFC 3339 timestamp | Captured once from the daemon clock |
| truncated | boolean | True when more matches exist than returned |
| results | SearchMatch array | At most the requested limit, deterministically ordered |

## SearchMatch

| Field | Type | Rules |
|---|---|---|
| kind | enum | task, group, failure, schedule, alert |
| objectID | string | Task, group, run, occurrence, or alert identity |
| taskID | optional string | Owning task identity when applicable |
| name | string | Safe object or task label |
| context | string | Safe bounded group, state, readiness, schedule, severity, or kind summary |
| occurredAt | optional RFC 3339 timestamp | Failure, schedule, or alert time |
| enabled | optional boolean | Present for task matches |
| actionHints | string set | Object-level possible actions before connection authority is applied |

SearchMatch never includes commands, arguments, working directories, environment, stdin, run output, alert messages, credentials, notification payloads, or secret-bearing endpoints.

## TargetSearchObservation

| Field | Type | Rules |
|---|---|---|
| generation | integer | Search generation that owns this observation |
| registration | SystemRegistration | Immutable local or saved-profile routing identity |
| daemonID | string | Current health-confirmed daemon identity when connected |
| state | connection state | Existing typed connection taxonomy |
| permissions | string array | Safe read, operate, manage projection |
| capabilities | string array | Current manifest capability names |
| observedAt | optional RFC 3339 timestamp | Daemon search observation time |
| response | optional DaemonSearchResponse | Present on success |
| failure | optional safe failure | State, explanation, and next action |

### Search generation transitions

```text
idle -> searching -> partial -> complete
searching -> superseded -> discarded
partial -> superseded -> discarded
registration_removed -> late_result -> discarded
target_searching -> connected_result
target_searching -> target_failure
```

## SearchResultView

| Field | Type | Rules |
|---|---|---|
| key | string | Registration key plus result kind and object identity |
| registration | SystemRegistration | Source identity shown everywhere the result appears |
| daemonID | string | Health-confirmed identity from the observation |
| match | SearchMatch | Safe daemon-produced object summary |
| freshness | string | Observation time and current or stale classification |
| actions | ActionAvailability map | Computed from kind, hints, connection state, capabilities, and permissions |

## ActionAvailability

| Field | Type | Rules |
|---|---|---|
| state | enum | available, unavailable, refresh_required |
| reason | string | Concise explanation without secret or raw transport detail |

## ActionIntent

| Field | Type | Rules |
|---|---|---|
| action | enum | acknowledge, enable, disable, run_now |
| selections | ActionSelection array | Non-empty deliberate selection; every item must support the same action |
| requestID | string | Desktop-generated retry identity shared only for this submission |

## ActionSelection

| Field | Type | Rules |
|---|---|---|
| registrationKey | string | `local` or exact saved profile ID |
| expectedDaemonID | string | Daemon identity observed by search |
| kind | enum | task or alert |
| objectID | string | Exact task or alert identifier |
| taskID | optional string | Task identity for task actions and alert context |
| displayName | string | Confirmation label only; never used for routing |
| observedAt | RFC 3339 timestamp | Search observation evidence |

## ActionOutcome

| Field | Type | Rules |
|---|---|---|
| selection | ActionSelection | Exact requested target and object |
| outcome | enum | accepted, rejected, unavailable, uncertain |
| message | string | Target-specific explanation and recovery action |
| currentDaemonID | optional string | Present after successful identity validation |

### Mutation transitions

```text
selected -> confirmed -> reconnecting -> identity_validated -> object_validated -> submitted -> accepted
reconnecting -> target_failure -> unavailable
identity_validated -> identity_mismatch -> rejected
object_validated -> object_changed_or_missing -> rejected
submitted -> transport_uncertain -> uncertain
submitted -> target_rejection -> rejected
```
