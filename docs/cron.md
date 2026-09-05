---
title: Cron interoperability
nav_order: 5
---

# Cron interoperability

**Audience:** anyone with an existing crontab, or anyone who has to read one\
**Applies to:** go-schedule 0.7.0 and later for import, explain, and export\
**Source of truth:** `internal/cron/`, this document describes what the converter does, and the fidelity table below is the contract it holds to.

> **Release status:** `cron convert`, direct cron input for `task add` and `task edit`, and cron entry in the desktop editor are currently unreleased changes planned for the first release after 0.8.0.

go-schedule accepts schedules written the way you would say them, such as `every 15 minutes` and `weekdays at 09:00`, and also accepts a supported subset of conventional five-field cron such as `0 9 * * 1-5`, plus a bounded Quartz-style six-field form such as `*/30 * * * * *`. Every accepted form compiles to the same stored recurrence model before execution.

But you probably already have a crontab, or need to translate one schedule. Cron is therefore supported for pure string conversion, task input, and as an **interchange format**: convert one string, supply it to `task add` or `task edit`, enter it in the desktop Schedule field, import a file, ask what a line means with upcoming runs, or export tasks back out.

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

gosched cron convert "0 9 * * 5L"
# last friday of the month at 09:00

gosched cron convert "last friday of the month at 09:00"
# 0 9 * * 5L
```

Automatic mode treats `@`-prefixed values and five or six cron-shaped fields as cron. Existing human forms such as `every 15 minutes from 9am` remain human input. Once classified, invalid cron is never retried as human text. `--to cron` forces human input; `--to human` forces cron input.

The normal result is one line and nothing else. Invalid or unfaithful input leaves stdout empty, names the reason on stderr, and exits 2. Global `--json` keeps the same stream and exit rule while returning stable syntax, input, output, and refusal fields. The command makes no daemon call and changes no task.

Not every human phrase contains enough phase information for cron. For example, `every 15 minutes` begins when its task is created, while `*/15` always begins at `:00`. Write `every 15 minutes starting at 00:00` when that is the intended cron phase; conversion refuses to invent it. The reverse translation retains that phase explicitly.

## Import a crontab

Always look first:

```sh
gosched cron import --file /etc/crontab --dry-run
```

The preview prints, for every line, the expression, phrase, effective timezone, shell command, and any run-as, environment, or stdin context. For a real import it also prints the task it created. Imported tasks retain the original normalized cron expression. Nothing is created while `--dry-run` is set.

```text
line 3: 0 2 * * *
  phrase:  every day at 02:00
  timezone: America/New_York
  command: /bin/sh -c /usr/local/bin/backup --full
line 5: */15 * * * *
  phrase:  every 15 minutes starting at 00:00
  command: /usr/local/bin/probe
line 6: @reboot
  phrase:  at scheduler startup
  command: /bin/sh -c /usr/local/bin/warm

