# Implementation Plan: Task Execution Safety and Diagnostics

**Branch**: `codex/051-task-execution-safety` | **Date**: 2026-09-05 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/051-task-execution-safety/spec.md`

## Summary

Deliver one operator-safety path across task creation, eligibility presentation, and failed-run diagnosis. Add an optional initial-enabled intent to task creation while preserving the existing default for callers that omit it; let the desktop send an explicit inactive value unless its creation-only checkbox is selected. Persist exact run correlation on failure alerts and output-truncation metadata, expose an exact-run lookup through the local API, and enrich Activity detail from that immutable identity. Extend the existing group-chain helper so the Tasks table can name the nearest disabled group without copying scheduler policy.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Go standard library, Fyne v2, modernc.org/sqlite, existing local HTTP/JSON IPC client/server

**Storage**: Existing SQLite store, forward-only schema migration v10

**Testing**: Go unit/integration tests, Fyne headless tests, race detector, coverage gate, golangci-lint, documentation and automation checks

**Target Platform**: Windows, Linux, and macOS daemon/client; Fyne desktop GUI

**Project Type**: Single Go module with daemon, CLI, desktop GUI, local API, and embedded persistence

**Performance Goals**: Task-row derivation remains linear in group-chain depth; exact run detail performs bounded point lookups; scheduler p99 dispatch budget remains below 100 ms

**Constraints**: No transient activation on inactive GUI creation; no direct SQLite access from GUI; additive wire changes; forward-only non-destructive migration; output remains byte-bounded; no new secret-bearing diagnostics; no new dependency; no visible Windows console process

**Scale/Scope**: Three GitHub issues, one additive migration, two additive API fields, one exact lookup endpoint, three GUI surfaces, and focused changes across domain/store/engine/API/client/GUI

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **I. Code Quality**: PASS. Existing ownership boundaries remain intact. The store owns persistence, the server owns compatibility defaults, the engine owns alert creation, and the GUI remains a client. New exported contracts will receive intent-bearing comments and all errors remain contextual.
- **II. Testing Standards**: PASS. Regression tests precede each behavior: migration rollback/data preservation, optional create semantics, exact run lookup, failure correlation, truncation, group-chain attribution, and headless GUI state/detail rendering. Canonical race and coverage gates remain mandatory.
- **III. User Experience Consistency**: PASS. Task state terms are separated, failure messages become actionable, all timestamps remain RFC 3339, and diagnostics retain task/run correlation without new secret disclosure.
- **IV. Performance Requirements**: PASS. No scheduler hot-path algorithm changes beyond assigning a known run ID to an alert. Group display traverses only the in-memory ancestor chain and exact-run retrieval is keyed by primary ID. No new benchmark is warranted.
- **V. Autonomous Build-Phase Execution**: PASS. S051 is traceable to #102, #118, and #120, uses a `codex/` review branch, executes the full Spec Kit sequence, and will publish one PR under the user's explicit authorization.
- **Engineering constraints**: PASS. Schema and wire changes are additive, existing callers retain behavior, inputs remain locally authorized by the IPC boundary, and no dependency or platform support changes.

## Phase 0 Research Decisions

1. Model initial activation as an optional create request boolean. Omission keeps the legacy active default; the GUI always expresses its explicit checkbox state. This is the only option that is both atomic and backward compatible.
2. Add migration v10 with nullable `alerts.run_id` text and `runs.output_truncated` boolean-default metadata. Alert correlation does not use a foreign key so deleting task-owned run history cannot erase the durable run identifier from a surviving alert.
3. Add `GET /v1/runs/{id}` rather than selecting a recent run by task or time. Exact primary-key lookup prevents close failures from being mismatched and gives the client an honest not-found state.
4. Keep combined stdout/stderr storage. Add explicit truncation metadata rather than inserting a marker into output bytes, which would either exceed the cap or replace captured content.
5. Extend the shared group-chain package with nearest-disabled-group discovery and keep `ChainEnabled` authoritative for cycle behavior. The GUI will show an unavailable-chain state when the policy is disabled but no named blocker exists.
6. Build Activity diagnostic text as a pure formatter and perform bounded background enrichment by exact IDs. The dialog opens immediately with durable IDs, then updates when optional run/task data arrives.

See [research.md](research.md) for alternatives and detailed rationale.

## Phase 1 Design

- Add optional initial activation to `TaskCreateRequest`; resolve omission to enabled at the server boundary before persistence and event publication.
- Persist `Run.OutputTruncated` and `Alert.RunID` through migration v10, CRUD, API JSON, and existing event streaming.
- Add exact `Store.GetRun`, server route, client method, and GUI backend method.
- Pass the persisted run ID into only the run-failure alert path. Other alert kinds remain uncorrelated and backward compatible.
- Add a creation-only Fyne checkbox to editor state, dirty snapshots, form building, and task submission. Existing edit dialogs omit it.
- Derive Tasks-table effective state from task enabled/lifecycle precedence and shared group-chain helpers; add a labeled Effective column with full value in disclosure text.
- Preserve alert task/run IDs in `logEntry`. Enrich run-failure detail from exact lookups and render task/run identity, trigger, outcome, exit/launch state, combined output, empty output, truncation, and unavailable fallbacks.

### Post-design constitution re-check

All gates remain PASS. Migration v10 is additive with total defaults. The point lookup and group traversal are bounded. GUI enrichment remains asynchronous. No sensitive task inputs are added to output. No pinned workflow, toolchain, release configuration, installer, or dependency file changes are planned.

## Project Structure

### Documentation (this feature)

```text
specs/051-task-execution-safety/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── task-execution-contract.md
├── checklists/
│   ├── requirements.md
│   └── operational-safety.md
└── tasks.md
```

### Source Code (repository root)

```text
gui/
├── app.go / app_test.go
├── editor.go / editor_test.go / editor_data_test.go
├── tasks.go / tasks_test.go / groupchoice_test.go
└── logs.go / logs_test.go

internal/
├── api/
│   ├── client/methods.go
│   └── server/{server.go,tasks.go,tasks_test.go,runs.go,runs_test.go}
├── domain/domain.go
├── engine/{engine.go,engine_extra_test.go}
├── executor/{executor.go,executor_test.go}
├── store/{store.go,crud.go,store_test.go,migration_v10_test.go}
└── task/{group.go,group_test.go}
```

**Structure Decision**: Extend the existing single-module layered structure. No new package is needed: each behavior belongs to an existing domain, store, engine, API, client, group-policy, or GUI boundary.

## Complexity Tracking

No constitution violation or exceptional complexity is required.
