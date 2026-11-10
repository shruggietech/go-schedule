# Implementation Plan: Ship v2.0.3

**Branch**: codex/142-windows-version-fix-publication | **Date**: 2026-11-09 | **Spec**: [spec.md](spec.md)

## Summary

Ship integrated #458 through existing release workflows and finish #462 without additional operator approvals. Reconcile implemented #355/#356 outcomes without claiming a historical root cause.

## Technical Context

Go 1.26.0, existing Wails/WiX and GitHub Release/Pages workflows. No new dependency, runtime architecture or storage change. Preserve exact historical snapshots and assets. Native metadata readback remains distinct from installed/attended UI observations.

## Constitution Check

PASS before and after design: normal source PRs, full foreground verification, required automated CI, immutable source and asset checks. The latest owner direction explicitly authorizes automatic PR merges, tagging, staging and promotion, superseding older approval pauses. The two-total review cap is retained across the two source checkpoints. No operator desktop focus.

## Project Structure

README.md, CHANGELOG.md, tag-specific release notes, site/data comparisons/reviews, current docs and Spec-Kit lifecycle records. Existing Release and standing-waiver promotion workflows are reused unchanged.

## Complexity Tracking

No constitutional violation or new release framework. Optional agent-context hooks are skipped; plan reference is maintained directly.
