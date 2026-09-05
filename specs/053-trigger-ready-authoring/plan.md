# Implementation Plan: Trigger-Ready Task Authoring

**Branch**: `codex/053-trigger-ready-authoring` | **Date**: 2026-09-05 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/053-trigger-ready-authoring/spec.md`

## Summary

Make schedule ownership optional for tasks, derive readiness from stored task and activation-source state, enforce readiness at mutation and dispatch boundaries, and expose the result consistently through API, CLI, Tasks, and Groups. Complete the authoring experience with two semantic sidebar sections and one-path keyboard group submission.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Go standard library, existing Cobra CLI, existing Fyne v2.7.4 GUI, and existing modernc SQLite driver

**Storage**: Existing SQLite database with a forward-only v11 migration making `tasks.schedule_id` nullable; no new persisted readiness column

**Testing**: Go unit and integration tests, migration fixtures, Fyne headless widget tests, race detector, coverage gate, documentation and automation checks

**Target Platform**: Windows, Linux, and macOS daemon, CLI, and Fyne desktop GUI

**Project Type**: Single Go module with daemon, local API, CLI, and desktop GUI

**Performance Goals**: Readiness checks use indexed task and completion-chain lookups; scheduler recomputation remains bounded by the existing task scan and performs no new background polling

**Constraints**: No new dependency, fake schedule, persisted malformed input, external trigger implementation, or second source of readiness truth

**Scale/Scope**: One schema migration, task create/update/detail/enable/run boundaries, completion-chain target transitions, CLI draft creation and clearing, two GUI views, navigation shell, and group dialog

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **I. Code Quality**: PASS. One domain readiness vocabulary is derived at boundaries; storage owns atomic eligibility changes and the GUI consumes shared results rather than inventing state.
- **II. Testing Standards**: PASS. Tests precede implementation across migration preservation, API transitions, manual execution, completion-chain changes, CLI behavior, and headless interaction.
- **III. User Experience Consistency**: PASS. API, CLI, Tasks, and Groups share `unnamed`, not runnable, manual only, ready, disabled, blocked, and terminal meanings.
- **IV. Performance Requirements**: PASS. Readiness evaluation adds bounded indexed reads on mutation and existing list refresh paths, with no hot-loop or scheduler-latency expansion requiring a benchmark.
- **V. Autonomous Build-Phase Execution**: PASS. S053 is traceable to #129, #130, and #131, uses a `codex/` review branch, and has explicit push and pull-request authorization.
- **Engineering constraints**: PASS. Migration is forward-only and preservation-tested; no dependency, network surface, secret, or platform-specific behavior is added.

## Phase 0 Research Decisions

1. Represent an absent schedule as `NULL`, not a fabricated event or placeholder schedule.
2. Derive readiness rather than persisting it, because command, lifecycle, schedule, completion-chain, group, and future trigger changes already supply the authoritative facts.
3. Keep manual execution independent of enabled state but reject it before dispatch when the command is absent.
4. Make enable validation and source-removal auto-disable operations transactional in the store to close request races.
5. Preserve invalid-input rejection and allow only omitted values to form drafts in S053.
6. Preserve wire compatibility by retaining existing JSON field shapes for configured tasks, returning `schedule: null` only for unscheduled tasks, and adding explicit clear intent to task updates.
7. Build navigation sections from destination metadata so #133 can insert Triggers without restructuring the shell.
8. Replace the stock group form dialog with a small shared submission controller so Enter and Create use one validation and duplicate-suppression path.

See [research.md](research.md) for alternatives and detailed rationale.

## Phase 1 Design

- Add a nullable schedule reference through a preservation-tested v11 table rebuild performed with foreign-key enforcement safely suspended and verified around the migration transaction.
- Add pure readiness vocabulary and display-name fallback helpers to `internal/task`, with storage queries contributing incoming completion-chain state.
- Return task detail with an optional schedule and derived readiness fields while preserving configured-task JSON.
- Extend create and update contracts for omitted values and explicit field clearing; incomplete create requests are forced disabled and readiness-removing updates atomically disable.
- Validate enable operations transactionally and validate Run now before dispatch.
- Make scheduler and calendar paths skip absent schedules as an expected state.
- Update the CLI to create nameless or partial drafts, clear mutable fields explicitly, and print absent schedule and readiness honestly.
- Update task editing and list views to preserve blank fields, show `unnamed`, and distinguish readiness from effective automatic eligibility.
- Group navigation destinations through stable section metadata and implement a validated group dialog whose Enter and Create actions converge.

### Post-design constitution re-check

All gates remain PASS. The migration and readiness changes are necessary to avoid placeholder schedules and split-brain state. Atomic storage methods are the smallest reliable boundary for concurrent mutation safety.

## Project Structure

### Documentation (this feature)

```text
specs/053-trigger-ready-authoring/
├── spec.md / plan.md / research.md / data-model.md / quickstart.md
├── contracts/task-readiness-contract.md
├── checklists/requirements.md / ux.md
└── tasks.md / verification.md
```

### Source Code (repository root)

```text
internal/
├── task/readiness.go / readiness_test.go
├── store/store.go / crud.go / chains.go / migration_v11_test.go
├── api/server/tasks.go / update.go / *_test.go
├── engine/engine.go / runnow.go / *_test.go
└── cli/task.go / task_test.go
gui/
├── editor.go / editor_test.go / editor_prefill_test.go
├── tasks.go / tasks_test.go
├── groups.go / groups_test.go
└── navigation.go / navigation_test.go
```

**Structure Decision**: Put pure vocabulary in `internal/task`, atomic source-aware enforcement in `internal/store`, boundary validation in API and engine, and presentation behavior in existing CLI and GUI packages. No new package or dependency is warranted.

## Complexity Tracking

| Decision | Why needed | Simpler alternative rejected |
| --- | --- | --- |
| Rebuild `tasks` in migration v11 | SQLite cannot remove `NOT NULL` from `schedule_id` in place | Placeholder schedules lie about user intent and poison scheduler semantics |
| Transactional readiness enforcement | Completion-chain and task mutations can race across requests | Server-side check then update leaves a time-of-check gap |
| Optional schedule in task detail | An unscheduled task has no truthful schedule object | A zero-value schedule object is ambiguous to clients |
