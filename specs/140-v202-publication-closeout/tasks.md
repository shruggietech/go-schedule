# Tasks: v2.0.2 publication closeout

**Input**: spec.md, plan.md, research.md, data-model.md, contracts/release-checkpoint.md and quickstart.md.

## Setup and foundation

- [X] T001 Verify existing ignore rules, hidden foreground launchers and immutable history/protected-evidence manifests in .gitignore and .test-temp/s140/ (FR-005, FR-008).
- [X] T002 Run the installed read-only blocking analysis over specs/140-v202-publication-closeout/ and record full requirement/story coverage before source implementation (FR-001 through FR-009).

## User Story 1 - Accurate maintenance comparison (P1)

**Independent observation**: all eight contexts, ten stable rows and 80 cells describe the selected products with honest dates and uncertainty.

- [X] T003 [US1] Prepare site/data/comparisons/v2.0.2.json from the reviewed 80-cell research, preserving identities/opportunities and five carried evidence dates (FR-001, FR-002, FR-003).
- [X] T004 [US1] Remove only the obsolete unreleased cron notice in docs/cron.md for #453, preserving documented dialect/fidelity (FR-009).
- [X] T005 [US1] Describe the comparison preparation and #453 correction in CHANGELOG.md and .github/release-notes/v2.0.2.md without claiming publication; consolidate related highlights within the existing six-bullet limit for #454 (FR-002, FR-006, FR-009).

## User Story 2 - Exact content-bound review (P1)

**Independent observation**: the existing guard rejects an unfinished review, then accepts ten completed decisions bound to all three final contents, while historical bytes are unchanged.

- [X] T006 [US2] Run site/scripts/review-release.mjs for v2.0.2, retain the generated placeholder and its expected rejection privately, then complete site/data/reviews/v2.0.2.json after individual row review (FR-004).
- [X] T007 [US2] Validate candidate binding, historical comparisons and release-source metadata with scripts/comparison-check and scripts/release-source-check (FR-004, FR-005).

## User Story 3 - Honest continuation (P2)

**Independent observation**: the continuation states every release prerequisite/outcome and current public selection remains v2.0.1.

- [X] T008 [US3] Finalize specs/140-v202-publication-closeout/contracts/release-checkpoint.md and delivery.md with the ordered remaining same-S140 publication states and issue outcomes (FR-006, FR-007, FR-009).
- [X] T009 [US3] Run existing production site/e2e contracts against the candidate and restored default output with distinct private headless-browser artifacts; confirm site/data/published.json remains v2.0.1 (FR-001, FR-006, FR-008).

## Integration and final records

- [X] T010 Run all eight scripts/verify.sh gates in the foreground and retain source/evidence integrity results under .test-temp/s140/ (FR-005, FR-008).
- [X] T011 Finalize specs/140-v202-publication-closeout/verification.md and factual lifecycle/task records; run GitHub-format, lifecycle, encoding, diff and immutable-history checks (FR-001 through FR-009).

## Dependencies and implementation strategy

Execute T001-T002 before source implementation. T003-T005 precede T006; T006 precedes T007. T008 defines the continuation without executing publication. T009 follows final dataset/review preparation; T010 follows all implementation. T011 follows actual successful observations. No new copy-mirroring tests are needed for declarative/prose changes; the existing structural, binding, website and canonical guards exercise their real contracts.

The official push, PR monitoring and owner handoff follow completed implementation records. Tagging, draft staging, inspected-asset promotion and the later discovery PR are subsequent same-S140 continuation states, not unchecked pre-tag source tasks. No independent tasks require further agent delegation.

## Post-publication continuation (2026-11-09)

- [X] T012 [US3] Promote the approved immutable inventory through the existing workflow and independently verify nine public assets/eight checksum entries under .test-temp/s140/promotion-01/ (FR-005, FR-007, FR-008, FR-010).
- [X] T013 [US3] Extend site/e2e/site.spec.ts with a failing-first published-baseline/dated-context regression for #459 (FR-011).
- [X] T014 [US3] Reconcile current README.md, docs/install.md, docs/INSTALL-windows.md, docs/gui-fields.md, docs/first-task.md, docs/website.md and site/data/published.json; clarify site/app/compare/page.tsx and update the first-task browser expectation (FR-010, FR-011).
- [X] T015 [US3] Record actual publication/inventory/known-defect limits in specs/140-v202-publication-closeout/ and verify all twelve production browser contracts without changing historical bindings (FR-005, FR-007, FR-008, FR-009, FR-010, FR-011).
- [X] T016 Run all eight canonical foreground gates for the final discovery source and retain source/protected-evidence integrity (FR-005, FR-008).
- [X] T017 Finalize lifecycle/tasks/specification inventory and run formatting, encoding, comparison/release-source and immutable-artifact checks before the already authorized official PR (FR-004, FR-005, FR-007, FR-008, FR-009, FR-010, FR-011).

T012 is the prerequisite for current availability claims. T013 fails against the existing page before T014. T014 precedes browser validation and T015/T016; factual T017 records follow actual observations. The pre-tag tasks remain completed history. The final owner merge/deployment is a delivery boundary, not an unchecked source task.
