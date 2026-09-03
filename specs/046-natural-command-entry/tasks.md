# Tasks: Natural Command Entry

**Input**: Design documents from `specs/046-natural-command-entry/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/portable-command-line.md, quickstart.md

**Tests**: Required by the feature specification and constitution. Behavioral tests are written and observed failing before their implementation tasks.

**Organization**: Tasks are grouped by user story. All three P1 stories form one coherent compatibility-safe command-authoring capability.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel with adjacent tasks because it changes a different file
- **[Story]**: User story from `spec.md`
- Every task names its concrete repository path

## Phase 1: Setup and Baseline

**Purpose**: Establish issue traceability and prove the pre-change behavior.

- [x] T001 Record the clean branch baseline, issue #110 evidence, current two-field behavior, and focused pre-change test result in `specs/046-natural-command-entry/verification.md`
- [x] T002 [P] Add S046 as Draft with pending delivery evidence in `specs/README.md`
- [x] T003 Verify existing Go and universal ignore patterns remain sufficient in `.gitignore` without unrelated edits

---

## Phase 2: Foundational Portable Command Grammar

**Purpose**: Establish the lossless platform-independent authoring boundary before changing the GUI.

**CRITICAL**: No editor integration begins until parser and formatter identity are green.

### Red tests

- [x] T004 Add failing lexical tests for whitespace, quote styles, empty values, adjacent segments, backslashes, Unicode, CR/LF/tab, shell punctuation, unmatched quotes, NUL, and position-aware errors in `internal/commandline/commandline_test.go`
- [x] T005 Add failing canonical formatting and parse-format-parse identity tables for existing unusual invocations in `internal/commandline/commandline_test.go`
- [x] T006 Add failing seeded fuzz properties for parser/formatter identity and canonical stability in `internal/commandline/commandline_test.go`
- [x] T007 Run the focused commandline selection and record expected red evidence in `specs/046-natural-command-entry/verification.md`

### Implementation

- [x] T008 Implement the Direct Invocation value and position-aware syntax errors in `internal/commandline/commandline.go`
- [x] T009 Implement the portable single-pass parser without shell evaluation in `internal/commandline/commandline.go`
- [x] T010 Implement the canonical lossless formatter and exact escaped display helper in `internal/commandline/commandline.go`
- [x] T011 Make all parser, formatter, and fuzz-seed tests pass and record focused green evidence in `specs/046-natural-command-entry/verification.md`

**Checkpoint**: Any process-valid stored invocation has one stable portable editor representation, and invalid text cannot produce a stale invocation.

---

## Phase 3: User Story 1 - Enter a Familiar Command Line (Priority: P1) MVP

**Goal**: Replace Command plus one-argument-per-line Arguments with one roomy, validatable Command line field.

**Independent Test**: Enter documented Windows/POSIX examples plus all edge values in the headless editor and observe exact structured create requests.

### Red tests

- [x] T012 [US1] Add failing headless tests for one Command line field, multiline configuration, at least six visible rows, and vertical growth in `gui/editor_test.go`
- [x] T013 [US1] Replace old split-argument tests with failing valid/invalid command-entry, Save-gating, and exact create-submission tests in `gui/editor_test.go`
- [x] T014 [US1] Run the focused editor construction and submission selection and record expected red evidence in `specs/046-natural-command-entry/verification.md`

### Implementation

- [x] T015 [US1] Replace taskEditor command/args state with one multiline Command line entry in `gui/editor.go`
- [x] T016 [US1] Add the six-row minimum and vertically expanding entry layout in `gui/editor.go`
- [x] T017 [US1] Parse the current draft once for validation and exact create submission in `gui/editor.go`
- [x] T018 [US1] Remove superseded one-argument-per-line splitting and obsolete field guidance in `gui/editor.go`
- [x] T019 [US1] Make the User Story 1 focused tests pass and record results in `specs/046-natural-command-entry/verification.md`

**Checkpoint**: New tasks can be authored naturally in one field and submit the exact direct invocation.

---

## Phase 4: User Story 2 - Understand Exactly What Will Run (Priority: P1)

**Goal**: Show an unambiguous program and numbered-argument preview, clear invalid state, and portable non-shell help.

**Independent Test**: For valid, empty, and invalid drafts, compare preview text and Save state with the parsed invocation and documented contract.

### Red tests

- [x] T020 [US2] Add failing pure preview tests for program labels, argument count/order, empty values, escaped quotes/backslashes/tabs/CR/LF, and no ambiguous reconstructed command in `gui/editor_data_test.go`
- [x] T021 [US2] Add failing editor tests for live valid preview, stale-preview clearing, actionable error location, shell-punctuation literals, and corrected recovery in `gui/editor_test.go`
- [x] T022 [US2] Add failing help-content tests for portable grammar, Windows/POSIX examples, no expansion, and explicit-shell responsibility in `gui/editor_test.go`

### Implementation

- [x] T023 [US2] Replace reconstructed command display with separately labeled exact preview composition in `gui/editor_data.go`
- [x] T024 [US2] Wire valid, empty, and invalid preview states from one parsed editor snapshot in `gui/editor.go`
- [x] T025 [US2] Replace old Command/Arguments help with portable grammar and explicit-shell guidance in `gui/editor.go`
- [x] T026 [US2] Make the User Story 2 focused tests pass and record results in `specs/046-natural-command-entry/verification.md`

**Checkpoint**: Users can audit every process boundary and understand that no shell is inferred.

---

## Phase 5: User Story 3 - Edit Existing Tasks Without Loss (Priority: P1)

**Goal**: Canonically represent every existing invocation and preserve its exact values through edit, API/store, and native execution boundaries.

**Independent Test**: Open and save unusual existing task fixtures, then exercise structured API/store and native helper-process round trips on the current host.

### Red tests

- [x] T027 [US3] Add failing prefill and unchanged-update identity tests for empty, quoted, Unicode, spaced-path, repeated, backslash, tab, CR/LF, and shell-punctuation arguments in `gui/editor_prefill_test.go`
- [x] T028 [P] [US3] Add failing exact unusual-argument API/store round-trip regression coverage in `internal/api/server/tasks_test.go` and `internal/store/store_test.go` only where existing coverage is insufficient
- [x] T029 [P] [US3] Add failing direct shell-punctuation executor regression coverage in `internal/executor/executor_test.go`
- [x] T030 [P] [US3] Add failing helper-process argument-vector proof for POSIX hosts in `internal/commandline/commandline_unix_test.go`
- [x] T031 [P] [US3] Add failing helper-process argument-vector and hidden-console compatibility proof for Windows in `internal/commandline/commandline_windows_test.go`

### Implementation

- [x] T032 [US3] Prefill existing tasks through canonical command formatting and submit exact update values in `gui/editor.go`
- [x] T033 [US3] Resolve any compatibility gaps found by API, store, executor, POSIX, or Windows tests without changing public or persistence shapes in the corresponding `internal/` packages
- [x] T034 [US3] Run the independent existing-task, API/store, executor, and native commandline tests and record results in `specs/046-natural-command-entry/verification.md`

**Checkpoint**: New and pre-S046 tasks share one lossless authoring projection and unchanged direct execution semantics.

---

## Phase 6: Polish and Cross-Cutting Verification

**Purpose**: Complete user documentation, lifecycle evidence, and every project gate.

- [x] T035 [P] Replace the obsolete two-field and one-argument-per-line documentation with the portable grammar, exact preview, and examples in `docs/gui-fields.md`
- [x] T036 [P] Update Windows installed-service task-authoring guidance for the combined field while retaining LocalSystem path/access cautions in `docs/INSTALL-windows.md`
- [x] T037 [P] Add the S046 feature and deliberate portable-grammar decision to the Unreleased section in `CHANGELOG.md`
- [x] T038 Run `gofmt` and the focused commandline/GUI/API/store/executor tests and record evidence in `specs/046-natural-command-entry/verification.md`
- [x] T039 Run the focused race selection and full `go test ./... -count=1` suite and record evidence in `specs/046-natural-command-entry/verification.md`
- [x] T040 Run `sh scripts/verify.sh all` through all eight gates and record evidence in `specs/046-natural-command-entry/verification.md`
- [x] T041 Audit requirements FR-001 through FR-026 and success criteria SC-001 through SC-008 against code/tests in `specs/046-natural-command-entry/verification.md`
- [x] T042 Validate changed text as UTF-8 without BOM, scan for mojibake, run `git diff --check`, and record results in `specs/046-natural-command-entry/verification.md`
- [x] T043 Mark every completed task and update S046 to Implemented with objective delivery evidence in `specs/046-natural-command-entry/spec.md`, `specs/046-natural-command-entry/tasks.md`, and `specs/README.md`
- [x] T044 Run `sh scripts/spec-lifecycle-check.sh .` and the documentation link gate after final artifact updates, recording the final result in `specs/046-natural-command-entry/verification.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Starts immediately.
- **Foundation (Phase 2)**: Depends on setup and blocks all user stories.
- **User Story 1 (Phase 3)**: Depends on the grammar foundation and is the MVP.
- **User Story 2 (Phase 4)**: Depends on the grammar foundation and the editor introduced by US1.
- **User Story 3 (Phase 5)**: Depends on the foundation and US1 editor boundary; native tests can be prepared alongside preview work.
- **Polish (Phase 6)**: Depends on all three user stories.

