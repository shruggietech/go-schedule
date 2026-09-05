# Research: Dual-Syntax Product Documentation

## Product terminology

Use **go-schedule's supported five-field cron subset**. Do not claim POSIX, standard-cron, Linux-crontab, or robfig compatibility. An expression contains five timing fields or one supported macro, not a command, username, environment assignment, or whole crontab file.

## Faithful recurrence shapes

Document only shapes the product preserves: every minute; evenly dividing wildcard-minute intervals; hourly at minute zero; evenly dividing wildcard-hour intervals; a fixed daily time; one weekday, weekdays, or weekends; one monthly day; and one month/day yearly date. Supported macros are `@hourly`, `@daily`, `@midnight`, `@weekly`, `@monthly`, `@yearly`, and `@annually`. `@reboot` is refused.

The lexical parser also understands numbers, lists, inclusive ascending ranges, steps, and three-letter case-insensitive month/weekday names. Recognition does not promise faithful scheduling. Joint day-of-month and day-of-week restrictions are refused because cron OR semantics cannot be represented by the product's recurrence intersection.

## Timezone and DST language

A cron expression has no timezone. The task timezone controls recurrence and anchor interpretation. Imports accept an explicit timezone; `CRON_TZ` and file environment assignments are not carried into task definitions. A nonexistent local time advances to the next valid instant; a repeated local time uses its first occurrence once. This is product behavior, not Linux cron DST parity.

## Historical policy

Update S001 as the authoritative product contract. Preserve S008's original human-only scope and add visible supersession notices to its spec, tasks, quickstart, and fidelity checklist.

## Policy automation

Use an explicit current-surface inventory and targeted phrases. A temporary fixture harness proves that aligned copy passes and targeted obsolete claims fail. Historical specifications and the changelog remain outside this scan.

## Reproduced defect

`0 9 */2 * *`, `0 9 * */2 *`, and `0 9 * * */2` currently normalize to "every day at 09:00". Refuse wildcard steps greater than one in those three calendar fields. Other dialect breadth remains out of scope and tracked by #22.

## Primary references

- https://pubs.opengroup.org/onlinepubs/9699919799/utilities/crontab.html
- https://man7.org/linux/man-pages/man5/crontab.5.html
- https://pkg.go.dev/github.com/robfig/cron/v3#hdr-CRON_Expression_Format
