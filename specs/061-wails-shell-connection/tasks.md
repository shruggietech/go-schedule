# Tasks: Production Wails Shell and Local Connection

**Input**: Design documents from `specs/061-wails-shell-connection/`

**Tests**: Required by the specification and constitution. Behavioral and concurrency tests precede implementation.

**Organization**: Tasks are chronological and grouped by independently testable user story. Publication and review are delivery evidence, not implementation tasks.

## Phase 1: Setup

**Purpose**: Establish the production module and lifecycle inventory without changing release inputs.

- [x] T001 Confirm clean synchronized `main`, create `codex/061-wails-shell-connection`, activate `.specify/feature.json`, and register Draft state in `specs/README.md`
- [x] T002 Create the isolated production Wails module metadata in `desktop/go.mod`, `desktop/wails.json`, and `desktop/frontend/package.json`
- [x] T003 [P] Copy approved local icon and font assets into `desktop/build/` and `desktop/frontend/public/` and register the consumers in `brand/repository-consumers.json`
- [x] T004 [P] Configure TypeScript, Vite, Vitest, Playwright, and ignored outputs in `desktop/frontend/` and `.gitignore`

---

## Phase 2: Foundational Models and Test Seams

**Purpose**: Define safe contracts and deterministic lifecycle seams that block every story.

- [x] T005 Define documented connection state, target, snapshot, event, action, backend, observer, and retry scheduler contracts in `desktop/connection/model.go`
- [x] T006 Define closed TypeScript connection, route, appearance, notification, and bridge models in `desktop/frontend/src/connection/model.ts`
- [x] T007 Add deterministic backend, scheduler, observer, and native-runtime fakes in `desktop/connection/test_helpers_test.go` and `desktop/app_test.go`

**Checkpoint**: Safe bridge vocabulary and no-sleep test seams are ready.

---

## Phase 3: User Story 1 - Open the Local Control Center Safely (Priority: P1)

**Goal**: Start one production-intent shell against This computer through protected local IPC.

**Independent Test**: Connected and unavailable fakes plus the real local adapter produce honest sanitized snapshots within the bounded attempt window and no network listener.

### Tests for User Story 1

- [x] T008 [P] [US1] Write failing target, compatibility, capability, initial snapshot, safe unavailable, and real local-adapter construction tests in `desktop/connection/local_test.go`
- [x] T009 [P] [US1] Write failing Wails facade startup, snapshot, native quit result, and sanitized bridge tests in `desktop/app_test.go`
- [x] T010 [P] [US1] Write failing frontend connection-store generation and browser fallback tests in `desktop/frontend/src/connection/store.test.ts`

### Implementation for User Story 1

- [x] T011 [US1] Implement the protected This computer adapter, compatibility policy, capability manifest, and safe failure mapping in `desktop/connection/local.go`
- [x] T012 [US1] Implement initial manager state, bounded health negotiation, immutable snapshots, and observer publication in `desktop/connection/manager.go`
- [x] T013 [US1] Implement the Wails facade and runtime adapter in `desktop/app.go` and `desktop/main.go`
- [x] T014 [US1] Implement the isolated production application metadata and embedded asset entry point in `desktop/wails.json` and `desktop/main.go`
- [x] T015 [US1] Implement the frontend bridge adapter and generation-aware connection store in `desktop/frontend/src/connection/bridge.ts` and `desktop/frontend/src/connection/store.ts`

**Checkpoint**: The production foundation starts, identifies This computer, and reports connected or unavailable state through a sanitized bridge.

---

## Phase 4: User Story 2 - Understand and Recover Connection Failures (Priority: P1)

**Goal**: Distinguish failures and own automatic retry, manual retry, event, stale-generation, and shutdown lifecycles.

**Independent Test**: Controlled fakes drive all eight states, exact retry progression, coalesced manual retry, event-only degradation, recovery, stale callbacks, and 100 shutdown cycles without sleeps or leaked work.

### Tests for User Story 2

- [x] T016 [P] [US2] Write failing state transition and safe category mapping tests in `desktop/connection/manager_test.go`
- [x] T017 [P] [US2] Write failing exact retry cadence, manual retry interruption, and coalescing tests in `desktop/connection/manager_test.go`
- [x] T018 [P] [US2] Write failing event degradation, recovery, stale-generation rejection, and sanitized event tests in `desktop/connection/manager_test.go`
- [x] T019 [P] [US2] Write failing cancellation, no-new-work-after-close, and 100-cycle lifecycle tests in `desktop/connection/manager_test.go`
- [x] T020 [P] [US2] Write failing frontend failure guidance, retry action, stale snapshot, and live announcement tests in `desktop/frontend/src/connection/store.test.ts`

### Implementation for User Story 2

- [x] T021 [US2] Implement the single-owner generation loop, child cancellation, and stale-result guard in `desktop/connection/manager.go`
- [x] T022 [US2] Implement bounded automatic retry, manual retry interruption, terminal failure waiting, and deterministic scheduler use in `desktop/connection/manager.go`
- [x] T023 [US2] Implement one generation-owned event stream, domain event sanitization, degradation, recovery, and retry reset in `desktop/connection/manager.go`
- [x] T024 [US2] Implement two-second orderly shutdown and rejected post-close actions in `desktop/connection/manager.go` and `desktop/app.go`
- [x] T025 [US2] Implement frontend connection-event subscription, stale-generation rejection, recovery actions, and announcements in `desktop/frontend/src/connection/store.ts`

