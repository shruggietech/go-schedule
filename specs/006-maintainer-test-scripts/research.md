# Phase 0 Research: Maintainer Test Scripts and Vendored Skills

**Feature**: 006-maintainer-test-scripts | **Date**: 2026-07-23

All unknowns from Technical Context are resolved below. Each entry records the decision, why
it was chosen, and what was rejected — per the autopilot decision policy, these were decided
rather than escalated.

---

## R1 — How is drift measured, given the scheduler supplies no context?

**Investigation**: `internal/executor/executor.go:42` builds the child environment as
`os.Environ()` plus the task's configured `Env` map. Nothing scheduler-generated is injected:
no scheduled time, no task ID, no run ID. Confirmed by reading the whole executor; there is
no other environment path.

**Decision**: Three-tier precedence for the expected moment, with the source recorded on
every beat.

1. `GOSCHED_SCHEDULED_TIME` from the environment, if present. Absent today; present-proofed.
2. Snap the beat's actual start to the nearest boundary of the caller-supplied interval.
3. None — record no expected moment and no drift.

**Rationale**: Tier 2 is the working path. For a wall-clock-aligned interval schedule
("every 1 minute" fires at `:00`), the expected moment *is* the nearest boundary, so
`actual - boundary` is genuine absolute dispatch latency, not jitter. It stays unambiguous
while drift is well under half the interval — 100ms against a 60s interval is three orders of
magnitude of headroom — and the regime where it breaks down is the same regime the missed-firing
query already flags.

**Alternatives rejected**:

- *Infer drift from the observed cadence.* Rejected: it measures jitter. A scheduler that is
  uniformly five seconds late has a perfect inter-beat interval and would score flawless. The
  one defect class this feature most needs to catch is exactly the one this method is blind to.
- *Inject `GOSCHED_*` variables from the executor.* Rejected here, though it is the better
  long-term answer. It modifies a safety-critical product surface for the benefit of maintainer
  tooling, and it would forfeit this release's provable "the binaries did not change" property.
  Recorded as follow-up work.

---

## R2 — Where do the databases live?

**Investigation**: `internal/platform/platform_unix.go:15` and `platform_windows.go:16-19`
put the daemon's data in `/var/lib/goschedule` and `%ProgramData%\goschedule`. Both are
system-wide and require elevation to write.

**Decision**: A separate, user-writable directory, resolved in this order — the
`GOSCHEDULE_TEST_DIR` environment variable; then the platform's per-user data location
(`%LOCALAPPDATA%\goschedule-test`, `$XDG_DATA_HOME/goschedule-test` or
`~/.local/share/goschedule-test`, `~/Library/Application Support/goschedule-test`). Every
script prints the resolved path it used.

**Rationale**: Test payloads must run unelevated — a task configured with `run_as` for an
unprivileged user is a case this feature explicitly wants to support. Writing into the
daemon's data directory would both demand elevation and mix disposable test output with live
scheduler state, in the one directory a maintainer least wants polluted.

**Consequence, documented rather than designed around**: a task running under `run_as` for a
different user resolves a *different* per-user directory. That is correct behavior and useful
signal, but it surprises. The `run_as` recipes therefore pass an explicit database path.

---

## R3 — Minimum `sqlite3` version, and why it matters

**Investigation**: local tool is 3.47.2. Two features are load-bearing:

- `.param set :name value` — bound parameters in the CLI, added in **3.32.0**.
- `.mode json` — machine-readable output, added in **3.33.0**.

**Decision**: require **3.33.0 or later**; detect by parsing `sqlite3 --version` and treat a
too-old tool as *not found*, continuing the search rather than failing.

**Rationale for `.param set` specifically — this is a correctness finding, not a convenience.**
The `sqlite3` CLI has no external parameter-binding API; the obvious implementation builds an
`INSERT` by string interpolation. The values being interpolated include the hostname, the
username, and network interface names — all attacker-influenceable on a machine someone else
administers, and all capable of containing a single quote. String-built SQL here is an
injection vector *and* an ordinary correctness bug for any user named `O'Brien`. `.param set`
binds properly and removes the whole class.

**Alternatives rejected**: hand-rolled quote doubling (works until it doesn't, and the failure
is silent data corruption); a here-doc per row (same interpolation problem); requiring 3.38+
for JSON operators (unnecessary — nothing here needs them, and it excludes still-supported
distributions for no gain).

---

## R4 — How the opt-in installer pins its download

**Decision**: `lib/sqlite-manifest.json` pins, per platform-architecture pair: the upstream
version, the download URL, the expected **SHA-256**, and the date pinned. Verification happens
on the downloaded file *before* unpacking; a mismatch deletes the download and exits.

**Note on hash algorithm**: sqlite.org publishes **SHA3-256** on its download page, but
neither `Get-FileHash` nor a baseline POSIX toolchain computes SHA3 without extra
dependencies. The manifest therefore pins a SHA-256 that we compute at pin time, recording
alongside it the upstream SHA3-256 that was cross-checked on the same bytes. This keeps
verification implementable everywhere while preserving a link back to the publisher's own
figure.

**Checksums are not invented.** The manifest is populated during implementation from bytes
actually fetched and hashed. If that fetch cannot be performed, the manifest ships empty and
the installer refuses to run — a manifest with a plausible-looking wrong hash is far worse
than no installer, because it converts a loud failure into a silent acceptance.

