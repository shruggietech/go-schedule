# Tasks: General Five-Field Cron Breadth

**Input**: Design documents from `specs/026-cron-expression-breadth/` **Tests**: Required by the specification and constitution; behavioral tasks follow test-first order.

## Phase 1: Setup

- [x] T001 Record the S026 baseline, active issue scope, and representative expression matrix in `specs/026-cron-expression-breadth/verification.md`
- [x] T002 [P] Add compiler and recurrence benchmark baselines for simple, broad, and sparse expressions in `internal/cron/composite_bench_test.go` and `internal/engine/engine_bench_test.go`

## Phase 2: Foundational

- [x] T003 Add red normalized-field, range-step, wildcard-step, overlap, name, and Sunday-alias tests in `internal/cron/composite_test.go`
- [x] T004 Add red direct compilation and durable schedule-shape tests in `internal/cron/compile_test.go`
- [x] T005 Implement reusable field restriction, normalization, and rendering helpers in `internal/cron/cron.go`
- [x] T006 Implement standard-field cron compilation into a constrained durable recurrence in `internal/cron/compile.go`

## Phase 3: User Story 1 - Run real-world composite cron schedules (P1)

**Goal**: One task executes standard five-field lists, ranges, steps, names, and safe conjunctions exactly.

**Independent Test**: Compare at least twenty expressions with hand-checked occurrences across ordinary, boundary, leap, and DST dates.

- [x] T007 [US1] Add red shared-input compilation and source-identity tests for composite expressions in `internal/scheduleinput/composite_test.go`
- [x] T008 [US1] Route standard cron input directly through the cron compiler while retaining focused modifier compatibility in `internal/scheduleinput/input.go`
- [x] T009 [US1] Add red next-run matrices for time sets, month/date sets, weekday sets, strict anchors, and DST transitions in `internal/schedule/composite_cron_test.go`
- [x] T010 [US1] Add policy-aware multi-date red tests for skip, last-valid, next-valid, and collision suppression in `internal/schedule/composite_cron_test.go`
- [x] T011 [US1] Extend policy-aware date-set recurrence resolution with wall-time normalization and anchor guards in `internal/schedule/missingdate.go` and `internal/schedule/recur.go`
- [x] T012 [US1] Add restart and catch-up red tests for composite schedules in `internal/store/composite_cron_test.go` and `internal/catchup/catchup_test.go`
- [x] T013 [US1] Complete restart and catch-up compatibility using the existing schedule persistence path in `internal/store/store.go` and `internal/catchup/catchup.go`

## Phase 4: User Story 2 - Inspect and import composite schedules safely (P1)

**Goal**: Explain, convert, preview, and import share an exact field-complete description and recurrence.

**Independent Test**: Dry-run and import a mixed crontab, then compare classification, descriptions, source identity, and upcoming runs.

- [x] T014 [US2] Add red concise-phrase compatibility and broad-description tests in `internal/cron/phrase_test.go` and `internal/cron/convert_test.go`
- [x] T015 [US2] Implement deterministic broad field descriptions while preserving existing concise phrases in `internal/cron/phrase.go` and `internal/cron/convert.go`
- [x] T016 [US2] Add red crontab classification and dry-run parity tests in `internal/cron/crontab_test.go` and `internal/cli/cron_test.go`
- [x] T017 [US2] Propagate composite descriptions and typed source through crontab conversion and import in `internal/cron/crontab.go` and `internal/cli/cron.go`

## Phase 5: User Story 3 - Edit, restart, and export without semantic drift (P1)

**Goal**: All task interfaces retain source and export a recurrence-equivalent canonical expression after edits and restart.

**Independent Test**: Create, edit, reload, export, and re-import composite tasks through CLI, API, and desktop coverage with identical runs.

- [x] T018 [US3] Add red canonical field compression, source-independent export, policy refusal, and round-trip tests in `internal/cron/composite_export_test.go`
- [x] T019 [US3] Implement constrained composite recurrence recognition and canonical five-field export in `internal/cron/export.go`
- [x] T020 [US3] Add red API preview/create/update/non-mutation and reload tests in `internal/api/server/composite_cron_test.go`
- [x] T021 [US3] Propagate composite schedule summaries and source identity through shared API task boundaries in `internal/api/server/tasks.go`
- [x] T022 [US3] Add red CLI task and JSON/text conversion/export tests in `internal/cli/composite_cron_test.go`
- [x] T023 [US3] Add red desktop preview, syntax selection, prefill, edit, and refusal tests in `gui/composite_cron_test.go`
- [x] T024 [US3] Preserve composite cron entry and preview behavior through the existing desktop editor path in `gui/editor.go`

## Phase 6: User Story 4 - Preserve honest boundaries (P2)

**Goal**: Unrepresentable cron semantics remain precise, stable, and non-mutating.

**Independent Test**: Exercise DOM/DOW OR, Quartz, boot, modifier-composite, malformed, and update-failure cases across parser and interfaces.

- [x] T025 [US4] Add regression tests for every retained parser, compiler, export, and task-boundary refusal in `internal/cron/composite_refusal_test.go` and `internal/api/server/composite_cron_test.go`
- [x] T026 [US4] Tighten field-specific errors and named semantic refusals without fallback in `internal/cron/cron.go`, `internal/cron/compile.go`, and `internal/scheduleinput/input.go`

## Phase 7: Polish and Cross-Cutting Concerns

- [x] T027 Update cron fidelity, conversion semantics, examples, and task authoring documentation in `docs/cron.md`, `docs/cli.md`, and `README.md`
- [x] T028 Update issue #22 inventory language, Unreleased changes, and the dated architectural decision in `CHANGELOG.md`
- [x] T029 Run focused suites and benchmarks, record before/after evidence in `specs/026-cron-expression-breadth/verification.md`, and resolve regressions
- [x] T030 Run all eight canonical gates in the foreground and record exact results in `specs/026-cron-expression-breadth/verification.md`
- [x] T031 Audit task completion, spec acceptance, `git diff --check`, UTF-8 without BOM, and mojibake across every changed file in `specs/026-cron-expression-breadth/verification.md`

## Dependencies and Execution Order

- Phase 1 establishes evidence; Phase 2 blocks every user story.
- US1 establishes execution and lifecycle behavior before US2 descriptions or US3 export can claim fidelity.
- US2 and US3 depend on the compiler but can otherwise proceed independently after US1 recurrence semantics stabilize.
- US4 regression coverage follows the accepted path so boundaries can be tested against final routing.
- Documentation and final verification follow all user stories.

## Parallel Opportunities

- T002 can run independently from T001.
- Within US3, API, CLI, and desktop red tests touch separate files after T019 defines export behavior.
- Documentation can be drafted after behavior stabilizes while focused benchmarks run, but final wording waits for verified outcomes.

## Implementation Strategy

1. Make parser normalization and direct compilation pass first.
2. Prove exact next-run and policy behavior before exposing the new breadth through interfaces.
3. Add descriptions and import, then source-independent export and lifecycle coverage.
4. Preserve explicit refusals, document the architectural deviation, and finish with full CI parity.
