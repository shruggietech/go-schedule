# Implementation Plan: S126 Dependency and Release Convergence

**Branch**: `codex/126-dependency-release-convergence` | **Date**: 2026-10-29 | **Spec**: [spec.md](spec.md)

## Summary

Converge ready dependency proposals #401-#403 into one coherent candidate. Defer #400/#404 because the documented corrected driver remains unpublished. Reconcile S125 owner merge, disclose selected versions in pending v1.7.0 records, qualify the assembled source and publish one PR with at most two Codex rounds. Keep #394 open for actual release delivery after owner merge and separate source/asset approvals.

## Technical Context

**Language/Version**: existing Go 1.26.0 modules, TypeScript/React, Node 26.
**Primary Dependencies**: API runtime 1.8.0, YAML 3.0.5, Vite 8.3.5, Next 16.4.3, Fumadocs core/UI 16.17.8, MDX 15.4.10; retain SQLite 1.60.1/libc 1.77.1.
**Storage**: unchanged schema v24 and configuration; versioned Markdown and existing release/comparison JSON.
**Testing**: existing API/YAML/storage contracts, both read-only Go module graphs, fresh Windows/WSL unattended QA, headless browser/site and foreground canonical aggregate.
**Target Platform**: Windows/Linux/macOS supported source and hosted matrix; local Windows/WSL without attended desktop operations.
**Project Type**: existing scheduler daemon, CLI, Wails desktop and static Next documentation site.
**Performance Goals**: unchanged scheduler budgets; one browser/QA worker and existing bounded CPU/memory/process/output controls.
**Constraints**: hidden noninteractive launchers, no installed-service/session changes, no weakened gates or original evidence rewrites, no merge/tag/promotion.
**Scale/Scope**: three ready dependency proposals, one explicitly deferred storage proposal/bug and one partially delivered release issue; no new product feature or framework.

## Constitution Check

Pre-research and post-design gates pass: minimal coherent graph changes (I); all existing security/tenancy/race/coverage contracts and fresh assembled QA retained (II); configuration/interfaces/accessibility and diagnostics preserved (III); no new hot-path behavior or unsupported performance claim (IV); installed Spec-Kit sequence, blocking analyze, review branch and explicit PR authorization (V). Functional issue completion remains distinct from verification. No constitution deviation.

## Phase 0: Research

[research.md](research.md) records three disjoint read-only research tracks and version-selection evidence. All unknowns are resolved. Ready dependency updates proceed; unavailable SQLite correction does not hold the release path.

## Phase 1: Design and Contracts

[data-model.md](data-model.md), [contracts/delivery.md](contracts/delivery.md) and [quickstart.md](quickstart.md) bind dependency selection, preserved evidence and owner transitions. Existing API/CLI/UI contracts are unchanged; no replacement external interface is introduced.

## Project Structure

- go.mod/go.sum and desktop/go.mod/go.sum: selected complete Go graphs.
- desktop/frontend/package.json/package-lock.json: Vite refresh.
- site/package.json/package-lock.json: coupled Next/Fumadocs refresh.
- CHANGELOG.md and .github/release-notes/v1.7.0.md: accurate pending maintenance content.
- site/data/reviews/v1.7.0.json: hashes bound to the unchanged capability dataset and current pending prose.
- specs/125-reliability-delivery-decision/verification.md, delivery.md, handoff.md: factual merged-source reconciliation.
- specs/126-dependency-release-convergence/: specification, analysis, verification and same-slice release continuation.

## Verification and Publication

Freeze historical inputs, regenerate manifests/locks with existing tools, inspect graph deltas and perform read-only graph builds. Run new unattended capture and review every exact message in context, finalizing fixed-shape reports. Run headless browser/site checks and canonical sh scripts/verify.sh all in the foreground. Format GitHub prose with go run ./scripts/github-format, audit UTF-8/no-BOM and lifecycle. Commit and push the branch, open the source PR, monitor exact-head CI and every bot comment, request at most a second Codex round and stop for owner merge review. Supersede dependency PRs only after their selected equivalent changes actually integrate. Keep #394/#400/#404 and original #355/#356 open.

## Complexity Tracking

No new framework, workflow policy, timeout, schema, toolchain pin or production architecture. Reuse the existing release handoff and record factual continuation rather than creating another preparatory mechanism.
