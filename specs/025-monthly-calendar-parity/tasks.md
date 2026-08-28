# Tasks: Monthly Calendar Cron Parity

**Input**: Design documents from `specs/025-monthly-calendar-parity/`

**Tests**: Required. Behavioral seams are recorded red before production implementation.

## Phase 1: Specification and Design

- [X] T001 Create the S025 review branch, validate baseline suites, and record both complete checklists
- [X] T002 Complete recurrence, cron-dialect, workflow, persistence, GUI, documentation, and test-surface research
- [X] T003 Record the typed adjustment, migration, execution ordering, fidelity boundary, and architecture deviation in design artifacts
- [X] T004 Run initial cross-artifact analysis with zero critical conflicts

## Phase 2: Foundational Domain and Persistence Contract

- [X] T005 Add failing domain serialization and invalid-adjustment contract tests in `internal/domain` and `internal/schedule`
- [X] T006 Add failing schema-v6 forward migration, default-value, reopen, and legacy-row tests in `internal/store/migration_v6_test.go`
- [X] T007 Add failing schedule CRUD, replacement, and restart round-trip tests for the adjustment in `internal/store`
- [X] T008 Record the expected red foundational results in `verification.md`
- [X] T009 Add the documented `CalendarAdjustment` enum and schedule field in `internal/domain/domain.go`
- [X] T010 Add schema version 6 and update all explicit schedule write/read queries in `internal/store`
- [X] T011 Run domain/store focused suites green and record evidence

## Phase 3: User Story 1 - Native Grammar and Calendar Execution (P1)

**Goal**: Author all three phrase families and execute exact monthly dates under every policy.

- [X] T012 [US1] Add failing phrase parser tests for canonical/optional-time `L`, `nW`, and `LW` forms in `internal/schedule/parse_test.go`
- [X] T013 [US1] Add failing all-seven-weekday and month-boundary matrices for `1W`, middle dates, and final dates in `internal/schedule/recur_test.go`
- [X] T014 [US1] Add failing short-month, leap-year, skip/last-valid/next-valid, collision, and duplicate-suppression tests in `internal/schedule/missingdate_test.go`
- [X] T015 [US1] Add failing last-day/last-weekday month-ending and DST wall-time matrices in `internal/schedule`
- [X] T016 [US1] Add failing unknown/incompatible persisted-adjustment error tests
- [X] T017 [US1] Record expected red grammar and recurrence results
- [X] T018 [US1] Compile last-day and last-weekday phrases to native RRULE selectors in `internal/schedule/parse.go`
- [X] T019 [US1] Compile nearest-weekday phrases to a monthly carrier RRULE plus typed adjustment in `internal/schedule/parse.go`
- [X] T020 [US1] Validate and execute adjusted monthly rules with policy-before-adjustment and DST normalization in `internal/schedule/recur.go` and `missingdate.go`
- [X] T021 [US1] Make `Describe` adjustment-aware without annotating policy-inert targets
- [X] T022 [US1] Run all schedule and store suites green and record evidence

## Phase 4: User Story 2 - Cron Import, Explain, and Refusal (P1)

**Goal**: Explain/import the exact five-field selector subset and reject surrounding shapes precisely.

- [X] T023 [US2] Add failing parsed-field and canonical phrase tests for `L`, case-insensitive `nW`, and `LW` in `internal/cron/cron_test.go`
- [X] T024 [US2] Add failing malformed, invalid-date, offset, list, range, step, mixed, restricted-field, Quartz, and regression refusal tests
- [X] T025 [US2] Add failing cron-to-human conversion and twelve-month execution parity tests in `internal/cron/convert_test.go`
- [X] T026 [US2] Record expected red cron results
- [X] T027 [US2] Add a bounded day-of-month selector representation and dedicated parser in `internal/cron/cron.go`
- [X] T028 [US2] Render canonical selector phrases in `internal/cron/phrase.go`
- [X] T029 [US2] Run parser, conversion, and scheduling parity suites green and record evidence

