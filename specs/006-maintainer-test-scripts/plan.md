# Implementation Plan: Maintainer Test Scripts and Vendored Skills

**Branch**: `main` (trunk-based) | **Date**: 2026-07-23 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/006-maintainer-test-scripts/spec.md`

## Summary

Ship three cross-platform script pairs under `test/scripts/` that let a maintainer prove an
installed `goschedd` actually fires on time, survives restarts, catches up after downtime,
and honors its overlap policies — recording the evidence into two local SQLite databases and
reading it back through a set of canned queries. Consolidate the documentation into
`docs/test-scripts.md`. Separately, narrow the `.gitignore` exclusion of `.claude/` so the
skills subtree is tracked, and vendor the house-standard skills plus a new project-native
verification skill.

The load-bearing technical decision is how drift is measured. The executor injects no
scheduler context into a spawned task (verified — `internal/executor/executor.go:42`), so the
expected firing moment is derived by snapping the beat's start to the nearest boundary of a
caller-declared interval, and every drift figure carries the source it was derived from.
This yields absolute dispatch latency with no product code change.

## Technical Context

**Language/Version**: PowerShell 7+ (`pwsh`), POSIX shell (`bash`), Go 1.25.0 for the tests

**Primary Dependencies**: `sqlite3` CLI ≥ 3.33.0 — external, detected, optionally installed;
never vendored as a binary

**Storage**: two SQLite databases in a user-writable per-user directory, WAL mode

**Testing**: `test/integration/testscripts_test.go`, driving the scripts as subprocesses
inside the existing `go test ./...` invocation

**Target Platform**: Windows, Linux, macOS

**Project Type**: maintainer tooling in an existing Go CLI/daemon repository

**Performance Goals**: none of its own. The feature *measures* the project's p99 < 100 ms
dispatch budget; it does not impose a budget on itself.

**Constraints**: no product code change; no CI workflow change; no network access without an
explicit opt-in flag; no third-party binaries committed

**Scale/Scope**: ~10 new files, ~2 modified. No change to the shipped binaries.

## Constitution Check

*GATE: passed before Phase 0. Re-checked after Phase 1 design — result unchanged.*

| Principle | Assessment |
|---|---|
| **I. Code Quality** | PASS. No Go product code changes. The one new Go file (the integration test) is subject to `gofmt`/`vet`/lint like any other. The shell deliverables get their own mechanical gates: the ShruggieTech PowerShell compliance checker and `shellcheck`. Error handling is explicit in both twins, and the exit-code contract is the machine-readable form of it. |
| **II. Testing Standards** | PASS, with a caveat recorded honestly. Behavioral tests ship with the code and cover the exit-code contract, the bounded loop, concurrent writers, and each canned query. No `time.Sleep`-based assertions: the tests drive the scripts directly rather than waiting on a scheduler. **Caveat**: this machine has no C compiler, so `-race` cannot be run locally (research §R8); it is deferred to CI and reported as such at the halt, never as a pass. Coverage: the gate measures six core packages, none of which this feature touches, so the gate is unaffected — but it is still run, because "should be unaffected" and "is unaffected" are different claims. |
| **III. UX Consistency** | PASS, and this is where the feature does most of its constitutional work. Results to stdout, diagnostics to stderr, conventional exit codes, RFC 3339 timestamps, both human and machine-readable output — the same contract the CLI already honors. Error messages name what failed and what to do (FR-017 requires the missing-tool message to name both remedies). |
| **IV. Performance** | PASS / not applicable. No hot path, no benchmark, no product code. |
| **V. Autonomous Build-Phase Execution** | PASS. Placed under autopilot by explicit operator request. Routine decisions were made and recorded (spec Clarifications, research §R1–R10) rather than escalated. One halt, before push. |

**Engineering constraints**: platform support is Windows + Linux + macOS, matching the
constitution's Linux-and-Windows requirement and exceeding it. No new Go dependency. No
secrets logged — and the repository-hygiene requirements (FR-026a/b) exist specifically to
keep the `.gitignore` change from *introducing* a secrets-tracking path.

**Complexity Tracking**: no violations; the table is omitted.

## Project Structure

### Documentation (this feature)

```text
specs/006-maintainer-test-scripts/
├── plan.md              # This file
├── spec.md
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1
├── contracts/
│   └── cli.md           # Phase 1
├── checklists/
│   ├── requirements.md
│   └── portability.md
└── tasks.md             # Phase 2 (/speckit-tasks)
```

### Source Code (repository root)

```text
test/scripts/
├── Test-GetSystemInfo.ps1        Test-GetSystemInfo.sh
├── Test-Heartbeat.ps1            Test-Heartbeat.sh
├── Test-ReadTestDB.ps1           Test-ReadTestDB.sh
├── lib/
│   ├── Sqlite.ps1                sqlite.sh          # one impl per twin (FR-021d)
│   └── sqlite-manifest.json                         # pinned versions + SHA-256
├── .bin/                                            # gitignored install target
└── README.md                                        # pointer to docs/test-scripts.md

