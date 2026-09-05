# Tasks: GUI Navigation and Information

**Input**: Design documents from `specs/017-gui-info-navigation/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, `quickstart.md`

**Tests**: Required by the feature specification and constitution. Behavioral regressions precede GUI implementation.

## Phase 1: Setup and baseline

**Purpose**: Capture the current navigation and local identity sources before changing the GUI.

- [X] T001 Record the four-item navigation baseline, existing Activity badge ownership, application mark, and build-version source in `specs/017-gui-info-navigation/verification.md`
- [X] T002 Run the current focused GUI tests and record the green pre-change baseline in `specs/017-gui-info-navigation/verification.md`

---

## Phase 2: User Story 1 - Follow a coherent workflow through the sidebar (Priority: P1)

**Goal**: Group Tasks and Groups before Schedule while retaining Activity and its dynamic badge.

**Independent Test**: A headless assertion proves the first four items are Tasks, Groups, Schedule, Activity and badge updates do not alter positions.

### Tests for User Story 1

- [X] T003 [US1] Add a failing management-navigation order regression to `gui/app_test.go`
- [X] T004 [US1] Extend the Activity badge regression with stable index and collection-order assertions in `gui/app_test.go`
- [X] T005 [US1] Run the focused navigation tests and record the expected pre-fix failure in `specs/017-gui-info-navigation/verification.md`

### Implementation for User Story 1

- [X] T006 [US1] Reorder the existing management views to Tasks, Groups, Schedule in `gui/app.go`
- [X] T007 [US1] Run the focused navigation tests and confirm User Story 1 is green

**Checkpoint**: Existing management and Activity views are coherently ordered without any Info implementation dependency.

---

## Phase 3: User Story 2 - Identify the installed application (Priority: P1)

**Goal**: Add a final local Info view containing the official mark, exact version, attribution, and canonical links.

**Independent Test**: Headless tests build Info without daemon data, inspect its mark and version, and match all descriptive link labels to exact destinations.

### Tests for User Story 2

- [X] T008 [P] [US2] Add failing local identity, long-version, brand-resource, and canonical-link regressions in `gui/info_test.go`
- [X] T009 [US2] Update the complete five-item tab contract and Info-last assertion in `gui/app_test.go`
- [X] T010 [US2] Run the focused Info and tab tests and record the expected pre-implementation failure in `specs/017-gui-info-navigation/verification.md`

### Implementation for User Story 2

- [X] T011 [US2] Implement the local scrollable Info hierarchy with aspect-preserving `appIcon`, exact `buildinfo.Version`, attribution, and standard hyperlinks in `gui/info.go`
- [X] T012 [US2] Append Info after the retained Activity tab item in `gui/app.go`
- [X] T013 [US2] Run the focused Info, navigation, and Activity badge tests and confirm both stories are green

**Checkpoint**: The final five-item sidebar and complete local Info contract are independently protected.

---

## Phase 4: Polish and cross-cutting verification

**Purpose**: Synchronize release notes, traceability, and repository quality.

- [X] T014 Update the Unreleased Added section in `CHANGELOG.md` with the navigation and Info outcomes and issue closure intent
- [X] T015 Validate both Spec-Kit checklists against the final implementation and record closure eligibility for #29 and #32 in `specs/017-gui-info-navigation/verification.md`
- [X] T016 Run `git diff --check`, strict UTF-8-without-BOM, and mojibake audits across all changed files
- [X] T017 Run `sh scripts/verify.sh all` in the foreground and record all eight gate results in `specs/017-gui-info-navigation/verification.md`
- [X] T018 Mark every completed task `[X]` in `specs/017-gui-info-navigation/tasks.md` and rerun `/speckit-analyze`
- [X] T019 Commit locally as `feat(017): add GUI information view` with the required co-author trailer

---

## Dependencies and execution order

### Phase dependencies

- Phase 1 has no dependencies.
- User Story 1 depends only on the recorded baseline.
- User Story 2 can begin its separate `info_test.go` work after Phase 1, but the complete tab assertion and integration follow User Story 1.
- Phase 4 depends on both user stories.

### User story dependencies

- **US1** is the navigation-order MVP and has no Info dependency.
- **US2** owns the new view and integrates it after US1's stable Activity item.

### Parallel opportunities

- T008 affects a new test file and may be prepared independently of US1's `app.go` work after baseline capture.
- Documentation and code are kept sequential where they share traceability or implementation files.

## Parallel example: User Story 2

```text
Task T008: define the isolated Info content contract in gui/info_test.go
Task T006: complete the management-view reorder in gui/app.go
Then T009-T012: integrate and prove the final five-item collection
```

## Implementation strategy

### MVP first

1. Capture the existing green baseline.
2. Write red navigation tests.
3. Deliver the independently useful Tasks, Groups, Schedule, Activity order.

### Incremental delivery

1. Navigation grouping completes issue #29's behavior while preserving current Activity terminology.
2. Local Info content and its final tab complete issue #32.
3. Full verification proves the combined slice without backend expansion.
