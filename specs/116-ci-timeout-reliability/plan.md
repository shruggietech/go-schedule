# Implementation Plan: S116 CI timeout reliability

**Branch**: `codex/116-ci-timeout-reliability` | **Date**: 2026-10-25 | **Spec**: [spec.md](spec.md)

## Summary

Implement the owner-directed increased completion policy across all six CI/publication workflows and nested runners, deliver #372 with exact-target assertions retained, exercise source-bound #355/#356 native paths under declared new allowances and reconcile #354 source handoff. The complete numeric/classification contract is [contracts/timeouts.md](contracts/timeouts.md).

## Technical Context

**Language/Version**: Go 1.26.0, TypeScript/React, POSIX shell and PowerShell7.

**Dependencies**: Existing Cobra, Wails2.17, Vitest 5.0.5, TestingLibrary, Playwright, Next 16.4.2 and existing YAML parser dependency; no new external dependency.

**Storage**: Owned temporary memory/WAL fixtures and immutable local receipts; no persisted product schema change.

**Testing**: Actual delayed-work before/after experiments, configuration/policy mutation regressions, native Go race and six-writer/real IPC tests, canonical scripts/verify.sh all, exact-head hosted CI and at most two Codex rounds.

**Platforms**: Hosted Linux/macOS/Windows; native Windows and existing nonroot Ubuntu-24.04 WSL.

**Performance/constraints**: Completion time increases are authorized; production scheduler latency, identity, auth, SQLite contention, resource/ownership and intentional expiry semantics remain. Windows processes use the verified CREATE_NO_WINDOW launcher. No installed-service/MSI/guest/login/logout operations on the owner host.

**Scope**: #372, #354/#355/#356 and a separately tracked cross-CI timeout-policy task under #354. Source work stops before merge/tag/staging/promotion.

## Constitution Check

Pre-research and post-design gates pass: idiomatic bounded code, no new dependencies, meaningful before/after regression discipline, all identity/concurrency/security assertions retained, coverage threshold unchanged, production performance unaffected, complete Spec-Kit analysis before source implementation, dated pinned-workflow decision, review branch and owner-authorized push/PR. The latest owner instruction supersedes older issue no-increase advice; this changes verification completion budgets, not result truth.

## Project Structure

- `.github/workflows/*.yml`, `scripts/verify.sh`, `scripts/coverage-gate.sh`, `scripts/automation-check.sh` and policy fixtures: explicit job360min, Go 60min, deployment/wait/promotion allowances, enforced policy.
- `scripts/ci-timeout-check/`: cross-workflow job/nested-runner budget guard with omitted/reduced-value mutation coverage; reuse an existing YAML parser only if needed.
- `desktop/frontend/vite.config.ts`, `src/test-setup.ts`, `src/App.test.tsx`, browser configs and site startup/build config: explicit increased nested limits and real identity regression.
- `internal/clievaluate/`, `scripts/cli-qa.ps1`, owned IPC and packaged fixture tests: supported defaults, aligned reserves, explicit/resumed values, targeted native metadata evidence.
- `test/windows/`, `scripts/windows-qualification-session/`, promotion helper: audited incidental completion allowances; deliberate expiry cases preserved.
- `CHANGELOG.md`, `.github/release-notes/v1.6.0.md`, `site/data/reviews/v1.6.0.json`, `docs/cli-qa.md`, `docs/build-autopilot.md` and this feature: actual policy, source/known-issue/release distinctions.

## Execution strategy

1. Complete setup, all requirements checklists and blocking cross-artifact analysis.
2. Add policy and delayed-runner regressions before changes; record original failures independently.
3. Implement shared CI policy and nested caps, then #372 exact-target regression and native supported-default/reserve changes. Keep intentional expiry tests exact.
4. Run focused Windows/Linux writer and actual CLI/owned API metadata cases, including controlled work beyond the former limit and an explicit short-budget failure. Preserve runtime/source/binary/budget/cleanup receipts and unresolved original causes.
5. Reconcile pending release content and comparison digests, run all eight canonical gates in the foreground to completion, finish lifecycle metadata and commit.
6. Push the review branch and open one official PR under current owner authorization. Monitor exact-head checks and resolve every warranted comment in at most two Codex rounds; stop unmerged.

## Agent context and hooks

Optional agent-context hooks were not required; update the mandated CLAUDE plan pointer directly as the installed plan command specifies. No mandatory extension hook is configured. Checklist prerequisite initially required a plan; setup-plan supplied it before checklist execution was retried. Helper BRANCH output was empty, so direct verified Git branch identity is authoritative; .specify/feature.json selects S116 explicitly.

## Complexity Tracking

No constitution deviation. A small cross-file policy guard is justified by the owner-requested project-wide contract and prevents later silent resets; no new retry system, process manager or production timeout architecture is introduced.

## Implementation audit adjustment

The direct developer Makefile test/benchmark targets also inherited the old Go package watchdog. FR-003/T004 includes their 60-minute parity and policy mutation coverage for Make variable commands. This is the same authorized completion policy, with no new product boundary. A stale service-helper 10-second error label is corrected to its actual 120-second allowance.
