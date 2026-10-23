# Implementation Plan: S112 cumulative feature release

**Branch**: `codex/112-cumulative-feature-release` | **Date**: 2026-10-22 | **Spec**: [spec.md](spec.md)

## Summary

Correct #352 and #353 before dependency convergence, review all eight pending proposals, align Wails build tooling, and prepare v1.6.0 documentation/comparison and exact-source candidate instructions. This PR references publication coordinator #354 without closing it. Owner authorization covers push and review fixes; merge, tags and public release remain separate.

## Technical Context

**Language/Version**: Go from go.mod, PowerShell 7, POSIX shell, TypeScript/React.

**Primary Dependencies**: Existing SQLite, Wails desktop and static Next.js website; review pinned dependency proposals without adding a new runtime subsystem.

**Storage**: Existing forward-only v23 notification migration, unchanged by S112. Fixture diagnostic JSON is bounded and secret-free.

**Testing**: Canonical `sh scripts/verify.sh all`, Go regression tests, script fixture tests and owned Windows/WSL QA.

**Target Platform**: Windows, Linux and macOS; real service and installer mutation only on disposable hosted CI.

**Project Type**: CLI, daemon, desktop and public documentation site.

**Performance Goals**: Retain scheduler dispatch and resource contracts; bound setup and runtime independently.

**Constraints**: Hidden noninteractive Windows children, no account/session/service disruption, no tag or publication.

**Scale/Scope**: Two independent CI defects, eight dependency proposals and one cumulative release preparation increment.

## Constitution Check

- I: Idiomatic small corrections, explicit failure diagnostics, no duplicate scheduler implementation. PASS before and after design.
- II: Behavioral regressions precede fixes; retain race, coverage, IPC and release integrity checks. PASS.
- III: Setup, execution, diagnostics and cleanup messages retain their distinct causes. PASS.
- IV: No hot-path scheduling change or performance claim; diagnostic retention remains bounded. PASS.
- V: Installed hyphenated Spec-Kit sequence and blocking analysis; explicit owner push authorization; stop before merge/release. PASS.

## Project Structure

- `test/windows/Invoke-ServiceCoreCI.ps1`: terminal selection and failure retention.
- `test/integration/`: executable PowerShell regression and existing release contracts.
- `internal/clientsecret/context_test.go`: separate bounded build and runtime contexts.
- `go.mod`, `desktop/go.mod`, frontend/site package and lockfiles: dependency convergence.
- `.github/workflows/`, `scripts/verify.sh`, `scripts/automation-check.sh`, policy fixtures: reviewed action/Wails pin changes.
- `docs/`, `CHANGELOG.md`, `site/data/comparisons/`, `specs/README.md`: release preparation and truthful publication state.

## Decisions and execution order

1. Diagnose the original source first. Local packaged-helper baseline passed in 2.50 seconds; this does not refute the hosted 45-second timeout. The isolated unchanged Wait-ProbeRun fixture with queued-first/completed-second observations reproduces masking without installing a service.
2. Select terminal records before taking the first match. Preserve unsuccessful terminal outcomes for existing contract checks rather than filtering only successes. Keep the original 45-second service timeout.
3. Give credential fixture construction a bounded setup context and create runtime context only after build success. Preserve real packaged executable, hidden launch, parent/image/peer assertions and owned process draining.
4. Capture owned service diagnostics before cleanup and persist cleanup results independently of primary failure. Avoid environment dumps or broad process inventories.
5. Review dependencies from current PR heads; align all Wails runtime/CLI pins as one graph. Review upstream Pages runtime and artifact changes with policy fixtures.
6. Keep v1.5.1 as current public release. Prepare v1.6.0 comparison and notes with clear pending state; do not produce a purported immutable release candidate from an uncommitted source or trigger a tag workflow.
7. After canonical gates and isolated QA, publish this PR under existing authorization, handle all review threads, request at most one second Codex round, then hand off exact-head CI evidence.

## Complexity Tracking

No constitutional exceptions or new subsystem. Existing release and QA infrastructure is reused. Actual release asset publication remains an explicitly outstanding outcome of #354.
