# Tasks: S114 release readiness

**Input**: [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md), contracts and quickstart.

## Phase 1: Setup

- [x] T001 Verify synchronized main, issue authority and S114 review branch in specs/114-release-readiness/spec.md.
- [x] T002 Complete installed specify, clarify, checklist and plan commands in specs/114-release-readiness/.
- [x] T003 Complete blocking analyze coverage review of specs/114-release-readiness/{spec,plan,tasks}.md before source changes.

## Phase 2: Shared evidence

- [x] T004 Preserve exact original observations and identify native Windows/Linux diagnostic inputs in specs/114-release-readiness/verification.md.
- [x] T005 Record S114 kickoff and any independently demonstrated new bug under #354 using formatted GitHub records.

## Phase 3: US1 payload diagnosis

Goal: retain independent native failure outcomes and demonstrate any new bounded diagnostic defects without relabelling the historical cause.

Independent criteria: six simultaneous writers retain exact records; candidate refusal/crash/launcher output and stages remain attributable.

- [x] T006 [US1] Add failing-wrapper diagnostic regression to test/integration/testscripts_probe_test.go before changing test/scripts/lib/sqlite.sh.
- [x] T007 [US1] Preserve capped sanitized wrapper stderr separately from candidate stderr/native status in test/scripts/lib/sqlite.sh.
- [x] T008 [US1] Add child-owned phase markers and bounded all-six outcome/row attribution on failure in test/integration/testscripts_binding_test.go.
- [x] T009 [US1] Investigate inherited-pipe deadline behavior with an owned synthetic candidate in .test-temp/s114; correct only a deterministic demonstrated gap with regression in test/integration/testscripts_probe_test.go and test/scripts/lib/Sqlite.ps1.
- [x] T010 [US1] Run current six-writer twins and candidate fault matrix on native Windows and native Linux, preserving original limits and causal disposition in specs/114-release-readiness/verification.md.

## Phase 4: US2 owned IPC diagnosis

Goal: attribute sustained/canceled requests over actual IPC while preserving command and partial observation outcomes.

Independent criteria: canceled request returns and is attributed, complete-state assertions fail for incomplete inspection, next fresh request succeeds and mutations occur exactly once.

- [x] T011 [US2] Exercise held-request cancellation and mixed typed API operations over a real owned IPC server/store in internal/clievaluate/ipc_diagnostics_test.go.
- [x] T012 [US2] Compare current/prior dependency diagnostics under existing deadlines on native Windows and Linux using .test-temp/s114 inputs.
- [x] T013 [US2] Record demonstrated or unconfirmed original IPC cause and independent bug outcome in specs/114-release-readiness/verification.md without speculative transport or storage changes.

## Phase 5: US3 release preparation

Goal: deliver coherent pending v1.6.0 notes, upgrade/rollback guidance and later exact-candidate publication instructions.

Independent criteria: actual prepared notes pass production read-only selection; metadata agrees and public latest remains v1.5.1.

- [x] T014 [US3] Add failing actual-note production resolver regression in scripts/release-promotion/promotion_test.go with zero mutation assertions.
- [x] T015 [US3] Correct standing-waiver disclosure and concise unresolved #355/#356 guidance in .github/release-notes/v1.6.0.md.
- [x] T016 [US3] Reconcile schema-v24 stopped-daemon/full-state backup and matching rollback guidance in docs/INSTALL-linux.md and docs/INSTALL-macos.md.
- [x] T017 [US3] Add exact later publication handoff in specs/114-release-readiness/delivery.md and update active links in docs/INSTALL-windows.md and test/windows/README.md.
- [x] T018 [US3] Update CHANGELOG.md, pending comparison review digests and specs/README.md while preserving site/data/published.json and historical v1.5.1 bytes.

## Phase 6: Verification

- [x] T019 Run foreground canonical eight gates, relevant release/comparison checks and encoding/format audit; record exact outcomes in specs/114-release-readiness/verification.md.

## Dependencies and parallel opportunities

T001-T003 precede all implementation. T004-T005 provide shared identity and issue tracking. US1 T006-T010 and US2 T011-T013 have independent files and can run independently after setup; US3 T014-T017 also has independent files. T018 follows actual source/notes changes so digests are final. T019 follows all applicable implementation and issue dispositions. Publication workflow is recorded separately below, outside implementation checkboxes. Within a story the failing regression precedes its correction.

## Implementation strategy

Implement demonstrated regressions first, preserve unknown causes visibly, then reconcile release source and exact content receipts. S114's source work ends at the owner review boundary; actual #354 public delivery remains pending. No test pass, archive or review receipt closes a causally unresolved bug.

## Publication history

Publication, review, merge and branch cleanup are workflow evidence, not implementation tasks (docs/build-autopilot.md). The original T020/T021 entries were incorrectly added as implementation checkboxes and are retained here as historical references. T020 records the authorized branch push, official PR #365 and attachment. T021 records all 22 terminal successful checks on d11646111be52851c5b32a468050b508c63954c4, the single automatic Codex round with no findings on 15e2578554a5ebc1b4b9b0ea7eec2f2113620575, zero review threads, and the owner squash merge at 7bce610396bf83210e09376a2343d2e6ad21d180. No second review was requested.

The unmerged documentation checkpoint 574d071 is preserved through refs/archive/s114-unmerged-docs and reconciled into this record. Final main CI, CodeQL and website deployment are successful on the actual merge SHA; no v1.6.0 tag or public release is created by this housekeeping. #354 remains the release coordinator and #355/#356 remain unresolved investigations.
