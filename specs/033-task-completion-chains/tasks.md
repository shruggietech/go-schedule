# Tasks: Task-Completion Chains

**Input**: Design documents in `specs/033-task-completion-chains/`

**Tests**: Required by the constitution and S033 FR-018. Every behavioral task follows a failing-test-first sequence.

**Organization**: Tasks are grouped by independently testable user story. Coordinator #72 and child issues #73-#77 ship in this one slice. Spec FR-020 is publication bookkeeping and is satisfied by the later PR description, not represented as an incomplete implementation task.

## Phase 1: Foundational Domain and Persistence

**Purpose**: Establish the forward-compatible model and durable atomic boundary required by every story.

- [x] T001 [P] Add completion outcome, delivery state, chain, delivery, run-correlation, and validation tests in `internal/domain/domain_test.go`
- [x] T002 Implement completion domain values and optional run correlation in `internal/domain/domain.go`
- [x] T003 [P] Add prior-v8, clean-database, legacy-trigger-history, and rollback migration tests in `internal/store/migration_v9_test.go`
- [x] T004 Implement forward-only v9 chain/delivery tables, indexes, and nullable run correlation columns in `internal/store/store.go`
- [x] T005 [P] Add chain CRUD, uniqueness, task validation, self-link, 100-task indirect-cycle, update, and cascade tests in `internal/store/chains_test.go`
- [x] T006 [P] Add atomic run/delivery insertion, duplicate suppression, claim, resolution, completion, and reopen recovery tests in `internal/store/delivery_test.go`
- [x] T007 Implement chain CRUD, graph validation, and enriched reads in `internal/store/chains.go`
- [x] T008 Implement atomic run recording, outgoing delivery creation, claim/recover/complete/resolve operations in `internal/store/chains.go` and update run scans in `internal/store/crud.go`
- [x] T009 Advance S033 to In Progress, run `go test ./internal/domain ./internal/store`, and record the foundational checkpoint in `specs/033-task-completion-chains/spec.md`, `specs/README.md`, and `verification.md`

