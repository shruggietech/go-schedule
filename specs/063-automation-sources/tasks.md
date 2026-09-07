# Tasks: Connected Automation Sources

**Input**: Design documents from `/specs/063-automation-sources/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/desktop-automation-bridge.md`

**Tests**: Required by FR-017 and SC-006. Tests are authored before or with the implementation they constrain and must pass before a task is complete.

## Phase 1: Setup and contracts

- [x] T001 Confirm issue #154 traceability and complete specification, clarification, plan, checklist, contract, and task artifacts in `specs/063-automation-sources/`
- [x] T002 [P] Add transport-neutral automation workspace, draft, operation, and ephemeral secret models in `desktop/automation/model.go`
- [x] T003 [P] Define and adapt the existing daemon client behind a bounded backend interface in `desktop/automation/local.go`

## Phase 2: Foundational desktop service

- [x] T004 Add complete secret-free workspace composition and safe error mapping in `desktop/automation/service.go`
- [x] T005 Add shared validation, stale-write detection, bounded calls, and mutation-refresh helpers in `desktop/automation/service.go`
- [x] T006 Add service tests for complete-snapshot failure, readiness and health mapping, stale writes, safe errors, and JSON secret exclusion in `desktop/automation/service_test.go`

## Phase 3: User Story 1 - Understand every automation source (Priority: P1)

**Goal**: Present every source type in one searchable, accessible source-to-target workspace with honest readiness and health.

**Independent Test**: Load mixed fixtures, filter at least 100 sources, and verify source, target, state, reason, focus retention, and secret absence.

- [x] T007 [US1] Bind the automation service through Wails in `desktop/main.go` and `desktop/app.go`
- [x] T008 [P] [US1] Add frontend automation contracts and native bridge in `desktop/frontend/src/automation/model.ts` and `desktop/frontend/src/automation/bridge.ts`
- [x] T009 [US1] Add request-sequenced workspace loading and last-complete-snapshot preservation in `desktop/frontend/src/automation/store.ts`
- [x] T010 [US1] Add one Automation Sources route and navigation entry in `desktop/frontend/src/connection/model.ts`, `desktop/frontend/src/components/Shell.tsx`, and `desktop/frontend/src/App.tsx`
- [x] T011 [US1] Render shared search, attention filters, and four semantic source sections in `desktop/frontend/src/automation/AutomationPage.tsx`
- [x] T012 [US1] Add responsive source cards, path overflow, status, editor, and secret-dialog styling in `desktop/frontend/src/styles.css`
- [x] T013 [US1] Add React tests for mixed readiness, missing relationships, read-only offline context, request ordering, keyboard focus, and 100-source filtering in `desktop/frontend/src/automation/AutomationPage.test.tsx` and `desktop/frontend/src/automation/store.test.ts`

## Phase 4: User Story 2 - Manage completion chains and watchers (Priority: P2)

**Goal**: Complete chain and watcher creation, editing, retargeting, toggling where supported, and deletion with accurate watcher health.

**Independent Test**: Exercise a full chain and watcher lifecycle, stale edit, long path, degraded health, and missing target.

- [x] T014 [US2] Implement completion-chain create, update, stale detection, and delete operations in `desktop/automation/service.go`
- [x] T015 [US2] Implement filesystem-watcher create, update, toggle, stale detection, and delete operations in `desktop/automation/service.go`
- [x] T016 [US2] Expose chain and watcher methods through `desktop/app.go` and `desktop/frontend/src/automation/bridge.ts`
- [x] T017 [US2] Add accessible type-specific chain and watcher editors and confirmations in `desktop/frontend/src/automation/AutomationPage.tsx`
- [x] T018 [US2] Add Go and React coverage for validation, stale overwrite, health accuracy, duplicate suppression, and failed-operation snapshot preservation

## Phase 5: User Story 3 - Manage secret-bearing trigger sources (Priority: P3)

**Goal**: Complete external-trigger and Trigger Set lifecycles while strictly containing raw secrets.

**Independent Test**: Create, edit, toggle, reveal, copy, rotate, fire, retarget, and delete trigger sources while proving raw keys only exist in the explicit secret dialog.

- [x] T019 [US3] Implement external-trigger create, update, toggle, stale detection, reveal, rotate, internal-key fire, and delete operations in `desktop/automation/service.go`
- [x] T020 [US3] Implement Trigger Set create, atomic retarget, toggle, reveal, rotate, and delete operations in `desktop/automation/service.go`
- [x] T021 [US3] Expose trigger and Trigger Set methods through `desktop/app.go` and `desktop/frontend/src/automation/bridge.ts`
- [x] T022 [US3] Add accessible trigger editors, action confirmations, copy feedback, and ephemeral secret dialogs in `desktop/frontend/src/automation/AutomationPage.tsx`
- [x] T023 [US3] Add Go and React tests for trigger firing, member ordering, duplicate suppression, dialog clearing, and absence of keys from ordinary state

## Phase 6: Integration and verification

- [x] T024 Update Unreleased change notes and architecture decision rationale in `CHANGELOG.md`
- [x] T025 Run focused Go, race, frontend unit, frontend build, and native Wails build verification from `quickstart.md`
- [x] T026 Run spec-kit analysis and remediate every finding until the report is clean
- [x] T027 Run canonical `sh scripts/verify.sh all` in the foreground and record objective evidence in `specs/063-automation-sources/verification.md`
- [x] T028 Mark required tasks complete, transition the specification to Implemented, and synchronize `specs/README.md`

## Dependencies and Execution Order

- Phase 1 establishes contracts used by every story.
- Phase 2 blocks all UI and lifecycle work.
- User Story 1 establishes the shared workspace required by User Stories 2 and 3.
- User Stories 2 and 3 use distinct lifecycle methods but share UI files, so they proceed chronologically to avoid conflicting edits.
- Phase 6 follows all required user stories.

## Implementation Strategy

Deliver the shared secret-free workspace first, then completion-chain and watcher management, then secret-bearing trigger lifecycles. Keep each phase independently testable, run focused verification after each story, and use the canonical repository gate as the only final definition of green.
