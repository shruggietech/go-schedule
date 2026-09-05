# Tasks: Structured Desktop Data Tables

**Input**: Design documents from `specs/045-structured-data-tables/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/structured-table.md, quickstart.md

**Tests**: Required by the feature specification and constitution. Behavioral tests are written and observed failing before their implementation tasks.

**Organization**: Tasks are grouped by user story so Tasks, Schedule List, and Activity remain independently testable while sharing one completed foundation.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel with adjacent tasks because it changes a different file
- **[Story]**: User story from `spec.md`
- Every task names its concrete repository path

## Phase 1: Setup and Baseline

**Purpose**: Establish traceability and prove the pre-change behavior before shared implementation begins.

- [x] T001 Record the clean baseline, issue links, existing concatenated row behavior, and initial focused test result in `specs/045-structured-data-tables/verification.md`
- [x] T002 [P] Add S045 as Draft with pending delivery evidence in `specs/README.md`
- [x] T003 Verify existing Go and universal ignore patterns remain sufficient in `.gitignore` without unrelated edits

---

## Phase 2: Foundational Structured-Table Presentation

**Purpose**: Build the shared, blocking row model and presentation behavior used by all three user stories.

**CRITICAL**: No view integration starts until the shared foundation passes its focused tests.

### Red tests

- [x] T004 Add failing tests for column validation, exact width conservation, narrow allocation, and shared header/body alignment in `gui/structured_table_test.go`
- [x] T005 Add failing tests for row cell count, ellipsis, full summary, stable identity, empty data, alternating rows, and 100-row virtualization in `gui/structured_table_test.go`
- [x] T006 [P] Add failing dark/light/follow-system contrast tests for base, alternate, hover, focus, selection, and semantic label states in `gui/theme_test.go`
- [x] T007 Run the focused foundation test selection and record the expected red evidence in `specs/045-structured-data-tables/verification.md`

### Implementation

- [x] T008 Implement `structuredColumn`, `structuredCell`, and `structuredRowModel` validation and normalization helpers in `gui/structured_table.go`
- [x] T009 Implement the bounded weighted responsive column allocator and shared header/row layout in `gui/structured_table.go`
- [x] T010 Implement the reusable truncating, semantic, tappable, double-tappable `structuredRow` widget in `gui/structured_table.go`
- [x] T011 Implement fixed-header, virtualized-body, disclosure, row-selection, and stable-identity reconciliation in `gui/structured_table.go`
- [x] T012 Add one subtle theme-aware alternate-row surface while reusing S044 hover/focus/selection and semantic roles in `gui/theme.go`
- [x] T013 Make the foundation tests pass and record green focused evidence in `specs/045-structured-data-tables/verification.md`

**Checkpoint**: One reusable structured presentation supports fixed headers, full-row interaction, responsive width, truncation disclosure, semantics, and scale.

---

## Phase 3: User Story 1 - Understand and Operate on Tasks (Priority: P1) MVP

**Goal**: Replace ambiguous task strings with labeled Task, Enabled, Lifecycle, Time zone, and Group columns while preserving every task operation.

**Independent Test**: Mixed-state tasks are understandable from headers and plain-language cells; selection, toolbar actions, keyboard navigation, refresh, and double activation target the correct stable ID.

### Red tests

- [x] T014 [P] [US1] Add failing task row-mapping tests for every lifecycle, enabled/disabled independence, group/timezone fallbacks, Unicode, and unknown values in `gui/tasks_test.go`
- [x] T015 [US1] Update failing task integration tests for fixed headers, five aligned cells, ellipsis disclosure, selection reorder/removal, and double activation in `gui/app_test.go`

### Implementation

- [x] T016 [US1] Implement pure Task row normalization and full-value summary mapping in `gui/tasks.go`
- [x] T017 [US1] Replace the task text list with the shared structured view and exact five-column contract in `gui/tasks.go`
- [x] T018 [US1] Preserve toolbar selection, keyboard behavior, stable-ID refresh reconciliation, and double-click editing in `gui/tasks.go`
- [x] T019 [US1] Remove the superseded task-only row widget while preserving general toolbar widgets in `gui/widgets.go`
- [x] T020 [US1] Update widget tests for the reusable structured-row activation contract in `gui/widgets_test.go`
- [x] T021 [US1] Run the independent Tasks focused tests and record results in `specs/045-structured-data-tables/verification.md`

**Checkpoint**: User Story 1 is independently functional and provides the S045 MVP.

---

## Phase 4: User Story 2 - Read the Schedule as Structured Events (Priority: P2)

**Goal**: Present Schedule List occurrences under When, Task, Event, and Outcome headers with normalized glyph/text semantics.

**Independent Test**: Future and all known/unknown past outcomes remain chronological, understandable, selectable for full disclosure, and compatible with range and Calendar controls.

### Red tests

- [x] T022 [US2] Add failing Schedule row-mapping tests for scheduled, success, failure, skipped, caught-up, queued, missing, unknown, duplicate, and Unicode occurrences in `gui/calendar_test.go`
- [x] T023 [US2] Add failing Schedule integration tests for fixed headers, responsive cells, disclosure identity, chronological refresh, range changes, and List/Calendar round trips in `gui/calendar_test.go`

### Implementation

- [x] T024 [US2] Implement pure occurrence identity, event/outcome normalization, glyph, importance, and summary mapping in `gui/schedule.go`
- [x] T025 [US2] Replace Schedule List text rows with the four-column shared structured view in `gui/schedule.go`
- [x] T026 [US2] Preserve chronological refresh, range selection, and Calendar switching while clearing or retaining disclosure safely in `gui/schedule.go`
- [x] T027 [US2] Run the independent Schedule focused tests and record results in `specs/045-structured-data-tables/verification.md`

**Checkpoint**: User Stories 1 and 2 work independently on the shared foundation.

---

## Phase 5: User Story 3 - Triage Activity Consistently (Priority: P2)

**Goal**: Present Activity under When, Severity, Source, and Summary headers with uppercase severity, paired glyph/color semantics, and correct detail activation.

**Independent Test**: Known, empty, and unknown severity rows remain newest-first and work with filtering, clearing, live refresh, alert acknowledgement, and exact detail opening.

### Red tests

- [x] T028 [US3] Add failing Activity row-mapping tests for INFO, WARNING, ERROR, empty, unknown, source/message fallback, alerts, logs, duplicates, and Unicode in `gui/logs_test.go`
- [x] T029 [US3] Add failing Activity integration tests for fixed headers, semantic importance, responsive cells, filtering, clear/acknowledge, refresh, and exact detail activation in `gui/logs_test.go`

### Implementation

- [x] T030 [US3] Extend merged Activity entries with stable source identity and implement pure severity/glyph/summary row mapping in `gui/logs.go`
- [x] T031 [US3] Replace Activity text rows with the four-column shared structured view in `gui/logs.go`
- [x] T032 [US3] Preserve newest-first filtering, Clear View acknowledgement, live refresh, and detail activation by current visible identity in `gui/logs.go`
- [x] T033 [US3] Remove superseded `severityMark` string formatting and make the Activity focused tests pass in `gui/logs.go` and `gui/logs_test.go`
- [x] T034 [US3] Run the independent Activity focused tests and record results in `specs/045-structured-data-tables/verification.md`

**Checkpoint**: All three user stories are independently functional.

---

## Phase 6: Polish and Cross-Cutting Verification

**Purpose**: Complete documentation, qualification guidance, lifecycle evidence, and every project gate.

- [x] T035 [P] Document the Tasks, Schedule List, and Activity column/fallback/semantic contracts in `docs/gui-fields.md`
- [x] T036 [P] Add populated dark/light S045 native qualification steps without claiming attended results in `test/windows/README.md`
- [x] T037 [P] Add the S045 feature and dated architectural decision entry to the unreleased section in `CHANGELOG.md`
- [x] T038 Run `gofmt` and focused GUI tests for all S045 surfaces and record evidence in `specs/045-structured-data-tables/verification.md`
- [x] T039 Run the focused race selection and full `go test ./... -count=1` suite and record evidence in `specs/045-structured-data-tables/verification.md`
- [x] T040 Run `sh scripts/verify.sh all` in the foreground through all eight gates and record evidence in `specs/045-structured-data-tables/verification.md`
- [x] T041 Audit requirements FR-001 through FR-028 and success criteria SC-001 through SC-008 against code/tests in `specs/045-structured-data-tables/verification.md`
- [x] T042 Validate all changed text as UTF-8 without BOM, scan for mojibake, and run `git diff --check`, recording results in `specs/045-structured-data-tables/verification.md`
- [x] T043 Mark every completed task and update S045 to Implemented with objective delivery evidence in `specs/045-structured-data-tables/spec.md`, `specs/045-structured-data-tables/tasks.md`, and `specs/README.md`
- [x] T044 Run `sh scripts/spec-lifecycle-check.sh .` and the documentation link gate after final artifact updates, recording the final result in `specs/045-structured-data-tables/verification.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Starts immediately.
- **Foundation (Phase 2)**: Depends on setup and blocks every user story.
- **User Story 1 (Phase 3)**: Depends only on the foundation and is the MVP.
- **User Story 2 (Phase 4)**: Depends only on the foundation; follows US1 in this single-agent execution to reduce integration churn.
- **User Story 3 (Phase 5)**: Depends only on the foundation; follows US2 for the same reason.
- **Polish (Phase 6)**: Depends on all selected user stories.

