# Tasks: External Trigger Lifecycle

**Input**: Design documents from `/specs/054-external-trigger-lifecycle/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Tests are required and precede implementation for each behavior.

## Phase 1: Foundation

**Purpose**: Establish the shared model, persistence, and readiness invariants required by every user story.

- [x] T001 Add failing domain serialization and run-provenance tests in `internal/domain/domain_test.go`
- [x] T002 Add failing migration and trigger CRUD tests in `internal/store/migrate_test.go` and `internal/store/triggers_test.go`
- [x] T003 Add failing automatic-source invariant tests in `internal/store/crud_test.go`, `internal/store/chains_test.go`, and `internal/task/readiness_test.go`
- [x] T004 Add trigger entities, readiness states, and external run provenance in `internal/domain/domain.go`
- [x] T005 Add SQLite migration v12 for triggers and run provenance in `internal/store/migrations/` and `internal/store/migrate.go`
- [x] T006 Implement secure key generation and trigger CRUD transactions in `internal/store/triggers.go`
- [x] T007 Extend automatic-source readiness and atomic task deactivation across schedules, chains, and triggers in `internal/task/readiness.go`, `internal/store/crud.go`, `internal/store/chains.go`, and `internal/store/triggers.go`

**Checkpoint**: Persistent triggers and automatic-source invariants work without an interface.

## Phase 2: User Story 1 - Invoke a Task from a Local Process (Priority: P1)

**Goal**: A local process can fire one eligible task by passing one opaque key.

**Independent Test**: Create a trigger for an eligible task, invoke its key, and observe one dispatch through the normal scheduler path.

- [x] T008 [US1] Add failing trigger eligibility, group blocking, overlap, and concurrency tests in `internal/engine/triggers_test.go`
- [x] T009 [US1] Add failing fire endpoint contract and redaction tests in `internal/api/server/triggers_test.go`
- [x] T010 [US1] Add failing fire client and CLI tests in `internal/api/client/client_test.go` and `internal/cli/trigger_test.go`
- [x] T011 [US1] Implement trigger eligibility errors and dispatch with source provenance in `internal/engine/triggers.go`, `internal/engine/overlap.go`, and `internal/engine/engine.go`
- [x] T012 [US1] Add the authenticated `POST /v1/triggers/fire` route and stable error mapping in `internal/api/server/server.go` and `internal/api/server/triggers.go`
- [x] T013 [US1] Add the typed fire client and `gosched trigger fire <key>` command in `internal/api/client/client.go`, `internal/cli/trigger.go`, and `internal/cli/cli.go`
- [x] T014 [US1] Add a trigger decision benchmark covering 100 concurrent calls in `internal/engine/triggers_test.go`

**Checkpoint**: The headless invocation path is independently usable and measured.

## Phase 3: User Story 2 - Administer Individual Triggers from the CLI (Priority: P1)

**Goal**: Users can manage the full trigger lifecycle through stable CLI commands.

**Independent Test**: Create, list, inspect, edit, disable, enable, rotate, reveal, and delete a trigger across a daemon restart.

- [x] T015 [US2] Add failing lifecycle API tests, including uniqueness and redaction boundaries, in `internal/api/server/triggers_test.go`
- [x] T016 [US2] Add failing lifecycle client and CLI tests in `internal/api/client/client_test.go` and `internal/cli/trigger_test.go`
- [x] T017 [US2] Implement trigger create, list, get, update, enable, disable, rotate, reveal, and delete handlers in `internal/api/server/triggers.go`
- [x] T018 [US2] Implement corresponding typed client methods in `internal/api/client/client.go`
- [x] T019 [US2] Implement lifecycle subcommands and human-readable output in `internal/cli/trigger.go`
- [x] T020 [US2] Verify daemon restart persistence and immediate old-key invalidation in `internal/api/server/triggers_test.go`

**Checkpoint**: All individual trigger administration is available headlessly.

## Phase 4: User Story 3 - Administer Triggers from the Desktop (Priority: P2)

**Goal**: Users can understand and manage triggers without leaving the desktop application.

**Independent Test**: Complete every lifecycle action through the Triggers view and observe live readiness changes.

- [x] T021 [US3] Add failing trigger snapshot and event-reduction tests in `gui/viewmodel/model_test.go`
- [x] T022 [US3] Add failing navigation, table, dialog, clipboard, confirmation, and refresh tests in `gui/triggers_test.go` and `gui/app_test.go`
- [x] T023 [US3] Add trigger lifecycle event types and redacted publisher payloads in `internal/events/broker.go` and `internal/api/server/publish.go`
- [x] T024 [US3] Extend the GUI backend and view model with trigger operations and live updates in `gui/app.go` and `gui/viewmodel/model.go`
- [x] T025 [US3] Implement the Triggers table, readiness presentation, and accessible actions in `gui/triggers.go`
- [x] T026 [US3] Implement create and edit dialogs with task selection and enabled state in `gui/triggers.go`
- [x] T027 [US3] Implement explicit reveal, copy key, copy command, rotate confirmation, and delete confirmation in `gui/triggers.go`
- [x] T028 [US3] Integrate Triggers after Chains in the Definitions section and register refresh behavior in `gui/app.go`

**Checkpoint**: Desktop lifecycle management satisfies issue #133.

## Phase 5: User Story 4 - Audit Triggered Runs without Leaking Keys (Priority: P2)

**Goal**: Triggered executions retain useful identity while all ordinary surfaces remain secret-free.

**Independent Test**: Fire and delete a trigger, then inspect history, Activity, logs, events, errors, and ordinary API output.

- [x] T029 [US4] Add failing run persistence and deleted-trigger provenance tests in `internal/store/runs_test.go` and `internal/engine/triggers_test.go`
- [x] T030 [US4] Add failing cross-surface secret redaction tests in `internal/api/server/triggers_test.go`, `internal/events/broker_test.go`, and `gui/viewmodel/model_test.go`
- [x] T031 [US4] Persist and expose `source_trigger_id` in run storage and API history in `internal/store/runs.go`, `internal/api/server/runs.go`, and related scan helpers
- [x] T032 [US4] Display external-trigger provenance in Activity without exposing keys in `gui/activity.go` and `gui/viewmodel/model.go`
- [x] T033 [US4] Audit and harden logging, events, and errors against raw-key disclosure across `internal/engine/`, `internal/api/server/`, and `internal/cli/`

**Checkpoint**: Trigger auditability remains intact after lifecycle changes without secret disclosure.

## Phase 6: Documentation and Verification

**Purpose**: Finish user guidance, traceability, and repository-wide evidence.

- [x] T034 Update CLI, GUI, architecture, and task-readiness documentation in `README.md` and `docs/`
- [x] T035 Update S054 issue traceability and verification evidence in `specs/054-external-trigger-lifecycle/quickstart.md` and `specs/README.md`
- [x] T036 Run focused tests for store, task, engine, API, CLI, events, view model, and GUI packages
- [x] T037 Run `sh scripts/verify.sh all` and record all eight mandatory gate results
- [x] T038 Perform UTF-8, mojibake, em-dash, and GitHub Markdown line-wrap audits on changed files
- [x] T039 Review acceptance criteria for issues #132 and #133 and report any incomplete item before closing either issue

## Dependencies and Execution Order

- Phase 1 blocks every user story.
- User Story 1 establishes fire semantics used by all later stories.
- User Story 2 depends on the lifecycle API introduced with User Story 1.
- User Story 3 depends on the typed client and lifecycle API from User Story 2.
- User Story 4 spans the established execution and interface paths and follows their completion.
- Documentation and verification follow all desired stories.

## Implementation Strategy

Complete the phases sequentially under autopilot, with each test task producing a failing test before its paired implementation. Keep trigger sets and filesystem watchers outside the diff. Stop only for a genuine external blocker, then push the complete slice, open one pull request, resolve every review finding, and request at most one second Codex review round.
