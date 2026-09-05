# Implementation Plan: Filesystem Watchers

**Branch**: `codex/056-filesystem-watchers` | **Date**: 2026-09-05 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/056-filesystem-watchers/spec.md` and GitHub issue #135.

## Summary

Add durable file and directory watcher definitions that feed stable matching files into the existing daemon's overlap-aware dispatcher. A daemon-owned watcher runtime uses one cross-platform native observer, explicitly registers recursive directories, coalesces event storms through an injected-clock state machine, exposes ephemeral health through the existing local API, and reloads atomically on lifecycle mutations. The CLI and desktop Triggers view provide complete administration, while schema v14 adds watcher configuration and durable watcher run provenance without replaying filesystem events.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Standard library, `github.com/fsnotify/fsnotify` v1.10.0, Cobra 1.10.2, Fyne 2.8.1, modernc.org/sqlite 1.57.0

**Storage**: Existing SQLite database with forward-only migration v14; runtime health and pending candidate state remain in memory

**Testing**: Go `testing`, injected project clock, fake observer tests, real temporary-filesystem integration tests, race detector, Fyne headless tests, cross-platform CI, and `scripts/verify.sh all`

**Target Platform**: Windows, Linux, and macOS daemon, CLI, and desktop application

**Project Type**: Local daemon with IPC API, CLI client, and desktop GUI

**Performance Goals**: Lifecycle and health reads below one second at 100 watchers; 100-trial write-storm and rename fidelity; existing dispatch p99 below 100 milliseconds after watcher acceptance

**Constraints**: No remote listener, no payload forwarding, no durable event replay, no link-directory traversal, no direct wall-clock use in the watcher state machine, no repeated identical health logs or events, and no matched path in durable run history

**Scale/Scope**: 100 configured watchers, arbitrary selected files bounded by native observer capacity and existing task overlap policies

## Constitution Check

*GATE: Passed before research and rechecked after design.*

| Principle | Design response | Status |
| --- | --- | --- |
| I. Code Quality | The runtime has one owned event loop, generation-scoped state, explicit observer closure, contextual errors, and documented exported contracts. | PASS |
| II. Testing Standards | Store migration, dispatch, coalescing, recovery, restart, race, CLI, API, and headless GUI tests are mandatory; timing uses the injected Clock. | PASS |
| III. User Experience Consistency | Watchers use the existing local API envelope, verb-noun CLI lifecycle, JSON parity, actionable health, and structured desktop controls. | PASS |
| IV. Performance Requirements | The design retains the existing dispatch budget, adds 100-watcher lifecycle budgets, and keeps one observer and one event-loop goroutine instead of one poller per watcher. | PASS |
| V. Autonomous Execution | S056 is issue-backed and runs the full Spec Kit sequence, analyze gate, review branch, authorized PR, and at most two review rounds. | PASS |

## Architecture

### Durable configuration

Schema v14 adds `filesystem_watchers` with a task foreign key and validated selection and timing fields. It also adds nullable `source_watcher_id` to runs. Runtime health is deliberately absent from SQLite because a stored `active` state would become false the instant the daemon stops.

### Native observer boundary

`internal/watcher` owns the event loop and depends on narrow Store, Dispatcher, Observer, Clock, Stat, and lifecycle-reporting interfaces. Production uses one fsnotify observer. File watchers observe the parent directory and filter the exact file name, which preserves atomic replacement workflows. Recursive directory watchers walk only real directories and add newly created real subdirectories.

### Generation and timing state machine

Every successful reload closes the old observer, clears pending candidates, increments a generation, loads current enabled definitions, and builds a new registration map. Candidate state is keyed by watcher ID and cleaned absolute file path. The first matching create or write signal sets a debounce deadline. Subsequent signals move that deadline. After debounce, the runtime captures size and modification time, waits one stability period, and dispatches only if the second snapshot matches. A mutation or degradation invalidates the generation before any pending dispatch.

### Health and recovery

Health is a mutex-protected runtime snapshot with `active`, `disabled`, or `degraded`, bounded reason text, and transition time. Missing or rejected roots and observer failures become degraded. A single two-second injected-clock recovery deadline rebuilds the complete observer. Identical state and reason pairs are suppressed; transitions publish one event and one structured log without a full path.

### Dispatch

The engine validates the watcher and target task at acceptance time using the same readiness and group checks as external triggers, then calls `dispatchWithOrigin` with trigger `filesystem_watcher` and `source_watcher_id`. Existing overlap, worker pool, execution, alerts, completion delivery, and history code remains authoritative.

### Interfaces

The local API provides `/v1/watchers` lifecycle routes and joins durable configuration to current runtime health. `gosched watcher` exposes the same lifecycle with duration parsing and JSON parity. The desktop Triggers destination gains a separate Filesystem Watchers table and toolbar so opaque keys and path rules are not visually conflated.

## Project Structure

### Documentation

```text
specs/056-filesystem-watchers/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/
│   └── filesystem-watcher-contract.md
└── checklists/
    ├── requirements.md
    └── watcher-contract.md
```

### Source Code

```text
cmd/goschedd/main.go
gui/
├── app.go
├── triggers.go
├── triggers_test.go
└── viewmodel/
    ├── viewmodel.go
    └── viewmodel_test.go
internal/
├── api/
│   ├── client/methods.go
│   └── server/
│       ├── server.go
│       ├── watchers.go
│       └── watchers_test.go
├── cli/
│   ├── root.go
│   ├── watcher.go
│   └── watcher_test.go
├── domain/domain.go
├── engine/
│   ├── engine.go
│   ├── triggers.go
│   └── watchers.go
├── events/broker.go
├── store/
│   ├── crud.go
│   ├── migration_v14_test.go
│   ├── store.go
│   ├── watchers.go
│   └── watchers_test.go
└── watcher/
    ├── manager.go
    ├── manager_test.go
    ├── observer.go
    └── observer_integration_test.go
```

**Structure Decision**: The native observation state machine is isolated in `internal/watcher`, while engine, store, API, CLI, and GUI retain their existing responsibilities. This avoids placing platform-event mechanics inside the scheduling loop and makes observer and clock behavior independently testable.

## Spec Kit Ordering Deviation

The installed checklist prerequisite requires `plan.md`, contradicting the mandated `specify -> clarify -> checklist -> plan` sequence. S056 used the active feature path returned by the successful clarification prerequisite call and generated `checklists/watcher-contract.md` before running plan setup. This is the established repository workaround and avoids unrelated workflow-tooling changes.

## Post-Design Constitution Recheck

The final design adds no constitution exception. The new long-lived goroutine has one context-owned termination path; deterministic timing is injected; persistence is forward-only; paths are normalized and not copied into durable run provenance; cross-platform limits are explicit; and every user-facing interface uses the existing authenticated local IPC boundary.
