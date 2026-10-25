# Tasks: S117 candidate QA and release handoff

**Input**: spec.md, plan.md, research.md, data-model.md, contracts/candidate-evidence.md and quickstart.md in this directory.

**Rules**: Every correction requires failing-before/passing-after regression evidence. All eight canonical gates precede any commit. No task-checkbox or source edits during frozen build/capture/review/finalization; record transient progress privately and update tasks afterward.

## Phase 1: Setup

- [x] T001 Verify branch, selector and issue scope, complete Spec-Kit artifacts and current plan pointer in `.specify/feature.json`, `CLAUDE.md` and `specs/117-candidate-qa-handoff/`. (FR-011)
- [x] T002 Complete blocking read-only analysis and checklist audit against `specs/117-candidate-qa-handoff/{spec,plan,tasks}.md`. (FR-011)

## Phase 2: Foundation

- [x] T003 Inventory existing headless tool/WSL prerequisites and ignore coverage; retain receipts under `.test-temp/s117/` without provisioning. (FR-003)
- [x] T004 Record S117 kickoff under #354, initialize private review-cap/publication ledger in `.test-temp/s117/publication-ledger.json` and add `specs/README.md` lifecycle entry before source freeze. (FR-008, FR-010, FR-012)

## Phase 3: US1 - Complete candidate evaluation (P1, first increment)

**Independent test**: The existing coordinator accounts for every applicable case/exclusion and actual reviewed message with validated canonical pass/fail/blocked reporting.

- [x] T005 [US1] Freeze source identity and build six immutable inputs plus receipt using `scripts/cli-qa.ps1` into `.test-temp/s117/qa-inputs-01`. (FR-001)
- [x] T006 [US1] Run one full capture with supported defaults using `scripts/cli-qa.ps1` into `.test-temp/s117/qa-run-01`; preserve all results and cleanup. (FR-002, FR-003, FR-005)
- [x] T007 [US1] Read every fresh corpus message and relevant case evidence; author exact-bound seven-rubric `.test-temp/s117/reviews-01.json` and classify execution failures. (FR-004, FR-005)
- [x] T008 [US1] Finalize, validate and render canonical reports with `scripts/cli-qa.ps1`, then preserve their hashes and true outcome in `specs/117-candidate-qa-handoff/verification.md`. (FR-001, FR-004, FR-005)

## Phase 4: US2 - Correct demonstrated defects (P1)

**Independent test**: Each demonstrated correction has an attributable failing-before/passing-after contract result; independent unresolved original causes and immutable capture remain visible.

- [x] T009 [US2] Inspect actual failures and perform only evidence-directed owned-fixture follow-ups in the demonstrated affected components (including `internal/store/`, `internal/api/server/`, `internal/mcpobserve/` and `internal/clievaluate/`); record exact findings in `specs/117-candidate-qa-handoff/verification.md`. (FR-005, FR-007)
- [x] T010 [US2] For each demonstrated component defect, add a failing regression beside its affected source in `internal/clievaluate/`, `internal/clifixture/` or `test/integration/`, retaining meaningful original assertions; otherwise explicitly record no demonstrated correction. (FR-006)
- [x] T011 [US2] Implement proportionate demonstrated fixes in those same affected components and obtain passing race-checked corrective acceptance, separately identified in `specs/117-candidate-qa-handoff/verification.md`. (FR-006)
- [x] T012 [US2] Reconcile #355/#356 individually and create/reuse distinct defect records under milestone 16; retain links, functional outcomes and unresolved causes in `specs/117-candidate-qa-handoff/delivery.md`. (FR-007, FR-008)

## Phase 5: US3 - Exact source release handoff (P2)

**Independent test**: The source record distinguishes actual candidate, final source, CI and review identities from later merge/tag/package/public release. Public version selection stays v1.5.1.

- [x] T013 [US3] Reconcile actual pending guidance in `.github/release-notes/v1.6.0.md`, `CHANGELOG.md` and `docs/INSTALL-windows.md` only if findings warrant changes, refreshing `site/data/reviews/v1.6.0.json` bindings when affected bytes change. (FR-009)
- [x] T014 [US3] Write exact candidate/final-source handoff and remaining owner actions in `specs/117-candidate-qa-handoff/{verification,delivery}.md`; validate source/comparison guards without tag or release operations. (FR-009, FR-010)

## Phase 6: Canonical validation

- [x] T015 Run UTF-8/mojibake and GitHub formatting checks plus all eight foreground gates through `scripts/verify.sh all`; record exact outcomes in `specs/117-candidate-qa-handoff/verification.md`. (FR-011)

## Dependencies and execution order

T001-T004 block US1. US1 T005-T008 is strictly sequential with source frozen. US2 follows the immutable original final report; T010 precedes any T011 behavioral correction. No demonstrated defect means an explicit no-correction disposition, not invented code churn. US3 follows actual dispositions. Canonical verification precedes every commit, publication and any in-scope review-fix push. Review count is automatic initial round plus at most one requested second round.

## Parallel opportunities

Within US1, reading independent retained case artifacts can be batched after capture; no source mutation may overlap. Within US2, independent evidence reads can be batched, while tests and affected source edits stay ordered. Within US3, independent release-content reads and identity audits can be batched. No implementation tasks are marked [P] because execution and source identity share state.

## Implementation strategy

Deliver US1 as the first meaningful increment, preserve it before US2 corrections, and then deliver US3 as one coherent source PR. A valid failed or blocked report is evidence to act on, not permission to weaken checks. No full rerun merely to seek a better outcome. #354 remains open for actual publication, and #355/#356 remain open unless their own corrective outcome is established.

## Historical workflow and reconciliation dispositions

Publication is external delivery evidence rather than implementation work, as required by docs/build-autopilot.md. The original task IDs remain traceable below without keeping completed publication steps as unchecked implementation tasks.

| Historical ID | Actual disposition |
| --- | --- |
| T016 | Completed through [PR #377](https://github.com/shruggietech/go-schedule/pull/377), reviewed head 7327a7c401b6140d0b6c02d2769ea3082d74a123; source PR published and later owner-merged |
| T017 | Exact-head checks passed and exactly two Codex rounds completed without findings; actual receipts are in [delivery.md](delivery.md) |
| T018 | S117's unmerged handoff and issue dispositions were delivered in PR #377, but source lifecycle bookkeeping was missed. S118 reconciles the spec, inventory and source delivery record after the owner merge; this is not attributed to the earlier S117 implementation session |

S117 implementation maturity is Implemented. Merge and public-release status are separate facts. #354 remains open for release/discovery and #355/#356 retain independently unconfirmed original causes.
