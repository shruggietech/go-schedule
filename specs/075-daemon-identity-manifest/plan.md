# Implementation Plan: Stable Daemon Identity and Capability Manifest

**Branch**: `codex/075-daemon-identity-manifest` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/075-daemon-identity-manifest/spec.md`

## Summary

Resolve issue #166 by persisting one UUID installation identity and editable display name in the daemon's SQLite store, exposing a privacy-bounded capability manifest through the protected local API, adding typed client and CLI operations for discovery, rename, and compare-and-confirm reset, and making the desktop consume daemon-owned target facts. Preserve the existing health contract and scheduler state while documenting restore and clone semantics.

## Technical Context

**Language/Version**: Go 1.25 module baseline with Go 1.27.1 canonical verification; Markdown and POSIX shell documentation gates

**Primary Dependencies**: Go standard library, `modernc.org/sqlite`, existing `github.com/google/uuid`, Cobra CLI, existing local IPC transport, existing Wails desktop connection package

**Storage**: SQLite schema migration v16 with a singleton `daemon_identity` row

**Testing**: Go unit and integration tests, migration tests, CLI and desktop connection tests, specification lifecycle gate, canonical `sh scripts/verify.sh all`

**Target Platform**: Linux and Windows daemon, CLI, and desktop clients over the existing protected local transport

**Project Type**: Existing Go daemon, local HTTP API, CLI, and Wails desktop application

**Performance Goals**: Constant-size singleton reads and writes; no additional polling, network listener, or unbounded manifest collection

**Constraints**: Forward-only schema migration; stable identity across restart, upgrade, and restore; atomic reset; no hostname, address, path, account, credential, command, or secret disclosure; health compatibility; no remote access or actor authorization

**Scale/Scope**: Issue #166 only, one singleton entity, three local API operations, three CLI operations, and one desktop connection integration

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Pre-research gate

- **I. Code Quality**: PASS. Identity validation, persistence, and API projection have distinct boundaries and use existing packages.
- **II. Testing Standards**: PASS. Migration, lifecycle, API, client, CLI, and desktop tests are authored before their implementations and canonical race verification remains mandatory.
- **III. User Experience Consistency**: PASS. The generic default avoids host disclosure, names remain editable, reset requires exact acknowledgement, and old health clients remain compatible.
- **IV. Performance Requirements**: PASS. Singleton database access and bounded static collections add no unbounded work.
- **V. Autonomous Build-Phase Execution**: PASS. Work traces to #166, follows the complete spec-kit lifecycle, uses a review branch and pull request, and has explicit publication authorization.
- **Engineering constraints**: PASS. No new dependency or remote listener is introduced, SQLite remains authoritative, errors fail closed, and Linux and Windows stay supported.

### Post-design gate

- **I. Code Quality**: PASS. The domain entity owns name validation, the store owns atomic identity lifecycle, and API projections own disclosure boundaries.
- **II. Testing Standards**: PASS. Every success criterion maps to focused automated evidence plus the full repository gate.
- **III. User Experience Consistency**: PASS. API, CLI, and desktop all consume the same daemon-owned identity and deterministic capabilities.
- **IV. Performance Requirements**: PASS. Manifest values are constant-size and deterministic, while reset is one transactional compare-and-swap operation.
- **V. Autonomous Build-Phase Execution**: PASS. S075 completes #166 without absorbing #167 authorization or later remote-access work.

## Project Structure

### Documentation

```text
docs/
├── api.md
├── architecture.md
└── daemon-identity.md

specs/075-daemon-identity-manifest/
├── checklists/
│   ├── identity-lifecycle.md
│   └── requirements.md
├── contracts/
│   └── local-manifest.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
├── tasks.md
└── verification.md
```

### Source Code

```text
internal/domain/
└── identity.go

internal/store/
├── identity.go
├── identity_test.go
├── migration_v16_test.go
└── store.go

internal/api/
├── client/
│   ├── client.go
│   └── client_test.go
└── server/
    ├── manifest.go
    ├── manifest_test.go
    └── server.go

internal/cli/
├── cli.go
├── daemon.go
└── daemon_test.go

desktop/connection/
├── local.go
├── local_test.go
├── manager.go
└── model.go
```

**Structure Decision**: Keep identity rules in the domain package, persistence and migration in the store, safe wire projections in the local API, user operations in the existing CLI, and target adoption in the desktop connection boundary. This extends established repository packages without adding a subsystem or dependency.

## Implementation Strategy

1. Write migration and store lifecycle tests first, demonstrate failure, then add schema v16 and the singleton identity store.
2. Write server and client contract tests first, then add the bounded manifest, rename, and compare-and-confirm reset endpoints while preserving health.
3. Write CLI and desktop integration tests first, then expose operator commands and replace desktop identity and capability placeholders.
4. Document clean install, upgrade, restore, clone, rename, and reset behavior; update changelog and lifecycle inventory.
5. Run focused tests, canonical verification, publication formatting, and repository integrity audits before authorized publication.

## Complexity Tracking

No constitutional violation or justified complexity exception is present.
