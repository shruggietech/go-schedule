# Tasks: Dual-Syntax Task Input Foundation

**Input**: Design documents from `specs/019-dual-syntax-input/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Required by the feature specification and constitution. Every story uses a recorded red phase before production changes.

**Organization**: Tasks are grouped by independently testable user story after the shared input boundary.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it affects a separate file and has no incomplete dependency.
- **[Story]**: Maps a task to US1, US2, or US3 in spec.md.

## Phase 1: Setup and Baseline

**Purpose**: Prove the starting behavior and capture the exact boundaries being superseded.

- [X] T001 Run the current `internal/cron`, `internal/schedule`, `internal/api/server`, `internal/cli`, and `gui` focused tests and record the green baseline in `specs/019-dual-syntax-input/verification.md`
- [X] T002 Record current detector, human-only API parsing, expression persistence, import phrase substitution, GUI validation, and no-migration evidence in `specs/019-dual-syntax-input/verification.md`

---

## Phase 2: Foundational Central Input Boundary

**Purpose**: Establish one dual-syntax parser and identity contract that blocks all product-surface stories.

**Critical**: No user story implementation begins until this phase is green.

### Tests

- [X] T003 Add failing automatic/forced classification, hint validation, no-fallback, retained-expression, named-refusal, human-compatibility, and DST/month parity tests in `internal/scheduleinput/input_test.go`
- [X] T004 Add failing exported detector parity and five-word-human regression tests in `internal/cron/convert_test.go`
- [X] T005 Run the foundational focused tests and record the expected pre-implementation failures in `specs/019-dual-syntax-input/verification.md`

### Implementation

- [X] T006 Export the single S018 structural detector and route converter auto mode through it in `internal/cron/convert.go`
- [X] T007 Implement typed syntax hints, single-parser compilation, retained source, named cron refusals, and source-identity derivation in `internal/scheduleinput/input.go`
- [X] T008 Run `internal/cron` and `internal/scheduleinput` focused tests green and record the boundary evidence in `specs/019-dual-syntax-input/verification.md`

**Checkpoint**: One deterministic parser can compile and identify either syntax without storage, API, CLI, GUI, or engine coupling.

---

## Phase 3: User Story 1 - Create and Preview with Cron (Priority: P1) MVP

**Goal**: Preview and create a recurring task with supported cron through the API/CLI request path, retaining cron source and returning source identity.

**Independent Test**: Preview and create `0 9 * * 1-5`, fetch the task, and compare RRULE/upcoming runs to `weekdays at 09:00` while proving named failures mutate nothing.

### Tests

- [X] T009 [P] [US1] Add failing cron preview/create, explicit hint, invalid hint field, no-fallback/refusal, human regression, and run-parity tests in `internal/api/server/tasks_test.go`
- [X] T010 [P] [US1] Add failing cron expression storage reload, derived source identity, one-off identity, and expressionless legacy compatibility tests in `internal/api/server/expression_test.go`
- [X] T011 [US1] Run the US1 server tests and record the expected pre-implementation failures in `specs/019-dual-syntax-input/verification.md`

### Implementation

- [X] T012 [US1] Add transient recurring `source_syntax` response metadata and update source-expression intent comments in `internal/domain/domain.go`
- [X] T013 [US1] Add optional `schedule_syntax`, central parse/error routing, preview identity, and derived task response identity in `internal/api/server/tasks.go`
- [X] T014 [US1] Run the US1 and foundational suites green and record preview/create/read/storage parity in `specs/019-dual-syntax-input/verification.md`

**Checkpoint**: Supported cron is independently useful for preview and task creation, with human input and one-offs unchanged.

---

## Phase 4: User Story 2 - Edit and Retrieve Original Syntax (Priority: P2)

**Goal**: Replace or preserve recurring schedules through update, expose honest CLI help, and prevent accidental GUI scope expansion.

**Independent Test**: Update cron to human and human to cron, perform an unrelated edit, and verify source/timing/policy changes only when intended.

### Tests

- [X] T015 [P] [US2] Add failing cron-to-human, human-to-cron, omitted-schedule preservation, invalid-hint, invalid-cron no-mutation, and missing-date-policy update tests in `internal/api/server/update_test.go`
- [X] T016 [P] [US2] Add failing dual-syntax `--schedule` help contract tests in `internal/cli/task_test.go`
- [X] T017 [P] [US2] Add failing GUI request assertions that recurring preview/create/update explicitly select human syntax in `gui/editor_test.go` and `gui/app_test.go`
- [X] T018 [US2] Run the US2 focused tests and record the expected pre-implementation failures in `specs/019-dual-syntax-input/verification.md`

### Implementation

- [X] T019 [US2] Route recurring replacement and syntax-hint validation through the central parser in `internal/api/server/update.go`
- [X] T020 [US2] Update task add/edit schedule help without adding local parsing in `internal/cli/task.go`
- [X] T021 [US2] Set explicit human syntax on existing recurring GUI preview/create/update requests in `gui/editor.go`
- [X] T022 [US2] Run API, CLI, GUI, and foundational suites green and record source/policy/GUI-containment evidence in `specs/019-dual-syntax-input/verification.md`

**Checkpoint**: Non-GUI edit/read round-trips both syntaxes, while GUI behavior remains deliberately unchanged.

---

## Phase 5: User Story 3 - Keep Cron Source Through Import (Priority: P3)

**Goal**: Preview and create imported jobs from the scanner's cron expression while retaining explanatory phrase reporting and existing partial-success rules.

**Independent Test**: Import a supported job through a recording preview/create double and verify both receive `Line.Expr` with cron identity, while dry-run, decline, unreachable-daemon, command, and partial-success results remain stable.

### Tests

- [X] T023 [US3] Replace the obsolete human-only import policy assertion with failing cron source/hint retention, preview/create input parity, and mixed partial-success regressions in `internal/cli/cron_test.go`
- [X] T024 [US3] Run the focused import tests and record the expected pre-implementation failures in `specs/019-dual-syntax-input/verification.md`

### Implementation

- [X] T025 [US3] Preview `Line.Expr` through explicit cron syntax and retain the existing unavailable-daemon report behavior in `internal/cli/cron.go`
- [X] T026 [US3] Create imported tasks from `Line.Expr` with explicit cron syntax while preserving phrase output, commands, counts, refusals, and continuation in `internal/cli/cron.go`
- [X] T027 [US3] Run the complete import/cron/API focused suites green and record retention and partial-success evidence in `specs/019-dual-syntax-input/verification.md`

**Checkpoint**: Imported users remain cron users after task creation.

---

## Phase 6: Polish and Cross-Cutting Verification

**Purpose**: Synchronize current truth, validate artifacts, and prove repository quality before the local commit.

- [X] T028 [P] Narrowly correct live dual-syntax task/input guidance in `README.md`, `docs/cli.md`, `docs/cron.md`, and the focused API contract without performing issue #52's broad rewrite (the repository has no `docs/api.md`)
- [X] T029 [P] Supersede current boundary-only package comments in `internal/cron/cron.go`, `internal/schedule/parse.go`, and `internal/schedule/recur.go`
- [X] T030 Update the Unreleased feature and dated architecture decision in `CHANGELOG.md` with `Refs #50`, no schema migration, RRULE authority, and GUI deferral
- [X] T031 Validate both Spec-Kit checklists and record focused quickstart, issue disposition, and no-security/no-migration evidence in `specs/019-dual-syntax-input/verification.md`
- [X] T032 Run `git diff --check`, strict UTF-8-without-BOM, and mojibake audits across every changed file
- [X] T033 Run `sh scripts/verify.sh all` in the foreground and record all eight gate results in `specs/019-dual-syntax-input/verification.md`
- [X] T034 Mark every completed task `[X]` in `specs/019-dual-syntax-input/tasks.md` and rerun `/speckit-analyze`
- [X] T035 Commit locally as `feat(019): add dual-syntax task input` with the required co-author trailer

