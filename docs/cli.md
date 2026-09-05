---
title: CLI reference
nav_order: 3
---

# `gosched` command reference

**Audience:** anyone using go-schedule from a terminal\
**Applies to:** go-schedule 0.7.0 and later; see the release-status note below\
**Source of truth:** `internal/cli/`, this document describes what the binary does. `specs/001-task-scheduler/contracts/cli.md` describes what it must do, and remains a specification artifact.

> **Release status:** dual-syntax `task add` / `task edit` input and `cron convert` are currently unreleased changes planned for the first release after 0.8.0. Existing import, explain, and export workflows remain applicable from 0.7.0.

`gosched` is a thin client. Every command below talks to the `goschedd` daemon over local IPC, a Unix socket on Linux and macOS, a named pipe on Windows, so the CLI and the desktop GUI act on identical state, and the schedule keeps running whether or not either is open.

## Contents

- [Conventions](#conventions)
- [Global flags](#global-flags)
- [Exit codes](#exit-codes)
- [`health`](#health)
- [`task`](#task)
- [`cron`](#cron)
- [`group`](#group)
- [`runs`](#runs)
- [`logs`](#logs)
- [`service`](#service)
- [`gui`](#gui)
- [Deprecated: `alerts`](#deprecated-alerts)

## Conventions

Commands are written bare, `gosched task list`, not a full path. On Windows that requires the installer's `PATH` entry, which is present from 0.6.0 onward and needs a **newly opened** shell to be visible. See the [Windows install guide](INSTALL-windows.md).

Task and group identifiers are UUIDs assigned by the daemon and printed when the object is created. Anywhere `<id>` appears, that is what it means.

Times you supply are RFC 3339 (`2026-08-04T09:00:00Z`). Times printed back are RFC 3339 as well. Internally everything is UTC; the per-task timezone decides when "09:00" happens, including across a Daylight Saving transition.

## Global flags

| Flag | Effect |
| --- | --- |
| `--json` | Emit machine-readable JSON instead of the table or summary. Available on every command that produces output. |
| `-v`, `--version` | Print the CLI version. This is the *client* version; see [`health`](#health) for the daemon's. |
| `-h`, `--help` | Help for any command or subcommand. |

## Exit codes

| Code | Meaning |
| --- | --- |
| `0` | Success. |
| `1` | Runtime failure, the daemon was unreachable, the operation was refused, or the OS denied the request. |
| `2` | Usage or validation failure, a missing required flag, a malformed timestamp, an invalid policy value. Server-side validation failures map here too, so bad input is reported the same way wherever it is caught. |

Results go to stdout; diagnostics and errors go to stderr. That split is what makes `gosched task list --json | ...` safe in a pipeline.

## `health`

Check that the daemon is running and report its version.

```sh
gosched health
```

```text
daemon ok (version 0.6.0)
```

Worth knowing: the version printed here is the **daemon's**, which can differ from `gosched --version` after a partial upgrade. If you are filing a bug report, both are asked for, and that is why.

## `task`

Create and manage tasks. A task is a command, a schedule, and the policies that decide what happens when runs collide or are missed.

### `task add [name]`

The name, command, and timing are optional at creation so unfinished work can be saved safely. An omitted name displays as `unnamed`; a task without a command is not runnable, while a task with a command but no schedule is manual-only and remains available through `task run-now`.

Create a complete task or save an incomplete draft. `--schedule` and `--at` are mutually exclusive when either is supplied.

| Flag | Meaning |
| --- | --- |
| `--command` | Program or script to run. Omit it to save a non-runnable draft. |
| `--arg` | An argument to the command. Repeatable; each use adds one argument, so values containing spaces stay intact. |
| `--cwd` | Working directory for the run. |
| `--env` | An environment variable as `KEY=VALUE`. Repeatable. |
| `--group` | Group ID to file the task under. |
| `--tz` | IANA timezone, e.g. `America/New_York`. Defaults to the system zone. |
| `--schedule` | Human-readable schedule or supported five- or six-field cron expression, including `at scheduler startup` / `@reboot`. |
| `--at` | One-off run time, RFC 3339. |
| `--overlap` | `queue_one` (default), `skip`, or `allow_concurrent`. |
| `--catchup` | `one` (default) or `none`. |
| `--missing-date` | `skip` (default), `last_valid`, or `next_valid`. |
| `--time-basis` | `wall_clock` (default), `elapsed`, or `utc`. |
| `--dst-gap` | `next_valid` (default) or `skip`. |
| `--dst-overlap` | `first` (default), `both`, or `last`. |

```sh
gosched task add nightly-backup \
  --command /usr/local/bin/backup.sh \
  --schedule "every day at 02:30" \
  --tz America/New_York
```

```sh
gosched task add release-announce \
  --command /usr/bin/curl --arg -X --arg POST --arg https://example.test/hook \
  --at 2026-08-04T09:00:00Z
```

The schedule can be written the way you would say it, such as `every 15 minutes`, `every weekday at 09:00`, or `3rd wednesday monthly at 14:00`. The documented cron subset is also accepted, such as `0 9 * * 1-5` or the seconds-precision `*/30 * * * * *`. Use `@reboot` or `at scheduler startup` for one run per daemon start. This is a daemon lifecycle event, so task detail prints no upcoming times and edits, imports, enables, and reloads wait until the next daemon start. On success the command echoes back how it understood you, and the next few run times. It also names the effective timing basis and transition behavior, so a misreading is visible immediately rather than at 02:30 tomorrow:

```text
created task 6f1c… (nightly-backup)
schedule: every day at 02:30 (America/New_York)
timing: Local wall clock; spring gap: next valid; fall overlap: first
next runs:
  2026-07-24T06:30:00Z
  2026-07-25T06:30:00Z
```

**Overlap policy** decides what happens when a run is still going as the next one comes due. `queue_one` holds exactly one pending run and drops any further ones, which is almost always what you want; `skip` discards the new run outright; `allow_concurrent` lets them run side by side.

**Catch-up policy** decides what happens when the machine was off. `one` fires a single catch-up run after downtime and then resumes the normal schedule, so a task that missed forty runs fires once, not forty times. `none` skips the missed window entirely.

**Missing-date policy** decides what happens in a period that has no matching date. It applies only to schedules that can actually miss one: the 29th, 30th or 31st of a month, a yearly rule on 29 February, and the fifth of a weekday. Everything else ignores it entirely.

`skip` is the default and is what cron does: `on the 31st of every month` runs in the seven months that have a 31st and not in the other five. `last_valid` falls back to the last date that does exist, the 30th, or the 28th in February. `next_valid` rolls forward into the following period, landing on the 1st, without displacing that period's own run.

Whichever you choose, the schedule describes itself honestly. A rule that skips months says so rather than claiming "every month":

```text
schedule: The 31st of every month at 09:00, or the last day of the month when
there is no such date
```

**Time basis** decides which clock anchors a recurrence. `wall_clock` keeps local readings fixed, so a six-hour local cycle can span five or seven elapsed hours when the offset changes. `elapsed` keeps the real interval fixed and lets the displayed local reading shift; it is accepted only for fixed-duration interval schedules. `utc` evaluates recurrence fields against UTC and uses the task timezone only for local display.

For `wall_clock`, **DST gap** decides whether a nonexistent spring-forward time runs at the first valid instant (`next_valid`) or is omitted (`skip`). **DST overlap** chooses the earlier (`first`), both (`both`), or later (`last`) instant when a fall-back wall reading occurs twice. The transition choices stay stored but are inert under `elapsed` and `utc`.

### `task edit <id>`

Modify a task. Only the fields you pass change; everything else is left alone. The flags are those of `task add`, with two differences worth knowing before you use them:

Pass `--name ""` or `--command ""` to clear that value, and pass `--clear-schedule` to remove automatic timing. Clearing command or the final automatic source disables the task atomically.

- `--arg` and `--env` **replace** the existing set rather than appending to it. Pass the full list you want.
- `--group` is three-way. Omit it and group membership is untouched; pass a group ID to move the task; pass an empty string (`--group ""`) to remove the task from its group.

```sh
gosched task edit 6f1c… --schedule "every weekday at 07:00"
```

At most one of `--schedule` or `--at` may be given, since they are two ways of answering the same question.

### `task list`

```sh
gosched task list
gosched task list --group 4b2e… --state active
```

| Flag | Meaning |
| --- | --- |
| `--group` | Show only tasks in this group. |
| `--state` | `active`, `completed`, or `disabled`. |

### `task show <id>`

Full detail for one task, including command, timezone, state, all effective scheduling policies, how its schedule was understood, and upcoming run times.

### `task enable <id>` · `task disable <id>`

Stop or resume scheduling without deleting anything. A disabled task keeps its history and its definition.

### `task rm <id>`

Delete a task.

### `task run-now <id>`

Trigger an immediate run, outside the schedule. The scheduled runs are unaffected; this is the "does it actually work" button.

## `cron`

Convert strings and crontab data locally. Supported cron can also be supplied to `task add` and `task edit` through `--schedule`; invalid or unfaithful forms are refused rather than retried as human text. The desktop Schedule field accepts the same two forms and retains cron text exactly when editing.

The full guide, including the table of what each direction can and cannot carry, is [Cron interoperability](cron.md). In brief:

### `cron convert [--to cron|human] <schedule-string>`

Translate exactly one string in either direction without contacting the daemon or changing a task:

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

gosched cron convert "0 9 L * *"
# last day of every month at 09:00

gosched cron convert "nearest weekday to the 15th of every month at 09:00"
# 0 9 15W * *

gosched cron convert "0 9 LW * *"
# last weekday of every month at 09:00

gosched cron convert "*/10 9-17 * * MON,WED,FRI"
# every 10 minutes during hours 9 through 17 on Monday, Wednesday, and Friday

gosched cron convert "@reboot"
# at scheduler startup
```

Automatic mode treats `@`-prefixed input and five or six cron-shaped fields as cron. Existing human forms such as `every 15 minutes from 9am` remain human input. Use `--to cron` to force human input or `--to human` to force cron input. Quoting is the same in POSIX shells and PowerShell: place a schedule containing spaces in single or double quotes so it arrives as one argument.

For a broad cron expression, cron-to-human output is an exact readable description, not necessarily text accepted by the human schedule grammar. Execution and storage use the compiled recurrence, while the original normalized cron remains the editable source.

Default success is exactly one converted line on stdout. Invalid or lossy input produces no stdout, a named stderr diagnostic, and exit code 2. With global `--json`, success writes the stable conversion object to stdout; refusal writes the same five fields (`input_syntax`, `output_syntax`, `input`, `output`, and `refusal_reason`) to stderr and still exits 2.

Use `convert` for one pure string, `explain` for a cron expression and any applicable upcoming runs, `import` for a crontab file, and `export` for stored tasks. `@reboot` explains successfully without upcoming times.

### `cron explain <expression>`

Print the plain-language phrase an expression maps to, and its next run times. Creates nothing.

```sh
gosched cron explain "0 9 * * 1-5"
gosched cron explain "0 9 * * 5#3"
gosched cron explain "0 9 * * 5L"
gosched cron explain "0 9 L * *"
gosched cron explain "0 9 15W * *"
gosched cron explain "0 9 LW * *"
gosched cron explain "0 9,17 * * *"
gosched cron explain "*/10 9-17 * * MON,WED,FRI"
gosched cron explain "*/30 * * * * *"
gosched cron explain "@reboot"
```

`--timezone` sets the zone the run times are shown in; `--count` how many to show. An expression that cannot be represented is reported by name and exits 0, a refusal is an answer. A malformed expression exits 2 naming the field.

### `cron import`

Read a crontab and create a task per line.

```sh
gosched cron import --file /etc/crontab --dry-run
gosched cron import --file /etc/crontab --system --dry-run
gosched cron import --file quartz.cron --dialect quartz --dry-run
```

A line such as `0 9 * * 5#3 /usr/local/bin/report` previews as the third Friday of each month and retains the original cron expression when the task is created. Likewise, `0 9 * * 5L /usr/local/bin/report` previews as the last Friday of each month and retains the `5L` source. Day-of-month `L`, `15W`, and `LW` lines likewise preview their last-day, nearest-weekday, or last-weekday meaning and retain the exact timing source. Standard lists, ranges, field-local steps, names, and safe cross-field combinations use the same preview and import path. A restricted day-of-month combined with a restricted day-of-week is still refused because cron applies OR semantics that this recurrence model cannot reproduce.

| Flag | Meaning |
| --- | --- |
| `--file` | Crontab to read, or `-` for standard input. **Required.** |
| `--dry-run` | Produce the identical report and create nothing. |
| `--dialect` | `unix` for five timing fields (default), or `quartz` for six. |
| `--system` | Consume the system-crontab user field and map it to run-as. |
| `--run-as` | Supply the owner account for a user crontab; cannot be combined with `--system`. |
| `--timezone` | IANA zone override for all tasks; otherwise `CRON_TZ` applies per line. |
| `--group` | Group ID to file them under. |
| `--count` | Upcoming runs shown per line. Default 3. |

Always run it with `--dry-run` first. Preview and creation use the same cron conversion result. The created task retains the normalized cron expression and its compiled recurrence; the displayed phrase is an explanation, not the value stored in place of the expression.

### `cron export`

Emit the task set as crontab lines.

```sh
gosched cron export
gosched cron export --task <id>
```

Every task appears exactly once: a crontab line where cron can carry the schedule, and a `# declined:` comment naming the task and the reason where it cannot. Nothing is approximated and nothing is omitted.

## `group`

Groups nest, and enabling or disabling one cascades through everything beneath it. That is the point of them: one command to silence a whole subtree.

### `group add <name>`

```sh
gosched group add backups
gosched group add databases --parent 4b2e…
```

`--parent` takes a group ID; omit it for a top-level group.

### `group list`

```sh
gosched group list
gosched group list --tree
```

`--tree` prints the hierarchy with disabled groups marked, rather than a flat table.

### `group enable <id>` · `group disable <id>`

Applies to the group **and its whole subtree**.

### `group rm <id>`

Delete a group. Child groups cascade; tasks are not deleted, they become ungrouped.

## `runs`

Run history: what was scheduled, what happened, and how it was triggered.

```sh
gosched runs
gosched runs --task 6f1c… --limit 20
```

| Flag | Meaning | Default |
| --- | --- | --- |
| `--task` | Filter to one task ID. | all tasks |
| `--limit` | Maximum rows. | `50` |

The `SOURCE TASK` and `SOURCE RUN` columns identify the upstream completion for a chained execution; they are `-` for schedule, startup, catch-up, and manual runs. The `EXIT` column is the process exit code, or `-` where there isn't one, a run that never started has no exit code, and printing `0` for it would be a lie.

## `chain`

Completion chains run a target task after a source task reaches a terminal result. They supplement both tasks' normal schedules.

```sh
gosched chain create --source <task-id> --target <task-id> --on success
gosched chain list
gosched chain show <chain-id>
gosched chain update <chain-id> --on any
gosched chain rm <chain-id>
```

`--on` accepts `success`, `failure`, or `any`. A source and target must be different tasks, duplicate relationships are refused, and the complete graph must remain acyclic. Add `--json` to create, list, show, or update for the API representation, including task names and stable IDs.

## `logs`

The CLI returns a bounded set of recent daemon records. Scheduler alerts appear alongside those records in the desktop GUI's Activity view.

```sh
gosched logs
gosched logs --severity error --limit 200
```

| Flag | Meaning | Default |
| --- | --- | --- |
| `--severity` | `info`, `warning`, or `error`. | all |
| `--limit` | Maximum rows. | `100` |

The Activity view identifies itself as a limited recent view and displays the exact configured path to the daemon's complete rotating JSONL log. Platform install guides list the default locations, but `log_file_path` overrides them; the path reported in Activity is authoritative for the running daemon.

## `service`

Manage the system-wide background service, so the scheduler starts on boot and runs whether or not anyone is logged in.

| Subcommand | Effect | Elevation |
| --- | --- | --- |
| `install` | Register the daemon with the system service manager. | **Required** |
| `uninstall` | Remove the registration. | **Required** |
| `start` | Start the service. | **Required** |
| `stop` | Stop the service. | **Required** |
| `restart` | Stop, then start. | **Required** |
| `status` | Report `running`, `stopped`, or that it is not installed. | Not required |

```sh
sudo gosched service install
sudo gosched service start
gosched service status
```

`status` is deliberately the one subcommand an ordinary user can run. It asks the operating system for no more access than a read needs, so it answers for an unprivileged caller wherever the service's own permissions allow a status query which, for a service installed by go-schedule, they do. Before 0.6.0 it requested start and stop rights it never used and failed with `Access is denied` for anyone not elevated, which reported that permission was withheld when in fact it was granted.

The other five genuinely change system state and genuinely require elevation. That is not being relaxed.

## `gui`

```sh
gosched gui
```

Launches the desktop application and detaches. On Windows no console window appears, which is why launching it this way is preferable to running the GUI binary from a shell.

The GUI must be present next to the `gosched` binary. If it is not, a server-only install, for instance, the command says so and names the path it looked in.

## Deprecated: `alerts`

`gosched alerts` and `gosched alerts ack <id>` still work but are deprecated and hidden from `--help`. Alerts were folded into the unified Activity view; use [`logs`](#logs) instead. They will be removed in a future release.
