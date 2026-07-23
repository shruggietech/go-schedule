# Tasks: Maintainer Test Scripts and Vendored Skills

**Feature**: 006-maintainer-test-scripts | **Date**: 2026-07-23
**Input**: [plan.md](plan.md), [spec.md](spec.md), [data-model.md](data-model.md),
[contracts/cli.md](contracts/cli.md), [research.md](research.md), [quickstart.md](quickstart.md)

Tests are **required**, not optional — constitution principle II is non-negotiable and
explicitly covers integration tests for concurrent execution and recovery.

**Parallelism note**: `[P]` marks tasks touching different files with no incomplete
dependency. The dominant constraint here is that the PowerShell and POSIX twins of any given
script are always parallelizable with each other, but both depend on their own `lib/`.

---

## Phase 1: Setup

- [ ] T001 Create the `test/scripts/` and `test/scripts/lib/` directory structure per plan.md
- [ ] T003 [P] Write `test/scripts/README.md` as a short pointer to `docs/test-scripts.md`, not a duplicate of it

> **T002 was merged into T041.** Both touched `.gitignore`, which is a pinned artifact; two
> edits in two phases meant two chances for only one of them to reach the changelog decision
> entry that a pinned change requires. One task, one edit, one record.

---

## Phase 2: Foundational (blocks every user story)

This is the shared library. FR-021d requires exactly one implementation per twin, so every
later task consumes these rather than re-deriving them.

- [ ] T004 Create `test/scripts/lib/sqlite-manifest.json` with the pinned upstream version, per-platform URLs, SHA-256 fields, cross-checked SHA3-256, and the pin date (research §R4)
- [ ] T005 Populate the manifest checksums from bytes actually fetched and hashed. If the fetch cannot be performed, leave them empty and make the installer refuse to run — never write a plausible-looking hash (research §R4)
- [ ] T006 Implement `Resolve-Sqlite` in `test/scripts/lib/Sqlite.ps1`: strict precedence (explicit → `.bin/` → PATH), version gate at 3.33.0, and treat an unusable or too-old candidate as not-found so the search continues (FR-016, FR-016a)
- [ ] T007 [P] Implement the same resolution in `test/scripts/lib/sqlite.sh` with identical precedence and version gating (FR-015, FR-016a)
- [ ] T008 Implement `Install-Sqlite` in `test/scripts/lib/Sqlite.ps1`: platform/arch detection, download to temp, SHA-256 verify **before** unpack, delete on mismatch, unpack to `.bin/`, post-install version check (FR-018, FR-018a, FR-018b)
- [ ] T009 [P] Implement the same installer in `test/scripts/lib/sqlite.sh` (FR-018)
- [ ] T010 Implement the data-directory resolver in `test/scripts/lib/Sqlite.ps1`: `GOSCHEDULE_TEST_DIR` → platform per-user location; never the daemon's system-wide dir; report the resolved path (FR-011a, research §R2)
- [ ] T011 [P] Implement the same resolver in `test/scripts/lib/sqlite.sh` (FR-011a)
- [ ] T012 Implement bound-parameter execution in `test/scripts/lib/Sqlite.ps1` using `.param set`, plus WAL and `busy_timeout=5000` on every connection, and a bounded 3-attempt retry that fails with exit 1 naming contention (FR-021, FR-021a, research §R3, §R5)
- [ ] T013 [P] Implement the same in `test/scripts/lib/sqlite.sh` (FR-021, FR-021a)
- [ ] T014 Implement schema creation for both databases in `test/scripts/lib/Sqlite.ps1` per data-model.md, including `meta.schema_version`, indexes, and CHECK constraints for the validation rules
- [ ] T015 [P] Implement the same schema creation in `test/scripts/lib/sqlite.sh`
- [ ] T016 Implement the shared exit-code contract and `Write-Log`-style structured stderr logging in both `lib/` files: 0 success, 1 runtime, 2 usage/prerequisite; results to stdout, diagnostics to stderr (FR-019, FR-020)
- [ ] T017 Implement the missing-tool failure message in both `lib/` files, naming **both** `-InstallSqlite`/`--install-sqlite` and the platform package-manager command, exiting 2 (FR-017)

**Checkpoint**: both `lib/` files resolve, install, connect, create schema, and log identically.

---

## Phase 3: User Story 1 — Prove the scheduler fires on time (P1) 🎯 MVP

**Goal**: a maintainer can register a recurring task and get a quantified on-time-firing
verdict with drift and missed-firing figures.

