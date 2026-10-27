# Implementation Plan: Historical Windows incident differential diagnosis

**Branch**: `codex/123-historical-windows-diagnosis` | **Date**: 2026-10-27 | **Spec**: [spec.md](spec.md)

## Summary

Deliver the full #355/#356 bounded historical comparison, correct demonstrated defects and reconcile merged S122. The current baseline is e8cefc1. Recover original source from base 2b827e0 plus retained dirty overlay and verify the ff4658 source digest before calling it exact. Retain untouched originals.

## Technical Context

**Language/Version**: Existing Go 1.26, PowerShell and Bash.

**Primary Dependencies**: Existing SQLite, native IPC, MCP and owned CLI fixture packages; no new dependencies.

**Storage**: Private ignored copies and owned WAL stores, with conditional current component memory/WAL contrasts only after a useful divergence.

**Testing**: Test-driven replay selection regressions, compiled race-enabled historical/current integration tests, targeted IPC replay and canonical sh scripts/verify.sh all.

**Target Platform**: Native Windows comparison; existing Windows/WSL integration acceptance after actual application changes.

**Project Type**: Maintainer diagnostics in the existing scheduler repository.

**Performance Goals**: One thirty-minute investigation execution deadline, sequential groups, unchanged six concurrent writers, bounded CPU/memory/process/output controls.

**Constraints**: Hidden noninteractive launchers, immutable original receipts, no service/session mutations, independent results and no invented causality.

**Scale/Scope**: #355 and #356 with separately filed demonstrated defects; no release or clustered-scheduling work.

## Constitution Check

Before research: PASS. Preserve idiomatic code, deterministic regression contracts, race/security gates, native results, bounded owned resources and review-branch delivery.

After design: PASS. No production behavior change is justified before an actual finding. Small opt-in diagnostic selections extend the existing replay, without changing ordinary QA defaults or relaxing checks. Source identity and comparison provenance are explicit. Owner push/PR authorization is already present; merge/tag/promotion remain outside it.

## Project Structure

- specs/123-historical-windows-diagnosis/: specification, plan, research, data model, quickstart, tasks, investigation and verification/delivery records.
- internal/clievaluate/incident_replay_test.go: selected case allowance, collector selection and explicit controller/binary provenance for opt-in replay.
- docs/maintainers/windows-incident-replay.md: historical comparison and interpretation guidance.
- specs/122-windows-incident-replay/: factual post-merge reconciliation.
- .test-temp/s123/: ignored reconstructed source, prepared binaries, original hashes, campaign and acceptance receipts.

## Decisions and chronology

1. Reuse the existing strict binary receipt and catalog. A new general QA runner would duplicate delivered functionality.
2. Runtime comparison uses exact original generated binding wrappers/adapters against current wrappers/adapters, not merely the unchanged heartbeat entrypoints. Recover complete original source and run its original binding test separately; current instrumentation is attributed only to current observations.
3. Packaged IPC comparison initially uses the current controller with historical/current packaged binary sets. service.restart remains the current recorder-backed command harness; the other targets use the selected CLI and daemon. Historical daemon collection is explicitly disabled/unavailable. Current protocol timing differs from historical timing, so this cannot inherit original-controller identity.
4. Extend only missing opt-in replay controls under test-first discipline: positive selected case duration, explicit incident selection and emitted actual provenance. Keep defaults at supported 120 seconds and current collection enabled.
5. Before campaign execution, freeze a matrix of original/current binding tests across retained/current PowerShell and Bash (five repetitions, both schema modes and six writers, 360 writer outcomes), plus historical/current packaged WAL candidates at 15/120-second case allowances (three repetitions, seven steps, 84 outcomes). Declare any compatibility prerequisite failure explicitly. Total planned primary observations: 444 across six runtime conditions and four IPC conditions.
6. Use one shared thirty-minute execution deadline after preparation. Stop unexpected conditions for analysis; preserve the clock and first failure. Conditional crossed adapters, dependency or memory/WAL experiments require an observed divergence and remain within that deadline.
7. No matching recurrence means an inconclusive original-cause disposition and no speculative expansion. Correct independently demonstrated defects under separate issues; rerun full source-bound QA only after application changes. Canonical local CI parity remains mandatory for this source PR.
8. Run the optional agent-context hook after plan; no mandatory hooks are registered. Clarify/checklist have no extension hooks.

## Complexity Tracking

No constitution deviation or new architecture. Private reconstruction and test selection are diagnostic inputs, not installed product changes.
