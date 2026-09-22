# Implementation Plan: Portable Automation Sources

**Branch**: `codex/099-portable-sources` | **Date**: 2026-09-22 | **Spec**: [spec.md](spec.md)

## Summary

Extend the S098 bundle contract to portable source and policy intent. Keep secrets and machine-local paths outside the document, use target-local path bindings in preview, and apply only reviewed operations. Existing v1 documents stay valid.

## Technical Context

Go daemon, SQLite store, Cobra CLI, Wails desktop, React frontend. Existing bundle API and authorization boundary remain authoritative. No new dependency or background process is required.

## Constitution Check

The design preserves code quality through one bundle domain model, testing through focused API and store cases, interface consistency across CLI and desktop, bounded deterministic comparisons, and the authorized autopilot PR workflow. No constitution deviation is planned.

## Decisions

1. Introduce schema v2 instead of changing the meaning of v1. This preserves existing document compatibility.
2. Keep watcher paths outside canonical bundles. A path is chosen for each target during preview and held in the single-use server plan.
3. Generate trigger credentials on target and create sources disabled. A transferred definition cannot silently become a live invocation surface.
4. Match notification channels by unique name. The target operator configures endpoint and authorization locally before policy references can bind.
5. Preserve target-only state and existing local paths and keys on updates.
6. Treat trigger-set rename and member-count changes as explicit conflicts. The store supports atomic retargeting but not atomic structural edits; pretending to apply those changes would make preview and apply disagree.

## Sequence

Implement model and canonical validation, then daemon export and preview, then apply and UI/CLI bindings, then cross-surface tests and CI parity.
