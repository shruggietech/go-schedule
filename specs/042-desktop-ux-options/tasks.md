# Tasks: Desktop UX and Options

**Input**: Design documents from `/specs/042-desktop-ux-options/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Regression tests are mandatory. Add each story's tests before its implementation and preserve the red-to-green evidence.

**Organization**: Tasks are grouped by independently testable user story. The five linked issues share the desktop shell, theme, and native visual evidence boundary, so they form one coherent slice.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it changes a different file after shared expectations are fixed
- **[Story]**: Maps a task to one user story from spec.md

## Phase 1: Setup and Baseline

**Purpose**: Fix the feature boundary, establish authoritative framework behavior, and preserve current GUI contracts.

- [x] T001 Verify clean synchronized `main`, create `codex/042-desktop-ux-options`, and record GitHub issues #101, #103, #104, #105, and #106 in `specs/042-desktop-ux-options/spec.md`
- [x] T002 Research Fyne 2.8.1 list gestures, selectable labels, preferences, settings, theme variants, and custom-theme delegation in `specs/042-desktop-ux-options/research.md`
- [x] T003 Record the native Windows visual evidence boundary and #94 dependency in `specs/042-desktop-ux-options/spec.md`, `plan.md`, and `contracts/desktop-ux.md`

---

## Phase 2: Foundational Desktop State

**Purpose**: Establish bounded appearance, storage, editor-ownership, and shutdown models before wiring views.

**CRITICAL**: Tests must fail against the S041 implementation before production state is added.

- [x] T004 [P] Add failing appearance default, validation, persistence, and reset tests in `gui/appearance_test.go`
- [x] T005 [P] Add failing theme matrix, palette, and font-resource tests in `gui/theme_test.go`
- [x] T006 [P] Add failing storage classification, exact-stat, unavailable-path, and lifecycle-copy tests in `gui/storage_locations_test.go`
- [x] T007 [P] Add failing one-editor ownership and concurrent one-shot shutdown tests in `gui/app_test.go`
- [x] T008 Run the foundational focused tests and record their expected failures in `specs/042-desktop-ux-options/verification.md`
- [x] T009 Implement validated appearance state and Fyne preference persistence in `gui/appearance.go`
- [x] T010 Extend the immutable brand theme for Dark, Light, Follow system, Brand, System, and Monospace in `gui/theme.go`
- [x] T011 Implement bounded read-only storage resolution in `gui/storage_locations.go`
- [x] T012 Implement editor ownership and idempotent shutdown coordinators in `gui/app.go`
- [x] T013 Re-run foundational tests and confirm the state layer turns green

**Checkpoint**: Shared state is deterministic, race-safe where concurrent, and independently testable without constructing the complete shell.

---

## Phase 3: User Story 1 - Choose and Retain a Readable Appearance (Priority: P1) MVP

**Goal**: Provide immediate, persistent, safe appearance choices and correct the two affected Info labels' semantic rendering contract.

**Independent Test**: Construct the headless Options and Info views over controlled preferences, select every appearance combination, simulate restart, and assert full-app theme replacement plus centered unwrapped labels.

### Tests for User Story 1

- [x] T014 [P] [US1] Add failing Options selection, immediate application, persistence, and reset interaction tests in `gui/options_test.go`
- [x] T015 [P] [US1] Add failing centered, unwrapped, exact-text, and dynamic-version assertions in `gui/info_test.go`

### Implementation for User Story 1

- [x] T016 [US1] Build the Appearance section and wire application-wide persistence/application in `gui/options.go` and `gui/app.go`
- [x] T017 [US1] Remove wrapping from the version and attribution labels without changing their copy or source in `gui/info.go`
- [x] T018 [US1] Run the focused appearance and Info tests and confirm every valid and invalid state passes

**Checkpoint**: Appearance and Info behavior work independently through the headless driver; native sharpness remains unclaimed pending #94.

---

## Phase 4: User Story 2 - Understand Local Application Storage (Priority: P1)

**Goal**: Present and copy a truthful, read-only inventory of application-owned storage and removal behavior.

**Independent Test**: Supply controlled platform, configuration, executable, app-storage, clipboard, and stat inputs and verify every displayed value and lifecycle statement.

### Tests for User Story 2

- [x] T019 [US2] Add failing storage-row rendering, selectable-text, exact-copy, scope, existence, and unavailable-state tests in `gui/options_test.go`

### Implementation for User Story 2

- [x] T020 [US2] Build the Storage section with selectable paths and explicit Copy controls in `gui/options.go`
- [x] T021 [US2] Wire active configuration, Fyne storage root, executable directory, and the shared cleanup-evidence resolver in `gui/app.go`, `gui/storage_locations.go`, and `internal/winuninstall/`
- [x] T022 [US2] Run focused storage resolver and Options interaction tests and confirm no path is created, traversed, or misclassified

**Checkpoint**: Storage transparency works independently and makes no deletion or external-ownership claim.

---

## Phase 5: User Story 3 - Navigate Comfortably and Exit Predictably (Priority: P2)

**Goal**: Replace the narrow tab strip with a spacious, ordered, accessible rail and bottom-right orderly Exit command.

**Independent Test**: Construct the complete headless shell at default and minimum dimensions, select every destination, update the Activity badge, and activate every close path repeatedly.

### Tests for User Story 3

- [x] T023 [P] [US3] Add failing destination order, selection, rail-width, minimum-window, Activity badge, and Exit-layout tests in `gui/navigation_test.go`
- [x] T024 [P] [US3] Extend failing close-entry integration assertions in `gui/app_test.go`

### Implementation for User Story 3

- [x] T025 [US3] Implement the package-owned destination/content/command shell in `gui/navigation.go`
- [x] T026 [US3] Replace AppTabs, insert Options above Info, retain Activity badge updates, and route all Exit controls through `requestClose` in `gui/app.go`
- [x] T027 [US3] Run focused navigation and shutdown tests at 1280 by 800 and 800 by 600

**Checkpoint**: Navigation and Exit work without depending on native rendering; their exact-candidate visual criteria remain in #94.

---

## Phase 6: User Story 4 - Edit the Intended Task Directly (Priority: P2)

**Goal**: Open exactly one correct task editor on row double activation while preserving selection, keyboard, toolbar, fallback, and live-refresh behavior.

**Independent Test**: Drive a bound row through single and double activation before and after task reorder/removal and through detail lookup failure.

### Tests for User Story 4

- [x] T028 [P] [US4] Add failing row single/double activation and identity-binding tests in `gui/widgets_test.go`
- [x] T029 [P] [US4] Add failing reordered, stale, degraded-detail, toolbar-sharing, and duplicate-editor tests in `gui/app_test.go`

### Implementation for User Story 4

- [x] T030 [US4] Implement stable-ID single-selection forwarding and double activation in the task-row widget in `gui/widgets.go`
- [x] T031 [US4] Resolve current task identity and share the guarded editor path in `gui/tasks.go` and `gui/editor.go`
- [x] T032 [US4] Run focused task interaction tests and confirm single selection and keyboard behavior remain unchanged

**Checkpoint**: Every valid row opens the intended editor once; stale identities fail closed.

---

## Phase 7: Polish and Cross-Cutting Verification

**Purpose**: Complete native instructions, traceability, lifecycle status, encoding integrity, and repository-wide evidence.

- [x] T033 [P] Document desktop defaults, storage boundaries, and exact-candidate native GUI checks in `docs/options.md` and `test/windows/README.md`
- [x] T034 [P] Add the S042 feature and dated desktop UX decision to `CHANGELOG.md`
- [x] T035 [P] Add S042 to the chronological feature index in `specs/README.md`
- [x] T036 Update `specs/042-desktop-ux-options/spec.md` to Implemented and create `specs/042-desktop-ux-options/verification.md` with red/green, focused, native-boundary, and gate evidence
- [x] T037 Run `C:\Program Files\Git\bin\bash.exe scripts/verify.sh all` and record all eight canonical gate results in `specs/042-desktop-ux-options/verification.md`
- [x] T038 Perform a systematic local requirements, code, concurrency, history, safety, and evidence review; fix every high-confidence finding
- [x] T039 Validate UTF-8 without BOM, scan changed files for mojibake, run `git diff --check`, and mark every completed task in `specs/042-desktop-ux-options/tasks.md`

---

## Dependencies and Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Complete.
- **Foundational State (Phase 2)**: Blocks all views and interactions.
- **Appearance (Phase 3)**: Depends on validated appearance and theme state.
- **Storage (Phase 4)**: Depends on the storage resolver; otherwise independent of appearance.
- **Navigation (Phase 5)**: Depends on Options content and shutdown coordination.
- **Task Editing (Phase 6)**: Depends only on editor ownership and can be validated independently.
- **Polish (Phase 7)**: Depends on all stories.

### User Story Dependencies

- **User Story 1** establishes appearance behavior consumed by the Options view.
- **User Story 2** adds independently testable storage content to the same view.
- **User Story 3** places that complete Options view in the final shell.
- **User Story 4** is interaction-independent but shares the foundational editor guard.

### Parallel Opportunities

- T004 through T007 touch independent test surfaces once expectations are fixed.
- T014 and T015 are independent appearance and Info regressions.
- T023 and T024 separate shell geometry from application close integration.
- T028 and T029 separate row mechanics from application-level task resolution.
- T033 through T035 touch independent documentation files after behavior stabilizes.
- No multi-agent delegation is used because current session policy reserves it for explicit operator requests.

## Implementation Strategy

1. Make foundational pure-state and headless contracts fail against S041.
2. Implement the minimum shared state needed for all stories and turn it green.
3. Deliver appearance, storage, navigation, and task editing in priority order,
   preserving a red-to-green checkpoint for each.
4. Run focused, full GUI, race, and canonical verification.
5. Review the complete diff against every functional requirement and evidence boundary.
6. Publish once authorized local verification is complete, then close automatic
   review, CI, and one manually triggered second review round.
7. Leave merge, tag, release, and exact-candidate issue closure to their existing authority gates.

## Notes

- The operator has already authorized push, PR creation, and verified in-scope review-fix pushes for S042.
- At most one manual `@Codex review` request may follow the automatic first review round.
- Exact-candidate native visual evidence remains required by #94 and is not replaced by headless tests.
- Merge, tag, release, and material scope expansion remain unauthorized.
