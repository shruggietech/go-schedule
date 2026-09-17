# Tasks: Final v1.4.0 boundary and qualification

**Input**: spec.md, plan.md, research.md, data-model.md, contracts/release-disposition.md and quickstart.md.

## Phase 1: Setup

- [x] T001 Verify repository, main and draft identities in research.md (FR-003, FR-004).
- [x] T002 Record capabilities, Sandbox-only authority and testing waiver in spec.md (FR-002, FR-010).

## Phase 2: Foundation

- [x] T003 Analyze spec.md, plan.md and tasks.md; resolve stale all-environments blocking language (FR-007, FR-009).

## Phase 3: User Story 1, Coherent Boundary (P1)

Independent test: existing metadata fixtures and direct cumulative-content audit agree with reviewed source.

- [x] T004 [US1] Reconcile desktop entries and chronological decisions in CHANGELOG.md without claiming native completion (FR-001).
- [x] T005 [US1] Update .github/release-notes/v1.4.0.md highlights and testing-waiver disclosure, retaining its anchor (FR-001).
- [x] T006 [US1] Document Sandbox and waiver disposition in test/windows/README.md without changing passing-evidence semantics (FR-002, FR-007, FR-010).

## Phase 4: User Story 2, Exact Candidate (P1)

Independent test: final assets agree; supported observations retain actual results and unsupported checks remain untested.

- [ ] T007 [US2] After preparation merge and exact-main CI, preserve and refresh the draft under release-operation authority using .github/workflows/release.yml (FR-003, FR-004).
- [ ] T008 [US2] Execute supported Sandbox fresh/upgrade, desktop and task observations using test/windows/README.md; record waived gaps in qualification-status.md (FR-005, FR-006, FR-007).

## Phase 5: User Story 3, Review Handoff (P2)

Independent test: stored body matches Markdown, findings are disposed, latest-head CI is green.

- [x] T009 [US3] Run scripts/verify.sh all and record actual results in verification.md (FR-008).
- [ ] T010 [US3] Publish from .git/s088-pr-body.md, verify stored body and preserve open issue criteria (FR-008, FR-009).
- [ ] T011 [US3] Handle all findings within two rounds and verify latest-head CI before merge handoff in verification.md (FR-008, FR-009).

## Phase 6: Cross-Cutting

- [x] T012 Update CLAUDE.md plan context and validate quickstart.md and publication formatting (FR-001, FR-008).

## Dependencies and Strategy

T001-T003 precede implementation. T004-T006 form the preparation increment; T009-T012 validate and publish it. T007-T008 require maintainer merge, not merely publication; leave them unchecked at that handoff. Do not declare full S088 qualification from the preparation PR.

Parallel opportunities: after analysis, cumulative-content and runbook audits affect distinct files; release-note review joins their findings. Execution is sequential locally. MVP: coherent release metadata, followed by reviewed staging and supported native observations.
