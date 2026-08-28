# Contract: Monthly Calendar Cron Selectors

## Accepted five-field input

Day-of-month accepts exactly `L`, `<day>W`, or `LW`, where day is 1 through 31. Month and day-of-week must be `*`.

## Canonical human output

```text
last day of every month at HH:MM
nearest weekday to the <ordinal> of every month at HH:MM
last weekday of every month at HH:MM
```

The same phrase families are accepted using existing optional-time and timezone rules.

## Canonical cron output

```text
<minute> <hour> L * *
<minute> <hour> <day>W * *
<minute> <hour> LW * *
```

`L`, `LW`, and targets 1 through 28 export under every policy. Targets 29 through 31 require skip.

## Calendar behavior

- Weekdays remain unchanged.
- Saturday moves to Friday, except day 1 moves to Monday day 3.
- Sunday moves to Monday, except a final-day Sunday moves to Friday.
- `LW` uses the last date when it is a weekday, otherwise the last Friday.
- An absent target is skipped or resolved under task policy before adjustment.
- Timezone and DST use the existing local wall-time contract.

## Source and interface behavior

Explain, conversion, preview, creation, editing, import, and export share this contract. Cron input retains its normalized expression and `cron` syntax. Failed or refused creation/update mutates nothing.

## Refusal boundary

Bare/invalid `W`, `0W`, `32W`, malformed `LW`, offsets, lists, ranges, steps, mixed/multiple modifiers, restricted month/day-of-week, Quartz `?`, six fields, and richer native selectors are not approximated.
