# Implementation Plan: Resilient Remote Connections

**Branch**: `codex/079-remote-connection-resilience` | **Date**: 2026-09-10 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/079-remote-connection-resilience/spec.md`

## Summary

Extend the existing single-owner desktop connection manager so remote profiles automatically recover safe negotiation and event subscriptions with bounded jittered backoff. Add explicit freshness and retry metadata to the bridge contract, classify trust, identity, authorization, compatibility, timeout, and reachability failures, persist successful contact timestamps, and mark retained frontend workspaces stale while all daemon-backed mutations remain single-attempt and disabled outside a current connection.

## Technical Context

**Language/Version**: Go 1.25.0; TypeScript 5.9 and React 19 in the existing Wails frontend

**Primary Dependencies**: Go standard library HTTP, TLS, x509, synchronization, contexts, and pseudorandom source; existing Wails bindings and React hooks

**Storage**: Existing versioned user-scoped connection profile JSON gains updates to its existing `last_successful_at` field; daemon SQLite is unchanged

**Testing**: Go unit and race tests, deterministic fake retry scheduling and jitter, `httptest` transport tests, frontend Vitest and Testing Library, Chromium accessibility checks, native Wails build, and `sh scripts/verify.sh all`

**Target Platform**: Linux, macOS, and Windows desktop clients using local IPC or authenticated TLS 1.3 HTTPS

**Project Type**: Existing Go daemon and CLI with a Wails desktop frontend

**Performance Goals**: First retry scheduled in 250 to 750 milliseconds, exponential growth capped at 30 seconds, one health attempt and one event stream per connection generation, no scheduler hot-path changes

**Constraints**: No automatic mutation replay, offline queue, durable event replay, insecure TLS, or target fallback; local IPC semantics remain unchanged; Markdown is UTF-8 without BOM, contains no Unicode em dash, and is not width-wrapped

**Scale/Scope**: One selected desktop target, one active generation, one retry timer, one stream, up to 100 profiles, and retained last-complete state in each mounted workspace

## Constitution Check

### I. Code Quality

PASS. Connection ownership remains in one manager goroutine, retry policy is isolated behind deterministic seams, and all new errors carry safe typed classification.

### II. Testing Standards

PASS. Tests precede retry, classification, freshness, persistence, mutation uncertainty, cancellation, accessibility, and lifecycle changes. The race and core coverage gates remain mandatory.

### III. User Experience Consistency

PASS. The bridge gains one shared recovery vocabulary; retained data is marked stale; local controls remain usable; remote mutations fail once with actionable uncertainty instead of being replayed.

### IV. Performance Requirements

PASS. Retry concurrency is bounded to one generation and one timer. Scheduler dispatch is untouched, so no scheduler benchmark change is required.

### V. Autonomous Build-Phase Execution

PASS. S079 is traceable to #172, runs every spec-kit phase, uses a review branch, and has explicit authorization to publish one PR and its in-scope review fixes.

## Project Structure

### Documentation

```text
specs/079-remote-connection-resilience/
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
internal/api/client/             # transport classification and uncertain mutation errors
internal/clientprofile/          # last-success persistence
desktop/connection/              # retry policy, state machine, failure classification, persistence callback
desktop/connections/             # active profile last-success updates
desktop/frontend/src/connection/ # freshness-aware bridge model and store
desktop/frontend/src/components/ # global stale and retry context
desktop/frontend/src/settings/   # detailed recovery guidance
desktop/frontend/src/tasks/      # retained workspace and draft behavior
desktop/frontend/src/operations/ # retained schedule and activity behavior
docs/                            # remote recovery and non-replay guidance
```

**Structure Decision**: Keep retry ownership in `desktop/connection` because it already serializes health, stream, switch, retry, and shutdown generations. Keep the shared API client single-attempt for all methods and add only typed uncertainty classification for remote non-idempotent transport failures. Treat the event stream as an invalidation source and trigger existing authoritative workspace reloads after a recovered generation instead of building a replay log.

## Design Phases

### Phase 0: Typed transport and recovery contract

Add safe trust and mutation-uncertainty error categories, a deterministic retry policy seam, expanded connection snapshot fields, and target-scoped successful-contact notification.

### Phase 1: Single-owner remote recovery

Enable remote automatic retry only for transient failures, stop on terminal failures, publish scheduled retry metadata, reset retry progression after stable activity, cancel all stale generations, and persist successful contact for the active profile.

### Phase 2: Honest retained-state desktop

Expose stale and retry metadata through Wails, retain complete feature workspaces and task drafts, disable daemon-backed controls outside current connection, announce recovery accessibly, and force authoritative reloads before resumed events can trigger incremental refreshes.

### Phase 3: Closure

Document retry and mutation boundaries, run lifecycle and secret-canary tests, regenerate bindings through the native build, complete canonical verification, and record delivery evidence.

## Key Decisions

1. **Refresh instead of replay**: The existing SSE stream carries invalidation signals, not authoritative state. A full read after reconnection closes event gaps with less protocol and storage complexity than durable cursor replay.
2. **Manager-owned retry**: Per-request automatic retries could duplicate mutations or overlap target changes. Only the connection manager retries health and stream setup; feature service calls remain one attempt.
3. **Typed uncertain mutations**: Remote non-GET transport failures become a typed uncertainty that services can present consistently. HTTP error responses remain authoritative rejections.
4. **Terminal trust failures**: Certificate and identity changes stop retries. Silent trust replacement is rejected because it would defeat profile pinning.
5. **Deterministic jitter injection**: Production uses bounded jitter, while tests supply an exact source. This preserves realistic retry spreading without wall-clock-dependent assertions.

## Complexity / Deviation

No constitutional deviation is planned. The specification says live recovery uses a refresh-before-stream contract rather than a server-side resume cursor. This deliberately refines the earlier informal S078 phrase “resume cursors” because the shipped event API is an invalidation channel with authoritative read endpoints, and adding durable replay would be disproportionate to #172's explicit prohibition on elaborate offline synchronization.

## Post-Design Constitution Check

PASS. The design keeps connection ownership bounded, preserves local compatibility, adds deterministic tests for every safety boundary, exposes actionable state, avoids unbounded queues and replay, and changes no pinned process artifact.
