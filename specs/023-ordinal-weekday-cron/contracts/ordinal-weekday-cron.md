# Contract: Ordinal-Weekday Cron Parity

## Accepted input

The parser accepts one five-field expression whose day-of-week field is:

```text
<weekday>#<ordinal>
```

- `<weekday>` is numeric 0 through 7 or an existing case-insensitive weekday name.
- `<ordinal>` is an integer 1 through 5.
- Day-of-month and month are unrestricted wildcard fields.

## Canonical output

Native monthly ordinal-weekday schedules export as:

```text
<minute> <hour> * * <numeric-weekday>#<ordinal>
```

Weekdays use Sunday 0 through Saturday 6. Numeric 7 and names are accepted on
input but never emitted.

## Human representation

Explain and cron-to-human conversion return the existing phrase:

```text
<ordinal-word> <weekday-name> monthly at HH:MM
```

## Task boundary

Successful cron input produces the existing monthly recurrence while retaining:

- the original cron expression as schedule source;
- `cron` as the explicit schedule syntax;
- the same generated run times as the corresponding human schedule.

Malformed or refused input creates or updates no task.

## Refusal boundary

The implementation does not approximate:

- ordinals outside 1 through 5;
- lists, ranges, steps, or multiple ordinal terms;
- `#` outside day-of-week;
- restricted day-of-month or month combinations;
- zero, negative, or last-weekday native ordinals;
- fifth-weekday rules whose effective missing-date policy is not skip;
- `L`, `W`, Quartz seconds or `?`, and macro forms.

## Compatibility

CLI text/JSON streams, crontab job classification, existing cron expressions,
and existing named refusals remain unchanged outside this accepted subset.
