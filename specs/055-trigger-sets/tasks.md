# Tasks: Trigger Sets

**Input**: Design documents from `specs/055-trigger-sets/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/trigger-set-contract.md`, and `quickstart.md`

**Tests**: Required by the project constitution. Behavioral tests precede implementation tasks within each story.

## Phase 1: Setup

**Purpose**: Establish traceability and the feature implementation surface.

- [x] T001 Register S055 as Draft in `specs/README.md` and confirm `.specify/feature.json` points to `specs/055-trigger-sets`
- [x] T002 [P] Add the S055 unreleased feature and architecture decision entries to `CHANGELOG.md`
- [x] T003 [P] Verify Go and universal artifact patterns remain covered in `.gitignore`

---

## Phase 2: Foundational Persistence

**Purpose**: Add the shared durable model and transaction primitives that block every user story.

- [x] T004 Add Trigger Set and member metadata types with documentation in `internal/domain/domain.go`
- [x] T005 [P] Add migration-v13 compatibility tests from schema version 12 in `internal/store/migration_v13_test.go`
- [x] T006 [P] Add store tests for count bounds, unique ordered positions, standalone compatibility, cascade deletion, and transaction rollback in `internal/store/trigger_sets_test.go`
- [x] T007 Add schema migration v13 for Trigger Sets and nullable member metadata in `internal/store/store.go`
- [x] T008 Implement transaction-safe Trigger Set create, get, list, reveal, retarget, enable, disable, rotate, and delete operations in `internal/store/trigger_sets.go`
- [x] T009 Update individual trigger scan, list, update, and delete behavior for optional set membership and final-member cleanup in `internal/store/triggers.go`
- [x] T010 Preserve automatic-source readiness across member and set mutations with regression coverage in `internal/store/trigger_sets_test.go`
- [x] T011 Run focused persistence and migration tests for `internal/store`

**Checkpoint**: Durable Trigger Sets exist, preserve S054 triggers, and enforce all invariants atomically.

---

## Phase 3: User Story 1 - Create and Copy a Trigger Set (Priority: P1)

**Goal**: Create 1 through 99 uniquely keyed members atomically and retrieve exact ordered commands.

**Independent Test**: Create sets at counts 1 and 99, reveal them, and prove exact membership, uniqueness, ordering, target equality, redaction, and rollback.

- [x] T012 [P] [US1] Add API contract tests for create bounds, ordered secret output, ordinary redaction, and invalid-request rollback in `internal/api/server/trigger_sets_test.go`
- [x] T013 [P] [US1] Add typed client contract tests for create, list, show, and reveal in `internal/api/client/client_test.go`
- [x] T014 [US1] Define redacted and secret Trigger Set request and response contracts in `internal/api/server/trigger_sets.go`
- [x] T015 [US1] Register Trigger Set create, list, show, and reveal routes in `internal/api/server/server.go`
- [x] T016 [US1] Implement Trigger Set create, list, show, and reveal client methods in `internal/api/client/methods.go`
- [x] T017 [US1] Run focused server and client tests for User Story 1

**Checkpoint**: Headless callers can atomically create and explicitly copy an ordered Trigger Set without ordinary key disclosure.

---

## Phase 4: User Story 2 - Administer Members Independently (Priority: P1)

**Goal**: Keep every existing individual lifecycle action isolated while enforcing set target and empty-set invariants.

**Independent Test**: Mutate one member in 2-member and 99-member sets and compare every sibling field before and after.

- [x] T018 [P] [US2] Add API regression tests for member rename, state change, rotation, deletion, forbidden target changes, and final-member cleanup in `internal/api/server/trigger_sets_test.go`
- [x] T019 [P] [US2] Add CLI regression tests for individual set-member isolation and actionable retarget guidance in `internal/cli/trigger_test.go`
- [x] T020 [US2] Enforce set-member target immutability through the existing individual API handler in `internal/api/server/triggers.go`
- [x] T021 [US2] Extend ordinary trigger responses and CLI list or show output with optional redacted set name and position in `internal/api/server/triggers.go` and `internal/cli/trigger.go`
- [x] T022 [US2] Run focused store, API, and CLI tests for User Story 2

**Checkpoint**: One member can be managed safely without sibling drift or a broken set invariant.

---

## Phase 5: User Story 3 - Administer a Complete Set (Priority: P2)

**Goal**: Retarget, enable, disable, rotate, and delete a complete set atomically.

**Independent Test**: Exercise every broad action at 99 members, inject failures, and prove all-or-nothing state plus old-key invalidation.

