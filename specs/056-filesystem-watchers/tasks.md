# Tasks: Filesystem Watchers

**Input**: Design documents from `specs/056-filesystem-watchers/`
**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/filesystem-watcher-contract.md`, `quickstart.md`
**Tests**: Mandatory under constitution principle II and the S056 specification.

## Phase 1: Setup and Governance

**Purpose**: Register the slice, dependency decision, and publication constraints before implementation.

- [x] T001 Confirm clean synchronized `main`, create `codex/056-filesystem-watchers`, activate `specs/056-filesystem-watchers` in `.specify/feature.json`, and register Draft state in `specs/README.md`
- [x] T002 [P] Add the S056 feature and architecture decisions to `CHANGELOG.md`
- [x] T003 [P] Promote `github.com/fsnotify/fsnotify` v1.10.0 to a direct dependency in `go.mod` and `go.sum`
- [x] T004 [P] Confirm generated databases, binaries, watcher fixtures, and temporary output remain covered by `.gitignore`
- [x] T005 Run the Spec Kit cross-artifact analysis, resolve every critical or high finding and coverage gap, and record metrics in `specs/056-filesystem-watchers/verification.md`

---

## Phase 2: Foundational Persistence and Contracts

**Purpose**: Establish the durable definition, migration, provenance, and transaction rules required by every interface and runtime story.

- [x] T006 Add documented watcher kind, configuration, health, and run-provenance fields in `internal/domain/domain.go`
- [x] T007 [P] Add migration-v14 tests from a real schema-version-13 fixture in `internal/store/migration_v14_test.go`
- [x] T008 [P] Add watcher CRUD, validation, task-cascade, activation-readiness, and 100-definition budget tests in `internal/store/watchers_test.go`
- [x] T009 Add schema migration v14 for `filesystem_watchers` and `runs.source_watcher_id` in `internal/store/store.go`
- [x] T010 Implement transaction-safe watcher create, get, list, update, state, and delete methods in `internal/store/watchers.go`
- [x] T011 Extend run scans and writes for watcher provenance in `internal/store/runs.go` and related store tests
- [x] T012 Extend automatic activation-source checks and task mutation cleanup for enabled watchers in `internal/store/crud.go`
- [x] T013 Run focused migration, watcher persistence, task readiness, and run provenance tests in `internal/store`

**Checkpoint**: Valid watcher definitions and provenance persist without changing existing data or trigger-ready task behavior.

---

## Phase 3: User Story 1 - Stable File Dispatch (Priority: P1)

**Goal**: Convert portable filesystem hints into one stable, provenance-bearing run request through the normal dispatcher.

**Independent Test**: Exercise exact-file and directory matching, write storms, rename into place, recursion, stability, generation cancellation, task eligibility, and overlap behavior with deterministic timing.

- [x] T014 [P] [US1] Add fake-observer and injected-clock tests for selection, debounce, stability, generation cancellation, and recursion in `internal/watcher/manager_test.go`
- [x] T015 [P] [US1] Add real temporary-filesystem integration tests for direct writes, atomic replacement, nested directories, and link exclusion in `internal/watcher/observer_integration_test.go`
- [x] T016 [P] [US1] Add engine tests for filesystem watcher eligibility, overlap reuse, trigger classification, and source watcher provenance in `internal/engine/watchers_test.go`
- [x] T017 [US1] Implement the production fsnotify adapter and portable event representation in `internal/watcher/observer.go`
- [x] T018 [US1] Implement the generation-scoped watcher event loop, matching, debounce, stability, recursion, and dispatcher interface in `internal/watcher/manager.go`
- [x] T019 [US1] Refactor shared external-source task validation and add filesystem watcher dispatch in `internal/engine/triggers.go` and `internal/engine/watchers.go`
- [x] T020 [US1] Start, reload, and terminate the watcher manager with the engine lifecycle in `internal/engine/engine.go`
- [x] T021 [US1] Wire filesystem run and watcher lifecycle publication through `internal/events/broker.go`, `internal/api/server/publish.go`, and `cmd/goschedd/main.go`
- [x] T022 [US1] Run focused watcher runtime, real filesystem, engine, and event tests including 100 write-storm and rename trials

**Checkpoint**: Supported stable filesystem outcomes request exactly one normal task run with truthful watcher provenance.

---

## Phase 4: User Story 2 - Complete Administration (Priority: P1)

**Goal**: Configure and administer watchers through local API, CLI, and desktop interfaces without daemon restart.

**Independent Test**: Complete create, list, show, update, enable, disable, and delete independently through API, CLI, and headless GUI while checking errors, JSON parity, events, health, and live reload.

- [x] T023 [P] [US2] Add API lifecycle, validation, health-join, event-count, and reload contract tests in `internal/api/server/watchers_test.go`
- [x] T024 [P] [US2] Add typed client lifecycle tests in `internal/api/client/watchers_test.go`
- [x] T025 [P] [US2] Add human and JSON CLI lifecycle, duration, validation, and exit-status tests in `internal/cli/watcher_test.go`
- [x] T026 [P] [US2] Add view-model refresh and live-event tests for watchers in `gui/viewmodel/viewmodel_test.go`
- [x] T027 [P] [US2] Add headless GUI table, form, action, health-detail, confirmation, accessibility, and light or dark appearance tests in `gui/triggers_test.go`
- [x] T028 [US2] Implement watcher request and response contracts, lifecycle handlers, health joins, and routes in `internal/api/server/watchers.go` and `internal/api/server/server.go`
- [x] T029 [US2] Implement typed watcher lifecycle methods in `internal/api/client/methods.go`
- [x] T030 [US2] Implement `gosched watcher` lifecycle commands and register them in `internal/cli/watcher.go` and `internal/cli/root.go`
- [x] T031 [US2] Extend the GUI backend and view model with watcher state in `gui/app.go` and `gui/viewmodel/viewmodel.go`
- [x] T032 [US2] Add a separate Filesystem Watchers table and complete toolbar and editor lifecycle to `gui/triggers.go`
- [x] T033 [US2] Refresh desktop state on watcher lifecycle and health events in `gui/app.go`
- [x] T034 [US2] Run focused API, client, CLI, view-model, and headless GUI lifecycle tests

**Checkpoint**: Both supported user interfaces completely administer current watcher definitions and runtime state.

---

## Phase 5: User Story 3 - Failure and Recovery (Priority: P2)

**Goal**: Degrade and recover predictably across missing, replaced, inaccessible, unsupported, or overflowed observation without flooding or replay.

**Independent Test**: Drive health transitions with a fake observer and real roots, including daemon restart, while proving bounded recovery, no replay, transition deduplication, and complete shutdown.

- [x] T035 [P] [US3] Add deterministic missing-root, replacement, observer-error, overflow, deduplication, and bounded-recovery tests in `internal/watcher/manager_test.go`
- [x] T036 [P] [US3] Add restart, no-replay, reload race, and shutdown leak tests in `internal/engine/watchers_test.go`
- [x] T037 [P] [US3] Add health event redaction and transition-only publication tests in `internal/events/watcher_test.go`
- [x] T038 [US3] Implement degraded health transitions, full observer rebuild, two-second recovery scheduling, reason bounding, and repeated-state suppression in `internal/watcher/manager.go`
- [x] T039 [US3] Expose engine health snapshots and publish transition events without paths in `internal/engine/watchers.go` and `internal/events/broker.go`
- [x] T040 [US3] Run focused race-enabled recovery, restart, redaction, and lifecycle termination tests

**Checkpoint**: Watchers truthfully expose failure, recover prospectively, and leave no runtime resources after shutdown.

---

## Phase 6: Documentation and Completion

**Purpose**: Close user guidance, traceability, analysis, formatting, and canonical verification.

- [x] T041 [P] Document watcher concepts, platform limits, API, CLI, GUI fields, run provenance, and architecture in `README.md`, `docs/api.md`, `docs/cli.md`, `docs/gui-fields.md`, and `docs/architecture.md`
- [x] T042 [P] Add implementation and deterministic evidence to `specs/056-filesystem-watchers/verification.md`
- [x] T043 Run all focused suites, race tests, cross-platform compile checks, lifecycle checks, and the 100-trial reliability suites and record exact results in `specs/056-filesystem-watchers/verification.md`
- [x] T044 Run `go run ./scripts/github-format` and correct every em dash or hard-wrapped Markdown defect before publication
- [x] T045 Run `sh scripts/verify.sh all` in the foreground and record all eight canonical gate results and coverage in `specs/056-filesystem-watchers/verification.md`
- [x] T046 Mark every required task complete, advance S056 to Implemented in `spec.md` and `specs/README.md`, rerun lifecycle and corruption checks, and record delivery evidence
- [x] T047 Commit the complete slice locally as `feat(056): add filesystem watchers` with the project co-author trailer

---

## Dependencies and Execution Order

- Phase 1 precedes all implementation.
- Phase 2 blocks every user story because runtime, API, CLI, and GUI all require one durable contract.
- User Story 1 establishes the runtime and shared dispatcher integration.
- User Story 2 depends on the runtime reload and health boundary but remains independently testable through each interface.
- User Story 3 extends the same event loop with degraded recovery and proves lifecycle safety.
- Phase 6 follows all user stories and is the only path to Implemented state.

## Parallel Opportunities

- T002 through T004 touch independent setup surfaces.
- T006 and T007 cover separate migration and CRUD fixtures.
- T013 through T015 cover runtime, native integration, and engine contracts independently.
- T022 through T026 cover API, client, CLI, view-model, and GUI contracts independently.
- T034 through T036 cover manager recovery, engine lifecycle, and event redaction independently.
- T040 and T041 touch user guidance and evidence independently.

## Implementation Strategy

Establish schema v14 and activation-readiness invariants first. Build the watcher state machine behind fake observer and injected-clock boundaries before attaching fsnotify and engine dispatch. Add API and client contracts, then CLI and desktop administration. Finish with degraded recovery, transition deduplication, platform contract tests, documentation, full analysis, formatting, and every canonical verification gate.
