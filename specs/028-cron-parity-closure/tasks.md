# Tasks: Cron Parity Closure

**Input**: Design documents from `specs/028-cron-parity-closure/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Required and written first for every behavior-bearing phase.

## Phase 1: Baseline and Setup

- [x] T001 Record focused suite and recurrence benchmark baselines in `specs/028-cron-parity-closure/research.md`
- [x] T002 Confirm no dependency, pinned artifact, or ignore change is required and update the Spec-Kit context pointer in `CLAUDE.md`

## Phase 2: Foundational Task Standard Input

**Goal**: Give cron percent payloads one durable, cross-platform execution representation.

- [x] T003 [P] Add red task stdin create, retain, replace, and clear API tests in `internal/api/server/tasks_test.go` and `internal/api/server/update_test.go`
- [x] T004 [P] Add red schema-v8 migration and CRUD/reopen tests in `internal/store/migration_v8_test.go` and `internal/store/store_test.go`
- [x] T005 [P] Add red exact stdin and environment override execution tests in `internal/executor/executor_test.go`
- [x] T006 Add `Stdin` to `domain.Task`, API create, and pointer-based update in `internal/domain/domain.go`, `internal/api/server/tasks.go`, and `internal/api/server/update.go`
- [x] T007 Add the schema-v8 empty default and CRUD propagation in `internal/store/store.go` and `internal/store/crud.go`
- [x] T008 Supply non-empty task stdin verbatim in `internal/executor/executor.go`

## Phase 3: User Story 1 - Operational crontab fidelity (Priority: P1)

**Goal**: Ordered crontab context, shell behavior, percent input, and system identity survive import.

**Independent Test**: Parse and import one mixed fixture, inspect every task request, and execute the resulting stdin payload.

- [x] T009 [US1] Add red scanner tests for assignment quoting, reassignment snapshots, `CRON_TZ` versus `TZ`, `SHELL`, mail warnings, and invalid lines in `internal/cron/crontab_test.go`
- [x] T010 [US1] Add red percent parser tests for escaped percent, quoted percent, empty stdin, newlines, and preserved backslashes in `internal/cron/crontab_test.go`
- [x] T011 [US1] Add red Unix/Quartz and user/system layout matrix tests in `internal/cron/crontab_test.go`
- [x] T012 [US1] Implement ordered scan context, assignment parsing, immutable snapshots, and per-line timezone/environment in `internal/cron/crontab.go`
- [x] T013 [US1] Implement shell command preservation and percent-to-stdin transformation in `internal/cron/crontab.go`
- [x] T014 [US1] Implement explicit scan options for timing dialect and system user fields in `internal/cron/crontab.go`
- [x] T015 [US1] Add red CLI request/report tests for effective timezone, env, shell, stdin, run-as, override, and layout flags in `internal/cli/cron_test.go`
- [x] T016 [US1] Expose `--dialect` and `--system`, pass exact job context to API creation, and replace obsolete summary claims in `internal/cli/cron.go`
- [x] T017 [US1] Add API import lifecycle coverage for timezone validation, Unix run-as, and non-mutating per-line failures in `internal/api/server/`

## Phase 4: User Story 2 - Seconds and Quartz timing (Priority: P1)

**Goal**: Supported six-field schedules use the same durable evaluator and round-trip without shifting seconds or weekdays.

**Independent Test**: Enumerate and round-trip seconds, value-step, Quartz `?`, and numeric weekday fixtures through every timing boundary.

- [x] T018 [US2] Add red parser tests for five/six fields, second sets, value-start steps, Unix/Quartz weekday numbers, `?`, year fields, and modifier boundaries in `internal/cron/seconds_test.go`
- [x] T019 [US2] Refactor cron field identity and implement dialect-aware parsing with synthetic second zero in `internal/cron/cron.go`
- [x] T020 [US2] Add red compile and description tests for second boundaries, strict-after-anchor behavior, readable second sets, and focused-selector refusals in `internal/cron/compile_test.go` and `internal/cron/phrase_test.go`
- [x] T021 [US2] Compile `BYSECOND` directly and make deterministic descriptions seconds-aware in `internal/cron/compile.go` and `internal/cron/phrase.go`
- [x] T022 [US2] Add red conversion/detection tests for six fields and retained five-field behavior in `internal/cron/convert_test.go` and `internal/scheduleinput/input_test.go`
- [x] T023 [US2] Detect six-field input structurally and update local conversion boundaries in `internal/cron/convert.go`
- [x] T024 [US2] Add red export tests for canonical five fields, six-field seconds, native divisible secondly intervals, Quartz weekday numbers, and operational-context refusals in `internal/cron/export_test.go`
- [x] T025 [US2] Generalize recurrence export for exact seconds while preserving five-field output in `internal/cron/export.go` and `internal/cli/cron.go`
- [x] T026 [US2] Add red missing-date and calendar-adjustment tests with non-zero and multiple seconds in `internal/schedule/missingdate_test.go`
- [x] T027 [US2] Remove second-zero assumptions from supported composite date resolution in `internal/schedule/missingdate.go`

## Phase 5: User Stories 2 and 4 - Cross-surface lifecycle (Priority: P1/P2)

**Goal**: Six-field source and exact occurrences survive authoring, persistence, restart, catch-up, GUI, and dispatch.

- [x] T028 [US2] Add API preview/create/update/reload and source-identity tests in `internal/api/server/composite_cron_test.go`
- [x] T029 [P] [US2] Add GUI preview/create/prefill tests for six-field source in `gui/composite_cron_test.go`
- [x] T030 [P] [US4] Add store restart and catch-up tests for seconds schedules in `internal/store/composite_cron_test.go` and `internal/catchup/catchup_test.go`
- [x] T031 [US2] Route any required source/help adjustments through API, CLI, and GUI without adding a second evaluator
- [x] T032 [US4] Add and compare seconds compile/next-run benchmarks in `internal/cron/composite_bench_test.go`, `internal/engine/engine_bench_test.go`, and `specs/028-cron-parity-closure/verification.md`

## Phase 6: User Story 3 - Audit closure and documentation (Priority: P1)

**Goal**: Every issue #22 row has one honest decision and all public claims match it.

- [x] T033 [US3] Replace the cron fidelity table with explicit A1-A12 and B1-B9 dispositions and update import/export examples in `docs/cron.md`
- [x] T034 [US3] Scope universal cron claims and align CLI, GUI field, and docs index wording in `README.md`, `docs/cli.md`, `docs/gui-fields.md`, and `docs/README.md`
- [x] T035 [US3] Record S028, schema v8, behavioral deviations, and issue #22 closure intent in `CHANGELOG.md`

## Phase 7: Verification and Local Delivery

- [x] T036 Run focused tests and quickstart scenarios and record results in `specs/028-cron-parity-closure/verification.md`
- [x] T037 Run all eight canonical gates in the foreground and record exact results in `specs/028-cron-parity-closure/verification.md`
- [x] T038 Audit all 32 requirements, 39 tasks, both checklists, issue rows, `git diff --check`, UTF-8 without BOM, and mojibake in `specs/028-cron-parity-closure/verification.md`
- [x] T039 Commit the verified slice locally with `Closes #22` reserved for the eventual pull request, then halt before publication

## Dependencies and Execution Order

- Baseline and Spec-Kit gates precede implementation.
- Durable stdin blocks percent import creation and execution.
- Scanner context and explicit layouts stabilize before CLI wiring.
- Cron field identity and parsing stabilize before compiler, description, export, and missing-date work.
- Cross-surface lifecycle follows the core recurrence contract.
- Documentation and final verification follow all behavior-bearing work.

## Parallel Opportunities

- T003, T004, and T005 are independent red-test surfaces.
- Scanner assignment, percent, and layout tests can be drafted independently before their shared implementation.
- GUI, store/catch-up, and documentation tests touch separate surfaces after the core parser stabilizes.

## Implementation Strategy

1. Land the smallest durable execution primitive, task stdin, with migration and lifecycle proof.
2. Replace lossy crontab scanning with explicit ordered context and deterministic layouts.
3. Extend the cron field model to seconds and Quartz semantics, then compile into RRULE.
4. Prove every authoring, persistence, execution, and export boundary.
5. Publish the complete audit matrix, run every gate, commit locally, and stop once before the branch leaves the machine.
