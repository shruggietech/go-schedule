# Tasks: S124 reliability release preparation

**Input**: [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md).

## Phase 1: Setup

- [x] T001 Confirm clean synced accepted base and S124 branch/selector in `.specify/feature.json`.
- [x] T002 Complete specify, autonomous clarify and requirements checklists in `specs/124-reliability-release-preparation/checklists/`.
- [x] T003 Complete research, design and context update in `specs/124-reliability-release-preparation/plan.md` and `CLAUDE.md`.

## Phase 2: Foundation

- [x] T004 Run blocking cross-artifact analyze using `specs/124-reliability-release-preparation/{spec,plan,tasks}.md`; resolve every finding before implementation.
- [x] T005 Freeze historical inputs and published v1.6.0 records in `.test-temp/s124/evidence-before.json`.

## Phase 3: US1, coherent candidate content

**Goal**: Select v1.7.0 and explain delivered changes and recovery without implying publication.

**Independent assessment**: Current preparation metadata agrees, full recovery guidance is actionable, existing latest release links remain unchanged.

- [x] T006 [US1] Prepare v1.7.0 notes and dated changelog with classification and known issues in `.github/release-notes/v1.7.0.md` and `CHANGELOG.md` (FR-001, FR-002, FR-008).
- [x] T007 [P] [US1] Align candidate examples and preserve latest badge in `README.md`, `docs/build-autopilot.md`, `docs/website.md` (FR-001, FR-004).
- [x] T008 [P] [US1] Complete platform upgrade/rollback and current handoff references in `docs/INSTALL-windows.md`, `docs/INSTALL-linux.md`, `docs/INSTALL-macos.md`, `test/windows/README.md` (FR-003, FR-008).
- [x] T009 [US1] Reconcile owner merge and accepted-main evidence in `specs/123-historical-windows-diagnosis/spec.md`, `verification.md` and `specs/README.md` (FR-009).

## Phase 4: US2, current qualified comparison

**Goal**: Refresh all 80 evidence cells without changing published v1.6.0 records.

**Independent assessment**: New snapshot is selectable through existing tooling; current sources and qualifiers are present and all ten semantic decisions bind settled content.

- [x] T010 [US2] Add current product contexts, dates and evidence notes in `site/data/comparisons/v1.7.0.json` (FR-005).
- [x] T011 [US2] Bind settled dataset/changelog/notes digests and all ten decisions in `site/data/reviews/v1.7.0.json` (FR-006).
- [x] T012 [US2] Run existing source/comparison guards and public-default checks; record source-only limits in `specs/124-reliability-release-preparation/verification.md` (FR-004, FR-006).

## Phase 5: US3, current source qualification and handoff

**Goal**: Deliver a traceable current source proposal with explicit remaining release operations.

**Independent assessment**: Run, input, message-review, exact-source CI and bot receipts have separate attributable identities and the PR remains unmerged.

- [x] T013 [US3] Freeze assembled source and build six actual dev-version QA inputs under `.test-temp/s124/qa-inputs-01` (FR-007).
- [x] T014 [US3] Capture current Windows/WSL operations with incidents in `.test-temp/s124/qa-run-01`, preserving first failures, overall deadline and owned-runtime cleanup (FR-007).
- [x] T015 [US3] Evaluate every current message against seven rubric criteria in `.test-temp/s124/reviews-01.json`; finalize, validate and render through the original coordinator (FR-007).
- [x] T016 [US3] Preserve demonstrated new defects in separate GitHub issues, correct #395 in `scripts/release-source-check/main_test.go` with failing-then-passing regression evidence, and record results in `specs/124-reliability-release-preparation/verification.md` (FR-007, FR-010).
- [x] T017 [US3] Run all eight canonical foreground gates against frozen implementation and record terminal receipt, complete timing and source/evidence drift in `.test-temp/s124/verification-receipt.json` (FR-007).
- [x] T018 [US3] Update #394/roadmap selected-version source status and source handoff in `specs/124-reliability-release-preparation/verification.md`, keeping actual release delivery and #355/#356 open (FR-010).

## Phase 6: Final implementation records

- [x] T019 Finalize lifecycle, tasks, inventory and factual evidence in `specs/124-reliability-release-preparation/`, `specs/README.md`; audit formatting, UTF-8/no-BOM, immutable inputs and source drift.

## Dependencies and execution

T001-T005 precede implementation. US1 and US2 content can proceed in separate files; T011 waits for settled T006/T010 bytes. US3 waits for all assembled content and formatting. Capture freezes source and binaries. T017 follows current QA and any fixes; final factual prose has a separate audit. Publication follows local qualification, and final owner handoff follows exact-head CI and satisfied reviews.

Parallel examples: T007 and T008 affect separate guides; research for US2 was independent of version/QA research. Execution is sequential where shared files, content digests or runtime budgets require it.

## Implementation strategy

Deliver US1 candidate content first, then US2 comparison, then US3 current qualification and source PR. Do not convert this source-only MVP into a release operation. Every required task remains tracked through the explicit owner-review boundary.

## Authorized publication protocol (separate workflow records)

The original T020-T022 publishing actions remain fully required and are tracked in .test-temp/s124/publication-ledger.json and the official PR receipts. They are workflow acceptance after local source implementation, rather than unfinished source tasks.

1. Commit verified source with the Codex trailer, automatically push this branch, create and attach the official PR referencing #394 and closing #395 after merge.
2. Monitor CI on the actual proposed commit, address every review comment and resolve warranted findings, using automatic Codex round one and at most one requested second round. Preserve each receipt before a summary can be overwritten.
3. Audit clean refs, head, terminal successful checks, threads, labels/milestone, issue links and release boundaries; publish the final owner handoff. Leave the PR unmerged and do not create a tag, draft release or promotion.

This bookkeeping separates 19 completed implementation tasks from three unchanged publishing steps; it removes no owner-authorized work and grants no additional release authority.
