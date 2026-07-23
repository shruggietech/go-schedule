# Maintainer Test Scripts

**Status:** current
**Feature:** [`specs/006-maintainer-test-scripts/`](../specs/006-maintainer-test-scripts/)
**Audience:** maintainers verifying an installed `goschedd`

Scripts that let you prove an installed scheduler actually fires when it said it
would — and that it survives restarts, catches up after downtime, and honors its
overlap policies. They are test payloads: you point a scheduled task at one, walk
away, and afterwards read the evidence out of a local SQLite database.

## Contents

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

`gosched runs` tells you a task ran. It does not tell you how far from its
scheduled moment it ran, nor that a firing you expected never happened at all.
Those are the two questions that matter when you have just installed a release on
a new machine, and neither was answerable before.

Every script comes as a matched pair — a PowerShell `.ps1` and a POSIX `.sh` —
with identical behavior and one-to-one options. A PowerShell `-FooBar` is
`--foo-bar` in the shell twin. There are no options on one side without a
counterpart on the other.

## Prerequisites

- **`sqlite3` 3.33.0 or later.** Required. The floor is not arbitrary: the
  scripts use `.param set` (3.32.0) for bound parameters and `.mode json`
  (3.33.0) for machine-readable output.
- **PowerShell 7+** (`pwsh`) for the `.ps1` twins — including on Linux and macOS,
  where they work fine.
- **`bash`** for the `.sh` twins.
- A running `goschedd` for anything involving the scheduler.

If `sqlite3` is missing, the scripts exit **2** and tell you both ways to fix it:

```bash
pwsh -File test/scripts/Test-Heartbeat.ps1 -InstallSqlite
```

That downloads the pinned build from sqlite.org, verifies its SHA-256 against
[`lib/sqlite-manifest.json`](../test/scripts/lib/sqlite-manifest.json), and
unpacks it into `test/scripts/.bin/` (gitignored). It is the only thing in this
feature that touches the network, and only when you pass the flag. A checksum
mismatch deletes the download and fails — it never installs an unverified binary.

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

Databases live in a **user-writable** directory, never the daemon's system-wide
data directory — test payloads must run unelevated, and disposable test output
does not belong beside live scheduler state:

