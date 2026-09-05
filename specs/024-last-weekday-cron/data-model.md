# Data Model: Last-Weekday Cron Parity

## Parsed cron field

`cron.Field` retains its current values, wildcard, step, bounds, and ordinal:

- `Ordinal == 0`: ordinary field behavior.
- `Ordinal == -1`: the last occurrence of one weekday in each month.
- `Ordinal == 1..5`: the existing numbered weekday behavior.

### Invariants

- A nonzero ordinal is valid only on day-of-week.
- A nonzero ordinal requires exactly one normalized weekday value.
- Weekday input 7 and `SUN` normalize to value 0.
- Day-of-month and month remain wildcard for a representable phrase.

## Existing last-weekday schedule

No domain entity changes. The existing monthly recurrence contains one numbered weekday with `N() == -1`, plus hour, minute, timezone, and existing missing-date policy.

### Export invariants

- Frequency and interval are exactly monthly and one.
- Exactly one `-1` weekday is present.
- No competing date, month, set-position, bound, multi-time, or sub-minute selector is present.
- Missing-date policy does not affect eligibility because the occurrence exists every month.

## State transitions

```text
five-field cron -> Field{weekday, Ordinal:-1} -> existing human phrase
-> existing monthly RRULE -> task boundary retaining original cron source

existing monthly RRULE -> fidelity checks -> canonical numeric weekdayL
-> parsed field and phrase -> equivalent existing monthly RRULE
```

Malformed or refused input terminates before persistence. There is no schema, migration, or stored-state transition.
