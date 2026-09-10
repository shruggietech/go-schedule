# Implementation Plan: Remote Connection Profiles and Target-Safe Clients

**Branch**: `codex/078-remote-connection-profiles` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/078-remote-connection-profiles/spec.md`

## Summary

Deliver one user-scoped, secret-free connection profile system shared by the Wails desktop and `gosched`; extend the shared API client with immutable local and authenticated HTTPS constructors; persist desktop selection; add target-aware desktop connection administration and capability gating; add explicit CLI profile and endpoint targeting; and publish safe JSON-client guidance. The slice completes #170 and #171 without absorbing the automatic recovery work in #172.

## Technical Context

**Language/Version**: Go 1.25.0; TypeScript 5.9 and React 19 in the existing Wails frontend

**Primary Dependencies**: Standard library JSON, filesystem, HTTP, TLS, and synchronization; existing `github.com/zalando/go-keyring`; existing Cobra; existing Wails bindings

**Storage**: Versioned JSON at the user configuration root for non-secret profiles; operating-system credential store for bearer values; existing daemon SQLite remains unchanged

**Testing**: Go unit and race tests, frontend Vitest and Testing Library, Wails native build, secret-canary checks, documentation integrity, and `sh scripts/verify.sh all`

**Target Platform**: Linux, macOS, and Windows desktop and CLI clients connecting through local IPC or TLS 1.3 HTTPS

**Project Type**: Existing Go daemon and CLI with a Wails desktop frontend

**Performance Goals**: Profile resolution and selection complete within one ordinary local filesystem and keyring operation; no scheduler hot path changes; bounded two-second initial desktop health attempt and existing ten-second CLI request deadline

**Constraints**: Local IPC stays the default; no bearer in JSON or command arguments; no insecure TLS or redirects; remote identity is pinned; Markdown uses UTF-8 without BOM, no Unicode em dash, and no width wrapping

**Scale/Scope**: Up to 100 user-scoped profiles, one active desktop target, one target per CLI invocation, and one connection generation at a time

## Constitution Check

### I. Code Quality

PASS. Profile, selection, remote transport, and desktop orchestration have explicit ownership boundaries. No new dependency is required.

### II. Testing Standards

PASS. Tests precede behavior for atomic persistence, concurrent access, secret exclusion, TLS trust, identity mismatch, profile lifecycle, target switching, UI accessibility, and CLI compatibility. Full race and coverage gates remain mandatory.

### III. User Experience Consistency

PASS. Existing local CLI behavior remains unchanged. Remote human diagnostics use stderr, JSON stdout stays parseable, and the desktop maintains one consistent target vocabulary.

### IV. Performance Requirements

PASS. Scheduler dispatch paths are untouched. Profile operations are bounded and connection work uses existing timeouts.

### V. Autonomous Build-Phase Execution

PASS. S078 is traceable to #170 and #171, runs through every spec-kit phase, uses a review branch, and will publish only under the user's explicit authorization.

## Project Structure

### Documentation

```text
specs/078-remote-connection-profiles/
├── checklists/
├── contracts/
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
internal/clientprofile/          # shared profile model, validation, and atomic user store
internal/api/client/             # immutable local and remote client transports
internal/cli/                    # target flags and profile lifecycle commands
desktop/connection/              # remote backend and switchable connection generations
desktop/remotepairing/           # pairing plus durable profile handoff
desktop/connections/             # desktop profile lifecycle service
desktop/frontend/src/connection/ # target-aware bridge and store
desktop/frontend/src/settings/   # connection list and lifecycle UI
docs/                            # CLI, JSON API, and remote-access guidance
```

**Structure Decision**: Put profile persistence under `internal/clientprofile` because desktop and CLI share the same interactive-user data and validation contract. Keep bearer storage in `internal/clientsecret`. Extend the existing concrete API client instead of duplicating every operation behind a second generated facade.

## Design Phases

### Phase 0: Profile and transport foundation

Create immutable client target constructors, versioned profile persistence, canonical endpoint and certificate validation, file locking, atomic replacement, and secret-exclusion tests.

### Phase 1: Desktop lifecycle

Persist profiles after pairing, add remote connection negotiation and target switching, expose profile lifecycle through the Wails bridge, and make target identity plus disabled capability explanations available throughout the shell and Connections view.

### Phase 2: CLI and JSON clients

Add invocation-scoped targeting flags, profile pair/list/show/rename/remove operations, safe remote diagnostics, and executable documentation for JSON enrollment and authenticated requests.

### Phase 3: Closure

Run cross-artifact analysis, focused security and race tests, full CI parity, generated binding checks, secret and encoding scans, then record verification.

## Complexity / Deviation

No constitutional deviation is planned. A shared mutable transport was rejected because a target switch could redirect an in-flight mutation. The desktop instead replaces an immutable target client only after canceling the prior connection generation, and daemon-backed services resolve the current client through a synchronized target provider.

## Post-Design Constitution Check

PASS. The design keeps local defaults, secrets, concurrency ownership, diagnostics, and review evidence explicit. It adds no dependency, daemon schema, remote route, scheduler hot-path work, or release claim.
