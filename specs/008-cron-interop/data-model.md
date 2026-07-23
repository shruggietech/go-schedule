# Phase 1 Data Model: Cron interoperability and calendar-anomaly policy

**Feature**: `008-cron-interop` | **Date**: 2026-07-23

## Changed entities

### `domain.Task` — gains `MissingDatePolicy`

```go
// MissingDatePolicy controls what a schedule does in a period that has no
// matching date — February with no 29th, a 30-day month with no 31st, a month
// with only four Fridays.
type MissingDatePolicy string

const (
    MissingDateSkip      MissingDatePolicy = "skip"       // default: no run that period
    MissingDateLastValid MissingDatePolicy = "last_valid" // the last date that does exist
    MissingDateNextValid MissingDatePolicy = "next_valid" // roll into the next period
)
```

Added to `Task` after `CatchupPolicy`, matching the existing policy fields in
placement, JSON tag style (`missing_date_policy`), and validation treatment.

**Why the task and not the schedule**: an edit that changes the phrase creates a
new schedule row (`internal/api/server/update.go`), so a schedule-borne policy
would silently reset on an unrelated edit. See research decision D1.

### `schedules` table — unchanged

No column is added. The v4 `expression` column is the last schedule-level
addition.

### `tasks` table — migration v5

```sql
ALTER TABLE tasks ADD COLUMN missing_date_policy TEXT NOT NULL DEFAULT 'skip';
```

Additive with a total default; forward-only; no existing column, row, or value
is read or rewritten. Every pre-v5 task takes `skip`, which is exactly the
behavior it had before (FR-020, FR-026).

## New entities (in-memory only, never persisted)

### `cron.Spec` — a parsed cron expression

The structural result of parsing one crontab timing field-set. Holds the five
field value-sets (minute, hour, day-of-month, month, day-of-week) plus the
shorthand it came from, if any. It exists so the converter can inspect field
structure — which a scheduling library's compiled schedule does not expose — and
decide representability.

### `cron.Unsupported` — a named refusal

Carries the input, a machine-readable reason code, and a human sentence. Returned
wherever the converter declines: an unsupported extension, a non-dividing step, a
day-of-month/day-of-week combination, or a schedule shape the export cannot
carry. Never an error value on its own — a refusal is an outcome (research
decision D4).

### `cron.Line` — one crontab line's conversion result

Pairs the source line with either a phrase plus payload (command, arguments,
working directory) or an `Unsupported`, plus any per-line warnings (`MAILTO`,
variable assignments). The import report is a slice of these; preview and real
import produce identical slices.

### `cron.Report` — the run summary

Counts of lines read, created, skipped as comment or blank, declined, and
errored, plus the fidelity statements the run must make (timezone applied and
that cron carries none; the catch-up, overlap, and restart-recovery defaults
imported tasks received).

## Changed function signatures

```go
// internal/schedule
func NextRun(sch domain.Schedule, tzName string, policy domain.MissingDatePolicy, after time.Time) (time.Time, bool, error)
func UpcomingRuns(sch domain.Schedule, tzName string, policy domain.MissingDatePolicy, after time.Time, n int) ([]time.Time, error)
```

Six call sites, all of which already hold the task:
`internal/engine/engine.go:146,235`, `internal/catchup/catchup.go:31` (via a new
parameter on `catchup.Evaluate`), `internal/api/server/calendar.go:74`,
`internal/api/server/tasks.go:209,216`.

## API contract changes

`TaskCreateRequest` and `TaskUpdateRequest` each gain
`missing_date_policy string`, validated exactly like `overlap_policy`: empty
means "unchanged" on update and "default" on create; an unrecognized value is a
validation error naming the field. `TaskResponse` carries the policy on the
embedded `Task`, so every existing reader gets it for free.

## Invariants

1. A task's stored timing values never change as a result of this feature.
2. The policy is inert for any schedule whose recurrence has no date component
   (interval, weekday, dayset rules) — it is consulted only by the by-date and
   ordinal resolution path.
3. The policy and the schedule phrase are independently mutable; neither write
   path touches the other.
4. Daylight-saving resolution runs on whatever instant the policy produces, with
   the same rules as today.
