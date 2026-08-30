# Contract: Scheduler Startup Schedule

## Canonical syntax

| Surface | Input | Canonical counterpart |
| --- | --- | --- |
| Cron | `@reboot` | `at scheduler startup` |
| Human | `at scheduler startup` | `@reboot` |

Human parsing is case-insensitive and both inputs trim surrounding whitespace.

## Schedule response

```json
{
  "kind": "event",
  "trigger_id": "scheduler_startup",
  "human_summary": "At scheduler startup",
  "expression": "@reboot",
  "source_syntax": "cron"
}
```

`source_syntax` is derived for event schedules with retained expressions. `next_runs` is empty. RRULE, anchor, elapsed epoch, and one-off time are absent.

## Explain and conversion

- `cron explain @reboot` succeeds, returns the phrase, and prints no `next` lines.
- `cron convert @reboot` outputs `at scheduler startup`.
- `cron convert "at scheduler startup"` outputs `@reboot`.

## Engine lifecycle

For each `Engine.Start(ctx)` invocation, load active eligible tasks, dispatch each startup schedule once at injected current time with origin `startup`, continue catch-up and normal scheduling, and treat reload as recomputation only. Cancellation drains startup work through the existing worker lifecycle.

## Crontab interoperability

`@reboot <command>` is supported in user and system layouts. Dry-run, duplicate detection, command context, group selection, and mutation accounting match timed jobs. Export produces `@reboot` only when existing task-context fidelity rules permit it.

