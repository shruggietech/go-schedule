# Contract: Last-Weekday Cron Parity

## Accepted input

The parser accepts one five-field expression whose day-of-week field is:

```text
<weekday>L
```

- `<weekday>` is numeric 0 through 7 or an existing case-insensitive weekday name.
- Day-of-month and month are unrestricted wildcard fields.

## Canonical output

Native monthly last-weekday schedules export as:

```text
<minute> <hour> * * <numeric-weekday>L
```

Weekdays use Sunday 0 through Saturday 6. Numeric 7 and names are accepted on input but never emitted.

## Human representation

Explain and cron-to-human conversion return:

```text
last <weekday-name> of the month at HH:MM
```

## Task boundary

Successful cron input produces the existing monthly recurrence while retaining the original expression as source and `cron` as syntax. Malformed or refused input creates or updates no task.

## Refusal boundary

The implementation does not approximate bare or malformed `L`, lists, ranges, steps, multiple terms, mixed `L`/`#`, day-of-month `L`, restricted dates or months, arbitrary negative ordinals, selector-rich native recurrences, `W`, Quartz seconds or `?`, and macro-only forms.

## Compatibility

CLI text/JSON streams, crontab job classification, existing cron expressions, and existing named refusals remain unchanged outside this accepted subset.
