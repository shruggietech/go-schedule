---
title: GUI field reference
nav_order: 4
---

# GUI task editor — field reference

> **Release status:** cron entry in the Schedule field is currently an
> unreleased change planned for the first release after 0.8.0. Released 0.8.0
> builds accept the human forms documented below.

This page explains every field in the desktop GUI's **New Task** / **Edit Task** dialog:
what it accepts, what's required, and what each option means. It's the GUI counterpart to the
CLI contract in [`specs/001-task-scheduler/contracts/cli.md`](https://github.com/shruggietech/go-schedule/blob/main/specs/001-task-scheduler/contracts/cli.md).

## Data view tables

The desktop's **Tasks**, **Schedule** List, and **Activity** views use fixed
headers above vertically scrolling rows. Columns contract proportionally when
the window narrows, so these views do not introduce horizontal scrolling. Long
values are shortened visually with an ellipsis. Selecting a Tasks or Schedule
row exposes its complete labeled values below the list; selecting an Activity
row opens its existing full-detail dialog.

| View | Columns | Missing-value behavior | Semantic labels |
| --- | --- | --- | --- |
| **Tasks** | Task, Enabled, Lifecycle, Time zone, Group | `Unnamed task`, `Unknown`, or `Not assigned` | Enabled and lifecycle are separate concepts. Enabled uses success treatment, while Disabled is muted. Active uses informational treatment. |
| **Schedule** List | When, Task, Event, Outcome | `Unnamed task` or `NOT AVAILABLE` | `▷ SCHEDULED`; completed outcomes use `✓ SUCCESS`, `✗ FAILURE`, `↷ SKIPPED`, `↻ CAUGHT UP`, or `⋯ QUEUED`. Unknown values remain neutral and readable. |
| **Activity** | When, Severity, Source, Summary | `daemon` or `No message` | Severity is always uppercase: `• INFO`, `⚠ WARNING`, or `✗ ERROR`. Unknown values are normalized to uppercase with a neutral `?` glyph. |

Alternating rows use a quiet theme-aware surface. Selection, hover, focus, and
semantic meaning retain a text or glyph cue in every appearance mode, so color
is never the only indicator. Row identity, rather than the current visual
index, controls selection and activation across live refreshes.

The dialog is a two-pane layout. The **left** pane holds the form, grouped into **What to run**
(Name and Command line), **When** (Timezone, Mode, the relevant time field), and a collapsible
**Advanced Settings** (Overlap, Catch-up, Missing dates, Time basis, Spring gap,
Fall overlap) that starts closed - its disclosure arrow points ▶ when
collapsed and ▼ when expanded. The **right** pane shows the live **Preview** by default, with a
**Help** button that swaps it to a field-by-field guide (and back). Required fields are marked with
a `*`, and the **Save** button (bottom-right, next to **Cancel**) stays disabled until every
required field is valid. Clicking **Cancel** after you've typed something asks for confirmation
before discarding.

## At a glance

| Field | Required | Format / options | Default |
|-------|----------|------------------|---------|
| **Name** | yes | any text label | — |
| **Command line** | yes | program followed by arguments in the portable direct-command syntax below | — |
| **Group** | no | `(none)`, or a group shown by its path (`Backups / Nightly`) | `(none)` |
| **Timezone** | no | searchable list of common zones, or any IANA name / `Local` | `Local` |
| **Mode** | yes | `Recurring` or `One-off` | `Recurring` |
| **Schedule** | when Recurring (create) | human phrase or supported five- or six-field cron (see below) | n/a |
| **Start at** | no | anchor time for sub-daily intervals, e.g. `09:00` | — |
| **One-off date / time** | when One-off (create) | date + time picked in the task's zone, must be future | — |
| **Overlap** *(Advanced)* | no | Queue one run · Skip this run · Allow concurrent runs | Queue one run |
| **Catch-up** *(Advanced)* | no | Run once to catch up · Skip missed runs | Run once to catch up |
| **Missing dates** *(Advanced)* | no | Skip that period · Use the last valid date · Roll into the next period | Skip that period |
| **Time basis** *(Advanced)* | no | Local wall clock · Fixed elapsed time · UTC clock | Local wall clock |
| **Spring gap** *(Advanced)* | no | Run at the next valid time · Skip that occurrence | Run at the next valid time |
| **Fall overlap** *(Advanced)* | no | First occurrence · Both occurrences · Last occurrence | First occurrence |

**Mode decides which time field is shown.** In `Recurring` mode the **Schedule** (and optional
**Start at**) field is shown and the one-off inputs are hidden; in `One-off` mode it's the reverse.
Switching Mode keeps whatever you already typed in either field. When editing an existing task,
leaving the time field blank keeps the task's current schedule.

**Editing shows the task as it actually is.** Opening an existing task fills in its real Mode and
either its retained schedule expression or its one-off date and time, shown in the task's own timezone, so you
can see what the task is currently set to before changing anything. Saving without touching those
fields leaves the schedule exactly as it was. If you *switch* Mode, the new mode's time fields
become required — there is no existing schedule of the new kind to fall back on.

Tasks created before this was added have no stored schedule phrase, so their Schedule field opens
blank. That is safe — a blank field keeps the existing schedule — and typing a new expression replaces
it.

**Live Preview.** The right pane's Preview shows two things at once: a plain-language summary of
the schedule with the next few run times, and the exact Program plus numbered Arguments in order.
Each value uses a quoted escaped display, so empty values, spaces, tabs, line breaks, quotes, and
backslashes are visible instead of being hidden inside an ambiguous reconstructed command string.
Changing **Missing dates** immediately refreshes this preview and sends the
selected policy to the daemon, so date-sensitive schedules show the same runs
that Save will create.

**Overlap, Catch-up and Missing dates** are shown with friendly labels but stored as the same
underlying policy values (`queue_one`/`skip`/`allow_concurrent`, `one`/`none`, and
`skip`/`last_valid`/`next_valid`) used by the CLI and API.

---

## Name

A label for the task — any text. Used only to identify the task in lists and the calendar.

## Command line

Type the executable and arguments together in the same roomy field:

```text
python -m http.server --bind "127.0.0.1"
"C:\Program Files\Tool\tool.exe" --name "Ada Lovelace"
/usr/bin/printf '%s\n' 'hello world'
```

The editor uses one portable direct-command grammar with identical value boundaries on Windows,
macOS, and Linux:

- Whitespace separates values outside quotation marks.
- Single or double quotation marks keep spaces and literal line breaks inside one value. The
  quotation marks are not part of the value.
- Empty quotes (`''` or `""`) create an intentional empty argument.
- Adjacent quoted and unquoted parts form one value.
- Inside single quotes every enclosed character is literal. Inside double quotes a backslash can
  escape a double quote; other backslashes remain literal.
- Outside quotes, a backslash can escape whitespace or a quotation mark. Before an ordinary
  character, another backslash, or the end of the field, it remains literal.
- Invalid UTF-8 text, an unmatched quote, or an unsupported NUL character is an error. The Preview
  identifies the problem (and the line and column for character-level syntax errors), and Save
  remains disabled until it is corrected.

This syntax only separates a program from its arguments. It does **not** expand environment
variables or wildcards, create pipelines, perform redirects, interpret comments, or run a shell.
Characters such as `$`, `%`, `*`, `|`, `>`, `;`, and `&` stay literal.

To request shell behavior, name the shell explicitly. For example, use
`cmd /c "echo hello > output.txt"` on Windows or `sh -c 'echo hello > output.txt'` on a POSIX
host. That named shell then owns its platform-specific quoting, expansion, redirection, and
security behavior. go-schedule still stores and launches one program plus its ordered arguments.

The field is required, displays at least six lines at the default dialog size, and grows when the
dialog gains vertical space. A quoted literal line break is part of one argument; an unquoted line
break separates values.

## Group

Which group the task belongs to, or `(none)` for no group. Groups can be enabled and disabled as a
unit, and disabling a group suppresses every task inside it and inside its subgroups.

- Groups are listed by their full path (`Backups / Nightly`), so two groups with the same name at
  different levels are distinguishable.
- Choose `(none)` to take a task back out of its group. It then appears under **Ungrouped** in the
  Groups tab.
- Create groups in the **Groups** tab; this field only assigns to existing ones.
- You can also move a task from the Groups tab: select it under its group and use
  **Move to group…**. Both paths offer the same choices.
- The equivalent CLI form is `gosched task add --group <id>` / `gosched task edit --group <id>`,
  and `gosched task edit --group ""` to un-group. Omitting `--group` leaves membership unchanged.

## Timezone

An [IANA time-zone name](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) or the
literal word `Local`. The field is a searchable dropdown seeded with common zones; you can pick
one from the list or type any other valid IANA name. Schedules are interpreted in this zone, with
correct Daylight Saving Time handling; the backend stores everything in UTC.

- Examples: `Local` (your system clock), `UTC`, `America/New_York`, `Europe/London`, `Asia/Tokyo`
- An empty field is treated as `Local`. An unknown name (e.g. `Mars/Phobos`) is rejected and
  blocks Save.

## Mode

- **Recurring** — the task fires repeatedly on a schedule. Fill in **Schedule**.
- **One-off** — the task fires exactly once at a specific time. Fill in **One-off time**.

## Schedule *(Recurring mode)*

Use a plain-language phrase (the approachable default) or a supported five- or six-field cron expression.
Human parsing is case-insensitive. Supported human forms include:

| Pattern | Examples |
|---------|----------|
| Fixed interval | `every 15 minutes`, `every 30s`, `every 2 hours`, `every 3 days`, `every week` |
| Daily with a time | `every day at 09:00` |
| Weekday / weekend sets | `weekdays at 9:00 AM`, `weekends at 18:00` |
| A single weekday | `every monday at 9am` |
| Monthly ordinal weekday | `3rd wednesday monthly at 14:00`, `last friday of the month` |
| Monthly calendar selector | `last day of every month at 09:00`, `nearest weekday to the 15th of every month at 09:00`, `last weekday of every month at 09:00` |
| Scheduler startup | `at scheduler startup`, `@reboot` |

**Units** (any spelling): `second`/`sec`/`s`, `minute`/`min`/`m`, `hour`/`hr`/`h`, `day`/`d`,
`week`/`w`.

**Ordinals:** `1st`–`5th`, `first`–`fifth`, or `last`. The monthly clause can be written as
`monthly`, `of the month`, `of each month`, or `of every month`.

**Time-of-day** accepts: `14:00`, `9:00`, `9:00 AM`, `9am`, or a bare hour like `9` (= 09:00).
Hours are 0–23, minutes 0–59.

Supported cron uses the familiar `minute hour day-of-month month day-of-week` shape. For example,
`0 9 * * 1-5` runs at 09:00 on weekdays, while `0 9 L * *`, `0 9 15W * *`,
and `0 9 LW * *` select the last day, nearest weekday to the 15th, and last
weekday. The editor detects the syntax from the current field text,
so replacing cron with a human phrase, or the reverse, updates both preview and save behavior. A
cron-shaped value that is invalid or cannot be represented faithfully is rejected as cron rather
than retried as a human phrase. See the [cron fidelity contract](cron.md#fidelity) for supported
fields, macros, and explicit refusals.

Seconds precision uses the Quartz-style `second minute hour day-of-month month
day-of-week` shape. For example, `*/30 * * * * *` runs twice per minute and
`0 0 12 ? * MON` runs at noon each Monday. Numeric weekday values differ by
dialect, so consult the fidelity contract before converting numeric six-field
input.

When editing, the exact retained cron expression is shown rather than translated
back from its readable preview. Standard lists, ranges, and field-local steps
such as `*/10 9-17 * * MON,WED,FRI` and `0 9 */2 * *` retain their exact timing.
The preview is display text; the compiled recurrence, not a generated English
phrase, is authoritative for execution.

Startup schedules are non-clock events. Their preview reads **At scheduler
startup** and intentionally has no next-run list. They run once for each daemon
process start, never merely because a task was edited, enabled, imported, moved,
or reloaded. This does not prove a physical host reboot, and other services may
not be ready when the command begins.

> ⚠️ **Sub-daily intervals can't take an `at` time.** Seconds/minutes/hours fire on a rolling
> interval, so `every 15 minutes at 09:00` is **rejected**. The `at <time>` clause is only valid
> for daily-or-coarser schedules (`every day`, `weekdays`, `every monday`, monthly ordinals).

As you type a valid Schedule, the **Preview** pane fills in with a plain-English summary plus the
next few run times — a quick way to confirm your expression parsed the way you meant. The **Help**
button (top of the right pane) shows the full list of supported phrasings and a guide to every
field.

### Start at *(sub-daily human intervals only)*

By default a fixed interval like `every 15 minutes` is anchored to the moment you create the task,
so it might fire at an awkward phase (6:07, 6:22, 6:37 …). To align the cycle, set a **Start at**
time — a separate field that appears only when the Schedule is a sub-daily interval. With
`every 15 minutes` and a Start at of `09:00`, runs fall on `:00 / :15 / :30 / :45`.

Equivalently, you can type the anchor directly in the Schedule using a `starting at` (or `from`)
clause — the GUI and CLI both understand it:

| Phrase | Effect |
|--------|--------|
| `every 15 minutes starting at 09:00` | aligns to `:00/:15/:30/:45` |
| `every 30 minutes from 9am` | aligns to `:00/:30` relative to 09:00 |
| `every 2 hours starting at 08:00` | fires at 08:00, 10:00, 12:00 … |

The anchor is interpreted in the task's **Timezone**. It applies only to sub-daily human intervals;
`every day starting at 09:00` is rejected (use `every day at 09:00`).
Cron expressions carry their phase in their own fields and do not show a separate Start at row.

## One-off date / time *(One-off mode)*

Pick the **Date** (`2026-08-04`) and **Time** (`09:00`) in two fields — no hand-typed RFC 3339
required. The instant is interpreted in the task's **Timezone**, and a line under the fields
echoes the resolved run time so you can confirm it.

- Must be in the **future**, or Save stays disabled.
- The backend still stores the moment as a UTC instant, exactly as before.

## Advanced Settings

The **Overlap**, **Catch-up**, **Missing dates**, **Time basis**, **Spring gap**,
and **Fall overlap** controls live in a collapsible **Advanced
Settings** section that starts closed. They are shown with human-readable labels; the stored policy values (used by the
CLI and API) are unchanged.

## Completion chains

The **Chains** view is separate from the task editor because a chain is a
relationship between two existing tasks. **New** selects an **After task**, a
**Run task**, and one plain-language condition:

- **Only when the source succeeds**
- **Only when the source fails**
- **Whenever the source finishes**

Task choices show both name and stable ID, so duplicate names stay
distinguishable. The list uses current task names and updates live after create,
edit, or delete. A chain adds another way to start the target; its timed or
startup schedule remains unchanged. Disabled or inactive targets are resolved
without execution rather than accumulating hidden work.

### Overlap

What to do when a task is still running at the moment its next run would start:

- **Queue one run** (`queue_one`, *default*) — queue exactly one pending run; drop any further
  triggers until the current run finishes. A warning is logged and surfaced as a GUI alert.
- **Skip this run** (`skip`) — skip the new trigger entirely; do nothing until the next scheduled
  time.
- **Allow concurrent runs** (`allow_concurrent`) — let multiple runs of the same task execute at
  the same time.

### Catch-up

What to do after downtime (the daemon was stopped) when one or more scheduled runs were missed:

- **Run once to catch up** (`one`, *default*) — run once to catch up, then resume the normal
  schedule.
- **Skip missed runs** (`none`) — skip all missed runs and resume the normal schedule.

### Missing dates

What to do in a period that has no matching date. This applies only to schedules that can
actually miss one — the 29th, 30th or 31st of a month, a yearly rule on 29 February, and the
fifth of a given weekday. For everything else the setting is inert and changes no run time.

- **Skip that period** (`skip`, *default*) — no run that period. This is what cron does, and what
  every task created before this setting existed already did.
- **Use the last valid date** (`last_valid`) — fall back to the last date that exists: the 31st
  becomes the 30th, or the 28th in February; a missing fifth Friday becomes the last Friday.
- **Roll into the next period** (`next_valid`) — run on the 1st of the following month instead,
  without displacing that month's own run.

The Preview names whichever you pick, so a schedule that skips months says so rather than
claiming "every month".

### Time basis and DST transitions

- **Local wall clock** (`wall_clock`, *default*) keeps calendar and clock
  readings fixed in the task timezone. Elapsed gaps can therefore be shorter or
  longer when the UTC offset changes.
- **Fixed elapsed time** (`elapsed`) keeps exact seconds between interval runs
  and lets their local display shift. Calendar-selected schedules are refused
  because a month or ordinal weekday is not a fixed duration.
- **UTC clock** (`utc`) keeps recurrence readings fixed in UTC. The task
  timezone is used only to display the resulting instants locally.

For wall-clock schedules, **Spring gap** either advances a nonexistent reading
to the first valid instant (`next_valid`) or omits it (`skip`). **Fall overlap**
selects the earlier (`first`), both (`both`), or later (`last`) instant. These
choices remain saved but have no effect while the basis is elapsed or UTC.

---

## A known-good example

A "heartbeat" task you can watch succeed within a couple of minutes:

| Field | Value |
|-------|-------|
| Name | `heartbeat` |
| Command line | `cmd /c "echo %DATE% %TIME% >> C:\Users\you\gosched-test.txt"` |
| Timezone | `Local` |
| Mode | `Recurring` |
| Schedule | `every 1 minute` |
| Overlap *(Advanced)* | Queue one run |
| Catch-up *(Advanced)* | Run once to catch up |
| Missing dates *(Advanced)* | Skip that period |
| Time basis *(Advanced)* | Local wall clock |
| Spring gap *(Advanced)* | Run at the next valid time |
| Fall overlap *(Advanced)* | First occurrence |

After saving, a new timestamp line should appear in the file about once a minute.