**Checkpoint**: The data model is forward-migrating, acyclic, transactionally durable, and independently verified (FR-001, FR-004-FR-005, FR-008, FR-010, FR-017; issue #73).

---

## Phase 2: User Story 1 - Run One Task After Another (Priority: P1)

**Goal**: Match terminal source outcomes, dispatch targets through existing policies, correlate history, and support finite cascades.

**Independent Test**: A success chain dispatches one correlated target, selectors match correctly, normal schedules remain intact, and A-to-B-to-C completes finitely.

- [x] T010 [P] [US1] Add success/failure/any matching, non-terminal exclusion, fan-out, and source-correlation engine tests in `internal/engine/completion_test.go`
- [x] T011 [P] [US1] Add finite cascade and schedule-plus-completion coexistence integration tests in `test/integration/completion_chains_test.go`
- [x] T012 [P] [US1] Add completion-origin queue-one, skip, allow-concurrent, and extra-trigger-drop tests in `internal/engine/overlap_test.go`
- [x] T013 [US1] Introduce internal correlated dispatch origins and preserve correlation through queued work in `internal/engine/engine.go` and `internal/engine/overlap.go`
- [x] T014 [US1] Record terminal runs and outgoing deliveries atomically, then claim and dispatch eligible completion work in `internal/engine/engine.go`
- [x] T015 [US1] Resolve overlap-skipped and collapsed completion deliveries without producing downstream completions in `internal/engine/overlap.go`
- [x] T016 [US1] Preserve current schedules, group eligibility, worker bounds, callbacks, and alerts across completion dispatch in `internal/engine/engine.go`
- [x] T017 [US1] Run targeted engine and integration suites and record the primary-flow checkpoint in `specs/033-task-completion-chains/verification.md`

**Checkpoint**: Core chaining works without management UI (FR-002-FR-003, FR-006-FR-007; issue #74).

---

## Phase 3: User Story 2 - Recover Delivery Safely (Priority: P1)

**Goal**: Recover pending and interrupted work under the explicit at-least-once contract while never replaying completed work.

**Independent Test**: Reopen pending, claimed, completed, and resolved deliveries and observe the specified state transition and dispatch behavior.

- [x] T018 [P] [US2] Add engine-start recovery tests for pending, claimed replay, and 100 completed-delivery repeated-restart suppressions in `internal/engine/completion_recovery_test.go`
- [x] T019 [P] [US2] Add directly disabled, ancestor-disabled, inactive, deleted-target, and deleted-chain claim tests in `internal/engine/completion_eligibility_test.go`
- [x] T020 [US2] Recover interrupted claims before initial delivery drain and emit replay diagnostics in `internal/engine/engine.go`
- [x] T021 [US2] Resolve ineligible or stale targets terminally with structured diagnostics and no retry loop in `internal/engine/engine.go`
- [x] T022 [US2] Prove cancellation drains completion work without goroutine leaks or polling in `internal/engine/completion_recovery_test.go`
- [x] T023 [US2] Run restart/recovery tests and record at-least-once evidence and the crash-window limitation in `specs/033-task-completion-chains/verification.md`

**Checkpoint**: Durable recovery is honest, bounded, and diagnosable (FR-005, FR-009, FR-016; issues #73-#74).

---

## Phase 4: User Story 3 - Manage Chains from API and CLI (Priority: P2)

**Goal**: Provide complete local API/client/CLI lifecycle management and validation.

**Independent Test**: Create, list, show, update, and delete a chain through API and CLI; invalid mutations leave state unchanged and return compatible errors.

- [x] T024 [P] [US3] Add API lifecycle, response-shape, invalid field, duplicate, cycle, and no-partial-mutation tests in `internal/api/server/chains_test.go`
- [x] T025 [P] [US3] Add typed client lifecycle and error-envelope tests in `internal/api/client/client_test.go`
- [x] T026 [P] [US3] Add CLI text, JSON, usage, validation-exit, and transport-error tests in `internal/cli/chain_test.go`
- [x] T027 [US3] Add versioned chain routes, request/response types, validation mapping, and live publication in `internal/api/server/chains.go` and `internal/api/server/server.go`
- [x] T028 [US3] Add typed chain client methods in `internal/api/client/methods.go`
- [x] T029 [US3] Add `gosched chain create|list|show|update|rm` and root registration in `internal/cli/chain.go` and `internal/cli/cli.go`
- [x] T030 [US3] Run API/client/CLI contract suites and record the management checkpoint in `specs/033-task-completion-chains/verification.md`

**Checkpoint**: Scriptable management is complete through the existing authorized local boundary (FR-011-FR-012; issue #75).

---

## Phase 5: User Story 4 - Manage Chains Visually (Priority: P2)

**Goal**: Add a focused, live desktop Chains view with complete lifecycle and error states.

**Independent Test**: Headless GUI tests create, update, delete, refresh, and render empty/stale/error states using task names.

- [x] T031 [P] [US4] Add chain event publish/drop-safety tests in `internal/events/broker_test.go`
- [x] T032 [P] [US4] Add initial chain refresh and created/updated/deleted folding tests in `gui/viewmodel/viewmodel_test.go`
- [x] T033 [P] [US4] Add chain row/form model, label identity, empty-state, validation, and mutation tests in `gui/chains_test.go`
- [x] T034 [US4] Add chain mutation event domain and broker publication in `internal/events/broker.go`
- [x] T035 [US4] Load and fold chain state in `gui/viewmodel/viewmodel.go`
- [x] T036 [US4] Extend the GUI backend contract and test fake with chain lifecycle methods in `gui/app.go` and `gui/app_test.go`
- [x] T037 [US4] Build the dedicated Chains navigation, list, create/edit/delete forms, empty state, and task-name labels in `gui/chains.go` and `gui/app.go`
- [x] T038 [US4] Run view-model and headless GUI suites and record the desktop checkpoint in `specs/033-task-completion-chains/verification.md`

**Checkpoint**: Desktop users can manage chains without task IDs or CLI use (FR-013-FR-014; issue #76).

---

## Phase 6: User Story 5 - Diagnose and Understand Chained Runs (Priority: P2)

**Goal**: Make every completion-triggered execution traceable and document the complete product boundary.

**Independent Test**: CLI/API/GUI activity identify completion source correlation, public walkthroughs work, and non-completion history stays compatible.

- [x] T039 [P] [US5] Add completion and non-completion text/JSON history compatibility tests in `internal/cli/runs_test.go` and `internal/api/server/runs_test.go`
- [x] T040 [P] [US5] Add activity correlation rendering tests in `gui/logs_test.go`
- [x] T041 [US5] Render source correlation consistently in run CLI output, desktop activity detail, and structured completion logs in `internal/cli/runs.go`, `gui/logs.go`, and `internal/engine/engine.go`
- [x] T042 [P] [US5] Document chain lifecycle and run correlation in `docs/cli.md`, `docs/api.md`, and `docs/gui-fields.md`
- [x] T043 [P] [US5] Document delivery flow, replay semantics, and deferred external-event boundary in `README.md`, `docs/architecture.md`, and `specs/001-task-scheduler/spec.md`
- [x] T044 [US5] Add the S033 feature entry and dated architecture decision to `CHANGELOG.md`
- [x] T045 [US5] Run documentation policy checks and record the diagnostic/documentation checkpoint in `specs/033-task-completion-chains/verification.md`

**Checkpoint**: Operators can explain every chained run and follow public workflows (FR-015, FR-019; issue #77).

---

## Phase 7: Performance, Compatibility, and Done Gate

**Purpose**: Protect cross-cutting constitutional guarantees and close the whole issue hierarchy coherently.

- [x] T046 [P] Extend engine fan-out benchmarks and p99 assertions for 100 completion deliveries in `internal/engine/engine_bench_test.go` and `internal/engine/latency_test.go`
- [x] T047 [P] Add end-to-end prior-database, fan-out, converging-path, deletion, and full lifecycle coverage in `test/integration/completion_chains_test.go`
- [x] T048 Run `go test ./internal/store ./internal/engine ./internal/api/... ./internal/cli ./internal/events ./gui/... ./test/integration` and resolve every failure
- [x] T049 Reconcile all completed tasks, issue acceptance criteria, and lifecycle inventory in `specs/033-task-completion-chains/spec.md`, `tasks.md`, and `specs/README.md`
- [x] T050 Run `sh scripts/verify.sh all` in the foreground and record format, vet, lint, race, GUI, coverage, docs, and automation evidence in `specs/033-task-completion-chains/verification.md`
- [x] T051 Audit changed files for UTF-8 without BOM, mojibake, placeholders, unchecked tasks, diff errors, and unintended pinned-artifact changes in `specs/033-task-completion-chains/verification.md`
- [x] T052 Advance S033 to Implemented with local delivery evidence in `specs/033-task-completion-chains/spec.md` and `specs/README.md`

## Dependencies and Execution Order

- Phase 1 blocks every user story.
- US1 establishes live delivery and blocks US2 recovery behavior.
- US3 depends only on Phase 1 and can be validated independently from the engine using store-backed routes.
- US4 depends on US3 client contracts and live events.
- US5 depends on US1 correlation and the completed user surfaces.
- Phase 7 requires all story checkpoints.

## Parallel Opportunities

- Foundational domain, migration, graph, and delivery tests touch separate files before implementation.
- Within US3, API, client, and CLI tests can be authored independently.
- Within US4, broker, view-model, and pure GUI model tests can be authored independently.
- Documentation files and benchmark/integration tests are independent until the final verification gate.

## Implementation Strategy

1. Complete the durable model and graph validation first.
2. Deliver one independently testable source-to-target chain through the engine.
3. Add restart recovery without changing primary-flow semantics.
4. Expose the stable lifecycle through API/CLI, then the desktop.
5. Finish correlation, docs, scale evidence, and the aggregate done gate.
