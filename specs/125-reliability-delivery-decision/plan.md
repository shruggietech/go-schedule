# Implementation Plan: S125 Reliability Delivery and Scheduling Decision

**Branch**: `codex/125-reliability-delivery-decision` | **Date**: 2026-10-28 | **Spec**: [spec.md](spec.md)

## Summary

Correct #397's browser fixture ordering and preserve failed-run diagnostics. Fold that correction into the prepared v1.7.0 source, reconcile the actual S124 merge, and provide the complete existing release continuation path without crossing merge/tag/promotion boundaries. Assess #186/#189/#20 as one decision chain and explicitly hold #187/#188 unless real demand and approved semantics permit them.

## Technical Context

**Language/Version**: existing Go modules; TypeScript/React frontend; Node 26; existing Playwright.
**Primary Dependencies**: existing Playwright, GitHub Actions upload-artifact@v7 and release/comparison guards; no new dependency.
**Storage**: Markdown decision/handoff records and existing comparison-review JSON; no product schema change.
**Testing**: deterministic adverse-order browser regression, full headless frontend browser contracts, and canonical `sh scripts/verify.sh all`.
**Target Platform**: supported desktop interfaces; hosted Linux browser contract and native Windows local verification.
**Project Type**: existing scheduler desktop/CLI/daemon repository.
**Performance Goals**: unchanged scheduler budgets; one local browser worker, bounded artifact retention.
**Constraints**: hidden noninteractive console launches, no installed-service/session changes, no weakened assertions/timeouts, no tag/release mutation.
**Scale/Scope**: seven authorized issue records plus independently discovered source-document metadata defect #398; delivered source work and conditional/partial functional outcomes are reported separately.

## Constitution Check

Pre-research and post-design gates pass: idiomatic minimal correction (I), meaningful failure regression and full canonical verification (II), existing keyboard focus and identity semantics preserved (III), no scheduler hot-path change or new performance claim (IV), installed Spec-Kit sequence and one authorized review branch (V). The owner explicitly authorizes push/PR and up to two review rounds. Merge/tag/promotion are excluded. Issue closure is functional, not a CI gate.

## Phase 0: Research and Decisions

The installed planning command explicitly requests research agents. Separate agents assessed focus sequencing, release continuation and distributed demand; findings are consolidated in [research.md](research.md).

- Synchronize route focus in the browser fixture rather than changing Shell accessibility behavior without a demonstrated product cause.
- Retain browser failure files with the existing artifact action, 14-day retention and no continue-on-error on the exercise.
- Keep S124 capability cells and official vendor evidence; refresh only content-bound review hashes when pending changelog/notes change.
- Reuse existing v1.7.0 release workflows. Future merged source and concrete assets cannot be selected in this unmerged proposal.
- Available public demand evidence supports a reversible hold, not a permanent no-go. No concrete workload meets #186's required inputs; #187/#188 stay gated and #186/#189/#20 remain incomplete.
- Preserve the owner's July 24 clustered-execution intent. Do not adopt a fleet-key architecture or cancel the proposal without evidence/approval.
- Reuse unchanged assembled application QA; browser/CI and documentation changes receive focused regressions plus canonical verification. No additional unchanged historical campaign.

## Phase 1: Design and Contracts

See [data-model.md](data-model.md), [contracts/delivery.md](contracts/delivery.md) and [quickstart.md](quickstart.md). Document findings and an exact demand-intake template in `docs/distributed-scheduling.md`. Record the release continuation in `handoff.md`; link it from current Windows maintainer guides.

## Project Structure

- `desktop/frontend/e2e/cross-daemon-search.spec.ts`: route readiness and browser regression.
- `.github/workflows/ci.yml`: existing browser job failure artifact retention.
- `CHANGELOG.md`, `.github/release-notes/v1.7.0.md`, `site/data/reviews/v1.7.0.json`: truthful pending release content.
- `docs/distributed-scheduling.md`, `docs/remote-access.md`: demand evidence, alternatives, intake and gates.
- `docs/INSTALL-windows.md`, `test/windows/README.md`: current continuation pointer.
- `specs/124-reliability-release-preparation/`: factual owner-merge/main-failure reconciliation.
- `specs/125-reliability-delivery-decision/`: specification, analysis, decision, verification and delivery records.

## Verification and Publication

Run the original/fixed adverse-order scenario with headless Chromium, preserve both receipts, then full browser suite and canonical foreground aggregate. Freeze prior QA artifacts/binaries to ensure no evidence is rewritten. Format all GitHub prose with `go run ./scripts/github-format`, verify UTF-8/no-BOM and no corruption. Audit final SHA, exact-head CI, all review comments and no more than two Codex rounds. Publish only closing references for delivered #397; other incomplete outcomes use Refs.

## Complexity Tracking

No constitution deviation. Pinned CI workflow and Windows guide edits are necessary for #397/#394; record their dated rationale in CHANGELOG. No new framework, database, production coordination backend or speculative cluster API.
