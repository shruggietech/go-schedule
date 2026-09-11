# Tasks: Task-First Notifications

**Input**: Design documents from `specs/085-task-first-notifications/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/task-first-notifications.md`, `quickstart.md`

**Tests**: Required by the constitution and issue #232. Behavioral work begins with failing focused regressions.

## Phase 1: Specification and Setup

- [x] T001 Register S085 and its Draft delivery state in `.specify/feature.json` and `specs/README.md`
- [x] T002 [P] Complete requirements and UX quality validation in `specs/085-task-first-notifications/checklists/requirements.md` and `specs/085-task-first-notifications/checklists/ux.md`
- [x] T003 [P] Complete research, view-state model, interaction contract, and validation guide in `specs/085-task-first-notifications/`
- [x] T004 Update the managed plan reference in `CLAUDE.md` through the agent-context extension

## Phase 2: Foundational Coverage Projection

**Purpose**: Make configured task and group coverage available to the initial overview without exposing secrets or changing daemon contracts.

- [x] T005 Add failing configured task, inherited task, direct group, empty, partial-failure, ordering, and secret-exclusion tests in `desktop/notifications/service_test.go`
- [x] T006 Extend secret-free coverage models in `desktop/notifications/model.go` and `desktop/frontend/src/notifications/model.ts`
- [x] T007 Implement bounded coverage projection through existing backend methods in `desktop/notifications/service.go`
- [x] T008 Run focused notification service race tests and resolve every regression

**Checkpoint**: One workspace load can truthfully summarize configured notification coverage.

## Phase 3: User Story 1 - Immediate Status and Recent Results (Priority: P1)

**Goal**: Answer whether notifications are active, what can notify, what recently happened, and what action is useful before setup forms.

**Independent Test**: Render empty, disabled, healthy, in-progress, retrying, failed, and incomplete-coverage workspaces and inspect the first surface at 800 by 600.

- [x] T009 [US1] Add failing overview hierarchy, count, guidance, coverage, bounded-result, and empty-state tests in `desktop/frontend/src/notifications/NotificationsPage.test.tsx`
- [x] T010 [US1] Implement deterministic overview and delivery-guidance helpers in `desktop/frontend/src/notifications/NotificationsPage.tsx`
- [x] T011 [US1] Implement the initial status, configured coverage, and bounded recent-results composition in `desktop/frontend/src/notifications/NotificationsPage.tsx`
- [x] T012 [US1] Add overview, coverage, result-list, and state styling in `desktop/frontend/src/styles.css`
- [x] T013 [US1] Run focused notification component tests and resolve every regression

**Checkpoint**: The first viewport independently satisfies the primary user-facing outcome.

## Phase 4: User Story 2 - Deliberate Setup (Priority: P1)

**Goal**: Keep destinations and assignment rules collapsed until requested while preserving every existing operation and secret boundary.

**Independent Test**: Expand each setup section by keyboard, exercise channel and policy operations, and confirm contextual explanations and write-only secrets.

- [x] T014 [US2] Add failing initial-collapse, disclosure, contextual-help, channel-operation, policy, and focus-order tests in `desktop/frontend/src/notifications/NotificationsPage.test.tsx`
- [x] T015 [US2] Move destination administration into a task-named Disclosure with contextual help in `desktop/frontend/src/notifications/NotificationsPage.tsx`
- [x] T016 [US2] Move assignment administration into a task-named Disclosure with policy and inheritance guidance in `desktop/frontend/src/notifications/NotificationsPage.tsx`
- [x] T017 [US2] Add advanced notification section, help, form, and responsive styling in `desktop/frontend/src/styles.css`
- [x] T018 [US2] Run preserved channel, policy, dialog, and store tests and resolve every regression

**Checkpoint**: Setup remains complete, secret-safe, and available without dominating Notifications.

## Phase 5: User Story 3 - On-Demand Diagnostics (Priority: P1)

**Goal**: Preserve full redacted delivery investigation behind one explicit diagnostic disclosure.

**Independent Test**: Expand diagnostics, filter every state and destination, select test and task outcomes, and inspect stale and redacted detail.

- [x] T019 [US3] Add failing diagnostic-collapse, filter, selection, missing-record, redaction, and long-history tests in `desktop/frontend/src/notifications/NotificationsPage.test.tsx`
- [x] T020 [US3] Move full history filters, table, and detail into Delivery diagnostics in `desktop/frontend/src/notifications/NotificationsPage.tsx`
- [x] T021 [US3] Bound diagnostics and long values responsively in `desktop/frontend/src/styles.css`
- [x] T022 [US3] Run focused diagnostic and polling tests and resolve every regression

**Checkpoint**: Detailed evidence remains fully available without burdening the initial page.

## Phase 6: Cross-Workflow Verification and Delivery Evidence

- [x] T023 Add S085 browser scenarios for first-time comprehension, all overview states, disclosures, preserved administration, 800 by 600, 200 percent zoom, keyboard use, long values, and axe coverage in `desktop/frontend/e2e/notifications.spec.ts`
- [x] T024 Run the focused Go, frontend, bundle, and Playwright checks from `specs/085-task-first-notifications/quickstart.md`
- [x] T025 Update the Unreleased changelog with the #232 outcome and architecture decision in `CHANGELOG.md`
- [x] T026 Record analysis, test-first evidence, focused results, native build evidence, canonical results, and release boundary in `specs/085-task-first-notifications/verification.md`
- [x] T027 Move the specification and inventory to Implemented with objective evidence in `specs/085-task-first-notifications/spec.md` and `specs/README.md`
- [x] T028 Run `sh scripts/verify.sh all` in the foreground and resolve every S085 failure

## Dependencies and Execution Order

- Phase 1 defines and validates scope before implementation.
- Phase 2 blocks the overview because configured coverage must be authoritative and secret-free.
- User Story 1 establishes the primary initial surface before advanced sections move.
- User Story 2 preserves setup under progressive disclosure after the overview is stable.
- User Story 3 preserves full evidence after recent results have a bounded initial representation.
- Phase 6 depends on every story and is the only point that may mark the specification Implemented.

## Parallel Opportunities

- T002 and T003 affect separate specification artifacts.
- Coverage model updates in Go and TypeScript can proceed together after T005 defines the contract.
- Browser fixture and documentation preparation affect separate files after component behavior stabilizes.

## Implementation Strategy

1. Prove configured coverage and its security boundary in Go before exposing it to the browser.
2. Establish the initial overview and state guidance before moving existing administration.
3. Move channel and policy markup with behavior-preservation regressions, then move detailed evidence.
4. Exercise all states and expanded sections at the supported viewport and zoom boundary.
5. Run the complete canonical gate after focused and native checks are green.
