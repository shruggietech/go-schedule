# Implementation Plan: External Trigger Lifecycle

**Branch**: `codex/054-external-trigger-lifecycle` | **Date**: 2026-09-05 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/054-external-trigger-lifecycle/spec.md`

## Summary

Deliver the first complete external-trigger capability from issue #17 by adding persistent trigger records, secure key lifecycle operations, one-shot dispatch through the existing scheduler, CLI and local API contracts, run provenance, and a desktop Triggers view. The implementation bundles issues #132 and #133 while deliberately excluding trigger sets and filesystem watchers from this slice.

## Technical Context

**Language/Version**: Go 1.25

**Primary Dependencies**: Go standard library, Cobra CLI, Fyne desktop toolkit, modernc SQLite

**Storage**: Existing local SQLite database with schema migration v12

**Testing**: Go unit, integration, API, CLI, engine, store, GUI, race, vet, format, documentation, and packaging verification through `scripts/verify.sh all`

**Target Platform**: Windows, Linux, and macOS desktop and headless environments

**Project Type**: Single Go desktop application with CLI, local IPC API, and scheduler engine

**Performance Goals**: Trigger dispatch decision p99 below 100 ms under 100 concurrent local fire requests

**Constraints**: Local-only authenticated IPC, no new listener or background service, no raw keys in ordinary responses or telemetry, and no bypass of existing overlap or worker limits

**Scale/Scope**: Up to 1,000 persisted triggers, 100 concurrent fire requests, two linked GitHub issues, and one end-to-end desktop plus CLI workflow

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- PASS: The feature remains local-first and does not introduce a remote network listener.
- PASS: Existing scheduler dispatch, overlap, concurrency, history, and event paths remain authoritative.
- PASS: Sensitive trigger keys are explicitly excluded from ordinary API responses, logs, history, and events.
- PASS: Persistence changes use a forward-only SQLite migration with compatibility coverage.
- PASS: CLI and GUI behavior share the same local API contracts rather than duplicating business rules.
- PASS: The plan includes focused unit and integration tests plus the repository-wide verification gate.
- PASS: The slice is bounded to single triggers and their GUI lifecycle; trigger sets and filesystem watchers remain separately tracked.

The Phase 1 design preserves every gate above. No exceptions require complexity tracking.

## Project Structure

### Documentation (this feature)

```text
specs/054-external-trigger-lifecycle/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── trigger-contract.md
├── checklists/
│   ├── requirements.md
│   └── security-ux.md
└── tasks.md
```

### Source Code (repository root)

```text
internal/
├── domain/             # Trigger model and run provenance
├── store/              # Migration and trigger persistence
├── task/               # Automatic-source readiness rules
├── engine/             # Key validation and scheduler dispatch
├── api/
│   ├── server/         # Trigger HTTP handlers and events
│   └── client/         # Typed local API client
├── events/             # Trigger lifecycle event payloads
└── cli/                # Trigger lifecycle and fire commands

gui/
├── app.go              # Navigation and refresh integration
├── triggers.go         # Trigger management view and dialogs
└── viewmodel/          # Trigger snapshots and event updates

docs/                   # User-facing CLI, GUI, and architecture guidance
specs/054-external-trigger-lifecycle/ # Specification and verification evidence
```

**Structure Decision**: Extend the existing single-module architecture in place. Trigger business rules live below both interfaces, and the GUI consumes the same local API used by the CLI.

## Complexity Tracking

No constitution violations or exceptional complexity are introduced.
