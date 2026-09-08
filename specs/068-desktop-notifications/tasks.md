# Tasks: Desktop Notification Management

**Input**: Design documents from `/specs/068-desktop-notifications/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/desktop-notifications-bridge.md`

**Tests**: Required by FR-025 and SC-006. Tests are authored before or with the implementation they constrain and must pass before a task is complete.

## Phase 1: Setup and contracts

- [x] T001 Confirm issue #160 traceability and complete specification, clarification, plan, checklist, contract, and task artifacts in `specs/068-desktop-notifications/`
- [x] T002 [P] Add secret-free channel, scope, policy, delivery, draft, workspace, and result models in `desktop/notifications/model.go`
- [x] T003 [P] Define and adapt existing daemon notification, task, and group clients behind a bounded backend interface in `desktop/notifications/local.go`

## Phase 2: Foundational desktop service

- [x] T004 Add complete workspace composition, ordering, redaction-by-type, and delivery-state mapping in `desktop/notifications/service.go`
- [x] T005 Add channel validation, create, repair, authorization replacement, enable, disable, test, delete, and complete refresh operations in `desktop/notifications/service.go`
- [x] T006 Add selected task/group direct and effective policy loading and atomic replacement in `desktop/notifications/service.go`
- [x] T007 Add failure-first service tests for redaction, bounded calls, complete-snapshot failure, validation, mutation refresh, inheritance, state mapping, and safe errors in `desktop/notifications/service_test.go`

## Phase 3: User Story 1 - Configure and repair channels (Priority: P1)

**Goal**: Deliver a safe complete channel-management workflow without CLI use or stored-secret disclosure.

**Independent Test**: Create, test, disable, repair, re-enable, and remove channels while proving protected fields are never returned or redisplayed.

- [x] T008 [US1] Bind the notification service and channel methods through Wails in `desktop/main.go` and `desktop/app.go`
- [x] T009 [P] [US1] Add frontend notification contracts and native bridge in `desktop/frontend/src/notifications/model.ts` and `desktop/frontend/src/notifications/bridge.ts`
- [x] T010 [US1] Add request-sequenced workspace loading, event debounce, pending-action suppression, and secret-draft clearing in `desktop/frontend/src/notifications/store.ts`
- [x] T011 [US1] Add accessible channel list, create/edit form, explicit replacement controls, test, state toggle, and confirmed removal in `desktop/frontend/src/notifications/NotificationsPage.tsx`
- [x] T012 [US1] Add failure-first Go and React tests for channel lifecycle, disabled tests, replace-versus-retain behavior, authorization clearing, secret absence, focus, and announcements in `desktop/notifications/service_test.go` and `desktop/frontend/src/notifications/NotificationsPage.test.tsx`

## Phase 4: User Story 2 - Assign understandable policies (Priority: P1)

**Goal**: Configure task and group outcomes while making override and inheritance meaning explicit.

**Independent Test**: Save a group failure policy, inspect ancestor inheritance on a task, add a direct override, deliberately enable success, and clear the override.

- [x] T013 [US2] Expose selected-scope policy load and replacement methods in `desktop/app.go` and `desktop/frontend/src/notifications/bridge.ts`
- [x] T014 [US2] Add stale-safe policy selection, draft initialization, replacement, and refreshed effective policy state in `desktop/frontend/src/notifications/store.ts`
- [x] T015 [US2] Add hierarchy-aware scope selection, per-channel outcomes, direct/effective explanation, clear override, and success-volume warning in `desktop/frontend/src/notifications/NotificationsPage.tsx`
- [x] T016 [US2] Add failure-first Go and React tests for task override, immediate and ancestor group inheritance, no policy, invalid assignments, disabled channels, and success warnings in `desktop/notifications/service_test.go` and `desktop/frontend/src/notifications/NotificationsPage.test.tsx`

## Phase 5: User Story 3 - Understand delivery progress and failures (Priority: P2)

**Goal**: Deliver bounded redacted evidence that unmistakably separates tests from task outcomes and explains every delivery state.

**Independent Test**: Filter and inspect 200 mixed records across all event and lifecycle states while preserving selection and the last complete snapshot.

- [x] T017 [US3] Add channel and state filters, typed rows, stable delivery selection, and safe diagnostic detail in `desktop/frontend/src/notifications/NotificationsPage.tsx`
- [x] T018 [US3] Add responsive workspace, form, policy, history, state, and detail styling in `desktop/frontend/src/styles.css`
- [x] T019 [US3] Add failure-first React store and page tests for test/production distinction, queued/retrying/sending/success/failure, filtering, stale responses, failed refresh preservation, stable selection, and 200 records in `desktop/frontend/src/notifications/store.test.ts` and `desktop/frontend/src/notifications/NotificationsPage.test.tsx`

## Phase 6: Integration and verification

- [x] T020 Add the Notifications route, navigation, application wiring, and facade safety coverage in `desktop/frontend/src/connection/model.ts`, `desktop/frontend/src/components/Shell.tsx`, `desktop/frontend/src/App.tsx`, `desktop/frontend/src/App.test.tsx`, and `desktop/app_test.go`
- [x] T021 Add Chromium keyboard, axe, secret-absence, full workflow, 200-record, and 80-through-200-percent zoom coverage in `desktop/frontend/e2e/notifications.spec.ts`
- [x] T022 Update Unreleased change notes and architecture rationale in `CHANGELOG.md`
- [x] T023 Run focused Go race, frontend unit/build, Chromium, and native Wails verification from `quickstart.md`
- [x] T024 Run spec-kit analysis and remediate every finding until the report is clean
- [x] T025 Run canonical `scripts/verify.sh all` in the foreground and record objective evidence in `specs/068-desktop-notifications/verification.md`
- [x] T026 Scan changed files for mojibake, protected-value leakage, Unicode em dashes, and publication formatting defects
- [x] T027 Mark required tasks complete, transition the specification to Implemented, and synchronize `specs/README.md`

## Dependencies and Execution Order

- Phase 1 establishes the response and backend contracts used by all workspaces.
- Phase 2 blocks every user story and proves the secret boundary before Wails binding.
- User Story 1 establishes workspace state and mutation mechanics used by policy administration.
- User Story 2 depends on channel summaries and the shared store but remains independently testable with fixtures.
- User Story 3 depends only on the complete workspace and can be tested independently from mutations.
- Phase 6 follows all required stories.

## Implementation Strategy

Build failure-first backend tests and the secret-free service boundary, then deliver the channel workflow, policy explanation, and history investigation in user-story order. Keep the API and persistence unchanged, validate each story with focused tests, and use the canonical repository gate as the final definition of green.
