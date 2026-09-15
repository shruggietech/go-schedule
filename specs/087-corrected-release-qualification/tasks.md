# Tasks: Corrected v1.4.0 qualification

**Input**: [plan.md](plan.md), [spec.md](spec.md), research, data model, and qualification contract.

## Phase 1: Setup

- [x] T001 Record scope and authorization in specs/087-corrected-release-qualification/spec.md and checklists/.
- [x] T002 Research exact reviewed tag refresh and native resources in specs/087-corrected-release-qualification/research.md.

## Phase 2: Foundational

- [x] T003 Preserve old draft metadata, annotated tag and all eight assets in dist/s087/obsolete-candidate/ before refreshing.

## Phase 3: User Story 1 - Corrected candidate

**Independent test**: Candidate identities and complete hosted staging agree.

- [x] T004 [US1] Replace the reviewed tag with an exact old-reference lease and record staging identity in specs/087-corrected-release-qualification/verification.md (FR-001, FR-002).
- [x] T005 [US1] Verify the completed staging asset set and downloaded candidate through scripts/windows-release-gate; retain files in dist/s087/candidate/ (FR-001).

## Phase 4: User Story 2 - Native observations

**Independent test**: Full fresh/upgrade and corrected-UI evidence validates through the unchanged gate.

- [x] T006 [US2] Acquire official offline prerequisites and build independent sessions with scripts/windows-qualification-session into dist/s087/sessions/ (FR-005).
- [ ] T007 [US2] Execute fresh-install and upgrade qualification from test/windows/README.md, preserving exact-candidate observations and exports in dist/s087/sessions/ (FR-003, FR-004).
- [ ] T008 [US2] Complete separate required native user/profile/display environments and finalize evidence through test/windows/Invoke-ReleaseCandidateAttended.ps1 (FR-003, FR-004, FR-006).

## Phase 5: User Story 3 - Honest handoff

**Independent test**: Tracked claims accurately describe validated evidence and excluded public promotion.

- [x] T009 [US3] Reconcile issue #226 and #228 status with actual release evidence and document results in specs/087-corrected-release-qualification/verification.md (FR-007, FR-009).
- [x] T010 [US3] Document review scope, exact source and remaining promotion decision in specs/087-corrected-release-qualification/verification.md (FR-008).

## Phase 6: Verification

- [x] T012 Reproduce the native selector-height finding from #231 in desktop/frontend/e2e/tasks.spec.ts before changing shared CSS (FR-010).
- [x] T013 Correct intrinsic field alignment in desktop/frontend/src/styles.css, run browser regressions and full foreground verification, and document the still-blocked release qualification (FR-010).

- [x] T011 Update specs/README.md, CHANGELOG.md and slice delivery fields with actual completed outcomes, then run full foreground scripts/verify.sh all before commit.

## Dependencies and execution

T003 precedes T004. T004 and T005 precede all candidate observations. T006 precedes T007; T007 and T008 precede qualification-complete claims in T009 and T010. T011 follows completed evidence. Backup, failed runs, and unavailable rows are retained rather than counted as passes.

Portable prerequisite retrieval can run alongside hosted staging, but no native qualification can begin before exact bytes are verified. No parallel agent implementation is required. MVP is the verified refreshed draft (US1); it is explicitly not complete native qualification. Push, PR, review and merge are publication workflow, not implementation tasks. At most two AI review rounds are permitted. Public promotion is not a task in this slice.

Qualification finding disposition: T012 and T013 were added after the fresh walkthrough reproduced the #231 selector defect. T007/T008 remain blocked pending reviewed repaired bytes; they are not silently removed or marked complete. T009/T010 now document the actual blocker rather than qualification completion. The bounded repair/evidence PR precedes restaging because the release workflow requires reviewed main. S087 remains In Progress until its native requirements are genuinely complete.
