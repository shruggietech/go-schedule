# Quickstart: Last-Weekday Cron Parity

## Explain and convert cron

```text
goschedule cron explain "0 9 * * 5L"
last friday of the month at 09:00

goschedule cron convert "30 14 * * WEDL"
last wednesday of the month at 14:30
```

## Convert a native schedule to canonical cron

```text
goschedule cron convert "last wednesday of the month at 14:00"
0 14 * * 3L
```

Sunday input may use `0L`, `7L`, or `SUNL`; export always uses `0L`.

## Import a crontab job

A line such as `0 9 * * 5L report-command` previews as a job with the last
Friday phrase. Creation retains `0 9 * * 5L` as the schedule expression and
marks its syntax as cron.

## Intentional boundaries

- Exactly one weekday followed by `L` is supported.
- Day-of-month and month must remain unrestricted.
- Bare `L`, lists, ranges, steps, multiple terms, mixed `L`/`#`, day-of-month
  `L`, `W`, six fields, and `@reboot` remain refused.
- Missing-date policy does not constrain export because every month contains a
  last occurrence of each weekday.

This is a documented go-schedule five-field extension. It does not claim that
every POSIX, Vixie-derived, Quartz, or hosted cron implementation accepts the
same expression.
