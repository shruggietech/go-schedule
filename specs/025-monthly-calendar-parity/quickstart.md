# Quickstart: Monthly Calendar Cron Parity

## Explain existing cron

```text
goschedule cron explain "0 9 L * *"
last day of every month at 09:00

goschedule cron explain "0 9 15W * *"
nearest weekday to the 15th of every month at 09:00

goschedule cron explain "0 9 LW * *"
last weekday of every month at 09:00
```

## Convert native schedules

```text
goschedule cron convert "last day of every month at 09:00"
0 9 L * *

goschedule cron convert "nearest weekday to the 15th of every month at 09:00"
0 9 15W * *

goschedule cron convert "last weekday of every month at 09:00"
0 9 LW * *
```

## Calendar examples

- `1W` on Saturday uses Monday the 3rd.
- `15W` on Saturday uses Friday the 14th.
- `31W` on Sunday uses Friday the 29th.
- `LW` in a month ending Sunday uses Friday.

Targets 29 through 31 can be absent. Native tasks apply selected missing-date policy first; imported cron uses skip.

## Intentional boundaries

Exactly one selector is accepted with wildcard month and day-of-week. Offsets, lists, ranges, steps, mixtures, Quartz `?`, and six fields remain refusals. DOM `LW` means last weekday; DOW `5L` remains the distinct last-Friday expression.