### Within Each User Story

- Tests are written and observed failing before implementation.
- Pure row mapping precedes view wiring.
- View wiring precedes regression verification.
- Stable identity is resolved against current visible data before any activation.

### Parallel Opportunities

- T002 and T003 can run independently of T001.
- T006 can run beside T004-T005 because it changes theme tests.
- T014 can run beside the integration-test preparation in T015.
- T035-T037 affect independent documentation files.
- After Phase 2, US1, US2, and US3 are structurally independent, but sequential execution is preferred because all consume `structured_table.go` and benefit from immediate shared feedback.

## Parallel Example: User Story 1

```text
Task A: Add pure Task row-mapping tests in gui/tasks_test.go.
Task B: Update task interaction integration tests in gui/app_test.go.
After both are red, implement the Task mapping and structured view in gui/tasks.go.
```

## Implementation Strategy

### MVP First

1. Complete setup and the shared structured-table foundation.
2. Deliver User Story 1's unambiguous Tasks table.
3. Run the independent Tasks checkpoint before touching Schedule or Activity.

### Incremental Delivery

1. Shared foundation: fixed header, virtualized rows, responsive sizing, disclosure, stable identity, theme roles.
2. Tasks: remove ambiguous status strings and preserve all task operations.
3. Schedule: normalize event/outcome semantics while preserving controls.
4. Activity: normalize severity and preserve diagnostics workflows.
5. Complete shared documentation and all verification gates.

## Notes

- #112 and #113 remain separately traceable; neither closes until its native Windows acceptance evidence is complete.
- #102 and #110 are explicit non-goals.
- The implementation uses no new module dependency, persistence migration, sorting control, or horizontal scroll path. First-round review required one bounded deviation: Calendar now returns the existing stored run ID on past occurrences because equal-time run rows cannot have stable identities without a source discriminator.
