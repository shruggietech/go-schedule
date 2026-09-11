# Tasks: Desktop Visual and Shell Foundations

**Input**: Design documents from `specs/083-desktop-foundations/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/desktop-foundations.md`, `quickstart.md`

**Tests**: Required by the constitution and issues #229 and #230. Behavioral changes begin with failing focused tests.

## Phase 1: Specification and Contracts

- [x] T001 Register S083 and its Draft delivery state in `.specify/feature.json` and `specs/README.md`
- [x] T002 [P] Validate specification quality in `specs/083-desktop-foundations/checklists/requirements.md` and `specs/083-desktop-foundations/checklists/ux.md`
- [x] T003 [P] Finalize research, transient view models, interaction contracts, and validation guidance in `specs/083-desktop-foundations/research.md`, `specs/083-desktop-foundations/data-model.md`, `specs/083-desktop-foundations/contracts/desktop-foundations.md`, and `specs/083-desktop-foundations/quickstart.md`

## Phase 2: Shared Component Foundation

- [x] T004 Add failing shared button, dismissible Notice, grouped Dialog action, and transient ToastRegion tests in `desktop/frontend/src/components/components.test.tsx`
- [x] T005 Implement shared semantic control variants and feedback behavior in `desktop/frontend/src/components/index.tsx`
- [x] T006 Define and consume compact semantic tokens, interaction states, card spacing, dialog actions, and overlay feedback styling in `desktop/frontend/src/styles.css`

## Phase 3: User Story 1 - Operate a Coherent Desktop (Priority: P1)

**Goal**: Make shared actions compact, visibly interactive, semantically distinct, and palette-consistent.

**Independent Test**: Exercise every shared variant and representative task actions with pointer and keyboard in each appearance.

- [x] T007 [US1] Add failing affirmative task-action and appearance-resolution tests in `desktop/frontend/src/tasks/TaskActions.test.tsx` and `desktop/frontend/src/App.test.tsx`
- [x] T008 [US1] Apply affirmative and destructive task actions plus pending semantics in `desktop/frontend/src/tasks/TaskActions.tsx`
- [x] T009 [US1] Resolve follow-system palette changes and expose consistent imagery state in `desktop/frontend/src/components/Shell.tsx`
- [x] T010 [US1] Extend automated contrast and state coverage in `desktop/frontend/src/accessibility.test.tsx` and `desktop/frontend/e2e/shell.spec.ts`

## Phase 4: User Story 2 - Keep the Shell Available (Priority: P1)

**Goal**: Keep navigation, target context, Appearance, and Exit reachable while only active page content scrolls.

**Independent Test**: Scroll every route at standard, compact, minimum, and zoomed viewports without body overflow or lost shell controls.

- [x] T011 [US2] Add failing viewport ownership, scroll-container, fixed-control, 800 by 600, and 200 percent zoom checks in `desktop/frontend/e2e/shell.spec.ts`
- [x] T012 [US2] Implement persistent shell regions and one active-page scroll owner in `desktop/frontend/src/components/Shell.tsx` and `desktop/frontend/src/styles.css`

## Phase 5: User Story 3 - Receive Non-Disruptive Feedback (Priority: P1)

**Goal**: Replace permanent concatenated status strings with bounded transient feedback and dismissible persistent errors.

**Independent Test**: Trigger replacement, timeout, hover pause, focus pause, explicit dismissal, and persistent error behavior without geometry changes.

- [x] T013 [US3] Add failing shell integration tests for message replacement and non-disruptive feedback in `desktop/frontend/src/App.test.tsx`
- [x] T014 [US3] Classify settings and connection outcomes instead of concatenating them in `desktop/frontend/src/App.tsx`
- [x] T015 [US3] Wire the transient toast and persistent message contracts through `desktop/frontend/src/components/Shell.tsx` and affected callers in `desktop/frontend/src/App.tsx`

## Phase 6: Verification and Publication Readiness

- [x] T016 [P] Record the shared-foundation architecture decision and user-visible correction in `CHANGELOG.md`
- [x] T017 Run focused frontend component, accessibility, bundle, and Playwright checks from `desktop/frontend/`
- [x] T018 Build the Wails desktop and record focused Windows observation against one exact source tree in `specs/083-desktop-foundations/verification.md`
- [x] T019 Run `sh scripts/verify.sh all` and record all eight canonical gate results in `specs/083-desktop-foundations/verification.md`
- [x] T020 Mark every task complete and transition S083 to Implemented with objective delivery evidence in `specs/083-desktop-foundations/spec.md`, `specs/083-desktop-foundations/tasks.md`, and `specs/README.md`

## Dependencies and Execution Order

- Phase 1 is complete and enables all implementation work.
- Phase 2 is foundational and blocks every user story.
- User Story 1 establishes visual primitives before User Stories 2 and 3 consume the shared shell and feedback layers.
- User Stories 2 and 3 can proceed independently after Phase 2, but changes to `Shell.tsx` and `styles.css` remain sequential.
- Verification and publication readiness depend on all three user stories.

## Parallel Opportunities

- T002 and T003 cover independent requirements artifacts.
- T016 can proceed independently of focused test execution once source behavior is final.
- Component tests and browser test authoring touch different files when explicitly marked, but implementation changes to shared source files remain sequential.

## Implementation Strategy

1. Establish the component foundation with failing tests and the smallest shared API changes.
2. Complete semantic action and theme resolution behavior.
3. Establish one viewport and scroll contract.
4. Replace concatenated footer feedback with one transient message and persistent actionable notices.
5. Validate the entire shared layer before allowing page-specific work in #231 through #233 to build on it.
