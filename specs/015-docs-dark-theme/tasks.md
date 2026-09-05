# Tasks: Documentation Dark-Theme Quality

**Input**: Design documents in `specs/015-docs-dark-theme/`

**Tests**: Required by the specification and constitution principle II. Each behavioral phase adds its offline regression assertion before the corresponding SCSS or Markdown correction.

**Task Reconciliation**: PR #47 merged after the recorded pre-push halt; the 2026-08-30 lifecycle audit resolves the stale publication marker.

## Phase 1: Specification and Design

- [x] T001 Create `specs/015-docs-dark-theme/spec.md` with traceability to issues #36, #37, and #35.
- [x] T002 Complete `specs/015-docs-dark-theme/checklists/requirements.md` and `specs/015-docs-dark-theme/checklists/accessibility.md`.
- [x] T003 Record research, data model, theme contract, plan, and quickstart in `specs/015-docs-dark-theme/`.
- [x] T004 Synchronize the active feature context in `CLAUDE.md`.

---

## Phase 2: User Story 1 - Read Every Code Token (Priority: P1)

**Goal**: Make every current and future syntax token readable on the dark code surface, including selected text and highlighted lines.

**Independent Test**: The offline theme contract rejects the current incomplete palette, then passes after safe fallback and role mappings are complete.

### Tests

- [x] T005 [US1] Add safe fallback, named role, highlighted-line, and selection contract assertions to `scripts/docs-check.sh`.
- [x] T006 [US1] Run `sh scripts/docs-check.sh` and record the expected failure against the current incomplete theme.

### Implementation

- [x] T007 [US1] Implement safe-default and complete dark token styling in `docs/_sass/custom/custom.scss`.
- [x] T008 [US1] Rerun the focused docs gate and confirm the syntax theme contract passes.

---

## Phase 3: User Story 2 - Encounter Consistent Code Examples (Priority: P2)

**Goal**: Require every published code block to declare one approved, accurate content category.

**Independent Test**: The fence scanner rejects the existing untagged example, then all 11 pages pass after it is classified as text.

### Tests

- [x] T009 [US2] Add actionable fence-vocabulary validation for published pages to `scripts/docs-check.sh`.
- [x] T010 [US2] Run `sh scripts/docs-check.sh` and record the expected untagged fence failure in `docs/gui-fields.md`.

### Implementation

- [x] T011 [US2] Classify the non-executable example fence in `docs/gui-fields.md` as `text` without changing its content.
- [x] T012 [US2] Rerun the focused docs gate and confirm every published fence uses `sh`, `bash`, `powershell`, or `text`.

---

## Phase 4: User Story 3 - Read a Comfortably Spaced Endorsement (Priority: P2)

**Goal**: Align the sidebar endorsement with navigation content on small and desktop layouts.

**Independent Test**: The offline style contract requires responsive theme gutter variables and non-zero bottom spacing, then passes after implementation.

### Tests

- [x] T013 [US3] Add responsive endorsement-spacing assertions to `scripts/docs-check.sh` and record the expected failure.

### Implementation

- [x] T014 [US3] Apply small and desktop navigation gutters plus vertical spacing in `docs/_sass/custom/custom.scss`.
- [x] T015 [US3] Rerun the focused docs gate and confirm the complete docs contract passes.

---

## Phase 5: Documentation, Verification, and Handoff

- [x] T016 Add the reader-facing improvement and dated pinned docs-gate decision to `CHANGELOG.md`.
- [x] T017 Run `sh scripts/verify.sh all` in the foreground and confirm all eight gates pass.
- [x] T018 Audit scope, exact contrast pairs, fence inventory, UTF-8 without BOM, mojibake, and `git diff --check`.
- [x] T019 Mark tasks and spec status complete, review the final diff, and commit with a truthful co-author trailer.
- [x] T020 Halt before pushing or opening the pull request and report that its body must contain `Closes #36`, `Closes #37`, and `Closes #35`.

## Dependencies and Execution Order

1. Phase 1 completes before the read-only analysis gate.
2. T005 and T006 precede T007 so the accessibility defect is observed red.
3. T009 and T010 precede T011 so fence normalization is observed red.
4. T013 precedes T014 so the spacing defect is observed red.
5. Focused docs validation passes before the full repository aggregate.
6. Publication waits at T020 for explicit operator authorization.

## Parallel Opportunities

The three user stories touch the same SCSS and validation files, so they execute chronologically rather than in parallel. This keeps each red-green signal attributable to one behavior and avoids file conflicts.