8 line(s) read: 3 would create, 3 skipped (comments, blanks, variables),
1 unsupported, 1 error(s)
```

The preview is not advisory. Preview and creation use the same conversion result. A created task retains the normalized cron expression and its compiled recurrence; the phrase is an explanation of that result, not a replacement schedule string.

When it does, drop the flag:

```sh
gosched cron import --file /etc/crontab --timezone America/New_York --group ops
gosched cron import --file /etc/crontab --system --dry-run
gosched cron import --file quartz.cron --dialect quartz --dry-run
```

| Flag | Meaning |
| --- | --- |
| `--file` | Crontab to read, or `-` for standard input. Required. |
| `--dry-run` | Produce the identical report and create nothing. |
| `--dialect` | `unix` consumes five timing fields (default); `quartz` consumes six beginning with seconds. |
| `--system` | Consume the system-crontab username field and map it to run-as. |
| `--run-as` | Supply the owner account for a user crontab; cannot be combined with `--system`. |
| `--timezone` | IANA zone override for every task. Without it, each line uses effective `CRON_TZ` or `Local`. |
| `--group` | Group ID to file the imported tasks under. |
| `--count` | How many upcoming runs to show per line. Default 3. |

A line that cannot be converted never stops the ones that can: the supported lines are still created, and the summary counts every line of the file. Reading the file successfully is a success, whatever the mix of outcomes. A crontab of nothing but `@reboot` lines previews or creates startup tasks. Only an unreadable file, an unknown timezone, or a failed creation is a failure.

Importing the same crontab twice creates two sets of tasks. There is no deduplication; the counts are how you notice.

Assignments apply from their line onward. `CRON_TZ` controls schedule timing; ordinary variables, including `TZ`, become task environment. `SHELL` also selects the executable used with `-c`. Cron's unescaped `%` split is preserved as task stdin, including newline conversion after the first percent. `MAILTO` and `MAILFROM` remain visible warnings because output delivery belongs to the notification feature, not the child environment. `LOGNAME` assignments are also ignored with a warning because cron does not permit overriding the executing account's login name.

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

Nothing is created or changed. An expression that cannot be represented is reported by name, and that is an answer rather than a failure, the exit code stays 0. A malformed expression *is* a failure, and exits 2 naming the field.

## Export back to cron

```sh
gosched cron export
gosched cron export --task 6f1c…
```

Every task appears exactly once, as a crontab line or as a commented refusal:

```text
# gosched cron export, 4 task(s)
0 9 * * 1-5 /usr/bin/report --daily
# declined: "nightly backup", cron cannot express a schedule that fires exactly once
# declined: "health probe", the task is disabled and cron has no disabled state
*/15 * * * * /usr/local/bin/probe
```

Nothing is silently omitted and nothing is approximated. A converter that quietly rounds a schedule to the nearest thing cron can say is worse than one that declines, because the difference only surfaces at 02:30 some morning.

## Fidelity

Five-field input uses `minute hour day-of-month month day-of-week`. Numeric bounds are `0-59`, `0-23`, `1-31`, `1-12`, and `0-7`, with both `0` and `7` meaning Sunday. Six-field input uses `second minute hour day-of-month month day-of-week`, seconds `0-59`, Quartz weekdays `1=Sunday` through `7=Saturday`, and a complete-field `?` in day-of-month or day-of-week. The parser accepts numbers, comma-separated lists, inclusive ascending ranges, wildcard and range steps, and case-insensitive month or weekday names by their first three letters. Overlapping terms and the two Sunday aliases normalize to one ordered set.

Standard field combinations compile directly into the same durable recurrence model used by human schedules. Their readable description is display output, not a phrase that is parsed again. This matters for expressions such as `*/10 9-17 * * MON,WED,FRI`, whose complete meaning has no equivalent in the human authoring grammar. Existing simple forms keep their concise phrases.

This is a product-specific subset, not a promise of POSIX, Linux, or robfig parity. For the upstream dialects and file behavior, see the [POSIX crontab utility](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/crontab.html), [Cronie crontab(5)](https://github.com/cronie-crond/cronie/blob/master/man/crontab.5), [robfig/cron expression format](https://pkg.go.dev/github.com/robfig/cron/v3#hdr-CRON_Expression_Format), and the [Quartz CronTrigger modifier semantics](https://www.quartz-scheduler.org/documentation/quartz-2.5.x/tutorials/crontrigger.html). The single `weekday#ordinal`, day-of-week `weekdayL`, and day-of-month `L`, `nW`, and `LW` forms are focused subsets. Six-field input adds seconds and `?` without adding Quartz's optional year or full modifier language. Verify the destination before treating exported text as portable.

### Supported

