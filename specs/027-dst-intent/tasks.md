# Tasks: Explicit DST Scheduling Intent

**Input**: Design documents from `specs/027-dst-intent/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Required and written first for every behavior-bearing phase.

## Phase 1: Baseline and Setup

- [x] T001 Record current New York spring/fall wall-clock outputs and next-run benchmark baselines in `specs/027-dst-intent/verification.md`
- [x] T002 Verify existing Go and universal ignore patterns require no changes in `.gitignore`

## Phase 2: Foundational Policy and Resolution

- [x] T003 [P] Add red enum/default tests for the task scheduling policy set in `internal/domain/domain_test.go`
- [x] T004 [P] Add red zero/one/two-instant resolver tests across New York and Lord Howe transitions in `internal/timezone/dst_policy_test.go`
- [x] T005 Implement time basis, gap, overlap, and effective SchedulePolicy types in `internal/domain/domain.go`
- [x] T006 Implement policy-aware wall-time resolution while preserving the compatibility wrapper in `internal/timezone/timezone.go`

## Phase 3: User Story 1 - Explicit recurrence basis (Priority: P1)

**Goal**: Wall-clock, elapsed, and UTC schedules have distinct, exact semantics.

**Independent Test**: Enumerate equivalent schedules across both 2026 New York transitions and reject calendar-selected elapsed shapes.

- [x] T007 [US1] Add red basis matrices, fixed-duration validation, and duplicate-suppression tests in `internal/schedule/dst_policy_test.go`
- [x] T008 [US1] Refactor recurrence entry points to accept SchedulePolicy and implement compatibility defaults in `internal/schedule/recur.go`
- [x] T009 [US1] Implement constant-time fixed-duration validation and elapsed next-run arithmetic in `internal/schedule/recur.go`
- [x] T010 [US1] Implement UTC-basis recurrence and missing-date evaluation in `internal/schedule/recur.go` and `internal/schedule/missingdate.go`
- [x] T011 [US1] Implement floating wall-clock recurrence intent and policy-aware candidate selection in `internal/schedule/recur.go`

## Phase 4: User Story 2 - Explicit transition behavior (Priority: P1)

**Goal**: Gap skip/next-valid and overlap first/both/last compose with all wall-clock recurrence paths.

**Independent Test**: Pin exact UTC occurrences at real spring and fall transitions, including both folds and missing-date/calendar-adjusted recurrences.

- [x] T012 [US2] Add red composition and between-fold cursor tests in `internal/schedule/dst_policy_test.go`
- [x] T013 [US2] Route calendar adjustment, composite date-set, and missing-date candidates through the shared resolver in `internal/schedule/missingdate.go`
- [x] T014 [US2] Preserve overlap-both ordering and suppress next-valid collisions in `internal/schedule/recur.go`
- [x] T015 [US2] Update catch-up tests for policy parity and catch-up-one bounds in `internal/catchup/catchup_test.go`
- [x] T016 [US2] Route SchedulePolicy through catch-up and live dispatch in `internal/catchup/catchup.go` and `internal/engine/engine.go`

## Phase 5: User Story 3 - Durable cross-surface policy (Priority: P1)

**Goal**: Policy values persist and behave identically through storage, API, CLI, GUI, detail, and calendar.

**Independent Test**: Create/edit/restart tasks through every boundary and compare next runs and policy fields.

- [x] T017 [US3] Add red v7 forward-migration and CRUD round-trip tests in `internal/store/migration_v7_test.go`
- [x] T018 [US3] Add schema v7 defaults and persist all policy fields in `internal/store/store.go` and `internal/store/crud.go`
- [x] T019 [US3] Add red API create/update/preview/detail/calendar validation and parity tests in `internal/api/server/dst_policy_test.go`
- [x] T020 [US3] Expose and validate policy fields on API create, update, preview, detail, and calendar in `internal/api/server/tasks.go`, `internal/api/server/update.go`, and `internal/api/server/calendar.go`
- [x] T021 [US3] Add red CLI flag, request, show-output, and invalid-value tests in `internal/cli/task_test.go`
- [x] T022 [US3] Expose `--time-basis`, `--dst-gap`, and `--dst-overlap` on add/edit and show effective values in `internal/cli/task.go`
- [x] T023 [US3] Add red GUI default, prefill, preview, dirty-state, and submission tests in `gui/editor_test.go` and `gui/editor_prefill_test.go`
- [x] T024 [US3] Add friendly Advanced Settings selectors and wire values through preview/create/edit in `gui/editor.go` and `gui/editor_data.go`

## Phase 6: User Story 4 - Compatibility and bounded execution (Priority: P2)

**Goal**: Existing behavior remains compatible, invalid input is non-mutating, and evaluation stays inside performance budgets.

**Independent Test**: Run legacy regression suites, invalid combinations, migration fixtures, and before/after benchmarks.

- [x] T025 [US4] Preserve missing-date compatibility wrappers, route all production occurrence callers through SchedulePolicy, and add omitted-value assertions across `internal/`
- [x] T026 [US4] Add policy-aware next-run benchmarks and compare against T001 in `internal/engine/engine_bench_test.go` and `specs/027-dst-intent/verification.md`
- [x] T027 [US4] Exercise field-specific non-mutating validation across create, preview, and update in `internal/api/server/dst_policy_test.go`

## Phase 7: Documentation and Verification

- [x] T028 Update task authoring, GUI fields, scheduling semantics, and examples in `docs/cli.md`, `docs/gui-fields.md`, `README.md`, and applicable scheduling docs
- [x] T029 Record S027 behavior and the dated architectural deviation in `CHANGELOG.md`
- [x] T030 Run focused suites and the quickstart scenarios and record results in `specs/027-dst-intent/verification.md`
- [x] T031 Run all eight canonical gates in the foreground and record exact results in `specs/027-dst-intent/verification.md`
- [x] T032 Audit task completion, acceptance traceability, `git diff --check`, UTF-8 without BOM, and mojibake across every changed file in `specs/027-dst-intent/verification.md`

## Dependencies and Execution Order

- Baseline evidence precedes implementation.
- Foundational policy and resolver work blocks both recurrence stories.
- US1 establishes basis routing before US2 composes transition choices.
- US3 begins after the evaluator contract stabilizes, then US4 verifies compatibility and bounds.
- Documentation and final verification follow all behavior-bearing work.

## Parallel Opportunities

- T003 and T004 touch separate packages and can be written independently.
- Storage, CLI, and GUI red tests touch separate files after the evaluator contract stabilizes.
- Documentation drafting can begin after cross-surface behavior stabilizes, while focused benchmarks run.

## Implementation Strategy

1. Pin the current behavior and define typed defaults.
2. Make basis and transition matrices pass in the core evaluator.
3. Route the stable contract through catch-up, engine, persistence, API, CLI, GUI, and calendar.
4. Prove migration, lifecycle, failure, and performance behavior.
5. Document the deliberate premise correction, run all gates, and close issue #8 through the eventual PR.