**Checkpoint**: Connection behavior is complete, deterministic, race-safe, and transport-neutral at the feature boundary.

---

## Phase 5: User Story 3 - Use a Consistent Accessible Application Frame (Priority: P2)

**Goal**: Provide the shared Calm Operations shell and component catalog needed by all migration slices.

**Independent Test**: Component and browser tests cover semantic contracts, keyboard behavior, dialog focus, appearance, reduced motion, every connection state, compact layout, zoom, long content, and zero serious or critical axe findings.

### Tests for User Story 3

- [x] T026 [P] [US3] Write failing primitive behavior tests for buttons, notices, states, fields, tables, disclosures, dialogs, and notifications in `desktop/frontend/src/components/components.test.tsx`
- [x] T027 [P] [US3] Write failing shell navigation, target identity, appearance, retry, exit, and placeholder honesty tests in `desktop/frontend/src/App.test.tsx`
- [x] T028 [P] [US3] Write failing axe, keyboard, focus return, live region, reduced-motion, and offline-asset tests in `desktop/frontend/src/accessibility.test.tsx`
- [x] T029 [P] [US3] Write failing responsive, long-content, connection-state, appearance, keyboard, zoom, and overflow browser contracts in `desktop/frontend/e2e/shell.spec.ts`

### Implementation for User Story 3

- [x] T030 [US3] Implement reusable Button, Link, StatusBadge, Notice, StatePanel, Field, DataTable, Disclosure, Dialog, and ToastRegion primitives in `desktop/frontend/src/components/`
- [x] T031 [US3] Implement approved tokens, local fonts, state styling, appearance modes, reduced motion, focus, compact reflow, and zoom behavior in `desktop/frontend/src/styles.css`
- [x] T032 [US3] Implement the Shell landmarks, navigation, target bar, connection panel, page framing, appearance control, notification region, and exact dialog focus lifecycle in `desktop/frontend/src/components/Shell.tsx`
- [x] T033 [US3] Implement honest placeholder routes and compose the production application in `desktop/frontend/src/App.tsx` and `desktop/frontend/src/main.tsx`

**Checkpoint**: Every required shared component contract exists and the production shell is independently usable and accessible.

---

## Phase 6: Cross-Platform Evidence and Delivery

**Purpose**: Make the new production boundary fail closed without altering shipped artifacts.

- [x] T034 Add exact production desktop Go, frontend, browser, offline-asset, and Windows/macOS/Linux native build jobs in `.github/workflows/ci.yml`
- [x] T035 Extend pinned automation policy and fixtures for the production desktop matrix in `scripts/automation-check.sh` and `test/scripts/automation-check_test.sh`
- [x] T036 Document production module ownership, dependency licenses, build prerequisites, connection trust boundary, and migration exclusions in `desktop/README.md`
- [x] T037 Update the S061 architecture decision and feature summary in `CHANGELOG.md`, `CLAUDE.md`, and `specs/README.md`
- [x] T038 Run nested Go race, dependency audit, frontend tests, frontend build, browser contract, exact native Windows build, brand check, lifecycle check, encoding check, and diff integrity and record results in `specs/061-wails-shell-connection/verification.md`
- [x] T039 Run `sh scripts/verify.sh all`, resolve all eight canonical gates, and record exact coverage and environment evidence in `specs/061-wails-shell-connection/verification.md`
- [x] T040 Audit every #151 and #152 acceptance criterion, resolve all tasks, advance S061 to Implemented, and commit as `feat(061): build Wails shell and local connection` with the required co-author trailer

---

## Dependencies and Execution Order

- Phase 1 precedes all production files.
- Phase 2 defines contracts and fakes before behavioral tests.
- US1 establishes negotiation and bridge behavior required by US2 and the shell.
- US2 completes lifecycle behavior before the shell presents recovery controls.
- US3 depends only on the safe connection contract and can develop component files in parallel with manager internals after Phase 2.
- Cross-platform workflow changes follow passing local nested gates.

## Parallel Opportunities

- T003 and T004 affect independent asset and tooling files.
- T008, T009, and T010 cover independent Go bridge and frontend store surfaces.
- T016 through T020 can be authored in separate test files or sections before manager integration.
- T026 through T029 cover component, shell, accessibility, and browser contracts independently.

## Implementation Strategy

The MVP is US1 plus the foundational models: one buildable production shell boundary with honest local state. US2 then makes the connection trustworthy under failure and concurrency. US3 turns the approved direction into reusable application infrastructure. The slice is complete only when all three stories and hosted production build contracts are present; no partial story closes #151 or #152.

## Completion Rule

All 40 tasks must be checked. Push, pull-request creation, hosted review, merge, and cleanup remain publication evidence and do not appear as implementation tasks. The operator has explicitly authorized this branch and pull request plus verified in-scope fixes through at most two Codex review rounds; merge remains the maintainer's decision.
