# Tasks: Desktop Settings, Information, and Recovery

**Input**: Design documents from `/specs/065-options-info/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/desktop-settings-bridge.md`

**Tests**: Required by FR-020 and SC-006. Tests are authored before or with the implementation they constrain and must pass before a task is complete.

## Phase 1: Setup and contracts

- [x] T001 Confirm issues #156 and #195 traceability and complete specification, clarification, plan, checklist, contract, and task artifacts in `specs/065-options-info/`
- [x] T002 [P] Add transport-neutral preference, storage, product, workspace, and result models in `desktop/settings/model.go`
- [x] T003 [P] Define injectable filesystem, runtime-info, clipboard, browser, clock, and platform boundaries in `desktop/settings/local.go`

## Phase 2: Foundational preference and settings service

- [x] T004 Add versioned preference loading, validation, one-time Fyne appearance migration, and deterministic transition state in `desktop/settings/preferences.go`
- [x] T005 Add same-directory atomic preference persistence, restrictive permissions, save, and restore behavior in `desktop/settings/preferences.go`
- [x] T006 Add authoritative local and daemon storage inventory plus product metadata in `desktop/settings/service.go`
- [x] T007 Add identifier-authorized clipboard and HTTPS product-link actions in `desktop/settings/service.go`
- [x] T008 Add Go tests for migration matrices, established preference authority, write failure, concurrent mutation, storage truth, offline behavior, and native-action rejection in `desktop/settings/service_test.go`

## Phase 3: User Story 1 - Upgrade with predictable preferences (Priority: P1)

**Goal**: Preserve valid appearance intent, retire obsolete controls explicitly, and provide durable appearance and restore controls.

**Independent Test**: Exercise legacy and current fixtures, change appearance, restart, and restore defaults without affecting scheduler data.

- [x] T009 [US1] Bind settings initialization, appearance save, and restore through Wails in `desktop/main.go` and `desktop/app.go`
- [x] T010 [P] [US1] Add frontend settings contracts and native bridge in `desktop/frontend/src/settings/model.ts` and `desktop/frontend/src/settings/bridge.ts`
- [x] T011 [US1] Initialize global appearance from the authoritative Settings workspace and preserve the prior active value on failed writes in `desktop/frontend/src/settings/store.ts` and `desktop/frontend/src/App.tsx`
- [x] T012 [US1] Add accessible Appearance and Preference transition sections with save, pending, error, and restore states in `desktop/frontend/src/settings/SettingsPage.tsx`
- [x] T013 [US1] Add Go and React tests for initialization, migration disclosure, retired controls, save, restore, duplicate suppression, and failed-write preservation

## Phase 4: User Story 2 - Understand storage and product information (Priority: P1)

**Goal**: Deliver accurate storage ownership and removal guidance plus trustworthy product metadata and native actions.

**Independent Test**: Render every storage source state, copy eligible exact paths, open fixed product links, and reject unknown identifiers.

- [x] T014 [US2] Add Storage and About sections with quiet unavailable states and connection recovery navigation in `desktop/frontend/src/settings/SettingsPage.tsx`
- [x] T015 [US2] Expose identifier-based copy and product-link actions through `desktop/app.go` and `desktop/frontend/src/settings/bridge.ts`
- [x] T016 [US2] Add responsive Settings layout, storage record, disclosure, copy status, and About styling in `desktop/frontend/src/styles.css`
- [x] T017 [US2] Add React tests for storage ownership, removal semantics, exact paths, unavailable daemon data, copy announcements, fixed links, and unknown action failures

## Phase 5: User Story 3 - Recover the local connection without interruption (Priority: P2)

**Goal**: Turn existing connection state into a dedicated inline diagnosis and retry workspace.

**Independent Test**: Present every diagnosis, retry safely, and verify route, focus, and local Settings remain stable.

- [x] T018 [US3] Add an accessible Connections page using the existing connection manager snapshot and retry action in `desktop/frontend/src/settings/ConnectionsPage.tsx`
- [x] T019 [US3] Replace the Connections placeholder and keep connection updates inline without automatically opening dialogs in `desktop/frontend/src/App.tsx` and `desktop/frontend/src/Shell.tsx`
- [x] T020 [US3] Add React tests for diagnosis-specific guidance, safe context, retry sequencing, duplicate suppression, route preservation, focus preservation, and no automatic modal

## Phase 6: Integration and verification

- [x] T021 Add Chromium keyboard, accessibility, native-action, offline, and 80 through 200 percent zoom coverage in `desktop/frontend/e2e/settings.spec.ts`
- [x] T022 Update preference transition and storage documentation in `docs/options.md` plus Unreleased change and architecture rationale in `CHANGELOG.md`
- [x] T023 Run focused Go, race, frontend unit, frontend build, Chromium, and native Wails verification from `quickstart.md`
- [x] T024 Run spec-kit analysis and remediate every finding until the report is clean
- [x] T025 Run canonical `sh scripts/verify.sh all` in the foreground and record objective evidence in `specs/065-options-info/verification.md`
- [x] T026 Mark required tasks complete, transition the specification to Implemented, and synchronize `specs/README.md`

## Dependencies and Execution Order

- Phase 1 establishes the shared contracts and injectable boundaries.
- Phase 2 blocks all frontend stories and native actions.
- User Story 1 establishes authoritative appearance state consumed by the shell.
- User Story 2 builds on the workspace and action contracts but is independently testable with fixtures.
- User Story 3 reuses the existing connection manager and can proceed after the Settings route integration.
- Phase 6 follows all required stories.

## Implementation Strategy

Build and test preference persistence and settings composition first, then integrate startup appearance, storage and About, and finally the dedicated Connections route. Preserve exact source authority at each boundary, run focused verification after every story, and use the canonical repository gate as the final definition of green.