- [x] T023 [P] [US3] Add server tests for set-level retarget, state, rotate, delete, rollback, and standard error envelopes in `internal/api/server/trigger_sets_test.go`
- [x] T024 [P] [US3] Add redacted Trigger Set event tests and one-event-per-operation coverage in `internal/events/trigger_set_test.go`
- [x] T025 [US3] Implement set-level API handlers and stable validation errors in `internal/api/server/trigger_sets.go`
- [x] T026 [US3] Add redacted Trigger Set lifecycle event identity and publishing in `internal/events/broker.go` and `internal/api/server/publish.go`
- [x] T027 [US3] Implement typed set-level client methods in `internal/api/client/methods.go`
- [x] T028 [US3] Add maximum-size set operation benchmarks and enforce the one-second budget in `internal/store/trigger_sets_test.go`
- [x] T029 [US3] Run focused store, API, client, and event tests for User Story 3

**Checkpoint**: Every broad lifecycle action is atomic, observable once, and remains below its nominal budget.

---

## Phase 6: User Story 4 - Use Trigger Sets from CLI and Desktop (Priority: P2)

**Goal**: Provide complete machine-friendly and desktop administration through the shared local API.

**Independent Test**: Complete the lifecycle independently through CLI and headless GUI tests, including exact output, confirmations, live refresh, accessibility labels, and light or dark readability contracts.

- [x] T030 [P] [US4] Add CLI tests for every nested `trigger set` command, JSON parity, exact line output, errors, and exit behavior in `internal/cli/trigger_set_test.go`
- [x] T031 [P] [US4] Add headless GUI tests for membership columns, create, bulk copy, broad confirmations, selected-set actions, and status feedback in `gui/triggers_test.go`
- [x] T032 [US4] Implement nested Trigger Set CLI commands and exact ordered output in `internal/cli/trigger_set.go` and register them in `internal/cli/trigger.go`
- [x] T033 [US4] Extend the GUI Backend and view model with Trigger Set lifecycle data in `gui/app.go`, `gui/viewmodel/viewmodel.go`, and their tests
- [x] T034 [US4] Add Set and Position columns plus standalone fallbacks to the Triggers table in `gui/triggers.go`
- [x] T035 [US4] Implement New Set, bulk copy, retarget, enable or disable, rotate, and delete desktop actions with scoped confirmations in `gui/triggers.go`
- [x] T036 [US4] Refresh desktop state once for redacted Trigger Set events and preserve individual trigger refresh behavior in `gui/app.go`
- [x] T037 [US4] Run focused CLI and headless GUI tests for User Story 4

**Checkpoint**: Operators can administer Trigger Sets completely through either supported interface.

---

## Phase 7: Documentation and Completion

**Purpose**: Close traceability, validate all requirements, and record reproducible evidence.

- [x] T038 [P] Update Trigger Set API, CLI, GUI, and architecture guidance in `docs/api.md`, `docs/cli.md`, `docs/gui-fields.md`, `docs/architecture.md`, and `README.md`
- [x] T039 [P] Add deterministic end-to-end validation evidence and redaction results to `specs/055-trigger-sets/verification.md`
- [x] T040 Run focused persistence, migration, API, client, CLI, event, and GUI tests and record results in `specs/055-trigger-sets/verification.md`
- [x] T041 Run `go run ./scripts/github-format` and correct all em dash or hard-wrap defects before publication
- [x] T042 Run `sh scripts/verify.sh all` in the foreground and record every mandatory gate in `specs/055-trigger-sets/verification.md`
- [x] T043 Re-run the spec lifecycle automation, mark every task complete, advance S055 to Implemented in `spec.md` and `specs/README.md`, and record delivery evidence
- [x] T044 Commit the complete slice locally as `feat(055): add Trigger Sets` with the project co-author trailer

---

## Dependencies and Execution Order

- Phase 1 precedes every implementation phase.
- Phase 2 blocks every user story because all interfaces require the shared persisted transaction model.
- User Story 1 establishes set identity and explicit secret retrieval.
- User Story 2 depends on membership metadata but remains independently testable through existing individual trigger endpoints.
- User Story 3 depends on the foundational transaction methods and adds broad lifecycle behavior.
- User Story 4 depends on the API and client contracts from User Stories 1 through 3.
- Phase 7 follows all user stories and is the only path to the Implemented state.

## Parallel Opportunities

- T002 and T003 can run while foundational test surfaces are prepared.
- T005 and T006 cover separate migration and store behavior files.
- T012 and T013 cover server and client contracts independently.
- T018 and T019 cover API and CLI member isolation independently.
- T023 and T024 cover handler and event contracts independently.
- T030 and T031 cover CLI and GUI behavior independently.
- T038 and T039 touch independent documentation and verification files.

## Implementation Strategy

First establish migration v13 and the transaction-safe domain model. Deliver creation and explicit bulk copy before member isolation and broad lifecycle operations. Add CLI and desktop interfaces only after the local API contract is stable. Finish with complete redaction, rollback, maximum-size performance, documentation, lifecycle, and eight-gate verification evidence.
