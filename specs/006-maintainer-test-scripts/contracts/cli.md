# Phase 1 Contract: Script Command-Line Interfaces

**Feature**: 006-maintainer-test-scripts | **Date**: 2026-07-23

The scripts are the feature's external interface. This is their contract. Every PowerShell
option has exactly one POSIX counterpart and vice versa — parity is a mechanical rule, not a
judgment call (FR-015).

## Naming rule

A PowerShell parameter `-FooBar` corresponds to the POSIX long option `--foo-bar`. Switch
parameters correspond to flags taking no value. PowerShell parameters additionally carry
single-letter aliases per the ShruggieTech standard; the POSIX twins accept the same letters
as short options. There are no options on one side without a counterpart on the other.

## Exit codes (all six scripts)

| Code | Meaning |
|---|---|
| `0` | Success. |
| `1` | Runtime failure — a probe that was required failed, contention exhausted its retries, a write failed, a query errored. |
| `2` | Usage error or unmet prerequisite — bad arguments, `sqlite3` absent or too old, unsupported platform for the installer. |

`0` and `2` are reserved: `--fail-with` rejects both (FR-006), so an induced failure can never
be confused with success or with a missing tool.

Results go to stdout, diagnostics and warnings to stderr (FR-020). This matches the project's
existing CLI behavior (constitution principle III) and is what lets a maintainer pipe reader
output into another tool without stripping log lines out of it first.

## Common options (all six scripts)

| PowerShell | POSIX | Type | Meaning |
|---|---|---|---|
| `-SqliteExe <path>` | `--sqlite-exe <path>` | string | Explicit tool path; highest precedence in the search order. |
| `-InstallSqlite` | `--install-sqlite` | switch | Download, verify, and install the pinned tool into the repo-local `.bin/`. The only option in the feature that touches the network. |
| `-Help` | `--help`, `-h` | switch | Usage to stdout; exit 0. |
| `-Quiet` | `--quiet`, `-q` | switch | Suppress informational stderr output. Warnings and errors are never suppressed. |

## `Test-Heartbeat.ps1` / `Test-Heartbeat.sh`

| PowerShell | POSIX | Default | Meaning |
|---|---|---|---|
| `-Database <path\|name>` | `--database` | `heartbeat` | Well-known name or explicit path. |
| `-IntervalSeconds <int>` | `--interval-seconds` | none | Declared schedule interval. Enables `boundary` drift and the gap query. |
| `-Label <string>` | `--label` | none | Tag recorded on each beat. |
| `-Loop` | `--loop` | off | Opt-in continuous mode. |
| `-MaxBeats <int>` | `--max-beats` | none | Loop bound by count. |
| `-DurationSeconds <int>` | `--duration-seconds` | `3600` when `--loop` is set and no bound given | Loop bound by time. |
| `-SleepSeconds <int>` | `--sleep-seconds` | `0` | Deliberately extend the run, to exercise overlap policies. |
| `-FailWith <int>` | `--fail-with` | none | Exit with this code after recording. Rejects `0` and `2`. |

**Behavior**

- Default (no `--loop`): record one beat, exit. The scheduler owns the cadence.
- `--loop`: whichever of `--max-beats` / `--duration-seconds` is reached first ends it. With
  neither supplied, the 3600-second default applies — there is no unbounded form (FR-004).
- The duration bound is checked *between* beats, so one slow run may overrun it (FR-004a).
- The beat is written once, at the end of the run (FR-021c).

## `Test-GetSystemInfo.ps1` / `Test-GetSystemInfo.sh`

| PowerShell | POSIX | Default | Meaning |
|---|---|---|---|
| `-Database <path\|name>` | `--database` | `system` | |
| `-InvocationSource <string>` | `--invocation-source` | `manual` | Recorded on the snapshot. |
| `-SkipPorts` | `--skip-ports` | off | Omit the port probe, which is the slowest and the one most often requiring elevation. |

Exits `0` when the snapshot is recorded, even if individual probes degraded to `NULL`
(FR-009); `1` only if the snapshot itself could not be written.

## `Test-ReadTestDB.ps1` / `Test-ReadTestDB.sh`

| PowerShell | POSIX | Default | Meaning |
|---|---|---|---|
| `-Database <path\|name>` | `--database` | required | `heartbeat`, `system`, or a path. |
| `-Query <name>` | `--query` | `summary` | Canned query to run. |
| `-List` | `--list` | off | List available queries with the question each answers; exit 0. |
| `-Format <Table\|Json\|Csv>` | `--format` | `Table` | Output form. |
| `-Limit <int>` | `--limit` | `20` | Row cap for queries that return rows. |
| `-IntervalSeconds <int>` | `--interval-seconds` | inferred | Expected interval for the gap and reliability checks; inferred from the modal observed interval when omitted, and the inference is stated in the output. |

### Canned queries

| Name | DB | Question answered |
|---|---|---|
| `summary` | both | How many records, spanning what period, from how many sessions and hosts? |
| `recent` | both | What are the most recent records? |
| `cadence` | heartbeat | What were the observed intervals — min, p50, p95, p99, max? |
| `drift` | heartbeat | How far from expected did firings land, **broken down by expected-source**? |
| `gaps` | heartbeat | Which expected firings were missed or badly delayed? |
| `overlaps` | heartbeat | Which runs overlapped in time? |
| `failures` | heartbeat | Which runs reported failure? |
| `restarts` | heartbeat | Where are the session boundaries, and did recording continue across them? |
| `hosts` | both | Which hosts and users produced records? |
| `listeners` | system | What is listening now, and what changed since the previous snapshot? |
| `schema` | both | What is the stored structure? |

**Reporting obligations** — these are contract, not presentation:

- Any query excluding rows states how many it excluded and why (FR-013a). A percentile over
  an unstated subset is a confident number drawn from unknown evidence.
- `drift` never pools `env`-sourced and `boundary`-sourced rows into one figure (FR-013b), and
  never emits a drift number without its source (FR-003).
- `drift` flags boundary-derived values exceeding a quarter of the interval as unreliable.
- `gaps` states whether its interval was supplied or inferred.

## Prerequisite resolution (shared by all six)

1. `--sqlite-exe` if given.
2. Repo-local `test/scripts/.bin/`.
3. `PATH`.

A candidate that exists but is not executable, or reports below **3.33.0**, is treated as not
found and the search continues (FR-016a) — a stale tool early in the order must not shadow a
good one later.

On exhaustion: exit `2`, naming both `--install-sqlite` and the platform's package manager
command (`winget install SQLite.SQLite`, `apt install sqlite3`, `brew install sqlite`).

## Installer contract

1. Resolve platform and architecture. Unsupported → exit `2` naming the package manager; no
   source builds (FR-018a).
2. Look up the pinned entry in `lib/sqlite-manifest.json`.
3. Download to a temporary path.
4. Compute SHA-256 and compare to the pin. **Mismatch → delete the download and exit `1`.**
   The artifact is never left on disk (FR-018b).
5. Unpack into `test/scripts/.bin/`, verify the result runs and reports an acceptable
   version, and report the installed path.

Verification precedes unpacking. Nothing unverified is ever placed where step 2 of the
resolution order would later find it.
