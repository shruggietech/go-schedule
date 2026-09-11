# Tasks: Responsive Task and Administration Workflows

**Input**: Design documents from `specs/084-responsive-workflows/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/responsive-workflows.md`, `quickstart.md`

**Tests**: Required by the constitution and issues #231 and #233. Behavioral work begins with failing focused regressions.

## Phase 1: Specification and Setup

- [x] T001 Register S084 and its Draft delivery state in `.specify/feature.json` and `specs/README.md`
- [x] T002 [P] Complete requirement-quality validation in `specs/084-responsive-workflows/checklists/requirements.md` and `specs/084-responsive-workflows/checklists/ux.md`
- [x] T003 [P] Complete research, transient data model, interaction contract, and validation guide in `specs/084-responsive-workflows/`
- [x] T004 Update the managed plan reference in `CLAUDE.md` through the agent-context extension

## Phase 2: Shared Responsive Foundations

**Purpose**: Establish the reusable semantics that every S084 user story consumes.

- [x] T005 Add failing component tests for FormGrid, DescriptionList, StatusLabel, CardSection, PathDisplay, dialog close labels, scroll structure, and unmount focus restoration in `desktop/frontend/src/components/components.test.tsx`
- [x] T006 Implement the shared administration primitives and extended Dialog contract in `desktop/frontend/src/components/index.tsx`
- [x] T007 Add responsive form, definition-list, path, disclosure, dialog-body, and action-footer styling in `desktop/frontend/src/styles.css`

**Checkpoint**: Shared presentation and modal contracts are independently covered and ready for domain workflows.

## Phase 3: User Story 1 - Focused Task Authoring (Priority: P1)

**Goal**: Create and edit tasks in one bounded, accessible modal without changing the Tasks page geometry.

**Independent Test**: Open Create task and Edit, exercise safe insertion, preview, validation, Advanced settings, Escape, Cancel, Save, and deletion confirmation while checking focus and geometry.

- [x] T008 [US1] Add failing modal create, edit, cancellation, save, geometry, and focus-restoration regressions in `desktop/frontend/src/tasks/TasksPage.test.tsx`
- [x] T009 [P] [US1] Add failing task editor tests for initial focus, concise safe example, Advanced settings default, preview, validation focus, and pending actions in `desktop/frontend/src/tasks/TaskEditor.test.tsx`
- [x] T010 [P] [US1] Extend deletion-dialog tests for safe action naming, spacing hook, pending protection, Escape, and focus restoration in `desktop/frontend/src/tasks/TaskActions.test.tsx`
- [x] T011 [US1] Refactor TaskEditor into the shared dialog body and actions contract in `desktop/frontend/src/tasks/TaskEditor.tsx`
- [x] T012 [US1] Track exact create and edit invokers and open the focused workflow from `desktop/frontend/src/tasks/TasksPage.tsx`
- [x] T013 [US1] Finalize task dialog and confirmation responsive styling in `desktop/frontend/src/styles.css`
- [x] T014 [US1] Run the focused task and shared component tests and resolve every regression

**Checkpoint**: #231 is complete and independently verifiable.

## Phase 4: User Story 2 - Agent Access Composition (Priority: P1)

**Goal**: Present safe agent-access status before a responsive, deliberately disclosed optional-network form.

**Independent Test**: Inspect off and active states, expand localhost configuration, traverse fields in order, and verify no secret or collision at constrained sizes.

- [x] T015 [US2] Add failing hierarchy, disclosure, field-label, definition-list, and active-lifecycle tests in `desktop/frontend/src/agentaccess/AgentAccessPage.test.tsx`
- [x] T016 [US2] Recompose Agent Access through shared sections, definitions, fields, status labels, and disclosure in `desktop/frontend/src/agentaccess/AgentAccessPage.tsx`
- [x] T017 [US2] Run the focused Agent Access tests and resolve every regression

**Checkpoint**: Agent Access satisfies its portion of #233 without changing access contracts.

## Phase 5: User Story 3 - Connection Diagnosis and Pairing (Priority: P1)

**Goal**: Present diagnosis and saved profiles before a responsive pairing workflow that is collapsed until needed.

**Independent Test**: Exercise connection states, long profile values, ordinary pairing disclosure, automatic repair disclosure, reading order, and keyboard focus.

