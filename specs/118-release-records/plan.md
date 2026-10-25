# Implementation Plan: S118 release record reconciliation

**Branch**: `codex/118-release-records` | **Date**: 2026-10-25 | **Spec**: [spec.md](spec.md)

## Summary

Reconcile completed S117 source delivery and supply one current release handoff for #354-#356. Preserve original failed captures, independently unresolved causes, reviewed pending content bindings and public v1.5.1 selection. This is a records-only increment through an unmerged reviewed PR.

## Technical Context

**Language/Version**: UTF-8 Markdown and existing JSON records; Go 1.26 verification tooling.

**Primary Dependencies**: Installed Spec-Kit 0.11.3, native GitHub Issues/PRs, existing release-source/comparison/lifecycle checks.

**Storage**: Versioned source records; ignored local receipts remain immutable.

**Testing**: Existing eight canonical gates, lifecycle and release-content guards, retained evidence hash audit. No new runtime tests or unchanged full candidate sweep.

**Target Platform**: Repository documentation consumed by maintainers on Windows/Linux/macOS.

**Project Type**: Release-process documentation increment for the existing scheduler.

**Performance Goals**: No application execution or resource-policy change.

**Constraints**: No merge, tag, staging, promotion, installed-service or DNS operation. At most two Codex rounds. No fabricated causal resolution or release identity.

**Scale/Scope**: S117 lifecycle records, current S118 handoff, existing Windows procedure links, changelog decision and three existing issue dispositions.

## Constitution Check

Before research: PASS. Existing architecture and security/data contracts remain intact; no new dependency, runtime surface, waiver or closure gate. Source identity, original failure attribution and explicit publication boundaries are required.

After design: PASS. Reuse existing validators and release flow. Updating the pinned Windows installation guide requires the dated changelog decision. All eight foreground gates precede every commit. Publication evidence stays outside implementation checkboxes. No complexity exception.

## Research and design

Read-only delegated lifecycle and release-content research completed under the installed plan command. Findings and alternatives are in [research.md](research.md). Record entities and transitions are in [data-model.md](data-model.md); the maintainer interface is [contracts/release-records.md](contracts/release-records.md). See [quickstart.md](quickstart.md) for validation.

## Project Structure

- `specs/117-candidate-qa-handoff/`: corrected maturity, task dispositions and actual delivery receipts.
- `specs/118-release-records/`: specification, design, tasks, current handoff and source verification.
- `specs/README.md`, `CLAUDE.md`, `.specify/feature.json`: aligned active feature and inventory.
- `test/windows/README.md`, `docs/INSTALL-windows.md`: current procedure references.
- `CHANGELOG.md`: Unreleased record correction and pinned-file decision; selected v1.6.0 section unchanged.
- `.github/release-notes/v1.6.0.md`, `site/data/reviews/v1.6.0.json`, `site/data/published.json`: audited, preserved content.

## Execution strategy

Complete lifecycle reconciliation first, then current release/issue handoff. Freeze completed source for existing canonical gates. Finalize factual verification and task metadata, then publish the authorized source PR. Hosted CI and individual reviews are external delivery receipts, with automatic initial review plus at most one requested second round.
