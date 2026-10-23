---
title: CLI reference
nav_order: 3
---

# `gosched` command reference

**Audience:** anyone using go-schedule from a terminal\
**Applies to:** current source; see the release-status note below\
**Source of truth:** `internal/cli/`, this document describes what the binary does. `specs/001-task-scheduler/contracts/cli.md` describes what it must do, and remains a specification artifact.

> **Release status:** S107 additions are unreleased: explicit local IPC selection, complete authoring/preview, structured notification policies, operational views and headless preferences. Existing features retain their prior release availability; see [the changelog](https://github.com/shruggietech/go-schedule/blob/main/CHANGELOG.md) and the installed binary's help.

S109's CLI contract corrections are also unreleased: parser/operand validation, independent webhook/daemon endpoints, complete failed-import status, contextual diagnostics, labeled lists and accurate credential/watcher responses. The descriptions below apply to current source.

`gosched` is a thin client. With no target flags it talks to `goschedd` over protected local IPC, so existing scripts remain local and unchanged. A deliberate `--profile` or complete explicit remote selection uses the same command surface through the authenticated HTTPS allowlist.

## Contents

Maintainers can evaluate the complete command inventory with the [unattended CLI coordinator](cli-qa.md), including human and JSON output, deliberate bad inputs, scoped host adapters and canonical after-action reports.

- [Conventions](#conventions)
- [Global flags](#global-flags)
- [Remote profiles](#remote-profiles)
- [Exit codes](#exit-codes)
- [`health`](#health)
- [`daemon`](#daemon)
- [`bundle`](#bundle)
- [`mcp`](#mcp)
- [`task`](#task)
- [`schedule`](#schedule)
- [`calendar`](#calendar)
- [`runtime`](#runtime)
- [`systems` and `search`](#systems-and-search)
- [`preferences`](#preferences)
- [`cron`](#cron)
- [`group`](#group)
- [`trigger`](#trigger)
- [`watcher`](#watcher)
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
| `--ipc-path PATH` | Explicit local socket or Windows named pipe, incompatible with remote target selection; invalid selectors never fall back (S107). |
| `--json` | Emit machine-readable JSON instead of the table or summary. Available on every command that produces output. |
| `--profile ID_OR_LABEL` | Use one saved remote profile for this invocation. Labels must be unambiguous. |
| `--daemon-endpoint`, `--daemon-id`, `--credential-id`, `--certificate-file` | Select one explicit remote target. All four are required together and cannot be combined with `--profile`. |
| `--endpoint` | Legacy alias for the global remote endpoint where a command has no local flag of that name. Notification channel add/update use this name for the webhook destination; use `--daemon-endpoint` for their remote daemon. Both global endpoint aliases together are refused. |
| `-v`, `--version` | Print the CLI version. This is the *client* version; see [`health`](#health) for the daemon's. |
| `-h`, `--help` | Help for any command or subcommand. |

Remote human-mode commands identify the profile, endpoint, and shortened daemon ID on stderr. JSON stdout remains valid command JSON without a diagnostic prefix. Bearer values and pairing phrases are never accepted as target flags.

For a disposable local target, use an absolute Unix socket such as `gosched --ipc-path /tmp/owned-gosched/daemon.sock health`, or a Windows pipe such as `gosched --ipc-path \\.\pipe\gosched-test health`. The socket/pipe must already belong to the daemon you intend to operate. Local preference/profile commands administer user persistence independently of daemon target selection.

## Remote profiles

Pair a CLI identity by piping or interactively entering the one-time phrase on standard input. The phrase and issued bearer value are not printed or stored in the profile document:

```sh
gosched profile pair production --address https://10.0.0.8:8443 --expected-daemon-id DAEMON_ID --pairing-id PAIRING_ID --trusted-certificate daemon.pem --client-name "Operations CLI" --capability operate
```

Administer saved profiles with `gosched profile list`, `gosched profile show PROFILE`, `gosched profile rename PROFILE LABEL`, and `gosched profile remove PROFILE --confirm PROFILE_ID`. List and show expose the certificate fingerprint and credential reference, never the certificate body or bearer value. Removal requires the exact stable profile ID and deletes the native credential before metadata.

Use a saved target for exactly one command:

```sh
gosched --profile production task list
gosched --profile PROFILE_ID --json runs list
```

## Exit codes

| Code | Meaning |
| --- | --- |
| `0` | Success. |
| `1` | Runtime failure, the daemon was unreachable, the operation was refused, or the OS denied the request. |
| `2` | Usage or validation failure, a missing required flag, a malformed timestamp, an invalid policy value. Server-side validation failures map here too, so bad input is reported the same way wherever it is caught. |

Results go to stdout; diagnostics and errors go to stderr. That split is what makes `gosched task list --json | ...` safe in a pipeline.

Unknown commands/flags, invalid typed flag values, missing or surplus operands and malformed local bundle documents return 2. File read/write, connection and OS failures return 1. Commands reject extra operands before contacting a daemon or invoking a service/GUI adapter. Human actor, audit, pairing, profile, channel and watcher lists include column headings; their JSON shapes remain unchanged.

## `health`

Check that the daemon is running and report its version.

```sh
gosched health
```

```text
daemon ok (version 0.6.0)
```

Worth knowing: the version printed here is the **daemon's**, which can differ from `gosched --version` after a partial upgrade. If you are filing a bug report, both are asked for, and that is why.

## `daemon`

Inspect the connected daemon's stable identity, compatibility facts, operating mode, and capabilities:

```sh
gosched daemon manifest
gosched daemon manifest --json
```

Change only its operator-facing label with `gosched daemon rename "Workshop scheduler"`. Names are trimmed, contain 1 through 80 Unicode characters, and cannot contain control characters.

After copying a daemon database for independent concurrent use, replace the clone's copied identity with `gosched daemon reset-identity --confirm <current-installation-id>`. The exact current identifier is required. Reset preserves the display name and every scheduler record. See [Daemon identity](daemon-identity.md) for clean-install, upgrade, restore, clone, and reset semantics.

## `bundle`

Move portable automation intent to a deliberately selected daemon without exporting credentials, notification endpoints, enabled state, or watcher paths. A bundle operation first reads that daemon's manifest. If bundle support or a feature required by the document is missing, validation reports each missing capability and preview or apply stops before forwarding the request. Watcher bundles require a Windows, Linux, or macOS target; watcher paths are supplied separately for the target to interpret.

```sh
gosched bundle export --output automation.json
gosched --profile destination bundle validate automation.json
gosched --profile destination bundle compare automation.json
gosched --profile destination bundle preview automation.json --watcher-path WATCHER_ID=ABSOLUTE_TARGET_PATH --output plan.json
gosched --profile destination bundle apply plan.json
```

Inspect validation findings, comparison drift, and the complete preview plan before applying. `compare` is read-only. A preview plan is single-use and bound to the selected daemon and its current state; applying after switching targets or after the plan becomes stale requires a fresh preview. Target-only objects are reported as drift, never inferred as removals. Version 1 and version 2 bundles remain readable, and the daemon performs final path and domain validation.

Each `--watcher-path` must name a watcher portable ID present in the selected bundle, and each ID can be supplied only once, even if repeated paths are identical. Unknown or duplicate bindings return validation failure before a usable plan is created. Use the bundle's watcher `portable_id`, rather than a task/group ID or a daemon-local watcher ID.

## Webhook destinations and daemon selection

Channel add/update use their local `--endpoint` for the webhook destination. It can be combined with local IPC, a saved profile or the unambiguous explicit daemon selector:

```sh
gosched --ipc-path /tmp/owned-gosched/daemon.sock notification channel add Operations --endpoint https://hooks.example.test/events
gosched --profile production notification channel update CHANNEL_ID --endpoint https://hooks.example.test/events
gosched --daemon-endpoint https://daemon.example.test:8443 --daemon-id DAEMON_ID --credential-id CREDENTIAL_ID --certificate-file daemon.pem notification channel add Operations --endpoint https://hooks.example.test/events
```

Cobra resolves the final command's flags regardless of whether they appear before or after its name. On these two channel commands, `--endpoint` therefore always means the destination. Use `--daemon-endpoint` to select the daemon instead of relying on argument placement.

## Actors and management audit

Use `gosched actor list`, `create`, `update`, and `revoke` to administer credential-independent client identities. Use `gosched audit list` and `export` to inspect deterministic, redacted management evidence. See [Actor Permissions and Management Audit](access-control.md) for the capability matrix, lifecycle rules, retention limits, filters, and examples.

`gosched pairing list` identifies enrollment sessions and remaining attempts; `pairing cancel PAIRING_ID` prevents that phrase from enrolling another client. `credential list` lists remote credential records without tokens. `credential rotate CREDENTIAL_ID` returns the replacement token once, with its actual actor and daemon association, and invalidates the previous token. Store the replacement before leaving the response; saved profiles are not automatically updated. `credential revoke CREDENTIAL_ID` invalidates that credential without deleting its actor or other credentials.

## `mcp`

Serve the scheduler's observe-only Model Context Protocol resources over standard input and standard output for a local host:

```sh
gosched mcp serve
```

This command is meant to be launched by Codex or another MCP host. It opens no network listener, advertises no mutation tools, exits when the host disconnects, and reaches the daemon only through the existing protected local IPC endpoint. Protocol messages are the only stdout content; diagnostics use stderr. See [Local MCP access](mcp.md) for setup, resources, bounds, and the untrusted-content boundary.

`gosched mcp http enable --port PORT [--name NAME] [--origin ORIGIN...]` creates one runtime-only named localhost client and prints its bearer credential once. Omitted names retain the compatibility label `Local MCP client`; the desktop requires an explicit name. Status adds only non-secret identity and aggregate successful-access evidence. Rotation preserves the name and listener policy while clearing evidence, and disable revokes the client by closing the listener.

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
| `--stdin`, `--stdin-file` | Task input as text or UTF-8 file; `--stdin-file -` reads this command's standard input (S107). |
| `--run-as` | Execution identity supported by the selected daemon's operating system (S107). Unsupported identities return its existing validation error. |
| `--disabled` | Create disabled in the creation request, without an enabled interval (S107). |
| `--syntax` | `auto` (default), `human`, or `cron`; explicit selection requires a nonempty `--schedule` (S107). |
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

Modify a task. Only the fields you pass change; everything else is left alone. The flags are those of `task add`, except disabled creation is not an edit operation.

Pass `--name ""` or `--command ""` to clear that value, and pass `--clear-schedule` to remove automatic timing. Clearing command or the final automatic source disables the task atomically.

- `--arg` and `--env` **replace** the existing set rather than appending to it. Pass the full list you want.
- `--group` is three-way. Omit it and group membership is untouched; pass a group ID to move the task; pass an empty string (`--group ""`) to remove the task from its group.

```sh
gosched task edit 6f1c… --schedule "every weekday at 07:00"
```

At most one of `--schedule` or `--at` may be given.

S107 adds `--clear-args`, `--clear-env`, `--clear-cwd`, `--clear-run-as` and `--clear-stdin`. A clear flag cannot be combined with its replacement flag. Clearing environment overrides preserves the inherited base process environment. `--arg ""` remains one empty argument; `--clear-args` removes every argument. `--stdin ""`, `--cwd ""` and `--run-as ""` express their corresponding empty/default intent. `task show --json` includes the persisted execution inputs.

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

`--parent` takes a group ID; omit it for a top-level group. S107 adds `--disabled` for atomic inactive creation.

### `group edit <id>`

S107 adds `--name` and `--parent` for direct editing. Use `--parent=` to move to the root, and omit a field to retain it. A group-cycle or missing-parent rejection leaves the entire combined update unchanged.

```sh
gosched group edit GROUP_ID --name "Nightly jobs" --parent PARENT_ID
gosched group edit GROUP_ID --parent=
```

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

The `SOURCE TASK` and `SOURCE RUN` columns identify the upstream completion for a chained execution. `SOURCE TRIGGER` identifies an external trigger without exposing its key. `SOURCE WATCHER` identifies a filesystem watcher without recording its matched path. These columns are `-` when they do not apply. The `EXIT` column is the process exit code, or `-` where there isn't one, a run that never started has no exit code, and printing `0` for it would be a lie.

S107 adds `runs active`, `runs get RUN_ID` and `runs --details`. Detail includes run/task IDs, lifecycle timestamps, captured output/truncation, start failure and source task/run/trigger/watcher correlations. Default history text and JSON arrays retain their formats. `--limit 0` preserves the existing daemon-default query behavior.

## `trigger`


External triggers let a local process request one run through the normal scheduler by presenting one opaque key.

```sh
gosched trigger create --name "Build hook" --task <task-id>
gosched trigger list
gosched trigger show <trigger-id>
gosched trigger show <trigger-id> --reveal-key
gosched trigger update <trigger-id> --name "Release hook" --task <task-id>
gosched trigger disable <trigger-id>
gosched trigger enable <trigger-id>
gosched trigger rotate <trigger-id>
gosched trigger fire <key>
gosched trigger rm <trigger-id>
```

Create defaults to enabled; pass `--disabled` to create a dormant trigger. Create, rotate, and explicit `show --reveal-key` are the only commands that display a raw key. Rotation invalidates the old key immediately. Firing respects the target task's current lifecycle, enabled state, ancestor groups, worker limit, and overlap policy.

Trigger Sets create and administer 1 through 99 ordinary trigger members that share one target:

```text
gosched trigger set create --name "Build agents" --task <task-id> --count 10
gosched trigger set list
gosched trigger set show <set-id>
gosched trigger set reveal <set-id>
gosched trigger set retarget <set-id> --task <task-id>
gosched trigger set disable <set-id>
gosched trigger set enable <set-id>
gosched trigger set rotate <set-id>
gosched trigger set rm <set-id>
```

Create, reveal, and rotate print exactly one complete `gosched trigger fire <key>` command per member in permanent position order with one final newline. Pass `--json` for structured identities, positions, keys, and commands. Ordinary list and show output remains redacted. Set-level mutations are atomic, while individual member operations continue to affect only the selected member.

## `watcher`

Filesystem watchers request one normal scheduler run after matching file activity settles. The complete lifecycle is available in human-readable and JSON modes:

```text
gosched watcher create --name "Incoming JSON" --kind directory --path ./incoming --pattern "*.json" --recursive --debounce 250ms --stability 500ms --task <task-id>
gosched watcher list
gosched watcher show <watcher-id>
gosched watcher update <watcher-id> --path ./other --pattern "*.json" --stability 1s
gosched watcher disable <watcher-id>
gosched watcher enable <watcher-id>
gosched watcher rm <watcher-id>
```

Use `--kind file` for one exact file or `--kind directory` for basename glob matching. Missing paths can be configured and appear as degraded until observation succeeds. Recursive directory observation excludes symbolic links and junctions. Network filesystems are best effort. The default debounce is 250 milliseconds and the default stability interval is 500 milliseconds; both accept values from 25 milliseconds through one hour. Changes while the daemon is stopped or degraded are not replayed.

Detail output separates `Saved enabled`, `Last observed health`, `Last observed at`, `Reconciliation` and `Target readiness`. An immediate enable/disable reply may still contain the previous health observation while reconciliation is pending; run `gosched watcher show WATCHER_ID` again to inspect the later observation. `current` means the observation matches the saved definition; `unknown` means the daemon did not supply revision metadata or its runtime is unavailable. JSON exposes the same distinction through additive definition/observation revisions and `reconciliation_pending`; existing health fields retain their meaning.

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

## `notification`

Manage reusable webhook or SMTP channels, task and group policies, and redacted delivery evidence. Add `--json` to every read or mutation that returns an object.

```sh
gosched notification channel add "Operations" --endpoint https://receiver.example/hook --authorization "Bearer value" --health-interval 5m
gosched notification channel list
gosched notification channel get <channel-id>
gosched notification channel update <channel-id> --name "Primary operations" --endpoint https://new.example/hook
gosched notification channel rotate <channel-id> --authorization "Bearer replacement"
gosched notification channel test <channel-id>
gosched notification channel disable <channel-id>
gosched notification channel enable <channel-id>
gosched notification task set <task-id> --channel <channel-id> --on failure
gosched notification task set <task-id> --channel <channel-id> --on failure --failure-threshold 3 --failure-to-start --duration-threshold 30m --recovery --reminder 2h
gosched notification group set <group-id> --channel <channel-id> --on success --quiet-period 24h
gosched notification task show <task-id>
gosched notification task effective <task-id>
gosched notification group set <group-id> --channel <channel-id> --on success,failure
gosched notification group show <group-id>
gosched notification deliveries --state failed --limit 20
gosched notification channel rm <channel-id>
```

`--authorization` is write-only and may be visible in shell history, so use the shell's protected input practices when that matters. An empty rotation value clears authorization. `--health-interval 0` disables healthy-presence heartbeats. Condition durations accept Go duration syntax such as `30m` or `2h`; zero disables the condition. Task assignments replace inherited group assignments; otherwise the nearest group with assignments wins. Omitting every `--channel` from `task set` or `group set` clears the direct policy and resumes inheritance. The complete condition, suppression, restart, heartbeat, and receiver contracts are in [Webhook notifications](notifications.md).

S107 adds structured per-destination replacement. Export a canonical request, edit its assignments and reapply it:

```sh
gosched notification task show TASK_ID --input > policy.json
gosched notification task set TASK_ID --file policy.json --json
gosched notification group set GROUP_ID --file -
```

`--file -` reads standard input. Input is one UTF-8 JSON object of at most 1 MiB, with a required non-null `assignments` array. Unknown fields, trailing values and conflicts with explicitly supplied shorthand flags are rejected. Every replacement is validated before persistence; an invalid/duplicate destination cannot partially replace the existing policy.

```text
{
  "assignments": [
    {
      "channel_id": "SUCCESS_CHANNEL_ID",
      "on_success": true,
      "on_failure": false
    },
    {
      "channel_id": "PROBLEM_CHANNEL_ID",
      "on_success": false,
      "on_failure": true,
      "failure_threshold": 3,
      "on_failure_to_start": true,
      "duration_threshold_seconds": 600,
      "on_recovery": true,
      "reminder_interval_seconds": 3600,
      "quiet_period_seconds": 600
    }
  ]
}
```

Structured duration fields are integer seconds. `{"assignments":[]}` removes the direct policy and resumes inheritance. Canonical export omits database IDs, scope metadata and timestamps while preserving each supported destination condition. Default `show` output remains compatible. `notification task effective` identifies the selected direct/inherited source and each destination's effective conditions.

The Activity view identifies itself as a limited recent view and displays the exact configured path to the daemon's complete rotating JSONL log. Platform install guides list the default locations, but `log_file_path` overrides them; the path reported in Activity is authoritative for the running daemon.

`logs --metadata` (S107) includes the selected daemon's log path; JSON becomes a `logs`/`log_path` envelope. Default JSON remains a bare record array. `logs --limit 0` retains every record available in the bounded ring.

## `service`

Manage the system-wide background service, so the scheduler starts on boot and runs whether or not anyone is logged in.

| Subcommand | Effect | Elevation |
| --- | --- | --- |
| `install [--config FILE]` | Register the daemon with the system service manager, optionally retaining one validated configuration path. | **Required** |
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

A flag-free installation loads `config.json` from the platform data directory when that file exists and otherwise uses safe built-in defaults. Use `--config` when the service must retain another file. The command verifies that the file exists and that its complete schema is valid, resolves it to an absolute path, and records exactly `--config <absolute-path>` in the service definition before registration. Relative data, IPC, log, certificate, and private-key paths inside a file-backed configuration resolve from that configuration file's directory, not from the service manager's working directory.

```sh
sudo gosched service install --config /etc/goschedule/daemon.json
```

The configuration contains paths and policy, not pairing phrases or bearer credentials. Ensure the service identity can read the configuration, certificate, and private key after installation. A missing or invalid retained file makes daemon startup fail closed; it never falls back to a network-enabled or partially parsed default.

`status` is deliberately the one subcommand an ordinary user can run. It asks the operating system for no more access than a read needs, so it answers for an unprivileged caller wherever the service's own permissions allow a status query which, for a service installed by go-schedule, they do. Before 0.6.0 it requested start and stop rights it never used and failed with `Access is denied` for anyone not elevated, which reported that permission was withheld when in fact it was granted.

The other five genuinely change system state and genuinely require elevation. That is not being relaxed.

## `gui`

```sh
gosched gui
```

Launches the Wails desktop application and detaches. On Windows no console window appears, which is why launching it this way is preferable to running the GUI binary from a shell.

The GUI must be present next to the `gosched` binary. If it is not, a server-only install, for instance, the command says so and names the path it looked in.

## Deprecated: `alerts`

`gosched alerts` and `gosched alerts ack <id>` still work but are deprecated and hidden from `--help`. Alerts were folded into the unified Activity view; use [`logs`](#logs) instead. They will be removed in a future release.

## Remote pairing and credentials

`gosched pairing create <display-name>` creates a ten-minute one-time phrase through protected local IPC. Select `--kind desktop|cli|json|mcp`, `--capability observe|operate|manage|enroll`, and `--grant-duration 1h|24h|7d|30d|non-expiring`. The duration controls the resulting actor after exchange, while the separate ten-minute enrollment deadline limits only the phrase. The command displays the phrase once together with the pairing ID, daemon ID, enrollment expiration, and intended grant expiration; do not place those values in scripts or ordinary shell history.

`gosched pairing list` shows metadata without phrases, and `gosched pairing cancel <pairing-id>` invalidates an active request. `gosched credential list` shows safe fingerprints and lifecycle metadata. `gosched credential rotate <credential-id>` displays one replacement token and immediately invalidates the old value. `gosched credential revoke <credential-id>` revokes both the credential and its actor relationship. Named CLI profiles use `gosched profile pair`, and one invocation selects a saved target with `gosched --profile <id-or-label> <command>`. See [Remote access](remote-access.md) for the complete operator sequence.

## `schedule`

S107 adds non-mutating full-policy preview for human-readable and supported cron expressions. It returns five predicted occurrences using the same evaluator as task creation.

```sh
gosched schedule preview "every month on the 31st at 09:00" --syntax human --tz America/New_York --missing-date last_valid --time-basis wall_clock --dst-gap skip --dst-overlap both
gosched schedule preview "0 9 * * 1-5" --syntax cron --tz UTC --json
```

`--syntax auto|human|cron` defaults to automatic interpretation. Preview accepts `--tz`, `--missing-date`, `--time-basis`, `--dst-gap` and `--dst-overlap` with the same meanings as task authoring. Text labels predictions, syntax, timezone and effective policy; JSON includes the next-run timestamps and policy summary. No task is created or edited.

## `calendar`

S107 adds a window of recorded executions and predicted occurrences. Text labels these separately and includes task/run identity; JSON retains the source occurrence fields.

```sh
gosched calendar
gosched calendar --from 2026-10-19T00:00:00Z --to 2026-10-26T00:00:00Z --json
```

The default interval starts now and ends seven days later. Explicit timestamps require RFC 3339 with a timezone. The end must follow the start by no more than 30 days. Predictions describe intended scheduling; recorded rows describe runs that exist.

## `runtime`

S107 adds local inspection of the connected daemon's actual data, database, configuration, log and lock paths, IPC endpoint and process ID.

```sh
gosched runtime
gosched --ipc-path /tmp/owned-gosched/daemon.sock --json runtime
```

Runtime path inspection is a local API capability. A selected remote target returns explicit unavailable guidance; the command never substitutes paths from this computer for remote data.

## `systems` and `search`

S107 adds identified operational summaries and search. Commands use the selected target by default; `--all` explicitly includes the local target and every saved registration.

```sh
gosched systems summary
gosched --profile production systems summary --json
gosched systems summary --all --json
gosched search "nightly" --kind task,group --limit 20
gosched systems search "failure" --all --kind failure,alert --json
```

Search kinds are `task`, `group`, `failure`, `schedule` and `alert`. Query text contains 1 through 200 characters, and `--limit` allows 1 through 50 results per scheduler. Results identify the registration, endpoint, expected/observed daemon identity and observation time. Summaries include active-task, recent-failure, unacknowledged-alert and notification-problem information.

At most four targets are observed concurrently, each with a five-second deadline inside the command's ten-second request deadline. Saved credentials resolve only when that target's worker runs. TLS, identity and current Observe authority are checked before feature reads.

Native credential setup participates in operational request cancellation through hidden owned helper processes. Secrets remain in the native store and peer-authenticated private local channels. Helper requests and responses never use CLI standard streams. Request budgets exclude the separate bounded cleanup drain (at most five seconds); a cleanup failure remains an error and does not become a successful target observation.

One failed target does not discard another target's successful response. JSON contains independent observations and failures; partial failure returns nonzero after writing usable results. `--all` cannot be combined with remote target flags. An explicit `--ipc-path` can select the local fixture included in an all-system query.

## `preferences`

S107 adds headless read, partial update and explicitly scoped reset of the desktop's durable appearance and popup choices.

```sh
gosched preferences show --json
gosched preferences set --appearance dark --popup-enabled=true
gosched preferences set --popup-conditions failure,alert --popup-severities warning,error
gosched preferences set --popup-daemons DAEMON_ID
gosched preferences set --popup-daemons=
gosched preferences reset --scope appearance
gosched preferences reset --scope popups
gosched preferences reset --scope all
```

Appearance accepts `system`, `light` or `dark`. Popup conditions accept `failure`, `success` and `alert`; severities accept `info`, `warning` and `error`. Empty source-daemon filters include all registered schedulers. Omitted fields remain unchanged; invalid input leaves the document unchanged. Condition/severity lists remain nonempty even while delivery is disabled.

The CLI and desktop share the version-1 schema, legacy migration, defaults, validation and atomic persistence. Concurrent writers coordinate scoped edits so an appearance change does not overwrite popup choices. Output identifies the user and file path; JSON contains `user`, `path` and `preferences`. Commands do not need a daemon, open the desktop, emit a popup or change daemon webhook policies. `reset --scope all` also repairs an unreadable preference document.

The saved desktop startup connection is separate:

```sh
gosched profile default show --json
gosched profile default set production
gosched profile default reset
```

The next desktop startup reads that choice; resetting selects this computer. These commands use existing registrations without reading native credentials or contacting a daemon. Per-invocation `--profile` and the ordinary CLI local default retain their existing meanings.

## SMTP channel administration (Unreleased S111)

The notification channel verbs also support `--kind smtp`. Supply `--smtp-host`, `--smtp-from`, `--smtp-to` and optional `--smtp-port`, `--smtp-tls starttls|tls|plain`, `--smtp-timeout 10s` and `--smtp-html`. Update merges explicitly supplied settings with the selected channel's current settings. Transport kind cannot be changed. SMTP settings flags do not select a daemon; use the existing profile or global `--daemon-endpoint` selector for that purpose.

Create and rotate accept `--smtp-auth-file PATH` or `--smtp-auth-file -` for bounded UTF-8 JSON from a file or standard input. The object is `{"username":"...","password":"..."}`; `{}` clears authentication during rotation. Values are never returned by get/list/history. Human history includes each recipient's state, numeric SMTP status and authored explanation. JSON keeps the same command response shape with additive SMTP settings and recipient fields.

See [Notifications](notifications.md#smtp-email-channels) for all settings, encryption choices, provider setup, retries and portable local-credential binding.
