# Data Model: Ordinal-Weekday Cron Parity

## Parsed cron field

`cron.Field` retains its current values, wildcard, step, and bounds. It gains:

- `Ordinal int`: zero for ordinary field behavior; 1 through 5 for one day-of-week occurrence within a month.

### Invariants

- A nonzero ordinal is valid only on the day-of-week field.
- A nonzero ordinal requires exactly one normalized weekday value.
- Weekday input 7 and `SUN` normalize to value 0.
- Ordinal input is restricted to 1 through 5.
- Day-of-month and month remain wildcard for a representable phrase.

## Existing ordinal-weekday schedule

No domain entity changes. The existing monthly recurrence contains one numbered weekday (`rrule.Weekday` with positive `N`) plus hour, minute, timezone, and the existing missing-date policy.

### Export invariants

- Frequency is monthly.
- Exactly one numbered weekday is present.
- Its occurrence is 1 through 5.
- No competing date, month, set-position, yearly, or Easter selector is present.
- Fifth occurrence requires effective skip behavior.

## State transitions

```text
five-field cron -> parsed Field{weekday, ordinal} -> human phrase
-> existing monthly RRULE -> task boundary retaining original cron source

existing monthly RRULE -> fidelity checks -> canonical numeric cron
-> parsed field and phrase -> equivalent existing monthly RRULE
```

Malformed or refused input terminates before task persistence. There is no schema, migration, or stored-state transition.
