# Tasks: Task Execution Safety and Diagnostics

**Input**: Design documents from `specs/051-task-execution-safety/`

**Tests**: Regression, persistence, interface, engine, executor, headless GUI, race, coverage, documentation, and automation tests are mandatory.

## Phase 1: Setup

**Purpose**: Establish S051 scope, traceability, and decision records.

- [x] T001 Create the S051 specification, clarification record, requirements checklists, plan, research, data model, interface contract, and quickstart under `specs/051-task-execution-safety/`
- [x] T002 Activate S051 in `.specify/feature.json`, `CLAUDE.md`, and `specs/README.md`

---

## Phase 2: Foundational Diagnostic Identity and Compatibility

**Purpose**: Add the durable, backward-compatible contracts required before failed-run detail or safe desktop creation can be implemented.

**Critical**: User-story implementation starts only after the additive model, migration, and exact lookup contracts are proven.

- [x] T003 [P] Write failing domain JSON compatibility tests for optional alert run identity and run truncation in `internal/domain/domain_test.go`
- [x] T004 [P] Write failing v9-to-v10 preservation and transactional rollback tests in `internal/store/migration_v10_test.go`
- [x] T005 Add `Run.OutputTruncated` and `Alert.RunID` to `internal/domain/domain.go` and migration v10 to `internal/store/store.go`
- [x] T006 Write failing run/alert round-trip and exact-not-found tests in `internal/store/store_test.go`
- [x] T007 Persist and scan the new fields and implement exact `GetRun` lookup in `internal/store/crud.go`
- [x] T008 Run focused domain and store tests and record the foundational checkpoint in `specs/051-task-execution-safety/verification.md`

**Checkpoint**: Existing data upgrades safely, old JSON remains compatible, and one exact run can be retrieved without recency inference.

---

## Phase 3: User Story 1 - Diagnose an Exact Failed Run (Priority: P1)

**Goal**: Make every new failed-run Activity entry resolve to bounded, selectable diagnostics for the exact persisted task run.

**Independent Test**: Produce nonzero-exit, launch-failure, empty-output, multiline-output, and truncated-output runs, then open each matching alert and verify exact identity, status, trigger, and privacy-safe output.

### Tests for User Story 1

- [x] T009 [P] [US1] Write failing bounded-output truncation and empty/multiline capture tests in `internal/executor/executor_test.go`
- [x] T010 [P] [US1] Write failing exact run-failure alert correlation tests across run triggers in `internal/engine/engine_extra_test.go`
- [x] T011 [P] [US1] Write failing `GET /v1/runs/{id}` success/not-found and alert JSON tests in `internal/api/server/runs_test.go`
- [x] T012 [P] [US1] Write failing exact-run client method tests in `internal/api/client/` test files
- [x] T013 [P] [US1] Write failing pure diagnostic-format and alert-correlation merge tests in `gui/logs_test.go`
- [x] T014 [US1] Write failing headless dialog enrichment tests with exact, missing, legacy, empty, and truncated run states in `gui/app_test.go` and `gui/logs_test.go`

### Implementation for User Story 1

- [x] T015 [US1] Track discarded output bytes without exceeding the configured cap in `internal/executor/executor.go`
- [x] T016 [US1] Attach the persisted run ID to failure alerts without changing other alert kinds in `internal/engine/engine.go`
- [x] T017 [US1] Register and implement exact run retrieval in `internal/api/server/server.go` and `internal/api/server/runs.go`
- [x] T018 [US1] Implement the exact-run client contract in `internal/api/client/methods.go`
- [x] T019 [US1] Extend the desktop backend and fake backend with bounded exact task/run enrichment in `gui/app.go` and `gui/app_test.go`
- [x] T020 [US1] Preserve task/run identity in Activity rows and render selectable combined-output diagnostics with honest fallback states in `gui/logs.go`
- [x] T021 [US1] Run focused executor, engine, API, client, and GUI tests and record the US1 checkpoint in `specs/051-task-execution-safety/verification.md`

**Checkpoint**: Failed-run detail is exact, bounded, actionable, selectable, backward compatible, and adds no secret-bearing task inputs.

---

## Phase 4: User Story 2 - Create a Task Without Accidental Execution (Priority: P2)

**Goal**: Default desktop-created tasks to inactive while preserving atomic opt-in activation and existing non-desktop behavior.

**Independent Test**: Exercise omitted, false, and true create requests plus fresh, validation-recovery, and edit desktop flows; no inactive request may be observed as scheduler eligible.

### Tests for User Story 2

