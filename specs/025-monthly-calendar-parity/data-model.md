# Data Model: Monthly Calendar Cron Parity

## Schedule calendar adjustment

`domain.Schedule` gains `CalendarAdjustment`, whose values are empty or `nearest_weekday`.

- Empty means RRULE completely defines recurrence execution.
- `nearest_weekday` means RRULE supplies one monthly positive day target and the scheduler adjusts that valid target to Monday through Friday.
- The field is execution-authoritative and independent of `Expression`.
- JSON omits empty; SQLite stores a non-null empty string.

### Invariants

- Frequency is monthly with interval one.
- Exactly one `BYMONTHDAY` from 1 through 31 is present.
- No competing weekday, month, set-position, bound, or richer selector exists.
- Exactly one hour and minute are recoverable; seconds are zero.
- Unknown or incompatible values are errors.

## Native selectors

- Last day uses monthly `BYMONTHDAY=-1` with no adjustment.
- Last weekday uses monthly Monday-through-Friday plus final set position with no adjustment.

## Persistence migration

Schema version 6 adds `calendar_adjustment TEXT NOT NULL DEFAULT ''` to schedules. Existing rows execute unchanged. Create, retrieve, list, replace, backup/recovery, and restart paths preserve the field.

## Execution transition

```text
monthly carrier RRULE + nearest_weekday
-> resolve numbered target under missing-date policy
-> adjust weekend within the resolved month
-> normalize local wall time under DST rules
-> require occurrence strictly after cursor
```

For `next_valid`, existing behavior maps an absent monthly target to the first day of the following month before adjustment. Strictly increasing iteration prevents duplicate output.

## Conversion transition

```text
cron nW -> DOM selector -> phrase -> carrier RRULE + adjustment
native phrase -> carrier RRULE + adjustment -> fidelity checks -> canonical nW
```

`L` and `LW` follow the same flow without metadata. Failure terminates before persistence mutation.
