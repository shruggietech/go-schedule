---
title: Cron interoperability
nav_order: 5
---

# Cron interoperability

**Audience:** anyone with an existing crontab, or anyone who has to read one\
**Applies to:** go-schedule 0.7.0 and later for import, explain, and export\
**Source of truth:** `internal/cron/` — this document describes what the
converter does, and the fidelity table below is the contract it holds to.

> **Release status:** `cron convert`, direct cron input for `task add` and
> `task edit`, and cron entry in the desktop editor are currently unreleased
> changes planned for the first release after 0.8.0.

go-schedule accepts schedules written the way you would say them, such as
`every 15 minutes` and `weekdays at 09:00`, and also accepts a supported subset
of five-field cron such as `0 9 * * 1-5`. Every accepted form compiles to the
same stored recurrence model before execution.

But you probably already have a crontab, or need to translate one schedule. Cron
is therefore supported for pure string conversion, task input, and as an
**interchange format**: convert one string, supply it to `task add` or `task
edit`, enter it in the desktop Schedule field, import a file, ask what a line
means with upcoming runs, or export tasks back out.

## Contents

- [Convert one string](#convert-one-string)
- [Import a crontab](#import-a-crontab)
- [Explain one expression](#explain-one-expression)
- [Export back to cron](#export-back-to-cron)
- [Fidelity](#fidelity)
- [What cron cannot say](#what-cron-cannot-say)
- [Expression versus crontab file](#expression-versus-crontab-file)
- [What this scheduler cannot say in cron](#what-this-scheduler-cannot-say-in-cron)

## Convert one string

`convert` is symmetric and entirely local:

```sh
gosched cron convert "0 9 * * 1-5"
# weekdays at 09:00

gosched cron convert "weekdays at 09:00"
# 0 9 * * 1-5

gosched cron convert "0 9 * * 5#3"
# 3rd friday monthly at 09:00

gosched cron convert "3rd friday monthly at 09:00"
# 0 9 * * 5#3
```

Automatic mode treats `@`-prefixed values and five fields with a cron-shaped
minute field as cron. Existing human forms such as
`every 15 minutes from 9am` remain human input. Once classified, invalid cron
is never retried as human text. `--to cron` forces human input; `--to human`
forces cron input.

The normal result is one line and nothing else. Invalid or unfaithful input
leaves stdout empty, names the reason on stderr, and exits 2. Global `--json`
keeps the same stream and exit rule while returning stable syntax, input,
output, and refusal fields. The command makes no daemon call and changes no
task.

Not every human phrase contains enough phase information for cron. For example,
`every 15 minutes` begins when its task is created, while `*/15` always begins
at `:00`. Write `every 15 minutes starting at 00:00` when that is the intended
cron phase; conversion refuses to invent it. The reverse translation retains
that phase explicitly.

## Import a crontab

Always look first:

```sh
gosched cron import --file /etc/crontab --dry-run
```

The preview prints, for every line, the expression, the phrase it maps to, the
resolved command, and, for a real import, the task it created. Imported tasks
retain the original normalized cron expression. Nothing is created while
`--dry-run` is set.

```text
line 3: 0 2 * * *
  phrase:  every day at 02:00
  command: /usr/local/bin/backup --full
line 5: */15 * * * *
  phrase:  every 15 minutes starting at 00:00
  command: /usr/local/bin/probe
line 6: @reboot
  unsupported: @reboot fires at boot rather than on a schedule, which has no
  equivalent here

8 line(s) read: 2 would create, 3 skipped (comments, blanks, variables),
2 unsupported, 1 error(s)
```

The preview is not advisory. Preview and creation use the same conversion
result. A created task retains the normalized cron expression and its compiled
recurrence; the phrase is an explanation of that result, not a replacement
schedule string.

When it does, drop the flag:

```sh
gosched cron import --file /etc/crontab --timezone America/New_York --group ops
```

| Flag | Meaning |
| --- | --- |
| `--file` | Crontab to read, or `-` for standard input. Required. |
| `--dry-run` | Produce the identical report and create nothing. |
| `--timezone` | IANA zone for the created tasks. Cron has none — see below. |
| `--group` | Group ID to file the imported tasks under. |
| `--count` | How many upcoming runs to show per line. Default 3. |

A line that cannot be converted never stops the ones that can: the supported
lines are still created, and the summary counts every line of the file. Reading
the file successfully is a success, whatever the mix of outcomes — a crontab of
nothing but `@reboot` lines converts to a report of refusals and exits 0. Only
an unreadable file, an unknown timezone, or a failed creation is a failure.

Importing the same crontab twice creates two sets of tasks. There is no
deduplication; the counts are how you notice.

## Explain one expression

```sh
gosched cron explain "0 9 * * 1-5"
```

```text
0 9 * * 1-5
  phrase: weekdays at 09:00
  next:   2026-07-24T13:00:00Z
          2026-07-27T13:00:00Z
          2026-07-28T13:00:00Z
```

Nothing is created or changed. An expression that cannot be represented is
reported by name, and that is an answer rather than a failure — the exit code
stays 0. A malformed expression *is* a failure, and exits 2 naming the field.

## Export back to cron

```sh
gosched cron export
gosched cron export --task 6f1c…
```

Every task appears exactly once, as a crontab line or as a commented refusal:

```text
# gosched cron export — 4 task(s)
0 9 * * 1-5 /usr/bin/report --daily
# declined: "nightly backup" — cron cannot express a schedule that fires exactly once
# declined: "health probe" — the task is disabled and cron has no disabled state
*/15 * * * * /usr/local/bin/probe
```

Nothing is silently omitted and nothing is approximated. A converter that
quietly rounds a schedule to the nearest thing cron can say is worse than one
that declines, because the difference only surfaces at 02:30 some morning.

## Fidelity

The supported input shape is `minute hour day-of-month month day-of-week`.
Numeric bounds are `0-59`, `0-23`, `1-31`, `1-12`, and `0-7` respectively,
with both `0` and `7` meaning Sunday. The parser recognizes numbers,
comma-separated lists, inclusive ascending ranges, wildcard steps, and
case-insensitive month or weekday names by their first three letters. Recognition
alone does not guarantee conversion: only the recurrence shapes listed below
are scheduled, and every other well-formed form is refused by name.

This is a product-specific subset, not a promise of POSIX, Linux, or robfig
parity. For the upstream dialects and file behavior, see the
[POSIX crontab utility](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/crontab.html),
[Linux crontab(5)](https://man7.org/linux/man-pages/man5/crontab.5.html), and
[robfig/cron expression format](https://pkg.go.dev/github.com/robfig/cron/v3#hdr-CRON_Expression_Format).
The single `weekday#ordinal` form is an explicit go-schedule extension to this
five-field subset. Some schedulers use `#` with different field layouts or
weekday numbering, so verify the destination before treating exported text as
portable.

### Supported

| Cron | Notes |
| --- | --- |
| Every minute | `* * * * *` |
| Minute intervals dividing one hour evenly | `*/5 * * * *`, `*/15 * * * *`, `*/30 * * * *` |
| Hourly and hour intervals dividing one day evenly | `0 * * * *`, `0 */6 * * *` |
| One fixed daily time | `0 9 * * *` |
| One weekday, weekdays, or weekends | `0 14 * * WED`, `0 9 * * 1-5`, `0 10 * * 0,6` |
| One ordinal weekday each month | `0 9 * * 5#3` means the third Friday. One weekday and ordinal 1 through 5 are supported when month and day-of-month are unrestricted. |
| One monthly day or one yearly month/day | `0 9 31 * *`, `0 0 4 7 *` |
| Month and weekday names | Names are case-insensitive; their first three letters are significant. |
| Sunday as `0` or `7` | Both are accepted |
| `@hourly`, `@daily`, `@midnight`, `@weekly`, `@monthly`, `@yearly`, `@annually` | Expanded to their documented five-field equivalents |

### Declined, by name

| Cron | Why |
| --- | --- |
| `@reboot` | Fires at boot rather than on a schedule. There is no equivalent, and there is no honest approximation of one. |
| Six-field (Quartz) expressions | Seconds-precision cron dialects are a different language. Sub-minute schedules are expressible here directly (`every 30 seconds`) — just not through cron. |
| `L` and `W` | Non-standard day specifiers outside the supported subset. Last-weekday forms remain declined. |
| Broader `#` combinations | Lists, ranges, steps, multiple ordinal terms, ordinals outside 1 through 5, and month/date restrictions are declined rather than approximated. |
| A step that does not divide its range | `*/7` on minutes fires at :00, :07 … :56 and then :00 again — a four-minute gap a fixed interval cannot reproduce. `*/5`, `*/15` and `*/30` are exact and are accepted. |
| A wildcard step in day-of-month, month, or day-of-week | Cron restarts these steps inside each calendar field. The recurrence model cannot retain that field-local behavior, so forms such as `0 9 */2 * *` are refused rather than simplified to daily. |
| Both day-of-month and day-of-week restricted | `0 0 13 * 5` means "the 13th **or** any Friday" in cron. This scheduler intersects the two, which would turn a weekly job into a handful of runs a year. |
| Lists in the minute, hour, day, or month field | `0 9,17 * * *` is two schedules wearing one expression. Create two tasks. |

## What cron cannot say

Everything below is a property your tasks gain on import, and the import summary
says so rather than leaving you to discover it:

- **A timezone.** A five-field expression contains no timezone. Every task here
  carries its own IANA zone; `--timezone` sets it during import and otherwise
  the default is used. Crontab environment assignments such as `CRON_TZ` are
  file context and are not carried into task definitions.
- **Catch-up.** If the machine is off when a cron job is due, that run is simply
  lost. Imported tasks get `catchup one`: a single catch-up run after downtime,
  then the normal schedule resumes.
- **An overlap policy.** Cron will happily start a second copy of a job that is
  still running. Imported tasks get `overlap queue_one`.
- **A missing-date policy.** `0 9 31 * *` runs seven months in twelve and cron
  never mentions it. Imported tasks get `skip`, which is exactly cron's
  behavior — see [the CLI reference](cli.md#task) to change it.
- **Restart recovery.** The daemon reconstructs its schedule from durable state
  on restart.

The task timezone also owns Daylight Saving resolution. A nonexistent wall time
moves to the next valid instant; a repeated wall time uses its first occurrence
once. Those are go-schedule rules and should not be read as Linux cron parity.

## Expression versus crontab file

A schedule expression is only five timing fields (or a supported macro). A
user crontab line adds a command. A system crontab can also add a username, and
a file can include comments, variables, shell settings, or `CRON_TZ`. Use
`task add --schedule` for one expression and command you provide separately.
Use `cron import` when the command and file context must first be inspected.

## What this scheduler cannot say in cron

The export declines these rather than approximating them:

- One-off schedules — cron has no way to fire exactly once.
- Sub-minute intervals — cron's resolution is one minute.
- Intervals that do not divide their period evenly (`every 3 days`,
  `every 2 weeks`) — cron repeats by calendar position, not elapsed time.
- Sub-daily intervals whose stored phase does not align with cron's field-local
  step — exporting `:05/:20/:35/:50` as `*/15` would silently move every run.
- Last-weekday and broader ordinal-weekday combinations. The focused
  `weekday#1..5` subset exports, but the `#` extension is not universal across
  cron implementations.
- A fifth-weekday or other date-bearing task using a non-default missing-date
  policy. Cron would silently skip a missing occurrence; first through fourth
  weekdays exist every month, so their policy setting does not change output.
- Disabled tasks — cron has no disabled state, and emitting a live line for a
  task you deliberately stopped would be the worst possible outcome.
