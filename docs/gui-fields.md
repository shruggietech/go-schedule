# GUI task editor — field reference

This page explains every field in the desktop GUI's **New Task** / **Edit Task** dialog:
what it accepts, what's required, and what each option means. It's the GUI counterpart to the
CLI contract in [`specs/001-task-scheduler/contracts/cli.md`](../specs/001-task-scheduler/contracts/cli.md).

## At a glance

| Field | Required | Format / options | Default |
|-------|----------|------------------|---------|
| **Name** | yes | any text label | — |
| **Command** | yes | a single executable (name or full path) | — |
| **Arguments** | no | one argument per line | empty |
| **Timezone** | no | IANA zone name or `Local` | `Local` |
| **Mode** | yes | `Recurring` or `One-off` | `Recurring` |
| **Schedule** | when Recurring | human-readable phrase (see below) | — |
| **One-off time** | when One-off | RFC 3339 timestamp, must be future | — |
| **Overlap** | no | `queue_one` · `skip` · `allow_concurrent` | `queue_one` |
| **Catch-up** | no | `one` · `none` | `one` |

**Mode decides which time field matters.** In `Recurring` mode the **Schedule** field is used and
**One-off time** is ignored; in `One-off` mode it's the reverse.

---

## Name

A label for the task — any text. Used only to identify the task in lists and the calendar.

## Command

The program to run — **just the executable, not a full command line.** Put any arguments in the
**Arguments** field below, not here.

- Examples: `cmd`, `python`, `notepad.exe`, `C:\Windows\System32\notepad.exe`, `/usr/bin/make-report`
- Required and must be non-empty.

## Arguments

**One argument per line.** This is the most common point of confusion: don't paste a whole
command line on one line. Each line becomes one separate argument passed to the command. Blank
lines and surrounding whitespace are ignored.

To run the equivalent of `cmd /c echo hello`:

```
/c
echo hello
```

## Timezone

An [IANA time-zone name](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) or the
literal word `Local`. Schedules are interpreted in this zone, with correct Daylight Saving Time
handling; the backend stores everything in UTC.

- Examples: `Local` (your system clock), `UTC`, `America/New_York`, `Europe/London`, `Asia/Tokyo`
- An empty field is treated as `Local`. An unknown name (e.g. `Mars/Phobos`) is rejected.

## Mode

- **Recurring** — the task fires repeatedly on a schedule. Fill in **Schedule**.
- **One-off** — the task fires exactly once at a specific time. Fill in **One-off time**.

## Schedule *(Recurring mode)*

A plain-language phrase — no cron syntax. Parsing is case-insensitive. Supported forms:

| Pattern | Examples |
|---------|----------|
| Fixed interval | `every 15 minutes`, `every 30s`, `every 2 hours`, `every 3 days`, `every week` |
| Daily with a time | `every day at 09:00` |
| Weekday / weekend sets | `weekdays at 9:00 AM`, `weekends at 18:00` |
| A single weekday | `every monday at 9am` |
| Monthly ordinal weekday | `3rd wednesday monthly at 14:00`, `last friday of the month` |

**Units** (any spelling): `second`/`sec`/`s`, `minute`/`min`/`m`, `hour`/`hr`/`h`, `day`/`d`,
`week`/`w`.

**Ordinals:** `1st`–`5th`, `first`–`fifth`, or `last`. The monthly clause can be written as
`monthly`, `of the month`, `of each month`, or `of every month`.

**Time-of-day** accepts: `14:00`, `9:00`, `9:00 AM`, `9am`, or a bare hour like `9` (= 09:00).
Hours are 0–23, minutes 0–59.

> ⚠️ **Sub-daily intervals can't take an `at` time.** Seconds/minutes/hours fire on a rolling
> interval, so `every 15 minutes at 09:00` is **rejected**. The `at <time>` clause is only valid
> for daily-or-coarser schedules (`every day`, `weekdays`, `every monday`, monthly ordinals).

As you type a valid Schedule, the **Preview** row fills in with a plain-English summary plus the
next few run times — a quick way to confirm your phrase parsed the way you meant.

## One-off time *(One-off mode)*

An [RFC 3339](https://datatracker.ietf.org/doc/html/rfc3339) timestamp, exactly like the
placeholder `2026-08-04T09:00:00Z`. The trailing `Z` means UTC; you can also use an explicit
offset such as `2026-08-04T09:00:00-04:00`.

- Must be in the **future**, or the task is rejected.
- Malformed input shows: *"one-off time must be RFC3339, e.g. 2026-08-04T09:00:00Z"*.

## Overlap

What to do when a task is still running at the moment its next run would start:

- **`queue_one`** *(default)* — queue exactly one pending run; drop any further triggers until the
  current run finishes. A warning is logged and surfaced as a GUI alert.
- **`skip`** — skip the new trigger entirely; do nothing until the next scheduled time.
- **`allow_concurrent`** — let multiple runs of the same task execute at the same time.

## Catch-up

What to do after downtime (the daemon was stopped) when one or more scheduled runs were missed:

- **`one`** *(default)* — run once to catch up, then resume the normal schedule.
- **`none`** — skip all missed runs and resume the normal schedule.

---

## A known-good example

A "heartbeat" task you can watch succeed within a couple of minutes:

| Field | Value |
|-------|-------|
| Name | `heartbeat` |
| Command | `cmd` |
| Arguments | `/c` (line 1) · `echo %DATE% %TIME% >> C:\Users\you\gosched-test.txt` (line 2) |
| Timezone | `Local` |
| Mode | `Recurring` |
| Schedule | `every 1 minute` |
| Overlap | `queue_one` |
| Catch-up | `one` |

After saving, a new timestamp line should appear in the file about once a minute.