---

## Dependencies and Execution Order

### Phase dependencies

- Phase 1 has no dependencies.
- Phase 2 follows the baseline and blocks every user story.
- US1 depends on the central input boundary.
- US2 depends on the US1 response/request contract and central input boundary.
- US3 depends on US1 preview/create contracts and the central input boundary.
- Phase 6 depends on all selected user stories.

### User story dependencies

- **US1** is the MVP and first independently usable behavior.
- **US2** extends the same API contract to update and preserves GUI scope.
- **US3** reuses the stable preview/create contract for import retention.

### Parallel opportunities

- T003 and T004 affect separate packages before shared implementation.
- T009 and T010 split API behavior from persistence/identity behavior.
- T015, T016, and T017 protect update, CLI, and GUI files independently.
- T028 and T029 update separate documentation/source-comment surfaces.

## Parallel Example: User Story 2

```text
Task T015: write update and no-mutation tests in internal/api/server/update_test.go
Task T016: write CLI help contract tests in internal/cli/task_test.go
Task T017: write GUI human-hint request tests in gui tests
Then T019-T021: implement server, CLI, and GUI preservation paths
```

## Implementation Strategy

### MVP first

1. Capture the existing green baseline.
2. Deliver the central input boundary test-first.
3. Complete US1 preview/create/read and validate it independently.

### Incremental delivery

1. Extend the same contract to update and CLI help while containing GUI scope.
2. Switch import preview/create to retained cron source.
3. Synchronize current docs/comments, run analysis and all verification gates, then commit locally and halt before publication.