test/integration/
└── testscripts_test.go                              # new

docs/
└── test-scripts.md                                  # new, consolidated

.claude/skills/                                      # newly tracked
├── speckit-*/                                       # 8 existing, become tracked
├── shruggie-powershell/  shruggie-markdown/
├── shruggie-speckit/     gh-fix-ci/
└── go-schedule-verify/                              # new, project-native

.gitignore                                           # modified (PINNED ARTIFACT)
CHANGELOG.md                                         # modified
README.md                                            # one-line pointer
CLAUDE.md                                            # active-feature pointer
```

**Structure Decision**: scripts live under `test/scripts/`, not `scripts/`. The existing
`scripts/` directory holds build and CI tooling (`coverage-gate.sh`) that CI invokes
directly; these are test artifacts and belong beside `test/integration/`, which is what
exercises them.

`lib/` exists because three scripts per twin each need prerequisite resolution, installation,
and database access. Three copies would be three chances to disagree, and disagreement there
presents as an intermittent platform-specific failure — hence FR-021d requiring exactly one
implementation per twin.

## Design Decisions

Full rationale in [research.md](research.md); the four that shape the code:

1. **Drift by boundary snapping** (§R1). Three-tier precedence — environment, boundary, none
   — with the source recorded on every beat. Rejected: cadence inference (measures jitter, so
   blind to a uniformly-late scheduler) and executor modification (changes a safety-critical
   surface for tooling's benefit, and forfeits the provably-unchanged-binaries property).
2. **Bound parameters via `sqlite3 .param set`** (§R3). The values written include hostname,
   username, and interface names — attacker-influenceable on a shared machine and capable of
   containing a quote. String-interpolated SQL here is both an injection vector and a plain
   bug for any user named `O'Brien`. This sets the 3.33.0 minimum, together with `.mode json`.
3. **One write per beat, at run end** (§R5, FR-021c). Halves contention. An interrupted run
   records nothing and appears as a missed firing — the honest signal, since a maintainer
   cannot act differently on a run that vanished mid-flight than on one that never started.
4. **User-writable database location** (§R2). The daemon's data directory is system-wide and
   needs elevation; test payloads must run unelevated, and disposable test output does not
   belong in the directory holding live scheduler state.

## Implementation Phases

**Phase A — shared library.** `lib/Sqlite.ps1` and `lib/sqlite.sh`: resolution order with
version gating, the checksum-verified installer, WAL/busy-timeout connection setup, bound
parameter execution, schema creation, the data-directory resolver, and structured stderr
logging. Everything else depends on this; it is where the parity risk concentrates.

**Phase B — recording scripts.** Heartbeat then system-info, both twins. Platform branching
in the PowerShell twin (`$IsWindows`) is the highest-risk code in the feature.

**Phase C — reader.** The eleven canned queries plus the reporting obligations that are
contract rather than presentation: excluded-row disclosure, per-source drift breakdown,
unreliability flagging, inferred-vs-supplied interval disclosure.

**Phase D — tests.** `testscripts_test.go`, with stated-reason skips.

**Phase E — documentation.** `docs/test-scripts.md` to the vendored Markdown house style,
plus the README and CLAUDE.md pointers.

**Phase F — repository configuration.** The `.gitignore` negation, skill vendoring, the new
`go-schedule-verify` skill, the pre-commit hygiene check (FR-026b), and the dated CHANGELOG
decision entry the pinned-artifact change requires (FR-029).

## Risks

| Risk | Mitigation |
|---|---|
| PowerShell twin assumes Windows-only cmdlets | `$IsWindows` branching; tests exercise the twin on the host platform; called out in research §R6 as the likeliest defect |
| Checksums cannot be fetched at implementation time | Manifest ships empty and the installer refuses to run. A plausible wrong hash is worse than no installer — it turns a loud failure into a silent acceptance |
| `.gitignore` negation sweeps in a credential file | FR-026b requires an explicit `git status` inspection before commit; verified at the halt |
| `-race` unrunnable locally | Reported explicitly at the halt, non-race suite run in its place, race gate deferred to CI. Never reported as passing |
| Twin drift over time | Parity is a mechanical naming rule (contracts/cli.md), not a judgment call, so a missing counterpart is visible on inspection |