**Independent test**: register one task against the heartbeat script, let it run several
intervals, run the `cadence`, `drift`, and `gaps` queries.

- [ ] T018 [US1] Write the failing integration test for single-shot recording in `test/integration/testscripts_test.go`: one invocation writes exactly one beat, exit 0, with `t.Skip` and a stated reason when `sqlite3`/`pwsh`/`bash` is absent (research §R7)
- [ ] T019 [US1] Implement `test/scripts/Test-Heartbeat.ps1` single-shot path: capture start in memory, do the work, write one beat at run end (FR-001, FR-021c)
- [ ] T020 [P] [US1] Implement the same single-shot path in `test/scripts/Test-Heartbeat.sh` (FR-015)
- [ ] T021 [US1] Implement the three-tier expected-moment precedence in both heartbeat twins — `GOSCHED_SCHEDULED_TIME` → interval boundary snap → none — recording `expected_source` on every beat and never emitting drift without it (FR-003, FR-003a, research §R1)
- [ ] T022 [US1] Implement `-SleepSeconds` / `--sleep-seconds` and `-FailWith` / `--fail-with` in both heartbeat twins, with `--fail-with` rejecting the reserved codes 0 and 2, and the beat still written on induced failure (FR-005, FR-006)
- [ ] T023 [US1] Implement bounded `-Loop` / `--loop` in both heartbeat twins: first-bound-wins, a default duration bound when neither is given, and the duration checked between beats (FR-004, FR-004a)
- [ ] T024 [US1] Write the failing integration test for the bounded loop and the exit-code contract in `test/integration/testscripts_test.go`: `--fail-with` exits as asked *and* records; `--max-beats` bounds the loop; a bogus `--sqlite-exe` exits **2**, not 1
- [ ] T025 [US1] Implement `test/scripts/Test-ReadTestDB.ps1` skeleton: database resolution, `-List`, `-Format`, `-Limit`, and the `summary`, `recent`, and `schema` queries (FR-011, FR-012, FR-014)
- [ ] T026 [P] [US1] Implement the same reader skeleton in `test/scripts/Test-ReadTestDB.sh` (FR-015)
- [ ] T027 [US1] Implement the `cadence`, `drift`, and `gaps` queries in both reader twins, including the per-source drift breakdown, the excluded-row disclosure, the quarter-interval unreliability flag, and the inferred-vs-supplied interval disclosure (FR-013, FR-013a, FR-013b)

**Checkpoint**: US1 is independently shippable — quickstart scenarios 1, 2, and 8 pass.

---

## Phase 4: User Story 2 — Restarts, downtime, and overlap (P2)

**Goal**: the harder guarantees become answerable from recorded history.

**Independent test**: stop the daemon across a firing, restart, confirm one make-up beat and
a visible session boundary.

- [ ] T028 [US2] Write the failing integration test for concurrent writers in `test/integration/testscripts_test.go`: two simultaneous invocations both land their beats, neither fails on contention (FR-021, SC-007)
- [ ] T029 [US2] Implement the `overlaps` query in both reader twins as an intersection of `[started_ms, finished_ms]` ranges — decidable only because both endpoints are stored (FR-013, data-model)
- [ ] T030 [P] [US2] Implement the `restarts` query in both reader twins: session and pid boundaries, confirming recording continued across them (FR-013)
- [ ] T031 [P] [US2] Implement the `failures` query in both reader twins over non-zero `exit_code` (FR-013)
- [ ] T032 [US2] Add a CHECK-backed assertion that `outcome` and `exit_code` never disagree, and an integration test covering it (FR-017)

**Checkpoint**: quickstart scenarios 3, 4, 5, and 6 pass.

---

## Phase 5: User Story 3 — Host snapshot as a realistic workload (P2)

**Goal**: a scheduled task that exercises subprocess spawning, platform tooling, and
multi-row writes — where cross-platform bugs actually surface.

**Independent test**: invoke once per platform; a complete snapshot lands with addresses and
ports attached.

