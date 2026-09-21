# Tasks: Cross-Daemon Search and Target-Safe Actions

**Input**: Design documents from `/specs/097-cross-daemon-search/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Tests are mandatory under the project constitution and are written before the behavior they cover.

**Organization**: Tasks are grouped by user story so discovery, exact-source opening, mutation safety, and accessible scale remain independently reviewable.

## Phase 1: Setup and Lifecycle

**Purpose**: Activate the reviewed S097 specification and traceability surfaces.

- [x] T001 Mark S097 Ready in `specs/097-cross-daemon-search/spec.md` and add its lifecycle row to `specs/README.md`
- [x] T002 Record S097 as In progress on issue #183's GitHub Project item and preserve parent #174 traceability
- [x] T003 Verify existing ignore files cover Go, Node, generated frontend, and Wails build artifacts in `.gitignore`, `desktop/frontend/.gitignore`, and applicable tool configuration

---

## Phase 2: Foundational Daemon Search Contract

**Purpose**: Build the bounded, Observe-authorized, transport-parity contract required by every story.

**Critical**: User-story implementation starts only after this phase passes focused tests.

- [x] T004 [P] Write safe search entity and deterministic ordering tests in `internal/domain/search_test.go`
- [x] T005 [P] Write bounded task, group, failure, schedule-source, and alert query tests with secret-exclusion assertions in `internal/store/search_test.go`
- [x] T006 Implement the allowlisted search models in `internal/domain/search.go`
- [x] T007 Implement normalized, escaped, deterministic, limit-plus-one SQLite search in `internal/store/search.go`
- [x] T008 [P] Write request validation, Observe authorization, occurrence, redaction, and error contract tests in `internal/api/server/search_test.go` and `internal/api/server/access_test.go`
- [x] T009 Implement `GET /v1/search` and schedule occurrence projection in `internal/api/server/search.go`, `internal/api/server/server.go`, `internal/authorization/catalog.go`, and `internal/api/server/manifest.go`
- [x] T010 [P] Write local client and remote generated-client contract tests in `internal/api/client/search_test.go` and `internal/remote/server_test.go`
- [x] T011 Add the search operation to `api/openapi/remote-v1.yaml`, regenerate `internal/remote/openapi.gen.go`, and implement the transport client in `internal/api/client/search.go`
- [x] T012 Run focused race tests for `internal/domain`, `internal/store`, `internal/api/server`, `internal/api/client`, `internal/authorization`, and `internal/remote`

**Checkpoint**: One daemon can answer a safe, bounded, versioned search over both supported transports.

---

## Phase 3: User Story 1 - Find Automation Across Registered Daemons (Priority: P1)

**Goal**: Submit one explicit query and receive progressive, source-labeled results and failures from every registration.

**Independent Test**: Search duplicate names across mixed reachable, unavailable, Observe-only, and older daemons, then inspect source identity, freshness, truncation, and target-specific failures.

### Tests for User Story 1

- [x] T013 [P] [US1] Write fan-out, eight-worker bound, three-second deadline, cancellation, edited-registration, duplicate-name, and partial-failure tests in `desktop/search/service_test.go`
- [x] T014 [P] [US1] Write query validation, filtering, truncation, progressive announcement, and duplicate-identity component tests in `desktop/frontend/src/search/SearchPage.test.tsx`

### Implementation for User Story 1

- [x] T015 [US1] Implement registration snapshots, temporary identity-pinned clients, generation-scoped fan-out, and action availability in `desktop/search/model.go` and `desktop/search/service.go`
- [x] T016 [US1] Expose `SearchAcrossSystems` and `search:event` through `desktop/app.go`, `desktop/main.go`, and `desktop/app_test.go`
- [x] T017 [US1] Add typed bridge models and event subscriptions in `desktop/frontend/src/search/model.ts`, `desktop/frontend/src/search/bridge.ts`, and `desktop/frontend/src/connection/model.ts`
- [x] T018 [US1] Implement the Search route, query controls, progressive target summaries, filters, deterministic result table, and safe empty or error states in `desktop/frontend/src/search/SearchPage.tsx`, `desktop/frontend/src/App.tsx`, and `desktop/frontend/src/components/Shell.tsx`

**Checkpoint**: User Story 1 provides useful read-only cross-daemon discovery without changing the selected connection.

---

## Phase 4: User Story 2 - Open the Exact Source Safely (Priority: P1)

**Goal**: Open a result only after selecting its immutable registration and confirming the expected daemon identity.

**Independent Test**: Open identically named task, schedule, failure, and alert results from different profiles, then introduce an identity mismatch and confirm navigation fails closed.

### Tests for User Story 2

- [x] T019 [P] [US2] Write exact-registration, expected-identity, stale-result, deleted-profile, and failed-selection tests in `desktop/frontend/src/App.test.tsx` and `desktop/frontend/src/search/SearchPage.test.tsx`

### Implementation for User Story 2

- [x] T020 [US2] Extend search result intents and the existing select-and-wait flow to require the expected daemon identity before navigating in `desktop/frontend/src/search/model.ts`, `desktop/frontend/src/search/SearchPage.tsx`, and `desktop/frontend/src/App.tsx`
- [x] T021 [US2] Preserve task or record context and source labels for Tasks, Schedule, and Activity destinations in `desktop/frontend/src/App.tsx`

**Checkpoint**: User Story 2 never routes by display name and leaves Search intact on identity or connection failure.

---

## Phase 5: User Story 3 - Perform Explicit Target-Safe Actions (Priority: P1)

**Goal**: Confirm and execute acknowledge, enable, disable, or run-now across explicitly selected compatible objects with per-object outcomes.

**Independent Test**: Execute one compatible action across multiple daemons while changing one identity, revoking one authority, deleting one object, and failing one transport after another target succeeds.

### Tests for User Story 3

- [x] T022 [P] [US3] Write selection validation, profile reconstruction, daemon identity, authority, object existence and state, uncertainty, and partial-outcome tests in `desktop/search/actions_test.go`
- [x] T023 [P] [US3] Write compatible selection, grouped confirmation, keyboard cancellation, disabled-action reason, and per-object outcome tests in `desktop/frontend/src/search/SearchPage.test.tsx`

### Implementation for User Story 3

- [x] T024 [US3] Implement exact-target task and alert revalidation plus audited existing-endpoint dispatch in `desktop/search/actions.go`
- [x] T025 [US3] Expose `ExecuteSearchAction` through `desktop/app.go`, wire the service dependencies in `desktop/main.go`, and cover the facade in `desktop/app_test.go`
- [x] T026 [US3] Implement single-action selection rules, target-grouped confirmation, cancellation, pending states, and per-object outcome summaries in `desktop/frontend/src/search/SearchPage.tsx`, `desktop/frontend/src/search/model.ts`, and `desktop/frontend/src/search/bridge.ts`
- [x] T027 [US3] Refresh affected search targets after definitive outcomes while preserving uncertain evidence and unrelated target results in `desktop/search/service.go` and `desktop/frontend/src/search/SearchPage.tsx`

**Checkpoint**: User Story 3 performs no mutation without current identity, authority, object revalidation, and explicit confirmation.

---

## Phase 6: User Story 4 - Search and Act Accessibly at Scale (Priority: P2)

**Goal**: Keep the complete workflow operable by keyboard and assistive technology under supported scale and viewport constraints.

**Independent Test**: Exercise 100 profiles, duplicate labels, progressive and partial outcomes, dialogs, 800 by 600, 200 percent zoom, reduced motion, and keyboard-only operation.

### Tests for User Story 4

- [x] T028 [P] [US4] Add Playwright keyboard, focus, duplicate-label, 100-profile, partial-failure, zoom, and horizontal-overflow coverage in `desktop/frontend/e2e/cross-daemon-search.spec.ts`
- [x] T029 [P] [US4] Add accessible-name, live-region, focus-return, non-color state, and reduced-motion component assertions in `desktop/frontend/src/search/SearchPage.test.tsx`

### Implementation for User Story 4

- [x] T030 [US4] Refine responsive result, confirmation, and outcome layouts plus visible focus and reduced-motion styles in `desktop/frontend/src/styles.css` and `desktop/frontend/src/search/SearchPage.tsx`
- [x] T031 [US4] Add stable focus restoration and concise progressive announcements in `desktop/frontend/src/search/SearchPage.tsx`

**Checkpoint**: All stories remain usable at the supported profile, viewport, zoom, and keyboard boundaries.

---

## Phase 7: Documentation, Reconciliation, and Verification

**Purpose**: Synchronize public contracts, lifecycle evidence, and canonical quality gates.

- [x] T032 [P] Document the daemon search contract and remote authority behavior in `docs/api.md` and `docs/remote-access.md`
- [x] T033 [P] Record S097 behavior and architecture decisions under Unreleased in `CHANGELOG.md`
- [x] T034 Run `go run ./scripts/github-format`, frontend formatting, focused race tests, frontend unit and build checks, Playwright search coverage, generated-contract checks, and `scripts/verify.sh all`
- [x] T035 Record commands and outcomes in `specs/097-cross-daemon-search/verification.md`
- [x] T036 Mark all tasks complete, set S097 Implemented with delivery evidence in `specs/097-cross-daemon-search/spec.md` and `specs/README.md`, and reconcile issue #183 and parent #174 functional checklists

---

## Dependencies and Execution Order

### Phase Dependencies

- **Phase 1**: Starts immediately.
- **Phase 2**: Depends on Phase 1 and blocks all user stories.
- **Phase 3 / US1**: Depends on the daemon contract and establishes shared search state.
- **Phase 4 / US2**: Depends on US1 result intents but remains independently testable.
- **Phase 5 / US3**: Depends on US1 result identity and action availability. It reuses, but does not depend on, global selected-connection navigation from US2.
- **Phase 6 / US4**: Depends on the complete interaction surface.
- **Phase 7**: Depends on all selected stories.

### Parallel Opportunities

- T004 and T005 can run in parallel before domain and store implementation.
- T008 and T010 can be written in parallel once the contract shape is fixed.
- T013 and T014 cover separate Go and TypeScript surfaces.
- T022 and T023 cover separate service and UI action contracts.
- T028 and T029 cover browser and component accessibility independently.
- T032 and T033 touch separate documentation files.

## Implementation Strategy

### MVP First

1. Complete lifecycle setup and the bounded daemon contract.
2. Deliver User Story 1 as read-only cross-daemon discovery.
3. Validate duplicate names, partial failures, cancellation, and result limits.
4. Add exact-source opening and mutations only after source identity is reliable.

### Incremental Delivery

1. Daemon contract establishes one safe search unit.
2. Fleet fan-out establishes progressive discovery without global connection changes.
3. Exact-source opening establishes identity-safe navigation.
4. Explicit actions add current revalidation and independent outcomes.
5. Accessibility and scale refinement validates the combined surface.

## Notes

- Every test task precedes the implementation it covers.
- Search and mutation use immutable registration keys and daemon IDs; display labels never route operations.
- Existing authorization and mutation endpoints remain authoritative.
- No task represents pull-request publication, CI waiting, or review bookkeeping because those are workflow evidence rather than implementation work.
