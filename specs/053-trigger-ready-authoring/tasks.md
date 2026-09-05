# Tasks: Trigger-Ready Task Authoring

**Input**: Design documents from `specs/053-trigger-ready-authoring/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Required by FR-028 and constitution principle II. Behavioral tests must fail before implementation.

## Phase 1: Setup and Baseline

**Purpose**: Confirm existing task, chain, group, navigation, and migration contracts before changing them.

- [x] T001 Record the clean baseline and focused test results in specs/053-trigger-ready-authoring/verification.md
- [x] T002 Verify existing ignore files remain sufficient for the Go and Fyne project and record the result in specs/053-trigger-ready-authoring/verification.md

---

## Phase 2: Foundational Readiness and Persistence

**Purpose**: Establish one derived vocabulary and truthful unscheduled storage used by every story.

- [x] T003 [P] Add failing command, activation-source, status-precedence, and `unnamed` helper tests in internal/task/readiness_test.go
- [x] T004 [P] Add a failing v10-to-v11 preservation test for nullable task schedules in internal/store/migration_v11_test.go
- [x] T005 [P] Add failing store tests for unscheduled round trips, guarded enable, update auto-disable, and completion-source transitions in internal/store/store_test.go and internal/store/chains_test.go
- [x] T006 Implement readiness vocabulary and display-name fallback in internal/task/readiness.go
- [x] T007 Implement the forward-only nullable-schedule migration and nullable CRUD scan/write behavior in internal/store/store.go and internal/store/crud.go
- [x] T008 Implement transactional enable, task-update readiness enforcement, and automatic-source queries in internal/store/crud.go
- [x] T009 Implement atomic direct-chain and source-task deletion readiness transitions in internal/store/chains.go and internal/store/crud.go

**Checkpoint**: Storage represents omitted schedules truthfully and every readiness mutation has one atomic authority.

---

## Phase 3: User Story 1 - Preserve Incomplete Task Work (Priority: P1)

**Goal**: Save, reload, list, and edit tasks with omitted name, command, or schedule while preserving invalid-input rejection.

**Independent Test**: Round-trip each omission and a fully blank draft through API, restart, CLI, and editor without placeholder persistence or partial mutation.

- [x] T010 [P] [US1] Add failing API create, detail, update-clear, malformed-input, and compatibility tests in internal/api/server/tasks_test.go and internal/api/server/update_test.go
- [x] T011 [P] [US1] Add failing CLI draft-create, explicit-clear, `unnamed`, and absent-schedule output tests in internal/cli/task_test.go
- [x] T012 [P] [US1] Add failing editor tests for blank draft save, blank prefill, exact fallback display, and removal edits in gui/editor_test.go and gui/editor_prefill_test.go
- [x] T013 [US1] Extend task API contracts for omitted create fields, explicit clear intent, optional schedule detail, and derived readiness in internal/api/server/tasks.go and internal/api/server/update.go
- [x] T014 [US1] Make calendar and scheduler reads treat absent schedules as expected in internal/api/server/calendar.go and internal/engine/engine.go
- [x] T015 [US1] Extend CLI create, edit, and output behavior for drafts and explicit clears in internal/cli/task.go
- [x] T016 [US1] Extend task editor validation, submission, and prefill for omitted draft values and explicit removal in gui/editor.go

**Checkpoint**: Drafts round-trip safely and configured callers retain their existing behavior.

---

## Phase 4: User Story 2 - Understand and Control Readiness (Priority: P1)

**Goal**: Enforce and display consistent not-runnable, manual-only, ready, disabled, blocked, and terminal states.

**Independent Test**: Exercise all readiness combinations and source transitions, then compare API, CLI, Tasks, Groups, enable, and Run now behavior.

- [x] T017 [P] [US2] Add failing API and engine tests for guarded enable, guarded Run now, manual-only execution, and no missing-schedule log noise in internal/api/server/tasks_test.go and internal/engine/engine_extra_test.go
- [x] T018 [P] [US2] Add failing Tasks and Groups presentation tests for readiness precedence, completion sources, and stable nameless identity in gui/tasks_test.go and gui/groups_test.go
- [x] T019 [US2] Enforce command readiness before Run now and map readiness errors to actionable API responses in internal/api/server/tasks.go and internal/engine/runnow.go
- [x] T020 [US2] Present readiness and effective eligibility consistently in CLI task list/show output in internal/cli/task.go
- [x] T021 [US2] Present readiness, fallback names, and stable identity consistently in Tasks and Groups in gui/tasks.go and gui/groups.go

**Checkpoint**: Every interface agrees on why a task can or cannot run automatically or manually.

---

## Phase 5: User Story 3 - Navigate by Product Concept (Priority: P2)

**Goal**: Render two semantic destination sections above the separately anchored Exit control.

**Independent Test**: Exercise section membership, destination selection, badge updates, theme interaction, keyboard order, and constrained layouts headlessly.

- [x] T022 [P] [US3] Add failing semantic-section, future-insertion, selection, badge, size, and theme tests in gui/navigation_test.go and gui/app_test.go
- [x] T023 [US3] Add destination section metadata and grouped navigation layout in gui/navigation.go and gui/app.go

**Checkpoint**: Navigation communicates concept groups without exposing unfinished Triggers behavior.

---

## Phase 6: User Story 4 - Create a Group from the Keyboard (Priority: P2)

**Goal**: Make Enter and Create share one validated, duplicate-safe new-group submission path.

**Independent Test**: Submit valid, blank, repeated, composition-resolved, and click-after-key cases through the headless group dialog controller.

- [x] T024 [P] [US4] Add failing group submission tests for Enter, button parity, blank retention, and duplicate suppression in gui/groups_test.go
- [x] T025 [US4] Implement the validated shared group submission controller and dialog in gui/groups.go

**Checkpoint**: Keyboard and pointer creation are behaviorally identical and create at most one group.

---

## Phase 7: Documentation and Verification

**Purpose**: Close user documentation, traceability, and all quality gates.

- [x] T026 [P] Update README.md, docs/cli.md, docs/gui-fields.md, CHANGELOG.md, and specs/README.md for S053 behavior and decisions
- [x] T027 Run focused package tests, focused race tests, full tests, and canonical eight-gate verification; record exact evidence in specs/053-trigger-ready-authoring/verification.md
- [x] T028 Audit FR-001 through FR-029 and SC-001 through SC-009, complete every task checkbox, and set S053 lifecycle evidence in specs/053-trigger-ready-authoring/spec.md and specs/README.md

---

## Dependencies and Execution Order

- Setup precedes every implementation task.
- T003 through T005 are independent failing-test tasks; T006 through T009 implement their shared foundation in order.
- User Story 1 depends on the foundation and establishes persisted draft behavior.
- User Story 2 depends on User Story 1's API and editor shape.
- User Stories 3 and 4 are independent after setup and may proceed after the readiness stories stabilize.
- Documentation may proceed separately; final verification depends on all code and documentation.

## Parallel Examples

```text
T003: Readiness helper tests in internal/task/readiness_test.go
T004: Migration preservation tests in internal/store/migration_v11_test.go
T005: Store transition tests in internal/store/store_test.go and internal/store/chains_test.go
```

```text
T022: Navigation section tests in gui/navigation_test.go and gui/app_test.go
T024: Group keyboard submission tests in gui/groups_test.go
```

## Implementation Strategy

1. Establish truthful storage and derived readiness under failing tests.
2. Deliver safe draft persistence as the functional MVP.
3. Enforce and present readiness consistently across every execution surface.
4. Complete the two contained authoring improvements.
5. Run full traceability and canonical verification before commit and publication.
