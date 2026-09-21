# Implementation Plan: Portable Automation Bundles and Drift Reporting

**Branch**: `codex/098-portable-bundles` | **Date**: 2026-09-21 | **Spec**: [spec.md](spec.md)

## Summary

S098 completes #184 with one versioned, canonical portable-bundle contract. A daemon exports safe declarative intent, validates it, creates a target-bound preview, compares drift without mutation, and applies only the reviewed operations in deterministic dependency order. CLI and desktop use that same contract. The design deliberately keeps platform-specific execution inputs and secret-bearing records outside the format.

## Technical Context

**Language/Version**: Go with TypeScript React desktop frontend

**Primary Dependencies**: Go standard library, existing Cobra CLI, Wails bridge, existing authenticated API client

**Storage**: Existing SQLite store, plus a durable portable identity field for every transferable definition

**Testing**: Go unit and API tests, race tests, Vitest component tests, Playwright workflow coverage, canonical project verification

**Target Platform**: Local and opt-in authenticated remote daemons on supported desktop and headless platforms

**Project Type**: Daemon API, CLI, Wails desktop application

**Performance Goals**: Deterministic plans for 100 objects without unbounded reads, worker fan-out, or background synchronization

**Constraints**: No secret export, no implicit removal, no hidden target choice, no cross-daemon transaction, no automatic retry after uncertain remote mutation

**Scale/Scope**: One selected daemon per export, validation, compare, preview, or apply request; v1 bundle format only

## Constitution Check

| Principle | Design response |
|---|---|
| Code quality | The bundle package owns canonicalization, validation, matching, planning, and errors rather than duplicating behavior among API, CLI, and desktop. |
| Testing | Golden canonicalization, redaction, reference, plan binding, partial-outcome, remote uncertainty, API contract, CLI, and desktop tests precede or accompany behavior. |
| UX consistency | CLI and desktop expose the same explicit target, plan identity, changes, conflicts, drift, and item outcomes. |
| Performance | Enumeration is bounded by existing list pagination; normalization uses stable sorts and linear reference maps. |
| Autonomous execution | This plan is traceable to #184 and follows the configured review-branch PR workflow. |

No constitution violation is required.

## Architecture Decisions

1. **One canonical domain package**: `internal/bundle` is the single implementation of format v1, secret filtering, validation, fingerprints, comparison, and plan generation. API, CLI, and desktop consume typed results. This prevents one interface from treating a bundle more permissively than another.
2. **Explicit portable identities**: Records receive stable `portable_id` values at first export or bundle-aware inspection. IDs are not display names and are serialized as the cross-daemon match key. Existing internal UUIDs and target-specific IDs never enter the document.
3. **Target-bound previews**: A preview contains the selected daemon ID, a deterministic target fingerprint, the canonical bundle digest, and its planned operations. Apply recomputes the fingerprint and rejects stale previews before any operation begins.
4. **Safe v1 coverage**: Groups, task metadata and schedule intent, and chains are represented. Trigger keys, watcher paths, trigger sets, notification-policy references, raw commands, directories, users, environments, keys, endpoints, and authorizations become explicit excluded findings.
5. **No destructive semantics**: Omission produces target-only drift only. Apply has create and update operations but no delete, disable, key rotation, credential, or destination operation.
6. **Authority is capability-aligned**: Bundle export, validation, preview, and comparison are Observe. Apply is Manage. The remote allowlist exposes those operation names only with their matching capability.

## Project Structure

```text
internal/bundle/                     # Canonical document, validation, planning, outcomes
internal/store/                      # Portable identity persistence and complete snapshots
internal/api/server/bundles.go       # HTTP contract and target-bound apply
internal/api/client/bundles.go       # Local and remote client methods
internal/cli/bundle.go               # Headless export, validate, compare, preview, apply
internal/remote/                     # Authenticated remote operation allowlist
desktop/bundles/                     # Target-aware desktop orchestration
desktop/frontend/src/bundles/        # Bundle workflow UI and bridge
api/openapi/remote-v1.yaml           # Remote contract
specs/098-portable-bundles/          # S098 specification artifacts
```

**Structure Decision**: The canonical bundle package is a dependency-light pure Go layer. It accepts a sanitized snapshot and produces bundle, validation, comparison, and plan values. Store and server code only adapt persistence and authorization, avoiding business rules embedded in HTTP handlers.

## Complexity Tracking

| Added complexity | Why needed | Simpler alternative rejected because |
|---|---|---|
| Versioned canonical bundle package | A transfer format needs one deterministic truth shared by all interfaces. | Ad hoc endpoint JSON would produce incompatible, difficult-to-review behavior. |
| Target fingerprint and plan identity | Apply must prove it operates on the exact reviewed target state. | Applying a raw bundle after a preview could silently act on changed state. |
| Durable portable identities | Display names are mutable and ambiguous across daemons. | Matching solely by names risks overwriting an unrelated definition. |