- [ ] T033 [US3] Write the failing integration test for snapshot recording in `test/integration/testscripts_test.go`: one snapshot with non-null required columns, plus child rows where the platform permits
- [ ] T034 [US3] Implement `test/scripts/Test-GetSystemInfo.ps1` core snapshot: timestamps, tz offset, hostname, username, process count, uptime, platform, `script_flavor`, `invocation_source` (FR-007, FR-010)
- [ ] T035 [P] [US3] Implement the same core snapshot in `test/scripts/Test-GetSystemInfo.sh` (FR-015)
- [ ] T036 [US3] Implement `$IsWindows` branching in `Test-GetSystemInfo.ps1` so the PowerShell twin never assumes Windows-only cmdlets and falls through to POSIX commands elsewhere — the highest-risk code in the feature (research §R6)
- [ ] T037 [US3] Implement the address probe in both system-info twins with the fixed documented fallback order, writing `snapshot_address` rows (FR-008, FR-018d)
- [ ] T038 [US3] Implement the listening-port probe in both system-info twins with the fixed fallback order and `-SkipPorts` / `--skip-ports`, writing `snapshot_port` rows (FR-008, FR-018d)
- [ ] T039 [US3] Implement graceful probe degradation in both system-info twins: record `NULL`, warn on stderr, never abort the snapshot, and never conflate `NULL` with a legitimate zero (FR-009)
- [ ] T040 [P] [US3] Implement the `hosts` and `listeners` queries in both reader twins, including the port-set diff against the previous snapshot (FR-013)

**Checkpoint**: quickstart scenario 7 passes on the host platform.

---

## Phase 6: User Story 4 — House standards on a fresh clone (P3)

**Goal**: a fresh clone arrives with spec-kit commands and house standards, and without
credentials.

**Independent test**: clone to a clean directory; skills present, no credential file.

- [ ] T041 [US4] Change `.gitignore` in one edit: `.claude/*` plus `!.claude/skills/` (exclude-then-narrowly-admit, not a denylist), and `test/scripts/.bin/` so a locally installed `sqlite3` is never tracked (FR-026, FR-026a, FR-027). **PINNED ARTIFACT — requires the dated CHANGELOG decision in T056**
- [ ] T042 [P] [US4] Vendor `shruggie-powershell` into `.claude/skills/` with a provenance note recording source and date (FR-028, research §R9)
- [ ] T043 [P] [US4] Vendor `shruggie-markdown` into `.claude/skills/` with a provenance note (FR-028)
- [ ] T044 [P] [US4] Vendor `shruggie-speckit` into `.claude/skills/` with a provenance note (FR-028)
- [ ] T045 [P] [US4] Vendor `gh-fix-ci` into `.claude/skills/` with a provenance note (FR-028)
- [ ] T046 [US4] Author the new project-native `.claude/skills/go-schedule-verify/SKILL.md`: the six CI-parity commands, the foreground-only rule, `-coverpkg` cross-package coverage semantics, and both local-environment traps (FR-028)
- [ ] T047 [US4] Run the repository-hygiene check: `git status --porcelain .claude` must show skills only, no settings or credential file swept in by the negation; **then clone the repository to a temporary directory and confirm `.claude/skills/` arrives populated with no post-clone step** — SC-006 says "fresh clone", and inspecting the working tree is not the same claim (FR-026b, SC-006)

**Checkpoint**: quickstart scenario 9 passes.

---

## Phase 7: Polish & Cross-Cutting

- [ ] T048 Write `docs/test-scripts.md` to the vendored Markdown house style: prerequisites, quickstart, per-script reference for both twins, schemas, query catalog, exit-code table, troubleshooting (FR-023)
- [ ] T049 Add the worked end-to-end recipes to `docs/test-scripts.md` — on-time firing, restart survival, downtime catch-up, each overlap policy, failure reporting (FR-024)
- [ ] T050 Add the POSIX-shell conventions section to `docs/test-scripts.md`, mirroring the PowerShell contract point for point, since no separate shell standard exists (FR-025)
- [ ] T051 Document in `docs/test-scripts.md`: the no-retention policy with approximate growth rates (FR-022a), the duration-bound overrun (FR-004a), and the `run_as` different-user directory consequence (research §R2)
- [ ] T052 [P] Add a one-line pointer to `docs/test-scripts.md` from `README.md`
- [ ] T053 [P] Add `-Help` / `--help` usage output to all six scripts (FR-022)
- [ ] T054 Run `gofmt`, `go vet`, `golangci-lint`, the test suite, and `scripts/coverage-gate.sh` in the foreground, watched to completion. Report the `-race` result honestly — this machine has no C compiler, so it cannot run locally and is deferred to CI (research §R8)
- [ ] T055 Run `shellcheck` against all `.sh` files and the ShruggieTech compliance checker against all `.ps1` files; report honestly if a linter is unavailable rather than reporting a pass
- [ ] T056 Update `CHANGELOG.md`: Added entries, plus the **dated Decisions entry for the `.gitignore` pinned-artifact change** (FR-029) and for the drift-derivation approach
- [ ] T057 Write the **twin-parity test** in `test/integration/testscripts_test.go`: run both twins of each script against one database and assert the recorded rows are equivalent field-for-field apart from `script_flavor`, and that every documented option exists on both sides. Research §R6 names twin divergence the likeliest defect in this feature and nothing else asserts against it (SC-004, FR-015)
- [ ] T058 Walk the whole of [quickstart.md](quickstart.md) end to end against a running daemon on this machine, recording the elapsed and hands-on time, and confirm each scenario's stated expectation. SC-001 and SC-003 are claims about a maintainer's experience and are only verifiable by having the experience (SC-001, SC-003)
- [ ] T059 Confirm and record the two honesty obligations before the halt: that `-race` did not run here (no C compiler) and is deferred to CI, and that any unavailable linter is reported as not-run rather than as a pass. A skipped gate reported as green is the one failure mode that silently invalidates every other claim in the report (research §R8, FR-021b)

