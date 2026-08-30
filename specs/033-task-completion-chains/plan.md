# Implementation Plan: Task-Completion Chains

**Branch**: `codex/033-task-completion-chains` | **Date**: 2026-08-30 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/033-task-completion-chains/spec.md`

## Summary

Add acyclic completion relationships as an additional task trigger, backed by a durable at-least-once delivery ledger. Record each terminal run and its matching deliveries in one transaction, claim pending deliveries through the bounded engine, correlate target history to the source run, and recover interrupted claims on daemon start. Expose the complete relationship lifecycle through the existing local API/client/CLI, a dedicated desktop Chains view, live events, and consistent documentation. Reuse current overlap, eligibility, worker, history, IPC, and event-stream boundaries rather than restoring the removed legacy trigger subsystem.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Standard library plus existing SQLite, UUID, Cobra, Fyne, broker, and local API dependencies. No new dependency.

**Storage**: Existing embedded SQLite store, forward-only schema migration v9, single-writer connection.

**Testing**: Go unit, API contract, CLI, headless Fyne, store migration/reopen, engine integration, race, coverage, benchmark, documentation, and automation gates.

**Target Platform**: Linux, macOS, and Windows daemon/CLI/desktop builds over existing local IPC.

**Project Type**: One Go repository containing daemon, CLI, native desktop GUI, local JSON API, and documentation.

**Performance Goals**: Preserve the p99 dispatch budget below 100 ms; keep 100-way fan-out bounded by the configured worker pool and indexed lookups.

**Constraints**: Deterministic injected clock, no polling goroutine, defined goroutine shutdown, no new network exposure, forward-compatible databases, honest at-least-once replay window, and core coverage at least 80 percent.

**Scale/Scope**: Local single-machine graphs of at least 100 tasks, 100-way fan-out, finite acyclic cascades, and one S033 slice covering issues #72 through #77.

## Constitution Check

*GATE: PASS before research and PASS after design.*

| Principle | Design compliance |
| --- | --- |
| I. Code Quality | Store transactions own atomicity; the engine retains one bounded worker lifecycle; errors are wrapped and surfaced. |
| II. Testing Standards | Tests precede migration, transactions, cycle detection, replay, overlap, API/CLI/GUI, event folding, and recovery implementation. Time remains injected. |
| III. User Experience | All surfaces share `chain`, `completion`, `source_task_id`, and `source_run_id`, with actionable field errors and compatible text/JSON output. |
| IV. Performance | Indexed source/outcome and delivery-state queries avoid history scans. Claims are bounded and execution keeps the worker semaphore. |
| V. Autonomous Execution | The installed hyphenated Spec-Kit sequence runs on the review branch with analyze, all eight gates, local commit, and pre-publication halt. |

No constitution exception or complexity waiver is required.

## Architecture and Delivery Flow

1. A terminal executor run reaches the engine with optional incoming-delivery correlation.
2. One store transaction inserts the run, completes its incoming delivery, and inserts unique pending deliveries for matching outgoing chains.
3. The engine claims pending deliveries in bounded batches and resolves current target eligibility before dispatch.
4. The overlap path accepts an internal origin carrying trigger and source correlation. Actual completion closes the delivery; queued bookkeeping does not.
5. Startup returns claimed unfinished deliveries to pending with replay evidence. Completed or resolved deliveries never return.
6. Relationship mutations validate task existence, uniqueness, and the prospective graph before commit, then publish one live chain event.

## Project Structure

### Documentation (this feature)

```text
specs/033-task-completion-chains/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── completion-delivery.md
│   └── management.md
├── checklists/
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
internal/domain/domain.go
internal/store/store.go
internal/store/chains.go
internal/engine/engine.go
internal/engine/overlap.go
internal/api/server/chains.go
internal/api/client/methods.go
internal/cli/chain.go
internal/events/broker.go
gui/app.go
gui/chains.go
gui/viewmodel/viewmodel.go
docs/cli.md
docs/api.md
docs/gui-fields.md
docs/architecture.md
```

**Structure Decision**: Extend existing layers in place. A standalone trigger service would duplicate eligibility, overlap, event, and lifecycle behavior and is rejected.

## Complexity Tracking

No constitution violations require justification.

## Post-Design Constitution Recheck

PASS. The design is bounded, deterministic, local-only, dependency-neutral, forward-migrating, race-testable, and compatible with existing worker, overlap, IPC, and event-stream contracts. The architecture-required at-least-once replay window is recorded in the specification, issues, research, and changelog decision.
