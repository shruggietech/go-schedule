---

description: "Task list for 008-cron-interop"
---

# Tasks: Cron interoperability and calendar-anomaly policy

**Input**: Design documents from `/specs/008-cron-interop/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/cli.md

**Tests**: REQUIRED. Constitution principle II (NON-NEGOTIABLE) — tests are
written alongside or before the code they verify, and this feature touches two
safety-critical surfaces (timezone/DST resolution, forward-only migrations).

**Organization**: grouped by user story so each is independently implementable
and testable. US2 (grammar + policy) precedes US1 (import) despite both being P1,
because US1's conversion has no target representation without US2's grammar —
stated as a dependency in spec §FR-015–§FR-017, not a priority inversion.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: can run in parallel (different files, no incomplete dependencies)
- **[Story]**: the user story this task serves
- Exact file paths are given in every task

## Phase 1: Setup

**Purpose**: nothing to scaffold — the module, test harness, and CI gates all
exist. This phase records the one shared decision the rest depends on.

- [X] T001 Add the `MissingDatePolicy` enum (`skip`, `last_valid`, `next_valid`) and its doc comment to `internal/domain/domain.go`, beside `OverlapPolicy` and `CatchupPolicy`
- [X] T002 Add the `MissingDatePolicy` field to `domain.Task` in `internal/domain/domain.go` with the `missing_date_policy` JSON tag, placed after `CatchupPolicy`

---

## Phase 2: Foundational (blocking prerequisites)

**Purpose**: persistence and signature changes that every later phase compiles
against. No user story can land until these do.

**⚠️ CRITICAL**: complete before any user-story phase.

- [X] T003 Add migration v5 to `internal/store/store.go`: `ALTER TABLE tasks ADD COLUMN missing_date_policy TEXT NOT NULL DEFAULT 'skip';` with a comment in the v4 style stating it is additive with a total default and rewrites nothing
- [X] T004 Write `internal/store/migration_v5_test.go` asserting a v4-era database opens, every existing task reads `skip`, and no stored `rrule`, `anchor`, or `run_at` value changes (FR-026, SC-006)
- [X] T005 Carry the new column through the task insert and scan in `internal/store/crud.go`, defaulting an empty stored value to `skip` on read
- [X] T006 Change `NextRun` and `UpcomingRuns` in `internal/schedule/recur.go` to take `policy domain.MissingDatePolicy`, and update the six call sites: `internal/engine/engine.go` (two), `internal/catchup/catchup.go` (via a new `Evaluate` parameter), `internal/api/server/calendar.go`, `internal/api/server/tasks.go` (two)
- [X] T007 Update `internal/engine/engine_bench_test.go` for the new `NextRun` signature and confirm `BenchmarkNextRun` still compiles and runs

**Checkpoint**: the tree builds and every existing test passes with the policy
threaded through but always `skip` — i.e. behavior is provably unchanged.

---

## Phase 3: User Story 2 — Schedule by calendar date, with a stated policy (Priority: P1) 🎯 MVP

**Goal**: by-date monthly, yearly, and month/year-interval phrases parse; each of
the three policies produces the stated run times; descriptions stop lying.

**Independent test**: create a task under each policy and compare its next runs
across a year containing a short month and a non-leap February; confirm the
description names the policy.

### Tests for US2

- [X] T008 [P] [US2] Add by-date monthly, yearly, and `every 12 months` parse cases to `internal/schedule/parse_test.go`, including the optional `at <time>` tail, the rejection of an out-of-range day, and a round-trip assertion that each new phrase is retained in `Schedule.Expression` and that its summary describes the same rule (FR-018)
- [X] T009 [P] [US2] Add policy resolution tests to `internal/schedule/recur_test.go` pinned to real dates: the 31st across 2026, 29 February in the non-leap 2027, and the 30th in February, each under all three policies; include an assertion that a `next_valid` roll-forward neither displaces nor duplicates the following period's own occurrence (SC-004, FR-019a)
- [X] T010 [P] [US2] Add a test asserting `5th friday monthly` under `last_valid` resolves to the last Friday of a four-Friday month, and that its summary names the policy (US2 scenario 5, FR-023)
- [X] T011 [P] [US2] Add a test asserting the policy is inert for interval, weekday, and dayset rules — identical run times under all three settings (FR-024)
- [X] T012 [P] [US2] Add a test asserting a policy-resolved date crossing the 2026-03-08 `America/New_York` transition still receives the existing DST normalization (FR-025)

### Implementation for US2

- [X] T013 [US2] Add the by-date monthly parse arm and its regex to `internal/schedule/parse.go`, producing `FREQ=MONTHLY;BYMONTHDAY=<n>` and a summary in the established phrasing
- [X] T014 [US2] Add the yearly parse arm to `internal/schedule/parse.go`, producing `FREQ=YEARLY;BYMONTH=<m>;BYMONTHDAY=<d>`, with month-name parsing
- [X] T015 [US2] Extend `unitToFreq` in `internal/schedule/parse.go` with month and year units so `every 12 months` and `every year` fall out of the existing interval arm
- [X] T016 [US2] Create `internal/schedule/missingdate.go`: detect date-bearing rules (`BYMONTHDAY`, ordinal `BYDAY`), and implement the bounded period walk resolving `last_valid` and `next_valid`, returning "no further run" rather than spinning at the cap
- [X] T017 [US2] Wire the resolution into `nextRecurring` in `internal/schedule/recur.go` *before* the existing `timezone.WallTime` normalization, leaving the `skip` path byte-identical to today's code
- [X] T018 [US2] Add the policy clause to `HumanSummary` for date-bearing rules in `internal/schedule/parse.go` (and wherever a summary is regenerated), so no description asserts "every month" for a rule that skips months (FR-023, SC-005)

**Checkpoint**: US2 is independently demonstrable through `gosched task add`
once Phase 5's flag lands; until then it is demonstrable through its tests.

---

## Phase 4: User Story 1 — Move an existing crontab across (Priority: P1)

**Goal**: preview and import a crontab, with every line either translated or
refused by name, and the fidelity stated in the summary.

**Independent test**: preview a sample crontab and confirm nothing is created;
import it and confirm the tasks match the preview.

### Tests for US1

- [X] T019 [P] [US1] Write `internal/cron/cron_test.go`: field parsing table covering wildcards, lists, ranges, names, and dividing steps, plus refusals for six-field input, `@reboot`, `L`, `W`, `#`, non-dividing steps, and day-of-month + day-of-week (FR-001, FR-002, FR-003b)
- [X] T020 [P] [US1] Write `internal/cron/phrase_test.go` seeded from the `cronparity_test.go` corpus, asserting each expression maps to the expected phrase and that the phrase parses to a schedule with matching run times (FR-003a)
- [X] T021 [P] [US1] Write `internal/cli/cron_test.go` covering the import report: comments and blanks skipped, `MAILTO` and variable assignments warned, counts correct, preview creating nothing, and a declined line exiting 0 (FR-005, FR-007, FR-010, FR-010a)
- [X] T021a [P] [US1] Add a test asserting the phrase a preview shows for a line is the phrase the task created from that same line reports back — the guarantee that makes the preview trustworthy (SC-002a, research D2)
- [X] T021b [P] [US1] Add a test asserting cron syntax is still rejected wherever a human phrase is accepted (`schedule.Parse("0 9 * * 1-5")` errors) and that no GUI input path accepts a cron expression (FR-014)