- [x] T018 [P] [US3] Add failing diagnosis hierarchy, long-value, profile-action, and pairing-disclosure tests in `desktop/frontend/src/settings/ConnectionsPage.test.tsx`
- [x] T019 [P] [US3] Add failing responsive field, semantic form, pending, and repair-mode tests in `desktop/frontend/src/remotepairing/PairingForm.test.tsx`
- [x] T020 [US3] Recompose diagnosis and saved targets through shared administration primitives in `desktop/frontend/src/settings/ConnectionsPage.tsx`
- [x] T021 [US3] Recompose PairingForm with shared Field and FormGrid semantics in `desktop/frontend/src/remotepairing/PairingForm.tsx`
- [x] T022 [US3] Run focused Connections and Pairing tests and resolve every regression

**Checkpoint**: Connections and Pairing satisfy their portion of #233 without changing network behavior.

## Phase 6: User Story 4 - Settings Paths and Scoped Copy Feedback (Priority: P1)

**Goal**: Make storage paths scannable and keep each copy action visually and temporally independent.

**Independent Test**: Present long paths, overlap copy operations, resolve them out of order, and verify only the matching action changes state.

- [x] T023 [P] [US4] Add failing complete-path, monospace, record-label, and isolated button-state tests in `desktop/frontend/src/settings/SettingsPage.test.tsx`
- [x] T024 [P] [US4] Add failing keyed-operation, overlap, out-of-order, success, and failure tests in `desktop/frontend/src/settings/store.test.ts`
- [x] T025 [US4] Implement record-keyed pending and copy-result state in `desktop/frontend/src/settings/store.ts` and integrate it in `desktop/frontend/src/App.tsx`
- [x] T026 [US4] Recompose Settings sections, storage definitions, paths, badges, and copy actions in `desktop/frontend/src/settings/SettingsPage.tsx`
- [x] T027 [US4] Run focused Settings component, store, and App tests and resolve every regression

**Checkpoint**: #233 is complete and independently verifiable across all three administration destinations.

## Phase 7: Cross-Workflow Verification and Delivery Evidence

- [x] T028 Add S084 browser workflows for modal geometry, disclosures, administration reflow, 200 percent zoom, themes, keyboard focus, long values, copy isolation, and axe coverage in `desktop/frontend/e2e/shell.spec.ts`
- [x] T029 Run the focused frontend suite, production bundle, and Playwright checks from `specs/084-responsive-workflows/quickstart.md`
- [x] T030 Update the Unreleased changelog with the two completed issue outcomes and the decision to extend existing primitives in `CHANGELOG.md`
- [x] T031 Record specification analysis, test-first evidence, focused results, native build evidence, canonical results, and release boundary in `specs/084-responsive-workflows/verification.md`
- [x] T032 Move the specification and inventory to Implemented with objective evidence in `specs/084-responsive-workflows/spec.md` and `specs/README.md`
- [x] T033 Run `sh scripts/verify.sh all` in the foreground and resolve every S084 failure

## Dependencies and Execution Order

- Phase 1 defines the reviewed scope and must finish before implementation.
- Phase 2 blocks every user story because each workflow consumes the shared primitives.
- User Story 1 can complete after Phase 2 and independently closes #231.
- User Stories 2 and 3 can begin after Phase 2, but both must complete before #233 can close.
- User Story 4 depends on the shared primitives and integrates keyed settings state through App.
- Phase 7 depends on all four user stories and is the only point that may mark the specification Implemented.

## Parallel Opportunities

- T002 and T003 affect separate specification artifacts.
- T009 and T010 cover separate task components after T008 defines the page-level contract.
- T018 and T019 cover separate connection and pairing components.
- T023 and T024 cover separate presentation and state-management layers.
- User Stories 2 and 3 touch separate domain modules after the shared foundation is stable.

## Implementation Strategy

1. Establish shared semantics and failing regressions before changing domain markup.
2. Complete the focused task workflow first so #231 remains an independently reviewable increment inside the bundled slice.
3. Apply the same administration primitives to Agent Access, Connections, Pairing, and Settings without changing their bridges.
4. Prove record identity and responsive behavior at component level before browser geometry checks.
5. Run the complete canonical gate once the focused suite and production bundle are green.