| Cron | Notes |
| --- | --- |
| Seconds sets and steps | `*/30 * * * * *` and `5/15 * * * * *`; six-field output is used only when seconds are required. |
| Quartz no-specific-value | `0 0 12 ? * MON`; `?` must occupy one complete day field. |
| Every minute | `* * * * *` |
| Minute lists, ranges, and field-local steps | `5,20,45 * * * *`, `10-20/2 * * * *`, and uneven `*/7 * * * *` restart within each hour. |
| Hour lists, ranges, and field-local steps | `0 9,17 * * *`, `30 8-17 * * *`, `0 */5 * * *` |
| One fixed daily time | `0 9 * * *` |
| Arbitrary weekday sets and ranges | `0 14 * * WED`, `0 9 * * 1-5`, `0 10 * * SUN,TUE,THU` |
| One ordinal weekday each month | `0 9 * * 5#3` means the third Friday. One weekday and ordinal 1 through 5 are supported when month and day-of-month are unrestricted. |
| The last selected weekday each month | `0 9 * * 5L` means the last Friday. One weekday followed by `L` is supported when month and day-of-month are unrestricted; `SUNL` and `7L` import as Sunday and export canonically as `0L`. |
| Last calendar day | `0 9 L * *` means the final date of every month. |
| Nearest weekday to one date | `0 9 15W * *` keeps the 15th on Monday-Friday, moves Saturday to Friday, and Sunday to Monday without crossing a month boundary. Targets 1 through 31 are supported. |
| Last weekday of the month | `0 9 LW * *` means the final Monday-through-Friday date. This differs from `0 9 * * 5L`, which means the last Friday. |
| Date and month lists, ranges, and steps | `0 0 1,15 JAN,MAR *`, `0 9 */2 * *`, `0 9 * */2 *` |
| Safe cross-field conjunctions | Minute, hour, month, and exactly one of day-of-month or day-of-week may all be restricted, as in `*/10 9-17 * * MON,WED,FRI`. |
| Month and weekday names | Names are case-insensitive; their first three letters are significant. |
| Sunday as `0` or `7` | Both are accepted |
| `@hourly`, `@daily`, `@midnight`, `@weekly`, `@monthly`, `@yearly`, `@annually` | Expanded to their documented five-field equivalents |
| `@reboot` | A non-clock event that runs once per scheduler daemon start; `at scheduler startup` is the symmetric human form. |

### Issue #22 audit decisions

Every audited row has one disposition. "Deferred" names a different product capability; "out of scope" means the converter refuses instead of approximating.

