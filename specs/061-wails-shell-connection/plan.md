# Implementation Plan: Production Wails Shell and Local Connection

**Branch**: `codex/061-wails-shell-connection` | **Date**: 2026-09-07 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/061-wails-shell-connection/spec.md`

## Summary

Graduate the approved S060 Wails and Calm Operations decisions into a production-intent, non-shipping `desktop/` module. Build a reusable React component catalog and application frame over a Go-owned connection manager that uses the existing protected local IPC client by default, owns request, retry, event, and shutdown lifecycles, exposes only sanitized generation-stamped snapshots and events, and remains structurally ready for future remote adapters. Preserve every current Fyne and release input until #157.

## Technical Context

**Language/Version**: Go 1.25.0, TypeScript 5.6.3, React 19.1.0

**Primary Dependencies**: Wails v2.14.0, Vite 7.3.6, Vitest 5.0.0, Testing Library 16.3.3, Playwright 1.63.0, axe-core 4.13.0

**Storage**: No new persisted application data; appearance is session-local until #156 defines preference migration

**Testing**: Go unit and race tests, Vitest component and accessibility tests, Playwright browser contract, native Wails builds on Windows, macOS, and Linux, canonical repository verification

**Target Platform**: Windows, macOS, and Linux desktop platforms already supported by the project

**Project Type**: Existing Go daemon, CLI, and Fyne desktop plus an isolated production Wails nested module

**Performance Goals**: Local connection outcome within two seconds; state events delivered without duplicate streams; responsive shell at 900 by 650 through 1440 by 900 and 200 percent zoom

**Constraints**: No network listener, no raw transport detail in the frontend, no remote assets or telemetry, no workflow migration, no current release-input change, no wall-clock sleeps in deterministic manager tests

**Scale/Scope**: One local target, eight connection states, one active connection generation, one event stream, five shell routes, and a documented primitive catalog for #153 through #156

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Gate | Result |
| --- | --- | --- |
| I. Code Quality | Explicit goroutine ownership, bounded errors, documented public contracts, no root dependency pollution | PASS |
| II. Testing Standards | Manager behavior is test-first under race; frontend primitives and accessibility receive direct tests; no clock sleeps | PASS |
| III. User Experience Consistency | S060 language, status, focus, appearance, and target rules become shared contracts | PASS |
| IV. Performance Requirements | Scheduler hot paths and persistence are untouched; connection timing and lifecycle counts are measured | PASS |
| V. Autonomous Execution | Full Spec Kit sequence, analysis gate, local verification, review branch, PR, and bounded review rounds | PASS |

No constitutional deviation is required. The user explicitly authorizes immediate branch publication after the local done gate, overriding only the ordinary pre-push halt timing, not any quality or review requirement.

## Project Structure

### Documentation

```text
specs/061-wails-shell-connection/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── connection-contract.md
│   └── shell-component-contract.md
├── checklists/
│   ├── requirements.md
│   └── foundation.md
├── tasks.md
└── verification.md
```

### Source Code

```text
desktop/
├── app.go
├── app_test.go
├── main.go
├── go.mod
├── go.sum
├── wails.json
├── build/appicon.png
├── connection/
│   ├── manager.go
│   ├── manager_test.go
│   ├── model.go
│   └── local.go
└── frontend/
    ├── e2e/shell.spec.ts
    ├── public/
    ├── src/
    │   ├── components/
    │   ├── connection/
    │   ├── App.tsx
    │   ├── App.test.tsx
    │   ├── accessibility.test.tsx
    │   ├── main.tsx
    │   └── styles.css
    ├── package.json
    ├── package-lock.json
    ├── playwright.config.ts
    ├── tsconfig.json
    └── vite.config.ts

.github/workflows/ci.yml
scripts/automation-check.sh
test/scripts/automation-check_test.sh
```

**Structure Decision**: Keep `desktop/` as a nested module with an explicit local replacement for the repository root. This isolates Wails and frontend dependencies from the shipping daemon, CLI, Fyne application, and installer during staged migration, while the module path remains within the repository import boundary and can use the protected local client. The S060 proof remains immutable historical evidence rather than becoming a mutable production directory.

## Phase 0: Research Decisions

See [research.md](research.md). The decisive choices are a separate production module, an actor-style single-owner connection manager, a scheduler seam for deterministic retry tests, stable sanitized bridge models, bounded version compatibility, and component-level contracts derived from S060 rather than copied page markup.

## Phase 1: Design and Contracts

See [data-model.md](data-model.md), [connection-contract.md](contracts/connection-contract.md), [shell-component-contract.md](contracts/shell-component-contract.md), and [quickstart.md](quickstart.md).

The connection manager owns one root context and one loop goroutine. Each attempt receives a monotonic generation and its own child context. Health negotiation uses a two-second timeout. A successful compatible negotiation starts exactly one event stream. Transport failure maps to safe state, waits on an injected scheduler using 250 millisecond, one second, then five second bounded delays, and can be interrupted by manual retry or shutdown. Event callbacks and results are accepted only for the active generation. Access denial and incompatibility wait for manual action instead of generating noisy automatic retries.

The Wails facade exposes current connection state and manual retry, translates manager notifications to one `desktop:event` channel, and never exposes endpoint names, credentials, transport errors, or backend objects. The frontend connection store consumes that contract and feeds a reusable Shell. Placeholder route panels prove the frame and primitives without implementing feature workflows.

## Phase 2: Implementation Strategy

1. Establish the nested production module and copy only approved local brand assets and exact dependency baselines.
2. Write manager state, transition, stale-generation, retry, safe-error, and shutdown tests before the manager implementation.
3. Add the real protected local adapter and Wails facade with bridge tests.
4. Write component, store, accessibility, responsive, and offline-asset tests before building the shared primitives and shell.
5. Add production three-platform builds and browser contracts to the existing fail-closed CI and automation policy.
6. Run nested and canonical verification, record evidence, update lifecycle inventory and changelog, and publish under the user's explicit authorization.

## Post-Design Constitution Check

| Principle | Evidence | Result |
| --- | --- | --- |
| I. Code Quality | One manager owns concurrency; adapters and bridge remain small; errors are categorized and sanitized | PASS |
| II. Testing Standards | Every state and concurrent lifecycle edge has deterministic race coverage before implementation | PASS |
| III. User Experience Consistency | Shared primitives encode the approved focus, status, language, appearance, and responsive rules | PASS |
| IV. Performance Requirements | No scheduler hot path changes; connection deadlines and leak cycles are explicit done gates | PASS |
| V. Autonomous Execution | Requirements, plan, tasks, analysis, implementation, verification, publication, and reviews remain chronological | PASS |

## Complexity Tracking

| Choice | Why Required | Simpler Alternative Rejected Because |
| --- | --- | --- |
| Nested `desktop/` module | Isolates migration dependencies from shipping root binaries until #157 | Adding Wails to the root module makes non-shipping migration dependencies part of every root dependency operation |
| Single-owner manager loop | Prevents overlapping retries and event streams and makes stale generations rejectable | Independent goroutines and locks reproduce the coordination defects already identified in the Fyne connection path |
| Injected retry scheduler | Proves exact retry behavior without wall-clock sleeps | Real timers make recovery tests slow and nondeterministic |
