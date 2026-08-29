# Research: Explicit DST Scheduling Intent

## Decision: preserve wall-clock as the compatibility default

**Rationale**: Existing tasks and the master specification already promise day-or-coarser local wall-clock behavior with next-valid spring and first-occurrence fall resolution. The cited five/seven-hour gap is expected for wall-clock intent, so changing it would silently reinterpret stored tasks.

**Alternatives considered**: Default sub-daily schedules to elapsed time (rejected because frequency would continue to infer unstated intent); add `auto` as a fourth basis (rejected because it recreates the ambiguity this slice removes).

## Decision: support elapsed only for exact-duration recurrence shapes

**Rationale**: Seconds, minutes, hours, single-occurrence days, and single-occurrence weeks have exact durations. Months, years, date sets, ordinal weekdays, and multi-time periods do not. Refusal is clearer than assigning an arbitrary duration to a calendar unit.

**Alternatives considered**: Evaluate all elapsed schedules as UTC RRULEs (rejected because a UTC month is still not an elapsed duration); define months as 30 days (rejected as silent approximation).

## Decision: evaluate UTC as a distinct calendar basis

**Rationale**: UTC means recurrence fields remain fixed to UTC clock/calendar readings. This differs from elapsed for variable calendar periods and from wall-clock for task-local readings.

**Alternatives considered**: Alias UTC to elapsed (rejected because monthly and date-selected UTC schedules are valid but not fixed-duration intervals).

## Decision: resolve local wall intent by enumerating concrete mappings

**Rationale**: Go's default local-time constructor does not expose ambiguity and may choose either offset. Enumerating candidate instants around the requested calendar reading yields zero, one, or two exact mappings and works for non-one-hour IANA transitions.

**Alternatives considered**: Add or subtract one hour (rejected because transition deltas vary); keep first-occurrence hardcoded (rejected because it cannot implement `both` or `last`).

## Decision: carry all scheduling choices in one value

**Rationale**: `NextRun` already receives timezone plus missing-date policy. Three more positional strings would be error-prone across engine, catch-up, preview, and calendar. One typed value makes defaults and composition explicit.

**Alternatives considered**: Put policies on Schedule (rejected because schedule replacement would reset task intent); global daemon configuration (rejected because tasks on one machine have different purposes).

## Decision: retain transition choices while inert

**Rationale**: Gap and overlap policies have no effect under elapsed or UTC, but clearing them would destroy operator intent when switching bases. Persistence is cheap and predictable.

**Alternatives considered**: Reject non-default transition values outside wall-clock (rejected because it couples unrelated edits and complicates basis changes).
