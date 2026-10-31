# Implementation Plan: S129 Prelaunch Baseline Publication

**Branch**: `codex/129-prelaunch-release` | **Date**: 2026-10-31 | **Spec**: [spec.md](spec.md)

**Input**: `specs/129-prelaunch-release/spec.md`, owner kickoff and [#412](https://github.com/shruggietech/go-schedule/issues/412).

## Summary

Prepare and review one coherent next release for the merged current-only prelaunch baseline. Reuse existing staging and standing-waiver promotion workflows, preserve actual public v1.7.0 discovery during preparation, and retain a concrete continuation in this same S129 slice after later owner approvals. The current authorized boundary is verified source push and one official PR, followed by owner review/merge.

## Technical Context

**Language/Version**: Go 1.26.0 in root and desktop modules; Node 26 for the existing Next.js site and Vite/Wails frontend.
**Primary Dependencies**: Existing pinned graph, including Wails 2.17.0, SQLite 1.60.1/libc 1.77.1. No dependency upgrades or new framework.
**Storage**: Current application schema 24, preferences version 1, metrics format 4 and qualification evidence revision 2 remain as implemented by S128. No conversion or historical upgrade work.
**Testing**: Foreground `sh scripts/verify.sh all` (all eight existing gates), focused release/comparison checks, clean npm restore and existing headless browser exercises; retained-input/source digests.
**Target Platform**: Existing Linux, macOS and Windows package targets. Local tests remain headless; attended native observations are untested under the standing waiver.
**Project Type**: Daemon/CLI, nested Wails desktop module and static website.
**Performance Goals**: No scheduler hot-path change. Existing dispatch budget and coverage/benchmark requirements remain unchanged.
**Constraints**: UTF-8 without BOM; no em dashes or width-wrapped GitHub prose; preserve 215 retained inputs, historical releases/notes/comparison bindings and accepted asset bytes. At most two Codex requests including automatic PR-open review. No scheduled monitoring or release mutation now.
**Scale/Scope**: One candidate release and one source PR; public delivery stays open in #412 until actually published/discoverable.
**Resolved design**: Prepare v2.0.0 with root /v2 and desktop /v2/desktop module namespaces, active imports and linker targets aligned. Add its comparison and content-bound review without changing public v1.7.0 selection. Reuse the existing eight-asset staging and nine-asset promotion contract; see research.md.

## Constitution Check

Pre-research gate: PASS. Principles I-IV retain existing code quality, deterministic/race-checked tests, current persistence/refusal/security and performance requirements. Principle V is followed through the installed Spec-Kit sequence and blocking analyze gate; the owner's explicit push/PR instruction already supplies the normal pre-publication authorization. It does not authorize merge, tag/draft or promotion. Current-only policy 4.0.0 is preserved. #412 tracks actual functional release delivery and does not make verification a closure condition. No unjustified deviation.

Post-design gate: PASS. The major module namespace and existing linker targets are mechanical release-identity changes; HTTP API and current-state formats remain unchanged. Record the pinned release-workflow/Windows-guide decision in CHANGELOG.md. All acceptance, evidence, review-cap and later owner boundaries remain explicit.

## Project Structure

### Documentation (this feature)

`specs/129-prelaunch-release/` contains the specification, requirements/release checklists, plan, research, data model, release contract, quickstart, tasks, analysis, verification and concrete handoff.

### Source Code (repository root)

- `CHANGELOG.md`, `README.md`, `.github/release-notes/`: candidate identity and truthful current guidance.
- `site/data/comparisons/`, `site/data/reviews/`: new candidate capability/review records; `site/data/published.json` keeps actual public selection.
- `go.mod`, `desktop/go.mod`, active Go imports and existing build version injection: align /v2 identity and preserve root/desktop internal access.
- `scripts/release-source-check/`, `scripts/comparison-check/`, existing release workflows and platform guides: reuse established identity/approval contracts; change only a demonstrated candidate mismatch.

**Structure Decision**: Extend the existing release and comparison surfaces. No new release framework, historical conversion or product subsystem.

## Complexity Tracking

No constitutional exception is proposed. The v2 module namespace and two existing release linker targets are required for coherent module-aware major-version tagging. No product runtime behavior or new framework is introduced. See research.md for rejected alternatives.

## Discovered bounded correction (#413)

After the first full aggregate passed with unchanged source, correct the demonstrated mismatch between full v-prefixed workflow injection and the source guard's stripped README example. A new regression must reject unprefixed output and accept the actual injected tag. Retain all existing badge/date/version/notes negative cases and historical tagged content. Runtime and workflow version injection remain unchanged. Repeat final-source verification after this bounded guard/prose change; preserve the earlier complete run separately.