### Implementation for US1

- [X] T022 [P] [US1] Create `internal/cron/cron.go` with the `Spec`, `Unsupported`, `Line`, and `Report` types and the five-field parser, refusing by name everything listed in FR-002
- [X] T023 [US1] Create `internal/cron/phrase.go` mapping a `Spec` to the phrase a user would have typed, refusing anything with no phrase — never emitting a recurrence directly (FR-003a, research D2)
- [X] T024 [US1] Add crontab-file scanning to `internal/cron`: comments, blanks, variable assignments, and schedule lines, splitting each line into timing and command payload (FR-006, FR-007)
- [X] T025 [US1] Create `internal/cli/cron.go` with `newCronCmd()` and the `import` subcommand per `contracts/cli.md`, including `--file`, `--dry-run`, `--timezone`, `--group`, `--count`, the per-line report, and the mandatory fidelity summary (FR-004, FR-005, FR-008, FR-009, FR-010)
- [X] T026 [US1] Register `newCronCmd()` in `newRoot()` in `internal/cli/cli.go`
- [X] T027 [US1] Implement partial-failure behavior in the import: created tasks remain, the failure is reported alongside the created count, exit 1 (FR-005a)

**Checkpoint**: a crontab can be previewed and imported end to end.

---

## Phase 5: Policy surfacing (serves US2; required by quickstart)

**Goal**: the policy is settable and visible everywhere a task is.

- [X] T028 [P] [US2] Add `missing_date_policy` to `TaskCreateRequest` in `internal/api/server/tasks.go` with validation copying the `overlap_policy` block, defaulting to `skip`
- [X] T029 [US2] Add `missing_date_policy` to `TaskUpdateRequest` in `internal/api/server/update.go`, leaving the stored value untouched when the field is empty, and confirm a schedule replacement does not reset it (FR-024a)
- [X] T030 [P] [US2] Add a server test asserting an edit that changes the phrase preserves the policy, and an edit that changes the policy preserves the phrase, in `internal/api/server/update_test.go`
- [X] T031 [P] [US2] Add `--missing-date` to `task add` and `task edit` in `internal/cli/task.go`, and print the policy in `task show`
- [X] T032 [P] [US2] Add `missingDateChoices`, `missingDateLabel`, and `missingDateValue` to `gui/editor_data.go` using the existing generic `policyChoice[T]`
- [X] T033 [US2] Add the **Missing dates** row to the Advanced Settings form in `gui/editor.go`, prefilled from the task and carried through both submit paths
- [X] T034 [P] [US2] Add a GUI test in `gui/editor_test.go` asserting the selector prefills from an existing task and round-trips on save

---

## Phase 6: User Story 3 — Explain one expression (Priority: P2)

**Goal**: a side-effect-free translator.