| Platform | Location |
|---|---|
| Windows | `%LOCALAPPDATA%\goschedule-test\` |
| Linux | `$XDG_DATA_HOME/goschedule-test/` or `~/.local/share/goschedule-test/` |
| macOS | `~/Library/Application Support/goschedule-test/` |

Override with `GOSCHEDULE_TEST_DIR`. Every script prints the path it resolved.

> **A task using `run_as` resolves a *different* per-user directory** — the one
> belonging to the user the task runs as. That is correct, and occasionally
> surprising. Pass an explicit `--database` path when `run_as` is in play.

## The scripts

### Test-Heartbeat

Records one beat per invocation and exits. **This is deliberate.** The point is
to verify that `goschedd` supplies the cadence; a script that sleeps in a loop
tests its own `sleep`, not the scheduler.

| PowerShell | POSIX | Default | Meaning |
|---|---|---|---|
| `-Database` / `-d` | `--database` / `-d` | `heartbeat` | Name or path |
| `-IntervalSeconds` / `-i` | `--interval-seconds` / `-i` | none | Declared schedule interval; enables drift and gap detection |
| `-Label` / `-l` | `--label` / `-l` | none | Tag recorded on each beat |
| `-Loop` / `-r` | `--loop` / `-r` | off | Opt-in bounded continuous mode |
| `-MaxBeats` / `-m` | `--max-beats` / `-m` | none | Stop after N beats |
| `-DurationSeconds` / `-t` | `--duration-seconds` / `-t` | 3600 under `-Loop` | Stop after N seconds |
| `-SleepSeconds` / `-s` | `--sleep-seconds` / `-s` | 0 | Deliberately extend the run |
| `-FailWith` / `-f` | `--fail-with` / `-f` | none | Exit non-zero after recording |

Notes that matter:

- **`-Loop` has no unbounded form.** With neither bound supplied the 3600-second
  default applies. A runaway loop launched under a scheduler is a resource
  incident, not a test.
- **The duration bound is checked between beats**, so one deliberately-slow run
  can overrun it by up to that run's length. Interrupting a run mid-write would
  corrupt the record the bound exists to protect.
- **`-FailWith` rejects 0 and 2.** Those are reserved for success and for unmet
  prerequisites, and an induced failure must never be mistakable for either.
- **The beat is written once, at the end of the run.** A run interrupted
  mid-flight therefore records nothing and shows up as a missed firing — which
  is the honest signal, since you cannot act differently on a run that vanished
  than on one that never started.

### Test-GetSystemInfo

Records a snapshot of the machine: identity, user, process count, uptime,
network addresses, listening ports. As a scheduler test it exercises subprocess
spawning, platform tooling, and multi-row writes — where cross-platform bugs
actually surface.

| PowerShell | POSIX | Default | Meaning |
|---|---|---|---|
| `-Database` / `-d` | `--database` / `-d` | `system` | Name or path |
| `-InvocationSource` / `-i` | `--invocation-source` / `-i` | `manual` | Tag recorded on the snapshot |
| `-SkipPorts` / `-s` | `--skip-ports` / `-s` | off | Skip the port probe |

**Probes degrade; they do not abort.** A probe that cannot run on this host
records `NULL`, warns on stderr, and the snapshot is still written with exit 0.
`NULL` means *could not determine*; it is never used for a legitimate zero,
because a process count of zero and an unavailable process count support opposite
conclusions.

`process_name` on ports is commonly `NULL` — most platforms want elevation to
attribute a socket to a process. That is normal, not a defect.

### Test-ReadTestDB

| PowerShell | POSIX | Default | Meaning |
|---|---|---|---|
| `-Database` / `-d` | `--database` / `-d` | `heartbeat` | `heartbeat`, `system`, or a path |
| `-Query` / `-k` | `--query` / `-k` | `summary` | Which canned query |
| `-List` / `-n` | `--list` / `-n` | off | List queries and exit |
| `-Format` / `-f` | `--format` / `-f` | `Table` | `Table`, `Json`, `Csv` |
| `-Limit` / `-m` | `--limit` / `-m` | 20 | Row cap |
| `-IntervalSeconds` / `-i` | `--interval-seconds` / `-i` | inferred | Interval for gap and reliability checks |

> `-Query` uses the alias `-k`, not `-Q`. PowerShell aliases are
> case-insensitive, so `-Q` would collide with `-Quiet`. `-Loop` uses `-r` for
> the same reason — `-L` would collide with `-Label`.

| Query | Database | Answers |
|---|---|---|
| `summary` | both | How many records, over what period, from how many sessions and hosts? |
| `recent` | both | What are the most recent records? |
| `cadence` | heartbeat | What were the observed intervals? |
| `drift` | heartbeat | How far from expected did firings land, by source? |
| `gaps` | heartbeat | Which expected firings were missed or badly delayed? |
| `overlaps` | heartbeat | Which runs overlapped in time? |
| `failures` | heartbeat | Which runs reported failure? |
| `restarts` | heartbeat | Where are the session boundaries? |
| `hosts` | both | Which hosts and users produced records? |
| `listeners` | system | What is listening, and what changed since the previous snapshot? |
| `schema` | both | What is the stored structure? |

**Three reporting rules are contract, not presentation.** Breaking any of them
would produce numbers that look more certain than they are:

1. A query that excludes rows says how many and why. A percentile over an
   unstated subset is a confident number drawn from unknown evidence.
2. Drift is never pooled across sources and never shown without its source. A
   measured value and a derived one are different kinds of number.
3. `gaps` says whether its interval was supplied or inferred.

## What gets recorded

Full schema: [`data-model.md`](../specs/006-maintainer-test-scripts/data-model.md).

`heartbeat.db` holds one `beat` row per completed run, carrying start and finish
moments (both, so overlap is decidable rather than guessed), session and
sequence, expected moment with its source and the resulting drift, exit code and
outcome.

`system.db` holds a `snapshot` row per invocation with `snapshot_address` and
`snapshot_port` children.

**Nothing is pruned, rotated, or size-capped.** These are disposable test
artifacts, and deleting the file is the documented reset — automatic retention
would silently destroy the history you are trying to inspect. For scale: a beat
is roughly 200 bytes, so a one-minute schedule is about 525,000 rows and ~100 MB
per year. A snapshot with children is ~2 KB, so hourly is ~18 MB per year.

## Recipes

### Does it fire on time?

```bash
gosched task add "hb-verify" --command pwsh --arg -File --arg test/scripts/Test-Heartbeat.ps1 --arg -IntervalSeconds --arg 60 --schedule "every 1 minute"
```

Wait ten minutes, then:

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Query drift -IntervalSeconds 60
```

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Query gaps -IntervalSeconds 60
```

Expect ~10 beats, `boundary` source, no gaps. Compare the drift against the
project's documented p99 < 100 ms dispatch budget. This *measures and compares*;
it does not certify.

### Does it survive a restart?

```bash
gosched service restart
```

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Query restarts
```

