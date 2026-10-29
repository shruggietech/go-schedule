# Tasks: S126 Dependency and Release Convergence

**Input**: spec.md, plan.md, research.md, data model and delivery contract.

## Phase 1: Setup

- [x] T001 Synchronize intended repository and create codex/126-dependency-release-convergence with matching .specify/feature.json.
- [x] T002 Complete installed specify, clarify, checklist and plan artifacts in specs/126-dependency-release-convergence/.

## Phase 2: Foundations

- [x] T003 Freeze prior evidence and binary inputs into .test-temp/s126/historical-inputs.json.
- [x] T004 Run blocking cross-artifact analysis of spec.md, plan.md and tasks.md; retain its result in analysis.md.

## Phase 3: US1 Coherent selected dependencies

- [x] T005 [US1] Apply #401 selected runtime/YAML updates and regenerate root and desktop go.mod/go.sum (FR-001).
- [x] T006 [US1] Apply #402 Vite and #403 coupled site versions and regenerate desktop/frontend and site package-lock.json; run clean restores (FR-001).
- [x] T007 [US1] Record published-driver deferral for #400/#404 in research.md and native planning records; preserve SQLite/libc selection (FR-002/FR-003).
- [x] T008 [US1] Exercise both read-only Go graphs and existing API/YAML/storage contracts and record results in verification.md (FR-003/FR-004).

## Phase 4: US2 Truthful candidate readiness

- [x] T009 [US2] Reconcile S125 actual merge/main acceptance in specs/125-reliability-delivery-decision/verification.md, delivery.md and handoff.md (FR-005).
- [x] T010 [US2] Update CHANGELOG.md and .github/release-notes/v1.7.0.md for selected/deferred versions and refresh site/data/reviews/v1.7.0.json content hashes (FR-005/FR-006).
- [x] T011 [US2] Build and capture fresh Windows/WSL QA under .test-temp/s126 with unchanged controls and source/binary receipts (FR-004/FR-006).
- [x] T012 [US2] Review all frozen corpus messages in context, finalize/validate/render QA and preserve cleanup/integrity receipts in verification.md (FR-004/FR-006).
- [x] T013 [US2] Run full headless browser contracts, static-site routes/search and responsive dark/light brand checks; record in verification.md (FR-004).

## Phase 5: US3 Reviewable delivery

- [x] T014 [US3] Record exact same-slice release continuation and remaining approvals in handoff.md without invented asset identities (FR-008).
- [x] T015 [US3] Run complete foreground sh scripts/verify.sh all and source/comparison guards; preserve receipts and first failures in verification.md (FR-004/FR-007).
- [x] T016 [US3] Finalize specification lifecycle, specs/README.md entry and delivery.md, format prose and audit encoding/evidence integrity (FR-005/FR-007).

## Phase 6: Authorized publishing protocol

After implementation tasks, commit/push the branch and open one PR under the owner's explicit authorization. Monitor exact-head CI and all review comments. Request no more than one additional Codex round after the automatic initial round; resolve findings autonomously. Reconcile native records and hand off for owner merge without merging, tagging, staging or promotion. Later approved release actions remain same-slice continuation, not source tasks falsely marked delivered.

## Dependencies

T004 blocks implementation. T005-T007 precede candidate acceptance. T009-T010 precede frozen QA capture. T011 precedes T012; complete T012-T015 before factual lifecycle finalization T016. US3 release continuation can be documented independently, but publication waits for all source acceptance.
