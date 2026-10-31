---
title: Maintainer test scripts
nav_order: 7
---

# Maintainer Test Scripts

**Status:** current **Feature:** [`specs/006-maintainer-test-scripts/`](https://github.com/shruggietech/go-schedule/tree/main/specs/006-maintainer-test-scripts/) **Audience:** maintainers verifying an installed `goschedd`

Scripts that let you prove an installed scheduler actually fires when it said it would, and that it survives restarts, catches up after downtime, and honors its overlap policies. They are test payloads: you point a scheduled task at one, walk away, and afterwards read the evidence out of a local SQLite database.

## Contents

For a complete disposable CLI evaluation, use the [unattended Windows and WSL coordinator](cli-qa.md) and its versioned maintainer prompts. The installation payloads below serve their existing installed-scheduler purpose.

- [Why these exist](#why-these-exist)
- [Prerequisites](#prerequisites)
- [Quickstart](#quickstart)
- [The scripts](#the-scripts)
  - [Test-Heartbeat](#test-heartbeat)
  - [Test-GetSystemInfo](#test-getsysteminfo)
  - [Test-ReadTestDB](#test-readtestdb)
- [What gets recorded](#what-gets-recorded)
- [Recipes](#recipes)
- [How drift is derived](#how-drift-is-derived)
- [Exit codes](#exit-codes)
- [Shell conventions](#shell-conventions)
- [Troubleshooting](#troubleshooting)

## Why these exist

`gosched runs` tells you a task ran. It does not tell you how far from its scheduled moment it ran, nor that a firing you expected never happened at all. Those are the two questions that matter when you have just installed a release on a new machine, and neither was answerable before.

Every script comes as a matched pair, a PowerShell `.ps1` and a POSIX `.sh`, with identical behavior and one-to-one options. A PowerShell `-FooBar` is `--foo-bar` in the shell twin. There are no options on one side without a counterpart on the other.

## Prerequisites

- **`sqlite3` 3.33.0 or later.** Required. The floor is not arbitrary: the scripts use `.param set` (3.32.0) for bound parameters and `.mode json` (3.33.0) for machine-readable output.
- **PowerShell 7+** (`pwsh`) for the `.ps1` twins, including on Linux and macOS, where they work fine.
- **`bash`** for the `.sh` twins.
- A running `goschedd` for anything involving the scheduler.

If `sqlite3` is missing, the scripts exit **2** and tell you both ways to fix it:

```bash
pwsh -File test/scripts/Test-Heartbeat.ps1 -InstallSqlite
```

That downloads the pinned build from sqlite.org, verifies its SHA-256 against [`lib/sqlite-manifest.json`](https://github.com/shruggietech/go-schedule/blob/main/test/scripts/lib/sqlite-manifest.json), and unpacks it into `test/scripts/.bin/` (gitignored). It is the only thing in this feature that touches the network, and only when you pass the flag. A checksum mismatch deletes the download and fails, it never installs an unverified binary.

Or install it yourself:

```bash
winget install SQLite.SQLite
```

```bash
sudo apt install sqlite3
```

```bash
brew install sqlite
```

## Quickstart

```bash
pwsh -File test/scripts/Test-Heartbeat.ps1 -Label smoke
```

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Database heartbeat -Query summary
```

One beat recorded, one beat read back. That is the whole loop.

Databases live in a **user-writable** directory, never the daemon's system-wide data directory, test payloads must run unelevated, and disposable test output does not belong beside live scheduler state:

| Platform | Location |
|---|---|
| Windows | `%LOCALAPPDATA%\goschedule-test\` |
| Linux | `$XDG_DATA_HOME/goschedule-test/` or `~/.local/share/goschedule-test/` |
| macOS | `~/Library/Application Support/goschedule-test/` |

Override with `GOSCHEDULE_TEST_DIR`. Every script prints the path it resolved.

> **A task using `run_as` resolves a *different* per-user directory**, the one belonging to the user the task runs as. That is correct, and occasionally surprising. Pass an explicit `--database` path when `run_as` is in play.

## The scripts

### Test-Heartbeat

Records one beat per invocation and exits. **This is deliberate.** The point is to verify that `goschedd` supplies the cadence; a script that sleeps in a loop tests its own `sleep`, not the scheduler.

| PowerShell | POSIX | Default | Meaning |
|---|---|---|---|
| `-Database` / `-d` | `--database` / `-d` | `heartbeat` | Name or path |
| `-IntervalSeconds` / `-i` | `--interval-seconds` / `-i` | none | Declared schedule interval; enables gap detection, and drift when paired with an anchor |
| `-AnchorIso` / `-a` | `--anchor-iso` / `-a` | none | One real firing time (RFC 3339) from `gosched task show`; with the interval this yields true latency |
| `-Label` / `-l` | `--label` / `-l` | none | Tag recorded on each beat |
| `-Loop` / `-r` | `--loop` / `-r` | off | Opt-in bounded continuous mode |
| `-MaxBeats` / `-m` | `--max-beats` / `-m` | none | Stop after N beats |
| `-DurationSeconds` / `-t` | `--duration-seconds` / `-t` | 3600 under `-Loop` | Stop after N seconds |
| `-SleepSeconds` / `-s` | `--sleep-seconds` / `-s` | 0 | Deliberately extend the run |
| `-FailWith` / `-f` | `--fail-with` / `-f` | none | Exit non-zero after recording |

Notes that matter:

- **`-Loop` has no unbounded form.** With neither bound supplied the 3600-second default applies. A runaway loop launched under a scheduler is a resource incident, not a test.
- **The duration bound is checked between beats**, so one deliberately-slow run can overrun it by up to that run's length. Interrupting a run mid-write would corrupt the record the bound exists to protect.
- **`-FailWith` rejects 0 and 2.** Those are reserved for success and for unmet prerequisites, and an induced failure must never be mistakable for either.
- **The beat is written once, at the end of the run.** A run interrupted mid-flight therefore records nothing and shows up as a missed firing, which is the honest signal, since you cannot act differently on a run that vanished than on one that never started.

### Test-GetSystemInfo

Records a snapshot of the machine: identity, user, process count, uptime, network addresses, listening ports. As a scheduler test it exercises subprocess spawning, platform tooling, and multi-row writes, where cross-platform bugs actually surface.

| PowerShell | POSIX | Default | Meaning |
|---|---|---|---|
| `-Database` / `-d` | `--database` / `-d` | `system` | Name or path |
| `-InvocationSource` / `-i` | `--invocation-source` / `-i` | `manual` | Tag recorded on the snapshot |
| `-SkipPorts` / `-s` | `--skip-ports` / `-s` | off | Skip the port probe |

**Probes degrade; they do not abort.** A probe that cannot run on this host records `NULL`, warns on stderr, and the snapshot is still written with exit 0. `NULL` means *could not determine*; it is never used for a legitimate zero, because a process count of zero and an unavailable process count support opposite conclusions.

Each snapshot also records **why** its address and port lists look the way they do, in `addresses_probe` and `ports_probe`:

| Status | Meaning |
|---|---|
| `ok` | The probe ran. Zero rows is a real answer. |
| `unavailable` | No tool on this host could answer. Zero rows means nothing. |
| `skipped` | `-SkipPorts` was given. |

This exists because an empty port list is otherwise ambiguous, and the two readings support opposite conclusions. `listeners` relies on it (below).

`process_name` on ports is commonly `NULL`, most platforms want elevation to attribute a socket to a process. That is normal, not a defect.

### Test-ReadTestDB

| PowerShell | POSIX | Default | Meaning |
|---|---|---|---|
| `-Database` / `-d` | `--database` / `-d` | `heartbeat` | `heartbeat`, `system`, or a path |
| `-Query` / `-k` | `--query` / `-k` | `summary` | Which canned query |
| `-List` / `-n` | `--list` / `-n` | off | List queries and exit |
| `-Format` / `-f` | `--format` / `-f` | `Table` | `Table`, `Json`, `Csv` |
| `-Limit` / `-m` | `--limit` / `-m` | 20 | Row cap |
| `-IntervalSeconds` / `-i` | `--interval-seconds` / `-i` | inferred | Interval for gap, jitter, and drift |
| `-AnchorIso` / `-a` | `--anchor-iso` / `-a` | none | One real firing time from `gosched task show`; with the interval makes `drift` report true latency |

> `-Query` uses the alias `-k`, not `-Q`. PowerShell aliases are case-insensitive, so `-Q` would collide with `-Quiet`. `-Loop` uses `-r` for the same reason, `-L` would collide with `-Label`.

| Query | Database | Answers |
|---|---|---|
| `summary` | both | How many records, over what period, from how many sessions and hosts? |
| `recent` | both | What are the most recent records? |
| `cadence` | heartbeat | What were the observed intervals? |
| `drift` | heartbeat | True dispatch latency, by source. Needs an anchor. |
| `jitter` | heartbeat | Variation around the schedule's own observed phase. No anchor needed. |
| `gaps` | heartbeat | Which expected firings were missed or badly delayed? |
| `overlaps` | heartbeat | Which runs overlapped in time? |
| `failures` | heartbeat | Which runs reported failure? |
| `restarts` | heartbeat | Where are the session boundaries? |
| `hosts` | both | Which hosts and users produced records? |
| `listeners` | system | What is listening on the most recent snapshot whose port probe ran, and what changed since the previous comparable one? |
| `schema` | both | What is the stored structure? |

**Three reporting rules are contract, not presentation.** Breaking any of them would produce numbers that look more certain than they are:

1. A query that excludes rows says how many and why. A percentile over an unstated subset is a confident number drawn from unknown evidence.
2. Drift is never pooled across sources and never shown without its source. A measured value and a derived one are different kinds of number.
3. `gaps` says whether its interval was supplied or inferred.

## What gets recorded

The unreleased S128 helpers initialize the complete shared metrics format 4 directly: metadata, heartbeat, snapshot, address and port tables. Both writers may use the same current-format database. Historical formats 2/3, foreign and partial layouts are refused before conversion or metadata overwrite. Select a new empty path or explicitly reset disposable development state after stopping its writers. The historical [original data model](https://github.com/shruggietech/go-schedule/blob/main/specs/006-maintainer-test-scripts/data-model.md) records its earlier design.

`heartbeat.db` holds one `beat` row per completed run, carrying start and finish moments (both, so overlap is decidable rather than guessed), session and sequence, expected moment with its source and the resulting drift, exit code and outcome.

`system.db` holds a `snapshot` row per invocation with `snapshot_address` and `snapshot_port` children.

**Nothing is pruned, rotated, or size-capped.** These are disposable test artifacts, and deleting the file is the documented reset, automatic retention would silently destroy the history you are trying to inspect. For scale: a beat is roughly 200 bytes, so a one-minute schedule is about 525,000 rows and ~100 MB per year. A snapshot with children is ~2 KB, so hourly is ~18 MB per year.

## Recipes

### Does it fire on time?

```bash
gosched task add "hb-verify" --command pwsh --arg -File --arg test/scripts/Test-Heartbeat.ps1 --arg -IntervalSeconds --arg 60 --arg -AnchorIso --arg 2026-07-23T12:08:06Z --schedule "every 1 minute"
```

Wait ten minutes, then:

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Query drift -IntervalSeconds 60
```

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Query gaps -IntervalSeconds 60
```

Expect ~10 beats, `boundary` source, no gaps. Compare the drift against the project's documented p99 < 100 ms dispatch budget. This *measures and compares*; it does not certify.

### Does it survive a restart?

Startup-event lifecycle is covered without a manual install/restart session:

```bash
go test ./internal/engine -run Startup
go test ./internal/api/server ./internal/cron ./internal/store -run 'Startup|Reboot'
```

The controlled tests open the store, start two independent engines, inject the startup instant, exercise a reload inside one lifecycle, and check disabled task/group, history origin, import/export, and persistence behavior. Installed daemon smoke testing may still use the service restart below, but it is not a release gate and no physical Windows VM observation is required.

```bash
gosched service restart
```

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Query restarts
```

If a sampled completion timestamp precedes its start, the heartbeat still records one beat with finished time equal to start and duration zero. An unsuppressed warning includes both observed timestamps and asks the maintainer to check the system clock. That zero duration records the timing anomaly, not a reliable elapsed-time measurement; the intended task exit code and other observations are retained.

Each invocation is its own session, so under a scheduler the session boundary marks every beat, the evidence is the uninterrupted `started_ms` sequence across the restart.

### Does catch-up work?

Stop the daemon past a scheduled firing, start it again, then check `gaps`. Under a make-up-once policy expect exactly one make-up beat before the cadence resumes. More than one, or none, is a real finding.

### Do overlap policies hold?

```bash
gosched task add "hb-overlap" --command pwsh --arg -File --arg test/scripts/Test-Heartbeat.ps1 --arg -SleepSeconds --arg 90 --schedule "every 1 minute" --overlap queue_one
```

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Query overlaps
```

Expect: `queue_one` → no overlaps, serialized; `skip` → no overlaps and visibly fewer beats than intervals; `allow_concurrent` → overlapping ranges. Delete and recreate the task between policies.

### Does failure handling work?

```bash
gosched task add "hb-fail" --command pwsh --arg -File --arg test/scripts/Test-Heartbeat.ps1 --arg -FailWith --arg 1 --schedule "every 1 minute"
```

Then check the `failures` query *and* `gosched alerts`. The beat proves the record survived the failure; the alert proves the daemon saw it.

## How drift is derived

**Drift is derived, not reported by the scheduler.** The executor passes a spawned task the inherited environment plus the task's own configured variables and nothing else, no scheduled time, no run ID. So the expected moment comes from one of three sources, in strict precedence, and the source is recorded on every beat:

| Source | Meaning |
|---|---|
| `env` | `GOSCHED_SCHEDULED_TIME` from the environment. Not set today; the tier exists so a future release is consumed with no change here. |
| `anchor` | `-AnchorIso` plus `-IntervalSeconds`. **This is the working path.** |
| `none` | Neither available. No drift recorded, and the reader excludes the beat. |

### Getting the anchor, pass it to the *reader*

Any single real firing time works: they all sit on the same `anchor + k × interval` grid.

**The anchor cannot be known before the task exists.** This scheduler derives an interval schedule's phase from the task's creation moment, so you only learn the anchor by creating the task and then asking. Passing it to the recorder is a chicken-and-egg problem, that is why drift is derived at *read* time:

```bash
gosched task add "hb-verify" --command pwsh --arg -File --arg test/scripts/Test-Heartbeat.ps1 --arg -IntervalSeconds --arg 60 --schedule "every 1 minute"
```

```bash
gosched task show <task-id>
```

Take any value from `next runs`, let the task run a while, then:

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Query drift -IntervalSeconds 60 -AnchorIso 2026-07-23T12:35:11Z
```

Now `drift` is genuine dispatch latency, how far each run's start landed from the moment the schedule said it would fire. Because it is computed from the raw start timestamps, it works on beats **already recorded**, and you can re-run it with a corrected anchor if you got it wrong the first time.

A measured example from this machine: `261–312 ms` mean `286 ms`, against the `6505 ms` the pre-0.5.1 code reported for the same kind of schedule.

`Test-Heartbeat` also accepts `-AnchorIso`, which records the drift at write time. Use it only when you genuinely know the firing grid in advance, a fixed-time schedule, say. For interval schedules, prefer the reader.

### Why there is no epoch-boundary tier

Versions before 0.5.1 accepted an interval alone and snapped the run's start to the nearest multiple of that interval **counted from the Unix epoch**. That is correct only if the schedule happens to sit on the epoch grid, and this scheduler anchors an interval schedule to the **task's creation time**. A task created at `:06` fires at `:06` forever.

The result was a constant phase offset reported as though it were lateness. A measured example: a task firing at `:06` produced drift of 6505, 6262, 6254 ms, apparently 64× over the project's 100 ms dispatch budget. The same run's `cadence` query showed intervals of 59757–60006 ms, meaning the scheduler was on time to within a quarter of a second. The 6.4 s was entirely the `:06` anchor.

That tier is gone. Without an anchor, no drift is recorded at all. **Reporting nothing is better than reporting a confident wrong number**, because nothing in the presentation of the wrong number tells you which one you got.

### When you have no anchor: `jitter`

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Query jitter -IntervalSeconds 60
```

`jitter` derives the schedule's phase from the data itself and reports the variation around it, `phase_ms` (the offset it found) plus the min, max, and spread around that offset.

**It cannot detect uniform lateness.** A scheduler consistently five seconds late has a perfectly stable phase and therefore zero jitter. Jitter bounds the *variability* of dispatch; only an anchor gives you the absolute figure. The reader states this on every `jitter` run rather than leaving you to remember it.

Current drift queries use the current heartbeat source and timing fields. Historical schema and nullable probe fallbacks are retired; no helper converts older recordings.

## Exit codes

| Code | Meaning |
|---|---|
| `0` | Success |
| `1` | Runtime failure, a required probe failed, a write failed, contention exhausted its retries |
| `2` | Usage error or unmet prerequisite, bad arguments, `sqlite3` absent or too old, unsupported platform |

The distinction between 1 and 2 is load-bearing: it is what makes a run row in `gosched runs` mean something. An unmet prerequisite is a usage-class failure, conflating it with a runtime failure sends you debugging the wrong thing.

Results go to stdout, diagnostics and warnings to stderr, so reader output pipes cleanly into another tool without log lines mixed in.

## Shell conventions

The `.ps1` twins follow the ShruggieTech PowerShell standard, enforced by [`.claude/skills/shruggie-powershell/scripts/Test-ScriptCompliance.ps1`](https://github.com/shruggietech/go-schedule/blob/main/.claude/skills/shruggie-powershell/scripts/Test-ScriptCompliance.ps1). There is **no equivalent shell standard yet**, so the `.sh` twins are governed by this section plus `shellcheck`:

- `#!/usr/bin/env bash` and `set -euo pipefail`.
- A header comment block covering purpose, the twin relationship, the exit-code contract, and a pointer here.
- Options parsed in a `while`/`case` loop; every PowerShell parameter has a `--kebab-case` counterpart and the same single-letter short option.
- `log LEVEL MESSAGE` for all diagnostics, to stderr. Results to stdout.
- `die_usage` exits 2; `die_runtime` exits 1. No bare `exit 1` for a usage error.
- All shared logic lives in `lib/sqlite.sh`, one implementation per twin. Three copies of the resolution order would be three chances for them to disagree, and disagreement there presents as an intermittent platform-specific failure.
- Quote every expansion. Values written to the database include hostnames and usernames, which can contain quotes and spaces.
- **Never build SQL by string interpolation.** Use `sqlite_exec`'s bound parameters. Interpolating a username into SQL is both an injection vector and an ordinary bug for anyone named `O'Brien`.

Both SQLite adapters bind non-null values as literal UTF-8 text through a CLI-quoted `CAST(X'...' AS TEXT)` expression. Numeric-looking session IDs, leading zeroes, quotes, newlines and non-ASCII text retain their exact values; explicit null values (the existing Bash `__NULL__` marker and PowerShell `$null`) remain SQL `NULL`. Bash reserves `__NULL__` for that purpose. The SQL statement performs any intended numeric conversion. This prevents distinct hexadecimal session IDs that resemble scientific notation from being evaluated as the same numeric value.

## Troubleshooting

**`exit 2` with "No usable sqlite3 found"**, none of the three search locations held a usable 3.33.0+. Use `--install-sqlite` or your package manager.

**`--sqlite-exe` exits 2 even though the file exists**, on Windows, pass a native path (`C:\bin\sqlite3.exe`), not a Git Bash path (`/c/bin/sqlite3`). PowerShell cannot resolve the latter. An explicit path that cannot be used is a hard error by design: silently falling back to a different `sqlite3` than the one you named is how you debug the wrong binary for an hour.

**"address probe unavailable" on Windows under Git Bash**, expected. Neither `ip` nor `ifconfig` exists there. Use the PowerShell twin, which uses `Get-NetIPAddress`. The snapshot is still recorded.

**Drift looks enormous**, check the anchor. `-AnchorIso` must be a real firing time of *that* schedule, and `-IntervalSeconds` must match it. A wrong anchor shifts the whole grid. If you are seeing a large, suspiciously *constant* drift, that is the signature of a phase error rather than a late scheduler, compare against `-Query cadence`, which is anchor-independent.

**`listeners` shows an older snapshot than expected**, by design. It reads the most recent snapshot whose `ports_probe` is `ok`, and says which one it picked. If your newest snapshot was taken with `-SkipPorts`, or on a host with no `ss` or `netstat`, that snapshot has no port data and showing it would present an empty list as though nothing were listening. The warning names the newest snapshot and its probe status so the substitution is never silent.

**`listeners` says "no snapshot with usable port data"**, every snapshot so far recorded `ports_probe` of `unavailable` or `skipped`. Re-run `Test-GetSystemInfo` without `-SkipPorts` on a host that has a port tool. The message is deliberate: an empty table here would read as a finding rather than an absence of data.

**`drift` returns no rows**, every beat had `expected_source = 'none'`, meaning no anchor was supplied at record time. The excluded-count warning says how many. Pass `-AnchorIso` *and* `-IntervalSeconds` when recording, not when reading. If the beats are already recorded and you cannot re-run, use `-Query jitter`.

**Contention errors under `allow_concurrent`**, the scripts use WAL with a 5-second busy timeout and at most three SQLite attempts. A diagnostic containing `locked` or `busy` takes the bounded retry path. Exhaustion reports that the record was not written and retains the final SQLite diagnostic. Other failures, such as a read-only database or a constraint violation, exit 1 immediately with the underlying error instead of being retried as contention.

The concurrent-writer integration check retains each failed child's actual output alongside its exit status and persisted-row count. Preserve that complete failure when investigating missing records. The fixture also records interpreter/runtime identity, elapsed time, invocation phase and decimal/hexadecimal native exit on failure. A later passing run does not establish the cause of an earlier failure whose diagnostic was not captured. SQLite candidate probes retain bounded stderr and distinguish a launch failure, nonzero native exit, ten-second probe timeout, malformed version and an old version. Both twins require a numeric three-part version. An unusable implicit candidate allows the documented search to continue; an explicit selection is never silently replaced.

## Automated tests

### Isolated CLI operation

Introduced in v1.6.0.

The reusable isolation and focused command adapters support the complete Windows/WSL unattended evaluator and fixed report contract described in the [maintainer CLI QA guide](cli-qa.md). Operator commands are documented in the [CLI reference](cli.md).

The root `internal/clifixture` package creates a UUID workspace with separate disposable runtime and retained evidence directories. Its guard checks installation identity, exact daemon PID, selected IPC and effective configuration/database/log/lock paths before a mutating packaged-CLI operation. It rejects changed ownership and symlink/junction escapes. Caller-selected watcher/task/export paths must also pass containment checks. Cleanup targets only the owned runtime path and exact process capability; evidence survives.

Real binaries connect using `--ipc-path` to the owned daemon. Windows children are hidden, noninteractive and suspended until assigned to their kill-on-close Job Object, inheriting only redirected stdin/stdout/stderr handles. Unix children use an owned process group. Cleanup waits for the exact process tree before removing runtime state and never terminates by executable name.

`cli.ExecuteWithDependencies` runs the actual handlers with scoped profile/preferences, secret, service, launcher, clock and I/O dependencies. Recording service/GUI/credential adapters exercise OS-facing success/failure semantics without calling the host manager, opening a window or touching a native keyring. Production defaults keep their existing native behavior; no environment-variable credential downgrade is introduced. Help/completion needs no daemon or remote credential, and cancellation reaches request contexts and explicit owned MCP streams.

Evidence kinds are `packaged_cli`, `command_harness` and `deterministic_engine`. They identify the execution boundary and do not substitute one kind of evidence for another. The focused qualification entry is:

```sh
go test ./internal/cli ./internal/userprefs ./internal/clientobserve ./internal/clifixture ./test/integration -run 'TestInvocation|TestExplicitIPC|TestRecording|TestGuard|TestOwnedProcessTree|TestPackagedCLIUsesOnlyOwnedDaemon' -count=1
```

Run it through the repository's verified hidden launcher on Windows. The normal aggregate remains `sh scripts/verify.sh all`; the packaged fixture is included in its root test coverage. The full scenario catalog, interrupted-run reporting and agent capture/review/finalize routine use the coordinator described in the maintainer CLI QA guide.

### Existing maintainer script contracts

```bash
go test ./test/integration/ -run TestScripts -v
```

Covers single-shot recording, boundary drift, the exit-code contract, the bounded loop, concurrent writers, exact literal-text binding, distinct numeric-looking session IDs, bounded retry and failure diagnostics, twin parity, snapshots, and every canned query. It **skips with a stated reason** when `sqlite3`, `pwsh`, or `bash` is missing, a skip is not a pass, and the reason is printed so the two cannot be confused.

### Concurrent-writer investigation

For the maintainer integration fixture only, set `GOSCHEDULE_PROBE_DIAGNOSTICS=1` when running `go test -race ./test/integration -run TestScriptsBindingsKeepConcurrentSessionIdentities -v`. The fixture emits child-owned interpreter, library, prerequisite and write phases; on failure, the test log retains all six native outcomes and exact owned rows. Default writer output is unchanged. This mode does not retry, serialize or alter the exact-record assertions. Original unknown-cause Windows observations remain tracked in [#355](https://github.com/shruggietech/go-schedule/issues/355).

### Outer process collection

The binding fixture starts its hidden child in the existing owned Job Object or process group, closes stdin, allows 120 seconds for execution and two seconds to drain inherited streams after root completion. Each stream retains at most one MiB. It retains the root native exit, collection failure and owned cleanup failure separately; exit 0 with incomplete collection fails the fixture. A descendant cannot keep the fixture waiting indefinitely after the root exits.

With GOSCHEDULE_PROBE_DIAGNOSTICS=1, outer provenance adds resolved executable and harness SHA-256 values, exact root PID, elapsed time and independent collection/cleanup status. It does not print environment values or harness contents. Existing child phase markers remain authoritative. Use TestScriptsRuntimeAndSelectionContrasts to compare metadata initialization with an actual writer and implicit discovery with explicit selection of the same SQLite executable. Passing contrasts do not explain a historical CLR crash or prerequisite refusal. Preserve the original failed log, then collect evidence from the exact owned PID if the crash recurs; do not replace the runtime or register a system-wide debugger as part of this routine.

## Opt-in owned crash incidents

For a focused integration experiment, set GOSCHEDULE_INCIDENT_DIR to an existing private absolute directory and GOSCHEDULE_PROBE_DIAGNOSTICS=1 before running the script binding/contrast tests. Each invocation creates an exclusive incident directory. Ordinary test runs leave native crash collection disabled, except the explicitly named controlled-crash regression. No installed interpreter, service, debugger registration or account session changes.

The owned PowerShell child enables CoreCLR Mini collection through child-only environment settings. Its original harness records PID, PowerShell/framework version and loaded runtime module hash before payload execution, then runs its body without an outer exit-code wrapper. Pre-payload identity remains distinct from crash-time evidence. A process-local error mode suppresses Windows failure dialogs in this diagnostic child. The hidden, noninteractive process-tree launcher contains the runtime collector child.

incident.json version 1 retains original native exit separately from collection status and artifact names, lengths and SHA-256 hashes. Native Mini dumps retain at most 32 MiB; runtime identity and collector logs each 64 KiB. Oversized generated files are discarded after the owned producer stops, with an explicit oversized observation. This bounds retained evidence, not peak temporary disk usage. Existing 120-second execution and two-second stream-drain limits remain. Windows dumps require a valid minidump signature; no dump after a nonzero exit is absent, not collected. Collection failures do not replace native exit, stream or cleanup outcomes.

Keep raw dumps and stacks private. Even a small dump can contain process memory. Use a native-capable debugger for native CLR faults; dotnet-dump does not display native stack frames. References: [CoreCLR dump generation](https://github.com/dotnet/runtime/blob/main/docs/design/coreclr/botr/xplat-minidump-generation.md), [Microsoft crash configuration](https://learn.microsoft.com/en-us/dotnet/core/diagnostics/collect-dumps-crash) and [dotnet-dump](https://learn.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-dump). A successful controlled crash demonstrates collection capability and does not explain the original intermittent #355 failure.

For the original Windows runtime/prerequisite and IPC investigation, use the [S122 runtime and IPC investigation recipe](https://github.com/shruggietech/go-schedule/blob/main/docs/maintainers/windows-incident-replay.md). It supplies a finite condition matrix, shared deadline, process/source attribution and fixed causal report while retaining the original outcomes and complete QA contract.