Each invocation is its own session, so under a scheduler the session boundary
marks every beat — the evidence is the uninterrupted `started_ms` sequence across
the restart.

### Does catch-up work?

Stop the daemon past a scheduled firing, start it again, then check `gaps`.
Under a make-up-once policy expect exactly one make-up beat before the cadence
resumes. More than one, or none, is a real finding.

### Do overlap policies hold?

```bash
gosched task add "hb-overlap" --command pwsh --arg -File --arg test/scripts/Test-Heartbeat.ps1 --arg -SleepSeconds --arg 90 --schedule "every 1 minute" --overlap queue_one
```

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Query overlaps
```

Expect: `queue_one` → no overlaps, serialized; `skip` → no overlaps and visibly
fewer beats than intervals; `allow_concurrent` → overlapping ranges. Delete and
recreate the task between policies.

### Does failure handling work?

```bash
gosched task add "hb-fail" --command pwsh --arg -File --arg test/scripts/Test-Heartbeat.ps1 --arg -FailWith --arg 1 --schedule "every 1 minute"
```

Then check the `failures` query *and* `gosched alerts`. The beat proves the
record survived the failure; the alert proves the daemon saw it.

## How drift is derived

**Drift here is derived, not reported by the scheduler.** The executor passes a
spawned task the inherited environment plus the task's own configured variables
and nothing else — no scheduled time, no run ID. So the expected moment comes
from one of three sources, in strict precedence, and the source is recorded on
every beat:

| Source | Meaning |
|---|---|
| `env` | `GOSCHED_SCHEDULED_TIME` from the environment. Not set today; the tier exists so a future release is consumed with no change here. |
| `boundary` | The run's start snapped to the nearest `-IntervalSeconds` boundary. **This is the working path.** |
| `none` | Neither available. No drift recorded, and the reader excludes the beat. |

Boundary snapping works because a wall-clock-aligned schedule ("every 1 minute"
fires at `:00`) has the boundary *as* its expected moment, so the difference is
genuine absolute dispatch latency rather than mere jitter. It stays unambiguous
while drift is well under half the interval — 100 ms against 60 s is three orders
of magnitude of headroom. Past a quarter of the interval the reader flags the
figure as unreliable rather than reporting it as fact, because beyond that a late
firing and an early next one are not distinguishable.

The rejected alternative is worth knowing: inferring drift from the *observed
cadence* measures jitter, so a scheduler that is uniformly five seconds late
scores perfectly. That is the single defect class this feature most needs to
catch.

## Exit codes

| Code | Meaning |
|---|---|
| `0` | Success |
| `1` | Runtime failure — a required probe failed, a write failed, contention exhausted its retries |
| `2` | Usage error or unmet prerequisite — bad arguments, `sqlite3` absent or too old, unsupported platform |

The distinction between 1 and 2 is load-bearing: it is what makes a run row in
`gosched runs` mean something. An unmet prerequisite is a usage-class failure —
conflating it with a runtime failure sends you debugging the wrong thing.

Results go to stdout, diagnostics and warnings to stderr, so reader output pipes
cleanly into another tool without log lines mixed in.

## Shell conventions

The `.ps1` twins follow the ShruggieTech PowerShell standard, enforced by
[`.claude/skills/shruggie-powershell/scripts/Test-ScriptCompliance.ps1`](../.claude/skills/shruggie-powershell/scripts/Test-ScriptCompliance.ps1).
There is **no equivalent shell standard yet**, so the `.sh` twins are governed by
this section plus `shellcheck`:

- `#!/usr/bin/env bash` and `set -euo pipefail`.
- A header comment block covering purpose, the twin relationship, the exit-code
  contract, and a pointer here.
