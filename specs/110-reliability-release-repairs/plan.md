# Implementation Plan: Reliability and release-promotion repairs

**Branch**: `codex/110-reliability-release-repairs` | **Date**: 2026-10-21 | **Spec**: [spec.md](spec.md)

## Summary

Deliver #262, #271 and #335 in full, preserving individual outcome traceability. Bind release operations to one authenticated ID through a small standard-library Go helper built from trusted workflow source before tagged-source checkout. Preserve same-revision watcher state through observer recovery. Correct demonstrated literal-text coercion in both SQLite CLI parameter adapters, retaining failure output and bounded retries.

## Technical Context

**Language/Version**: Go 1.26.0 module, Bash and PowerShell 7; existing YAML workflows.

**Primary Dependencies**: Existing project dependencies only; net/http, encoding/json, crypto/sha256 and httptest suffice for release identity handling.

**Storage**: One versioned temporary release-selection document; existing watcher in-memory maps; existing disposable test SQLite schemas without migration.

**Testing**: Failing regression first, focused Windows/native Linux race tests, deterministic fake-clock recovery tests, real script twins, HTTP fixtures and workflow ordering checks; canonical sh scripts/verify.sh all.

**Target Platform**: Linux, Windows and existing native macOS observer CI. No installed-service changes.

**Project Type**: Daemon maintenance and release tooling.

**Performance Goals**: Retain watcher O(watcher count + pending paths) state and existing timer/debounce budgets; no per-poll unbounded dispatch history. Streaming asset I/O, finite HTTP timeouts and no parallel stress beyond bounded fixtures.

**Constraints**: No publication, live release mutation, tag, new dependency, schema migration or authorization-policy change. Preserve all eight verification gates and issue-level functional outcomes.

**Scale/Scope**: Three independently testable issue boundaries, one reviewable maintenance commit/slice.

## Constitution Check

- I: Focused helpers with contextual errors and documented lifecycle; no copied tag-resolution defect, no custom protocol/framework.
- II: New defects get failing regressions; keep current concurrent row-count and watcher exactly-once assertions; retain race and core coverage gates.
- III: Plain-language ID/tag/path and failed-writer diagnostics, documented operations and stable script exit codes.
- IV: Recovery state bounded by current definitions and pending candidates; native hints keep ordinary subsequent-write behavior; no indefinite retry.
- V: Installed specify, clarify, checklist, plan, tasks, analyze and implement commands drive the sequence; local branch and commit, then one pre-publication halt.

Pre-research gate: PASS. Post-design gate: PASS. No constitutional deviation is proposed. Pinned standing-waiver workflow and its automation policy may change within #262, with a dated CHANGELOG decision; full-evidence workflow stays unchanged.

## Project Structure

- specs/110-reliability-release-repairs/: specification, research, contracts, model, quickstart, tasks, analysis and local delivery evidence.
- scripts/release-promotion/: small CLI with resolve, download, ensure-checksum and publish commands plus HTTP fixture tests.
- .github/workflows/promote-release-waived.yml: trusted helper build, tagged validation, ID-bound helper calls.
- scripts/automation-check.sh and test/scripts/automation-check_test.sh: recognize actual helper mutation boundaries while retaining comparison ordering.
- test/integration/: standing-waiver workflow contract and real payload regressions.
- internal/watcher/manager.go and recovery tests: state retention for unchanged enabled revisions and bounded duplicate suppression.
- test/scripts/lib/{sqlite.sh,Sqlite.ps1}: exact UTF-8 text parameter expressions and retained null semantics.

## Phases and decisions

1. Preserve trusted automation executable before checkout of older tags. Resolve all authenticated pages, freeze release identity/metadata/asset identities and digests, and refuse ambiguity. Keep staging, exact candidate, allowlist, comparison, checksum and immutable-tag checks in workflow. No gh release command may independently resolve the tag.
2. Distinguish explicit Reload from observer recovery. Retain prior identity, pending candidates and last dispatched identity only for same enabled revision; baseline newly changed revisions. Exclude inactive pending candidates from timers/processing until registration recovers. Reconcile retained state and suppress recovery/late-hint duplicates without permanent native-write suppression.
3. Bind non-null script values using CLI-quoted CAST(X'UTF8HEX' AS TEXT), preserving explicit null handling and existing SQL numeric casts. Distinct valid hex session IDs must remain distinct, including the demonstrated scientific-notation collision. No speculative contention change is justified by passing stress runs.
4. Integrate, run focused regressions and canonical verification, account for each issue independently, update lifecycle/changelog and commit locally.

## Complexity Tracking

No deviations. The temporary selection document is required to bind separate workflow steps to one identity; a small compiled helper avoids re-resolving tags or loading missing code from old release tags. No general release platform or new application service is introduced.