## Phase 5: User Story 3 - Native Export and Round Trip (P1)

**Goal**: Export only lossless native selector schedules and preserve runs for at least one leap/DST year.

- [X] T030 [US3] Add failing exporter tests for canonical `L`, `1W` through `31W`, `LW`, and full policy eligibility in `internal/cron/export_test.go`
- [X] T031 [US3] Add failing richer-RRULE, malformed-adjustment, bounds, multi-time, competing-selector, and source-inference refusal tests
- [X] T032 [US3] Add failing twelve-month bidirectional run parity over leap day, weekend targets, short months, and both DST transitions
- [X] T033 [US3] Record expected red export results
- [X] T034 [US3] Export native last-day and last-weekday selector shapes with existing fidelity guards
- [X] T035 [US3] Export marked nearest-weekday shapes before ordinary month-day handling and enforce policy boundaries
- [X] T036 [US3] Run export/conversion round-trip suites green and record evidence

## Phase 6: User Story 4 - Shared Interfaces and Atomic Failure (P1/P2)

**Goal**: Preserve identical behavior/source across crontab, CLI, API, desktop, restart, and scheduling consumers.

- [X] T037 [US4] Add shared-input cron/human identity, adjustment, and run-parity tests in `internal/scheduleinput/input_test.go`
- [X] T038 [US4] Add crontab preview/import, CLI text/JSON, exact-expression retention, and regression tests in `internal/cli`
- [X] T039 [US4] Add API preview/create/reload/update matrices and whole-request non-mutation tests in `internal/api/server`
- [X] T040 [US4] Expose the adjustment in preview responses so recurrence metadata is complete
- [X] T041 [US4] Add desktop validation, preview, create/update, retained-prefill, policy-change refresh, and refusal tests in `gui`
- [X] T042 [US4] Send selected missing-date policy in desktop previews and refresh preview when it changes in `gui/editor.go`
- [X] T043 [US4] Add engine/catch-up/calendar integration coverage proving persisted adjusted schedules use shared execution
- [X] T044 [US4] Run shared-boundary and integration suites green and record evidence

## Phase 7: Documentation and Compatibility

- [X] T045 Document the Quartz-derived five-field `L`/`nW`/`LW` subset, exact distinctions, policy limits, and refusals in `docs/cron.md`
- [X] T046 Add conversion, explain, import, and policy examples in `docs/cli.md`
- [X] T047 Align desktop schedule help and `docs/gui-fields.md` with selector and preview behavior
- [X] T048 Update issue #22's parity inventory without closing it and add the newest chronological Unreleased changelog entry with `Refs #22`
- [X] T049 Run documentation integrity checks and record results

## Phase 8: Analysis, Verification, and Local Delivery

- [X] T050 Revalidate both checklists and run final Spec-Kit analysis to zero critical findings
- [X] T051 Run focused domain/store/schedule/cron/input/CLI/API/GUI/integration suites and record final results
- [X] T052 Run `git diff --check`, strict UTF-8-without-BOM decoding, and mojibake audits
- [X] T053 Run all eight canonical `scripts/verify.sh` gates in the foreground and record each result honestly
- [X] T054 Mark completed tasks, rerun analysis, and commit locally with the required co-author trailer
- [X] T055 Halt before push and PR creation with a concise delivery/compliance summary; eventual PR body must use `Refs #22`

## Dependencies and Execution Order

Design precedes red persistence and recurrence contracts. Domain/store foundations precede runtime semantics. Native execution precedes cron import/export, which precede shared interface propagation. Documentation follows stable behavior. Final analysis and all verification precede the local commit and mandatory publication halt.

## Implementation Strategy

1. Prove the durable model and migration first.
2. Establish exact calendar semantics in the scheduler before exposing syntax.
3. Add import and export only after runtime truth is stable.
4. Exercise every interface through shared boundaries, including atomic failure and restart.
5. Document, analyze, verify, commit locally, and halt once before publication.