---

## Dependencies

```text
Setup (T001–T003)
   └─> Foundational (T004–T017)   ← blocks everything
          ├─> US1 (T018–T027)     ← MVP
          ├─> US2 (T028–T032)     depends on US1's heartbeat + reader
          ├─> US3 (T033–T040)     independent of US1/US2
          └─> US4 (T041–T047)     fully independent of all script work
                 └─> Polish (T048–T056)
```

- **US4 shares no code with US1–US3.** It could be done first, last, or by a different person.
- **US3 is independent of US1 and US2** beyond the shared `lib/`.
- **US2 genuinely depends on US1** — it queries the beats US1 records.

## Parallel Opportunities

- **Twin pairs**: every `[P]`-marked `.sh` task pairs with the `.ps1` task above it.
- **Foundational**: T007, T009, T011, T013, T015 all parallel their PowerShell counterparts.
- **US3 probes**: T037 and T038 touch different code paths.
- **US4 vendoring**: T042–T045 are four independent directory copies.
- **Docs**: T052 and T053 are independent of the `docs/test-scripts.md` body.

## Implementation Strategy

**MVP = Phase 1 + Phase 2 + Phase 3 (US1).** That alone delivers the feature's core promise:
a maintainer can prove the scheduler fires on time and quantify the drift. Everything after
it is additive.

Then US3 (different execution profile, catches cross-platform bugs), then US2 (the harder
guarantees), then US4 (no runtime behavior, no risk), then Polish.

**Task count**: 58 — Setup 2, Foundational 14, US1 10, US2 5, US3 8, US4 7, Polish 12.
(T002 merged into T041; T057–T059 added by the analyze gate.)

## Analyze Gate Record (2026-07-23)

The gate passed with **zero CRITICAL findings**, so no early halt. Eight findings were
raised; the five that mattered are resolved above rather than deferred:

- **F1 (HIGH) — silent deviation from an approved artifact.** The POSIX twins had been
  renamed to kebab-case (`test-heartbeat.sh`) on POSIX-idiom grounds. But the operator's
  request and the approved plan both name them `Test-Heartbeat.sh`, and quietly improving an
  approved decision is not the agent's call to make. Reverted across plan, tasks, and
  contracts. The paired naming also makes the twins visually adjacent in a directory listing,
  which suits a design whose central risk is the two halves drifting apart.
- **C1 (HIGH) — twin parity had no test.** Research §R6 identifies twin divergence as the
  likeliest defect in the feature, and no task asserted against it: the tests exercised one
  twin and assumed the other. Added T057.
- **C2 (HIGH) — SC-001 and SC-003 had no owner.** Both are claims about what a maintainer
  can do in a given time, and neither is provable from unit tests. Added T058 to actually
  walk the quickstart and record the timings.
- **C3 (MEDIUM) — "fresh clone" was being checked by looking at the working tree**, which is
  a different claim. T047 now performs the clone.
- **I1 (MEDIUM) — two tasks edited `.gitignore`** in different phases. One pinned artifact,
  one edit, one changelog record. Merged into T041.

T059 was added on top of the findings: the verification report's honesty obligations are
themselves easy to lose between running the gates and writing the halt breakdown, and a
skipped gate reported as green would silently invalidate every other claim in that report.

I2 (test-after-implementation within US1) and L1 (`--help` scheduled in Polish) were assessed
and accepted — constitution principle II permits tests written *alongside* the code, and
neither blocks anything.
