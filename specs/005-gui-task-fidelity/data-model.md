# Phase 1 Data Model: GUI task fidelity

**Feature**: [spec.md](spec.md) · **Research**: [research.md](research.md) ·
**Date**: 2026-07-22

Only one persisted entity changes. Everything else on this page is stated to
pin what must *not* change.

## Schedule *(modified)*

`internal/domain.Schedule` gains one field.

| Field | Type | Change | Meaning |
|---|---|---|---|
| `ID` | string | — | identity |
| `Kind` | ScheduleKind | — | `one_off` \| `recurring` \| `event` |
| `RRULE` | string | — | **authoritative recurrence**; the only timing input the engine reads |
| `Anchor` | *time.Time | — | first-cycle instant (UTC); an explicit anchor and the creation-time default are indistinguishable once stored |
| `RunAt` | *time.Time | — | one-off instant (UTC) |
| `TriggerID` | string | — | unused since migration v3 |
| `HumanSummary` | string | — | system-generated description, e.g. `Every weekday at 09:00`; **not re-parseable** |
| `Expression` | string | **new** | the operator's phrase, e.g. `weekdays at 09:00`; re-parseable; inert with respect to execution (FR-011a) |

### The three timing strings, distinguished

Confusing these is the main risk this feature introduces, so:

- **`RRULE`** — what the engine evaluates. Authoritative. Unchanged by this
  feature.
- **`HumanSummary`** — what the system says back to the operator. Display only.
  Unchanged by this feature.
- **`Expression`** — what the operator typed. Round-trip only, so a client can
  put their own words back in the input they typed them into. **Nothing on the
  execution path may read it** (FR-011a).

### Rules

- Set by `schedule.Parse` from its trimmed input, for recurring schedules only.
- `schedule.NewOneOff` leaves it empty: a one-off's date and time are recovered
  from `RunAt`, so there is nothing to round-trip.
- May legitimately be empty: one-off schedules, and rows written before
  migration v4. Nothing reconstructs a phrase for the latter — an empty field on
  edit means "keep the existing schedule", which is harmless (R1).
- Two schedules with identical `RRULE` and different `Expression` are equivalent
  for execution. `Expression` is never used for comparison, dispatch, or
  ordering.

### Storage

`schedules` table, migration **v4** — one statement, appended to the existing
ordered `migrations` slice:

```sql
ALTER TABLE schedules ADD COLUMN expression TEXT NOT NULL DEFAULT '';
```

Forward-only and non-destructive: adds one column with a total default, reads
and rewrites nothing (FR-002). Transactional rollback and re-run idempotency come
from the existing runner (R4). `CreateSchedule` and `GetSchedule` in
`internal/store/crud.go` add `expression` to their column lists.

The migration is kept rather than folding `expression` into the v1
`CREATE TABLE`: an existing database records `schema_version = 3`, so v1 would
never re-run and the column would silently not exist, failing every schedule
query. A wipe is a fine outcome for the maintainers' non-functional databases;
a silent hard failure is not.

## Task *(unchanged shape, one reachability change)*

No field is added or removed. `GroupID` keeps its type and meaning; what changes
is that "no group" becomes an *expressible update intent* rather than an
unreachable state — see [contracts/task-update.md](contracts/task-update.md).
An unresolvable `GroupID` is presented as no group and never raises an error
(FR-019a).

## Group *(unchanged)*

No change. Deletion semantics stand: children cascade, member tasks are released
to ungrouped by the existing `ON DELETE SET NULL` foreign key. Those released
tasks must now be visible in the ungrouped area (FR-019).

## Derived, not stored

- The **group hierarchy path** shown in group choice lists (`Parent / Child`) is
  computed from the group list at render time. Not persisted.
- The **ungrouped area** in the hierarchy view is a synthetic node with no
  backing record. Always present, empty or not (FR-019).
