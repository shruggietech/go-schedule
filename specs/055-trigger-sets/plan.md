# Implementation Plan: Trigger Sets

**Branch**: `codex/055-trigger-sets` | **Date**: 2026-09-05 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/055-trigger-sets/spec.md`

## Summary

Extend the S054 external-trigger lifecycle with persistent Trigger Sets that atomically create and administer 1 through 99 ordinary triggers targeting one task. The slice adds ordered bulk secret output, transactional set operations, CLI and desktop administration, and migration, redaction, rollback, and performance coverage while preserving standalone triggers and the existing dispatcher.

## Technical Context

**Language/Version**: Go 1.25

**Primary Dependencies**: Go standard library, Cobra CLI, Fyne desktop toolkit, modernc SQLite

**Storage**: Existing local SQLite database with forward-only schema migration v13

**Testing**: Go store, API, client, CLI, event, GUI, migration, transaction-failure, benchmark, race, format, documentation, and automation verification through `scripts/verify.sh all`

**Target Platform**: Windows, Linux, and macOS desktop and headless environments

**Project Type**: Single Go desktop application with CLI, local IPC API, and scheduler engine

**Performance Goals**: Every maximum-size 99-member set mutation completes below one second under nominal local load

**Constraints**: Existing local IPC only, no new dependency, no raw keys in ordinary responses or telemetry, all broad mutations transactional, standalone trigger compatibility, and no bypass of task readiness or dispatch controls

**Scale/Scope**: Up to 99 members per set, permanent sparse positions after deletion, and one GitHub issue delivered end to end

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- PASS: One nullable membership relation extends ordinary triggers instead of creating a second dispatch model.
- PASS: SQLite transactions own every broad mutation and automatic-source side effect, with migration and injected-failure coverage required.
- PASS: Secret-bearing create, reveal, and rotate responses remain explicit while ordinary API, event, log, Activity, and history surfaces remain redacted.
- PASS: CLI and GUI consume the same local API and stable error contract.
- PASS: No network listener, goroutine, filesystem watcher, remote capability, or dependency is introduced.
- PASS: Maximum-size operations have a measured one-second budget and deterministic benchmark coverage.
- PASS: The complete local CI-parity aggregate remains mandatory before publication.

The Phase 1 design preserves every gate above. No constitution exception requires complexity tracking.

## Project Structure

### Documentation (this feature)

```text
specs/055-trigger-sets/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── trigger-set-contract.md
├── checklists/
│   ├── requirements.md
│   └── security-atomicity.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
internal/
├── domain/             # Trigger Set identity and member metadata
├── store/              # Migration v13 and transactional set persistence
├── api/
│   ├── server/         # Versioned Trigger Set handlers and redacted responses
│   └── client/         # Typed local API methods
├── events/             # Set lifecycle event identity without secrets
└── cli/                # Nested trigger set lifecycle commands and ordered output

gui/
├── app.go              # Backend capability and live refresh integration
├── triggers.go         # Membership column and selected-set administration
└── *_test.go           # Headless set lifecycle and redaction contracts

docs/                   # CLI, GUI, API, and architecture guidance
```

**Structure Decision**: Extend the current trigger packages in place. Trigger Set operations coordinate existing trigger records through the store, while firing remains unaware of set membership and continues through the S054 dispatcher.

## Design Decisions

- Persist a Trigger Set row plus nullable `set_id` and `set_position` fields on ordinary triggers. A separate membership table would add joins and lifecycle states without supporting many-to-many membership, which the specification forbids.
- Keep positions permanent and sparse after deletion. Renumbering would mutate unaffected siblings and make copied caller identities unstable.
- Make every set-level mutation fully transactional. Precise partial-success objects are unnecessary when no partial state is permitted, and transaction failure returns the existing API error envelope.
- Reject individual target changes for set members. Allowing them would violate the set's single-target invariant; the set retarget operation updates all members and affected task readiness in one transaction.
- Publish one redacted Trigger Set lifecycle event per broad operation instead of one event per member. This prevents a 99-member action from causing 99 desktop refreshes while preserving live convergence.
- Use `gosched trigger set ...` as the CLI namespace. It extends the established trigger noun without introducing a competing top-level vocabulary.

## Complexity Tracking

No constitution violations or exceptional complexity are introduced.
