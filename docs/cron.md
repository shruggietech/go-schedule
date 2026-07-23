# Cron interoperability

**Audience:** anyone with an existing crontab, or anyone who has to read one\
**Applies to:** go-schedule 0.7.0 and later\
**Source of truth:** `internal/cron/` — this document describes what the
converter does, and the fidelity table below is the contract it holds to.

go-schedule does not use cron syntax. Schedules are written the way you would
say them — `every 15 minutes`, `weekdays at 09:00` — and that is deliberate:
`0 9 * * 1-5` is a fine thing for a computer to read and a poor thing for a
person to maintain.

But you probably already have a crontab. So cron is supported as an
**interchange format**: you can import one, ask what a line means, and export
your tasks back out. What you cannot do is *write* cron here. There is no field
in the GUI that takes an expression, and `--schedule "0 9 * * 1-5"` is an error.
Conversion happens at the boundary; it does not become the interface.

## Contents

- [Import a crontab](#import-a-crontab)
- [Explain one expression](#explain-one-expression)
- [Export back to cron](#export-back-to-cron)
- [Fidelity](#fidelity)
- [What cron cannot say](#what-cron-cannot-say)
- [What this scheduler cannot say in cron](#what-this-scheduler-cannot-say-in-cron)

## Import a crontab

Always look first:

```sh
gosched cron import --file /etc/crontab --dry-run
```

The preview prints, for every line, the expression, the phrase it maps to, the
resolved command, and — for a real import — the task it created. Nothing is
created while `--dry-run` is set.

```text
line 3: 0 2 * * *
  phrase:  every day at 02:00
  command: /usr/local/bin/backup --full
line 5: */15 * * * *
  phrase:  every 15 minutes
  command: /usr/local/bin/probe
line 6: @reboot
  unsupported: @reboot fires at boot rather than on a schedule, which has no
  equivalent here

8 line(s) read: 2 would create, 3 skipped (comments, blanks, variables),
2 unsupported, 1 error(s)
```

The preview is not advisory. The phrase it shows you is the string that is
parsed and stored when you run it for real — there is no second conversion path
— so if the preview reads correctly, the import is correct.

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
0 9 * * 1,2,3,4,5 /usr/bin/report --daily
# declined: "nightly backup" — cron cannot express a schedule that fires exactly once
# declined: "health probe" — the task is disabled and cron has no disabled state
*/15 * * * * /usr/local/bin/probe
```

Nothing is silently omitted and nothing is approximated. A converter that
quietly rounds a schedule to the nearest thing cron can say is worse than one
that declines, because the difference only surfaces at 02:30 some morning.

## Fidelity

### Supported

| Cron | Notes |
| --- | --- |
| `*` in any field | |
| Single values, lists, and ranges | `0 9 * * 1-5`, `0 10 * * 0,6` |
| Step values | Only where the step divides its field's range evenly — see below |
| Month and weekday names | `JAN`, `MON`, and their long forms |
| Sunday as `0` or `7` | Both are accepted |
| `@hourly`, `@daily`, `@midnight`, `@weekly`, `@monthly`, `@yearly`, `@annually` | Expanded to their documented five-field equivalents |

### Declined, by name

| Cron | Why |
| --- | --- |
| `@reboot` | Fires at boot rather than on a schedule. There is no equivalent, and there is no honest approximation of one. |
| Six-field (Quartz) expressions | Seconds-precision cron dialects are a different language. Sub-minute schedules are expressible here directly (`every 30 seconds`) — just not through cron. |
| `L`, `W`, `#` | Non-standard day specifiers. `#` in particular (`5#3`, "the third Friday") is expressible here as `3rd friday monthly`, so write it that way rather than importing it. |
| A step that does not divide its range | `*/7` on minutes fires at :00, :07 … :56 and then :00 again — a four-minute gap a fixed interval cannot reproduce. `*/5`, `*/15` and `*/30` are exact and are accepted. |
| Both day-of-month and day-of-week restricted | `0 0 13 * 5` means "the 13th **or** any Friday" in cron. This scheduler intersects the two, which would turn a weekly job into a handful of runs a year. |
| Lists in the minute, hour, day, or month field | `0 9,17 * * *` is two schedules wearing one expression. Create two tasks. |

## What cron cannot say

Everything below is a property your tasks gain on import, and the import summary
says so rather than leaving you to discover it:

- **A timezone.** Cron inherits the daemon's. Every task here carries its own
  IANA zone and resolves Daylight Saving transitions explicitly. `--timezone`
  sets it; without it, imported tasks take the default.
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

## What this scheduler cannot say in cron

The export declines these rather than approximating them:

- One-off schedules — cron has no way to fire exactly once.
- Sub-minute intervals — cron's resolution is one minute.
- Intervals that do not divide their period evenly (`every 3 days`,
  `every 2 weeks`) — cron repeats by calendar position, not elapsed time.
- Ordinal-weekday rules (`3rd wednesday monthly`) — expressible only through the
  non-standard `#` extension, which not every cron implementation has.
- Any task using a non-default missing-date policy — cron would silently skip
  the periods the task is specifically configured to run in.
- Disabled tasks — cron has no disabled state, and emitting a live line for a
  task you deliberately stopped would be the worst possible outcome.
