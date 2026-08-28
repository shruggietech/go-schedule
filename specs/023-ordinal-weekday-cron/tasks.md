# Tasks: Ordinal-Weekday Cron Parity

**Input**: Design documents from `specs/023-ordinal-weekday-cron/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Required. Behavioral seams are recorded red before production implementation.

## Phase 1: Setup and Baseline

- [X] T001 Create the review branch, record focused and documentation baselines in `specs/023-ordinal-weekday-cron/verification.md`, and validate both checklists
- [X] T002 Complete parser, export, workflow, dialect, documentation, and test-surface research in `specs/023-ordinal-weekday-cron/research.md`
- [X] T003 Run cross-artifact Spec-Kit analysis for `specs/023-ordinal-weekday-cron/spec.md`, `plan.md`, and `tasks.md` with zero critical findings

---

## Phase 2: Foundational Test Contract

- [X] T004 Add failing parsed-field and phrase tests for numeric/name/case/Sunday aliases and all 35 weekday-ordinal combinations in `internal/cron/cron_test.go`
- [X] T005 Add failing malformed and named-refusal tests for invalid ordinals, token shapes, non-DOW placement, and restricted DOM/month combinations in `internal/cron/cron_test.go`
- [X] T006 Add failing two-way conversion tests for canonical phrases and expressions in `internal/cron/convert_test.go`
- [X] T007 Run focused foundational tests and record expected pre-implementation failures in `specs/023-ordinal-weekday-cron/verification.md`

---

## Phase 3: User Story 1 - Understand and Import Ordinal Cron (Priority: P1)

**Goal**: Explain, convert, preview, and import one ordinal weekday while preserving original cron source and run behavior.

**Independent Test**: Import numeric and named forms and compare generated runs with the corresponding native monthly rule.

### Tests

- [X] T008 [US1] Add failing crontab classification, readable phrase, and original-expression retention tests in `internal/cli/cron_test.go`
- [X] T009 [P] [US1] Add failing shared input tests for cron source, syntax, RRULE, and generated-run parity in `internal/scheduleinput/input_test.go`
- [X] T010 [P] [US1] Add failing API create/update acceptance tests and malformed/refused non-mutation regressions in `internal/api/server/tasks_test.go` and `internal/api/server/update_test.go`
- [X] T011 [US1] Record the expected pre-implementation failures for shared input, CLI, and API boundaries in `specs/023-ordinal-weekday-cron/verification.md`

### Implementation

- [X] T012 [US1] Add the bounded ordinal field representation and dedicated day-of-week parser in `internal/cron/cron.go`
- [X] T013 [US1] Render supported ordinal weekday input through the existing monthly phrase grammar in `internal/cron/phrase.go`
- [X] T014 [US1] Run parser, conversion, shared-input, CLI, and API suites green and record evidence in `specs/023-ordinal-weekday-cron/verification.md`

---

## Phase 4: User Story 2 - Export an Ordinal-Weekday Schedule (Priority: P1)

**Goal**: Export the existing native monthly ordinal recurrence to canonical numeric cron without changing run times.

**Independent Test**: Export and re-import first through fifth weekdays across DST, month boundaries, and a missing fifth occurrence.

### Tests

- [X] T015 [US2] Add failing exporter tests for all 35 canonical expressions, Sunday zero, selector shape, and missing-policy behavior in `internal/cron/export_test.go`
- [X] T016 [US2] Add failing runtime round-trip tests across DST, three month boundaries, and an absent fifth occurrence in `internal/cron/export_test.go` and `internal/cron/convert_test.go`
- [X] T017 [US2] Record expected pre-implementation exporter and round-trip failures in `specs/023-ordinal-weekday-cron/verification.md`

### Implementation

- [X] T018 [US2] Export exactly one lossless positive numbered weekday from a monthly recurrence in `internal/cron/export.go`
- [X] T019 [US2] Run exporter, conversion, and round-trip suites green and record evidence in `specs/023-ordinal-weekday-cron/verification.md`

---

## Phase 5: User Story 3 - Precise Boundary Feedback (Priority: P2)

**Goal**: Reject every malformed or unsupported surrounding shape clearly and without task mutation.

**Independent Test**: Exercise each contract refusal through the parser and at least one shared task boundary.

- [X] T020 [US3] Verify `L`, `W`, six-field, `@reboot`, existing list, and lossy-combination refusals remain compatible in `internal/cron/cron_test.go`
- [X] T021 [US3] Refine ordinal errors and named refusals to identify the day-of-week field, valid range, or one-term boundary in `internal/cron/cron.go` and `internal/cron/phrase.go`
- [X] T022 [US3] Run boundary and task non-mutation suites green and record evidence in `specs/023-ordinal-weekday-cron/verification.md`

---

## Phase 6: Documentation and Compatibility

- [X] T023 [US1] [US2] Document the supported single-`#` subset, dialect caveat, canonical output, and remaining refusals in `docs/cron.md`
- [X] T024 [US1] [US2] Add concise explain, convert, and crontab examples in `docs/cli.md`
- [X] T025 [US1] [US2] Add the newest chronological Unreleased changelog entry that references but does not close issue #22 in `CHANGELOG.md`
- [X] T026 Run documentation integrity checks and record results in `specs/023-ordinal-weekday-cron/verification.md`

---

## Phase 7: Polish and Cross-Cutting Verification

- [X] T027 Validate both checklists and rerun Spec-Kit analysis to zero critical findings
- [X] T028 Run focused cron, shared-input, CLI, and API suites and record final results in `specs/023-ordinal-weekday-cron/verification.md`
- [X] T029 Run `git diff --check`, strict UTF-8-without-BOM, and mojibake audits across changed files and record results in `specs/023-ordinal-weekday-cron/verification.md`
- [X] T030 Run `sh scripts/verify.sh all` in the foreground and record all eight gate results in `specs/023-ordinal-weekday-cron/verification.md`
- [X] T031 Mark every completed task `[X]`, rerun analysis, and commit locally with the required co-author trailer
- [X] T032 Halt before push and pull-request creation with a concise compliance and delivery summary requiring `Refs #22` in the PR body

## Dependencies and Execution Order

- Setup, clarification, checklists, research, and initial analysis precede implementation.
- Foundational parser and conversion tests establish the contract for all stories.
- US1 provides import and shared-boundary support before US2 adds export.
- US3 validates refusal compatibility after both accepted paths are stable.
- Documentation follows the final syntax and fidelity behavior.
- Final analysis, audits, and all canonical gates precede the local commit and publication halt.

## Parallel Opportunities

- T009 and T010 cover independent shared-input and API files after the foundational contract is red.
- Documentation edits T023 and T024 affect separate files once behavior is stable.
- Production files remain sequential because parser metadata, phrase rendering, and exporter fidelity share one contract.

## Implementation Strategy

1. Establish the exact accepted and refused syntax in failing tests.
2. Implement import through the existing phrase and task-input pipeline.
3. Implement canonical export and prove time equivalence.
4. Confirm unsupported surrounding syntax still fails honestly and non-mutatingly.
5. Document, analyze, verify, commit locally, and halt before publication.
