# Quickstart: Ordinal-Weekday Cron Parity

## Explain and convert cron

```text
goschedule cron explain "0 9 * * 5#3"
3rd friday monthly at 09:00

goschedule cron convert "30 14 * * WED#2"
2nd wednesday monthly at 14:30
```

## Convert a native schedule to canonical cron

```text
goschedule cron convert "3rd wednesday monthly at 14:00"
0 14 * * 3#3
```

Sunday input may use `0`, `7`, or `SUN`; export always uses `0`.

## Import a crontab job

A line such as `0 9 * * 5#3 report-command` previews as a job with the third Friday phrase. Creation retains `0 9 * * 5#3` as the schedule expression and marks its syntax as cron.

## Intentional boundaries

- Ordinal must be 1 through 5.
- Exactly one weekday and one ordinal are supported.
- Day-of-month and month must remain unrestricted.
- Lists, ranges, steps, multiple terms, `L`, `W`, six fields, and `@reboot` remain refused.
- Fifth-weekday export requires skip behavior because some months have no fifth occurrence.

This is a documented go-schedule five-field extension. It does not claim that every POSIX, Vixie-derived, Quartz, or hosted cron implementation accepts the same expression.