**Independent test**: explain several expressions; confirm output correctness and
that nothing was created.

- [X] T035 [P] [US3] Add explain tests to `internal/cli/cron_test.go`: a supported expression, a named refusal exiting 0, and a malformed expression exiting 2 with the offending field named (US3 scenarios 1–3, FR-010a)
- [X] T036 [US3] Add the `explain` subcommand to `internal/cli/cron.go` with `--timezone` and `--count`, per `contracts/cli.md`

---

## Phase 7: User Story 4 — Export back out (Priority: P3)

**Goal**: every task appears as a crontab line or a named refusal.

**Independent test**: export a mixed task set; confirm one line or one refusal
per task.

- [X] T037 [P] [US4] Write `internal/cron/export_test.go`: expressible recurrences produce lines with matching run times; one-off, sub-minute, and disabled tasks produce refusals; an empty set produces the header only (FR-011a, FR-012, FR-012a)
- [X] T038 [P] [US4] Write the round-trip test in `internal/schedule/cronparity_test.go` (or a sibling): cron → phrase → schedule → cron, asserting run times match across the 2026-03-08 DST transition and a month boundary (FR-013, SC-003)
- [X] T039 [US4] Create `internal/cron/export.go` inverting a `domain.Schedule` into a crontab line, or a named refusal
- [X] T040 [US4] Add the `export` subcommand to `internal/cli/cron.go` with `--task`, emitting the header comment, one entry per task, and refusals as `# declined:` comments

---

## Phase 8: User Story 5 — Attribution link (Priority: P3)

- [X] T041 [P] [US5] Link the ShruggieTech attribution to `https://shruggie.tech` in the `README.md` footer
- [X] T042 [P] [US5] Apply the same link wherever the organization is named in `CONTRIBUTING.md`, `SECURITY.md`, and `CODE_OF_CONDUCT.md`

---

## Phase 9: Polish & cross-cutting

- [X] T043 [P] Write `docs/cron.md`: the conversion guide and the supported/declined fidelity table, covering both directions and what cron cannot carry
- [X] T044 [P] Link `docs/cron.md` from `docs/README.md` and `docs/cli.md`, and document the `cron` command group and the `--missing-date` flag in `docs/cli.md`
- [X] T045 [P] Document the Missing dates field in `docs/gui-fields.md`
- [X] T046 [P] Add the cron migration capability to the feature list in `README.md`
- [X] T047 Update `CHANGELOG.md`: an Added entry for the feature, and dated 2026-07-23 decision records for storing the policy on the task rather than the schedule (research D1) and for writing the cron parser in-tree (research D3)
- [X] T047a Run `go test -run '^$' -bench 'BenchmarkNextRun|BenchmarkDispatch' -benchmem ./internal/engine/...` before and after the change and confirm neither regresses by more than 10%; record the figures in the commit message, and record an explicit justification if either does (constitution principle IV — this feature adds work to the next-run hot path)
- [X] T048 Run the six CI-parity gates in the foreground via the `go-schedule-verify` skill and fix anything red
- [X] T049 Walk `specs/008-cron-interop/quickstart.md` end to end against built binaries and record the result

---

## Dependencies

```text
Phase 1 (T001–T002)
      ↓
Phase 2 (T003–T007)  ← blocking for everything
      ↓
Phase 3 US2 grammar + policy (T008–T018)
      ↓
      ├── Phase 4 US1 import (T019–T027)   ← needs the grammar
      ├── Phase 5 surfacing (T028–T034)
      ├── Phase 6 US3 explain (T035–T036)  ← needs US1's converter
      └── Phase 7 US4 export (T037–T040)   ← needs US1's converter
              ↓
Phase 8 US5 docs link (T041–T042)  ← independent of everything; may run any time
      ↓
Phase 9 polish (T043–T049)
```

**Story independence**: US5 depends on nothing. US2 depends only on the
foundational phase. US1, US3, and US4 share `internal/cron` and are ordered by
priority within that shared dependency.

## Parallel execution

- **Phase 3 tests**: T008–T012 are five separate test concerns and can be written
  together before T013–T018.
- **Phase 4**: T019–T021b in parallel; then T022 in parallel with nothing (it is
  the type foundation), then T023–T027 in sequence within `internal/cron` and
  `internal/cli`.
- **Phase 5**: T028, T031, T032, T034 touch four different files and parallelize;
  T029 and T033 must follow T028 and T032 respectively.
- **Phase 8 and Phase 9 docs**: T041–T046 are six independent files.

## Implementation strategy

**MVP** is Phase 1 → Phase 2 → Phase 3 (US2): calendar-date scheduling with a
stated missing-date policy. That alone closes the half of issue #8 this feature
takes on, and it is what everything else stands on.

**Increment 2** adds Phase 4 + Phase 5: the crontab import and the policy's
visible controls — the point at which issue #12's headline claim is true.

**Increment 3** adds Phases 6–7: explain and export, completing #12.

**Increment 4** is Phases 8–9: issue #9 and the documentation set.

Every increment ends green on the six CI-parity gates; a red gate halts rather
than accumulating.
