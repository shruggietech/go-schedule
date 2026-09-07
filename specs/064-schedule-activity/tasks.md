# Tasks: Operational Schedule and Activity

**Input**: Design documents from `/specs/064-schedule-activity/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/desktop-operations-bridge.md`

**Tests**: Required by FR-020 and SC-006. Tests are authored before or with the implementation they constrain and must pass before a task is complete.

## Phase 1: Setup and contracts

- [x] T001 Confirm issue #155 traceability and complete specification, clarification, plan, checklist, contract, and task artifacts in `specs/064-schedule-activity/`
- [x] T002 [P] Add transport-neutral Schedule and Activity models in `desktop/operations/model.go`
- [x] T003 [P] Define and adapt existing daemon clients behind a bounded backend interface in `desktop/operations/local.go`

## Phase 2: Foundational desktop service

- [x] T004 Add schedule-window mapping with explicit prediction and recorded-run identities in `desktop/operations/service.go`
- [x] T005 Add complete Activity workspace composition, deterministic attribute formatting, and safe error mapping in `desktop/operations/service.go`
- [x] T006 Add individual and visible-set acknowledgement operations in `desktop/operations/service.go`
- [x] T007 Add service tests for mapping, ordering, state vocabulary, complete-snapshot failure, exact log path, and acknowledgement in `desktop/operations/service_test.go`

## Phase 3: User Story 1 - Understand scheduled work (Priority: P1)

**Goal**: Deliver fast, accurate agenda and calendar understanding for predictions and recorded runs.

**Independent Test**: Load mixed occurrences, switch ranges and views, select a record, and refresh without losing its identity or focus.

- [x] T008 [US1] Bind the operations service and schedule method through Wails in `desktop/main.go` and `desktop/app.go`
- [x] T009 [P] [US1] Add frontend operations contracts and native bridge in `desktop/frontend/src/operations/model.ts` and `desktop/frontend/src/operations/bridge.ts`
- [x] T010 [US1] Add request-sequenced Schedule loading and debounced live updates in `desktop/frontend/src/operations/store.ts`
- [x] T011 [US1] Replace the Schedule placeholder with accessible agenda, calendar, range, and stable-detail behavior in `desktop/frontend/src/operations/SchedulePage.tsx` and `desktop/frontend/src/App.tsx`
- [x] T012 [US1] Add responsive Schedule table, state, calendar, and detail styling in `desktop/frontend/src/styles.css`
- [x] T013 [US1] Add React tests for prediction distinction, all run states, range preservation, calendar access, stale responses, stable selection, focus, empty state, and 100 rows

## Phase 4: User Story 2 - Investigate activity and failures (Priority: P1)

**Goal**: Deliver one fast recent-activity workspace that preserves the identity and diagnostics of runs, logs, and alerts.

**Independent Test**: Filter mixed records, inspect each type, and verify complete failed-run detail and exact log metadata.

- [x] T014 [US2] Add request-sequenced Activity loading and debounced live updates in `desktop/frontend/src/operations/store.ts`
- [x] T015 [US2] Replace the Activity placeholder with search, type, severity, outcome, typed rows, and stable detail in `desktop/frontend/src/operations/ActivityPage.tsx` and `desktop/frontend/src/App.tsx`
- [x] T016 [US2] Add responsive Activity filters, table, state, and diagnostic styling in `desktop/frontend/src/styles.css`
- [x] T017 [US2] Add React tests for typed records, filters, run diagnostics, exact log path, inaccessible references, stable selection, focus, complete-snapshot preservation, and 100 rows

## Phase 5: User Story 3 - Control alerts and current view safely (Priority: P2)

**Goal**: Acknowledge alerts and clear only the visible local view with explicit non-destructive semantics.

**Independent Test**: Acknowledge one alert and clear filtered results, proving mutation scope, local cutoff, duplicate suppression, and later-event visibility.

- [x] T018 [US3] Expose acknowledgement methods through `desktop/app.go` and `desktop/frontend/src/operations/bridge.ts`
- [x] T019 [US3] Add accessible acknowledgement and Clear View controls with visible non-destructive explanation in `desktop/frontend/src/operations/ActivityPage.tsx`
- [x] T020 [US3] Add Go and React coverage for one-alert scope, visible-only acknowledgement, duplicate suppression, persisted-data wording, and post-cutoff activity

## Phase 6: Integration and verification

- [x] T021 Add Chromium keyboard, accessibility, stable-selection, and 100-row interaction coverage in `desktop/frontend/e2e/operations.spec.ts`
- [x] T022 Update Unreleased change notes and architecture rationale in `CHANGELOG.md`
- [x] T023 Run focused Go, race, frontend unit, frontend build, Chromium, and native Wails verification from `quickstart.md`
- [x] T024 Run spec-kit analysis and remediate every finding until the report is clean
- [x] T025 Run canonical `sh scripts/verify.sh all` in the foreground and record objective evidence in `specs/064-schedule-activity/verification.md`
- [x] T026 Mark required tasks complete, transition the specification to Implemented, and synchronize `specs/README.md`

## Dependencies and Execution Order

- Phase 1 establishes contracts used by all workspaces.
- Phase 2 blocks both frontend stories and safe alert actions.
- User Story 1 and User Story 2 share contracts and state mechanics but remain independently testable routes.
- User Story 3 depends on the Activity workspace from User Story 2.
- Phase 6 follows all required stories.

## Implementation Strategy

Build and test the complete backend snapshots first, then Schedule, then Activity inspection, then acknowledgement and Clear View. Keep daemon behavior unchanged, run focused verification after each story, and use the canonical repository gate as the final definition of green.
