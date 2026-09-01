# Implementation Plan: IPC Access-Denied Recovery

**Branch**: `codex/036-ipc-access-denied-recovery` | **Date**: 2026-08-31 |
**Spec**: [spec.md](spec.md)

## Summary

Introduce typed client failure classification, then replace independent modal
reporting with one application-owned connection incident rendered above the
existing tabs. A single cancelable reconnect coordinator provides immediate
Retry and bounded 2-to-30-second background backoff. Windows-only, read-only
diagnostics determine whether local-group membership exists but the process
token is stale; no authorization or installer policy changes.

## Technical Context

**Language/Version**: Go 1.25, Markdown, PowerShell 7 for native evidence

**Primary Dependencies**: Go standard library, Fyne v2, `golang.org/x/sys/windows`

**Storage**: In-memory incident state only; Markdown verification evidence

**Testing**: Go unit/headless GUI tests, Windows build-tag tests, native Windows walkthrough, canonical eight-gate suite

**Target Platform**: Cross-platform desktop GUI with Windows-specific named-pipe diagnosis

**Project Type**: Go desktop application and shared local API client

**Performance Goals**: One reconnect loop, no unbounded goroutines, 2-second initial and 30-second maximum delay

**Constraints**: Preserve named-pipe ACLs; no elevation, group mutation, service mutation, modal startup errors, real-time sleeps in deterministic tests, or new dependency

**Scale/Scope**: One process-wide incident, current startup refreshers and event stream, five failure classes

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | Typed wrapped errors, mutex-owned incident state, cancelable loop, documented lifecycle. |
| II. Testing Standards | PASS | Classification and concurrent three-source regression tests are written first; race and GUI gates remain mandatory. |
| III. UX Consistency | PASS | Accurate actionable copy, reachable frame, Retry and Exit, no OK-only storm. |
| IV. Performance | PASS | One bounded backoff loop replaces repeated refresh fan-out; no hot scheduling path changes. |
| V. Autonomous Execution | PASS | Full Spec Kit sequence, analyze gate, eight gates, review branch, authorized PR. |

### Post-design re-check

All gates remain PASS. `docs/INSTALL-windows.md` is pinned and changes because
the issue requires a fresh-install diagnostic walkthrough, so the dated
decision is recorded in `CHANGELOG.md`. The IPC listener and ACL implementation
remain untouched. Native Windows proves installed service, group, account,
token, pipe, and guidance behavior; deterministic testing proves the authorized
recovery transition without requiring the operator to terminate the session.

## Architecture and Decision Log

### Classify at the shared client boundary

Wrap transport failures in a typed error carrying a stable category and the
original cause. GUI string matching was rejected as fragile and duplicated;
changing IPC dial signatures was rejected as unnecessary security-boundary churn.
API `StatusError` remains distinct.

### Own one incident in the GUI

The application owns a mutex-protected incident snapshot and renders one
persistent panel above tabs. Every transport source reports through that
coordinator. Repeated reports update the snapshot rather than allocate dialogs.
Unrelated operation errors still use existing modals after connectivity is healthy.

### Coordinate retry rather than recursively refresh

The stream goroutine is the sole background reconnect loop. It attempts the
stream, records failures, waits with cancelable exponential backoff, and on
successful reconnection clears the incident and performs one coordinated
refresh. Retry interrupts the wait and requests an immediate attempt. Per-tab
retries and fixed polling were rejected because they reproduce amplification.

### Diagnose Windows state without mutation

A small build-tagged diagnostic reports service state, group existence,
account membership, and token membership as verified values or unknown. Only
verified membership plus missing token SID yields sign-out guidance. Diagnostic
failure never broadens access or claims a root cause.

## Project Structure

```text
internal/api/client/{client,methods,events,errors}_test.go
internal/api/client/errors.go
gui/{app,connection,connection_test}.go
gui/access_diagnosis_{windows,other}.go
docs/INSTALL-windows.md
specs/036-ipc-access-denied-recovery/
CHANGELOG.md
CLAUDE.md
specs/README.md
```

**Structure Decision**: Extend the existing shared client and GUI packages,
keeping OS inspection behind build tags. No new package or dependency is needed.

## Complexity Tracking

No constitution violation requires justification.