| Row | Feature | Decision and rationale |
| --- | --- | --- |
| A1 | `@reboot` | **Supported startup event.** It runs once per daemon process start, not on reload, and has no upcoming clock occurrence. |
| A2 | Arbitrary field combinations | **Supported.** Safe combinations compile directly to one durable recurrence. |
| A3 | Minute/hour/date/month lists | **Supported.** Values remain ordered and deduplicated. |
| A4 | Arbitrary weekday sets | **Supported.** Names, lists, ranges, and Unix Sunday aliases normalize exactly. |
| A5 | Uneven steps | **Supported.** Steps reset within their field rather than pretending to be elapsed intervals. |
| A6 | Range steps | **Supported.** The selected set is compiled exactly. |
| A7 | Restricted day-of-month plus day-of-week | **Out of scope.** Cron uses OR while the single recurrence model uses conjunction; input is refused by name. |
| A8 | Six-field seconds | **Supported subset.** Seconds lists, ranges, and steps map to `BYSECOND`; Quartz year remains out of scope. |
| A9 | `L` | **Supported focused subset.** `L`, `LW`, and one `weekdayL` are supported; offsets and composites are refused. |
| A10 | `W` | **Supported focused subset.** One `1W` through `31W` and `LW` are supported; composites are refused. |
| A11 | `#` | **Supported focused subset.** One `weekday#1..5` is supported; mixtures are refused. |
| A12 | Quartz `?` | **Supported subset.** It must occupy one complete six-field day position. |
| B1 | `CRON_TZ` / `TZ` | **Supported with distinct meanings.** `CRON_TZ` controls following schedules; `TZ` remains child environment. |
| B2 | Environment assignments | **Supported.** Ordered assignments are snapshotted onto following tasks; matching outer quotes preserve boundary space. `LOGNAME` overrides are visibly ignored to match cron. |
| B3 | `MAILTO` / `MAILFROM` | **Deferred to notifications (#19).** Import warns rather than claiming delivery. |
| B4 | System-crontab user field | **Supported with `--system`.** Explicit layout prevents command corruption. |
| B5 | Run-as user | **Supported where the platform executor supports it.** System files use each user field; user files accept `--run-as`. Windows creation keeps its existing explicit refusal. |
| B6 | Percent stdin | **Supported.** Escaped percent stays literal; the first unescaped percent starts persisted stdin and later ones become newlines. |
| B7 | Shell command semantics | **Supported.** Imported text runs through effective `SHELL -c` without whitespace splitting. |
| B8 | run-parts directories | **Out of scope.** Directory discovery is a separate import format rather than a crontab line. |
| B9 | anacron | **Out of scope.** Its delay and catch-up file format requires a separate importer. |

### Remaining named refusals

- Quartz optional year, `C`, `L-n`, and unsupported modifier mixtures.
- Both day-of-month and day-of-week restricted under cron OR semantics.
- Boot events and notification delivery.
- Operational export for tasks carrying environment, run-as, or stdin, because a standalone line cannot serialize that context faithfully.

## What a timing expression cannot say

Everything below is a property your tasks gain on import, and the import summary says so rather than leaving you to discover it:

- **A timezone by itself.** A timing expression contains no zone. Crontab import applies ordered `CRON_TZ` context per line, while `--timezone` explicitly overrides every line. Direct task input uses the task timezone field.
- **Catch-up.** If the machine is off when a cron job is due, that run is simply lost. Imported tasks get `catchup one`: a single catch-up run after downtime, then the normal schedule resumes.
- **An overlap policy.** Cron will happily start a second copy of a job that is still running. Imported tasks get `overlap queue_one`.
- **A missing-date policy.** `0 9 31 * *` runs seven months in twelve and cron never mentions it. Imported tasks get `skip`, which is exactly cron's behavior, see [the CLI reference](cli.md#task) to change it. The same applies to `29W` through `31W`: imported cron skips absent target dates. Native nearest-weekday tasks can instead resolve the target under `last_valid` or `next_valid` before weekday adjustment.
- **Restart recovery.** The daemon reconstructs its schedule from durable state on restart.

The task timezone also owns Daylight Saving resolution. A nonexistent wall time moves to the next valid instant; a repeated wall time uses its first occurrence once. Those are go-schedule rules and should not be read as Linux cron parity.

## Expression versus crontab file

A schedule expression is five Unix timing fields, six Quartz-style timing fields, or a supported macro. A user crontab line adds a command. A system crontab also adds a username, and a file can include comments, variables, shell settings, or `CRON_TZ`. Use `task add --schedule` for one expression and command you provide separately. Use `cron import` when the command and file context must first be inspected.

## What this scheduler cannot say in cron

The export declines these rather than approximating them:

- One-off schedules, cron has no way to fire exactly once.
- Secondly intervals that do not divide a minute evenly. Divisible intervals such as every 30 seconds export through the supported six-field form.
- Intervals that do not divide their period evenly (`every 3 days`, `every 2 weeks`), cron repeats by calendar position, not elapsed time.
- Sub-daily intervals whose stored phase does not align with cron's field-local step, exporting `:05/:20/:35/:50` as `*/15` would silently move every run.
- Broader ordinal and calendar-selector combinations. The focused `weekday#1..5`, day-of-week `weekdayL`, and day-of-month `L`/`nW`/`LW` subsets export, but these extensions are not universal cron syntax.
- A fifth-weekday or other date-bearing task using a non-default missing-date policy. Cron would silently skip a missing occurrence; first through fourth weekdays exist every month, so their policy setting does not change output.
- Missing-date policy never changes a last-weekday schedule because every month contains a last occurrence of every weekday; all policies therefore export to the same canonical `weekdayL` expression.
- Last-day, `LW`, and `1W` through `28W` also exist in every month and export under every policy. `29W` through `31W` export only with effective `skip`; `last_valid` and `next_valid` are refused because five-field cron cannot carry those additional dates.
- Disabled tasks, cron has no disabled state, and emitting a live line for a task you deliberately stopped would be the worst possible outcome.
- Tasks carrying environment, run-as, or stdin. A standalone exported line cannot preserve that operational context, so export names the refusal.
