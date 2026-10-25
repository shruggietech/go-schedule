# Tasks: S116 CI timeout reliability

**Input**: spec.md, plan.md, research.md, data-model.md and contracts/timeouts.md.

**Delivery**: Implemented on review branch codex/116-ci-timeout-reliability with all eight canonical gates passed. Publication/CI/review workflow is tracked separately from implementation checkboxes.

## Phase1: Setup

- [x] T001 Create S116 branch, selector and design artifacts in specs/116-ci-timeout-reliability/.

## Phase2: Foundation

- [x] T002 Complete blocking cross-artifact analysis of specs/116-ci-timeout-reliability/spec.md, plan.md and tasks.md before source implementation.

## Phase3: US1 increased verification completion

- [x] T003 [US1] Add omitted/reduced policy mutation tests in scripts/ci-timeout-check/main_test.go and retain an actual old 5s frontend watchdog failure experiment under .test-temp/s116/.
- [x] T004 [US1] Implement all 22 workflow job caps and Go 60m parity in .github/workflows/, scripts/verify.sh, scripts/coverage-gate.sh, Makefile, scripts/ci-timeout-check/ and automation policy fixtures.
- [x] T005 [US1] Increase frontend/browser/site runner limits in desktop/frontend/vite.config.ts, src/test-setup.ts, both playwright.config.ts files, site/scripts/browser-test.mjs and site/next.config.mjs.
- [x] T006 [US1] Increase audited preparation/service/installer/guest/promotion completion allowances in test/integration/, test/windows/, scripts/windows-qualification-session/ and scripts/release-promotion/.
- [x] T007 [US1] Document the full numeric/classified policy and pinned-workflow decision in docs/ci-timeouts.md, docs/build-autopilot.md and CHANGELOG.md.

## Phase4: US2 exact search identity

- [x] T008 [US2] Retain and strengthen exact selected-profile and no-wrong-context assertions in desktop/frontend/src/App.test.tsx, removing its isolated 5s lookup override.
- [x] T009 [US2] Demonstrate delayed frontend completion beyond the old whole-test limit and run actual identity/search suites; record results in specs/116-ci-timeout-reliability/verification.md.

## Phase5: US3 native budget and attribution

- [x] T010 [US3] Add new-default, unchanged explicit/resumed-budget and aligned cleanup reserve regressions in internal/clievaluate/timing_policy_test.go and environment_test.go before changes.
- [x] T011 [US3] Implement 120s/2h QA defaults, aligned completion reserves and wrapper preparation allowances in internal/clievaluate/ and scripts/cli-qa.ps1; reconcile docs/cli-qa.md.
- [x] T012 [US3] Add actual CLI/owned API metadata cases with controlled beyond-old-limit and explicit-short-budget evidence in internal/clievaluate/metadata_deadline_test.go.
- [x] T013 [US3] Run one Windows/Linux diagnostic six-writer and memory/WAL/metadata follow-up, preserve exact source/runtime/binary/budget/cleanup identities in specs/116-ci-timeout-reliability/verification.md.

## Phase6: US4 release source handoff

- [x] T014 [US4] Reconcile actual S115 merged-main failure, known issues and later release identities in specs/116-ci-timeout-reliability/handoff.md.
- [x] T015 [US4] Include delivered S115/S116 candidate scope in CHANGELOG.md and .github/release-notes/v1.6.0.md, refresh site/data/reviews/v1.6.0.json governed digests and retain published v1.5.1.

## Phase7: Cross-cutting completion

- [x] T016 Run all eight canonical gates through scripts/verify.sh all in the foreground, plus formatting/encoding checks, and record exact platform/results in specs/116-ci-timeout-reliability/verification.md.
- [x] T017 Complete spec/tasks lifecycle and source handoff in specs/116-ci-timeout-reliability/ and specs/README.md; commit on the review branch.

## Dependencies and implementation strategy

T001 -> T002 blocks all source changes. T003 precedes T004-T007 and includes the actual old frontend watchdog experiment before T005. US2 T008/T009 follows the shared runner policy and preserves production identity code. US3 T010 precedes T011/T012; T013 follows all native source changes. US4 T014/T015 follows actual current observations. T016 follows all functional code/content; T017 follows canonical acceptance. The full owner-authorized bundle is delivered together, not split into one issue per slice.

## Parallel opportunities

Three read-only research agents completed disjoint CI/native/search-release inventories under the installed plan command. Source implementation is coordinated sequentially to avoid overlapping policy/fixture changes; native experiments are serialized for attributable resource use.

## Requirement coverage

| Requirement | Tasks |
| --- | --- |
| FR-001 | T003, T004, T005, T006, T007 |
| FR-002 | T003, T004, T005, T006 |
| FR-003 | T004, T005, T016 |
| FR-004 | T003, T008, T010, T012, T013, T016 |
| FR-005 | T008, T009 |
| FR-006 | T010, T011, T012, T013 |
| FR-007 | T010, T012, T013 |
| FR-008 | T014, T015 |
| FR-009 | T001, T007, T014, T017 |
| FR-010 | T017 and publication workflow |

## Publication workflow

Owner explicitly authorized branch push, one official PR, exact-head CI monitoring and at most two Codex rounds. Resolve all warranted comments, record individual source/CI/review identities and issue dispositions, then stop before merge, tag, staging or promotion. #354 stays open for public release; original #355/#356 remain open unless their own functional outcomes are actually delivered. No extra approval pause is needed for this already-authorized PR publication.
