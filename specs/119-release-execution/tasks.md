# Tasks: S119 release execution

**Input**: spec.md, plan.md, research.md, data-model.md, contracts/release-discovery.md and quickstart.md.

## Phase 1: Setup and foundation

- [x] T001 Align .specify/feature.json, CLAUDE.md and specs/README.md with the S119 specification/design/checklists. (FR-009)
- [x] T002 Complete blocking read-only analysis and both checklists before release/source implementation. (FR-009)
- [x] T003 Audit approved source, exact-main acceptance, original report/content hashes and read-only release guards; preserve receipts under .test-temp/s119/. (FR-001, FR-003, FR-007)

## Phase 2: US1 - Exact candidate records

**Independent validation**: Actual approved tag/source, terminal staging run, successful Windows attempt, draft and eight original assets agree.

- [x] T004 [US1] Write actual candidate/run/attempt/asset identities and owner authorities in specs/119-release-execution/delivery.md, with original private manifest/hash receipts in .test-temp/s119/. (FR-001, FR-002, FR-003)

## Phase 3: US2 - Public discovery source

**Independent validation**: Existing responsive browser assertions match actual prepared comparison selection and canonical current status reflects the real published release.

- [x] T005 [US2] Reproduce the stale site/e2e/site.spec.ts comparison-version assertion on an owned v1.6.0 build, retain before evidence and file/reuse its attributable GitHub defect under #354 if demonstrated. (FR-008)
- [x] T006 [US2] Correct only the demonstrated release-selection assertion in site/e2e/site.spec.ts and obtain passing existing browser acceptance without weakening exact visibility/table assertions. (FR-006, FR-008)
- [x] T007 [US2] After actual authorized promotion, reconcile site/data/published.json, README.md and current publication-status paragraphs in docs/{INSTALL-windows,INSTALL-linux,INSTALL-macos,mcp,notifications,cli-qa,remote-access,test-scripts}.md plus current-status/datestamped decision under CHANGELOG.md Unreleased, correct its current v1.6 introduction and refresh only site/data/reviews/v1.6.0.json changelog digest/date; preserve immutable tagged notes/dataset/original receipt and upgrade/rollback contracts. (FR-004, FR-005, FR-006)

## Phase 4: US3 - Traceability and source acceptance

- [x] T008 [US3] Record actual staging/promotion/publication and independent #354/#355/#356 outcomes in specs/119-release-execution/{delivery,verification}.md and existing GitHub records; preserve original failures and distinguish later source PR identity. (FR-003, FR-004, FR-005, FR-007, FR-008, FR-010)
- [x] T009 Run meaningful targeted browser acceptance, existing format/lifecycle/release/comparison guards and strict UTF-8/no-BOM/mojibake/hash audits; record actual results in specs/119-release-execution/verification.md. (FR-006, FR-007, FR-009)
- [x] T010 Freeze finished source and complete all eight foreground gates through scripts/verify.sh all, preserving source-input hashes and exact results under .test-temp/s119/. (FR-009)
- [x] T011 Finalize factual task/verification/Delivery and Implemented lifecycle in specs/119-release-execution/ and specs/README.md; recheck metadata before source commit. (FR-009)

## Dependencies and execution order

T001-T002 block implementation. T003 precedes the external approved tag/staging operation and its source receipt T004. T005 precedes T006; this local owned browser reproduction may proceed independently of staging. T007 requires actual public promotion following concrete-asset owner approval. T008 consolidates actual external receipts and independent issue outcomes. T009 precedes frozen T010; only factual lifecycle/task/verification finalization follows in T011. All eight gates precede any source commit.

## Parallel opportunities

Independent metadata/evidence reads may be batched. No [P] source tasks share concurrent mutation; browser builds, source identities and canonical checks require ordered execution.

## External delivery workflow

The owner explicitly approved immutable v1.6.0 tagging/draft staging at 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b. The owner separately approved the exact staged candidate on 2026-10-27; guarded promotion completed as run 40450478494. Release actions and approvals are external delivery receipts, not unchecked implementation tasks. Source push/one PR and maximum two Codex rounds are authorized; record exact-head CI and every review disposition externally and stop before owner merge. #354 closes only for actual delivered release/discovery; #355/#356 remain independently open without a demonstrated causal outcome.
