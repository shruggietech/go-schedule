# Implementation Plan: S117 candidate QA and release handoff

**Branch**: `codex/117-candidate-qa-handoff` | **Date**: 2026-10-25 | **Spec**: [spec.md](spec.md)

## Summary

Run one complete existing Windows and nonroot WSL candidate evaluation, review every captured message in context, correct demonstrated project defects, and deliver an exact source release handoff through one reviewed PR. Scope is #354, #355 and #356, plus demonstrated child defects #375 and #376. The owner explicitly authorizes push and PR publication, with at most two Codex rounds; merge, tagging, staging and promotion remain outside this slice.

## Technical Context

**Language/Version**: Existing Go 1.26 module, PowerShell 7, Bash, Node.js and TypeScript toolchains.

**Primary Dependencies**: Existing SQLite, named-pipe/Unix-socket IPC, evaluator and frontend dependencies. No new dependency or QA framework.

**Storage**: Owned temporary SQLite fixtures and schema-versioned JSON evidence under ignored `.test-temp/s117`; durable source disposition under this feature directory.

**Testing**: Existing race-checked Go regressions, full CLI manifest and contextual message review, canonical eight-gate driver, source release/comparison guards and exact-head hosted CI.

**Target Platform**: Native Windows and existing nonroot Ubuntu-24.04 WSL. Hosted CI provides additional supported-platform evidence.

**Project Type**: Daemon, CLI and desktop application with an existing unattended evaluation coordinator.

**Performance Goals**: Preserve S116 supported defaults: 120 seconds per case, two hours total, 25% CPU, 1 GiB memory, Windows 32 processes and Linux 128 tasks. Preserve intentional expiry tests and all original data/concurrency/authority assertions.

**Constraints**: Hidden noninteractive processes; no installed service, session, installation, DNS or release mutation. No source edits or commits during build, capture, resume or report finalization. Original failed evidence remains immutable.

**Scale/Scope**: One full current manifest and every required message review, bounded focused follow-ups only for observed failures, proportionate corrections and one official source PR.

## Constitution Check

Pre-research and post-design checks pass. Quality and error handling remain enforced; behavioral corrections require failing-before/passing-after regressions under race. Ownership, authentication, migrations, clocks and concurrency assertions cannot be weakened. No measured performance-sensitive production change is presumed. All eight canonical gates run in the foreground before each commit. Functional issue completion remains distinct from engineering verification. The owner has already supplied the normal publication authorization; stop before merge, tag, staging and promotion.

## Decisions and Phase Order

1. Complete and analyze Spec-Kit artifacts before execution. Reuse the established evaluator, schemas, review prompt and report renderer.
2. Freeze the accepted main commit plus the explicitly dirty S117 design snapshot. This avoids an otherwise unnecessary intermediate commit and full canonical run solely to change a truthful dirty flag to clean. Record the actual commit, dirty flag, source digest and six executable hashes; never claim the snapshot is the later final PR.
3. Inventory existing tools and WSL controls without provisioning. Build inputs into `.test-temp/s117/qa-inputs-01`, capture into `qa-run-01` with no budget override, and retain the native outcome. Do not edit any source or task checkboxes until review/finalization is complete.
4. Read every fresh corpus message and relevant occurrences. Produce exact-bound seven-rubric reviews and actionable findings, then finalize, validate and render through the existing coordinator. Exit 2/3 is a valid failed/blocked report, not evidence corruption.
5. Inspect observed failures before attribution. Preserve original #355 runtime/prerequisite and #356 IPC causes independently. All completed requests plus slow non-request phases do not prove a transport failure. Successful repetitions do not close unknown causes.
6. After freezing the original final report, reproduce and fix demonstrated defects with tests. Use new immutable input/run directories if another source evaluation becomes necessary; do not rerun a full capture just to seek a preferred outcome. Track distinct defects using existing taxonomy and milestone 16.
7. Reconcile source handoff and only materially changed pending release guidance. Keep public v1.5.1 selection and governed capability decisions. Refresh precise notes/changelog bindings only when their bytes change.
8. Run canonical eight gates and release/comparison guards, commit, push and open one PR. Monitor exact-head checks and individual feedback, allow at most two Codex rounds, and hand the unmerged PR to the owner.

## Project Structure

- `specs/117-candidate-qa-handoff/`: spec, plan, research, data model, contracts, quickstart, tasks, verification and delivery.
- `scripts/cli-qa.ps1`, `cmd/gosched-qa/`, `internal/clievaluate/`: existing capture, identity, controls, review validation and reporting.
- `internal/clifixture/`, `test/integration/`, `test/scripts/lib/`: owned fixtures and potential demonstrated diagnostic/native corrections.
- `.github/release-notes/v1.6.0.md`, `CHANGELOG.md`, `site/data/reviews/v1.6.0.json`: governed pending source claims, changed only when findings warrant it.
- `.test-temp/s117/`: ignored inputs, immutable capture/reviews and private verification/publication receipts.

## Demonstrated implementation boundaries

The immutable full capture revealed missing human schedule descriptions in bounded observations (#375). Supplemental state inspection revealed valid MCP authority arrays discarded by diagnostic attachment and a reproduced null-state panic (#376). Correct the store/API projection and evaluator respectively, with real persisted/authenticated MCP and state-preservation regressions. These are permitted FR-006 discoveries, not inferred fixes for the original CLR or IPC causes.

## Complexity Tracking

No constitution deviation or new architecture. Truthfully identified dirty-source evaluation is supported by the existing receipt contract. Read-only research agents were required by the installed planning command; execution remains sequential and bounded.