- Options parsed in a `while`/`case` loop; every PowerShell parameter has a
  `--kebab-case` counterpart and the same single-letter short option.
- `log LEVEL MESSAGE` for all diagnostics, to stderr. Results to stdout.
- `die_usage` exits 2; `die_runtime` exits 1. No bare `exit 1` for a usage error.
- All shared logic lives in `lib/sqlite.sh` — one implementation per twin. Three
  copies of the resolution order would be three chances for them to disagree,
  and disagreement there presents as an intermittent platform-specific failure.
- Quote every expansion. Values written to the database include hostnames and
  usernames, which can contain quotes and spaces.
- **Never build SQL by string interpolation.** Use `sqlite_exec`'s bound
  parameters. Interpolating a username into SQL is both an injection vector and
  an ordinary bug for anyone named `O'Brien`.

## Troubleshooting

**`exit 2` with "No usable sqlite3 found"** — none of the three search locations
held a usable 3.33.0+. Use `--install-sqlite` or your package manager.

**`--sqlite-exe` exits 2 even though the file exists** — on Windows, pass a
native path (`C:\bin\sqlite3.exe`), not a Git Bash path (`/c/bin/sqlite3`).
PowerShell cannot resolve the latter. An explicit path that cannot be used is a
hard error by design: silently falling back to a different `sqlite3` than the one
you named is how you debug the wrong binary for an hour.

**"address probe unavailable" on Windows under Git Bash** — expected. Neither
`ip` nor `ifconfig` exists there. Use the PowerShell twin, which uses
`Get-NetIPAddress`. The snapshot is still recorded.

**Drift looks enormous** — check that `-IntervalSeconds` matches the schedule you
actually registered. A mismatched interval snaps to the wrong boundary. The
reader flags values past a quarter of the interval for this reason.

**`drift` returns no rows** — every beat had `expected_source = 'none'`, meaning
no interval was declared. The excluded-count warning tells you how many. Pass
`-IntervalSeconds` when recording, not just when reading.

**Contention errors under `allow_concurrent`** — the scripts use WAL with a
5-second busy timeout and three retries. Exhausting that is reported as a
*harness* failure, explicitly, so it is not mistaken for a scheduler defect.

## Automated tests

```bash
go test ./test/integration/ -run TestScripts -v
```

Covers single-shot recording, boundary drift, the exit-code contract, the
bounded loop, concurrent writers, twin parity, snapshots, and every canned query.
It **skips with a stated reason** when `sqlite3`, `pwsh`, or `bash` is missing —
a skip is not a pass, and the reason is printed so the two cannot be confused.