**Coverage**: Windows x64, Linux x64, macOS x64 and arm64 — the combinations sqlite.org
publishes prebuilt tools for. Everything else exits 2 naming the package manager. No
building from source.

---

## R5 — Concurrent writers

**Decision**: every database connection opens with `PRAGMA journal_mode=WAL` and
`PRAGMA busy_timeout=5000`, and each write is attempted up to 3 times before failing with
exit 1 and a message naming contention.

**Rationale**: overlap-policy testing creates simultaneous writers by definition — that is the
point of `--sleep-seconds`. Under the default rollback journal, a concurrent writer gets
`SQLITE_BUSY` immediately, which would surface as a test-harness failure that reads exactly
like a scheduler defect. WAL plus a 5-second timeout makes ordinary contention invisible;
the retry covers the pathological case; the bounded failure keeps a genuinely stuck database
from hanging a scheduled task forever.

**Write model**: one write per beat, at the end of the run, carrying a start moment captured
in memory when the run began. Two writes per beat would double contention to buy a
mid-flight-crash record that no maintainer can act on differently from a missed firing.

---

## R6 — Host inspection, per platform, with fixed fallback order

Ordering is fixed and documented so that provenance does not silently vary between hosts.

| Datum | Windows | Linux | macOS |
|---|---|---|---|
| Addresses | `Get-NetIPAddress` | `ip -j addr` → `ip addr` → `ifconfig` | `ifconfig` |
| Listening ports | `Get-NetTCPConnection -State Listen` + `Get-NetUDPEndpoint` | `ss -lntup` → `netstat -lntup` → `lsof` | `lsof -iTCP -sTCP:LISTEN` → `netstat -an` |
| Process count | `(Get-Process).Count` | `ps -e` line count | `ps -ax` line count |
| Uptime | `Get-Uptime` | `/proc/uptime` | `sysctl -n kern.boottime` |
| User | `$env:USERNAME` | `id -un` → `$USER` | `id -un` → `$USER` |

**The portability trap, stated explicitly**: `Get-NetIPAddress`, `Get-NetTCPConnection`,
`Get-NetUDPEndpoint`, and `Get-Uptime` ship in Windows-only modules. PowerShell 7 itself runs
on all three platforms, so `Test-GetSystemInfo.ps1` **must** branch on `$IsWindows` and fall
through to the POSIX commands otherwise. This is the single likeliest source of a
works-on-my-machine defect in the feature, and it is why the automated tests exercise the
PowerShell twin on whatever platform CI runs.

Every probe degrades to a recorded `NULL` plus a warning on stderr. A `NULL` means "could not
determine"; a `0` means "determined, and it was zero". The distinction is preserved because a
process count of zero and an unavailable process count support opposite conclusions.

---

## R7 — Testing shell scripts from Go

**Decision**: `test/integration/testscripts_test.go` drives the scripts as subprocesses
against a `t.TempDir()` database and asserts on rows read back through `sqlite3`. Each test
`t.Skip`s with a specific reason when its interpreter or `sqlite3` is missing.

**Rationale**: it satisfies constitution principle II with real behavioral coverage, and it
runs inside the existing `go test ./...` invocation, so **no CI workflow change is needed** —
avoiding a second pinned-artifact modification for no benefit.

**Skipping is not passing.** A skip prints why it skipped. A silently-passing test on a
machine that never ran it is how a cross-platform bug reaches a release.

---

## R8 — Local verification tooling gaps (affects the halt, not the design)

Probed on this machine:

| Tool | Status |
|---|---|
| `sqlite3` 3.47.2 | present |
| `pwsh` 7.6.3 | present |
| `bash` 5.3.15 | present |
| `go` 1.25.0, `go.mod` go 1.25.0 | present — the golangci-lint toolchain trap does **not** apply |
| `shellcheck` | **absent**; `winget` available to install it |
| C compiler (`gcc`/`cc`/`clang`/mingw) | **absent** |

The missing C compiler means `CGO_ENABLED=1 go test -race` **cannot run on this machine** —
the exact trap CLAUDE.md documents. Per that document's instruction, this will be reported
explicitly at the halt rather than papered over: the non-race suite will be run in its place,
and the race gate deferred to CI. Reporting the suite as green on the strength of a run that
never happened is the specific failure the instruction exists to prevent.

---

## R9 — Vendoring mechanics

**Decision**: `.gitignore` becomes `.claude/*` plus `!.claude/skills/`. Skills are copied
from `~/.claude/skills/`, each gaining a provenance note recording source and date.
Plugin-managed `superpowers:*` skills are **not** copied — two registrations of one skill
name, diverging silently, is worse than the inconvenience it would solve.

A `git status` inspection before commit confirms no settings or credential file rides along
on the negation pattern. The `.gitignore` edit touches a pinned artifact and therefore
requires a dated `CHANGELOG.md` decision entry.

---

## R10 — What is deliberately *not* being built

- **No `shruggie-bash` skill.** The POSIX twins get a conventions section in the docs plus
  shellcheck. Authoring a full shell standard is its own piece of work and would balloon this
  feature.
- **No CI workflow change.** See R7.
- **No product code change.** The value of "the binaries are provably identical" is high and
  the cost of forfeiting it is not repaid by anything in scope here.
- **No automatic retention.** The databases are disposable; deleting one is the reset.
