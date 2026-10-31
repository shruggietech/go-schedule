# Tasks: S131 v2.0.0 Publication

**Input**: [spec.md](spec.md), [plan.md](plan.md), research.md, data-model.md and contracts/publication.md.

## Phase 1: Setup

- [x] T001 Initialize S131 feature/branch identity and installed Spec-Kit artifacts in specs/131-v2-publication/ and .specify/feature.json (FR-008).
- [x] T002 Preserve the 221-input inventory and private publication ledger location in .test-temp/s131/ (FR-005).

## Phase 2: Foundational

- [x] T003 Complete blocking cross-artifact analysis in specs/131-v2-publication/analysis.md before executing release work (FR-001 through FR-008).
- [x] T004 Recheck accepted main, exact-head CI/CodeQL/website, source/comparison guards and absent tag/draft; retain .test-temp/s131/source-selection.json (FR-001, FR-002).

## Phase 3: User Story 1 - Stage an attributable candidate

**Goal**: One complete draft on the exact approved source. **Independent test**: Inspect actual draft inventory without publishing it.

- [x] T005 [US1] Obtain exact-source tag/draft approval and record its scope in .test-temp/s131/publication-ledger.json (FR-001).
- [x] T006 [US1] Create/push the approved immutable tag through the existing .github/workflows/release.yml; retain actual run/attempt and all seven terminal jobs in specs/131-v2-publication/delivery.md (FR-001, FR-002).
- [x] T007 [US1] Download and inspect all eight actual draft assets, six archives and MSI/manifest identity into .test-temp/s131/candidate/; record exact inventory and native limitations in specs/131-v2-publication/verification.md (FR-002, FR-005, FR-006).

## Phase 4: User Story 2 - Publish the approved bytes

**Goal**: Actual public/latest v2.0.0 with immutable originals. **Independent test**: Resolve public assets and checksum against the approved inventory.

- [x] T008 [US2] Obtain separate concrete-asset promotion approval naming actual release/source/run/attempt/inventory in .test-temp/s131/publication-ledger.json (FR-003).
- [x] T009 [US2] Dispatch existing .github/workflows/promote-release-waived.yml and retain actual terminal outcome in specs/131-v2-publication/delivery.md (FR-003).
- [x] T010 [US2] Independently verify public/latest identities, unchanged eight originals and nine-file/eight-entry checksum set; preserve receipts in .test-temp/s131/publication-ledger.json (FR-003, FR-005).

## Phase 5: User Story 3 - Discover the published baseline

**Goal**: Actual public selection and current source guidance agree. **Independent test**: Read public routes and inspect reconciliation diff.

- [x] T011 [US3] Verify actual Pages source/deployment and HTTPS home, comparison, installation and representative downloads in specs/131-v2-publication/verification.md (FR-004).
- [x] T012 [US3] Reconcile site/data/published.json, README.md and affected current installation wording only after actual publication; record any pinned-document decision in CHANGELOG.md (FR-004, FR-006).
- [x] T013 [US3] Reconcile #412/milestone 19 with factual delivery and remaining final source integration; keep unrelated incidents/deferrals open, recorded in specs/131-v2-publication/delivery.md (FR-007).

## Phase 6: Final verification and official PR

- [x] T014 Run sh scripts/verify.sh all in the foreground, targeted public-selection/browser cases and source/frozen-input/encoding/format/lifecycle audits; retain source-bound results in specs/131-v2-publication/verification.md and .test-temp/s131/ (FR-005, FR-006, FR-008).
- [x] T015 Finalize specs/131-v2-publication/{spec,tasks,analysis,delivery,verification}.md, specs/README.md and active CLAUDE.md context with truthful functional/evidence disposition (FR-004, FR-007, FR-008).
## Authorized publication workflow (not implementation tasks)

W001 (FR-008): Commit, automatically push codex/131-v2-publication and open the official PR targeting main; retain URL/head in .test-temp/s131/publication-ledger.json (FR-008).

W002 (FR-008): Inspect all CI/review comments/threads and PR-body reactions, resolve warranted findings with at most one manual second review request, then report final ready head/check/review receipts in the private .test-temp/s131/publication-ledger.json and PR description/comment, with committed delivery.md retaining stable links and owner-merge disposition (FR-008).

## Dependencies and execution order

T001-T004 precede all release mutations. T005 approval precedes T006-T007. T008 concrete-candidate approval precedes T009-T010. T011-T013 follow actual publication, then final source verification and official PR. Story 1 draft inspection is useful independently; Story 2 requires that inventory, and Story 3 requires actual public delivery. S131 is not complete after planning or a draft alone. Final source integration is an owner boundary.

## Parallel opportunities

Read-only source/CI/tag/release checks may run concurrently. Within US1, archive inventories may be read concurrently after the downloads are frozen. Within US2, independent digest readback may run concurrently. Within US3, public route readbacks may run concurrently. Mutations, approvals, candidate selection, spec updates, commits and review requests remain sequential; no parallel tasks modify shared files.

## Implementation strategy

Execute the existing release process through its explicit approvals, preserving first failures and exact identities. Do not invent filler source code or repeat unchanged historical QA. Run existing security/identity negative tests and all canonical gates on final reconciliation source. Automatic PR-open review consumes slot one; only one manual recovery/final request remains. Any subsequent correction is separately attributed if no review slot remains.
