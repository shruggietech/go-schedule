# Implementation Plan: Scheduler Startup Trigger

**Branch**: `codex/031-reboot-startup-trigger` | **Date**: 2026-08-30 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/031-reboot-startup-trigger/spec.md`

## Summary

Make `@reboot` a first-class once-per-daemon-start event across cron, task authoring, engine, persistence, run history, desktop, and documentation. Reuse the retained event schedule shape (`kind=event`, `trigger_id=scheduler_startup`) without restoring task-completion triggers. The engine takes one startup snapshot per `Start` invocation, dispatches eligible startup tasks through the normal overlap/worker path with a distinct `startup` run origin, and never repeats that boundary on reload.

## Technical Context

**Language/Version**: Go 1.25

**Primary Dependencies**: Standard library, Cobra, Fyne, rrule-go, modernc.org/sqlite

**Storage**: Existing SQLite `schedules` and `runs` tables; no schema migration

**Testing**: Go unit/integration tests, deterministic fake clock, controlled engine contexts, shell verification suite

**Target Platform**: Linux, macOS, and Windows daemon, CLI, and desktop GUI

**Project Type**: Cross-platform local daemon with CLI, local HTTP/IPC API, and desktop client

**Performance Goals**: Preserve dispatch-latency p99 below 100 ms under nominal load

**Constraints**: Exactly once per engine start, never on reload; no direct `time.Now()` in engine code; no catch-up; no trigger-table revival; existing export context refusals remain intact

**Scale/Scope**: One event identity spanning schedule parsing, engine startup, run history, cron import/export, API, CLI, GUI, persistence compatibility, and user documentation

## Constitution Check

*GATE: Passed before research and re-checked after design.*

- **I. Code Quality and Architecture**: PASS. One typed startup identity flows through existing package boundaries; no generic dispatcher or new service is introduced.
- **II. Testing Standards**: PASS. Lifecycle tests use injected time and controlled cancellation, authoring tests precede implementation, prior-database persistence is exercised, and the coverage floor remains unchanged.
- **III. User Experience Consistency**: PASS. Both syntaxes round-trip, all task clients share `scheduleinput`, and event previews contain no clock times.
- **IV. Performance**: PASS. Startup scans the already loaded active-task snapshot once; recurring hot-path computation and dispatch remain unchanged.
- **V. Build-Phase Autopilot**: PASS. S031 follows specify through local commit and halts before publication.
- **Pinned artifacts**: PASS. No pinned process artifact is planned.

Post-design re-check: PASS. Research and contracts preserve all gates with no exception or complexity waiver.

## Project Structure

### Documentation (this feature)

```text
specs/031-reboot-startup-trigger/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/startup-schedule.md
├── checklists/
└── tasks.md
```

### Source Code (repository root)

```text
internal/
├── domain/           # event identity and startup run origin
├── schedule/         # canonical human phrase and no-clock semantics
├── scheduleinput/    # shared human/cron authoring boundary
├── cron/             # parse, explain, compile, convert, import, and export
├── engine/           # once-per-Start snapshot and dispatch
├── store/            # existing event schedule and run persistence
├── api/server/       # preview/create/update contracts
└── cli/              # explain, convert, import/export, task help
gui/                  # Schedule field prefill, preview, help, and tests
test/integration/     # restart and persistence scenarios
docs/                 # CLI, GUI, cron, and test-script references
```

**Structure Decision**: Extend existing packages in place. The event schedule representation and database columns are sufficient; a new trigger package, event bus, schema table, or editor mode would be unnecessary architecture.

## Complexity Tracking

No constitution violations or complexity exceptions.

