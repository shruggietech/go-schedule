# Data Model: Scheduler Startup Trigger

## Startup Schedule

| Field | Value | Rule |
| --- | --- | --- |
| `Kind` | `event` | Distinguishes non-clock scheduling |
| `TriggerID` | `scheduler_startup` | Only supported event identity in this slice |
| `HumanSummary` | `At scheduler startup` | Canonical display text |
| `Expression` | `@reboot` or `at scheduler startup` | Retained source, inert for execution |
| `RRULE`, `Anchor`, `ElapsedEpoch`, `RunAt` | empty | Startup has no clock occurrence |

Unknown event trigger identifiers remain readable but are not startup-eligible.

## Task

Startup eligibility requires active state, task enabled, all ancestor groups enabled, and a referenced `scheduler_startup` event schedule. The task remains active after firing. Existing overlap and execution-context fields apply. Catch-up and recurrence anomaly policies do not create event timing.

## Run

The append-only run entity gains origin `startup`. `ScheduledFor` is injected engine time at startup dispatch. Outcome, timestamps, exit code, and output retain existing executor behavior. Legacy `event` values remain readable.

## Lifecycle

```text
persisted eligible task
        |
        v
Engine.Start initial recompute
        |
        v
startup snapshot -> dispatch once -> append startup run
        |
        +-- reload/mutation --------------------> no startup dispatch
        |
        +-- stop, new Engine.Start -------------> eligible again
```

No fired flag, boot identity, trigger row, deduplication ledger, or completion transition is introduced.

