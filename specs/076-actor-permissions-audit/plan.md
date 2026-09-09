# Implementation Plan: Actor Permissions and Management Audit

**Branch**: `codex/076-actor-permissions-audit` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/076-actor-permissions-audit/spec.md`

## Summary

Resolve issue #167 by adding a durable actor model, one shared fail-closed operation catalog, per-request authorization, intent-first management audit with bounded retention, and actor plus audit administration across the local API, typed client, and CLI. Preserve frictionless IPC use through a protected built-in local actor and avoid implementing the remote transport or credentials owned by #168 and #169.

## Technical Context

**Language/Version**: Go 1.25 module baseline with repository canonical verification

**Primary Dependencies**: Go standard library, `modernc.org/sqlite`, existing `github.com/google/uuid`, Cobra CLI, existing local IPC transport

**Storage**: SQLite schema migration v17 with actors and audit events

**Testing**: Go unit and integration tests, migration tests, API and CLI tests, specification lifecycle gate, canonical `bash scripts/verify.sh all`

**Target Platform**: Linux and Windows daemon and CLI over the existing protected local transport

**Performance Goals**: One bounded actor lookup per request, indexed bounded audit queries, no unbounded history, no new polling or listener

**Constraints**: Forward-only migration, fail-closed authorization, durable pre-side-effect audit intent, 10,000-event and 90-day retention, secret-free records, no new local authentication or remote access

**Scale/Scope**: Issue #167, two persisted entities, a shared catalog for all current management routes, six local API operations, and two CLI command families

## Constitution Check

### Pre-research gate

- **I. Code Quality**: PASS. Domain validation, persistence, authorization metadata, request enforcement, and user surfaces have explicit boundaries.
- **II. Testing Standards**: PASS. Migration, lifecycle, retention, authorization, catalog completeness, middleware, client, and CLI tests precede implementation.
- **III. User Experience Consistency**: PASS. Current local workflows remain credential-free, while new commands share existing human and JSON conventions.
- **IV. Performance Requirements**: PASS. Per-request lookups are indexed and audit history is bounded by age and count.
- **V. Autonomous Build-Phase Execution**: PASS. Work traces to #167, follows the complete spec-kit lifecycle, and has explicit push and pull-request authorization.
- **Engineering constraints**: PASS. No dependency, network listener, credential system, or general policy language is added.

### Post-design gate

- **I. Code Quality**: PASS. One operation catalog is authoritative and both actor and audit models reject invalid states.
- **II. Testing Standards**: PASS. Every success criterion maps to focused automated evidence plus canonical repository verification.
- **III. User Experience Consistency**: PASS. Local API, typed client, and CLI expose the same lifecycle and audit vocabulary.
- **IV. Performance Requirements**: PASS. Audit pruning occurs in the intent transaction and queries have deterministic bounded ordering.
- **V. Autonomous Build-Phase Execution**: PASS. S076 completes #167 without absorbing #168 or #169.

## Project Structure

```text
internal/domain/access.go
internal/authorization/catalog.go
internal/store/access.go
internal/store/store.go
internal/api/server/access.go
internal/api/server/server.go
internal/api/client/access.go
internal/cli/access.go
docs/access-control.md
specs/076-actor-permissions-audit/
```

**Structure Decision**: Keep portable authority types in domain, fail-closed operation metadata in a dedicated authorization package, durable state in store, transport enforcement in server middleware, and operator workflows in existing client and CLI packages.

## Implementation Strategy

1. Write domain, migration, actor lifecycle, audit lifecycle, retention, and filtering tests before adding schema v17 and store methods.
2. Write authorization hierarchy and catalog completeness tests before implementing the shared catalog and authorizer.
3. Write server middleware and API contract tests before enforcing the catalog and adding actor and audit routes.
4. Write typed-client and CLI tests before implementing actor administration and audit inspection commands.
5. Document the trust boundary, capability matrix, storage and redaction contract, migration, retention, export, and follow-on integration.
6. Run focused, canonical, formatting, encoding, lifecycle, and scope verification before authorized publication.

## Complexity Tracking

No constitutional violation or justified complexity exception is present.
