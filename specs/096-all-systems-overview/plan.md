# Implementation Plan: All Systems Operational Overview

**Branch**: `codex/096-all-systems-overview` | **Date**: 2026-09-21 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/096-all-systems-overview/spec.md`

## Summary

Add a read-only All Systems desktop destination that fans out across This computer and every saved connection profile, retrieves one bounded operational summary per daemon, preserves partial and session-stale results, and safely drills into the exact source daemon. The daemon exposes one additive observe-authorized summary endpoint backed by bounded SQLite queries, while the desktop performs generation-scoped refreshes with four-worker concurrency, five-second target deadlines, and no persistent overview cache.

## Technical Context

**Language/Version**: Go 1.26, TypeScript 5.6, React 19

**Primary Dependencies**: Go standard library, modernc SQLite, existing local IPC and HTTPS API clients, Wails desktop bridge, React, Vite

**Storage**: Existing SQLite task, run, alert, and notification tables; existing JSON connection profiles and native credential storage; no new persisted overview state

**Testing**: Go unit and integration tests, race detector, Vitest and Testing Library, Playwright accessibility checks, repository CI-parity scripts

**Target Platform**: Windows and Linux daemon, Windows desktop application, local IPC and enabled remote HTTPS transports

**Project Type**: Go daemon and API with a React desktop frontend

**Performance Goals**: Bound every daemon summary independently of history size, cap refresh fan-out at four targets, classify a nonresponsive target within five seconds, and keep the 100-profile desktop usable at 800 by 600 and 200 percent zoom

**Constraints**: Read-only observe authority, no secrets or task execution details in responses, no cross-session cache, no cluster semantics, exact target selection before drill-down, and no new third-party dependency

**Scale/Scope**: This computer plus up to 100 saved profiles; one aggregate endpoint, one desktop service, one new route, and focused supporting tests and documentation

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **Code Quality**: New concurrency has explicit generation cancellation, worker ownership, bounded target contexts, and race-tested service coverage. Public Go symbols receive intent-focused documentation.
- **Testing Standards**: Store aggregation, API authorization and redaction, target refresh failures, stale fallback, generation replacement, frontend sorting and filtering, and exact-target drill-down receive deterministic tests. No wall-clock sleeps are required for assertions.
- **User Experience Consistency**: Connection states reuse the existing typed taxonomy, times remain RFC 3339 at the API boundary, errors identify the affected registration and recovery action, and the page uses existing shell, card, focus, status, and responsive conventions.
- **Performance Requirements**: SQL counts and representative rows use time predicates and limits instead of loading history. Fan-out and response cardinality are explicitly bounded. No scheduler hot path changes or benchmark-sensitive behavior are introduced.
- **Autonomous Build-Phase Execution**: S096 is traceable to issue #182 and follows the full spec-kit sequence on a review branch. The operator explicitly authorized publication and the standard two-round review process in the kickoff, satisfying the publication halt in advance.
- **Engineering Constraints**: The contract is additive, works through both transports, requires only Observe authority, validates remote identity through the existing client, and returns no secret-bearing fields.
- **Outcome-based issue completion**: Issue #182 closes only when its functional overview and drill-down outcomes are delivered. Engineering checks remain PR gates rather than issue acceptance criteria.

No constitutional deviation is required.

## Project Structure

### Documentation (this feature)

```text
specs/096-all-systems-overview/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── system-summary.md
├── checklists/
│   ├── operational-overview.md
│   └── requirements.md
└── tasks.md
```

### Source Code (repository root)

```text
internal/
├── api/
│   ├── client/
│   └── server/
├── authorization/
├── domain/
└── store/

desktop/
├── systems/
├── frontend/src/
│   ├── components/
│   └── systems/
├── app.go
└── main.go

docs/
├── api.md
└── remote-access.md
```

**Structure Decision**: Extend the existing daemon API, store, desktop service, and feature-folder frontend architecture. The summary contract belongs in `internal/domain`, bounded reads belong in `internal/store`, the HTTP boundary belongs in `internal/api`, and cross-profile orchestration belongs in a new `desktop/systems` package. No parallel persistence or networking abstraction is introduced.

## Design Decisions

### One daemon-owned aggregate endpoint

Each daemon computes its own summary from its authoritative clock and database. This avoids transferring raw histories or reconstructing scheduling semantics in the desktop. The response schema is additive and carries counts plus at most one safe representative record per category.

### Independent clients, unchanged selected connection

The overview constructs temporary clients from saved profiles and native credentials without selecting them in the global connection router. Only a user drill-down calls the existing exact-profile selection path. Background observation therefore cannot redirect mutation controls.

### Generation-scoped bounded fan-out

The desktop service cancels the previous generation when a refresh begins, snapshots current registrations, and runs at most four workers. Each target gets a five-second child context. Results from older generations and removed registrations are discarded before publication.

### Session-only stale continuity

The service retains the latest successful summary per registration in process memory. A later failure may pair that summary with a stale state and original observation time. Restarting the desktop intentionally clears the cache.

### Context-preserving drill-down

Overview actions carry the registration key, destination, and optional task or record identifiers. The desktop first selects the exact local or saved profile, then changes route and supplies a source-context hint. Selection failure leaves the overview and its results intact.

## Complexity Tracking

No constitution violations or complexity exceptions are required.
