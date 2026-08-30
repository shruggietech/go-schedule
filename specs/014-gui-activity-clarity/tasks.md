# Tasks: GUI Activity Clarity

**Input**: `specs/014-gui-activity-clarity/`

**Tests**: Required by FR-008 and constitution principle II. Production changes
follow observed failing focused tests.

**Task Reconciliation**: PR #46 merged after the recorded pre-push halt; the 2026-08-30 lifecycle audit resolves the stale publication marker.

## Phase 1: Specification and Design

- [x] T001 Create `specs/014-gui-activity-clarity/spec.md` with traceability to
  issues #26, #28, and #30.
- [x] T002 Complete the requirements-quality review in
  `specs/014-gui-activity-clarity/checklists/ux.md`.
- [x] T003 Record research, data model, implementation plan, and quickstart in
  `specs/014-gui-activity-clarity/`.
- [x] T004 Synchronize the active feature context in `CLAUDE.md`.

---

## Phase 2: User Story 1 and 2 - Activity Label and Bounded Badge (P1)

**Goal**: Accurately name the mixed view and keep its alert count compact.

**Independent Test**: Focused GUI tests cover the initial tab label and counts
0, 1, 99, 100, and above 100.

### Tests

- [x] T005 [US1] Update the tab-construction expectation from Logs to Activity
  in `gui/app_test.go` and `gui/logs_test.go`.
- [x] T006 [US2] Add table-driven badge-label boundary tests in
  `gui/app_test.go`, then run focused tests and record the expected failure.

### Implementation

- [x] T007 [US1] Use Activity as the initial user-facing tab label in
  `gui/app.go`.
- [x] T008 [US2] Add and apply the bounded Activity badge-label helper in
  `gui/app.go`.
- [x] T009 [US1] Run focused GUI tests and confirm both stories pass.

---

## Phase 3: User Story 3 - Non-Destructive Clear Presentation (P1)

**Goal**: Make the current-view action accurate and understandable without
hover while preserving its existing behavior.

**Independent Test**: A headless GUI test finds `Clear View`, the clear-content
icon, and explanatory copy; existing cutoff tests remain green.

### Tests

- [x] T010 [US3] Add failing control-presentation assertions to
  `gui/logs_test.go` before changing production code.

### Implementation

- [x] T011 [US3] Replace Dismiss All/delete presentation and add persistent
  explanatory copy in `gui/logs.go`.
- [x] T012 [US3] Update internal comments to describe Activity and clearing
  without renaming low-value identifiers in `gui/logs.go`.
- [x] T013 [US3] Run focused GUI tests and confirm presentation plus existing
  merge/filter/cutoff behavior passes.

---

## Phase 4: Documentation and Verification

- [x] T014 Add the Activity clarity change to `CHANGELOG.md`.
- [x] T015 Run `gofmt` on changed Go files and `go test ./gui/...`.
- [x] T016 Run `sh scripts/verify.sh all` in the foreground and confirm all
  eight gates pass.
- [x] T017 Audit the diff for scope, stale user-facing Logs/Dismiss All text,
  UTF-8 without BOM, and mojibake.
- [x] T018 Mark tasks and spec status complete, review the final diff, and
  commit with a truthful co-author trailer.
- [x] T019 Halt before pushing or opening the pull request and report that its
  body must contain `Closes #26`, `Closes #28`, and `Closes #30`.

## Dependencies and Execution Order

1. Phase 1 is complete before implementation analysis.
2. T005 and T006 precede T007 and T008 so the label/badge change is test-first.
3. T010 precedes T011 and T012 so the control change is test-first.
4. Focused tests pass before the full verification aggregate.
5. Publication waits at T019 for explicit operator authorization.
