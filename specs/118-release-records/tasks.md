# Tasks: S118 release record reconciliation

**Input**: spec.md, plan.md, research.md, data-model.md, contracts/release-records.md and quickstart.md.

## Phase 1: Setup and foundation

- [x] T001 Align branch, feature selector, CLAUDE.md plan pointer and specs/README.md inventory for S118; complete required planning artifacts. (FR-009)
- [x] T002 Complete blocking read-only analysis and both requirements checklists in specs/118-release-records/ before source implementation. (FR-009)
- [x] T003 Audit immutable retained S115/S117 report and pending-content hashes; preserve identity receipts in .test-temp/s118/. (FR-003, FR-006, FR-008)

## Phase 2: US1 - Completed implementation records

**Independent validation**: S117 maturity, task dispositions, inventory and actual PR #377 receipts agree without altering original candidate outcomes.

- [x] T004 [US1] Reconcile specs/117-candidate-qa-handoff/spec.md and tasks.md; retain T016/T017 as external receipts and explicitly attribute T018 reconciliation to S118. (FR-001, FR-002)
- [x] T005 [US1] Reconcile specs/117-candidate-qa-handoff/delivery.md and append chronological actual delivery evidence to verification.md; update specs/README.md. (FR-001, FR-003)

## Phase 3: US2 - Current release handoff

**Independent validation**: Current source and issue records agree on delivered scope, remaining causes, pending/public versions and later publication identities.

- [x] T006 [US2] Write specs/118-release-records/handoff.md with individual #354/#355/#356 dispositions and exact-source staging/promotion sequence; preserve unchanged pending content and public selection. (FR-003, FR-004, FR-005, FR-006, FR-007, FR-008)
- [x] T007 [US2] Update test/windows/README.md and docs/INSTALL-windows.md current procedure links and add the dated pinned-file decision plus Unreleased record correction in CHANGELOG.md. (FR-006, FR-007)
- [x] T008 [US2] Record corresponding current dispositions on existing GitHub #354/#355/#356 and in specs/118-release-records/delivery.md without closing unresolved outcomes. (FR-004, FR-005, FR-010)

## Phase 4: Source validation and finalization

- [x] T009 Run existing lifecycle, release-source, comparison and GitHub-format guards; audit links, UTF-8/no-BOM/mojibake and retained hashes; record actual results in specs/118-release-records/verification.md. (FR-003, FR-006, FR-008, FR-009)
- [x] T010 Freeze completed source inputs and run all eight foreground gates through scripts/verify.sh all, retaining the exact outcome and no-drift receipt under .test-temp/s118/. (FR-009)
- [x] T011 Finalize factual verification, task dispositions and Implemented lifecycle in specs/118-release-records/{spec,tasks,verification}.md and specs/README.md; recheck formatting and immutable inputs before commit. (FR-001, FR-009)

## Dependencies and execution order

T001-T002 precede implementation. T003 establishes identity preservation before T004-T008. US1 T004-T005 precedes the current handoff T006-T008. T009 precedes frozen T010, then only factual verification/task/lifecycle finalization in T011. No source edits overlap frozen validation. Every commit requires the canonical eight gates.

## Parallel opportunities

Independent record and evidence reads may be batched. No [P] implementation tasks: shared lifecycle, issue and source identities require ordered edits. No new runtime tests are warranted for this records-only change; existing validators and canonical checks exercise the relevant contracts.

## Delivery workflow

Push, one official PR, exact-head CI and individual feedback resolution are authorized external workflow, recorded in Delivery and the PR rather than implementation checkboxes. Count automatic initial Codex review plus at most one requested second round. Stop before merge, tagging, staging or promotion. Every included issue retains its individual functional disposition.
