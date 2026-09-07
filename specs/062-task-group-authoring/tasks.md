# Tasks: Task and Group Authoring

**Input**: Design documents from `specs/062-task-group-authoring/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, `quickstart.md`

**Tests**: Required by the specification and constitution. Contract tests are written before their corresponding implementation.

**Organization**: Tasks are grouped by independently testable user story after shared contracts are established.

## Phase 1: Setup and Traceability

**Purpose**: Establish the S062 lifecycle, shared paths, and traceability before behavior changes.

- [x] T001 Register S062 as Draft in `specs/README.md` and point `.specify/feature.json` plus `CLAUDE.md` at `specs/062-task-group-authoring/`
- [x] T002 [P] Audit current Fyne task and group behavior against #153 and #192 in `specs/062-task-group-authoring/research.md`
- [x] T003 [P] Define detailed-list, desktop bridge, safe-example, and stale-write contracts in `specs/062-task-group-authoring/contracts/`
- [x] T004 Validate all requirements checklists and record the resolved clarification decisions in `specs/062-task-group-authoring/spec.md`

---

## Phase 2: Foundational Contracts

**Purpose**: Build the authoritative read/write and safe-example seams that block every story.

**Critical**: No story implementation begins until this phase passes focused tests.

- [x] T005 [P] Write failing exact platform mapping, allowlist, parsing, and bounded native execution tests in `internal/commandexample/suggestion_test.go`
- [x] T006 Implement the canonical safe platform suggestion contract in `internal/commandexample/suggestion.go`
- [x] T007 [P] Write failing backward-compatibility and detailed task-list server tests in `internal/api/server/tasks_test.go`
- [x] T008 Implement opt-in `details=true` task responses in `internal/api/server/tasks.go`
- [x] T009 [P] Write failing detailed-list and group-update client tests in `internal/api/client/tasks_test.go` and `internal/api/client/groups_test.go`
- [x] T010 Add `ListTaskDetails` and `UpdateGroup` client methods in `internal/api/client/methods.go` and `internal/api/client/groups.go`
- [x] T011 [P] Write failing working-directory and run-identity clear plus inactive group creation contract tests in `internal/api/server/update_test.go` and `internal/api/server/groups_test.go`
- [x] T012 Add explicit safe task clear semantics and optional inactive group creation in `internal/api/server/update.go` and `internal/api/server/groups.go`
- [x] T013 [P] Define transport-neutral task, group, draft, preview, workspace, and operation models in `desktop/taskgroup/model.go`
- [x] T014 [P] Write failing local adapter mapping tests in `desktop/taskgroup/local_test.go`
- [x] T015 Implement the protected local task/group backend adapter in `desktop/taskgroup/local.go`
- [x] T016 [P] Write failing service tests for safe errors, full paths, effective state, exact parsing, differential updates, stale rejection, overwrite, bounded calls, and one hundred repeated mutation/cancellation cycles in `desktop/taskgroup/service_test.go`
- [x] T017 Implement daemon-authoritative task/group service behavior in `desktop/taskgroup/service.go`
- [x] T018 Wire the service into typed Wails facade methods and facade tests in `desktop/app.go`, `desktop/main.go`, and `desktop/app_test.go`

**Checkpoint**: The desktop can load and mutate task/group contracts without exposing daemon internals or duplicating domain logic.

---

## Phase 3: User Story 1 - Understand Tasks and Groups at a Glance (Priority: P1)

**Goal**: Deliver an accurate, scalable, keyboard-accessible Tasks overview with stable identity and clear effective state.

**Independent Test**: Load empty, mixed-state, one-hundred-task, and twenty-level-group fixtures, then exercise search, filters, selection, refresh, and long values.

- [x] T019 [P] [US1] Write failing workspace bridge and live-refresh store tests in `desktop/frontend/src/tasks/store.test.ts`
- [x] T020 [P] [US1] Write failing overview, filtering, selection, state, hierarchy, empty, and scale component tests in `desktop/frontend/src/tasks/TasksPage.test.tsx`
- [x] T021 [US1] Implement typed task/group bridge models and calls in `desktop/frontend/src/tasks/model.ts` and `desktop/frontend/src/tasks/bridge.ts`
- [x] T022 [US1] Implement authority refresh, event coalescing, stable selection, and draft-preservation state in `desktop/frontend/src/tasks/store.ts`
- [x] T023 [US1] Implement task table, filters, detail summary, group hierarchy, and empty/loading states in `desktop/frontend/src/tasks/TasksPage.tsx` and `desktop/frontend/src/tasks/GroupsPanel.tsx`
- [x] T024 [US1] Add reusable selectable-table and inline state primitives in `desktop/frontend/src/components/index.tsx`
- [x] T025 [US1] Replace the Tasks placeholder while preserving other honest route placeholders in `desktop/frontend/src/App.tsx`

**Checkpoint**: The Tasks route independently explains the current task/group inventory and remains usable at the required scale.

---

## Phase 4: User Story 2 - Create a Known-Good Task (Priority: P1)

**Goal**: Let a new operator create an inactive safe example with exact preview and keyboard guidance.

**Independent Test**: Exercise each execution-host platform mapping, one-shot empty Tab insertion, traversal, basic save, Run now handoff, and retired-guidance absence.

- [x] T026 [P] [US2] Write failing new-task defaults, platform suggestion, command preview, Tab, traversal, clearing, and reinsertion tests in `desktop/frontend/src/tasks/TaskEditor.test.tsx`
- [x] T027 [P] [US2] Write failing Fyne shared mapping and focused Tab insertion tests in `gui/command_suggestion_test.go`
- [x] T028 [US2] Implement the accessible basic task editor and exact command preview in `desktop/frontend/src/tasks/TaskEditor.tsx`
- [x] T029 [US2] Implement platform-host suggestion guidance and one-shot Tab behavior in `desktop/frontend/src/tasks/CommandField.tsx`
- [x] T030 [US2] Adapt the shipping Fyne new-task command entry to the canonical mapping in `gui/command_suggestion.go` and `gui/editor.go`
- [x] T031 [US2] Add the inactive save and guided Run now to Activity handoff in `desktop/frontend/src/tasks/TasksPage.tsx`

**Checkpoint**: A first-time operator can create and intentionally run a harmless platform-native task without advanced settings.

---

## Phase 5: User Story 3 - Edit Complete Task Intent (Priority: P1)

**Goal**: Round-trip every task field and scheduling policy without silent change or partial mutation.

**Independent Test**: Edit recurring, one-off, manual-only, advanced-policy, draft, and legacy-empty tasks; cover invalid fields, preview failure, cancellation, stale reload, and explicit overwrite.

- [x] T032 [P] [US3] Write failing complete field round-trip, disclosure, validation focus, schedule preview, dirty, and stale-draft tests in `desktop/frontend/src/tasks/TaskEditor.test.tsx`
- [x] T033 [P] [US3] Write failing task service differential-update and field-error mapping cases in `desktop/taskgroup/service_test.go`
- [x] T034 [US3] Implement complete basic and advanced task fields with deterministic draft conversion in `desktop/frontend/src/tasks/TaskEditor.tsx`
- [x] T035 [US3] Implement recurring, one-off, and manual-only preview states with five occurrences in `desktop/frontend/src/tasks/ScheduleEditor.tsx`
- [x] T036 [US3] Implement field-linked safe validation, first-error focus, pending suppression, cancellation, and stale reload/overwrite choices in `desktop/frontend/src/tasks/TaskEditor.tsx`
- [x] T037 [US3] Preserve exact values and selection after authoritative saves and event refreshes in `desktop/frontend/src/tasks/store.ts`

**Checkpoint**: Every current task value can be inspected and intentionally changed with authoritative outcomes.

---

## Phase 6: User Story 4 - Organize Work with Nested Groups (Priority: P2)

**Goal**: Deliver complete nested group creation, editing, movement, cascade state, task assignment, and deletion.

**Independent Test**: Use duplicate names, twenty levels, disabled ancestors, draft disabled groups, cycle attempts, reparenting, task movement, cancellation, and descendant deletion.

- [x] T038 [P] [US4] Write failing hierarchy derivation, descendant exclusion, duplicate path, cascade count, and stale group service tests in `desktop/taskgroup/service_test.go`
- [x] T039 [P] [US4] Write failing group editor, hierarchy keyboard, task movement, cascade confirmation, cycle failure, and deletion tests in `desktop/frontend/src/tasks/GroupsPanel.test.tsx`
- [x] T040 [US4] Implement group creation, rename, reparent, root/descendant choices, and disabled draft behavior in `desktop/frontend/src/tasks/GroupEditor.tsx`
- [x] T041 [US4] Implement nested hierarchy disclosure, effective-state reasons, counts, stable selection, and task movement in `desktop/frontend/src/tasks/GroupsPanel.tsx`
- [x] T042 [US4] Implement target-aware group toggle and delete confirmations with authoritative refresh in `desktop/frontend/src/tasks/GroupsPanel.tsx`

**Checkpoint**: Nested group lifecycle is complete and understandable without ambiguous names or hidden cascade effects.

---

## Phase 7: User Story 5 - Perform Safe Task Actions (Priority: P2)

**Goal**: Deliver bounded, target-aware Run now, enable, disable, and delete actions with honest availability.

**Independent Test**: Exercise ready, manual-only, disabled-group, invalid, disconnected, concurrently deleted, repeated-activation, accepted, and failed outcomes.

- [x] T043 [P] [US5] Write failing action availability, target confirmation, duplicate suppression, safe failure, selection, and Activity handoff tests in `desktop/frontend/src/tasks/TaskActions.test.tsx`
- [x] T044 [P] [US5] Add repeated-action and missing-entity service contracts in `desktop/taskgroup/service_test.go`
- [x] T045 [US5] Implement task action availability and target-aware confirmation controls in `desktop/frontend/src/tasks/TaskActions.tsx`
- [x] T046 [US5] Implement accepted, rejected, stale, disconnected, and late-result handling in `desktop/frontend/src/tasks/store.ts`
- [x] T047 [US5] Integrate safe announcements and preserve selection across action refresh in `desktop/frontend/src/tasks/TasksPage.tsx`

**Checkpoint**: Every task action is safely reachable, bounded, and honest about target and outcome.

---

## Phase 8: Cross-Cutting Polish and Verification

**Purpose**: Close accessibility, documentation, performance, platform, and lifecycle evidence across all stories.

- [x] T048 [P] Replace unsafe first-run examples with the canonical platform walkthrough in `docs/gui-fields.md`
- [x] T049 [P] Add documentation/source consistency checks for all three exact suggestions in `internal/commandexample/suggestion_test.go`
- [x] T050 [P] Expand Chromium journeys for accessibility, keyboard, narrow-window, high-zoom, reduced-motion, and one-hundred-row behavior in `desktop/frontend/e2e/tasks.spec.ts`
- [x] T051 Add an end-to-end create, Run now, captured Activity record, and recognizable output test in `test/integration/guided_example_test.go`, then run it in the existing Windows, macOS, and Linux production desktop matrix in `.github/workflows/ci.yml`
- [x] T052 Audit event and result payloads for command, environment, input, identity, path, and backend detail leakage in `desktop/taskgroup/service_test.go` and `desktop/connection/local_test.go`
- [x] T053 Run focused root, desktop race, frontend unit/build, Chromium, and native Wails validation from `specs/062-task-group-authoring/quickstart.md`
- [x] T054 Run `go run ./scripts/github-format`, specification lifecycle checks, and `/speckit-analyze` follow-up against all S062 artifacts
- [x] T055 Run the canonical eight-gate `scripts/verify.sh all` workflow and record results in `specs/062-task-group-authoring/verification.md`
- [x] T056 Advance S062 to Implemented, check every task, update `specs/README.md`, `CHANGELOG.md`, and Delivery evidence, then run UTF-8 and mojibake sanity checks

---

## Dependencies and Execution Order

- Phase 1 has no dependencies.
- Phase 2 depends on Phase 1 and blocks every user story.
- User Story 1 depends on the Phase 2 workspace contract.
- User Story 2 depends on the Phase 2 command suggestion and save contract, but is independently testable from the overview with a focused fixture.
- User Story 3 depends on the shared editor established by User Story 2.
- User Story 4 depends only on the Phase 2 group service and can be tested separately from task editing.
- User Story 5 depends on the User Story 1 selection surface and Phase 2 action methods.
- Phase 8 depends on all selected user stories.

## Parallel Opportunities

- Server, client, command suggestion, desktop model, and adapter test files in Phase 2 can be authored independently before implementation.
- Store and component tests for User Story 1 can be authored independently.
- Fyne adapter tests and React editor tests for User Story 2 can be authored independently.
- Group service and group component tests can be authored independently in User Story 4.
- Documentation, Chromium journeys, and security payload tests can proceed in parallel after feature behavior stabilizes.

## Implementation Strategy

1. Establish backward-compatible daemon and desktop contracts under failing tests.
2. Deliver the read-only overview as the first independently demonstrable increment.
3. Add safe basic creation, then complete task editing without changing the overview contract.
4. Add groups and task actions through the same authoritative store.
5. Close cross-platform, accessibility, documentation, security, and canonical verification together.

## Notes

- `[P]` marks tasks that affect different files and have no dependency on an incomplete task in the same phase.
- Every user-story task includes its `[US#]` traceability label and exact file path.
- Publication, pull-request review, merge, and housekeeping are workflow evidence, not implementation tasks.
