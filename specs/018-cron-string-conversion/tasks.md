# Tasks: Pure Schedule String Conversion

**Input**: Design documents from `specs/018-cron-string-conversion/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, `quickstart.md`

**Tests**: Required by the feature specification, constitution, and autopilot protocol. Behavioral regressions precede implementation.

## Phase 1: Setup and baseline

**Purpose**: Capture existing cron fidelity, export behavior, and CLI health before introducing the symmetric converter.

- [X] T001 Run the current `internal/cron` and `internal/cli` focused tests and record the green baseline in `specs/018-cron-string-conversion/verification.md`
- [X] T002 Record current cron explanation, schedule export, task-authoring rejection, daemon dependency, and CLI output conventions in `specs/018-cron-string-conversion/verification.md`

---

## Phase 2: User Story 1 - Translate cron without a daemon (Priority: P1)

**Goal**: Convert one cron expression to one human schedule locally with exact text output and no fallback.

**Independent Test**: With no daemon dependency, `0 9 * * 1-5` produces only `weekdays at 09:00` plus newline; malformed and unsupported cron produce named refusals and are never parsed as human input.

### Tests for User Story 1

- [X] T003 [P] [US1] Add failing cron-input classification, success, malformed, unsupported, whitespace, and no-fallback tests in `internal/cron/convert_test.go`
- [X] T004 [US1] Add failing exact default-text, `--to human`, invalid-target, and daemon-independent cron conversion command tests in `internal/cli/cron_test.go`
- [X] T005 [US1] Run the focused conversion tests and record the expected pre-implementation failures in `specs/018-cron-string-conversion/verification.md`

### Implementation for User Story 1

- [X] T006 [US1] Implement syntax identities, stable conversion results, automatic classification, and cron-to-human conversion in `internal/cron/convert.go`
- [X] T007 [US1] Register the local one-argument `cron convert` text command, destination validation, and cron-to-human rendering in `internal/cli/cron.go`
- [X] T008 [US1] Run the focused tests and confirm the independently useful cron-to-human story is green

**Checkpoint**: One cron expression can be translated locally without preview, labels, mutation, or daemon access.

---

## Phase 3: User Story 2 - Produce one faithful cron expression (Priority: P1)

**Goal**: Convert supported human schedules to canonical five-field cron while refusing implicit-anchor and otherwise lossy schedules.

**Independent Test**: `weekdays at 09:00` produces `0 9 * * 1-5`; the result round-trips across DST and a month boundary, while implicit creation-aligned and unsupported schedules are refused by name.

### Tests for User Story 2

- [X] T009 [P] [US2] Add failing human-to-cron canonical, implicit-anchor, carve-out, forced-direction, and calendar-parity tests in `internal/cron/convert_test.go`
- [X] T010 [P] [US2] Add schedule-only renderer and task-export preservation regressions in `internal/cron/export_test.go`
- [X] T011 [US2] Run the focused conversion and export tests and record the expected pre-implementation failures in `specs/018-cron-string-conversion/verification.md`

### Implementation for User Story 2

- [X] T012 [US2] Extract the schedule-only recurrence renderer while preserving task-state checks in `internal/cron/export.go`
- [X] T013 [US2] Implement deterministic human parsing, explicit-anchor fidelity checks, and human-to-cron conversion in `internal/cron/convert.go`
- [X] T014 [US2] Complete human-to-cron text integration in `internal/cli/cron.go` and run the focused suites green

**Checkpoint**: Both conversion directions are local, canonical, and protected against silent timing approximation.

---

## Phase 4: User Story 3 - Automate conversion reliably (Priority: P2)

**Goal**: Add forced destinations plus stable structured success/refusal streams and exit behavior.

**Independent Test**: Success JSON is emitted on stdout; malformed and unfaithful JSON is emitted on stderr with empty stdout and exit class 2; all five fields are present and no duplicate plain diagnostic appears.

### Tests for User Story 3

- [X] T015 [US3] Add failing structured-field, exact-stream, and reported-error exit tests in `internal/cli/cron_test.go` and `internal/cli/cli_test.go`

### Implementation for User Story 3

- [X] T016 [US3] Implement destination validation and stable structured conversion rendering through Cobra writers in `internal/cli/cron.go`
- [X] T017 [US3] Add reported structured-usage error suppression without changing existing CLI error behavior in `internal/cli/cli.go`
- [X] T018 [US3] Run all focused conversion, CLI, JSON, stream, and exit tests and confirm all three stories are green

**Checkpoint**: Text and machine consumers receive deterministic conversion results with conventional streams and exit codes.

---

## Phase 5: Polish and cross-cutting verification

**Purpose**: Synchronize help, user documentation, traceability, and repository quality.

- [X] T019 [P] Document convert, forced direction, quoting, JSON, and exit behavior in `docs/cli.md`
- [X] T020 [P] Distinguish convert, explain, import, export, and task-authoring boundaries in `docs/cron.md`
- [X] T021 Update the Unreleased Added section in `CHANGELOG.md` with issue #51 closure and #50 relationship
- [X] T022 Validate both Spec-Kit checklists and record issue #51 closure eligibility, #50 reference status, and final focused evidence in `specs/018-cron-string-conversion/verification.md`
- [X] T023 Run `git diff --check`, strict UTF-8-without-BOM, and mojibake audits across all changed files
- [X] T024 Run `sh scripts/verify.sh all` in the foreground and record all eight gate results in `specs/018-cron-string-conversion/verification.md`
- [X] T025 Mark every completed task `[X]` in `specs/018-cron-string-conversion/tasks.md` and rerun `/speckit-analyze`
- [X] T026 Commit locally as `feat(018): add pure schedule conversion` with the required co-author trailer

---

## Dependencies and execution order

### Phase dependencies

- Phase 1 has no dependencies.
- User Story 1 follows the baseline and establishes the shared result and CLI surface.
- User Story 2 depends on User Story 1's result model but is independently protected at the pure conversion boundary.
- User Story 3 depends on both successful conversion directions.
- Phase 5 depends on all user stories.

### User story dependencies

- **US1** is the MVP: pure cron-to-human conversion without the daemon.
- **US2** extends the same operation symmetrically and introduces no new CLI surface.
- **US3** adds automation contracts after both directions have stable outcomes.

### Parallel opportunities

- T003 and T004 affect separate pure-conversion and CLI test files before their shared implementation.
- T009 and T010 split conversion and export-preservation regressions.
- T019 and T020 update separate documentation files after behavior is stable.

## Parallel example: User Story 2

```text
Task T009: define human conversion and calendar parity in internal/cron/convert_test.go
Task T010: protect schedule-only and task export behavior in internal/cron/export_test.go
Then T012-T014: extract once, integrate human conversion, and run green
```

## Implementation strategy

### MVP first

1. Capture the existing green baseline.
2. Write red cron classification and exact-output tests.
3. Deliver local cron-to-human conversion through the new command.

### Incremental delivery

1. Add faithful human-to-cron conversion without changing task input.
2. Add forced and structured automation behavior.
3. Synchronize documentation, run all eight gates, and close #51 while only referencing #50.