- [x] T022 [P] [US2] Write failing omitted/false/true creation and published-state tests in `internal/api/server/tasks_test.go` and `internal/api/server/publish_test.go`
- [x] T023 [P] [US2] Write failing new-task checkbox default, dirty-state, validation-retention, submit-intent, fresh-reset, and edit-preservation tests in `gui/editor_test.go`, `gui/editor_data_test.go`, and `gui/editor_prefill_test.go`

### Implementation for User Story 2

- [x] T024 [US2] Add an optional enabled intent and resolve its compatible default atomically in `internal/api/server/tasks.go`
- [x] T025 [US2] Add the creation-only activation checkbox to editor construction, snapshots, form data, and submission in `gui/editor.go`
- [x] T026 [US2] Run focused server and headless editor tests and record the US2 checkpoint in `specs/051-task-execution-safety/verification.md`

**Checkpoint**: Fresh desktop tasks remain inactive unless explicitly activated, while legacy callers and existing task edits retain their prior contract.

---

## Phase 5: User Story 3 - Understand Effective Task Eligibility (Priority: P2)

**Goal**: Separate configured, lifecycle, and effective task state and identify the nearest disabled group responsible for suppression.

**Independent Test**: Render and live-refresh ungrouped, own-disabled, lifecycle-inactive, direct-blocked, ancestor-blocked, missing-chain, and cyclic tasks while preserving stable identity and full-value disclosure.

### Tests for User Story 3

- [x] T027 [P] [US3] Write failing nearest-disabled ancestor, enabled, missing, and cycle tests in `internal/task/group_test.go`
- [x] T028 [US3] Write failing effective-state precedence, table header, disclosure, narrow-mode, and live group-event refresh tests in `gui/tasks_test.go`, `gui/groupchoice_test.go`, and `gui/app_test.go`

### Implementation for User Story 3

- [x] T029 [US3] Add nearest disabled-group discovery beside authoritative chain eligibility in `internal/task/group.go`
- [x] T030 [US3] Derive and display the labeled Effective task column with full-value disclosure in `gui/tasks.go`
- [x] T031 [US3] Run focused group-policy, view-model event, and headless Tasks-table tests and record the US3 checkpoint in `specs/051-task-execution-safety/verification.md`

**Checkpoint**: Every task row gives one unambiguous, accessible explanation of configured and effective eligibility without duplicating scheduling policy.

---

## Phase 6: Cross-Cutting Verification and Delivery Record

**Purpose**: Prove the bundled slice as one backward-compatible operator journey.

- [x] T032 Run focused race tests for affected non-Fyne packages and `gui/viewmodel` and record results in `specs/051-task-execution-safety/verification.md`
- [x] T033 Run the end-to-end quickstart scenarios and full Go test suite and record results in `specs/051-task-execution-safety/verification.md`
- [x] T034 Run canonical `scripts/verify.sh all` in the foreground and record all eight gates plus core-package coverage in `specs/051-task-execution-safety/verification.md`
- [x] T035 Update the Unreleased feature and dated decision records in `CHANGELOG.md`, then set S051 lifecycle evidence in `specs/051-task-execution-safety/spec.md` and `specs/README.md`
- [x] T036 Audit UTF-8 without BOM, mojibake, placeholders, task ID uniqueness, issue acceptance coverage, backward compatibility, diagnostic privacy, diff scope, and constitution compliance in `specs/051-task-execution-safety/verification.md`

---

## Dependencies and Execution Order

- Setup precedes every design and implementation action.
- Foundational model/migration/lookup work blocks User Story 1.
- User Story 1 is the P1 diagnostic defect and lands before the P2 GUI stories.
- User Stories 2 and 3 share GUI files and therefore execute sequentially even though their product outcomes remain independently testable.
- Cross-cutting verification depends on all three user-story checkpoints.

## Parallel Opportunities

- Foundational domain and migration tests can be authored independently.
- Executor, engine, server, client, and pure GUI diagnostic tests can be authored independently before their implementation convergence.
- Server create-contract and editor tests can be authored independently.
- Group-policy tests can be authored independently from GUI row tests.

## Implementation Strategy

1. Establish additive persistent and wire identity without changing behavior.
2. Deliver exact failed-run diagnosis as the P1 independently testable outcome.
3. Deliver atomic safe desktop creation without changing legacy caller defaults.
4. Deliver effective eligibility explanation using shared group policy.
5. Run focused, race, full-suite, and canonical verification before committing.

## Format Validation

All 36 tasks use the required checkbox and sequential ID format. User-story tasks carry `[US1]`, `[US2]`, or `[US3]`; `[P]` appears only where file and dependency boundaries permit independent work.
