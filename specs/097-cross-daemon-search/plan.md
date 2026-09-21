# Implementation Plan: Cross-Daemon Search and Target-Safe Actions

**Branch**: `codex/097-cross-daemon-search` | **Date**: 2026-09-21 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/097-cross-daemon-search/spec.md`

## Summary

Add an explicit Search destination that fans one bounded query across This computer and every saved daemon profile, displays progressive source-labeled matches and target failures, and executes confirmed acknowledge, enable, disable, or run-now requests only after reconnecting and revalidating the exact daemon and object. Each daemon exposes one additive Observe-authorized search endpoint backed by bounded store queries. The desktop owns fan-out, capability and authority presentation, exact-target revalidation, independent per-object outcomes, and an accessible React workflow. No result cache, cluster state, distributed transaction, or automatic rollback is introduced.

## Technical Context

**Language/Version**: Go 1.26, TypeScript 5.6, React 19

**Primary Dependencies**: Go standard library, modernc SQLite, existing local IPC and identity-pinned HTTPS clients, Wails desktop bridge, React, Vite

**Storage**: Existing SQLite task, group, schedule, run, and alert tables; existing JSON connection profiles and native credential storage; no new schema or persisted search state

**Testing**: Go unit and API integration tests, race detector, Vitest and Testing Library, Playwright accessibility and scale checks, repository CI-parity scripts

**Target Platform**: Windows and Linux daemon, Windows desktop application, local IPC and enabled remote HTTPS transports

**Project Type**: Go daemon and API with a React desktop frontend

**Performance Goals**: Bound each daemon response to 50 results, cap fan-out at eight targets, classify each target within three seconds of its worker start, stream progressive results, and keep 100 registered profiles usable at 800 by 600 and 200 percent zoom

**Constraints**: Explicit query submission, Observe-authorized secret-free reads, Operate-authorized audited mutations, exact registration and daemon identity revalidation, no persisted result cache, no cluster semantics, no distributed rollback, and no new third-party dependency

**Scale/Scope**: This computer plus up to 100 saved profiles; five result kinds; four mutation actions; one daemon search endpoint; one desktop orchestration service; one new primary route; focused API, service, frontend, documentation, and regression coverage

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **Code Quality**: Search fan-out has generation-scoped cancellation, bounded worker ownership, child deadlines, and deterministic ordering. Public Go types and functions receive intent-focused comments. The desktop search service centralizes target construction and identity-safe mutation instead of duplicating it in the frontend.
- **Testing Standards**: Store search, API authorization and redaction, transport parity, fan-out cancellation, duplicate names, authority differences, identity changes, object deletion, partial mutation failure, and accessible frontend workflows receive deterministic tests written before implementation. No assertion depends on wall-clock sleeps.
- **User Experience Consistency**: Connection states reuse the typed taxonomy, timestamps remain RFC 3339, errors name one target and recovery action, confirmations name every daemon and object, and all controls reuse the established shell, dialog, notice, focus, and responsive patterns.
- **Performance Requirements**: SQL applies normalized matching and hard limits before returning data. Target and result cardinality are bounded. No scheduler dispatch hot path changes, persisted index, or benchmark-sensitive code is introduced.
- **Autonomous Build-Phase Execution**: S097 is traceable to issue #183 and follows the complete spec-kit sequence on a review branch. The operator explicitly authorized push, pull-request publication, review fixes, and up to two review rounds in the kickoff, satisfying the publication halt in advance.
- **Engineering Constraints**: The API addition is backward-compatible, available through both transports, validates all query inputs, preserves pinned identity, and returns only allowlisted observation fields. Mutations reuse existing audited endpoints and retry-safety behavior.
- **Outcome-based issue completion**: Issue #183 closes when the functional search and target-safe action outcomes are delivered. Local checks, CI, and reviews remain engineering gates rather than issue acceptance criteria.

No constitutional deviation is required.

## Project Structure

### Documentation (this feature)

```text
specs/097-cross-daemon-search/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── daemon-search.md
│   └── desktop-actions.md
├── checklists/
│   ├── requirements.md
│   └── target-safety.md
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
├── remote/
└── store/

api/openapi/
└── remote-v1.yaml

desktop/
├── search/
├── frontend/src/
│   ├── components/
│   ├── connection/
│   └── search/
├── app.go
└── main.go

docs/
├── api.md
└── remote-access.md
```

**Structure Decision**: Extend the existing daemon domain, store, authorization catalog, HTTP API, generated remote contract, and client. Add one `desktop/search` package beside `desktop/systems` because search has independent fan-out state and mutation orchestration. Add one feature-folder frontend route. Reuse connection profile, credential, backend health, exact-profile selection, operation endpoints, shared UI primitives, and event patterns. Do not alter the global selected connection during background search or direct search actions.

## Design Decisions

### One bounded daemon-owned search endpoint

Each daemon searches its authoritative tables and computes matching schedule occurrences from matching scheduled tasks. The response is limited before crossing transport boundaries and contains only safe labels, identifiers, states, timestamps, readiness, and action hints. This avoids downloading full workspaces or histories to every desktop query.

### Independent clients and immutable registration keys

The search service snapshots This computer and saved profiles, loads each credential through the native secret store, and constructs temporary identity-pinned clients. Background queries and direct actions do not switch the global router. Each result carries the immutable registration key and observed daemon identity used to locate the exact client again.

### Generation-scoped eight-worker fan-out

A new search cancels the previous generation and snapshots the current registrations. At most eight workers issue one query per target with a three-second child deadline. Progressive snapshots replace only their registration's observation. Removed or edited registrations and superseded generations discard late results.

### Revalidation immediately before every mutation

The action service reconstructs the recorded registration, performs health, manifest, pinned-identity, and authority checks, reloads the exact task or alert, and compares kind and current state. Only then does it call the existing enable, disable, run-now, or acknowledgement endpoint. No display label participates in routing.

### Compatible selection and independent outcomes

The frontend allows one action at a time and requires every selected result to advertise that action. Confirmation groups results by target and names authority. Execution is sequential within a target and bounded across targets. Every object receives its own accepted, rejected, unavailable, or uncertain outcome. There is no cross-target rollback.

### Open uses the existing selected-connection path

Opening a result remains a navigation action. It selects the exact registration through the established connection manager, waits for a connected snapshot with the expected daemon identity, and then supplies task or record context to Tasks, Schedule, or Activity. A mismatch leaves Search visible.

## Complexity Tracking

No constitution violations or complexity exceptions are required.