### Within Each User Story

- Tests are written and observed failing before implementation.
- The parser result is the only boundary used by validation, preview, and submission.
- Canonical formatting precedes existing-task prefill.
- Native process tests compare received values, not reconstructed command strings.

### Parallel Opportunities

- T002 and T003 can run independently of T001.
- T020-T022 affect separate pure/display/help concerns once US1 is complete.
- T028-T031 target independent API/store, executor, POSIX, and Windows evidence surfaces.
- T035-T037 affect independent documentation files.

## Parallel Example: User Story 3

```text
Task A: Add existing-task editor identity tests in gui/editor_prefill_test.go.
Task B: Add API/store exact-value coverage.
Task C: Add executor literal-punctuation coverage.
Task D: Add POSIX and Windows helper-process argument proofs.
After all are red, integrate canonical prefill and resolve only demonstrated compatibility gaps.
```

## Implementation Strategy

### MVP First

1. Complete the parser/formatter foundation.
2. Replace the two fields with one validated multiline editor.
3. Prove exact create requests before adding preview polish.

### Incremental Delivery

1. Portable grammar and canonical identity.
2. Natural one-field entry and layout.
3. Exact preview and help.
4. Existing-task, API/store, executor, and native-host compatibility.
5. Documentation and all canonical gates.

## Notes

- #110 is the only issue in S046 and closes only when every acceptance criterion and required verification is complete.
- #102 remains explicitly out of scope.
- The user authorized automatic push, PR publication, and at most two external Codex review rounds. Final merge remains with the user.
