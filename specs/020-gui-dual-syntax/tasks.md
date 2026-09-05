# Tasks: GUI Dual-Syntax Scheduling

**Input**: Design documents from `specs/020-gui-dual-syntax/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Required by the specification and constitution. Behavioral tests are recorded red before production changes.

## Phase 1: Setup and Baseline

- [X] T001 Run focused GUI, schedule-input, and API server tests and record the green baseline in `specs/020-gui-dual-syntax/verification.md`
- [X] T002 Record current human-only validation/request hints, exact expression prefill, legacy preservation, and one-off isolation evidence

---

## Phase 2: User Story 1 - Create and Preview Cron (Priority: P1)

**Goal**: A supported cron expression validates, previews, and creates through the existing Schedule field with cron identity.

### Tests

- [X] T003 [US1] Add failing cron/human automatic classification, normalized preview identity, create identity, five-word-human, invalid cron, and named fidelity-refusal tests in GUI tests
- [X] T004 [US1] Run the focused US1 tests and record the expected pre-implementation failures

### Implementation

- [X] T005 [US1] Add one editor helper over `scheduleinput.Parse` and use it for local recurring validity
- [X] T006 [US1] Send the helper's normalized expression and selected syntax in preview and recurring create requests
- [X] T007 [US1] Run US1 and central input focused tests green and record the request/error evidence

---

## Phase 3: User Story 2 - Edit in Original or Switched Syntax (Priority: P2)

**Goal**: Retained cron prefill remains exact, while save and preview follow the current text after either syntax switch.

### Tests

- [X] T008 [US2] Add failing exact cron prefill/save, cron-to-human, human-to-cron, whitespace, expressionless legacy, one-off, and human-anchor regressions
- [X] T009 [US2] Run the focused US2 tests and record the expected pre-implementation failures

### Implementation

- [X] T010 [US2] Carry selected recurring syntax through `taskForm` and recurring update requests
- [X] T011 [US2] Preserve current-expression reclassification, one-off omission, legacy blank preservation, and human Start-at behavior
- [X] T012 [US2] Run the complete GUI and schedule-input focused suites green and record round-trip evidence

---

## Phase 4: User Story 3 - Explain the Dual-Syntax Field (Priority: P3)

**Goal**: GUI-local help tells the truth without absorbing issue #52's broad documentation work.

- [X] T013 [US3] Add or update focused help/content assertions for human-first examples, five-field cron, and the fidelity guide
- [X] T014 [US3] Update the field hint and in-editor help in `gui/editor.go`
- [X] T015 [US3] Narrowly update `docs/gui-fields.md`, the newly false GUI sentence in `docs/cron.md`, and the chronological Unreleased feature/decision entries in `CHANGELOG.md`
- [X] T016 [US3] Run focused GUI tests and documentation checks green

---

## Phase 5: Polish and Cross-Cutting Verification

- [X] T017 Validate both Spec-Kit checklists and rerun `/speckit-analyze`
- [X] T018 Run `git diff --check`, strict UTF-8-without-BOM, and mojibake audits across every changed file
- [X] T019 Run `sh scripts/verify.sh all` in the foreground and record all eight gate results
- [X] T020 Mark every completed task `[X]`, rerun analysis, and commit locally as `feat(020): add GUI dual-syntax scheduling` with the required co-author trailer

## Dependencies and Execution Order

- Baseline precedes all behavioral changes.
- US1 establishes the shared GUI boundary and is the MVP.
- US2 consumes the US1 boundary for round-trip editing.
- US3 follows the stable behavior and stays independently reviewable as guidance.
- Final analysis and verification follow all stories.

## Implementation Strategy

1. Capture the green baseline.
2. Add US1 tests red, implement the smallest shared editor path, and turn green.
3. Add US2 compatibility tests red, extend the private form path, and turn green.
4. Synchronize GUI-local guidance.
5. Analyze, verify all eight gates, commit locally, and halt before push.
