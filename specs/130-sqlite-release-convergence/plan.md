# Implementation Plan: SQLite convergence and release delivery

**Branch**: `codex/130-sqlite-release-convergence` | **Date**: 2026-10-31 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from specs/130-sqlite-release-convergence/spec.md; #404 and #412.

## Summary

Adopt published SQLite v1.61.1 with its exact libc v1.79.0 pair, regenerate both Go module graphs and preserve normal read-only builds. Add focused real-driver cancellation/prepared-operation evidence and update only current candidate dependency facts and their content bindings. Continue actual v2.0.0 delivery through the existing S129/#412 path after owner integration and separate approvals. The current kickoff ends at owner final review/merge.

## Technical Context

**Language/Version**: Go 1.26.0; existing TypeScript frontend, no frontend dependency changes.

**Primary Dependencies**: modernc.org/sqlite v1.61.1 and modernc.org/libc v1.79.0; existing Wails v2.17.0 and all unrelated direct dependencies preserved.

**Storage**: Existing application schema 24, foreign keys, one production connection and WAL. No production pool/pragmas, SQL schema or API changes planned.

**Testing**: Existing Go/race/bootstrap/security/audit suites plus focused real SQLite busy/commit cancellation, prepared-query correctness and same-machine repeated benchmarks. Canonical sh scripts/verify.sh all supplies all eight gates.

**Target Platform**: Existing Linux/macOS/Windows daemon, CLI and desktop targets.

**Project Type**: Daemon/CLI and locally replaced Wails desktop Go modules.

**Performance Goals**: Preserve existing scheduling p99 target and no prepared-query regression above 10 percent without measured explanation. Use identical benchmark code/workloads before and after, record medians and noise rather than claiming equality from one run.

**Constraints**: No migrations, historical upgrade campaign, source review cap above two, release mutation, active-workstation native interaction or scheduled monitor. All 221 retained inputs frozen. UTF-8/no-BOM, no em dash or width-wrapped GitHub prose.

**Scale/Scope**: Four dependency metadata files, a small storage regression/benchmark file, current prepared notes/changelog/review bindings, existing handoff and S130 specification artifacts. No new dependency-management or release framework.

## Constitution Check

Pre-research and post-design gates pass: idiomatic existing architecture; behavioral tests precede dependency change; root/desktop race suites and all eight canonical gates remain; real driver timing is bounded by context cancellation and retained lock ownership, not engine clock changes; performance gets before/after benchmarks; current access/refusal/atomicity/audit tests remain; owner authorizes source push/PR explicitly and retains merge/tag/promotion. No principle or security test is weakened. Existing dependency licenses and transitive notices are reviewed. No constitutional deviation.

## Project Structure

### Documentation (this feature)

spec.md, plan.md, research.md, data-model.md, contracts/storage-release.md, quickstart.md, tasks.md, checklists/requirements.md, checklists/dependency-release.md, analysis.md, verification.md and delivery.md.

### Source Code (repository root)

- go.mod, go.sum, desktop/go.mod and desktop/go.sum: coherent selected graph.
- internal/store/sqlite_driver_test.go: real-driver regression and prepared-query benchmark.
- internal/store/store_test.go: correct two misleading historical-migration descriptions and the obsolete weak schema-version assertion while preserving current behavior assertions.
- .github/release-notes/v2.0.0.md, CHANGELOG.md and site/data/reviews/v2.0.0.json: current dependency facts and exact content binding.
- specs/129-prelaunch-release/handoff.md: factual dependency-hold supersession, retaining the existing staging/promotion continuation.
- desktop/README.md: current Wails build command and authoritative dependency-manifest references.
- .specify/feature.json, CLAUDE.md and specs/README.md: exact S130 selector/context/lifecycle.

**Structure Decision**: Keep the two existing modules and store architecture. Actual -mod=readonly resolution/tests demonstrate graph completeness; do not create a generalized manifest checker solely for this update.

## Decisions and Alternatives

1. Choose published v1.61.1/libc v1.79.0 over held v1.61.0, an unpublished fork or indefinite deferral because the owner approved this candidate lane and upstream published its correction. Actual module metadata and product tests decide acceptance.
2. Preserve the production one-connection/WAL policy. Test rollback-journal commit cancellation as a driver boundary on owned test connections; do not change product journaling or pool size to manufacture a new feature.
3. Reproduce #404 by updating root first and recording the expected normal desktop read-only failure before tidying desktop. Then resolve both graphs/read-only tests/checksums and demonstrate zero metadata drift.
4. New tests exercise repeated prepared queries on memory and disk/WAL and actual lock cancellation. Retain old-driver failure/benchmark receipts and compare identical final workloads. Existing current-state and audit tests cover durable product behavior.
5. Keep v2.0.0, /v2 imports and /api/v1. Update candidate notes and their review digests without changing comparison capability data or historical releases. Actual published selection remains v1.7.0.
6. Source implementation of #404 and public delivery of #412 remain separate functional outcomes. On integration #404 may close and #400 may be superseded; #412 stays open until actual publication/discovery. Continue its existing workflows after exact-source staging and concrete-asset approvals, never create another preparation-only framework.
7. Three bounded read-only research agents were dispatched as required by the installed /speckit-plan command: upstream candidate/license contract, current store tests and cancellation gaps, and module/release identity. No agent edits implementation.

## Implementation and Verification Order

First preserve source/input baseline and add tests/benchmarks on the old pair. Retain expected cancellation red and baseline benchmark results. Update root to the exact candidate, retain desktop graph red, then tidy desktop and verify graph/cancellation green. Update candidate facts/bindings and issue traceability. Freeze implementation and run all canonical gates in the foreground, watching terminal completion. Inspect scoped failures and fix without weakening gates. Final factual records and format/integrity/lifecycle checks precede the explicit-owner-authorized push/PR. Hosted CI and at most two review requests follow; stop for owner review/merge.

## Complexity Tracking

No constitutional violations or architecture expansion. Real I/O cancellation tests use owned SQLite connections and explicit cancellation timer lifecycle; elapsed bounds are watchdogs around actual lock behavior, not scheduling-clock assertions.

## Factual Guide Correction Discovered During Verification

The current desktop README still names Wails2.15.0 and obsolete frontend versions despite module/tooling Wails2.17.0. Correct its command to the tool already used by the passing native gate and replace duplicated version lists with authoritative manifest/lock links. This proportional dependency-guide correction adds no runtime or dependency changes. Rerun documentation and source/format/lifecycle checks after that edit; preserve the completed aggregate and its quiescent source identity separately.
